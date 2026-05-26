"""
generate_targets_day304.py — DAY-304 新增目標物精靈圖生成
T131 電鰻魚、T132 巨型安康魚、T133 黑洞魚、T134 賞金獵人魚、T135 海嘯魚
"""
import struct, zlib, os, math, random

OUT_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
SIZE = 64

def make_png(pixels):
    """pixels: list of (r,g,b,a) tuples, row by row"""
    raw = b""
    for row in range(SIZE):
        raw += b"\x00"
        for col in range(SIZE):
            r, g, b, a = pixels[row * SIZE + col]
            raw += bytes([r, g, b, a])
    compressed = zlib.compress(raw, 9)
    def chunk(name, data):
        c = name + data
        return struct.pack(">I", len(data)) + c + struct.pack(">I", zlib.crc32(c) & 0xFFFFFFFF)
    ihdr = struct.pack(">IIBBBBB", SIZE, SIZE, 8, 2, 0, 0, 0)
    # RGBA
    ihdr = struct.pack(">II", SIZE, SIZE) + bytes([8, 6, 0, 0, 0])
    png = b"\x89PNG\r\n\x1a\n"
    png += chunk(b"IHDR", ihdr)
    png += chunk(b"IDAT", compressed)
    png += chunk(b"IEND", b"")
    return png

def px(pixels, x, y, r, g, b, a=255):
    if 0 <= x < SIZE and 0 <= y < SIZE:
        pixels[y * SIZE + x] = (r, g, b, a)

def fill_circle(pixels, cx, cy, radius, r, g, b, a=255, shaded=True):
    for dy in range(-radius, radius + 1):
        for dx in range(-radius, radius + 1):
            if dx*dx + dy*dy <= radius*radius:
                nx, ny = cx + dx, cy + dy
                if 0 <= nx < SIZE and 0 <= ny < SIZE:
                    if shaded:
                        # 光源左上
                        lx, ly = -0.707, -0.707
                        dist = math.sqrt(dx*dx + dy*dy)
                        nx2 = dx / max(dist, 0.001)
                        ny2 = dy / max(dist, 0.001)
                        dot = max(0, -(nx2*lx + ny2*ly))
                        light = 0.6 + 0.4 * dot
                        nr = min(255, int(r * light))
                        ng = min(255, int(g * light))
                        nb = min(255, int(b * light))
                        pixels[ny * SIZE + nx] = (nr, ng, nb, a)
                    else:
                        pixels[ny * SIZE + nx] = (r, g, b, a)

def fill_ellipse(pixels, cx, cy, rx, ry, r, g, b, a=255):
    for dy in range(-ry, ry + 1):
        for dx in range(-rx, rx + 1):
            if (dx/rx)**2 + (dy/ry)**2 <= 1.0:
                nx, ny = cx + dx, cy + dy
                if 0 <= nx < SIZE and 0 <= ny < SIZE:
                    pixels[ny * SIZE + nx] = (r, g, b, a)

