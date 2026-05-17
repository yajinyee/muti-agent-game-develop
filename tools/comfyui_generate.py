# -*- coding: utf-8 -*-
"""
ComfyUI API 整合腳本
透過 ComfyUI HTTP API 生成吉伊卡哇像素藝術 Sprite

使用方式：
  py tools/comfyui_generate.py --character chiikawa --pose idle
  py tools/comfyui_generate.py --character hachiware --pose attack
  py tools/comfyui_generate.py --all

ComfyUI 需先啟動：
  C:\ComfyUI\ComfyUI_windows_portable\run_nvidia_gpu.bat
"""
import os
import sys
import json
import time
import uuid
import argparse
import urllib.request
import urllib.parse

COMFYUI_URL = "http://127.0.0.1:8188"
OUTPUT_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\ai_generated"

# 角色提示詞（洋紅色背景策略，來自 agent-sprite-forge）
CHARACTER_PROMPTS = {
    "chiikawa": {
        "positive": "chiikawa character, white fluffy round creature, tiny cute chibi, big black eyes with white highlight, small pink blush marks, pixel art sprite, 16-bit retro game style, solid magenta background #FF00FF, white fur, simple clean design, game sprite, centered in frame",
        "negative": "realistic, 3d render, blurry, white background, transparent background, gradient background, text, watermark, human, tall, thin, adult, multiple characters",
        "color_hint": "white body, black outline, pink blush"
    },
    "hachiware": {
        "positive": "hachiware character, white cat with blue stripes on head, pointed ears, cute chibi, big black eyes, pixel art sprite, 16-bit retro game style, solid magenta background #FF00FF, white fur with blue markings, simple clean design, game sprite, centered in frame",
        "negative": "realistic, 3d render, blurry, white background, transparent background, gradient background, text, watermark, human, multiple characters",
        "color_hint": "white body, blue stripes, black outline"
    },
    "usagi": {
        "positive": "usagi rabbit character, white rabbit with long ears, red eyes, cute chibi, pixel art sprite, 16-bit retro game style, solid magenta background #FF00FF, white fur, pink inner ears, simple clean design, game sprite, centered in frame",
        "negative": "realistic, 3d render, blurry, white background, transparent background, gradient background, text, watermark, human, multiple characters",
        "color_hint": "white body, long ears, red eyes, black outline"
    }
}

# 動作提示詞
POSE_PROMPTS = {
    "idle": "standing pose, arms at sides, relaxed, neutral expression",
    "attack": "attacking pose, swinging weapon upward, dynamic action, determined expression, motion blur on weapon",
    "bigwin": "celebrating pose, jumping, arms raised high, happy expression, stars around"
}

def build_workflow(character, pose, seed=None):
    """建立 ComfyUI workflow JSON（SD 1.5 + Pixel Art LoRA）"""
    char_data = CHARACTER_PROMPTS[character]
    pose_text = POSE_PROMPTS[pose]
    
    positive = f"{char_data['positive']}, {pose_text}"
    negative = char_data["negative"]
    
    if seed is None:
        seed = int(time.time()) % 1000000
    
    workflow = {
        "1": {
            "inputs": {"ckpt_name": "v1-5-pruned-emaonly.safetensors"},
            "class_type": "CheckpointLoaderSimple"
        },
        "2": {
            "inputs": {
                "lora_name": "pixel_art_lora.safetensors",
                "strength_model": 0.85,
                "strength_clip": 0.85,
                "model": ["1", 0],
                "clip": ["1", 1]
            },
            "class_type": "LoraLoader"
        },
        "3": {
            "inputs": {
                "seed": seed,
                "steps": 28,
                "cfg": 7.5,
                "sampler_name": "euler_ancestral",
                "scheduler": "normal",
                "denoise": 1.0,
                "model": ["2", 0],
                "positive": ["6", 0],
                "negative": ["7", 0],
                "latent_image": ["5", 0]
            },
            "class_type": "KSampler"
        },
        "5": {
            "inputs": {"width": 512, "height": 512, "batch_size": 1},
            "class_type": "EmptyLatentImage"
        },
        "6": {
            "inputs": {"text": positive, "clip": ["2", 1]},
            "class_type": "CLIPTextEncode"
        },
        "7": {
            "inputs": {"text": negative, "clip": ["2", 1]},
            "class_type": "CLIPTextEncode"
        },
        "8": {
            "inputs": {"samples": ["3", 0], "vae": ["1", 2]},
            "class_type": "VAEDecode"
        },
        "9": {
            "inputs": {
                "filename_prefix": f"{character}_{pose}",
                "images": ["8", 0]
            },
            "class_type": "SaveImage"
        }
    }
    return workflow

