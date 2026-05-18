# 待辦清單（Backlog）

> 由 Game Director Agent 維護。所有未排入今日計畫的任務都在這裡。定期審查並移入 today-plan.md。

---

## 優先級說明
- 🔴 P0：阻擋性，需盡快排入計畫
- 🟠 P1：重要，本週內完成
- 🟡 P2：一般，本月內完成
- 🟢 P3：優化，有時間再做
- ⚪ 待評估：需要更多資訊才能決定優先級

---

## 美術優化

### 🔴 P0（阻擋 merge，明日必須完成）
- [x] **修復 chiikawa/hachiware attack 幀一致性**（2026-05-17 夜間完成）
  - 結果：chiikawa 0px/0px ✅，hachiware 0px/0px ✅
  - 工具：tools/fix_char_attack.py

- [ ] ~~**修復 hachiware hurt 動畫**~~ → **降級 P3**（遊戲不使用 hurt 狀態）
- [ ] ~~**修復 usagi bigwin 動畫**~~ → **已修復**（1px/1px ✅）
- [ ] ~~**修復 usagi hurt 動畫**~~ → **降級 P3**（遊戲不使用 hurt 狀態）
- [ ] ~~**補齊 hachiware skill 動畫**~~ → **降級 P3**（遊戲不使用 skill 狀態）
- [ ] ~~**補齊 hachiware fail 動畫**~~ → **降級 P3**（遊戲不使用 fail 狀態）
- [ ] ~~**補齊 usagi fail 動畫**~~ → **降級 P3**（遊戲不使用 fail 狀態）

### 🟠 P1
- [x] **角色動畫幀數提升**：吉伊卡哇 idle 動畫從 4 幀提升到 8 幀（已完成，spritesheet 768x288 = 8cols）
  - 驗證：chiikawa/hachiware/usagi 全部 8 幀 ✅

- [x] **目標物游泳動畫優化**：T001-T030 普通魚類動畫更自然（2026-05-18 DAY-010 完成）
  - 升級：Y軸搖擺 + 旋轉傾斜（±2-5度）+ 特殊目標縮放呼吸感 + 隨機相位

### 🟡 P2
- [ ] **BOSS 進場特效**：B001 BOSS 進場時的粒子特效
  - 負責：Godot Client Agent + Animation Agent
  - 預計工時：2 小時
  - 驗收：視覺衝擊感評分 >= 85

- [ ] **命中特效升級**：目標物被擊中時的閃白效果更明顯
  - 負責：Godot Client Agent
  - 預計工時：1 小時
  - 驗收：Gameplay Feel >= 90

- [ ] **UI 美化**：下注面板、分數顯示的像素風格統一
  - 負責：Art Director + Godot Client Agent
  - 預計工時：3 小時
  - 驗收：Visual Consistency >= 95

### 🟢 P3
- [ ] **背景動態效果**：海底背景的水波紋動畫
  - 負責：Animation Agent
  - 預計工時：2 小時
  - 驗收：視覺豐富度提升

- [ ] **角色升級特效**：玩家升級時的慶祝動畫
  - 負責：Animation Agent + Sprite Generation Agent
  - 預計工時：3 小時

---

## 音效優化

### 🟡 P2
- [ ] **Bonus 遊戲音效升級**：更有興奮感的 Bonus 觸發音效
  - 負責：Audio Director
  - 預計工時：1 小時
  - 驗收：Audio Sync >= 95

- [ ] **BOSS 戰 BGM**：專屬的 BOSS 戰背景音樂
  - 負責：Audio Director + Research Agent（尋找素材）
  - 預計工時：2 小時

### 🟢 P3
- [ ] **環境音效**：海底環境音（水泡聲、海浪聲）
  - 負責：Audio Director
  - 預計工時：1 小時

---

## 遊戲功能

### 🟠 P1
- [x] **HTML5 效能優化**：排除開發資源、Spritesheet AtlasTexture、資源預載入（2026-05-18 完成）
- [x] **Server 記憶體優化**：safeAfterFunc、graceful shutdown、/stats 端點（2026-05-18 完成）
- [x] **自動射擊模式優化**：評分系統（2026-05-17 完成）
- [x] **連線斷線提示優化**：DisconnectOverlay（2026-05-17 完成）
- [x] **HTML5 低階設備 30 FPS**：PerformanceMonitor 自動效能降級（2026-05-18 DAY-010 完成）
  - 三個等級：HIGH/MEDIUM/LOW，自動偵測並降級
  - LOW 模式：關閉游泳動畫/震動/outline shader，鎖 30 FPS

### 🟡 P2
- [x] **部署指南**：完整的 Server + Client 部署說明（2026-05-18 完成）
  - 輸出：docs/deployment-guide.md

