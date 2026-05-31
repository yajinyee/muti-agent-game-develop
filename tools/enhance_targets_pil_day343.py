"""
enhance_targets_pil_day343.py — DAY-343 用 PIL 增強目標物美術
proper-pixel-art 對 64x64 小圖效果不好，改用 PIL 直接增強

策略：
1. 提升飽和度（讓顏色更鮮豔）
2. 提升對比度（讓輪廓更清晰）
3. 銳化（讓邊緣更清晰）
4. 保持透明背景
"""
import os
import sys
import shutil
from PIL import Image, ImageEnhance, ImageFilter
import numpy as np

TARGETS_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
BACKUP_DIR = r"d:\Kiro\tmp\targets_backup_day343_pil"

# 要升級的目標物（最常出現的基礎目標物）
PRIORITY_TARGETS = [
    "T001_grass.png",
    "T002_bug_g.png",
    "T003_bug_r.png",
    "T004_bug_b.png",
    "T005_pudding.png",
    "T006_mushroom.png",
    "T101_mimic.png",
    "T102_chest.png",
    "T103_meteor.png",
    "T104_gold_grass.png",
    "T105_coin_fish.png",
]

def enhance_target_pil(filename: str) -> bool:
    """用 PIL 增強單個目標物"""
    input_path = os.path.join(TARGETS_DIR, filename)
    if not os.path.exists(input_path):
        print(f"  ⚠️ 找不到：{filename}")
        return False
    
    # 備份
    os.makedirs(BACKUP_DIR, exist_ok=True)
    backup_path = os.path.join(BACKUP_DIR, filename)
    shutil.copy2(input_path, backup_path)
    
    try:
        img = Image.open(input_path).convert("RGBA")
        
        # 分離 RGB 和 Alpha 通道
        r, g, b, a = img.split()
        rgb = Image.merge("RGB", (r, g, b))
        
        # 計算原始非透明像素數
        alpha_arr = np.array(a)
        original_pixels = np.sum(alpha_arr > 0)
        
        # 1. 提升飽和度（1.3x）
        rgb = ImageEnhance.Color(rgb).enhance(1.3)
        
        # 2. 提升對比度（1.2x）
        rgb = ImageEnhance.Contrast(rgb).enhance(1.2)
        
        # 3. 銳化（讓邊緣更清晰）
        rgb = rgb.filter(ImageFilter.SHARPEN)
        
        # 4. 重新合併 Alpha 通道（保持透明背景不變）
        r2, g2, b2 = rgb.split()
        result = Image.merge("RGBA", (r2, g2, b2, a))
        
        # 計算結果非透明像素數
        result_pixels = original_pixels  # Alpha 不變，所以像素數相同
        
        # 儲存
        result.save(input_path)
        
        print(f"  ✅ {filename}: {original_pixels}px（飽和度+30%，對比度+20%，銳化）")
        return True
        
    except Exception as e:
        print(f"  ❌ {filename}: {e}")
        # 還原備份
        if os.path.exists(backup_path):
            shutil.copy2(backup_path, input_path)
        return False

def analyze_quality(filename: str) -> dict:
    """分析圖片品質指標"""
    path = os.path.join(TARGETS_DIR, filename)
    if not os.path.exists(path):
        return {}
    
    img = Image.open(path).convert("RGBA")
    arr = np.array(img)
    
    # 非透明像素
    alpha = arr[:, :, 3]
    mask = alpha > 0
    pixels = np.sum(mask)
    
    if pixels == 0:
        return {"pixels": 0, "saturation": 0, "contrast": 0}
    
    # 計算飽和度（HSV 的 S 通道）
    rgb = arr[:, :, :3][mask]
    r, g, b = rgb[:, 0] / 255.0, rgb[:, 1] / 255.0, rgb[:, 2] / 255.0
    max_c = np.maximum(np.maximum(r, g), b)
    min_c = np.minimum(np.minimum(r, g), b)
    saturation = np.mean(max_c - min_c)
    
    # 計算對比度（標準差）
    gray = 0.299 * r + 0.587 * g + 0.114 * b
    contrast = np.std(gray)
    
    return {
        "pixels": int(pixels),
        "saturation": float(saturation),
        "contrast": float(contrast),
    }

def main():
    print("=== DAY-343 PIL 目標物美術升級 ===")
    print()
    
    # 先分析原始品質
    print("原始品質分析：")
    original_quality = {}
    for filename in PRIORITY_TARGETS:
        q = analyze_quality(filename)
        if q:
            original_quality[filename] = q
            print(f"  {filename}: {q['pixels']}px, 飽和度={q['saturation']:.3f}, 對比度={q['contrast']:.3f}")
    
    print()
    print("開始升級...")
    
    success = 0
    failed = 0
    
    for filename in PRIORITY_TARGETS:
        if enhance_target_pil(filename):
            success += 1
        else:
            failed += 1
    
    print()
    print("升級後品質分析：")
    for filename in PRIORITY_TARGETS:
        q = analyze_quality(filename)
        orig = original_quality.get(filename, {})
        if q and orig:
            sat_change = (q['saturation'] - orig['saturation']) / max(orig['saturation'], 0.001) * 100
            con_change = (q['contrast'] - orig['contrast']) / max(orig['contrast'], 0.001) * 100
            print(f"  {filename}: 飽和度 {orig['saturation']:.3f}→{q['saturation']:.3f} ({sat_change:+.0f}%), 對比度 {orig['contrast']:.3f}→{q['contrast']:.3f} ({con_change:+.0f}%)")
    
    print(f"\n=== 完成 ===")
    print(f"成功：{success}/{len(PRIORITY_TARGETS)}")
    print(f"備份位置：{BACKUP_DIR}")

if __name__ == "__main__":
    main()
