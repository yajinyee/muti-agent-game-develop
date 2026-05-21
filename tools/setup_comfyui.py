# -*- coding: utf-8 -*-
"""
ComfyUI 安裝和設定腳本
1. 解壓縮 ComfyUI portable
2. 下載像素藝術模型
3. 建立生成腳本
"""
import os
import subprocess
import urllib.request
import json

COMFYUI_ARCHIVE = r"C:\Users\yajinyee0306\Downloads\ComfyUI_portable.7z"
COMFYUI_DIR = r"C:\ComfyUI"
MODELS_DIR = os.path.join(COMFYUI_DIR, "ComfyUI_windows_portable", "ComfyUI", "models", "checkpoints")

HEADERS = {"User-Agent": "Mozilla/5.0 Chrome/120.0.0.0"}

# 像素藝術模型（SD 1.5 based，4GB VRAM 可跑）
PIXEL_ART_MODELS = [
    {
        "name": "pixel_art_sprite_diffusion.safetensors",
        "url": "https://huggingface.co/Onodofthenorth/SD_PixelArt_SpriteSheet_Generator/resolve/main/PixelArtV4.safetensors",
        "description": "Pixel Art Sprite Sheet Generator - 專門生成像素藝術 Sprite"
    }
]

def check_download_status():
    """確認 ComfyUI 是否下載完成"""
    if os.path.exists(COMFYUI_ARCHIVE):
        size = os.path.getsize(COMFYUI_ARCHIVE)
        print(f"ComfyUI archive: {size/1024/1024:.0f} MB")
        return size > 1000 * 1024 * 1024  # > 1GB 才算完成
    return False

def extract_comfyui():
    """解壓縮 ComfyUI"""
    seven_zip = r"C:\Program Files\7-Zip\7z.exe"
    if not os.path.exists(seven_zip):
        print("7-zip not found, trying alternative...")
        seven_zip = r"C:\Program Files (x86)\7-Zip\7z.exe"
    
    if not os.path.exists(seven_zip):
        print("ERROR: 7-zip not installed")
        return False
    
    os.makedirs(COMFYUI_DIR, exist_ok=True)
    cmd = [seven_zip, "x", COMFYUI_ARCHIVE, f"-o{COMFYUI_DIR}", "-y"]
    result = subprocess.run(cmd, capture_output=True, text=True)
    if result.returncode == 0:
        print(f"Extracted to: {COMFYUI_DIR}")
        return True
    else:
        print(f"Extraction failed: {result.stderr}")
        return False

def download_model(model_info):
    """下載模型"""
    os.makedirs(MODELS_DIR, exist_ok=True)
    save_path = os.path.join(MODELS_DIR, model_info["name"])
    
    if os.path.exists(save_path):
        print(f"  Already exists: {model_info['name']}")
        return True
    
    print(f"  Downloading: {model_info['name']}...")
    try:
        req = urllib.request.Request(model_info["url"], headers=HEADERS)
        with urllib.request.urlopen(req, timeout=300) as r:
            total = int(r.headers.get("Content-Length", 0))
            downloaded = 0
            with open(save_path, "wb") as f:
                while True:
                    chunk = r.read(1024 * 1024)  # 1MB chunks
                    if not chunk:
                        break
                    f.write(chunk)
                    downloaded += len(chunk)
                    if total > 0:
                        pct = downloaded * 100 // total
                        print(f"\r  Progress: {pct}% ({downloaded//1024//1024}MB/{total//1024//1024}MB)", end="")
        print(f"\n  Saved: {save_path}")
        return True
    except Exception as e:
        print(f"\n  Error: {e}")
        return False

def create_generation_workflow():
    """建立像素藝術生成 workflow"""
    workflow = {
        "description": "Chiikawa pixel art sprite generation workflow",
        "prompt_template": "chiikawa character, white fluffy creature, cute chibi, pixel art sprite, 16-bit retro game style, transparent background, simple design, {pose}",
        "negative_prompt": "realistic, 3d, blurry, complex background, text, watermark",
        "settings": {
            "width": 512,
            "height": 512,
            "steps": 20,
            "cfg_scale": 7.5,
            "sampler": "euler_a",
            "model": "v1-5-pruned-emaonly.safetensors",
            "lora": "pixel_art_lora.safetensors",
            "lora_strength": 0.8
        },
        "poses": {
            "idle": "standing pose, arms at sides",
            "attack": "attacking pose, swinging weapon, dynamic action",
            "bigwin": "jumping pose, arms raised, celebrating"
        }
    }
    
    workflow_path = r"D:\Kiro\tools\comfyui_workflow.json"
    with open(workflow_path, "w", encoding="utf-8") as f:
        json.dump(workflow, f, indent=2, ensure_ascii=False)
    print(f"Workflow saved: {workflow_path}")

def main():
    print("=== ComfyUI Setup ===")
    
    # 1. 確認下載狀態
    print("\n[1] Checking download status...")
    if check_download_status():
        print("  ComfyUI download complete!")
    else:
        print("  ComfyUI still downloading... check back later")
        print(f"  Archive: {COMFYUI_ARCHIVE}")
        return
    
    # 2. 解壓縮
    print("\n[2] Extracting ComfyUI...")
    if not extract_comfyui():
        return
    
    # 3. 下載模型
    print("\n[3] Downloading pixel art models...")
    for model in PIXEL_ART_MODELS:
        print(f"  {model['description']}")
        download_model(model)
    
    # 4. 建立 workflow
    print("\n[4] Creating generation workflow...")
    create_generation_workflow()
    
    print("\n=== Setup Complete ===")
    print(f"ComfyUI location: {COMFYUI_DIR}")
    print("To start ComfyUI:")
    print(f"  {COMFYUI_DIR}\\ComfyUI\\run_nvidia_gpu.bat")

if __name__ == "__main__":
    main()
