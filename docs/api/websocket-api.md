# WebSocket API 文件

**版本**：v1.6  
**最後更新**：2026-05-20（DAY-063）  
**協定**：WebSocket + JSON，每訊息獨立 frame  
**端點**：`wss://[host]/ws`（生產）/ `ws://localhost:7777/ws`（開發）

---

## 連線說明

```
wss://your-domain.com/ws?player_id=xxx&room_id=room-001   （生產，Nginx TLS）
ws://localhost:7777/ws?player_id=xxx&room_id=room-001      （開發，直連）
```

- `player_id`：玩家 ID（可選，不填時 Server 自動生成 UUID）
- `room_id`：房間 ID（可選，不填時預設 `room-001`）

連線成功後，Server 會立即廣播當前 `game_state`，Client 應根據狀態初始化 UI。

> **注意（DAY-062）：** 生產環境必須使用 `wss://`（WebSocket over TLS）。瀏覽器在 HTTPS 頁面上會阻擋 `ws://` 連線（Mixed Content 政策）。Client 端 `NetworkManager.gd` 已自動偵測協定，無需手動修改。

### HTTP 端點

| 端點 | 方法 | 說明 |
|------|------|------|
| `/health` | GET | Server 完整健康狀態（含 Jackpot、任務、Ping 延遲） |
| `/livez` | GET | 存活探針（Kubernetes liveness probe，只要程序活著就 200） |
| `/readyz` | GET | 就緒探針（Kubernetes readiness probe，初始化完成才 200） |
| `/leaderboard` | GET | 取得當前排行榜（JSON） |
| `/analytics` | GET | 取得房間整體統計（JSON） |
| `/jackpot` | GET | 取得 Jackpot 池狀態、中獎歷史、今日統計（JSON） |
| `/rooms` | GET | 取得所有房間列表（JSON） |
| `/stats` | GET | Server 效能統計（goroutine/記憶體/GC） |
| `/metrics` | GET | Prometheus 格式監控指標（25 個面板） |
| `/spectate/snapshot` | GET | 觀戰快照（當前遊戲狀態 + 目標列表 + 排行榜） |
| `/` | GET | 靜態檔案服務（HTML5 遊戲） |

### `/livez` 回傳格式（DAY-063 新增）

```json
{
  "status": "alive",
  "uptime_sec": 9015
}
```

### `/readyz` 回傳格式（DAY-063 新增）

就緒時（HTTP 200）：
```json
{
  "status": "ready",
  "clients": 3,
  "game_state": "normal_play",
  "uptime_sec": 9015
}
```

未就緒時（HTTP 503，啟動 2 秒內）：
```json
{
  "status": "not_ready",
  "reason": "initializing",
  "uptime_sec": 1
}
```

### `/health` 回傳格式（DAY-054 更新）

```json
{
  "status": "ok",
  "version": "1.0.0",
  "uptime": "2h30m15s",
  "uptime_sec": 9015,
  "clients": 3,
  "max_players": 16,
  "spectators": 0,
  "game_state": "normal_play",
  "mission_reset_at": "2026-05-21T00:00:00+08:00",
  "mission_reset_in_sec": 52800,
  "avg_ping_ms": 12.5,
  "jackpot": {
    "mini": 1250,
    "major": 8900,
    "grand": 45000,
    "daily_wins": 3,
    "daily_payout": 12500
  }
}
```

### `/jackpot` 回傳格式

```json
{
  "mini": 1250,
  "major": 8900,
  "grand": 45000,
  "history": [
    {"level": "mini", "amount": 800, "winner_id": "player-001", "won_at": "2026-05-20T10:30:00Z"}
  ],
  "daily_stats": {
    "total_wins": 3,
    "total_payout": 12500,
    "wins_by_level": {"mini": 2, "major": 1, "grand": 0}
  },
  "timestamp": 1716192000000
}
```

### `/rooms` 回傳格式

