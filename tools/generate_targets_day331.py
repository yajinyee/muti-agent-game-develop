"""
generate_targets_day331.py — DAY-331 T239-T243 精靈圖生成
業界依據：BGaming Shark & Spark Hold & Win（2026-05-30 最新）
"""
import os
from PIL import Image, ImageDraw
import math

OUT_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
SIZE = 64

def make_img():
    return Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))

def draw_rays(draw, cx, cy, count, r_inner, r_outer, color, width=1):
    for i in range(count):
        angle = math.radians(i * 360 / count)
        x1 = cx + r_inner * math.cos(angle)
        y1 = cy + r_inner * math.sin(angle)
        x2 = cx + r_outer * math.cos(angle)
        y2 = cy + r_outer * math.sin(angle)
        draw.line([(x1, y1), (x2, y2)], fill=color, width=width)

def draw_ring(draw, cx, cy, r, color, width=2):
    draw.ellipse([cx-r, cy-r, cx+r, cy+r], outline=color, width=width)

def fill_circle(img, cx, cy, r, color):
    draw = ImageDraw.Draw(img)
    draw.ellipse([cx-r, cy-r, cx+r, cy+r], fill=color)

def fill_ellipse(img, cx, cy, rx, ry, color):
    draw = ImageDraw.Draw(img)
    draw.ellipse([cx-rx, cy-ry, cx+rx, cy+ry], fill=color)

# ── T239 幸運鯊魚閃電魚（Shark & Spark）──────────────────────
def gen_t239():
    img = make_img()
    draw = ImageDraw.Draw(img)
    cx, cy = 32, 32

    # 深海藍魚身（橢圓）
    fill_ellipse(img, cx, cy, 22, 16, (0, 80, 160, 255))
    fill_ellipse(img, cx, cy, 18, 12, (0, 120, 200, 255))

    # 鯊魚背鰭（三角形）
    draw.polygon([(cx, cy-16), (cx-8, cy-8), (cx+8, cy-8)], fill=(0, 60, 130, 255))

    # 閃電符號（黃色）
    lightning = [(cx-4, cy-8), (cx+2, cy-2), (cx-2, cy-2), (cx+4, cy+8), (cx-2, cy+2), (cx+2, cy+2)]
    draw.polygon(lightning, fill=(255, 220, 0, 255))

    # 珍珠光點（白色小圓）
    for angle_deg in range(0, 360, 45):
        angle = math.radians(angle_deg)
        px = int(cx + 20 * math.cos(angle))
        py = int(cy + 14 * math.sin(angle))
        draw.ellipse([px-2, py-2, px+2, py+2], fill=(255, 255, 255, 200))

    # 8 道光芒
    draw_rays(draw, cx, cy, 8, 24, 30, (0, 200, 255, 180), 1)

    # 輪廓
    draw_ring(draw, cx, cy, 22, (0, 40, 100, 255), 2)

    img.save(os.path.join(OUT_DIR, "T239_shark_spark.png"))
    print(f"T239 saved ({sum(1 for p in img.getdata() if p[3] > 0)} px)")

# ── T240 幸運冬季冰釣魚（Winter Ice Fishing）─────────────────
def gen_t240():
    img = make_img()
    draw = ImageDraw.Draw(img)
    cx, cy = 32, 32

    # 冰藍魚身
    fill_ellipse(img, cx, cy, 20, 14, (100, 180, 240, 255))
    fill_ellipse(img, cx, cy, 16, 10, (150, 210, 255, 255))

    # 雪花符號（6 道）
    for i in range(6):
        angle = math.radians(i * 60)
        x1 = cx + 4 * math.cos(angle)
        y1 = cy + 4 * math.sin(angle)
        x2 = cx + 10 * math.cos(angle)
        y2 = cy + 10 * math.sin(angle)
        draw.line([(x1, y1), (x2, y2)], fill=(255, 255, 255, 255), width=2)
        # 小分支
        for branch_angle in [-30, 30]:
            ba = math.radians(i * 60 + branch_angle)
            bx1 = cx + 7 * math.cos(angle)
            by1 = cy + 7 * math.sin(angle)
            bx2 = bx1 + 3 * math.cos(ba)
            by2 = by1 + 3 * math.sin(ba)
            draw.line([(bx1, by1), (bx2, by2)], fill=(255, 255, 255, 200), width=1)

    # 冰晶光環
    draw_ring(draw, cx, cy, 22, (200, 240, 255, 200), 2)
    draw_ring(draw, cx, cy, 26, (150, 200, 255, 120), 1)

    # 53格輪盤提示（小圓點）
    for i in range(8):
        angle = math.radians(i * 45)
        px = int(cx + 28 * math.cos(angle))
        py = int(cy + 28 * math.sin(angle))
        if 0 <= px < SIZE and 0 <= py < SIZE:
            draw.ellipse([px-1, py-1, px+1, py+1], fill=(255, 255, 200, 200))

    img.save(os.path.join(OUT_DIR, "T240_winter_ice.png"))
    print(f"T240 saved ({sum(1 for p in img.getdata() if p[3] > 0)} px)")

