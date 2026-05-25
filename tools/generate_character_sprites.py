#!/usr/bin/env python3
"""
generate_character_sprites.py — 生成角色像素圖（idle/attack/bigwin 三狀態）
character-pixel-agent + character-animation-agent 負責維護
使用純 Python 標準庫（struct + zlib），不需要 Pillow
輸出到 client/chiikawa-pixel/assets/sprites/characters/
"""

import struct
import zlib
import os
import math

OUTPUT_DIR = os.path.join(
    os.path.dirname(os.path.dirname(os.path.abspath(__file__))),
    "client", "chiikawa-pixel", "assets", "sprites", "characters"
)

# 每個角色 3 個狀態（idle/attack/bigwin），每個狀態 1 幀（32x32 像素）
SPRITE_W = 32
SPRITE_H = 32


def write_png(filename: str, pixels: list[list[tuple]]) -> None:
    """
    pixels: H x W 的 RGBA tuple 列表
    寫入 PNG 檔案
    """
    os.makedirs(OUTPUT_DIR, exist_ok=True)
    path = os.path.join(OUTPUT_DIR, filename)
    h = len(pixels)
    w = len(pixels[0]) if h > 0 else 0

    # PNG signature
    sig = b'\x89PNG\r\n\x1a\n'

    # IHDR chunk
    ihdr_data = struct.pack('>IIBBBBB', w, h, 8, 2, 0, 0, 0)  # 8-bit RGB
    # Actually use RGBA (color type 6)
    ihdr_data = struct.pack('>II', w, h) + bytes([8, 6, 0, 0, 0])
    ihdr = _make_chunk(b'IHDR', ihdr_data)

    # IDAT chunk (image data)
    raw_data = b''
    for row in pixels:
        raw_data += b'\x00'  # filter type None
        for r, g, b, a in row:
            raw_data += bytes([r, g, b, a])
    compressed = zlib.compress(raw_data, 9)
    idat = _make_chunk(b'IDAT', compressed)

    # IEND chunk
    iend = _make_chunk(b'IEND', b'')

    with open(path, 'wb') as f:
        f.write(sig + ihdr + idat + iend)
    print(f"  ✅ {filename} ({w}x{h})")


def _make_chunk(chunk_type: bytes, data: bytes) -> bytes:
    length = struct.pack('>I', len(data))
    crc = struct.pack('>I', zlib.crc32(chunk_type + data) & 0xFFFFFFFF)
    return length + chunk_type + data + crc


def empty_canvas(w: int = SPRITE_W, h: int = SPRITE_H) -> list[list[tuple]]:
    return [[(0, 0, 0, 0)] * w for _ in range(h)]


def set_pixel(canvas, x: int, y: int, color: tuple) -> None:
    if 0 <= x < len(canvas[0]) and 0 <= y < len(canvas):
        canvas[y][x] = color


def fill_rect(canvas, x: int, y: int, w: int, h: int, color: tuple) -> None:
    for dy in range(h):
        for dx in range(w):
            set_pixel(canvas, x + dx, y + dy, color)


def draw_circle(canvas, cx: int, cy: int, r: int, color: tuple, fill: bool = True) -> None:
    for dy in range(-r, r + 1):
        for dx in range(-r, r + 1):
            if dx * dx + dy * dy <= r * r:
                if fill:
                    set_pixel(canvas, cx + dx, cy + dy, color)
            elif abs(dx * dx + dy * dy - r * r) <= r:
                set_pixel(canvas, cx + dx, cy + dy, color)


# ── 吉伊卡哇（Chiikawa）────────────────────────────────────────
# 特徵：圓臉、小耳朵、粉色系

CHIIKAWA_SKIN = (255, 220, 200, 255)
CHIIKAWA_EAR  = (255, 180, 170, 255)
CHIIKAWA_EYE  = (40, 30, 30, 255)
CHIIKAWA_BLUSH = (255, 160, 160, 200)
CHIIKAWA_BODY = (255, 230, 210, 255)
CHIIKAWA_OUTLINE = (60, 40, 40, 255)