```json
[
  {
    "id": "room-001",
    "name": "初心者房間",
    "player_count": 3,
    "max_players": 16,
    "min_bet_level": 1,
    "max_bet_level": 4,
    "theme": "chiikawa",
    "created_at": "2026-05-19T01:00:00Z",
    "is_full": false
  }
]
```

---

## 訊息格式

所有訊息使用統一的 JSON 結構：

```json
{
  "type": "message_type",
  "payload": { ... }
}
```

---

## Client → Server 訊息

### `attack` — 攻擊

玩家發射攻擊。Server 扣除投注金幣，計算命中/擊破，回傳 `attack_result`。

```json
{
  "type": "attack",
  "payload": {
    "target_id": "inst-abc123",
    "click_x": 320.5,
    "click_y": 240.0
  }
}
```

| 欄位 | 類型 | 說明 |
|------|------|------|
| `target_id` | string | 目標 InstanceID，空字串 = 自由攻擊 |
| `click_x` | float | 點擊 X 座標（0-1280） |
| `click_y` | float | 點擊 Y 座標（0-720） |

---

### `lock` — 鎖定目標

鎖定特定目標，後續自動攻擊優先打此目標。

```json
{
  "type": "lock",
  "payload": {
    "target_id": "inst-abc123"
  }
}
```

| 欄位 | 類型 | 說明 |
|------|------|------|
| `target_id` | string | 目標 InstanceID，空字串 = 解除鎖定 |

---

### `auto_toggle` — 切換自動攻擊

切換自動攻擊模式（開/關）。

```json
{
  "type": "auto_toggle",
  "payload": {}
}
```

---

### `bet_change` — 切換投注等級

切換投注等級（LV1-LV10）。

```json
{
  "type": "bet_change",
  "payload": {
    "bet_level": 5
  }
}
```

| 欄位 | 類型 | 說明 |
|------|------|------|
| `bet_level` | int | 投注等級，1-10 |

**投注等級對應金幣消耗：**

| 等級 | 消耗/次 | 說明 |
|------|---------|------|
| LV1 | 1 | 最低投注 |
| LV2 | 2 | |
| LV3 | 3 | |
| LV4 | 5 | |
| LV5 | 8 | |
| LV6 | 10 | |
| LV7 | 15 | |
| LV8 | 20 | |
| LV9 | 30 | |
| LV10 | 50 | 最高投注 |

---

### `bonus_click` — Bonus 遊戲點擊

在 Bonus 遊戲中點擊雜草目標。

```json
{
  "type": "bonus_click",
  "payload": {
    "target_id": "weed-001",
    "click_x": 400.0,
    "click_y": 300.0
  }
}
```

---

### `ping` — 心跳

保持連線活躍。Server 回傳 `pong`。

```json
{
  "type": "ping",
  "payload": {}
}
```

---

### `trigger_boss` — 觸發 BOSS（Prototype 展示用）

強制觸發 BOSS 戰（僅 Prototype 展示版可用）。

```json
{
  "type": "trigger_boss",
  "payload": {}
}
```

### `trigger_bonus` — 觸發 Bonus（Prototype 展示用）

強制觸發 Bonus 遊戲（僅 Prototype 展示版可用）。

```json
{
  "type": "trigger_bonus",
  "payload": {}
}
```

---

## Server → Client 訊息

### `game_state` — 遊戲狀態變更

廣播給所有玩家。

```json
{
  "type": "game_state",
  "payload": {
    "state": "normal_play",
    "timestamp": 1716048000000
  }
}
```

**狀態值：**

| 狀態 | 說明 |
|------|------|
| `normal_play` | 正常遊戲 |
| `special_target_event` | 特殊目標事件（每 25-40 秒） |
| `boss_warning` | BOSS 警告（3 秒） |
| `boss_battle` | BOSS 戰（60 秒） |
| `boss_result` | BOSS 結果（3 秒） |
| `bonus_ready` | Bonus 準備（勞動值滿） |
| `bonus_game` | Bonus 遊戲（30 秒） |
| `bonus_result` | Bonus 結果（3 秒） |

