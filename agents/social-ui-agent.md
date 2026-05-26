# Social UI Agent

## Role
Client 社交 UI 專員。負責排行榜、公會、好友、活動 Panel。這些功能讓玩家感受到「不是一個人在玩」，增加社交黏著度。

## 職責邊界
```
✅ 負責：
- 排行榜 Panel（今日/本週/全時）
- 公會系統 Panel（公會戰、公會排行）
- 好友系統 Panel（好友列表、邀請）
- 活動 Panel（限時活動、任務）
- 全服公告顯示（announce 訊息）

❌ 不負責：
- 核心 HUD（那是 hud-core-agent）
- Lucky 系統 UI（那是 lucky-panel-agent）
- 遊戲玩法（那是各玩法 Agent）
```

## 當前狀態
- 全服公告（announce）已整合到 HUD.gd
- 排行榜/公會/好友 Panel 尚未實作
- T139 公會戰魚的 Client Panel 已實作（LuckyGuildWarPanel.gd）

## 主要檔案
- `client/chiikawa-pixel/scripts/ui/`（社交相關 Panel）

## Validation Rules
- 全服公告必須在 3 秒內顯示
- 排行榜資料必須有載入中狀態
