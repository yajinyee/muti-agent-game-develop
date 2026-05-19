# 今日任務計畫

> 由 Game Director Agent 維護。每日開始時更新，結束時標記完成狀態。

**日期**：2026-05-19（DAY-033）  
**整體目標**：高倍率目標光暈效果 + 能力評估更新 + GitHub 上傳

---

## 優先級說明
- 🔴 P0：阻擋性問題，必須今日解決
- 🟠 P1：重要任務，今日完成
- 🟡 P2：一般任務，盡量今日完成
- 🟢 P3：優化任務，有時間再做

---

## 今日任務清單

### ✅ DAY-033 啟動檢查

- [x] 讀取 docs/progress.md 確認上次完成狀態（100%，DAY-032 目標物倍率標籤 + BackgroundManager 修復，美術 98/100）
- [x] 讀取 .kiro/skills/knowhow-log.md 確認已知問題
- [x] go build ./... 確認 Server 編譯狀態（BUILD OK）
- [x] go vet ./... 確認無警告（VET OK）
- [x] go test ./... 確認測試通過（全部 OK）
- [x] git log 確認最後 commit（DAY-032，已 push）

### ✅ 高倍率目標光暈效果（P2 → 完成）

- [x] **TargetManager.gd**：加入 `_add_high_value_glow()` 函數
  - 30x+：金色光暈（Color 1.0, 0.85, 0.0, 0.25），脈動週期 0.6s
  - 50x+：橙紅光暈（Color 1.0, 0.5, 0.1, 0.35），脈動週期 0.4s + 縮放脈動
  - 光暈 ColorRect 放在 Sprite 後面（z_index = -1）
  - 脈動閃爍動畫（呼吸感）
  - 50x 額外縮放脈動（更強烈視覺衝擊）
- [x] `_create_target_node` 整合：倍率 >= 30 且非 BOSS 時觸發

### ✅ 能力評估更新（P3 → 完成）

- [x] 更新 `docs/ability-score.md`（評估 #24，DAY-032 + DAY-033）

### ✅ 上傳 GitHub（P1 → 完成）

- [x] git add
- [x] git commit（DAY-033 高倍率目標光暈效果 + 能力評估更新）
- [x] git push origin main

---

## 每日自問（Game Director 必填）

> **「這遊戲完成度多少？美術質量滿分100分給幾分？玩法跟規格書呈現有100%一致了嗎？」**

- 完成度：**100%**
- 美術質量：**99/100**（高倍率目標有光暈，視覺層次更豐富）
- 規格一致性：**100%**
- 架構成熟度：**完整，Docker + Redis 就緒**

**今日改善摘要：**
1. 高倍率目標（30x+）加入金色/橙紅光暈閃爍效果
2. 50x 目標有縮放脈動，視覺衝擊更強
3. 能力評估 #24 更新

---

## 明日預覽（DAY-034）

### 🟢 P3
1. **BOSS AI 圖生成**（B001 完整動畫集，ComfyUI — 需手動啟動）
2. **Server Docker Compose 部署測試**（驗證完整部署流程）
3. **Backlog 清理**（標記已完成項目）
