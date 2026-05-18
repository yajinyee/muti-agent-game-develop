"""
generate_bonus_weeds.py
生成 Bonus Game 的 5 種像素雜草 Sprite（BG001-BG005）
尺寸：32×48 px（寬×高），輸出到 assets/sprites/targets/
"""

import os
from PIL import Image, ImageDraw

# 輸出目錄
OUT_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"

# ── 顏色定義 ──────────────────────────────────────────────────

# BG001 普通雜草（綠色）
BG001_COLORS = {
    "stem":   (60, 140, 40),
    "stem_d": (40, 100, 25),
    "leaf":   (80, 180, 50),
    "leaf_d": (55, 140, 35),
    "leaf_l": (120, 210, 80),
    "outline":(20, 60, 10),
    "soil":   (120, 80, 40),
    "soil_d": (90, 55, 25),
}

# BG002 硬雜草（深綠，粗壯）
BG002_COLORS = {
    "stem":   (30, 100, 20),
    "stem_d": (20, 70, 12),
    "leaf":   (50, 140, 30),
    "leaf_d": (35, 100, 20),
    "leaf_l": (80, 170, 55),
    "outline":(10, 40, 5),
    "thorn":  (180, 60, 20),
    "soil":   (100, 65, 30),
    "soil_d": (75, 45, 18),
}

# BG003 發光雜草（亮綠，帶光暈）
BG003_COLORS = {
    "stem":   (80, 200, 60),
    "stem_d": (55, 160, 40),
    "leaf":   (120, 240, 80),
    "leaf_d": (85, 190, 55),
    "leaf_l": (180, 255, 130),
    "outline":(30, 100, 15),
    "glow":   (200, 255, 150),
    "soil":   (130, 90, 45),
    "soil_d": (100, 65, 30),
}

# BG004 金色雜草（金黃色，閃亮）
BG004_COLORS = {
    "stem":   (200, 160, 20),
    "stem_d": (160, 120, 10),
    "leaf":   (240, 200, 40),
    "leaf_d": (190, 150, 20),
    "leaf_l": (255, 235, 100),
    "outline":(100, 70, 5),
    "shine":  (255, 250, 180),
    "soil":   (140, 95, 50),
    "soil_d": (110, 70, 35),
}

# BG005 搗亂怪草（紅色，帶眼睛）
BG005_COLORS = {
    "stem":   (160, 40, 40),
    "stem_d": (120, 25, 25),
    "leaf":   (200, 60, 60),
    "leaf_d": (155, 40, 40),
    "leaf_l": (230, 100, 100),
    "outline":(60, 10, 10),
    "eye_w":  (255, 255, 255),
    "eye_p":  (20, 20, 20),
    "soil":   (110, 70, 35),
    "soil_d": (85, 50, 22),
}


def px(img, x, y, color, alpha=255):
    """安全設置像素"""
    if 0 <= x < img.width and 0 <= y < img.height:
        if len(color) == 3:
            img.putpixel((x, y), (color[0], color[1], color[2], alpha))
        else:
            img.putpixel((x, y), color)


