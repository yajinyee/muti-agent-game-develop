#!/usr/bin/env python3
"""
upgrade_basic_targets_day334.py — DAY-334 基礎目標物視覺升級
target-pixel-agent 負責維護

升級目標：
- T001-T006：32x32 → 64x64，提升密度和視覺細節
- T103 隕石：密度 22.6% → 目標 35%+
- T105 金幣魚：密度 28.4% → 目標 35%+

品質標準：
- 尺寸：64x64（與 T101-T105 一致）
- 密度：> 30%（視覺上清晰可辨）
- 陰影：3色漸層（LIGHT/MID/DARK）
- 輪廓：深色輪廓（辨識度）
"""

from PIL import Image, ImageDraw
import os
import math

OUT_DIR = "client/chiikawa-pixel/assets/sprites/targets"
SIZE = 64

# ── 顏色定義 ─────────────────────────────────────────────────
def shade(base_rgb, factor):
    """生成陰影色"""
    r, g, b = base_rgb
    return (int(r * factor), int(g * factor), int(b * factor), 255)

def with_alpha(rgb, a=255):
    return (rgb[0], rgb[1], rgb[2], a)

OUTLINE = (30, 25, 20, 255)  # 深棕輪廓

# ── 繪圖工具 ─────────────────────────────────────────────────
def fill_circle(img, cx, cy, r, color):
    """填充實心圓"""
    pixels = img.load()
    w, h = img.size
    for y in range(max(0, cy - r), min(h, cy + r + 1)):
        for x in range(max(0, cx - r), min(w, cx + r + 1)):
            if (x - cx) ** 2 + (y - cy) ** 2 <= r * r:
                pixels[x, y] = color

def fill_circle_shaded(img, cx, cy, r, light, mid, dark):
    """填充帶陰影的圓（左上亮，右下暗）"""
    pixels = img.load()
    w, h = img.size
    for y in range(max(0, cy - r), min(h, cy + r + 1)):
        for x in range(max(0, cx - r), min(w, cx + r + 1)):
            dx, dy = x - cx, y - cy
            if dx * dx + dy * dy <= r * r:
                # 左上亮，右下暗
                if dx < -r * 0.2 and dy < -r * 0.2:
                    pixels[x, y] = light
                elif dx > r * 0.2 and dy > r * 0.2:
                    pixels[x, y] = dark
                else:
                    pixels[x, y] = mid

def draw_outline_circle(img, cx, cy, r, color, thickness=1):
    """畫圓形輪廓"""
    pixels = img.load()
    w, h = img.size
    for t in range(thickness):
        rr = r - t
        for angle in range(360):
            rad = math.radians(angle)
            x = int(cx + rr * math.cos(rad))
            y = int(cy + rr * math.sin(rad))
            if 0 <= x < w and 0 <= y < h:
                pixels[x, y] = color

def fill_rect_shaded(img, x1, y1, x2, y2, light, mid, dark):
    """填充帶陰影的矩形"""
    pixels = img.load()
    w, h = img.size
    cx = (x1 + x2) / 2
    cy = (y1 + y2) / 2
    for y in range(max(0, y1), min(h, y2 + 1)):
        for x in range(max(0, x1), min(w, x2 + 1)):
            dx = x - cx
            dy = y - cy
            if dx < -(x2 - x1) * 0.2 and dy < -(y2 - y1) * 0.2:
                pixels[x, y] = light
            elif dx > (x2 - x1) * 0.2 and dy > (y2 - y1) * 0.2:
                pixels[x, y] = dark
            else:
                pixels[x, y] = mid

def draw_pixel(img, x, y, color):
    """安全繪製單個像素"""
    w, h = img.size
    if 0 <= x < w and 0 <= y < h:
        img.load()[x, y] = color

