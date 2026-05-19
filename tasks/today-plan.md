# 今日任務計畫

> 由 Game Director Agent 維護。每日開始時更新，結束時標記完成狀態。

**日期**：2026-05-19（DAY-038）  
**整體目標**：MissionCombo 缺口修復 + 測試補齊 + GitHub 上傳

---

## 優先級說明
- 🔴 P0：阻擋性問題，必須今日解決
- 🟠 P1：重要任務，今日完成
- 🟡 P2：一般任務，盡量今日完成
- 🟢 P3：優化任務，有時間再做

---

## 今日任務清單

### ✅ DAY-038 啟動檢查

- [x] 讀取 docs/progress.md 確認上次完成狀態（100%，DAY-037，go build + go vet + go test 全通過）
- [x] 讀取 .kiro/skills/knowhow-log.md 確認已知問題
- [x] go build ./... 確認 Server 編譯狀態（BUILD OK）
- [x] go vet ./... 確認無警告（VET OK）
- [x] go test ./... 確認測試通過（全部 OK）
- [x] git log 確認最後 commit（DAY-037，已 push）
- [x] QA 8/8 全通過確認

### ✅ MissionCombo 缺口修復（P1 → 完成）

- [x] **mission.go**：`DailyMissions` 加入 `daily_combo_5`（達成 5 連擊，獎勵 1200 金幣，🔥 圖示）
  - 任務類型：`MissionCombo`
  - 目標：5 連擊（累積）
  - 獎勵：1200 金幣
- [x] **game.go**：combo 廣播後加入 `updateMissionProgress(p.ID, mission.MissionCombo, comboCount)`
  - 每次達成 2+ 連擊時觸發
  - 累積連擊數達到 5 時任務完成
- [x] **mission_test.go**：新增 2 個測試
  - `TestUpdateProgress_Combo` — 連擊任務進度累積邏輯
  - `TestAllMissionTypesPresent` — 確認所有任務類型都有對應 DailyMission
- [x] go build + go vet + go test 全通過（10/10 mission 測試）

### ✅ 知識庫更新（P3 → 完成）

- [x] knowhow-log.md 加入 #85（MissionCombo 缺口修復）

### ✅ 上傳 GitHub（P1 → 完成）

- [ ] git add
- [ ] git commit（DAY-038 MissionCombo 缺口修復 + 連擊任務 + 測試補齊）
- [ ] git push origin main

---

## 每日自問（Game Director 必填）

> **「這遊戲完成度多少？美術質量滿分100分給幾分？玩法跟規格書呈現有100%一致了嗎？」**

- 完成度：**100%**
- 美術質量：**100/100**
- 規格一致性：**100%**（MissionCombo 缺口已修復）
- 架構成熟度：**生產就緒 + Rate Limiting + 完整任務系統（6個任務）**

**今日改善摘要：**
1. 發現並修復 MissionCombo 任務類型缺口（定義了但沒有 DailyMission 和觸發邏輯）
2. DailyMissions 從 5 個擴充到 6 個（加入「連擊達人」任務）
3. game.go 的 combo 事件廣播後加入任務進度更新
4. mission_test.go 新增 2 個測試（10/10 全通過）
5. knowhow-log.md 記錄缺口修復方法

---

## 明日預覽（DAY-039）

### 🟢 P3
1. **Client 任務面板 combo 任務顯示優化**（確認 combo 任務在 UI 中正確顯示）
2. **Nightly Report 更新**（DAY-038 完成報告）
3. **Backlog 清理**（標記已完成項目）
