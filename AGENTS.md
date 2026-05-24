# Multi-Agent Game Studio v2.0 — 重構架構說明

## 專案概覽

**遊戲名稱**：吉伊卡哇：像素大討伐（捕魚機）  
**技術棧**：Go + WebSocket（Port 7777）/ Godot 4.6.2（HTML5 匯出）  
**架構版本**：v2.0（2026-05-24 重構，從技術棧分工改為體驗循環分工）

---

## 為什麼重構

v1.0 的 12 個 Agent 按技術棧分工（Go Agent、Godot Agent），導致：
1. Godot Client Agent 一個人負責 UI、玩法、WebSocket、效能、資產整合——5 個人的工作量
2. 沒有任何 Agent 的職責是「確認 Server↔Client 端對端真的通」
3. 沒有任何 Agent 負責「玩起來好不好玩」
4. 250 天的功能堆疊，沒有一次被玩家視角驗證過

v2.0 改為**按體驗循環分工**：每個功能必須走完「設計→實作→整合→玩家驗證→優化」才算完成。

---

## 新架構圖（20 個 Agent）

```
╔══════════════════════════════════════════════════════════════╗
║                    決策層（2個）                              ║
║  ┌─────────────────┐    ┌──────────────────────────────┐    ║
║  │  Game Director  │    │  Player Experience Director  │    ║
║  │  （技術決策）    │    │  （玩家體驗決策）             │    ║
║  └────────┬────────┘    └──────────────┬───────────────┘    ║
╚═══════════╪═════════════════════════════╪════════════════════╝
            │                             │
╔═══════════╪═════════════════════════════╪════════════════════╗
║           ▼         設計層（3個）        ▼                    ║
║  ┌──────────────┐  ┌──────────────┐  ┌──────────────────┐   ║
║  │    Spec      │  │   Balance    │  │  Gameplay Design │   ║
║  │  Architect   │  │    Agent     │  │     Agent        │   ║
║  └──────────────┘  └──────────────┘  └──────────────────┘   ║
╚════════════════════════════════════════════════════════════╝
            │
╔═══════════╪════════════════════════════════════════════════╗
║           ▼         實作層（5個）                           ║
║  ┌──────────────┐  ┌──────────────┐  ┌──────────────────┐  ║
║  │  Go Server   │  │   Gameplay   │  │    UI / HUD      │  ║
║  │    Agent     │  │    Agent     │  │     Agent        │  ║
║  └──────────────┘  └──────────────┘  └──────────────────┘  ║
║  ┌──────────────┐  ┌──────────────┐                         ║
║  │  Art/Sprite  │  │    Audio     │                         ║
║  │    Agent     │  │    Agent     │                         ║
║  └──────────────┘  └──────────────┘                         ║
╚════════════════════════════════════════════════════════════╝
            │
╔═══════════╪════════════════════════════════════════════════╗
║           ▼         整合層（3個）                           ║
║  ┌──────────────┐  ┌──────────────┐  ┌──────────────────┐  ║
║  │ Integration  │  │   Protocol   │  │   Build &        │  ║
║  │  Test Agent  │  │  Sync Agent  │  │  Export Agent    │  ║
║  └──────────────┘  └──────────────┘  └──────────────────┘  ║
╚════════════════════════════════════════════════════════════╝
            │
╔═══════════╪════════════════════════════════════════════════╗
║           ▼         驗證層（4個）                           ║
║  ┌──────────────┐  ┌──────────────┐  ┌──────────────────┐  ║
║  │  QA Playtest │  │   Player     │  │  Video Analysis  │  ║
║  │    Agent     │  │  Experience  │  │     Agent        │  ║
║  │              │  │    Agent     │  │                  │  ║
║  └──────────────┘  └──────────────┘  └──────────────────┘  ║
║  ┌──────────────┐                                           ║
║  │  Regression  │                                           ║
║  │  Guard Agent │                                           ║
║  └──────────────┘                                           ║
╚════════════════════════════════════════════════════════════╝
            │
╔═══════════╪════════════════════════════════════════════════╗
║           ▼         知識層（3個）                           ║
║  ┌──────────────┐  ┌──────────────┐  ┌──────────────────┐  ║
║  │   Research   │  │    Skill     │  │   Animation      │  ║
║  │    Agent     │  │  Librarian   │  │     Agent        │  ║
║  └──────────────┘  └──────────────┘  └──────────────────┘  ║
╚════════════════════════════════════════════════════════════╝
```

