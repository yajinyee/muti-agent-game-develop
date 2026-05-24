# Multi-Agent Game Studio v3.0 — 完整架構

## 專案概覽
**遊戲**：吉伊卡哇：像素大討伐（捕魚機）  
**技術棧**：Go + WebSocket（Port 7777）/ Godot 4.6.2（HTML5）  
**架構版本**：v3.0（2026-05-24，術業有專攻——每個子系統有自己的 Agent）

---

## 設計原則

> **術業有專攻。** 每個子系統、每個美術環節、每個製作方式，都有自己的 Agent。
> Agent 越專精，輸出品質越高，問題越容易定位。

---

## 完整 Agent 清單（38 個）

### 🎯 決策層（2個）
| Agent | 職責 |
|-------|------|
| `game-director` | 技術架構決策、任務優先級、風險管理 |
| `player-experience-director` | 玩家體驗決策、爽感設計、體驗循環完整性 |

---

### 📐 設計層（4個）
| Agent | 職責 |
|-------|------|
| `spec-architect` | Server↔Client 協定一致性、規格文件 |
| `balance-agent` | RTP 模擬、數值平衡、獎勵結構 |
| `gameplay-design-agent` | 核心玩法設計、特殊機制設計、手感規格 |
| `target-design-agent` | 目標物設計：倍率、HP、行為、視覺主題 |

---

### ⚙️ Server 實作層（4個）
| Agent | 職責 | 主要檔案 |
|-------|------|---------|
| `server-core-agent` | 遊戲主循環、狀態機、玩家管理 | `game.go`, `hub.go` |
| `server-combat-agent` | 擊破判定、RTP 計算、獎勵分配 | `combat.go`, `target.go` |
| `server-event-agent` | BOSS 系統、Bonus 系統、特殊事件 | `boss.go`, `bonus.go`, `lucky_*.go` |
| `server-infra-agent` | WebSocket Hub、Store、Config、部署 | `ws/`, `store/`, `config/` |

---

### 🎮 Client 玩法層（5個）
| Agent | 職責 | 主要檔案 |
|-------|------|---------|
| `cannon-agent` | 射擊系統：投射物、AUTO、手感、Hit Stop | `Cannon.gd`, `BulletPool.gd` |
| `target-system-agent` | 目標物生命週期：生成、移動、受擊、擊破 | `TargetManager.gd`, `TargetPool.gd` |
| `boss-battle-agent` | BOSS 戰客戶端：進場、Phase 2、死亡動畫 | `TargetManager.gd` BOSS 部分 |
| `bonus-game-agent` | Bonus 遊戲：拔草場景、計時、結算 | `BonusGame.gd` |
| `game-state-agent` | 遊戲狀態機、訊號分發、GameManager | `GameManager.gd` |

---

### 🖥️ Client UI 層（4個）
| Agent | 職責 | 主要檔案 |
|-------|------|---------|
| `hud-core-agent` | 核心 HUD：金幣、BET、勞動值、AUTO、LOCK | `HUD.gd` |
| `lucky-panel-agent` | 150+ LuckyXxxPanel 重構、BaseLuckyPanel | `scripts/ui/Lucky*.gd` |
| `social-ui-agent` | 排行榜、公會、好友、活動 Panel | 社交相關 Panel |
| `screen-recorder-agent` | 側錄功能、REC 按鈕 | `ScreenRecorder.gd` |

---

### ✨ Client 特效層（3個）
| Agent | 職責 | 主要檔案 |
|-------|------|---------|
| `hit-effect-agent` | 命中特效、擊破粒子、獎勵跳字 | `HitEffect.gd`, `effects/` |
| `screen-effect-agent` | 螢幕震動、Hit Stop、水下 Shader、像素化過場 | `ScreenShake.gd`, `UnderwaterOverlay.gd` |
| `environment-agent` | 背景管理、氣泡層、環境動畫 | `BackgroundManager.gd`, `BubbleLayer.gd` |

---

### 🌐 Client 網路層（2個）
| Agent | 職責 | 主要檔案 |
|-------|------|---------|
| `network-agent` | WebSocket 連線、重連、心跳、訊息收發 | `NetworkManager.gd` |
| `protocol-sync-agent` | Server↔Client 訊息對應驗證 | `protocol.go` ↔ `GameManager.gd` |

---

