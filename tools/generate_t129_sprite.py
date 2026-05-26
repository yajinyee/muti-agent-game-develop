#!/usr/bin/env python3
"""
generate_t129_sprite.py — T129 幸運連鎖隕石魚精靈圖生成
target-pixel-agent 負責維護
"""
import struct, zlib, os

def write_png(path, width, height, pixels):
    """pixels: list of (r,g,b,a) tuples, row-major"""
    def chunk(name, data):
        c = name + data
        return struct.pack('>I', len(data)) + c + struct.pack('>I', zlib.crc32(c) & 0xFFFFFFFF)
    
    raw = b''
    for y in range(height):
        raw += b'\x00'
        for x in range(width):
            r, g, b, a = pixels[y * width + x]
            raw += bytes([r, g, b, a])
    
    compressed = zlib.compress(raw, 9)
    png = b'\x89PNG\r\n\x1a\n'
    png += chunk(b'IHDR', struct.pack('>IIBBBBB', width, height, 8, 2, 0, 0, 0)[:13])
    # Fix: IHDR needs color type 6 (RGBA)
    ihdr_data = struct.pack('>II', width, height) + bytes([8, 6, 0, 0, 0])
    png = b'\x89PNG\r\n\x1a\n'
    png += chunk(b'IHDR', ihdr_data)
    png += chunk(b'IDAT', compressed)
    png += chunk(b'IEND', b'')
    
    with open(path, 'wb') as f:
        f.write(png)

def dist(x1, y1, x2, y2):
    return ((x1-x2)**2 + (y1-y2)**2) ** 0.5

