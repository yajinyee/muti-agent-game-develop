#!/usr/bin/env python3
"""
QA 驗證腳本 — DAY-299
驗證 Server + Client 的完整性和一致性
"""
import os
import re
import sys

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
SERVER_GAME = os.path.join(ROOT, "server", "internal", "game")
SERVER_DATA = os.path.join(ROOT, "server", "internal", "data")
SERVER_PROTOCOL = os.path.join(ROOT, "server", "internal", "protocol")
CLIENT_SCRIPTS = os.path.join(ROOT, "client", "chiikawa-pixel", "scripts")
CLIENT_SPRITES = os.path.join(ROOT, "client", "chiikawa-pixel", "assets", "sprites", "targets")

PASS = "✅"
FAIL = "❌"
WARN = "⚠️"

results = []

def check(name, condition, detail=""):
    status = PASS if condition else FAIL
    results.append((status, name, detail))
    return condition

def warn(name, condition, detail=""):
    status = PASS if condition else WARN
    results.append((status, name, detail))
    return condition

# ── 1. Server Lucky Handler 數量 ─────────────────────────────
lucky_handlers = [f for f in os.listdir(SERVER_GAME) if f.startswith("lucky_") and f.endswith("_handler.go")]
check("Server Lucky Handler 數量 == 20", len(lucky_handlers) == 20, f"實際: {len(lucky_handlers)}")

# ── 2. Server tables.go 目標物數量 ───────────────────────────
with open(os.path.join(SERVER_DATA, "tables.go"), encoding="utf-8") as f:
    tables_content = f.read()
target_ids = re.findall(r'"(T\d{3}|B\d{3})"', tables_content)
target_ids = list(set(target_ids))
check("Server 目標物數量 >= 32", len(target_ids) >= 32, f"實際: {len(target_ids)} 個")

# ── 3. Protocol 訊息類型 ─────────────────────────────────────
with open(os.path.join(SERVER_PROTOCOL, "messages.go"), encoding="utf-8") as f:
    protocol_content = f.read()
lucky_msgs = re.findall(r'MsgLucky\w+\s*=\s*"lucky_\w+"', protocol_content)
check("Protocol Lucky 訊息類型 == 20", len(lucky_msgs) == 20, f"實際: {len(lucky_msgs)}")

# ── 4. Client 精靈圖完整性 ───────────────────────────────────
expected_sprites = [f"T{i:03d}" for i in range(1, 7)] + \
                   [f"T{i:03d}" for i in range(101, 126)] + ["B001"]
missing_sprites = []
for tid in expected_sprites:
    files = [f for f in os.listdir(CLIENT_SPRITES) if f.startswith(tid + "_")]
    if not files:
        missing_sprites.append(tid)
check("Client 精靈圖完整性", len(missing_sprites) == 0,
      f"缺少: {missing_sprites}" if missing_sprites else "全部存在")

# ── 5. GameManager.gd Lucky 訊號 ────────────────────────────
gm_path = os.path.join(CLIENT_SCRIPTS, "game", "GameManager.gd")
with open(gm_path, encoding="utf-8") as f:
    gm_content = f.read()
lucky_signals = re.findall(r'signal lucky_\w+', gm_content)
check("GameManager Lucky 訊號 == 20", len(lucky_signals) == 20, f"實際: {len(lucky_signals)}")

# ── 6. HUD.gd Lucky 事件處理函數 ────────────────────────────
hud_path = os.path.join(CLIENT_SCRIPTS, "ui", "HUD.gd")
with open(hud_path, encoding="utf-8") as f:
    hud_content = f.read()
lucky_handlers_client = re.findall(r'func _on_lucky_\w+', hud_content)
check("HUD Lucky 事件處理函數 == 20", len(lucky_handlers_client) == 20, f"實際: {len(lucky_handlers_client)}")

