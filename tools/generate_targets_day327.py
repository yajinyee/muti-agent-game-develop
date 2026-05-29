"""
generate_targets_day327.py — DAY-327 T224-T228 精靈圖生成
T224 幸運黃金鍋魚：深金色魚身 + 黃金鍋符號 + 12格網格 + 光芒
T225 幸運瀑布鎖定魚：深藍色魚身 + 瀑布連鎖符號 + 鎖定圖示 + 光芒
T226 幸運傳說覺醒魚：火橙色魚身 + 龍形符號 + 8次獎勵指示 + 光芒
T227 幸運崩潰收割魚：深紅色魚身 + 崩潰箭頭 + 收割鐮刀 + 光芒
T228 幸運宇宙大融合魚：宇宙紫色超大型魚身 + 4 Phase 符號 + 45道光芒 + 六層光環
"""
import os
import math
import random
from PIL import Image, ImageDraw

OUT_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
SIZE = 64
rng = random.Random(327)

def fill_circle(draw, cx, cy, r, color):
    for y in range(cy - r, cy + r + 1):
        for x in range(cx - r, cx + r + 1):
            if (x - cx) ** 2 + (y - cy) ** 2 <= r * r:
                if 0 <= x < SIZE and 0 <= y < SIZE:
                    draw.point((x, y), fill=color)

def fill_circle_shaded(draw, cx, cy, r, base_color):
    br, bg, bb = base_color[:3]
    for y in range(cy - r, cy + r + 1):
        for x in range(cx - r, cx + r + 1):
            dist = math.sqrt((x - cx) ** 2 + (y - cy) ** 2)
            if dist <= r:
                if 0 <= x < SIZE and 0 <= y < SIZE:
                    shade = 1.0 - (dist / r) * 0.35
                    if x < cx and y < cy:
                        shade = min(1.0, shade + 0.15)
                    elif x > cx and y > cy:
                        shade = max(0.5, shade - 0.15)
                    draw.point((x, y), fill=(
                        min(255, int(br * shade)),
                        min(255, int(bg * shade)),
                        min(255, int(bb * shade)),
                        255
                    ))

def draw_rays(draw, cx, cy, n_rays, r_inner, r_outer, color):
    for i in range(n_rays):
        angle = 2 * math.pi * i / n_rays
        x1 = int(cx + r_inner * math.cos(angle))
        y1 = int(cy + r_inner * math.sin(angle))
        x2 = int(cx + r_outer * math.cos(angle))
        y2 = int(cy + r_outer * math.sin(angle))
        draw.line([(x1, y1), (x2, y2)], fill=color, width=2)

def draw_ring(draw, cx, cy, r, color, width=2):
    for angle_deg in range(0, 360, 2):
        angle = math.radians(angle_deg)
        for w in range(width):
            rx = r + w - width // 2
            x = int(cx + rx * math.cos(angle))
            y = int(cy + rx * math.sin(angle))
            if 0 <= x < SIZE and 0 <= y < SIZE:
                draw.point((x, y), fill=color)

def draw_ellipse_shaded(draw, cx, cy, rx, ry, base_color):
    br, bg, bb = base_color[:3]
    for y in range(cy - ry, cy + ry + 1):
        for x in range(cx - rx, cx + rx + 1):
            if (x - cx) ** 2 / max(1, rx * rx) + (y - cy) ** 2 / max(1, ry * ry) <= 1.0:
                if 0 <= x < SIZE and 0 <= y < SIZE:
                    dx = (x - cx) / max(1, rx)
                    dy = (y - cy) / max(1, ry)
                    shade = 1.0 - (abs(dx) + abs(dy)) * 0.2
                    if dx < 0 and dy < 0:
                        shade = min(1.0, shade + 0.15)
                    draw.point((x, y), fill=(
                        min(255, int(br * shade)),
                        min(255, int(bg * shade)),
                        min(255, int(bb * shade)),
                        255
                    ))

def add_eye(draw, ex, ey):
    fill_circle(draw, ex, ey, 3, (255, 255, 255, 255))
    fill_circle(draw, ex, ey, 2, (20, 20, 20, 255))
    draw.point((ex - 1, ey - 1), fill=(255, 255, 255, 255))

