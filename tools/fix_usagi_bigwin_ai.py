"""
fix_usagi_bigwin_ai.py — 修復 usagi bigwin AI 生成圖的一致性問題（DAY-115）
從 usagi idle 幀生成 bigwin（上移 + 金色色調 + 星星），確保尺寸一致
"""
from PIL import Image, ImageEnhance
import os

CHARS_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\characters"
TARGET_SIZE = 96

def add_gold_tint(img: Image.Image, strength: float = 0.3) -> Image.Image:
    """加金色色調"""
    pixels = img.load()
    w, h = img.size
    result = img.copy()
    rpx = result.load()
    for y in range(h):
        for x in range(w):
            r, g, b, a = pixels[x, y]
            if a > 10:
                # 金色 = 提高 R 和 G，降低 B
                nr = min(255, int(r + (255 - r) * strength * 0.8))
                ng = min(255, int(g + (255 - g) * strength * 0.5))
                nb = max(0, int(b * (1 - strength * 0.4)))
                rpx[x, y] = (nr, ng, nb, a)
    return result

def add_stars(img: Image.Image) -> Image.Image:
    """加星星光點"""
    result = img.copy()
    pixels = result.load()
    w, h = img.size
    
    # 在角色周圍加幾個小星星
    star_positions = [
        (15, 10), (75, 8), (10, 50), (80, 45), (20, 80), (70, 75)
    ]
    star_color = (255, 240, 100, 200)
    
    for sx, sy in star_positions:
        if 0 <= sx < w and 0 <= sy < h:
            # 十字星
            for dx, dy in [(0,0), (1,0), (-1,0), (0,1), (0,-1)]:
                nx, ny = sx+dx, sy+dy
                if 0 <= nx < w and 0 <= ny < h:
                    pixels[nx, ny] = star_color
    return result

def main():
    idle_path = os.path.join(CHARS_DIR, "usagi_idle.png")
    bigwin_out = os.path.join(CHARS_DIR, "usagi_bigwin.png")
    
    # 從 idle 幀生成 bigwin
    idle = Image.open(idle_path).convert("RGBA")
    
    # 上移 4px（跳起感）
    shifted = Image.new("RGBA", (TARGET_SIZE, TARGET_SIZE), (0, 0, 0, 0))
    shifted.paste(idle, (0, -4))
    
    # 加金色色調
    gold = add_gold_tint(shifted, 0.25)
    
    # 加星星
    result = add_stars(gold)
    
    result.save(bigwin_out)
    
    # 驗證
    saved = Image.open(bigwin_out).convert("RGBA")
    pixels = saved.load()
    non_trans = sum(1 for y in range(TARGET_SIZE) for x in range(TARGET_SIZE) if pixels[x, y][3] > 10)
    
    # 計算 bbox
    min_x, min_y = TARGET_SIZE, TARGET_SIZE
    max_x, max_y = 0, 0
    for y in range(TARGET_SIZE):
        for x in range(TARGET_SIZE):
            if pixels[x, y][3] > 10:
                min_x = min(min_x, x)
                min_y = min(min_y, y)
                max_x = max(max_x, x)
                max_y = max(max_y, y)
    
    bbox_w = max_x - min_x + 1 if max_x >= min_x else 0
    bbox_h = max_y - min_y + 1 if max_y >= min_y else 0
    
    print(f"✅ usagi_bigwin: {non_trans}/{TARGET_SIZE*TARGET_SIZE} ({100*non_trans//(TARGET_SIZE*TARGET_SIZE)}%)")
    print(f"   bbox: {bbox_w}x{bbox_h}")
    
    # 比較 idle 的 bbox
    idle_pixels = idle.load()
    idle_min_x, idle_min_y = TARGET_SIZE, TARGET_SIZE
    idle_max_x, idle_max_y = 0, 0
    for y in range(TARGET_SIZE):
        for x in range(TARGET_SIZE):
            if idle_pixels[x, y][3] > 10:
                idle_min_x = min(idle_min_x, x)
                idle_min_y = min(idle_min_y, y)
                idle_max_x = max(idle_max_x, x)
                idle_max_y = max(idle_max_y, y)
    
    idle_bbox_w = idle_max_x - idle_min_x + 1
    idle_bbox_h = idle_max_y - idle_min_y + 1
    
    print(f"   idle bbox: {idle_bbox_w}x{idle_bbox_h}")
    print(f"   height diff: {abs(bbox_h - idle_bbox_h)}px")
    print(f"   width diff: {abs(bbox_w - idle_bbox_w)}px")

if __name__ == "__main__":
    main()
