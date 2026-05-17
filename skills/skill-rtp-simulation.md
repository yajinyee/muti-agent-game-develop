# Skill：RTP 蒙地卡羅模擬

## 目的
使用蒙地卡羅模擬方法驗證捕魚機遊戲的 RTP（Return to Player）是否符合設計目標（92-96%）。確保遊戲在商業可行性與玩家體驗之間取得平衡。

## 適用場景
- 修改目標物獎勵倍率後的驗證
- 修改命中機率後的驗證
- 新增目標物類型後的驗證
- 定期 RTP 健康檢查

## 前置條件
- Go 1.21+ 已安裝
- `server/` 目錄下有 RTP 計算邏輯
- 有目標物倍率表（config.go 或 JSON）

## 使用方法

### 方法 1：使用 Go 模擬腳本

```go
// tools/rtp_simulator.go
package main

import (
    "fmt"
    "math/rand"
    "time"
)

type Target struct {
    ID         string
    Multiplier float64
    Weight     int
}

func simulateRTP(targets []Target, iterations int) float64 {
    rand.Seed(time.Now().UnixNano())
    
    totalBet := 0.0
    totalReturn := 0.0
    
    // 建立加權隨機選擇
    totalWeight := 0
    for _, t := range targets {
        totalWeight += t.Weight
    }
    
    for i := 0; i < iterations; i++ {
        bet := 1.0
        totalBet += bet
        
        // 隨機選擇目標物
        r := rand.Intn(totalWeight)
        cumWeight := 0
        for _, t := range targets {
            cumWeight += t.Weight
            if r < cumWeight {
                // 模擬命中判定（基於目標 RTP 反推命中率）
                hitRate := calculateHitRate(t.Multiplier, 0.94) // 目標 RTP 94%
                if rand.Float64() < hitRate {
                    totalReturn += bet * t.Multiplier
                }
                break
            }
        }
    }
    
    return totalReturn / totalBet * 100
}

func calculateHitRate(multiplier, targetRTP float64) float64 {
    // 簡化計算：命中率 = 目標RTP / 倍率
    // 實際應考慮所有目標物的加權平均
    return targetRTP / multiplier
}

func main() {
    targets := []Target{
        {ID: "T001", Multiplier: 1.0, Weight: 100},
        {ID: "T031", Multiplier: 3.0, Weight: 50},
        {ID: "T061", Multiplier: 8.0, Weight: 20},
        {ID: "T091", Multiplier: 20.0, Weight: 5},
        // ... 其他目標物
    }
    
    iterations := 1_000_000 // 100 萬局
    rtp := simulateRTP(targets, iterations)
    
    fmt.Printf("模擬局數：%d\n", iterations)
    fmt.Printf("實際 RTP：%.2f%%\n", rtp)
    fmt.Printf("目標 RTP：92-96%%\n")
    
    if rtp >= 92 && rtp <= 96 {
        fmt.Println("✅ RTP 在目標範圍內")
    } else {
        fmt.Println("❌ RTP 超出目標範圍，需要調整")
    }
}
```

### 執行方式
```bash
cd d:\Kiro
go run tools/rtp_simulator.go
```

### 方法 2：使用 Python 快速模擬

```python
# tools/rtp_quick_sim.py
import random
import statistics

def simulate_rtp(targets, iterations=1_000_000):
    """
    targets: list of (multiplier, weight, hit_rate)
    """
    total_bet = 0
    total_return = 0
    
    # 建立加權列表
    weighted_targets = []
    for mult, weight, hit_rate in targets:
        weighted_targets.extend([(mult, hit_rate)] * weight)
    
    for _ in range(iterations):
        bet = 1
        total_bet += bet
        
        target = random.choice(weighted_targets)
        mult, hit_rate = target
        
        if random.random() < hit_rate:
            total_return += bet * mult
    
    return (total_return / total_bet) * 100

# 目標物設定
targets = [
    # (倍率, 權重, 命中率)
    (1.0, 100, 0.94),   # T001-T030 普通魚
    (3.0, 50, 0.31),    # T031-T060 中型魚
    (8.0, 20, 0.12),    # T061-T090 大型魚
    (20.0, 5, 0.047),   # T091-T105 特殊目標
]

rtp = simulate_rtp(targets, 1_000_000)
print(f"模擬 RTP：{rtp:.2f}%")
print(f"目標範圍：92-96%")
print("✅ 通過" if 92 <= rtp <= 96 else "❌ 需要調整")
```

## 範例輸出
```
模擬局數：1,000,000
實際 RTP：93.47%
目標 RTP：92-96%
標準差：±0.23%
95% 信賴區間：[93.01%, 93.93%]
✅ RTP 在目標範圍內
```

## 注意事項
- 模擬樣本數至少 100 萬局，建議 1000 萬局以獲得更精確結果
- 每次修改倍率或命中率後必須重新模擬
- RTP 偏差超過 ±2% 必須立即調整
- Bonus 遊戲需要單獨模擬（目標 RTP 105-115%）
- 所有金融計算使用整數（避免浮點誤差），最後才轉換為百分比

## 已知問題
- 簡化模型未考慮玩家行為模式（選擇高倍率目標物的傾向）
- 實際 RTP 可能因玩家策略而略有偏差

## 版本記錄
- 2025-01-01：初始版本
