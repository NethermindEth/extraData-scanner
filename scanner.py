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

    print(f"🔍 Scanning blocks from {start_block} to {end_block}...")
    for block_number in range(start_block, end_block + 1):
        try:
            block = w3.eth.get_block(block_number)
            extra_data = block.extraData.hex()
            extra_data_counter[extra_data] += 1
            print(f"Block {block_number}: {extra_data}")
        except Exception as e:
            print(f"⚠️ Error at block {block_number}: {e}")

    return extra_data_counter, end_block - start_block + 1

def print_summary(counter, total):
    print("\n📊 Summary of `extraData` values:")
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

