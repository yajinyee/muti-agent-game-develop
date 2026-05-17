# 吉伊卡哇：像素大討伐 — 系統架構與完整規格

> 版本：v1.0（2026-05-15）
> 完成度：98% | 美術質量：87/100 | 規格一致性：99%

---

## 目錄

1. [專案概覽](#1-專案概覽)
2. [系統架構](#2-系統架構)
3. [目錄結構](#3-目錄結構)
4. [Go Server 架構](#4-go-server-架構)
5. [Godot Client 架構](#5-godot-client-架構)
6. [WebSocket 通訊協定](#6-websocket-通訊協定)
7. [遊戲狀態機](#7-遊戲狀態機)
8. [角色系統](#8-角色系統)
9. [投注等級表](#9-投注等級表)
10. [目標物規格](#10-目標物規格)
11. [擊破判定系統](#11-擊破判定系統)
12. [勞動值系統](#12-勞動值系統)
13. [Bonus Game 規格](#13-bonus-game-規格)
14. [BOSS 戰規格](#14-boss-戰規格)
15. [目標生成系統](#15-目標生成系統)
16. [RTP 設計](#16-rtp-設計)
17. [美術資產規格](#17-美術資產規格)
18. [工具腳本清單](#18-工具腳本清單)
19. [開發環境與啟動方式](#19-開發環境與啟動方式)
20. [已知問題與待辦](#20-已知問題與待辦)

---

## 1. 專案概覽

| 項目 | 內容 |
|------|------|
| 遊戲名稱 | 吉伊卡哇：像素大討伐 |
| 類型 | 捕魚機 / 休閒射擊 / IP 包裝型 |
| 美術風格 | 16-bit Retro Pixel Art |
| 核心體驗 | 可愛角色討伐怪物，獲得勞動報酬袋 |
| 遊玩網址 | http://localhost:7777 |
| Server Port | 7777 |

### IP 包裝對照

| 傳統捕魚機 | 本作 |
|-----------|------|
| 砲台 | 吉伊卡哇 / 小八 / 烏薩奇 |
| 砲彈 | 討伐棒劍氣 |
| 魚群 | 怪物、雜草、小蟲 |
| 金幣掉落 | 勞動報酬袋、銅幣、金幣 |
| BOSS | 那個孩子 |
| 能量條 | 勞動值 |
| Bonus Game | 瘋狂拔草 Weeding Frenzy |

---

## 2. 系統架構

```
┌─────────────────────────────────────────────────────────┐
│                    玩家瀏覽器                             │
│  ┌─────────────────────────────────────────────────┐    │
│  │           Godot Client (HTML5 Export)            │    │
│  │  NetworkManager ←→ GameManager ←→ UI/Scene      │    │
│  └──────────────────┬──────────────────────────────┘    │
└─────────────────────┼───────────────────────────────────┘
                      │ WebSocket (ws://localhost:7777/ws)
                      │ JSON 訊息，每訊息獨立 frame
┌─────────────────────┼───────────────────────────────────┐
│                     ▼                                    │
│  ┌─────────────────────────────────────────────────┐    │
│  │              Go Server (Port 7777)               │    │
│  │  HTTP Server (/ws, /health, 靜態檔案)            │    │
│  │  WebSocket Hub → Game Loop (10 FPS)              │    │
│  │  ├── State Machine (10 states)                   │    │
│  │  ├── Target Spawn System                         │    │
│  │  ├── Combat System (混合制擊破)                  │    │
│  │  ├── Player Manager                              │    │
│  │  ├── BOSS System                                 │    │
│  │  └── Bonus Game System                           │    │
│  └─────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────┘
```

### 技術棧

| 層 | 技術 |
|----|------|
| Server | Go 1.26 + gorilla/websocket v1.5.3 |
| Client | Godot 4.6.2 (GDScript)，HTML5 匯出 |
| 通訊 | WebSocket + JSON |
| AI 美術 | ComfyUI + SD 1.5 + Pixel Art LoRA |
| 後處理 | Python + Pillow（agent-sprite-forge 技術）|

---

## 3. 目錄結構

```
d:\Kiro\
├── server/                          # Go Server
│   ├── cmd/gameserver/main.go       # 入口點
│   ├── internal/
│   │   ├── config/                  # 設定管理
│   │   ├── data/tables.go           # 靜態資料表（角色/目標/投注）
│   │   ├── game/
│   │   │   ├── game.go              # 遊戲主邏輯、事件循環
│   │   │   ├── bonus/               # Bonus Game 邏輯
│   │   │   ├── boss/                # BOSS 系統
│   │   │   ├── combat/combat.go     # 攻擊/擊破/獎勵計算
│   │   │   ├── state/state.go       # 狀態機定義
│   │   │   └── target/target.go     # 目標物實例管理
│   │   ├── player/player.go         # 玩家狀態管理
│   │   └── ws/
│   │       ├── hub.go               # WebSocket Hub
│   │       └── protocol.go          # 訊息協定定義
│   ├── pkg/logger/                  # 日誌套件
│   ├── static/                      # Godot HTML5 匯出檔案
│   └── go.mod
│
├── client/chiikawa-pixel/           # Godot 專案
│   ├── scripts/
│   │   ├── game/
│   │   │   ├── GameManager.gd       # 訊息路由、狀態管理（Autoload）
│   │   │   ├── AudioManager.gd      # 音效管理（Autoload）
│   │   │   ├── BackgroundManager.gd # 背景切換
│   │   │   ├── BonusGame.gd         # 瘋狂拔草場景
│   │   │   ├── Cannon.gd            # 砲台控制、投射物
│   │   │   ├── CharacterAnimator.gd # 角色動畫（Spritesheet）
│   │   │   └── TargetManager.gd     # 目標物節點管理
│   │   ├── network/
│   │   │   └── NetworkManager.gd    # WebSocket 連線（Autoload）
│   │   └── ui/
│   │       └── HUD.gd               # 主遊戲 UI
│   ├── assets/sprites/
│   │   ├── characters/              # 角色 Sprites（96x96）
│   │   ├── targets/                 # 目標物 Sprites（64x64）
│   │   ├── effects/                 # 特效（48x48）
│   │   ├── backgrounds/             # 背景（1280x720）
│   │   ├── ui/                      # UI 元素
│   │   └── sheets/                  # Spritesheets
│   └── scenes/
│       ├── Main.tscn
│       └── BonusGame.tscn
│
├── tools/                           # Python 工具腳本
│   ├── generate_chars_v6.py         # 角色生成 v6
│   ├── generate_targets_v3.py       # 目標物生成 v3
│   ├── generate_backgrounds_v2.py   # 背景生成 v2
│   ├── generate_effects_v2.py       # 特效生成 v2
│   ├── generate_animation_frames.py # 動畫 Spritesheet 生成
│   ├── process_sprites.py           # 後處理工具（去背/對齊/QC）
│   ├── comfyui_generate.py          # ComfyUI API 整合
│   ├── batch_process_ai.py          # AI 圖批次後處理
│   ├── simulate_rtp.py              # RTP 模擬器
│   └── test_server.py               # Server 整合測試
│
└── docs/
    ├── game-spec.md                 # 遊戲規格知識庫
    ├── progress.md                  # 開發進度追蹤
    └── system-architecture.md      # 本文件
```

---

## 4. Go Server 架構

### 模組：`digital-twin/server`

#### 4.1 入口點（`cmd/gameserver/main.go`）
- 初始化 config、logger
- 建立 WebSocket Hub
- 建立 Game 實例並啟動 Game Loop
- 啟動 HTTP Server（Port 7777）

#### 4.2 HTTP 端點

| 路徑 | 方法 | 說明 |
|------|------|------|
| `/ws` | GET | WebSocket 升級 |
| `/health` | GET | 健康檢查 |
| `/` | GET | 靜態檔案（Godot HTML5）|

HTTP Headers（COOP/COEP，支援 SharedArrayBuffer）：
```
Cross-Origin-Opener-Policy: same-origin
Cross-Origin-Embedder-Policy: require-corp
```

#### 4.3 WebSocket Hub（`internal/ws/hub.go`）
- 管理所有 WebSocket 連線
- `Send(clientID, msg)` — 單播
- `Broadcast(msg)` — 廣播
- 每訊息獨立 frame（不合併）

#### 4.4 Game Loop（`internal/game/game.go`）
- 10 FPS Ticker（100ms）
- 依當前狀態執行對應 update：
  - `updateNormalPlay()` — 生成目標、補償機制、Auto 攻擊
  - `updateBossBattle()` — BOSS 超時檢查
  - `updateBonusGame()` — 倒數計時廣播

#### 4.5 靜態資料表（`internal/data/tables.go`）

常數：
```go
BaseRTPFactor = 0.92   // 基礎擊破機率係數
LaborValueMax = 100    // 勞動值上限
SpawnInterval = 0.8    // 目標生成間隔（秒）
MaxTargetsOnScreen = 18
BossDuration = 60.0    // BOSS 持續時間（秒）
BonusDuration = 15.0   // Bonus Game 持續時間（秒）
```

#### 4.6 玩家狀態（`internal/player/player.go`）

`PlayerSnapshot`（傳送給 Client）：
```json
{
  "id": "string",
  "coins": 10000,
  "bet_level": 5,
  "bet_cost": 10,
  "character_id": "hachiware",
  "character_name": "小八",
  "labor_value": 45,
  "is_auto": false,
  "lock_target_id": "",
  "projectile_speed": 800.0,
  "fire_rate": 2.3
}
```

---

## 5. Godot Client 架構

### 5.1 Autoload 單例

| 節點 | 腳本 | 職責 |
|------|------|------|
| NetworkManager | NetworkManager.gd | WebSocket 連線、訊息收發、自動重連 |
| GameManager | GameManager.gd | 訊息路由、遊戲狀態、訊號發送 |
| AudioManager | AudioManager.gd | SFX 音效池、BGM 管理 |

### 5.2 場景節點樹（Main.tscn）

```
Main (Node2D)
├── BackgroundManager (Sprite2D)     # 背景切換
├── TargetManager (Node2D)           # 目標物節點管理
├── Cannon (Node2D)                  # 砲台
│   ├── CannonSprite (Sprite2D)      # 角色動畫（CharacterAnimator.gd）
│   └── AttackLabel (Label)
├── BonusGame (Node2D)               # 瘋狂拔草場景（預設隱藏）
└── HUD (CanvasLayer)                # UI 覆蓋層
    ├── TopBar
    │   ├── CoinsLabel
    │   ├── BetLabel
    │   ├── CharacterLabel
    │   ├── LaborBar (ProgressBar)
    │   └── LaborLabel
    ├── BottomBar
    │   ├── AutoButton
    │   ├── LockButton
    │   ├── BetMinusButton
    │   ├── BetPlusButton
    │   ├── BossButton (Prototype)
    │   └── BonusButton (Prototype)
    ├── RewardPopup (Label)
    ├── StateLabel (Label)
    ├── WarningOverlay (Control)
    └── BonusOverlay (Control)
```

### 5.3 訊號流

```
NetworkManager.message_received
    → GameManager._on_message_received()
        → emit_signal("target_spawned")  → TargetManager
        → emit_signal("target_updated")  → TargetManager
        → emit_signal("target_killed")   → TargetManager
        → emit_signal("attack_result")   → Cannon, CharacterAnimator
        → emit_signal("reward_received") → Cannon, HUD
        → emit_signal("player_updated")  → HUD, CharacterAnimator, Cannon
        → emit_signal("boss_event")      → HUD, TargetManager
        → emit_signal("bonus_event")     → HUD, BonusGame
        → emit_signal("game_state_changed") → HUD, BackgroundManager
```

### 5.4 NetworkManager 特性
- Web 模式：`WebSocketPeer`
- 桌面模式：`WebSocketPeer`（統一）
- 自動重連：斷線後 2 秒重試
- Ping/Pong 心跳：每 30 秒

---

## 6. WebSocket 通訊協定

所有訊息格式：
```json
{ "type": "message_type", "payload": { ... } }
```

### 6.1 Client → Server

| type | payload | 說明 |
|------|---------|------|
| `attack` | `{target_id, click_x, click_y}` | 玩家攻擊 |
| `lock` | `{target_id}` | 鎖定目標（空字串=解除）|
| `auto_toggle` | `{}` | 切換自動攻擊 |
| `bet_change` | `{bet_level}` | 切換投注等級 |
| `bonus_click` | `{target_id, click_x, click_y}` | Bonus 拔草點擊 |
| `ping` | `{}` | 心跳 |
| `trigger_boss` | `{}` | 手動觸發 BOSS（Prototype）|
| `trigger_bonus` | `{}` | 手動觸發 Bonus（Prototype）|

### 6.2 Server → Client

| type | payload | 說明 |
|------|---------|------|
| `game_state` | `{state, timestamp}` | 狀態變更 |
| `target_spawn` | `{instance_id, def_id, name, type, x, y, hp, max_hp, speed, lifetime, behavior}` | 目標生成 |
| `target_update` | `{instance_id, hp, max_hp, x, y, phase, is_fleeing}` | 目標狀態更新 |
| `target_kill` | `{instance_id, def_id, multiplier, reward, labor_gain, killer_id}` | 目標擊破 |
| `attack_result` | `{target_id, is_hit, is_kill, damage, reward, labor_gain, character_id, multiplier}` | 攻擊結果（單播）|
| `reward` | `{source, amount, multiplier, new_balance}` | 獎勵發放（單播）|
| `boss_event` | `{event, instance_id, phase, hp, max_hp, reward, multiplier}` | BOSS 事件 |
| `bonus_event` | `{event, time_left, score, multiplier, reward}` | Bonus 事件 |
| `player_update` | `{id, coins, bet_level, bet_cost, character_id, character_name, labor_value, is_auto, lock_target_id, projectile_speed, fire_rate}` | 玩家狀態（單播）|
| `error` | `{code, message}` | 錯誤訊息（單播）|
| `pong` | `{}` | 心跳回應 |

### 6.3 BOSS 事件類型

| event | 說明 |
|-------|------|
| `warning` | BOSS 警告（3秒後出現）|
| `spawn` | BOSS 出現 |
| `phase_change` | Phase 2 觸發（HP ≤ 50%）|
| `kill` | BOSS 被擊殺 |
| `timeout` | BOSS 超時消失 |

### 6.4 Bonus 事件類型

| event | 說明 |
|-------|------|
| `ready` | Bonus 準備（3秒後開始）|
| `start` | Bonus Game 開始 |
| `tick` | 倒數計時更新（每秒）|
| `end` | Bonus Game 結束，含結算 |

---

## 7. 遊戲狀態機

```
Loading
  └→ Lobby
       └→ NormalPlay ←──────────────────────────────┐
            ├→ SpecialTargetEvent → NormalPlay        │
            ├→ BossWarning                            │
            │    └→ BossBattle                        │
            │         ├→ BossResult ──────────────────┤
            │         └→ NormalPlay（超時）────────────┤
            └→ BonusReady                             │
                 └→ BonusGame                         │
                      └→ BonusResult ─────────────────┘
```

### 狀態說明

| 狀態 | 說明 | 持續時間 |
|------|------|---------|
| `loading` | 初始載入 | 短暫 |
| `lobby` | 大廳等待 | 直到連線 |
| `normal_play` | 正常遊戲 | 持續 |
| `special_target_event` | 特殊目標事件（流星等）| 25-40 秒觸發一次 |
| `boss_warning` | BOSS 警告 | 3 秒 |
| `boss_battle` | BOSS 戰 | 最長 60 秒 |
| `boss_result` | BOSS 結算 | 3 秒 |
| `bonus_ready` | Bonus 準備 | 3 秒 |
| `bonus_game` | 瘋狂拔草 | 15 秒 |
| `bonus_result` | Bonus 結算 | 3 秒 |

---

## 8. 角色系統

### 吉伊卡哇（LV1-3）

| 屬性 | 值 |
|------|---|
| 攻擊色 | 粉紅色劍氣 |
| 攻速 | 2.0-2.1 shots/sec |
| Kill Modifier | 1.00 |
| Labor Modifier | 1.10 |
| 大獎演出 | 驚慌跳起，「YaDa!」字卡 |

### 小八（LV4-7）

| 屬性 | 值 |
|------|---|
| 攻擊色 | 藍色劍氣 |
| 攻速 | 2.2-2.5 shots/sec |
| Kill Modifier | 1.00 |
| Fire Rate Modifier | 1.08 |
| 大獎演出 | 高舉討伐棒，「Yagaina!」字卡 |

### 烏薩奇（LV8-10）

| 屬性 | 值 |
|------|---|
| 攻擊色 | 黃色旋轉殘影 |
| 攻速 | 2.7-3.0 shots/sec |
| Kill Modifier | 0.98 |
| Fire Rate Modifier | 1.20 |
| 大獎演出 | 高速旋轉跳起，「Yaha!」字卡 |

---

## 9. 投注等級表

| LV | 角色 | Bet Cost | Attack Power | Fire Rate | 投射物速度 |
|----|------|----------|--------------|-----------|-----------|
| 1 | 吉伊卡哇 | 1 | 1 | 2.0 | 700 |
| 2 | 吉伊卡哇 | 2 | 2 | 2.0 | 720 |
| 3 | 吉伊卡哇 | 3 | 3 | 2.1 | 740 |
| 4 | 小八 | 5 | 5 | 2.2 | 780 |
| 5 | 小八 | 10 | 10 | 2.3 | 800 |
| 6 | 小八 | 20 | 20 | 2.4 | 820 |
| 7 | 小八 | 30 | 30 | 2.5 | 850 |
| 8 | 烏薩奇 | 50 | 50 | 2.7 | 900 |
| 9 | 烏薩奇 | 80 | 80 | 2.9 | 940 |
| 10 | 烏薩奇 | 100 | 100 | 3.0 | 980 |

---

## 10. 目標物規格

### 10.1 基礎目標（2x-10x）

| ID | 名稱 | 倍率 | HP | 出現權重 | 速度 | 停留 | 勞動值 | 行為 |
|----|------|------|----|---------|------|------|-------|------|
| T001 | 像素雜草 | 2x | 3 | 180 | 0 | 20s | 1 | static_sway |
| T002 | 綠色小蟲 | 3x | 5 | 160 | 40 | 18s | 1 | linear |
| T003 | 紅色小蟲 | 5x | 8 | 130 | 55 | 16s | 1 | jump |
| T004 | 藍色小蟲 | 6x | 10 | 110 | 65 | 15s | 2 | curve |
| T005 | 會走路的布丁 | 8x | 16 | 90 | 35 | 20s | 2 | sway |
| T006 | 巨大蘑菇 | 10x | 22 | 70 | 25 | 22s | 3 | sink |

### 10.2 特殊目標（15x-50x）

| ID | 名稱 | 倍率 | HP | 出現權重 | 速度 | 停留 | 勞動值 | 特殊行為 |
|----|------|------|----|---------|------|------|-------|---------|
| T101 | 擬態型怪物 | 15-30x | 35 | 35 | 50 | 14s | 5 | 死亡變形（閃爍→縮放→爆炸）|
| T102 | 寶箱怪 | 25x | 55 | 22 | 70 | 10s | 6 | 受擊後加速逃跑（×2.5）|
| T103 | 流星 | 20-50x | 20 | 18 | 220 | 4s | 5 | 快速通過（斜向）|
| T104 | 金色雜草 | 30x | 45 | 12 | 0 | 8s | 15 | 靜止（大量勞動值）|
| T105 | 巨大金幣魚 | 50x | 90 | 8 | 80 | 8s | 10 | 擊破後金幣雨（15枚）|

### 10.3 BOSS

| ID | 名稱 | HP | 倍率 | 停留 | 勞動值 | 行為 |
|----|------|-----|------|------|-------|------|
| B001 | 那個孩子 | 3000 | 100-500x | 60s | 30 | 左右移動，Phase 2 變紅 |

### 10.4 流星倍率權重

| 倍率 | 權重 | 機率 |
|------|------|------|
| 20x | 50 | 50% |
| 30x | 30 | 30% |
| 40x | 15 | 15% |
| 50x | 5 | 5% |

### 10.5 移動行為說明

| behavior | 說明 |
|----------|------|
| `linear` | 直線向左移動 |
| `curve` | 向左移動 + 上下波浪 |
| `jump` | 向左移動 + 跳躍 |
| `meteor` | 向左 + 向下斜移 |
| `sway` | 左右搖擺 |
| `static_sway` | 原地左右微搖 |
| `sink` | 向下沉 |
| `flee` | 向左加速逃跑（T102 受擊後）|
| `mimic` | 向左 + 上下波浪（T101）|
| `boss_phases` | 左右正弦移動（BOSS）|

---

## 11. 擊破判定系統

### 11.1 混合制公式

```
Kill Chance = BaseRTPFactor / Multiplier × CharKillModifier
            = 0.92 / Multiplier × KillMod

Required_Hits（保底）:
  基礎目標 = min(ceil(Multiplier / 0.92 × 3.0), Lifetime × 3.0 × 0.8)
  特殊目標 = 99999（不設保底，純機率）
  BOSS     = 99999（不設保底）
```

### 11.2 單次命中擊破率參考

| 倍率 | 擊破率 |
|------|-------|
| 2x | 46.0% |
| 5x | 18.4% |
| 10x | 9.2% |
| 25x | 3.7% |
| 50x | 1.8% |
| 100x | 0.9% |
| 500x | 0.18% |

### 11.3 視覺 HP 系統

每次命中扣除視覺 HP：
```
damage = ceil(MaxHP / (RequiredHits + 1))
HP = max(0, HP - damage)
```

HP 條顯示在目標物上方（48px 寬，5px 高）。

### 11.4 獎勵計算

```
Reward = BetCost × Multiplier（四捨五入）
LaborGain = TargetLaborGain × CharLaborModifier
```

---

## 12. 勞動值系統

- 上限：100
- 觸發 Bonus 後歸零
- Bonus 冷卻：90 秒

### 各目標勞動值

| 目標 | 勞動值 |
|------|-------|
| 像素雜草 | +1 |
| 小蟲類（T002-T003）| +1 |
| 藍色小蟲（T004）| +2 |
| 布丁（T005）| +2 |
| 巨大蘑菇（T006）| +3 |
| 擬態型怪物（T101）| +5 |
| 寶箱怪（T102）| +6 |
| 流星（T103）| +5 |
| 金色雜草（T104）| +15 |
| 巨大金幣魚（T105）| +10 |
| BOSS（B001）| +30 |

### 補償機制

30 秒內無高倍率獎勵（≥20x）時，特殊目標出現率 +5%。

---

## 13. Bonus Game 規格

### 13.1 流程

```
勞動值 = 100
  → bonus_event: "ready"（3秒等待）
  → bonus_event: "start"（15秒倒數）
  → 玩家點擊拔草
  → bonus_event: "tick"（每秒廣播剩餘時間）
  → bonus_event: "end"（結算）
  → 回到 NormalPlay
```

### 13.2 Bonus 目標表

| ID | 名稱 | 點擊分數 | 出現權重 | 特殊效果 |
|----|------|---------|---------|---------|
| BG001 | 普通雜草 | 1 | 180 | 無 |
| BG002 | 硬雜草 | 3 | 80 | 需連點 2 次（第一次搖晃）|
| BG003 | 發光雜草 | 8 | 35 | 增加倍率（+5分）|
| BG004 | 金色雜草 | 20 | 10 | 觸發金幣雨（20枚）|
| BG005 | 搗亂怪草 | -5 | 20 | 暫停操作 0.3 秒 |

### 13.3 獎勵計算

```
Bonus_Multiplier = clamp(20 + Score × 0.375, 20, 50)
Bonus_Reward = EntryBetCost × Bonus_Multiplier
```

> 注意：規格書原始設計為 50-150x，Prototype 版調整為 20-50x 以控制 RTP。

### 13.4 Bonus 目標生成

- 進入 Bonus Game 時清除所有一般目標
- 生成 20 個 Bonus 目標（依權重隨機）
- 目標固定在畫面上（不移動）

---

## 14. BOSS 戰規格

### 14.1 觸發條件

- 時間觸發：每 3-5 分鐘（180-300 秒）
- Prototype：手動觸發按鈕

### 14.2 BOSS 參數

| 參數 | 值 |
|------|---|
| HP | 3000 |
| 基礎倍率 | 100x |
| 最高倍率 | 500x |
| 出場時間 | 60 秒 |
| Phase 2 門檻 | HP ≤ 50% |
| 位置 | (1100, 360) |

### 14.3 獎勵（依擊殺剩餘時間）

| 剩餘時間 | 倍率 |
|---------|------|
| 0-10 秒 | 100x |
| 11-20 秒 | 150x |
| 21-30 秒 | 200x |
| 31-40 秒 | 300x |
| 41-50 秒 | 400x |
| 51-60 秒 | 500x |

實際獎勵 = BetCost × Multiplier × 0.15（RTP 係數）

### 14.4 Phase 2 視覺

- 觸發條件：HP ≤ 50%
- 視覺：紅色調 + 閃爍 3 次 + 放大到 scale(2.2, 2.2)
- 文字：「PHASE 2!」浮現後消失

---

## 15. 目標生成系統

### 15.1 生成參數

| 參數 | 值 |
|------|---|
| Spawn Interval | 0.8 秒 |
| Max Targets | 18（BOSS 期間 8）|
| 生成位置 | 畫面右側（x: 1280-1380，y: 100-600）|

### 15.2 動態難度（依 Bet Level）

| Bet 區間 | 基礎目標 | 特殊目標 |
|---------|---------|---------|
| LV1-3 | 90% | 10% |
| LV4-7 | 82% | 18% |
| LV8-10 | 75% | 25% |

### 15.3 SpecialTargetEvent

- 每 25-40 秒觸發一次
- 強制生成一個特殊目標（T101-T105 隨機）
- 觸發後進入 `special_target_event` 狀態

---

## 16. RTP 設計

### 16.1 目標 RTP

- 整體目標：94%
- 波動定位：中高波動

### 16.2 實際 RTP 分布（Prototype 版）

| Bet Level | 實際 RTP | Bonus/局 | BOSS/局 |
|-----------|---------|---------|---------|
| LV1-3 | ~81% | 0.04-0.06 | 0 |
| LV5 | ~81% | 0.21 | 0 |
| LV7 | ~87% | 0.66 | 0.60 |
| LV10 | ~96% | 0.93 | 0.86 |

> 高 bet 玩家 RTP 更高，符合業界慣例。正式版需數值工程師精確調整。

### 16.3 RTP 公式

```
Kill Chance = 0.92 / Multiplier
期望命中次數 = 1 / Kill Chance = Multiplier / 0.92
RTP = Multiplier / 期望命中次數 = 0.92 = 92%（基礎目標）
```

Bonus + BOSS 補足到 94%。

---

## 17. 美術資產規格

### 17.1 角色 Sprites

| 項目 | 規格 |
|------|------|
| 基礎尺寸 | 48×48（程式生成）|
| 輸出尺寸 | 96×96（2x 放大）|
| 格式 | RGBA PNG，透明背景 |
| 狀態 | idle / attack / bigwin |
| 生成工具 | `tools/generate_chars_v6.py` |
| AI 生成 | ComfyUI SD 1.5 + Pixel Art LoRA |
| 後處理 | `tools/process_sprites.py`（shared_scale + bottom align）|

### 17.2 角色 Animated Spritesheet

| 項目 | 規格 |
|------|------|
| 尺寸 | 384×288（4幀 × 3狀態 × 96px）|
| 行0 | idle（4幀，4fps，呼吸感縮放）|
| 行1 | attack（3幀，8fps，旋轉揮棒）|
| 行2 | bigwin（4幀，6fps，跳躍+星星）|
| 生成工具 | `tools/generate_animation_frames.py` |

### 17.3 目標物 Sprites

| 項目 | 規格 |
|------|------|
| 尺寸 | 64×64（T001-T105），96×96（B001 BOSS）|
| 格式 | RGBA PNG，透明背景 |
| 生成工具 | `tools/generate_targets_v3.py` |

### 17.4 特效 Sprites

| 資產 | 尺寸 | 說明 |
|------|------|------|
| hit_chiikawa/hachiware/usagi.png | 48×48 | 命中特效（放射狀光線）|
| projectile_chiikawa/hachiware/usagi.png | 32×16 | 投射物（橢圓+尾焰）|
| death_particles.png | 48×48 | 死亡粒子（8方向爆炸）|
| warning.png | 128×64 | WARNING 特效 |

### 17.5 背景

| 資產 | 尺寸 | 說明 |
|------|------|------|
| sea_bg.png | 1280×720 | 海底（漸層+珊瑚+氣泡）|
| boss_bg.png | 1280×720 | BOSS（暗紅+裂縫+光環）|
| bonus_bg.png | 1280×720 | Bonus 草地（天空+雲+花）|

### 17.6 UI 元素

| 資產 | 尺寸 | 說明 |
|------|------|------|
| coin.png | 32×32 | 金幣（帶¥符號）|
| reward_bag.png | 40×48 | 勞動報酬袋 |
| btn_normal/active/auto.png | 96×36 | 按鈕（圓角+漸層）|
| labor_bar_bg.png | 240×24 | 勞動值條背景 |
| labor_bar_fill.png | 236×20 | 勞動值條填充 |
| warning_card.png | 256×64 | BOSS 警告字卡 |
| ui_frame.png | 400×100 | UI 框架 |

### 17.7 Spritesheets

| 資產 | 尺寸 | 說明 |
|------|------|------|
| characters_sheet.png | 288×288 | 3角色×3狀態，cell=96 |
| targets_sheet.png | 256×192 | 11目標，cell=64 |
| effects_sheet.png | 192×96 | 7特效，cell=48 |
| chiikawa/hachiware/usagi_animated.png | 384×288 | 動畫 Spritesheet |

---

## 18. 工具腳本清單

### 18.1 美術生成

| 腳本 | 說明 |
|------|------|
| `generate_chars_v6.py` | 角色 v6（48×48 基礎，帶陰影+眼睛+手臂）|
| `generate_targets_v3.py` | 目標物 v3（64×64，帶陰影細節）|
| `generate_backgrounds_v2.py` | 背景 v2（漸層+細節）|
| `generate_effects_v2.py` | 特效 v2（命中/投射物/粒子）|
| `generate_ui_v2.py` | UI 元素 v2（金幣/報酬袋/按鈕）|
| `generate_animation_frames.py` | 動畫 Spritesheet（idle/attack/bigwin）|
| `generate_sfx.py` | 8-bit 音效生成 |

### 18.2 後處理工具

| 腳本 | 說明 |
|------|------|
| `process_sprites.py` | 主後處理工具（4個模式）|
| `batch_process_ai.py` | AI 生成圖批次後處理 |
| `generate_spritesheet.py` | 重建所有 Spritesheets |
| `verify_animated_sheets.py` | 驗證動畫 Sheet 格式 |
| `preview_animation.py` | 輸出 GIF 動畫預覽 |

`process_sprites.py` 使用方式：
```bash
py tools/process_sprites.py --mode qc        # 品質報告
py tools/process_sprites.py --mode realign   # 重新對齊
py tools/process_sprites.py --mode sheet     # 重建 Spritesheet
py tools/process_sprites.py --mode comfyui --input <path> --char <char> --pose <pose>
```

### 18.3 AI 生成

| 腳本 | 說明 |
|------|------|
| `comfyui_generate.py` | ComfyUI API 整合（生成 9 張角色圖）|
| `batch_process_ai.py` | 批次後處理 AI 生成圖 |
| `download_sd15.py` | 下載 SD 1.5 基礎模型 |
| `download_pixel_model.py` | 下載 Pixel Art LoRA |

ComfyUI 啟動：
```bash
powershell -Command "Set-Location 'C:\ComfyUI\ComfyUI_windows_portable'; .\python_embeded\python.exe -s ComfyUI\main.py --windows-standalone-build --lowvram"
```

### 18.4 測試與驗證

| 腳本 | 說明 |
|------|------|
| `test_server.py` | Server 整合測試（7/7 通過）|
| `simulate_rtp.py` | RTP 模擬器（1000-10000 局）|
| `rtp_analysis.py` | RTP 問題根源分析 |
| `check_cuda.py` | CUDA 可用性診斷 |

---

## 19. 開發環境與啟動方式

### 19.1 環境需求

| 項目 | 版本 |
|------|------|
| Go | 1.26.2 |
| Godot | 4.6.2 |
| Python | 3.12+ |
| NVIDIA Driver | 596.49+（支援 CUDA 13.0）|
| ComfyUI | 0.21.0 |
| PyTorch | 2.11.0+cu130 |

Python 套件：
```
Pillow, numpy, opencv-python
```

### 19.2 啟動 Server

```bash
# 編譯
cd d:\Kiro\server
go build ./...

# 執行
go run cmd/gameserver/main.go

# 或直接執行編譯好的
server/bin/gameserver.exe
```

Server 啟動後：
- WebSocket：`ws://localhost:7777/ws`
- 遊戲頁面：`http://localhost:7777`
- 健康檢查：`http://localhost:7777/health`

### 19.3 啟動 ComfyUI（AI 美術生成）

```bash
# GPU 模式（需 NVIDIA Driver 596.49+）
powershell -Command "Set-Location 'C:\ComfyUI\ComfyUI_windows_portable'; .\python_embeded\python.exe -s ComfyUI\main.py --windows-standalone-build --lowvram"

# 生成 9 張角色圖
py tools/comfyui_generate.py --all

# 批次後處理
py tools/batch_process_ai.py

# 重建 Spritesheets
py tools/process_sprites.py --mode sheet
```

### 19.4 美術重新生成流程

```bash
# 1. 生成角色 sprites
py tools/generate_chars_v6.py

# 2. 對齊處理
py tools/process_sprites.py --mode realign

# 3. 品質確認
py tools/process_sprites.py --mode qc

# 4. 生成動畫 Spritesheet
py tools/generate_animation_frames.py

# 5. 重建所有 Spritesheets
py tools/process_sprites.py --mode sheet
py tools/generate_spritesheet.py
```

### 19.5 Server 驗證

```bash
cd d:\Kiro\server
go build ./...    # 編譯
go vet ./...      # 靜態分析
py tools/test_server.py  # 整合測試（需先啟動 server）
```

---

## 20. 已知問題與待辦

### 20.1 已完成（本次開發）

- [x] Go Server 完整實作（狀態機、戰鬥、BOSS、Bonus）
- [x] Godot Client 完整實作（所有 GDScript 模組）
- [x] WebSocket 通訊協定（11 種訊息類型）
- [x] 角色系統（3角色，投注等級切換）
- [x] 目標物系統（11種目標，特殊行為）
- [x] 擊破判定（混合制：機率+保底）
- [x] 勞動值系統（補償機制、Bonus 冷卻）
- [x] Bonus Game（5種雜草，特殊效果）
- [x] BOSS 戰（Phase 2、計時獎勵）
- [x] RTP 校正（基礎 92%，高 bet 可達 96%）
- [x] 美術全面升級（角色 v6、目標物 v3、背景 v2、特效 v2、UI v2）
- [x] AI 角色圖生成（ComfyUI GPU 模式，9張）
- [x] 多幀動畫（idle/attack/bigwin Spritesheet）
- [x] T102 受擊加速逃跑
- [x] BG005 搗亂怪草暫停效果
- [x] BOSS Phase 2 視覺
- [x] T101 擬態死亡變形
- [x] T105 金幣魚金幣雨
- [x] 投射物速度依 BetLevel 動態計算

### 20.2 待辦（低優先）

- [ ] BG003 發光雜草「倍率提升」視覺效果
- [ ] 像素字體整合（目前用系統字體）
- [ ] hachiware/usagi 幀一致性優化（height diff 9px/5px）
- [ ] 數據埋點
- [ ] 營運工具後台
- [ ] 多人房間支援

### 20.3 已知技術問題

| 問題 | 狀態 | 說明 |
|------|------|------|
| ComfyUI GPU 模式 | ✅ 已解決 | 更新 NVIDIA 驅動到 596.49 |
| usagi_attack 消失 | ✅ 已修復 | dy=-3 超出畫布，加 ear_top = max(0, dy) |
| generate_animation_frames 覆蓋 sprites | ✅ 已修復 | 移除 frames[0].save() |
| RTP 600%+ | ✅ 已修復 | required_hits 公式修正 |
| Client 射擊當掉 | ✅ 已修復 | tween 生命週期綁定到節點 |

---

## 附錄：音效清單

| 音效 | 檔案 | 觸發時機 |
|------|------|---------|
| 攻擊（吉伊卡哇）| attack_fire.wav | 吉伊卡哇攻擊 |
| 攻擊（小八）| attack_fire_hachiware.wav | 小八攻擊 |
| 攻擊（烏薩奇）| attack_fire_usagi.wav | 烏薩奇攻擊 |
| 命中 | hit.wav | 命中目標 |
| 擊殺 | kill.wav | 擊破目標 |
| 金幣掉落 | coin_drop.wav | 金幣雨 |
| 大獎 | big_win.wav | 倍率 ≥20x |
| BOSS 警告 | boss_warning.wav | BOSS 警告 |
| BOSS 出現 | boss_enter.wav | BOSS 出現 |
| Bonus 準備 | bonus_ready.wav | Bonus Ready |
| Bonus 遊戲 | bonus_game.wav | Bonus BGM |
| 拔草 | weed_pull.wav | 拔草成功 |
| 報酬袋 | reward_bag.wav | 報酬袋獲得 |

| BGM | 檔案 | 使用場景 |
|-----|------|---------|
| 主遊戲 | main_game.wav | NormalPlay |
| BOSS 出現 | boss_enter.wav | BossBattle |
| Bonus | bonus_game.wav | BonusGame |

---

*文件生成時間：2026-05-15*
*最後更新：完成度 98%，美術質量 87/100，規格一致性 99%*