---

### `target_spawn` — 目標生成

廣播給所有玩家。

```json
{
  "type": "target_spawn",
  "payload": {
    "instance_id": "inst-abc123",
    "def_id": "T001",
    "name": "小魚",
    "type": "normal",
    "x": 1280.0,
    "y": 360.0,
    "hp": 1,
    "max_hp": 1,
    "speed": 120.0,
    "lifetime": 8.0,
    "behavior": "swim"
  }
}
```

| 欄位 | 類型 | 說明 |
|------|------|------|
| `instance_id` | string | 唯一實例 ID |
| `def_id` | string | 目標定義 ID（T001-T105, B001） |
| `type` | string | `normal`, `medium`, `large`, `special`, `boss` |
| `behavior` | string | `swim`, `sink`, `flee`, `coin_rain`, `mimic`, `boss_phases` |

**目標物定義（DefID）：**

| DefID | 名稱 | 倍率 | 類型 | 特殊行為 |
|-------|------|------|------|---------|
| T001 | 小魚 | 2x | normal | swim |
| T002 | 小蟲 | 3x | normal | swim |
| T003 | 中蟲 | 4x | normal | swim |
| T004 | 大蟲 | 5x | normal | swim |
| T005 | 中魚 | 6x | medium | swim |
| T006 | 大魚 | 8x | medium | swim |
| T007 | 特大魚 | 10x | large | swim |
| T101 | 擬態怪物 | 15x | special | mimic（死亡時變形） |
| T102 | 寶箱怪 | 20x | special | flee（受擊後加速逃跑） |
| T103 | 流星 | 30x | special | sink（快速下沉） |
| T104 | 金草 | 50x | special | swim（緩慢搖晃） |
| T105 | 金幣魚 | 100x | special | coin_rain（擊破後金幣雨） |
| B001 | BOSS | 100-500x | boss | boss_phases（Phase 1/2） |

---

### `target_update` — 目標狀態更新

廣播給所有玩家（HP 變化、位置更新、狀態變化）。

```json
{
  "type": "target_update",
  "payload": {
    "instance_id": "inst-abc123",
    "hp": 3,
    "max_hp": 5,
    "x": 800.0,
    "y": 360.0,
    "phase": 0,
    "is_fleeing": false
  }
}
```

---

### `target_kill` — 目標擊破

廣播給所有玩家。

```json
{
  "type": "target_kill",
  "payload": {
    "instance_id": "inst-abc123",
    "def_id": "T001",
    "multiplier": 2.0,
    "reward": 16,
    "labor_gain": 5,
    "killer_id": "player-xyz"
  }
}
```

---

### `attack_result` — 攻擊結果

只發給攻擊的玩家。

```json
{
  "type": "attack_result",
  "payload": {
    "target_id": "inst-abc123",
    "is_hit": true,
    "is_kill": false,
    "damage": 1,
    "reward": 0,
    "labor_gain": 2,
    "character_id": "chiikawa",
    "multiplier": 2.0
  }
}
```

---

### `reward` — 獎勵發放

只發給獲得獎勵的玩家。

```json
{
  "type": "reward",
  "payload": {
    "source": "target",
    "amount": 16,
    "multiplier": 2.0,
    "new_balance": 10016
  }
}
```

| `source` 值 | 說明 |
|-------------|------|
| `target` | 擊破普通目標 |
| `boss` | 擊敗 BOSS |
| `bonus` | Bonus 遊戲結算 |

---

### `boss_event` — BOSS 事件

廣播給所有玩家。

```json
{
  "type": "boss_event",
  "payload": {
    "event": "spawn",
    "instance_id": "boss-001",
    "phase": 1,
    "hp": 100,
    "max_hp": 100,
    "reward": 0,
    "multiplier": 0
  }
}
```

