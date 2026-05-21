# -*- coding: utf-8 -*-
"""
Sprite 後處理工具 — 移植自 agent-sprite-forge 的核心技術
功能：
  1. 洋紅色背景去除（比白色背景更可靠，不會誤刪白色毛皮）
  2. shared_scale：所有幀統一縮放比例，解決 attack 幀比 idle 大的問題
  3. component_mode=largest：只保留最大連通區域，去除 FX 噪點
  4. 輸出透明 PNG + 動畫 GIF
  5. 重新對齊現有 sprites（即使沒有 ComfyUI 也能改善品質）

使用方式：
  # 處理現有角色 sprites（重新對齊 + shared_scale）
  py tools/process_sprites.py --mode realign

  # 處理 ComfyUI 生成的原始圖（洋紅色背景）
  py tools/process_sprites.py --mode comfyui --input path/to/raw.png --char chiikawa --pose idle

  # 重建所有 Spritesheet
  py tools/process_sprites.py --mode sheet
"""

from __future__ import annotations

import argparse
import json
import math
import os
import sys
from collections import deque
from pathlib import Path

import numpy as np
from PIL import Image, ImageEnhance

# ── 路徑設定 ──────────────────────────────────────────────────────────────────
SPRITES_DIR   = Path(r"D:\Kiro\client\chiikawa-pixel\assets\sprites")
CHARS_DIR     = SPRITES_DIR / "characters"
SHEETS_DIR    = SPRITES_DIR / "sheets"
AI_GEN_DIR    = SPRITES_DIR / "ai_generated"

CELL_SIZE     = 96   # 輸出每幀大小
FIT_SCALE     = 0.82 # 角色佔 cell 的比例（留邊）
ALIGN         = "bottom"  # 腳底對齊，讓角色站在同一條線上

CHARS  = ["chiikawa", "hachiware", "usagi"]
STATES = ["idle", "attack", "bigwin"]

# ── 核心後處理函數（移植自 agent-sprite-forge）────────────────────────────────

def remove_bg_magenta(img: Image.Image, threshold: int = 100, edge_threshold: int = 150) -> Image.Image:
    """去除洋紅色背景（#FF00FF）— 兩階段：全局 threshold + BFS flood fill"""
    img = img.convert("RGBA")
    pixels = img.load()
    width, height = img.size

    def dist_magenta(r, g, b):
        return math.sqrt((r - 255)**2 + g**2 + (b - 255)**2)

    # 第一階段：全局 threshold
    for x in range(width):
        for y in range(height):
            r, g, b, a = pixels[x, y]
            if a == 0:
                continue
            if dist_magenta(r, g, b) < threshold:
                pixels[x, y] = (0, 0, 0, 0)

    # 第二階段：BFS flood fill 從邊緣擴展
    visited = set()
    queue = deque()
    for x in range(width):
        queue.append((x, 0))
        queue.append((x, height - 1))
    for y in range(height):
        queue.append((0, y))
        queue.append((width - 1, y))

    while queue:
        x, y = queue.popleft()
        if (x, y) in visited or x < 0 or x >= width or y < 0 or y >= height:
            continue
        visited.add((x, y))
        r, g, b, a = pixels[x, y]
        if a == 0:
            for dx in (-1, 0, 1):
                for dy in (-1, 0, 1):
                    if dx == 0 and dy == 0:
                        continue
                    nx, ny = x + dx, y + dy
                    if (nx, ny) not in visited:
                        queue.append((nx, ny))
        elif dist_magenta(r, g, b) < edge_threshold:
            pixels[x, y] = (0, 0, 0, 0)
            for dx in (-1, 0, 1):
                for dy in (-1, 0, 1):
                    if dx == 0 and dy == 0:
                        continue
                    nx, ny = x + dx, y + dy
                    if (nx, ny) not in visited:
                        queue.append((nx, ny))
    return img


