"""
generate_boss_sprites.py — B001 BOSS 完整動畫集生成
那個孩子（The Child）— 吉伊卡哇宇宙中的神秘 BOSS

動畫狀態：
  idle    (4幀) — 緩慢漂浮，眼睛閃爍
  phase2  (4幀) — 紅色憤怒狀態，快速震動
  death   (4幀) — 爆炸消散

輸出：
  assets/sprites/targets/B001_boss_sheet.png  (512x128, 4幀×4狀態×128px)
  assets/sprites/targets/B001_boss.png        (128x128, 靜態 idle 幀0，覆蓋原版)
"""

import os
import math
from PIL import Image, ImageDraw, ImageFilter, ImageEnhance

# ---- 設定 ----
FRAME_SIZE = 128
COLS = 4   # 每個狀態 4 幀
ROWS = 3   # idle / phase2 / death
OUT_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"

# ---- 官方吉伊卡哇顏色 ----
WHITE      = (255, 255, 247, 255)
OUTLINE    = (41, 42, 43, 255)
PINK       = (239, 165, 201, 255)
PINK_DARK  = (200, 100, 150, 255)
PINK_LIGHT = (255, 210, 230, 255)
GRAY_LIGHT = (220, 220, 215, 255)
GRAY_MID   = (180, 180, 175, 255)
GRAY_DARK  = (130, 130, 125, 255)
RED_EYE    = (255, 91, 86, 255)
RED_DARK   = (180, 30, 30, 255)
RED_GLOW   = (255, 60, 60, 180)
GOLD       = (255, 200, 50, 255)
GOLD_DARK  = (200, 140, 20, 255)
TRANSPARENT = (0, 0, 0, 0)


def px(img, x, y, color):
    """安全設定像素"""
    if 0 <= x < img.width and 0 <= y < img.height:
        img.putpixel((x, y), color)


def fill_circle(img, cx, cy, r, color, outline_color=None):
    """填充圓形（帶可選輪廓）"""
    if color is not None:
        for dy in range(-r - 1, r + 2):
            for dx in range(-r - 1, r + 2):
                dist = math.sqrt(dx * dx + dy * dy)
                if dist <= r:
                    px(img, cx + dx, cy + dy, color)
    if outline_color:
        for dy in range(-r - 2, r + 3):
            for dx in range(-r - 2, r + 3):
                dist = math.sqrt(dx * dx + dy * dy)
                if r < dist <= r + 1.2:
                    px(img, cx + dx, cy + dy, outline_color)


def fill_circle_shaded(img, cx, cy, r, base_color, light_dir=(-1, -1)):
    """帶陰影的圓形"""
    lx, ly = light_dir
    llen = math.sqrt(lx * lx + ly * ly)
    lx, ly = lx / llen, ly / llen

    for dy in range(-r - 1, r + 2):
        for dx in range(-r - 1, r + 2):
            dist = math.sqrt(dx * dx + dy * dy)
            if dist <= r:
                # 光照計算
                nx, ny = dx / max(dist, 0.1), dy / max(dist, 0.1)
                dot = -(nx * lx + ny * ly)
                light = 0.5 + 0.5 * max(0, dot)

                r_c = int(base_color[0] * (0.7 + 0.6 * light))
                g_c = int(base_color[1] * (0.7 + 0.6 * light))
                b_c = int(base_color[2] * (0.7 + 0.6 * light))
                r_c = min(255, r_c)
                g_c = min(255, g_c)
                b_c = min(255, b_c)
                px(img, cx + dx, cy + dy, (r_c, g_c, b_c, 255))


