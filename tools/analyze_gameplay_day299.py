"""
analyze_gameplay_day299.py — 分析 DAY-299 錄影幀
"""
import cv2
import numpy as np
import os

FRAMES_DIR = r"d:\Kiro\tmp\frames_day299"
frames = sorted([f for f in os.listdir(FRAMES_DIR) if f.endswith('.jpg')])

def analyze_frame(path, fname):
    img = cv2.imread(path)
    if img is None:
        return None
    h, w = img.shape[:2]
    hsv = cv2.cvtColor(img, cv2.COLOR_BGR2HSV)
    gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)

    result = {"file": fname, "size": f"{w}x{h}"}
    result["brightness_mean"] = float(np.mean(gray))
    result["brightness_std"] = float(np.std(gray))

    # 非背景像素（深藍海底）
    lower_bg = np.array([90, 30, 20])
    upper_bg = np.array([140, 255, 120])
    bg_mask = cv2.inRange(hsv, lower_bg, upper_bg)
    result["non_bg_ratio"] = round(1.0 - float(np.sum(bg_mask > 0)) / (h * w), 3)

    # 特效像素（高飽和高亮度）
    effect_mask = (hsv[:,:,1] > 150) & (hsv[:,:,2] > 150)
    result["effect_ratio"] = round(float(np.sum(effect_mask)) / (h * w), 3)

    # 金色像素
    gold_mask = cv2.inRange(hsv, np.array([15, 100, 150]), np.array([35, 255, 255]))
    result["gold_ratio"] = round(float(np.sum(gold_mask > 0)) / (h * w), 3)

    # 紅色像素（BOSS/警告）
    red1 = cv2.inRange(hsv, np.array([0, 100, 100]), np.array([10, 255, 255]))
    red2 = cv2.inRange(hsv, np.array([160, 100, 100]), np.array([180, 255, 255]))
    result["red_ratio"] = round(float(np.sum((red1 | red2) > 0)) / (h * w), 3)

    # 白色像素（角色/目標物）
    white_mask = cv2.inRange(hsv, np.array([0, 0, 200]), np.array([180, 40, 255]))
    result["white_ratio"] = round(float(np.sum(white_mask > 0)) / (h * w), 3)

    return result

print("=" * 60)
print("  DAY-299 錄影分析報告")
print("=" * 60)

all_results = []
for fname in frames:
    path = os.path.join(FRAMES_DIR, fname)
    r = analyze_frame(path, fname)
    if r is None:
        continue
    all_results.append(r)
    t = fname.replace("frame_", "").replace("s.jpg", "")
    print(f"\n[{t}秒] {r['size']}")
    print(f"  亮度: {r['brightness_mean']:.1f} (std={r['brightness_std']:.1f})")
    print(f"  非背景: {r['non_bg_ratio']*100:.1f}%  白色: {r['white_ratio']*100:.1f}%")
    print(f"  特效: {r['effect_ratio']*100:.1f}%  金色: {r['gold_ratio']*100:.1f}%  紅色: {r['red_ratio']*100:.1f}%")

if not all_results:
    print("❌ 無法分析幀")
    exit(1)

print("\n" + "=" * 60)
print("  整體評估")
print("=" * 60)

avg_non_bg = np.mean([r['non_bg_ratio'] for r in all_results])
avg_effect = np.mean([r['effect_ratio'] for r in all_results])
avg_gold = np.mean([r['gold_ratio'] for r in all_results])
max_red = max([r['red_ratio'] for r in all_results])

print(f"目標物密度: {avg_non_bg*100:.1f}%", end="  ")
if avg_non_bg > 0.25: print("✅ 良好")
elif avg_non_bg > 0.15: print("⚠️ 偏低")
else: print("❌ 極低")

print(f"特效密度:   {avg_effect*100:.1f}%", end="  ")
if avg_effect > 0.05: print("✅ 良好")
elif avg_effect > 0.02: print("⚠️ 偏少")
else: print("❌ 極少")

print(f"金色元素:   {avg_gold*100:.2f}%", end="  ")
if avg_gold > 0.005: print("✅ 存在")
else: print("⚠️ 偏少")

print(f"紅色警告:   {max_red*100:.2f}%", end="  ")
if max_red > 0.01: print("⚠️ BOSS/警告出現")
else: print("ℹ️ 無 BOSS 場景")

print("\n" + "=" * 60)
print("  體驗評估（基於像素分析）")
print("=" * 60)
print(f"  射擊手感：無法從靜態幀評估（需要動態分析）")
print(f"  視覺清晰度：{'良好' if avg_non_bg > 0.2 else '需改善'} — 目標物密度 {avg_non_bg*100:.1f}%")
print(f"  特效豐富度：{'良好' if avg_effect > 0.03 else '需改善'} — 特效像素 {avg_effect*100:.1f}%")
print(f"  高倍率目標：{'可見' if avg_gold > 0.003 else '不明顯'} — 金色像素 {avg_gold*100:.2f}%")
print()
print("  建議：")
if avg_non_bg < 0.2:
    print("  🔴 目標物密度不足，考慮增加生成頻率或目標物大小")
if avg_effect < 0.03:
    print("  🟡 特效密度偏低，考慮強化命中/擊破特效")
if avg_gold < 0.003:
    print("  🟡 高倍率目標不夠顯眼，Lucky 視覺識別升級有助改善")
print()
print("  ✅ 分析完成")