def queue_prompt(workflow):
    """送出 workflow 到 ComfyUI"""
    client_id = str(uuid.uuid4())
    data = json.dumps({"prompt": workflow, "client_id": client_id}).encode("utf-8")
    req = urllib.request.Request(
        f"{COMFYUI_URL}/prompt",
        data=data,
        headers={"Content-Type": "application/json"}
    )
    try:
        with urllib.request.urlopen(req, timeout=30) as r:
            result = json.loads(r.read())
            return result.get("prompt_id"), client_id
    except Exception as e:
        print(f"  ERROR queuing prompt: {e}")
        return None, None

def wait_for_result(prompt_id, timeout=120):
    """等待生成完成"""
    start = time.time()
    while time.time() - start < timeout:
        try:
            with urllib.request.urlopen(f"{COMFYUI_URL}/history/{prompt_id}", timeout=10) as r:
                history = json.loads(r.read())
                if prompt_id in history:
                    outputs = history[prompt_id].get("outputs", {})
                    for node_id, node_output in outputs.items():
                        if "images" in node_output:
                            return node_output["images"]
        except Exception:
            pass
        time.sleep(2)
        elapsed = int(time.time() - start)
        print(f"\r  Waiting... {elapsed}s", end="", flush=True)
    return None

def download_image(filename, subfolder=""):
    """從 ComfyUI 下載生成的圖片"""
    params = urllib.parse.urlencode({
        "filename": filename,
        "subfolder": subfolder,
        "type": "output"
    })
    url = f"{COMFYUI_URL}/view?{params}"
    try:
        with urllib.request.urlopen(url, timeout=30) as r:
            return r.read()
    except Exception as e:
        print(f"  ERROR downloading: {e}")
        return None

def post_process_sprite(img_data, character, pose):
    """後處理：使用 process_sprites.py 的洋紅色去背 + shared_scale 技術"""
    try:
        import io
        import sys
        import math
        from collections import deque
        from PIL import Image, ImageEnhance

        img = Image.open(io.BytesIO(img_data)).convert("RGBA")
        print(f"  Raw size: {img.size}")

        # 縮小到 192x192 再處理（保持細節）
        if img.size[0] > 192:
            img = img.resize((192, 192), Image.NEAREST)

        # 去洋紅色背景（agent-sprite-forge 技術）
        pixels = img.load()
        w, h = img.size

        def dist_magenta(r, g, b):
            return math.sqrt((r - 255)**2 + g**2 + (b - 255)**2)

        # 全局 threshold
        for x in range(w):
            for y in range(h):
                r, g, b, a = pixels[x, y]
                if a == 0:
                    continue
                if dist_magenta(r, g, b) < 100:
                    pixels[x, y] = (0, 0, 0, 0)

        # BFS flood fill 從邊緣
        visited = set()
        queue = deque()
        for x in range(w):
            queue.append((x, 0)); queue.append((x, h-1))
        for y in range(h):
            queue.append((0, y)); queue.append((w-1, y))
        while queue:
            x, y = queue.popleft()
            if (x, y) in visited or x < 0 or x >= w or y < 0 or y >= h:
                continue
            visited.add((x, y))
            r, g, b, a = pixels[x, y]
            if a == 0:
                for dx in (-1, 0, 1):
                    for dy2 in (-1, 0, 1):
                        if dx == 0 and dy2 == 0: continue
                        if (x+dx, y+dy2) not in visited:
                            queue.append((x+dx, y+dy2))
            elif dist_magenta(r, g, b) < 150:
                pixels[x, y] = (0, 0, 0, 0)
                for dx in (-1, 0, 1):
                    for dy2 in (-1, 0, 1):
                        if dx == 0 and dy2 == 0: continue
                        if (x+dx, y+dy2) not in visited:
                            queue.append((x+dx, y+dy2))

        # 縮放到 96x96，bottom 對齊
        bbox = img.getbbox()
        if bbox:
            img = img.crop(bbox)
        cw, ch = img.size
        cell = 96
        fit = 0.82
        scale = min(cell/cw, cell/ch) * fit
        nw = max(1, int(cw * scale))
        nh = max(1, int(ch * scale))
        img = img.resize((nw, nh), Image.NEAREST)
        canvas = Image.new("RGBA", (cell, cell), (0, 0, 0, 0))
        px_x = (cell - nw) // 2
        pad = max(0, int(cell * (1 - fit) * 0.4))
        px_y = cell - nh - pad
        canvas.paste(img, (px_x, px_y))

        # 增強飽和度
        canvas = ImageEnhance.Color(canvas).enhance(1.3)
        canvas = ImageEnhance.Contrast(canvas).enhance(1.15)

        return canvas
    except ImportError:
        print("  PIL not available, saving raw image")
        return None

