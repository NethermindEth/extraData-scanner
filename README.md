# extraData-scanner

A command-line tool for scanning and analyzing the `extraData` field in an EVM based Chain blocks. This tool helps track and summarize the different `extraData` values used by validators across blocks.

## Features

- Scan any range of blocks on an EVM based Chain
- Decode `extraData` hex values to UTF-8 strings
- Generate statistical summaries of `extraData` usage
- Display results in a clear tabulated format

## Installation

```bash
# Clone the repository
git clone https://github.com/NethermindEth/extraData-scanner.git
cd extraData-scanner

# Install dependencies
go mod tidy
```

## Usage

```bash
go run cmd/extradata-scanner/main.go --rpc=<rpc_url> --start=<start_block> [--end=<end_block>] 
```

### Arguments

- `--rpc`: Required. RPC URL for an EVM based Chain (must be provided, no default)
- `--start`: Optional. The block number to start scanning from (defaults to block 0)
- `--end`: Optional. The block number to end scanning at (defaults to latest block if not greater than start block)

### Example

```bash
go run cmd/extradata-scanner/main.go --rpc=<your_rpc_url> --start=1000000 --end=1000100
```

## GitHub Actions

This repository includes a GitHub Action workflow that can run the scanner automatically:

1. Go to the "Actions" tab in your repository
2. Select "ExtraData Scanner"
3. Click "Run workflow"
4. Fill in the parameters:
   - RPC URL (required)
   - Start block (optional)
   - End block (optional)
