# -*- coding: utf-8 -*-
"""
目標物品質提升 v2
1. 修復 T105 金幣魚（魚尾超出邊界，魚身更飽滿）
2. 對所有目標物做後處理增強（飽和度、對比度）
3. 輸出品質報告
"""
from PIL import Image, ImageEnhance
import numpy as np
import os
import math

OUTPUT_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
SIZE = 64


def px(img, x, y, c):
    if 0 <= x < SIZE and 0 <= y < SIZE:
        img.putpixel((x, y), c)


def fill_circle(img, cx, cy, r, color):
    for y in range(max(0, cy-r), min(SIZE, cy+r+1)):
        for x in range(max(0, cx-r), min(SIZE, cx+r+1)):
            if (x-cx)**2 + (y-cy)**2 <= r**2:
                px(img, x, y, color)


def gen_T105_coin_fish_v2():
    """巨大金幣魚 50x — 修復版（魚身更飽滿，魚尾在畫布內）"""
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    GOLD   = (220, 180, 40)
    GOLD_L = (255, 230, 80)
    GOLD_D = (160, 120, 10)
    OUTLINE= (80, 50, 5, 255)
    EYE_W  = (255, 255, 255, 255)
    EYE_B  = (20, 20, 20, 255)
    COIN   = (255, 200, 30)
    SCALE  = (200, 160, 30, 200)

    # 魚身（大橢圓，中心偏左讓魚尾有空間）
    cx, cy = 24, 34
    rx, ry = 22, 18
    for y in range(max(0, cy-ry), min(SIZE, cy+ry+1)):
        for x in range(max(0, cx-rx), min(SIZE, cx+rx+1)):
            if ((x-cx)/rx)**2 + ((y-cy)/ry)**2 <= 1.0:
                nx_ = (x-cx)/rx
                ny_ = (y-cy)/ry
                dot = -(nx_*(-0.7)+ny_*(-0.7))
                r_v = GOLD_L[0] if dot > 0.3 else (GOLD_D[0] if dot < -0.1 else GOLD[0])
                g_v = GOLD_L[1] if dot > 0.3 else (GOLD_D[1] if dot < -0.1 else GOLD[1])
                b_v = GOLD_L[2] if dot > 0.3 else (GOLD_D[2] if dot < -0.1 else GOLD[2])
                px(img, x, y, (r_v, g_v, b_v, 255))

    # 魚身輪廓
    for y in range(max(0, cy-ry-1), min(SIZE, cy+ry+2)):
        for x in range(max(0, cx-rx-1), min(SIZE, cx+rx+2)):
            d = ((x-cx)/rx)**2 + ((y-cy)/ry)**2
            if 0.9 <= d <= 1.15:
                px(img, x, y, OUTLINE)

    # 魚鱗（圓弧）
    for (sx, sy, sr) in [(18, 28, 7), (26, 38, 7), (12, 38, 6)]:
        for angle in range(200, 340, 8):
            rad = math.radians(angle)
            x = int(sx + sr * math.cos(rad))
            y = int(sy + sr * math.sin(rad))
            px(img, x, y, SCALE)

    # 魚尾（扇形，x=46-62）
    for i in range(16):
        tail_x = 46 + i
        spread = 2 + i * 3 // 4
        for j in range(-spread, spread+1):
            if 0 <= tail_x < SIZE and 0 <= cy+j < SIZE:
                alpha = max(0, 255 - i*12)
                r_v = GOLD[0] if abs(j) < spread else GOLD_D[0]
                g_v = GOLD[1] if abs(j) < spread else GOLD_D[1]
                b_v = GOLD[2] if abs(j) < spread else GOLD_D[2]
                px(img, tail_x, cy+j, (r_v, g_v, b_v, alpha))
        if 0 <= tail_x < SIZE:
            if 0 <= cy-spread-1 < SIZE:
                px(img, tail_x, cy-spread-1, OUTLINE)
            if 0 <= cy+spread+1 < SIZE:
                px(img, tail_x, cy+spread+1, OUTLINE)

    # 魚鰭（上方）
    for i in range(8):
        for j in range(i+1):
            px(img, 14+j, cy-ry-i, (*GOLD_L, 200))
        px(img, 14+i, cy-ry-i-1, OUTLINE)

    # 眼睛（左側）
    fill_circle(img, 8, cy-2, 5, EYE_W)
    fill_circle(img, 8, cy-2, 3, EYE_B)
    px(img, 7, cy-3, EYE_W)

    # 嘴巴
    for i in range(4):
        px(img, cx-rx+i, cy+2+i//2, OUTLINE)

    # 金幣符號（身上，大圓）
    coin_cx, coin_cy = 24, 34
    for y in range(coin_cy-8, coin_cy+9):
        for x in range(coin_cx-8, coin_cx+9):
            if ((x-coin_cx)/8)**2 + ((y-coin_cy)/8)**2 <= 1.0:
                px(img, x, y, (*COIN, 255))
    for y in range(coin_cy-9, coin_cy+10):
        for x in range(coin_cx-9, coin_cx+10):
            if 0 <= x < SIZE and 0 <= y < SIZE:
                d = ((x-coin_cx)/8)**2 + ((y-coin_cy)/8)**2
                if 0.9 <= d <= 1.2:
                    px(img, x, y, OUTLINE)
    # ¥ 符號
    for y in range(coin_cy-5, coin_cy+6):
        px(img, coin_cx, y, OUTLINE)
    for x in range(coin_cx-5, coin_cx+6):
        px(img, x, coin_cy-1, OUTLINE)
        px(img, x, coin_cy+1, OUTLINE)

    return img


def enhance_target(img: Image.Image) -> Image.Image:
    """後處理增強：飽和度 + 對比度"""
    img = ImageEnhance.Color(img).enhance(1.35)
    img = ImageEnhance.Contrast(img).enhance(1.2)
    img = ImageEnhance.Brightness(img).enhance(1.05)
    return img


def process_all_targets():
    """對所有目標物做後處理增強"""
    print("=== 目標物品質提升 v2 ===\n")

    targets_dir = OUTPUT_DIR
    target_files = [f for f in os.listdir(targets_dir) if f.endswith('.png') and not f.endswith('.import')]

    results = []
    for fname in sorted(target_files):
        path = os.path.join(targets_dir, fname)
        img = Image.open(path).convert("RGBA")

        # 特殊處理 T105
        if 'T105' in fname:
            print(f"  🔧 重新生成 {fname}（修復魚尾超出邊界）")
            new_img = gen_T105_coin_fish_v2()
            new_img = enhance_target(new_img)
            new_img.save(path)
            arr = np.array(new_img)
            non_t = int((arr[:,:,3] > 10).sum())
            total = new_img.width * new_img.height
            pct = non_t * 100 // total
            print(f"  ✅ {fname}: {new_img.size}, {non_t}px ({pct}%)")
            results.append((fname, pct))
            continue

        # 其他目標物：後處理增強
        enhanced = enhance_target(img)
        enhanced.save(path)
        arr = np.array(enhanced)
        non_t = int((arr[:,:,3] > 10).sum())
        total = enhanced.width * enhanced.height
        pct = non_t * 100 // total
        status = "✅" if pct >= 55 else "⚠️ "
        print(f"  {status} {fname}: {enhanced.size}, {non_t}px ({pct}%)")
        results.append((fname, pct))

    avg = sum(p for _, p in results) // len(results) if results else 0
    print(f"\n  平均非透明像素: {avg}%")
    print(f"  最低: {min(p for _, p in results)}% ({min(results, key=lambda x: x[1])[0]})")
    print(f"  最高: {max(p for _, p in results)}% ({max(results, key=lambda x: x[1])[0]})")

    if avg >= 65:
        print(f"\n  🎉 目標物品質優秀（{avg}%）")
    elif avg >= 55:
        print(f"\n  ✅ 目標物品質良好（{avg}%）")
    else:
        print(f"\n  ⚠️  目標物品質需要改善（{avg}%）")

    return avg


if __name__ == "__main__":
    avg = process_all_targets()
    print(f"\nDone! 平均品質: {avg}%")
