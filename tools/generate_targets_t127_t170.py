# -*- coding: utf-8 -*-
"""
目標物 T127-T170 Sprite 生成（DAY-212 補齊）
每個目標物 64x64 像素，帶陰影和細節
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
            nx_ = (x-cx)/max(r,1); ny_ = (y-cy)/max(r,1)
            dot = -(nx_*(-0.7) + ny_*(-0.7))
            if dot > 0.25:
                c = (min(255,r_v+40), min(255,g_v+40), min(255,b_v+40), 255)
            elif dot < -0.1:
                c = (max(0,r_v-45), max(0,g_v-45), max(0,b_v-45), 255)
            else:
                c = (r_v, g_v, b_v, 255)
            px(img, x, y, c)

def fill_rect_shaded(img, x1, y1, x2, y2, base_rgb):
    r_v, g_v, b_v = base_rgb
    w = max(1, x2-x1); h = max(1, y2-y1)
    for y in range(max(0,y1), min(SIZE,y2)):
        for x in range(max(0,x1), min(SIZE,x2)):
            nx_ = (x-x1)/w*2-1; ny_ = (y-y1)/h*2-1
            dot = -(nx_*(-0.7)+ny_*(-0.7))
            if dot > 0.25:
                c = (min(255,r_v+35), min(255,g_v+35), min(255,b_v+35), 255)
            elif dot < -0.1:
                c = (max(0,r_v-40), max(0,g_v-40), max(0,b_v-40), 255)
            else:
                c = (r_v, g_v, b_v, 255)
            px(img, x, y, c)

def fill_circle(img, cx, cy, r, color):
    for y in range(max(0,cy-r), min(SIZE,cy+r+1)):
        for x in range(max(0,cx-r), min(SIZE,cx+r+1)):
            if (x-cx)**2 + (y-cy)**2 <= r**2:
                px(img, x, y, color)

def outline_all(img, color=(20,20,20,255)):
    orig = img.copy()
    w, h = orig.size
    for y in range(h):
        for x in range(w):
            if orig.getpixel((x,y))[3] > 10:
                for dx,dy in [(-1,0),(1,0),(0,-1),(0,1)]:
                    nx,ny = x+dx, y+dy
                    if 0<=nx<w and 0<=ny<h and orig.getpixel((nx,ny))[3] < 10:
                        img.putpixel((nx,ny), color)
    return img

def draw_eye(img, cx, cy, pupil=(20,20,20,255)):
    fill_circle(img, cx, cy, 3, (255,255,255,255))
    fill_circle(img, cx, cy, 2, pupil)
    px(img, cx-1, cy-1, (255,255,255,200))

def save(img, name):
    path = os.path.join(OUTPUT_DIR, name)
    img.save(path)
    print(f"  saved: {name}")

# ─── T170 時間凍結魚 ─────────────────────────────────────────────────────────
def gen_t170():
    img = new_img()
    # 冰藍色魚身（橢圓）
    fill_circle_shaded(img, 32, 32, 18, (0, 191, 255))  # 冰藍主體
    # 魚尾（三角形）
    for y in range(22, 42):
        w = int((y-22) * 0.6) if y < 32 else int((42-y) * 0.6)
        for x in range(50, 50+w+6):
            px(img, x, y, (0, 150, 220, 255))
    # 魚鰭（上下）
    fill_circle_shaded(img, 28, 18, 6, (135, 206, 235))  # 上鰭
    fill_circle_shaded(img, 28, 46, 5, (135, 206, 235))  # 下鰭
    # 眼睛（冰白色）
    draw_eye(img, 22, 28, (200, 240, 255, 255))
    # 冰晶效果（身體上的雪花點）
    for cx, cy in [(30,25),(38,30),(25,35),(35,38),(20,30)]:
        px(img, cx, cy, (255,255,255,220))
        px(img, cx+1, cy, (200,240,255,180))
        px(img, cx, cy+1, (200,240,255,180))
    # 冰晶光暈（外圈淡藍）
    for angle in range(0, 360, 30):
        rad = math.radians(angle)
        gx = int(32 + 22 * math.cos(rad))
        gy = int(32 + 22 * math.sin(rad))
        px(img, gx, gy, (135, 206, 235, 120))
    outline_all(img, (0, 80, 150, 255))
    save(img, "T170_time_freeze_fish.png")

# ─── T127-T169 批次生成（使用主題色區分）────────────────────────────────────
TARGETS_T127_T169 = [
    # (id, name, body_rgb, fin_rgb, eye_pupil, outline_rgb, special)
    ("T127", "vortex_fish",       (100,180,220), (60,140,200),  (20,60,100,255),  (20,60,120,255),  "vortex"),
    ("T128", "freeze_bomb",       (180,220,255), (120,180,240), (20,80,160,255),  (20,60,140,255),  "ice"),
    ("T129", "ice_fish",          (200,235,255), (150,200,240), (30,100,180,255), (20,70,150,255),  "ice"),
    ("T130", "lucky_egg_fish",    (255,220,100), (220,180,60),  (120,60,20,255),  (140,80,20,255),  "egg"),
    ("T131", "rainbow_lucky_fish",(255,100,180), (200,60,140),  (80,20,60,255),   (120,20,80,255),  "rainbow"),
    ("T132", "sea_anemone",       (255,120,80),  (220,80,40),   (60,20,10,255),   (140,40,20,255),  "tentacle"),
    ("T133", "lucky_dice_fish",   (255,255,255), (200,200,200), (20,20,20,255),   (60,60,60,255),   "dice"),
    ("T134", "fire_storm_fish",   (255,80,20),   (220,40,0),    (60,10,0,255),    (120,20,0,255),   "fire"),
    ("T135", "golden_treasure",   (255,200,0),   (220,160,0),   (100,60,0,255),   (140,80,0,255),   "treasure"),
    ("T136", "mermaid",           (100,220,180), (60,180,140),  (20,80,60,255),   (20,100,80,255),  "mermaid"),
    ("T137", "lucky_clover_fish", (80,200,80),   (40,160,40),   (20,60,20,255),   (20,80,20,255),   "clover"),
    ("T138", "rainbow_shark",     (180,100,255), (140,60,220),  (60,20,100,255),  (80,20,140,255),  "shark"),
    ("T139", "thunder_shark",     (255,220,0),   (220,180,0),   (80,60,0,255),    (120,80,0,255),   "shark"),
    ("T140", "vampire_fish",      (180,20,60),   (140,0,40),    (200,0,0,255),    (80,0,20,255),    "vampire"),
    ("T141", "lightning_fish",    (255,255,100), (220,220,40),  (80,80,0,255),    (120,100,0,255),  "lightning"),
    ("T142", "meteor_fish",       (255,140,60),  (220,100,20),  (80,40,0,255),    (120,60,0,255),   "meteor"),
    ("T143", "phoenix_fish",      (255,160,0),   (220,120,0),   (100,40,0,255),   (140,60,0,255),   "phoenix"),
    ("T144", "dragon_turtle",     (80,160,80),   (40,120,40),   (20,60,20,255),   (20,80,20,255),   "turtle"),
    ("T145", "chain_bomb",        (60,60,60),    (40,40,40),    (200,200,200,255),(20,20,20,255),   "bomb"),
    ("T146", "croc_hunter",       (100,160,60),  (60,120,20),   (20,60,10,255),   (20,80,10,255),   "croc"),
    ("T147", "time_bomb_fish",    (255,60,60),   (220,20,20),   (60,0,0,255),     (100,0,0,255),    "bomb"),
    ("T148", "triple_lucky",      (255,180,60),  (220,140,20),  (80,40,0,255),    (120,60,0,255),   "triple"),
    ("T149", "school_leader",     (60,180,220),  (20,140,180),  (0,60,100,255),   (0,80,120,255),   "school"),
    ("T150", "rock_skeleton",     (200,200,200), (160,160,160), (20,20,20,255),   (60,60,60,255),   "skull"),
    ("T151", "electric_jellyfish",(180,100,255), (140,60,220),  (60,20,100,255),  (80,20,140,255),  "jellyfish"),
    ("T152", "chainlong_king",    (255,160,0),   (220,120,0),   (100,40,0,255),   (140,60,0,255),   "dragon"),
    ("T153", "drill_bit_lobster", (255,100,60),  (220,60,20),   (80,20,0,255),    (120,40,0,255),   "lobster"),
    ("T154", "anglerfish_elec",   (60,60,120),   (40,40,100),   (200,200,255,255),(20,20,80,255),   "anglerfish"),
    ("T155", "mystic_dragon",     (120,60,200),  (80,20,160),   (200,100,255,255),(60,0,120,255),   "dragon"),
    ("T156", "ghost_fish",        (220,220,255), (180,180,240), (100,100,200,255),(80,80,160,255),  "ghost"),
    ("T157", "thunder_lobster_v2",(255,200,0),   (220,160,0),   (80,60,0,255),    (120,80,0,255),   "lobster"),
    ("T158", "ice_phoenix",       (100,200,255), (60,160,220),  (0,80,160,255),   (0,60,140,255),   "phoenix"),
    ("T159", "serial_bomb_crab",  (255,120,40),  (220,80,0),    (80,30,0,255),    (120,50,0,255),   "crab"),
    ("T160", "abyss_vortex",      (60,0,120),    (40,0,100),    (180,100,255,255),(30,0,80,255),    "vortex"),
    ("T161", "humpback_whale",    (60,100,160),  (20,60,120),   (0,40,100,255),   (0,40,100,255),   "whale"),
    ("T162", "free_spin_fish",    (0,200,200),   (0,160,160),   (0,60,80,255),    (0,80,100,255),   "spin"),
    ("T163", "jackpot_dragon",    (255,180,0),   (220,140,0),   (100,50,0,255),   (140,70,0,255),   "dragon"),
    ("T164", "comet_fish",        (255,140,60),  (220,100,20),  (80,40,0,255),    (120,60,0,255),   "comet"),
    ("T165", "golden_wave_fish",  (255,200,0),   (220,160,0),   (100,60,0,255),   (140,80,0,255),   "wave"),
    ("T166", "dragon_king",       (200,0,0),     (160,0,0),     (255,200,200,255),(100,0,0,255),    "dragon"),
    ("T167", "fortune_coin_fish", (255,215,0),   (220,175,0),   (100,60,0,255),   (140,80,0,255),   "coin"),
    ("T168", "lucky_hot_zone",    (255,100,0),   (220,60,0),    (80,20,0,255),    (120,40,0,255),   "fire"),
    ("T169", "lucky_trident",     (150,80,220),  (110,40,180),  (60,0,100,255),   (80,0,140,255),   "trident"),
]

def gen_generic_fish(tid, name, body_rgb, fin_rgb, eye_pupil, outline_rgb, special):
    img = new_img()
    # 魚身主體
    fill_circle_shaded(img, 30, 32, 18, body_rgb)
    # 魚尾
    for y in range(22, 42):
        w = int(abs(y-32) * 0.4) + 4
        for x in range(50, 50+w):
            px(img, x, y, (fin_rgb[0], fin_rgb[1], fin_rgb[2], 255))
    # 魚鰭
    fill_circle_shaded(img, 26, 18, 6, fin_rgb)
    fill_circle_shaded(img, 26, 46, 5, fin_rgb)
    # 眼睛
    draw_eye(img, 20, 28, eye_pupil)
    # 特殊標記
    if special == "ice":
        for cx,cy in [(32,26),(38,32),(26,36)]:
            px(img, cx, cy, (255,255,255,220))
    elif special == "fire":
        for cx,cy in [(34,22),(38,28),(30,22)]:
            px(img, cx, cy, (255,255,100,200))
    elif special == "lightning":
        for cx,cy in [(34,22),(36,28),(32,34)]:
            px(img, cx, cy, (255,255,0,220))
    elif special == "dragon":
        fill_circle(img, 32, 14, 4, (fin_rgb[0], fin_rgb[1], fin_rgb[2], 200))
        fill_circle(img, 28, 14, 3, (fin_rgb[0], fin_rgb[1], fin_rgb[2], 200))
    elif special == "ghost":
        for cx,cy in [(28,26),(34,26),(28,32),(34,32)]:
            px(img, cx, cy, (255,255,255,180))
    elif special == "whale":
        fill_circle_shaded(img, 30, 32, 22, body_rgb)  # 更大的身體
    elif special == "shark":
        fill_rect_shaded(img, 26, 14, 36, 22, fin_rgb)  # 背鰭
    elif special == "skull":
        fill_circle(img, 20, 28, 3, (20,20,20,255))
        fill_circle(img, 26, 28, 3, (20,20,20,255))
    elif special == "coin":
        fill_circle(img, 32, 32, 8, (255,215,0,200))
        px(img, 32, 32, (255,255,255,220))
    elif special == "trident":
        for y in range(14, 22):
            px(img, 32, y, (255,215,0,220))
        px(img, 28, 14, (255,215,0,220))
        px(img, 36, 14, (255,215,0,220))
    outline_all(img, outline_rgb)
    fname = f"{tid}_{name}.png"
    save(img, fname)

if __name__ == "__main__":
    os.makedirs(OUTPUT_DIR, exist_ok=True)
    print("生成 T127-T169...")
    for tid, name, body, fin, eye, outline, special in TARGETS_T127_T169:
        fname = f"{tid}_{name}.png"
        fpath = os.path.join(OUTPUT_DIR, fname)
        if not os.path.exists(fpath):
            gen_generic_fish(tid, name, body, fin, eye, outline, special)
        else:
            print(f"  skip (exists): {fname}")
    print("生成 T170...")
    gen_t170()
    print("完成！")
