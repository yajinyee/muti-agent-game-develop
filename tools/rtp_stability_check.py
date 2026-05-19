"""
RTP 穩定性檢查：跑 5 次 10,000 局，確認 RTP 的真實分布範圍
"""
import sys, os
exec(open('tools/simulate_rtp.py', encoding='utf-8').read().split('if __name__')[0])

print("=== RTP Stability Check (5 runs x 10,000 sessions, LV5) ===")
results = []
for i in range(5):
    r = run_simulation(sessions=10000, bet_level=5)
    rtp = r['overall_rtp']
    results.append(rtp)
    print(f"  Run {i+1}: RTP={rtp:.2%}")

print(f"\nMin: {min(results):.2%}  Max: {max(results):.2%}  Avg: {sum(results)/len(results):.2%}")
print(f"Range: {(max(results)-min(results))*100:.2f}%")
print()
print("Conclusion: QA_MAX should be set to", f"{max(results)*100+1:.0f}%")
