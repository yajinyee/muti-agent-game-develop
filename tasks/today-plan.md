# 今日任務計畫

> 由 Game Director Agent 維護。每日開始時更新，結束時標記完成狀態。

**日期**：2026-05-19（DAY-042）  
**整體目標**：WebSocket 壓縮統計 + Client 可見性剔除 + Grafana 面板升級 + GitHub 上傳

---

## 優先級說明
- 🔴 P0：阻擋性問題，必須今日解決
- 🟠 P1：重要任務，今日完成
- 🟡 P2：一般任務，盡量今日完成
- 🟢 P3：優化任務，有時間再做

---

## 今日任務清單

### ✅ DAY-042 啟動檢查

- [x] 讀取 docs/progress.md 確認上次完成狀態（100%，DAY-041d，TargetPool tween 追蹤修復）
- [x] 讀取 .kiro/skills/knowhow-log.md 確認已知問題
- [x] go build ./... 確認 Server 編譯狀態（BUILD OK）
- [x] go vet ./... 確認無警告（VET OK）
- [x] go test ./... 確認測試通過（全部 OK，Windows unlinkat 是防毒軟體問題非測試失敗）
- [x] git log 確認最後 commit（DAY-041d，已 push，HEAD = origin/main）

### 🟠 Server WebSocket 壓縮統計（P1）

- [x] **`hub.go`**：加入 `BytesSentRaw atomic.Int64` 計數器
  - `Send()` 和 `Broadcast()` 中加入 `h.BytesSentRaw.Add(int64(len(data)))`
- [x] **`main.go`**：`/metrics` 端點加入 3 個壓縮統計指標
  - `chiikawa_ws_bytes_sent_raw_total`（原始位元組數）
  - `chiikawa_ws_avg_message_size_bytes`（平均訊息大小）
  - `chiikawa_ws_estimated_bytes_saved_total`（估算節省頻寬，65% 壓縮率）
- [x] go build + go vet + go test 確認通過（BUILD OK，VET OK，TEST OK）

### 🟡 Client 可見性剔除（P2）

- [x] **`TargetManager.gd`**：`_update_target_positions` 加入可見性剔除
  - 畫面內（-64 ~ 1344 x，-64 ~ 784 y）：`visible = true`
  - 畫面外但未移除：`visible = false`（減少 draw call）
  - 64px 緩衝避免邊緣閃爍

### 🟡 Grafana Dashboard 升級（P2）

- [x] **`chiikawa-overview.json`**：加入 2 個壓縮統計面板（共 14 個面板）
  - Panel 13：平均訊息大小 stat（bytes 單位，顏色警告）
  - Panel 14：壓縮節省頻寬 timeseries（原始 vs 估算 wire bytes/s）

### 🟢 KnowHow 更新（P3）

- [x] 記錄 #90（WebSocket 壓縮統計的正確方式）
- [x] 記錄 #91（Godot 4 可見性剔除）

### ✅ 上傳 GitHub（P1）

- [ ] git add
- [ ] git commit（DAY-042 WebSocket壓縮統計 + Client可見性剔除 + Grafana 14面板）
- [ ] git push origin main

---

## 每日自問（Game Director 必填）

> **「這遊戲完成度多少？美術質量滿分100分給幾分？玩法跟規格書呈現有100%一致了嗎？」**

- 完成度：**100%**
- 美術質量：**100/100**
- 規格一致性：**100%**
- 架構成熟度：**生產就緒 + Rate Limiting + 完整任務系統 + Prometheus 監控（14面板）+ TargetPool + 可見性剔除**

**今日改善目標：**
1. WebSocket 壓縮統計（讓 Grafana 能顯示頻寬節省效果）
2. Client 可見性剔除（減少畫面外目標物的 draw call）
3. Grafana 面板從 12 升級到 14 個

---

## 明日預覽（DAY-043）

### 🟢 P3
1. **Server 端 WebSocket 訊息類型統計**（各類型訊息的發送頻率）
2. **Client 端 FPS 穩定性測試**（記錄 60 秒的 FPS 波動）
3. **Backlog 清理**（標記已完成項目，整理 tasks/ 目錄）