def draw_line(img, x1, y1, x2, y2, color, thickness=1):
    """畫線段"""
    dx = abs(x2 - x1)
    dy = abs(y2 - y1)
    steps = max(dx, dy, 1)
    for i in range(steps + 1):
        t = i / steps
        x = int(x1 + t * (x2 - x1))
        y = int(y1 + t * (y2 - y1))
        for tx in range(-thickness // 2, thickness // 2 + 1):
            for ty in range(-thickness // 2, thickness // 2 + 1):
                draw_pixel(img, x + tx, y + ty, color)

# ── T001 草（升級版）────────────────────────────────────────
def gen_T001_grass():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    pixels = img.load()

    # 草的顏色
    GRASS_LIGHT = (120, 220, 60, 255)
    GRASS_MID   = (80, 180, 40, 255)
    GRASS_DARK  = (50, 130, 25, 255)
    STEM_COLOR  = (60, 140, 30, 255)
    SOIL_LIGHT  = (180, 140, 80, 255)
    SOIL_MID    = (150, 110, 60, 255)
    SOIL_DARK   = (110, 80, 40, 255)
    FLOWER_Y    = (255, 230, 50, 255)   # 黃色小花
    FLOWER_W    = (255, 255, 255, 255)  # 白色花瓣

    # 土壤底部（橢圓，更大）
    for y in range(46, 62):
        for x in range(8, 56):
            dx = (x - 32) / 24.0
            dy = (y - 54) / 8.0
            if dx * dx + dy * dy <= 1.0:
                if dx < -0.3 and dy < -0.3:
                    pixels[x, y] = SOIL_LIGHT
                elif dx > 0.3 and dy > 0.3:
                    pixels[x, y] = SOIL_DARK
                else:
                    pixels[x, y] = SOIL_MID

    # 草莖（5根，更密集）
    stems = [
        (16, 44, 10, 10),   # 最左
        (22, 44, 18, 6),    # 左
        (32, 42, 32, 4),    # 中
        (42, 44, 46, 6),    # 右
        (48, 44, 54, 10),   # 最右
    ]
    for x1, y1, x2, y2 in stems:
        draw_line(img, x1, y1, x2, y2, STEM_COLOR, 3)

    # 草葉（每根莖上的葉片，更寬更密）
    leaf_configs = [
        (10, 10, [(-10, 6), (6, 10)]),
        (18, 6,  [(-8, 8), (8, 8)]),
        (32, 4,  [(-10, 8), (10, 8)]),
        (46, 6,  [(-8, 8), (8, 8)]),
        (54, 10, [(-6, 10), (10, 6)]),
    ]
    for sx, sy, leaves in leaf_configs:
        for (lx, ly) in leaves:
            for t in range(14):
                tt = t / 13.0
                x = int(sx + lx * tt)
                y = int(sy + ly * tt)
                w_half = max(1, int(4 * (1 - tt)))
                for dx in range(-w_half, w_half + 1):
                    c = GRASS_LIGHT if tt < 0.3 else (GRASS_MID if tt < 0.7 else GRASS_DARK)
                    draw_pixel(img, x + dx, y, c)

    # 小花（增加視覺趣味）
    flowers = [(32, 4), (18, 6), (46, 6)]
    for fx, fy in flowers:
        # 花瓣
        for angle in range(0, 360, 60):
            rad = math.radians(angle)
            px = int(fx + 4 * math.cos(rad))
            py = int(fy + 4 * math.sin(rad))
            fill_circle(img, px, py, 2, FLOWER_W)
        # 花心
        fill_circle(img, fx, fy, 2, FLOWER_Y)

    # 輪廓
    result = _add_outline(img)
    return result

# ── T002 綠蟲（升級版）──────────────────────────────────────
def gen_T002_bug_g():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))

    BUG_LIGHT = (100, 220, 80, 255)
    BUG_MID   = (60, 180, 40, 255)
    BUG_DARK  = (30, 130, 20, 255)
    EYE_WHITE = (255, 255, 255, 255)
    EYE_PUPIL = (20, 20, 20, 255)
    EYE_SHINE = (255, 255, 255, 255)
    LEG_COLOR = (40, 140, 20, 255)
    ANTENNA   = (50, 160, 30, 255)

    # 身體（橢圓，橫向）
    fill_circle_shaded(img, 32, 36, 18, BUG_LIGHT, BUG_MID, BUG_DARK)
    # 頭部（圓形）
    fill_circle_shaded(img, 48, 30, 12, BUG_LIGHT, BUG_MID, BUG_DARK)

    # 觸角
    draw_line(img, 48, 18, 44, 8, ANTENNA, 2)
    draw_line(img, 48, 18, 54, 8, ANTENNA, 2)
    fill_circle(img, 44, 7, 3, BUG_MID)
    fill_circle(img, 54, 7, 3, BUG_MID)

    # 眼睛
    fill_circle(img, 52, 26, 5, EYE_WHITE)
    fill_circle(img, 53, 27, 3, EYE_PUPIL)
    draw_pixel(img, 52, 25, EYE_SHINE)

    # 腿（3對）
    legs = [(20, 44, 12, 52), (28, 46, 22, 56), (36, 46, 30, 56),
            (20, 44, 12, 36), (28, 46, 22, 38), (36, 46, 30, 38)]
    for x1, y1, x2, y2 in legs:
        draw_line(img, x1, y1, x2, y2, LEG_COLOR, 2)

    # 輪廓
    result = _add_outline(img)
    return result

