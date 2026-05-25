"""
extract_frames_day299.py — 從 DAY-299 錄影提取關鍵幀
"""
import cv2
import os

VIDEO_PATH = r"d:\Kiro\錄製內容 2026-05-24 221337.mp4"
OUT_DIR = r"d:\Kiro\tmp\frames_day299"
os.makedirs(OUT_DIR, exist_ok=True)

cap = cv2.VideoCapture(VIDEO_PATH)
fps = cap.get(cv2.CAP_PROP_FPS)
total = int(cap.get(cv2.CAP_PROP_FRAME_COUNT))
duration = total / fps if fps > 0 else 0

print(f"影片資訊: {fps:.1f}fps, {total} 幀, {duration:.1f} 秒")

# 每 5 秒取一幀
sample_times = list(range(0, int(duration), 5))
if not sample_times:
    sample_times = [0]

saved = []
for t in sample_times:
    frame_idx = int(t * fps)
    cap.set(cv2.CAP_PROP_POS_FRAMES, frame_idx)
    ret, frame = cap.read()
    if ret:
        path = os.path.join(OUT_DIR, f"frame_{t:05d}s.jpg")
        cv2.imwrite(path, frame, [cv2.IMWRITE_JPEG_QUALITY, 85])
        saved.append(path)
        print(f"  saved: {path}")

cap.release()
print(f"\n共提取 {len(saved)} 幀到 {OUT_DIR}")
