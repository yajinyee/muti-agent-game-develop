"""
DAY-347: T221-T253 目標物美術升級
最高階 Lucky 系統（T221-T253）— 最鮮豔、最強烈、最有辨識度

升級策略：
- 飽和度 +55%（最強等級，超越 T191-T220 的 +50%）
- 對比度 +45%（輪廓最清晰）
- 亮度 +15%（顏色最飽滿）
- 五重銳化（細節最清楚）
- 分段光暈：
  - T221-T230：深紫/洋紅光暈（宇宙/神話感，R+10%, G-8%, B+20%）
  - T231-T240：金橙/火焰光暈（能量/爆炸感，R+18%, G+8%, B-15%）
  - T241-T253：冰藍/電光光暈（冰雪/電擊感，R-10%, G+10%, B+22%）
"""

import os
import shutil
from PIL import Image, ImageEnhance, ImageFilter
import numpy as np

SPRITES_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
BACKUP_DIR = r"d:\Kiro\data\tmp\targets_backup_day347"

# T221-T253 的實際檔案名稱對照表
TARGET_NAMES = {
    221: "T221_dice_bonus",
    222: "T222_dual_bonus",
    223: "T223_coin_respin",
    224: "T224_golden_pot",
    225: "T225_cascade_lock",
    226: "T226_legend_awaken",
    227: "T227_crash_harvest",
    228: "T228_cosmic_fusion",
    229: "T229_magnetic_attraction",
    230: "T230_super_chain",
    231: "T231_holy_pillar",
    232: "T232_time_stop",
    233: "T233_cosmic_restart",
    234: "T234_fever_boost_ultimate",
    235: "T235_rapid_riches_ultimate",
    236: "T236_ice_fishing_master",
    237: "T237_cosmic_miracle",
    238: "T238_genesis_ultimate",
    239: "T239_shark_spark",
    240: "T240_winter_ice",
    241: "T241_atlantis_frenzy",
    242: "T242_fishing_time_wheel",
    243: "T243_ultimate_shark",
    244: "T244_wild_collector",
    245: "T245_lightning_eel_ultra",
    246: "T246_domino_chain",
    247: "T247_immortal_boss_ultra",
    248: "T248_quad_fusion",
    249: "T249_electrical_frame",
    250: "T250_magnetic_respin",
    251: "T251_fisherman_trail",
    252: "T252_golden_gills",
    253: "T253_penta_fusion",
}

def backup_targets():
    """備份 T221-T253"""
    os.makedirs(BACKUP_DIR, exist_ok=True)
    backed_up = 0
    for i, name in TARGET_NAMES.items():
        src = os.path.join(SPRITES_DIR, f"{name}.png")
        if os.path.exists(src):
            dst = os.path.join(BACKUP_DIR, f"{name}.png")
            shutil.copy2(src, dst)
            backed_up += 1
    print(f"備份完成：{backed_up} 個檔案 → {BACKUP_DIR}")
    return backed_up

def get_halo_params(target_id):
    """依目標 ID 決定光暈參數"""
    if 221 <= target_id <= 230:
        # 深紫/洋紅光暈（宇宙/神話感）
        return {"r": 10, "g": -8, "b": 20, "name": "深紫洋紅"}
    elif 231 <= target_id <= 240:
        # 金橙/火焰光暈（能量/爆炸感）
        return {"r": 18, "g": 8, "b": -15, "name": "金橙火焰"}
    else:  # 241-253
        # 冰藍/電光光暈（冰雪/電擊感）
        return {"r": -10, "g": 10, "b": 22, "name": "冰藍電光"}

def enhance_target(img_path, target_id):
    """升級單個目標物"""
    img = Image.open(img_path).convert("RGBA")
    
    # 分離 alpha 通道
    r, g, b, a = img.split()
    rgb = Image.merge("RGB", (r, g, b))
    
    # 1. 飽和度 +55%（最強等級）
    rgb = ImageEnhance.Color(rgb).enhance(1.55)
    
    # 2. 對比度 +45%
    rgb = ImageEnhance.Contrast(rgb).enhance(1.45)
    
    # 3. 亮度 +15%
    rgb = ImageEnhance.Brightness(rgb).enhance(1.15)
    
    # 4. 五重銳化
    for _ in range(5):
        rgb = rgb.filter(ImageFilter.SHARPEN)
    
    # 5. 光暈色調調整
    halo = get_halo_params(target_id)
    arr = np.array(rgb, dtype=np.float32)
    
    # 只對非透明像素應用光暈
    alpha_arr = np.array(a)
    mask = alpha_arr > 10
    
    arr[:, :, 0] = np.clip(arr[:, :, 0] + np.where(mask, halo["r"] * 2.5, 0), 0, 255)
    arr[:, :, 1] = np.clip(arr[:, :, 1] + np.where(mask, halo["g"] * 2.5, 0), 0, 255)
    arr[:, :, 2] = np.clip(arr[:, :, 2] + np.where(mask, halo["b"] * 2.5, 0), 0, 255)
    
    rgb = Image.fromarray(arr.astype(np.uint8))
    
    # 合併回 RGBA
    result = Image.merge("RGBA", (*rgb.split(), a))
    return result

def process_all():
    """處理 T221-T253"""
    print("=== DAY-347 T221-T253 美術升級 ===")
    
    # 備份
    backed_up = backup_targets()
    if backed_up == 0:
        print("警告：沒有找到任何目標物檔案！")
    
    success = 0
    failed = 0
    
    for i, name in TARGET_NAMES.items():
        img_path = os.path.join(SPRITES_DIR, f"{name}.png")
        if not os.path.exists(img_path):
            print(f"  跳過 {name}（檔案不存在）")
            continue
        
        try:
            result = enhance_target(img_path, i)
            result.save(img_path, "PNG")
            
            halo = get_halo_params(i)
            print(f"  ✅ {name} 升級完成（{halo['name']}光暈）")
            success += 1
        except Exception as e:
            print(f"  ❌ {name} 失敗：{e}")
            failed += 1
    
    print(f"\n=== 完成 ===")
    print(f"成功：{success} 個")
    print(f"失敗：{failed} 個")
    print(f"備份位置：{BACKUP_DIR}")
    
    return success, failed

if __name__ == "__main__":
    process_all()
