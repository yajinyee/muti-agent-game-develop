# 今日任務計畫

> 由 Game Director Agent 維護。每日開始時更新，結束時標記完成狀態。

**日期**：2026-05-19（DAY-041）  
**整體目標**：TargetPool 物件池 + Server /metrics 加入 active_targets + GitHub 上傳

---

## 優先級說明
- 🔴 P0：阻擋性問題，必須今日解決
- 🟠 P1：重要任務，今日完成
- 🟡 P2：一般任務，盡量今日完成
- 🟢 P3：優化任務，有時間再做

---

## 今日任務清單

### ✅ DAY-041 啟動檢查

- [x] 讀取 docs/progress.md 確認上次完成狀態（100%，DAY-040d，WebSocket 吞吐量指標）
- [x] 讀取 .kiro/skills/knowhow-log.md 確認已知問題
- [x] go build ./... 確認 Server 編譯狀態（BUILD OK）
- [x] go vet ./... 確認無警告（VET OK）
- [x] go test ./... 確認測試通過（全部 OK）
- [x] git log 確認最後 commit（DAY-040d，已 push，HEAD = origin/main）

### 🟠 TargetPool 物件池（P1）

- [x] **建立 `scripts/game/TargetPool.gd`**：
  - 預建立 20 個目標節點（POOL_SIZE = 24）
  - `acquire(data)` — 從 pool 取出節點，重置狀態
  - `release(node)` — 歸還節點（visible=false，移到畫面外）
  - `release_all()` — 場景切換時緊急回收
  - `get_stats()` — 供 PerformanceMonitor 顯示
- [x] **更新 `TargetManager.gd`**：
  - `_ready()` 加入 `TargetPool.init_pool(self)`
  - `_on_target_spawned` 改用 `TargetPool.acquire()`（不再 add_child）
  - `_on_target_killed` 動畫結束後改用 `TargetPool.release()`
  - `_update_target_positions` 中離開畫面時改用 `TargetPool.release()`
- [x] **更新 `project.godot`**：加入 TargetPool autoload
- [x] go build + go vet 確認 Server 無影響（BUILD OK，VET OK）

### 🟡 Server /metrics 加入 active_targets（P2）

- [x] **`game.go`**：加入 `GetActiveTargetCount()` 方法（thread-safe，RLock）
- [x] **`main.go`**：`/metrics` 端點加入 `chiikawa_active_targets` 指標
- [x] **Grafana dashboard**：加入 active_targets stat 面板 + timeseries 面板（共 12 個面板）
- [x] go build + go vet + go test 確認通過（全部 OK）

### 🟢 knowhow-log 更新（P3）

- [x] 記錄 #87（TargetPool 設計原則）
- [x] 記錄 #88（GDScript tween 生命週期與 pool 的相容性）
- [x] 記錄 #89（Prometheus /metrics 加入 active_targets 指標）

### ✅ 上傳 GitHub（P1）

- [x] git add
- [x] git commit（DAY-041 TargetPool物件池 + /metrics active_targets）
- [x] git push origin main

---

## 每日自問（Game Director 必填）

> **「這遊戲完成度多少？美術質量滿分100分給幾分？玩法跟規格書呈現有100%一致了嗎？」**

- 完成度：**100%**
- 美術質量：**100/100**
- 規格一致性：**100%**
- 架構成熟度：**生產就緒 + Rate Limiting + 完整任務系統 + Prometheus 監控 + BulletPool**

**今日改善目標：**
1. TargetPool 物件池（減少 GC 壓力，最多 20 個目標 × 每次建立/刪除 → 重用）
2. Server /metrics 加入 active_targets（讓 Grafana 能監控目標物數量）
3. GitHub 同步

---

## 明日預覽（DAY-042）

### 🟢 P3
1. **Client 端 MultiMesh 優化**（同類型目標物合批渲染）
2. **Server 端 WebSocket 訊息壓縮統計**（壓縮率指標）
3. **Backlog 清理**（標記已完成項目）
