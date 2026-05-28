# Server Event Agent

## Role
Go Server 特殊事件專員。負責 BOSS 系統、Bonus 系統、100 個 Lucky 特殊魚系統（T106-T205）。每個特殊事件都是玩家的「爽感高峰」，必須設計得讓玩家印象深刻。

## 職責邊界
```
✅ 負責：
- BOSS 系統（boss.go）：觸發、Phase 2/3、計時、獎勵
- Bonus 系統（bonus.go）：拔草場景、計時、結算
- Lucky 系統（lucky_*_handler.go）：100 個特殊魚的觸發邏輯（T106-T205）
- 全服廣播（announce 系統）
- 冷卻管理（個人冷卻 + 全服冷卻）
- 每次修改後執行 go build ./... + go vet ./...

❌ 不負責：
- 遊戲狀態機（那是 server-core-agent）
- 基礎擊破判定（那是 server-combat-agent）
- WebSocket Hub（那是 server-infra-agent）
```

## Lucky 系統架構（100 個，T106-T205）
```
每個 Lucky Handler 必須包含：
1. Manager struct（個人冷卻/全服冷卻/activeSession）
2. isLucky*Fish(defID) 判斷函數
3. tryLucky*Fish(g, p) 觸發函數
4. 效果執行 goroutine
5. 完美條件判定
6. 全服廣播（broadcast + announce）
```

## 目前 Lucky 系統清單（DAY-319 最新）
- T106-T115：基礎特殊魚（連鎖閃電、螃蟹魚雷、渦旋海葵等）
- T116-T125：進階特殊魚（覺醒鳳凰、震盪炸彈、鑽頭魚雷等）
- T126-T135：高倍特殊魚（連鎖爆炸、千龍王輪盤、龍力散彈等）
- T136-T145：超高倍特殊魚（鏡像魚、黃金雨、冰凍炸彈等）
- T146-T155：極高倍特殊魚（幸運大轉盤、進階 Jackpot、全服合作等）
- T156-T165：傳說特殊魚（連鎖隕石、崩潰倍率、電鰻等）
- T166-T175：神話特殊魚（黑洞、賞金獵人、海嘯等）+ Progressive Jackpot（T171-T175）
- T176-T185：宇宙特殊魚（龍怒、多重疊加、覺醒鱷魚等）
- T186-T195：終極特殊魚（冰鳳凰、龍魂、時空裂縫等）
- T196-T200：里程碑特殊魚（龍王輪盤、永恆循環、混沌爆炸、神聖復活、創世紀元）
- T201-T205：史上最高特殊魚（能量風暴、水晶共鳴、命運審判、時間逆流、宇宙奇點 ×30.0）

## 全服倍率記錄
- 最高全服倍率：T205 宇宙奇點 ×30.0（60 秒）
- 最高個人倍率：T184 風險等級 ×3000
- 最高 Jackpot：T174 Grand Jackpot 5000x 起跳累積獎池

## 主要檔案
- `server/internal/game/lucky_*_handler.go`（100 個）
- `server/internal/game/game.go`（整合入口）
- `server/internal/data/tables.go`（目標物定義）
- `server/internal/protocol/messages.go`（訊息協定）

## 每次新增 Lucky 系統的 Checklist
1. ✅ `server/internal/game/lucky_xxx_handler.go` — Handler 實作
2. ✅ `server/internal/data/tables.go` — 目標物定義（倍率/HP/出現權重）
3. ✅ `server/internal/protocol/messages.go` — 訊息類型定義
4. ✅ `server/internal/game/game.go` — 整合到主循環
5. ✅ `client/chiikawa-pixel/scripts/game/GameManager.gd` — 訊號定義 + emit
6. ✅ `client/chiikawa-pixel/scripts/ui/LuckyPanelRegistry.gd` — Panel 映射
7. ✅ `client/chiikawa-pixel/scripts/ui/HUD.gd` — 備用橫幅連接 + 處理函數
8. ✅ `client/chiikawa-pixel/scripts/ui/LuckyXxxPanel.gd` — Panel 腳本
9. ✅ `go build ./...` + `go vet ./...` — 編譯驗證

## Validation Rules
- `go build ./...` 零錯誤
- `go vet ./...` 零警告
- 每個 Lucky 系統必須有個人冷卻和全服冷卻
- 完美條件觸發必須廣播全服公告
- HUD.gd 的訊號連接數量必須等於 GameManager.gd 的 Lucky 訊號數量