- [x] **排行榜功能**：顯示當日最高分玩家（2026-05-18 DAY-010 完成）
  - Server：`/leaderboard` 端點、每 10 秒 WebSocket 廣播、Player 統計追蹤
  - Client：HUD 右上角面板，前 5 名，可折疊，自己高亮

- [ ] **成就系統**：首次擊敗 BOSS、連續命中等成就
  - 負責：Go Server Agent + Godot Client Agent
  - 預計工時：6 小時
  - ~~待辦~~ → **✅ 完成（2026-05-18 DAY-011）**
  - 實作：12 個成就、Tracker 模式、佇列式通知 UI

### 🟢 P3
- [ ] **觀戰模式**：允許玩家觀看其他玩家的遊戲
  - 負責：Go Server Agent + Godot Client Agent
  - 預計工時：8 小時
  - 依賴：多人房間系統

---

## 技術優化

### 🟠 P1
- [x] **Server 記憶體優化**：長時間運行後的記憶體洩漏檢查（2026-05-18 DAY-014 完成）
  - 壓力測試結果：Heap 增長 1.2 MB（< 50 MB 門檻）✅，Goroutine 增長 0 ✅，錯誤率 0% ✅
  - 修正 stress_test.py 判斷邏輯：改用穩定後趨勢斜率，而非從零開始的增長比例

### 🟡 P2
- [x] **WebSocket 壓縮**：啟用 permessage-deflate 壓縮（已完成，hub.go `EnableCompression: true`）
  - 確認：hub.go 第 24 行已有 `EnableCompression: true`，backlog 漏標記

- [ ] **資產預載入優化**：減少遊戲初始載入時間
  - 負責：Godot Client Agent
  - 預計工時：2 小時
  - 驗收：載入時間 < 5 秒

### 🟢 P3
- [ ] **Server 水平擴展**：支援多個 Server 實例
  - 負責：Go Server Agent
  - 預計工時：8 小時
  - 依賴：Redis 或其他共享狀態方案

---

## 文件與知識

### 🟡 P2
- [x] **玩家操作手冊**：遊戲規則與操作說明文件（2026-05-18 DAY-014 完成）
  - 輸出：docs/player-manual.md v1.1（效能設定 + 畫面說明 + FAQ 補充）

- [ ] **部署指南**：如何部署 Server 與 Client 的完整說明
  - 負責：Go Server Agent + Godot Client Agent
  - 預計工時：2 小時

### 🟢 P3
- [ ] **API 文件**：完整的 WebSocket API 文件（含範例）
  - 負責：Spec Architect
  - 預計工時：3 小時

---

## 已完成（歸檔）

- [x] AI 角色圖生成（吉伊卡哇、小八、烏薩奇）
- [x] 目標物 AI 生成（T001-T105）
- [x] RTP 校正
- [x] 多幀動畫實作
- [x] WebSocket 基礎通訊
- [x] Bonus 遊戲基礎流程
- [x] BOSS 戰基礎流程（B001）
- [x] HTML5 匯出設定
- [x] Multi-Agent Studio Scaffold（Phase 1）
- [x] Phase 2：規格文件（animation/audio/qa/daily-build spec）
- [x] Phase 3：Animation Pipeline（animation_pipeline.py）
- [x] Phase 4：Audio Pipeline（audio-map.json、sfx-list、bgm-layer-plan、sync-table）
- [x] Phase 5：Daily Build + QA 自動化（daily_build.ps1、qa_check.py）
- [x] Phase 6：Self-Improvement Loop（skills、failed-attempts）
- [x] Phase 7：Full Autonomous Studio 整合（nightly report、quality score）

---

## Backlog 統計

| 優先級 | 數量 | 預計總工時 |
|-------|------|---------|
| 🔴 P0 | 6 | ~8.5 小時 |
| 🟠 P1 | 7 | ~15 小時 |
| 🟡 P2 | 9 | ~25 小時 |
| 🟢 P3 | 7 | ~25 小時 |
| ⚪ 待評估 | 0 | - |
| **合計** | **29** | **~73.5 小時** |

*最後更新：2026-05-17*

---

## Phase 8 以後的計畫

### Phase 8：完整自主每日循環測試
- [ ] 測試完整的 Agent 協作循環（Game Director → 各 Agent → QA → Nightly Report）
  - 負責：Game Director + 全體 Agent
  - 預計工時：8 小時
  - 驗收：完整循環無人工介入執行一次

### 動畫幀數提升
- [ ] **chiikawa idle 幀數提升**（4 幀 → 8 幀）
  - 負責：Animation Agent + Sprite Generation Agent
  - 預計工時：3 小時
  - 驗收：Animation Quality >= 90

### BOSS AI 圖生成
- [ ] **BOSS AI 圖生成完成**（B001 完整動畫集）
  - 負責：Sprite Generation Agent + Animation Agent
  - 預計工時：4 小時
  - 驗收：BOSS 所有動畫狀態完整
