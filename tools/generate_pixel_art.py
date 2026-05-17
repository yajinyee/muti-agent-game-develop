"""
像素美術生成器
為吉伊卡哇：像素大討伐生成 16-bit 風格佔位圖
規格書 8章：美術風格 16-bit Retro Pixel Art
"""
from PIL import Image, ImageDraw, ImageFont
import os
import math

OUTPUT_BASE = r"D:\Kiro\client\chiikawa-pixel\assets\sprites"

# 調色盤（16-bit 復古風格）
PALETTE = {
    "pink":       (255, 153, 204),
    "pink_dark":  (204, 102, 153),
    "blue":       (102, 153, 255),
    "blue_dark":  (51, 102, 204),
    "yellow":     (255, 230, 51),
    "yellow_dark":(204, 179, 0),
    "green":      (102, 204, 102),
    "green_dark": (51, 153, 51),
    "red":        (255, 77, 77),
    "red_dark":   (204, 26, 26),
    "white":      (255, 255, 255),
    "black":      (0, 0, 0),
    "bg_blue":    (13, 26, 64),
    "bg_sea":     (26, 77, 153),
    "gold":       (255, 204, 0),
    "gold_dark":  (204, 153, 0),
    "purple":     (153, 51, 204),
    "brown":      (153, 102, 51),
    "gray":       (128, 128, 128),
    "transparent": (0, 0, 0, 0),
}

def new_img(w, h, bg=(0,0,0,0)):
    """建立透明背景圖"""
    return Image.new("RGBA", (w, h), bg)

def draw_pixel_border(draw, x, y, w, h, color, thickness=2):
    """畫像素邊框"""
    for t in range(thickness):
        draw.rectangle([x+t, y+t, x+w-1-t, y+h-1-t], outline=color)

def draw_pixel_char(img, char_type, size=32):
    """畫角色像素圖"""
    draw = ImageDraw.Draw(img)
    w, h = size, size
    cx, cy = w//2, h//2

    if char_type == "chiikawa":
        # 吉伊卡哇：圓頭、小耳朵、粉紅色
        c = PALETTE["pink"]
        cd = PALETTE["pink_dark"]
        # 身體（圓形）
        draw.ellipse([cx-10, cy-8, cx+10, cy+10], fill=c, outline=cd)
        # 耳朵
        draw.ellipse([cx-12, cy-14, cx-6, cy-8], fill=c, outline=cd)
        draw.ellipse([cx+6, cy-14, cx+12, cy-8], fill=c, outline=cd)
        # 眼睛（閉著，努力感）
        draw.line([cx-5, cy-2, cx-2, cy-2], fill=PALETTE["black"], width=2)
        draw.line([cx+2, cy-2, cx+5, cy-2], fill=PALETTE["black"], width=2)
        # 嘴巴
        draw.arc([cx-3, cy+1, cx+3, cy+5], 0, 180, fill=PALETTE["black"], width=1)
        # 討伐棒
        draw.line([cx+8, cy-4, cx+14, cy-12], fill=PALETTE["pink_dark"], width=3)
        draw.ellipse([cx+12, cy-14, cx+16, cy-10], fill=PALETTE["yellow"])

    elif char_type == "hachiware":
        # 小八：藍白條紋貓、精神感
        c = PALETTE["white"]
        cd = PALETTE["blue"]
        draw.ellipse([cx-10, cy-8, cx+10, cy+10], fill=c, outline=cd)
        # 條紋
        for i in range(-8, 10, 4):
            draw.line([cx+i, cy-6, cx+i, cy+8], fill=cd, width=1)
        # 耳朵（尖的）
        draw.polygon([cx-12, cy-8, cx-8, cy-16, cx-4, cy-8], fill=c, outline=cd)
        draw.polygon([cx+4, cy-8, cx+8, cy-16, cx+12, cy-8], fill=c, outline=cd)
        # 眼睛（精神）
        draw.ellipse([cx-6, cy-3, cx-2, cy+1], fill=cd)
        draw.ellipse([cx+2, cy-3, cx+6, cy+1], fill=cd)
        # 討伐棒
        draw.line([cx+8, cy-4, cx+15, cy-13], fill=PALETTE["blue_dark"], width=3)
        draw.ellipse([cx+13, cy-15, cx+17, cy-11], fill=PALETTE["blue"])

    elif char_type == "usagi":
        # 烏薩奇：長耳兔、黃色、興奮
        c = PALETTE["white"]
        cd = PALETTE["yellow_dark"]
        draw.ellipse([cx-9, cy-6, cx+9, cy+10], fill=c, outline=cd)
        # 長耳朵
        draw.ellipse([cx-10, cy-20, cx-4, cy-6], fill=c, outline=cd)
        draw.ellipse([cx+4, cy-20, cx+10, cy-6], fill=c, outline=cd)
        # 耳朵內側（粉紅）
        draw.ellipse([cx-8, cy-18, cx-6, cy-8], fill=PALETTE["pink"])
        draw.ellipse([cx+6, cy-18, cx+8, cy-8], fill=PALETTE["pink"])
        # 眼睛（興奮）
        draw.ellipse([cx-6, cy-2, cx-2, cy+2], fill=PALETTE["red"])
        draw.ellipse([cx+2, cy-2, cx+6, cy+2], fill=PALETTE["red"])
        # 討伐棒（旋轉感）
        draw.line([cx+7, cy-5, cx+16, cy-14], fill=PALETTE["yellow_dark"], width=3)
        draw.ellipse([cx+14, cy-16, cx+18, cy-12], fill=PALETTE["yellow"])

    return img

