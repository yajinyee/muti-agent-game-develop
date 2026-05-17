# -*- coding: utf-8 -*-
"""診斷 LV1 RTP 78% 的原因"""
import random
import math

TARGET_RTP = 0.94
BASE_RTP = 0.92

def kill_chance(multiplier, kill_mod=1.0):
    return min(BASE_RTP / multiplier * kill_mod, 0.95)

def required_hits_lv1(multiplier, bet_cost=1):
    safety = 4.0
    expected = multiplier / (TARGET_RTP * bet_cost)
    return max(2, math.ceil(expected * safety))

# 模擬 LV1 單局
bet_cost = 1
fire_rate = 2.0
duration = 300.0
attack_interval = 1.0 / fire_rate

total_bet = 0
total_reward = 0
target_count = 0
hit_counts = []

random.seed(42)
time_elapsed = 0.0

while time_elapsed < duration:
    time_elapsed += attack_interval  # 主迴圈時間推進

    # 選目標（簡化：全部 T001 雜草 2x）
    mult = 2.0
    chance = kill_chance(mult)
    req = required_hits_lv1(mult, bet_cost)

    hit_count = 0
    killed = False

    while hit_count < req and not killed:
        hit_count += 1
        total_bet += bet_cost
        time_elapsed += attack_interval
        if time_elapsed >= duration:
            break
        if random.random() < chance:
            killed = True

    if killed or hit_count >= req:
        total_reward += int(bet_cost * mult)
        target_count += 1
        hit_counts.append(hit_count)

rtp = total_reward / total_bet if total_bet > 0 else 0
avg_hits = sum(hit_counts) / len(hit_counts) if hit_counts else 0

print(f"LV1 模擬結果（5分鐘）：")
print(f"  總投注：{total_bet}")
print(f"  總獎勵：{total_reward}")
print(f"  RTP：{rtp:.2%}")
print(f"  擊破目標數：{target_count}")
print(f"  平均命中次數：{avg_hits:.2f}")
print(f"  required_hits(T001,bet=1)：{required_hits_lv1(2.0, 1)}")
print(f"  kill_chance(T001)：{kill_chance(2.0):.3f}")
print()

# 理論值
kc = kill_chance(2.0)
req = required_hits_lv1(2.0, 1)
# 幾何分布截斷期望值
p = kc
n = req
expected_hits = sum(k * p * (1-p)**(k-1) for k in range(1, n+1)) + n * (1-p)**n
print(f"理論期望命中次數：{expected_hits:.3f}")
print(f"理論 RTP：{2.0 / expected_hits:.2%}")
print()
print(f"問題：主迴圈 time_elapsed += attack_interval 重複計算了時間！")
print(f"每個目標：主迴圈 +0.5s，內迴圈每次 +0.5s")
print(f"→ 第一次命中實際花了 1.0s（主迴圈 + 內迴圈各一次）")
print(f"→ 這導致時間消耗過快，5分鐘內攻擊次數減少，RTP 偏低")
