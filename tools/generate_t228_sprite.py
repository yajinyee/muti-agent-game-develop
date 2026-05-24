"""
generate_t228_sprite.py — 生成 T228 幸運鏡像對決魚精靈圖（DAY-270）
主題：紫金鏡像魚，左右對稱，帶鏡面反射效果
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

def fill_circle_shaded(img, cx, cy, r, light, mid, dark):
    for dy in range(-r, r+1):
        for dx in range(-r, r+1):
            if dx*dx + dy*dy <= r*r:
                if dx < -r//3 and dy < -r//3:
                    c = light
                elif dx > r//3 and dy > r//3:
                    c = dark
                else:
                    c = mid
                px(img, cx+dx, cy+dy, c)

def draw_fish_body(img, cx, cy, flip=False):
    """畫一條魚的身體（橢圓形）"""
    # 主體（紫色）
    LIGHT = (200, 150, 255, 255)
    MID   = (155, 89, 182, 255)
    DARK  = (100, 50, 130, 255)
    
    # 橢圓身體（加大）
    for dy in range(-10, 11):
        for dx in range(-16, 17):
            if (dx/16)**2 + (dy/10)**2 <= 1.0:
                # 陰影
                if dx < -5 and dy < -3:
                    c = LIGHT
                elif dx > 5 and dy > 3:
                    c = DARK
                else:
                    c = MID
                px(img, cx+dx, cy+dy, c)
    
    # 魚尾（三角形，加大）
    tail_dir = -1 if not flip else 1
    TAIL_MID = (130, 70, 160, 255)
    TAIL_DARK = (90, 40, 110, 255)
    for i in range(12):
        spread = max(1, i * 6 // 10)
        for j in range(-spread, spread+1):
            c = TAIL_DARK if i > 8 else TAIL_MID
            px(img, cx + tail_dir*(16+i), cy+j, c)
    
    # 魚眼（金色，加大）
    eye_x = cx + (-9 if not flip else 9)
    EYE_GOLD = (255, 215, 0, 255)
    EYE_DARK = (50, 30, 0, 255)
    fill_circle(img, eye_x, cy-2, 4, EYE_GOLD)
    fill_circle(img, eye_x, cy-2, 2, EYE_DARK)
    px(img, eye_x-1, cy-3, (255, 255, 200, 255))  # 高光
    
    # 魚鰭（上方，加大）
    FIN_LIGHT = (220, 170, 255, 255)
    FIN_MID   = (180, 120, 220, 255)
    for i in range(8):
        for j in range(i+1):
            c = FIN_LIGHT if j < 3 else FIN_MID
            px(img, cx + tail_dir*(-4+i), cy-10-j, c)
    
    # 魚鱗（細節）
    SCALE = (170, 100, 200, 200)
    for i in range(3):
        for j in range(2):
            sx = cx + tail_dir*(-4 + i*4)
            sy = cy + j*4 - 2
            fill_circle(img, sx, sy, 2, SCALE)

def draw_mirror_effect(img, cx, cy):
    """畫鏡面反射效果（中間的鏡面線）"""
    MIRROR_COLOR = (255, 255, 255, 220)
    MIRROR_GLOW  = (200, 180, 255, 150)
    MIRROR_OUTER = (160, 140, 220, 80)
    
    # 垂直鏡面線（加粗）
    for y in range(cy-14, cy+15):
        px(img, cx, y, MIRROR_COLOR)
        px(img, cx-1, y, MIRROR_GLOW)
        px(img, cx+1, y, MIRROR_GLOW)
        px(img, cx-2, y, MIRROR_OUTER)
        px(img, cx+2, y, MIRROR_OUTER)
    
    # 鏡面光暈（菱形，加大）
    for i in range(7):
        for j in range(-i, i+1):
            alpha = int(180 * (1 - i/7))
            px(img, cx+j, cy-i, (255, 255, 255, alpha))
            px(img, cx+j, cy+i, (255, 255, 255, alpha))
    
    # 鏡面框（上下端點）
    FRAME = (255, 215, 0, 255)
    fill_circle(img, cx, cy-14, 3, FRAME)
    fill_circle(img, cx, cy+14, 3, FRAME)

def draw_sparkles(img):
    """畫星星光點（金色）"""
    GOLD = (255, 215, 0, 255)
    GOLD_DIM = (200, 160, 0, 200)
    
    sparkle_positions = [
        (8, 8), (56, 8), (8, 56), (56, 56),
        (16, 20), (48, 20), (16, 44), (48, 44),
        (32, 6), (32, 58),
    ]
    for (sx, sy) in sparkle_positions:
        px(img, sx, sy, GOLD)
        px(img, sx-1, sy, GOLD_DIM)
        px(img, sx+1, sy, GOLD_DIM)
        px(img, sx, sy-1, GOLD_DIM)
        px(img, sx, sy+1, GOLD_DIM)

def draw_outline(img):
    """加深色輪廓"""
    OUTLINE = (40, 20, 60, 255)
    pixels = img.load()
    w, h = img.size
    result = img.copy()
    
    for y in range(1, h-1):
        for x in range(1, w-1):
            if pixels[x, y][3] > 0:
                # 檢查四個方向是否有透明像素
                for dx, dy in [(0,1),(0,-1),(1,0),(-1,0)]:
                    nx, ny = x+dx, y+dy
                    if 0 <= nx < w and 0 <= ny < h and pixels[nx, ny][3] == 0:
                        result.putpixel((x, y), OUTLINE)
                        break
    return result

def generate():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    
    cx = SIZE // 2  # 32
    cy = SIZE // 2  # 32
    
    # 左側魚（原版，朝右）
    draw_fish_body(img, cx - 14, cy, flip=False)
    
    # 右側魚（鏡像，朝左）
    draw_fish_body(img, cx + 14, cy, flip=True)
    
    # 中間鏡面效果
    draw_mirror_effect(img, cx, cy)
    
    # 星星光點
    draw_sparkles(img)
    
    # 加輪廓
    img = draw_outline(img)
    
    # 儲存
    out_path = os.path.join(OUT_DIR, "T228_mirrorduel.png")
    img.save(out_path)
    
    # 統計非透明像素
    pixels = img.load()
    count = sum(1 for y in range(SIZE) for x in range(SIZE) if pixels[x,y][3] > 0)
    pct = count / (SIZE * SIZE) * 100
    print(f"T228 生成完成：{out_path}")
    print(f"非透明像素：{count} ({pct:.1f}%)")

if __name__ == "__main__":
    generate()
