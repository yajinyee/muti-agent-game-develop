# -*- coding: utf-8 -*-
"""
目標物 T171-T227 批次生成 — Lucky 系列特殊目標物
DAY-109+ 補充：為 T171-T227 生成 Sprite（64x64 像素）
每個目標物用獨特顏色+圖案區分
"""
from PIL import Image, ImageDraw
import os, math

OUTPUT_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
SIZE = 64

def new_img():
    return Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))

def px(img, x, y, c):
    if 0 <= x < SIZE and 0 <= y < SIZE:
        img.putpixel((x, y), c)

def fill_circle_shaded(img, cx, cy, r, base_rgb):
    r_v, g_v, b_v = base_rgb
    for y in range(max(0,cy-r), min(SIZE,cy+r+1)):
        for x in range(max(0,cx-r), min(SIZE,cx+r+1)):
            if (x-cx)**2 + (y-cy)**2 > r**2:
                continue
            nx_ = (x-cx)/max(r,1)
            ny_ = (y-cy)/max(r,1)
            dot = -(nx_*(-0.7) + ny_*(-0.7))
            if dot > 0.25:
                c = (min(255,r_v+45), min(255,g_v+45), min(255,b_v+45), 255)
            elif dot < -0.1:
                c = (max(0,r_v-50), max(0,g_v-50), max(0,b_v-50), 255)
            else:
                c = (r_v, g_v, b_v, 255)
            px(img, x, y, c)

def outline_circle(img, cx, cy, r, color):
    for y in range(max(0,cy-r-2), min(SIZE,cy+r+3)):
        for x in range(max(0,cx-r-2), min(SIZE,cx+r+3)):
            d = math.sqrt((x-cx)**2 + (y-cy)**2)
            if r+0.1 <= d <= r+1.5:
                px(img, x, y, color)

def draw_star(img, cx, cy, r_out, r_in, points, color):
    """畫星形"""
    draw = ImageDraw.Draw(img)
    pts = []
    for i in range(points * 2):
        angle = math.pi * i / points - math.pi / 2
        r = r_out if i % 2 == 0 else r_in
        pts.append((cx + r * math.cos(angle), cy + r * math.sin(angle)))
    draw.polygon(pts, fill=color)

def draw_diamond(img, cx, cy, r, color):
    draw = ImageDraw.Draw(img)
    pts = [(cx, cy-r), (cx+r, cy), (cx, cy+r), (cx-r, cy)]
    draw.polygon(pts, fill=color)

def draw_lightning(img, cx, cy, color):
    """閃電符號"""
    pts = [(cx+4, cy-14), (cx-2, cy-2), (cx+4, cy-2), (cx-4, cy+14), (cx+2, cy+2), (cx-4, cy+2)]
    draw = ImageDraw.Draw(img)
    draw.polygon(pts, fill=color)

def draw_spiral(img, cx, cy, r, color):
    """螺旋點"""
    for i in range(20):
        angle = i * 0.5
        rr = r * i / 20
        x = int(cx + rr * math.cos(angle))
        y = int(cy + rr * math.sin(angle))
        px(img, x, y, color)
        px(img, x+1, y, color)

def draw_cross(img, cx, cy, r, w, color):
    draw = ImageDraw.Draw(img)
    draw.rectangle([cx-w, cy-r, cx+w, cy+r], fill=color)
    draw.rectangle([cx-r, cy-w, cx+r, cy+w], fill=color)

def draw_eye(img, cx, cy, color):
    """眼睛符號"""
    fill_circle_shaded(img, cx, cy, 6, (255,255,255))
    fill_circle_shaded(img, cx, cy, 3, color)
    px(img, cx-1, cy-1, (255,255,255,200))

def add_glow_ring(img, cx, cy, r, color, alpha=80):
    """外圈光暈"""
    for y in range(max(0,cy-r-4), min(SIZE,cy+r+5)):
        for x in range(max(0,cx-r-4), min(SIZE,cx+r+5)):
            d = math.sqrt((x-cx)**2 + (y-cy)**2)
            if r+1 <= d <= r+4:
                fade = int(alpha * (1 - (d-r-1)/3))
                if fade > 0:
                    existing = img.getpixel((x,y))
                    if existing[3] == 0:
                        img.putpixel((x,y), (color[0], color[1], color[2], fade))