def draw_eye(img, cx, cy, is_open=True, is_angry=False):
    """繪製眼睛（3x3 眼白 + 2x2 瞳孔 + 高光）"""
    if not is_open:
        # 閉眼：一條橫線
        for dx in range(-2, 3):
            px(img, cx + dx, cy, OUTLINE)
        return

    eye_color = RED_EYE if is_angry else (60, 40, 30, 255)
    pupil_color = (20, 10, 10, 255) if is_angry else (10, 5, 5, 255)

    # 眼白
    for dy in range(-1, 2):
        for dx in range(-2, 3):
            px(img, cx + dx, cy + dy, (240, 240, 235, 255))

    # 瞳孔
    for dy in range(-1, 1):
        for dx in range(-1, 1):
            px(img, cx + dx, cy + dy, eye_color)

    # 高光
    px(img, cx - 1, cy - 1, (255, 255, 255, 255))

    # 輪廓
    for dy in range(-2, 3):
        for dx in range(-3, 4):
            dist = abs(dx) + abs(dy)
            if dist == 4 or (abs(dx) == 3 and abs(dy) <= 1) or (abs(dy) == 2 and abs(dx) <= 2):
                px(img, cx + dx, cy + dy, OUTLINE)


def draw_blush(img, cx, cy, r=4):
    """腮紅（橢圓漸層）"""
    for dy in range(-r // 2, r // 2 + 1):
        for dx in range(-r, r + 1):
            dist = math.sqrt((dx / r) ** 2 + (dy / (r // 2)) ** 2)
            if dist <= 1.0:
                alpha = int(120 * (1.0 - dist))
                px(img, cx + dx, cy + dy, (239, 165, 201, alpha))


def draw_boss_idle(frame_idx):
    """BOSS idle 幀 — 緩慢漂浮，偶爾眨眼"""
    img = Image.new("RGBA", (FRAME_SIZE, FRAME_SIZE), TRANSPARENT)

    # 漂浮偏移（上下搖擺）
    float_dy = int(math.sin(frame_idx * math.pi / 2) * 3)

    cx, cy = 64, 60 + float_dy

    # 身體（大圓，帶陰影）
    fill_circle_shaded(img, cx, cy, 38, (245, 245, 240))
    # 輪廓
    fill_circle(img, cx, cy, 38, None, OUTLINE)

    # 耳朵（圓形，左右各一）
    ear_r = 12
    fill_circle_shaded(img, cx - 28, cy - 28, ear_r, (240, 240, 235))
    fill_circle(img, cx - 28, cy - 28, ear_r, None, OUTLINE)
    fill_circle_shaded(img, cx + 28, cy - 28, ear_r, (240, 240, 235))
    fill_circle(img, cx + 28, cy - 28, ear_r, None, OUTLINE)

    # 耳朵內側（粉紅）
    fill_circle(img, cx - 28, cy - 28, 7, PINK)
    fill_circle(img, cx + 28, cy - 28, 7, PINK)

    # 眼睛（frame 2 閉眼，其他開眼）
    is_open = (frame_idx != 2)
    draw_eye(img, cx - 12, cy - 5, is_open=is_open)
    draw_eye(img, cx + 12, cy - 5, is_open=is_open)

    # 腮紅
    draw_blush(img, cx - 20, cy + 5)
    draw_blush(img, cx + 20, cy + 5)

    # 嘴巴（V 形微笑）
    px(img, cx - 2, cy + 12, OUTLINE)
    px(img, cx - 1, cy + 13, OUTLINE)
    px(img, cx,     cy + 14, OUTLINE)
    px(img, cx + 1, cy + 13, OUTLINE)
    px(img, cx + 2, cy + 12, OUTLINE)

    # 光暈（頂部高光）
    for dx in range(-15, 16):
        for dy in range(-15, 0):
            dist = math.sqrt(dx * dx + dy * dy)
            if 30 <= dist <= 38:
                alpha = int(60 * (1.0 - abs(dist - 34) / 4))
                if alpha > 0:
                    px(img, cx + dx, cy + dy, (255, 255, 255, alpha))

    return img


def draw_boss_phase2(frame_idx):
    """BOSS Phase 2 — 憤怒紅色，震動，眼睛變紅"""
    img = Image.new("RGBA", (FRAME_SIZE, FRAME_SIZE), TRANSPARENT)

    # 震動偏移
    shake_x = [0, 3, -3, 2][frame_idx]
    shake_y = [0, -2, 2, -1][frame_idx]

    cx, cy = 64 + shake_x, 60 + shake_y

    # 憤怒光暈（紅色外圈）
    for r in range(42, 48):
        alpha = int(150 * (1.0 - (r - 42) / 6))
        fill_circle(img, cx, cy, r, (255, 50, 50, alpha))

    # 身體（紅色調）
    fill_circle_shaded(img, cx, cy, 38, (255, 200, 200))
    fill_circle(img, cx, cy, 38, None, RED_DARK)

    # 耳朵（紅色調）
    ear_r = 12
    fill_circle_shaded(img, cx - 28, cy - 28, ear_r, (255, 190, 190))
    fill_circle(img, cx - 28, cy - 28, ear_r, None, RED_DARK)
    fill_circle_shaded(img, cx + 28, cy - 28, ear_r, (255, 190, 190))
    fill_circle(img, cx + 28, cy - 28, ear_r, None, RED_DARK)

    # 耳朵內側（深紅）
    fill_circle(img, cx - 28, cy - 28, 7, (220, 80, 80, 255))
    fill_circle(img, cx + 28, cy - 28, 7, (220, 80, 80, 255))

    # 憤怒眼睛（全部開眼，紅色）
    draw_eye(img, cx - 12, cy - 5, is_open=True, is_angry=True)
    draw_eye(img, cx + 12, cy - 5, is_open=True, is_angry=True)

    # 憤怒眉毛（斜線）
    for dx in range(-5, 1):
        px(img, cx - 12 + dx, cy - 10 - dx // 2, RED_DARK)
    for dx in range(0, 6):
        px(img, cx + 12 + dx, cy - 10 - (5 - dx) // 2, RED_DARK)

    # 腮紅（紅色）
    for dy in range(-2, 3):
        for dx in range(-4, 5):
            dist = math.sqrt((dx / 4) ** 2 + (dy / 2) ** 2)
            if dist <= 1.0:
                alpha = int(150 * (1.0 - dist))
                px(img, cx - 20 + dx, cy + 5 + dy, (255, 100, 100, alpha))
                px(img, cx + 20 + dx, cy + 5 + dy, (255, 100, 100, alpha))

    # 嘴巴（憤怒，倒 V）
    px(img, cx - 2, cy + 14, RED_DARK)
    px(img, cx - 1, cy + 13, RED_DARK)
    px(img, cx,     cy + 12, RED_DARK)
    px(img, cx + 1, cy + 13, RED_DARK)
    px(img, cx + 2, cy + 14, RED_DARK)

    # 閃爍效果（frame 1,3 加強紅色）
    if frame_idx % 2 == 1:
        enhancer = ImageEnhance.Color(img)
        img = enhancer.enhance(1.3)

    return img


def draw_boss_death(frame_idx):
    """BOSS death — 爆炸消散"""
    img = Image.new("RGBA", (FRAME_SIZE, FRAME_SIZE), TRANSPARENT)

    cx, cy = 64, 60

    # 消散進度（0.0 → 1.0）
    progress = frame_idx / 3.0

    # 縮小的身體
    body_r = int(38 * (1.0 - progress * 0.7))
    if body_r > 0:
        alpha = int(255 * (1.0 - progress))
        # 身體（逐漸透明）
        for dy in range(-body_r - 1, body_r + 2):
            for dx in range(-body_r - 1, body_r + 2):
                dist = math.sqrt(dx * dx + dy * dy)
                if dist <= body_r:
                    light = 0.7 + 0.3 * max(0, -(dx / max(dist, 1)) * 0.7 - (dy / max(dist, 1)) * 0.7)
                    r_c = min(255, int(245 * light))
                    g_c = min(255, int(245 * light))
                    b_c = min(255, int(240 * light))
                    px(img, cx + dx, cy + dy, (r_c, g_c, b_c, alpha))

    # 爆炸粒子（隨 progress 擴散）
    particle_count = 16
    for i in range(particle_count):
        angle = (i / particle_count) * 2 * math.pi
        dist_r = int(20 + progress * 50)
        px_x = cx + int(math.cos(angle) * dist_r)
        px_y = cy + int(math.sin(angle) * dist_r)
        alpha = int(255 * (1.0 - progress))
        size = max(1, int(4 * (1.0 - progress)))

        # 粒子顏色（交替金色和粉紅）
        if i % 2 == 0:
            color = (255, 200, 50, alpha)
        else:
            color = (239, 165, 201, alpha)

        for dy in range(-size, size + 1):
            for dx in range(-size, size + 1):
                if dx * dx + dy * dy <= size * size:
                    px(img, px_x + dx, px_y + dy, color)

    # 星形光芒（frame 0,1 最亮）
    if frame_idx <= 1:
        star_alpha = int(200 * (1.0 - frame_idx * 0.5))
        for angle_deg in range(0, 360, 45):
            angle = math.radians(angle_deg)
            for r in range(5, 35):
                sx = cx + int(math.cos(angle) * r)
                sy = cy + int(math.sin(angle) * r)
                ray_alpha = int(star_alpha * (1.0 - r / 35))
                px(img, sx, sy, (255, 255, 200, ray_alpha))

    return img


def generate_boss_sheet():
    """生成完整 BOSS spritesheet"""
    sheet_w = FRAME_SIZE * COLS
    sheet_h = FRAME_SIZE * ROWS
    sheet = Image.new("RGBA", (sheet_w, sheet_h), TRANSPARENT)

    print("生成 B001 BOSS 動畫集...")

    # Row 0: idle
    print("  [idle] 4幀...")
    for i in range(COLS):
        frame = draw_boss_idle(i)
        sheet.paste(frame, (i * FRAME_SIZE, 0))

    # Row 1: phase2
    print("  [phase2] 4幀...")
    for i in range(COLS):
        frame = draw_boss_phase2(i)
        sheet.paste(frame, (i * FRAME_SIZE, FRAME_SIZE))

    # Row 2: death
    print("  [death] 4幀...")
    for i in range(COLS):
        frame = draw_boss_death(i)
        sheet.paste(frame, (i * FRAME_SIZE, FRAME_SIZE * 2))

    # 儲存 spritesheet
    sheet_path = os.path.join(OUT_DIR, "B001_boss_sheet.png")
    sheet.save(sheet_path, "PNG")
    print(f"  ✅ Spritesheet: {sheet_path} ({sheet_w}x{sheet_h})")

    # 更新靜態 B001_boss.png（idle 幀0）
    idle_frame = draw_boss_idle(0)
    # 縮小到 64x64（與其他目標物一致）
    idle_small = idle_frame.resize((64, 64), Image.NEAREST)
    static_path = os.path.join(OUT_DIR, "B001_boss.png")
    idle_small.save(static_path, "PNG")
    print(f"  ✅ Static: {static_path} (64x64)")

    # 品質報告
    pixels = list(idle_frame.getdata())
    non_transparent = sum(1 for p in pixels if p[3] > 10)
    total = FRAME_SIZE * FRAME_SIZE
    print(f"\n品質報告:")
    print(f"  idle 幀0: {non_transparent}/{total} 非透明像素 ({non_transparent/total*100:.1f}%)")
    print(f"  動畫狀態: idle(4幀) + phase2(4幀) + death(4幀) = 12幀")
    print(f"  Spritesheet: {sheet_w}x{sheet_h}px")

    return sheet_path


if __name__ == "__main__":
    os.makedirs(OUT_DIR, exist_ok=True)
    generate_boss_sheet()
    print("\n✅ B001 BOSS 動畫集生成完成！")