def draw_line(pixels, x0, y0, x1, y1, r, g, b, a=255, width=1):
    dx = abs(x1 - x0)
    dy = abs(y1 - y0)
    sx = 1 if x0 < x1 else -1
    sy = 1 if y0 < y1 else -1
    err = dx - dy
    while True:
        for w in range(-width//2, width//2 + 1):
            px(pixels, x0 + w, y0, r, g, b, a)
            px(pixels, x0, y0 + w, r, g, b, a)
        if x0 == x1 and y0 == y1:
            break
        e2 = 2 * err
        if e2 > -dy:
            err -= dy
            x0 += sx
        if e2 < dx:
            err += dx
            y0 += sy

# ── T131 電鰻魚 ──────────────────────────────────────────────
def gen_t131():
    pixels = [(0,0,0,0)] * (SIZE * SIZE)
    # 電鰻身體：細長橢圓，電黃色
    fill_ellipse(pixels, 32, 32, 26, 12, 255, 240, 0)
    # 電鰻條紋：深黃色橫條
    for i in range(5):
        x = 10 + i * 9
        for dy in range(-8, 9):
            if abs(dy) < 8:
                px(pixels, x, 32 + dy, 200, 180, 0, 200)
    # 電鰻頭部：圓形
    fill_circle(pixels, 52, 32, 10, 255, 220, 0, shaded=True)
    # 眼睛：白色 + 黑色瞳孔
    fill_circle(pixels, 55, 29, 3, 255, 255, 255, shaded=False)
    fill_circle(pixels, 56, 29, 1, 0, 0, 0, shaded=False)
    # 電弧：青藍色閃電
    for i in range(4):
        x = 15 + i * 10
        draw_line(pixels, x, 24, x+3, 28, 0, 220, 255, 200, 1)
        draw_line(pixels, x+3, 28, x+1, 32, 0, 220, 255, 200, 1)
        draw_line(pixels, x+1, 32, x+4, 36, 0, 220, 255, 200, 1)
    # 電光暈：外圈淡黃
    for dy in range(-14, 15):
        for dx in range(-28, 29):
            if 26*26 < dx*dx*1.5 + dy*dy*4 <= 30*30:
                nx, ny = 32 + dx, 32 + dy
                if 0 <= nx < SIZE and 0 <= ny < SIZE:
                    pixels[ny * SIZE + nx] = (255, 255, 100, 60)
    return pixels

# ── T132 巨型安康魚 ──────────────────────────────────────────
def gen_t132():
    pixels = [(0,0,0,0)] * (SIZE * SIZE)
    # 安康魚身體：深海青藍色，圓形
    fill_circle(pixels, 30, 34, 20, 0, 80, 120, shaded=True)
    # 大嘴：深色橢圓
    fill_ellipse(pixels, 30, 42, 16, 6, 0, 30, 50)
    # 牙齒：白色小三角
    for i in range(5):
        x = 18 + i * 6
        px(pixels, x, 40, 255, 255, 255)
        px(pixels, x, 41, 255, 255, 255)
        px(pixels, x+1, 42, 255, 255, 255)
    # 眼睛：大圓眼，黃色
    fill_circle(pixels, 24, 28, 6, 255, 220, 0, shaded=False)
    fill_circle(pixels, 24, 28, 3, 0, 0, 0, shaded=False)
    fill_circle(pixels, 36, 28, 6, 255, 220, 0, shaded=False)
    fill_circle(pixels, 36, 28, 3, 0, 0, 0, shaded=False)
    # 誘餌燈：頭頂發光球
    draw_line(pixels, 30, 14, 30, 8, 0, 180, 200, 200, 1)
    fill_circle(pixels, 30, 6, 5, 0, 255, 220, shaded=False)
    # 誘餌光暈
    for dy in range(-8, 9):
        for dx in range(-8, 9):
            if 5*5 < dx*dx + dy*dy <= 8*8:
                nx, ny = 30 + dx, 6 + dy
                if 0 <= nx < SIZE and 0 <= ny < SIZE:
                    pixels[ny * SIZE + nx] = (0, 255, 200, 80)
    # 胸鰭
    fill_ellipse(pixels, 12, 36, 8, 5, 0, 60, 100)
    fill_ellipse(pixels, 48, 36, 8, 5, 0, 60, 100)
    return pixels

# ── T133 黑洞魚 ──────────────────────────────────────────────
def gen_t133():
    pixels = [(0,0,0,0)] * (SIZE * SIZE)
    # 黑洞核心：深紫黑色圓形
    fill_circle(pixels, 32, 32, 22, 10, 0, 20, shaded=True)
    # 事件視界：深紫色環
    for dy in range(-24, 25):
        for dx in range(-24, 25):
            d2 = dx*dx + dy*dy
            if 22*22 < d2 <= 24*24:
                nx, ny = 32 + dx, 32 + dy
                if 0 <= nx < SIZE and 0 <= ny < SIZE:
                    pixels[ny * SIZE + nx] = (80, 0, 120, 220)
    # 吸積盤：橙色環
    for dy in range(-28, 29):
        for dx in range(-28, 29):
            d2 = dx*dx + dy*dy
            if 24*24 < d2 <= 28*28:
                nx, ny = 32 + dx, 32 + dy
                if 0 <= nx < SIZE and 0 <= ny < SIZE:
                    angle = math.atan2(dy, dx)
                    brightness = int(128 + 127 * math.sin(angle * 3))
                    pixels[ny * SIZE + nx] = (brightness, brightness//3, 0, 180)
    # 奇點：中心白色亮點
    fill_circle(pixels, 32, 32, 4, 255, 255, 255, shaded=False)
    fill_circle(pixels, 32, 32, 2, 200, 150, 255, shaded=False)
    # 引力線：4方向螺旋
    for i in range(4):
        angle = i * math.pi / 2
        for r in range(5, 22):
            spiral_angle = angle + r * 0.15
            x = int(32 + r * math.cos(spiral_angle))
            y = int(32 + r * math.sin(spiral_angle))
            if 0 <= x < SIZE and 0 <= y < SIZE:
                pixels[y * SIZE + x] = (150, 50, 200, 160)
    return pixels

# ── T134 賞金獵人魚 ──────────────────────────────────────────
def gen_t134():
    pixels = [(0,0,0,0)] * (SIZE * SIZE)
    # 魚身：火橙色橢圓
    fill_ellipse(pixels, 32, 32, 24, 14, 220, 100, 0)
    # 魚鱗：深橙色格紋
    for i in range(4):
        for j in range(3):
            cx = 16 + i * 10
            cy = 24 + j * 8
            fill_circle(pixels, cx, cy, 4, 180, 70, 0, shaded=False)
    # 魚頭：圓形
    fill_circle(pixels, 50, 32, 12, 240, 120, 0, shaded=True)
    # 眼睛：紅色（賞金獵人的眼神）
    fill_circle(pixels, 53, 29, 3, 255, 50, 0, shaded=False)
    fill_circle(pixels, 54, 29, 1, 255, 255, 0, shaded=False)
    # 賞金標記：金色星星
    for angle_deg in [0, 72, 144, 216, 288]:
        angle = math.radians(angle_deg)
        x = int(32 + 8 * math.cos(angle))
        y = int(32 + 8 * math.sin(angle))
        fill_circle(pixels, x, y, 2, 255, 215, 0, shaded=False)
    # 中心金色圓
    fill_circle(pixels, 32, 32, 4, 255, 215, 0, shaded=False)
    # 魚尾：橙色三角
    for dy in range(-10, 11):
        for dx in range(-8, 1):
            if abs(dy) <= abs(dx) + 2:
                nx, ny = 8 + dx, 32 + dy
                if 0 <= nx < SIZE and 0 <= ny < SIZE:
                    pixels[ny * SIZE + nx] = (200, 80, 0, 200)
    # 瞄準線：白色十字
    draw_line(pixels, 28, 32, 36, 32, 255, 255, 255, 180, 1)
    draw_line(pixels, 32, 28, 32, 36, 255, 255, 255, 180, 1)
    return pixels

# ── T135 海嘯魚 ──────────────────────────────────────────────
def gen_t135():
    pixels = [(0,0,0,0)] * (SIZE * SIZE)
    # 魚身：深藍色橢圓
    fill_ellipse(pixels, 32, 32, 26, 14, 0, 80, 180)
    # 波浪紋路：青藍色弧線
    for i in range(3):
        y_base = 24 + i * 8
        for x in range(8, 56):
            y = y_base + int(3 * math.sin((x - 8) * 0.4 + i * 2))
            if 0 <= y < SIZE:
                pixels[y * SIZE + x] = (0, 180, 255, 180)
    # 魚頭：圓形
    fill_circle(pixels, 50, 32, 12, 0, 100, 200, shaded=True)
    # 眼睛：白色 + 深藍瞳孔
    fill_circle(pixels, 53, 29, 3, 255, 255, 255, shaded=False)
    fill_circle(pixels, 54, 29, 1, 0, 0, 100, shaded=False)
    # 三波浪光環
    for wave in range(3):
        r = 16 + wave * 6
        alpha = 150 - wave * 40
        for dy in range(-r-1, r+2):
            for dx in range(-r-1, r+2):
                d2 = dx*dx + dy*dy
                if r*r < d2 <= (r+2)*(r+2):
                    nx, ny = 32 + dx, 32 + dy
                    if 0 <= nx < SIZE and 0 <= ny < SIZE:
                        pixels[ny * SIZE + nx] = (0, 150 + wave*30, 255, alpha)
    # 魚尾：藍色三角
    for dy in range(-12, 13):
        for dx in range(-10, 1):
            if abs(dy) <= abs(dx) + 2:
                nx, ny = 6 + dx, 32 + dy
                if 0 <= nx < SIZE and 0 <= ny < SIZE:
                    pixels[ny * SIZE + nx] = (0, 60, 160, 200)
    return pixels

# ── 生成並儲存 ────────────────────────────────────────────────
targets = [
    ("T131_electric_eel", gen_t131()),
    ("T132_angler_fish", gen_t132()),
    ("T133_black_hole", gen_t133()),
    ("T134_bounty_hunter", gen_t134()),
    ("T135_tsunami", gen_t135()),
]

for name, pixels in targets:
    path = os.path.join(OUT_DIR, f"{name}.png")
    with open(path, "wb") as f:
        f.write(make_png(pixels))
    non_transparent = sum(1 for p in pixels if p[3] > 0)
    pct = non_transparent / (SIZE * SIZE) * 100
    print(f"✅ {name}.png — {non_transparent} 非透明像素 ({pct:.1f}%)")

print("\n✅ DAY-304 目標物精靈圖生成完成！")
