# -*- coding: utf-8 -*-
"""批次後處理所有 AI 生成的圖片"""
import os
import sys
sys.path.insert(0, 'd:/Kiro/tools')

AI_DIR   = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\ai_generated"
CHARS    = ["chiikawa", "hachiware", "usagi"]
POSES    = ["idle", "attack", "bigwin"]

# 直接 import process_sprites 的函數
import math
from collections import deque
from PIL import Image, ImageEnhance
from pathlib import Path

CHARS_DIR = Path(r"D:\Kiro\client\chiikawa-pixel\assets\sprites\characters")
CELL_SIZE = 96
FIT_SCALE = 0.82
ALIGN     = "bottom"

def remove_bg_magenta(img, threshold=100, edge_threshold=150):
    img = img.convert("RGBA")
    pixels = img.load()
    w, h = img.size
    def dist_m(r, g, b):
        return math.sqrt((r-255)**2 + g**2 + (b-255)**2)
    for x in range(w):
        for y in range(h):
            r, g, b, a = pixels[x, y]
            if a == 0: continue
            if dist_m(r, g, b) < threshold:
                pixels[x, y] = (0, 0, 0, 0)
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
        elif dist_m(r, g, b) < edge_threshold:
            pixels[x, y] = (0, 0, 0, 0)
            for dx in (-1, 0, 1):
                for dy2 in (-1, 0, 1):
                    if dx == 0 and dy2 == 0: continue
                    if (x+dx, y+dy2) not in visited:
                        queue.append((x+dx, y+dy2))
    return img

def remove_bg_white(img, threshold=200):
    img = img.convert("RGBA")
    pixels = img.load()
    w, h = img.size
    def is_white(px):
        r, g, b, a = px
        return r > threshold and g > threshold and b > threshold and a > 10
    queue = deque()
    visited = [[False]*h for _ in range(w)]
    for x in range(w):
        for y in [0, h-1]:
            if not visited[x][y] and is_white(pixels[x, y]):
                queue.append((x, y)); visited[x][y] = True
    for y in range(h):
        for x in [0, w-1]:
            if not visited[x][y] and is_white(pixels[x, y]):
                queue.append((x, y)); visited[x][y] = True
    while queue:
        x, y = queue.popleft()
        pixels[x, y] = (0, 0, 0, 0)
        for dx, dy in [(0,1),(0,-1),(1,0),(-1,0)]:
            nx, ny = x+dx, y+dy
            if 0 <= nx < w and 0 <= ny < h and not visited[nx][ny] and is_white(pixels[nx, ny]):
                visited[nx][ny] = True; queue.append((nx, ny))
    return img

def fit_to_cell(img, cell_size, fit_scale, align):
    bbox = img.getbbox()
    if not bbox:
        return Image.new("RGBA", (cell_size, cell_size), (0,0,0,0))
    cropped = img.crop(bbox)
    cw, ch = cropped.size
    scale = min(cell_size/cw, cell_size/ch) * fit_scale
    nw = max(1, int(cw*scale))
    nh = max(1, int(ch*scale))
    resized = cropped.resize((nw, nh), Image.NEAREST)
    canvas = Image.new("RGBA", (cell_size, cell_size), (0,0,0,0))
    px_x = (cell_size - nw) // 2
    if align == "bottom":
        pad = max(0, int(cell_size*(1-fit_scale)*0.4))
        px_y = cell_size - nh - pad
    else:
        px_y = (cell_size - nh) // 2
    canvas.paste(resized, (px_x, px_y))
    return canvas

def process_ai_image(src_path, char, pose):
    img = Image.open(src_path).convert("RGBA")
    print(f"  原始: {img.size}")

    # 縮小到 192x192 再處理
    if img.size[0] > 192:
        img = img.resize((192, 192), Image.NEAREST)

    # 嘗試洋紅色去背，失敗則用白色去背
    img_magenta = remove_bg_magenta(img.copy())
    bbox_m = img_magenta.getbbox()
    non_t_m = sum(1 for px in img_magenta.getdata() if px[3] > 10) if bbox_m else 0

    img_white = remove_bg_white(img.copy())
    bbox_w = img_white.getbbox()
    non_t_w = sum(1 for px in img_white.getdata() if px[3] > 10) if bbox_w else 0

    # 選擇去背效果更好的（非透明像素更多）
    if non_t_m >= non_t_w:
        img_clean = img_magenta
        print(f"  去背: 洋紅色（{non_t_m}px）")
    else:
        img_clean = img_white
        print(f"  去背: 白色（{non_t_w}px）")

    # 縮放到 96x96
    img_fit = fit_to_cell(img_clean, CELL_SIZE, FIT_SCALE, ALIGN)

    # 增強
    img_fit = ImageEnhance.Color(img_fit).enhance(1.3)
    img_fit = ImageEnhance.Contrast(img_fit).enhance(1.15)

    # 儲存
    out_path = CHARS_DIR / f"{char}_{pose}.png"
    img_fit.save(out_path)

    non_t_final = sum(1 for px in img_fit.getdata() if px[3] > 10)
    print(f"  儲存: {out_path.name} ({non_t_final}px, {non_t_final*100//CELL_SIZE//CELL_SIZE}%)")
    return non_t_final

def main():
    print("=== 批次後處理 AI 生成圖片 ===\n")
    results = []
    for char in CHARS:
        print(f"[{char}]")
        for pose in POSES:
            src = os.path.join(AI_DIR, f"{char}_{pose}.png")
            if not os.path.exists(src):
                print(f"  SKIP: {char}_{pose}.png not found")
                continue
            non_t = process_ai_image(src, char, pose)
            results.append((char, pose, non_t))
        print()

    print("=== 結果摘要 ===")
    for char, pose, non_t in results:
        quality = "✅" if non_t > 3000 else ("⚠️" if non_t > 1000 else "❌")
        print(f"  {quality} {char}_{pose}: {non_t}px")

if __name__ == "__main__":
    main()
