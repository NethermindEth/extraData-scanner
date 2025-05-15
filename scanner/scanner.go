package scanner

import (
	"context"
	"log/slog"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
)

type BlockData struct {
	BlockNumber uint64
	ExtraHex    string
	Extra       string
}

type ScanResult struct {
	TotalBlocks    uint64
	TotalProcessed uint64
	Data           map[string]uint64
}

func Scan(
	ctx context.Context,
	ethClient *ethclient.Client,
	startBlock uint64,
	endBlock uint64,
	workers uint,
) (ScanResult, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	blockChan := make(chan uint64, workers)
	resultsChan := make(chan BlockData)

	var wg sync.WaitGroup
	wg.Add(int(workers))

	for i := uint(0); i < workers; i++ {
		go worker(ctx, ethClient, blockChan, resultsChan, &wg)
	}

	go func() {
		for blockNumber := startBlock; blockNumber <= endBlock; blockNumber++ {
			blockChan <- blockNumber
		}
		close(blockChan)
	}()

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	data := make(map[string]uint64, 1000)
	processedBlocks := uint64(0)

	for result := range resultsChan {
		processedBlocks++
		slog.Debug("Processed block", "block", result.BlockNumber, "extraData", result.Extra)
		count, ok := data[result.ExtraHex]
		if ok {
			data[result.ExtraHex] = count + 1
		} else {
			data[result.ExtraHex] = 1
		}
	}

	return ScanResult{
		TotalBlocks:    endBlock - startBlock + 1,
		TotalProcessed: processedBlocks,
		Data:           data,
	}, nil
}

func worker(
	ctx context.Context,
	ethClient *ethclient.Client,
	blockChan <-chan uint64,
	resultsChan chan<- BlockData,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	for {
		select {
		case blockNumber, ok := <-blockChan:
			if !ok {
				return
			}
			block, err := ethClient.BlockByNumber(ctx, big.NewInt(int64(blockNumber)))
			if err != nil {
				slog.Error("Error getting block", "block", blockNumber, "error", err)
				resultsChan <- BlockData{
					BlockNumber: blockNumber,
					ExtraHex:    "(error)",
					Extra:       "",
				}
				continue
			}
			resultsChan <- BlockData{
				BlockNumber: blockNumber,
				ExtraHex:    hexutil.Encode(block.Extra()),
				Extra:       string(block.Extra()),
			}
		case <-ctx.Done():
			return
		}
	}
}