# ── T003 紅蟲（升級版）──────────────────────────────────────
def gen_T003_bug_r():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))

    BUG_LIGHT = (255, 120, 80, 255)
    BUG_MID   = (220, 60, 40, 255)
    BUG_DARK  = (160, 30, 20, 255)
    EYE_WHITE = (255, 255, 255, 255)
    EYE_PUPIL = (20, 20, 20, 255)
    EYE_SHINE = (255, 255, 255, 255)
    LEG_COLOR = (180, 40, 20, 255)
    ANTENNA   = (200, 50, 30, 255)
    SPOT      = (255, 200, 180, 255)  # 斑點

    # 身體
    fill_circle_shaded(img, 32, 36, 18, BUG_LIGHT, BUG_MID, BUG_DARK)
    # 頭部
    fill_circle_shaded(img, 48, 30, 12, BUG_LIGHT, BUG_MID, BUG_DARK)

    # 斑點（紅蟲特色）
    fill_circle(img, 26, 32, 4, SPOT)
    fill_circle(img, 36, 40, 3, SPOT)

    # 觸角
    draw_line(img, 48, 18, 44, 8, ANTENNA, 2)
    draw_line(img, 48, 18, 54, 8, ANTENNA, 2)
    fill_circle(img, 44, 7, 3, BUG_MID)
    fill_circle(img, 54, 7, 3, BUG_MID)

    # 眼睛
    fill_circle(img, 52, 26, 5, EYE_WHITE)
    fill_circle(img, 53, 27, 3, EYE_PUPIL)
    draw_pixel(img, 52, 25, EYE_SHINE)

    # 腿
    legs = [(20, 44, 12, 52), (28, 46, 22, 56), (36, 46, 30, 56),
            (20, 44, 12, 36), (28, 46, 22, 38), (36, 46, 30, 38)]
    for x1, y1, x2, y2 in legs:
        draw_line(img, x1, y1, x2, y2, LEG_COLOR, 2)

    result = _add_outline(img)
    return result