def generate_sprite(character, pose):
    """生成單個 Sprite"""
    print(f"\n[GENERATE] {character} - {pose}")
    
    # 建立 workflow
    workflow = build_workflow(character, pose)
    
    # 送出請求
    prompt_id, client_id = queue_prompt(workflow)
    if not prompt_id:
        print("  Failed to queue prompt")
        return False
    
    print(f"  Prompt ID: {prompt_id}")
    
    # 等待結果
    images = wait_for_result(prompt_id)
    if not images:
        print("\n  Timeout waiting for result")
        return False
    
    print(f"\n  Generated {len(images)} image(s)")
    
    # 下載並儲存
    os.makedirs(OUTPUT_DIR, exist_ok=True)
    for img_info in images:
        img_data = download_image(img_info["filename"], img_info.get("subfolder", ""))
        if img_data:
            # 後處理
            processed = post_process_sprite(img_data, character, pose)
            
            save_path = os.path.join(OUTPUT_DIR, f"{character}_{pose}.png")
            if processed:
                processed.save(save_path)
            else:
                with open(save_path, "wb") as f:
                    f.write(img_data)
            print(f"  Saved: {save_path}")
            return True
    
    return False

def check_comfyui_running():
    """確認 ComfyUI 是否在跑"""
    try:
        with urllib.request.urlopen(f"{COMFYUI_URL}/system_stats", timeout=5) as r:
            stats = json.loads(r.read())
            print(f"  ComfyUI running: {stats.get('system', {}).get('python_version', 'unknown')}")
            return True
    except Exception:
        return False

def main():
    parser = argparse.ArgumentParser(description="Generate Chiikawa pixel art sprites via ComfyUI")
    parser.add_argument("--character", choices=["chiikawa", "hachiware", "usagi"], help="Character to generate")
    parser.add_argument("--pose", choices=["idle", "attack", "bigwin"], default="idle", help="Pose to generate")
    parser.add_argument("--all", action="store_true", help="Generate all characters and poses")
    args = parser.parse_args()
    
    print("=== ComfyUI Sprite Generator ===")
    
    # 確認 ComfyUI 在跑
    print("\n[CHECK] ComfyUI status...")
    if not check_comfyui_running():
        print("  ComfyUI is NOT running!")
        print(f"  Start it: C:\\ComfyUI\\ComfyUI_windows_portable\\run_nvidia_gpu.bat")
        sys.exit(1)
    
    if args.all:
        characters = ["chiikawa", "hachiware", "usagi"]
        poses = ["idle", "attack", "bigwin"]
        total = len(characters) * len(poses)
        done = 0
        for char in characters:
            for pose in poses:
                ok = generate_sprite(char, pose)
                done += 1
                print(f"  Progress: {done}/{total}")
    elif args.character:
        generate_sprite(args.character, args.pose)
    else:
        parser.print_help()
        print("\nExample:")
        print("  py tools/comfyui_generate.py --character chiikawa --pose idle")
        print("  py tools/comfyui_generate.py --all")

if __name__ == "__main__":
    main()
