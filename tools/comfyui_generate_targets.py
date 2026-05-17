# -*- coding: utf-8 -*-
"""
ComfyUI 目標物 AI 生成腳本
生成 11 個目標物的像素藝術圖片
"""
import os, sys, json, time, uuid, argparse
import urllib.request, urllib.parse

COMFYUI_URL = "http://127.0.0.1:8188"
OUTPUT_DIR  = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\targets"

# 目標物提示詞（洋紅色背景策略）
TARGET_PROMPTS = {
    "T001_grass": {
        "positive": "pixel art grass plant, three green leaves, simple cute design, 16-bit retro game sprite, solid magenta background #FF00FF, centered, game enemy sprite",
        "negative": "realistic, 3d, blurry, text, watermark, white background, gradient, character, animal"
    },
    "T002_bug_g": {
        "positive": "pixel art green bug insect, cute round body, big eyes, antennae, 16-bit retro game sprite, solid magenta background #FF00FF, centered, game enemy",
        "negative": "realistic, 3d, blurry, text, watermark, white background"
    },
    "T003_bug_r": {
        "positive": "pixel art red bug insect, cute round body, big eyes, antennae, 16-bit retro game sprite, solid magenta background #FF00FF, centered, game enemy",
        "negative": "realistic, 3d, blurry, text, watermark, white background"
    },
    "T004_bug_b": {
        "positive": "pixel art blue bug insect, cute round body, big eyes, antennae, 16-bit retro game sprite, solid magenta background #FF00FF, centered, game enemy",
        "negative": "realistic, 3d, blurry, text, watermark, white background"
    },
    "T005_pudding": {
        "positive": "pixel art cute pudding dessert with face, yellow jelly pudding, caramel top, big eyes smile, 16-bit retro game sprite, solid magenta background #FF00FF, centered",
        "negative": "realistic, 3d, blurry, text, watermark, white background"
    },
    "T006_mushroom": {
        "positive": "pixel art giant mushroom, red cap with white spots, thick stem, cute face, 16-bit retro game sprite, solid magenta background #FF00FF, centered, game enemy",
        "negative": "realistic, 3d, blurry, text, watermark, white background"
    },
    "T101_mimic": {
        "positive": "pixel art mimic monster disguised as grass plant, hidden evil eyes glowing purple, creepy smile, 16-bit retro game sprite, solid magenta background #FF00FF, centered",
        "negative": "realistic, 3d, blurry, text, watermark, white background"
    },
    "T102_chest": {
        "positive": "pixel art treasure chest monster, wooden chest with gold metal trim, large glowing eyes on lid, sharp teeth visible, dark thick outline, 16-bit retro game sprite, solid bright magenta background #FF00FF, centered, clear silhouette, high contrast",
        "negative": "realistic, 3d, blurry, text, watermark, white background, dark background, gradient, multiple objects"
    },
    "T103_meteor": {
        "positive": "pixel art glowing meteor shooting star, golden core with fire trail, sparkles, 16-bit retro game sprite, solid magenta background #FF00FF, centered",
        "negative": "realistic, 3d, blurry, text, watermark, white background"
    },
    "T104_gold_grass": {
        "positive": "pixel art golden shiny grass plant, three golden leaves, sparkle effects, 16-bit retro game sprite, solid magenta background #FF00FF, centered, rare item",
        "negative": "realistic, 3d, blurry, text, watermark, white background"
    },
    "T105_coin_fish": {
        "positive": "pixel art giant golden coin fish, round fat fish body, bright gold color, yen symbol on body, large fins, big cute eye, dark outline, 16-bit retro game sprite, solid bright magenta background #FF00FF, centered, clear silhouette",
        "negative": "realistic, 3d, blurry, text, watermark, white background, dark background, gradient background, multiple fish"
    },
}

