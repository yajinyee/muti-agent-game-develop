"""
T130 幸運崩潰魚精靈圖生成
視覺主題：深紅漸層魚身 + 崩潰裂縫紋路 + 上升倍率數字 + 爆炸光芒 + 警告符號
"""
import struct
import zlib
import os

SIZE = 64
OUT = r'd:\Kiro\client\chiikawa-pixel\assets\sprites\targets\T130_crash_fish.png'

def make_png(pixels, w, h):
    def chunk(name, data):
        c = struct.pack('>I', len(data)) + name + data
        return c + struct.pack('>I', zlib.crc32(c[4:]) & 0xFFFFFFFF)
    
    raw = b''
    for y in range(h):
        raw += b'\x00'
        for x in range(w):
            r, g, b, a = pixels[y][x]
            raw += bytes([r, g, b, a])
    
    ihdr = struct.pack('>IIBBBBB', w, h, 8, 2, 0, 0, 0)
    # RGBA
    ihdr = struct.pack('>II', w, h) + bytes([8, 6, 0, 0, 0])
    
    png = b'\x89PNG\r\n\x1a\n'
    png += chunk(b'IHDR', ihdr)
    png += chunk(b'IDAT', zlib.compress(raw))
    png += chunk(b'IEND', b'')
    return png

def make_png_rgba(pixels, w, h):
    """生成 RGBA PNG"""
    def adler32(data):
        return zlib.adler32(data) & 0xFFFFFFFF
    
    def crc32(data):
        return zlib.crc32(data) & 0xFFFFFFFF
    
    def chunk(tag, data):
        length = struct.pack('>I', len(data))
        body = tag + data
        checksum = struct.pack('>I', crc32(body))
        return length + body + checksum
    
    # IHDR
    ihdr_data = struct.pack('>IIBBBBB', w, h, 8, 6, 0, 0, 0)
    
    # IDAT
    raw_rows = b''
    for y in range(h):
        raw_rows += b'\x00'  # filter type
        for x in range(w):
            r, g, b, a = pixels[y][x]
            raw_rows += bytes([r, g, b, a])
    
    compressed = zlib.compress(raw_rows, 9)
    
    png = b'\x89PNG\r\n\x1a\n'
    png += chunk(b'IHDR', ihdr_data)
    png += chunk(b'IDAT', compressed)
    png += chunk(b'IEND', b'')
    return png

# 初始化透明畫布
pixels = [[(0, 0, 0, 0)] * SIZE for _ in range(SIZE)]

def px(x, y, r, g, b, a=255):
    if 0 <= x < SIZE and 0 <= y < SIZE:
        pixels[y][x] = (r, g, b, a)

def fill_circle(cx, cy, radius, r, g, b, a=255):
    for dy in range(-radius, radius+1):
        for dx in range(-radius, radius+1):
            if dx*dx + dy*dy <= radius*radius:
                px(cx+dx, cy+dy, r, g, b, a)

def fill_circle_shaded(cx, cy, radius, base_r, base_g, base_b):
    """帶陰影的圓形"""
    for dy in range(-radius, radius+1):
        for dx in range(-radius, radius+1):
            dist = (dx*dx + dy*dy) ** 0.5
            if dist <= radius:
                # 左上亮，右下暗
                shade = 1.0 - (dx + dy) / (radius * 2.5)
                shade = max(0.6, min(1.3, shade))
                r2 = min(255, int(base_r * shade))
                g2 = min(255, int(base_g * shade))
                b2 = min(255, int(base_b * shade))
                px(cx+dx, cy+dy, r2, g2, b2)

def draw_line(x0, y0, x1, y1, r, g, b, a=255, width=1):
    """Bresenham 直線"""
    dx = abs(x1 - x0)
    dy = abs(y1 - y0)
    sx = 1 if x0 < x1 else -1
    sy = 1 if y0 < y1 else -1
    err = dx - dy
    while True:
        for w in range(width):
            px(x0+w, y0, r, g, b, a)
            px(x0, y0+w, r, g, b, a)
        if x0 == x1 and y0 == y1:
            break
        e2 = 2 * err
        if e2 > -dy:
            err -= dy
            x0 += sx
        if e2 < dx:
            err += dx
            y0 += sy