def remove_bg_white(img: Image.Image, threshold: int = 200) -> Image.Image:
    """去除白色背景（flood fill，保留角色內部白色）"""
    img = img.convert("RGBA")
    pixels = img.load()
    w, h = img.size

    def is_white(px):
        r, g, b, a = px
        return r > threshold and g > threshold and b > threshold and a > 10

    queue = deque()
    visited = [[False] * h for _ in range(w)]

    for x in range(w):
        for y in [0, h - 1]:
            if not visited[x][y] and is_white(pixels[x, y]):
                queue.append((x, y))
                visited[x][y] = True
    for y in range(h):
        for x in [0, w - 1]:
            if not visited[x][y] and is_white(pixels[x, y]):
                queue.append((x, y))
                visited[x][y] = True

    while queue:
        x, y = queue.popleft()
        pixels[x, y] = (0, 0, 0, 0)
        for dx, dy in [(0, 1), (0, -1), (1, 0), (-1, 0)]:
            nx, ny = x + dx, y + dy
            if 0 <= nx < w and 0 <= ny < h and not visited[nx][ny] and is_white(pixels[nx, ny]):
                visited[nx][ny] = True
                queue.append((nx, ny))
    return img


def connected_components(img: Image.Image, min_area: int = 1) -> list[dict]:
    """找所有連通區域，按面積排序（最大的在前）"""
    alpha = img.getchannel("A")
    pixels = alpha.load()
    width, height = img.size
    visited = [[False] * width for _ in range(height)]
    components = []

    for y in range(height):
        for x in range(width):
            if pixels[x, y] == 0 or visited[y][x]:
                continue
            queue = deque([(x, y)])
            visited[y][x] = True
            area = 0
            min_x = max_x = x
            min_y = max_y = y

            while queue:
                cx, cy = queue.popleft()
                area += 1
                min_x = min(min_x, cx)
                min_y = min(min_y, cy)
                max_x = max(max_x, cx)
                max_y = max(max_y, cy)
                for dx, dy in ((1, 0), (-1, 0), (0, 1), (0, -1)):
                    nx, ny = cx + dx, cy + dy
                    if 0 <= nx < width and 0 <= ny < height and pixels[nx, ny] > 0 and not visited[ny][nx]:
                        visited[ny][nx] = True
                        queue.append((nx, ny))

            if area >= min_area:
                components.append({
                    "area": area,
                    "bbox": (min_x, min_y, max_x + 1, max_y + 1),
                })

    components.sort(key=lambda c: c["area"], reverse=True)
    return components


def keep_largest_component(img: Image.Image, min_area_ratio: float = 0.05) -> Image.Image:
    """
    只保留最大連通區域（去除 FX 噪點）
    min_area_ratio：最大連通區域面積 < 總像素 * ratio 時，跳過（圖片可能壞掉）
    """
    total_pixels = img.width * img.height
    components = connected_components(img, min_area=50)
    if not components:
        return img

    largest = components[0]
    # 如果最大連通區域太小，說明圖片可能壞掉，直接返回原圖
    if largest["area"] < total_pixels * min_area_ratio:
        print(f"  WARNING: largest component only {largest['area']} pixels, skipping cleanup")
        return img

    bbox = largest["bbox"]
    x0, y0 = bbox[0], bbox[1]
    alpha = img.getchannel("A")
    pixels_a = alpha.load()
    width, height = img.size
    visited = [[False] * width for _ in range(height)]

    start = None
    for y in range(y0, bbox[3]):
        for x in range(x0, bbox[2]):
            if pixels_a[x, y] > 0:
                start = (x, y)
                break
        if start:
            break

    if not start:
        return img

    keep_pixels = set()
    queue = deque([start])
    visited[start[1]][start[0]] = True
    while queue:
        cx, cy = queue.popleft()
        keep_pixels.add((cx, cy))
        for dx, dy in ((1, 0), (-1, 0), (0, 1), (0, -1)):
            nx, ny = cx + dx, cy + dy
            if 0 <= nx < width and 0 <= ny < height and pixels_a[nx, ny] > 0 and not visited[ny][nx]:
                visited[ny][nx] = True
                queue.append((nx, ny))

    result = img.copy()
    pixels_r = result.load()
    for y in range(height):
        for x in range(width):
            if pixels_r[x, y][3] > 0 and (x, y) not in keep_pixels:
                pixels_r[x, y] = (0, 0, 0, 0)

    return result


