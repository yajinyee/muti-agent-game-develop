"""
enhance_targets_proper_pixel_day343.py — DAY-343 用 proper-pixel-art 升級目標物美術
把現有的程式生成目標物圖轉換為更乾淨的像素藝術

使用 proper-pixel-art 套件：
  pip install proper-pixel-art
  
算法：Canny 邊緣偵測 → Hough 變換找格線 → 量化顏色 → 每格取最常見顏色
"""
import os
import sys
from PIL import Image

try:
    from proper_pixel_art.pixelate import pixelate
    print("✅ proper-pixel-art 套件已安裝")
except ImportError:
    print("❌ 請先安裝：pip install proper-pixel-art")
    sys.exit(1)

TARGETS_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
BACKUP_DIR = r"d:\Kiro\tmp\targets_backup_day343"

# 要升級的目標物（最常出現的 10 個基礎目標物）
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

def enhance_target(filename: str, num_colors: int = 16) -> bool:
    """用 proper-pixel-art 升級單個目標物"""
    input_path = os.path.join(TARGETS_DIR, filename)
    if not os.path.exists(input_path):
        print(f"  ⚠️ 找不到：{filename}")
        return False
    
    # 備份原始檔案
    os.makedirs(BACKUP_DIR, exist_ok=True)
    backup_path = os.path.join(BACKUP_DIR, filename)
    
    try:
        img = Image.open(input_path).convert("RGBA")
        original_size = img.size
        original_pixels = sum(1 for p in img.getdata() if p[3] > 0)
        
        # 備份
        img.save(backup_path)
        
        # 用 proper-pixel-art 處理
        result = pixelate(
            img,
            num_colors=num_colors,
            transparent_background=False,  # 不做透明背景處理（圖片已有透明背景）
            scale_result=1  # 保持原始大小
        )
        
        # 確保輸出大小與原始一致
        if result.size != original_size:
            result = result.resize(original_size, Image.NEAREST)
        
        # 計算改善程度
        result_pixels = sum(1 for p in result.getdata() if p[3] > 0)
        
        # 儲存
        result.save(input_path)
        
        improvement = "✅" if result_pixels >= original_pixels * 0.8 else "⚠️"
        print(f"  {improvement} {filename}: {original_pixels}px → {result_pixels}px ({num_colors} 色)")
        return True
        
    except Exception as e:
        print(f"  ❌ {filename}: {e}")
        # 還原備份
        if os.path.exists(backup_path):
            import shutil
            shutil.copy2(backup_path, input_path)
        return False

def main():
    print("=== DAY-343 proper-pixel-art 目標物升級 ===")
    print(f"目標目錄：{TARGETS_DIR}")
    print(f"備份目錄：{BACKUP_DIR}")
    print()
    
    success = 0
    failed = 0
    
    for filename in PRIORITY_TARGETS:
        print(f"處理 {filename}...")
        if enhance_target(filename, num_colors=16):
            success += 1
        else:
            failed += 1
    
    print(f"\n=== 完成 ===")
    print(f"成功：{success}/{len(PRIORITY_TARGETS)}")
    print(f"失敗：{failed}/{len(PRIORITY_TARGETS)}")
    
    if failed > 0:
        print(f"\n備份位置：{BACKUP_DIR}")
        print("如需還原，請從備份目錄複製回來")

if __name__ == "__main__":
    main()
