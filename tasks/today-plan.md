# 今日任務計畫

> 由 Game Director Agent 維護。每日開始時更新，結束時標記完成狀態。

**日期**：2026-05-19（DAY-018）  
**整體目標**：資產預載入優化 + 角色升級特效 + 持續品質優化

---

## 優先級說明
- 🔴 P0：阻擋性問題，必須今日解決
- 🟠 P1：重要任務，今日完成
- 🟡 P2：一般任務，盡量今日完成
- 🟢 P3：優化任務，有時間再做

---

## 今日任務清單

### ✅ DAY-018 啟動檢查

- [x] 讀取 docs/progress.md 確認上次完成狀態（100%）
- [x] 讀取 .kiro/skills/knowhow-log.md 確認已知問題
- [x] go build ./... 確認 Server 編譯狀態（BUILD OK）
- [x] go vet ./... 確認無警告（VET OK）
- [x] go test ./... 確認測試通過（全部通過 ✅）
- [x] py tools/process_sprites.py --mode qc（全部 0px/0px ✅）
- [x] py tools/qa_check.py（8/8 全部通過 ✅，RTP 95.93%）

### ✅ 確認 BubbleLayer 氣泡音效（DAY-017 已完成）

- [x] 確認 BubbleLayer.gd 已有 bubble_pop 音效邏輯（第 113-116 行）
- [x] 確認 bubble_pop.wav 存在（13274 bytes）
- [x] 確認 bubble_pop.wav.import 存在（494 bytes）
- **結論**：DAY-017 已完成，無需重複實作

### ✅ 資產預載入系統（P2 → 完成）

- [x] **建立 LoadingManager.gd**（scripts/game/LoadingManager.gd）
  - 背景預載入：27 個 Textures + 14 個 Audio + 7 個 Shaders = 48 個資產
  - `ResourceLoader.load_threaded_request()` 非阻塞背景載入
  - `get_progress()` / `is_ready()` / `get_texture()` / `get_audio()` / `get_shader()` API
  - `loading_progress` / `loading_complete` 訊號
  - 載入完成後自動停止 `_process` 輪詢（節省 CPU）
- [x] **加入 project.godot Autoload**（LoadingManager）
- [x] **GameManager._ready() 啟動預載入**（call_deferred 避免初始化順序問題）

### ✅ 角色升級特效（P3 → 完成）

- [x] **HitEffect.spawn_level_up(pos, char_id)**
  - 全畫面角色色閃光（0.35 alpha）
  - 從砲台位置噴射 16 個金色星星粒子（向上扇形噴射 + 重力下落）
  - 大閃光環（角色顏色 + 金色，雙層）
  - 衝擊波
  - 升級文字動畫（"BONUS READY!" 彈入 → 停留 → 淡出，BACK 彈性）
  - 副標題「★ ★ ★」金色星星
- [x] **HUD.gd 整合**
  - `_last_labor_value` 追蹤上次勞動值
  - 偵測 labor >= 100 且上次 < 100 時觸發 `HitEffect.spawn_level_up()`
  - 同時觸發 `ScreenShake.add_trauma(0.3)`

### ✅ KnowHow 更新

- [x] KnowHow #85：Godot 4 ResourceLoader.load_threaded_request 背景載入
- [x] KnowHow #86：升級特效設計原則

---

## 每日自問（Game Director 必填）

> **「這遊戲完成度多少？美術質量滿分100分給幾分？玩法跟規格書呈現有100%一致了嗎？」**

- 完成度：**100%**
- 美術質量：**100/100**
- 規格一致性：**100%**

**今日改善摘要：**
1. 資產預載入系統：LoadingManager.gd（48 個資產背景預載入，避免首次使用卡頓）
2. 角色升級特效：spawn_level_up()（勞動值滿 100 時觸發，金色星星 + 彈入文字 + 閃光環）
3. HUD 整合：偵測勞動值達到 100 自動觸發升級特效

---

## 明日預覽（DAY-019）

### 🟢 P3
1. 上傳 GitHub（今日完成後執行）
2. 效能監控面板優化（FPS + 記憶體顯示）
3. 多人房間架構設計文件（未來功能規劃）