def fit_to_cell(img: Image.Image, cell_size: int, fit_scale: float, align: str = "center") -> Image.Image:
    """把 sprite 縮放並置入 cell，支援 bottom 對齊"""
    bbox = img.getbbox()
    if not bbox:
        return Image.new("RGBA", (cell_size, cell_size), (0, 0, 0, 0))

    cropped = img.crop(bbox)
    cw, ch = cropped.size

    scale = min(cell_size / cw, cell_size / ch) * fit_scale
    new_w = max(1, int(cw * scale))
    new_h = max(1, int(ch * scale))
    resized = cropped.resize((new_w, new_h), Image.NEAREST)

    canvas = Image.new("RGBA", (cell_size, cell_size), (0, 0, 0, 0))
    paste_x = (cell_size - new_w) // 2

    if align == "bottom":
        pad = max(0, int(cell_size * (1 - fit_scale) * 0.4))
        paste_y = cell_size - new_h - pad
    else:
        paste_y = (cell_size - new_h) // 2

    canvas.paste(resized, (paste_x, paste_y))
    return canvas


def process_frames_shared_scale(frames: list[Image.Image], cell_size: int, fit_scale: float, align: str) -> list[Image.Image]:
    """
    shared_scale：所有幀用同一個縮放比例
    基於 idle 幀（第一幀）的大小計算縮放，確保角色大小一致
    解決 attack 幀比 idle 幀大、bigwin 幀被裁切縮小的問題
    """
    if not frames:
        return []

    # 優先用第一幀（idle）的 bbox 計算縮放比例
    # 如果第一幀太小（可能壞掉），才用所有幀的最大值
    ref_bbox = frames[0].getbbox() if frames else None
    if ref_bbox:
        ref_w = ref_bbox[2] - ref_bbox[0]
        ref_h = ref_bbox[3] - ref_bbox[1]
    else:
        ref_w = ref_h = 0

    # 如果 idle 幀太小，fallback 到所有幀最大值
    if ref_w < 20 or ref_h < 20:
        for frame in frames:
            bbox = frame.getbbox()
            if bbox:
                ref_w = max(ref_w, bbox[2] - bbox[0])
                ref_h = max(ref_h, bbox[3] - bbox[1])

    if ref_w == 0 or ref_h == 0:
        return [Image.new("RGBA", (cell_size, cell_size), (0, 0, 0, 0))] * len(frames)

    # 用參考幀計算統一縮放比例
    common_scale = min(cell_size / ref_w, cell_size / ref_h) * fit_scale

    result = []
    for frame in frames:
        bbox = frame.getbbox()
        if not bbox:
            result.append(Image.new("RGBA", (cell_size, cell_size), (0, 0, 0, 0)))
            continue

        cropped = frame.crop(bbox)
        cw, ch = cropped.size
        new_w = max(1, int(cw * common_scale))
        new_h = max(1, int(ch * common_scale))
        resized = cropped.resize((new_w, new_h), Image.NEAREST)

        canvas = Image.new("RGBA", (cell_size, cell_size), (0, 0, 0, 0))
        paste_x = (cell_size - new_w) // 2

        if align == "bottom":
            pad = max(0, int(cell_size * (1 - fit_scale) * 0.4))
            paste_y = cell_size - new_h - pad
        else:
            paste_y = (cell_size - new_h) // 2

        canvas.paste(resized, (paste_x, paste_y))
        result.append(canvas)

    return result


