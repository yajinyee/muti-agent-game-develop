# -*- coding: utf-8 -*-
"""精確診斷 LV1 RTP 82% 的原因"""
import random
import math

BASE_RTP = 0.92
TARGET_RTP = 0.94

def kill_chance_fn(mult, kill_mod=1.0):
    return min(BASE_RTP / mult * kill_mod, 0.95)

def required_hits_fn(mult):
    kc = kill_chance_fn(mult)
    expected = 1.0 / kc
    return max(2, math.ceil(expected * 3.0))

# 模擬 LV1 單局，詳細追蹤
bet_cost = 1
fire_rate = 2.0
duration = 300.0
attack_interval = 1.0 / fire_rate

total_bet = 0
total_reward = 0
target_count = 0
total_hits = 0

random.seed(42)
time_elapsed = 0.0

while time_elapsed < duration:
    time_elapsed += attack_interval  # 主迴圈推進

    # 選目標（全部 T001 雜草 2x）
    mult = 2.0
    chance = kill_chance_fn(mult)
    req = required_hits_fn(mult)

    hit_count = 0
    killed = False

    while hit_count < req and not killed:
        hit_count += 1
        total_bet += bet_cost
        if hit_count > 1:
            time_elapsed += attack_interval
        if time_elapsed >= duration:
            break
        if random.random() < chance:
            killed = True

    total_hits += hit_count

    if killed or hit_count >= req:
        total_reward += int(bet_cost * mult)
        target_count += 1

rtp = total_reward / total_bet if total_bet > 0 else 0
avg_hits = total_hits / target_count if target_count > 0 else 0

print(f"LV1 詳細模擬（5分鐘）：")
print(f"  總投注：{total_bet}")
print(f"  總獎勵：{total_reward}")
print(f"  RTP：{rtp:.2%}")
print(f"  擊破目標數：{target_count}")
print(f"  總命中次數：{total_hits}")
print(f"  平均命中次數/目標：{avg_hits:.2f}")
print(f"  kill_chance(T001)：{kill_chance_fn(2.0):.3f}")
print(f"  required_hits(T001)：{required_hits_fn(2.0)}")
print()

# 理論計算
kc = kill_chance_fn(2.0)
req = required_hits_fn(2.0)
print(f"理論分析：")
print(f"  kill_chance = {kc:.4f}")
print(f"  required_hits = {req}")
print(f"  期望命中次數（無保底）= {1/kc:.3f}")

# 幾何分布截斷期望值
p = kc
n = req
# E[min(X,n)] = sum_{k=1}^{n-1} k*p*(1-p)^(k-1) + n*(1-p)^(n-1)
# 更準確：E[X | X<=n] * P(X<=n) + n * P(X>n)
# P(X<=n) = 1 - (1-p)^n
# E[X | X<=n] = (1 - (1-p)^n * (1 + n*p)) / (p * (1-(1-p)^n))  # 截斷幾何分布
prob_exceed = (1-p)**n
expected_truncated = (1 - (1-p)**n) / p  # 這是 E[min(X,n)] 的近似
# 更精確：
expected_exact = sum(k * p * (1-p)**(k-1) for k in range(1, n+1)) + n * (1-p)**n
print(f"  期望命中次數（有保底）= {expected_exact:.3f}")
print(f"  理論 RTP = {2.0 / expected_exact:.2%}")
print()
print(f"  注意：主迴圈每次都推進一個 attack_interval")
print(f"  但第一次命中不在內迴圈計時（hit_count>1 才計時）")
print(f"  → 第一次命中的時間已在主迴圈計算，正確")
print()

# 計算實際每個目標的平均時間消耗
avg_time_per_target = duration / target_count if target_count > 0 else 0
print(f"  實際每目標平均時間：{avg_time_per_target:.3f}s")
print(f"  理論每目標平均時間：{expected_exact * attack_interval:.3f}s")
print(f"  主迴圈額外時間：{attack_interval:.3f}s")
print(f"  實際每目標時間（含主迴圈）：{(expected_exact + 1) * attack_interval:.3f}s")
print()
print(f"問題找到了！")
print(f"主迴圈每次推進 {attack_interval}s，但這次時間沒有對應的攻擊（只是時間流逝）")
print(f"實際上每個目標消耗 {(expected_exact + 1) * attack_interval:.3f}s，但只有 {expected_exact:.3f} 次攻擊")
print(f"→ 主迴圈的時間推進是多餘的，導致時間消耗過快，攻擊次數減少")