def draw_chiikawa_idle() -> list[list[tuple]]:
    c = empty_canvas()
    # 身體
    fill_rect(c, 10, 18, 12, 10, CHIIKAWA_BODY)
    # 頭（圓形）
    draw_circle(c, 16, 13, 9, CHIIKAWA_SKIN)
    # 耳朵
    draw_circle(c, 9, 6, 3, CHIIKAWA_EAR)
    draw_circle(c, 23, 6, 3, CHIIKAWA_EAR)
    # 眼睛
    fill_rect(c, 12, 11, 3, 3, CHIIKAWA_EYE)
    fill_rect(c, 18, 11, 3, 3, CHIIKAWA_EYE)
    # 眼睛高光
    set_pixel(c, 13, 11, (255, 255, 255, 200))
    set_pixel(c, 19, 11, (255, 255, 255, 200))
    # 腮紅
    fill_rect(c, 9, 14, 4, 3, CHIIKAWA_BLUSH)
    fill_rect(c, 20, 14, 4, 3, CHIIKAWA_BLUSH)
    # 嘴巴
    set_pixel(c, 15, 16, CHIIKAWA_EYE)
    set_pixel(c, 16, 17, CHIIKAWA_EYE)
    set_pixel(c, 17, 16, CHIIKAWA_EYE)
    # 手臂
    fill_rect(c, 7, 19, 3, 6, CHIIKAWA_BODY)
    fill_rect(c, 22, 19, 3, 6, CHIIKAWA_BODY)
    # 腳
    fill_rect(c, 11, 27, 4, 3, CHIIKAWA_BODY)
    fill_rect(c, 17, 27, 4, 3, CHIIKAWA_BODY)
    return c


def draw_chiikawa_attack() -> list[list[tuple]]:
    c = draw_chiikawa_idle()
    # 攻擊：右手舉起，身體前傾
    fill_rect(c, 22, 14, 3, 6, CHIIKAWA_BODY)  # 右手舉高
    fill_rect(c, 24, 12, 4, 4, CHIIKAWA_BODY)  # 拳頭
    # 動作線
    set_pixel(c, 27, 11, (255, 200, 50, 200))
    set_pixel(c, 28, 10, (255, 200, 50, 200))
    set_pixel(c, 29, 9, (255, 200, 50, 150))
    # 眼睛變成決心狀
    fill_rect(c, 12, 11, 3, 2, CHIIKAWA_EYE)
    fill_rect(c, 18, 11, 3, 2, CHIIKAWA_EYE)
    return c


def draw_chiikawa_bigwin() -> list[list[tuple]]:
    c = draw_chiikawa_idle()
    # 大獎：雙手舉起，嘴巴大開
    fill_rect(c, 7, 14, 3, 6, CHIIKAWA_BODY)   # 左手舉高
    fill_rect(c, 22, 14, 3, 6, CHIIKAWA_BODY)  # 右手舉高
    fill_rect(c, 5, 12, 4, 4, CHIIKAWA_BODY)   # 左拳
    fill_rect(c, 24, 12, 4, 4, CHIIKAWA_BODY)  # 右拳
    # 嘴巴大開（驚喜）
    fill_rect(c, 14, 15, 5, 3, (200, 100, 100, 255))
    # 星星特效
    set_pixel(c, 3, 5, (255, 220, 0, 255))
    set_pixel(c, 4, 4, (255, 220, 0, 255))
    set_pixel(c, 5, 5, (255, 220, 0, 255))
    set_pixel(c, 28, 5, (255, 220, 0, 255))
    set_pixel(c, 27, 4, (255, 220, 0, 255))
    set_pixel(c, 29, 4, (255, 220, 0, 255))
    return c


# ── 小八（Hachiware）──────────────────────────────────────────
# 特徵：藍白條紋、尖耳朵

HACHIWARE_SKIN  = (220, 235, 255, 255)
HACHIWARE_STRIPE = (100, 140, 220, 255)
HACHIWARE_EYE   = (30, 30, 60, 255)
HACHIWARE_BODY  = (200, 220, 255, 255)
HACHIWARE_NOSE  = (255, 150, 150, 255)