def save_with_import(img, name):
    path = os.path.join(OUT_DIR, name)
    img.save(path)
    import_path = path + ".import"
    if not os.path.exists(import_path):
        with open(import_path, "w") as f:
            f.write(f"""[remap]
importer="texture"
type="CompressedTexture2D"
uid="uid://dag{abs(hash(name)) % 999999999}"
path="res://.godot/imported/{name}-{abs(hash(name)) % 999999999}.ctex"
metadata={{"vram_texture": false}}

[deps]
source_file="res://assets/sprites/targets/{name}"
dest_files=["res://.godot/imported/{name}-{abs(hash(name)) % 999999999}.ctex"]

[params]
compress/mode=0
compress/high_quality=false
compress/lossy_quality=0.7
compress/normal_map=0
compress/channel_pack=0
mipmaps/generate=false
mipmaps/limit=-1
roughness/mode=0
roughness/src_normal=""
process/fix_alpha_border=true
process/premult_alpha=false
process/normal_map_invert_y=false
process/hdr_as_srgb=false
process/hdr_clamp_exposure=false
process/size_limit=0
detect_3d/compress_to=1
svg/scale=1.0
editor/scale_with_editor_scale=false
editor/convert_colors_with_editor_theme=false
""")
    print(f"  Saved: {name}")

# ── T224 幸運黃金鍋魚 ─────────────────────────────────────────
def gen_t224():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = 32, 32

    # 深金色橢圓魚身
    draw_ellipse_shaded(draw, cx, cy, 22, 16, (200, 140, 0))

    # 黃金鍋符號（圓形鍋子）
    fill_circle_shaded(draw, cx, cy - 2, 10, (220, 160, 0))
    # 鍋口
    for x in range(cx - 10, cx + 11):
        if 0 <= x < SIZE:
            draw.point((x, cy - 12), fill=(255, 200, 0, 255))
    # 鍋把手
    draw.line([(cx - 12, cy - 4), (cx - 16, cy - 8)], fill=(180, 120, 0, 255), width=2)
    draw.line([(cx + 12, cy - 4), (cx + 16, cy - 8)], fill=(180, 120, 0, 255), width=2)

    # 12格網格（3×4）
    grid_x, grid_y = cx - 9, cy - 6
    for row in range(3):
        for col in range(4):
            gx = grid_x + col * 5
            gy = grid_y + row * 5
            draw.rectangle([(gx, gy), (gx + 3, gy + 3)], outline=(255, 220, 0, 200))

    # 光芒
    draw_rays(draw, cx, cy, 16, 24, 30, (255, 200, 0, 180))

    # 光環
    draw_ring(draw, cx, cy, 26, (255, 180, 0, 150))
    draw_ring(draw, cx, cy, 29, (255, 220, 0, 100))

    # 眼睛
    add_eye(draw, cx - 8, cy - 4)

    # 魚尾
    for i in range(8):
        tx = cx + 20 + i
        ty = cy - 4 + i
        if 0 <= tx < SIZE and 0 <= ty < SIZE:
            draw.point((tx, ty), fill=(180, 120, 0, 200))
        ty2 = cy + 4 - i
        if 0 <= tx < SIZE and 0 <= ty2 < SIZE:
            draw.point((tx, ty2), fill=(180, 120, 0, 200))

    save_with_import(img, "T224_golden_pot.png")

# ── T225 幸運瀑布鎖定魚 ───────────────────────────────────────
def gen_t225():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = 32, 32

    # 深藍色橢圓魚身
    draw_ellipse_shaded(draw, cx, cy, 22, 16, (0, 80, 180))

    # 瀑布連鎖符號（向下箭頭疊加）
    for i in range(4):
        ay = cy - 8 + i * 5
        draw.line([(cx - 5, ay), (cx, ay + 4), (cx + 5, ay)], fill=(0, 200, 255, 220), width=2)

    # 鎖定圖示（右側）
    fill_circle(draw, cx + 8, cy - 2, 4, (0, 150, 220, 200))
    draw.rectangle([(cx + 5, cy + 1), (cx + 11, cy + 6)], fill=(0, 120, 200, 220))

    # Pearl 符號（左側小圓）
    fill_circle_shaded(draw, cx - 10, cy, 4, (200, 230, 255))

    # 光芒
    draw_rays(draw, cx, cy, 14, 24, 30, (0, 180, 255, 160))

    # 光環
    draw_ring(draw, cx, cy, 26, (0, 150, 255, 150))
    draw_ring(draw, cx, cy, 29, (0, 200, 255, 100))

    # 眼睛
    add_eye(draw, cx - 8, cy - 4)

    # 魚尾
    for i in range(8):
        tx = cx + 20 + i
        ty = cy - 4 + i
        if 0 <= tx < SIZE and 0 <= ty < SIZE:
            draw.point((tx, ty), fill=(0, 60, 150, 200))
        ty2 = cy + 4 - i
        if 0 <= tx < SIZE and 0 <= ty2 < SIZE:
            draw.point((tx, ty2), fill=(0, 60, 150, 200))

    save_with_import(img, "T225_cascade_lock.png")

