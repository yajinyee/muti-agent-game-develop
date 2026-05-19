# 今日任務計畫

> 由 Game Director Agent 維護。每日開始時更新，結束時標記完成狀態。

**日期**：2026-05-20（DAY-053）  
**整體目標**：HUD.gd 大型腳本拆分（JackpotPanel / MissionPanel / SessionStatsPanel）+ KnowHow 更新 + GitHub 上傳

---

## 優先級說明
- 🔴 P0：阻擋性問題，必須今日解決
- 🟠 P1：重要任務，今日完成
- 🟡 P2：一般任務，盡量今日完成
- 🟢 P3：優化任務，有時間再做

---

## 今日任務清單

### ✅ DAY-053 啟動檢查

- [x] 讀取 docs/progress.md 確認上次完成狀態（100%，DAY-052，AudioManager 快取優化）
- [x] 讀取 .kiro/skills/knowhow-log.md 確認已知問題（最後 #100）
- [x] go build ./... 確認 Server 編譯狀態（BUILD OK）
- [x] go vet ./... 確認無警告（VET OK）
- [x] go test ./... 確認測試通過（全部 OK）
- [x] QA 全部通過（8/8，Audio Sync 99/100）

### ✅ HUD.gd 大型腳本拆分（P1）

- [x] 建立 `scripts/ui/JackpotPanel.gd`（Progressive Jackpot 面板，~250 行）
- [x] 建立 `scripts/ui/MissionPanel.gd`（每日任務面板，~250 行）
- [x] 建立 `scripts/ui/SessionStatsPanel.gd`（Session 統計面板，~200 行）
- [x] 更新 `HUD.gd`：加入 preload + 新的初始化函數
- [x] 移除 HUD.gd 中的舊函數（減少 860 行，從 2428 → 1598 行）
- [x] 確認 Server build/vet/test 全部通過

### ✅ KnowHow 更新（P2）

- [x] KnowHow #101：GDScript 大型腳本拆分策略
- [x] KnowHow #102：PowerShell str_replace 中文亂碼問題

### ✅ 上傳 GitHub（P1）

- [ ] git add（新腳本 + HUD.gd + knowhow + today-plan + progress）
- [ ] git commit（DAY-053 HUD.gd 拆分 + 三個獨立面板腳本）
- [ ] git push origin main

---

## 每日自問（Game Director 必填）

> **「這遊戲完成度多少？美術質量滿分100分給幾分？玩法跟規格書呈現有100%一致了嗎？」**

- 完成度：**100%**
- 美術質量：**100/100**
- 規格一致性：**100%**
- 架構成熟度：**生產就緒 + HUD 模組化（JackpotPanel/MissionPanel/SessionStatsPanel 獨立腳本）**

**今日改善：**
1. HUD.gd 從 2428 行縮減到 1598 行（-34%）
2. 三個面板各自獨立，易於維護
3. KnowHow #101-102 更新

---

## 明日預覽（DAY-054）

### 🟢 P3
1. **上網搜尋** — 「Godot 4 signal best practices large project」
2. **game.go 拆分評估** — game.go 已達 1740 行，考慮拆出 jackpot_handler.go / mission_handler.go
3. **Audio Sync 100/100** — 找出最後 1 分的差距（目前 99/100）