def enhance_pixel_art(img: Image.Image) -> Image.Image:
    """增強像素藝術品質：提升飽和度和對比度"""
    img = ImageEnhance.Color(img).enhance(1.3)
    img = ImageEnhance.Contrast(img).enhance(1.15)
    return img


def save_gif(frames: list[Image.Image], out_path: Path, duration: int = 200) -> None:
    """儲存透明 GIF（移植自 agent-sprite-forge 的 save_transparent_gif）"""
    if not frames:
        return

    key = (255, 0, 254)  # 接近洋紅但不完全相同的 key color
    width, height = frames[0].size
    stacked = Image.new("RGB", (width, height * len(frames)), key)

    for i, frame in enumerate(frames):
        r, g, b, a = frame.split()
        hard_mask = a.point(lambda v: 255 if v >= 128 else 0)
        rgb = Image.merge("RGB", (r, g, b))
        stacked.paste(rgb, (0, i * height), hard_mask)

    paletted = stacked.convert("P", palette=Image.Palette.ADAPTIVE, colors=255, dither=Image.Dither.NONE)
    palette = list(paletted.getpalette() or [])
    while len(palette) < 256 * 3:
        palette.append(0)

    # 找 key color 的 palette index
    key_index = None
    for i in range(256):
        if palette[i*3:i*3+3] == list(key):
            key_index = i
            break
    if key_index is None:
        best_dist = None
        best_i = 0
        for i in range(256):
            r2, g2, b2 = palette[i*3], palette[i*3+1], palette[i*3+2]
            d = (r2-key[0])**2 + (g2-key[1])**2 + (b2-key[2])**2
            if best_dist is None or d < best_dist:
                best_dist = d
                best_i = i
        key_index = best_i

    if key_index != 0:
        arr = np.array(paletted)
        lut = np.arange(256, dtype=np.uint8)
        lut[0], lut[key_index] = key_index, 0
        arr = lut[arr]
        paletted = Image.fromarray(arr, mode="P")
        for ch in range(3):
            palette[ch], palette[key_index*3+ch] = palette[key_index*3+ch], palette[ch]
        paletted.putpalette(palette)

    out_frames = [paletted.crop((0, i*height, width, (i+1)*height)) for i in range(len(frames))]
    out_frames[0].save(
        out_path, format="GIF", save_all=True,
        append_images=out_frames[1:],
        duration=duration, loop=0, disposal=2,
        transparency=0, background=0,
    )
    print(f"  GIF: {out_path.name} ({len(frames)} frames)")


# ── 模式：重新對齊現有 sprites ────────────────────────────────────────────────

def mode_realign():
    """
    重新對齊現有角色 sprites：
    1. 去除白色背景（flood fill）
    2. 用 shared_scale 統一各幀縮放
    3. bottom 對齊（腳底在同一條線）
    4. 增強飽和度/對比度
    5. 輸出 GIF 動畫預覽
    """
    print("=== 重新對齊現有 Sprites（shared_scale + bottom align）===\n")

    gifs_dir = CHARS_DIR / "gifs"
    gifs_dir.mkdir(exist_ok=True)

    for char in CHARS:
        print(f"[{char}]")
        frames_raw = []
        for state in STATES:
            path = CHARS_DIR / f"{char}_{state}.png"
            if not path.exists():
                print(f"  MISSING: {path.name}")
                frames_raw.append(None)
                continue
            img = Image.open(path).convert("RGBA")
            # 如果是 64x64，先放大到 96x96（NEAREST 保持像素感）
            if img.size == (64, 64):
                img = img.resize((96, 96), Image.NEAREST)
            # 去白色背景
            img = remove_bg_white(img)
            # 保留最大連通區域（去噪點）
            img = keep_largest_component(img)
            frames_raw.append(img)

        # 過濾掉 None
        valid_frames = [f for f in frames_raw if f is not None]
        if not valid_frames:
            continue

        # shared_scale：所有幀統一縮放
        processed = process_frames_shared_scale(valid_frames, CELL_SIZE, FIT_SCALE, ALIGN)

        # 增強像素藝術品質
        processed = [enhance_pixel_art(f) for f in processed]

        # 儲存各幀
        for i, (state, frame) in enumerate(zip(STATES, processed)):
            out_path = CHARS_DIR / f"{char}_{state}.png"
            frame.save(out_path)
            bbox = frame.getbbox()
            if bbox:
                w = bbox[2] - bbox[0]
                h = bbox[3] - bbox[1]
                print(f"  {state}: {w}x{h}px (saved)")

        # 儲存 GIF 預覽（idle loop）
        idle_frames = [processed[0]] * 2 + [processed[1]] + [processed[0]] * 2
        save_gif(idle_frames, gifs_dir / f"{char}_idle.gif", duration=250)

        # 儲存 attack GIF
        attack_frames = [processed[0], processed[1], processed[0]]
        save_gif(attack_frames, gifs_dir / f"{char}_attack.gif", duration=120)

        # 儲存 bigwin GIF
        bigwin_frames = [processed[2], processed[0], processed[2], processed[0]]
        save_gif(bigwin_frames, gifs_dir / f"{char}_bigwin.gif", duration=180)

        print()

    print("✅ 重新對齊完成！")
    print(f"GIF 預覽: {gifs_dir}")


