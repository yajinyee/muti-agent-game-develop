import cv2, os
VIDEO = r"d:\Kiro\錄製內容 2026-05-24 221337.mp4"
OUT = r"d:\Kiro\tmp\frames2"
os.makedirs(OUT, exist_ok=True)
cap = cv2.VideoCapture(VIDEO)
fps = cap.get(cv2.CAP_PROP_FPS)
total = int(cap.get(cv2.CAP_PROP_FRAME_COUNT))
dur = total/fps
print(f"{fps:.0f}fps {total}幀 {dur:.1f}秒")
times = [t for t in [0,1,2,3,5,7,9,11,13,15,17,19,21,23,25] if t < dur]
for t in times:
    cap.set(cv2.CAP_PROP_POS_FRAMES, int(t*fps))
    ret, frame = cap.read()
    if ret:
        p = f"{OUT}/f{t:04.1f}.jpg"
        cv2.imwrite(p, frame, [cv2.IMWRITE_JPEG_QUALITY, 88])
        print(f"  {p}")
cap.release()
