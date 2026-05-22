# -*- coding: utf-8 -*-
"""
生成 T118-T126 目標物 sprite（64x64 像素藝術）
T118 皇家閃電鰻 / T119 黃金海龜 / T120 幸運星魚
T121 黃金鯊魚 / T122 金幣魚王 / T123 船長魚
T124 深淵巨鯨 / T125 黃金輪盤螃蟹 / T126 獅子舞魚
"""
from PIL import Image, ImageDraw
import os, math

OUT = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
S = 64

def img():
    return Image.new("RGBA", (S, S), (0,0,0,0))

def px(im, x, y, c):
    if 0 <= x < S and 0 <= y < S:
        im.putpixel((x,y), c)

def circle(im, cx, cy, r, c):
    for y in range(max(0,cy-r), min(S,cy+r+1)):
        for x in range(max(0,cx-r), min(S,cx+r+1)):
            if (x-cx)**2+(y-cy)**2 <= r**2:
                px(im,x,y,c)

def circle_shaded(im, cx, cy, r, base):
    rv,gv,bv = base
    for y in range(max(0,cy-r), min(S,cy+r+1)):
        for x in range(max(0,cx-r), min(S,cx+r+1)):
            if (x-cx)**2+(y-cy)**2 > r**2: continue
            nx_=(x-cx)/max(r,1); ny_=(y-cy)/max(r,1)
            dot=-(nx_*(-0.7)+ny_*(-0.7))
            if dot>0.25: c=(min(255,rv+40),min(255,gv+40),min(255,bv+40),255)
            elif dot<-0.1: c=(max(0,rv-45),max(0,gv-45),max(0,bv-45),255)
            else: c=(rv,gv,bv,255)
            px(im,x,y,c)

def outline(im, cx, cy, r, c):
    for y in range(max(0,cy-r-2), min(S,cy+r+3)):
        for x in range(max(0,cx-r-2), min(S,cx+r+3)):
            d=math.sqrt((x-cx)**2+(y-cy)**2)
            if r+0.1<=d<=r+1.5: px(im,x,y,c)

def rect(im, x1,y1,x2,y2, c):
    for y in range(max(0,y1),min(S,y2)):
        for x in range(max(0,x1),min(S,x2)):
            px(im,x,y,c)

def save(im, name):
    path = os.path.join(OUT, name)
    im.save(path)
    print(f"  saved {name}")

# ── T118 皇家閃電鰻（電藍色，細長身體，電弧紋路）
def gen_T118():
    im = img()
    BODY=(30,80,220,255); LIGHT=(100,180,255,255); OUTLINE=(10,30,120,255)
    GOLD=(255,200,0,255); EYE=(255,255,255,255)
    # 細長身體（S形）
    for i in range(40):
        cx=12+i; cy=32+int(8*math.sin(i*0.3))
        circle(im,cx,cy,4,BODY)
    # 輪廓
    for i in range(40):
        cx=12+i; cy=32+int(8*math.sin(i*0.3))
        outline(im,cx,cy,4,OUTLINE)
    # 電弧紋路（黃色閃電）
    for i in range(0,40,6):
        cx=12+i; cy=32+int(8*math.sin(i*0.3))
        px(im,cx,cy-3,GOLD); px(im,cx+1,cy-2,GOLD); px(im,cx,cy-1,GOLD)
    # 眼睛
    circle(im,14,30,3,LIGHT); circle(im,14,30,2,EYE)
    # 皇冠
    for x in range(10,20): px(im,x,22,GOLD)
    for x in [10,13,16,19]: px(im,x,20,GOLD); px(im,x,21,GOLD)
    save(im,"T118_royal_lightning_eel.png")