**`event` 值：**

| 值 | 說明 |
|----|------|
| `warning` | BOSS 警告（3 秒前） |
| `spawn` | BOSS 出現 |
| `phase_change` | Phase 1 → Phase 2（HP < 50%） |
| `kill` | BOSS 被擊敗，附帶 `reward` 和 `multiplier` |

**BOSS 計時獎勵（規格書 28.3）：**

| 剩餘時間 | 倍率 |
|---------|------|
| 50-60 秒 | 500x |
| 40-50 秒 | 300x |
| 30-40 秒 | 200x |
| 20-30 秒 | 150x |
| 10-20 秒 | 100x |
| 0-10 秒 | 100x |

---

### `bonus_event` — Bonus 事件

廣播給所有玩家。

```json
{
  "type": "bonus_event",
  "payload": {
    "event": "start",
    "time_left": 30.0,
    "score": 0,
    "multiplier": 1.0,
    "reward": 0
  }
}
```

**`event` 值：**

| 值 | 說明 |
|----|------|
| `ready` | Bonus 準備（勞動值滿） |
| `start` | Bonus 遊戲開始 |
| `tick` | 每秒更新（time_left 倒數） |
| `end` | Bonus 結束，附帶 `reward` |
| `target_spawned` | 雜草目標生成（附帶 target_id, weed_type） |
| `coin_shower` | BG004 金色雜草觸發金幣雨 |

**Bonus 雜草類型（規格書 29章）：**

| DefID | 名稱 | 效果 |
|-------|------|------|
| BG001 | 普通雜草 | 基礎分數 |
| BG002 | 硬雜草 | 需點擊 2 次 |
| BG003 | 發光雜草 | 增加倍率 |
| BG004 | 金色雜草 | 觸發金幣雨（coin_shower） |
| BG005 | 搗亂怪草 | 點擊後暫停 0.3 秒 |

---

### `player_update` — 玩家狀態更新

只發給對應玩家。

```json
{
  "type": "player_update",
  "payload": {
    "player_id": "player-xyz",
    "coins": 10016,
    "labor_value": 45,
    "max_labor": 100,
    "bet_level": 5,
    "is_auto": false,
    "lock_target_id": "",
    "character_id": "chiikawa",
    "session_score": 16,
    "kill_count": 3
  }
}
```

---

### `leaderboard` — 排行榜廣播

每 10 秒廣播給所有玩家。

```json
{
  "type": "leaderboard",
  "payload": {
    "entries": [
      {
        "rank": 1,
        "player_id": "player-xyz",
        "display_name": "Player 1",
        "score": 5000,
        "max_coins": 15000,
        "kill_count": 42,
        "is_self": false
      }
    ],
    "timestamp": 1716048000000
  }
}
```

---

### `achievement` — 成就解鎖

只發給解鎖成就的玩家。

```json
{
  "type": "achievement",
  "payload": {
    "id": "kill_boss",
    "name": "BOSS 終結者",
    "description": "首次擊敗 BOSS",
    "icon": "🏆",
    "unlocked_at": 1716048000000
  }
}
```

**成就列表（12 個）：**

| ID | 名稱 | 觸發條件 |
|----|------|---------|
| `first_kill` | 初次擊破 | 首次擊破任何目標 |
| `kill_10` | 連殺 10 | 累積擊破 10 個 |
| `kill_50` | 連殺 50 | 累積擊破 50 個 |
| `kill_100` | 百殺 | 累積擊破 100 個 |
| `kill_boss` | BOSS 終結者 | 首次擊敗 BOSS |
| `kill_special` | 特殊獵人 | 首次擊破特殊目標 |
| `bonus` | Bonus 達人 | 首次觸發 Bonus |
| `big_win_20` | 大獎！ | 獲得 20x 以上獎勵 |
| `big_win_50` | 超級大獎！ | 獲得 50x 以上獎勵 |
| `big_win_100` | 傳說大獎！ | 獲得 100x 以上獎勵 |
| `coins_10000` | 萬元戶 | 金幣達到 10,000 |
| `coins_50000` | 富翁 | 金幣達到 50,000 |