# ── 模式：處理 ComfyUI 生成的洋紅色背景圖 ────────────────────────────────────

def mode_comfyui(input_path: str, char: str, pose: str):
    """處理 ComfyUI 生成的原始圖（洋紅色背景）"""
    print(f"=== 處理 ComfyUI 圖片: {input_path} ===")

    img = Image.open(input_path).convert("RGBA")
    print(f"  原始尺寸: {img.size}")

    # 去洋紅色背景
    img = remove_bg_magenta(img)
    print("  去背完成（洋紅色）")

    # 保留最大連通區域
    img = keep_largest_component(img)
    print("  去噪點完成")

    # 縮放到 cell_size
    img = fit_to_cell(img, CELL_SIZE, FIT_SCALE, ALIGN)
    print(f"  縮放到 {CELL_SIZE}x{CELL_SIZE}")

    # 增強
    img = enhance_pixel_art(img)

    # 儲存
    out_path = CHARS_DIR / f"{char}_{pose}.png"
    img.save(out_path)
    print(f"  儲存: {out_path}")

    # 同時存到 ai_generated
    ai_path = AI_GEN_DIR / f"{char}_{pose}.png"
    AI_GEN_DIR.mkdir(exist_ok=True)
    img.save(ai_path)
    print(f"  備份: {ai_path}")


# ── 模式：重建 Spritesheet ────────────────────────────────────────────────────

