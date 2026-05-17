# Autonomous Multi-Agent Game Studio 規格書

> 專案：吉伊卡哇：像素大討伐 / Godot + Go + AI Asset Pipeline  
> 目標：將現有遊戲雛型升級為可自主開發、測試、修正、生成素材、沉澱 Skill 的 Multi-Agent 遊戲開發工廠  
> 版本：v1.0  
> 日期：2026-05-17  
> 適用環境：Kiro / Vibe Coding / Windows / Go Server / Godot Client / Python Tools / ComfyUI / Git Repo

---

## 0. Executive Summary

本規格書定義一套 **Autonomous Multi-Agent Game Studio**，目標不是單純增加幾個 AI 助手，而是把目前已完成度很高的遊戲專案，升級成一套可以每日自主運作的 AI 研發工廠。

系統必須能做到：

1. 自主讀取現有 repo、docs、progress、測試報告與素材狀態。
2. 自主拆解每日開發任務。
3. 由多個專職 Agent 平行開發：規格、研究、美術、動畫、音效、Godot、Go Server、數值、QA。
4. 自主產生 sprite、spritesheet、動畫 preview、音效、BGM layer、Godot 接入 patch、Go Server patch。
5. 自主跑 build、測試、RTP 模擬、遊戲錄影、截圖驗收。
6. 自主根據品質分數決定重做、修正、rollback 或 merge。
7. 每日輸出 playable build、demo video、quality report、next-day plan。
8. 每日沉澱 lessons learned、skills、memory、checklist，讓系統越做越好。

核心原則：

> 不是 prompt 寫得多漂亮，而是 agent 每天能不能比昨天更好一點。

---

## 1. 專案背景與現況

現有專案已具備完整遊戲雛型與基礎工具鏈，包含：

- Go Server
- Godot HTML5 Client
- WebSocket 通訊
- 遊戲狀態機
- 角色系統
- 目標物系統
- BOSS 戰
- Bonus Game
- RTP / Balance 模擬
- ComfyUI / Python 美術生成工具
- Spritesheet 生成工具
- 音效生成工具
- Server 測試工具
- 基礎 docs 與 progress 記錄

因此本專案不是從零重寫，而是在既有架構外層新增一個 **Multi-Agent Autonomous Development Layer**。

---

## 2. 專案最終目標

### 2.1 優先順序

使用者指定優先順序：

```text
D > B > C > A
```

意義如下：

1. **D：可複製的遊戲開發工廠**  
   不只做吉伊卡哇，也能複用到其他魚機、老虎機、Boss、小遊戲、Prototype。

2. **B：接近正式商業品質的遊戲**  
   追求接近 Steam 遊戲品質的美術、動畫、音效、操作手感、完整度。

3. **C：公司內部 AI 製程示範**  
   能展現 AI 如何縮短開發流程、建立方法論、累積 know-how。

4. **A：快速 Prototype**  
   可快速產出可玩 Demo，但不是最高優先。

### 2.2 成功定義

本系統成功不是指「一次生成完美遊戲」，而是：

- 每天能自主開發。
- 每天能產出可展示結果。
- 每天能量化品質。
- 每天能自我修正。
- 每天能沉澱 skill。
- 長期能支撐不同遊戲題材複製開發。

---

## 3. 核心設計原則

### 3.1 Autonomous Studio 模式

本系統採用 **Autonomous Studio** 模式。

Human 只提供：

- 高層方向
- 產品偏好
- 最終風格判斷
- 必要時的重大決策修正

Agent 負責：

- 任務拆解
- 規格擴充
- 研究參考
- 程式修改
- 素材生成
- 動畫製作
- 音效製作
- 測試驗收
- 自我修正
- skill 沉澱

### 3.2 Vibe Coding 友善

因為使用者採用 Vibe Coding，且偏好 Kiro + 高 token 預算，因此規格設計應允許：

- 大量上下文讀取
- 長規格文件
- 長 implementation prompt
- 多 agent 平行討論
- 反覆生成與重做
- 自動整理工作紀錄
- 自動沉澱 know-how

### 3.3 不怕重做，但不能亂改

本系統允許大量迭代，但不能失控。

必須具備：

- branch isolation
- test before merge
- quality gate
- rollback plan
- failed-attempt log
- risk classification

### 3.4 程式碼與測試結果優先於文件

Source of Truth 原則：

```text
repo 實際程式碼 + 測試結果 > docs > agent 記憶 > 對話紀錄
```

當文件與程式碼衝突時：

1. 以程式碼與測試結果為準。
2. 建立 discrepancy report。
3. 由 Spec Architect Agent 更新文件。

