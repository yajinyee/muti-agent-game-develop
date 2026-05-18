# -*- coding: utf-8 -*-
"""
為角色 idle 動畫加入眨眼效果
在 8 幀 idle 中，第 5-6 幀加入眼睛閉合效果
讓角色更有生命感

技術：
- 分析角色眼睛位置（找深色像素群）
- 在特定幀用水平線覆蓋眼睛（模擬閉眼）
- 重建 Spritesheet
"""
import os
import json
import numpy as np
from PIL import Image, ImageEnhance
import math

CHARS_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\characters"
SHEET_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\sheets"
FRAME_SIZE = 96

# 各角色的眼睛顏色（深色輪廓）
EYE_COLORS = {
    "chiikawa":  (41, 42, 43),   # 深黑色輪廓
    "hachiware": (41, 42, 43),
    "usagi":     (17, 17, 17),   # 更深的黑
}

# 各角色的眼睛大致位置（相對於 96x96 畫布）
# 這些是估計值，程式會自動微調
EYE_REGIONS = {
    "chiikawa":  {"y_range": (28, 45), "x_range": (20, 75)},
    "hachiware": {"y_range": (28, 45), "x_range": (20, 75)},
    "usagi":     {"y_range": (28, 45), "x_range": (20, 75)},
}


def find_eye_pixels(img_arr, char_name):
    """找到眼睛像素的位置（深色像素群）"""
    region = EYE_REGIONS.get(char_name, {"y_range": (25, 50), "x_range": (15, 80)})
    y_min, y_max = region["y_range"]
    x_min, x_max = region["x_range"]
    
    eye_pixels = []
    for y in range(y_min, y_max):
        for x in range(x_min, x_max):
            if y >= img_arr.shape[0] or x >= img_arr.shape[1]:
                continue
            r, g, b, a = img_arr[y, x]
            if a < 50:  # 透明像素跳過
                continue
            # 深色像素 = 眼睛輪廓
            if r < 80 and g < 80 and b < 80:
                eye_pixels.append((x, y))
    
    return eye_pixels


def get_eye_rows(eye_pixels):
    """取得眼睛所在的行（y 座標）"""
    if not eye_pixels:
        return []
    y_coords = [p[1] for p in eye_pixels]
    y_min, y_max = min(y_coords), max(y_coords)
    return list(range(y_min, y_max + 1))


def apply_blink(img_arr, eye_pixels, blink_amount):
    """
    套用眨眼效果
    blink_amount: 0.0 = 全開, 1.0 = 全閉
    """
    if not eye_pixels or blink_amount <= 0:
        return img_arr.copy()
    
    result = img_arr.copy()
    
    # 找眼睛的 y 範圍
    y_coords = sorted(set(p[1] for p in eye_pixels))
    if not y_coords:
        return result
    
    y_min, y_max = y_coords[0], y_coords[-1]
    eye_height = y_max - y_min + 1
    
    # 要遮蓋的行數（從上往下）
    rows_to_cover = int(eye_height * blink_amount)
    
    # 找眼睛的 x 範圍
    x_coords = [p[0] for p in eye_pixels]
    x_min, x_max = min(x_coords), max(x_coords)
    
    # 用角色的皮膚色（白色）覆蓋眼睛上半部
    # 找眼睛周圍的皮膚色
    skin_color = None
    for y in range(max(0, y_min - 3), y_min):
        for x in range(x_min, x_max + 1):
            if y < img_arr.shape[0] and x < img_arr.shape[1]:
                r, g, b, a = img_arr[y, x]
                if a > 100 and r > 200 and g > 200 and b > 200:
                    skin_color = (r, g, b, a)
                    break
        if skin_color:
            break
    
    if skin_color is None:
        skin_color = (255, 255, 247, 255)  # 預設白色
    
    # 覆蓋眼睛上半部
    for i in range(rows_to_cover):
        y = y_min + i
        if y >= result.shape[0]:
            break
        for x in range(x_min - 1, x_max + 2):
            if 0 <= x < result.shape[1]:
                # 只覆蓋有像素的地方
                if result[y, x, 3] > 50:
                    result[y, x] = skin_color
    
    # 在眼睛中間加一條深色線（閉眼線）
    if rows_to_cover > 0:
        close_y = y_min + rows_to_cover - 1
        if close_y < result.shape[0]:
            for x in range(x_min, x_max + 1):
                if 0 <= x < result.shape[1] and result[close_y, x, 3] > 50:
                    result[close_y, x] = (41, 42, 43, 255)
    
    return result


