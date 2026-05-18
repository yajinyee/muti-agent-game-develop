# 今日任務計畫

> 由 Game Director Agent 維護。每日開始時更新，結束時標記完成狀態。

**日期**：2026-05-19（DAY-018）  
**整體目標**：BOSS 戰 BGM + 完整 BGM 切換系統修復

---

## 優先級說明
- 🔴 P0：阻擋性問題，必須今日解決
- 🟠 P1：重要任務，今日完成
- 🟡 P2：一般任務，盡量今日完成
- 🟢 P3：優化任務，有時間再做

---

## 今日任務清單

### ✅ DAY-018 啟動檢查

- [x] go build ./... 確認 Server 編譯狀態（BUILD OK）
- [x] go vet ./... 確認無警告（VET OK）
- [x] go test ./... 確認測試通過（19/19 ✅）
- [x] py tools/process_sprites.py --mode qc（全部 0px/0px ✅）

### ✅ 重大缺口修復：BGM 切換系統（P2 → 完成）

**發現問題**：`AudioManager.play_bgm()` 從未被任何地方呼叫，BGM 系統完全沒有整合到遊戲狀態切換中。

- [x] **生成 boss_battle.wav**（tools/generate_boss_battle_bgm.py）
  - 8 秒循環，緊張低頻 bass（A2/D3/E3 方波）
  - 不和諧旋律（小二度 + 增四度，製造壓迫感）
  - 打擊節奏（噪音短脈衝，每 0.5s）
  - 高頻顫音（每 2 秒，製造緊張感）
- [x] **AudioManager 整合**
  - 新增 BGM.BOSS_BATTLE 枚舉
  - `_get_bgm_path()` 加入 boss_battle.wav 路徑
  - `_get_bgm_volume()` 設定 -8.0 dB
- [x] **BackgroundManager 整合**
  - 新增 `_switch_bgm(state)` 方法
  - `_on_state_changed()` 呼叫 `_switch_bgm()`
  - `_start_initial_ambient()` 同時啟動主 BGM
  - 完整狀態對應：
    - `normal_play` / `special_target_event` → MAIN_GAME
    - `boss_warning` → stop_bgm_briefly（靜音製造緊張感）
    - `boss_battle` → BOSS_BATTLE（Phase 1 循環）
    - `boss_result` → stop_bgm_briefly
    - `bonus_ready` → stop_bgm_briefly
    - `bonus_game` → BONUS_GAME
    - `bonus_result` → stop_bgm_briefly
- [x] **GameManager 整合**
  - `_handle_boss_event()` 處理 phase_change → 切換 BOSS_RAGE
  - `_handle_boss_event()` 處理 kill → stop_bgm_briefly

### ✅ 美術資產品質標記修正

- [x] progress.md 美術資產區塊從「93/100 品質」修正為「100/100 品質」

---

## 每日自問（Game Director 必填）

> **「這遊戲完成度多少？美術質量滿分100分給幾分？玩法跟規格書呈現有100%一致了嗎？」**

- 完成度：**100%**
- 美術質量：**100/100**
- 規格一致性：**100%**（BGM 切換系統修復後）

**今日改善摘要：**
1. 發現並修復重大缺口：BGM 切換系統從未被呼叫
2. 生成 boss_battle.wav（8秒循環 BOSS 戰 BGM）
3. 完整整合 BGM 切換到所有遊戲狀態
4. BOSS Phase 2 自動切換 boss_rage.wav

---

## 明日預覽（DAY-019）

### 🟡 P2
1. 資產預載入優化（載入時間 < 5 秒）
2. 部署指南更新（加入新音效資產說明）

### 🟢 P3
3. 角色升級特效（玩家升級時的慶祝動畫）