def generate_t129():
    W, H = 64, 64
    pixels = [(0, 0, 0, 0)] * (W * H)
    
    def px(x, y, r, g, b, a=255):
        if 0 <= x < W and 0 <= y < H:
            pixels[y * W + x] = (r, g, b, a)
    
    def fill_circle(cx, cy, radius, r, g, b, a=255):
        for dy in range(-radius, radius+1):
            for dx in range(-radius, radius+1):
                if dx*dx + dy*dy <= radius*radius:
                    px(cx+dx, cy+dy, r, g, b, a)
    
    def fill_ellipse(cx, cy, rx, ry, r, g, b, a=255):
        for dy in range(-ry, ry+1):
            for dx in range(-rx, rx+1):
                if (dx/rx)**2 + (dy/ry)**2 <= 1.0:
                    px(cx+dx, cy+dy, r, g, b, a)
    
    cx, cy = 32, 34
    
    # 1. 外層火橙光環（隕石魚的特徵光環）
    for dy in range(-26, 27):
        for dx in range(-26, 27):
            d = (dx*dx + dy*dy) ** 0.5
            if 22 <= d <= 26:
                alpha = int(180 * (1 - abs(d - 24) / 2))
                px(cx+dx, cy+dy, 230, 100, 25, alpha)
    
    # 2. 主體：橢圓形火橙漸層魚身
    for dy in range(-18, 19):
        for dx in range(-22, 23):
            if (dx/22)**2 + (dy/18)**2 <= 1.0:
                # 漸層：中心亮橙，邊緣深橙
                d_ratio = ((dx/22)**2 + (dy/18)**2) ** 0.5
                if d_ratio < 0.4:
                    r, g, b = 255, 140, 40   # 亮橙核心
                elif d_ratio < 0.7:
                    r, g, b = 230, 100, 20   # 中橙
                else:
                    r, g, b = 180, 60, 10    # 深橙邊緣
                px(cx+dx, cy+dy, r, g, b)
    
    # 3. 隕石紋路（斜向裂縫，代表隕石特徵）
    # 主裂縫：左上到右下
    for i in range(-15, 16):
        px(cx + i, cy + i//2, 255, 200, 80, 200)
        px(cx + i, cy + i//2 + 1, 255, 180, 60, 150)
    # 次裂縫：右上到左下
    for i in range(-10, 11):
        px(cx + i, cy - i//2 + 3, 255, 200, 80, 180)
    
    # 4. 火焰光芒（4方向，代表隕石燃燒）
    for i in range(1, 12):
        alpha = int(220 * (1 - i/12))
        px(cx + i + 22, cy, 255, 120, 30, alpha)      # 右
        px(cx - i - 22, cy, 255, 120, 30, alpha)      # 左
        px(cx, cy + i + 18, 255, 120, 30, alpha)      # 下
        px(cx, cy - i - 18, 255, 120, 30, alpha)      # 上
    
    # 5. 斜向火焰光芒
    for i in range(1, 9):
        alpha = int(180 * (1 - i/9))
        px(cx + i + 16, cy + i + 13, 255, 150, 50, alpha)
        px(cx - i - 16, cy + i + 13, 255, 150, 50, alpha)
        px(cx + i + 16, cy - i - 13, 255, 150, 50, alpha)
        px(cx - i - 16, cy - i - 13, 255, 150, 50, alpha)
    
    # 6. 隕石核心（金色發光核心）
    fill_circle(cx, cy, 5, 255, 220, 80)
    fill_circle(cx, cy, 3, 255, 240, 120)
    fill_circle(cx, cy, 1, 255, 255, 200)
    
    # 7. 魚眼（橙紅色，代表燃燒的眼睛）
    fill_circle(cx - 6, cy - 4, 4, 200, 50, 10)
    fill_circle(cx - 6, cy - 4, 2, 255, 80, 20)
    px(cx - 7, cy - 5, 255, 200, 100)  # 高光
    
    # 8. 魚尾（深橙色，向右延伸）
    for i in range(8):
        h = 10 - i
        for j in range(-h, h+1):
            alpha = int(220 * (1 - i/8))
            px(cx + 22 + i, cy + j, 180, 60, 10, alpha)
    
    # 9. 連鎖標記（5個小隕石點，代表連鎖5顆）
    chain_positions = [
        (cx - 18, cy - 20),
        (cx - 8, cy - 22),
        (cx + 2, cy - 22),
        (cx + 12, cy - 20),
        (cx + 20, cy - 16),
    ]
    for i, (mx, my) in enumerate(chain_positions):
        # 小隕石
        fill_circle(mx, my, 3, 255, 140, 40)
        fill_circle(mx, my, 1, 255, 220, 80)
        # 連接線
        if i < len(chain_positions) - 1:
            nx, ny = chain_positions[i+1]
            steps = max(abs(nx-mx), abs(ny-my))
            if steps > 0:
                for s in range(steps+1):
                    lx = mx + (nx-mx)*s//steps
                    ly = my + (ny-my)*s//steps
                    px(lx, ly, 255, 180, 60, 150)
    
    # 10. 散落火花粒子
    import random
    rng = random.Random(129)
    for _ in range(20):
        sx = rng.randint(5, 58)
        sy = rng.randint(5, 58)
        d = dist(sx, sy, cx, cy)
        if d > 20:
            alpha = rng.randint(100, 200)
            r_val = rng.randint(200, 255)
            g_val = rng.randint(80, 160)
            px(sx, sy, r_val, g_val, 30, alpha)
    
    # 11. 輪廓（深橙色輪廓）
    outline_pixels = set()
    for y in range(H):
        for x in range(W):
            if pixels[y * W + x][3] > 0:
                for dy2, dx2 in [(-1,0),(1,0),(0,-1),(0,1)]:
                    nx2, ny2 = x+dx2, y+dy2
                    if 0 <= nx2 < W and 0 <= ny2 < H and pixels[ny2*W+nx2][3] == 0:
                        outline_pixels.add((x, y))
    for (ox, oy) in outline_pixels:
        pixels[oy * W + ox] = (120, 40, 5, 255)
    
    return pixels

if __name__ == "__main__":
    out_dir = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
    os.makedirs(out_dir, exist_ok=True)
    
    pixels = generate_t129()
    out_path = os.path.join(out_dir, "T129_chain_meteor.png")
    write_png(out_path, 64, 64, pixels)
    
    # 統計非透明像素
    non_transparent = sum(1 for p in pixels if p[3] > 0)
    total = 64 * 64
    print(f"T129 連鎖隕石魚精靈圖生成完成：{out_path}")
    print(f"非透明像素：{non_transparent}/{total} ({non_transparent/total*100:.1f}%)")