# ── T226 幸運傳說覺醒魚 ───────────────────────────────────────
def gen_t226():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = 32, 32

    # 火橙色橢圓魚身
    draw_ellipse_shaded(draw, cx, cy, 22, 16, (200, 60, 0))

    # 龍形符號（S形曲線）
    pts = []
    for t in range(20):
        angle = t * math.pi / 10
        px = int(cx + 6 * math.sin(angle))
        py = int(cy - 8 + t)
        pts.append((px, py))
    for i in range(len(pts) - 1):
        draw.line([pts[i], pts[i + 1]], fill=(255, 150, 0, 220), width=2)

    # 8次獎勵指示（8個小圓點）
    for i in range(8):
        angle = 2 * math.pi * i / 8
        dx = int(cx + 14 * math.cos(angle))
        dy = int(cy + 14 * math.sin(angle))
        fill_circle(draw, dx, dy, 2, (255, 200, 0, 200))

    # 光芒
    draw_rays(draw, cx, cy, 16, 24, 31, (255, 120, 0, 160))

    # 光環
    draw_ring(draw, cx, cy, 26, (255, 100, 0, 150))
    draw_ring(draw, cx, cy, 29, (255, 160, 0, 100))

    # 眼睛
    add_eye(draw, cx - 8, cy - 4)

    # 魚尾
    for i in range(8):
        tx = cx + 20 + i
        ty = cy - 4 + i
        if 0 <= tx < SIZE and 0 <= ty < SIZE:
            draw.point((tx, ty), fill=(160, 40, 0, 200))
        ty2 = cy + 4 - i
        if 0 <= tx < SIZE and 0 <= ty2 < SIZE:
            draw.point((tx, ty2), fill=(160, 40, 0, 200))

    save_with_import(img, "T226_legend_awaken.png")

# ── T227 幸運崩潰收割魚 ───────────────────────────────────────
def gen_t227():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = 32, 32

    # 深紅色橢圓魚身
    draw_ellipse_shaded(draw, cx, cy, 22, 16, (180, 20, 0))

    # 崩潰箭頭（向上急升）
    draw.line([(cx - 8, cy + 6), (cx - 2, cy - 4), (cx + 4, cy + 2), (cx + 10, cy - 8)],
              fill=(255, 80, 0, 220), width=2)
    # 箭頭頭部
    draw.polygon([(cx + 10, cy - 8), (cx + 7, cy - 4), (cx + 13, cy - 4)],
                 fill=(255, 60, 0, 220))

    # 收割鐮刀符號
    for angle_deg in range(200, 360, 10):
        angle = math.radians(angle_deg)
        x1 = int(cx - 6 + 5 * math.cos(angle))
        y1 = int(cy + 4 + 5 * math.sin(angle))
        x2 = int(cx - 6 + 7 * math.cos(angle))
        y2 = int(cy + 4 + 7 * math.sin(angle))
        draw.line([(x1, y1), (x2, y2)], fill=(200, 0, 0, 200), width=1)

    # 崩潰裂縫
    draw.line([(cx - 10, cy - 8), (cx - 4, cy - 2), (cx - 8, cy + 4)],
              fill=(255, 100, 0, 180), width=1)

    # 光芒
    draw_rays(draw, cx, cy, 12, 24, 30, (255, 50, 0, 160))

    # 光環
    draw_ring(draw, cx, cy, 26, (200, 0, 0, 150))
    draw_ring(draw, cx, cy, 29, (255, 80, 0, 100))

    # 眼睛
    add_eye(draw, cx - 8, cy - 4)

    # 魚尾
    for i in range(8):
        tx = cx + 20 + i
        ty = cy - 4 + i
        if 0 <= tx < SIZE and 0 <= ty < SIZE:
            draw.point((tx, ty), fill=(140, 10, 0, 200))
        ty2 = cy + 4 - i
        if 0 <= tx < SIZE and 0 <= ty2 < SIZE:
            draw.point((tx, ty2), fill=(140, 10, 0, 200))

    save_with_import(img, "T227_crash_harvest.png")

