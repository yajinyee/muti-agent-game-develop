# -*- coding: utf-8 -*-
"""
捕魚機 RTP 正確模型分析

捕魚機的核心公式：
  期望 RTP = Σ (擊破機率 × 獎勵) / 每次攻擊成本

正確理解：
  - 玩家每次攻擊花費 bet_cost
  - 目標有 HP，需要多次命中才能擊破
  - 擊破後獲得 bet_cost × multiplier 的獎勵

所以：
  期望命中次數 = 1 / kill_chance_per_hit
  期望攻擊成本 = bet_cost × 期望命中次數
  期望獎勵 = bet_cost × multiplier（擊破後）

  RTP = 期望獎勵 / 期望攻擊成本
      = (bet_cost × multiplier) / (bet_cost × 期望命中次數)
      = multiplier / 期望命中次數
      = multiplier × kill_chance_per_hit

  要讓 RTP = 0.94：
  kill_chance_per_hit = 0.94 / multiplier  ← 這是正確的！

但問題在於「保底機制」：
  如果設定 required_hits = ceil(multiplier / bet_cost × difficulty)
  那麼最多打 required_hits 次必定擊破
  
  實際期望命中次數 = min(1/kill_chance, required_hits)
  
  當 required_hits 很小時（例如 T001 雜草 2x，bet=10，required=ceil(2/10×0.3)=1）
  → 1 次必定擊破！
  → 實際 RTP = multiplier × 1 = 200%（遠超目標）

修正方案：
  required_hits 要足夠大，讓保底不會太快觸發
  正確公式：required_hits = ceil(multiplier / (RTP × bet_cost))
  這樣期望命中次數 ≈ multiplier / (RTP × bet_cost)
  期望成本 = bet_cost × multiplier / (RTP × bet_cost) = multiplier / RTP
  RTP = multiplier / (multiplier / RTP) = RTP ✓
"""

TARGET_RTP = 0.94

def analyze_target(name, multiplier, bet_cost, difficulty_factor, kill_mod=1.0):
    """分析單個目標的實際 RTP"""
    # 目前的 kill_chance
    kill_chance = TARGET_RTP / multiplier * kill_mod
    kill_chance = min(kill_chance, 0.95)
    
    # 目前的 required_hits（保底）
    required_hits_current = max(1, int(multiplier / bet_cost * difficulty_factor + 0.999))
    
    # 正確的 required_hits（讓保底不會太快觸發）
    # 期望命中次數 = 1 / kill_chance
    expected_hits = 1.0 / kill_chance
    required_hits_correct = max(int(expected_hits * 2), 1)  # 保底設為期望的 2 倍
    
    # 計算實際期望命中次數（考慮保底）
    # 幾何分布截斷：E[min(X, n)] = (1 - (1-p)^n) / p
    p = kill_chance
    n = required_hits_current
    if p >= 1.0:
        actual_expected_hits = 1.0
    else:
        actual_expected_hits = (1 - (1-p)**n) / p
        # 加上保底的貢獻
        # P(X > n) = (1-p)^n，這些情況下命中次數 = n
        actual_expected_hits = actual_expected_hits + n * (1-p)**n
    
    # 實際 RTP
    reward = bet_cost * multiplier
    cost = bet_cost * actual_expected_hits
    actual_rtp = reward / cost if cost > 0 else 0
    
    return {
        "name": name,
        "multiplier": multiplier,
        "bet_cost": bet_cost,
        "kill_chance": kill_chance,
        "required_hits_current": required_hits_current,
        "required_hits_correct": required_hits_correct,
        "expected_hits_no_cap": expected_hits,
        "actual_expected_hits": actual_expected_hits,
        "actual_rtp": actual_rtp,
    }

# 分析各目標在 LV5（bet=10）的 RTP
print("=" * 80)
print("捕魚機 RTP 分析 — LV5 (bet_cost=10)")
print("=" * 80)
print(f"{'目標':<12} {'倍率':>5} {'擊破率':>8} {'保底(現)':>8} {'保底(正)':>8} {'期望命中':>8} {'實際RTP':>8}")
print("-" * 80)

targets = [
    ("T001雜草",   2,  0.3),
    ("T002綠蟲",   3,  0.3),
    ("T003紅蟲",   5,  0.35),
    ("T004藍蟲",   6,  0.4),
    ("T005布丁",   8,  0.4),
    ("T006蘑菇",   10, 0.5),
    ("T101擬態",   22, 0.7),  # 平均倍率
    ("T102寶箱",   25, 0.8),
    ("T103流星",   30, 0.6),  # 平均倍率
    ("T104金草",   30, 0.7),
    ("T105金魚",   50, 0.9),
]

bet_cost = 10
for name, mult, diff in targets:
    r = analyze_target(name, mult, bet_cost, diff)
    print(f"{r['name']:<12} {r['multiplier']:>5.0f}x {r['kill_chance']:>7.1%} "
          f"{r['required_hits_current']:>8} {r['required_hits_correct']:>8} "
          f"{r['actual_expected_hits']:>8.1f} {r['actual_rtp']:>7.1%}")

print()
print("=" * 80)
print("問題診斷：")
print("  T001 雜草 2x，bet=10：required_hits = ceil(2/10×0.3) = 1")
print("  → 1 次必定擊破，實際 RTP = 200%（遠超目標）")
print()
print("修正方案：")
print("  required_hits 應該 = ceil(1 / kill_chance) × 安全係數")
print("  讓保底不會在期望命中次數之前觸發")
print()

# 計算正確的 difficulty_factor
print("=" * 80)
print("正確的 difficulty_factor 計算（讓保底 = 期望命中次數 × 1.5）")
print("=" * 80)
print(f"{'目標':<12} {'倍率':>5} {'期望命中':>8} {'正確保底':>8} {'正確diff':>10}")
print("-" * 80)

for name, mult, diff in targets:
    kill_chance = min(TARGET_RTP / mult, 0.95)
    expected_hits = 1.0 / kill_chance
    correct_required = max(int(expected_hits * 1.5) + 1, 2)
    # 反推 difficulty_factor：required = ceil(mult / bet × diff)
    # diff = required × bet / mult
    correct_diff = correct_required * bet_cost / mult
    print(f"{name:<12} {mult:>5.0f}x {expected_hits:>8.1f} {correct_required:>8} {correct_diff:>10.3f}")

print()
print("結論：difficulty_factor 需要大幅提高（從 0.3-0.9 提高到 3-15）")
print("這樣保底才不會在 1-2 次就觸發")
