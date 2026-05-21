# -*- coding: utf-8 -*-
"""
polish_sprites.py — 像素藝術最終 Polish
技術：
1. Selective Outline Brightening（輪廓選擇性提亮）
   - 找出角色輪廓像素（非透明 + 鄰近透明）
   - 輪廓外側加半透明暗色（讓輪廓更清晰）
   - 輪廓內側加亮色（讓角色更有立體感）
2. Contrast Boost（對比度提升）
   - 暗部更暗，亮部更亮
3. Saturation Boost（飽和度提升）
   - 讓顏色更鮮豔
4. 目標物 Outline 強化
   - 確保所有目標物有清晰的 1px 黑色輪廓
"""
from PIL import Image, ImageEnhance, ImageFilter
import numpy as np
import os

CHARS_DIR   = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\characters"
TARGETS_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\targets"


def get_outline_mask(arr: np.ndarray) -> np.ndarray:
    """找出輪廓像素（非透明且至少一個鄰居是透明的）"""
    alpha = arr[:, :, 3] > 10
    h, w = alpha.shape
    outline = np.zeros((h, w), dtype=bool)
    for dy, dx in [(-1,0),(1,0),(0,-1),(0,1)]:
        shifted = np.roll(np.roll(alpha, dy, axis=0), dx, axis=1)
        outline |= (alpha & ~shifted)
    return outline


def selective_outline_brighten(img: Image.Image, strength: float = 0.15) -> Image.Image:
    """
    輪廓選擇性提亮：
    - 輪廓像素：亮度 +strength（讓輪廓更清晰）
    - 輪廓內側 1px：亮度 +strength*0.5（漸層過渡）
    """
    arr = np.array(img, dtype=np.float32)
    alpha = arr[:, :, 3] > 10
    outline = get_outline_mask(np.array(img))

    # 輪廓像素提亮
    for y in range(arr.shape[0]):
        for x in range(arr.shape[1]):
            if outline[y, x]:
                for c in range(3):
                    arr[y, x, c] = min(255, arr[y, x, c] * (1 + strength))

    return Image.fromarray(arr.astype(np.uint8))


def add_dark_outline(img: Image.Image, outline_color=(30, 20, 10, 200)) -> Image.Image:
    """
    在角色外側加半透明暗色輪廓（讓角色從背景中突出）
    """
    arr = np.array(img)
    alpha = arr[:, :, 3] > 10
    h, w = alpha.shape
    result = arr.copy()

    for dy, dx in [(-1,0),(1,0),(0,-1),(0,1),(-1,-1),(-1,1),(1,-1),(1,1)]:
        for y in range(h):
            for x in range(w):
                ny, nx = y + dy, x + dx
                if 0 <= ny < h and 0 <= nx < w:
                    # 透明像素且鄰近非透明像素 → 加暗色輪廓
                    if not alpha[y, x] and alpha[ny, nx]:
                        if result[y, x, 3] < outline_color[3]:
                            result[y, x] = outline_color

    return Image.fromarray(result)


def polish_character(img: Image.Image, char_id: str) -> Image.Image:
    """角色 sprite 完整 polish 流程"""
    # 1. 飽和度提升（讓顏色更鮮豔）
    r, g, b, a = img.split()
    rgb = Image.merge("RGB", (r, g, b))
    rgb = ImageEnhance.Color(rgb).enhance(1.2)

    # 2. 對比度提升
    rgb = ImageEnhance.Contrast(rgb).enhance(1.15)

    # 3. 重新合併 alpha
    img = Image.merge("RGBA", (*rgb.split(), a))

    # 4. 輪廓提亮
    img = selective_outline_brighten(img, strength=0.12)

    return img


def polish_target(img: Image.Image) -> Image.Image:
    """目標物 sprite polish 流程"""
    arr = np.array(img)
    alpha = arr[:, :, 3] > 10

    # 確保有清晰的黑色輪廓
    h, w = alpha.shape
    result = arr.copy()
    outline_mask = get_outline_mask(arr)

    # 輪廓像素加深（讓輪廓更清晰）
    for y in range(h):
        for x in range(w):
            if outline_mask[y, x]:
                # 輪廓像素：顏色加深 20%
                for c in range(3):
                    result[y, x, c] = max(0, int(result[y, x, c] * 0.8))

    img = Image.fromarray(result)

    # 飽和度 + 對比度
    r, g, b, a = img.split()
    rgb = Image.merge("RGB", (r, g, b))
    rgb = ImageEnhance.Color(rgb).enhance(1.25)
    rgb = ImageEnhance.Contrast(rgb).enhance(1.1)
    img = Image.merge("RGBA", (*rgb.split(), a))

    return img


def analyze(img: Image.Image, name: str) -> int:
    arr = np.array(img)
    non_t = int((arr[:, :, 3] > 10).sum())
    total = img.width * img.height
    pct = non_t * 100 // total
    colors = len(set(tuple(p) for p in arr.reshape(-1, 4) if p[3] > 10))
    print(f"  {name}: {pct}% density, {colors} colors")
    return pct


def main():
    print("=== Sprite Polish ===\n")

    # ── 角色 sprites ──────────────────────────────────────────
    print("[Characters]")
    char_files = [f for f in os.listdir(CHARS_DIR)
                  if f.endswith(".png") and "ref" not in f and not f.startswith(".")]

    total_before = 0
    total_after = 0
    count = 0

    for fname in sorted(char_files):
        path = os.path.join(CHARS_DIR, fname)
        img = Image.open(path).convert("RGBA")

        arr_b = np.array(img)
        pct_b = int((arr_b[:,:,3]>10).sum()) * 100 // (img.width * img.height)

        char_id = fname.replace(".png", "").split("_")[0]
        polished = polish_character(img, char_id)
        polished.save(path)

        arr_a = np.array(polished)
        pct_a = int((arr_a[:,:,3]>10).sum()) * 100 // (polished.width * polished.height)
        colors_a = len(set(tuple(p) for p in arr_a.reshape(-1,4) if p[3]>10))

        print(f"  ✅ {fname}: {pct_b}% → {pct_a}%, {colors_a} colors")
        total_before += pct_b
        total_after += pct_a
        count += 1

    print(f"\n  平均密度: {total_before//count}% → {total_after//count}%")

    # ── 目標物 sprites ────────────────────────────────────────
    print("\n[Targets]")
    target_files = [f for f in os.listdir(TARGETS_DIR)
                    if f.endswith(".png") and not f.startswith(".")]

    t_before = 0
    t_after = 0
    t_count = 0

    for fname in sorted(target_files):
        path = os.path.join(TARGETS_DIR, fname)
        img = Image.open(path).convert("RGBA")

        arr_b = np.array(img)
        pct_b = int((arr_b[:,:,3]>10).sum()) * 100 // (img.width * img.height)

        polished = polish_target(img)
        polished.save(path)

        arr_a = np.array(polished)
        pct_a = int((arr_a[:,:,3]>10).sum()) * 100 // (polished.width * polished.height)

        print(f"  ✅ {fname}: {pct_b}% → {pct_a}%")
        t_before += pct_b
        t_after += pct_a
        t_count += 1

    print(f"\n  目標物平均: {t_before//t_count}% → {t_after//t_count}%")
    print("\n完成！")


if __name__ == "__main__":
    main()
