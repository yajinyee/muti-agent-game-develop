# 今日任務計畫

> 由 Game Director Agent 維護。每日開始時更新，結束時標記完成狀態。

**日期**：2026-05-22（DAY-007）  
**整體目標**：Gameplay Feel 提升（目標 92+）

---

## 優先級說明
- 🔴 P0：阻擋性問題，必須今日解決
- 🟠 P1：重要任務，今日完成
- 🟡 P2：一般任務，盡量今日完成
- 🟢 P3：優化任務，有時間再做

---

## 今日任務清單

### ✅ FEEL-001：Gameplay Juice 系統

- [x] 研究 Gameplay Feel 最佳實踐（Trauma Screen Shake、Hit Stop）
- [x] 建立 `ScreenShake.gd`（Autoload，Trauma-based 震動）
- [x] 建立 `HitEffect.gd`（Autoload，命中/擊殺/大獎特效）
- [x] 更新 `Cannon.gd`（命中震動 + Hit Stop + 子彈拖尾）
- [x] 更新 `TargetManager.gd`（擊殺爆炸 + BOSS 震動）
- [x] 更新 `BonusGame.gd`（Bonus 觸發特效）
- [x] 建立 `skills/skill-gameplay-juice.md`
- [x] Commit + Merge + Push

---

## 今日決策記錄

| 時間 | 決策 | 理由 |
|------|------|------|
| 自動觸發 | 優先改善 Gameplay Feel（88 → 92+）| 最低分項目，玩法完整性優先 |
| 設計 | ScreenShake 用 Node 繼承而非 Camera2D | Autoload 不能直接是 Camera2D |
| 設計 | Hit Stop 0.04s | 太長會讓玩家感覺卡頓，0.04s 剛好 |

---

## 明日預覽（DAY-008）

### 🔴 P0
1. 美術質量提升（目標 95+）
2. 上網搜尋「pixel art chiikawa style sprite optimization」

### 🟠 P1
3. 更新 `docs/ability-score.md`
4. 執行完整 QA 確認 Gameplay Feel 分數提升

### 🟡 P2
5. 規格一致性補齊（97% → 100%）
6. 研究 Godot 4 HTML5 效能優化

---

## 每日自問（Game Director 必填）

> **「這遊戲完成度多少？美術質量滿分100分給幾分？玩法跟規格書呈現有100%一致了嗎？」**

- 完成度：99%（Gameplay Juice 系統完成）
- 美術質量：92/100（目標 95+，明日重點）
- 規格一致性：97%（目標 100%）

**最低分項目**：美術質量（92）→ 明日重點提升
