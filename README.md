# 吉伊卡哇：像素大討伐 🐾

> **Multi-Agent Game Studio** — 由 12 個 AI Agent 協作開發的像素風格捕魚機遊戲

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/yajinyee/muti-agent-game-develop)
[![Quality Gates](https://img.shields.io/badge/quality%20gates-8%2F8-brightgreen)](https://github.com/yajinyee/muti-agent-game-develop/wiki/Quality-Gates)
[![Gameplay Feel](https://img.shields.io/badge/gameplay%20feel-100%2F100-brightgreen)](https://github.com/yajinyee/muti-agent-game-develop/wiki/Quality-Gates)
[![Art Quality](https://img.shields.io/badge/art%20quality-100%2F100-brightgreen)](https://github.com/yajinyee/muti-agent-game-develop/wiki/Quality-Gates)
[![RTP](https://img.shields.io/badge/RTP-95.93%25-blue)](https://github.com/yajinyee/muti-agent-game-develop/wiki/Game-Spec)
[![License](https://img.shields.io/badge/license-private-lightgrey)](LICENSE)

---

## 🎮 遊戲介紹

以日本人氣 IP「吉伊卡哇」為主題的像素風格捕魚機遊戲。玩家選擇吉伊卡哇、小八或烏薩奇，射擊畫面上的各種目標物獲得倍率獎勵，觸發 Bonus Game 或挑戰 BOSS 贏取大獎。

| 項目 | 內容 |
|------|------|
| 遊戲類型 | 捕魚機（Fish Shooting Game）|
| 目標平台 | Web（HTML5）|
| 開發狀態 | **100% 完成** ✅ |
| Server | Go + WebSocket，Port 7777 |
| Client | Godot 4.6.2，HTML5 匯出 |
| 美術 | AI 生成像素風格（ComfyUI + SD 1.5 + Pixel Art LoRA）|
| 架構 | RedisStore + Docker + 多房間支援，生產就緒 |

---

## ✨ 遊戲特色

- **3 位角色**：吉伊卡哇 / 小八（ハチワレ）/ 烏薩奇（うさぎ），各有專屬攻擊音效與動畫
- **11 種目標物**：T001–T105 普通目標 + B001 BOSS，倍率 2x–500x
- **Bonus Game**：瘋狂拔草小遊戲，5 種雜草類型各有特殊行為
- **BOSS 戰**：B001 兩階段 BOSS，Phase 2 變紅加速，真實計時獎勵 100-500x
- **Gameplay Juice**：Screen Shake、Hit Stop、命中特效、子彈拖尾、全畫面閃光、連擊系統
- **精準 RTP**：蒙地卡羅模擬校正至 95.93%（目標 92–96%）
- **多房間支援**：RoomManager + LobbyManager，支援多人同時遊玩
- **Redis 水平擴展**：RedisStore 完整實作，支援 Docker 部署
- **排行榜系統**：即時排行榜，每 10 秒廣播更新
- **成就系統**：多種成就類型，彩色通知動畫
- **數據埋點**：完整 Analytics 系統，JSONL 日誌 + HTTP 端點

---

## 🏗️ 技術架構

```
Browser (HTML5)
└── Godot 4.6.2 Client
    ├── Main.tscn（主遊戲場景）
    ├── BonusGame.tscn（Bonus 場景）
    └── WebSocket ←→ Go Server (Port 7777)
                      ├── GameServer（連線管理）
                      ├── GameRoom（房間狀態）
                      ├── RTPEngine（95.93%）
                      ├── BonusEngine
                      └── BOSSEngine
```

---

## 🚀 快速開始

### 環境需求

- Go 1.21+
- Godot 4.6.2
- Python 3.x（QA 工具）
- Docker（可選，用於生產部署）

### 啟動 Server（開發模式）

```bash
cd server
go run ./cmd/gameserver/main.go
# Server 在 ws://localhost:7777 啟動
# 健康檢查：http://localhost:7777/health
```

### 啟動 Server（Docker 模式）

```bash
# 一鍵啟動 Game Server + Redis
docker-compose up -d

# 查看日誌
docker-compose logs -f game-server
```

### 執行 QA 檢查

```powershell
$env:PYTHONUTF8="1"
py tools/qa_check.py
```

### 執行 Go 測試

```bash
cd server
go test ./...
go vet ./...
```

### 執行每日 Build

```powershell
powershell -File tools/daily_build.ps1
```

---

## 📁 專案結構

```
.
├── agents/          # 12 個 AI Agent 定義文件
├── audio/           # 音效資產與設定（14 個音效）
├── builds/          # 建置產出（daily / release）
├── client/          # Godot 4.6.2 專案
│   └── chiikawa-pixel/
│       ├── scenes/  # Main.tscn, BonusGame.tscn
│       ├── scripts/ # GDScript（GameManager, Cannon, TargetManager...）
│       └── assets/  # 像素圖像、音效
├── docs/            # 核心設計文件與規格
├── memory/          # 專案記憶（project-memory.md）
├── reports/         # QA、動畫、音效、每日報告
├── server/          # Go 伺服器原始碼
├── skills/          # 12 個可重用技能文件
├── tasks/           # 任務計畫（today-plan.md, backlog.md）
├── tools/           # 自動化工具
│   ├── qa_check.py          # QA 自動化（8 項指標）
│   ├── animation_pipeline.py # 動畫品質檢查
│   └── daily_build.ps1      # 每日 Build 自動化
└── AGENTS.md        # Multi-Agent Studio 架構說明
```

---

## 🤖 Multi-Agent Studio

本專案採用 **12 個 AI Agent 協作**的開發架構，每個 Agent 有明確職責：

| Agent | 職責 |
|-------|------|
| Game Director | 最高決策者，每日任務計畫 |
| Spec Architect | 規格文件、協定管理 |
| Research Agent | 技術研究、最佳實踐 |
| Art Director | 美術審核、視覺風格 |
| Sprite Generation | AI 像素美術生成 |
| Animation Agent | 角色動畫、Spritesheet |
| Audio Director | 音效審核、BGM Layer |
| Godot Client Agent | GDScript、HTML5 Build |
| Go Server Agent | Go Server、WebSocket |
| Balance Agent | RTP 校正、數值平衡 |
| QA/Playtest Agent | 自動化測試、品質分數 |
| Skill Librarian | 知識庫管理、Skill 索引 |

每 4 小時人類時間 = 1 AI 日，自主執行完整開發循環。

---

## 📊 品質分數（DAY-034 最終確認）

| 指標 | 分數 | 門檻 | 狀態 |
|------|------|------|------|
| Build Stability | 100 | ≥95 | ✅ |
| Visual Consistency | 100 | ≥90 | ✅ |
| Animation Quality | 100 | ≥88 | ✅ |
| Balance Health | 96 | ≥90 | ✅ |
| Audio Sync | 97 | ≥90 | ✅ |
| Gameplay Feel | 100 | ≥85 | ✅ |
| Spec Completeness | 100 | ≥95 | ✅ |
| Regression Risk | 5 | ≤10 | ✅ |

**8/8 全部通過 🎉 — 完成度 100%，美術質量 100/100**

---

## 📖 Wiki

完整文件請參考 [GitHub Wiki](https://github.com/yajinyee/muti-agent-game-develop/wiki)：

- [🏠 Home](https://github.com/yajinyee/muti-agent-game-develop/wiki) — 專案概覽
- [🏗️ Architecture](https://github.com/yajinyee/muti-agent-game-develop/wiki/Architecture) — 系統架構
- [🎮 Game Spec](https://github.com/yajinyee/muti-agent-game-develop/wiki/Game-Spec) — 遊戲規格
- [🤖 Agent System](https://github.com/yajinyee/muti-agent-game-develop/wiki/Agent-System) — Agent 架構
- [🔀 Git Workflow](https://github.com/yajinyee/muti-agent-game-develop/wiki/Git-Workflow) — 開發流程
- [📅 Development Log](https://github.com/yajinyee/muti-agent-game-develop/wiki/Development-Log) — 每日記錄
- [✅ Quality Gates](https://github.com/yajinyee/muti-agent-game-develop/wiki/Quality-Gates) — 品質門檻
- [📚 Skills & Knowledge](https://github.com/yajinyee/muti-agent-game-develop/wiki/Skills-Knowledge) — 技能庫

---

## 🏷️ Issue Labels

| 分類 | Labels |
|------|--------|
| 類型 | `type: feat` `type: fix` `type: chore` `type: docs` `type: art` `type: balance` |
| 優先級 | `priority: P0` `priority: P1` `priority: P2` `priority: P3` |
| Agent | `agent: go-server` `agent: godot-client` `agent: art-director` 等 |
| 狀態 | `status: in-progress` `status: blocked` `status: done` |

---

## 📝 開發日誌摘要

| DAY | 日期 | 主要成就 |
|-----|------|---------|
| DAY-001 | 2026-05-17 | Multi-Agent Studio Phase 1–7 完成，8/8 品質門檻首次通過 |
| DAY-007 | 2026-05-22 | Gameplay Juice 系統（Screen Shake + Hit Stop + 特效），第 12 個 Skill |
| DAY-010 | 2026-05-18 | 排行榜系統（Server + Client 完整）|
| DAY-013 | 2026-05-18 | BOSS 進場預覽 UI + Go 單元測試（9個）+ 像素風格遊戲邊框 |
| DAY-016 | 2026-05-18 | 數據埋點系統（Analytics + JSONL 日誌）|
| DAY-018 | 2026-05-19 | BOSS 戰 BGM + 完整 BGM 切換系統 + 資產預載入 |
| DAY-019 | 2026-05-19 | 效能監控面板升級 + 多房間架構 Phase 1 |
| DAY-020 | 2026-05-19 | 大廳 UI + Client 多房間支援 |
| DAY-026 | 2026-05-19 | Redis 水平擴展架構 + Store 模組（MemoryStore + RedisStore）|
| DAY-028 | 2026-05-19 | RedisStore 完整實作 + B001 BOSS 完整動畫集 |
| DAY-031 | 2026-05-19 | Docker 部署配置 + 目標物游泳動畫 |
| DAY-032 | 2026-05-19 | 目標物倍率標籤 + Server 協定擴展 |
| DAY-033 | 2026-05-19 | 高倍率目標光暈效果（30x+ 金色，50x+ 橙紅）|
| DAY-033b | 2026-05-19 | 目標物進場動畫 + 美術質量 100/100 |
| **DAY-034** | **2026-05-19** | **最終整合確認，完成度 100%，GitHub 同步** |

---

*由 Multi-Agent Game Studio 自主開發維護 · 最後更新：DAY-034（2026-05-19）*
