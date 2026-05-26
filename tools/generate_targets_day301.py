"""
generate_targets_day301.py — DAY-301 新增目標物精靈圖生成
T126 幸運進階 Jackpot 魚 / T127 幸運全服合作魚 / T128 幸運時間扭曲魚
"""
import struct, zlib, os

OUT_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
SIZE = 64

def make_png(pixels):
    """pixels: list of (r,g,b,a) tuples, SIZE*SIZE"""
    raw = b""
    for y in range(SIZE):
        raw += b"\x00"
        for x in range(SIZE):
            r, g, b, a = pixels[y * SIZE + x]
            raw += bytes([r, g, b, a])
    compressed = zlib.compress(raw, 9)
    def chunk(name, data):
        c = name + data
        return struct.pack(">I", len(data)) + c + struct.pack(">I", zlib.crc32(c) & 0xFFFFFFFF)
    ihdr = struct.pack(">IIBBBBB", SIZE, SIZE, 8, 2, 0, 0, 0)
    # RGBA = color type 6
    ihdr = struct.pack(">II", SIZE, SIZE) + bytes([8, 6, 0, 0, 0])
    png = b"\x89PNG\r\n\x1a\n"
    png += chunk(b"IHDR", ihdr)
    png += chunk(b"IDAT", compressed)
    png += chunk(b"IEND", b"")
    return png

def px(pixels, x, y, r, g, b, a=255):
    if 0 <= x < SIZE and 0 <= y < SIZE:
        pixels[y * SIZE + x] = (r, g, b, a)

def fill_circle(pixels, cx, cy, radius, r, g, b, a=255):
    for dy in range(-radius, radius + 1):
        for dx in range(-radius, radius + 1):
            if dx*dx + dy*dy <= radius*radius:
                px(pixels, cx+dx, cy+dy, r, g, b, a)

def fill_circle_shaded(pixels, cx, cy, radius, light, mid, dark):
    for dy in range(-radius, radius + 1):
        for dx in range(-radius, radius + 1):
            if dx*dx + dy*dy <= radius*radius:
                if dx < -radius//3 and dy < -radius//3:
                    c = light
                elif dx > radius//3 and dy > radius//3:
                    c = dark
                else:
                    c = mid
                px(pixels, cx+dx, cy+dy, *c)

