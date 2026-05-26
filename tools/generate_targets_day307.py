"""
generate_targets_day307.py — DAY-307 T146-T150 精靈圖生成
target-pixel-agent 負責維護
生成 5 個新 Lucky 魚精靈圖（64x64 像素）
"""
import os
from PIL import Image, ImageDraw

OUT_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
SIZE = 64

def fill_circle(draw, cx, cy, r, color):
    for y in range(cy - r, cy + r + 1):
        for x in range(cx - r, cx + r + 1):
            if (x - cx) ** 2 + (y - cy) ** 2 <= r * r:
                if 0 <= x < SIZE and 0 <= y < SIZE:
                    draw.point((x, y), fill=color)

def fill_circle_shaded(draw, cx, cy, r, base_color):
    br, bg, bb = base_color[:3]
    for y in range(cy - r, cy + r + 1):
        for x in range(cx - r, cx + r + 1):
            dist = ((x - cx) ** 2 + (y - cy) ** 2) ** 0.5
            if dist <= r:
                if 0 <= x < SIZE and 0 <= y < SIZE:
                    # 左上亮，右下暗
                    shade = 1.0 - 0.3 * ((x - cx + y - cy) / (2 * r))
                    shade = max(0.5, min(1.3, shade))
                    c = (min(255, int(br * shade)), min(255, int(bg * shade)), min(255, int(bb * shade)), 255)
                    draw.point((x, y), fill=c)

def draw_fish_body(draw, cx, cy, rx, ry, color):
    """橢圓魚身"""
    for y in range(cy - ry, cy + ry + 1):
        for x in range(cx - rx, cx + rx + 1):
            if (x - cx) ** 2 / (rx * rx) + (y - cy) ** 2 / (ry * ry) <= 1.0:
                if 0 <= x < SIZE and 0 <= y < SIZE:
                    # 陰影
                    shade = 1.0 - 0.25 * ((x - cx + y - cy) / (rx + ry))
                    shade = max(0.6, min(1.2, shade))
                    r, g, b = color[:3]
                    c = (min(255, int(r * shade)), min(255, int(g * shade)), min(255, int(b * shade)), 255)
                    draw.point((x, y), fill=c)

def draw_outline(draw, cx, cy, rx, ry, outline_color, thickness=1):
    """橢圓輪廓"""
    for y in range(cy - ry - thickness, cy + ry + thickness + 1):
        for x in range(cx - rx - thickness, cx + rx + thickness + 1):
            dist = (x - cx) ** 2 / ((rx + thickness) ** 2) + (y - cy) ** 2 / ((ry + thickness) ** 2)
            inner = (x - cx) ** 2 / (rx ** 2) + (y - cy) ** 2 / (ry ** 2)
            if dist <= 1.0 and inner > 1.0:
                if 0 <= x < SIZE and 0 <= y < SIZE:
                    draw.point((x, y), fill=outline_color)

def draw_eye(draw, ex, ey):
    """魚眼"""
    fill_circle(draw, ex, ey, 3, (255, 255, 255, 255))
    fill_circle(draw, ex, ey, 2, (20, 20, 20, 255))
    draw.point((ex - 1, ey - 1), fill=(255, 255, 255, 255))

def draw_tail(draw, tx, ty, color):
    """魚尾（三角形）"""
    for y in range(ty - 8, ty + 9):
        width = int(8 * abs(y - ty) / 8)
        for x in range(tx, tx + width + 1):
            if 0 <= x < SIZE and 0 <= y < SIZE:
                draw.point((x, y), fill=color)

