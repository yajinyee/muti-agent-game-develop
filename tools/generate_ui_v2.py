# -*- coding: utf-8 -*-
"""
UI 元素 v2 - 更大更清晰的像素藝術 UI
coin: 32x32, reward_bag: 40x48, btn: 96x36, warning: 256x64
"""
from PIL import Image, ImageDraw
import os
import math

UI_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\ui"
os.makedirs(UI_DIR, exist_ok=True)

def px(img, x, y, c):
    if 0 <= x < img.width and 0 <= y < img.height:
        img.putpixel((x, y), c)

def fill_circle(img, cx, cy, r, color):
    for y in range(max(0,cy-r), min(img.height,cy+r+1)):
        for x in range(max(0,cx-r), min(img.width,cx+r+1)):
            if (x-cx)**2+(y-cy)**2 <= r**2:
                px(img, x, y, color)

def outline_circle(img, cx, cy, r, color):
    for y in range(max(0,cy-r-1), min(img.height,cy+r+2)):
        for x in range(max(0,cx-r-1), min(img.width,cx+r+2)):
            d = math.sqrt((x-cx)**2+(y-cy)**2)
            if r+0.1 <= d <= r+1.4:
                px(img, x, y, color)

def gen_coin():
    """金幣 32x32 - 金色圓形，帶光澤和¥符號"""
    SIZE = 32
    img = Image.new("RGBA", (SIZE, SIZE), (0,0,0,0))
    cx, cy = SIZE//2, SIZE//2

    GOLD   = (220, 175, 30)
    GOLD_L = (255, 230, 80)
    GOLD_D = (160, 120, 10)
    OUTLINE= (100, 60, 5, 255)
    SHINE  = (255, 255, 200, 255)

    # 金幣主體（帶陰影）
    for y in range(cy-12, cy+13):
        for x in range(cx-12, cx+13):
            if (x-cx)**2+(y-cy)**2 <= 144:
                nx_ = (x-cx)/12
                ny_ = (y-cy)/12
                dot = -(nx_*(-0.7)+ny_*(-0.7))
                if dot > 0.25:
                    c = GOLD_L
                elif dot < -0.1:
                    c = GOLD_D
                else:
                    c = GOLD
                px(img, x, y, (*c, 255))

    # 輪廓
    outline_circle(img, cx, cy, 12, OUTLINE)

    # 內圈（裝飾）
    for y in range(cy-9, cy+10):
        for x in range(cx-9, cx+10):
            d = math.sqrt((x-cx)**2+(y-cy)**2)
            if 8.5 <= d <= 9.5:
                px(img, x, y, (*GOLD_D, 200))

    # ¥ 符號（像素）
    # 垂直線
    for y in range(cy-4, cy+6):
        px(img, cx, y, (*GOLD_D, 255))
        px(img, cx+1, y, (*GOLD_D, 255))
    # 橫線 x2
    for x in range(cx-4, cx+6):
        px(img, x, cy, (*GOLD_D, 255))
        px(img, x, cy+2, (*GOLD_D, 255))
    # Y 上半部
    for i in range(5):
        px(img, cx-4+i, cy-4-i//2, (*GOLD_D, 255))
        px(img, cx+5-i, cy-4-i//2, (*GOLD_D, 255))

    # 高光
    px(img, cx-5, cy-5, SHINE)
    px(img, cx-4, cy-5, SHINE)
    px(img, cx-5, cy-4, SHINE)

    return img


def gen_reward_bag():
    """勞動報酬袋 40x48 - 布袋形狀，帶¥符號"""
    W, H = 40, 48
    img = Image.new("RGBA", (W, H), (0,0,0,0))

    BROWN  = (160, 100, 40)
    BROWN_L= (200, 140, 70)
    BROWN_D= (100, 60, 15)
    GOLD   = (220, 175, 30)
    OUTLINE= (60, 30, 5, 255)
    ROPE   = (120, 70, 20, 255)

    cx = W//2

    # 袋子主體（梨形）
    for y in range(12, H-2):
        t = (y-12) / (H-14)
        # 寬度從上到下先增後減
        if t < 0.6:
            w = int(8 + t * 25)
        else:
            w = int(23 - (t-0.6) * 15)
        w = max(4, w)
        for x in range(cx-w, cx+w+1):
            nx_ = (x-cx)/max(w,1)
            ny_ = (y-12)/(H-14)
            dot = -(nx_*(-0.7)+ny_*(-0.7))
            if dot > 0.25:
                c = BROWN_L
            elif dot < -0.1:
                c = BROWN_D
            else:
                c = BROWN
            px(img, x, y, (*c, 255))

    # 袋口（繩子）
    for x in range(cx-6, cx+7):
        for y in range(8, 14):
            px(img, x, y, ROPE)
    for y in range(8, 14):
        px(img, cx-6, y, OUTLINE)
        px(img, cx+6, y, OUTLINE)
    for x in range(cx-6, cx+7):
        px(img, x, 8, OUTLINE)
        px(img, x, 13, OUTLINE)

    # 袋頂（打結）
    fill_circle(img, cx, 6, 5, (*BROWN, 255))
    outline_circle(img, cx, 6, 5, OUTLINE)

    # 輪廓
    for y in range(12, H-2):
        t = (y-12) / (H-14)
        if t < 0.6:
            w = int(8 + t * 25)
        else:
            w = int(23 - (t-0.6) * 15)
        w = max(4, w)
        px(img, cx-w-1, y, OUTLINE)
        px(img, cx+w+1, y, OUTLINE)
    for x in range(cx-8, cx+9):
        px(img, x, H-2, OUTLINE)

    # ¥ 符號（金色）
    for y in range(24, 36):
        px(img, cx, y, (*GOLD, 255))
        px(img, cx+1, y, (*GOLD, 255))
    for x in range(cx-5, cx+7):
        px(img, x, 28, (*GOLD, 255))
        px(img, x, 30, (*GOLD, 255))
    for i in range(5):
        px(img, cx-5+i, 24-i//2, (*GOLD, 255))
        px(img, cx+6-i, 24-i//2, (*GOLD, 255))

    return img


def gen_button(state="normal"):
    """按鈕 96x36 - 圓角矩形，帶光澤"""
    W, H = 96, 36
    img = Image.new("RGBA", (W, H), (0,0,0,0))
    draw = ImageDraw.Draw(img)

    if state == "normal":
        BG     = (60, 80, 140)
        BG_L   = (90, 120, 190)
        BG_D   = (30, 50, 100)
        BORDER = (150, 180, 255, 255)
    elif state == "active":
        BG     = (40, 160, 80)
        BG_L   = (70, 200, 110)
        BG_D   = (20, 110, 50)
        BORDER = (150, 255, 180, 255)
    else:  # auto
        BG     = (140, 80, 160)
        BG_L   = (180, 120, 200)
        BG_D   = (90, 40, 110)
        BORDER = (220, 160, 255, 255)

    R = 8  # 圓角半徑

    # 主體（帶漸層）
    for y in range(H):
        for x in range(W):
            # 圓角遮罩
            in_rect = True
            if x < R and y < R and (x-R)**2+(y-R)**2 > R**2:
                in_rect = False
            if x > W-R-1 and y < R and (x-(W-R-1))**2+(y-R)**2 > R**2:
                in_rect = False
            if x < R and y > H-R-1 and (x-R)**2+(y-(H-R-1))**2 > R**2:
                in_rect = False
            if x > W-R-1 and y > H-R-1 and (x-(W-R-1))**2+(y-(H-R-1))**2 > R**2:
                in_rect = False
            if not in_rect:
                continue
            # 漸層（上亮下暗）
            t = y / H
            r_v = int(BG_L[0] + (BG_D[0]-BG_L[0])*t)
            g_v = int(BG_L[1] + (BG_D[1]-BG_L[1])*t)
            b_v = int(BG_L[2] + (BG_D[2]-BG_L[2])*t)
            img.putpixel((x, y), (r_v, g_v, b_v, 255))

    # 邊框
    draw.rounded_rectangle([0, 0, W-1, H-1], radius=R, outline=BORDER[:3], width=2)

    # 上方高光線
    for x in range(R, W-R):
        img.putpixel((x, 2), (255, 255, 255, 80))
        img.putpixel((x, 3), (255, 255, 255, 40))

    return img


def gen_labor_bar_bg():
    """勞動值條背景 240x24"""
    W, H = 240, 24
    img = Image.new("RGBA", (W, H), (0,0,0,0))
    draw = ImageDraw.Draw(img)

    # 背景（深色圓角）
    draw.rounded_rectangle([0, 0, W-1, H-1], radius=6, fill=(30, 30, 50, 220))
    draw.rounded_rectangle([0, 0, W-1, H-1], radius=6, outline=(80, 80, 120, 255), width=2)

    # 內部凹陷效果
    draw.rounded_rectangle([2, 2, W-3, H-3], radius=5, outline=(15, 15, 30, 150), width=1)

    return img


def gen_labor_bar_fill():
    """勞動值條填充 236x20 - 綠色漸層"""
    W, H = 236, 20
    img = Image.new("RGBA", (W, H), (0,0,0,0))
    draw = ImageDraw.Draw(img)

    # 漸層（左綠右黃綠）
    for x in range(W):
        t = x / W
        r = int(50 + t * 100)
        g = int(180 + t * 50)
        b = int(50 - t * 30)
        for y in range(H):
            img.putpixel((x, y), (r, g, b, 255))

    # 圓角遮罩
    draw.rounded_rectangle([0, 0, W-1, H-1], radius=5, outline=(0,0,0,0), width=0)

    # 上方高光
    for x in range(W):
        img.putpixel((x, 1), (255, 255, 255, 60))
        img.putpixel((x, 2), (255, 255, 255, 30))

    # 邊框
    draw.rounded_rectangle([0, 0, W-1, H-1], radius=5, outline=(100, 220, 80, 255), width=1)

    return img


def gen_warning_card():
    """WARNING 字卡 256x64 - 紅色警告"""
    W, H = 256, 64
    img = Image.new("RGBA", (W, H), (0,0,0,0))
    draw = ImageDraw.Draw(img)

    # 背景（半透明紅）
    draw.rounded_rectangle([0, 0, W-1, H-1], radius=8,
                           fill=(180, 20, 20, 200))
    draw.rounded_rectangle([0, 0, W-1, H-1], radius=8,
                           outline=(255, 80, 80, 255), width=3)

    # 內框
    draw.rounded_rectangle([4, 4, W-5, H-5], radius=6,
                           outline=(255, 150, 150, 150), width=1)

    # ⚠ 符號（左側）
    # 三角形
    tri_pts = [(20, 48), (36, 48), (28, 16)]
    draw.polygon(tri_pts, fill=(255, 220, 50, 255), outline=(200, 150, 20, 255))
    # ! 符號
    for y in range(24, 38):
        draw.point((28, y), fill=(30, 20, 5, 255))
        draw.point((29, y), fill=(30, 20, 5, 255))
    draw.ellipse([27, 40, 31, 44], fill=(30, 20, 5, 255))

    # WARNING 文字（大）
    draw.text((50, 12), "WARNING", fill=(255, 255, 255, 255))
    draw.text((50, 36), "BOSS IS COMING!", fill=(255, 220, 100, 255))

    return img


def gen_ui_frame():
    """UI 框架 400x100 - 底部控制欄背景"""
    W, H = 400, 100
    img = Image.new("RGBA", (W, H), (0,0,0,0))
    draw = ImageDraw.Draw(img)

    # 半透明深色背景
    draw.rounded_rectangle([0, 0, W-1, H-1], radius=12,
                           fill=(20, 25, 45, 210))
    draw.rounded_rectangle([0, 0, W-1, H-1], radius=12,
                           outline=(80, 100, 160, 255), width=2)

    # 上方高光線
    for x in range(12, W-12):
        img.putpixel((x, 2), (150, 180, 255, 60))

    # 裝飾點（四角）
    for (cx, cy) in [(10, 10), (W-10, 10), (10, H-10), (W-10, H-10)]:
        fill_circle(img, cx, cy, 3, (100, 140, 220, 200))

    return img


def main():
    print("=== UI 元素 v2 生成 ===\n")

    assets = [
        ("coin.png",           gen_coin()),
        ("reward_bag.png",     gen_reward_bag()),
        ("btn_normal.png",     gen_button("normal")),
        ("btn_active.png",     gen_button("active")),
        ("btn_auto.png",       gen_button("auto")),
        ("labor_bar_bg.png",   gen_labor_bar_bg()),
        ("labor_bar_fill.png", gen_labor_bar_fill()),
        ("warning_card.png",   gen_warning_card()),
        ("ui_frame.png",       gen_ui_frame()),
    ]

    for filename, img in assets:
        path = os.path.join(UI_DIR, filename)
        img.save(path)
        print(f"  ✅ {filename}: {img.size}")

    print(f"\n✅ UI v2 完成！")

if __name__ == "__main__":
    main()
