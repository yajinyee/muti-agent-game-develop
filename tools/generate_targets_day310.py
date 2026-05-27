"""
generate_targets_day310.py — DAY-310 T161-T165 精靈圖生成
T161 幸運連擊爆發魚（火橙色 + 連擊火焰 + 爆炸光芒）
T162 幸運時間炸彈魚（深紅色 + 炸彈外殼 + 倒數符號）
T163 幸運元素融合魚（深紫色 + 三元素符號 + 融合光環）
T164 幸運寶藏獵人魚（金色 + 寶石 + 獵人標記）
T165 幸運神話覺醒魚（神聖金色 + 12方向光芒 + 神話符號）
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
    return [[( 0, 0, 0, 0)] * SIZE for _ in range(SIZE)]

def dist(x1, y1, x2, y2):
    return math.sqrt((x1-x2)**2 + (y1-y2)**2)

def fill_circle(pixels, cx, cy, r, color, alpha=255):
    for y in range(SIZE):
        for x in range(SIZE):
            if dist(x, y, cx, cy) <= r:
                pixels[y][x] = (*color, alpha)

def fill_circle_shaded(pixels, cx, cy, r, base_color):
    light = tuple(min(255, int(c * 1.4)) for c in base_color)
    dark = tuple(max(0, int(c * 0.6)) for c in base_color)
    for y in range(SIZE):
        for x in range(SIZE):
            d = dist(x, y, cx, cy)
            if d <= r:
                if x < cx and y < cy:
                    col = light
                elif x > cx and y > cy:
                    col = dark
                else:
                    col = base_color
                pixels[y][x] = (*col, 255)

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

def draw_star(pixels, cx, cy, r_outer, r_inner, points, color, alpha=230):
    for i in range(points * 2):
        angle = math.radians(i * 180 / points - 90)
        r = r_outer if i % 2 == 0 else r_inner
        x = cx + r * math.cos(angle)
        y = cy + r * math.sin(angle)
        # 填充三角形（簡化版）
        for ty in range(SIZE):
            for tx in range(SIZE):
                d = dist(tx, ty, cx, cy)
                if d <= r_outer:
                    a = math.atan2(ty - cy, tx - cx)
                    # 判斷是否在星形內
                    sector = (a + math.pi) / (2 * math.pi) * points
                    frac = sector - int(sector)
                    r_at_angle = r_inner + (r_outer - r_inner) * abs(2 * frac - 1)
                    if d <= r_at_angle:
                        pixels[ty][tx] = (*color, alpha)

# ── T161 幸運連擊爆發魚（火橙色 + 連擊火焰 + 爆炸光芒）────────────────
def gen_t161():
    p = empty()
    cx, cy = 32, 32
    # 主體：火橙橢圓魚身
    for y in range(SIZE):
        for x in range(SIZE):
            dx = (x - cx) / 20.0
            dy = (y - cy) / 14.0
            if dx*dx + dy*dy <= 1.0:
                # 火焰漸層
                t = (y - (cy - 14)) / 28.0
                r = int(255)
                g = int(120 - t * 80)
                b = int(20)
                p[y][x] = (r, g, b, 255)
    # 連擊火焰光芒（8方向）
    for angle in range(0, 360, 45):
        draw_ray(p, cx, cy, angle, 14, 1, (255, 200, 50), 200)
    # 中心爆炸圓
    fill_circle(p, cx, cy, 8, (255, 240, 100), 230)
    # 眼睛
    p[cy-3][cx-5] = (255, 255, 255, 255)
    p[cy-3][cx-4] = (255, 255, 255, 255)
    p[cy-3][cx+4] = (255, 255, 255, 255)
    p[cy-3][cx+5] = (255, 255, 255, 255)
    # 連擊符號（×）
    for i in range(-3, 4):
        if 0 <= cx+i < SIZE and 0 <= cy+i < SIZE:
            p[cy+i][cx+i] = (255, 255, 255, 220)
        if 0 <= cx+i < SIZE and 0 <= cy-i < SIZE:
            p[cy-i][cx+i] = (255, 255, 255, 220)
    return p

# ── T162 幸運時間炸彈魚（深紅色 + 炸彈外殼 + 倒數符號）────────────────
def gen_t162():
    p = empty()
    cx, cy = 32, 32
    # 主體：深紅圓形炸彈身
    fill_circle_shaded(p, cx, cy, 18, (180, 20, 20))
    # 炸彈導線（右上角）
    for i in range(8):
        nx = cx + 12 + i
        ny = cy - 12 - i
        if 0 <= nx < SIZE and 0 <= ny < SIZE:
            p[ny][nx] = (200, 150, 50, 255)
    # 炸彈頂部圓
    fill_circle(p, cx+18, cy-18, 4, (200, 150, 50), 255)
    # 倒數符號（!）
    for i in range(-5, 2):
        if 0 <= cy+i < SIZE:
            p[cy+i][cx] = (255, 255, 100, 230)
    p[cy+4][cx] = (255, 255, 100, 230)
    # 裂縫紋路
    for i in range(-8, 9):
        if 0 <= cx+i < SIZE and 0 <= cy+i//2 < SIZE:
            p[cy+i//2][cx+i] = (100, 0, 0, 200)
    # 光環
    draw_ring(p, cx, cy, 19, 22, (255, 100, 50), 150)
    return p

# ── T163 幸運元素融合魚（深紫色 + 三元素符號 + 融合光環）────────────────
def gen_t163():
    p = empty()
    cx, cy = 32, 32
    # 主體：深紫橢圓魚身
    for y in range(SIZE):
        for x in range(SIZE):
            dx = (x - cx) / 19.0
            dy = (y - cy) / 13.0
            if dx*dx + dy*dy <= 1.0:
                t = (y - (cy - 13)) / 26.0
                r = int(100 + t * 40)
                g = int(20 + t * 20)
                b = int(200 - t * 50)
                p[y][x] = (r, g, b, 255)
    # 三元素符號（火/冰/雷）
    # 火（上方）
    for i in range(-3, 4):
        if 0 <= cy-8+i < SIZE and 0 <= cx+i//2 < SIZE:
            p[cy-8+i][cx+i//2] = (255, 100, 20, 220)
    # 冰（左下）
    for i in range(-3, 4):
        if 0 <= cy+5 < SIZE and 0 <= cx-8+i < SIZE:
            p[cy+5][cx-8+i] = (100, 200, 255, 220)
    # 雷（右下）
    for i in range(-3, 4):
        if 0 <= cy+5 < SIZE and 0 <= cx+5+i < SIZE:
            p[cy+5][cx+5+i] = (255, 255, 50, 220)
    # 融合光環
    draw_ring(p, cx, cy, 20, 24, (180, 100, 255), 180)
    # 中心融合點
    fill_circle(p, cx, cy, 5, (220, 180, 255), 230)
    return p

# ── T164 幸運寶藏獵人魚（金色 + 寶石 + 獵人標記）────────────────
def gen_t164():
    p = empty()
    cx, cy = 32, 32
    # 主體：金色橢圓魚身
    for y in range(SIZE):
        for x in range(SIZE):
            dx = (x - cx) / 20.0
            dy = (y - cy) / 14.0
            if dx*dx + dy*dy <= 1.0:
                t = (y - (cy - 14)) / 28.0
                r = int(220 + t * 35)
                g = int(180 - t * 30)
                b = int(20 + t * 10)
                p[y][x] = (min(255,r), g, b, 255)
    # 5個寶石標記（圍繞中心）
    for i in range(5):
        angle = math.radians(i * 72 - 90)
        gx = int(cx + 10 * math.cos(angle))
        gy = int(cy + 10 * math.sin(angle))
        fill_circle(p, gx, gy, 3, (100, 200, 255), 230)
        # 寶石高光
        if 0 <= gx-1 < SIZE and 0 <= gy-1 < SIZE:
            p[gy-1][gx-1] = (255, 255, 255, 200)
    # 獵人瞄準線（十字）
    for i in range(-12, 13):
        if 0 <= cx+i < SIZE:
            p[cy][cx+i] = (255, 220, 50, 150)
        if 0 <= cy+i < SIZE:
            p[cy+i][cx] = (255, 220, 50, 150)
    # 金色光環
    draw_ring(p, cx, cy, 21, 25, (255, 200, 50), 180)
    return p

# ── T165 幸運神話覺醒魚（神聖金色 + 12方向光芒 + 神話符號）────────────────
def gen_t165():
    p = empty()
    cx, cy = 32, 32
    # 主體：神聖金色大型圓形魚身
    fill_circle_shaded(p, cx, cy, 16, (220, 190, 50))
    # 12方向神聖光芒
    for angle in range(0, 360, 30):
        draw_ray(p, cx, cy, angle, 16, 1, (255, 240, 150), 200)
    # 外圈光環
    draw_ring(p, cx, cy, 17, 21, (255, 220, 100), 180)
    draw_ring(p, cx, cy, 22, 25, (255, 200, 50), 120)
    # 神話符號（星形）
    draw_star(p, cx, cy, 10, 5, 6, (255, 255, 200), 200)
    # 眼睛（紅色龍眼）
    p[cy-3][cx-5] = (255, 50, 50, 255)
    p[cy-3][cx-4] = (255, 50, 50, 255)
    p[cy-3][cx+4] = (255, 50, 50, 255)
    p[cy-3][cx+5] = (255, 50, 50, 255)
    # 中心神聖光點
    fill_circle(p, cx, cy, 4, (255, 255, 220), 240)
    return p

targets = [
    ("T161_combo_burst", gen_t161),
    ("T162_time_bomb", gen_t162),
    ("T163_elemental_fusion", gen_t163),
    ("T164_treasure_hunter", gen_t164),
    ("T165_myth_awaken", gen_t165),
]

os.makedirs(OUT_DIR, exist_ok=True)
for name, gen_fn in targets:
    pixels = gen_fn()
    data = make_png(pixels)
    path = os.path.join(OUT_DIR, f"{name}.png")
    with open(path, "wb") as f:
        f.write(data)
    # 計算非透明像素
    non_transparent = sum(1 for row in pixels for r,g,b,a in row if a > 0)
    pct = non_transparent / (SIZE * SIZE) * 100
    print(f"✅ {name}.png ({non_transparent} 非透明像素, {pct:.1f}%)")

print("\n🎉 DAY-310 T161-T165 精靈圖生成完成！")