def generate_character_sprites():
    """生成三個角色的 Sprite"""
    chars = ["chiikawa", "hachiware", "usagi"]
    states = ["idle", "attack", "bigwin"]
    out_dir = os.path.join(OUTPUT_BASE, "characters")
    os.makedirs(out_dir, exist_ok=True)

    for char in chars:
        for state in states:
            img = new_img(32, 32)
            draw = ImageDraw.Draw(img)

            if state == "idle":
                draw_pixel_char(img, char, 32)
            elif state == "attack":
                # 攻擊狀態：往右傾斜
                base = new_img(32, 32)
                draw_pixel_char(base, char, 32)
                img = base.rotate(-15, expand=False)
                img = img.convert("RGBA")
            elif state == "bigwin":
                # 大獎：跳起來
                base = new_img(32, 32)
                draw_pixel_char(base, char, 32)
                # 往上移動
                shifted = new_img(32, 32)
                shifted.paste(base, (0, -4))
                img = shifted

            path = os.path.join(out_dir, f"{char}_{state}.png")
            img.save(path)
            print(f"  ✅ {path}")

def generate_target_sprites():
    """生成目標物 Sprite"""
    out_dir = os.path.join(OUTPUT_BASE, "targets")
    os.makedirs(out_dir, exist_ok=True)

    targets = {
        "T001_grass":   {"color": PALETTE["green"],  "shape": "rect",    "size": (16, 20)},
        "T002_bug_g":   {"color": PALETTE["green"],  "shape": "oval",    "size": (20, 12)},
        "T003_bug_r":   {"color": PALETTE["red"],    "shape": "oval",    "size": (20, 12)},
        "T004_bug_b":   {"color": PALETTE["blue"],   "shape": "oval",    "size": (20, 12)},
        "T005_pudding": {"color": PALETTE["yellow"], "shape": "pudding", "size": (24, 28)},
        "T006_mushroom":{"color": PALETTE["brown"],  "shape": "mushroom","size": (28, 28)},
        "T101_mimic":   {"color": PALETTE["purple"], "shape": "star",    "size": (28, 28)},
        "T102_chest":   {"color": PALETTE["gold"],   "shape": "chest",   "size": (28, 24)},
        "T103_meteor":  {"color": PALETTE["yellow"], "shape": "star",    "size": (24, 24)},
        "T104_gold_grass":{"color": PALETTE["gold"], "shape": "rect",    "size": (16, 20)},
        "T105_coin_fish":{"color": PALETTE["gold"],  "shape": "fish",    "size": (32, 24)},
        "B001_boss":    {"color": PALETTE["red"],    "shape": "boss",    "size": (64, 64)},
    }

    for name, cfg in targets.items():
        w, h = cfg["size"]
        img = new_img(w, h)
        draw = ImageDraw.Draw(img)
        c = cfg["color"]
        cd = tuple(max(0, v-50) for v in c[:3]) + (255,)
        shape = cfg["shape"]

        if shape == "rect":
            draw.rectangle([2, 2, w-3, h-3], fill=c, outline=cd)
            # 草的細節
            for i in range(3, w-3, 4):
                draw.line([i, 2, i, 6], fill=cd, width=1)

        elif shape == "oval":
            draw.ellipse([1, 2, w-2, h-3], fill=c, outline=cd)
            # 蟲的眼睛
            draw.ellipse([3, 3, 7, 7], fill=PALETTE["white"])
            draw.ellipse([4, 4, 6, 6], fill=PALETTE["black"])
            # 觸角
            draw.line([4, 2, 2, -1], fill=cd, width=1)
            draw.line([8, 2, 6, -1], fill=cd, width=1)

        elif shape == "pudding":
            # 布丁形狀
            draw.ellipse([2, 8, w-3, h-3], fill=c, outline=cd)
            draw.ellipse([4, 2, w-5, 12], fill=c, outline=cd)
            # 眼睛
            draw.ellipse([6, 10, 10, 14], fill=PALETTE["black"])
            draw.ellipse([w-10, 10, w-6, 14], fill=PALETTE["black"])
            # 腳
            draw.rectangle([5, h-5, 9, h-1], fill=cd)
            draw.rectangle([w-9, h-5, w-5, h-1], fill=cd)

        elif shape == "mushroom":
            # 蘑菇
            draw.ellipse([2, 2, w-3, h//2+4], fill=c, outline=cd)
            draw.rectangle([w//2-4, h//2, w//2+4, h-2], fill=PALETTE["white"], outline=cd)
            # 斑點
            draw.ellipse([6, 6, 10, 10], fill=PALETTE["white"])
            draw.ellipse([w-10, 8, w-6, 12], fill=PALETTE["white"])

        elif shape == "star":
            # 星形（流星/擬態怪物）
            cx, cy = w//2, h//2
            points = []
            for i in range(10):
                angle = math.pi * i / 5 - math.pi/2
                r = (cx-3) if i % 2 == 0 else (cx-3)//2
                points.append((cx + r*math.cos(angle), cy + r*math.sin(angle)))
            draw.polygon(points, fill=c, outline=cd)

        elif shape == "chest":
            # 寶箱
            draw.rectangle([2, h//3, w-3, h-2], fill=c, outline=cd)
            draw.rectangle([2, 2, w-3, h//3+2], fill=cd, outline=cd)
            # 鎖
            draw.rectangle([w//2-3, h//3-2, w//2+3, h//3+4], fill=PALETTE["gold"])
            # 眼睛（寶箱怪）
            draw.ellipse([5, h//3+4, 9, h//3+8], fill=PALETTE["white"])
            draw.ellipse([w-9, h//3+4, w-5, h//3+8], fill=PALETTE["white"])

        elif shape == "fish":
            # 金幣魚
            draw.ellipse([2, 4, w-8, h-5], fill=c, outline=cd)
            # 魚尾
            draw.polygon([w-10, h//2, w-2, 4, w-2, h-4], fill=c, outline=cd)
            # 眼睛
            draw.ellipse([5, h//2-4, 11, h//2+2], fill=PALETTE["white"])
            draw.ellipse([6, h//2-3, 10, h//2+1], fill=PALETTE["black"])
            # 金幣紋路
            draw.ellipse([w//3, h//4, w//3*2, h//4*3], outline=cd, width=1)

        elif shape == "boss":
            # BOSS：那個孩子（大型、可愛但詭異）
            cx, cy = w//2, h//2
            # 身體
            draw.ellipse([8, 12, w-9, h-5], fill=c, outline=cd)
            # 大頭
            draw.ellipse([4, 2, w-5, h//2+8], fill=c, outline=cd)
            # 耳朵（小圓耳）
            draw.ellipse([4, 2, 16, 14], fill=c, outline=cd)
            draw.ellipse([w-16, 2, w-4, 14], fill=c, outline=cd)
            # 眼睛（空洞感）
            draw.ellipse([cx-14, cy-10, cx-4, cy], fill=PALETTE["white"])
            draw.ellipse([cx+4, cy-10, cx+14, cy], fill=PALETTE["white"])
            draw.ellipse([cx-12, cy-8, cx-6, cy-2], fill=PALETTE["black"])
            draw.ellipse([cx+6, cy-8, cx+12, cy-2], fill=PALETTE["black"])
            # 嘴巴（詭異微笑）
            draw.arc([cx-10, cy+2, cx+10, cy+14], 0, 180, fill=PALETTE["black"], width=2)
            # 紅色警告邊框
            draw_pixel_border(draw, 0, 0, w, h, PALETTE["red"], 2)

        path = os.path.join(out_dir, f"{name}.png")
        img.save(path)
        print(f"  ✅ {path}")

def generate_ui_sprites():
    """生成 UI 元素"""
    out_dir = os.path.join(OUTPUT_BASE, "ui")
    os.makedirs(out_dir, exist_ok=True)

    # 金幣
    img = new_img(16, 16)
    draw = ImageDraw.Draw(img)
    draw.ellipse([1, 1, 14, 14], fill=PALETTE["gold"], outline=PALETTE["gold_dark"])
    draw.ellipse([4, 4, 11, 11], fill=PALETTE["yellow"], outline=PALETTE["gold_dark"])
    img.save(os.path.join(out_dir, "coin.png"))
    print(f"  ✅ coin.png")

    # 報酬袋
    img = new_img(20, 24)
    draw = ImageDraw.Draw(img)
    draw.ellipse([2, 6, 17, 22], fill=PALETTE["brown"], outline=PALETTE["gold_dark"])
    draw.rectangle([7, 2, 12, 8], fill=PALETTE["brown"], outline=PALETTE["gold_dark"])
    draw.line([9, 0, 9, 4], fill=PALETTE["gold_dark"], width=2)
    # 袋子上的 ¥ 符號
    draw.line([6, 12, 13, 12], fill=PALETTE["gold"], width=1)
    draw.line([6, 15, 13, 15], fill=PALETTE["gold"], width=1)
    img.save(os.path.join(out_dir, "reward_bag.png"))
    print(f"  ✅ reward_bag.png")

    # 勞動值條（空）
    img = new_img(200, 20)
    draw = ImageDraw.Draw(img)
    draw.rectangle([0, 0, 199, 19], fill=PALETTE["bg_blue"], outline=PALETTE["white"])
    draw.rectangle([2, 2, 197, 17], fill=(30, 30, 60), outline=None)
    img.save(os.path.join(out_dir, "labor_bar_bg.png"))
    print(f"  ✅ labor_bar_bg.png")

    # 勞動值條（填充）
    img = new_img(196, 16)
    draw = ImageDraw.Draw(img)
    # 漸層效果（像素風）
    for x in range(196):
        ratio = x / 196
        r = int(51 + ratio * 204)
        g = int(204 - ratio * 51)
        b = 51
        draw.line([(x, 0), (x, 15)], fill=(r, g, b))
    img.save(os.path.join(out_dir, "labor_bar_fill.png"))
    print(f"  ✅ labor_bar_fill.png")

    # WARNING 字卡
    img = new_img(200, 48)
    draw = ImageDraw.Draw(img)
    draw.rectangle([0, 0, 199, 47], fill=(180, 0, 0), outline=PALETTE["yellow"])
    draw_pixel_border(draw, 0, 0, 200, 48, PALETTE["yellow"], 3)
    # 文字（用像素方塊模擬）
    try:
        font = ImageFont.truetype("C:/Windows/Fonts/Arial.ttf", 20)
        draw.text((20, 12), "!! WARNING !!", fill=PALETTE["yellow"], font=font)
    except:
        draw.text((20, 16), "!! WARNING !!", fill=PALETTE["yellow"])
    img.save(os.path.join(out_dir, "warning_card.png"))
    print(f"  ✅ warning_card.png")

    # 按鈕背景（像素風）
    for name, color in [("btn_normal", PALETTE["bg_blue"]),
                         ("btn_active", PALETTE["green_dark"]),
                         ("btn_auto", PALETTE["blue_dark"])]:
        img = new_img(80, 32)
        draw = ImageDraw.Draw(img)
        draw.rectangle([0, 0, 79, 31], fill=color)
        draw_pixel_border(draw, 0, 0, 80, 32, PALETTE["white"], 2)
        # 高光
        draw.line([(2, 2), (77, 2)], fill=tuple(min(255, v+60) for v in color[:3]))
        img.save(os.path.join(out_dir, f"{name}.png"))
        print(f"  ✅ {name}.png")

    # 像素框（UI 邊框）
    img = new_img(400, 100)
    draw = ImageDraw.Draw(img)
    draw.rectangle([0, 0, 399, 99], fill=(10, 20, 50, 200))
    draw_pixel_border(draw, 0, 0, 400, 100, PALETTE["blue"], 3)
    # 角落裝飾
    for corner in [(0,0), (392,0), (0,92), (392,92)]:
        draw.rectangle([corner[0], corner[1], corner[0]+7, corner[1]+7], fill=PALETTE["yellow"])
    img.save(os.path.join(out_dir, "ui_frame.png"))
    print(f"  ✅ ui_frame.png")

def generate_effects():
    """生成特效圖"""
    out_dir = os.path.join(OUTPUT_BASE, "effects")
    os.makedirs(out_dir, exist_ok=True)

    # 命中特效（星光）
    for char, color in [("chiikawa", PALETTE["pink"]),
                         ("hachiware", PALETTE["blue"]),
                         ("usagi", PALETTE["yellow"])]:
        img = new_img(24, 24)
        draw = ImageDraw.Draw(img)
        cx, cy = 12, 12
        # 十字星光
        for angle in range(0, 360, 45):
            rad = math.radians(angle)
            x1 = cx + 2*math.cos(rad)
            y1 = cy + 2*math.sin(rad)
            x2 = cx + 10*math.cos(rad)
            y2 = cy + 10*math.sin(rad)
            draw.line([(x1,y1),(x2,y2)], fill=color, width=2)
        draw.ellipse([cx-3, cy-3, cx+3, cy+3], fill=PALETTE["white"])
        img.save(os.path.join(out_dir, f"hit_{char}.png"))
        print(f"  ✅ hit_{char}.png")

    # 死亡粒子
    img = new_img(32, 32)
    draw = ImageDraw.Draw(img)
    import random
    random.seed(42)
    for _ in range(12):
        x = random.randint(2, 28)
        y = random.randint(2, 28)
        s = random.randint(2, 5)
        c = random.choice([PALETTE["gold"], PALETTE["yellow"], PALETTE["white"]])
        draw.rectangle([x, y, x+s, y+s], fill=c)
    img.save(os.path.join(out_dir, "death_particles.png"))
    print(f"  ✅ death_particles.png")

    # 投射物（三種角色）
    for char, color in [("chiikawa", PALETTE["pink"]),
                         ("hachiware", PALETTE["blue"]),
                         ("usagi", PALETTE["yellow"])]:
        img = new_img(12, 8)
        draw = ImageDraw.Draw(img)
        # 劍氣形狀
        draw.polygon([(0,4),(4,0),(12,4),(4,8)], fill=color)
        draw.polygon([(2,4),(5,2),(10,4),(5,6)], fill=PALETTE["white"])
        img.save(os.path.join(out_dir, f"projectile_{char}.png"))
        print(f"  ✅ projectile_{char}.png")

def generate_background():
    """生成背景圖"""
    out_dir = os.path.join(OUTPUT_BASE, "backgrounds")
    os.makedirs(out_dir, exist_ok=True)

    # 海底背景（1280x720）
    img = new_img(1280, 720, PALETTE["bg_sea"])
    draw = ImageDraw.Draw(img)

    # 漸層海底
    for y in range(720):
        ratio = y / 720
        r = int(26 * (1-ratio) + 13 * ratio)
        g = int(77 * (1-ratio) + 26 * ratio)
        b = int(153 * (1-ratio) + 77 * ratio)
        draw.line([(0, y), (1280, y)], fill=(r, g, b))

    # 珊瑚（簡化像素）
    import random
    random.seed(123)
    for _ in range(15):
        x = random.randint(50, 1200)
        h = random.randint(40, 100)
        c = random.choice([PALETTE["red"], PALETTE["pink"], PALETTE["purple"]])
        draw.rectangle([x, 720-h, x+8, 720], fill=c)
        draw.ellipse([x-6, 720-h-8, x+14, 720-h+8], fill=c)

    # 海草
    for _ in range(20):
        x = random.randint(0, 1260)
        for i in range(5):
            offset = int(8 * math.sin(i * 0.8))
            draw.rectangle([x+offset, 720-i*15-15, x+offset+6, 720-i*15], fill=PALETTE["green_dark"])

    # 氣泡
    for _ in range(30):
        x = random.randint(0, 1270)
        y = random.randint(100, 600)
        r = random.randint(3, 8)
        draw.ellipse([x, y, x+r*2, y+r*2], outline=PALETTE["white"])

    img.save(os.path.join(out_dir, "sea_bg.png"))
    print(f"  ✅ sea_bg.png (1280x720)")

    # BOSS 背景（暗化版）
    boss_img = img.copy()
    dark = Image.new("RGBA", (1280, 720), (0, 0, 0, 140))
    boss_img = Image.alpha_composite(boss_img.convert("RGBA"), dark)
    boss_draw = ImageDraw.Draw(boss_img)
    # 紅色邊框
    draw_pixel_border(boss_draw, 0, 0, 1280, 720, PALETTE["red"], 8)
    boss_img.save(os.path.join(out_dir, "boss_bg.png"))
    print(f"  ✅ boss_bg.png")

    # Bonus 草地背景
    bonus_img = new_img(1280, 720)
    bonus_draw = ImageDraw.Draw(bonus_img)
    # 天空
    for y in range(400):
        ratio = y / 400
        r = int(135 + ratio * 50)
        g = int(206 + ratio * 20)
        b = int(235 - ratio * 50)
        bonus_draw.line([(0, y), (1280, y)], fill=(r, g, b))
    # 草地
    for y in range(400, 720):
        ratio = (y-400) / 320
        r = int(34 * (1-ratio) + 20 * ratio)
        g = int(139 * (1-ratio) + 80 * ratio)
        b = int(34 * (1-ratio) + 20 * ratio)
        bonus_draw.line([(0, y), (1280, y)], fill=(r, g, b))
    # 草叢
    for _ in range(40):
        x = random.randint(0, 1260)
        for i in range(3):
            bonus_draw.polygon([
                (x+i*6, 400),
                (x+i*6+3, 380-random.randint(0,20)),
                (x+i*6+6, 400)
            ], fill=PALETTE["green"])
    bonus_img.save(os.path.join(out_dir, "bonus_bg.png"))
    print(f"  ✅ bonus_bg.png")

if __name__ == "__main__":
    print("🎨 生成像素美術資產...")
    print("\n[角色 Sprites]")
    generate_character_sprites()
    print("\n[目標物 Sprites]")
    generate_target_sprites()
    print("\n[UI 元素]")
    generate_ui_sprites()
    print("\n[特效]")
    generate_effects()
    print("\n[背景]")
    generate_background()
    print("\n✅ 所有像素美術資產生成完畢！")
    print(f"輸出目錄: {OUTPUT_BASE}")
