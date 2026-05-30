#!/usr/bin/env python3
"""
QA Check DAY-339
- 多人投射物顯示（other_player_attack）
- 目標物移動模式改善（wave/zigzag/spiral）
- Server build/vet 驗證
"""
import os
import subprocess
import sys

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
SERVER_DIR = os.path.join(ROOT, "server")
CLIENT_DIR = os.path.join(ROOT, "client", "chiikawa-pixel")

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

def read_file(path):
    try:
        with open(path, "r", encoding="utf-8") as f:
            return f.read()
    except:
        return ""

print("=" * 60)
print("DAY-339 QA Check")
print("=" * 60)

# ── 1. Server 編譯 ────────────────────────────────────────────
print("\n[1] Server 編譯驗證")
result = subprocess.run(["go", "build", "./..."], cwd=SERVER_DIR, capture_output=True, text=True)
check("go build ./...", result.returncode == 0, result.stderr[:200] if result.stderr else "")

result = subprocess.run(["go", "vet", "./..."], cwd=SERVER_DIR, capture_output=True, text=True)
check("go vet ./...", result.returncode == 0, result.stderr[:200] if result.stderr else "")

# ── 2. 多人投射物顯示 — Server 端 ────────────────────────────
print("\n[2] 多人投射物顯示 — Server 端")
messages_go = read_file(os.path.join(SERVER_DIR, "internal", "protocol", "messages.go"))
check("MsgOtherPlayerAttack 常數定義", "MsgOtherPlayerAttack" in messages_go)
check("other_player_attack 字串值", '"other_player_attack"' in messages_go)
check("OtherPlayerAttackPayload 結構定義", "OtherPlayerAttackPayload" in messages_go)
check("PlayerID 欄位", "PlayerID" in messages_go and "player_id" in messages_go)
check("CharacterID 欄位", "CharacterID" in messages_go and "character_id" in messages_go)
check("TargetX/TargetY 欄位", "TargetX" in messages_go and "TargetY" in messages_go)
check("IsHit 欄位", "IsHit" in messages_go and "is_hit" in messages_go)

hub_go = read_file(os.path.join(SERVER_DIR, "internal", "ws", "hub.go"))
check("BroadcastExcept 方法定義", "BroadcastExcept" in hub_go)
check("BroadcastExcept 排除邏輯", "excludeID" in hub_go)

game_go = read_file(os.path.join(SERVER_DIR, "internal", "game", "game.go"))
check("game.go 呼叫 BroadcastExcept", "BroadcastExcept" in game_go)
check("game.go 廣播 MsgOtherPlayerAttack", "MsgOtherPlayerAttack" in game_go)

# ── 3. 多人投射物顯示 — Client 端 ────────────────────────────
print("\n[3] 多人投射物顯示 — Client 端")
gm_gd = read_file(os.path.join(CLIENT_DIR, "scripts", "game", "GameManager.gd"))
check("GameManager.gd other_player_attack 訊號定義", "signal other_player_attack" in gm_gd)
check("GameManager.gd other_player_attack 訊息處理", '"other_player_attack"' in gm_gd)
check("GameManager.gd emit_signal other_player_attack", 'emit_signal("other_player_attack"' in gm_gd)

cannon_gd = read_file(os.path.join(CLIENT_DIR, "scripts", "game", "Cannon.gd"))
check("Cannon.gd OTHER_PLAYER_COLORS 定義", "OTHER_PLAYER_COLORS" in cannon_gd)
check("Cannon.gd _on_other_player_attack 函數", "_on_other_player_attack" in cannon_gd)
check("Cannon.gd 連接 other_player_attack 訊號", "other_player_attack.connect" in cannon_gd)
check("Cannon.gd 其他玩家投射物 z_index=15", "z_index = 15" in cannon_gd)
check("Cannon.gd 烏薩奇旋轉效果", "usagi" in cannon_gd and "720.0" in cannon_gd)

# ── 4. 目標物移動模式改善 — Server 端 ────────────────────────
print("\n[4] 目標物移動模式改善 — Server 端")
tables_go = read_file(os.path.join(SERVER_DIR, "internal", "data", "tables.go"))
check("BehaviorWave 常數定義", "BehaviorWave" in tables_go)
check("BehaviorZigzag 常數定義", "BehaviorZigzag" in tables_go)
check("BehaviorSpiral 常數定義", "BehaviorSpiral" in tables_go)
check("T002 使用 wave 移動", '"T002"' in tables_go and "BehaviorWave" in tables_go)
check("T003 使用 zigzag 移動", '"T003"' in tables_go and "BehaviorZigzag" in tables_go)
check("T004 使用 wave 移動", '"T004"' in tables_go and "BehaviorWave" in tables_go)
check("T105 使用 wave 移動", '"T105"' in tables_go and "BehaviorWave" in tables_go)

# ── 5. 目標物移動模式改善 — Client 端 ────────────────────────
print("\n[5] 目標物移動模式改善 — Client 端")
tm_gd = read_file(os.path.join(CLIENT_DIR, "scripts", "game", "TargetManager.gd"))
check("TargetManager.gd wave 移動邏輯", '"wave"' in tm_gd)
check("TargetManager.gd zigzag 移動邏輯", '"zigzag"' in tm_gd)
check("TargetManager.gd spiral 移動邏輯", '"spiral"' in tm_gd)
check("TargetManager.gd wave_amp meta", "wave_amp" in tm_gd)
check("TargetManager.gd wave_freq meta", "wave_freq" in tm_gd)
check("TargetManager.gd wave_phase meta", "wave_phase" in tm_gd)
check("TargetManager.gd wave_elapsed meta", "wave_elapsed" in tm_gd)
check("TargetManager.gd zz_amp meta", "zz_amp" in tm_gd)
check("TargetManager.gd zz_period meta", "zz_period" in tm_gd)
check("TargetManager.gd sp_amp meta", "sp_amp" in tm_gd)
check("TargetManager.gd base_y meta", "base_y" in tm_gd)
check("TargetManager.gd sin 函數使用", "sin(" in tm_gd)
check("TargetManager.gd clamp 防止超出邊界", "clamp(" in tm_gd)

# ── 6. 協定一致性 ─────────────────────────────────────────────
print("\n[6] 協定一致性驗證")
check("Server MsgOtherPlayerAttack 與 Client 字串一致",
      '"other_player_attack"' in messages_go and '"other_player_attack"' in gm_gd)
check("Server BehaviorWave 與 Client wave 字串一致",
      '"wave"' in tables_go and '"wave"' in tm_gd)
check("Server BehaviorZigzag 與 Client zigzag 字串一致",
      '"zigzag"' in tables_go and '"zigzag"' in tm_gd)

# ── 結果 ─────────────────────────────────────────────────────
print("\n" + "=" * 60)
total = passed + failed
print(f"結果：{passed}/{total} 通過")
if failed == 0:
    print("✅ 全部通過！DAY-339 QA 驗證完成")
else:
    print(f"❌ {failed} 項失敗，請修復後重新執行")
print("=" * 60)

sys.exit(0 if failed == 0 else 1)
