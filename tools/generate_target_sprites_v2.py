#!/usr/bin/env python3
"""
generate_target_sprites_v2.py — 改善版目標物像素圖生成
target-pixel-agent 負責維護
重新生成 T001-T006 基礎目標物，讓它們更有辨識度
使用純 Python 標準庫（struct + zlib）
"""

import struct
import zlib
import os
import math

OUTPUT_DIR = os.path.join(
    os.path.dirname(os.path.dirname(os.path.abspath(__file__))),
    "client", "chiikawa-pixel", "assets", "sprites", "targets"
)

SPRITE_W = 32
SPRITE_H = 32


def write_png(filename: str, pixels: list) -> None:
    os.makedirs(OUTPUT_DIR, exist_ok=True)
    path = os.path.join(OUTPUT_DIR, filename)
    h = len(pixels)
    w = len(pixels[0]) if h > 0 else 0
    sig = b'\x89PNG\r\n\x1a\n'
    ihdr_data = struct.pack('>II', w, h) + bytes([8, 6, 0, 0, 0])
    ihdr = _make_chunk(b'IHDR', ihdr_data)
    raw_data = b''
    for row in pixels:
        raw_data += b'\x00'
        for r, g, b, a in row:
            raw_data += bytes([r, g, b, a])
    compressed = zlib.compress(raw_data, 9)
    idat = _make_chunk(b'IDAT', compressed)
    iend = _make_chunk(b'IEND', b'')
    with open(path, 'wb') as f:
        f.write(sig + ihdr + idat + iend)
    print(f"  ✅ {filename}")


def _make_chunk(chunk_type: bytes, data: bytes) -> bytes:
    length = struct.pack('>I', len(data))
    crc = struct.pack('>I', zlib.crc32(chunk_type + data) & 0xFFFFFFFF)
    return length + chunk_type + data + crc


def empty_canvas(w=SPRITE_W, h=SPRITE_H):
    return [[(0, 0, 0, 0)] * w for _ in range(h)]


def sp(canvas, x, y, color):
    if 0 <= x < len(canvas[0]) and 0 <= y < len(canvas):
        canvas[y][x] = color


def fill(canvas, x, y, w, h, color):
    for dy in range(h):
        for dx in range(w):
            sp(canvas, x + dx, y + dy, color)


def circle(canvas, cx, cy, r, color):
    for dy in range(-r, r + 1):
        for dx in range(-r, r + 1):
            if dx * dx + dy * dy <= r * r:
                sp(canvas, cx + dx, cy + dy, color)


# ── T001 像素雜草（綠色，2x）────────────────────────────────
def draw_T001():
    c = empty_canvas()
    # 莖
    fill(c, 15, 12, 2, 16, (30, 120, 30, 255))
    # 葉子
    fill(c, 8, 8, 8, 6, (50, 180, 50, 255))
    fill(c, 16, 6, 8, 6, (40, 160, 40, 255))
    fill(c, 10, 14, 6, 5, (60, 200, 60, 255))
    # 高光
    sp(c, 10, 9, (100, 230, 100, 200))
    sp(c, 18, 7, (100, 230, 100, 200))
    # 根部
    fill(c, 12, 27, 8, 3, (80, 50, 20, 255))
    return c


# ── T002 綠色小蟲（3x）──────────────────────────────────────
def draw_T002():
    c = empty_canvas()
    # 身體（橢圓）
    circle(c, 16, 18, 8, (60, 180, 60, 255))
    circle(c, 16, 18, 6, (80, 210, 80, 255))
    # 頭
    circle(c, 16, 10, 5, (50, 160, 50, 255))
    # 眼睛
    fill(c, 13, 8, 2, 2, (20, 20, 20, 255))
    fill(c, 17, 8, 2, 2, (20, 20, 20, 255))
    sp(c, 13, 8, (255, 255, 255, 200))
    sp(c, 17, 8, (255, 255, 255, 200))
    # 觸角
    sp(c, 13, 5, (40, 140, 40, 255))
    sp(c, 12, 4, (40, 140, 40, 255))
    sp(c, 19, 5, (40, 140, 40, 255))
    sp(c, 20, 4, (40, 140, 40, 255))
    # 腳
    for i in range(3):
        fill(c, 8, 16 + i * 3, 3, 2, (40, 140, 40, 255))
        fill(c, 21, 16 + i * 3, 3, 2, (40, 140, 40, 255))
    return c


# ── T003 紅色小蟲（5x）──────────────────────────────────────
def draw_T003():
    c = empty_canvas()
    circle(c, 16, 18, 8, (180, 40, 40, 255))
    circle(c, 16, 18, 6, (220, 60, 60, 255))
    circle(c, 16, 10, 5, (160, 30, 30, 255))
    fill(c, 13, 8, 2, 2, (20, 20, 20, 255))
    fill(c, 17, 8, 2, 2, (20, 20, 20, 255))
    sp(c, 13, 8, (255, 255, 255, 200))
    sp(c, 17, 8, (255, 255, 255, 200))
    # 紅色觸角（更長）
    sp(c, 13, 5, (200, 50, 50, 255))
    sp(c, 12, 4, (200, 50, 50, 255))
    sp(c, 11, 3, (200, 50, 50, 255))
    sp(c, 19, 5, (200, 50, 50, 255))
    sp(c, 20, 4, (200, 50, 50, 255))
    sp(c, 21, 3, (200, 50, 50, 255))
    for i in range(3):
        fill(c, 8, 16 + i * 3, 3, 2, (160, 40, 40, 255))
        fill(c, 21, 16 + i * 3, 3, 2, (160, 40, 40, 255))
    # 斑點
    sp(c, 14, 19, (255, 100, 100, 200))
    sp(c, 18, 17, (255, 100, 100, 200))
    return c


