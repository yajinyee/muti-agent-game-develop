"""
generate_swim_animation.py
為目標物生成 2 幀游泳動畫 spritesheet

技術：
- 幀 0：原始圖（向上彎曲）
- 幀 1：垂直翻轉 + 輕微縮放（向下彎曲）
- 輸出：{target_id}_swim.png（128x64，2幀橫排）

用途：讓目標物有真正的幀動畫，比 Tween 位移更有生命感
"""

import os
import sys
from PIL import Image, ImageEnhance, ImageFilter
import numpy as np

# 目標物清單（不包含 BOSS 和 Bonus 雜草）
TARGETS = [
    "T001_grass",
    "T002_bug_g",
    "T003_bug_r",
    "T004_bug_b",
    "T005_pudding",
    "T006_mushroom",
    "T101_mimic",
    "T102_chest",
    "T103_meteor",
    "T104_gold_grass",
    "T105_coin_fish",
]

SPRITE_DIR = "client/chiikawa-pixel/assets/sprites/targets"
OUTPUT_DIR = "client/chiikawa-pixel/assets/sprites/targets"
FRAME_SIZE = 64


def generate_swim_frame1(img: Image.Image) -> Image.Image:
    """
    幀 1：模擬游泳向下彎曲
    技術：
    - 上半部分向右偏移 2px（頭部方向）
    - 下半部分向左偏移 2px（尾部方向）
    - 輕微亮度提升（模擬水面反光）
    """
    w, h = img.size
    arr = np.array(img, dtype=np.float32)
    result = np.zeros_like(arr)

    # 分段位移（模擬魚身彎曲）
    for y in range(h):
        # 計算位移量：上半部分向右，下半部分向左
        # 使用 sin 曲線讓彎曲更自然
        t = y / h  # 0.0 (頂部) → 1.0 (底部)
        dx = int(2.0 * np.sin(t * np.pi))  # 最大位移 2px，中間最大

        for x in range(w):
            src_x = x - dx
            if 0 <= src_x < w:
                result[y, x] = arr[y, src_x]
            # 超出邊界的像素保持透明（0）

    # 轉回 PIL Image
    frame1 = Image.fromarray(result.astype(np.uint8), 'RGBA')

    # 輕微亮度提升（模擬水面反光）
    enhancer = ImageEnhance.Brightness(frame1)
    frame1 = enhancer.enhance(1.05)

    return frame1


def generate_swim_frame2(img: Image.Image) -> Image.Image:
    """
    幀 2：模擬游泳向上彎曲（與幀 1 相反方向）
    技術：
    - 上半部分向左偏移 2px
    - 下半部分向右偏移 2px
    - 輕微暗度（模擬水下陰影）
    """
    w, h = img.size
    arr = np.array(img, dtype=np.float32)
    result = np.zeros_like(arr)

    for y in range(h):
        t = y / h
        dx = -int(2.0 * np.sin(t * np.pi))  # 反向位移

        for x in range(w):
            src_x = x - dx
            if 0 <= src_x < w:
                result[y, x] = arr[y, src_x]

    frame2 = Image.fromarray(result.astype(np.uint8), 'RGBA')

    # 輕微暗度
    enhancer = ImageEnhance.Brightness(frame2)
    frame2 = enhancer.enhance(0.95)

    return frame2


def create_swim_sheet(target_id: str) -> bool:
    """為單個目標物生成 2 幀游泳動畫 spritesheet"""
    src_path = os.path.join(SPRITE_DIR, f"{target_id}.png")
    if not os.path.exists(src_path):
        print(f"  ⚠️  {target_id}.png 不存在，跳過")
        return False

    img = Image.open(src_path).convert('RGBA')
    if img.size != (FRAME_SIZE, FRAME_SIZE):
        img = img.resize((FRAME_SIZE, FRAME_SIZE), Image.NEAREST)

    # 生成兩幀
    frame1 = generate_swim_frame1(img)
    frame2 = generate_swim_frame2(img)

    # 合成 spritesheet（2幀橫排，128x64）
    sheet = Image.new('RGBA', (FRAME_SIZE * 2, FRAME_SIZE), (0, 0, 0, 0))
    sheet.paste(frame1, (0, 0))
    sheet.paste(frame2, (FRAME_SIZE, 0))

    # 儲存
    out_path = os.path.join(OUTPUT_DIR, f"{target_id}_swim.png")
    sheet.save(out_path, 'PNG')
    print(f"  ✅ {target_id}_swim.png ({FRAME_SIZE*2}x{FRAME_SIZE})")
    return True


def main():
    print("=" * 60)
    print("目標物游泳動畫生成器")
    print("=" * 60)

    success = 0
    for target_id in TARGETS:
        print(f"\n處理 {target_id}...")
        if create_swim_sheet(target_id):
            success += 1

    print(f"\n完成！成功生成 {success}/{len(TARGETS)} 個游泳動畫")
    print(f"輸出目錄：{OUTPUT_DIR}")
    print("\n下一步：")
    print("1. 在 TargetManager.gd 中載入 *_swim.png")
    print("2. 用 AtlasTexture 切換幀（幀0 和 幀1 交替）")
    print("3. 游泳 FPS：4fps（每 0.25 秒切換一次）")


if __name__ == "__main__":
    main()
