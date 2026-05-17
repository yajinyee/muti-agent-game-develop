# -*- coding: utf-8 -*-
"""
fix_attack_sprites.py
修復 attack 幀品質問題：
1. chiikawa_attack：密度 43% → 60%+，顏色 1742 → 300 以內
2. hachiware_attack：密度 46% → 60%+
策略：
  - 從對應的 idle 幀做變換生成 attack（旋轉 + 光效），確保密度一致
  - 顏色量化到 256 色以內（去除 AI 噪點）
  - 保留原有的 usagi_attack（已修復，57% 合格）
"""
from PIL import Image, ImageEnhance, ImageFilter
import numpy as np
import os

CHARS_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\characters"
SIZE = 96


def quantize_sprite(img: Image.Image, max_colors: int = 256) -> Image.Image:
    """顏色量化：去除 AI 噪點，保留主要顏色"""
    # 分離 alpha
    r, g, b, a = img.split()
    rgb = Image.merge("RGB", (r, g, b))

    # 量化
    quantized = rgb.quantize(colors=max_colors, method=Image.Quantize.FASTOCTREE, dither=0)
    rgb_back = quantized.convert("RGB")

    # 重新合併 alpha
    result = Image.merge("RGBA", (*rgb_back.split(), a))
    return result


def gen_attack_from_idle(char_id: str, attack_color: tuple, rotate_deg: float = -15) -> Image.Image:
    """
    從 idle 幀生成 attack 幀：
    - 旋轉角色（揮棒動作）
    - 加攻擊光效（角色顏色的光暈）
    - 量化顏色
    """
    idle_path = os.path.join(CHARS_DIR, f"{char_id}_idle.png")
    img = Image.open(idle_path).convert("RGBA")
    arr = np.array(img)

    # 找角色 bbox（非透明區域）
    mask = arr[:, :, 3] > 10
    rows = np.any(mask, axis=1)
    cols = np.any(mask, axis=0)
    rmin, rmax = np.where(rows)[0][[0, -1]]
    cmin, cmax = np.where(cols)[0][[0, -1]]

    # 裁切角色區域
    char_crop = img.crop((cmin, rmin, cmax + 1, rmax + 1))

    # 旋轉（模擬揮棒）
    rotated = char_crop.rotate(rotate_deg, expand=False, fillcolor=(0, 0, 0, 0),
                                resample=Image.NEAREST)

    # 稍微提高亮度（攻擊狀態更亮）
    r2, g2, b2, a2 = rotated.split()
    rgb2 = Image.merge("RGB", (r2, g2, b2))
    rgb2 = ImageEnhance.Brightness(rgb2).enhance(1.15)
    rotated = Image.merge("RGBA", (*rgb2.split(), a2))

    # 建立輸出畫布
    canvas = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))

    # 置中貼上旋轉後的角色
    cw, ch = rotated.size
    paste_x = (SIZE - cw) // 2
    paste_y = (SIZE - ch) // 2
    canvas.paste(rotated, (paste_x, paste_y), rotated)

    # 加攻擊光效（右上角光暈，模擬劍氣）
    canvas_arr = np.array(canvas)
    ar, ag, ab = attack_color

    # 在角色右上方加光點
    for y in range(SIZE):
        for x in range(SIZE):
            if canvas_arr[y, x, 3] > 50:
                # 右上角像素加光暈
                if x > SIZE * 0.55 and y < SIZE * 0.45:
                    # 混入攻擊色
                    blend = 0.25
                    canvas_arr[y, x, 0] = min(255, int(canvas_arr[y, x, 0] * (1 - blend) + ar * blend))
                    canvas_arr[y, x, 1] = min(255, int(canvas_arr[y, x, 1] * (1 - blend) + ag * blend))
                    canvas_arr[y, x, 2] = min(255, int(canvas_arr[y, x, 2] * (1 - blend) + ab * blend))

    # 加幾個光點粒子
    import random
    rng = random.Random(42)
    for _ in range(8):
        px = rng.randint(SIZE // 2, SIZE - 8)
        py = rng.randint(4, SIZE // 2)
        size = rng.randint(2, 4)
        alpha = rng.randint(180, 255)
        for dy in range(-size, size + 1):
            for dx in range(-size, size + 1):
                if dx * dx + dy * dy <= size * size:
                    nx, ny = px + dx, py + dy
                    if 0 <= nx < SIZE and 0 <= ny < SIZE:
                        canvas_arr[ny, nx] = [ar, ag, ab, alpha]

    result = Image.fromarray(canvas_arr.astype(np.uint8))

    # 量化顏色（去除噪點）
    result = quantize_sprite(result, max_colors=200)

    return result


def analyze(img: Image.Image, name: str):
    arr = np.array(img)
    non_t = int((arr[:, :, 3] > 10).sum())
    total = img.width * img.height
    colors = len(set(tuple(p) for p in arr.reshape(-1, 4) if p[3] > 10))
    pct = non_t * 100 // total
    status = "✅" if pct >= 55 else "⚠️ "
    print(f"  {status} {name}: {img.size}, {non_t}px ({pct}%), {colors} colors")
    return pct


def main():
    print("=== Attack Sprite 修復 ===\n")

    # 角色攻擊色
    attack_colors = {
        "chiikawa":  (255, 150, 200),  # 粉紅劍氣
        "hachiware": (100, 150, 255),  # 藍色劍氣
    }

    for char_id, color in attack_colors.items():
        print(f"[{char_id}]")
        out_path = os.path.join(CHARS_DIR, f"{char_id}_attack.png")

        # 分析原始
        orig = Image.open(out_path).convert("RGBA")
        analyze(orig, f"{char_id}_attack (原始)")

        # 生成新版
        new_img = gen_attack_from_idle(char_id, color)
        new_img.save(out_path)

        # 分析新版
        analyze(new_img, f"{char_id}_attack (修復後)")
        print()

    # 同時對所有 sprites 做顏色量化（去除 AI 噪點）
    print("[顏色量化 — 所有角色 sprites]")
    for fname in sorted(os.listdir(CHARS_DIR)):
        if not fname.endswith(".png") or fname.startswith(".") or "ref" in fname:
            continue
        if "attack" in fname:
            continue  # 已處理
        path = os.path.join(CHARS_DIR, fname)
        img = Image.open(path).convert("RGBA")
        arr = np.array(img)
        colors_before = len(set(tuple(p) for p in arr.reshape(-1, 4) if p[3] > 10))

        if colors_before > 400:
            quantized = quantize_sprite(img, max_colors=256)
            quantized.save(path)
            arr2 = np.array(quantized)
            colors_after = len(set(tuple(p) for p in arr2.reshape(-1, 4) if p[3] > 10))
            print(f"  🔧 {fname}: {colors_before} → {colors_after} colors")
        else:
            print(f"  ✅ {fname}: {colors_before} colors (OK)")

    print("\n完成！")


if __name__ == "__main__":
    main()