# ── 魚身（橢圓，深紅漸層）────────────────────────────────────
# 主體：深紅色橢圓
for dy in range(-18, 19):
    for dx in range(-26, 27):
        if (dx/26)**2 + (dy/18)**2 <= 1.0:
            # 漸層：中心深紅，邊緣更暗
            dist = ((dx/26)**2 + (dy/18)**2) ** 0.5
            # 左上亮，右下暗
            shade = 1.0 - (dx + dy) / 60.0
            shade = max(0.5, min(1.2, shade))
            r2 = min(255, int(180 * shade))
            g2 = min(255, int(20 * shade))
            b2 = min(255, int(20 * shade))
            px(32+dx, 32+dy, r2, g2, b2)

# ── 崩潰裂縫紋路（Z字形閃電裂縫）────────────────────────────
# 主裂縫：從左上到右下的 Z 字形
crack_points = [
    (18, 20), (22, 24), (26, 22), (30, 26), (34, 24), (38, 28), (42, 26)
]
for i in range(len(crack_points)-1):
    x0, y0 = crack_points[i]
    x1, y1 = crack_points[i+1]
    draw_line(x0, y0, x1, y1, 255, 200, 50, 255, 2)  # 金色裂縫

# 次裂縫：分叉
draw_line(26, 22, 24, 18, 255, 180, 30, 200, 1)
draw_line(34, 24, 36, 20, 255, 180, 30, 200, 1)
draw_line(30, 26, 28, 30, 255, 180, 30, 200, 1)

# ── 崩潰光環（橙紅色）────────────────────────────────────────
for angle_deg in range(0, 360, 8):
    import math
    angle = math.radians(angle_deg)
    for r in range(20, 24):
        x = int(32 + r * math.cos(angle))
        y = int(32 + r * math.sin(angle))
        alpha = 180 if r == 21 else 100
        px(x, y, 255, 80, 20, alpha)

# ── 4方向爆炸光芒（火橙色）───────────────────────────────────
for i in range(8):
    angle = math.radians(i * 45)
    for r in range(22, 30):
        x = int(32 + r * math.cos(angle))
        y = int(32 + r * math.sin(angle))
        alpha = max(0, 220 - (r - 22) * 28)
        px(x, y, 255, 100 + (r-22)*10, 0, alpha)

# ── 倍率上升指示（右上角小數字符號）─────────────────────────
# 用像素繪製 "×" 符號
for d in range(-3, 4):
    px(48+d, 14+d, 255, 220, 50, 220)
    px(48+d, 14-d, 255, 220, 50, 220)

# ── 警告符號（左上角三角形）──────────────────────────────────
# 小三角形
for i in range(8):
    px(14+i, 14+8-i, 255, 50, 50, 200)
    px(14+i, 14+8-i-1, 255, 50, 50, 200)
px(18, 14, 255, 50, 50, 200)
px(18, 15, 255, 50, 50, 200)
px(18, 16, 255, 50, 50, 200)
px(18, 17, 255, 50, 50, 200)

# ── 魚眼（金色）──────────────────────────────────────────────
fill_circle(38, 28, 3, 255, 220, 50)
fill_circle(38, 28, 2, 200, 160, 20)
px(39, 27, 255, 255, 200)  # 高光

# ── 魚尾（深紅色三角形）──────────────────────────────────────
for i in range(10):
    for j in range(-i, i+1):
        px(6+i, 32+j, 150, 15, 15)

# ── 崩潰核心（中心金色圓點）──────────────────────────────────
fill_circle(32, 32, 4, 255, 200, 50)
fill_circle(32, 32, 2, 255, 240, 100)
px(32, 32, 255, 255, 200)

# ── 輸出 PNG ─────────────────────────────────────────────────
os.makedirs(os.path.dirname(OUT), exist_ok=True)
png_data = make_png_rgba(pixels, SIZE, SIZE)
with open(OUT, 'wb') as f:
    f.write(png_data)

# 統計非透明像素
non_transparent = sum(1 for row in pixels for p in row if p[3] > 10)
total = SIZE * SIZE
print(f"T130_crash_fish: {SIZE}x{SIZE}, non-transparent={non_transparent}/{total} ({non_transparent/total*100:.1f}%)")
print(f"Saved to: {OUT}")
