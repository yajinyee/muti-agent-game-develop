# 今日任務計畫

> 由 Game Director Agent 維護。每日開始時更新，結束時標記完成狀態。

**日期**：2026-05-31（DAY-341）
**整體目標**：Combo 里程碑音效 + Combo UI 升級 ✅

---

## 今日任務清單

### ✅ DAY-341 啟動檢查

- [x] go build ./... 確認 Server 編譯狀態（BUILD OK）
- [x] go vet ./... 確認無警告（VET OK）
- [x] 讀取 docs/progress.md 確認上次進度（DAY-340）
- [x] 讀取 knowhow-log.md 確認已知問題
- [x] 上網研究 2026 最新捕魚機趨勢

### ✅ DAY-341 Combo 里程碑音效

- [x] 生成 combo_5.wav（5連擊輕快音階）
- [x] 生成 combo_10.wav（10連擊和弦爆發）
- [x] 生成 combo_20.wav（20連擊電子感）
- [x] 生成 combo_30.wav（30連擊勝利號角）
- [x] 生成 4 個 .import 檔案
- [x] AudioManager.gd 新增 COMBO_5/10/20/30 枚舉
- [x] AudioManager.gd 新增 play_combo_milestone() 函數

### ✅ DAY-341 Combo UI 升級

- [x] HUD.gd _create_combo_label() 升級（帶背景面板）
- [x] HUD.gd _update_combo_display() 整合里程碑音效觸發
- [x] HUD.gd _spawn_combo_milestone_effect() 新函數（粒子 + 大文字）
- [x] 震動強度隨連擊數增加

### ✅ DAY-341 知識庫更新

- [x] knowhow-log 條目 197（Combo 里程碑音效設計原則）
- [x] knowhow-log 條目 198（Godot 4 CanvasLayer visible 控制）
- [x] knowhow-log 條目 199（Combo 里程碑特效位置計算）
- [x] docs/progress.md 更新
- [x] tasks/today-plan.md 更新

### ✅ DAY-341 QA 驗證

- [x] qa_check_day341.py（41 項驗證，41/41 全部通過）
- [x] go build + vet 最終確認（零錯誤零警告）

### ⏳ DAY-341 GitHub 同步

- [ ] git add + commit + push

---

## 每日自問（Game Director 必填）

**最後一次玩這個遊戲是什麼時候？**
→ 尚未在 Godot 實際遊玩（需要在 Godot 確認 Combo 音效和特效效果）

**玩的時候最讓我不爽的是什麼？**
→ 無法確認，因為沒有實際遊玩

**我修了嗎？**
→ 已加 Combo 里程碑音效（DAY-341）
→ 已升級 Combo UI（DAY-341）
→ 下一步：在 Godot 實際遊玩一局，確認效果

---

## 里程碑記錄

- **Combo 里程碑音效**（DAY-341）— 4個里程碑各有獨特音效
- **Combo UI 升級**（DAY-341）— 帶背景面板 + 粒子特效 + 大文字
- **QA 驗證：** 41/41 全部通過（DAY-341）