---

## 4. 系統總體架構

```text
Human Vision
  ↓
Game Director Agent
  ↓
Spec / Research / Art / Animation / Audio / Code / Balance / QA Agents
  ↓
Execution Tool Layer
  - Kiro
  - Git
  - Godot
  - Go Server
  - Python Tools
  - ComfyUI
  - Browser Research
  - Test Runner
  ↓
Build / Playtest / Quality Gate
  ↓
Memory / Skills / Reports / Lessons Learned
  ↓
Next Iteration
```

---

## 5. 系統分層

### L0：Human Vision Layer

人類只輸入：

- 專案方向
- 題材偏好
- 產品感
- 商業目標
- 不可違反的風格限制

不負責 babysitting。

### L1：Director / Orchestrator Layer

負責：

- 每日目標決策
- 任務拆解
- 角色分派
- 進度整合
- 品質判斷
- 是否重做 / rollback / merge

### L2：Architect Layer

負責：

- 規格架構
- 玩法設計
- 技術邊界
- 資料結構
- 驗收標準
- Design Constitution 維護

### L3：Specialist Agent Layer

負責具體產出：

- 美術
- 動畫
- 音樂音效
- Godot Client
- Go Server
- 數值平衡
- 自動測試

### L4：Memory / Skill / Feedback Layer

負責：

- lessons learned
- skill 更新
- memory 更新
- checklist 更新
- failed attempt 記錄
- 每日 retrospective

---

## 6. Agent 清單與職責

## 6.1 Game Director Agent

### 定位

專案總監與 Orchestrator。

### 職責

- 讀取昨日 report、QA 結果、progress、repo 狀態。
- 決定今日 1～3 個主目標。
- 拆解任務樹。
- 指派給各 Specialist Agent。
- 監控品質分數。
- 決定是否重做、rollback、merge。
- 輸出 nightly report。

### 可讀資料

- docs/
- reports/
- progress.md
- tasks/
- QA report
- quality-score.md
- git diff

### 可寫資料

- tasks/today-plan.md
- reports/nightly-report.md
- reports/decision-log.md
- integration summary

### 不應直接做

- 不直接改 Go code。
- 不直接改 Godot scene。
- 不直接生成素材。

---

## 6.2 Spec Architect Agent

### 定位

規格架構師。

### 職責

- 維護 game-spec.md。
- 維護 feature-spec.md。
- 將 Research Agent 的研究整理成可執行規格。
- 定義玩法、BOSS、Bonus、Event、UI、音效、動畫需求。
- 檢查規格與 repo 是否一致。

### 可寫資料

