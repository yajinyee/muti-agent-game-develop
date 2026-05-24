"""
generate_targets_v4.py — 目標物像素圖生成（重建版）
target-pixel-agent 負責維護

每個目標物 64x64 px，透明背景，清楚可辨識的剪影
"""
from PIL import Image, ImageDraw
import os

OUT = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
os.makedirs(OUT, exist_ok=True)

SIZE = 64

def new_img():
    return Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))

def circle(draw, cx, cy, r, color, outline=None):
    draw.ellipse([cx-r, cy-r, cx+r, cy+r], fill=color, outline=outline)

def rect(draw, x, y, w, h, color, outline=None):
    draw.rectangle([x, y, x+w, y+h], fill=color, outline=outline)

def save(img, name):
    path = os.path.join(OUT, name)
    img.save(path)
    pixels = sum(1 for p in img.getdata() if p[3] > 0)
    pct = pixels / (SIZE*SIZE) * 100
    print(f"  {name}: {pct:.1f}% non-transparent")

# ── T001 像素雜草（綠色，靜止）────────────────────────────────
def gen_T001():
    img = new_img(); d = ImageDraw.Draw(img)
    # 草莖
    for i in range(3):
        x = 16 + i * 12
        d.line([(x, 48), (x-4, 20)], fill=(30, 140, 30, 255), width=3)
        d.line([(x, 48), (x+4, 22)], fill=(50, 170, 50, 255), width=2)
        # 草葉
        d.ellipse([x-8, 14, x+4, 26], fill=(60, 200, 60, 255))
    # 地面
    rect(d, 8, 48, 48, 8, (80, 60, 30, 255))
    save(img, "T001_grass.png")

# ── T002 綠色小蟲 ─────────────────────────────────────────────
def gen_T002():
    img = new_img(); d = ImageDraw.Draw(img)
    # 身體
    circle(d, 32, 36, 16, (50, 200, 50, 255), (20, 120, 20, 255))
    # 頭
    circle(d, 32, 20, 10, (70, 220, 70, 255), (20, 120, 20, 255))
    # 眼睛
    circle(d, 28, 18, 3, (0, 0, 0, 255))
    circle(d, 36, 18, 3, (0, 0, 0, 255))
    circle(d, 27, 17, 1, (255, 255, 255, 255))
    circle(d, 35, 17, 1, (255, 255, 255, 255))
    # 觸角
    d.line([(28, 12), (22, 4)], fill=(20, 120, 20, 255), width=2)
    d.line([(36, 12), (42, 4)], fill=(20, 120, 20, 255), width=2)
    # 腳
    for i in range(3):
        y = 32 + i * 5
        d.line([(16, y), (8, y+4)], fill=(20, 120, 20, 255), width=2)
        d.line([(48, y), (56, y+4)], fill=(20, 120, 20, 255), width=2)
    save(img, "T002_bug_g.png")

# ── T003 紅色小蟲 ─────────────────────────────────────────────
def gen_T003():
    img = new_img(); d = ImageDraw.Draw(img)
    circle(d, 32, 36, 16, (220, 50, 50, 255), (140, 20, 20, 255))
    circle(d, 32, 20, 10, (240, 70, 70, 255), (140, 20, 20, 255))
    circle(d, 28, 18, 3, (0, 0, 0, 255))
    circle(d, 36, 18, 3, (0, 0, 0, 255))
    circle(d, 27, 17, 1, (255, 255, 255, 255))
    circle(d, 35, 17, 1, (255, 255, 255, 255))
    d.line([(28, 12), (22, 4)], fill=(140, 20, 20, 255), width=2)
    d.line([(36, 12), (42, 4)], fill=(140, 20, 20, 255), width=2)
    for i in range(3):
        y = 32 + i * 5
        d.line([(16, y), (8, y+4)], fill=(140, 20, 20, 255), width=2)
        d.line([(48, y), (56, y+4)], fill=(140, 20, 20, 255), width=2)
    save(img, "T003_bug_r.png")

