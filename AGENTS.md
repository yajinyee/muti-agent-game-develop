# Multi-Agent Game Studio — 協作架構說明

## 專案概覽

**遊戲名稱**：吉伊卡哇：像素大討伐（捕魚機）  
**技術棧**：Go + WebSocket（Port 7777）/ Godot 4.6.2（HTML5 匯出）  
**當前狀態**：完成度 100%，美術質量 100/100，規格一致性 100%，RTP 95.98%，DAY-058

---

## Agent 架構圖

```
                    ┌─────────────────┐
                    │  Game Director  │  ← 最高決策者
                    └────────┬────────┘
                             │ 指揮 / 審核
          ┌──────────────────┼──────────────────┐
          │                  │                  │
   ┌──────▼──────┐   ┌───────▼──────┐   ┌──────▼──────┐
   │   Spec      │   │   Research   │   │  QA/Playtest│
   │  Architect  │   │    Agent     │   │    Agent    │
   └──────┬──────┘   └───────┬──────┘   └──────┬──────┘
          │ 規格              │ 知識              │ 品質報告
          │                  ▼                  │
          │          ┌───────────────┐          │
          │          │ Skill         │          │
          │          │ Librarian     │          │
          │          └───────────────┘          │
          │                                     │
   ┌──────▼──────────────────────────────────────▼──────┐
   │                  開發層                              │
   │  ┌─────────────┐  ┌─────────────┐  ┌────────────┐  │
   │  │ Go Server   │  │   Godot     │  │  Balance   │  │
   │  │   Agent     │  │ Client Agent│  │   Agent    │  │
   │  └─────────────┘  └─────────────┘  └────────────┘  │
   └─────────────────────────────────────────────────────┘
          │                  │
   ┌──────▼──────┐   ┌───────▼──────┐
   │    Art      │   │    Audio     │
   │  Director   │   │  Director    │
   └──────┬──────┘   └──────────────┘
          │
   ┌──────▼──────┐   ┌──────────────┐
   │   Sprite    │   │  Animation   │
   │ Generation  │   │    Agent     │
   │   Agent     │   └──────────────┘
   └─────────────┘
```

---

## 12 個 Agent 職責摘要

| Agent | 角色 | 主要輸出 |
|-------|------|---------|
| **game-director** | 遊戲總監 | 每日任務計畫、決策記錄 |
| **spec-architect** | 規格架構師 | 協定文件、規格一致性報告 |
| **research-agent** | 研究員 | Skill 文件、研究筆記 |
| **art-director** | 美術總監 | 美術審核報告、視覺風格指南 |
| **sprite-generation-agent** | 精靈圖生成 | PNG 圖像資產、生成日誌 |
| **animation-agent** | 動畫專員 | .tscn 動畫場景、動畫報告 |
| **audio-director** | 音效總監 | 音效審核報告、音效設定 |
| **godot-client-agent** | Godot 客戶端 | GDScript、HTML5 Build |
| **go-server-agent** | Go 伺服器 | Go 原始碼、Server Binary |
| **balance-agent** | 數值平衡 | RTP 報告、數值設定 |
| **qa-playtest-agent** | QA 測試 | QA 報告、品質分數 |
| **skill-librarian** | 技能圖書館員 | Skill 索引、知識整合 |

---

## 協作流程

### 日常開發循環

```
1. Game Director 讀取昨日 nightly report
2. Game Director 更新 tasks/today-plan.md
3. 各 Agent 依任務計畫執行工作
4. QA Agent 執行測試，更新品質分數
5. 各 Agent 輸出報告到 reports/
6. Skill Librarian 整合新知識
7. Game Director 審閱，決定明日計畫
```

### 美術生成流程

```
Art Director → 定義需求
    ↓
Sprite Generation Agent → 生成圖像（暫存）
    ↓
Art Director → 審核（Visual Consistency >= 90？）
    ↓ 通過
Animation Agent → 製作動畫
    ↓
QA Agent → 動畫品質測試（Animation Quality >= 88？）
    ↓ 通過
Godot Client Agent → 整合到遊戲
```

### 協定變更流程

```
任何 Agent 提出協定變更需求
    ↓
Spec Architect → 評估影響範圍
    ↓
Game Director → 審核批准
    ↓
Spec Architect → 更新雙側文件
    ↓
Go Server Agent → 更新 Server 端
    ↓
Godot Client Agent → 更新 Client 端
    ↓
QA Agent → 回歸測試
```

---

## 品質門檻（硬規則）

| 指標 | 門檻 | 違反後果 |
|------|------|---------|
| Spec Completeness | >= 95 | 停止新功能開發 |
| Build Stability | >= 95 | 禁止產出展示版 |
| Visual Consistency | >= 90 | 禁止替換正式素材 |
| Animation Quality | >= 88 | 禁止 merge |
| Audio Sync | >= 90 | 重新調整觸發時機 |
| Gameplay Feel | >= 85 | 優先修復玩法問題 |
| Balance Health | >= 90 | 重新模擬數值 |
| Regression Risk | <= 10 | 自動 rollback |

---

## 檔案系統規範

### 目錄用途
- `agents/` — Agent 定義文件（本目錄）
- `tasks/` — 任務計畫與待辦清單
- `reports/` — 各類報告（按類型分子目錄）
- `memory/` — 專案記憶與狀態摘要
- `skills/` — 可重用技能文件
- `docs/` — 核心設計文件
- `references/` — 研究筆記與參考資料
- `failed-attempts/` — 失敗記錄（避免重複踩坑）
- `builds/` — 建置產出（daily/release）
- `client/` — Godot 4 專案
- `server/` — Go 伺服器原始碼

### 命名規範
- 報告檔案：`[type]-report-YYYY-MM-DD.md`
- Skill 檔案：`skill-[kebab-case-name].md`
- 失敗記錄：`failed-[topic]-YYYY-MM-DD.md`

---

## 緊急處理程序

### Build 崩潰
1. Go Server Agent 立即執行 `go build ./...` 確認問題
2. Godot Client Agent 確認 HTML5 匯出狀態
3. QA Agent 記錄 Regression Risk 分數
4. Game Director 決定是否 rollback

### RTP 異常
1. Balance Agent 立即執行模擬確認
2. Go Server Agent 暫停遊戲邏輯
3. Game Director 通知相關人員

### 美術品質下降
1. Art Director 發出阻擋指令
2. Sprite Generation Agent 重新生成
3. 未通過審核的素材不得進入正式目錄
