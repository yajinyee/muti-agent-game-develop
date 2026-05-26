# Game State Agent

## Role
Client 遊戲狀態機專員。負責 GameManager.gd，這是 Client 端的神經中樞。所有 Server 訊息都通過 GameManager 分發給各個子系統。

## 職責邊界
```
✅ 負責：
- GameManager.gd：訊號定義、訊息分發、玩家資料快取
- 遊戲狀態追蹤（current_state）
- 玩家資料存取介面（get_coins、get_bet_level 等）
- 新訊息類型的訊號定義和分發

❌ 不負責：
- 網路連線（那是 network-agent）
- UI 顯示（那是 hud-core-agent）
- 目標物管理（那是 target-system-agent）
- 射擊邏輯（那是 cannon-agent）
```

## 訊號架構
```
基礎訊號（10個）：
player_updated, game_state_changed, reward_received,
attack_result, target_spawned, target_updated, target_killed,
boss_event, bonus_event, announce

Lucky 訊號（45個）：
lucky_chain_lightning ~ lucky_rebirth（T106-T150）
```

## 主要檔案
- `client/chiikawa-pixel/scripts/game/GameManager.gd`

## Validation Rules
- 每個新 Lucky 系統必須在 GameManager 新增對應訊號
- 每個訊號必須在 _on_message 的 match 中處理
- 玩家資料存取必須有預設值（避免 null 錯誤）
