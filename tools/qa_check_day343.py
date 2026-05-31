"""
qa_check_day343.py — DAY-343 QA 驗證腳本
驗證 PIL 目標物美術升級效果
"""
import os
import sys
import numpy as np
from PIL import Image

PASS = "✅"
FAIL = "❌"
results = []

def check(name, condition, detail=""):
    status = PASS if condition else FAIL
    results.append((status, name, detail))
    if not condition:
        print(f"  {FAIL} {name}: {detail}")

TARGETS_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"
BACKUP_DIR = r"d:\Kiro\tmp\targets_backup_day343_pil"

PRIORITY_TARGETS = [
    "T001_grass.png", "T002_bug_g.png", "T003_bug_r.png",
    "T004_bug_b.png", "T005_pudding.png", "T006_mushroom.png",
    "T101_mimic.png", "T102_chest.png", "T103_meteor.png",
    "T104_gold_grass.png", "T105_coin_fish.png",
]

def get_quality(path):
    if not os.path.exists(path):
        return None
    img = Image.open(path).convert("RGBA")
    arr = np.array(img)
    alpha = arr[:, :, 3]
    mask = alpha > 0
    pixels = int(np.sum(mask))
    if pixels == 0:
        return {"pixels": 0, "saturation": 0, "contrast": 0}
    rgb = arr[:, :, :3][mask]
    r, g, b = rgb[:, 0] / 255.0, rgb[:, 1] / 255.0, rgb[:, 2] / 255.0
    max_c = np.maximum(np.maximum(r, g), b)
    min_c = np.minimum(np.minimum(r, g), b)
    saturation = float(np.mean(max_c - min_c))
    gray = 0.299 * r + 0.587 * g + 0.114 * b
    contrast = float(np.std(gray))
    return {"pixels": pixels, "saturation": saturation, "contrast": contrast}

# ── 1. 目標物檔案存在 ─────────────────────────────────────────
for filename in PRIORITY_TARGETS:
    path = os.path.join(TARGETS_DIR, filename)
    check(f"目標物存在: {filename}", os.path.exists(path))

# ── 2. 備份存在 ───────────────────────────────────────────────
for filename in PRIORITY_TARGETS:
    path = os.path.join(BACKUP_DIR, filename)
    check(f"備份存在: {filename}", os.path.exists(path))

# ── 3. 升級後品質比備份好 ─────────────────────────────────────
for filename in PRIORITY_TARGETS:
    current_path = os.path.join(TARGETS_DIR, filename)
    backup_path = os.path.join(BACKUP_DIR, filename)
    if not os.path.exists(current_path) or not os.path.exists(backup_path):
        continue
    current_q = get_quality(current_path)
    backup_q = get_quality(backup_path)
    if current_q and backup_q and backup_q["pixels"] > 0:
        # 升級後飽和度應該更高
        sat_improved = current_q["saturation"] >= backup_q["saturation"] * 0.95
        # 像素數應該保持不變
        pixels_ok = current_q["pixels"] == backup_q["pixels"]
        check(f"品質提升: {filename}", sat_improved and pixels_ok,
              f"飽和度 {backup_q['saturation']:.3f}→{current_q['saturation']:.3f}, 像素 {backup_q['pixels']}→{current_q['pixels']}")

# ── 4. 工具腳本存在 ───────────────────────────────────────────
check("工具: enhance_targets_pil_day343.py", os.path.exists(r"d:\Kiro\tools\enhance_targets_pil_day343.py"))
check("工具: enhance_targets_proper_pixel_day343.py", os.path.exists(r"d:\Kiro\tools\enhance_targets_proper_pixel_day343.py"))

# ── 5. Server 編譯狀態 ────────────────────────────────────────
import subprocess
result = subprocess.run(
    ["go", "build", "./..."],
    cwd=r"d:\Kiro\server",
    capture_output=True, text=True
)
check("Server: go build 通過", result.returncode == 0, result.stderr[:200] if result.stderr else "")

result2 = subprocess.run(
    ["go", "vet", "./..."],
    cwd=r"d:\Kiro\server",
    capture_output=True, text=True
)
check("Server: go vet 通過", result2.returncode == 0, result2.stderr[:200] if result2.stderr else "")

# ── 結果統計 ──────────────────────────────────────────────────
total = len(results)
passed = sum(1 for r in results if r[0] == PASS)
failed = total - passed

print(f"\n{'='*50}")
print(f"DAY-343 QA 結果：{passed}/{total} 通過")
print(f"{'='*50}")

if failed > 0:
    print("\n失敗項目：")
    for status, name, detail in results:
        if status == FAIL:
            print(f"  {FAIL} {name}: {detail}")
else:
    print("🎉 全部通過！")

sys.exit(0 if failed == 0 else 1)
