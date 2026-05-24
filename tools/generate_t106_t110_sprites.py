"""
generate_t106_t110_sprites.py — DAY-292 新增特殊目標物精靈圖生成
T106 幸運連鎖閃電魚（青藍色，閃電紋路）
T107 幸運螃蟹魚雷（橙紅色，螃蟹形狀）
T108 幸運渦旋海葵（紫色，渦旋觸手）
T109 幸運黃金龍魚（金色，龍鱗紋路）
T110 幸運雷霆龍蝦（火紅色，龍蝦形狀）
"""
import os
from PIL import Image, ImageDraw

OUT_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
SIZE = 64

def px(img, x, y, color):
    if 0 <= x < SIZE and 0 <= y < SIZE:
        img.putpixel((x, y), color)

def fill_circle(img, cx, cy, r, color):
    for y in range(max(0, cy-r), min(SIZE, cy+r+1)):
        for x in range(max(0, cx-r), min(SIZE, cx+r+1)):
            if (x-cx)**2 + (y-cy)**2 <= r*r:
                img.putpixel((x, y), color)

def fill_circle_shaded(img, cx, cy, r, light, mid, dark):
    for y in range(max(0, cy-r), min(SIZE, cy+r+1)):
        for x in range(max(0, cx-r), min(SIZE, cx+r+1)):
            if (x-cx)**2 + (y-cy)**2 <= r*r:
                if x < cx and y < cy:
                    img.putpixel((x, y), light)
                elif x > cx and y > cy:
                    img.putpixel((x, y), dark)
                else:
                    img.putpixel((x, y), mid)

def fill_ellipse(img, cx, cy, rx, ry, color):
    for y in range(max(0, cy-ry), min(SIZE, cy+ry+1)):
        for x in range(max(0, cx-rx), min(SIZE, cx+rx+1)):
            if ((x-cx)/rx)**2 + ((y-cy)/ry)**2 <= 1.0:
                img.putpixel((x, y), color)

def draw_outline(img, color=(41, 42, 43, 255)):
    """為非透明像素加輪廓"""
    pixels = img.load()
    w, h = img.size
    outline_pixels = set()
    for y in range(h):
        for x in range(w):
            if pixels[x, y][3] > 0:
                for dx, dy in [(-1,0),(1,0),(0,-1),(0,1)]:
                    nx, ny = x+dx, y+dy
                    if 0 <= nx < w and 0 <= ny < h and pixels[nx, ny][3] == 0:
                        outline_pixels.add((nx, ny))
    for (x, y) in outline_pixels:
        img.putpixel((x, y), color)

# ── T106 幸運連鎖閃電魚 ───────────────────────────────────────
def gen_t106():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    # 魚身（青藍漸層橢圓）
    LIGHT = (0, 255, 255, 255)   # 亮青
    MID   = (0, 180, 220, 255)   # 中青藍
    DARK  = (0, 100, 160, 255)   # 深藍
    fill_ellipse(img, 32, 32, 22, 14, MID)
    # 陰影
    for y in range(20, 46):
        for x in range(10, 54):
            if ((x-32)/22)**2 + ((y-32)/14)**2 <= 1.0:
                if x < 28 and y < 28:
                    img.putpixel((x, y), LIGHT)
                elif x > 36 and y > 36:
                    img.putpixel((x, y), DARK)
    # 閃電紋路（Z字形）
    BOLT = (255, 255, 0, 255)  # 黃色閃電
    for i in range(8):
        px(img, 28+i, 26, BOLT)
    for i in range(6):
        px(img, 34-i, 32, BOLT)
    for i in range(8):
        px(img, 28+i, 38, BOLT)
    # 閃電光暈
    GLOW = (180, 255, 255, 180)
    for i in range(6):
        px(img, 26+i, 25, GLOW)
        px(img, 26+i, 39, GLOW)
    # 魚眼（白色+黑色瞳孔）
    fill_circle(img, 46, 28, 4, (255, 255, 255, 255))
    fill_circle(img, 47, 28, 2, (20, 20, 20, 255))
    px(img, 46, 27, (255, 255, 255, 255))  # 高光
    # 魚尾（青藍三角）
    for i in range(10):
        for j in range(-i, i+1):
            px(img, 8+i, 32+j, MID)
    # 光環（外圈青藍）
    RING = (0, 220, 255, 120)
    for angle_deg in range(0, 360, 8):
        import math
        angle = math.radians(angle_deg)
        rx, ry = int(32 + 26*math.cos(angle)), int(32 + 16*math.sin(angle))
        if 0 <= rx < SIZE and 0 <= ry < SIZE:
            img.putpixel((rx, ry), RING)
    draw_outline(img)
    return img