# ── T146 幸運量子魚 ──────────────────────────────────────────
def gen_t146():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    
    # 量子魚：青藍色橢圓魚身 + 量子粒子環 + 青藍光暈
    BODY = (0, 180, 220, 255)
    OUTLINE = (0, 80, 120, 255)
    GLOW = (0, 230, 255, 180)
    PARTICLE = (100, 255, 255, 255)
    
    # 光暈
    for r in range(26, 30):
        for angle_deg in range(0, 360, 3):
            import math
            angle = math.radians(angle_deg)
            x = int(32 + r * math.cos(angle))
            y = int(32 + r * math.sin(angle))
            if 0 <= x < SIZE and 0 <= y < SIZE:
                draw.point((x, y), fill=GLOW)
    
    # 魚尾
    draw_tail(draw, 46, 32, (0, 120, 160, 255))
    
    # 魚身
    draw_fish_body(draw, 28, 32, 18, 12, BODY)
    draw_outline(draw, 28, 32, 18, 12, OUTLINE)
    
    # 量子粒子環（8個粒子）
    import math
    for i in range(8):
        angle = math.radians(i * 45)
        px = int(28 + 22 * math.cos(angle))
        py = int(32 + 22 * math.sin(angle))
        fill_circle(draw, px, py, 2, PARTICLE)
    
    # 眼睛
    draw_eye(draw, 18, 28)
    
    # 量子符號（中心）
    draw.text((22, 26), "⚛", fill=(255, 255, 255, 200))
    
    img.save(os.path.join(OUT_DIR, "T146_quantum.png"))
    pixels = sum(1 for p in img.getdata() if p[3] > 0)
    print(f"T146 量子魚: {pixels} 非透明像素 ({pixels/SIZE/SIZE*100:.1f}%)")

# ── T147 幸運超新星魚 ────────────────────────────────────────
def gen_t147():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    
    BODY = (220, 80, 20, 255)
    OUTLINE = (120, 30, 0, 255)
    GLOW = (255, 120, 40, 160)
    CORE = (255, 220, 100, 255)
    
    import math
    
    # 爆炸光芒（8方向）
    for i in range(8):
        angle = math.radians(i * 45)
        for r in range(20, 30):
            x = int(32 + r * math.cos(angle))
            y = int(32 + r * math.sin(angle))
            if 0 <= x < SIZE and 0 <= y < SIZE:
                alpha = int(200 * (30 - r) / 10)
                draw.point((x, y), fill=(255, 150, 50, alpha))
    
    # 外光暈
    for r in range(24, 28):
        for angle_deg in range(0, 360, 2):
            angle = math.radians(angle_deg)
            x = int(32 + r * math.cos(angle))
            y = int(32 + r * math.sin(angle))
            if 0 <= x < SIZE and 0 <= y < SIZE:
                draw.point((x, y), fill=GLOW)
    
    # 魚尾
    draw_tail(draw, 46, 32, (160, 40, 0, 255))
    
    # 魚身
    draw_fish_body(draw, 28, 32, 18, 12, BODY)
    draw_outline(draw, 28, 32, 18, 12, OUTLINE)
    
    # 核心
    fill_circle_shaded(draw, 28, 32, 7, (255, 200, 80))
    
    # 眼睛
    draw_eye(draw, 18, 28)
    
    img.save(os.path.join(OUT_DIR, "T147_supernova.png"))
    pixels = sum(1 for p in img.getdata() if p[3] > 0)
    print(f"T147 超新星魚: {pixels} 非透明像素 ({pixels/SIZE/SIZE*100:.1f}%)")

# ── T148 幸運無限魚 ──────────────────────────────────────────
def gen_t148():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    
    BODY = (120, 40, 200, 255)
    OUTLINE = (60, 10, 100, 255)
    GLOW = (160, 80, 255, 150)
    SYMBOL = (255, 220, 255, 255)
    
    import math
    
    # 無限光暈
    for r in range(26, 30):
        for angle_deg in range(0, 360, 3):
            angle = math.radians(angle_deg)
            x = int(32 + r * math.cos(angle))
            y = int(32 + r * math.sin(angle))
            if 0 <= x < SIZE and 0 <= y < SIZE:
                draw.point((x, y), fill=GLOW)
    
    # 魚尾
    draw_tail(draw, 46, 32, (80, 20, 140, 255))
    
    # 魚身
    draw_fish_body(draw, 28, 32, 18, 12, BODY)
    draw_outline(draw, 28, 32, 18, 12, OUTLINE)
    
    # 無限符號（∞）— 用兩個圓圈模擬
    fill_circle(draw, 23, 32, 5, (200, 150, 255, 200))
    fill_circle(draw, 33, 32, 5, (200, 150, 255, 200))
    fill_circle(draw, 23, 32, 3, (255, 220, 255, 255))
    fill_circle(draw, 33, 32, 3, (255, 220, 255, 255))
    
    # 眼睛
    draw_eye(draw, 16, 28)
    
    img.save(os.path.join(OUT_DIR, "T148_infinite.png"))
    pixels = sum(1 for p in img.getdata() if p[3] > 0)
    print(f"T148 無限魚: {pixels} 非透明像素 ({pixels/SIZE/SIZE*100:.1f}%)")

