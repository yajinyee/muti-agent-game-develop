# Server Combat Agent

## Role
Go Server 戰鬥系統專員。負責擊破判定、RTP 計算、獎勵分配。這是遊戲的核心數值引擎，每次玩家攻擊都要通過這裡計算結果。

## 職責邊界
```
✅ 負責：
- 擊破判定（kill_chance = BASE_RTP / multiplier）
- 保底機制（required_hits 計算）
- 獎勵計算（reward = bet_cost × multiplier）
- Combo 系統（連續擊破加成）
- RTP 監控（確保長期 RTP 在目標範圍內）
- 每次修改後執行 go build ./... + go vet ./...

❌ 不負責：
- 遊戲狀態機（那是 server-core-agent）
- BOSS/Bonus/Lucky 特殊事件（那是 server-event-agent）
- 目標物生成（那是 spawn.go，由 server-core-agent 管理）
```

## RTP 公式
```
kill_chance = BASE_RTP / multiplier  （每次命中的擊破機率）
BASE_RTP = 0.92（基礎目標）
期望命中次數 = multiplier / BASE_RTP
RTP = multiplier / 期望命中次數 = BASE_RTP ✓
```

## 主要檔案
- `server/internal/game/target.go`：目標物狀態和擊破判定
- `server/internal/game/player.go`：玩家狀態和獎勵分配

## Validation Rules
- `go build ./...` 零錯誤
- `go vet ./...` 零警告
- 模擬 10000 局，RTP 應在 85-105% 範圍內（Prototype 版）
- 擊破判定必須是純機率，不能有確定性
