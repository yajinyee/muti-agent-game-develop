"""
generate_t229_sprite.py — 生成 T229 幸運倍率重擲魚精靈圖（DAY-271）
主題：橙金骰子魚，帶骰子點數和重擲箭頭效果
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
    """畫重擲魚主體（橙金色）"""
    cx, cy = 28, 32

    LIGHT = (255, 200, 80, 255)   # 亮金橙
    MID   = (255, 140, 0, 255)    # 橙 #FF8C00
    DARK  = (200, 90, 0, 255)     # 深橙

    # 主體橢圓（加大）
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

    # 魚鱗
    SCALE = (230, 120, 0, 200)
    for i in range(4):
        for j in range(3):
            sx = cx - 8 + i*5
            sy = cy - 4 + j*5
            fill_circle(img, sx, sy, 2, SCALE)

    # 魚尾（加大）
    TAIL_MID  = (220, 110, 0, 255)
    TAIL_DARK = (160, 70, 0, 255)
    for i in range(14):
        spread = max(1, i * 8 // 10)
        for j in range(-spread, spread+1):
            c = TAIL_DARK if i > 10 else TAIL_MID
            px(img, cx+18+i, cy+j, c)

    # 魚眼（金色，加大）
    EYE_GOLD = (255, 215, 0, 255)
    EYE_DARK = (60, 30, 0, 255)
    fill_circle(img, cx-11, cy-2, 5, EYE_GOLD)
    fill_circle(img, cx-11, cy-2, 3, EYE_DARK)
    px(img, cx-12, cy-3, (255, 255, 200, 255))

    # 背鰭（加大）
    FIN_LIGHT = (255, 230, 100, 255)
    FIN_MID   = (255, 180, 50, 255)
    for i in range(10):
        for j in range(i+1):
            c = FIN_LIGHT if j < 3 else FIN_MID
            px(img, cx-6+i, cy-12-j, c)

    # 腹鰭
    for i in range(6):
        for j in range(i+1):
            px(img, cx-3+i, cy+12+j, FIN_MID)

def draw_dice_face(img, cx, cy, value):
    """在魚身上畫骰子面（顯示重擲值）"""
    # 骰子背景（白色圓角方形）
    DICE_BG   = (255, 255, 240, 230)
    DICE_DOT  = (50, 20, 0, 255)
    DICE_GOLD = (255, 215, 0, 255)

    # 骰子框（8x8）
    for dy in range(-4, 5):
        for dx in range(-4, 5):
            if abs(dx) <= 3 and abs(dy) <= 3:
                px(img, cx+dx, cy+dy, DICE_BG)

    # 骰子邊框（金色）
    for i in range(-4, 5):
        px(img, cx+i, cy-4, DICE_GOLD)
        px(img, cx+i, cy+4, DICE_GOLD)
        px(img, cx-4, cy+i, DICE_GOLD)
        px(img, cx+4, cy+i, DICE_GOLD)

    # 骰子點數（根據 value 1-6）
    dot_positions = {
        1: [(0, 0)],
        2: [(-2, -2), (2, 2)],
        3: [(-2, -2), (0, 0), (2, 2)],
        4: [(-2, -2), (2, -2), (-2, 2), (2, 2)],
        5: [(-2, -2), (2, -2), (0, 0), (-2, 2), (2, 2)],
        6: [(-2, -2), (2, -2), (-2, 0), (2, 0), (-2, 2), (2, 2)],
    }
    dots = dot_positions.get(value, dot_positions[6])
    for (dx, dy) in dots:
        fill_circle(img, cx+dx, cy+dy, 1, DICE_DOT)

def draw_reroll_arrows(img):
    """畫重擲箭頭（循環箭頭，表示重擲）"""
    ARROW = (255, 215, 0, 220)
    ARROW_DIM = (200, 160, 0, 150)

    # 左上角循環箭頭（簡化版）
    # 上弧
    for i in range(8):
        angle_x = int(6 * (1 - (i/7)**2))
        px(img, 8+i, 8-angle_x//2, ARROW)
        px(img, 8+i, 8-angle_x//2-1, ARROW_DIM)
    # 右弧
    for i in range(8):
        angle_y = int(6 * (1 - (i/7)**2))
        px(img, 16+angle_y//2, 8+i, ARROW)
        px(img, 16+angle_y//2+1, 8+i, ARROW_DIM)
    # 箭頭頭部
    for i in range(3):
        px(img, 8+i, 6+i, ARROW)
        px(img, 8+i, 10-i, ARROW)

def draw_sparkles(img):
    """畫星星光點"""
    GOLD = (255, 215, 0, 255)
    GOLD_DIM = (200, 160, 0, 180)

    positions = [(5, 50), (58, 12), (58, 52), (5, 12), (32, 5), (32, 59)]
    for (sx, sy) in positions:
        px(img, sx, sy, GOLD)
        for dx, dy in [(1,0),(-1,0),(0,1),(0,-1)]:
            px(img, sx+dx, sy+dy, GOLD_DIM)

def draw_outline(img):
    OUTLINE = (60, 30, 0, 255)
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
    draw_dice_face(img, 28, 32, 6)  # 骰子顯示 6（最高值）
    draw_reroll_arrows(img)
    draw_sparkles(img)
    img = draw_outline(img)

    out_path = os.path.join(OUT_DIR, "T229_reroll.png")
    img.save(out_path)

    pixels = img.load()
    count = sum(1 for y in range(SIZE) for x in range(SIZE) if pixels[x,y][3] > 0)
    pct = count / (SIZE * SIZE) * 100
    print(f"T229 生成完成：{out_path}")
    print(f"非透明像素：{count} ({pct:.1f}%)")

if __name__ == "__main__":
    generate()
