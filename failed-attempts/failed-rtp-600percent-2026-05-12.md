# 失敗記錄：RTP 600% 異常

**日期**：2026-05-12  
**記錄者**：Balance Agent  
**狀態**：✅ 已解決

---

## 問題描述

執行 RTP 模擬時，發現實際 RTP 高達 600%，遠超目標的 92-96%。

### 症狀

```
RTP 模擬結果（10,000 局）：
- 總下注：10,000
- 總獲獎：60,234
- 實際 RTP：602.34%  ← 嚴重異常！
- 目標 RTP：94%
```

---

## 根本原因分析

### 問題 1：required_hits 公式錯誤

**錯誤的公式**：

```go
// 錯誤：required_hits 計算方式導致保底太快觸發
func calculateRequiredHits(multiplier int, betAmount float64) int {
    // 錯誤：用 multiplier 直接除以 betAmount
    return int(float64(multiplier) / betAmount)  // 這是錯的！
}
```

**問題**：當 `betAmount = 1`，`multiplier = 100` 時，`required_hits = 100`。
但這個公式的意思是「100 次攻擊後必定命中」，而不是「命中率 1/100」。

**正確的公式**：

```go
// 正確：required_hits 應該基於 RTP 目標計算
func calculateHitRate(multiplier int, targetRTP float64) float64 {
    // 命中率 = 目標RTP / 倍率
    // 例如：倍率 100x，目標 RTP 94%
    // 命中率 = 0.94 / 100 = 0.0094 ≈ 0.94%
    return targetRTP / float64(multiplier)
}
```

### 問題 2：特殊目標設了保底機制

**錯誤的設計**：

```go
// 錯誤：特殊目標（20x-50x）設了保底
type Target struct {
    Multiplier  int
    HitRate     float64
    MaxMisses   int  // 保底：最多連續 N 次未命中後必定命中
}

// 特殊目標設定
specialTarget := Target{
    Multiplier: 50,
    HitRate:    0.02,  // 2% 命中率
    MaxMisses:  10,    // 最多 10 次未命中後必定命中！← 這是問題
}
```

**問題**：`MaxMisses = 10` 意味著每 10 次攻擊必定命中一次 50x 目標。
實際 RTP 貢獻 = 50 / 10 = 5.0（500%），遠超目標。

**正確的設計**：

```go
// 正確：特殊目標不設保底，或保底值要合理計算
type Target struct {
    Multiplier int
    HitRate    float64
    // 移除 MaxMisses，或設定合理的值
}

// 如果要設保底，必須確保 RTP 貢獻合理
// 保底 N 次 = 倍率 / N 的 RTP 貢獻
// 要讓 RTP 貢獻 <= 目標 RTP 的 X%
// N >= 倍率 / (目標RTP * X%)
// 例如：50x 目標，目標 RTP 94%，允許貢獻 5%
// N >= 50 / (0.94 * 0.05) = 50 / 0.047 ≈ 1064
// 所以 MaxMisses 至少要 1064，而不是 10！
```

---

## 解決過程

### Step 1：確認問題範圍

```go
// 加入詳細日誌，追蹤每個目標的 RTP 貢獻
func simulateRTP(rounds int) {
    contributions := make(map[int]float64)  // multiplier -> RTP contribution
    
    for i := 0; i < rounds; i++ {
        target := selectTarget()
        if isHit(target) {
            contributions[target.Multiplier] += float64(target.Multiplier)
        }
    }
    
    // 輸出各倍率的 RTP 貢獻
    for mult, contrib := range contributions {
        fmt.Printf("倍率 %dx：RTP 貢獻 %.2f%%\n", mult, contrib/float64(rounds)*100)
    }
}
```

輸出結果：
```
倍率 1x：RTP 貢獻 25.3%
倍率 2x：RTP 貢獻 18.7%
...
倍率 50x：RTP 貢獻 498.2%  ← 問題在這裡！
```

### Step 2：修正公式

```go
// 修正後的目標設定
var TargetConfig = []TargetDef{
    // 普通目標（1-3x）：高命中率
    {Multiplier: 1, HitRate: 0.25},
    {Multiplier: 2, HitRate: 0.12},
    {Multiplier: 3, HitRate: 0.08},
    
    // 中型目標（4-8x）：中命中率
    {Multiplier: 4, HitRate: 0.06},
    {Multiplier: 5, HitRate: 0.05},
    {Multiplier: 8, HitRate: 0.03},
    
    // 大型目標（10-20x）：低命中率
    {Multiplier: 10, HitRate: 0.02},
    {Multiplier: 15, HitRate: 0.012},
    {Multiplier: 20, HitRate: 0.008},
    
    // 特殊目標（30-50x）：極低命中率，不設保底
    {Multiplier: 30, HitRate: 0.004},
    {Multiplier: 50, HitRate: 0.002},
    
    // BOSS（100-500x）：稀有，不設保底
    {Multiplier: 100, HitRate: 0.001},
    {Multiplier: 500, HitRate: 0.0002},
}
```

### Step 3：驗證修正結果

```
修正後 RTP 模擬（100,000 局）：
- 總下注：100,000
- 總獲獎：93,847
- 實際 RTP：93.85%  ✅ 在目標範圍 92-96% 內
```

---

## 教訓

### 核心教訓：捕魚機 RTP 公式的正確理解

**RTP 公式**：
```
RTP = Σ(倍率 × 命中率 × 出現頻率)

其中：
- 倍率：擊殺目標的獎勵倍數
- 命中率：子彈命中目標的機率
- 出現頻率：目標在畫面上出現的機率
```

**設計原則**：
1. **特殊目標不設保底**：保底機制會大幅提高 RTP，必須謹慎使用
2. **先計算理論 RTP，再調整命中率**：不要憑感覺設定命中率
3. **每次修改後立即模擬**：至少 10,000 局，確認 RTP 在目標範圍

**RTP 貢獻計算工具**：

```python
def calculate_rtp_contribution(targets: list) -> float:
    """計算所有目標的理論 RTP 貢獻"""
    total_rtp = 0
    for target in targets:
        multiplier = target['multiplier']
        hit_rate = target['hit_rate']
        appearance_rate = target['appearance_rate']
        
        contribution = multiplier * hit_rate * appearance_rate
        total_rtp += contribution
        
        print(f"倍率 {multiplier}x：貢獻 {contribution*100:.3f}%")
    
    print(f"\n理論 RTP：{total_rtp*100:.2f}%")
    return total_rtp

# 使用範例
targets = [
    {'multiplier': 1, 'hit_rate': 0.85, 'appearance_rate': 0.25},
    {'multiplier': 50, 'hit_rate': 0.002, 'appearance_rate': 0.004},
    # ...
]
calculate_rtp_contribution(targets)
```

---

## 相關 Skill

- `skills/skill-rtp-simulation.md` — RTP 模擬完整技術文件

---

*記錄時間：2026-05-12*  
*解決時間：2026-05-12（當日解決）*
