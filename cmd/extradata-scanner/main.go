package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"extradata-scanner/scanner"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jedib0t/go-pretty/v6/progress"
	"github.com/jedib0t/go-pretty/v6/table"
)

func main() {
	var rpcURL string
	var startBlock uint64
	var endBlock uint64
	var workers uint

	flag.StringVar(&rpcURL, "rpc", "https://rpc.gnosischain.com", "Ethereum RPC URL")
	flag.Uint64Var(&startBlock, "start", 0, "Start block number")
	flag.Uint64Var(&endBlock, "end", 0, "End block number. If not provided or invalid, the latest block will be used.")
	flag.UintVar(&workers, "workers", 10, "Number of workers to use")
	flag.Parse()

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		slog.Error("Error dialing RPC URL", "error", err)
		os.Exit(1)
	}

	if endBlock <= startBlock {
		slog.Info("End block must be greater than start block. Using latest block as end block.")
		latestBlock, err := client.BlockByNumber(context.Background(), nil)
		if err != nil {
			slog.Error("Error getting latest block", "error", err)
			os.Exit(1)
		}
		endBlock = latestBlock.NumberU64()
		slog.Info("Using latest block as end block", "block", endBlock)
	}

	pw := progress.NewWriter()
	pw.SetStyle(progress.StyleCircle)
	go pw.Render()

	slog.Info("Scanning blocks", "start", startBlock, "end", endBlock, "rpc", rpcURL)
	scanResults, err := scanner.Scan(
		context.Background(),
		client,
		startBlock,
		endBlock,
		workers,
		pw,
	)
	pw.Stop()
	if err != nil {
		slog.Error("Error scanning blocks", "error", err)
		os.Exit(1)
	}

	slog.Info("Blocks scanned", "total", scanResults.TotalBlocks)
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Hex", "Decoded", "Count", "Percentage"})
	for hex, count := range scanResults.Data {
		decoded, err := hexutil.Decode(hex)
		if err != nil {
			slog.Error("Error decoding hex", "error", err)
			decoded = []byte("(decode error)")
		}
		t.AppendRow(table.Row{hex, string(decoded), count, fmt.Sprintf("%.2f%%", float64(count)/float64(scanResults.TotalBlocks)*100)})
	}
	t.SortBy([]table.SortBy{
		{Name: "Count", Mode: table.Dsc},
	})
	t.SetStyle(table.StyleLight)
	t.Render()
}