# ── T228 幸運宇宙大融合魚 ─────────────────────────────────────
def gen_t228():
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    cx, cy = 32, 32

    # 宇宙紫色超大型橢圓魚身
    draw_ellipse_shaded(draw, cx, cy, 24, 18, (120, 0, 180))

    # 4 Phase 符號（四象限）
    # Phase 1 (右上): 金幣
    fill_circle(draw, cx + 8, cy - 8, 4, (200, 150, 0, 200))
    # Phase 2 (左上): 瀑布
    draw.line([(cx - 12, cy - 10), (cx - 8, cy - 6), (cx - 12, cy - 2)],
              fill=(0, 150, 255, 200), width=2)
    # Phase 3 (右下): 龍
    draw.arc([(cx + 4, cy + 2), (cx + 14, cy + 12)], 0, 270, fill=(255, 100, 0, 200), width=2)
    # Phase 4 (左下): 爆炸
    for i in range(6):
        angle = math.radians(i * 60)
        x1 = int(cx - 8 + 3 * math.cos(angle))
        y1 = int(cy + 8 + 3 * math.sin(angle))
        x2 = int(cx - 8 + 6 * math.cos(angle))
        y2 = int(cy + 8 + 6 * math.sin(angle))
        draw.line([(x1, y1), (x2, y2)], fill=(255, 0, 200, 200), width=1)

    # 45道光芒（最多）
    draw_rays(draw, cx, cy, 45, 26, 31, (200, 0, 255, 120))

    # 六層光環
    for r, alpha in [(27, 180), (29, 150), (31, 120), (33, 90), (35, 60), (37, 40)]:
        draw_ring(draw, cx, cy, r, (180, 0, 255, alpha))

    # 宇宙粒子
    for _ in range(20):
        px = rng.randint(2, SIZE - 3)
        py = rng.randint(2, SIZE - 3)
        dist = math.sqrt((px - cx) ** 2 + (py - cy) ** 2)
        if dist > 20:
            draw.point((px, py), fill=(255, 100, 255, 180))

    # 眼睛
    add_eye(draw, cx - 10, cy - 5)

    # 魚尾
    for i in range(10):
        tx = cx + 22 + i
        ty = cy - 5 + i
        if 0 <= tx < SIZE and 0 <= ty < SIZE:
            draw.point((tx, ty), fill=(100, 0, 160, 200))
        ty2 = cy + 5 - i
        if 0 <= tx < SIZE and 0 <= ty2 < SIZE:
            draw.point((tx, ty2), fill=(100, 0, 160, 200))

    save_with_import(img, "T228_cosmic_fusion.png")

# ── 主程式 ────────────────────────────────────────────────────
if __name__ == "__main__":
    os.makedirs(OUT_DIR, exist_ok=True)
    print("DAY-327 精靈圖生成中...")
    gen_t224()
    gen_t225()
    gen_t226()
    gen_t227()
    gen_t228()
    print("完成！T224-T228 精靈圖已生成。")

    # 統計非透明像素
    for name in ["T224_golden_pot.png", "T225_cascade_lock.png", "T226_legend_awaken.png",
                 "T227_crash_harvest.png", "T228_cosmic_fusion.png"]:
        path = os.path.join(OUT_DIR, name)
        img = Image.open(path).convert("RGBA")
        pixels = list(img.getdata())
        non_transparent = sum(1 for p in pixels if p[3] > 10)
        total = SIZE * SIZE
        print(f"  {name}: {non_transparent}/{total} = {non_transparent/total*100:.1f}%")
