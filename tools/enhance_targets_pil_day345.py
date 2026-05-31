"""
DAY-345: PIL 目標物美術升級 — T131-T160（中高階 Lucky 系統目標物）
策略：飽和度+40%、對比度+30%、銳化、分層光暈效果
目標：讓中高階 Lucky 目標物有更強的視覺衝擊力和辨識度
- T131-T140: 深藍/紫色光暈（神秘感）
- T141-T150: 橙紅色光暈（火焰感）
- T151-T160: 青綠色光暈（電光感）
"""
import os
import shutil
from PIL import Image, ImageEnhance, ImageFilter
import numpy as np

BASE = r'd:\Kiro\client\chiikawa-pixel\assets\sprites\targets'
BACKUP = r'd:\Kiro\tmp\targets_backup_day345'

os.makedirs(BACKUP, exist_ok=True)

# 找出 T131-T160 的所有 PNG
targets = []
for f in os.listdir(BASE):
    if not f.endswith('.png') or f.endswith('.import'):
        continue
    if not f.startswith('T'):
        continue
    try:
        num = int(f[1:4])
        if 131 <= num <= 160:
            targets.append(f)
    except ValueError:
        continue

targets.sort()
print(f"找到 {len(targets)} 個目標物（T131-T160）")

def add_tier_glow(img: Image.Image, target_num: int) -> Image.Image:
    """
    分層光暈效果：
    T131-T140: 深藍/紫色光暈（神秘高階感）
    T141-T150: 橙紅色光暈（火焰爆發感）
    T151-T160: 青綠色光暈（電光科技感）
    """
    arr = np.array(img, dtype=np.float32)
    alpha = arr[:, :, 3]
    mask = alpha > 50

    if not mask.any():
        return img

    if 131 <= target_num <= 140:
        # 深藍/紫色光暈：B+12%, R+6%, G-5%
        arr[:, :, 0] = np.where(mask, np.clip(arr[:, :, 0] * 1.06, 0, 255), arr[:, :, 0])
        arr[:, :, 1] = np.where(mask, np.clip(arr[:, :, 1] * 0.95, 0, 255), arr[:, :, 1])
        arr[:, :, 2] = np.where(mask, np.clip(arr[:, :, 2] * 1.12, 0, 255), arr[:, :, 2])
    elif 141 <= target_num <= 150:
        # 橙紅色光暈：R+12%, G+3%, B-10%
        arr[:, :, 0] = np.where(mask, np.clip(arr[:, :, 0] * 1.12, 0, 255), arr[:, :, 0])
        arr[:, :, 1] = np.where(mask, np.clip(arr[:, :, 1] * 1.03, 0, 255), arr[:, :, 1])
        arr[:, :, 2] = np.where(mask, np.clip(arr[:, :, 2] * 0.90, 0, 255), arr[:, :, 2])
    elif 151 <= target_num <= 160:
        # 青綠色光暈：G+10%, B+8%, R-5%
        arr[:, :, 0] = np.where(mask, np.clip(arr[:, :, 0] * 0.95, 0, 255), arr[:, :, 0])
        arr[:, :, 1] = np.where(mask, np.clip(arr[:, :, 1] * 1.10, 0, 255), arr[:, :, 1])
        arr[:, :, 2] = np.where(mask, np.clip(arr[:, :, 2] * 1.08, 0, 255), arr[:, :, 2])

    return Image.fromarray(arr.astype(np.uint8), 'RGBA')

def enhance_high_tier_target(img: Image.Image, target_name: str) -> Image.Image:
    """
    中高階 Lucky 目標物增強：
    - 飽和度 +40%（比 T106-T130 更強，突出高階感）
    - 對比度 +30%（輪廓更清晰）
    - 亮度 +8%（顏色更飽滿）
    - 銳化（細節更清楚）
    - 分層光暈（依目標編號決定顏色）
    """
    if img.mode != 'RGBA':
        img = img.convert('RGBA')

    r, g, b, a = img.split()
    rgb = Image.merge('RGB', (r, g, b))

    # 飽和度 +40%（高階目標更鮮豔）
    rgb = ImageEnhance.Color(rgb).enhance(1.40)

    # 對比度 +30%
    rgb = ImageEnhance.Contrast(rgb).enhance(1.30)

    # 亮度 +8%
    rgb = ImageEnhance.Brightness(rgb).enhance(1.08)

    # 銳化（兩次，讓細節更清楚）
    rgb = rgb.filter(ImageFilter.SHARPEN)
    rgb = rgb.filter(ImageFilter.SHARPEN)

    # 重新合併 alpha
    r2, g2, b2 = rgb.split()
    result = Image.merge('RGBA', (r2, g2, b2, a))

    # 加分層光暈
    try:
        num = int(target_name[1:4])
        result = add_tier_glow(result, num)
    except ValueError:
        pass

    return result

# 執行升級
enhanced_count = 0

for filename in targets:
    src_path = os.path.join(BASE, filename)
    backup_path = os.path.join(BACKUP, filename)

    # 備份
    shutil.copy2(src_path, backup_path)

    # 載入
    img = Image.open(src_path).convert('RGBA')
    original_pixels = sum(1 for px in img.getdata() if px[3] > 50)

    # 增強
    enhanced = enhance_high_tier_target(img, filename)

    # 驗證
    enhanced_pixels = sum(1 for px in enhanced.getdata() if px[3] > 50)

    # 儲存
    enhanced.save(src_path)

    enhanced_count += 1
    tier = "藍紫" if int(filename[1:4]) <= 140 else ("橙紅" if int(filename[1:4]) <= 150 else "青綠")
    print(f"  ✅ {filename} [{tier}光暈]: {original_pixels} → {enhanced_pixels} 像素")

print(f"\n完成！升級了 {enhanced_count} 個中高階 Lucky 目標物（T131-T160）")
print(f"備份位置：{BACKUP}")
print(f"策略：飽和度+40%、對比度+30%、亮度+8%、雙重銳化")
print(f"T131-T140：深藍/紫色光暈（神秘感）")
print(f"T141-T150：橙紅色光暈（火焰感）")
print(f"T151-T160：青綠色光暈（電光感）")
