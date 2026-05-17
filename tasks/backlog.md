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

### 🟠 P1
- [ ] **角色動畫幀數提升**：吉伊卡哇 idle 動畫從 4 幀提升到 8 幀
  - 負責：Animation Agent + Sprite Generation Agent
  - 預計工時：3 小時
  - 驗收：Animation Quality >= 90

- [ ] **目標物游泳動畫優化**：T001-T030 普通魚類動畫更自然
  - 負責：Animation Agent
  - 預計工時：4 小時
  - 驗收：Animation Quality >= 90

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
- [ ] **自動射擊模式優化**：自動模式的目標選擇邏輯更智慧
  - 負責：Godot Client Agent
  - 預計工時：2 小時
  - 驗收：Gameplay Feel >= 90

- [ ] **連線斷線提示優化**：更友善的斷線/重連 UI
  - 負責：Godot Client Agent
  - 預計工時：1 小時

### 🟡 P2
- [ ] **排行榜功能**：顯示當日最高分玩家
  - 負責：Go Server Agent + Godot Client Agent
  - 預計工時：4 小時
  - 依賴：協定變更（需 Spec Architect 審核）

- [ ] **成就系統**：首次擊敗 BOSS、連續命中等成就
  - 負責：Go Server Agent + Godot Client Agent
  - 預計工時：6 小時

### 🟢 P3
- [ ] **觀戰模式**：允許玩家觀看其他玩家的遊戲
  - 負責：Go Server Agent + Godot Client Agent
  - 預計工時：8 小時
  - 依賴：多人房間系統

---

## 技術優化

### 🟠 P1
- [ ] **HTML5 效能優化**：確保低階設備也能達到 30 FPS
  - 負責：Godot Client Agent
  - 預計工時：3 小時
  - 驗收：最低 FPS >= 30（低階設備）

- [ ] **Server 記憶體優化**：長時間運行後的記憶體洩漏檢查
  - 負責：Go Server Agent
  - 預計工時：2 小時
  - 驗收：24 小時運行後記憶體增長 < 10%

### 🟡 P2
- [ ] **WebSocket 壓縮**：啟用 permessage-deflate 壓縮
  - 負責：Go Server Agent + Godot Client Agent
  - 預計工時：2 小時
  - 依賴：協定變更審核

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
- [ ] **玩家操作手冊**：遊戲規則與操作說明文件
  - 負責：Spec Architect
  - 預計工時：2 小時

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

---

## Backlog 統計

| 優先級 | 數量 | 預計總工時 |
|-------|------|---------|
| 🔴 P0 | 0 | - |
| 🟠 P1 | 5 | ~12 小時 |
| 🟡 P2 | 9 | ~25 小時 |
| 🟢 P3 | 7 | ~25 小時 |
| ⚪ 待評估 | 0 | - |
| **合計** | **21** | **~62 小時** |

*最後更新：2025-01-01*
