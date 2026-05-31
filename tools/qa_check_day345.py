"""
DAY-345 QA 驗證腳本
驗證項目：
1. T131-T160 美術升級（30個目標物）
2. T141-T145, T156-T160 損壞 PNG 修復（10個）
3. Server 每日任務系統（daily_quest.go）
4. Protocol 每日任務訊息（messages.go）
5. GameManager.gd 每日任務訊號
6. DailyQuestPanel.gd 存在
7. HUD.gd 每日任務按鈕
8. Server 編譯通過
"""
import os
import subprocess
import sys
from PIL import Image

BASE = r'd:\Kiro'
TARGETS = r'd:\Kiro\client\chiikawa-pixel\assets\sprites\targets'
SERVER = r'd:\Kiro\server'

passed = 0
failed = 0

def check(name, condition, detail=""):
    global passed, failed
    if condition:
        print(f"  ✅ {name}")
        passed += 1
    else:
        print(f"  ❌ {name}" + (f": {detail}" if detail else ""))
        failed += 1

print("=" * 60)
print("DAY-345 QA 驗證")
print("=" * 60)

# ── 1. T131-T160 美術升級 ──────────────────────────────────────
print("\n[1] T131-T160 美術升級（30個目標物）")
for i in range(131, 161):
    found = False
    for f in os.listdir(TARGETS):
        if f.startswith(f'T{i:03d}') and f.endswith('.png') and not f.endswith('.import'):
            found = True
            # 確認可以正常讀取
            try:
                img = Image.open(os.path.join(TARGETS, f)).convert('RGBA')
                pixels = sum(1 for px in img.getdata() if px[3] > 50)
                check(f"T{i:03d} 可讀取（{pixels} 像素）", pixels > 100, f"像素數不足: {pixels}")
            except Exception as e:
                check(f"T{i:03d} 可讀取", False, str(e))
            break
    if not found:
        check(f"T{i:03d} 存在", False, "檔案不存在")

# ── 2. Server 每日任務系統 ──────────────────────────────────────
print("\n[2] Server 每日任務系統")
daily_quest_path = os.path.join(SERVER, 'internal', 'game', 'daily_quest.go')
check("daily_quest.go 存在", os.path.exists(daily_quest_path))

if os.path.exists(daily_quest_path):
    content = open(daily_quest_path, encoding='utf-8').read()
    check("DailyQuestSystem 結構定義", "type DailyQuestSystem struct" in content)
    check("QuestKillTargets 任務類型", "QuestKillTargets" in content)
    check("QuestComboReach 任務類型", "QuestComboReach" in content)
    check("QuestTriggerBonus 任務類型", "QuestTriggerBonus" in content)
    check("OnKillTarget 函數", "func (dqs *DailyQuestSystem) OnKillTarget" in content)
    check("OnComboReach 函數", "func (dqs *DailyQuestSystem) OnComboReach" in content)
    check("OnTriggerBonus 函數", "func (dqs *DailyQuestSystem) OnTriggerBonus" in content)
    check("ClaimReward 函數", "func (dqs *DailyQuestSystem) ClaimReward" in content)
    check("GetPlayerStatus 函數", "func (dqs *DailyQuestSystem) GetPlayerStatus" in content)
    check("每日重置機制", "checkReset" in content)

# ── 3. Protocol 每日任務訊息 ──────────────────────────────────────
print("\n[3] Protocol 每日任務訊息")
messages_path = os.path.join(SERVER, 'internal', 'protocol', 'messages.go')
if os.path.exists(messages_path):
    content = open(messages_path, encoding='utf-8').read()
    check("MsgDailyQuestUpdate 常數", "MsgDailyQuestUpdate" in content)
    check("MsgDailyQuestComplete 常數", "MsgDailyQuestComplete" in content)
    check("MsgDailyQuestClaim 常數", "MsgDailyQuestClaim" in content)
    check("MsgDailyQuestRequest 常數", "MsgDailyQuestRequest" in content)
    check("DailyQuestUpdatePayload 結構", "DailyQuestUpdatePayload" in content)
    check("DailyQuestCompletePayload 結構", "DailyQuestCompletePayload" in content)
    check("DailyQuestClaimRequest 結構", "DailyQuestClaimRequest" in content)
    check("QuestStatusPayload 結構", "QuestStatusPayload" in content)