def gen_idle_8frames_with_blink(base_img, char_name):
    """
    8幀 idle 動畫（正弦波呼吸感 + 眨眼）
    幀0-3: 正常（0°-135°）
    幀4: 半閉眼（180°，開始眨眼）
    幀5: 全閉眼（225°）
    幀6: 半閉眼（270°，睜開中）
    幀7: 正常（315°）
    """
    base_arr = np.array(base_img)
    eye_pixels = find_eye_pixels(base_arr, char_name)
    
    print(f"  找到 {len(eye_pixels)} 個眼睛像素")
    
    # 眨眼時機（幀索引 -> 眨眼程度）
    blink_schedule = {
        0: 0.0,   # 全開
        1: 0.0,   # 全開
        2: 0.0,   # 全開
        3: 0.0,   # 全開
        4: 0.5,   # 半閉（開始眨眼）
        5: 1.0,   # 全閉
        6: 0.5,   # 半閉（睜開中）
        7: 0.0,   # 全開
    }
    
    frames = []
    
    for i in range(8):
        angle = i * math.pi / 4
        sin_val = math.sin(angle)
        
        # 位移：最大 ±2px
        offset_y = int(-sin_val * 2)
        
        # 縮放：1.0 到 1.02
        scale = 1.0 + max(0, sin_val) * 0.02
        
        # 先做縮放/位移
        frame = Image.new("RGBA", (FRAME_SIZE, FRAME_SIZE), (0, 0, 0, 0))
        
        if abs(scale - 1.0) > 0.001:
            new_size = int(FRAME_SIZE * scale)
            scaled = base_img.resize((new_size, new_size), Image.NEAREST)
            paste_x = (FRAME_SIZE - new_size) // 2
            paste_y = (FRAME_SIZE - new_size) // 2 + offset_y
            paste_x = max(0, min(FRAME_SIZE - new_size, paste_x))
            paste_y = max(0, min(FRAME_SIZE - new_size, paste_y))
            frame.paste(scaled, (paste_x, paste_y))
        else:
            if offset_y != 0:
                temp = Image.new("RGBA", (FRAME_SIZE, FRAME_SIZE), (0, 0, 0, 0))
                temp.paste(base_img, (0, offset_y))
                frame = temp
            else:
                frame = base_img.copy()
        
        # 套用眨眼效果
        blink_amount = blink_schedule.get(i, 0.0)
        if blink_amount > 0 and len(eye_pixels) > 0:
            frame_arr = np.array(frame)
            # 調整眼睛像素位置（考慮位移）
            adjusted_eye_pixels = [(x, y + offset_y) for x, y in eye_pixels
                                   if 0 <= y + offset_y < FRAME_SIZE]
            frame_arr = apply_blink(frame_arr, adjusted_eye_pixels, blink_amount)
            frame = Image.fromarray(frame_arr)
        
        frames.append(frame)
    
    return frames