# ── T149 幸運創世魚 ──────────────────────────────────────────
def gen_t149():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    
    BODY = (200, 160, 0, 255)
    OUTLINE = (100, 70, 0, 255)
    GLOW = (255, 220, 50, 160)
    CROWN = (255, 200, 0, 255)
    
    import math
    
    # 神聖光芒（12方向）
    for i in range(12):
        angle = math.radians(i * 30)
        for r in range(22, 32):
            x = int(32 + r * math.cos(angle))
            y = int(32 + r * math.sin(angle))
            if 0 <= x < SIZE and 0 <= y < SIZE:
                alpha = int(180 * (32 - r) / 10)
                draw.point((x, y), fill=(255, 220, 50, alpha))
    
    # 外光暈
    for r in range(24, 28):
        for angle_deg in range(0, 360, 2):
            angle = math.radians(angle_deg)
            x = int(32 + r * math.cos(angle))
            y = int(32 + r * math.sin(angle))
            if 0 <= x < SIZE and 0 <= y < SIZE:
                draw.point((x, y), fill=GLOW)
    
    # 魚尾
    draw_tail(draw, 48, 32, (140, 100, 0, 255))
    
    # 魚身（大型）
    draw_fish_body(draw, 28, 32, 20, 14, BODY)
    draw_outline(draw, 28, 32, 20, 14, OUTLINE)
    
    # 王冠（頂部）
    for i in range(5):
        x = 18 + i * 5
        draw.line([(x, 18), (x, 14)], fill=CROWN, width=2)
    fill_circle(draw, 28, 13, 3, CROWN)
    
    # 眼睛（紅色）
    fill_circle(draw, 17, 28, 3, (255, 255, 255, 255))
    fill_circle(draw, 17, 28, 2, (200, 0, 0, 255))
    
    img.save(os.path.join(OUT_DIR, "T149_genesis.png"))
    pixels = sum(1 for p in img.getdata() if p[3] > 0)
    print(f"T149 創世魚: {pixels} 非透明像素 ({pixels/SIZE/SIZE*100:.1f}%)")

# ── T150 幸運重生魚 ──────────────────────────────────────────
def gen_t150():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    
    BODY = (200, 80, 20, 255)
    OUTLINE = (100, 30, 0, 255)
    GLOW = (255, 120, 40, 150)
    FLAME = (255, 180, 0, 255)
    
    import math
    
    # 鳳凰火焰光暈
    for r in range(25, 30):
        for angle_deg in range(0, 360, 3):
            angle = math.radians(angle_deg)
            x = int(32 + r * math.cos(angle))
            y = int(32 + r * math.sin(angle))
            if 0 <= x < SIZE and 0 <= y < SIZE:
                draw.point((x, y), fill=GLOW)
    
    # 火焰（頂部）
    for i in range(3):
        fx = 22 + i * 8
        for fy in range(10, 20):
            alpha = int(200 * (20 - fy) / 10)
            draw.point((fx, fy), fill=(255, 150 + i * 30, 0, alpha))
    
    # 魚尾
    draw_tail(draw, 46, 32, (140, 40, 0, 255))
    
    # 魚身
    draw_fish_body(draw, 28, 32, 18, 12, BODY)
    draw_outline(draw, 28, 32, 18, 12, OUTLINE)
    
    # 鳳凰紋路（金色條紋）
    for i in range(3):
        y = 26 + i * 4
        for x in range(16, 40):
            if (x + y) % 4 == 0:
                draw.point((x, y), fill=FLAME)
    
    # 眼睛
    draw_eye(draw, 18, 28)
    
    img.save(os.path.join(OUT_DIR, "T150_rebirth.png"))
    pixels = sum(1 for p in img.getdata() if p[3] > 0)
    print(f"T150 重生魚: {pixels} 非透明像素 ({pixels/SIZE/SIZE*100:.1f}%)")

if __name__ == "__main__":
    os.makedirs(OUT_DIR, exist_ok=True)
    print("=== DAY-307 T146-T150 精靈圖生成 ===")
    gen_t146()
    gen_t147()
    gen_t148()
    gen_t149()
    gen_t150()
    print("=== 完成 ===")
