"""
preview_bonus_weeds.py
預覽 Bonus 雜草 Sprites，輸出 4x 放大的合併預覽圖
"""
import os
from PIL import Image

SPRITES_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
OUT_PATH = r"d:\Kiro\docs\bonus_weeds_preview.png"

NAMES = [
    "BG001_weed_normal",
    "BG002_weed_hard",
    "BG003_weed_glow",
    "BG004_weed_gold",
    "BG005_weed_evil",
]

SCALE = 4
PADDING = 8
BG_COLOR = (40, 60, 40, 255)  # 草地背景色

def main():
    imgs = []
    for name in NAMES:
        path = os.path.join(SPRITES_DIR, name + ".png")
        img = Image.open(path).convert("RGBA")
        # 4x 放大（NEAREST 保持像素感）
        scaled = img.resize((img.width * SCALE, img.height * SCALE), Image.NEAREST)
        imgs.append(scaled)

    # 合併成一排
    total_w = sum(img.width for img in imgs) + PADDING * (len(imgs) + 1)
    max_h = max(img.height for img in imgs) + PADDING * 2

    canvas = Image.new("RGBA", (total_w, max_h), BG_COLOR)

    x = PADDING
    for img in imgs:
        y = (max_h - img.height) // 2
        canvas.paste(img, (x, y), img)
        x += img.width + PADDING

    canvas.save(OUT_PATH)
    print(f"✅ 預覽圖已儲存：{OUT_PATH}")
    print(f"   尺寸：{canvas.width}x{canvas.height}px（{SCALE}x 放大）")

if __name__ == "__main__":
    main()