def draw_hachiware_idle() -> list[list[tuple]]:
    c = empty_canvas()
    # 身體
    fill_rect(c, 10, 18, 12, 10, HACHIWARE_BODY)
    # 頭
    draw_circle(c, 16, 13, 9, HACHIWARE_SKIN)
    # 條紋（特徵）
    for i in range(3):
        fill_rect(c, 8 + i * 3, 8, 2, 10, HACHIWARE_STRIPE)
    # 尖耳朵
    fill_rect(c, 8, 2, 4, 6, HACHIWARE_SKIN)
    fill_rect(c, 20, 2, 4, 6, HACHIWARE_SKIN)
    set_pixel(c, 9, 1, HACHIWARE_SKIN)
    set_pixel(c, 10, 0, HACHIWARE_SKIN)
    set_pixel(c, 21, 1, HACHIWARE_SKIN)
    set_pixel(c, 22, 0, HACHIWARE_SKIN)
    # 眼睛
    fill_rect(c, 12, 11, 3, 3, HACHIWARE_EYE)
    fill_rect(c, 18, 11, 3, 3, HACHIWARE_EYE)
    set_pixel(c, 13, 11, (255, 255, 255, 200))
    set_pixel(c, 19, 11, (255, 255, 255, 200))
    # 鼻子
    set_pixel(c, 16, 15, HACHIWARE_NOSE)
    # 嘴
    set_pixel(c, 15, 16, HACHIWARE_EYE)
    set_pixel(c, 16, 17, HACHIWARE_EYE)
    set_pixel(c, 17, 16, HACHIWARE_EYE)
    # 手腳
    fill_rect(c, 7, 19, 3, 6, HACHIWARE_BODY)
    fill_rect(c, 22, 19, 3, 6, HACHIWARE_BODY)
    fill_rect(c, 11, 27, 4, 3, HACHIWARE_BODY)
    fill_rect(c, 17, 27, 4, 3, HACHIWARE_BODY)
    return c


def draw_hachiware_attack() -> list[list[tuple]]:
    c = draw_hachiware_idle()
    fill_rect(c, 22, 14, 3, 6, HACHIWARE_BODY)
    fill_rect(c, 24, 12, 4, 4, HACHIWARE_BODY)
    # 藍色能量
    set_pixel(c, 27, 11, (100, 180, 255, 220))
    set_pixel(c, 28, 10, (100, 180, 255, 200))
    set_pixel(c, 29, 9, (100, 180, 255, 150))
    return c


def draw_hachiware_bigwin() -> list[list[tuple]]:
    c = draw_hachiware_idle()
    fill_rect(c, 7, 14, 3, 6, HACHIWARE_BODY)
    fill_rect(c, 22, 14, 3, 6, HACHIWARE_BODY)
    fill_rect(c, 5, 12, 4, 4, HACHIWARE_BODY)
    fill_rect(c, 24, 12, 4, 4, HACHIWARE_BODY)
    fill_rect(c, 14, 15, 5, 3, (150, 100, 200, 255))
    # 藍色星星
    set_pixel(c, 3, 5, (100, 180, 255, 255))
    set_pixel(c, 4, 4, (100, 180, 255, 255))
    set_pixel(c, 28, 5, (100, 180, 255, 255))
    set_pixel(c, 27, 4, (100, 180, 255, 255))
    return c


# ── 烏薩奇（Usagi）────────────────────────────────────────────
# 特徵：長耳朵、黃色系、兔子

USAGI_SKIN  = (255, 245, 200, 255)
USAGI_EAR   = (255, 200, 150, 255)
USAGI_EYE   = (50, 30, 20, 255)
USAGI_BODY  = (255, 240, 190, 255)
USAGI_BLUSH = (255, 180, 100, 200)


