"""
analyze_gameplay.py — 分析遊玩幀，提取可量化的視覺資訊
不需要 AI 看圖，直接用像素分析
"""
import cv2
import numpy as np
import os
import json

FRAMES_DIR = r"d:\Kiro\tmp\frames"
frames = sorted([f for f in os.listdir(FRAMES_DIR) if f.endswith('.jpg')])

def analyze_frame(path, fname):
    img = cv2.imread(path)
    h, w = img.shape[:2]
    hsv = cv2.cvtColor(img, cv2.COLOR_BGR2HSV)
    
    result = {"file": fname, "size": f"{w}x{h}"}
    
    # 1. 畫面亮度（整體）
    gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
    result["brightness_mean"] = float(np.mean(gray))
    result["brightness_std"] = float(np.std(gray))
    
    # 2. 顏色分布（判斷畫面豐富度）
    # 把畫面分成 9 個區域，分析各區域主色
    regions = {}
    for ry in range(3):
        for rx in range(3):
            y1, y2 = ry*h//3, (ry+1)*h//3
            x1, x2 = rx*w//3, (rx+1)*w//3
            region = img[y1:y2, x1:x2]
            mean_bgr = region.mean(axis=(0,1))
            regions[f"r{ry}{rx}"] = [round(float(c),1) for c in mean_bgr]
    result["region_colors"] = regions
    
    # 3. 偵測 UI 元素（頂部和底部條）
    top_bar = img[0:40, :, :]
    bottom_bar = img[h-60:h, :, :]
    result["top_bar_brightness"] = float(np.mean(cv2.cvtColor(top_bar, cv2.COLOR_BGR2GRAY)))
    result["bottom_bar_brightness"] = float(np.mean(cv2.cvtColor(bottom_bar, cv2.COLOR_BGR2GRAY)))
    
    # 4. 偵測動態物件（非背景色的像素）
    # 背景是深藍色（海底），找非深藍的像素
    # 深藍: H=100-130, S>50, V<100
    lower_bg = np.array([90, 30, 20])
    upper_bg = np.array([140, 255, 120])
    bg_mask = cv2.inRange(hsv, lower_bg, upper_bg)
    non_bg_ratio = 1.0 - (np.sum(bg_mask > 0) / (h * w))
    result["non_background_ratio"] = round(non_bg_ratio, 3)
    
    # 5. 偵測白色像素（角色/目標物輪廓）
    lower_white = np.array([0, 0, 200])
    upper_white = np.array([180, 40, 255])
    white_mask = cv2.inRange(hsv, lower_white, upper_white)
    white_ratio = np.sum(white_mask > 0) / (h * w)
    result["white_pixel_ratio"] = round(white_ratio, 3)
    
    # 6. 偵測亮色特效（高飽和度像素 = 特效/UI）
    high_sat = hsv[:,:,1] > 150
    high_val = hsv[:,:,2] > 150
    effect_mask = high_sat & high_val
    effect_ratio = np.sum(effect_mask) / (h * w)
    result["effect_pixel_ratio"] = round(float(effect_ratio), 3)
    
    # 7. 偵測紅色（BOSS 警告、血條低）
    lower_red1 = np.array([0, 100, 100])
    upper_red1 = np.array([10, 255, 255])
    lower_red2 = np.array([160, 100, 100])
    upper_red2 = np.array([180, 255, 255])
    red_mask = cv2.inRange(hsv, lower_red1, upper_red1) | cv2.inRange(hsv, lower_red2, upper_red2)
    red_ratio = np.sum(red_mask > 0) / (h * w)
    result["red_pixel_ratio"] = round(float(red_ratio), 3)
    
    # 8. 偵測金色/黃色（獎勵、高倍率目標）
    lower_gold = np.array([15, 100, 150])
    upper_gold = np.array([35, 255, 255])
    gold_mask = cv2.inRange(hsv, lower_gold, upper_gold)
    gold_ratio = np.sum(gold_mask > 0) / (h * w)
    result["gold_pixel_ratio"] = round(float(gold_ratio), 3)
    
    # 9. 偵測底部 UI 按鈕區域的文字（亮度分析）
    ui_zone = img[h-60:h, 0:500, :]
    ui_bright = np.mean(cv2.cvtColor(ui_zone, cv2.COLOR_BGR2GRAY))
    result["ui_button_brightness"] = round(float(ui_bright), 1)
    
    # 10. 右上角是否有 REC 按鈕（找特定位置的亮色）
    rec_zone = img[0:50, w-200:w, :]
    rec_brightness = np.mean(cv2.cvtColor(rec_zone, cv2.COLOR_BGR2GRAY))
    result["rec_zone_brightness"] = round(float(rec_brightness), 1)
    
    return result

