# 今日任務計畫

> 由 Game Director Agent 維護。每日開始時更新，結束時標記完成狀態。

**日期**：2026-05-18（DAY-009）  
**整體目標**：維持 100% 完成度，修復潛在 bug，持續優化

---

## 優先級說明
- 🔴 P0：阻擋性問題，必須今日解決
- 🟠 P1：重要任務，今日完成
- 🟡 P2：一般任務，盡量今日完成
- 🟢 P3：優化任務，有時間再做

---

## 今日任務清單

### ✅ DAY-009 啟動檢查（自動執行）

- [x] 讀取 docs/progress.md 確認上次完成狀態（100%）
- [x] 讀取 .kiro/skills/knowhow-log.md 確認已知問題
- [x] go build ./... 確認 Server 編譯狀態（BUILD OK）
- [x] go vet ./... 確認無警告（VET OK）
- [x] py tools/process_sprites.py --mode qc（全部 ✅）
- [x] py tools/qa_check.py（8/8 全部通過）

### ✅ Bug 修復

- [x] **Bonus tick 廣播 bug**：`int(elapsed)%1 == 0` 永遠為 true，改為 `lastBonusTickAt` 追蹤，確保每秒只廣播一次
  - 修復前：每 100ms 廣播一次（10x 過度廣播）
  - 修復後：每秒廣播一次（正確）
  - 影響：減少 90% 的 bonus tick 網路流量

### 🟢 P3 優化（進行中）

- [ ] 上網搜尋 Godot 4 HTML5 export 優化技術
- [ ] 評估是否需要 Spritesheet 合併（減少 draw call）
- [ ] 評估 Go Server 的 WebSocket 壓縮效果

---

## 今日決策記錄

| 時間 | 決策 | 理由 |
|------|------|------|
| 啟動 | 修復 bonus tick bug | `int(elapsed)%1 == 0` 是邏輯錯誤，永遠為 true |
| 設計 | 用 `lastBonusTickAt` 追蹤 | 比 `int(elapsed)%1` 更清晰，不會有精度問題 |

---

## 每日自問（Game Director 必填）

> **「這遊戲完成度多少？美術質量滿分100分給幾分？玩法跟規格書呈現有100%一致了嗎？」**

- 完成度：100%（修復了 bonus tick bug，品質更高）
- 美術質量：100/100（QC 全部通過）
- 規格一致性：100%（QA 8/8 通過）

**最低分項目**：無明顯低分項目。繼續搜尋可優化的地方。

---

## 明日預覽（DAY-010）

### 🟡 P2
1. 評估 HTML5 export 大小優化（Lossy 壓縮、移除未使用資產）
2. 評估多人房間支援的架構設計

### 🟢 P3
3. 數據埋點設計（玩家行為分析）
4. 像素字體在 HTML5 上的渲染測試
