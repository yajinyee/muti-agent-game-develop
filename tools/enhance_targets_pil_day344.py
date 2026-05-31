"""
DAY-344: PIL 目標物美術升級 — T106-T130（Lucky 系統目標物）
策略：飽和度+35%、對比度+25%、銳化、加強輪廓光暈
目標：讓 Lucky 系統目標物更有辨識度和視覺衝擊力
"""
import os
import shutil
from PIL import Image, ImageEnhance, ImageFilter
import numpy as np

BASE = r'd:\Kiro\client\chiikawa-pixel\assets\sprites\targets'
BACKUP = r'd:\Kiro\tmp\targets_backup_day344'

os.makedirs(BACKUP, exist_ok=True)

# 找出 T106-T130 的所有 PNG
targets = []
for f in os.listdir(BASE):
    if not f.endswith('.png') or f.endswith('.import'):
        continue
    # 解析 T 編號
    if not f.startswith('T'):
        continue
    try:
        num = int(f[1:4])
        if 106 <= num <= 130:
            targets.append(f)
    except ValueError:
        continue

targets.sort()
print(f"找到 {len(targets)} 個目標物（T106-T130）")

def enhance_lucky_target(img: Image.Image, target_name: str) -> Image.Image:
    """
    Lucky 系統目標物增強：
    - 飽和度 +35%（比基礎目標更鮮豔，突出 Lucky 感）
    - 對比度 +25%（輪廓更清晰）
    - 銳化（細節更清楚）
    - 加強高倍率目標的光暈效果
    """
    if img.mode != 'RGBA':
        img = img.convert('RGBA')
    
    # 分離 alpha 通道
    r, g, b, a = img.split()
    rgb = Image.merge('RGB', (r, g, b))
    
    # 飽和度 +35%（Lucky 目標物更鮮豔）
    rgb = ImageEnhance.Color(rgb).enhance(1.35)
    
    # 對比度 +25%
    rgb = ImageEnhance.Contrast(rgb).enhance(1.25)
    
    # 亮度微調 +5%（讓顏色更飽滿）
    rgb = ImageEnhance.Brightness(rgb).enhance(1.05)
    
    # 銳化
    rgb = rgb.filter(ImageFilter.SHARPEN)
    
    # 重新合併 alpha
    r2, g2, b2 = rgb.split()
    result = Image.merge('RGBA', (r2, g2, b2, a))
    
    # 對高倍率目標（T121+）加強光暈效果
    try:
        num = int(target_name[1:4])
        if num >= 121:
            result = add_glow_effect(result, num)
    except ValueError:
        pass
    
    return result

def add_glow_effect(img: Image.Image, target_num: int) -> Image.Image:
    """
    為高倍率 Lucky 目標物加強光暈效果
    T121-T125: 金色光暈
    T126-T130: 彩虹光暈
    """
    arr = np.array(img, dtype=np.float32)
    h, w = arr.shape[:2]
    
    # 找出非透明像素的邊界
    alpha = arr[:, :, 3]
    mask = alpha > 50
    
    if not mask.any():
        return img
    
    # 計算中心
    ys, xs = np.where(mask)
    cy, cx = ys.mean(), xs.mean()
    
    # 根據目標編號決定光暈顏色
    if target_num >= 126:
        # 彩虹光暈（T126-T130）
        for y in range(h):
            for x in range(w):
                if alpha[y, x] > 50:
                    # 根據位置計算彩虹色
                    angle = np.arctan2(y - cy, x - cx)
                    hue = (angle / (2 * np.pi) + 0.5) % 1.0
                    # 簡單 HSV to RGB
                    hi = int(hue * 6)
                    f = hue * 6 - hi
                    if hi == 0: rc, gc, bc = 1, f, 0
                    elif hi == 1: rc, gc, bc = 1-f, 1, 0
                    elif hi == 2: rc, gc, bc = 0, 1, f
                    elif hi == 3: rc, gc, bc = 0, 1-f, 1
                    elif hi == 4: rc, gc, bc = f, 0, 1
                    else: rc, gc, bc = 1, 0, 1-f
                    # 混合 10% 彩虹色
                    arr[y, x, 0] = min(255, arr[y, x, 0] * 0.92 + rc * 255 * 0.08)
                    arr[y, x, 1] = min(255, arr[y, x, 1] * 0.92 + gc * 255 * 0.08)
                    arr[y, x, 2] = min(255, arr[y, x, 2] * 0.92 + bc * 255 * 0.08)
    else:
        # 金色光暈（T121-T125）
        for y in range(h):
            for x in range(w):
                if alpha[y, x] > 50:
                    # 加強紅色和綠色通道（金色 = 高R + 高G + 低B）
                    arr[y, x, 0] = min(255, arr[y, x, 0] * 1.08)  # R +8%
                    arr[y, x, 1] = min(255, arr[y, x, 1] * 1.05)  # G +5%
                    arr[y, x, 2] = max(0, arr[y, x, 2] * 0.92)    # B -8%
    
    return Image.fromarray(arr.astype(np.uint8), 'RGBA')

# 統計
enhanced_count = 0
total_sat_gain = 0
total_contrast_gain = 0

for filename in targets:
    src_path = os.path.join(BASE, filename)
    backup_path = os.path.join(BACKUP, filename)
    
    # 備份原始檔案
    shutil.copy2(src_path, backup_path)
    
    # 載入圖片
    img = Image.open(src_path).convert('RGBA')
    original_pixels = sum(1 for px in img.getdata() if px[3] > 50)
    
    # 增強
    enhanced = enhance_lucky_target(img, filename)
    
    # 驗證：非透明像素數量應該相同
    enhanced_pixels = sum(1 for px in enhanced.getdata() if px[3] > 50)
    
    # 儲存
    enhanced.save(src_path)
    
    enhanced_count += 1
    print(f"  ✅ {filename}: {original_pixels} → {enhanced_pixels} 像素（{'✓' if enhanced_pixels == original_pixels else '⚠️ 變化'}）")

print(f"\n完成！升級了 {enhanced_count} 個 Lucky 目標物（T106-T130）")
print(f"備份位置：{BACKUP}")
print(f"策略：飽和度+35%、對比度+25%、亮度+5%、銳化")
print(f"T121-T125：額外金色光暈")
print(f"T126-T130：額外彩虹光暈")