# ── T107 幸運螃蟹魚雷 ─────────────────────────────────────────
def gen_t107():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    # 螃蟹身體（橙紅圓形）
    LIGHT = (255, 140, 80, 255)
    MID   = (220, 80, 30, 255)
    DARK  = (160, 40, 10, 255)
    fill_circle_shaded(img, 32, 34, 16, LIGHT, MID, DARK)
    # 螃蟹殼紋路
    SHELL = (180, 50, 10, 255)
    for i in range(5):
        px(img, 24+i*2, 30, SHELL)
        px(img, 24+i*2, 36, SHELL)
    # 螃蟹大螯（左右各一）
    CLAW = (200, 70, 20, 255)
    # 左螯
    fill_ellipse(img, 12, 28, 8, 5, CLAW)
    fill_circle(img, 8, 26, 4, CLAW)
    # 右螯
    fill_ellipse(img, 52, 28, 8, 5, CLAW)
    fill_circle(img, 56, 26, 4, CLAW)
    # 螃蟹腳（6條）
    LEG = (180, 60, 15, 255)
    for i in range(3):
        # 左腳
        for j in range(6):
            px(img, 18-j, 30+i*4, LEG)
        # 右腳
        for j in range(6):
            px(img, 46+j, 30+i*4, LEG)
    # 魚雷尾焰
    FLAME = (255, 200, 50, 255)
    for i in range(8):
        px(img, 8-i, 34, FLAME)
        if i < 5:
            px(img, 8-i, 33, FLAME)
            px(img, 8-i, 35, FLAME)
    # 眼睛
    fill_circle(img, 38, 28, 3, (255, 255, 255, 255))
    fill_circle(img, 39, 28, 2, (20, 20, 20, 255))
    px(img, 38, 27, (255, 255, 255, 255))
    draw_outline(img)
    return img

# ── T108 幸運渦旋海葵 ─────────────────────────────────────────
def gen_t108():
    import math
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    # 海葵身體（紫色橢圓）
    LIGHT = (200, 150, 255, 255)
    MID   = (140, 80, 200, 255)
    DARK  = (80, 30, 140, 255)
    fill_ellipse(img, 32, 38, 14, 10, MID)
    # 陰影
    for y in range(28, 50):
        for x in range(18, 46):
            if ((x-32)/14)**2 + ((y-38)/10)**2 <= 1.0:
                if x < 28:
                    img.putpixel((x, y), LIGHT)
                elif x > 36:
                    img.putpixel((x, y), DARK)
    # 渦旋觸手（8條，螺旋狀）
    TENTACLE = (180, 100, 240, 255)
    TENTACLE_TIP = (255, 200, 255, 255)
    for i in range(8):
        angle = math.radians(i * 45)
        for r in range(6, 22):
            spiral_angle = angle + r * 0.15
            tx = int(32 + r * math.cos(spiral_angle))
            ty = int(32 + r * math.sin(spiral_angle))
            if 0 <= tx < SIZE and 0 <= ty < SIZE:
                color = TENTACLE_TIP if r > 18 else TENTACLE
                img.putpixel((tx, ty), color)
    # 渦旋中心（白色光點）
    fill_circle(img, 32, 32, 5, (220, 180, 255, 255))
    fill_circle(img, 32, 32, 3, (255, 255, 255, 255))
    # 渦旋光環
    RING = (180, 100, 255, 100)
    for angle_deg in range(0, 360, 6):
        angle = math.radians(angle_deg)
        rx, ry = int(32 + 28*math.cos(angle)), int(32 + 28*math.sin(angle))
        if 0 <= rx < SIZE and 0 <= ry < SIZE:
            img.putpixel((rx, ry), RING)
    draw_outline(img)
    return img