# ── 7. HUD.gd 重複函數檢查 ───────────────────────────────────
show_lucky_banner_count = hud_content.count("func _show_lucky_banner")
check("HUD _show_lucky_banner 無重複定義", show_lucky_banner_count == 1,
      f"定義次數: {show_lucky_banner_count}")

# ── 8. LuckyEventSystem.gd 主題配置 ─────────────────────────
les_path = os.path.join(CLIENT_SCRIPTS, "ui", "LuckyEventSystem.gd")
with open(les_path, encoding="utf-8") as f:
    les_content = f.read()
lucky_configs = re.findall(r'"(chain_lightning|crab_torpedo|vortex|golden_dragon|thunder_lobster|awakened_phoenix|shockwave_bomb|drill_torpedo|time_freeze|chain_explosion|chain_long_king|dragon_shotgun|rocket_cannon|deep_whirlpool|vampire_mult|mirror_fish|golden_rain|freeze_bomb|thunder_storm|lucky_wheel)":', les_content)
lucky_configs = list(set(lucky_configs))
check("LuckyEventSystem 主題配置 == 20", len(lucky_configs) == 20, f"實際: {len(lucky_configs)}")

# ── 9. TargetManager.gd T106-T125 映射 ──────────────────────
tm_path = os.path.join(CLIENT_SCRIPTS, "game", "TargetManager.gd")
with open(tm_path, encoding="utf-8") as f:
    tm_content = f.read()
t106_125 = [f'"T{i}"' for i in range(106, 126)]
missing_tm = [t for t in t106_125 if t not in tm_content]
check("TargetManager T106-T125 映射完整", len(missing_tm) == 0,
      f"缺少: {missing_tm}" if missing_tm else "全部存在")

# ── 10. NetworkManager.gd collect_golden_coin ───────────────
nm_path = os.path.join(CLIENT_SCRIPTS, "network", "NetworkManager.gd")
with open(nm_path, encoding="utf-8") as f:
    nm_content = f.read()
check("NetworkManager send_collect_golden_coin 已實作",
      "send_collect_golden_coin" in nm_content)

# ── 11. game.go Lucky Handler 整合 ──────────────────────────
with open(os.path.join(SERVER_GAME, "game.go"), encoding="utf-8") as f:
    game_content = f.read()
lucky_cases = re.findall(r'case isLucky\w+Fish', game_content)
check("game.go Lucky Handler 整合 >= 20", len(lucky_cases) >= 20, f"實際: {len(lucky_cases)}")

# ── 12. 音效資產完整性 ───────────────────────────────────────
sfx_dir = os.path.join(ROOT, "client", "chiikawa-pixel", "assets", "audio", "sfx")
bgm_dir = os.path.join(ROOT, "client", "chiikawa-pixel", "assets", "audio", "bgm")
sfx_count = len([f for f in os.listdir(sfx_dir) if f.endswith(".wav")])
bgm_count = len([f for f in os.listdir(bgm_dir) if f.endswith(".wav")])
check("SFX 音效 >= 12", sfx_count >= 12, f"實際: {sfx_count}")
check("BGM 音樂 >= 4", bgm_count >= 4, f"實際: {bgm_count}")

# ── 輸出結果 ─────────────────────────────────────────────────
print("\n" + "="*60)
print("  QA 驗證報告 — DAY-299")
print("="*60)
passed = sum(1 for s, _, _ in results if s == PASS)
failed = sum(1 for s, _, _ in results if s == FAIL)
warned = sum(1 for s, _, _ in results if s == WARN)

for status, name, detail in results:
    line = f"  {status} {name}"
    if detail:
        line += f"\n       → {detail}"
    print(line)

print("="*60)
print(f"  通過: {passed}/{len(results)}  失敗: {failed}  警告: {warned}")
print("="*60)

if failed > 0:
    print("\n❌ QA 未通過，請修復失敗項目")
    sys.exit(1)
else:
    print("\n✅ QA 全部通過！")
    sys.exit(0)