# ── T004 藍蟲（升級版）──────────────────────────────────────
def gen_T004_bug_b():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))

    BUG_LIGHT = (100, 160, 255, 255)
    BUG_MID   = (50, 100, 220, 255)
    BUG_DARK  = (20, 60, 160, 255)
    EYE_WHITE = (255, 255, 255, 255)
    EYE_PUPIL = (20, 20, 20, 255)
    EYE_SHINE = (255, 255, 255, 255)
    LEG_COLOR = (30, 80, 180, 255)
    ANTENNA   = (40, 90, 200, 255)
    WING      = (150, 200, 255, 180)  # 半透明翅膀

    # 翅膀（橢圓，半透明）
    for y in range(10, 40):
        for x in range(8, 30):
            dx = (x - 19) / 11.0
            dy = (y - 25) / 15.0
            if dx * dx + dy * dy <= 1.0:
                img.load()[x, y] = WING

    # 身體
    fill_circle_shaded(img, 32, 36, 18, BUG_LIGHT, BUG_MID, BUG_DARK)
    # 頭部
    fill_circle_shaded(img, 48, 30, 12, BUG_LIGHT, BUG_MID, BUG_DARK)

    # 觸角
    draw_line(img, 48, 18, 44, 8, ANTENNA, 2)
    draw_line(img, 48, 18, 54, 8, ANTENNA, 2)
    fill_circle(img, 44, 7, 3, BUG_MID)
    fill_circle(img, 54, 7, 3, BUG_MID)

    # 眼睛
    fill_circle(img, 52, 26, 5, EYE_WHITE)
    fill_circle(img, 53, 27, 3, EYE_PUPIL)
    draw_pixel(img, 52, 25, EYE_SHINE)

    # 腿
    legs = [(20, 44, 12, 52), (28, 46, 22, 56), (36, 46, 30, 56),
            (20, 44, 12, 36), (28, 46, 22, 38), (36, 46, 30, 38)]
    for x1, y1, x2, y2 in legs:
        draw_line(img, x1, y1, x2, y2, LEG_COLOR, 2)

    result = _add_outline(img)
    return result

# ── T005 布丁（升級版）──────────────────────────────────────
def gen_T005_pudding():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))

    PUDI_LIGHT = (255, 220, 100, 255)
    PUDI_MID   = (240, 180, 60, 255)
    PUDI_DARK  = (200, 140, 30, 255)
    CARAMEL_L  = (200, 120, 40, 255)
    CARAMEL_M  = (170, 90, 20, 255)
    CARAMEL_D  = (130, 60, 10, 255)
    SHINE      = (255, 255, 200, 200)
    EYE_WHITE  = (255, 255, 255, 255)
    EYE_PUPIL  = (60, 40, 20, 255)
    MOUTH      = (180, 80, 40, 255)
    BLUSH      = (255, 180, 160, 180)

    # 焦糖頂部（橢圓）
    for y in range(8, 22):
        for x in range(16, 48):
            dx = (x - 32) / 16.0
            dy = (y - 15) / 7.0
            if dx * dx + dy * dy <= 1.0:
                if dx < -0.3 and dy < -0.3:
                    img.load()[x, y] = CARAMEL_L
                elif dx > 0.3 and dy > 0.3:
                    img.load()[x, y] = CARAMEL_D
                else:
                    img.load()[x, y] = CARAMEL_M

    # 布丁主體（梯形）
    for y in range(18, 56):
        t = (y - 18) / 38.0
        x_left = int(14 + t * 4)
        x_right = int(50 - t * 4)
        for x in range(x_left, x_right):
            dx = x - 32
            dy = y - 37
            if dx < -8 and dy < -8:
                img.load()[x, y] = PUDI_LIGHT
            elif dx > 8 and dy > 8:
                img.load()[x, y] = PUDI_DARK
            else:
                img.load()[x, y] = PUDI_MID

    # 高光
    fill_circle(img, 24, 28, 5, SHINE)

    # 眼睛
    fill_circle(img, 26, 36, 4, EYE_WHITE)
    fill_circle(img, 27, 37, 2, EYE_PUPIL)
    fill_circle(img, 38, 36, 4, EYE_WHITE)
    fill_circle(img, 39, 37, 2, EYE_PUPIL)

    # 腮紅
    fill_circle(img, 22, 42, 4, BLUSH)
    fill_circle(img, 42, 42, 4, BLUSH)

    # 嘴巴（V形）
    draw_pixel(img, 30, 44, MOUTH)
    draw_pixel(img, 31, 45, MOUTH)
    draw_pixel(img, 32, 45, MOUTH)
    draw_pixel(img, 33, 45, MOUTH)
    draw_pixel(img, 34, 44, MOUTH)

    result = _add_outline(img)
    return result

