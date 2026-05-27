"""
generate_targets_day309.py — DAY-309 T156-T160 精靈圖生成
T156 冰鳳凰魚、T157 龍怒能量魚、T158 倍率瀑布魚、T159 覺醒 BOSS v2 魚、T160 終極審判魚
"""
import os
import struct
import zlib
import math

OUT_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
SIZE = 64

def make_png(pixels, w, h):
    def chunk(name, data):
        c = zlib.crc32(name + data) & 0xFFFFFFFF
        return struct.pack('>I', len(data)) + name + data + struct.pack('>I', c)
    sig = b'\x89PNG\r\n\x1a\n'
    ihdr = chunk(b'IHDR', struct.pack('>IIBBBBB', w, h, 8, 2, 0, 0, 0))
    raw = b''
    for y in range(h):
        raw += b'\x00'
        for x in range(w):
            r, g, b, a = pixels[y][x]
            raw += bytes([r, g, b, a])
    idat = chunk(b'IDAT', zlib.compress(raw))
    iend = chunk(b'IEND', b'')
    return sig + ihdr + idat + iend

def new_canvas(w, h, bg=(0,0,0,0)):
    return [[list(bg) for _ in range(w)] for _ in range(h)]

def draw_circle(pixels, cx, cy, r, color, w, h):
    for y in range(max(0, cy-r), min(h, cy+r+1)):
        for x in range(max(0, cx-r), min(w, cx+r+1)):
            if (x-cx)**2 + (y-cy)**2 <= r*r:
                pixels[y][x] = list(color)

def draw_circle_shaded(pixels, cx, cy, r, base_color, w, h):
    br, bg, bb, ba = base_color
    for y in range(max(0, cy-r), min(h, cy+r+1)):
        for x in range(max(0, cx-r), min(w, cx+r+1)):
            dist = math.sqrt((x-cx)**2 + (y-cy)**2)
            if dist <= r:
                # 光源左上
                dx = (x - cx) / r
                dy = (y - cy) / r
                light = 1.0 - 0.35 * (dx + dy)
                light = max(0.6, min(1.3, light))
                pixels[y][x] = [
                    min(255, int(br * light)),
                    min(255, int(bg * light)),
                    min(255, int(bb * light)),
                    ba
                ]

def draw_ring(pixels, cx, cy, r_outer, r_inner, color, w, h):
    for y in range(max(0, cy-r_outer), min(h, cy+r_outer+1)):
        for x in range(max(0, cx-r_outer), min(w, cx+r_outer+1)):
            dist2 = (x-cx)**2 + (y-cy)**2
            if r_inner*r_inner <= dist2 <= r_outer*r_outer:
                pixels[y][x] = list(color)