def draw_stem(img, cx, base_y, top_y, color_mid, color_dark, color_outline, width=3):
    """畫莖幹（垂直，帶陰影）"""
    for y in range(top_y, base_y + 1):
        # 左側輪廓
        px(img, cx - width // 2 - 1, y, color_outline)
        # 右側輪廓
        px(img, cx + width // 2 + 1, y, color_outline)
        # 莖幹主體（左暗右亮）
        for dx in range(-width // 2, width // 2 + 1):
            if dx < 0:
                px(img, cx + dx, y, color_dark)
            else:
                px(img, cx + dx, y, color_mid)


def draw_leaf(img, cx, cy, direction, size, c_light, c_mid, c_dark, c_outline):
    """畫葉片（橢圓形，帶陰影）
    direction: 'left' or 'right'
    """
    w = size
    h = size // 2 + 2

    sign = -1 if direction == 'left' else 1

    for dy in range(-h, h + 1):
        for dx in range(-w, w + 1):
            # 橢圓判定
            if (dx / w) ** 2 + (dy / h) ** 2 <= 1.0:
                nx = cx + dx * sign
                ny = cy + dy
                # 陰影：左上亮，右下暗
                if dx < -w // 3 and dy < 0:
                    c = c_light
                elif dx > w // 3 and dy > 0:
                    c = c_dark
                else:
                    c = c_mid
                px(img, nx, ny, c)

    # 輪廓
    for dy in range(-h - 1, h + 2):
        for dx in range(-w - 1, w + 2):
            nx = cx + dx * sign
            ny = cy + dy
            if (dx / (w + 0.5)) ** 2 + (dy / (h + 0.5)) ** 2 <= 1.0:
                if (dx / w) ** 2 + (dy / h) ** 2 > 1.0:
                    px(img, nx, ny, c_outline)


def generate_bg001():
    """BG001 普通雜草 — 簡單三葉草形"""
    img = Image.new("RGBA", (32, 48), (0, 0, 0, 0))
    c = BG001_COLORS

    # 土壤底部
    for x in range(8, 24):
        px(img, x, 44, c["soil"])
        px(img, x, 45, c["soil_d"])
    for x in range(10, 22):
        px(img, x, 43, c["soil"])

    # 主莖（中央，從 y=43 到 y=20）
    draw_stem(img, 16, 43, 20, c["stem"], c["stem_d"], c["outline"], width=2)

    # 左側小莖（從 y=38 到 y=28，向左傾斜）
    for y in range(28, 39):
        t = (y - 28) / 10.0
        x = int(16 - t * 5)
        px(img, x, y, c["stem"])
        px(img, x - 1, y, c["stem_d"])
        px(img, x - 2, y, c["outline"])

    # 右側小莖（從 y=35 到 y=25，向右傾斜）
    for y in range(25, 36):
        t = (y - 25) / 10.0
        x = int(16 + t * 5)
        px(img, x, y, c["stem"])
        px(img, x + 1, y, c["stem_d"])
        px(img, x + 2, y, c["outline"])

    # 三片葉子
    draw_leaf(img, 16, 18, 'left', 6, c["leaf_l"], c["leaf"], c["leaf_d"], c["outline"])   # 頂葉（偏左）
    draw_leaf(img, 11, 26, 'left', 5, c["leaf_l"], c["leaf"], c["leaf_d"], c["outline"])   # 左葉
    draw_leaf(img, 21, 23, 'right', 5, c["leaf_l"], c["leaf"], c["leaf_d"], c["outline"])  # 右葉

    return img


def generate_bg002():
    """BG002 硬雜草 — 粗壯，帶刺"""
    img = Image.new("RGBA", (32, 48), (0, 0, 0, 0))
    c = BG002_COLORS

    # 土壤底部（更寬）
    for x in range(6, 26):
        px(img, x, 44, c["soil"])
        px(img, x, 45, c["soil_d"])
    for x in range(8, 24):
        px(img, x, 43, c["soil"])

    # 粗主莖（寬 4px）
    draw_stem(img, 16, 43, 16, c["stem"], c["stem_d"], c["outline"], width=3)

    # 左側粗莖
    for y in range(26, 38):
        t = (y - 26) / 12.0
        x = int(16 - t * 7)
        for dx in range(-1, 2):
            px(img, x + dx, y, c["stem"] if dx >= 0 else c["stem_d"])
        px(img, x - 2, y, c["outline"])
        px(img, x + 2, y, c["outline"])

    # 右側粗莖
    for y in range(22, 34):
        t = (y - 22) / 12.0
        x = int(16 + t * 7)
        for dx in range(-1, 2):
            px(img, x + dx, y, c["stem"] if dx >= 0 else c["stem_d"])
        px(img, x - 2, y, c["outline"])
        px(img, x + 2, y, c["outline"])

    # 三片大葉子
    draw_leaf(img, 16, 14, 'left', 7, c["leaf_l"], c["leaf"], c["leaf_d"], c["outline"])
    draw_leaf(img, 9, 24, 'left', 6, c["leaf_l"], c["leaf"], c["leaf_d"], c["outline"])
    draw_leaf(img, 23, 20, 'right', 6, c["leaf_l"], c["leaf"], c["leaf_d"], c["outline"])

    # 刺（橙色小三角）
    for (tx, ty) in [(14, 30), (18, 26), (12, 22)]:
        px(img, tx, ty, c["thorn"])
        px(img, tx - 1, ty + 1, c["thorn"])
        px(img, tx + 1, ty + 1, c["thorn"])
        px(img, tx, ty + 2, c["thorn"])

    return img


def generate_bg003():
    """BG003 發光雜草 — 亮綠，帶光點"""
    img = Image.new("RGBA", (32, 48), (0, 0, 0, 0))
    c = BG003_COLORS

    # 土壤底部
    for x in range(8, 24):
        px(img, x, 44, c["soil"])
        px(img, x, 45, c["soil_d"])
    for x in range(10, 22):
        px(img, x, 43, c["soil"])

    # 主莖（稍細，發光感）
    draw_stem(img, 16, 43, 18, c["stem"], c["stem_d"], c["outline"], width=2)

    # 左右莖
    for y in range(28, 38):
        t = (y - 28) / 10.0
        x = int(16 - t * 6)
        px(img, x, y, c["stem"])
        px(img, x - 1, y, c["stem_d"])
        px(img, x - 2, y, c["outline"])

    for y in range(24, 34):
        t = (y - 24) / 10.0
        x = int(16 + t * 6)
        px(img, x, y, c["stem"])
        px(img, x + 1, y, c["stem_d"])
        px(img, x + 2, y, c["outline"])

    # 三片葉子（更亮）
    draw_leaf(img, 16, 16, 'left', 6, c["leaf_l"], c["leaf"], c["leaf_d"], c["outline"])
    draw_leaf(img, 10, 26, 'left', 5, c["leaf_l"], c["leaf"], c["leaf_d"], c["outline"])
    draw_leaf(img, 22, 22, 'right', 5, c["leaf_l"], c["leaf"], c["leaf_d"], c["outline"])

    # 發光光點（散落在葉子上）
    glow_points = [(15, 14), (18, 15), (9, 24), (12, 27), (21, 20), (24, 23)]
    for (gx, gy) in glow_points:
        px(img, gx, gy, c["glow"])
        px(img, gx, gy - 1, c["glow"], 180)
        px(img, gx - 1, gy, c["glow"], 180)

    # 頂部星形光點
    px(img, 16, 10, c["glow"])
    px(img, 15, 11, c["glow"], 200)
    px(img, 17, 11, c["glow"], 200)
    px(img, 16, 12, c["glow"], 150)

    return img


def generate_bg004():
    """BG004 金色雜草 — 金黃，帶閃光"""
    img = Image.new("RGBA", (32, 48), (0, 0, 0, 0))
    c = BG004_COLORS

    # 土壤底部（金色土壤）
    for x in range(7, 25):
        px(img, x, 44, c["soil"])
        px(img, x, 45, c["soil_d"])
    for x in range(9, 23):
        px(img, x, 43, c["soil"])

    # 主莖（金色）
    draw_stem(img, 16, 43, 15, c["stem"], c["stem_d"], c["outline"], width=2)

    # 左右莖（更彎曲）
    for y in range(26, 38):
        t = (y - 26) / 12.0
        x = int(16 - t * 7)
        px(img, x, y, c["stem"])
        px(img, x - 1, y, c["stem_d"])
        px(img, x - 2, y, c["outline"])

    for y in range(22, 34):
        t = (y - 22) / 12.0
        x = int(16 + t * 7)
        px(img, x, y, c["stem"])
        px(img, x + 1, y, c["stem_d"])
        px(img, x + 2, y, c["outline"])

    # 三片金色葉子（更大）
    draw_leaf(img, 16, 13, 'left', 7, c["leaf_l"], c["leaf"], c["leaf_d"], c["outline"])
    draw_leaf(img, 9, 24, 'left', 6, c["leaf_l"], c["leaf"], c["leaf_d"], c["outline"])
    draw_leaf(img, 23, 20, 'right', 6, c["leaf_l"], c["leaf"], c["leaf_d"], c["outline"])

    # 閃光（十字星形）
    shine_points = [(16, 10), (10, 22), (22, 18)]
    for (sx, sy) in shine_points:
        px(img, sx, sy, c["shine"])
        px(img, sx - 1, sy, c["shine"], 200)
        px(img, sx + 1, sy, c["shine"], 200)
        px(img, sx, sy - 1, c["shine"], 200)
        px(img, sx, sy + 1, c["shine"], 200)
        px(img, sx - 2, sy, c["leaf_l"], 120)
        px(img, sx + 2, sy, c["leaf_l"], 120)
        px(img, sx, sy - 2, c["leaf_l"], 120)
        px(img, sx, sy + 2, c["leaf_l"], 120)

    return img


def generate_bg005():
    """BG005 搗亂怪草 — 紅色，帶眼睛和表情"""
    img = Image.new("RGBA", (32, 48), (0, 0, 0, 0))
    c = BG005_COLORS

    # 土壤底部
    for x in range(8, 24):
        px(img, x, 44, c["soil"])
        px(img, x, 45, c["soil_d"])
    for x in range(10, 22):
        px(img, x, 43, c["soil"])

    # 主莖（紅色，扭曲感）
    for y in range(20, 44):
        t = (y - 20) / 24.0
        # 輕微 S 形扭曲
        import math
        dx = int(math.sin(t * math.pi * 2) * 1.5)
        x = 16 + dx
        px(img, x, y, c["stem"])
        px(img, x - 1, y, c["stem_d"])
        px(img, x + 1, y, c["stem"])
        px(img, x - 2, y, c["outline"])
        px(img, x + 2, y, c["outline"])

    # 左右莖
    for y in range(28, 38):
        t = (y - 28) / 10.0
        x = int(16 - t * 6)
        px(img, x, y, c["stem"])
        px(img, x - 1, y, c["stem_d"])
        px(img, x - 2, y, c["outline"])

    for y in range(24, 34):
        t = (y - 24) / 10.0
        x = int(16 + t * 6)
        px(img, x, y, c["stem"])
        px(img, x + 1, y, c["stem_d"])
        px(img, x + 2, y, c["outline"])

    # 三片紅色葉子
    draw_leaf(img, 16, 18, 'left', 6, c["leaf_l"], c["leaf"], c["leaf_d"], c["outline"])
    draw_leaf(img, 10, 26, 'left', 5, c["leaf_l"], c["leaf"], c["leaf_d"], c["outline"])
    draw_leaf(img, 22, 22, 'right', 5, c["leaf_l"], c["leaf"], c["leaf_d"], c["outline"])

    # 眼睛（在頂部葉子上）
    # 左眼
    for dy in range(-2, 3):
        for dx in range(-2, 3):
            if dx * dx + dy * dy <= 4:
                px(img, 13 + dx, 16 + dy, c["eye_w"])
    px(img, 13, 16, c["eye_p"])
    px(img, 14, 16, c["eye_p"])

    # 右眼
    for dy in range(-2, 3):
        for dx in range(-2, 3):
            if dx * dx + dy * dy <= 4:
                px(img, 19 + dx, 16 + dy, c["eye_w"])
    px(img, 19, 16, c["eye_p"])
    px(img, 20, 16, c["eye_p"])

    # 眉毛（憤怒）
    for dx in range(-2, 2):
        px(img, 13 + dx, 13 - dx // 2, c["outline"])
        px(img, 19 + dx, 13 + dx // 2, c["outline"])

    # 嘴巴（邪惡微笑）
    for dx in range(-3, 4):
        py_val = 20 + abs(dx) // 2
        px(img, 16 + dx, py_val, c["outline"])

    return img


def save_with_import(img, name):
    """儲存 PNG 並建立 .import 檔案"""
    path = os.path.join(OUT_DIR, name + ".png")
    img.save(path)
    print(f"  ✅ {name}.png ({img.width}x{img.height}px)")

    # 建立 .import 檔案
    import_path = path + ".import"
    import_content = f"""[remap]

importer="texture"
type="CompressedTexture2D"
uid="uid://bonus_{name.lower()}"
path="res://.godot/imported/{name}.png-{name.lower()}.ctex"
metadata={{
"vram_texture": false
}}

[deps]

source_file="res://assets/sprites/targets/{name}.png"
dest_files=["res://.godot/imported/{name}.png-{name.lower()}.ctex"]

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
"""
    with open(import_path, 'w') as f:
        f.write(import_content)


def analyze_density(img):
    """計算非透明像素密度"""
    total = img.width * img.height
    non_transparent = sum(1 for x in range(img.width) for y in range(img.height)
                          if img.getpixel((x, y))[3] > 10)
    return non_transparent, total, non_transparent / total * 100


def main():
    print("=== 生成 Bonus 雜草 Sprites（BG001-BG005）===\n")

    generators = [
        ("BG001_weed_normal", generate_bg001),
        ("BG002_weed_hard",   generate_bg002),
        ("BG003_weed_glow",   generate_bg003),
        ("BG004_weed_gold",   generate_bg004),
        ("BG005_weed_evil",   generate_bg005),
    ]

    for name, gen_func in generators:
        img = gen_func()
        non_t, total, density = analyze_density(img)
        print(f"  生成 {name}: {non_t}/{total} px ({density:.1f}%)")
        save_with_import(img, name)

    print("\n✅ 全部完成！")
    print(f"輸出目錄：{OUT_DIR}")
    print("\n下一步：更新 BonusGame.gd 使用真正的 Sprite 替代 ColorRect")


if __name__ == "__main__":
    main()
