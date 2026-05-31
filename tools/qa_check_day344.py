"""
DAY-344 QA 驗證腳本
驗證項目：
1. T106-T130 目標物美術升級（25個）
2. Server PlayerUpdatePayload 新增 DisplayName 欄位
3. GameManager.gd 新增 get_display_name() 方法
4. Cannon.gd 玩家名稱顯示邏輯
5. Server 編譯狀態
"""
import os
import subprocess
import sys
from PIL import Image

BASE = r'd:\Kiro\client\chiikawa-pixel\assets\sprites\targets'
BACKUP = r'd:\Kiro\tmp\targets_backup_day344'
SERVER = r'd:\Kiro\server'
CLIENT_SCRIPTS = r'd:\Kiro\client\chiikawa-pixel\scripts'

passed = 0
failed = 0

def check(name, condition, detail=""):
    global passed, failed
    if condition:
        print(f"  ✅ {name}")
        passed += 1
    else:
        print(f"  ❌ {name}" + (f" — {detail}" if detail else ""))
        failed += 1

print("=" * 60)
print("DAY-344 QA 驗證")
print("=" * 60)

# ── 1. T106-T130 目標物美術升級 ──────────────────────────────
print("\n[1] T106-T130 目標物美術升級（25個）")
target_files = []
for f in os.listdir(BASE):
    if not f.endswith('.png') or f.endswith('.import'):
        continue
    if not f.startswith('T'):
        continue
    try:
        num = int(f[1:4])
        if 106 <= num <= 130:
            target_files.append(f)
    except ValueError:
        continue

check("找到 25 個 T106-T130 目標物", len(target_files) == 25, f"找到 {len(target_files)} 個")

for f in sorted(target_files):
    src = os.path.join(BASE, f)
    bak = os.path.join(BACKUP, f)
    # 確認備份存在
    check(f"{f} 備份存在", os.path.exists(bak))
    # 確認圖片可以開啟
    try:
        img = Image.open(src)
        check(f"{f} 圖片格式正確", img.mode == 'RGBA', f"mode={img.mode}")
    except Exception as e:
        check(f"{f} 圖片可開啟", False, str(e))

# ── 2. Server PlayerUpdatePayload 新增 DisplayName ──────────
print("\n[2] Server PlayerUpdatePayload DisplayName 欄位")
messages_go = os.path.join(SERVER, 'internal', 'protocol', 'messages.go')
with open(messages_go, 'r', encoding='utf-8') as f:
    content = f.read()
check("PlayerUpdatePayload 有 DisplayName 欄位", 'DisplayName    string  `json:"display_name"`' in content)
check("DAY-344 注釋存在", 'DAY-344' in content)

# ── 3. Server game.go 填充 DisplayName ──────────────────────
print("\n[3] Server game.go 填充 DisplayName")
game_go = os.path.join(SERVER, 'internal', 'game', 'game.go')
with open(game_go, 'r', encoding='utf-8') as f:
    game_content = f.read()
check("game.go 填充 DisplayName", 'DisplayName:     p.GetDisplayName()' in game_content)

# ── 4. GameManager.gd 新增 get_display_name ─────────────────
print("\n[4] GameManager.gd get_display_name 方法")
gm_path = os.path.join(CLIENT_SCRIPTS, 'game', 'GameManager.gd')
with open(gm_path, 'r', encoding='utf-8') as f:
    gm_content = f.read()
check("GameManager.gd 有 get_display_name()", 'func get_display_name()' in gm_content)
check("get_display_name 讀取 display_name 欄位", '"display_name"' in gm_content)

# ── 5. Cannon.gd 玩家名稱顯示 ───────────────────────────────
print("\n[5] Cannon.gd 玩家名稱顯示")
cannon_path = os.path.join(CLIENT_SCRIPTS, 'game', 'Cannon.gd')
with open(cannon_path, 'r', encoding='utf-8') as f:
    cannon_content = f.read()
check("Cannon.gd 有 display_name 邏輯", 'display_name' in cannon_content)
check("Cannon.gd 有名稱格式化", '[" + char_name + "]"' in cannon_content)

# ── 6. Server 編譯 ───────────────────────────────────────────
print("\n[6] Server 編譯驗證")
result = subprocess.run(['go', 'build', './...'], cwd=SERVER, capture_output=True, text=True)
check("go build ./... 通過", result.returncode == 0, result.stderr[:200] if result.stderr else "")
result2 = subprocess.run(['go', 'vet', './...'], cwd=SERVER, capture_output=True, text=True)
check("go vet ./... 通過", result2.returncode == 0, result2.stderr[:200] if result2.stderr else "")

# ── 結果 ─────────────────────────────────────────────────────
print("\n" + "=" * 60)
total = passed + failed
print(f"結果：{passed}/{total} 通過")
if failed == 0:
    print("✅ 全部通過！DAY-344 QA 完成")
else:
    print(f"❌ {failed} 項失敗，需要修復")
print("=" * 60)

sys.exit(0 if failed == 0 else 1)
