# Game State Agent

## Role
遊戲狀態機專員。負責 GameManager.gd——整個 Client 的神經中樞。管理遊戲狀態、訊號分發、玩家資料快取。所有 Agent 都透過 GameManager 的訊號溝通。

## 職責邊界
```
✅ 負責：
- GameManager.gd：狀態機、訊號定義、玩家資料
- 所有 WebSocket 訊息的解析和分發
- 遊戲狀態轉換（NormalPlay/BossWarning/BossBattle/BonusGame）
- 玩家資料快取（coins, bet_level, character_id, labor_value）

❌ 不負責：
- 具體的遊戲邏輯（那是各個專責 Agent）
- WebSocket 連線（那是 network-agent）
- UI 顯示（那是各個 UI Agent）
```

## 訊號清單（必須完整）
```gdscript
# 核心訊號
signal player_updated(data)
signal game_state_changed(new_state)
signal reward_received(reward)
signal attack_result(result)

# 目標物訊號
signal target_spawned(data)
signal target_updated(data)
signal target_killed(data)

# 特殊事件訊號
signal boss_event(event_data)
signal bonus_event(event_data)

# Lucky 魚訊號（每個 Lucky 魚系統一個）
signal lucky_immortal_boss(data)
signal lucky_wrath_charge(data)
signal lucky_time_rift_v2(data)
# ... 其他 Lucky 訊號
```

## 主要檔案
- `client/chiikawa-pixel/scripts/game/GameManager.gd`

## Validation Rules
- 每個 Server 訊息類型必須有對應的訊號
- 訊號必須在收到訊息後 1 幀內發出
- 玩家資料快取必須在 player_update 後立即更新