def save(img, name):
    path = os.path.join(OUTPUT_DIR, name)
    img.save(path)
    print(f"  saved: {name}")

# ---- 目標物定義：ID -> (名稱, 主色, 輔色, 圖案類型) ----
TARGETS = {
    "T171": ("彩虹稜鏡魚",  (180,100,220), (255,200,50),  "prism"),
    "T172": ("黃金累積魚",  (220,170,30),  (255,220,80),  "accumulate"),
    "T173": ("幸運鏡像魚",  (100,180,220), (200,240,255), "mirror"),
    "T174": ("詛咒毒魚",    (80,180,80),   (150,255,100), "poison"),
    "T175": ("幸運拍賣魚",  (220,120,50),  (255,180,80),  "auction"),
    "T176": ("幸運進化魚",  (50,200,150),  (100,255,200), "evolution"),
    "T177": ("幸運連鎖感染魚",(180,50,180),(240,100,240), "infection"),
    "T178": ("幸運反彈魚",  (50,150,220),  (100,200,255), "ricochet"),
    "T179": ("幸運黑洞魚",  (30,20,60),    (120,80,200),  "blackhole"),
    "T180": ("幸運共鳴魚",  (200,100,50),  (255,160,80),  "resonance"),
    "T181": ("幸運傳送魚",  (80,200,200),  (150,255,255), "teleport"),
    "T182": ("幸運分裂魚",  (220,80,80),   (255,140,140), "split"),
    "T183": ("幸運充能魚",  (50,100,220),  (100,160,255), "charge"),
    "T184": ("幸運鏈鎖爆炸魚",(220,100,30),(255,160,60),  "chainbomb"),
    "T185": ("幸運鏡像時空魚",(100,80,200),(180,160,255), "mirrortime"),
    "T186": ("幸運量子魚",  (0,200,180),   (80,255,240),  "quantum"),
    "T187": ("幸運寄生魚",  (120,180,50),  (180,240,80),  "parasite"),
    "T188": ("幸運風暴魚",  (60,120,200),  (120,180,255), "storm"),
    "T189": ("幸運迴旋鏢魚",(200,150,50),  (255,210,80),  "boomerang"),
    "T190": ("幸運磁力魚",  (180,50,100),  (240,100,160), "magnet"),
    "T191": ("幸運回聲魚",  (100,200,150), (160,255,200), "echo"),
    "T192": ("幸運漩渦魚",  (50,80,200),   (100,140,255), "vortex"),
    "T193": ("幸運時間炸彈魚",(200,60,60), (255,120,80),  "timebomb"),
    "T194": ("幸運鏡面世界魚",(80,160,200),(160,220,255), "mirrorworld"),
    "T195": ("幸運冰凍世界魚",(100,180,220),(180,230,255),"freezeworld"),
    "T196": ("幸運重力反轉魚",(150,80,200),(210,140,255), "gravity"),
    "T197": ("幸運共鳴爆發魚",(220,120,80),(255,180,120), "synergy"),
    "T198": ("幸運賭注魚",  (200,160,30),  (255,220,60),  "bet"),
    "T199": ("幸運連鎖反應魚",(180,80,180),(240,140,240), "chainreact"),
    "T200": ("幸運分身魚",  (80,180,100),  (140,240,160), "clone"),
    "T201": ("幸運預言魚",  (120,80,200),  (180,140,255), "prophecy"),
    "T202": ("幸運奪旗魚",  (220,60,60),   (255,120,100), "flag"),
    "T203": ("幸運幽靈魚",  (160,160,200), (220,220,255), "phantom"),
    "T204": ("幸運水晶球魚",(100,200,220), (160,240,255), "crystalball"),
    "T205": ("幸運時光倒流魚",(200,100,180),(255,160,230),"timerewind"),
    "T206": ("幸運龍捲風魚",(80,140,200),  (140,200,255), "tornado"),
    "T207": ("幸運黑洞爆炸魚",(40,20,80), (140,80,220),  "bhexplosion"),
    "T208": ("幸運鏡像分裂魚",(100,160,220),(160,220,255),"mirrorsplit"),
    "T209": ("幸運量子糾纏魚",(0,180,160), (60,240,220),  "quantum2"),
    "T210": ("幸運武器進化魚",(200,140,40),(255,200,80),  "weaponevo"),
    "T211": ("幸運星際隕石魚",(180,80,40), (240,140,80),  "meteor"),
    "T212": ("幸運龍王降臨魚",(200,50,50), (255,100,80),  "dragonking"),
    "T213": ("幸運時空裂縫魚",(80,40,160), (160,100,240), "rift"),
    "T214": ("幸運全服充能魚",(40,160,200),(100,220,255), "servercharge"),
    "T215": ("幸運公會戰魚", (200,100,40), (255,160,80),  "guildwar"),
    "T216": ("幸運閃電風暴魚",(60,100,220),(120,160,255), "lightningstorm"),
    "T217": ("幸運星座命運魚",(160,60,200),(220,120,255), "zodiac"),
    "T218": ("幸運寶藏獵人魚",(200,160,40),(255,220,80),  "treasure"),
    "T219": ("幸運時間膠囊魚",(80,180,180),(140,240,240), "timecapsule"),
    "T220": ("幸運累積大獎池魚",(220,160,30),(255,220,60),"progjakpot"),
    "T221": ("幸運元素融合魚",(100,200,100),(160,255,160),"elemfusion"),
    "T222": ("幸運命運輪迴魚",(180,80,160),(240,140,220), "karmacycle"),
    "T224": ("幸運連鎖爆炸魚",(220,80,40), (255,140,80),  "chainexp"),
    "T225": ("幸運倍率疊加魚",(200,160,20),(255,220,50),  "multstk"),
    "T226": ("幸運倒數炸彈魚",(220,60,30), (255,120,60),  "cntbomb"),
    "T227": ("幸運輪盤魚",   (160,40,200), (220,100,255), "spinwheel"),
}

