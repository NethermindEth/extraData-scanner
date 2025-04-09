import argparse
from web3 import Web3
from collections import Counter
from tabulate import tabulate
import sys

def connect_web3(rpc_url):
    w3 = Web3(Web3.HTTPProvider(rpc_url))
    if not w3.is_connected():
        print("❌ Failed to connect to RPC:", rpc_url)
        sys.exit(1)
    return w3


def decode_extra_data(hex_str):
    try:
        return bytes.fromhex(hex_str[2:]).decode('utf-8', errors='replace')
    except Exception:
        return "(decode error)"


def scan_extra_data(w3, start_block, end_block):
    if end_block is None:
        end_block = w3.eth.get_block_number()

    extra_data_counter = Counter()
    total_blocks = end_block - start_block + 1
    progress_interval = max(1000, total_blocks // 100)  # Show progress every 1000 blocks or 1% of total
    errors = []

    print(f"🔍 Scanning blocks from {start_block} to {end_block}...")
    for block_number in range(start_block, end_block + 1):
        try:
            block = w3.eth.get_block(block_number)
            extra_data = block.extraData.hex()
            extra_data_counter[extra_data] += 1
            
            # Show progress in batches
            if (block_number - start_block + 1) % progress_interval == 0:
                progress = (block_number - start_block + 1) / total_blocks * 100
                print(f"Progress: {progress:.1f}% ({block_number:,}/{end_block:,} blocks)")
        except Exception as e:
            errors.append(f"Block {block_number}: {e}")
            if len(errors) <= 5:  # Show only first 5 errors
                print(f"⚠️ Error at block {block_number}: {e}")
            elif len(errors) == 6:  # Show message when suppressing errors
                print("⚠️ Additional errors suppressed...")

    # Show final error count if any
    if errors:
        print(f"\n⚠️ Total errors encountered: {len(errors)}")

    return extra_data_counter, total_blocks

def print_summary(counter, total):
    print(f"\n📊 Summary of {total:,} blocks processed:")
    table = []
    for ed_hex, count in counter.most_common():
        ed_str = decode_extra_data(ed_hex)
        percent = (count / total) * 100
        table.append([ed_hex, ed_str, count, f"{percent:.2f}%"])
    print(tabulate(table, headers=["Hex", "Decoded", "Count", "Percentage"], tablefmt="github"))


def main():
    parser = argparse.ArgumentParser(description="Scan Gnosis Chain blocks for unique extraData values.")
    parser.add_argument("--rpc", type=str, default="https://rpc.gnosischain.com", help="RPC URL for Gnosis Chain")
    parser.add_argument("--start", type=int, required=True, help="Start block number")
    parser.add_argument("--end", type=int, help="End block number (optional; defaults to latest)")

    args = parser.parse_args()

    w3 = connect_web3(args.rpc)
    extra_data_counter, total = scan_extra_data(w3, args.start, args.end)
    print_summary(extra_data_counter, total)

if __name__ == "__main__":
    main()