def mode_sheet():
    """重建角色 Spritesheet（3角色 × 3狀態）"""
    print("=== 重建角色 Spritesheet ===\n")

    SHEETS_DIR.mkdir(exist_ok=True)

    cols = len(STATES)   # 3
    rows = len(CHARS)    # 3
    sheet = Image.new("RGBA", (CELL_SIZE * cols, CELL_SIZE * rows), (0, 0, 0, 0))
    metadata = {
        "cell_size": CELL_SIZE,
        "cols": cols,
        "rows": rows,
        "layout": "row=character, col=state",
        "characters": CHARS,
        "states": STATES,
        "sprites": {}
    }

    missing = []
    for row, char in enumerate(CHARS):
        for col, state in enumerate(STATES):
            path = CHARS_DIR / f"{char}_{state}.png"
            if not path.exists():
                missing.append(f"{char}_{state}")
                continue
            sprite = Image.open(path).convert("RGBA")
            if sprite.size != (CELL_SIZE, CELL_SIZE):
                sprite = sprite.resize((CELL_SIZE, CELL_SIZE), Image.NEAREST)
            sheet.paste(sprite, (col * CELL_SIZE, row * CELL_SIZE))
            metadata["sprites"][f"{char}_{state}"] = {
                "x": col * CELL_SIZE,
                "y": row * CELL_SIZE,
                "w": CELL_SIZE,
                "h": CELL_SIZE,
                "row": row,
                "col": col,
            }
            print(f"  ✅ {char}_{state}")

    if missing:
        print(f"\n  ⚠️  Missing: {missing}")

    sheet_path = SHEETS_DIR / "characters_sheet.png"
    sheet.save(sheet_path)
    print(f"\n  Sheet: {sheet_path} ({sheet.width}x{sheet.height})")

    meta_path = SHEETS_DIR / "characters_sheet.json"
    with open(meta_path, "w", encoding="utf-8") as f:
        json.dump(metadata, f, indent=2, ensure_ascii=False)
    print(f"  Meta:  {meta_path}")

    # 驗證
    loaded = Image.open(sheet_path)
    print(f"\n  驗證: {loaded.size}, mode={loaded.mode}")
    non_transparent = sum(
        1 for x in range(loaded.width) for y in range(loaded.height)
        if loaded.getpixel((x, y))[3] > 10
    )
    total = loaded.width * loaded.height
    print(f"  非透明像素: {non_transparent}/{total} ({non_transparent*100//total}%)")
    print("\n✅ Spritesheet 重建完成！")


# ── 模式：QC 報告 ─────────────────────────────────────────────────────────────

def mode_qc():
    """輸出品質報告"""
    print("=== Sprite QC 報告 ===\n")
    for char in CHARS:
        print(f"[{char}]")
        bboxes = {}
        for state in STATES:
            path = CHARS_DIR / f"{char}_{state}.png"
            if not path.exists():
                print(f"  {state}: MISSING")
                continue
            img = Image.open(path).convert("RGBA")
            bbox = img.getbbox()
            bboxes[state] = bbox
            if bbox:
                w = bbox[2] - bbox[0]
                h = bbox[3] - bbox[1]
                non_t = sum(1 for x in range(img.width) for y in range(img.height) if img.getpixel((x,y))[3] > 10)
                print(f"  {state}: {w}x{h}px, non-transparent={non_t}/{img.width*img.height}")

        if len(bboxes) == 3:
            heights = [bboxes[s][3]-bboxes[s][1] for s in STATES if bboxes.get(s)]
            widths  = [bboxes[s][2]-bboxes[s][0] for s in STATES if bboxes.get(s)]
            h_diff = max(heights) - min(heights)
            w_diff = max(widths) - min(widths)
            status = "✅" if h_diff <= 2 and w_diff <= 4 else "⚠️ "
            print(f"  {status} 一致性: height diff={h_diff}px, width diff={w_diff}px")
        print()


# ── 主程式 ────────────────────────────────────────────────────────────────────

def main():
    parser = argparse.ArgumentParser(description="Sprite 後處理工具（agent-sprite-forge 技術移植）")
    parser.add_argument("--mode", choices=["realign", "comfyui", "sheet", "qc"], required=True,
                        help="realign=重新對齊現有sprites, comfyui=處理ComfyUI輸出, sheet=重建Spritesheet, qc=品質報告")
    parser.add_argument("--input", help="ComfyUI 原始圖路徑（mode=comfyui 時必填）")
    parser.add_argument("--char", choices=CHARS, help="角色名稱（mode=comfyui 時必填）")
    parser.add_argument("--pose", choices=STATES, help="動作名稱（mode=comfyui 時必填）")
    args = parser.parse_args()

    if args.mode == "realign":
        mode_realign()
    elif args.mode == "comfyui":
        if not args.input or not args.char or not args.pose:
            print("ERROR: --mode comfyui 需要 --input, --char, --pose")
            sys.exit(1)
        mode_comfyui(args.input, args.char, args.pose)
    elif args.mode == "sheet":
        mode_sheet()
    elif args.mode == "qc":
        mode_qc()


if __name__ == "__main__":
    main()
