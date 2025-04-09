# extraData-scanner

A command-line tool for scanning and analyzing the `extraData` field in Gnosis Chain blocks. This tool helps track and summarize the different `extraData` values used by validators across blocks.

## Features

- Scan any range of blocks on Gnosis Chain
- Decode `extraData` hex values to UTF-8 strings
- Generate statistical summaries of `extraData` usage
- Display results in a clear tabulated format

## Installation

```bash
# Clone the repository
git clone https://github.com/NethermindEth/extraData-scanner.git
cd extraData-scanner

# Install dependencies
pip install -r requirements.txt
```

## Usage

```bash
python scanner.py --start <start_block> [--end <end_block>] [--rpc <rpc_url>]
```

### Arguments

- `--start`: Required. The block number to start scanning from
- `--end`: Optional. The block number to end scanning at (defaults to latest block)
- `--rpc`: Optional. RPC URL for Gnosis Chain (defaults to https://rpc.gnosischain.com)

### Example

```bash
python scanner.py --start 1000000 --end 1000100
```

## Dependencies

- web3.py: For interacting with the Gnosis Chain
- tabulate: For formatting output tables

## GitHub Actions

This repository includes a GitHub Action workflow that can run the scanner automatically:

1. Go to the "Actions" tab in your repository
2. Select "ExtraData Scanner"
3. Click "Run workflow"
4. Fill in the parameters:
   - Start block (required)
   - End block (optional)
   - RPC URL (optional)

### Setting up the Fallback RPC URL

If no RPC URL is provided when running the action, it will try to use a fallback URL from repository secrets:

1. Go to your repository Settings
2. Navigate to Secrets and variables > Actions
3. Add a new secret named `GNOSIS_RPC_URL`
4. Set its value to your fallback RPC URL
