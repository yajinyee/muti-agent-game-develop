"""
DAY-346: T191-T220 目標物美術升級
超高階 Lucky 系統目標物 — 最強視覺效果

分層光暈策略：
- T191-T200：深紫/洋紅光暈（PvP/宇宙感，R+8%, B+18%, G-6%）
- T201-T210：金橙/火焰光暈（能量/爆炸感，R+15%, G+5%, B-12%）
- T211-T220：冰藍/電光光暈（冰雪/電擊感，B+18%, G+8%, R-8%）

增強參數（最強等級）：
- 飽和度 +50%（最鮮豔）
- 對比度 +40%（最清晰）
- 亮度 +12%（最飽滿）
- 四重銳化（最清楚）
"""
import os
import sys
import shutil
from PIL import Image, ImageEnhance, ImageFilter

PYTHON_PATH = r"C:\Users\yajinyee0306\AppData\Local\Programs\Python\Python312\python.exe"
TARGETS_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
BACKUP_DIR = r"d:\Kiro\tmp\targets_backup_day346"

# 分層光暈設定
GLOW_CONFIGS = {
    # T191-T200: 深紫/洋紅光暈（PvP/宇宙感）
    range(191, 201): {"r": 8, "g": -6, "b": 18, "name": "深紫宇宙"},
    # T201-T210: 金橙/火焰光暈（能量/爆炸感）
    range(201, 211): {"r": 15, "g": 5, "b": -12, "name": "金橙火焰"},
    # T211-T220: 冰藍/電光光暈（冰雪/電擊感）
    range(211, 221): {"r": -8, "g": 8, "b": 18, "name": "冰藍電光"},
}

def get_glow_config(target_num):
    for r, cfg in GLOW_CONFIGS.items():
        if target_num in r:
            return cfg
    return {"r": 0, "g": 0, "b": 0, "name": "無光暈"}

def enhance_target(png_path, target_num):
    """增強單個目標物圖片"""
    try:
        img = Image.open(png_path).convert("RGBA")
    except Exception as e:
        print(f"  ❌ 無法讀取 {os.path.basename(png_path)}: {e}")
        return False
    
    # 分離 RGB 和 Alpha
    r, g, b, a = img.split()
    rgb = Image.merge("RGB", (r, g, b))
    
    # 1. 飽和度 +50%
    rgb = ImageEnhance.Color(rgb).enhance(1.50)
    
    # 2. 對比度 +40%
    rgb = ImageEnhance.Contrast(rgb).enhance(1.40)
    
    # 3. 亮度 +12%
    rgb = ImageEnhance.Brightness(rgb).enhance(1.12)
    
    # 4. 四重銳化
    for _ in range(4):
        rgb = rgb.filter(ImageFilter.SHARPEN)
    
    # 5. 分層光暈（逐像素調整）
    glow = get_glow_config(target_num)
    if glow["r"] != 0 or glow["g"] != 0 or glow["b"] != 0:
        pixels = rgb.load()
        alpha_pixels = a.load()
        w, h = rgb.size
        for y in range(h):
            for x in range(w):
                if alpha_pixels[x, y] > 10:  # 只處理非透明像素
                    pr, pg, pb = pixels[x, y]
                    pr = max(0, min(255, pr + glow["r"] * 2))
                    pg = max(0, min(255, pg + glow["g"] * 2))
                    pb = max(0, min(255, pb + glow["b"] * 2))
                    pixels[x, y] = (pr, pg, pb)
    
    # 合併回 RGBA
    r2, g2, b2 = rgb.split()
    result = Image.merge("RGBA", (r2, g2, b2, a))
    
    # 儲存
    result.save(png_path, "PNG")
    return True

def main():
    # 建立備份目錄
    os.makedirs(BACKUP_DIR, exist_ok=True)
    
    # 找出 T191-T220 的 PNG 檔案
    targets_to_enhance = []
    for num in range(191, 221):
        pattern = f"T{num}_"
        for fname in os.listdir(TARGETS_DIR):
            if fname.startswith(pattern) and fname.endswith(".png") and not fname.endswith(".import"):
                targets_to_enhance.append((num, os.path.join(TARGETS_DIR, fname), fname))
                break
    
    print(f"找到 {len(targets_to_enhance)} 個目標物需要升級")
    print()
    
    # 先備份
    print("備份原始檔案...")
    for num, path, fname in targets_to_enhance:
        backup_path = os.path.join(BACKUP_DIR, fname)
        shutil.copy2(path, backup_path)
    print(f"  ✅ 備份完成 → {BACKUP_DIR}")
    print()
    
    # 驗證 PNG 完整性
    print("驗證 PNG 完整性...")
    broken = []
    for num, path, fname in targets_to_enhance:
        try:
            img = Image.open(path)
            img.load()
        except Exception as e:
            broken.append((num, path, fname, str(e)))
            print(f"  ❌ 損壞: {fname} — {e}")
    
    if broken:
        print(f"\n發現 {len(broken)} 個損壞的 PNG，需要先修復")
        # 嘗試重新生成損壞的 PNG（用簡單的佔位圖）
        for num, path, fname, err in broken:
            print(f"  🔧 重新生成: {fname}")
            # 建立一個簡單的 64x64 佔位圖
            placeholder = Image.new("RGBA", (64, 64), (0, 0, 0, 0))
            # 畫一個彩色圓形作為佔位
            from PIL import ImageDraw
            draw = ImageDraw.Draw(placeholder)
            # 根據目標編號選擇顏色
            if num in range(191, 201):
                color = (180, 50, 220, 255)  # 紫色
            elif num in range(201, 211):
                color = (220, 120, 30, 255)  # 橙色
            else:
                color = (50, 150, 220, 255)  # 藍色
            draw.ellipse([8, 8, 56, 56], fill=color, outline=(255, 255, 255, 200))
            # 加文字
            draw.text((20, 25), f"T{num}", fill=(255, 255, 255, 255))
            placeholder.save(path, "PNG")
            print(f"    ✅ 已重新生成佔位圖")
    
    print()
    
    # 執行增強
    print("開始增強...")
    success = 0
    fail = 0
    
    for num, path, fname in targets_to_enhance:
        glow = get_glow_config(num)
        print(f"  T{num} ({glow['name']})...", end=" ")
        if enhance_target(path, num):
            print("✅")
            success += 1
        else:
            print("❌")
            fail += 1
    
    print()
    print(f"完成！成功 {success} 個，失敗 {fail} 個")
    
    # 統計結果
    print()
    print("=== 增強結果統計 ===")
    for num, path, fname in targets_to_enhance:
        try:
            img = Image.open(path).convert("RGBA")
            r, g, b, a = img.split()
            alpha_arr = list(a.getdata())
            non_transparent = sum(1 for v in alpha_arr if v > 10)
            total = len(alpha_arr)
            pct = non_transparent / total * 100
            glow = get_glow_config(num)
            print(f"  T{num}: {pct:.1f}% 非透明 [{glow['name']}]")
        except Exception as e:
            print(f"  T{num}: 讀取失敗 — {e}")

if __name__ == "__main__":
    main()
