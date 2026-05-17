# -*- coding: utf-8 -*-
"""
下載 SD 1.5 基礎模型到 ComfyUI checkpoints
來源：Comfy-Org 官方 archive（hash-identical 原版）
"""
import os
import time
import urllib.request

URL = "https://huggingface.co/Comfy-Org/stable-diffusion-v1-5-archive/resolve/main/v1-5-pruned-emaonly.safetensors"
SAVE_PATH = r"C:\ComfyUI\ComfyUI_windows_portable\ComfyUI\models\checkpoints\v1-5-pruned-emaonly.safetensors"
HEADERS = {"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"}

def main():
    print("=== Downloading SD 1.5 Base Model ===")
    print(f"URL: {URL}")
    print(f"Save: {SAVE_PATH}")

    if os.path.exists(SAVE_PATH):
        size = os.path.getsize(SAVE_PATH)
        print(f"Already exists: {size//1024//1024} MB")
        return

    os.makedirs(os.path.dirname(SAVE_PATH), exist_ok=True)
    tmp = SAVE_PATH + ".tmp"

    try:
        req = urllib.request.Request(URL, headers=HEADERS)
        start = time.time()
        with urllib.request.urlopen(req, timeout=600) as r:
            total = int(r.headers.get("Content-Length", 0))
            downloaded = 0
            with open(tmp, "wb") as f:
                while True:
                    chunk = r.read(1024 * 1024)
                    if not chunk:
                        break
                    f.write(chunk)
                    downloaded += len(chunk)
                    elapsed = time.time() - start
                    speed = downloaded / elapsed / 1024 / 1024 if elapsed > 0 else 0
                    if total > 0:
                        pct = downloaded * 100 // total
                        print(f"\r  {pct}% | {downloaded//1024//1024}MB/{total//1024//1024}MB | {speed:.1f} MB/s", end="", flush=True)

        os.rename(tmp, SAVE_PATH)
        print(f"\nDone! {os.path.getsize(SAVE_PATH)//1024//1024} MB")
    except Exception as e:
        print(f"\nERROR: {e}")
        if os.path.exists(tmp):
            os.remove(tmp)

if __name__ == "__main__":
    main()