def draw_star(pixels, cx, cy, size, r, g, b):
    """畫星形"""
    for i in range(size):
        px(pixels, cx+i, cy, r, g, b)
        px(pixels, cx-i, cy, r, g, b)
        px(pixels, cx, cy+i, r, g, b)
        px(pixels, cx, cy-i, r, g, b)
    for i in range(size//2):
        px(pixels, cx+i, cy+i, r, g, b)
        px(pixels, cx-i, cy+i, r, g, b)
        px(pixels, cx+i, cy-i, r, g, b)
        px(pixels, cx-i, cy-i, r, g, b)

# ── T126 幸運進階 Jackpot 魚 ──────────────────────────────────
# 設計：金色魚身 + 四層獎池光環（Mini/Minor/Major/Grand）+ 王冠 + 閃光
def gen_t126():
    pixels = [(0,0,0,0)] * (SIZE * SIZE)
    cx, cy = 32, 32

    # 魚身（金色橢圓）
    for dy in range(-14, 15):
        for dx in range(-20, 21):
            if (dx/20)**2 + (dy/14)**2 <= 1.0:
                # 金色漸層
                if dx < -6 and dy < -4:
                    c = (255, 240, 100)  # 亮金
                elif dx > 6 and dy > 4:
                    c = (180, 130, 0)    # 暗金
                else:
                    c = (255, 215, 0)    # 金色
                px(pixels, cx+dx, cy+dy, *c)

    # 四層光環（由外到內：Grand/Major/Minor/Mini）
    # Grand 光環（最外層，金色）
    for angle_deg in range(0, 360, 3):
        import math
        angle = math.radians(angle_deg)
        rx, ry = int(cx + 28 * math.cos(angle)), int(cy + 28 * math.sin(angle))
        px(pixels, rx, ry, 255, 215, 0, 200)
    # Major 光環（橙色）
    for angle_deg in range(0, 360, 4):
        angle = math.radians(angle_deg)
        rx, ry = int(cx + 24 * math.cos(angle)), int(cy + 24 * math.sin(angle))
        px(pixels, rx, ry, 255, 140, 0, 180)
    # Minor 光環（銀色）
    for angle_deg in range(0, 360, 5):
        angle = math.radians(angle_deg)
        rx, ry = int(cx + 20 * math.cos(angle)), int(cy + 20 * math.sin(angle))
        px(pixels, rx, ry, 200, 200, 220, 160)
    # Mini 光環（銅色）
    for angle_deg in range(0, 360, 6):
        angle = math.radians(angle_deg)
        rx, ry = int(cx + 16 * math.cos(angle)), int(cy + 16 * math.sin(angle))
        px(pixels, rx, ry, 180, 100, 50, 140)

    # 王冠（頂部）
    crown_y = cy - 16
    for i in range(-6, 7):
        px(pixels, cx+i, crown_y, 255, 215, 0)
    for i in range(-6, 7, 3):
        for j in range(4):
            px(pixels, cx+i, crown_y-j, 255, 215, 0)
    # 王冠寶石
    fill_circle(pixels, cx-4, crown_y-3, 2, 255, 50, 50)
    fill_circle(pixels, cx, crown_y-4, 2, 50, 200, 255)
    fill_circle(pixels, cx+4, crown_y-3, 2, 50, 255, 50)

    # 魚眼（金色眼睛）
    fill_circle(pixels, cx+12, cy-4, 4, 255, 255, 255)
    fill_circle(pixels, cx+12, cy-4, 2, 255, 215, 0)
    px(pixels, cx+13, cy-5, 255, 255, 255)

    # 魚尾（金色）
    for dy in range(-8, 9):
        for dx in range(0, 10):
            if abs(dy) <= 8 - dx//2:
                px(pixels, cx+20+dx, cy+dy, 200, 160, 0)

    # 閃光星星（4個角落）
    draw_star(pixels, cx-18, cy-18, 4, 255, 255, 200)
    draw_star(pixels, cx+18, cy-18, 3, 255, 255, 200)
    draw_star(pixels, cx-18, cy+18, 3, 255, 255, 200)
    draw_star(pixels, cx+18, cy+18, 4, 255, 255, 200)

    # 輪廓
    for dy in range(-14, 15):
        for dx in range(-20, 21):
            if abs((dx/20)**2 + (dy/14)**2 - 1.0) < 0.15:
                px(pixels, cx+dx, cy+dy, 150, 100, 0)

    return pixels

# ── T127 幸運全服合作魚 ───────────────────────────────────────
# 設計：青藍色魚身 + 多個小魚圍繞（合作感）+ 握手圖案 + 光環
def gen_t127():
    pixels = [(0,0,0,0)] * (SIZE * SIZE)
    cx, cy = 32, 32
    import math

    # 主魚身（青藍色橢圓）
    for dy in range(-13, 14):
        for dx in range(-18, 19):
            if (dx/18)**2 + (dy/13)**2 <= 1.0:
                if dx < -5 and dy < -4:
                    c = (100, 240, 255)  # 亮青藍
                elif dx > 5 and dy > 4:
                    c = (0, 150, 180)    # 暗青藍
                else:
                    c = (0, 200, 230)    # 青藍
                px(pixels, cx+dx, cy+dy, *c)

    # 圍繞的小魚（4個，代表合作）
    small_fish_positions = [
        (cx-22, cy-10), (cx+22, cy-10),
        (cx-22, cy+10), (cx+22, cy+10)
    ]
    small_colors = [
        (255, 150, 50),   # 橙
        (150, 255, 100),  # 綠
        (255, 100, 150),  # 粉
        (200, 150, 255),  # 紫
    ]
    for (sx, sy), sc in zip(small_fish_positions, small_colors):
        fill_circle(pixels, sx, sy, 5, *sc)
        # 小魚眼
        px(pixels, sx+3, sy-1, 255, 255, 255)

    # 連接線（合作感）
    for i in range(8):
        t = i / 7.0
        lx = int(cx-22 + t * 44)
        ly = int(cy-10 + t * 0)
        px(pixels, lx, cy-10, 0, 230, 255, 120)
        px(pixels, lx, cy+10, 0, 230, 255, 120)

    # 中心握手圖案（簡化版：兩個圓圈相交）
    fill_circle(pixels, cx-4, cy, 5, 255, 255, 255, 180)
    fill_circle(pixels, cx+4, cy, 5, 255, 255, 255, 180)
    fill_circle(pixels, cx, cy, 4, 0, 230, 255, 200)

    # 魚眼
    fill_circle(pixels, cx+10, cy-3, 4, 255, 255, 255)
    fill_circle(pixels, cx+10, cy-3, 2, 0, 180, 220)
    px(pixels, cx+11, cy-4, 255, 255, 255)

    # 魚尾
    for dy in range(-7, 8):
        for dx in range(0, 9):
            if abs(dy) <= 7 - dx//2:
                px(pixels, cx+18+dx, cy+dy, 0, 160, 190)

    # 光環（青藍色）
    for angle_deg in range(0, 360, 4):
        angle = math.radians(angle_deg)
        rx, ry = int(cx + 30 * math.cos(angle)), int(cy + 30 * math.sin(angle))
        px(pixels, rx, ry, 0, 230, 255, 150)

    # 輪廓
    for dy in range(-13, 14):
        for dx in range(-18, 19):
            if abs((dx/18)**2 + (dy/13)**2 - 1.0) < 0.15:
                px(pixels, cx+dx, cy+dy, 0, 100, 130)

    return pixels

# ── T128 幸運時間扭曲魚 ───────────────────────────────────────
# 設計：紫色魚身 + 時鐘圖案 + 扭曲光環 + 時間裂縫紋路
def gen_t128():
    pixels = [(0,0,0,0)] * (SIZE * SIZE)
    cx, cy = 32, 32
    import math

    # 魚身（深紫色橢圓）
    for dy in range(-13, 14):
        for dx in range(-19, 20):
            if (dx/19)**2 + (dy/13)**2 <= 1.0:
                if dx < -6 and dy < -4:
                    c = (200, 100, 255)  # 亮紫
                elif dx > 6 and dy > 4:
                    c = (80, 0, 140)     # 暗紫
                else:
                    c = (140, 50, 220)   # 紫色
                px(pixels, cx+dx, cy+dy, *c)

    # 時鐘圖案（中心）
    fill_circle(pixels, cx, cy, 8, 255, 255, 255, 180)
    fill_circle(pixels, cx, cy, 7, 200, 150, 255, 200)
    # 時鐘刻度（12個）
    for i in range(12):
        angle = math.radians(i * 30 - 90)
        tx = int(cx + 6 * math.cos(angle))
        ty = int(cy + 6 * math.sin(angle))
        px(pixels, tx, ty, 80, 0, 140)
    # 時鐘指針（扭曲的）
    px(pixels, cx, cy-4, 80, 0, 140)
    px(pixels, cx, cy-3, 80, 0, 140)
    px(pixels, cx+3, cy, 80, 0, 140)
    px(pixels, cx+2, cy, 80, 0, 140)

    # 扭曲光環（螺旋感）
    for i in range(60):
        angle = math.radians(i * 6)
        r = 16 + 4 * math.sin(angle * 3)
        rx = int(cx + r * math.cos(angle))
        ry = int(cy + r * math.sin(angle))
        px(pixels, rx, ry, 200, 100, 255, 180)

    # 時間裂縫紋路（Z字形）
    for i in range(-8, 9):
        px(pixels, cx+i, cy-8+abs(i)//2, 255, 200, 255, 150)
        px(pixels, cx+i, cy+8-abs(i)//2, 255, 200, 255, 150)

    # 魚眼（紫色眼睛）
    fill_circle(pixels, cx+11, cy-3, 4, 255, 255, 255)
    fill_circle(pixels, cx+11, cy-3, 2, 140, 50, 220)
    px(pixels, cx+12, cy-4, 255, 255, 255)

    # 魚尾（紫色）
    for dy in range(-7, 8):
        for dx in range(0, 9):
            if abs(dy) <= 7 - dx//2:
                px(pixels, cx+19+dx, cy+dy, 100, 30, 180)

    # 外圈光環（紫色）
    for angle_deg in range(0, 360, 3):
        angle = math.radians(angle_deg)
        rx, ry = int(cx + 29 * math.cos(angle)), int(cy + 29 * math.sin(angle))
        px(pixels, rx, ry, 200, 100, 255, 160)

    # 輪廓
    for dy in range(-13, 14):
        for dx in range(-19, 20):
            if abs((dx/19)**2 + (dy/13)**2 - 1.0) < 0.15:
                px(pixels, cx+dx, cy+dy, 60, 0, 100)

    return pixels

def save_png(filename, pixels):
    data = make_png(pixels)
    with open(filename, 'wb') as f:
        f.write(data)
    print(f"Saved: {filename} ({len([p for p in pixels if p[3] > 0])} non-transparent pixels)")

if __name__ == "__main__":
    import math  # 確保 math 在全域可用
    os.makedirs(OUT_DIR, exist_ok=True)

    print("Generating T126-T128 sprites...")
    save_png(os.path.join(OUT_DIR, "T126_jackpot_fish.png"), gen_t126())
    save_png(os.path.join(OUT_DIR, "T127_coop_fish.png"), gen_t127())
    save_png(os.path.join(OUT_DIR, "T128_time_warp.png"), gen_t128())
    print("Done!")
