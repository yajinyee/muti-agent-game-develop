# -*- coding: utf-8 -*-
"""
背景下載像素藝術模型
目標：PixelArtV4.safetensors (SD 1.5 based, ~2GB)
"""
import os
import sys
import urllib.request
import time

MODELS_DIR = r"C:\ComfyUI\ComfyUI_windows_portable\ComfyUI\models\checkpoints"
LORAS_DIR = r"C:\ComfyUI\ComfyUI_windows_portable\ComfyUI\models\loras"

MODELS = [
    {
        "name": "PixelArtV4.safetensors",
        "url": "https://huggingface.co/Onodofthenorth/SD_PixelArt_SpriteSheet_Generator/resolve/main/PixelArtV4.safetensors",
        "dir": MODELS_DIR,
        "description": "Pixel Art Sprite Sheet Generator v4"
    },
    {
        "name": "pixel_art_lora.safetensors",
        "url": "https://huggingface.co/nerijs/pixel-art-xl/resolve/main/pixel-art-xl.safetensors",
        "dir": LORAS_DIR,
        "description": "Pixel Art XL LoRA"
    }
]

HEADERS = {
    "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
}

def download_file(url, save_path, description):
    print(f"\n[DOWNLOAD] {description}")
    print(f"  URL: {url}")
    print(f"  Save: {save_path}")
    
    if os.path.exists(save_path):
        size = os.path.getsize(save_path)
        print(f"  Already exists: {size//1024//1024} MB")
        return True
    
    os.makedirs(os.path.dirname(save_path), exist_ok=True)
    
    try:
        req = urllib.request.Request(url, headers=HEADERS)
        start = time.time()
        with urllib.request.urlopen(req, timeout=600) as r:
            total = int(r.headers.get("Content-Length", 0))
            downloaded = 0
            with open(save_path + ".tmp", "wb") as f:
                while True:
                    chunk = r.read(1024 * 1024)  # 1MB
                    if not chunk:
                        break
                    f.write(chunk)
                    downloaded += len(chunk)
                    elapsed = time.time() - start
                    speed = downloaded / elapsed / 1024 / 1024 if elapsed > 0 else 0
                    if total > 0:
                        pct = downloaded * 100 // total
                        print(f"\r  {pct}% | {downloaded//1024//1024}MB/{total//1024//1024}MB | {speed:.1f} MB/s", end="", flush=True)
        
        os.rename(save_path + ".tmp", save_path)
        print(f"\n  Done! {os.path.getsize(save_path)//1024//1024} MB")
        return True
    except Exception as e:
        print(f"\n  ERROR: {e}")
        if os.path.exists(save_path + ".tmp"):
            os.remove(save_path + ".tmp")
        return False

def wait_for_comfyui_dir(timeout=300):
    """等待 ComfyUI 解壓縮完成"""
    print(f"[WAIT] Waiting for C:\\ComfyUI\\ComfyUI_windows_portable\\ComfyUI\\models to be ready...")
    start = time.time()
    while time.time() - start < timeout:
        if os.path.exists(r"C:\ComfyUI\ComfyUI_windows_portable\ComfyUI\models"):
            print(f"  ComfyUI models dir ready!")
            return True
        time.sleep(5)
        print(f"\r  Waiting... {int(time.time()-start)}s", end="", flush=True)
    print(f"\n  Timeout! Creating dirs manually...")
    return False

def main():
    print("=== Pixel Art Model Downloader ===")
    
    # 等待 ComfyUI 解壓縮
    if not os.path.exists(r"C:\ComfyUI\ComfyUI_windows_portable\ComfyUI\models"):
        wait_for_comfyui_dir(timeout=600)
    
    # 確保目錄存在
    os.makedirs(MODELS_DIR, exist_ok=True)
    os.makedirs(LORAS_DIR, exist_ok=True)
    
    # 下載模型
    results = []
    for model in MODELS:
        save_path = os.path.join(model["dir"], model["name"])
        ok = download_file(model["url"], save_path, model["description"])
        results.append((model["name"], ok))
    
    print("\n=== Download Summary ===")
    for name, ok in results:
        status = "OK" if ok else "FAILED"
        print(f"  [{status}] {name}")

if __name__ == "__main__":
    main()
