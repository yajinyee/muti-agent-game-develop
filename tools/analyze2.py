"""深度分析第二段影片 + 對照規格書"""
import cv2, numpy as np, os

FRAMES = r"d:\Kiro\tmp\frames2"
files = sorted(os.listdir(FRAMES))

for fname in files:
    img = cv2.imread(os.path.join(FRAMES, fname))
    h, w = img.shape[:2]
    hsv = cv2.cvtColor(img, cv2.COLOR_BGR2HSV)
    gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)

    # 非深藍背景像素（目標物+UI）
    bg_mask = cv2.inRange(hsv, np.array([90,30,20]), np.array([140,255,120]))
    non_bg = 1.0 - np.sum(bg_mask>0)/(h*w)

    # 白色像素（角色/目標物）
    white = cv2.inRange(hsv, np.array([0,0,200]), np.array([180,40,255]))
    white_r = np.sum(white>0)/(h*w)

    # 特效（高飽和高亮）
    eff = (hsv[:,:,1]>150) & (hsv[:,:,2]>150)
    eff_r = np.sum(eff)/(h*w)

    # 金色
    gold = cv2.inRange(hsv, np.array([15,100,150]), np.array([35,255,255]))
    gold_r = np.sum(gold>0)/(h*w)

    # 底部 UI 區域（按鈕）
    ui = img[h-70:h, :, :]
    ui_bright = np.mean(cv2.cvtColor(ui, cv2.COLOR_BGR2GRAY))

    # 右上角（REC 按鈕位置）
    rec = img[0:55, w-220:w, :]
    rec_bright = np.mean(cv2.cvtColor(rec, cv2.COLOR_BGR2GRAY))

    # 中間遊戲區域（目標物應該在這裡）
    game_area = img[50:h-80, 50:w-50, :]
    game_bright = np.mean(cv2.cvtColor(game_area, cv2.COLOR_BGR2GRAY))
    game_std = np.std(cv2.cvtColor(game_area, cv2.COLOR_BGR2GRAY))

    # 偵測移動物件（找非背景色的連通區域）
    non_bg_img = cv2.bitwise_not(bg_mask)
    # 去掉 UI 區域
    non_bg_img[h-80:, :] = 0
    non_bg_img[:50, :] = 0
    contours, _ = cv2.findContours(non_bg_img, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
    # 過濾太小的（< 100px²）
    valid_contours = [c for c in contours if cv2.contourArea(c) > 100]

    t = fname.replace("f","").replace(".jpg","")
    print(f"\n[{t}秒]")
    print(f"  遊戲區亮度={game_bright:.1f} std={game_std:.1f}")
    print(f"  非背景={non_bg*100:.1f}% 白色={white_r*100:.2f}% 特效={eff_r*100:.2f}% 金色={gold_r*100:.2f}%")
    print(f"  底部UI亮度={ui_bright:.1f}  右上角亮度={rec_bright:.1f}")
    print(f"  遊戲區可見物件數: {len(valid_contours)} 個")
    if valid_contours:
        areas = sorted([cv2.contourArea(c) for c in valid_contours], reverse=True)[:5]
        print(f"  最大物件面積: {[int(a) for a in areas]}")

print("\n" + "="*50)
print("規格書對照診斷")
print("="*50)
print("""
規格書要求 vs 實際觀察：

【目標物】
規格: Max 18 個目標，每 0.8 秒生成一個
實際: 遊戲區幾乎沒有可見物件 → 目標物沒有出現或太小

【角色（砲台）】
規格: 吉伊卡哇在畫面下方，可以射擊
實際: 截圖可見角色在底部，但沒有射擊動作

【AUTO 模式】
規格: AUTO 開啟後自動選擇目標並射擊
實際: 按 AUTO 後沒有自動射擊 → Server 沒連上或 AUTO 邏輯有問題

【畫面品質】
規格: 16-bit Retro Pixel Art，有海底背景、珊瑚礁、氣泡
實際: 背景存在，但目標物缺失，整體空曠

【REC 按鈕】
規格: 右上角應有 ⏺ REC 按鈕
實際: 右上角亮度偏低，按鈕沒有出現
""")
