# -*- coding: utf-8 -*-
"""診斷 LV1 RTP 為何只有 80%"""
import math

TARGET_RTP = 0.94
BASE_RTP = 0.92

# LV1: bet=1, fire_rate=2.0
bet_cost = 1
fire_rate = 2.0

# 5分鐘 = 300秒，每 0.5 秒攻擊一次
attacks_per_session = 300 * fire_rate  # 600 次
total_bet = attacks_per_session * bet_cost  # 600

print(f"LV1 分析：bet={bet_cost}, fire_rate={fire_rate}")
print(f"5分鐘攻擊次數：{attacks_per_session}")
print(f"5分鐘總投注：{total_bet}")
print()

# 目標分布（LV1-3：90% basic, 10% special）
# 主要是 T001 雜草（2x）
# kill_chance = 0.92 / 2 = 0.46
# required_hits = ceil(2 / (0.94 × 1) × 4.0) = ceil(8.51) = 9

mult = 2.0
kill_chance = BASE_RTP / mult
safety = 4.0
expected = mult / (TARGET_RTP * bet_cost)
required = max(2, math.ceil(expected * safety))

print(f"T001 雜草 (2x)：")
print(f"  kill_chance = {kill_chance:.3f}")
print(f"  expected_hits = {expected:.2f}")
print(f"  required_hits = {required}")
print()

# 期望命中次數（幾何分布截斷）
p = kill_chance
n = required
# E[min(X,n)] where X ~ Geometric(p)
# = sum_{k=1}^{n} k * p * (1-p)^(k-1) + n * (1-p)^n
# 簡化：= (1 - (1-p)^n) / p
expected_actual = (1 - (1-p)**n) / p + n * (1-p)**n
print(f"  實際期望命中次數：{expected_actual:.2f}")
print(f"  期望成本：{expected_actual * bet_cost:.2f}")
print(f"  獎勵：{bet_cost * mult:.2f}")
print(f"  單目標 RTP：{(bet_cost * mult) / (expected_actual * bet_cost):.2%}")
print()

# 問題：每次攻擊都選一個新目標，但每個目標需要多次命中
# 實際上玩家每次攻擊都在打不同目標（除非 lock）
# 所以每次攻擊的期望獎勵 = kill_chance × reward = 0.46 × 2 = 0.92
# 每次攻擊的成本 = 1
# 單次攻擊 RTP = 0.92 / 1 = 92%

print("=" * 50)
print("關鍵洞察：")
print(f"  每次攻擊期望獎勵 = kill_chance × reward = {kill_chance:.3f} × {mult} = {kill_chance * mult:.3f}")
print(f"  每次攻擊成本 = {bet_cost}")
print(f"  單次攻擊 RTP = {kill_chance * mult / bet_cost:.2%}")
print()
print("  但問題是：玩家每次攻擊都打新目標（不 lock）")
print("  所以每次攻擊都是獨立的，RTP = kill_chance × multiplier = 92%")
print()
print("  如果玩家 lock 同一個目標打到保底：")
print(f"  期望命中次數 = {expected_actual:.2f}")
print(f"  RTP = {mult / expected_actual:.2%}")
print()
print("  模擬器的問題：simulate_session 每次攻擊都選新目標")
print("  但 required_hits 是針對同一個目標的保底")
print("  → 模擬器邏輯有誤！")
