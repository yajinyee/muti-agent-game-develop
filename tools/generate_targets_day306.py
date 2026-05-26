"""
generate_targets_day306.py — DAY-306 T141-T145 精靈圖生成
target-pixel-agent 負責維護
生成 T141 龍捲風魚、T142 地震魚、T143 火山魚、T144 星際魚、T145 神龍魚
"""
import os
import struct
import zlib

OUT_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
SIZE = 64

def make_png(pixels, w, h):
    """從 RGBA 像素陣列生成 PNG bytes"""
    def chunk(name, data):
        c = zlib.crc32(name + data) & 0xFFFFFFFF
        return struct.pack('>I', len(data)) + name + data + struct.pack('>I', c)
    
    ihdr = struct.pack('>IIBBBBB', w, h, 8, 2, 0, 0, 0)
    raw = b''
    for y in range(h):
        raw += b'\x00'
        for x in range(w):
            r, g, b, a = pixels[y][x]
            raw += bytes([r, g, b, a])
    
    compressed = zlib.compress(raw, 9)
    return (b'\x89PNG\r\n\x1a\n' +
            chunk(b'IHDR', ihdr) +
            chunk(b'IDAT', compressed) +
            chunk(b'IEND', b''))

def save_png(filename, pixels, w=SIZE, h=SIZE):
    data = make_png(pixels, w, h)
    with open(filename, 'wb') as f:
        f.write(data)
    print(f"  Saved: {os.path.basename(filename)} ({len(data)} bytes)")

def empty_pixels(w=SIZE, h=SIZE):
    return [[(0,0,0,0)] * w for _ in range(h)]

def fill_circle(pixels, cx, cy, r, color, w=SIZE, h=SIZE):
    for y in range(max(0, cy-r), min(h, cy+r+1)):
        for x in range(max(0, cx-r), min(w, cx+r+1)):
            if (x-cx)**2 + (y-cy)**2 <= r*r:
                pixels[y][x] = color

def fill_circle_shaded(pixels, cx, cy, r, light, mid, dark, w=SIZE, h=SIZE):
    for y in range(max(0, cy-r), min(h, cy+r+1)):
        for x in range(max(0, cx-r), min(w, cx+r+1)):
            if (x-cx)**2 + (y-cy)**2 <= r*r:
                if x < cx and y < cy:
                    pixels[y][x] = light
                elif x > cx and y > cy:
                    pixels[y][x] = dark
                else:
                    pixels[y][x] = mid

def draw_line(pixels, x0, y0, x1, y1, color, w=SIZE, h=SIZE):
    dx = abs(x1-x0); dy = abs(y1-y0)
    sx = 1 if x0 < x1 else -1
    sy = 1 if y0 < y1 else -1
    err = dx - dy
    while True:
        if 0 <= x0 < w and 0 <= y0 < h:
            pixels[y0][x0] = color
        if x0 == x1 and y0 == y1:
            break
        e2 = 2 * err
        if e2 > -dy:
            err -= dy; x0 += sx
        if e2 < dx:
            err += dx; y0 += sy

# ── T141 龍捲風魚 ─────────────────────────────────────────────
def gen_t141():
    """青綠色龍捲風魚：橢圓魚身 + 旋轉龍捲風紋路 + 青綠光環"""
    p = empty_pixels()
    BODY_L = (100, 240, 200, 255)
    BODY_M = (0, 200, 160, 255)
    BODY_D = (0, 140, 110, 255)
    OUTLINE = (0, 80, 60, 255)
    GLOW = (150, 255, 230, 200)
    WHITE = (255, 255, 255, 255)
    
    # 橢圓魚身
    cx, cy = 32, 32
    for y in range(10, 54):
        for x in range(8, 56):
            dx = (x - cx) / 24.0
            dy = (y - cy) / 20.0
            if dx*dx + dy*dy <= 1.0:
                if x < cx and y < cy:
                    p[y][x] = BODY_L
                elif x > cx and y > cy:
                    p[y][x] = BODY_D
                else:
                    p[y][x] = BODY_M
    
    # 輪廓
    for y in range(9, 55):
        for x in range(7, 57):
            dx = (x - cx) / 24.0
            dy = (y - cy) / 20.0
            if 0.85 <= dx*dx + dy*dy <= 1.1:
                p[y][x] = OUTLINE
    
    # 龍捲風旋轉紋路（螺旋線）
    import math
    for t in range(0, 360, 8):
        rad = math.radians(t)
        r = 8 + t / 30.0
        x = int(cx + r * math.cos(rad))
        y = int(cy + r * math.sin(rad))
        if 0 <= x < SIZE and 0 <= y < SIZE:
            p[y][x] = GLOW
    
    # 眼睛
    fill_circle(p, 38, 28, 3, WHITE)
    fill_circle(p, 38, 28, 2, (0, 60, 50, 255))
    p[27][39] = WHITE
    
    # 魚尾
    for y in range(22, 42):
        for x in range(52, 62):
            if abs(y - 32) < (62 - x) * 0.6:
                p[y][x] = BODY_M
    
    save_png(os.path.join(OUT_DIR, "T141_tornado.png"), p)

