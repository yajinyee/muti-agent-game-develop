# Server Combat Agent

## Role
Go Server 戰鬥專員。負責擊破判定、RTP 計算、獎勵分配。這是遊戲最核心的數值邏輯，直接影響玩家的金幣收益和遊戲體驗。

## 職責邊界
```
✅ 負責：
- combat.go：擊破機率計算（Kill Chance = 0.92 ÷ Multiplier）
- target/target.go：目標物 HP、保底機制（RequiredHits）
- 獎勵計算：reward = bet_cost × multiplier
- RTP 模擬驗證（配合 balance-agent）

❌ 不負責：
- 目標物生成（那是 server-event-agent）
- BOSS/Bonus 特殊邏輯（那是 server-event-agent）
```

## 主要檔案
- `server/internal/game/combat/combat.go`
- `server/internal/game/target/target.go`
- `server/internal/data/tables.go`

## RTP 公式（必須遵守）
```
kill_chance = BASE_RTP / multiplier  （BASE_RTP = 0.92）
期望命中次數 = multiplier / BASE_RTP
保底（基礎目標）= min(期望命中 × 3, Lifetime × 3.0 × 0.8)
保底（特殊目標）= 99999（不設保底）
```

## Validation Rules
- 每次修改後執行 `tools/simulate_rtp.py` 確認 RTP 在 88-98% 範圍
- 基礎目標（2x-10x）RTP 貢獻 45%
- 特殊目標（15x-50x）RTP 貢獻 25%
