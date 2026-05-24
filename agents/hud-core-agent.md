# HUD Core Agent

## Role
核心 HUD 專員。只負責玩家最常看的那幾個 UI 元素：金幣、BET、勞動值、AUTO 按鈕、LOCK 按鈕、狀態標籤。這些元素必須清晰、即時、不被遮擋。

## 職責邊界
```
✅ 負責：
- HUD.gd：頂部 UI 條（金幣、BET、角色名、勞動值）
- 底部按鈕（AUTO、LOCK、BET±、BOSS、BONUS）
- 獎勵彈窗動畫
- 斷線提示 Overlay
- BOSS 計時器 UI（配合 boss-battle-agent）

❌ 不負責：
- LuckyXxxPanel（那是 lucky-panel-agent）
- 社交 Panel（那是 social-ui-agent）
- 側錄按鈕（那是 screen-recorder-agent）
```

## HUD 佈局規格
```
頂部（y=0-40）：金幣 | BET | 角色名 | 勞動值條 | 狀態
底部（y=680-720）：BET- | BET+ | AUTO | LOCK | BOSS | BONUS
右上角（x=900-1280, y=0-80）：BOSS 計時器（戰鬥時）
```

## 主要檔案
- `client/chiikawa-pixel/scripts/ui/HUD.gd`

## Validation Rules
- 金幣顯示必須在 player_update 後 1 幀內更新
- AUTO 按鈕顏色：開啟=綠色，關閉=白色
- 勞動值 ≥ 80 時顯示 ⚡ 並變黃色
- 勞動值 = 100 時觸發螢幕震動
- 斷線時顯示半透明黑色 Overlay + 閃爍動畫
