"""
generate_targets_day311_t161_fix.py — DAY-311 T161 精靈圖密度修復
問題：T161 連擊爆發魚非透明像素只有 21.5%，視覺稀疏
目標：提升到 35%+ 非透明像素，增加火焰填充密度
"""
import struct, zlib, os, math, random

OUT_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
SIZE = 64

def make_png(pixels):
    def chunk(name, data):
        c = zlib.crc32(name + data) & 0xFFFFFFFF
        return struct.pack(">I", len(data)) + name + data + struct.pack(">I", c)
    raw = b""
    for row in pixels:
        raw += b"\x00"
        for r, g, b, a in row:
            raw += bytes([r, g, b, a])
    compressed = zlib.compress(raw, 9)
    return (b"\x89PNG\r\n\x1a\n"
            + chunk(b"IHDR", struct.pack(">IIBBBBB", SIZE, SIZE, 8, 2, 0, 0, 0))
            + chunk(b"IDAT", compressed)
            + chunk(b"IEND", b""))

def empty():
    return [[(0, 0, 0, 0)] * SIZE for _ in range(SIZE)]

def dist(x1, y1, x2, y2):
    return math.sqrt((x1-x2)**2 + (y1-y2)**2)

def fill_circle(pixels, cx, cy, r, color, alpha=255):
    for y in range(SIZE):
        for x in range(SIZE):
            if dist(x, y, cx, cy) <= r:
                pixels[y][x] = (*color, alpha)

def draw_ring(pixels, cx, cy, r_inner, r_outer, color, alpha=200):
    for y in range(SIZE):
        for x in range(SIZE):
            d = dist(x, y, cx, cy)
            if r_inner <= d <= r_outer:
                pixels[y][x] = (*color, alpha)

def draw_ray(pixels, cx, cy, angle_deg, length, width, color, alpha=220):
    angle = math.radians(angle_deg)
    for t in range(int(length)):
        px = int(cx + t * math.cos(angle))
        py = int(cy + t * math.sin(angle))
        for dx in range(-width, width+1):
            for dy in range(-width, width+1):
                nx, ny = px+dx, py+dy
                if 0 <= nx < SIZE and 0 <= ny < SIZE:
                    pixels[ny][nx] = (*color, alpha)

# ── T161 修復版：火橙色 + 密集火焰填充 + 爆炸光芒 ──────────────────────
def gen_t161_fixed():
    p = empty()
    cx, cy = 32, 32

    # 1. 外層火焰光暈（大橢圓，低透明度）
    for y in range(SIZE):
        for x in range(SIZE):
            dx = (x - cx) / 28.0
            dy = (y - cy) / 22.0
            d = dx*dx + dy*dy
            if d <= 1.0:
                # 火焰漸層：中心亮黃，外圍深橙
                t = math.sqrt(d)
                r = 255
                g = int(200 - t * 160)
                b = int(30 - t * 30)
                alpha = int(180 - t * 120)
                if alpha > 0:
                    p[y][x] = (r, max(0,g), max(0,b), alpha)

    # 2. 主體：火橙橢圓魚身（密實填充）
    for y in range(SIZE):
        for x in range(SIZE):
            dx = (x - cx) / 20.0
            dy = (y - cy) / 14.0
            if dx*dx + dy*dy <= 1.0:
                t = (y - (cy - 14)) / 28.0
                r = 255
                g = int(140 - t * 100)
                b = int(20)
                p[y][x] = (r, max(0,g), b, 255)

    # 3. 8方向火焰光芒（加粗，長度增加）
    for angle in range(0, 360, 45):
        draw_ray(p, cx, cy, angle, 18, 2, (255, 200, 50), 220)

    # 4. 45度偏移的次要光芒（填補空隙）
    for angle in range(22, 360, 45):
        draw_ray(p, cx, cy, angle, 12, 1, (255, 160, 30), 180)

    # 5. 中心爆炸圓（加大）
    fill_circle(p, cx, cy, 10, (255, 240, 100), 240)
    fill_circle(p, cx, cy, 6, (255, 255, 200), 255)

    # 6. 火焰粒子（隨機散佈在魚身周圍）
    random.seed(161)
    for _ in range(20):
        angle = random.uniform(0, math.pi * 2)
        r_dist = random.uniform(14, 22)
        fx = int(cx + r_dist * math.cos(angle))
        fy = int(cy + r_dist * math.sin(angle))
        if 0 <= fx < SIZE and 0 <= fy < SIZE:
            fill_circle(p, fx, fy, random.randint(2, 4), (255, random.randint(100, 200), 20), 200)

    # 7. 眼睛（白色高光）
    for ex, ey in [(cx-5, cy-3), (cx-4, cy-3), (cx+4, cy-3), (cx+5, cy-3)]:
        if 0 <= ex < SIZE and 0 <= ey < SIZE:
            p[ey][ex] = (255, 255, 255, 255)

    # 8. 連擊符號（×，加粗）
    for i in range(-4, 5):
        for w in range(-1, 2):
            if 0 <= cx+i < SIZE and 0 <= cy+i+w < SIZE:
                p[cy+i+w][cx+i] = (255, 255, 255, 200)
            if 0 <= cx+i < SIZE and 0 <= cy-i+w < SIZE:
                p[cy-i+w][cx+i] = (255, 255, 255, 200)

    # 9. 外圈火焰環
    draw_ring(p, cx, cy, 22, 26, (255, 120, 20), 160)
    draw_ring(p, cx, cy, 26, 29, (200, 80, 10), 100)

    return p

if __name__ == "__main__":
    os.makedirs(OUT_DIR, exist_ok=True)
    pixels = gen_t161_fixed()
    data = make_png(pixels)
    path = os.path.join(OUT_DIR, "T161_combo_burst.png")
    with open(path, "wb") as f:
        f.write(data)
    non_transparent = sum(1 for row in pixels for r,g,b,a in row if a > 0)
    pct = non_transparent / (SIZE * SIZE) * 100
    print(f"✅ T161_combo_burst.png 修復完成")
    print(f"   非透明像素：{non_transparent} ({pct:.1f}%)")
    if pct >= 35:
        print(f"   ✅ 密度達標（目標 35%+）")
    else:
        print(f"   ⚠️ 密度不足，需要進一步調整")