# ── T119 黃金海龜（金色龜殼，圓形，慢速）
def gen_T119():
    im = img()
    SHELL=(200,160,20,255); SHELL_L=(240,200,60,255); SHELL_D=(140,100,10,255)
    SKIN=(80,140,60,255); OUTLINE=(60,40,0,255); EYE=(255,255,255,255)
    # 龜殼（大圓）
    circle_shaded(im,32,32,22,(200,160,20))
    outline(im,32,32,22,OUTLINE)
    # 龜殼紋路（六邊形格子）
    for dy in range(-2,3):
        for dx in range(-2,3):
            cx=32+dx*9; cy=32+dy*8
            if (cx-32)**2+(cy-32)**2 < 18**2:
                outline(im,cx,cy,3,SHELL_D)
    # 頭部
    circle_shaded(im,50,28,7,(80,140,60))
    outline(im,50,28,7,OUTLINE)
    # 眼睛
    circle(im,53,26,2,EYE); px(im,53,26,(0,0,0,255))
    # 四肢
    for pos in [(14,18),(14,46),(50,18),(50,46)]:
        circle_shaded(im,pos[0],pos[1],5,(80,140,60))
    save(im,"T119_golden_turtle.png")

# ── T120 幸運星魚（金色星形，閃亮）
def gen_T120():
    im = img()
    GOLD=(255,200,0,255); GOLD_L=(255,240,100,255); GOLD_D=(180,130,0,255)
    OUTLINE=(100,60,0,255); EYE=(255,255,255,255)
    # 星形身體（5角星）
    cx,cy=32,32; R=20; r=9
    pts=[]
    for i in range(10):
        angle=math.pi/2+i*math.pi/5
        rad=R if i%2==0 else r
        pts.append((cx+rad*math.cos(angle),cy-rad*math.sin(angle)))
    # 填充星形
    draw=ImageDraw.Draw(im)
    draw.polygon([(int(p[0]),int(p[1])) for p in pts], fill=GOLD)
    # 高光
    for i in range(0,10,2):
        px_=int(pts[i][0]); py_=int(pts[i][1])
        circle(im,px_,py_,3,GOLD_L)
    # 輪廓
    draw.polygon([(int(p[0]),int(p[1])) for p in pts], outline=(100,60,0,255))
    # 眼睛
    circle(im,36,30,4,EYE); px(im,36,30,(0,0,0,255))
    # 星星光點
    for pos in [(10,10),(54,10),(10,54),(54,54)]:
        px(im,pos[0],pos[1],GOLD_L); px(im,pos[0]+1,pos[1],GOLD_L)
    save(im,"T120_lucky_star_fish.png")