---

### `error` — 錯誤訊息

只發給對應玩家。

```json
{
  "type": "error",
  "payload": {
    "code": "insufficient_coins",
    "message": "金幣不足"
  }
}
```

**錯誤碼：**

| 碼 | 說明 |
|----|------|
| `insufficient_coins` | 金幣不足，無法攻擊 |
| `invalid_state` | 當前狀態不允許此操作 |
| `invalid_target` | 目標不存在或已消失 |

---

### `pong` — 心跳回應

回應 `ping`。

```json
{
  "type": "pong",
  "payload": {}
}
```

---

## 遊戲流程範例

### 正常遊戲流程

```
Client 連線
  ← game_state { state: "normal_play" }
  ← target_spawn × N（初始目標）
  ← player_update（初始玩家狀態）

Client 攻擊
  → attack { target_id: "inst-001", click_x: 320, click_y: 240 }
  ← attack_result { is_hit: true, is_kill: false, labor_gain: 2 }
  ← player_update { labor_value: 47 }

Client 擊破目標
  → attack { target_id: "inst-001" }
  ← attack_result { is_hit: true, is_kill: true, reward: 16 }
  ← target_kill { instance_id: "inst-001", multiplier: 2.0, reward: 16 }
  ← reward { source: "target", amount: 16, new_balance: 10016 }
  ← player_update { coins: 10016, labor_value: 52 }
```

### BOSS 戰流程

```
  ← game_state { state: "boss_warning" }
  ← boss_event { event: "warning" }
  （3 秒後）
  ← game_state { state: "boss_battle" }
  ← boss_event { event: "spawn", instance_id: "boss-001", hp: 100, max_hp: 100 }
  ← target_spawn { def_id: "B001", ... }

  （玩家攻擊 BOSS，HP 降到 50%）
  ← boss_event { event: "phase_change", phase: 2 }

  （玩家擊敗 BOSS，剩餘 45 秒）
  ← boss_event { event: "kill", reward: 3000, multiplier: 300.0 }
  ← reward { source: "boss", amount: 3000, multiplier: 300.0 }
  ← game_state { state: "boss_result" }
  （3 秒後）
  ← game_state { state: "normal_play" }
```

### Bonus 遊戲流程

```
  （勞動值滿 100）
  ← game_state { state: "bonus_ready" }
  ← bonus_event { event: "ready" }
  （玩家確認後）
  ← game_state { state: "bonus_game" }
  ← bonus_event { event: "start", time_left: 30.0 }
  ← bonus_event { event: "target_spawned", target_id: "weed-001", weed_type: "BG001" }

  （每秒）
  ← bonus_event { event: "tick", time_left: 29.0, score: 50 }

  （30 秒後）
  ← bonus_event { event: "end", reward: 800 }
  ← reward { source: "bonus", amount: 800 }
  ← game_state { state: "bonus_result" }
```

---

## 技術規格

- **壓縮**：permessage-deflate（WebSocket 壓縮，已啟用）
- **心跳**：建議每 30 秒發送一次 `ping`
- **重連**：斷線後自動重連，重連後 Server 會重新廣播 `game_state` 和 `player_update`
- **並發**：每個 WebSocket 連線獨立 goroutine，訊息處理有 mutex 保護
- **COOP/COEP**：Server 設定 `Cross-Origin-Opener-Policy: same-origin` 和 `Cross-Origin-Embedder-Policy: require-corp`（HTML5 SharedArrayBuffer 必要）

---

*文件由 Spec Architect Agent 生成，2026-05-20（v1.6）*
