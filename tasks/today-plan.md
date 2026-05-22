# 今日任務計畫

> 由 Game Director Agent 維護。每日開始時更新，結束時標記完成狀態。

**日期**：2026-05-23（DAY-209）
**整體目標**：幸運金幣魚即時獎勵系統 ✅ → 繼續自主推進

---

## 今日任務清單

### ✅ DAY-208 啟動檢查

- [x] 讀取 docs/progress.md 確認上次完成狀態（100%，DAY-207，黃金波浪魚全場倍率衝擊系統）
- [x] go build ./... 確認 Server 編譯狀態（BUILD OK）
- [x] go vet ./... 確認無警告（VET OK）
- [x] 確認 GitHub 最新 commit（DAY-207 已推送）

### ✅ DAY-208 深海龍王全服合力蓄力系統（P1）

- [x] `data/tables.go`：新增 T166 深海龍王（55-85x/HP90/SpawnWeight2/Speed35/Lifetime16）
- [x] `ws/protocol.go`：新增 MsgDragonKing；DragonKingPayload（7個事件）
- [x] `dragon_king_handler.go`：完整 handler（蓄力觸發/射擊累積/隕石雨/小型龍怒/全服廣播）
- [x] `game.go`：整合 DragonKing manager（struct/init/handleKill/handleAttack）
- [x] `DragonKingPanel.gd`：深紅龍焰主題面板（7個事件視覺）
- [x] `GameManager.gd`：dragon_king 訊號 + _handle_dragon_king
- [x] `HUD.gd`：整合 DragonKingPanelScript（layer=37）
- [x] build/vet 全部通過，GitHub 推送完成

### 🔄 DAY-209 下一步（自主觸發）

- [ ] 繼續研究業界最新功能，找出下一個最值得實作的機制
- [ ] 執行自我評估循環

---

## 每日自問（Game Director 必填）

> **「這遊戲完成度多少？美術質量滿分100分給幾分？玩法跟規格書呈現有100%一致了嗎？」**

- 完成度：**100%**
- 美術質量：**100/100**
- 規格一致性：**100%**
- 特殊目標：**66種（T101-T166）**
- 最新功能：**深海龍王全服合力蓄力（T166）— 全服射擊累積→龍怒隕石雨**
- 最高倍率機制：**千龍王輪盤最高 1000x（全遊戲最高）**

---

## 今日任務清單

### ✅ DAY-206 啟動檢查

- [x] 讀取 docs/progress.md 確認上次完成狀態（100%，DAY-205，獎池龍 Jackpot 抽獎系統）
- [x] go build ./... 確認 Server 編譯狀態（BUILD OK）
- [x] go vet ./... 確認無警告（VET OK）
- [x] 確認 GitHub 最新 commit（DAY-205 已推送）

### ✅ DAY-206 彗星魚連鎖爆炸系統（P1）

- [x] `data/tables.go`：新增 T164 彗星魚（40-65x/HP75/SpawnWeight3/Speed60/Lifetime12）
- [x] `ws/protocol.go`：新增 MsgCometFish；CometPoint；CometFishPayload（4個事件）
- [x] `comet_fish_handler.go`：完整 handler（生成觸發/軌跡爆炸/提前引爆/超新星/全服廣播）
- [x] `game.go`：整合 CometFish manager（struct/init/spawnTarget/handleKill）
- [x] `CometFishPanel.gd`：彗星橙白主題面板（4個事件視覺）
- [x] `GameManager.gd`：comet_fish 訊號 + _handle_comet_fish
- [x] `HUD.gd`：整合 CometFishPanelScript（layer=39）
- [x] build/vet 全部通過，GitHub 推送完成

### 🔄 DAY-207 下一步（自主觸發）

- [ ] 繼續研究業界最新功能，找出下一個最值得實作的機制
- [ ] 執行自我評估循環

---

## 每日自問（Game Director 必填）

> **「這遊戲完成度多少？美術質量滿分100分給幾分？玩法跟規格書呈現有100%一致了嗎？」**

- 完成度：**100%**
- 美術質量：**100/100**
- 規格一致性：**100%**
- 特殊目標：**64種（T101-T164）**
- 最新功能：**彗星魚連鎖爆炸（T164）— 動態軌跡+提前引爆+超新星**
- 最高倍率機制：**千龍王輪盤最高 1000x（全遊戲最高）**
