"""
把幀圖片縮小後轉成 base64，輸出成文字檔供 AI 讀取
"""
import cv2, base64, os, json

FRAMES_DIR = r"d:\Kiro\tmp\frames"
OUT_FILE = r"d:\Kiro\tmp\frames_b64.json"

frames = sorted([f for f in os.listdir(FRAMES_DIR) if f.endswith('.jpg')])
result = {}

for fname in frames:
    path = os.path.join(FRAMES_DIR, fname)
    img = cv2.imread(path)
    # 縮小到 640x360 減少 token
    h, w = img.shape[:2]
    scale = min(640/w, 360/h)
    new_w, new_h = int(w*scale), int(h*scale)
    img_small = cv2.resize(img, (new_w, new_h))
    _, buf = cv2.imencode('.jpg', img_small, [cv2.IMWRITE_JPEG_QUALITY, 75])
    b64 = base64.b64encode(buf).decode('utf-8')
    result[fname] = b64
    print(f"  {fname}: {len(b64)//1024}KB")

with open(OUT_FILE, 'w') as f:
    json.dump(result, f)

print(f"\n輸出到 {OUT_FILE}")
