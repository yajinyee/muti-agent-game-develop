import sys, os
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

# 直接 import 並執行模擬，避免 print emoji 問題
import importlib.util
spec = importlib.util.spec_from_file_location("simulate_rtp", "tools/simulate_rtp.py")
mod = importlib.util.load_from_spec = None

# 直接執行模擬邏輯
exec(open('tools/simulate_rtp.py', encoding='utf-8').read().split('if __name__')[0])

print("=== RTP Simulation ===")
for bet_level in [1, 3, 5, 7, 10]:
    result = run_simulation(sessions=2000, bet_level=bet_level)
    rtp = result["overall_rtp"]
    print(f"LV{bet_level:2d} | RTP: {rtp:.2%} | Bonus/session: {result['avg_bonus_per_session']:.2f} | BOSS/session: {result['avg_boss_per_session']:.2f}")

print()
detail = run_simulation(sessions=10000, bet_level=5)
print(f"LV5 (10000 sessions): RTP={detail['overall_rtp']:.2%}")
