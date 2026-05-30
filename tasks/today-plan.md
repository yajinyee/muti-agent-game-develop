# 今日任務計畫

> 由 Game Director Agent 維護。每日開始時更新，結束時標記完成狀態。

**日期**：2026-05-30（DAY-335）
**整體目標**：深度優先策略——T001-T006 視覺升級 + HUD.gd 重構 + GitHub 同步 ✅

---

## 今日任務清單

### ✅ DAY-335 啟動檢查

- [x] go build ./... 確認 Server 編譯狀態（BUILD OK）
- [x] go vet ./... 確認無警告（VET OK）
- [x] 讀取 docs/progress.md 確認上次進度（DAY-334）
- [x] 讀取 knowhow-log.md 確認已知問題

### ✅ DAY-335 深度優先策略（P0）

- [x] T001 草：多根草莖 + 露珠 + 光澤高光（64x64）
- [x] T002 綠蟲：觸角 + 6隻腳 + 眼睛 + 高光（64x64）
- [x] T003 紅蟲：速度線尾跡 + 流線型身體 + 兇眼（64x64）
- [x] T004 藍蟲：電光尾跡 + 發光藍眼 + 電光觸角（64x64）
- [x] T005 布丁：奶油頂 + 草莓 + 光澤 + 微笑（64x64）
- [x] T006 蘑菇：白色斑點 + 莖陰影 + 可愛眼睛（64x64）
- [x] 工具：tools/upgrade_basic_targets_day335.py

### ✅ DAY-335 HUD.gd 重構（技術債）

- [x] 建立 HUDLuckySignals.gd（拆分 148 個 Lucky 訊號連接）
- [x] 委派模式：_show_banner / _show_event / _show_reward
- [x] DAY-292~303 所有 Lucky 訊號處理函數移入

### ✅ DAY-335 Agent 更新

- [x] target-design-agent.md 更新至 DAY-335 狀態
- [x] visual-clarity-agent.md 更新至 DAY-335 狀態

### ✅ DAY-335 知識庫更新

- [x] knowhow-log 條目 180（HUD.gd 技術債解法）
- [x] knowhow-log 條目 181（Windows 多 Python 版本衝突）
- [x] knowhow-log 條目 182（T001-T006 視覺升級策略）
- [x] docs/progress.md 更新
- [x] tasks/today-plan.md 更新

### ✅ DAY-335 QA 驗證

- [x] qa_check_day335.py（58 項驗證，58/58 全部通過）
- [x] go build + vet 最終確認（零錯誤零警告）

### ✅ DAY-335 GitHub 同步

- [x] git add + commit + push

---

## 每日自問（Game Director 必填）

**最後一次玩這個遊戲是什麼時候？**
→ 尚未在 Godot 實際遊玩（視覺清晰度 7.5/10 的根本原因）

**玩的時候最讓我不爽的是什麼？**
→ 無法確認，因為沒有實際遊玩

**我修了嗎？**
→ 已修 T001-T006 視覺（深度優先策略第一步）
→ 下一步：在 Godot 實際遊玩一局，確認視覺效果

---

## 里程碑記錄

- **T001-T006 視覺升級**：深度優先策略第一步（DAY-335）
- **HUDLuckySignals.gd**：HUD.gd 技術債重構開始（DAY-335）
- **QA 驗證**：58/58 全部通過（DAY-335）