# ── T241 幸運大西洋狂潮魚（Big Atlantis Frenzy）──────────────
def gen_t241():
    img = make_img()
    draw = ImageDraw.Draw(img)
    cx, cy = 32, 32

    # 亞特蘭提斯藍魚身
    fill_ellipse(img, cx, cy, 21, 15, (20, 80, 200, 255))
    fill_ellipse(img, cx, cy, 17, 11, (40, 120, 240, 255))

    # 波浪紋路（3 條）
    for i, y_off in enumerate([-4, 0, 4]):
        pts = []
        for x in range(cx-14, cx+15, 2):
            y = cy + y_off + int(2 * math.sin((x - cx) * 0.5))
            pts.append((x, y))
        if len(pts) >= 2:
            draw.line(pts, fill=(100, 200, 255, 200), width=1)

    # Fish 符號（小魚形）
    draw.ellipse([cx-5, cy-3, cx+3, cy+3], fill=(255, 220, 100, 255))
    draw.polygon([(cx+3, cy), (cx+7, cy-3), (cx+7, cy+3)], fill=(255, 220, 100, 255))

    # 7 波連鎖光芒
    draw_rays(draw, cx, cy, 7, 22, 29, (0, 180, 255, 180), 1)

    # 輪廓
    draw_ring(draw, cx, cy, 21, (0, 50, 150, 255), 2)
    draw_ring(draw, cx, cy, 25, (0, 100, 200, 120), 1)

    img.save(os.path.join(OUT_DIR, "T241_atlantis_frenzy.png"))
    print(f"T241 saved ({sum(1 for p in img.getdata() if p[3] > 0)} px)")

# ── T242 幸運釣魚時間魚（Fishing Time Wheel）─────────────────
def gen_t242():
    img = make_img()
    draw = ImageDraw.Draw(img)
    cx, cy = 32, 32

    # 金橙色魚身
    fill_ellipse(img, cx, cy, 20, 14, (200, 120, 0, 255))
    fill_ellipse(img, cx, cy, 16, 10, (255, 180, 0, 255))

    # 輪盤（5 色扇形）
    wheel_colors = [
        (255, 50, 50, 200),   # 紅
        (255, 200, 0, 200),   # 黃
        (0, 200, 100, 200),   # 綠
        (0, 100, 255, 200),   # 藍
        (200, 0, 255, 200),   # 紫
    ]
    for i, wc in enumerate(wheel_colors):
        start_angle = i * 72 - 90
        draw.pieslice([cx-10, cy-10, cx+10, cy+10], start=start_angle, end=start_angle+72, fill=wc)

    # 輪盤中心
    draw.ellipse([cx-3, cy-3, cx+3, cy+3], fill=(255, 255, 255, 255))

    # 5 次旋轉光芒
    draw_rays(draw, cx, cy, 5, 22, 30, (255, 200, 0, 200), 2)

    # 輪廓
    draw_ring(draw, cx, cy, 20, (150, 80, 0, 255), 2)
    draw_ring(draw, cx, cy, 24, (255, 150, 0, 150), 1)

    img.save(os.path.join(OUT_DIR, "T242_fishing_time_wheel.png"))
    print(f"T242 saved ({sum(1 for p in img.getdata() if p[3] > 0)} px)")

# ── T243 幸運終極鯊魚魚（Ultimate Shark，里程碑 ×53.0）────────
def gen_t243():
    img = make_img()
    draw = ImageDraw.Draw(img)
    cx, cy = 32, 32

    # 超大型橙紅魚身
    fill_ellipse(img, cx, cy, 24, 18, (180, 40, 0, 255))
    fill_ellipse(img, cx, cy, 20, 14, (255, 80, 20, 255))

    # 鯊魚背鰭（大）
    draw.polygon([(cx, cy-18), (cx-10, cy-8), (cx+10, cy-8)], fill=(140, 20, 0, 255))

    # 鯊魚牙齒（下方）
    for i in range(-2, 3):
        tx = cx + i * 5
        draw.polygon([(tx-2, cy+12), (tx, cy+18), (tx+2, cy+12)], fill=(255, 255, 255, 220))

    # 14 道光芒（里程碑）
    draw_rays(draw, cx, cy, 14, 25, 31, (255, 150, 0, 200), 1)

    # 多層光環
    draw_ring(draw, cx, cy, 24, (200, 50, 0, 255), 2)
    draw_ring(draw, cx, cy, 28, (255, 100, 0, 180), 2)
    draw_ring(draw, cx, cy, 30, (255, 200, 0, 120), 1)

    # 里程碑符文（中心 S 形）
    draw.text((cx-4, cy-5), "S", fill=(255, 255, 0, 255))

    img.save(os.path.join(OUT_DIR, "T243_ultimate_shark.png"))
    print(f"T243 saved ({sum(1 for p in img.getdata() if p[3] > 0)} px)")

if __name__ == "__main__":
    os.makedirs(OUT_DIR, exist_ok=True)
    gen_t239()
    gen_t240()
    gen_t241()
    gen_t242()
    gen_t243()
    print("DAY-331 T239-T243 精靈圖生成完成！")
