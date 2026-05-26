"""
generate_targets_day305.py — DAY-305 新增目標物精靈圖生成
T136 龍怒蓄積魚、T137 座頭鯨魚、T138 傳說龍魚、T139 公會戰魚、T140 品質魚
"""
import struct, zlib, os, math

OUT_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
SIZE = 64

def make_png(pixels):
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

# ── T136 龍怒蓄積魚 ──────────────────────────────────────────
def gen_t136():
    pixels = [(0,0,0,0)] * (SIZE * SIZE)
    # 魚身：火橙色橢圓
    fill_ellipse(pixels, 32, 32, 24, 13, 220, 80, 0)
    # 怒氣火焰光環
    for dy in range(-16, 17):
        for dx in range(-26, 27):
            d2 = dx*dx*1.2 + dy*dy*3
            if 22*22 < d2 <= 26*26:
                nx, ny = 32 + dx, 32 + dy
                if 0 <= nx < SIZE and 0 <= ny < SIZE:
                    pixels[ny * SIZE + nx] = (255, 100, 0, 120)
    # 怒氣計量條（魚身上方）
    for x in range(12, 52):
        for y in range(16, 20):
            pct = (x - 12) / 40.0
            r = int(255 * pct)
            g = int(100 * (1 - pct))
            pixels[y * SIZE + x] = (r, g, 0, 200)
    # 魚頭
    fill_circle(pixels, 50, 32, 12, 240, 100, 0, shaded=True)
    # 眼睛：火紅色
    fill_circle(pixels, 53, 29, 3, 255, 50, 0, shaded=False)
    fill_circle(pixels, 54, 29, 1, 255, 200, 0, shaded=False)
    # 4方向火焰光芒
    for i in range(4):
        angle = i * math.pi / 2
        for r in range(14, 22):
            x = int(32 + r * math.cos(angle))
            y = int(32 + r * math.sin(angle))
            if 0 <= x < SIZE and 0 <= y < SIZE:
                pixels[y * SIZE + x] = (255, 80, 0, 180)
    # 魚尾
    for dy in range(-10, 11):
        for dx in range(-8, 1):
            if abs(dy) <= abs(dx) + 2:
                nx, ny = 8 + dx, 32 + dy
                if 0 <= nx < SIZE and 0 <= ny < SIZE:
                    pixels[ny * SIZE + nx] = (180, 60, 0, 200)
    return pixels

# ── T137 座頭鯨魚 ──────────────────────────────────────────
def gen_t137():
    pixels = [(0,0,0,0)] * (SIZE * SIZE)
    # 鯨魚身體：深藍色大橢圓
    fill_ellipse(pixels, 30, 34, 26, 18, 0, 60, 150)
    # 鯨魚腹部：淡藍色
    fill_ellipse(pixels, 30, 40, 18, 10, 100, 160, 220)
    # 鯨魚頭部
    fill_circle(pixels, 50, 32, 14, 0, 80, 170, shaded=True)
    # 眼睛
    fill_circle(pixels, 54, 28, 3, 255, 255, 255, shaded=False)
    fill_circle(pixels, 55, 28, 1, 0, 0, 80, shaded=False)
    # 鯨歌聲波：同心圓
    for wave in range(3):
        r = 20 + wave * 8
        alpha = 120 - wave * 30
        for dy in range(-r-1, r+2):
            for dx in range(-r-1, r+2):
                d2 = dx*dx + dy*dy
                if r*r < d2 <= (r+2)*(r+2):
                    nx, ny = 32 + dx, 32 + dy
                    if 0 <= nx < SIZE and 0 <= ny < SIZE:
                        pixels[ny * SIZE + nx] = (0, 180, 255, alpha)
    # 胸鰭
    fill_ellipse(pixels, 14, 38, 10, 6, 0, 50, 130)
    # 尾鰭
    for dy in range(-12, 13):
        for dx in range(-10, 1):
            if abs(dy) <= abs(dx) + 3:
                nx, ny = 6 + dx, 34 + dy
                if 0 <= nx < SIZE and 0 <= ny < SIZE:
                    pixels[ny * SIZE + nx] = (0, 40, 120, 200)
    return pixels

# ── T138 傳說龍魚 ──────────────────────────────────────────
def gen_t138():
    pixels = [(0,0,0,0)] * (SIZE * SIZE)
    # 龍魚身體：金橙色橢圓
    fill_ellipse(pixels, 32, 32, 24, 14, 200, 120, 0)
    # 龍鱗：深橙色格紋
    for i in range(5):
        for j in range(3):
            cx = 14 + i * 9
            cy = 24 + j * 8
            fill_circle(pixels, cx, cy, 3, 160, 80, 0, shaded=False)
    # 龍頭：圓形
    fill_circle(pixels, 50, 32, 13, 220, 140, 0, shaded=True)
    # 龍眼：金色
    fill_circle(pixels, 53, 28, 4, 255, 200, 0, shaded=False)
    fill_circle(pixels, 54, 28, 2, 200, 100, 0, shaded=False)
    # 龍角：兩個尖角
    draw_line(pixels, 48, 20, 44, 12, 255, 180, 0, 220, 2)
    draw_line(pixels, 52, 20, 56, 12, 255, 180, 0, 220, 2)
    # 火焰光環
    for dy in range(-16, 17):
        for dx in range(-26, 27):
            d2 = dx*dx*1.2 + dy*dy*3
            if 22*22 < d2 <= 26*26:
                nx, ny = 32 + dx, 32 + dy
                if 0 <= nx < SIZE and 0 <= ny < SIZE:
                    pixels[ny * SIZE + nx] = (255, 150, 0, 100)
    # 龍尾
    for dy in range(-12, 13):
        for dx in range(-10, 1):
            if abs(dy) <= abs(dx) + 2:
                nx, ny = 8 + dx, 32 + dy
                if 0 <= nx < SIZE and 0 <= ny < SIZE:
                    pixels[ny * SIZE + nx] = (180, 90, 0, 200)
    return pixels

