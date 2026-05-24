"""
extract_frames.py — 從影片提取關鍵幀供分析
用法: py tools/extract_frames.py
"""
import cv2
import os

VIDEO_PATH = r"d:\Kiro\錄製內容 2026-05-24 211548.mp4"
OUT_DIR = r"d:\Kiro\tmp\frames"
os.makedirs(OUT_DIR, exist_ok=True)

cap = cv2.VideoCapture(VIDEO_PATH)
fps = cap.get(cv2.CAP_PROP_FPS)
total = int(cap.get(cv2.CAP_PROP_FRAME_COUNT))
duration = total / fps if fps > 0 else 0

print(f"影片資訊: {fps:.1f}fps, {total} 幀, {duration:.1f} 秒")

# 每 3 秒取一幀，另外取第 0、1、2 秒（觀察開場）
sample_times = [0, 1, 2, 3, 5, 8, 12, 16, 20, 25, 30, 35, 40, 45, 50, 55, 60]
sample_times = [t for t in sample_times if t < duration]

saved = []
for t in sample_times:
    frame_idx = int(t * fps)
    cap.set(cv2.CAP_PROP_POS_FRAMES, frame_idx)
    ret, frame = cap.read()
    if ret:
        path = os.path.join(OUT_DIR, f"frame_{t:05.1f}s.jpg")
        cv2.imwrite(path, frame, [cv2.IMWRITE_JPEG_QUALITY, 85])
        saved.append(path)
        print(f"  saved: {path}")

cap.release()
print(f"\n共提取 {len(saved)} 幀到 {OUT_DIR}")