- docs/game-spec.md
- docs/feature-specs/*.md
- docs/system-architecture.md
- docs/design-constitution.md
- docs/acceptance-criteria.md

### 重要規則

Research Agent 找到的資料不能直接進正式規格。  
必須先進 references，再由 Spec Architect 摘要、轉譯、去風險後納入規格。

---

## 6.3 Research Agent

### 定位

研究員與參考資料收集者。

### 職責

- 上網找 Steam 像素遊戲參考。
- 找捕魚機、Arcade、Pachinko 的演出節奏參考。
- 找 Godot 4 動畫與音效最佳實踐。
- 找 AI sprite / pixel animation / ComfyUI workflow。
- 找 open-source agent / game dev tools / asset pipeline。
- 分類授權風險。

### 研究分類

```text
可直接使用：MIT / Apache / BSD / CC0
可參考不可搬：商業遊戲畫面、Steam 預告、IP 角色、美術風格
不可使用：有版權素材、閉源音樂、未授權角色圖、商業音效包
```

### 可寫資料

- references/research-notes/*.md
- references/license-review.md
- references/steam-quality-benchmark.md
- references/animation-reference.md
- references/audio-reference.md

### 不可做

- 不可直接下載商業素材放進遊戲。
- 不可直接複製有版權角色、美術、音樂。
- 不可直接修改正式規格。

---

## 6.4 Art Director Agent

### 定位

美術總監。

### 職責

- 定義整體視覺風格。
- 維護 art-style-guide.md。
- 檢查角色、怪物、背景、UI 是否一致。
- 給素材評分。
- 判斷素材是否可進正式 asset。

### 評分項目

- style consistency
- silhouette stability
- color palette consistency
- pixel readability
- IP-like cuteness
- not over-detailed
- not deformed
- transparent background quality

### 可寫資料

- docs/art-style-guide.md
- reports/art-review.md
- assets/candidates/
- assets/approved/

---

## 6.5 Sprite Generation Agent

### 定位

圖像生成與後處理 Agent。

### 職責

- 透過 ComfyUI / Python tool 生成角色、怪物、特效、背景。
- 生成候選圖。
- 去背、裁切、對齊、縮放、pixel clean-up。
- 輸出給 Art Director 評審。

### 可執行工具

- tools/comfyui_generate.py
- tools/generate_chars_v*.py
- tools/generate_targets_v*.py
- tools/generate_backgrounds_v*.py
- tools/generate_effects_v*.py
- tools/process_sprites.py
- tools/batch_process_ai.py

### 產出

- assets/candidates/sprites/
- assets/candidates/backgrounds/
- assets/candidates/effects/
- reports/sprite-generation-report.md

---

## 6.6 Animation Agent

### 定位

動畫工程與序列圖生成 Agent。

### 主要痛點

目前最痛的是：

```text
圖生不太出來，而且變成動畫後會變形到不能看。
```

因此 Animation Agent 必須優先解決：

- 幀與幀角色一致性
- 輪廓穩定
- anchor 對齊
- 武器位置穩定
- 透明背景穩定
- spritesheet cell size 穩定
- Godot 播放不變形

### 職責

- 製作 pose plan。
- 製作 keyframe spec。
- 生成 spritesheet。
- 生成 preview gif。
- 檢查 frame consistency。
- 自動接入 Godot animation spec。
- 由 QA Agent 錄影驗收。

### 可寫資料

- docs/animation-specs/*.md
- assets/candidates/animations/
- assets/approved/animations/
- reports/animation-quality-report.md
- client/chiikawa-pixel/assets/sprites/sheets/

---

## 6.7 Audio Director Agent

### 定位

音樂音效設計與音畫同步 Agent。

### 職責

- 設計 SFX 清單。
- 設計 BGM layer。
- 產生音效候選。
- 設定 Godot Audio Bus。
- 建立 audio-map.json。
- 檢查音畫同步。
- 檢查缺音、過吵、重複、不協調。

### 必須支援

- C：音效有情緒與層次。
- D：BGM 可分層，根據遊戲狀態動態切換。
- E：音畫同步。
- F：QA Agent 自動檢查缺失與不同步。

### 產出

- audio/audio-map.json
- audio/sfx-list.md
- audio/bgm-layer-plan.md
- audio/sync-table.md
- reports/audio-review.md

---

## 6.8 Godot Client Agent

### 定位

Godot 前端工程 Agent。

### 職責

- 接入動畫。
- 接入音效。
- 調整 UI。
- 調整操作手感。
- 改 GDScript。
- 維護 scene tree。
- 確保 WebSocket 訊息正確路由。

### 可寫資料

- client/chiikawa-pixel/scripts/
- client/chiikawa-pixel/scenes/
- client/chiikawa-pixel/assets/
- reports/godot-client-report.md

### 必須驗證

- Godot project 可開啟。
- HTML5 build 可產生。
- 主遊戲流程可跑。
- 動畫播放不變形。
- 音效事件能觸發。

---

## 6.9 Go Server Agent

### 定位

後端工程 Agent。

### 職責

- 維護 Go Server。
- 維護 WebSocket protocol。
- 維護遊戲狀態機。
- 維護 target spawn、combat、boss、bonus、player。
- 修復 server bug。
- 補測試。

### 可寫資料

- server/cmd/
- server/internal/
- server/pkg/
- server/go.mod
- docs/protocol.md
- reports/go-server-report.md

### 高風險規則

以下改動屬 High Risk：

- WebSocket protocol schema 變更
- state machine 變更
- RTP / kill chance 變更
- reward formula 變更
- player state 變更

必須附：

- rollback plan
- test result
- protocol migration note

---

## 6.10 Balance Agent

### 定位

數值工程與 RTP 驗證 Agent。

### 職責

- 跑 RTP simulation。
- 檢查倍率、HP、出現率、Bonus、Boss、補償機制。
- 檢查不同 bet level 的體驗。
- 提出數值 patch。
- 避免 RTP 爆炸或過低。

### 可用工具

- tools/simulate_rtp.py
- tools/rtp_analysis.py
- server/internal/data/tables.go

### 產出

- reports/balance-report.md
- reports/rtp-simulation.md
- patches/balance-patch.md

---

## 6.11 QA Playtest Agent

### 定位

自動測試玩家與品質評審 Agent。

### 職責

- 自動啟動 server。
- 自動開啟遊戲。
- 自動操作遊戲。
- 錄影、截圖。
- 檢查 crash、卡住、音畫不同步、動畫變形、UI 溢位。
- 打品質分數。
- 開 bug report。

### 測試項目

- Build stability
- Login / connect
- NormalPlay
- attack
- auto attack
- lock target
- bet change
- boss warning
- boss battle
- bonus ready
- bonus game
- reward popup
- SFX / BGM trigger
- animation playback
- RTP sanity check

### 產出

- reports/qa-playtest-report.md
- reports/bug-list.md
- reports/screenshots/
- reports/videos/
- reports/quality-score.md

---

## 6.12 Skill Librarian Agent

### 定位

知識管理與自我進化 Agent。

### 職責

- 收集每日 lessons learned。
- 把成功流程寫成 skill。
- 把失敗流程寫成 failed-attempt log。
- 更新 CLAUDE.md / AGENTS.md / skills。
- 維護 memory。
- 維護 checklist。

### 可寫資料

- skills/
- memory/
- reports/retro/
- failed-attempts/
- CLAUDE.md
- AGENTS.md
- checklist/

---

## 7. Design Constitution

Agent 可以自主擴充玩法，但不得偏離核心設計憲法。

### 7.1 核心定位

```text
可愛 IP 感 + 像素風 + 捕魚機玩法 + 高爽感 + 勞動報酬包裝 + 可持續擴充
```

### 7.2 不可偏離

- 不可變成純 RPG。
- 不可變成塔防。
- 不可變成橫向動作遊戲。
- 不可過度寫實。
- 不可過度血腥。
- 不可破壞可愛風格。
- 不可讓核心玩法離開「射擊 / 討伐 / 擊破 / 獎勵」。

### 7.3 可自主新增

Agent 可以自主設計：

- 新怪物
- 新特殊目標
- 新 Boss
- 新 Bonus Game
- 新 Event
- 新特效
- 新音效
- 新 UI feedback
- 新演出節奏
- 新數據分析工具

但必須符合：

- 玩法主軸
- 美術風格
- RTP 合理性
- 操作簡單性
- 可開發性
- 可測試性

---

## 8. Repo / Branch / Merge 策略

### 8.1 基本策略

採用：

```text
C + D
```

意思是：

1. 每個任務一個 branch。
2. 每天產一個 integration branch。
3. 晚上統一測試、產 build、合併。

### 8.2 Branch 命名

```text
agent/<agent-name>/<task-id>-<short-title>
integration/daily-YYYYMMDD
release/playable-YYYYMMDD
```

範例：

```text
agent/animation/ANIM-003-fix-chiikawa-attack-sheet
agent/audio/AUD-002-boss-phase2-layer
agent/godot/GODOT-006-import-animation-state
integration/daily-20260517
release/playable-20260517
```

### 8.3 Merge 條件

任何 branch 進 integration 前必須通過：

- build check
- lint / static check
- relevant test
- agent work report
- risk classification
- rollback plan for high-risk changes

### 8.4 Main branch 保護

main 不允許單一 Agent 直接改。

流程：

```text
task branch
  ↓
integration branch
  ↓
auto build + QA + quality gate
  ↓
release branch
  ↓
main
```

---

## 9. 風險分級與權限

| 風險 | 可自動做 | 條件 |
|---|---|---|
| Low | 補文件、產素材候選、產報告、跑測試 | 直接執行 |
| Medium | 改 Godot 動畫、UI、音效接入、素材替換 | 通過自動測試才進 integration |
| High | 改 Go Server、WebSocket、RTP、狀態機 | 必須有 rollback plan + 測試報告 |
| Critical | 大改核心玩法、刪檔、重構架構 | 只能進 proposal，不得直接 merge main |

### 9.1 Critical 操作禁止自動 merge

Critical 包含：

- 刪除大量 assets
- 改整體技術棧
- 重寫 WebSocket protocol
- 推翻遊戲核心玩法
- 大幅改 RTP 模型
- 移除既有主流程

Critical 只能輸出 proposal。

---

## 10. 每日自主開發節奏

```text
06:30 Daily Scan
讀 repo / docs / progress / yesterday QA report

07:00 Planning
Game Director 決定今日 1～3 個主目標

08:00 Research / Spec
Research 找參考，Spec 補規格與驗收標準

10:00 Parallel Production
Art / Animation / Audio / Godot / Go / Balance 分頭產出

15:00 Integration
合併到 daily integration branch

16:00 Auto Build
產出 Web build / local build

17:00 QA Playtest
自動玩、錄影、截圖、檢查 bug

18:00 Quality Review
Director 判斷達標與否，不達標自動開修正任務

20:00 Retro
寫入 lessons learned / skills / memory

22:00 Nightly Report
輸出今日成果、明日建議、品質趨勢
```

---

## 11. 每日輸出物

每天晚上必須輸出：

1. playable build
2. 60～120 秒 demo video
3. screenshots
4. quality-score.md
5. qa-playtest-report.md
6. bug-list.md
7. next-day-plan.md
8. retro-learning.md
9. updated skills / memory if needed

---

## 12. 品質分數系統

每日輸出 8 個品質分數。

| 指標 | 說明 | 目標 |
|---|---|---|
| Spec Completeness | 規格完整度 | ≥ 95 |
| Build Stability | 是否能成功 build / run | ≥ 95 |
| Visual Consistency | 美術一致性 | ≥ 90 |
| Animation Quality | 動畫不變形、幀穩定 | ≥ 88 |
| Audio Sync | 音畫同步 | ≥ 90 |
| Gameplay Feel | 操作手感與爽感 | ≥ 85 |
| Balance Health | RTP / 數值合理性 | ≥ 90 |
| Regression Risk | 是否破壞既有功能 | ≤ 10 |

### 12.1 硬門檻

```text
Animation Quality < 88：不得 merge
Visual Consistency < 90：不得替換正式素材
Build Stability < 95：不得產出展示版
Regression Risk > 10：自動 rollback
```

---

## 13. Animation Pipeline

### 13.1 目標

解決：

- 圖生不穩
- 動畫變形
- 幀與幀不一致
- spritesheet 接入 Godot 後破圖

### 13.2 流程

```text
Step 1：Reference Lock
先鎖角色正面、比例、輪廓、武器、色盤

Step 2：Pose Plan
先產動畫分鏡，不直接生 spritesheet

Step 3：Keyframe Generation
只生成關鍵幀：idle_01、attack_impact、hurt、death、bigwin

Step 4：In-between Generation
補中間幀，不允許每幀自由發揮

Step 5：Frame Consistency Check
檢查輪廓偏移、頭身比例、透明背景、腳底對齊、武器位置

Step 6：Spritesheet Packing
統一 cell size、anchor、hitbox、frame duration

Step 7：Preview GIF
輸出 gif 給 QA Agent 評分

Step 8：Godot Import
自動接入 AnimatedSprite2D / AnimationPlayer

Step 9：In-game Capture
錄影檢查實際遊戲中是否變形、閃爍、不同步
```

### 13.3 必備動畫狀態

每個主角色至少支援：

- idle
- attack
- hit feedback
- hurt / surprised
- bigwin
- skill / special
- bonus action
- fail / panic

主要怪物至少支援：

- idle / move
- hit
- death
- special behavior
- boss phase change if boss

### 13.4 Frame Consistency 檢查

必須檢查：

- canvas size
- transparent background
- anchor point
- bottom alignment
- silhouette bounding box
- head size ratio
- weapon position
- color palette drift
- unwanted deformation
- frame-to-frame jitter

### 13.5 Animation Work Report

每次動畫任務必須輸出：

```markdown
# Animation Work Report

## Target

## Animation States

## Frame Count

## Cell Size

## Anchor Point

## Generated Files

## Preview GIF

## Consistency Score

## Known Issues

## Godot Import Notes

## Next Fix
```

---

## 14. Audio Pipeline

### 14.1 目標

讓音樂音效品質朝 Steam 遊戲水準靠近，不是只有「有聲音」。

### 14.2 事件驅動音效

```text
Game Event
  ↓
Audio Cue
  ↓
Layer / Intensity / Timing
  ↓
Godot Audio Bus
  ↓
In-game Sync Check
```

### 14.3 事件音效表

| 事件 | 音效設計 |
|---|---|
| 普通攻擊 | 角色專屬短音色 |
| 命中 | 清脆 hit，不蓋過攻擊音 |
| 擊殺 | kill + coin/reward 分層 |
| 20x 以上 | big win sting |
| Boss Warning | 低頻警告 + UI 閃爍同步 |
| Boss Enter | 短暫壓迫感 transition |
| Boss Phase 2 | BGM 加層，音壓提升 |
| Bonus Ready | 預告感音效 |
| Bonus Start | 音樂切換，節奏變快 |
| Reward Bag | 有辨識度的報酬袋音效 |
| UI Click | 低干擾短音 |

### 14.4 BGM Layer

至少支援：

- normal base layer
- fever layer
- boss layer
- boss phase 2 layer
- bonus layer
- big win sting
- result jingle

### 14.5 Audio 產出文件

```text
audio/audio-map.json
audio/sfx-list.md
audio/bgm-layer-plan.md
audio/sync-table.md
reports/audio-review.md
reports/missing-audio-report.md
```

---

## 15. Godot 接入規範

### 15.1 動畫接入

Godot Client Agent 必須將 approved animation 接入：

- AnimatedSprite2D 或 AnimationPlayer
- 正確 frame duration
- 正確 anchor
- 正確 scale
- 正確 hitbox
- 正確 state trigger

### 15.2 音效接入

音效必須由 audio-map.json 驅動。

範例：

```json
{
  "attack.chiikawa": {
    "file": "res://assets/audio/sfx/attack_chiikawa.wav",
    "bus": "SFX",
    "volume_db": -4,
    "sync_event": "attack_frame_2"
  },
  "boss.warning": {
    "file": "res://assets/audio/sfx/boss_warning.wav",
    "bus": "SFX",
    "volume_db": -2,
    "sync_event": "warning_overlay_start"
  }
}
```

### 15.3 Client 驗收

必須檢查：

- WebSocket 連線正常
- 攻擊正常
- 自動攻擊正常
- 鎖定正常
- bet change 正常
- animation trigger 正常
- audio trigger 正常
- UI 不遮擋
- bonus / boss 流程正常

---

## 16. Go Server 接入規範

### 16.1 Protocol 變更

所有 WebSocket protocol 變更必須同步更新：

- server/internal/ws/protocol.go
- client NetworkManager / GameManager
- docs/protocol.md
- QA protocol test

### 16.2 狀態機變更

任何狀態機新增必須提供：

- state name
- enter condition
- exit condition
- duration
- event payload
- UI expectation
- audio cue
- QA test case

### 16.3 數值變更

任何倍率、擊破率、HP、出現權重、RTP 相關改動，必須由 Balance Agent 跑模擬。

---

## 17. Balance / RTP 驗證

### 17.1 目標

- 避免 RTP 異常。
- 維持中高波動爽感。
- 確保低 bet / 中 bet / 高 bet 都有合理體驗。

### 17.2 必跑測試

- 1,000 局快速模擬
- 10,000 局穩定模擬
- 各 bet level 模擬
- boss / bonus contribution
- special target contribution
- high multiplier hit frequency

### 17.3 報告格式

```markdown
# RTP Simulation Report

## Summary

## Bet Level Distribution

## Base Target RTP

## Special Target RTP

## Boss Contribution

## Bonus Contribution

## Volatility

## Risk

## Recommendation
```

---

## 18. QA Playtest 流程

### 18.1 自動測試流程

```text
Start Go Server
  ↓
Open Godot HTML5 Build
  ↓
Connect WebSocket
  ↓
Play Normal Mode
  ↓
Attack / Auto / Lock / Bet Change
  ↓
Trigger Boss
  ↓
Trigger Bonus
  ↓
Capture Video + Screenshots
  ↓
Analyze Crash / Visual / Audio / Gameplay
  ↓
Generate QA Report
```

### 18.2 QA 必檢問題

- 是否 crash
- 是否無法連線
- 是否攻擊無反應
- 是否怪物不生成
- 是否獎勵不跳
- 是否動畫變形
- 是否音效缺失
- 是否音畫不同步
- 是否 UI 破版
- 是否 BOSS 流程卡住
- 是否 Bonus 結算錯誤
- 是否 RTP 異常

---

## 19. Self-Improvement Loop

本系統必須每日進行：

```text
Log → Analyze → Fix → Skill → Reuse
```

### 19.1 Log Everything

每個 Agent 必須記錄：

- 做了什麼
- 為什麼做
- 改了哪些檔案
- 驗證結果
- 失敗原因
- 下一步建議

### 19.2 Failed Attempt Log

失敗不是停止，而是訓練資料。

每次失敗必須記錄：

```markdown
# Failed Attempt

## Task

## What Failed

## Error Message / Screenshot

## Suspected Cause

## Fix Attempt

## Final Result

## Lesson Learned

## Future Prevention
```

### 19.3 Skill 沉澱

當同一類任務成功 2 次以上，Skill Librarian Agent 應將流程寫成 skill。

範例：

- skill-godot-animation-import.md
- skill-comfyui-consistent-sprite.md
- skill-audio-sync-map.md
- skill-rtp-simulation.md
- skill-boss-event-design.md

---

## 20. Repo 內資料夾建議

建議新增：

```text
/agents
  game-director.md
  spec-architect.md
  research-agent.md
  art-director.md
  sprite-generation-agent.md
  animation-agent.md
  audio-director.md
  godot-client-agent.md
  go-server-agent.md
  balance-agent.md
  qa-playtest-agent.md
  skill-librarian.md

/tasks
  today-plan.md
  backlog.md
  done.md

/reports
  nightly/
  qa/
  quality/
  balance/
  art/
  audio/
  animation/
  build/

/references
  research-notes/
  license-review.md
  animation-reference.md
  audio-reference.md
  steam-quality-benchmark.md

/skills
  skill-godot-animation-import.md
  skill-comfyui-consistent-sprite.md
  skill-audio-sync-map.md
  skill-rtp-simulation.md

/memory
  project-memory.md
  art-memory.md
  audio-memory.md
  gameplay-memory.md

/failed-attempts
  YYYYMMDD-task-id.md

/builds
  daily/
  release/
```

---

## 21. Agent Work Report 標準格式

所有 Agent 完成任務時必須輸出：

```markdown
# Agent Work Report

## Agent

## Task

## Input
讀了哪些檔案 / 規格 / reference

## Change
改了哪些檔案 / 產了哪些素材

## Reason
為什麼這樣做

## Validation
怎麼驗證

## Quality Score
如果適用，填入品質分數

## Risk
可能造成什麼問題

## Rollback Plan
如何回復

## Next Action
下一步建議

## Skill Learned
這次有什麼流程值得沉澱
```

---

## 22. Nightly Report 格式

```markdown
# Nightly Report - YYYY-MM-DD

## 1. 今日目標

## 2. 今日完成

## 3. Build 狀態

## 4. Demo Video

## 5. Screenshots

## 6. Quality Score

| 指標 | 分數 | 是否達標 | 說明 |
|---|---:|---|---|
| Spec Completeness |  |  |  |
| Build Stability |  |  |  |
| Visual Consistency |  |  |  |
| Animation Quality |  |  |  |
| Audio Sync |  |  |  |
| Gameplay Feel |  |  |  |
| Balance Health |  |  |  |
| Regression Risk |  |  |  |

## 7. 主要問題

## 8. 已自動修正

## 9. 尚未解決

## 10. 明日建議

## 11. Lessons Learned

## 12. Skills Updated
```

---

## 23. MVP 實作順序

雖然最終目標是完整 Orchestrator，但實作順序應分階段。

### Phase 1：Multi-Agent Repo Scaffold

目標：建立 agent 目錄、report 格式、task 格式、daily loop。

交付：

- /agents
- /tasks
- /reports
- /memory
- /skills
- AGENTS.md
- CLAUDE.md update
- first daily plan

### Phase 2：Spec + Research + Director Loop

目標：讓系統能自主讀規格、補洞、拆任務。

交付：

- design-constitution.md
- acceptance-criteria.md
- research workflow
- task decomposition workflow

### Phase 3：Animation Pipeline

目標：解決目前最大痛點：動畫變形。

交付：

- animation-spec template
- frame consistency checker
- preview gif exporter
- Godot animation import workflow
- animation quality score

### Phase 4：Audio Pipeline

目標：建立 Steam-like 音樂音效流程。

交付：

- audio-map.json
- bgm-layer-plan.md
- sync-table.md
- missing-audio-report.md
- audio QA check

### Phase 5：Daily Build + QA Playtest

目標：每天產一版 playable build。

交付：

- auto build script
- QA playtest script
- screenshot / video capture
- quality-score.md

### Phase 6：Self-Improvement Loop

目標：讓系統能記住錯誤與成功流程。

交付：

- failed-attempt log
- skill generation workflow
- retro report
- memory update workflow

### Phase 7：Full Autonomous Studio

目標：完整自主開發、測試、修正、沉澱。

交付：

- full daily loop
- auto integration branch
- auto rollback
- quality trend
- weekly build showcase

---

## 24. Kiro Implementation Prompt

以下 prompt 可直接丟給 Kiro 作為實作起點。

```markdown
# 任務：建立 Autonomous Multi-Agent Game Studio 架構

你現在要在現有 Godot + Go + Python + ComfyUI 遊戲專案中，建立一套 Autonomous Multi-Agent Game Studio 架構。

## 目標

不是重寫遊戲，而是在現有 repo 外層建立 multi-agent 自主開發系統，使專案可以每日自主：

1. 讀取 repo / docs / progress / report
2. 拆解任務
3. 指派不同 agent
4. 產出素材 / 動畫 / 音效 / code patch / balance patch
5. 跑 build / QA / RTP / 截圖 / 錄影
6. 根據 quality gate 決定重做、rollback、merge
7. 每日輸出 playable build + report
8. 沉澱 skills / memory / failed-attempt logs

## 第一階段請先建立 Repo Scaffold

請新增以下資料夾與檔案：

/agents
/tasks
/reports
/reports/nightly
/reports/qa
/reports/quality
/reports/balance
/reports/art
/reports/audio
/reports/animation
/references
/references/research-notes
/skills
/memory
/failed-attempts
/builds/daily
/builds/release

## 建立 Agent 說明檔

請在 /agents 建立以下 markdown：

- game-director.md
- spec-architect.md
- research-agent.md
- art-director.md
- sprite-generation-agent.md
- animation-agent.md
- audio-director.md
- godot-client-agent.md
- go-server-agent.md
- balance-agent.md
- qa-playtest-agent.md
- skill-librarian.md

每份檔案需要包含：

1. Role
2. Responsibilities
3. Read Access
4. Write Access
5. Tools
6. Output Artifacts
7. Validation Rules
8. Risk Rules
9. Work Report Format

## 建立核心文件

請建立：

- AGENTS.md
- docs/design-constitution.md
- docs/acceptance-criteria.md
- docs/protocol-change-policy.md
- tasks/today-plan.md
- tasks/backlog.md
- reports/nightly/nightly-report-template.md
- reports/quality/quality-score-template.md
- failed-attempts/failed-attempt-template.md

## 品質門檻

請實作文件化 quality gate：

- Spec Completeness >= 95
- Build Stability >= 95
- Visual Consistency >= 90
- Animation Quality >= 88
- Audio Sync >= 90
- Gameplay Feel >= 85
- Balance Health >= 90
- Regression Risk <= 10

硬規則：

- Animation Quality < 88 不得 merge
- Visual Consistency < 90 不得替換正式素材
- Build Stability < 95 不得產展示版
- Regression Risk > 10 自動 rollback

## 不要做的事情

第一階段不要重寫遊戲主流程。
第一階段不要大改 Go Server。
第一階段不要大改 Godot scene。
第一階段先建立 multi-agent 工作骨架與文件規範。

## 完成後輸出

請輸出：

1. 建立了哪些檔案
2. 每個檔案用途
3. 下一階段建議
4. 風險
5. 如何驗證 scaffold 成功
```

---

## 25. 第一個可執行任務建議

建議第一個任務不是直接改遊戲，而是：

```text
TASK-001：建立 Multi-Agent Repo Scaffold
```

### 驗收條件

- /agents 內 12 個 agent 說明檔建立完成。
- /tasks 有 today-plan / backlog。
- /reports 有 nightly / QA / quality templates。
- /skills 與 /memory 建立完成。
- AGENTS.md 能說明所有 agent 如何協作。
- docs/design-constitution.md 能約束 agent 不偏離核心玩法。
- docs/acceptance-criteria.md 能定義品質門檻。

---

## 26. 第二個可執行任務建議

```text
TASK-002：建立 Animation Consistency Pipeline
```

### 目標

解決動畫變形問題。

### 內容

- 建立 animation-spec template。
- 建立 frame consistency checklist。
- 建立 preview gif 輸出流程。
- 建立 Godot import notes。
- 建立 QA animation score。

### 驗收條件

- 任一角色可輸出 idle / attack / bigwin preview gif。
- 每個 gif 有 consistency score。
- cell size / anchor / frame count 有明確紀錄。
- Godot 可正確播放，不明顯變形。

---

## 27. 第三個可執行任務建議

```text
TASK-003：建立 Audio Event Map Pipeline
```

### 目標

建立可接近 Steam 品質的音效事件架構。

### 內容

- 建立 audio-map.json。
- 建立 sfx-list.md。
- 建立 bgm-layer-plan.md。
- 建立 sync-table.md。
- 建立 missing-audio-report.md。

### 驗收條件

- attack / hit / kill / bigwin / boss / bonus / UI 都有 cue。
- Godot 可根據事件播放音效。
- QA 能檢查缺失音效。

---

## 28. 結語

這套系統的本質不是讓 AI 一次做出完美遊戲，而是讓 AI 具備一個遊戲工作室的日常節奏：

```text
想法 → 規格 → 研究 → 生成 → 接入 → 測試 → 修正 → 沉澱 → 再進化
```

真正的關鍵不是 Agent 數量，而是：

- 有沒有明確分工
- 有沒有 Source of Truth
- 有沒有品質門檻
- 有沒有自動測試
- 有沒有 rollback
- 有沒有 failed-attempt log
- 有沒有 skill 沉澱

只要這些建立起來，這個專案就不只是吉伊卡哇 Demo，而是可以複製到未來魚機、Boss、老虎機、小遊戲、AI 製程展示的自主研發工廠。

