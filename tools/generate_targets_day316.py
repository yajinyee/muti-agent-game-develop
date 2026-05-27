"""
generate_targets_day316.py — DAY-316 T186-T190 精靈圖生成
T186 幸運鏡像宇宙魚（2150x）：深藍橢圓魚身 + 鏡像反射光環 + 雙重輪廓
T187 幸運引力場魚（2200x）：深紫橢圓魚身 + 引力漩渦線 + 吸積盤
T188 幸運時間加速魚（2300x）：火橙橢圓魚身 + 閃電紋路 + 速度線
T189 幸運星雲漩渦魚（2400x）：深紫大型魚身 + 星雲漩渦 + 20 道光芒
T190 幸運宇宙審判魚（2500x）：深紅大型魚身 + 天秤符號 + 24 道光芒 + 四層光環
"""
import os
import math
from PIL import Image, ImageDraw

OUT_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
SIZE = 64

def px(img, x, y, color):
    if 0 <= x < SIZE and 0 <= y < SIZE:
        img.putpixel((x, y), color)

def fill_ellipse_shaded(img, cx, cy, rx, ry, base_color, light_color, dark_color):
    for y in range(max(0, cy - ry), min(SIZE, cy + ry + 1)):
        for x in range(max(0, cx - rx), min(SIZE, cx + rx + 1)):
            dx = (x - cx) / rx
            dy = (y - cy) / ry
            if dx*dx + dy*dy <= 1.0:
                # 光源左上
                if dx < -0.2 and dy < -0.2:
                    c = light_color
                elif dx > 0.3 and dy > 0.3:
                    c = dark_color
                else:
                    c = base_color
                img.putpixel((x, y), c)

def draw_outline(img, cx, cy, rx, ry, color, thickness=1):
    for t in range(thickness):
        for angle in range(360):
            rad = math.radians(angle)
            x = int(cx + (rx + t) * math.cos(rad))
            y = int(cy + (ry + t) * math.sin(rad))
            if 0 <= x < SIZE and 0 <= y < SIZE:
                img.putpixel((x, y), color)

def draw_rays(img, cx, cy, n_rays, length, color, alpha_start=200):
    for i in range(n_rays):
        angle = math.radians(i * 360 / n_rays)
        for r in range(8, length):
            x = int(cx + r * math.cos(angle))
            y = int(cy + r * math.sin(angle))
            if 0 <= x < SIZE and 0 <= y < SIZE:
                alpha = int(alpha_start * (1 - (r - 8) / (length - 8)))
                c = (*color[:3], alpha)
                img.putpixel((x, y), c)

def draw_ring(img, cx, cy, radius, color, width=2):
    for angle in range(360):
        rad = math.radians(angle)
        for w in range(width):
            x = int(cx + (radius + w) * math.cos(rad))
            y = int(cy + (radius + w) * math.sin(rad))
            if 0 <= x < SIZE and 0 <= y < SIZE:
                img.putpixel((x, y), color)

