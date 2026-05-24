"""
generate_t223_sprite.py — 生成 T223 幸運競速賽魚精靈圖（DAY-265）
主題：橙紅競速魚，帶速度線條和旗幟效果
"""
from PIL import Image, ImageDraw
import os

OUT_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
SIZE = 64

def px(img, x, y, color):
    if 0 <= x < SIZE and 0 <= y < SIZE:
        img.putpixel((x, y), color)

def fill_circle(img, cx, cy, r, color):
    for dy in range(-r, r+1):
        for dx in range(-r, r+1):
            if dx*dx + dy*dy <= r*r:
                px(img, cx+dx, cy+dy, color)

def draw_fish_body(img):
    """畫競速魚主體（橙紅色，流線型）"""
    cx, cy = 28, 32

    # 主體橢圓（橙紅色，流線型，加大）
    LIGHT = (255, 160, 80, 255)   # 亮橙
    MID   = (255, 107, 53, 255)   # 橙紅 #FF6B35
    DARK  = (200, 60, 20, 255)    # 深橙紅

    for dy in range(-12, 13):
        for dx in range(-18, 19):
            if (dx/18)**2 + (dy/12)**2 <= 1.0:
                if dx < -6 and dy < -4:
                    c = LIGHT
                elif dx > 6 and dy > 4:
                    c = DARK
                else:
                    c = MID
                px(img, cx+dx, cy+dy, c)

    # 魚鱗（細節）
    SCALE = (230, 90, 40, 200)
    for i in range(4):
        for j in range(3):
            sx = cx - 8 + i*5
            sy = cy - 4 + j*5
            fill_circle(img, sx, sy, 2, SCALE)

    # 魚尾（向右，三角形，加大）
    TAIL_MID  = (220, 80, 30, 255)
    TAIL_DARK = (160, 50, 10, 255)
    for i in range(14):
        spread = max(1, i * 8 // 10)
        for j in range(-spread, spread+1):
            c = TAIL_DARK if i > 10 else TAIL_MID
            px(img, cx+18+i, cy+j, c)

    # 魚眼（金色，加大）
    EYE_GOLD = (255, 215, 0, 255)
    EYE_DARK = (50, 20, 0, 255)
    fill_circle(img, cx-11, cy-2, 5, EYE_GOLD)
    fill_circle(img, cx-11, cy-2, 3, EYE_DARK)
    px(img, cx-12, cy-3, (255, 255, 200, 255))  # 高光

    # 背鰭（上方，加大）
    FIN_LIGHT = (255, 200, 100, 255)
    FIN_MID   = (255, 150, 50, 255)
    for i in range(10):
        for j in range(i+1):
            c = FIN_LIGHT if j < 3 else FIN_MID
            px(img, cx-6+i, cy-12-j, c)

    # 腹鰭（下方）
    for i in range(6):
        for j in range(i+1):
            px(img, cx-3+i, cy+12+j, FIN_MID)

def draw_speed_lines(img):
    """畫速度線條（向左延伸，表示高速）"""
    SPEED_GOLD   = (255, 215, 0, 200)
    SPEED_ORANGE = (255, 140, 0, 150)
    SPEED_WHITE  = (255, 255, 255, 120)

    # 三條速度線（從魚尾向右延伸）
    speed_lines = [
        (50, 26, 8, SPEED_GOLD),
        (52, 30, 10, SPEED_ORANGE),
        (50, 34, 7, SPEED_WHITE),
    ]
    for (sx, sy, length, color) in speed_lines:
        for i in range(length):
            alpha = int(color[3] * (1 - i/length))
            c = (color[0], color[1], color[2], alpha)
            px(img, sx+i, sy, c)
            if i < length - 2:
                px(img, sx+i, sy-1, (color[0], color[1], color[2], alpha//2))
                px(img, sx+i, sy+1, (color[0], color[1], color[2], alpha//2))

def draw_checkered_flag(img):
    """畫格子旗（左上角，表示競速）"""
    FLAG_WHITE = (255, 255, 255, 230)
    FLAG_BLACK = (30, 30, 30, 230)
    FLAG_POLE  = (180, 120, 50, 255)

    # 旗桿
    for y in range(4, 18):
        px(img, 4, y, FLAG_POLE)

    # 旗面（4x6 格子旗）
    flag_x, flag_y = 5, 4
    for row in range(3):
        for col in range(4):
            c = FLAG_WHITE if (row + col) % 2 == 0 else FLAG_BLACK
            for dy in range(2):
                for dx in range(2):
                    px(img, flag_x + col*2 + dx, flag_y + row*2 + dy, c)

def draw_rank_badge(img):
    """畫排名徽章（🏆 金色圓形）"""
    BADGE_GOLD  = (255, 215, 0, 255)
    BADGE_DARK  = (180, 140, 0, 255)
    BADGE_TEXT  = (100, 60, 0, 255)

    # 金色圓形徽章
    fill_circle(img, 14, 50, 7, BADGE_GOLD)
    fill_circle(img, 14, 50, 5, BADGE_DARK)

    # 「1」字樣（簡化像素）
    for y in range(47, 54):
        px(img, 14, y, (255, 255, 200, 255))
    px(img, 13, 48, (255, 255, 200, 255))

def draw_outline(img):
    """加深色輪廓"""
    OUTLINE = (60, 20, 0, 255)
    pixels = img.load()
    w, h = img.size
    result = img.copy()

    for y in range(1, h-1):
        for x in range(1, w-1):
            if pixels[x, y][3] > 0:
                for dx2, dy2 in [(0,1),(0,-1),(1,0),(-1,0)]:
                    nx, ny = x+dx2, y+dy2
                    if 0 <= nx < w and 0 <= ny < h and pixels[nx, ny][3] == 0:
                        result.putpixel((x, y), OUTLINE)
                        break
    return result

def generate():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))

    draw_fish_body(img)
    draw_speed_lines(img)
    draw_checkered_flag(img)
    draw_rank_badge(img)
    img = draw_outline(img)

    out_path = os.path.join(OUT_DIR, "T223_speedrace.png")
    img.save(out_path)

    pixels = img.load()
    count = sum(1 for y in range(SIZE) for x in range(SIZE) if pixels[x,y][3] > 0)
    pct = count / (SIZE * SIZE) * 100
    print(f"T223 生成完成：{out_path}")
    print(f"非透明像素：{count} ({pct:.1f}%)")

if __name__ == "__main__":
    generate()
