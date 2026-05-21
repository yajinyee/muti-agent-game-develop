# -*- coding: utf-8 -*-
"""生成 BOSS B001 的 AI 圖（96x96）"""
import os, sys, json, time, uuid, math
import urllib.request, urllib.parse
from collections import deque

COMFYUI_URL = "http://127.0.0.1:8188"
OUTPUT_DIR  = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\targets"

BOSS_PROMPT = {
    "positive": "pixel art chiikawa boss monster, white fluffy creature with dark evil aura, red glowing eyes, sinister smile with sharp teeth, dark purple energy surrounding, 16-bit retro game boss sprite, solid magenta background #FF00FF, centered, large imposing figure, dark outline",
    "negative": "realistic, 3d, blurry, text, watermark, white background, cute, friendly, small"
}

def build_workflow(positive, negative, seed=None, size=512):
    if seed is None:
        seed = int(time.time()) % 1000000
    return {
        "1": {"inputs": {"ckpt_name": "v1-5-pruned-emaonly.safetensors"}, "class_type": "CheckpointLoaderSimple"},
        "2": {"inputs": {"lora_name": "pixel_art_lora.safetensors", "strength_model": 0.85, "strength_clip": 0.85, "model": ["1",0], "clip": ["1",1]}, "class_type": "LoraLoader"},
        "3": {"inputs": {"seed": seed, "steps": 28, "cfg": 7.5, "sampler_name": "euler_ancestral", "scheduler": "normal", "denoise": 1.0, "model": ["2",0], "positive": ["6",0], "negative": ["7",0], "latent_image": ["5",0]}, "class_type": "KSampler"},
        "5": {"inputs": {"width": size, "height": size, "batch_size": 1}, "class_type": "EmptyLatentImage"},
        "6": {"inputs": {"text": positive, "clip": ["2",1]}, "class_type": "CLIPTextEncode"},
        "7": {"inputs": {"text": negative, "clip": ["2",1]}, "class_type": "CLIPTextEncode"},
        "8": {"inputs": {"samples": ["3",0], "vae": ["1",2]}, "class_type": "VAEDecode"},
        "9": {"inputs": {"filename_prefix": "boss", "images": ["8",0]}, "class_type": "SaveImage"},
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

def wait_for_result(prompt_id, timeout=150):
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

def post_process_boss(img_data):
    """後處理：去背 + 縮放到 96x96"""
    from PIL import Image, ImageEnhance
    import io

    img = Image.open(io.BytesIO(img_data)).convert("RGBA")
    if img.size[0] > 256:
        img = img.resize((256, 256), Image.NEAREST)

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

    # 縮放到 96x96
    bbox = img.getbbox()
    if bbox:
        img = img.crop(bbox)
    cw, ch = img.size
    cell = 96
    fit = 0.88
    scale = min(cell/cw, cell/ch) * fit
    nw = max(1, int(cw*scale))
    nh = max(1, int(ch*scale))
    img = img.resize((nw, nh), Image.NEAREST)
    canvas = Image.new("RGBA", (cell, cell), (0,0,0,0))
    px_x = (cell - nw) // 2
    px_y = (cell - nh) // 2
    canvas.paste(img, (px_x, px_y))
    canvas = ImageEnhance.Color(canvas).enhance(1.3)
    canvas = ImageEnhance.Contrast(canvas).enhance(1.2)
    return canvas

def main():
    print("=== BOSS AI 生成 ===")
    try:
        with urllib.request.urlopen(f"{COMFYUI_URL}/system_stats", timeout=5) as r:
            print("ComfyUI OK")
    except:
        print("ComfyUI NOT running!"); sys.exit(1)

    workflow = build_workflow(BOSS_PROMPT["positive"], BOSS_PROMPT["negative"], size=512)
    prompt_id, _ = queue_prompt(workflow)
    if not prompt_id:
        print("Failed"); sys.exit(1)

    print(f"  Prompt ID: {prompt_id}")
    images = wait_for_result(prompt_id)
    if not images:
        print("\n  Timeout"); sys.exit(1)

    print(f"\n  Generated!")
    img_data = download_image(images[0]["filename"], images[0].get("subfolder", ""))
    if not img_data:
        sys.exit(1)

    processed = post_process_boss(img_data)
    out_path = os.path.join(OUTPUT_DIR, "B001_boss.png")
    processed.save(out_path)
    non_t = sum(1 for px in processed.getdata() if px[3] > 10)
    print(f"  Saved: {out_path} ({non_t}px, {non_t*100//96//96}%)")

if __name__ == "__main__":
    main()
