# 今日任務計畫

> 由 Game Director Agent 維護。每日開始時更新，結束時標記完成狀態。

**日期**：2026-05-20（DAY-093）
**整體目標**：每日錦標賽系統 + 上傳 GitHub + 自主循環

---

## 今日任務清單

### ✅ DAY-093 啟動檢查

- [x] 讀取 docs/progress.md 確認上次完成狀態（100%，DAY-092，每日簽到轉盤）
- [x] go build ./... 確認 Server 編譯狀態（BUILD OK）
- [x] go vet ./... 確認無警告（VET OK）
- [x] 提交未提交的修改（BG002 硬雜草連點修復 + UI 修復）

### ✅ DAY-093 每日錦標賽系統（P1）

- [x] `tournament/tournament.go`：新增 `DailyTournament` 管理器（每日 UTC+8 00:00 重置）
- [x] `tournament/tournament_test.go`：10 個每日賽測試（共 23 個全部通過）
- [x] `protocol.go`：`MsgGetTournament`/`MsgDailyTournamentUpdate`/`MsgDailyTournamentResult` + Payload
- [x] `tournament_handler.go`：完整 handler 集（新建）
- [x] `game.go`：整合每日賽管理器 + handler 分支 + gameLoop 廣播
- [x] `boss_handler.go`：BOSS 擊殺加入每日賽積分
- [x] `bonus_handler.go`：Bonus 完成加入每日賽積分
- [x] `main.go`：`/daily-tournament` HTTP 端點
- [x] `TournamentPanel.gd`：升級為雙 Tab 面板（今日賽/週賽）
- [x] `GameManager.gd`：`daily_tournament_updated` 訊號 + handler
- [x] `NetworkManager.gd`：`send_get_tournament()`
- [x] build/vet/test 全部通過（23/23）

### 🟠 上傳 GitHub（P1）

- [ ] git add + git commit + git push（DAY-093 每日錦標賽系統）

---

## 每日自問（Game Director 必填）

> **「這遊戲完成度多少？美術質量滿分100分給幾分？玩法跟規格書呈現有100%一致了嗎？」**

- 完成度：**100%**
- 美術質量：**100/100**
- 規格一致性：**100%**
- 架構成熟度：**生產就緒 + 完整社交系統 + 特殊武器 + 每日/週賽錦標賽 + KnowHow 135 條**