def build_workflow(positive, negative, seed=None):
    if seed is None:
        seed = int(time.time()) % 1000000
    return {
        "1": {"inputs": {"ckpt_name": "v1-5-pruned-emaonly.safetensors"}, "class_type": "CheckpointLoaderSimple"},
        "2": {"inputs": {"lora_name": "pixel_art_lora.safetensors", "strength_model": 0.9, "strength_clip": 0.9, "model": ["1",0], "clip": ["1",1]}, "class_type": "LoraLoader"},
        "3": {"inputs": {"seed": seed, "steps": 25, "cfg": 7.5, "sampler_name": "euler_ancestral", "scheduler": "normal", "denoise": 1.0, "model": ["2",0], "positive": ["6",0], "negative": ["7",0], "latent_image": ["5",0]}, "class_type": "KSampler"},
        "5": {"inputs": {"width": 512, "height": 512, "batch_size": 1}, "class_type": "EmptyLatentImage"},
        "6": {"inputs": {"text": positive, "clip": ["2",1]}, "class_type": "CLIPTextEncode"},
        "7": {"inputs": {"text": negative, "clip": ["2",1]}, "class_type": "CLIPTextEncode"},
        "8": {"inputs": {"samples": ["3",0], "vae": ["1",2]}, "class_type": "VAEDecode"},
        "9": {"inputs": {"filename_prefix": "target", "images": ["8",0]}, "class_type": "SaveImage"},
    }

def queue_prompt(workflow):
    client_id = str(uuid.uuid4())
    data = json.dumps({"prompt": workflow, "client_id": client_id}).encode("utf-8")
    req = urllib.request.Request(f"{COMFYUI_URL}/prompt", data=data, headers={"Content-Type": "application/json"})
    try:
        with urllib.request.urlopen(req, timeout=30) as r:
            return json.loads(r.read()).get("prompt_id"), client_id
    except Exception as e:
        print(f"  ERROR: {e}"); return None, None

def wait_for_result(prompt_id, timeout=120):
    start = time.time()
    while time.time() - start < timeout:
        try:
            with urllib.request.urlopen(f"{COMFYUI_URL}/history/{prompt_id}", timeout=10) as r:
                history = json.loads(r.read())
                if prompt_id in history:
                    for node_output in history[prompt_id].get("outputs", {}).values():
                        if "images" in node_output:
                            return node_output["images"]
        except Exception:
            pass
        time.sleep(2)
        print(f"\r  Waiting... {int(time.time()-start)}s", end="", flush=True)
    return None

def download_image(filename, subfolder=""):
    params = urllib.parse.urlencode({"filename": filename, "subfolder": subfolder, "type": "output"})
    try:
        with urllib.request.urlopen(f"{COMFYUI_URL}/view?{params}", timeout=30) as r:
            return r.read()
    except Exception as e:
        print(f"  ERROR: {e}"); return None

