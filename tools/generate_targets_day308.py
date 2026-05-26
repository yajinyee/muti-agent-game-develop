#!/usr/bin/env python3
"""
DAY-308 目標物精靈圖生成
T151-T155：覺醒鱷魚、吸血鬼升級魚、超級覺醒魚、巨型獎勵魚、不死 BOSS 魚
"""
import os
from PIL import Image, ImageDraw

OUT_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
SIZE = 64

def new_img():
    return Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))

def px(img, x, y, c):
    if 0 <= x < SIZE and 0 <= y < SIZE:
        img.putpixel((x, y), c)

def fill_circle(img, cx, cy, r, color):
    for dy in range(-r, r+1):
        for dx in range(-r, r+1):
            if dx*dx + dy*dy <= r*r:
                px(img, cx+dx, cy+dy, color)

def fill_ellipse(img, cx, cy, rx, ry, color):
    for dy in range(-ry, ry+1):
        for dx in range(-rx, rx+1):
            if (dx/rx)**2 + (dy/ry)**2 <= 1.0:
                px(img, cx+dx, cy+dy, color)

def draw_eye(img, ex, ey, color=(20, 20, 20, 255)):
    fill_circle(img, ex, ey, 2, (255, 255, 255, 255))
    px(img, ex, ey, color)
    px(img, ex-1, ey-1, (255, 255, 255, 200))

def save(img, name):
    path = os.path.join(OUT_DIR, name)
    img.save(path)
    non_transparent = sum(1 for x in range(SIZE) for y in range(SIZE) if img.getpixel((x, y))[3] > 0)
    pct = non_transparent / (SIZE * SIZE) * 100
    print(f"  ✅ {name}: {non_transparent} px ({pct:.1f}%)")

# ── T151 幸運覺醒鱷魚 ─────────────────────────────────────────
def gen_t151():
    img = new_img()
    # 鱷魚身體（深綠橢圓）
    BODY = (0, 120, 40, 255)
    DARK = (0, 70, 20, 255)
    LIGHT = (0, 180, 60, 255)
    EYE = (255, 50, 50, 255)
    TOOTH = (255, 255, 220, 255)
    # 身體
    fill_ellipse(img, 32, 36, 22, 14, BODY)
    # 陰影
    fill_ellipse(img, 36, 40, 18, 10, DARK)
    # 高光
    fill_ellipse(img, 26, 30, 10, 6, LIGHT)
    # 頭部
    fill_ellipse(img, 20, 30, 12, 10, BODY)
    # 嘴巴（大嘴）
    for x in range(10, 30):
        px(img, x, 36, DARK)
    # 牙齒
    for i, tx in enumerate([12, 16, 20, 24]):
        px(img, tx, 35, TOOTH)
        px(img, tx, 37, TOOTH)
    # 眼睛（紅色）
    draw_eye(img, 16, 24, EYE)
    draw_eye(img, 24, 24, EYE)
    # 尾巴
    fill_ellipse(img, 52, 36, 8, 5, BODY)
    fill_ellipse(img, 58, 32, 5, 3, DARK)
    # 腳
    fill_ellipse(img, 22, 48, 5, 3, DARK)
    fill_ellipse(img, 38, 48, 5, 3, DARK)
    # 覺醒光環（綠色）
    for angle_deg in range(0, 360, 30):
        import math
        angle = math.radians(angle_deg)
        gx = int(32 + 28 * math.cos(angle))
        gy = int(36 + 20 * math.sin(angle))
        fill_circle(img, gx, gy, 2, (0, 255, 100, 180))
    save(img, "T151_awakened_croc.png")

# ── T152 幸運吸血鬼升級魚 ─────────────────────────────────────
def gen_t152():
    img = new_img()
    BODY = (80, 0, 120, 255)
    DARK = (40, 0, 60, 255)
    LIGHT = (140, 0, 200, 255)
    EYE = (255, 0, 0, 255)
    FANG = (255, 255, 255, 255)
    # 魚身（深紫橢圓）
    fill_ellipse(img, 32, 34, 20, 14, BODY)
    fill_ellipse(img, 36, 38, 16, 10, DARK)
    fill_ellipse(img, 26, 28, 8, 5, LIGHT)
    # 頭部
    fill_ellipse(img, 18, 30, 10, 9, BODY)
    # 紅眼
    draw_eye(img, 14, 26, EYE)
    draw_eye(img, 22, 26, EYE)
    # 獠牙
    px(img, 14, 34, FANG)
    px(img, 14, 35, FANG)
    px(img, 22, 34, FANG)
    px(img, 22, 35, FANG)
    # 魚尾
    fill_ellipse(img, 50, 34, 10, 6, BODY)
    fill_ellipse(img, 56, 28, 6, 4, DARK)
    fill_ellipse(img, 56, 40, 6, 4, DARK)
    # 吸血光環（紫色脈動）
    import math
    for i in range(8):
        angle = math.radians(i * 45)
        gx = int(32 + 26 * math.cos(angle))
        gy = int(34 + 18 * math.sin(angle))
        fill_circle(img, gx, gy, 2, (180, 0, 255, 200))
    # ×10 標記
    for x, y in [(44, 10), (45, 10), (46, 10), (44, 12), (45, 12), (46, 12)]:
        px(img, x, y, (255, 200, 0, 255))
    save(img, "T152_vampire_v2.png")