# ── T004 藍色小蟲 ─────────────────────────────────────────────
def gen_T004():
    img = new_img(); d = ImageDraw.Draw(img)
    circle(d, 32, 36, 16, (50, 100, 220, 255), (20, 50, 140, 255))
    circle(d, 32, 20, 10, (70, 130, 240, 255), (20, 50, 140, 255))
    circle(d, 28, 18, 3, (0, 0, 0, 255))
    circle(d, 36, 18, 3, (0, 0, 0, 255))
    circle(d, 27, 17, 1, (255, 255, 255, 255))
    circle(d, 35, 17, 1, (255, 255, 255, 255))
    d.line([(28, 12), (22, 4)], fill=(20, 50, 140, 255), width=2)
    d.line([(36, 12), (42, 4)], fill=(20, 50, 140, 255), width=2)
    for i in range(3):
        y = 32 + i * 5
        d.line([(16, y), (8, y+4)], fill=(20, 50, 140, 255), width=2)
        d.line([(48, y), (56, y+4)], fill=(20, 50, 140, 255), width=2)
    save(img, "T004_bug_b.png")

# ── T005 會走路的布丁 ─────────────────────────────────────────
def gen_T005():
    img = new_img(); d = ImageDraw.Draw(img)
    # 布丁身體（圓形）
    circle(d, 32, 34, 20, (255, 220, 80, 255), (180, 140, 20, 255))
    # 布丁頂部（焦糖）
    circle(d, 32, 18, 10, (200, 120, 30, 255), (140, 80, 10, 255))
    # 眼睛
    circle(d, 26, 30, 4, (0, 0, 0, 255))
    circle(d, 38, 30, 4, (0, 0, 0, 255))
    circle(d, 25, 29, 2, (255, 255, 255, 255))
    circle(d, 37, 29, 2, (255, 255, 255, 255))
    # 嘴巴
    d.arc([26, 36, 38, 44], 0, 180, fill=(180, 100, 20, 255), width=2)
    # 腳
    d.ellipse([14, 50, 26, 58], fill=(255, 200, 60, 255))
    d.ellipse([38, 50, 50, 58], fill=(255, 200, 60, 255))
    save(img, "T005_pudding.png")

# ── T006 巨大蘑菇 ─────────────────────────────────────────────
def gen_T006():
    img = new_img(); d = ImageDraw.Draw(img)
    # 蘑菇柄
    rect(d, 22, 38, 20, 20, (220, 200, 180, 255), (160, 140, 120, 255))
    # 蘑菇帽
    d.ellipse([8, 16, 56, 44], fill=(180, 60, 30, 255), outline=(120, 30, 10, 255))
    # 白色斑點
    circle(d, 22, 26, 5, (255, 255, 255, 200))
    circle(d, 38, 22, 4, (255, 255, 255, 200))
    circle(d, 46, 32, 3, (255, 255, 255, 200))
    # 眼睛
    circle(d, 28, 34, 3, (0, 0, 0, 255))
    circle(d, 36, 34, 3, (0, 0, 0, 255))
    save(img, "T006_mushroom.png")

# ── T101 擬態型怪物 ───────────────────────────────────────────
def gen_T101():
    img = new_img(); d = ImageDraw.Draw(img)
    # 不規則形狀（擬態感）
    d.polygon([(32,8),(52,20),(56,40),(44,56),(20,56),(8,40),(12,20)],
              fill=(120, 120, 140, 255), outline=(60, 60, 80, 255))
    # 問號（擬態標誌）
    d.text((26, 22), "?", fill=(200, 200, 220, 255))
    # 眼睛（詭異）
    circle(d, 24, 32, 5, (200, 50, 50, 255))
    circle(d, 40, 32, 5, (200, 50, 50, 255))
    circle(d, 24, 32, 2, (0, 0, 0, 255))
    circle(d, 40, 32, 2, (0, 0, 0, 255))
    save(img, "T101_mimic.png")

# ── T102 寶箱怪 ───────────────────────────────────────────────
def gen_T102():
    img = new_img(); d = ImageDraw.Draw(img)
    # 箱子主體
    rect(d, 10, 28, 44, 28, (180, 130, 40, 255), (100, 70, 10, 255))
    # 箱蓋
    rect(d, 10, 18, 44, 14, (200, 150, 50, 255), (100, 70, 10, 255))
    # 金屬扣
    rect(d, 26, 22, 12, 10, (220, 180, 60, 255), (140, 100, 20, 255))
    # 眼睛（在箱蓋上）
    circle(d, 24, 24, 4, (255, 255, 255, 255))
    circle(d, 40, 24, 4, (255, 255, 255, 255))
    circle(d, 24, 24, 2, (0, 0, 0, 255))
    circle(d, 40, 24, 2, (0, 0, 0, 255))
    # 牙齒
    for i in range(5):
        x = 14 + i * 8
        rect(d, x, 28, 5, 6, (255, 255, 255, 255))
    save(img, "T102_chest.png")

