# 今日任務計畫

> 由 Game Director Agent 維護。每日開始時更新，結束時標記完成狀態。

**日期**：2026-05-31（DAY-339）
**整體目標**：多人投射物顯示 + 目標物移動模式改善 ✅

---

## 今日任務清單

### ✅ DAY-339 啟動檢查

- [x] go build ./... 確認 Server 編譯狀態（BUILD OK）
- [x] go vet ./... 確認無警告（VET OK）
- [x] 讀取 docs/progress.md 確認上次進度（DAY-338）
- [x] 讀取 knowhow-log.md 確認已知問題

### ✅ DAY-339 多人投射物顯示

- [x] Server：hub.go 新增 BroadcastExcept 方法
- [x] Server：messages.go 新增 MsgOtherPlayerAttack + OtherPlayerAttackPayload
- [x] Server：game.go 在 handleAttackLocked 廣播 other_player_attack
- [x] Client：GameManager.gd 新增 other_player_attack 訊號 + 訊息處理
- [x] Client：Cannon.gd 新增 OTHER_PLAYER_COLORS + _on_other_player_attack 函數
- [x] Client：Cannon.gd 連接 other_player_attack 訊號

### ✅ DAY-339 目標物移動模式改善

- [x] Server：data/tables.go 新增 BehaviorWave/BehaviorZigzag/BehaviorSpiral
- [x] Server：T002 → wave，T003 → zigzag，T004 → wave，T105 → wave
- [x] Client：TargetManager.gd 新增 wave/zigzag/spiral 移動邏輯
- [x] Client：_create_target_node 初始化波浪/Z字形/螺旋移動 meta

### ✅ DAY-339 知識庫更新

- [x] knowhow-log 條目 190（多人投射物顯示的正確架構）
- [x] knowhow-log 條目 191（Godot 4 目標物波浪/Z字形移動實作）
- [x] knowhow-log 條目 192（捕魚機多人感設計原則）
- [x] docs/progress.md 更新
- [x] tasks/today-plan.md 更新

### ✅ DAY-339 QA 驗證

- [x] qa_check_day339.py（44 項驗證，44/44 全部通過）
- [x] go build + vet 最終確認（零錯誤零警告）

### ✅ DAY-339 GitHub 同步

- [x] git add + commit + push

---

## 每日自問（Game Director 必填）

**最後一次玩這個遊戲是什麼時候？**
→ 尚未在 Godot 實際遊玩（需要在 Godot 確認波浪移動效果）

**玩的時候最讓我不爽的是什麼？**
→ 無法確認，因為沒有實際遊玩

**我修了嗎？**
→ 已修 AUTO 死鎖（DAY-338）
→ 已修 打擊感（DAY-338）
→ 已加 多人投射物顯示（DAY-339）
→ 已加 目標物波浪移動（DAY-339）
→ 下一步：在 Godot 實際遊玩一局，確認效果

---

## 里程碑記錄

- **多人投射物顯示**（DAY-339）— 玩家可以看到其他玩家的投射物
- **目標物移動多樣化**（DAY-339）— wave/zigzag/spiral 三種新移動模式
- **QA 驗證：** 44/44 全部通過（DAY-339）