# ── T121 黃金鯊魚（橙金色，三角背鰭，兇猛）
def gen_T121():
    im = img()
    BODY=(220,140,20,255); BODY_L=(255,180,60,255); BELLY=(240,220,180,255)
    OUTLINE=(80,40,0,255); EYE=(255,255,255,255); TEETH=(255,255,255,255)
    # 身體（橢圓）
    for y in range(20,50):
        for x in range(8,58):
            if ((x-33)/24)**2+((y-35)/13)**2<=1:
                nx_=(x-33)/24; ny_=(y-35)/13
                if ny_>0.3: c=BELLY
                elif nx_<-0.5: c=BODY_L
                else: c=BODY
                px(im,x,y,c)
    # 輪廓
    for y in range(20,50):
        for x in range(8,58):
            if 0.95<=((x-33)/24)**2+((y-35)/13)**2<=1.15:
                px(im,x,y,OUTLINE)
    # 背鰭（三角）
    for i in range(12):
        for x in range(28-i//2,38+i//2):
            px(im,x,20-i,BODY)
    # 尾鰭
    for i in range(8):
        px(im,8-i,28+i,BODY); px(im,8-i,42-i,BODY)
    # 眼睛
    circle(im,48,30,4,EYE); px(im,48,30,(0,0,0,255))
    # 牙齒
    for x in range(50,58,3): px(im,x,38,TEETH); px(im,x,39,TEETH)
    save(im,"T121_golden_shark.png")

# ── T122 金幣魚王（金色圓形魚，帶金幣圖案）
def gen_T122():
    im = img()
    GOLD=(220,170,20,255); GOLD_L=(255,210,80,255); OUTLINE=(100,60,0,255)
    EYE=(255,255,255,255); COIN=(255,200,0,255)
    # 圓形身體
    circle_shaded(im,30,32,20,(220,170,20))
    outline(im,30,32,20,OUTLINE)
    # 魚尾
    for i in range(10):
        px(im,10-i,22+i,GOLD); px(im,10-i,42-i,GOLD)
    # 金幣圖案（¥符號）
    for y in range(26,38):
        for x in range(24,36):
            if abs(x-30)<4 and abs(y-32)<5: px(im,x,y,GOLD_L)
    px(im,30,27,OUTLINE); px(im,30,28,OUTLINE)
    for x in range(26,35): px(im,x,30,OUTLINE); px(im,x,32,OUTLINE)
    # 眼睛
    circle(im,40,28,4,EYE); px(im,40,28,(0,0,0,255))
    # 皇冠
    for x in range(26,36): px(im,x,10,COIN)
    for x in [26,29,32,35]: px(im,x,8,COIN); px(im,x,9,COIN)
    save(im,"T122_money_fish.png")

# ── T123 船長魚（藍色，戴船長帽，威嚴）
def gen_T123():
    im = img()
    BODY=(40,100,200,255); BODY_L=(80,150,240,255); BELLY=(180,210,240,255)
    OUTLINE=(10,30,100,255); EYE=(255,255,255,255)
    HAT=(20,20,60,255); HAT_BRIM=(40,40,100,255); GOLD=(255,200,0,255)
    # 身體
    for y in range(22,52):
        for x in range(10,56):
            if ((x-33)/21)**2+((y-37)/13)**2<=1:
                nx_=(x-33)/21; ny_=(y-37)/13
                c=BELLY if ny_>0.3 else (BODY_L if nx_<-0.4 else BODY)
                px(im,x,y,c)
    for y in range(22,52):
        for x in range(10,56):
            if 0.95<=((x-33)/21)**2+((y-37)/13)**2<=1.15:
                px(im,x,y,OUTLINE)
    # 船長帽
    rect(im,24,8,44,18,HAT); rect(im,20,18,48,22,HAT_BRIM)
    for x in range(24,44): px(im,x,8,OUTLINE); px(im,x,18,OUTLINE)
    # 帽徽（金色錨）
    circle(im,33,13,3,GOLD)
    # 眼睛
    circle(im,44,30,4,EYE); px(im,44,30,(0,0,0,255))
    # 鬍子
    for x in range(46,54): px(im,x,36,OUTLINE)
    save(im,"T123_captain_fish.png")

# ── T124 深淵巨鯨（深藍黑色，巨大，威壓感）
def gen_T124():
    im = img()
    BODY=(10,20,80,255); BODY_L=(30,60,140,255); BELLY=(60,100,160,255)
    OUTLINE=(0,0,30,255); EYE=(0,200,255,255); GLOW=(0,100,200,100)
    # 巨大身體（填滿大部分畫面）
    for y in range(10,56):
        for x in range(4,60):
            if ((x-32)/27)**2+((y-33)/20)**2<=1:
                nx_=(x-32)/27; ny_=(y-33)/20
                if ny_>0.4: c=BELLY
                elif nx_<-0.6: c=BODY_L
                else: c=BODY
                px(im,x,y,c)
    for y in range(10,56):
        for x in range(4,60):
            if 0.95<=((x-32)/27)**2+((y-33)/20)**2<=1.15:
                px(im,x,y,OUTLINE)
    # 尾鰭（大）
    for i in range(14):
        px(im,4-i//2,18+i,BODY_L); px(im,4-i//2,48-i,BODY_L)
    # 背鰭
    for i in range(10):
        for x in range(26-i//3,40+i//3):
            px(im,x,10-i,BODY_L)
    # 發光眼睛
    circle(im,50,26,5,EYE); circle(im,50,26,3,(0,255,255,255))
    # 深淵光暈（邊緣藍光）
    for y in range(8,58):
        for x in range(2,62):
            d=math.sqrt((x-32)**2+(y-33)**2)
            if 22<=d<=26: px(im,x,y,GLOW)
    save(im,"T124_abyss_whale.png")

# ── T125 黃金輪盤螃蟹（金色螃蟹，帶輪盤圖案）
def gen_T125():
    im = img()
    BODY=(200,150,20,255); BODY_L=(240,190,60,255); OUTLINE=(80,40,0,255)
    EYE=(255,255,255,255); WHEEL=(255,100,0,255); CLAW=(180,120,10,255)
    # 圓形身體
    circle_shaded(im,32,36,16,(200,150,20))
    outline(im,32,36,16,OUTLINE)
    # 輪盤圖案（8條輻射線）
    for i in range(8):
        angle=i*math.pi/4
        for r in range(4,14):
            x=int(32+r*math.cos(angle)); y=int(36+r*math.sin(angle))
            px(im,x,y,OUTLINE)
    circle(im,32,36,4,WHEEL)
    # 大螯（左右）
    for side in [-1,1]:
        cx=32+side*22; cy=30
        circle_shaded(im,cx,cy,8,(180,120,10))
        outline(im,cx,cy,8,OUTLINE)
        # 連接臂
        for i in range(8):
            px(im,32+side*(14+i),32,CLAW)
    # 眼睛（眼柄）
    for side in [-1,1]:
        for i in range(5):
            px(im,32+side*6,20-i,BODY)
        circle(im,32+side*6,15,4,EYE); px(im,32+side*6,15,(0,0,0,255))
    # 腳（6隻）
    for i in range(3):
        for side in [-1,1]:
            angle=math.pi/2+side*(0.3+i*0.4)
            for r in range(16,26):
                x=int(32+r*math.cos(angle+side*0.5))
                y=int(36+r*math.sin(angle+side*0.5))
                px(im,x,y,CLAW)
    save(im,"T125_roulette_crab.png")

# ── T126 獅子舞魚（橙紅色，獅子頭造型，節慶感）
def gen_T126():
    im = img()
    BODY=(220,80,20,255); BODY_L=(255,130,60,255); MANE=(200,40,0,255)
    OUTLINE=(80,20,0,255); EYE=(255,255,255,255); GOLD=(255,200,0,255)
    # 魚身（橢圓）
    for y in range(20,50):
        for x in range(14,56):
            if ((x-35)/20)**2+((y-35)/13)**2<=1:
                nx_=(x-35)/20
                c=BODY_L if nx_<-0.3 else BODY
                px(im,x,y,c)
    for y in range(20,50):
        for x in range(14,56):
            if 0.95<=((x-35)/20)**2+((y-35)/13)**2<=1.15:
                px(im,x,y,OUTLINE)
    # 獅子鬃毛（頭部周圍）
    for i in range(12):
        angle=i*math.pi/6
        for r in range(18,24):
            x=int(46+r*math.cos(angle)); y=int(32+r*math.sin(angle))
            px(im,x,y,MANE)
    # 獅子頭（大圓）
    circle_shaded(im,46,32,14,(220,80,20))
    outline(im,46,32,14,OUTLINE)
    # 眼睛
    circle(im,50,28,4,EYE); px(im,50,28,(0,0,0,255))
    circle(im,42,28,4,EYE); px(im,42,28,(0,0,0,255))
    # 鼻子
    circle(im,46,34,3,MANE)
    # 嘴巴（笑臉）
    for x in range(42,51): px(im,x,38,OUTLINE)
    px(im,41,37,OUTLINE); px(im,51,37,OUTLINE)
    # 魚尾（節慶彩帶）
    for i in range(10):
        c=GOLD if i%2==0 else MANE
        px(im,14-i,25+i,c); px(im,14-i,45-i,c)
    # 金色裝飾
    for x in range(30,46,4): px(im,x,20,GOLD); px(im,x,50,GOLD)
    save(im,"T126_lion_dance.png")

# ── 主程式 ──────────────────────────────────────────────────────────────────
if __name__ == "__main__":
    os.makedirs(OUT, exist_ok=True)
    print("生成 T118-T126 sprite...")
    gen_T118()
    gen_T119()
    gen_T120()
    gen_T121()
    gen_T122()
    gen_T123()
    gen_T124()
    gen_T125()
    gen_T126()
    print("完成！共生成 9 個 sprite。")