# ── T103 流星 ─────────────────────────────────────────────────
def gen_T103():
    img = new_img(); d = ImageDraw.Draw(img)
    # 流星主體（橢圓）
    d.ellipse([20, 24, 56, 44], fill=(255, 255, 220, 255), outline=(200, 200, 150, 255))
    # 尾跡
    for i in range(4):
        alpha = 200 - i * 40
        d.ellipse([4-i*2, 28+i, 24-i*2, 40+i],
                  fill=(255, 240, 180, alpha))
    # 光芒
    d.line([(38, 20), (38, 8)], fill=(255, 255, 200, 200), width=2)
    d.line([(50, 28), (60, 22)], fill=(255, 255, 200, 200), width=2)
    save(img, "T103_meteor.png")

# ── T104 金色雜草 ─────────────────────────────────────────────
def gen_T104():
    img = new_img(); d = ImageDraw.Draw(img)
    # 金色草莖（更粗更亮）
    for i in range(3):
        x = 14 + i * 14
        d.line([(x, 52), (x-6, 16)], fill=(200, 160, 0, 255), width=4)
        d.line([(x, 52), (x+6, 18)], fill=(220, 180, 20, 255), width=3)
        # 金色草葉
        d.ellipse([x-10, 10, x+6, 24], fill=(255, 200, 0, 255), outline=(180, 140, 0, 255))
    # 地面
    rect(d, 6, 52, 52, 8, (120, 90, 20, 255))
    # 金色光暈
    for r in [28, 32, 36]:
        d.ellipse([32-r, 32-r, 32+r, 32+r], outline=(255, 220, 0, 60))
    save(img, "T104_gold_grass.png")

# ── T105 巨大金幣魚 ───────────────────────────────────────────
def gen_T105():
    img = new_img(); d = ImageDraw.Draw(img)
    # 魚身（橢圓）
    d.ellipse([8, 20, 52, 48], fill=(255, 200, 30, 255), outline=(180, 130, 0, 255))
    # 魚鱗
    for row in range(2):
        for col in range(3):
            x = 16 + col * 12
            y = 26 + row * 10
            d.ellipse([x, y, x+8, y+6], outline=(200, 150, 0, 180))
    # 魚尾
    d.polygon([(52, 34), (62, 22), (62, 46)], fill=(220, 170, 0, 255))
    # 眼睛
    circle(d, 18, 30, 5, (255, 255, 255, 255))
    circle(d, 18, 30, 3, (0, 0, 0, 255))
    circle(d, 17, 29, 1, (255, 255, 255, 255))
    # ¥ 符號
    d.text((26, 26), "¥", fill=(180, 120, 0, 255))
    save(img, "T105_coin_fish.png")

# ── B001 BOSS（那個孩子）─────────────────────────────────────
def gen_B001():
    img = new_img(); d = ImageDraw.Draw(img)
    # 大圓形身體
    circle(d, 32, 32, 28, (60, 20, 80, 255), (30, 0, 50, 255))
    # 眼睛（大而詭異）
    circle(d, 22, 26, 8, (255, 50, 50, 255))
    circle(d, 42, 26, 8, (255, 50, 50, 255))
    circle(d, 22, 26, 4, (0, 0, 0, 255))
    circle(d, 42, 26, 4, (0, 0, 0, 255))
    circle(d, 20, 24, 2, (255, 255, 255, 255))
    circle(d, 40, 24, 2, (255, 255, 255, 255))
    # 嘴巴（邪惡微笑）
    d.arc([18, 36, 46, 52], 0, 180, fill=(200, 0, 0, 255), width=3)
    # 牙齒
    for i in range(4):
        x = 22 + i * 6
        d.polygon([(x, 44), (x+3, 44), (x+1, 50)], fill=(255, 255, 255, 255))
    # 光環
    for r in [30, 34]:
        d.ellipse([32-r, 32-r, 32+r, 32+r], outline=(150, 0, 200, 100))
    save(img, "B001_boss.png")

# ── 生成所有目標物 ────────────────────────────────────────────
print("生成目標物像素圖...")
gen_T001()
gen_T002()
gen_T003()
gen_T004()
gen_T005()
gen_T006()
gen_T101()
gen_T102()
gen_T103()
gen_T104()
gen_T105()
gen_B001()
print(f"\n✅ 完成！輸出到 {OUT}")