# ── T142 地震魚 ───────────────────────────────────────────────
def gen_t142():
    """棕橙色地震魚：圓形魚身 + 裂縫紋路 + 震波光環"""
    p = empty_pixels()
    BODY_L = (220, 140, 60, 255)
    BODY_M = (180, 100, 30, 255)
    BODY_D = (120, 60, 10, 255)
    OUTLINE = (60, 30, 5, 255)
    CRACK = (255, 200, 100, 255)
    
    # 圓形魚身
    fill_circle_shaded(p, 32, 32, 22, BODY_L, BODY_M, BODY_D)
    
    # 輪廓
    for y in range(SIZE):
        for x in range(SIZE):
            if (x-32)**2 + (y-32)**2 <= 22*22:
                if (x-32)**2 + (y-32)**2 >= 19*19:
                    p[y][x] = OUTLINE
    
    # 裂縫紋路（Z字形）
    draw_line(p, 20, 20, 32, 32, CRACK)
    draw_line(p, 32, 32, 44, 20, CRACK)
    draw_line(p, 20, 44, 32, 32, CRACK)
    draw_line(p, 32, 32, 44, 44, CRACK)
    
    # 震波光環
    import math
    for t in range(0, 360, 5):
        rad = math.radians(t)
        x = int(32 + 26 * math.cos(rad))
        y = int(32 + 26 * math.sin(rad))
        if 0 <= x < SIZE and 0 <= y < SIZE:
            p[y][x] = (255, 180, 80, 180)
    
    # 眼睛
    fill_circle(p, 38, 26, 3, (255, 255, 255, 255))
    fill_circle(p, 38, 26, 2, (60, 30, 5, 255))
    p[25][39] = (255, 255, 255, 255)
    
    save_png(os.path.join(OUT_DIR, "T142_earthquake.png"), p)

# ── T143 火山魚 ───────────────────────────────────────────────
def gen_t143():
    """火紅色火山魚：橢圓魚身 + 火焰光芒 + 熔岩核心"""
    p = empty_pixels()
    BODY_L = (255, 120, 50, 255)
    BODY_M = (220, 60, 10, 255)
    BODY_D = (160, 30, 5, 255)
    OUTLINE = (80, 15, 0, 255)
    LAVA = (255, 200, 50, 255)
    CORE = (255, 240, 180, 255)
    
    # 橢圓魚身
    cx, cy = 32, 32
    for y in range(12, 52):
        for x in range(10, 54):
            dx = (x - cx) / 22.0
            dy = (y - cy) / 18.0
            if dx*dx + dy*dy <= 1.0:
                if x < cx and y < cy:
                    p[y][x] = BODY_L
                elif x > cx and y > cy:
                    p[y][x] = BODY_D
                else:
                    p[y][x] = BODY_M
    
    # 輪廓
    for y in range(11, 53):
        for x in range(9, 55):
            dx = (x - cx) / 22.0
            dy = (y - cy) / 18.0
            if 0.85 <= dx*dx + dy*dy <= 1.1:
                p[y][x] = OUTLINE
    
    # 熔岩核心
    fill_circle(p, 32, 32, 8, LAVA)
    fill_circle(p, 32, 32, 5, CORE)
    
    # 火焰光芒（4方向）
    import math
    for angle in [0, 90, 180, 270]:
        rad = math.radians(angle)
        for r in range(24, 32):
            x = int(32 + r * math.cos(rad))
            y = int(32 + r * math.sin(rad))
            if 0 <= x < SIZE and 0 <= y < SIZE:
                p[y][x] = LAVA
    
    # 眼睛
    fill_circle(p, 38, 26, 3, (255, 255, 255, 255))
    fill_circle(p, 38, 26, 2, (80, 15, 0, 255))
    p[25][39] = (255, 255, 255, 255)
    
    # 魚尾
    for y in range(24, 40):
        for x in range(52, 62):
            if abs(y - 32) < (62 - x) * 0.5:
                p[y][x] = BODY_M
    
    save_png(os.path.join(OUT_DIR, "T143_volcano.png"), p)