# ── T153 幸運超級覺醒魚 ───────────────────────────────────────
def gen_t153():
    img = new_img()
    BODY = (200, 80, 0, 255)
    DARK = (120, 40, 0, 255)
    LIGHT = (255, 160, 0, 255)
    EYE = (255, 255, 0, 255)
    AURA = (255, 120, 0, 200)
    # 大型魚身（火橙）
    fill_ellipse(img, 32, 34, 22, 16, BODY)
    fill_ellipse(img, 36, 38, 18, 12, DARK)
    fill_ellipse(img, 24, 26, 10, 7, LIGHT)
    # 頭部
    fill_ellipse(img, 16, 30, 12, 10, BODY)
    # 黃色眼睛
    draw_eye(img, 12, 26, EYE)
    draw_eye(img, 20, 26, EYE)
    # 魚尾
    fill_ellipse(img, 52, 34, 10, 7, BODY)
    fill_ellipse(img, 58, 28, 6, 4, DARK)
    fill_ellipse(img, 58, 40, 6, 4, DARK)
    # 超級覺醒光芒（12方向）
    import math
    for i in range(12):
        angle = math.radians(i * 30)
        for r in range(20, 30):
            gx = int(32 + r * math.cos(angle))
            gy = int(34 + r * math.sin(angle))
            if 0 <= gx < SIZE and 0 <= gy < SIZE:
                px(img, gx, gy, AURA)
    # 中心爆炸
    fill_circle(img, 32, 34, 5, (255, 200, 0, 255))
    save(img, "T153_super_awaken.png")

# ── T154 幸運巨型獎勵魚 ───────────────────────────────────────
def gen_t154():
    img = new_img()
    BODY = (200, 160, 0, 255)
    DARK = (120, 90, 0, 255)
    LIGHT = (255, 220, 50, 255)
    EYE = (0, 0, 0, 255)
    STAR = (255, 255, 100, 255)
    # 大型金色魚身
    fill_ellipse(img, 32, 34, 22, 16, BODY)
    fill_ellipse(img, 36, 38, 18, 12, DARK)
    fill_ellipse(img, 24, 26, 10, 7, LIGHT)
    # 頭部
    fill_ellipse(img, 16, 30, 12, 10, BODY)
    # 眼睛
    draw_eye(img, 12, 26, EYE)
    draw_eye(img, 20, 26, EYE)
    # 魚尾
    fill_ellipse(img, 52, 34, 10, 7, BODY)
    fill_ellipse(img, 58, 28, 6, 4, DARK)
    fill_ellipse(img, 58, 40, 6, 4, DARK)
    # 5個星星（代表5次大獎）
    star_positions = [(10, 8), (20, 5), (32, 4), (44, 5), (54, 8)]
    for sx, sy in star_positions:
        fill_circle(img, sx, sy, 3, STAR)
        px(img, sx, sy-4, STAR)
        px(img, sx, sy+4, STAR)
        px(img, sx-4, sy, STAR)
        px(img, sx+4, sy, STAR)
    # 金色光環
    import math
    for i in range(8):
        angle = math.radians(i * 45)
        gx = int(32 + 26 * math.cos(angle))
        gy = int(34 + 18 * math.sin(angle))
        fill_circle(img, gx, gy, 2, (255, 200, 0, 180))
    save(img, "T154_giant_prize.png")

# ── T155 幸運不死 BOSS 魚 ─────────────────────────────────────
def gen_t155():
    img = new_img()
    BODY = (160, 20, 20, 255)
    DARK = (80, 0, 0, 255)
    LIGHT = (220, 60, 60, 255)
    EYE = (255, 200, 0, 255)
    SKULL = (255, 255, 255, 255)
    # 大型深紅魚身
    fill_ellipse(img, 32, 34, 24, 17, BODY)
    fill_ellipse(img, 36, 38, 20, 13, DARK)
    fill_ellipse(img, 24, 26, 12, 8, LIGHT)
    # 頭部
    fill_ellipse(img, 16, 30, 13, 11, BODY)
    # 金色眼睛
    draw_eye(img, 11, 26, EYE)
    draw_eye(img, 21, 26, EYE)
    # 魚尾
    fill_ellipse(img, 54, 34, 10, 7, BODY)
    fill_ellipse(img, 60, 28, 6, 4, DARK)
    fill_ellipse(img, 60, 40, 6, 4, DARK)
    # 骷髏標記（5條命）
    skull_x, skull_y = 32, 10
    fill_circle(img, skull_x, skull_y, 5, SKULL)
    px(img, skull_x-2, skull_y+1, DARK)
    px(img, skull_x+2, skull_y+1, DARK)
    px(img, skull_x, skull_y+3, DARK)
    # 5個心形（生命）
    for i in range(5):
        hx = 8 + i * 12
        hy = 56
        fill_circle(img, hx-1, hy-1, 2, (255, 50, 50, 255))
        fill_circle(img, hx+1, hy-1, 2, (255, 50, 50, 255))
        fill_circle(img, hx, hy+1, 2, (255, 50, 50, 255))
    # 不死光環（深紅）
    import math
    for i in range(8):
        angle = math.radians(i * 45)
        gx = int(32 + 28 * math.cos(angle))
        gy = int(34 + 20 * math.sin(angle))
        fill_circle(img, gx, gy, 2, (200, 0, 0, 200))
    save(img, "T155_immortal_boss.png")

if __name__ == "__main__":
    print("生成 DAY-308 目標物精靈圖...")
    gen_t151()
    gen_t152()
    gen_t153()
    gen_t154()
    gen_t155()
    print("完成！")