# ── 4. game.go 整合 ──────────────────────────────────────────
print("\n[4] game.go 每日任務整合")
game_path = os.path.join(SERVER, 'internal', 'game', 'game.go')
if os.path.exists(game_path):
    content = open(game_path, encoding='utf-8').read()
    check("dailyQuest 欄位", "dailyQuest *DailyQuestSystem" in content)
    check("newDailyQuestSystem 初始化", "newDailyQuestSystem()" in content)
    check("MsgDailyQuestRequest 處理", "MsgDailyQuestRequest" in content)
    check("MsgDailyQuestClaim 處理", "MsgDailyQuestClaim" in content)
    check("handleDailyQuestRequest 函數", "func (g *Game) handleDailyQuestRequest" in content)
    check("handleDailyQuestClaim 函數", "func (g *Game) handleDailyQuestClaim" in content)
    check("notifyQuestProgress 函數", "func (g *Game) notifyQuestProgress" in content)
    check("擊破觸發任務", 'notifyQuestProgress(playerID, "kill"' in content)
    check("Bonus 觸發任務", 'notifyQuestProgress(playerID, "bonus"' in content)
    check("連擊觸發任務", 'notifyQuestProgress(playerID, "combo"' in content)

# ── 5. Client GameManager.gd ──────────────────────────────────────
print("\n[5] Client GameManager.gd 每日任務訊號")
gm_path = r'd:\Kiro\client\chiikawa-pixel\scripts\game\GameManager.gd'
if os.path.exists(gm_path):
    content = open(gm_path, encoding='utf-8').read()
    check("daily_quest_update 訊號", "signal daily_quest_update" in content)
    check("daily_quest_complete 訊號", "signal daily_quest_complete" in content)
    check("daily_quest_update 處理", '"daily_quest_update"' in content)
    check("daily_quest_complete 處理", '"daily_quest_complete"' in content)
    check("request_daily_quests 函數", "func request_daily_quests" in content)
    check("claim_daily_quest 函數", "func claim_daily_quest" in content)

# ── 6. DailyQuestPanel.gd ──────────────────────────────────────
print("\n[6] DailyQuestPanel.gd")
panel_path = r'd:\Kiro\client\chiikawa-pixel\scripts\ui\DailyQuestPanel.gd'
check("DailyQuestPanel.gd 存在", os.path.exists(panel_path))
if os.path.exists(panel_path):
    content = open(panel_path, encoding='utf-8').read()
    check("_build_ui 函數", "func _build_ui" in content)
    check("_on_quest_update 函數", "func _on_quest_update" in content)
    check("_on_quest_complete 函數", "func _on_quest_complete" in content)
    check("_claim_quest 函數", "func _claim_quest" in content)
    check("_show_complete_notification 函數", "func _show_complete_notification" in content)
    check("_toggle_panel 函數", "func _toggle_panel" in content)

# ── 7. HUD.gd 每日任務按鈕 ──────────────────────────────────────
print("\n[7] HUD.gd 每日任務按鈕")
hud_path = r'd:\Kiro\client\chiikawa-pixel\scripts\ui\HUD.gd'
if os.path.exists(hud_path):
    content = open(hud_path, encoding='utf-8').read()
    check("_quest_btn 變數", "_quest_btn: Button" in content)
    check("_daily_quest_panel 變數", "_daily_quest_panel: Node" in content)
    check("_create_quest_button 呼叫", "_create_quest_button()" in content)
    check("_create_quest_button 函數", "func _create_quest_button" in content)

# ── 8. Server 編譯 ──────────────────────────────────────────
print("\n[8] Server 編譯驗證")
result = subprocess.run(
    ["go", "build", "./..."],
    cwd=SERVER,
    capture_output=True,
    text=True
)
check("go build ./... 通過", result.returncode == 0, result.stderr[:200] if result.stderr else "")

result_vet = subprocess.run(
    ["go", "vet", "./..."],
    cwd=SERVER,
    capture_output=True,
    text=True
)
check("go vet ./... 通過", result_vet.returncode == 0, result_vet.stderr[:200] if result_vet.stderr else "")

# ── 結果 ──────────────────────────────────────────────────────
print("\n" + "=" * 60)
total = passed + failed
print(f"結果：{passed}/{total} 通過")
if failed == 0:
    print("🎉 全部通過！DAY-345 驗證完成")
else:
    print(f"⚠️  {failed} 項失敗，需要修復")
print("=" * 60)

sys.exit(0 if failed == 0 else 1)