# ── T144 星際魚 ───────────────────────────────────────────────
def gen_t144():
    """紫色星際魚：橢圓魚身 + 8方向光束 + 星際光環"""
    p = empty_pixels()
    BODY_L = (180, 100, 255, 255)
    BODY_M = (120, 40, 200, 255)
    BODY_D = (70, 10, 140, 255)
    OUTLINE = (30, 5, 70, 255)
    RAY = (220, 180, 255, 200)
    STAR = (255, 240, 255, 255)
    
    # 橢圓魚身
    cx, cy = 32, 32
    for y in range(12, 52):
        for x in range(10, 54):
            dx = (x - cx) / 22.0
            dy = (y - cy) / 18.0
            if dx*dx + dy*dy <= 1.0:
                if x < cx and y < cy:
                    p[y][x] = BODY_L
                elif x > cx and y > cy:
                    p[y][x] = BODY_D
                else:
                    p[y][x] = BODY_M
    
    # 輪廓
    for y in range(11, 53):
        for x in range(9, 55):
            dx = (x - cx) / 22.0
            dy = (y - cy) / 18.0
            if 0.85 <= dx*dx + dy*dy <= 1.1:
                p[y][x] = OUTLINE
    
    # 8方向光束
    import math
    for angle in range(0, 360, 45):
        rad = math.radians(angle)
        for r in range(22, 32):
            x = int(32 + r * math.cos(rad))
            y = int(32 + r * math.sin(rad))
            if 0 <= x < SIZE and 0 <= y < SIZE:
                p[y][x] = RAY
    
    # 中心星形
    fill_circle(p, 32, 32, 6, STAR)
    fill_circle(p, 32, 32, 4, BODY_L)
    
    # 眼睛
    fill_circle(p, 38, 26, 3, (255, 255, 255, 255))
    fill_circle(p, 38, 26, 2, (30, 5, 70, 255))
    p[25][39] = (255, 255, 255, 255)
    
    # 魚尾
    for y in range(24, 40):
        for x in range(52, 62):
            if abs(y - 32) < (62 - x) * 0.5:
                p[y][x] = BODY_M
    
    save_png(os.path.join(OUT_DIR, "T144_cosmic_ray.png"), p)

# ── T145 神龍魚 ───────────────────────────────────────────────
def gen_t145():
    """金色神龍魚：大型圓形魚身 + 龍鱗紋路 + 金色光環 + 王冠"""
    p = empty_pixels()
    BODY_L = (255, 230, 100, 255)
    BODY_M = (220, 180, 30, 255)
    BODY_D = (160, 120, 10, 255)
    OUTLINE = (80, 50, 0, 255)
    SCALE = (255, 200, 50, 200)
    CROWN = (255, 240, 150, 255)
    RED_EYE = (255, 50, 50, 255)
    
    # 大型圓形魚身
    fill_circle_shaded(p, 32, 34, 24, BODY_L, BODY_M, BODY_D)
    
    # 輪廓
    for y in range(SIZE):
        for x in range(SIZE):
            if (x-32)**2 + (y-34)**2 <= 24*24:
                if (x-32)**2 + (y-34)**2 >= 21*21:
                    p[y][x] = OUTLINE
    
    # 龍鱗紋路（菱形格）
    for y in range(14, 54, 6):
        for x in range(12, 52, 6):
            if (x-32)**2 + (y-34)**2 <= 20*20:
                p[y][x] = SCALE
                if x+3 < SIZE and y+3 < SIZE:
                    p[y+3][x+3] = SCALE
    
    # 王冠（頂部）
    for i in range(5):
        cx_c = 20 + i * 6
        for h in range(4):
            if 0 <= cx_c < SIZE and 8 + h < SIZE:
                p[8 + h][cx_c] = CROWN
    
    # 金色光環
    import math
    for t in range(0, 360, 4):
        rad = math.radians(t)
        x = int(32 + 28 * math.cos(rad))
        y = int(34 + 28 * math.sin(rad))
        if 0 <= x < SIZE and 0 <= y < SIZE:
            p[y][x] = (255, 220, 80, 160)
    
    # 紅色龍眼
    fill_circle(p, 38, 28, 4, (255, 255, 255, 255))
    fill_circle(p, 38, 28, 3, RED_EYE)
    fill_circle(p, 38, 28, 1, (80, 0, 0, 255))
    p[27][39] = (255, 255, 255, 255)
    
    # 龍鬚（兩側）
    draw_line(p, 20, 28, 8, 20, CROWN)
    draw_line(p, 20, 28, 8, 36, CROWN)
    
    save_png(os.path.join(OUT_DIR, "T145_divine_dragon.png"), p)

if __name__ == "__main__":
    os.makedirs(OUT_DIR, exist_ok=True)
    print("=== DAY-306 T141-T145 精靈圖生成 ===")
    gen_t141()
    gen_t142()
    gen_t143()
    gen_t144()
    gen_t145()
    print("=== 完成 ===")