---

## 20 個 Agent 完整職責表

### 決策層

| Agent | 職責核心 | 主要輸出 | 觸發條件 |
|-------|---------|---------|---------|
| **game-director** | 技術架構決策、任務優先級、風險管理 | 每日任務計畫、決策記錄 | 每日開始、重大決策點 |
| **player-experience-director** | 玩家體驗決策、爽感設計、體驗循環完整性 | 體驗評估報告、優化指令 | 每個功能完成後、收到玩家錄影後 |

### 設計層

| Agent | 職責核心 | 主要輸出 | 觸發條件 |
|-------|---------|---------|---------|
| **spec-architect** | Server↔Client 協定一致性、規格文件 | 協定文件、規格差異報告 | 協定變更前後 |
| **balance-agent** | RTP 模擬、數值平衡、獎勵結構 | RTP 報告、數值設定 | 新目標物加入、數值調整 |
| **gameplay-design-agent** | 核心玩法設計、特殊機制設計、手感規格 | 玩法設計文件、手感規格 | 新功能設計階段 |

### 實作層

| Agent | 職責核心 | 主要輸出 | 觸發條件 |
|-------|---------|---------|---------|
| **go-server-agent** | Go Server 開發、遊戲邏輯、WebSocket 後端 | Go 原始碼、Server Binary | 後端功能開發 |
| **gameplay-agent** | Godot 核心玩法：Cannon、TargetManager、碰撞、AUTO | GDScript 玩法腳本 | 玩法功能開發 |
| **ui-hud-agent** | Godot UI：HUD、Panel、視覺回饋、CanvasLayer | GDScript UI 腳本 | UI 功能開發 |
| **art-sprite-agent** | 精靈圖生成、美術審核、視覺風格維護 | PNG 資產、美術報告 | 新美術需求 |
| **audio-agent** | 音效設計、BGM、音效同步 | WAV 資產、音效設定 | 新音效需求 |

### 整合層

| Agent | 職責核心 | 主要輸出 | 觸發條件 |
|-------|---------|---------|---------|
| **integration-test-agent** | 端對端驗證：Server 發訊息→Client 正確顯示 | 整合測試報告 | 每個功能完成後（必須） |
| **protocol-sync-agent** | 確認 Server 協定定義與 Client 處理完全對應 | 協定同步報告 | 協定變更後 |
| **build-export-agent** | HTML5 匯出、build 驗證、部署 | HTML5 Build、部署報告 | 每日 build |

### 驗證層

| Agent | 職責核心 | 主要輸出 | 觸發條件 |
|-------|---------|---------|---------|
| **qa-playtest-agent** | 功能測試、回歸測試、效能測試 | QA 報告、品質分數 | 每次 build 後 |
| **player-experience-agent** | 玩家視角評估：手感、清晰度、爽感密度 | 體驗評估報告 | 每個功能完成後 |
| **video-analysis-agent** | 分析玩家錄影：停頓點、困惑點、爽感點 | 影片分析報告、優化建議 | 收到玩家錄影時 |
| **regression-guard-agent** | 防止新功能破壞既有功能、自動回滾判斷 | 回歸風險報告 | 每次程式碼變更後 |

### 知識層