# ── T109 幸運黃金龍魚 ─────────────────────────────────────────
def gen_t109():
    import math
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    # 龍魚身體（金色橢圓）
    LIGHT = (255, 240, 150, 255)
    MID   = (220, 180, 50, 255)
    DARK  = (160, 120, 20, 255)
    fill_ellipse(img, 32, 32, 24, 14, MID)
    # 陰影
    for y in range(18, 46):
        for x in range(8, 56):
            if ((x-32)/24)**2 + ((y-32)/14)**2 <= 1.0:
                if x < 26 and y < 26:
                    img.putpixel((x, y), LIGHT)
                elif x > 38 and y > 38:
                    img.putpixel((x, y), DARK)
    # 龍鱗紋路
    SCALE = (180, 130, 20, 255)
    for row in range(3):
        for col in range(5):
            sx = 14 + col*8 + (row%2)*4
            sy = 24 + row*6
            fill_circle(img, sx, sy, 3, SCALE)
    # 龍角（金色）
    HORN = (255, 220, 80, 255)
    for i in range(6):
        px(img, 44+i, 20-i, HORN)
        px(img, 44+i, 22-i, HORN)
    # 龍鬚（金色細線）
    WHISKER = (255, 200, 50, 200)
    for i in range(10):
        px(img, 52+i//2, 28+i, WHISKER)
        px(img, 52+i//2, 36-i, WHISKER)
    # 魚眼（金色+黑瞳）
    fill_circle(img, 48, 28, 5, (255, 230, 100, 255))
    fill_circle(img, 49, 28, 3, (20, 20, 20, 255))
    px(img, 48, 27, (255, 255, 255, 255))
    # 魚尾（金色扇形）
    for i in range(12):
        for j in range(-i//2, i//2+1):
            px(img, 6+i, 32+j, MID)
    # 金色光環
    RING = (255, 220, 50, 120)
    for angle_deg in range(0, 360, 6):
        angle = math.radians(angle_deg)
        rx, ry = int(32 + 30*math.cos(angle)), int(32 + 18*math.sin(angle))
        if 0 <= rx < SIZE and 0 <= ry < SIZE:
            img.putpixel((rx, ry), RING)
    # 4方向光芒
    BEAM = (255, 240, 100, 200)
    for i in range(8):
        px(img, 32, 2+i, BEAM)
        px(img, 32, 62-i, BEAM)
        px(img, 2+i, 32, BEAM)
        px(img, 62-i, 32, BEAM)
    draw_outline(img)
    return img

# ── T110 幸運雷霆龍蝦 ─────────────────────────────────────────
def gen_t110():
    import math
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    # 龍蝦身體（火紅橢圓）
    LIGHT = (255, 120, 80, 255)
    MID   = (220, 60, 20, 255)
    DARK  = (150, 20, 5, 255)
    fill_ellipse(img, 32, 34, 18, 12, MID)
    # 陰影
    for y in range(22, 46):
        for x in range(14, 50):
            if ((x-32)/18)**2 + ((y-34)/12)**2 <= 1.0:
                if x < 26 and y < 28:
                    img.putpixel((x, y), LIGHT)
                elif x > 38 and y > 40:
                    img.putpixel((x, y), DARK)
    # 龍蝦殼節（分段）
    SEGMENT = (180, 40, 10, 255)
    for i in range(4):
        for j in range(-8, 9):
            px(img, 20+i*6, 34+j, SEGMENT)
    # 龍蝦大螯（左右）
    CLAW = (200, 50, 15, 255)
    fill_ellipse(img, 10, 26, 8, 6, CLAW)
    fill_circle(img, 6, 24, 5, CLAW)
    fill_ellipse(img, 54, 26, 8, 6, CLAW)
    fill_circle(img, 58, 24, 5, CLAW)
    # 龍蝦腳（8條）
    LEG = (180, 50, 10, 255)
    for i in range(4):
        for j in range(5):
            px(img, 18-j, 28+i*4, LEG)
            px(img, 46+j, 28+i*4, LEG)
    # 龍蝦觸角（長）
    ANTENNA = (255, 100, 50, 255)
    for i in range(14):
        px(img, 46+i, 20-i//2, ANTENNA)
        px(img, 46+i, 22-i//2, ANTENNA)
    # 雷霆光環（橙紅）
    RING = (255, 100, 30, 120)
    for angle_deg in range(0, 360, 6):
        angle = math.radians(angle_deg)
        rx, ry = int(32 + 28*math.cos(angle)), int(32 + 18*math.sin(angle))
        if 0 <= rx < SIZE and 0 <= ry < SIZE:
            img.putpixel((rx, ry), RING)
    # 閃電符文
    BOLT = (255, 220, 50, 255)
    for i in range(5):
        px(img, 28+i, 30, BOLT)
    for i in range(4):
        px(img, 32-i, 34, BOLT)
    for i in range(5):
        px(img, 28+i, 38, BOLT)
    # 眼睛
    fill_circle(img, 42, 26, 4, (255, 255, 255, 255))
    fill_circle(img, 43, 26, 2, (20, 20, 20, 255))
    px(img, 42, 25, (255, 255, 255, 255))
    draw_outline(img)
    return img

def count_pixels(img):
    total = img.size[0] * img.size[1]
    non_transparent = sum(1 for x in range(img.size[0]) for y in range(img.size[1]) if img.getpixel((x, y))[3] > 0)
    return non_transparent, total, non_transparent / total * 100

if __name__ == "__main__":
    os.makedirs(OUT_DIR, exist_ok=True)
    
    generators = [
        ("T106_chain_lightning", gen_t106),
        ("T107_crab_torpedo", gen_t107),
        ("T108_vortex_anemone", gen_t108),
        ("T109_golden_dragon", gen_t109),
        ("T110_thunder_lobster", gen_t110),
    ]
    
    for name, gen_func in generators:
        img = gen_func()
        # 放大到 64x64（已是 64x64，確認）
        out_path = os.path.join(OUT_DIR, f"{name}.png")
        img.save(out_path)
        non_t, total, pct = count_pixels(img)
        print(f"✅ {name}.png — 非透明像素 {non_t}/{total} ({pct:.1f}%)")
    
    print("\n全部完成！")