# ── T004 藍色小蟲（6x）──────────────────────────────────────
def draw_T004():
    c = empty_canvas()
    circle(c, 16, 18, 8, (40, 80, 200, 255))
    circle(c, 16, 18, 6, (60, 110, 230, 255))
    circle(c, 16, 10, 5, (30, 60, 180, 255))
    fill(c, 13, 8, 2, 2, (20, 20, 20, 255))
    fill(c, 17, 8, 2, 2, (20, 20, 20, 255))
    sp(c, 13, 8, (255, 255, 255, 200))
    sp(c, 17, 8, (255, 255, 255, 200))
    sp(c, 13, 5, (60, 100, 220, 255))
    sp(c, 12, 4, (60, 100, 220, 255))
    sp(c, 19, 5, (60, 100, 220, 255))
    sp(c, 20, 4, (60, 100, 220, 255))
    for i in range(3):
        fill(c, 8, 16 + i * 3, 3, 2, (40, 80, 180, 255))
        fill(c, 21, 16 + i * 3, 3, 2, (40, 80, 180, 255))
    # 藍色光澤
    sp(c, 15, 16, (150, 200, 255, 180))
    sp(c, 17, 14, (150, 200, 255, 180))
    return c


# ── T005 會走路的布丁（8x）──────────────────────────────────
def draw_T005():
    c = empty_canvas()
    # 布丁主體（黃色圓頂）
    circle(c, 16, 16, 10, (220, 180, 60, 255))
    circle(c, 16, 16, 8, (255, 220, 80, 255))
    # 焦糖頂部
    fill(c, 10, 6, 12, 5, (180, 100, 20, 255))
    circle(c, 16, 8, 4, (200, 120, 30, 255))
    # 眼睛
    fill(c, 12, 13, 3, 3, (40, 20, 10, 255))
    fill(c, 17, 13, 3, 3, (40, 20, 10, 255))
    sp(c, 12, 13, (255, 255, 255, 220))
    sp(c, 17, 13, (255, 255, 255, 220))
    # 嘴巴（微笑）
    sp(c, 14, 18, (40, 20, 10, 255))
    sp(c, 15, 19, (40, 20, 10, 255))
    sp(c, 16, 19, (40, 20, 10, 255))
    sp(c, 17, 18, (40, 20, 10, 255))
    # 腳
    fill(c, 10, 25, 4, 4, (200, 160, 50, 255))
    fill(c, 18, 25, 4, 4, (200, 160, 50, 255))
    # 高光
    sp(c, 13, 11, (255, 255, 200, 180))
    return c


# ── T006 巨大蘑菇（10x）─────────────────────────────────────
def draw_T006():
    c = empty_canvas()
    # 蘑菇傘（棕色）
    circle(c, 16, 12, 12, (120, 70, 30, 255))
    circle(c, 16, 12, 10, (160, 90, 40, 255))
    # 白色斑點
    circle(c, 11, 9, 2, (240, 240, 240, 255))
    circle(c, 20, 8, 2, (240, 240, 240, 255))
    circle(c, 16, 6, 2, (240, 240, 240, 255))
    circle(c, 13, 14, 2, (240, 240, 240, 255))
    # 蘑菇柄（白色）
    fill(c, 13, 22, 6, 8, (230, 230, 220, 255))
    # 眼睛
    fill(c, 13, 11, 2, 2, (40, 20, 10, 255))
    fill(c, 18, 11, 2, 2, (40, 20, 10, 255))
    sp(c, 13, 11, (255, 255, 255, 200))
    sp(c, 18, 11, (255, 255, 255, 200))
    # 嘴巴
    sp(c, 15, 15, (40, 20, 10, 255))
    sp(c, 16, 16, (40, 20, 10, 255))
    sp(c, 17, 15, (40, 20, 10, 255))
    return c


def main():
    print("🎨 生成改善版基礎目標物像素圖...")
    print(f"   輸出目錄：{OUTPUT_DIR}")
    print()
    sprites = [
        ("T001_grass.png",    draw_T001()),
        ("T002_bug_g.png",    draw_T002()),
        ("T003_bug_r.png",    draw_T003()),
        ("T004_bug_b.png",    draw_T004()),
        ("T005_pudding.png",  draw_T005()),
        ("T006_mushroom.png", draw_T006()),
    ]
    for filename, pixels in sprites:
        write_png(filename, pixels)
    print()
    print(f"✅ 完成！共生成 {len(sprites)} 個目標物精靈圖")


if __name__ == "__main__":
    main()