# ── T139 公會戰魚 ──────────────────────────────────────────
def gen_t139():
    pixels = [(0,0,0,0)] * (SIZE * SIZE)
    # 魚身：金色橢圓
    fill_ellipse(pixels, 32, 32, 24, 13, 200, 160, 0)
    # 公會徽章：中心金色盾牌
    for dy in range(-10, 11):
        for dx in range(-8, 9):
            if abs(dx) + abs(dy) <= 14:
                nx, ny = 32 + dx, 32 + dy
                if 0 <= nx < SIZE and 0 <= ny < SIZE:
                    pixels[ny * SIZE + nx] = (255, 200, 0, 200)
    # 盾牌輪廓
    for dy in range(-10, 11):
        for dx in range(-8, 9):
            if abs(abs(dx) + abs(dy) - 14) <= 1:
                nx, ny = 32 + dx, 32 + dy
                if 0 <= nx < SIZE and 0 <= ny < SIZE:
                    pixels[ny * SIZE + nx] = (200, 100, 0, 255)
    # 劍：白色十字
    draw_line(pixels, 28, 32, 36, 32, 255, 255, 255, 220, 2)
    draw_line(pixels, 32, 26, 32, 38, 255, 255, 255, 220, 2)
    # 魚頭
    fill_circle(pixels, 50, 32, 12, 220, 170, 0, shaded=True)
    # 眼睛
    fill_circle(pixels, 53, 29, 3, 255, 255, 255, shaded=False)
    fill_circle(pixels, 54, 29, 1, 100, 50, 0, shaded=False)
    # 魚尾
    for dy in range(-10, 11):
        for dx in range(-8, 1):
            if abs(dy) <= abs(dx) + 2:
                nx, ny = 8 + dx, 32 + dy
                if 0 <= nx < SIZE and 0 <= ny < SIZE:
                    pixels[ny * SIZE + nx] = (180, 130, 0, 200)
    return pixels

# ── T140 品質魚 ──────────────────────────────────────────
def gen_t140():
    pixels = [(0,0,0,0)] * (SIZE * SIZE)
    # 魚身：彩虹漸層橢圓
    for dy in range(-14, 15):
        for dx in range(-24, 25):
            if (dx/24)**2 + (dy/14)**2 <= 1.0:
                nx, ny = 32 + dx, 32 + dy
                if 0 <= nx < SIZE and 0 <= ny < SIZE:
                    # 彩虹漸層：從左到右
                    hue = (dx + 24) / 48.0
                    r = int(255 * abs(math.sin(hue * math.pi)))
                    g = int(255 * abs(math.sin((hue + 0.33) * math.pi)))
                    b = int(255 * abs(math.sin((hue + 0.67) * math.pi)))
                    pixels[ny * SIZE + nx] = (r, g, b, 220)
    # 品質星星：5顆金色星星
    for i in range(5):
        angle = i * 2 * math.pi / 5 - math.pi / 2
        x = int(32 + 10 * math.cos(angle))
        y = int(32 + 10 * math.sin(angle))
        fill_circle(pixels, x, y, 3, 255, 215, 0, shaded=False)
    # 中心鑽石
    fill_circle(pixels, 32, 32, 5, 255, 255, 255, shaded=False)
    fill_circle(pixels, 32, 32, 3, 200, 200, 255, shaded=False)
    # 魚頭
    fill_circle(pixels, 50, 32, 12, 180, 100, 200, shaded=True)
    # 眼睛
    fill_circle(pixels, 53, 29, 3, 255, 255, 255, shaded=False)
    fill_circle(pixels, 54, 29, 1, 100, 0, 150, shaded=False)
    # 魚尾
    for dy in range(-10, 11):
        for dx in range(-8, 1):
            if abs(dy) <= abs(dx) + 2:
                nx, ny = 8 + dx, 32 + dy
                if 0 <= nx < SIZE and 0 <= ny < SIZE:
                    hue = (dy + 10) / 20.0
                    r = int(200 * abs(math.sin(hue * math.pi)))
                    g = int(100 * abs(math.sin((hue + 0.33) * math.pi)))
                    b = int(200 * abs(math.sin((hue + 0.67) * math.pi)))
                    pixels[ny * SIZE + nx] = (r, g, b, 200)
    return pixels

# ── 生成並儲存 ────────────────────────────────────────────────
targets = [
    ("T136_dragon_wrath_v2", gen_t136()),
    ("T137_humpback_whale", gen_t137()),
    ("T138_legend_dragon", gen_t138()),
    ("T139_guild_war", gen_t139()),
    ("T140_quality_fish", gen_t140()),
]

for name, pixels in targets:
    path = os.path.join(OUT_DIR, f"{name}.png")
    with open(path, "wb") as f:
        f.write(make_png(pixels))
    non_transparent = sum(1 for p in pixels if p[3] > 0)
    pct = non_transparent / (SIZE * SIZE) * 100
    print(f"✅ {name}.png — {non_transparent} 非透明像素 ({pct:.1f}%)")

print("\n✅ DAY-305 目標物精靈圖生成完成！")
