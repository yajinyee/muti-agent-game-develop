"""
generate_characters_v1.py — 角色像素圖生成（重建版）
character-pixel-agent 負責維護

三個角色各生成 idle 靜態圖，96x96 px
"""
from PIL import Image, ImageDraw
import os

OUT = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\characters"
os.makedirs(OUT, exist_ok=True)

SIZE = 96

def new_img():
    return Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))

def circle(draw, cx, cy, r, fill, outline=None):
    draw.ellipse([cx-r, cy-r, cx+r, cy+r], fill=fill, outline=outline)

def save(img, name):
    path = os.path.join(OUT, name)
    img.save(path)
    pixels = sum(1 for p in img.getdata() if p[3] > 0)
    pct = pixels / (SIZE*SIZE) * 100
    print(f"  {name}: {pct:.1f}%")

# 官方顏色
WHITE = (255, 255, 247, 255)
OUTLINE = (41, 42, 43, 255)
BLUSH = (239, 165, 201, 255)
HACHIWARE_STRIPE = (51, 112, 192, 255)
USAGI_EYE = (255, 91, 86, 255)

def draw_eye(draw, cx, cy, color=(0,0,0,255)):
    """3x3 眼白 + 2x2 瞳孔 + 高光"""
    draw.ellipse([cx-3, cy-3, cx+3, cy+3], fill=(255,255,255,255))
    draw.ellipse([cx-2, cy-2, cx+2, cy+2], fill=color)
    draw.ellipse([cx-1, cy-2, cx, cy-1], fill=(255,255,255,255))

def draw_smile(draw, cx, cy):
    """3像素 V 形嘴巴"""
    draw.point([cx-1, cy], fill=OUTLINE)
    draw.point([cx, cy+1], fill=OUTLINE)
    draw.point([cx+1, cy], fill=OUTLINE)

# ── 吉伊卡哇 ─────────────────────────────────────────────────
def gen_chiikawa():
    img = new_img(); d = ImageDraw.Draw(img)
    cx, cy = 48, 48

    # 耳朵（圓形）
    circle(d, cx-18, cy-22, 10, WHITE, OUTLINE)
    circle(d, cx+18, cy-22, 10, WHITE, OUTLINE)

    # 頭（大圓）
    circle(d, cx, cy-8, 26, WHITE, OUTLINE)

    # 身體（橢圓）
    d.ellipse([cx-16, cy+14, cx+16, cy+38], fill=WHITE, outline=OUTLINE)

    # 手臂
    circle(d, cx-20, cy+20, 7, WHITE, OUTLINE)
    circle(d, cx+20, cy+20, 7, WHITE, OUTLINE)

    # 腳
    d.ellipse([cx-14, cy+34, cx-2, cy+44], fill=WHITE, outline=OUTLINE)
    d.ellipse([cx+2, cy+34, cx+14, cy+44], fill=WHITE, outline=OUTLINE)

    # 眼睛
    draw_eye(d, cx-8, cy-10)
    draw_eye(d, cx+8, cy-10)

    # 腮紅
    circle(d, cx-14, cy-4, 5, (*BLUSH[:3], 180))
    circle(d, cx+14, cy-4, 5, (*BLUSH[:3], 180))

    # 嘴巴
    draw_smile(d, cx, cy+2)

    save(img, "chiikawa_idle.png")

# ── 小八 ──────────────────────────────────────────────────────
def gen_hachiware():
    img = new_img(); d = ImageDraw.Draw(img)
    cx, cy = 48, 48

    # 耳朵（尖耳）
    d.polygon([(cx-26, cy-18), (cx-18, cy-36), (cx-10, cy-18)], fill=WHITE, outline=OUTLINE)
    d.polygon([(cx+10, cy-18), (cx+18, cy-36), (cx+26, cy-18)], fill=WHITE, outline=OUTLINE)

    # 頭
    circle(d, cx, cy-8, 26, WHITE, OUTLINE)

    # 藍色條紋（小八特徵）
    d.line([(cx-10, cy-20), (cx-10, cy+4)], fill=HACHIWARE_STRIPE, width=4)
    d.line([(cx+10, cy-20), (cx+10, cy+4)], fill=HACHIWARE_STRIPE, width=4)

    # 身體
    d.ellipse([cx-16, cy+14, cx+16, cy+38], fill=WHITE, outline=OUTLINE)

    # 手臂
    circle(d, cx-20, cy+20, 7, WHITE, OUTLINE)
    circle(d, cx+20, cy+20, 7, WHITE, OUTLINE)

    # 腳
    d.ellipse([cx-14, cy+34, cx-2, cy+44], fill=WHITE, outline=OUTLINE)
    d.ellipse([cx+2, cy+34, cx+14, cy+44], fill=WHITE, outline=OUTLINE)

    # 眼睛
    draw_eye(d, cx-8, cy-10)
    draw_eye(d, cx+8, cy-10)

    # 嘴巴
    draw_smile(d, cx, cy+2)

    save(img, "hachiware_idle.png")

# ── 烏薩奇 ────────────────────────────────────────────────────
def gen_usagi():
    img = new_img(); d = ImageDraw.Draw(img)
    cx, cy = 48, 52

    # 長耳朵（烏薩奇特徵）
    d.ellipse([cx-22, cy-52, cx-10, cy-16], fill=WHITE, outline=OUTLINE)
    d.ellipse([cx+10, cy-52, cx+22, cy-16], fill=WHITE, outline=OUTLINE)
    # 耳朵內側（粉紅）
    d.ellipse([cx-20, cy-50, cx-12, cy-20], fill=(*BLUSH[:3], 200))
    d.ellipse([cx+12, cy-50, cx+20, cy-20], fill=(*BLUSH[:3], 200))

    # 頭
    circle(d, cx, cy-8, 24, WHITE, OUTLINE)

    # 身體
    d.ellipse([cx-14, cy+12, cx+14, cy+34], fill=WHITE, outline=OUTLINE)

    # 手臂
    circle(d, cx-18, cy+18, 6, WHITE, OUTLINE)
    circle(d, cx+18, cy+18, 6, WHITE, OUTLINE)

    # 腳
    d.ellipse([cx-12, cy+30, cx-2, cy+40], fill=WHITE, outline=OUTLINE)
    d.ellipse([cx+2, cy+30, cx+12, cy+40], fill=WHITE, outline=OUTLINE)

    # 眼睛（紅色，烏薩奇特徵）
    draw_eye(d, cx-7, cy-10, USAGI_EYE)
    draw_eye(d, cx+7, cy-10, USAGI_EYE)

    # 嘴巴
    draw_smile(d, cx, cy+2)

    save(img, "usagi_idle.png")

print("生成角色像素圖...")
gen_chiikawa()
gen_hachiware()
gen_usagi()
print(f"\n✅ 完成！輸出到 {OUT}")
