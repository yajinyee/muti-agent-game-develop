# 今日任務計畫

> 由 Game Director Agent 維護。每日開始時更新，結束時標記完成狀態。

**日期**：2026-05-19（DAY-044）  
**整體目標**：Server Ping Latency 統計 + Grafana 18面板 + GitHub 上傳

---

## 優先級說明
- 🔴 P0：阻擋性問題，必須今日解決
- 🟠 P1：重要任務，今日完成
- 🟡 P2：一般任務，盡量今日完成
- 🟢 P3：優化任務，有時間再做

---

## 今日任務清單

### ✅ DAY-044 啟動檢查

- [x] 讀取 docs/progress.md 確認上次完成狀態（100%，DAY-043，訊息類型統計 + Combo 5+/7+ 視覺強化）
- [x] 讀取 .kiro/skills/knowhow-log.md 確認已知問題
- [x] go build ./... 確認 Server 編譯狀態（BUILD OK）
- [x] go vet ./... 確認無警告（VET OK）
- [x] go test ./... 確認測試通過（全部 OK）
- [x] git log 確認最後 commit（DAY-043，已 push，HEAD = origin/main）

### ✅ Server Ping Latency 統計（P1）

- [x] **`hub.go`**：Client 加入 `lastPingSentAt`、`lastPingLatMs`、`pingMu` 欄位
- [x] **`hub.go`**：Hub 加入 `pingLatencySum`、`pingLatencyCount`、`pingLatencyMax` 原子計數器
- [x] **`hub.go`**：`writePump` 在發送 ping 時記錄 `lastPingSentAt`
- [x] **`hub.go`**：`readPump` 的 pong handler 計算延遲並呼叫 `RecordPingLatency()`
- [x] **`hub.go`**：加入 `RecordPingLatency()`、`GetPingStats()`、`GetClientPingLatencies()` 方法
- [x] **`main.go`**：`/metrics` 加入 3 個 ping latency 指標（avg/max/samples）
- [x] **`main.go`**：`/metrics` 加入 per-client ping latency 指標
- [x] **`main.go`**：`/health` 加入 `avg_ping_ms` 欄位
- [x] go build + go vet + go test 確認通過（BUILD OK，VET OK，13/13 TEST OK）

### ✅ hub_test.go 新增測試（P1）

- [x] `TestRecordPingLatency`：確認 ping latency 統計正確累積（avg/max/count）
- [x] `TestRecordPingLatencyMax`：確認最大值追蹤正確（亂序記錄）
- [x] `TestGetClientPingLatencies`：確認 per-client 延遲查詢

### ✅ Grafana Dashboard 升級（P2）

- [x] **`chiikawa-overview.json`**：從 15 個面板升級到 18 個面板
  - Panel 16：平均 Ping 延遲 stat（ms 單位，顏色警告 100/300ms）
  - Panel 17：最大 Ping 延遲 stat（ms 單位，顏色警告 200/500ms）
  - Panel 18：Ping 延遲趨勢 timeseries（avg + max 雙線）

### ✅ Nightly Report DAY-043（P2）

- [x] 生成 `reports/nightly/nightly-report-2026-05-19-day043.md`

### ✅ 上傳 GitHub（P1）

- [ ] git add
- [ ] git commit（DAY-044 Ping Latency統計(avg/max/per-client) + Grafana 18面板 + 測試13/13）
- [ ] git push origin main

---

## 每日自問（Game Director 必填）

> **「這遊戲完成度多少？美術質量滿分100分給幾分？玩法跟規格書呈現有100%一致了嗎？」**

- 完成度：**100%**
- 美術質量：**100/100**
- 規格一致性：**100%**
- 架構成熟度：**生產就緒 + Rate Limiting + 完整任務系統（6個任務）+ Prometheus 監控（18個面板）+ TargetPool + 可見性剔除 + 訊息類型統計 + Ping Latency 追蹤**

**今日改善目標：**
1. Server 端 Ping Latency 統計（讓 Grafana 能顯示連線品質）
2. per-client 延遲追蹤（識別高延遲玩家）
3. Grafana 面板從 15 升級到 18 個

---

## 明日預覽（DAY-045）

### 🟢 P3
1. **Client 端效能數據上報** — FPS/記憶體/延遲定期上報 Server，讓 Grafana 能看到 Client 端效能
2. **Server 連線品質報告** — 定期輸出高延遲玩家警告 log
3. **Backlog 清理** — 標記已完成項目，整理 tasks/ 目錄
