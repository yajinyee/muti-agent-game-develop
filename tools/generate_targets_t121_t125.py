"""
generate_targets_t121_t125.py — 生成 T121-T125 目標物精靈圖
DAY-296 新增：鏡像魚、黃金雨魚、冰凍炸彈魚、雷暴魚、大轉盤魚
"""
import os
from PIL import Image, ImageDraw

SIZE = 64
OUT_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
os.makedirs(OUT_DIR, exist_ok=True)

def new_img():
    return Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))

def px(img, x, y, c):
    if 0 <= x < SIZE and 0 <= y < SIZE:
        img.putpixel((x, y), c)

def fill_circle(img, cx, cy, r, color, shade=True):
    for y in range(max(0, cy-r), min(SIZE, cy+r+1)):
        for x in range(max(0, cx-r), min(SIZE, cx+r+1)):
            dx, dy = x-cx, y-cy
            if dx*dx + dy*dy <= r*r:
                if shade:
                    # 3色陰影
                    dist = (dx*dx + dy*dy) ** 0.5
                    light = 1.0 - (dist / r) * 0.35
                    dot = -(dx * (-0.7) + dy * (-0.7)) / (r + 0.001)
                    lf = 0.5 + 0.5 * max(0, dot)
                    factor = light * (0.7 + 0.6 * lf)
                    r2 = min(255, int(color[0] * factor))
                    g2 = min(255, int(color[1] * factor))
                    b2 = min(255, int(color[2] * factor))
                    img.putpixel((x, y), (r2, g2, b2, 255))
                else:
                    img.putpixel((x, y), color)

def fill_ellipse(img, cx, cy, rx, ry, color):
    for y in range(max(0, cy-ry), min(SIZE, cy+ry+1)):
        for x in range(max(0, cx-rx), min(SIZE, cx+rx+1)):
            dx, dy = (x-cx)/rx, (y-cy)/ry
            if dx*dx + dy*dy <= 1.0:
                img.putpixel((x, y), color)

def draw_outline(img, color=(20, 20, 20, 255)):
    """8方向輪廓"""
    pixels = img.load()
    outline = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    op = outline.load()
    for y in range(SIZE):
        for x in range(SIZE):
            if pixels[x, y][3] > 128:
                for dx, dy in [(-1,0),(1,0),(0,-1),(0,1),(-1,-1),(1,-1),(-1,1),(1,1)]:
                    nx, ny = x+dx, y+dy
                    if 0 <= nx < SIZE and 0 <= ny < SIZE and pixels[nx, ny][3] < 128:
                        op[nx, ny] = color
    img.paste(outline, mask=outline)

# ── T121 幸運鏡像魚 ──────────────────────────────────────────
def gen_t121():
    img = new_img()
    # 主體：淡紫色橢圓魚身
    fill_ellipse(img, 32, 32, 22, 14, (200, 150, 255, 255))
    fill_ellipse(img, 32, 32, 18, 10, (220, 180, 255, 255))
    # 鏡像效果：左右對稱光線
    for i in range(8):
        x = 10 + i * 6
        for dy in range(-2, 3):
            px(img, x, 32+dy, (255, 220, 255, 200))
            px(img, SIZE-1-x, 32+dy, (255, 220, 255, 200))
    # 中心鏡面（菱形）
    for d in range(8):
        px(img, 32, 32-d, (255, 255, 255, 220))
        px(img, 32, 32+d, (255, 255, 255, 220))
        px(img, 32-d, 32, (255, 255, 255, 220))
        px(img, 32+d, 32, (255, 255, 255, 220))
    # 眼睛
    fill_circle(img, 38, 28, 3, (255, 255, 255, 255), shade=False)
    px(img, 38, 28, (80, 0, 120, 255))
    px(img, 39, 27, (255, 255, 255, 200))
    # 魚尾（左側）
    fill_ellipse(img, 10, 32, 6, 10, (180, 120, 240, 255))
    # 光環
    draw = ImageDraw.Draw(img)
    draw.ellipse([4, 4, 60, 60], outline=(200, 150, 255, 120), width=2)
    draw_outline(img, (80, 0, 120, 255))
    img.save(os.path.join(OUT_DIR, "T121_mirror_fish.png"))
    print("T121 saved")

# ── T122 幸運黃金雨魚 ────────────────────────────────────────
def gen_t122():
    img = new_img()
    # 主體：金色橢圓魚身
    fill_ellipse(img, 32, 32, 22, 14, (255, 200, 0, 255))
    fill_ellipse(img, 32, 32, 18, 10, (255, 220, 50, 255))
    # 金幣雨效果：上方散落金幣
    coin_positions = [(20, 10), (32, 6), (44, 10), (15, 18), (49, 18)]
    for cx, cy in coin_positions:
        fill_circle(img, cx, cy, 4, (255, 215, 0, 255), shade=False)
        px(img, cx, cy, (200, 150, 0, 255))
        px(img, cx-1, cy-1, (255, 255, 200, 200))
    # 雨滴線條
    for i in range(5):
        x = 12 + i * 10
        for dy in range(4):
            px(img, x, 22+dy, (255, 215, 0, 180))
    # 眼睛
    fill_circle(img, 38, 28, 3, (255, 255, 255, 255), shade=False)
    px(img, 38, 28, (150, 100, 0, 255))
    px(img, 39, 27, (255, 255, 255, 200))
    # 魚尾
    fill_ellipse(img, 10, 32, 6, 10, (200, 150, 0, 255))
    # 光環
    draw = ImageDraw.Draw(img)
    draw.ellipse([4, 4, 60, 60], outline=(255, 215, 0, 150), width=2)
    draw_outline(img, (150, 100, 0, 255))
    img.save(os.path.join(OUT_DIR, "T122_golden_rain.png"))
    print("T122 saved")

