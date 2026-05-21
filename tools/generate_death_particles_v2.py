"""
generate_death_particles_v2.py
升級死亡粒子特效 — 從 14% 提升到 50%+ 非透明像素
設計：8方向爆炸粒子 + 中心爆炸圓 + 星形光芒 + 碎片
"""
from PIL import Image, ImageDraw
import math
import os

OUT_PATH = "d:/Kiro/client/chiikawa-pixel/assets/sprites/effects/death_particles.png"
SIZE = 48

def draw_pixel(draw, x, y, color):
    """畫一個像素（邊界檢查）"""
    if 0 <= x < SIZE and 0 <= y < SIZE:
        draw.point((x, y), fill=color)

def draw_circle_filled(draw, cx, cy, r, color):
    """畫實心圓（像素風格）"""
    for dy in range(-r, r+1):
        for dx in range(-r, r+1):
            if dx*dx + dy*dy <= r*r:
                draw_pixel(draw, cx+dx, cy+dy, color)

def draw_line_thick(draw, x0, y0, x1, y1, color, thickness=2):
    """畫粗線（像素風格）"""
    steps = max(abs(x1-x0), abs(y1-y0), 1)
    for i in range(steps+1):
        t = i / steps
        x = int(x0 + (x1-x0)*t)
        y = int(y0 + (y1-y0)*t)
        for dy in range(-thickness//2, thickness//2+1):
            for dx in range(-thickness//2, thickness//2+1):
                draw_pixel(draw, x+dx, y+dy, color)

def generate():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    
    cx, cy = SIZE // 2, SIZE // 2
    
    # 顏色定義（金色系爆炸）
    GOLD_BRIGHT = (255, 220, 50, 255)
    GOLD_MID    = (255, 180, 20, 255)
    GOLD_DARK   = (200, 120, 10, 255)
    WHITE       = (255, 255, 255, 255)
    ORANGE      = (255, 120, 30, 255)
    RED_ORANGE  = (255, 80, 20, 200)
    
    # 1. 中心爆炸圓（最亮）
    draw_circle_filled(draw, cx, cy, 8, WHITE)
    draw_circle_filled(draw, cx, cy, 6, GOLD_BRIGHT)
    draw_circle_filled(draw, cx, cy, 4, WHITE)
    
    # 2. 8方向爆炸射線（主要粒子）
    for i in range(8):
        angle = i * math.pi / 4
        # 長射線
        x1 = int(cx + math.cos(angle) * 18)
        y1 = int(cy + math.sin(angle) * 18)
        draw_line_thick(draw, cx, cy, x1, y1, GOLD_MID, 2)
        
        # 射線末端粒子
        draw_circle_filled(draw, x1, y1, 2, GOLD_BRIGHT)
    
    # 3. 4方向次要射線（45度偏移）
    for i in range(4):
        angle = i * math.pi / 2 + math.pi / 8
        x1 = int(cx + math.cos(angle) * 14)
        y1 = int(cy + math.sin(angle) * 14)
        draw_line_thick(draw, cx, cy, x1, y1, GOLD_DARK, 1)
        draw_circle_filled(draw, x1, y1, 1, ORANGE)
    
    # 4. 星形光芒（4個尖角）
    for i in range(4):
        angle = i * math.pi / 2
        # 長尖角
        x_tip = int(cx + math.cos(angle) * 22)
        y_tip = int(cy + math.sin(angle) * 22)
        # 尖角兩側
        perp = angle + math.pi / 2
        x_l = int(cx + math.cos(angle) * 8 + math.cos(perp) * 3)
        y_l = int(cy + math.sin(angle) * 8 + math.sin(perp) * 3)
        x_r = int(cx + math.cos(angle) * 8 - math.cos(perp) * 3)
        y_r = int(cy + math.sin(angle) * 8 - math.sin(perp) * 3)
        
        # 畫三角形尖角
        draw_line_thick(draw, x_l, y_l, x_tip, y_tip, GOLD_BRIGHT, 1)
        draw_line_thick(draw, x_r, y_r, x_tip, y_tip, GOLD_BRIGHT, 1)
        draw_pixel(draw, x_tip, y_tip, WHITE)
    
    # 5. 散落碎片（隨機分布的小點）
    import random
    rng = random.Random(42)  # 固定種子，確保每次生成一致
    for _ in range(35):  # 增加到 35 個碎片
        angle = rng.uniform(0, math.pi * 2)
        dist = rng.uniform(5, 21)
        x = int(cx + math.cos(angle) * dist)
        y = int(cy + math.sin(angle) * dist)
        color = [GOLD_BRIGHT, GOLD_MID, ORANGE, WHITE, RED_ORANGE][rng.randint(0, 4)]
        size = rng.randint(1, 3)
        draw_circle_filled(draw, x, y, size, color)
    
    # 6. 外圈光環（爆炸衝擊波）
    for angle_deg in range(0, 360, 3):
        angle = math.radians(angle_deg)
        r = 21
        x = int(cx + math.cos(angle) * r)
        y = int(cy + math.sin(angle) * r)
        alpha = 180 if angle_deg % 6 == 0 else 100
        draw_pixel(draw, x, y, (255, 200, 50, alpha))
    
    # 儲存
    img.save(OUT_PATH)
    
    # 驗證
    check = Image.open(OUT_PATH).convert("RGBA")
    pixels = list(check.getdata())
    non_transparent = sum(1 for p in pixels if p[3] > 10)
    total = SIZE * SIZE
    pct = non_transparent / total * 100
    print(f"[OK] death_particles.png: {SIZE}x{SIZE}, {pct:.0f}% non-transparent ({non_transparent}/{total})")
    
    return pct > 40

if __name__ == "__main__":
    success = generate()
    if success:
        print("✅ death_particles.png 升級成功！")
    else:
        print("⚠️ 非透明像素比例仍偏低，需要進一步調整")
