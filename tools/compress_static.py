#!/usr/bin/env python3
"""
compress_static.py
預先壓縮 Godot HTML5 靜態檔案（wasm、pck、js）
Server 會優先提供 .gz 版本，讓瀏覽器下載更快

使用方式：
    py tools/compress_static.py

輸出：
    server/static/index.wasm.gz  (36MB -> ~9MB, -75%)
    server/static/index.pck.gz   (1MB  -> ~0.5MB, -50%)
    server/static/index.js.gz    (309KB -> ~100KB, -67%)
"""

import gzip
import os
import sys

STATIC_DIR = os.path.join(os.path.dirname(__file__), '..', 'server', 'static')
TARGETS = ['index.wasm', 'index.pck', 'index.js',
           'index.audio.worklet.js', 'index.audio.position.worklet.js']

def compress_file(src_path: str, dst_path: str, level: int = 6) -> tuple[int, int]:
    """壓縮單一檔案，回傳 (原始大小, 壓縮後大小)"""
    with open(src_path, 'rb') as f:
        data = f.read()
    compressed = gzip.compress(data, compresslevel=level)
    with open(dst_path, 'wb') as f:
        f.write(compressed)
    return len(data), len(compressed)

def main():
    print("🗜️  壓縮 Godot HTML5 靜態檔案...")
    print(f"📁 目錄：{os.path.abspath(STATIC_DIR)}")
    print()

    total_original = 0
    total_compressed = 0

    for filename in TARGETS:
        src = os.path.join(STATIC_DIR, filename)
        dst = src + '.gz'

        if not os.path.exists(src):
            print(f"  ⚠️  {filename} 不存在，跳過")
            continue

        orig_size, comp_size = compress_file(src, dst)
        reduction = 100 - (comp_size * 100 // orig_size)
        total_original += orig_size
        total_compressed += comp_size

        print(f"  ✅ {filename}")
        print(f"     {orig_size/1024/1024:.1f}MB → {comp_size/1024/1024:.1f}MB  (-{reduction}%)")

    print()
    total_reduction = 100 - (total_compressed * 100 // total_original)
    print(f"📊 總計：{total_original/1024/1024:.1f}MB → {total_compressed/1024/1024:.1f}MB  (-{total_reduction}%)")
    print()
    print("✅ 完成！Server 啟動後會自動提供 .gz 版本")
    print("   瀏覽器需支援 Accept-Encoding: gzip（所有現代瀏覽器都支援）")

if __name__ == '__main__':
    main()
