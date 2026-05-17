# -*- coding: utf-8 -*-
"""
修復 usagi 所有幀的一致性
目標：所有幀的 bbox 寬度和高度 diff <= 2px
策略：以 idle 幀的 bbox 為基準，用 shared_scale 重新縮放所有幀
"""
from PIL import Image, ImageEnhance
import numpy as np
import os

CHARS_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\characters"
CELL_SIZE = 96
FIT_SCALE = 0.82

def fix_usagi_consistency():
    states = ["idle", "attack", "bigwin"]
    
    # 載入所有幀
    frames = {}
    for state in states:
        path = os.path.join(CHARS_DIR, f"usagi_{state}.png")
        img = Image.open(path).convert("RGBA")
        frames[state] = img
        bbox = img.getbbox()
        if bbox:
            print(f"usagi_{state}: bbox={bbox}, size={bbox[2]-bbox[0]}x{bbox[3]-bbox[1]}")
    
    # 以 idle 幀的 bbox 為基準計算 shared_scale
    idle_bbox = frames["idle"].getbbox()
    ref_w = idle_bbox[2] - idle_bbox[0]
    ref_h = idle_bbox[3] - idle_bbox[1]
    print(f"\nidle content: {ref_w}x{ref_h}px")
    
    # 計算統一縮放比例（基於 idle）
    common_scale = min(CELL_SIZE / ref_w, CELL_SIZE / ref_h) * FIT_SCALE
    print(f"common_scale: {common_scale:.4f}")
    
    # 對每個幀應用 shared_scale
    for state in states:
        img = frames[state]
        bbox = img.getbbox()
        if not bbox:
            continue
        
        cropped = img.crop(bbox)
        cw, ch = cropped.size
        new_w = max(1, int(cw * common_scale))
        new_h = max(1, int(ch * common_scale))
        resized = cropped.resize((new_w, new_h), Image.NEAREST)
        
        canvas = Image.new("RGBA", (CELL_SIZE, CELL_SIZE), (0, 0, 0, 0))
        paste_x = (CELL_SIZE - new_w) // 2
        # bottom align
        pad = max(0, int(CELL_SIZE * (1 - FIT_SCALE) * 0.4))
        paste_y = CELL_SIZE - new_h - pad
        canvas.paste(resized, (paste_x, paste_y))
        
        # 增強飽和度和對比度
        canvas = ImageEnhance.Color(canvas).enhance(1.3)
        canvas = ImageEnhance.Contrast(canvas).enhance(1.15)
        
        # 儲存
        out_path = os.path.join(CHARS_DIR, f"usagi_{state}.png")
        canvas.save(out_path)
        
        new_bbox = canvas.getbbox()
        if new_bbox:
            nw = new_bbox[2] - new_bbox[0]
            nh = new_bbox[3] - new_bbox[1]
            print(f"  usagi_{state} -> {nw}x{nh}px (saved)")
    
    # 最終 QC
    print("\n=== 最終 QC ===")
    bboxes = {}
    for state in states:
        path = os.path.join(CHARS_DIR, f"usagi_{state}.png")
        img = Image.open(path).convert("RGBA")
        bbox = img.getbbox()
        bboxes[state] = bbox
        if bbox:
            w = bbox[2] - bbox[0]
            h = bbox[3] - bbox[1]
            arr = np.array(img)
            non_t = int((arr[:,:,3] > 10).sum())
            print(f"  {state}: {w}x{h}px, non-transparent={non_t}/{CELL_SIZE*CELL_SIZE}")
    
    heights = [bboxes[s][3]-bboxes[s][1] for s in states if bboxes.get(s)]
    widths  = [bboxes[s][2]-bboxes[s][0] for s in states if bboxes.get(s)]
    h_diff = max(heights) - min(heights)
    w_diff = max(widths) - min(widths)
    status = "✅" if h_diff <= 2 and w_diff <= 4 else "⚠️ "
    print(f"  {status} 一致性: height diff={h_diff}px, width diff={w_diff}px")

if __name__ == "__main__":
    fix_usagi_consistency()
    print("\nDone!")
