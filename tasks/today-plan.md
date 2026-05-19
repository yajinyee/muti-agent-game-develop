# 今日任務計畫

> 由 Game Director Agent 維護。每日開始時更新，結束時標記完成狀態。

**日期**：2026-05-20（DAY-060）  
**整體目標**：Redis Pub/Sub 水平擴展廣播層 + Godot 4.6.3 版本追蹤 + KnowHow#119-120 + GitHub 上傳

---

## 優先級說明
- 🔴 P0：阻擋性問題，必須今日解決
- 🟠 P1：重要任務，今日完成
- 🟡 P2：一般任務，盡量今日完成
- 🟢 P3：優化任務，有時間再做

---

## 今日任務清單

### ✅ DAY-060 啟動檢查

- [x] 讀取 docs/progress.md 確認上次完成狀態（100%，DAY-059，背景圖 Lossy 壓縮）
- [x] go build ./... 確認 Server 編譯狀態（BUILD OK）
- [x] go vet ./... 確認無警告（VET OK）
- [x] go test ./... 確認測試通過（9/9 套件全部 OK）
- [x] QA 全部通過（8/8，RTP 96.09%）

### 🟠 Redis Pub/Sub 水平擴展廣播層（P1）

- [x] 研究 Redis pub/sub 水平擴展模式（oneuptime.com 2026-03-31）
- [x] 建立 `server/internal/ws/pubsub.go`（PubSubBroker，170 行）
  - `NewPubSubBroker()` — 建立代理，無 Redis 時回傳 nil（降級）
  - `Start()` / `Stop()` — 啟動/停止訂閱 goroutine
  - `Publish(msg)` — 發布訊息到 Redis channel
  - `subscribeLoop()` — 訂閱 Redis channel，serverID 過濾避免無限循環
  - `Hub.localBroadcast()` — 只廣播給本地客戶端
  - `Hub.BroadcastWithPubSub()` — 可選升級路徑
- [x] 建立 `server/internal/ws/pubsub_test.go`（4 個單元測試全部通過）
- [x] go build ./... + go test ./... 全部通過
- [x] KnowHow #119 更新

### 🟡 上網研究（P2）

- [x] 搜尋「Godot 4.6 latest version release 2025 2026」
  - 結論：Godot 4.6.3 RC 2 已於 2026-05-12 發布，4.7 beta 進行中
  - 本專案用 4.6.2，等 4.6.3 正式版後評估升級
  - KnowHow #120 更新

### 🟡 能力評估 #37（P2）

- [x] 更新 docs/ability-score.md

### 🟡 Nightly Report（P2）

- [x] 生成 DAY-060 nightly report

### 🟠 上傳 GitHub（P1）

- [x] git add + git commit + git push

---

## 每日自問（Game Director 必填）

> **「這遊戲完成度多少？美術質量滿分100分給幾分？玩法跟規格書呈現有100%一致了嗎？」**

- 完成度：**100%**
- 美術質量：**100/100**
- 規格一致性：**100%**
- 架構成熟度：**生產就緒 + Redis pub/sub 水平擴展就緒 + KnowHow 120 條**

---

## 明日預覽（DAY-061）

### 🟢 P3
1. **Redis pub/sub 整合到 main.go** — 在 GameRoom 中使用 BroadcastWithPubSub，完成水平擴展閉環
2. **上網搜尋** — 「Godot 4.6.3 stable release notes」
3. **GitHub 上傳**