def gen_attack_3frames(base_img):
    """3幀 attack"""
    frames = []
    
    # 幀0：舉棒準備
    f0 = base_img.rotate(-18, expand=False, fillcolor=(0,0,0,0))
    f0_canvas = Image.new("RGBA", (FRAME_SIZE, FRAME_SIZE), (0,0,0,0))
    f0_canvas.paste(f0, (0, -3))
    f0_bright = ImageEnhance.Brightness(f0_canvas).enhance(1.12)
    frames.append(f0_bright)
    
    # 幀1：揮下衝擊
    f1 = base_img.rotate(12, expand=False, fillcolor=(0,0,0,0))
    f1_arr = np.array(f1)
    for i in range(8):
        x = FRAME_SIZE - 6 - i*3
        y = 6 + i*3
        if 0 <= x < FRAME_SIZE and 0 <= y < FRAME_SIZE:
            alpha = max(0, 220 - i*25)
            f1_arr[y, x] = [255, 150, 200, alpha]
            if x+1 < FRAME_SIZE:
                f1_arr[y, x+1] = [255, 180, 220, alpha//2]
    frames.append(Image.fromarray(f1_arr))
    
    # 幀2：收回
    f2 = Image.new("RGBA", (FRAME_SIZE, FRAME_SIZE), (0,0,0,0))
    f2.paste(base_img, (0, 1))
    frames.append(f2)
    
    return frames


def gen_bigwin_4frames(base_img):
    """4幀 bigwin"""
    frames = []
    configs = [
        (0,   1.00, False),
        (-10, 1.05, False),
        (-14, 1.08, True),
        (-5,  1.02, False),
    ]
    
    for offset, scale, add_stars in configs:
        frame = Image.new("RGBA", (FRAME_SIZE, FRAME_SIZE), (0, 0, 0, 0))
        if scale != 1.0:
            new_size = int(FRAME_SIZE * scale)
            scaled = base_img.resize((new_size, new_size), Image.NEAREST)
            paste_x = (FRAME_SIZE - new_size) // 2
            paste_y = (FRAME_SIZE - new_size) // 2 + offset
            paste_x = max(0, min(FRAME_SIZE - new_size, paste_x))
            paste_y = max(0, min(FRAME_SIZE - new_size, paste_y))
            frame.paste(scaled, (paste_x, paste_y))
        else:
            paste_y = max(0, offset)
            frame.paste(base_img, (0, paste_y))
        
        if add_stars:
            frame_arr = np.array(frame)
            star_positions = [
                (8, 8), (FRAME_SIZE-10, 8),
                (6, FRAME_SIZE//3), (FRAME_SIZE-8, FRAME_SIZE//3)
            ]
            for sx, sy in star_positions:
                if 0 <= sx < FRAME_SIZE and 0 <= sy < FRAME_SIZE:
                    frame_arr[sy, sx] = [255, 235, 50, 255]
                    for dx, dy in [(1,0),(-1,0),(0,1),(0,-1)]:
                        nx, ny = sx+dx, sy+dy
                        if 0 <= nx < FRAME_SIZE and 0 <= ny < FRAME_SIZE:
                            frame_arr[ny, nx] = [255, 235, 50, 180]
            frame = Image.fromarray(frame_arr)
        
        frames.append(frame)
    
    return frames


def build_spritesheet_with_blink(char_name):
    """建立帶眨眼效果的 8×3 Spritesheet"""
    idle_path = os.path.join(CHARS_DIR, f"{char_name}_idle.png")
    if not os.path.exists(idle_path):
        print(f"  SKIP: {idle_path} not found")
        return False
    
    base = Image.open(idle_path).convert("RGBA")
    if base.width != FRAME_SIZE:
        base = base.resize((FRAME_SIZE, FRAME_SIZE), Image.NEAREST)
    
    print(f"  生成 idle 8幀（含眨眼）...")
    idle_frames   = gen_idle_8frames_with_blink(base, char_name)
    print(f"  生成 attack 3幀...")
    attack_frames = gen_attack_3frames(base)
    print(f"  生成 bigwin 4幀...")
    bigwin_frames = gen_bigwin_4frames(base)
    
    # 建立 Spritesheet（8 cols × 3 rows）
    COLS = 8
    ROWS = 3
    sheet = Image.new("RGBA", (FRAME_SIZE * COLS, FRAME_SIZE * ROWS), (0, 0, 0, 0))
    
    all_frames = [idle_frames, attack_frames, bigwin_frames]
    for row, frames in enumerate(all_frames):
        for col, frame in enumerate(frames[:COLS]):
            sheet.paste(frame, (col * FRAME_SIZE, row * FRAME_SIZE))
    
    # 儲存
    sheet_path = os.path.join(SHEET_DIR, f"{char_name}_animated.png")
    sheet.save(sheet_path)
    
    # 更新 metadata（保持不變，CharacterAnimator.gd 不需要修改）
    meta = {
        "char": char_name,
        "frame_size": FRAME_SIZE,
        "cols": COLS,
        "rows": ROWS,
        "animations": {
            "idle":   {"row": 0, "frames": 8, "fps": 8.0},
            "attack": {"row": 1, "frames": 3, "fps": 8.0},
            "bigwin": {"row": 2, "frames": 4, "fps": 6.0},
        }
    }
    meta_path = os.path.join(SHEET_DIR, f"{char_name}_animated.json")
    with open(meta_path, "w", encoding="utf-8") as f:
        json.dump(meta, f, indent=2)
    
    print(f"  ✅ {char_name}: {sheet.width}x{sheet.height} 已儲存（含眨眼效果）")
    return True


def main():
    print("=== 角色 idle 動畫升級：加入眨眼效果 ===")
    print("技術：第 5-6 幀（180°-270°）加入眼睛閉合效果\n")
    
    for char in ["chiikawa", "hachiware", "usagi"]:
        print(f"[{char}]")
        build_spritesheet_with_blink(char)
        print()
    
    print("✅ 完成！眨眼效果已加入所有角色的 idle 動畫")
    print("   CharacterAnimator.gd 不需要修改（格式相同）")


if __name__ == "__main__":
    main()