# ── T006 蘑菇（升級版）──────────────────────────────────────
def gen_T006_mushroom():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))

    CAP_LIGHT  = (255, 80, 60, 255)
    CAP_MID    = (220, 40, 30, 255)
    CAP_DARK   = (160, 20, 15, 255)
    SPOT_COLOR = (255, 240, 220, 255)
    STEM_LIGHT = (255, 240, 200, 255)
    STEM_MID   = (240, 210, 160, 255)
    STEM_DARK  = (200, 170, 120, 255)
    EYE_WHITE  = (255, 255, 255, 255)
    EYE_PUPIL  = (40, 30, 20, 255)
    MOUTH      = (180, 60, 40, 255)
    BLUSH      = (255, 160, 140, 180)

    # 菌柄（矩形）
    fill_rect_shaded(img, 20, 40, 44, 58, STEM_LIGHT, STEM_MID, STEM_DARK)

    # 菌傘（半圓）
    for y in range(8, 44):
        for x in range(6, 58):
            dx = (x - 32) / 26.0
            dy = (y - 36) / 28.0
            if dx * dx + dy * dy <= 1.0 and y <= 36:
                if dx < -0.3 and dy < -0.3:
                    img.load()[x, y] = CAP_LIGHT
                elif dx > 0.3 and dy > 0.3:
                    img.load()[x, y] = CAP_DARK
                else:
                    img.load()[x, y] = CAP_MID

    # 白色斑點
    spots = [(22, 16, 5), (42, 14, 4), (32, 22, 6), (16, 28, 4), (46, 26, 4)]
    for sx, sy, sr in spots:
        fill_circle(img, sx, sy, sr, SPOT_COLOR)

    # 眼睛
    fill_circle(img, 26, 36, 4, EYE_WHITE)
    fill_circle(img, 27, 37, 2, EYE_PUPIL)
    fill_circle(img, 38, 36, 4, EYE_WHITE)
    fill_circle(img, 39, 37, 2, EYE_PUPIL)

    # 腮紅
    fill_circle(img, 22, 42, 4, BLUSH)
    fill_circle(img, 42, 42, 4, BLUSH)

    # 嘴巴
    draw_pixel(img, 30, 44, MOUTH)
    draw_pixel(img, 31, 45, MOUTH)
    draw_pixel(img, 32, 45, MOUTH)
    draw_pixel(img, 33, 45, MOUTH)
    draw_pixel(img, 34, 44, MOUTH)

    result = _add_outline(img)
    return result

# ── 輔助：加輪廓 ─────────────────────────────────────────────
def _add_outline(img):
    """在非透明像素邊緣加深色輪廓"""
    pixels = img.load()
    result = img.copy()
    rp = result.load()
    w, h = img.size
    for y in range(1, h - 1):
        for x in range(1, w - 1):
            if pixels[x, y][3] > 10:
                for dx, dy in [(-1, 0), (1, 0), (0, -1), (0, 1)]:
                    if pixels[x + dx, y + dy][3] < 10:
                        rp[x, y] = OUTLINE
                        break
    return result

# ── 主程式 ────────────────────────────────────────────────────
def main():
    os.makedirs(OUT_DIR, exist_ok=True)

    generators = [
        ("T001_grass",    gen_T001_grass,    "草"),
        ("T002_bug_g",    gen_T002_bug_g,    "綠蟲"),
        ("T003_bug_r",    gen_T003_bug_r,    "紅蟲"),
        ("T004_bug_b",    gen_T004_bug_b,    "藍蟲"),
        ("T005_pudding",  gen_T005_pudding,  "布丁"),
        ("T006_mushroom", gen_T006_mushroom, "蘑菇"),
    ]

    print("DAY-334 基礎目標物視覺升級")
    print("=" * 60)

    for fname, gen_func, name in generators:
        img = gen_func()
        path = os.path.join(OUT_DIR, fname + ".png")
        img.save(path)

        # 驗證
        w, h = img.size
        pixels = img.load()
        non_transparent = sum(1 for y in range(h) for x in range(w) if pixels[x, y][3] > 10)
        density = non_transparent / (w * h) * 100
        status = "✅" if density > 30 else "⚠️"
        print(f"{status} {fname} ({name}): {w}x{h}, 密度={density:.1f}%")

    print("\n升級完成！")
    print("注意：需要在 Godot 中重新匯入這些 PNG 才能看到效果")

if __name__ == "__main__":
    main()
