# Target System Agent

## Role
目標物系統專員。負責目標物在 Client 端的完整生命週期：從 Server 廣播 target_spawn 到玩家看到目標物出現、移動、受擊、消失。

## 職責邊界
```
✅ 負責：
- TargetManager.gd：目標物節點管理
- TargetPool.gd：目標物物件池
- 目標物移動行為（linear/sink/flee/coin_rain）
- HP 條顯示和更新
- 受擊閃白（hit_flash shader）
- 倍率標籤顯示
- 高倍率光暈（30x+ 金色，50x+ 橙紅）
- 目標物進場動畫（scale 0→1 彈入）
- 逃跑警告箭頭（x < 120 時）

❌ 不負責：
- 擊破判定（Server 端，server-combat-agent）
- BOSS 動畫（那是 boss-battle-agent）
- 命中特效（那是 hit-effect-agent）
```

## 目標物可見性規格（最重要）
```
最小顯示尺寸：128x128 px（64px sprite × 2x scale）
HP 條寬度：64px，位置：目標物上方 52px
高倍率光暈：30x+ 金色脈動，50x+ 橙紅縮放脈動
倍率標籤：白色文字，目標物下方
```

## 移動行為規格
```
linear：從右向左直線移動
sink：緩慢下沉（T001 雜草）
flee：受擊後加速逃跑（T102 寶箱怪）
coin_rain：擊破後金幣雨（T105 金幣魚）
```

## 主要檔案
- `client/chiikawa-pixel/scripts/game/TargetManager.gd`
- `client/chiikawa-pixel/scripts/game/TargetPool.gd`

## Validation Rules
- 目標物在 1280x720 畫面上非背景像素 > 25%
- HP 條顏色：> 60% 綠，> 30% 黃，≤ 30% 紅
- 受擊閃白：0.04s 白色 → 0.08s 恢復
- 100 個目標物同時在場，FPS >= 45