### 🎨 美術製作層（7個）
| Agent | 職責 | 工具 |
|-------|------|------|
| `character-pixel-agent` | 角色像素圖：吉伊卡哇/小八/烏薩奇，3狀態 | `generate_pixel_art_v5.py` |
| `character-animation-agent` | 角色動畫幀：idle/attack/bigwin Spritesheet | `generate_animation_frames.py` |
| `target-pixel-agent` | 目標物像素圖：T001-T249 程式生成 | `generate_targets_v3.py` |
| `target-ai-agent` | 目標物 AI 生成：ComfyUI + SD 1.5 | `comfyui_generate_targets.py` |
| `boss-art-agent` | BOSS 動畫：B001 三狀態 Spritesheet | `generate_boss_sheet.py` |
| `background-art-agent` | 背景圖：海底/BOSS/Bonus 三種場景 | `generate_backgrounds_v2.py` |
| `ui-art-agent` | UI 元素：按鈕、圖示、字體、特效 Sprite | `generate_ui_assets.py` |

---

### 🔊 音效層（2個）
| Agent | 職責 | 工具 |
|-------|------|------|
| `sfx-agent` | 音效設計：14個 SFX 的音量、音調、同步 | `AudioManager.gd`, WAV 工具 |
| `bgm-agent` | BGM 設計：4首 BGM 的循環、切換、淡入淡出 | `AudioManager.gd`, BGM 生成 |

---

### 🔍 整合驗證層（4個）
| Agent | 職責 |
|-------|------|
| `integration-test-agent` | 端對端驗證：Server 發→Client 顯示 |
| `regression-guard-agent` | 防止新功能破壞舊功能 |
| `build-export-agent` | HTML5 匯出、build 驗證 |
| `performance-agent` | FPS、記憶體、Draw Call 監控 |

---

### 👁️ 玩家驗證層（3個）
| Agent | 職責 |
|-------|------|
| `qa-playtest-agent` | 功能測試、回歸測試 |
| `player-experience-agent` | 玩家視角評估：手感、清晰度、爽感 |
| `video-analysis-agent` | 分析玩家錄影：停頓點、爽感點、時機問題 |

---

### 📚 知識層（3個）
| Agent | 職責 |
|-------|------|
| `research-agent` | 搜尋最新技術、業界最佳實踐 |
| `skill-librarian` | 管理知識庫、維護 Skill 索引 |
| `animation-agent` | 動畫系統：AnimationPlayer、SpriteFrames |

---

## 架構圖

```
決策層          game-director ←→ player-experience-director
                      │
設計層    spec-architect  balance-agent  gameplay-design-agent  target-design-agent
                      │
Server層  server-core  server-combat  server-event  server-infra
                      │
Client層  cannon  target-system  boss-battle  bonus-game  game-state
                      │
UI層      hud-core  lucky-panel  social-ui  screen-recorder
                      │
特效層    hit-effect  screen-effect  environment
                      │
網路層    network  protocol-sync
                      │
美術層    character-pixel  character-animation  target-pixel  target-ai
          boss-art  background-art  ui-art
                      │
音效層    sfx  bgm
                      │
整合層    integration-test  regression-guard  build-export  performance
                      │
驗證層    qa-playtest  player-experience  video-analysis
                      │
知識層    research  skill-librarian  animation
```

---

## 體驗循環（每個功能必須走完）

```
1. 設計層：target-design-agent 或 gameplay-design-agent 出設計文件
2. 實作層：對應的 Server + Client Agent 實作
3. 整合層：integration-test-agent 端對端驗證（不可跳過）
4. 驗證層：player-experience-agent 體驗評估（不可跳過）
5. 通過才算完成
```

---

## 品質門檻

| 指標 | 門檻 | 負責 Agent |
|------|------|-----------|
| Spec Completeness | >= 95 | spec-architect |
| Build Stability | >= 95 | build-export-agent |
| Visual Consistency | >= 90 | character-pixel-agent |
| Animation Quality | >= 88 | character-animation-agent |
| Audio Sync | >= 90 | sfx-agent |
| Gameplay Feel | >= 85 | player-experience-agent |
| Balance Health | >= 90 | balance-agent |
| Regression Risk | <= 10 | regression-guard-agent |
| Integration Pass | 100% | integration-test-agent |
| FPS (HTML5) | >= 30 | performance-agent |

---

## 已廢棄的 Agent（v1.0/v2.0）

以下 Agent 已被拆分或合併，保留文件供參考：
- `godot-client-agent.md` → 拆分為 cannon/target-system/boss-battle/bonus-game/game-state
- `art-director.md` + `sprite-generation-agent.md` → 拆分為 7 個美術 Agent
- `art-sprite-agent.md` → 進一步拆分為 7 個美術 Agent
- `gameplay-agent.md` → 拆分為 cannon/target-system/boss-battle/bonus-game
- `ui-hud-agent.md` → 拆分為 hud-core/lucky-panel/social-ui/screen-recorder
- `audio-director.md` → 拆分為 sfx-agent/bgm-agent