def draw_usagi_idle() -> list[list[tuple]]:
    c = empty_canvas()
    # 身體
    fill_rect(c, 10, 18, 12, 10, USAGI_BODY)
    # 頭
    draw_circle(c, 16, 14, 9, USAGI_SKIN)
    # 長耳朵（特徵）
    fill_rect(c, 10, 0, 4, 10, USAGI_EAR)
    fill_rect(c, 18, 0, 4, 10, USAGI_EAR)
    fill_rect(c, 11, 1, 2, 8, (255, 220, 180, 255))  # 耳內
    fill_rect(c, 19, 1, 2, 8, (255, 220, 180, 255))
    # 眼睛（大眼）
    fill_rect(c, 11, 12, 4, 4, USAGI_EYE)
    fill_rect(c, 17, 12, 4, 4, USAGI_EYE)
    set_pixel(c, 12, 12, (255, 255, 255, 220))
    set_pixel(c, 18, 12, (255, 255, 255, 220))
    # 腮紅
    fill_rect(c, 8, 15, 4, 3, USAGI_BLUSH)
    fill_rect(c, 20, 15, 4, 3, USAGI_BLUSH)
    # 嘴
    set_pixel(c, 15, 17, USAGI_EYE)
    set_pixel(c, 16, 18, USAGI_EYE)
    set_pixel(c, 17, 17, USAGI_EYE)
    # 手腳
    fill_rect(c, 7, 19, 3, 6, USAGI_BODY)
    fill_rect(c, 22, 19, 3, 6, USAGI_BODY)
    fill_rect(c, 10, 27, 5, 3, USAGI_BODY)
    fill_rect(c, 17, 27, 5, 3, USAGI_BODY)
    return c


def draw_usagi_attack() -> list[list[tuple]]:
    c = draw_usagi_idle()
    fill_rect(c, 22, 13, 3, 7, USAGI_BODY)
    fill_rect(c, 24, 11, 5, 5, USAGI_BODY)
    # 黃色能量
    set_pixel(c, 28, 10, (255, 220, 50, 220))
    set_pixel(c, 29, 9, (255, 220, 50, 200))
    set_pixel(c, 30, 8, (255, 220, 50, 150))
    return c


def draw_usagi_bigwin() -> list[list[tuple]]:
    c = draw_usagi_idle()
    fill_rect(c, 7, 13, 3, 7, USAGI_BODY)
    fill_rect(c, 22, 13, 3, 7, USAGI_BODY)
    fill_rect(c, 4, 11, 5, 5, USAGI_BODY)
    fill_rect(c, 24, 11, 5, 5, USAGI_BODY)
    fill_rect(c, 13, 16, 6, 4, (200, 120, 80, 255))
    # 金色星星
    set_pixel(c, 2, 5, (255, 200, 0, 255))
    set_pixel(c, 3, 4, (255, 200, 0, 255))
    set_pixel(c, 4, 5, (255, 200, 0, 255))
    set_pixel(c, 28, 5, (255, 200, 0, 255))
    set_pixel(c, 29, 4, (255, 200, 0, 255))
    set_pixel(c, 30, 5, (255, 200, 0, 255))
    return c


def main():
    print("🎨 生成角色像素圖...")
    print(f"   輸出目錄：{OUTPUT_DIR}")
    print()

    sprites = [
        ("chiikawa_idle.png",    draw_chiikawa_idle()),
        ("chiikawa_attack.png",  draw_chiikawa_attack()),
        ("chiikawa_bigwin.png",  draw_chiikawa_bigwin()),
        ("hachiware_idle.png",   draw_hachiware_idle()),
        ("hachiware_attack.png", draw_hachiware_attack()),
        ("hachiware_bigwin.png", draw_hachiware_bigwin()),
        ("usagi_idle.png",       draw_usagi_idle()),
        ("usagi_attack.png",     draw_usagi_attack()),
        ("usagi_bigwin.png",     draw_usagi_bigwin()),
    ]

    for filename, pixels in sprites:
        write_png(filename, pixels)

    print()
    print(f"✅ 完成！共生成 {len(sprites)} 個角色精靈圖")
    print("   每個角色：idle / attack / bigwin 三狀態")


if __name__ == "__main__":
    main()
