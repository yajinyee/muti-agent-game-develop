# 今日任務計畫

> 由 Game Director Agent 維護。每日開始時更新，結束時標記完成狀態。

**日期**：2026-05-30（DAY-336）
**整體目標**：HUDLuckySignals.gd 擴充 DAY-304~319 + GitHub 同步 ✅

---

## 今日任務清單

### ✅ DAY-336 啟動檢查

- [x] go build ./... 確認 Server 編譯狀態（BUILD OK）
- [x] go vet ./... 確認無警告（VET OK）
- [x] 讀取 docs/progress.md 確認上次進度（DAY-335）
- [x] 讀取 knowhow-log.md 確認已知問題

### ✅ DAY-336 HUD.gd 重構（技術債）

- [x] HUDLuckySignals.gd 擴充 DAY-304 訊號連接（5個）+ 處理函數
- [x] HUDLuckySignals.gd 擴充 DAY-305 訊號連接（5個）+ 處理函數
- [x] HUDLuckySignals.gd 擴充 DAY-306 訊號連接（5個）+ 處理函數
- [x] HUDLuckySignals.gd 擴充 DAY-307 訊號連接（5個）+ 處理函數
- [x] HUDLuckySignals.gd 擴充 DAY-308 訊號連接（5個）+ 處理函數
- [x] HUDLuckySignals.gd 擴充 DAY-309 訊號連接（5個）+ 處理函數
- [x] HUDLuckySignals.gd 擴充 DAY-310 訊號連接（5個）+ 處理函數
- [x] HUDLuckySignals.gd 擴充 DAY-312 訊號連接（5個）+ 處理函數
- [x] HUDLuckySignals.gd 擴充 DAY-313 訊號連接（1個）+ 處理函數
- [x] HUDLuckySignals.gd 擴充 DAY-314 訊號連接（5個）+ 處理函數
- [x] HUDLuckySignals.gd 擴充 DAY-315 訊號連接（5個）+ 處理函數
- [x] HUDLuckySignals.gd 擴充 DAY-316 訊號連接（5個）+ 處理函數
- [x] HUDLuckySignals.gd 擴充 DAY-317 訊號連接（5個）+ 處理函數
- [x] HUDLuckySignals.gd 擴充 DAY-318~319 fallback 連接（10個）
- [x] HUD.gd 加入 _init_lucky_signals() 初始化函數
- [x] HUD.gd 移除 DAY-304~319 重複訊號連接

### ✅ DAY-336 知識庫更新

- [x] knowhow-log 條目 183（GDScript 模組化拆分的正確方式）
- [x] knowhow-log 條目 184（PowerShell 行數計算與 Python splitlines() 的差異）
- [x] knowhow-log 條目 185（HUDLuckySignals 委派架構的設計原則）
- [x] docs/progress.md 更新
- [x] tasks/today-plan.md 更新

### ✅ DAY-336 QA 驗證

- [x] qa_check_day336.py（121 項驗證，121/121 全部通過）
- [x] go build + vet 最終確認（零錯誤零警告）

### ✅ DAY-336 GitHub 同步

- [ ] git add + commit + push

---

## 每日自問（Game Director 必填）

**最後一次玩這個遊戲是什麼時候？**
→ 尚未在 Godot 實際遊玩（視覺清晰度 7.5/10 的根本原因）

**玩的時候最讓我不爽的是什麼？**
→ 無法確認，因為沒有實際遊玩

**我修了嗎？**
→ 已修 T001-T006 視覺（DAY-335）
→ 已修 HUD.gd 技術債（DAY-335~336）
→ 下一步：在 Godot 實際遊玩一局，確認視覺效果

---

## 里程碑記錄

- **HUDLuckySignals.gd 接管 DAY-292~319 全部 148 個訊號**（DAY-336）
- **HUD.gd 行數：** 2369 → 2293（減少 76 行）
- **QA 驗證：** 121/121 全部通過（DAY-336）