| Agent | 職責核心 | 主要輸出 | 觸發條件 |
|-------|---------|---------|---------|
| **research-agent** | 搜尋最新技術、業界最佳實踐、免費素材 | Skill 文件、研究筆記 | 遇到知識缺口時 |
| **skill-librarian** | 管理知識庫、維護 Skill 索引、整合失敗記錄 | Skill 索引、知識整合報告 | 定期、新 Skill 加入時 |
| **animation-agent** | 多幀動畫製作、AnimationPlayer 設定 | 動畫場景、動畫報告 | 新動畫需求 |

---

## 體驗循環（每個功能必須走完）

```
1. 設計層確認
   ├── Gameplay Design Agent：玩法設計文件
   ├── Spec Architect：協定規格
   └── Balance Agent：數值設計

2. 實作層開發
   ├── Go Server Agent：後端邏輯
   ├── Gameplay Agent：玩法 GDScript
   └── UI/HUD Agent：視覺回饋

3. 整合層驗證（不可跳過）
   ├── Integration Test Agent：端對端測試
   └── Protocol Sync Agent：協定對應確認

4. 驗證層評估（不可跳過）
   ├── QA Playtest Agent：功能測試
   ├── Player Experience Agent：體驗評估
   └── Regression Guard Agent：回歸風險

5. 只有通過步驟 3+4，功能才算「完成」
```

---

## 品質門檻（硬規則）

| 指標 | 門檻 | 負責 Agent | 違反後果 |
|------|------|-----------|---------|
| Spec Completeness | >= 95 | Spec Architect | 停止新功能開發 |
| Build Stability | >= 95 | Build Export Agent | 禁止產出展示版 |
| Visual Consistency | >= 90 | Art Sprite Agent | 禁止替換正式素材 |
| Animation Quality | >= 88 | Animation Agent | 禁止 merge |
| Audio Sync | >= 90 | Audio Agent | 重新調整觸發時機 |
| **Gameplay Feel** | **>= 85** | **Player Experience Agent** | **優先修復，停止加新功能** |
| Balance Health | >= 90 | Balance Agent | 重新模擬數值 |
| Regression Risk | <= 10 | Regression Guard Agent | 自動 rollback |
| **Integration Pass** | **100%** | **Integration Test Agent** | **功能不算完成** |
| **Video Analysis** | **每週至少 1 次** | **Video Analysis Agent** | **觸發強制錄影要求** |

---

## 新增 Agent 的核心設計原則

### Player Experience Director（新）
> 「這個功能讓玩家更爽了嗎？」是唯一的評判標準。
> 不看程式碼，不看 build 狀態，只看玩家體驗。

### Integration Test Agent（新）
> 每個功能完成後，必須驗證：
> 1. Server 發出正確的 WebSocket 訊息
> 2. Client 收到並正確顯示
> 3. 玩家操作觸發正確的 Server 邏輯
> 缺少任何一步，功能不算完成。

### Video Analysis Agent（新）
> 分析玩家錄影的三個問題：
> 1. 玩家在哪裡停頓或猶豫？（操作不直覺）
> 2. 玩家在哪裡有明顯正面反應？（值得強化）
> 3. 特效和音效是否在正確時機出現？（時機問題）

### Gameplay Agent（從 Godot Client Agent 拆出）
> 只負責核心玩法：Cannon、TargetManager、AUTO 邏輯、碰撞偵測。
> 不碰 UI，不碰 WebSocket 協定，不碰美術資產。

### UI/HUD Agent（從 Godot Client Agent 拆出）
> 只負責 HUD、Panel、視覺回饋、CanvasLayer。
> 不碰遊戲邏輯，不碰 Server 通訊。

---

## 禁止行為（全體 Agent）

1. **禁止宣稱「完成度 100%」** — 改用具體指標
2. **禁止跳過整合層** — 每個功能必須通過 Integration Test Agent
3. **禁止跳過驗證層** — 每個功能必須通過 Player Experience Agent
4. **禁止只用 `go build` 通過作為完成標準** — 這只是第一層
5. **禁止連續加功能超過 3 個而不做一次完整體驗循環**