print("=" * 60)
print("遊玩影片幀分析報告")
print("=" * 60)

all_results = []
for fname in frames:
    path = os.path.join(FRAMES_DIR, fname)
    r = analyze_frame(path, fname)
    all_results.append(r)
    
    t = fname.replace("frame_", "").replace("s.jpg", "")
    print(f"\n[{t}秒]")
    print(f"  畫面亮度: {r['brightness_mean']:.1f} (std={r['brightness_std']:.1f})")
    print(f"  非背景像素: {r['non_background_ratio']*100:.1f}%  白色像素: {r['white_pixel_ratio']*100:.1f}%")
    print(f"  特效像素: {r['effect_pixel_ratio']*100:.1f}%  金色: {r['gold_pixel_ratio']*100:.1f}%  紅色: {r['red_pixel_ratio']*100:.1f}%")
    print(f"  底部UI亮度: {r['ui_button_brightness']:.1f}  右上角亮度: {r['rec_zone_brightness']:.1f}")

# 整體分析
print("\n" + "=" * 60)
print("整體分析")
print("=" * 60)

avg_non_bg = np.mean([r['non_background_ratio'] for r in all_results])
avg_effect = np.mean([r['effect_pixel_ratio'] for r in all_results])
avg_gold = np.mean([r['gold_pixel_ratio'] for r in all_results])
max_red = max([r['red_pixel_ratio'] for r in all_results])
avg_brightness = np.mean([r['brightness_mean'] for r in all_results])

print(f"平均非背景像素: {avg_non_bg*100:.1f}% (目標物密度指標)")
print(f"平均特效像素: {avg_effect*100:.1f}% (視覺豐富度指標)")
print(f"平均金色像素: {avg_gold*100:.1f}% (獎勵/高倍率目標出現頻率)")
print(f"最高紅色像素: {max_red*100:.1f}% (BOSS/警告出現過嗎？)")
print(f"平均畫面亮度: {avg_brightness:.1f}")

# 判斷 REC 按鈕是否出現
rec_brightnesses = [r['rec_zone_brightness'] for r in all_results]
print(f"\nREC 按鈕區域亮度: min={min(rec_brightnesses):.1f} max={max(rec_brightnesses):.1f} avg={np.mean(rec_brightnesses):.1f}")
print("  (> 50 代表右上角有 UI 元素，可能是 REC 按鈕)")

# 問題診斷
print("\n" + "=" * 60)
print("問題診斷")
print("=" * 60)

if avg_non_bg < 0.15:
    print("❌ 目標物密度極低 — 畫面幾乎只有背景，目標物太少或太小")
elif avg_non_bg < 0.25:
    print("⚠️  目標物密度偏低 — 畫面有些目標物但不夠豐富")
else:
    print("✅ 目標物密度正常")

if avg_effect < 0.02:
    print("❌ 特效幾乎不存在 — 射擊/擊破特效太弱或沒有觸發")
elif avg_effect < 0.05:
    print("⚠️  特效偏少 — 視覺回饋不夠強烈")
else:
    print("✅ 特效密度正常")

if avg_gold < 0.005:
    print("❌ 幾乎沒有金色元素 — 高倍率目標/獎勵顯示太少")
else:
    print(f"✅ 金色元素存在 ({avg_gold*100:.2f}%)")

if max_red < 0.01:
    print("ℹ️  整段影片沒有明顯紅色 — BOSS 沒有出現，或警告沒有觸發")

# 儲存結果
with open(r"d:\Kiro\tmp\analysis_result.json", 'w', encoding='utf-8') as f:
    json.dump(all_results, f, ensure_ascii=False, indent=2)
print(f"\n詳細結果已儲存到 d:\\Kiro\\tmp\\analysis_result.json")
