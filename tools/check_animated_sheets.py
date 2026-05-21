# -*- coding: utf-8 -*-
"""檢查 animated sheets 的品質：各幀是否有內容、是否一致"""
from PIL import Image
import os

SHEET_DIR = "D:/Kiro/client/chiikawa-pixel/assets/sprites/sheets"
FRAME_SIZE = 96

for char in ["chiikawa", "hachiware", "usagi"]:
    path = os.path.join(SHEET_DIR, f"{char}_animated.png")
    if not os.path.exists(path):
        print(f"{char}: NOT FOUND")
        continue

    sheet = Image.open(path).convert("RGBA")
    print(f"\n{char} ({sheet.width}x{sheet.height}):")

    state_names = ["idle", "attack", "bigwin"]
    for row, state in enumerate(state_names):
        frames_info = []
        for col in range(4):
            x = col * FRAME_SIZE
            y = row * FRAME_SIZE
            frame = sheet.crop((x, y, x + FRAME_SIZE, y + FRAME_SIZE))
            non_t = sum(1 for px in frame.getdata() if px[3] > 10)
            bbox = frame.getbbox()
            frames_info.append(f"f{col}:{non_t}px")
        print(f"  {state:8s}: {' | '.join(frames_info)}")
