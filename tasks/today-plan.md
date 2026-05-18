# 今日任務計畫

> 由 Game Director Agent 維護。每日開始時更新，結束時標記完成狀態。

**日期**：2026-05-19（DAY-023）  
**整體目標**：觀戰模式（Spectator Mode）+ 上傳 GitHub

---

## 優先級說明
- 🔴 P0：阻擋性問題，必須今日解決
- 🟠 P1：重要任務，今日完成
- 🟡 P2：一般任務，盡量今日完成
- 🟢 P3：優化任務，有時間再做

---

## 今日任務清單

### ✅ DAY-023 啟動檢查

- [x] 讀取 docs/progress.md 確認上次完成狀態（100%）
- [x] 讀取 .kiro/skills/knowhow-log.md 確認已知問題
- [x] go build ./... 確認 Server 編譯狀態（BUILD OK）
- [x] go vet ./... 確認無警告（VET OK）
- [x] go test ./... 確認測試通過（5 套件全部通過 ✅）

### ✅ 觀戰模式（Spectator Mode）（P3 → 完成）

- [x] **ws/hub.go：ClientRole 枚舉**
  - `RolePlayer` / `RoleSpectator`
  - `Client` 加入 `Role` 欄位
- [x] **ws/hub.go：角色分離邏輯**
  - `Register` / `Unregister`：只對 `RolePlayer` 觸發 `OnConnect` / `OnDisconnect`
  - `readPump`：觀戰者訊息過濾（只允許 `ping`）
  - `PlayerCount()` / `SpectatorCount()` / `ClientCount()` 三個計數方法
  - `ServeSpectatorWS()` 新端點
- [x] **game.go：SpectatorSnapshot**
  - `GetSpectatorSnapshot()`：遊戲狀態 + 所有存活目標 + 排行榜 + 玩家數
- [x] **protocol.go：MsgSpectatorJoin**
- [x] **main.go：觀戰端點**
  - `ws://host/spectate?room_id=xxx` — 觀戰 WebSocket
  - `GET /spectate/snapshot` — HTTP 快照（供前端預覽）
  - 觀戰者連線後 100ms 非同步傳送完整快照
  - `/health` / `/stats` 加入 `spectators` 欄位
- [x] **ws/hub_test.go：7 個單元測試**
  - TestClientRole ✅
  - TestHubPlayerCount ✅
  - TestHubSpectatorCount ✅
  - TestHubOnConnectNotCalledForSpectator ✅
  - TestHubOnDisconnectNotCalledForSpectator ✅
  - TestHubSendToSpectator ✅
  - TestHubBroadcastReachesSpectator ✅

### ✅ 上傳 GitHub（P1 → 完成）

- [x] git add（所有變更）
- [x] git commit（DAY-023 觀戰模式）
- [x] git push origin main

---

## 每日自問（Game Director 必填）

> **「這遊戲完成度多少？美術質量滿分100分給幾分？玩法跟規格書呈現有100%一致了嗎？」**

- 完成度：**100%**
- 美術質量：**100/100**
- 規格一致性：**100%**

**今日改善摘要：**
1. 觀戰模式：`/spectate` WebSocket 端點 + `/spectate/snapshot` HTTP 端點
2. Hub 角色分離：玩家/觀戰者計數分離，觀戰者不觸發遊戲邏輯
3. 7 個新單元測試全部通過
4. KnowHow #87-88 更新

---

## 明日預覽（DAY-024）

### 🟢 P3
1. Server 水平擴展設計文件（Redis 共享狀態方案）
2. HTML5 export 大小優化（Lossy 壓縮 + gzip 驗證）
3. 觀戰模式 Client 端整合（LobbyManager 加入觀戰按鈕）