def draw_ray(pixels, cx, cy, angle_deg, length, width, color, w, h):
    angle = math.radians(angle_deg)
    for t in range(length):
        for dw in range(-width//2, width//2+1):
            px = int(cx + t * math.cos(angle) + dw * math.sin(angle))
            py = int(cy + t * math.sin(angle) - dw * math.cos(angle))
            if 0 <= px < w and 0 <= py < h:
                pixels[py][px] = list(color)

def draw_star(pixels, cx, cy, r_outer, r_inner, points, color, w, h):
    for i in range(points * 2):
        angle = math.radians(i * 180 / points - 90)
        r = r_outer if i % 2 == 0 else r_inner
        x1 = cx + int(r * math.cos(angle))
        y1 = cy + int(r * math.sin(angle))
        angle2 = math.radians((i+1) * 180 / points - 90)
        r2 = r_inner if i % 2 == 0 else r_outer
        x2 = cx + int(r2 * math.cos(angle2))
        y2 = cy + int(r2 * math.sin(angle2))
        # 畫線
        steps = max(abs(x2-x1), abs(y2-y1)) + 1
        for s in range(steps):
            px = x1 + int((x2-x1) * s / max(1, steps-1))
            py = y1 + int((y2-y1) * s / max(1, steps-1))
            if 0 <= px < w and 0 <= py < h:
                pixels[py][px] = list(color)

def draw_outline(pixels, w, h, color=(0,0,0,255)):
    for y in range(h):
        for x in range(w):
            if pixels[y][x][3] > 0:
                for dx, dy in [(-1,0),(1,0),(0,-1),(0,1)]:
                    nx, ny = x+dx, y+dy
                    if 0 <= nx < w and 0 <= ny < h and pixels[ny][nx][3] == 0:
                        pixels[ny][nx] = list(color)

def save_png(pixels, path, w, h):
    data = make_png(pixels, w, h)
    with open(path, 'wb') as f:
        f.write(data)
    count = sum(1 for y in range(h) for x in range(w) if pixels[y][x][3] > 0)
    pct = count / (w*h) * 100
    print(f"  {os.path.basename(path)}: {count} 非透明像素 ({pct:.1f}%)")

# ── T156 冰鳳凰魚 ─────────────────────────────────────────────
def gen_t156():
    p = new_canvas(SIZE, SIZE)
    cx, cy = SIZE//2, SIZE//2
    # 冰藍色魚身（橢圓）
    for y in range(SIZE):
        for x in range(SIZE):
            dx = (x - cx) / 22
            dy = (y - cy) / 16
            if dx*dx + dy*dy <= 1.0:
                # 冰藍漸層
                light = 1.0 - 0.3 * (dx + dy)
                light = max(0.7, min(1.2, light))
                r = min(255, int(100 * light))
                g = min(255, int(200 * light))
                b = min(255, int(255 * light))
                p[y][x] = [r, g, b, 255]
    # 冰晶光環（白色環）
    draw_ring(p, cx, cy, 28, 25, (200, 240, 255, 180), SIZE, SIZE)
    # 鳳凰翅膀（左右各一個弧形）
    for angle in range(-60, 60, 3):
        rad = math.radians(angle)
        for r in range(20, 30):
            px = cx - 15 + int(r * math.cos(rad))
            py = cy - 5 + int(r * math.sin(rad))
            if 0 <= px < SIZE and 0 <= py < SIZE:
                p[py][px] = [150, 220, 255, 200]
    for angle in range(120, 240, 3):
        rad = math.radians(angle)
        for r in range(20, 30):
            px = cx + 15 + int(r * math.cos(rad))
            py = cy - 5 + int(r * math.sin(rad))
            if 0 <= px < SIZE and 0 <= py < SIZE:
                p[py][px] = [150, 220, 255, 200]
    # 眼睛（紅色）
    draw_circle(p, cx+8, cy-3, 3, (255, 50, 50, 255), SIZE, SIZE)
    draw_circle(p, cx+9, cy-4, 1, (255, 255, 255, 255), SIZE, SIZE)
    # 冰晶裝飾（6個小菱形）
    for i in range(6):
        angle = math.radians(i * 60)
        ix = cx + int(20 * math.cos(angle))
        iy = cy + int(20 * math.sin(angle))
        draw_circle(p, ix, iy, 2, (200, 240, 255, 220), SIZE, SIZE)
    draw_outline(p, SIZE, SIZE, (0, 80, 150, 255))
    save_png(p, os.path.join(OUT_DIR, "T156_ice_phoenix.png"), SIZE, SIZE)

# ── T157 龍怒能量魚 ───────────────────────────────────────────
def gen_t157():
    p = new_canvas(SIZE, SIZE)
    cx, cy = SIZE//2, SIZE//2
    # 火橙色魚身
    for y in range(SIZE):
        for x in range(SIZE):
            dx = (x - cx) / 20
            dy = (y - cy) / 15
            if dx*dx + dy*dy <= 1.0:
                light = 1.0 - 0.3 * (dx + dy)
                light = max(0.7, min(1.2, light))
                r = min(255, int(220 * light))
                g = min(255, int(80 * light))
                b = min(255, int(10 * light))
                p[y][x] = [r, g, b, 255]
    # 能量光環（橙色）
    draw_ring(p, cx, cy, 28, 24, (255, 150, 0, 180), SIZE, SIZE)
    # 龍鱗紋路（深色弧線）
    for i in range(4):
        for angle in range(-30, 30, 2):
            rad = math.radians(angle)
            r = 8 + i * 4
            px = cx - 5 + int(r * math.cos(rad))
            py = cy + int(r * math.sin(rad))
            if 0 <= px < SIZE and 0 <= py < SIZE:
                p[py][px] = [150, 40, 0, 255]
    # 8方向能量射線
    for i in range(8):
        angle = i * 45
        draw_ray(p, cx, cy, angle, 14, 1, (255, 200, 50, 200), SIZE, SIZE)
    # 眼睛（金色）
    draw_circle(p, cx+7, cy-2, 3, (255, 200, 0, 255), SIZE, SIZE)
    draw_circle(p, cx+8, cy-3, 1, (255, 255, 255, 255), SIZE, SIZE)
    # 中心能量核心
    draw_circle(p, cx, cy, 5, (255, 220, 100, 255), SIZE, SIZE)
    draw_outline(p, SIZE, SIZE, (100, 20, 0, 255))
    save_png(p, os.path.join(OUT_DIR, "T157_dragon_fury.png"), SIZE, SIZE)

# ── T158 倍率瀑布魚 ───────────────────────────────────────────
def gen_t158():
    p = new_canvas(SIZE, SIZE)
    cx, cy = SIZE//2, SIZE//2
    # 深藍色魚身
    for y in range(SIZE):
        for x in range(SIZE):
            dx = (x - cx) / 21
            dy = (y - cy) / 15
            if dx*dx + dy*dy <= 1.0:
                light = 1.0 - 0.3 * (dx + dy)
                light = max(0.7, min(1.2, light))
                r = min(255, int(20 * light))
                g = min(255, int(80 * light))
                b = min(255, int(220 * light))
                p[y][x] = [r, g, b, 255]
    # 瀑布流線（白色波浪）
    for i in range(5):
        for x in range(10, 54):
            y_wave = cy - 8 + i * 4 + int(3 * math.sin(x * 0.4 + i))
            if 0 <= y_wave < SIZE:
                p[y_wave][x] = [180, 220, 255, 200]
    # 倍率符號（×）
    for i in range(-6, 7):
        px1 = cx + i
        py1 = cy + i
        px2 = cx + i
        py2 = cy - i
        if 0 <= px1 < SIZE and 0 <= py1 < SIZE:
            p[py1][px1] = [255, 255, 255, 255]
        if 0 <= px2 < SIZE and 0 <= py2 < SIZE:
            p[py2][px2] = [255, 255, 255, 255]
    # 眼睛（白色）
    draw_circle(p, cx+8, cy-2, 3, (255, 255, 255, 255), SIZE, SIZE)
    draw_circle(p, cx+8, cy-2, 1, (0, 50, 150, 255), SIZE, SIZE)
    # 外圈光環
    draw_ring(p, cx, cy, 29, 26, (100, 180, 255, 160), SIZE, SIZE)
    draw_outline(p, SIZE, SIZE, (0, 30, 120, 255))
    save_png(p, os.path.join(OUT_DIR, "T158_mult_cascade.png"), SIZE, SIZE)

# ── T159 覺醒 BOSS v2 魚 ──────────────────────────────────────
def gen_t159():
    p = new_canvas(SIZE, SIZE)
    cx, cy = SIZE//2, SIZE//2
    # 金橙色大型魚身
    for y in range(SIZE):
        for x in range(SIZE):
            dx = (x - cx) / 24
            dy = (y - cy) / 18
            if dx*dx + dy*dy <= 1.0:
                light = 1.0 - 0.3 * (dx + dy)
                light = max(0.7, min(1.2, light))
                r = min(255, int(240 * light))
                g = min(255, int(140 * light))
                b = min(255, int(0 * light))
                p[y][x] = [r, g, b, 255]
    # 覺醒光環（金色）
    draw_ring(p, cx, cy, 30, 26, (255, 200, 0, 200), SIZE, SIZE)
    # 8次 Power Up 標記（8個小圓）
    for i in range(8):
        angle = math.radians(i * 45 - 90)
        mx = cx + int(22 * math.cos(angle))
        my = cy + int(22 * math.sin(angle))
        draw_circle(p, mx, my, 3, (255, 240, 100, 255), SIZE, SIZE)
    # 閃電符號（⚡）
    pts = [(cx-3, cy-8), (cx+2, cy-2), (cx-1, cy-2), (cx+4, cy+8), (cx-1, cy+2), (cx+2, cy+2)]
    for i in range(len(pts)-1):
        x1, y1 = pts[i]
        x2, y2 = pts[i+1]
        steps = max(abs(x2-x1), abs(y2-y1)) + 1
        for s in range(steps):
            px = x1 + int((x2-x1) * s / max(1, steps-1))
            py = y1 + int((y2-y1) * s / max(1, steps-1))
            if 0 <= px < SIZE and 0 <= py < SIZE:
                p[py][px] = [255, 255, 100, 255]
    # 眼睛（紅色）
    draw_circle(p, cx+10, cy-3, 4, (255, 50, 0, 255), SIZE, SIZE)
    draw_circle(p, cx+11, cy-4, 2, (255, 255, 255, 255), SIZE, SIZE)
    draw_outline(p, SIZE, SIZE, (120, 60, 0, 255))
    save_png(p, os.path.join(OUT_DIR, "T159_awaken_boss_v2.png"), SIZE, SIZE)

# ── T160 終極審判魚 ───────────────────────────────────────────
def gen_t160():
    p = new_canvas(SIZE, SIZE)
    cx, cy = SIZE//2, SIZE//2
    # 深紅色大型魚身
    for y in range(SIZE):
        for x in range(SIZE):
            dx = (x - cx) / 26
            dy = (y - cy) / 20
            if dx*dx + dy*dy <= 1.0:
                light = 1.0 - 0.3 * (dx + dy)
                light = max(0.7, min(1.2, light))
                r = min(255, int(200 * light))
                g = min(255, int(10 * light))
                b = min(255, int(10 * light))
                p[y][x] = [r, g, b, 255]
    # 終極光環（深紅+金色）
    draw_ring(p, cx, cy, 31, 27, (200, 0, 0, 200), SIZE, SIZE)
    draw_ring(p, cx, cy, 27, 25, (255, 200, 0, 180), SIZE, SIZE)
    # 12方向審判光芒
    for i in range(12):
        angle = i * 30
        draw_ray(p, cx, cy, angle, 16, 1, (255, 150, 0, 220), SIZE, SIZE)
    # 天秤符號（⚖）
    # 橫桿
    for x in range(cx-10, cx+11):
        if 0 <= x < SIZE:
            p[cy-5][x] = [255, 220, 0, 255]
    # 中心柱
    for y in range(cy-5, cy+6):
        if 0 <= y < SIZE:
            p[y][cx] = [255, 220, 0, 255]
    # 左盤
    draw_circle(p, cx-8, cy+3, 4, (255, 200, 0, 200), SIZE, SIZE)
    # 右盤
    draw_circle(p, cx+8, cy+3, 4, (255, 200, 0, 200), SIZE, SIZE)
    # 眼睛（金色）
    draw_circle(p, cx+12, cy-4, 4, (255, 200, 0, 255), SIZE, SIZE)
    draw_circle(p, cx+13, cy-5, 2, (255, 255, 255, 255), SIZE, SIZE)
    draw_outline(p, SIZE, SIZE, (80, 0, 0, 255))
    save_png(p, os.path.join(OUT_DIR, "T160_ultimate_judgment.png"), SIZE, SIZE)

if __name__ == "__main__":
    print("=== DAY-309 T156-T160 精靈圖生成 ===")
    gen_t156()
    gen_t157()
    gen_t158()
    gen_t159()
    gen_t160()
    print("=== 完成 ===")
