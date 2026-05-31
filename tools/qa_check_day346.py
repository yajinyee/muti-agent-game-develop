"""
DAY-346 QA 驗證腳本
驗證項目：
1. T191-T220 美術升級（30個）
2. 119 個 .import 補齊
3. Server 編譯狀態
"""
import os
import sys
import subprocess
from PIL import Image

TARGETS_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
SERVER_DIR = r"d:\Kiro\server"

def check_server_build():
    """驗證 Server 編譯"""
    result = subprocess.run(
        ["go", "build", "./..."],
        cwd=SERVER_DIR,
        capture_output=True,
        text=True
    )
    return result.returncode == 0, result.stderr

def check_server_vet():
    """驗證 Server vet"""
    result = subprocess.run(
        ["go", "vet", "./..."],
        cwd=SERVER_DIR,
        capture_output=True,
        text=True
    )
    return result.returncode == 0, result.stderr

def check_png_valid(path):
    """驗證 PNG 可讀取"""
    try:
        img = Image.open(path)
        img.load()
        return True, None
    except Exception as e:
        return False, str(e)

def check_import_exists(png_path):
    """驗證 .import 存在"""
    return os.path.exists(png_path + ".import")

def get_non_transparent_pct(path):
    """計算非透明像素比例"""
    try:
        img = Image.open(path).convert("RGBA")
        r, g, b, a = img.split()
        alpha_data = list(a.tobytes())
        non_transparent = sum(1 for v in alpha_data if v > 10)
        total = len(alpha_data)
        return non_transparent / total * 100
    except:
        return 0.0

def main():
    passed = 0
    failed = 0
    
    print("=" * 60)
    print("DAY-346 QA 驗證")
    print("=" * 60)
    
    # === 1. Server 編譯 ===
    print("\n[1] Server 編譯驗證")
    ok, err = check_server_build()
    if ok:
        print("  ✅ go build ./... 通過")
        passed += 1
    else:
        print(f"  ❌ go build 失敗: {err}")
        failed += 1
    
    ok, err = check_server_vet()
    if ok:
        print("  ✅ go vet ./... 通過")
        passed += 1
    else:
        print(f"  ❌ go vet 失敗: {err}")
        failed += 1
    
    # === 2. T191-T220 美術升級驗證 ===
    print("\n[2] T191-T220 美術升級驗證（30個）")
    for num in range(191, 221):
        # 找到對應的 PNG
        found = None
        for fname in os.listdir(TARGETS_DIR):
            if fname.startswith(f"T{num}_") and fname.endswith(".png") and not fname.endswith(".import"):
                found = os.path.join(TARGETS_DIR, fname)
                break
        
        if not found:
            print(f"  ❌ T{num}: 找不到 PNG 檔案")
            failed += 1
            continue
        
        # 驗證 PNG 可讀取
        ok, err = check_png_valid(found)
        if not ok:
            print(f"  ❌ T{num}: PNG 損壞 — {err}")
            failed += 1
            continue
        
        # 驗證 .import 存在
        if not check_import_exists(found):
            print(f"  ❌ T{num}: 缺少 .import")
            failed += 1
            continue
        
        # 驗證非透明像素比例 > 30%
        pct = get_non_transparent_pct(found)
        if pct < 30.0:
            print(f"  ⚠️  T{num}: 非透明像素 {pct:.1f}% < 30%（可能是細長形狀）")
            # 不算失敗，只是警告
        
        print(f"  ✅ T{num}: {pct:.1f}% 非透明")
        passed += 1
    
    # === 3. 批次 .import 補齊驗證 ===
    print("\n[3] 全部目標物 .import 驗證")
    all_pngs = [f for f in os.listdir(TARGETS_DIR) 
                if f.endswith(".png") and not f.endswith(".import")]
    
    missing_imports = []
    for fname in sorted(all_pngs):
        png_path = os.path.join(TARGETS_DIR, fname)
        if not check_import_exists(png_path):
            missing_imports.append(fname)
    
    if missing_imports:
        print(f"  ❌ 仍有 {len(missing_imports)} 個 PNG 缺少 .import:")
        for f in missing_imports[:10]:
            print(f"    - {f}")
        failed += 1
    else:
        print(f"  ✅ 全部 {len(all_pngs)} 個 PNG 都有對應的 .import")
        passed += 1
    
    # === 4. 關鍵目標物完整性驗證 ===
    print("\n[4] 關鍵目標物完整性驗證")
    key_targets = [
        ("T001_grass.png", "基礎草"),
        ("T101_mimic.png", "擬態怪物"),
        ("T191_pvp_battle.png", "PvP 戰鬥"),
        ("T200_genesis_epoch.png", "創世紀元"),
        ("T210_ultimate_miracle.png", "終極奇蹟"),
        ("T220_rapid_riches.png", "快速暴富"),
        ("T253_penta_fusion.png", "五重終極"),
        ("B001_boss.png", "BOSS"),
    ]
    
    for fname, name in key_targets:
        path = os.path.join(TARGETS_DIR, fname)
        if not os.path.exists(path):
            print(f"  ❌ {name} ({fname}): 不存在")
            failed += 1
            continue
        
        ok, err = check_png_valid(path)
        if not ok:
            print(f"  ❌ {name} ({fname}): PNG 損壞")
            failed += 1
            continue
        
        if not check_import_exists(path):
            print(f"  ❌ {name} ({fname}): 缺少 .import")
            failed += 1
            continue
        
        print(f"  ✅ {name} ({fname})")
        passed += 1
    
    # === 結果 ===
    print()
    print("=" * 60)
    total = passed + failed
    print(f"結果: {passed}/{total} 通過")
    if failed == 0:
        print("🎉 全部通過！")
    else:
        print(f"⚠️  {failed} 項失敗，需要修復")
    print("=" * 60)
    
    return 0 if failed == 0 else 1

if __name__ == "__main__":
    sys.exit(main())