def make_fish(tid, name, main_rgb, accent_rgb, pattern):
    """生成一條魚的 sprite"""
    img = new_img()
    cx, cy = 32, 32
    mr, mg, mb = main_rgb
    ar, ag, ab = accent_rgb

    # 魚身（橢圓形）
    for y in range(cy-14, cy+15):
        for x in range(cx-20, cx+21):
            if (x-cx)**2/400 + (y-cy)**2/196 <= 1.0:
                nx_ = (x-cx)/20
                ny_ = (y-cy)/14
                dot = -(nx_*(-0.7) + ny_*(-0.7))
                if dot > 0.3:
                    c = (min(255,mr+50), min(255,mg+50), min(255,mb+50), 255)
                elif dot < -0.1:
                    c = (max(0,mr-55), max(0,mg-55), max(0,mb-55), 255)
                else:
                    c = (mr, mg, mb, 255)
                px(img, x, y, c)

    # 魚尾
    for y in range(cy-10, cy+11):
        for x in range(cx+18, cx+28):
            dist_from_body = x - (cx+18)
            half_h = int(10 * (1 - dist_from_body/10))
            if abs(y-cy) <= half_h:
                px(img, x, y, (max(0,mr-30), max(0,mg-30), max(0,mb-30), 255))

    # 魚眼
    fill_circle_shaded(img, cx-10, cy-4, 5, (240,240,240))
    fill_circle_shaded(img, cx-10, cy-4, 3, (30,30,30))
    px(img, cx-11, cy-5, (255,255,255,220))

    # 魚鰭（上）
    for y in range(cy-20, cy-12):
        for x in range(cx-8, cx+6):
            if y > cy-20 + (x-(cx-8))*0.5:
                px(img, x, y, (max(0,mr-20), max(0,mg-20), max(0,mb-20), 200))

    # 輪廓
    outline_circle(img, cx, cy, 15, (max(0,mr-80), max(0,mg-80), max(0,mb-80), 255))

    # 圖案裝飾（依 pattern 類型）
    _add_pattern(img, cx, cy, ar, ag, ab, pattern)

    # 光暈
    add_glow_ring(img, cx, cy, 18, accent_rgb, 60)

    return img