# ── T123 幸運冰凍炸彈魚 ──────────────────────────────────────
def gen_t123():
    img = new_img()
    # 主體：冰藍色圓形魚身
    fill_circle(img, 32, 32, 20, (0, 200, 230, 255))
    fill_circle(img, 32, 32, 16, (100, 230, 255, 255))
    # 冰晶紋路（六角形）
    for angle_step in range(6):
        import math
        angle = angle_step * math.pi / 3
        ex = int(32 + 12 * math.cos(angle))
        ey = int(32 + 12 * math.sin(angle))
        for t in range(10):
            ix = int(32 + (ex-32) * t / 10)
            iy = int(32 + (ey-32) * t / 10)
            px(img, ix, iy, (200, 240, 255, 200))
    # 炸彈導火線（右上）
    for i in range(6):
        px(img, 48+i, 16-i, (100, 80, 60, 255))
    px(img, 54, 10, (255, 200, 0, 255))  # 火花
    px(img, 55, 9, (255, 150, 0, 255))
    # 眼睛
    fill_circle(img, 36, 28, 3, (255, 255, 255, 255), shade=False)
    px(img, 36, 28, (0, 100, 150, 255))
    px(img, 37, 27, (255, 255, 255, 200))
    # 冰凍光環
    draw = ImageDraw.Draw(img)
    draw.ellipse([4, 4, 60, 60], outline=(0, 200, 255, 150), width=2)
    draw_outline(img, (0, 100, 150, 255))
    img.save(os.path.join(OUT_DIR, "T123_freeze_bomb.png"))
    print("T123 saved")

# ── T124 幸運雷暴魚 ──────────────────────────────────────────
def gen_t124():
    img = new_img()
    # 主體：深黃色橢圓魚身
    fill_ellipse(img, 32, 32, 22, 14, (200, 160, 0, 255))
    fill_ellipse(img, 32, 32, 18, 10, (255, 200, 0, 255))
    # 閃電紋路（Z字形）
    lightning_pts = [(28, 14), (36, 14), (30, 28), (38, 28), (28, 44), (36, 44)]
    for i in range(len(lightning_pts)-1):
        x1, y1 = lightning_pts[i]
        x2, y2 = lightning_pts[i+1]
        steps = max(abs(x2-x1), abs(y2-y1))
        for t in range(steps+1):
            ix = int(x1 + (x2-x1) * t / max(steps, 1))
            iy = int(y1 + (y2-y1) * t / max(steps, 1))
            px(img, ix, iy, (255, 255, 100, 255))
    # 雷暴雲（上方）
    fill_circle(img, 24, 12, 6, (80, 80, 100, 200), shade=False)
    fill_circle(img, 32, 10, 7, (100, 100, 120, 200), shade=False)
    fill_circle(img, 40, 12, 6, (80, 80, 100, 200), shade=False)
    # 眼睛
    fill_circle(img, 38, 28, 3, (255, 255, 255, 255), shade=False)
    px(img, 38, 28, (100, 80, 0, 255))
    px(img, 39, 27, (255, 255, 255, 200))
    # 魚尾
    fill_ellipse(img, 10, 32, 6, 10, (180, 140, 0, 255))
    draw_outline(img, (100, 80, 0, 255))
    img.save(os.path.join(OUT_DIR, "T124_thunder_storm.png"))
    print("T124 saved")

# ── T125 幸運大轉盤魚 ────────────────────────────────────────
def gen_t125():
    img = new_img()
    # 主體：粉紅色圓形魚身
    fill_circle(img, 32, 32, 20, (255, 100, 180, 255))
    fill_circle(img, 32, 32, 16, (255, 150, 200, 255))
    # 轉盤格子（8格，交替顏色）
    import math
    colors_wheel = [
        (255, 50, 50, 200),   # 紅
        (255, 150, 0, 200),   # 橙
        (255, 220, 0, 200),   # 黃
        (50, 200, 50, 200),   # 綠
        (0, 150, 255, 200),   # 藍
        (100, 0, 200, 200),   # 紫
        (255, 100, 200, 200), # 粉
        (200, 200, 200, 200), # 白
    ]
    for i in range(8):
        angle_start = i * math.pi / 4
        angle_end = (i + 1) * math.pi / 4
        for r in range(6, 16):
            for a_step in range(20):
                a = angle_start + (angle_end - angle_start) * a_step / 20
                x = int(32 + r * math.cos(a))
                y = int(32 + r * math.sin(a))
                if 0 <= x < SIZE and 0 <= y < SIZE:
                    px(img, x, y, colors_wheel[i])
    # 中心圓
    fill_circle(img, 32, 32, 5, (255, 215, 0, 255), shade=False)
    px(img, 32, 32, (200, 150, 0, 255))
    # 指針（上方）
    for i in range(5):
        px(img, 32, 14+i, (255, 50, 50, 255))
    px(img, 31, 14, (255, 50, 50, 255))
    px(img, 33, 14, (255, 50, 50, 255))
    # 外圈
    draw = ImageDraw.Draw(img)
    draw.ellipse([12, 12, 52, 52], outline=(255, 215, 0, 200), width=2)
    draw_outline(img, (150, 50, 100, 255))
    img.save(os.path.join(OUT_DIR, "T125_lucky_wheel.png"))
    print("T125 saved")

if __name__ == "__main__":
    gen_t121()
    gen_t122()
    gen_t123()
    gen_t124()
    gen_t125()
    print("All T121-T125 sprites generated!")