def make_T186():
    """T186 幸運鏡像宇宙魚：深藍橢圓魚身 + 鏡像反射光環 + 雙重輪廓"""
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    cx, cy = 32, 32
    # 外層光暈（低透明度）
    for r in range(28, 32):
        draw_ring(img, cx, cy, r, (100, 150, 255, 40))
    # 魚身（深藍橢圓）
    fill_ellipse_shaded(img, cx, cy, 22, 16,
        (30, 60, 180),   # 基底深藍
        (80, 130, 255),  # 高光藍
        (10, 30, 120))   # 暗藍
    # 鏡像反射線（水平中線）
    for x in range(cx - 20, cx + 21):
        if 0 <= x < SIZE:
            img.putpixel((x, cy), (200, 220, 255, 180))
    # 雙重輪廓
    draw_outline(img, cx, cy, 22, 16, (150, 200, 255, 220), 1)
    draw_outline(img, cx, cy, 24, 18, (80, 120, 220, 120), 1)
    # 鏡像符號（中心）
    for dx, dy in [(-2,0),(-1,0),(0,0),(1,0),(2,0),(0,-2),(0,-1),(0,1),(0,2)]:
        px(img, cx+dx, cy+dy, (255, 255, 255, 230))
    # 眼睛
    px(img, cx-8, cy-4, (255, 255, 255, 255))
    px(img, cx-7, cy-4, (255, 255, 255, 255))
    px(img, cx-8, cy-3, (100, 150, 255, 255))
    # 魚尾
    for i in range(8):
        px(img, cx+22+i, cy-4+i//2, (30, 60, 180, 200-i*20))
        px(img, cx+22+i, cy+4-i//2, (30, 60, 180, 200-i*20))
    return img

def make_T187():
    """T187 幸運引力場魚：深紫橢圓魚身 + 引力漩渦線 + 吸積盤"""
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    cx, cy = 32, 32
    # 吸積盤（橢圓光環）
    for r in range(24, 30):
        for angle in range(0, 360, 2):
            rad = math.radians(angle)
            x = int(cx + r * 1.4 * math.cos(rad))
            y = int(cy + r * 0.5 * math.sin(rad))
            if 0 <= x < SIZE and 0 <= y < SIZE:
                alpha = 80 + int(60 * math.sin(math.radians(angle * 3)))
                img.putpixel((x, y), (180, 100, 255, alpha))
    # 魚身（深紫橢圓）
    fill_ellipse_shaded(img, cx, cy, 20, 15,
        (80, 20, 160),   # 基底深紫
        (150, 80, 255),  # 高光紫
        (40, 10, 100))   # 暗紫
    # 引力漩渦線（螺旋）
    for i in range(3):
        for t in range(60):
            angle = math.radians(t * 6 + i * 120)
            r = 5 + t * 0.2
            x = int(cx + r * math.cos(angle))
            y = int(cy + r * math.sin(angle))
            if 0 <= x < SIZE and 0 <= y < SIZE:
                alpha = max(0, 180 - t * 2)
                img.putpixel((x, y), (200, 150, 255, alpha))
    # 輪廓
    draw_outline(img, cx, cy, 20, 15, (200, 150, 255, 220), 1)
    # 眼睛
    px(img, cx-7, cy-4, (255, 255, 255, 255))
    px(img, cx-6, cy-4, (255, 255, 255, 255))
    px(img, cx-7, cy-3, (150, 80, 255, 255))
    # 魚尾
    for i in range(8):
        px(img, cx+20+i, cy-4+i//2, (80, 20, 160, 200-i*20))
        px(img, cx+20+i, cy+4-i//2, (80, 20, 160, 200-i*20))
    return img

def make_T188():
    """T188 幸運時間加速魚：火橙橢圓魚身 + 閃電紋路 + 速度線"""
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    cx, cy = 32, 32
    # 速度線（水平）
    for i in range(5):
        y_off = -8 + i * 4
        for x in range(4, cx - 18):
            alpha = int(150 * (x / (cx - 18)))
            img.putpixel((x, cy + y_off), (255, 200, 50, alpha))
    # 魚身（火橙橢圓）
    fill_ellipse_shaded(img, cx, cy, 21, 15,
        (220, 100, 20),  # 基底火橙
        (255, 180, 60),  # 高光橙
        (150, 60, 10))   # 暗橙
    # 閃電紋路（Z字）
    lightning = [(cx-8, cy-6), (cx-2, cy-2), (cx-6, cy+2), (cx+2, cy+6)]
    for i in range(len(lightning)-1):
        x1, y1 = lightning[i]
        x2, y2 = lightning[i+1]
        steps = max(abs(x2-x1), abs(y2-y1))
        for s in range(steps+1):
            lx = int(x1 + (x2-x1)*s/max(steps,1))
            ly = int(y1 + (y2-y1)*s/max(steps,1))
            px(img, lx, ly, (255, 255, 100, 220))
            px(img, lx+1, ly, (255, 255, 100, 150))
    # 輪廓
    draw_outline(img, cx, cy, 21, 15, (255, 200, 80, 220), 1)
    # 眼睛
    px(img, cx-8, cy-4, (255, 255, 255, 255))
    px(img, cx-7, cy-4, (255, 255, 255, 255))
    px(img, cx-8, cy-3, (255, 150, 0, 255))
    # 魚尾（火焰）
    for i in range(10):
        px(img, cx+21+i, cy-5+i//2, (255, 100+i*10, 0, 200-i*18))
        px(img, cx+21+i, cy+5-i//2, (255, 100+i*10, 0, 200-i*18))
    return img

def make_T189():
    """T189 幸運星雲漩渦魚：深紫大型魚身 + 星雲漩渦 + 20 道光芒"""
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    cx, cy = 32, 32
    # 外層星雲光暈
    for r in range(26, 32):
        for angle in range(0, 360, 3):
            rad = math.radians(angle)
            x = int(cx + r * math.cos(rad))
            y = int(cy + r * math.sin(rad))
            if 0 <= x < SIZE and 0 <= y < SIZE:
                alpha = 30 + int(40 * abs(math.sin(math.radians(angle * 2))))
                img.putpixel((x, y), (150, 50, 220, alpha))
    # 20 道光芒
    draw_rays(img, cx, cy, 20, 28, (180, 80, 255), 160)
    # 魚身（深紫大型橢圓）
    fill_ellipse_shaded(img, cx, cy, 23, 17,
        (100, 30, 180),  # 基底深紫
        (180, 100, 255), # 高光紫
        (50, 10, 120))   # 暗紫
    # 漩渦紋路
    for i in range(4):
        for t in range(80):
            angle = math.radians(t * 4.5 + i * 90)
            r = 3 + t * 0.22
            x = int(cx + r * math.cos(angle))
            y = int(cy + r * math.sin(angle))
            if 0 <= x < SIZE and 0 <= y < SIZE:
                alpha = max(0, 200 - t * 2)
                img.putpixel((x, y), (220, 160, 255, alpha))
    # 輪廓
    draw_outline(img, cx, cy, 23, 17, (220, 150, 255, 230), 1)
    # 眼睛
    px(img, cx-9, cy-5, (255, 255, 255, 255))
    px(img, cx-8, cy-5, (255, 255, 255, 255))
    px(img, cx-9, cy-4, (180, 80, 255, 255))
    # 魚尾
    for i in range(9):
        px(img, cx+23+i, cy-5+i//2, (100, 30, 180, 200-i*20))
        px(img, cx+23+i, cy+5-i//2, (100, 30, 180, 200-i*20))
    return img

def make_T190():
    """T190 幸運宇宙審判魚：深紅大型魚身 + 天秤符號 + 24 道光芒 + 四層光環"""
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    cx, cy = 32, 32
    # 四層光環
    for r, alpha in [(30, 30), (27, 50), (24, 70), (21, 90)]:
        draw_ring(img, cx, cy, r, (255, 80, 80, alpha), 2)
    # 24 道光芒
    draw_rays(img, cx, cy, 24, 30, (255, 100, 100), 180)
    # 魚身（深紅大型橢圓）
    fill_ellipse_shaded(img, cx, cy, 24, 18,
        (180, 20, 20),   # 基底深紅
        (255, 80, 80),   # 高光紅
        (100, 10, 10))   # 暗紅
    # 天秤符號（中心）
    # 橫桿
    for x in range(cx-10, cx+11):
        px(img, x, cy-4, (255, 220, 100, 230))
    # 中心柱
    for y in range(cy-8, cy+5):
        px(img, cx, y, (255, 220, 100, 230))
    # 左盤
    for dx in range(-8, 1):
        for dy in range(-2, 3):
            if dx*dx + dy*dy <= 16:
                px(img, cx+dx-2, cy+4+dy, (255, 220, 100, 180))
    # 右盤
    for dx in range(0, 9):
        for dy in range(-2, 3):
            if dx*dx + dy*dy <= 16:
                px(img, cx+dx+2, cy+4+dy, (255, 220, 100, 180))
    # 輪廓
    draw_outline(img, cx, cy, 24, 18, (255, 120, 120, 230), 1)
    # 眼睛
    px(img, cx-10, cy-6, (255, 255, 255, 255))
    px(img, cx-9, cy-6, (255, 255, 255, 255))
    px(img, cx-10, cy-5, (255, 80, 80, 255))
    # 魚尾
    for i in range(10):
        px(img, cx+24+i, cy-6+i//2, (180, 20, 20, 200-i*18))
        px(img, cx+24+i, cy+6-i//2, (180, 20, 20, 200-i*18))
    return img

def count_pixels(img):
    total = 0
    for y in range(SIZE):
        for x in range(SIZE):
            if img.getpixel((x, y))[3] > 10:
                total += 1
    return total

def main():
    os.makedirs(OUT_DIR, exist_ok=True)
    targets = [
        ("T186_mirror_universe", make_T186, "幸運鏡像宇宙魚"),
        ("T187_gravity_field",   make_T187, "幸運引力場魚"),
        ("T188_time_acceleration", make_T188, "幸運時間加速魚"),
        ("T189_nebula_vortex",   make_T189, "幸運星雲漩渦魚"),
        ("T190_cosmic_judgment", make_T190, "幸運宇宙審判魚"),
    ]
    for file_id, make_fn, name in targets:
        img = make_fn()
        tid = file_id.split("_")[0]
        out_path = os.path.join(OUT_DIR, f"{tid}_{file_id.split('_', 1)[1]}.png")
        img.save(out_path)
        count = count_pixels(img)
        pct = count / (SIZE * SIZE) * 100
        print(f"[OK] {tid} {name}: {count} 非透明像素 ({pct:.1f}%) → {out_path}")

if __name__ == "__main__":
    main()