def _add_pattern(img, cx, cy, ar, ag, ab, pattern):
    ac = (ar, ag, ab, 255)
    ac2 = (min(255,ar+40), min(255,ag+40), min(255,ab+40), 200)

    if pattern == "prism":
        # 彩虹稜鏡：多色條紋
        colors = [(255,80,80,180),(255,200,50,180),(80,255,80,180),(80,150,255,180),(200,80,255,180)]
        for i, c in enumerate(colors):
            for y in range(cy-10, cy+11):
                x = cx - 8 + i*4
                if abs(y-cy) < 10:
                    px(img, x, y, c)
    elif pattern == "accumulate":
        # 黃金累積：金幣堆疊
        for i in range(3):
            fill_circle_shaded(img, cx-4+i*4, cy+2, 4, (220,180,30))
            outline_circle(img, cx-4+i*4, cy+2, 4, (180,140,20,255))
    elif pattern == "mirror":
        # 鏡像：對稱線
        for y in range(cy-12, cy+13):
            px(img, cx, y, ac)
            px(img, cx+1, y, ac2)
    elif pattern == "poison":
        # 毒：氣泡
        for ox, oy in [(-6,-6),(6,-4),(-4,6),(8,2),(-8,2)]:
            fill_circle_shaded(img, cx+ox, cy+oy, 3, (80,200,80))
    elif pattern == "auction":
        # 拍賣：錘子形狀
        draw_cross(img, cx+4, cy-4, 5, 2, ac)
    elif pattern == "evolution":
        # 進化：箭頭向上
        draw = ImageDraw.Draw(img)
        draw.polygon([(cx,cy-12),(cx-6,cy-4),(cx+6,cy-4)], fill=ac)
        draw.rectangle([cx-2,cy-4,cx+2,cy+4], fill=ac)
    elif pattern == "infection":
        # 感染：擴散圓
        for r in [4, 8, 12]:
            outline_circle(img, cx, cy, r, (ar, ag, ab, 120))
    elif pattern == "ricochet":
        # 反彈：Z 形
        draw = ImageDraw.Draw(img)
        draw.line([(cx-8,cy-8),(cx+8,cy-8),(cx-8,cy+8),(cx+8,cy+8)], fill=ac, width=2)
    elif pattern == "blackhole":
        # 黑洞：螺旋
        draw_spiral(img, cx, cy, 12, ac)
    elif pattern in ("resonance","synergy"):
        # 共鳴：同心圓
        for r in [5, 9, 13]:
            outline_circle(img, cx, cy, r, ac)
    elif pattern == "teleport":
        # 傳送：菱形
        draw_diamond(img, cx, cy, 10, ac)
    elif pattern == "split":
        # 分裂：X 形
        draw = ImageDraw.Draw(img)
        draw.line([(cx-8,cy-8),(cx+8,cy+8)], fill=ac, width=2)
        draw.line([(cx+8,cy-8),(cx-8,cy+8)], fill=ac, width=2)
    elif pattern in ("charge","servercharge"):
        # 充能：閃電
        draw_lightning(img, cx, cy, ac)
    elif pattern in ("chainbomb","chainexp","cntbomb"):
        # 炸彈：圓+引線
        fill_circle_shaded(img, cx, cy+2, 8, (60,60,60))
        for i in range(4):
            px(img, cx+4+i, cy-6-i, ac)
    elif pattern in ("mirrortime","mirrorworld","mirrorsplit"):
        # 鏡像時空：雙菱形
        draw_diamond(img, cx-5, cy, 6, ac)
        draw_diamond(img, cx+5, cy, 6, (ar,ag,ab,150))
    elif pattern in ("quantum","quantum2"):
        # 量子：原子軌道
        outline_circle(img, cx, cy, 10, ac)
        outline_circle(img, cx, cy, 6, ac2)
        fill_circle_shaded(img, cx, cy, 3, (ar,ag,ab))
    elif pattern == "parasite":
        # 寄生：觸手
        for angle in range(0, 360, 60):
            rad = math.radians(angle)
            ex = int(cx + 12*math.cos(rad))
            ey = int(cy + 12*math.sin(rad))
            draw = ImageDraw.Draw(img)
            draw.line([(cx,cy),(ex,ey)], fill=ac, width=1)
    elif pattern == "storm":
        # 風暴：旋轉箭頭
        for angle in [0, 120, 240]:
            rad = math.radians(angle)
            ex = int(cx + 10*math.cos(rad))
            ey = int(cy + 10*math.sin(rad))
            px(img, ex, ey, ac)
            px(img, ex+1, ey, ac)
    elif pattern == "boomerang":
        # 迴旋鏢：弧形
        draw = ImageDraw.Draw(img)
        draw.arc([cx-10,cy-10,cx+10,cy+10], 0, 180, fill=ac, width=3)
    elif pattern == "magnet":
        # 磁力：U 形
        draw = ImageDraw.Draw(img)
        draw.arc([cx-8,cy-8,cx+8,cy+8], 180, 360, fill=ac, width=3)
        draw.line([(cx-8,cy),(cx-8,cy+6)], fill=ac, width=2)
        draw.line([(cx+8,cy),(cx+8,cy+6)], fill=ac, width=2)
    elif pattern == "echo":
        # 回聲：同心弧
        draw = ImageDraw.Draw(img)
        for r in [5, 9, 13]:
            draw.arc([cx-r,cy-r,cx+r,cy+r], -60, 60, fill=ac, width=1)
    elif pattern == "vortex":
        # 漩渦：螺旋
        draw_spiral(img, cx, cy, 14, ac)
    elif pattern == "timebomb":
        # 時間炸彈：時鐘
        outline_circle(img, cx, cy, 10, ac)
        draw = ImageDraw.Draw(img)
        draw.line([(cx,cy),(cx,cy-7)], fill=ac, width=2)
        draw.line([(cx,cy),(cx+5,cy+3)], fill=ac, width=2)
    elif pattern == "freezeworld":
        # 冰凍：雪花
        draw = ImageDraw.Draw(img)
        for angle in [0, 60, 120]:
            rad = math.radians(angle)
            draw.line([(cx-int(10*math.cos(rad)),cy-int(10*math.sin(rad))),
                       (cx+int(10*math.cos(rad)),cy+int(10*math.sin(rad)))], fill=ac, width=2)
    elif pattern == "gravity":
        # 重力：向下箭頭
        draw = ImageDraw.Draw(img)
        draw.polygon([(cx,cy+12),(cx-6,cy+4),(cx+6,cy+4)], fill=ac)
        draw.rectangle([cx-2,cy-8,cx+2,cy+4], fill=ac)
    elif pattern == "bet":
        # 賭注：骰子
        draw = ImageDraw.Draw(img)
        draw.rectangle([cx-7,cy-7,cx+7,cy+7], fill=ac, outline=(ar-40,ag-40,ab-40,255))
        for ox, oy in [(-3,-3),(3,-3),(0,0),(-3,3),(3,3)]:
            fill_circle_shaded(img, cx+ox, cy+oy, 1, (255,255,255))
    elif pattern == "chainreact":
        # 連鎖反應：鏈條
        for i in range(4):
            outline_circle(img, cx-9+i*6, cy, 3, ac)
    elif pattern == "clone":
        # 分身：多個小圓
        for ox, oy in [(-8,-4),(8,-4),(0,8)]:
            fill_circle_shaded(img, cx+ox, cy+oy, 4, (ar,ag,ab))
    elif pattern == "prophecy":
        # 預言：眼睛
        draw_eye(img, cx, cy, (ar,ag,ab))
    elif pattern == "flag":
        # 奪旗：旗幟
        draw = ImageDraw.Draw(img)
        draw.line([(cx-4,cy-12),(cx-4,cy+8)], fill=ac, width=2)
        draw.polygon([(cx-4,cy-12),(cx+8,cy-8),(cx-4,cy-4)], fill=ac)
    elif pattern == "phantom":
        # 幽靈：半透明輪廓
        outline_circle(img, cx, cy, 12, (ar,ag,ab,150))
        outline_circle(img, cx, cy, 8, (ar,ag,ab,100))
    elif pattern == "crystalball":
        # 水晶球：球+光點
        outline_circle(img, cx, cy, 12, ac)
        fill_circle_shaded(img, cx-4, cy-4, 2, (255,255,255))
    elif pattern == "timerewind":
        # 時光倒流：逆時針箭頭
        draw = ImageDraw.Draw(img)
        draw.arc([cx-10,cy-10,cx+10,cy+10], 90, 360, fill=ac, width=2)
        draw.polygon([(cx-10,cy),(cx-6,cy-4),(cx-6,cy+4)], fill=ac)
    elif pattern == "tornado":
        # 龍捲風：螺旋三角
        draw = ImageDraw.Draw(img)
        draw.polygon([(cx,cy-12),(cx-8,cy+8),(cx+8,cy+8)], outline=ac, fill=None)
        draw.polygon([(cx,cy-6),(cx-4,cy+4),(cx+4,cy+4)], fill=ac)
    elif pattern == "bhexplosion":
        # 黑洞爆炸：黑洞+爆炸線
        fill_circle_shaded(img, cx, cy, 8, (20,10,40))
        for angle in range(0, 360, 45):
            rad = math.radians(angle)
            ex = int(cx + 14*math.cos(rad))
            ey = int(cy + 14*math.sin(rad))
            draw = ImageDraw.Draw(img)
            draw.line([(cx,cy),(ex,ey)], fill=ac, width=1)
    elif pattern == "weaponevo":
        # 武器進化：劍形
        draw = ImageDraw.Draw(img)
        draw.polygon([(cx,cy-14),(cx-3,cy+6),(cx+3,cy+6)], fill=ac)
        draw.rectangle([cx-4,cy+6,cx+4,cy+10], fill=(ar-40,ag-40,ab-40,255))
    elif pattern == "meteor":
        # 隕石：橢圓+尾跡
        fill_circle_shaded(img, cx+4, cy-4, 7, (ar,ag,ab))
        draw = ImageDraw.Draw(img)
        draw.line([(cx+4,cy-4),(cx-8,cy+8)], fill=ac, width=2)
        draw.line([(cx+2,cy-2),(cx-6,cy+6)], fill=(ar,ag,ab,150), width=1)
    elif pattern == "dragonking":
        # 龍王：龍頭輪廓
        fill_circle_shaded(img, cx, cy-4, 10, (ar,ag,ab))
        for ox, oy in [(-8,-12),(8,-12)]:
            fill_circle_shaded(img, cx+ox, cy+oy, 4, (ar,ag,ab))
        draw_eye(img, cx-4, cy-6, (255,50,50))
        draw_eye(img, cx+4, cy-6, (255,50,50))
    elif pattern == "rift":
        # 時空裂縫：鋸齒線
        draw = ImageDraw.Draw(img)
        pts = [(cx-12,cy),(cx-6,cy-8),(cx,cy),(cx+6,cy+8),(cx+12,cy)]
        draw.line(pts, fill=ac, width=2)
    elif pattern == "guildwar":
        # 公會戰：盾牌
        draw = ImageDraw.Draw(img)
        draw.polygon([(cx,cy-12),(cx-10,cy-4),(cx-10,cy+6),(cx,cy+12),(cx+10,cy+6),(cx+10,cy-4)], fill=ac)
        draw.line([(cx,cy-10),(cx,cy+10)], fill=(255,255,255,200), width=2)
    elif pattern == "lightningstorm":
        # 閃電風暴：多閃電
        draw_lightning(img, cx-6, cy, ac)
        draw_lightning(img, cx+6, cy, (ar,ag,ab,180))
    elif pattern == "zodiac":
        # 星座：星形
        draw_star(img, cx, cy, 12, 5, 6, ac)
    elif pattern == "treasure":
        # 寶藏：寶箱
        draw = ImageDraw.Draw(img)
        draw.rectangle([cx-9,cy-4,cx+9,cy+8], fill=ac)
        draw.rectangle([cx-9,cy-8,cx+9,cy-4], fill=(ar-30,ag-30,ab-30,255))
        draw.rectangle([cx-2,cy-2,cx+2,cy+2], fill=(255,220,50,255))
    elif pattern == "timecapsule":
        # 時間膠囊：膠囊形
        draw = ImageDraw.Draw(img)
        draw.ellipse([cx-6,cy-12,cx+6,cy+12], fill=ac)
        draw.line([(cx-6,cy),(cx+6,cy)], fill=(255,255,255,200), width=2)
    elif pattern == "progjakpot":
        # 累積大獎池：金幣+光芒
        fill_circle_shaded(img, cx, cy, 10, (220,180,30))
        draw_star(img, cx, cy, 14, 8, 8, (255,220,50,120))
        outline_circle(img, cx, cy, 10, (180,140,20,255))
    elif pattern == "elemfusion":
        # 元素融合：四色象限
        colors = [(255,80,80,200),(80,200,80,200),(80,80,255,200),(255,200,50,200)]
        for i, c in enumerate(colors):
            angle_start = i * 90
            draw = ImageDraw.Draw(img)
            draw.pieslice([cx-12,cy-12,cx+12,cy+12], angle_start, angle_start+90, fill=c)
    elif pattern == "karmacycle":
        # 命運輪迴：太極
        fill_circle_shaded(img, cx, cy, 12, (ar,ag,ab))
        fill_circle_shaded(img, cx, cy-6, 5, (255,255,255))
        fill_circle_shaded(img, cx, cy+6, 5, (30,30,30))
        fill_circle_shaded(img, cx, cy-6, 2, (30,30,30))
        fill_circle_shaded(img, cx, cy+6, 2, (255,255,255))
    elif pattern == "multstk":
        # 倍率疊加：數字堆疊
        draw = ImageDraw.Draw(img)
        for i, mult in enumerate(["×1","×3","×5"]):
            draw.text((cx-8, cy-10+i*7), mult, fill=ac)
    elif pattern == "spinwheel":
        # 輪盤：扇形
        draw = ImageDraw.Draw(img)
        colors_w = [ac, (255,255,255,200), ac, (255,255,255,200), ac, (255,255,255,200)]
        for i, c in enumerate(colors_w):
            draw.pieslice([cx-12,cy-12,cx+12,cy+12], i*60, (i+1)*60, fill=c)
        outline_circle(img, cx, cy, 12, (ar-40,ag-40,ab-40,255))
        fill_circle_shaded(img, cx, cy, 3, (200,200,200))
    else:
        # 預設：星形
        draw_star(img, cx, cy, 10, 5, 5, ac)

    return img

def main():
    os.makedirs(OUTPUT_DIR, exist_ok=True)
    # 只生成尚未存在的 sprite
    generated = 0
    skipped = 0
    for tid, (name, main_rgb, accent_rgb, pattern) in TARGETS.items():
        # 命名格式：T171_rainbow_prism_fish.png（用 pattern 作為後綴）
        fname = f"{tid}_{pattern}.png"
        fpath = os.path.join(OUTPUT_DIR, fname)
        if os.path.exists(fpath):
            skipped += 1
            continue
        try:
            img = make_fish(tid, name, main_rgb, accent_rgb, pattern)
            save(img, fname)
            generated += 1
        except Exception as e:
            print(f"  ERROR {tid}: {e}")
    print(f"\n完成！生成 {generated} 個，跳過 {skipped} 個（已存在）")

if __name__ == "__main__":
    main()
