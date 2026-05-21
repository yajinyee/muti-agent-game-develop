# -*- coding: utf-8 -*-
"""從 animated spritesheet 提取各動畫並輸出 GIF 預覽"""
from PIL import Image
import os

SHEET_DIR = "D:/Kiro/client/chiikawa-pixel/assets/sprites/sheets"
GIF_DIR   = "D:/Kiro/client/chiikawa-pixel/assets/sprites/characters/gifs"
FRAME_SIZE = 96

os.makedirs(GIF_DIR, exist_ok=True)

ANIM_CONFIG = {
    "idle":   {"row": 0, "frames": 4, "duration": 250},
    "attack": {"row": 1, "frames": 3, "duration": 125},
    "bigwin": {"row": 2, "frames": 4, "duration": 167},
}

for char in ["chiikawa", "hachiware", "usagi"]:
    sheet_path = os.path.join(SHEET_DIR, f"{char}_animated.png")
    if not os.path.exists(sheet_path):
        print(f"SKIP: {char}")
        continue

    sheet = Image.open(sheet_path).convert("RGBA")

    for state, cfg in ANIM_CONFIG.items():
        row = cfg["row"]
        n_frames = cfg["frames"]
        duration = cfg["duration"]

        frames = []
        for col in range(n_frames):
            x = col * FRAME_SIZE
            y = row * FRAME_SIZE
            frame = sheet.crop((x, y, x + FRAME_SIZE, y + FRAME_SIZE))
            # 放大 2x 讓 GIF 更清晰
            frame = frame.resize((FRAME_SIZE * 2, FRAME_SIZE * 2), Image.NEAREST)
            frames.append(frame)

        # 儲存 GIF
        gif_path = os.path.join(GIF_DIR, f"{char}_{state}.gif")
        frames[0].save(
            gif_path,
            save_all=True,
            append_images=frames[1:],
            duration=duration,
            loop=0,
            disposal=2,
        )
        print(f"  {char}_{state}.gif ({n_frames} frames, {duration}ms/frame)")

print("\n✅ GIF 預覽生成完成！")
print(f"位置: {GIF_DIR}")