def post_process(img_data, target_id):
    """後處理：去背 + 縮放到 64x64"""
    import math
    from collections import deque
    from PIL import Image, ImageEnhance
    import io

    img = Image.open(io.BytesIO(img_data)).convert("RGBA")

    # 縮小到 192x192
    if img.size[0] > 192:
        img = img.resize((192, 192), Image.NEAREST)

    # 洋紅色去背
    pixels = img.load()
    w, h = img.size
    def dist_m(r, g, b):
        return math.sqrt((r-255)**2 + g**2 + (b-255)**2)

    for x in range(w):
        for y in range(h):
            r, g, b, a = pixels[x, y]
            if a == 0: continue
            if dist_m(r, g, b) < 100:
                pixels[x, y] = (0, 0, 0, 0)

    visited = set()
    queue = deque()
    for x in range(w):
        queue.append((x, 0)); queue.append((x, h-1))
    for y in range(h):
        queue.append((0, y)); queue.append((w-1, y))
    while queue:
        x, y = queue.popleft()
        if (x, y) in visited or x < 0 or x >= w or y < 0 or y >= h: continue
        visited.add((x, y))
        r, g, b, a = pixels[x, y]
        if a == 0:
            for dx in (-1,0,1):
                for dy2 in (-1,0,1):
                    if dx==0 and dy2==0: continue
                    if (x+dx, y+dy2) not in visited: queue.append((x+dx, y+dy2))
        elif dist_m(r, g, b) < 150:
            pixels[x, y] = (0, 0, 0, 0)
            for dx in (-1,0,1):
                for dy2 in (-1,0,1):
                    if dx==0 and dy2==0: continue
                    if (x+dx, y+dy2) not in visited: queue.append((x+dx, y+dy2))

    # 白色去背（fallback）
    non_t_m = sum(1 for px in img.getdata() if px[3] > 10)
    if non_t_m < 500:
        img2 = Image.open(io.BytesIO(img_data)).convert("RGBA")
        if img2.size[0] > 192:
            img2 = img2.resize((192, 192), Image.NEAREST)
        pixels2 = img2.load()
        w2, h2 = img2.size
        def is_white(px):
            r, g, b, a = px
            return r > 200 and g > 200 and b > 200 and a > 10
        queue2 = deque()
        visited2 = [[False]*h2 for _ in range(w2)]
        for x in range(w2):
            for y in [0, h2-1]:
                if not visited2[x][y] and is_white(pixels2[x, y]):
                    queue2.append((x, y)); visited2[x][y] = True
        for y in range(h2):
            for x in [0, w2-1]:
                if not visited2[x][y] and is_white(pixels2[x, y]):
                    queue2.append((x, y)); visited2[x][y] = True
        while queue2:
            x, y = queue2.popleft()
            pixels2[x, y] = (0, 0, 0, 0)
            for dx, dy in [(0,1),(0,-1),(1,0),(-1,0)]:
                nx, ny = x+dx, y+dy
                if 0<=nx<w2 and 0<=ny<h2 and not visited2[nx][ny] and is_white(pixels2[nx, ny]):
                    visited2[nx][ny] = True; queue2.append((nx, ny))
        img = img2

    # 縮放到 64x64
    bbox = img.getbbox()
    if bbox:
        img = img.crop(bbox)
    cw, ch = img.size
    cell = 64
    fit = 0.85
    scale = min(cell/cw, cell/ch) * fit
    nw = max(1, int(cw*scale))
    nh = max(1, int(ch*scale))
    img = img.resize((nw, nh), Image.NEAREST)
    canvas = Image.new("RGBA", (cell, cell), (0,0,0,0))
    px_x = (cell - nw) // 2
    px_y = (cell - nh) // 2
    canvas.paste(img, (px_x, px_y))

    # 增強
    canvas = ImageEnhance.Color(canvas).enhance(1.3)
    canvas = ImageEnhance.Contrast(canvas).enhance(1.15)
    return canvas

def generate_target(target_id):
    print(f"\n[GENERATE] {target_id}")
    prompts = TARGET_PROMPTS.get(target_id)
    if not prompts:
        print(f"  No prompt for {target_id}"); return False

    workflow = build_workflow(prompts["positive"], prompts["negative"])
    prompt_id, _ = queue_prompt(workflow)
    if not prompt_id:
        print("  Failed to queue"); return False

    print(f"  Prompt ID: {prompt_id}")
    images = wait_for_result(prompt_id)
    if not images:
        print("\n  Timeout"); return False

    print(f"\n  Generated!")
    img_data = download_image(images[0]["filename"], images[0].get("subfolder", ""))
    if not img_data:
        return False

    processed = post_process(img_data, target_id)
    out_path = os.path.join(OUTPUT_DIR, f"{target_id}.png")
    processed.save(out_path)
    non_t = sum(1 for px in processed.getdata() if px[3] > 10)
    print(f"  Saved: {out_path} ({non_t}px)")
    return True

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--target", help="單個目標物 ID")
    parser.add_argument("--all", action="store_true", help="生成全部")
    args = parser.parse_args()

    print("=== ComfyUI Target Generator ===")
    try:
        with urllib.request.urlopen(f"{COMFYUI_URL}/system_stats", timeout=5) as r:
            print("ComfyUI OK")
    except:
        print("ComfyUI NOT running!"); sys.exit(1)

    if args.all:
        targets = list(TARGET_PROMPTS.keys())
        for i, t in enumerate(targets):
            ok = generate_target(t)
            print(f"  Progress: {i+1}/{len(targets)}")
    elif args.target:
        generate_target(args.target)
    else:
        parser.print_help()

if __name__ == "__main__":
    main()
