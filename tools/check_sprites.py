# -*- coding: utf-8 -*-
from PIL import Image
import os

d = 'D:/Kiro/client/chiikawa-pixel/assets/sprites/characters'
for char in ['chiikawa', 'hachiware', 'usagi']:
    bboxes = {}
    for state in ['idle', 'attack', 'bigwin']:
        img = Image.open(os.path.join(d, f'{char}_{state}.png')).convert('RGBA')
        bbox = img.getbbox()
        bboxes[state] = bbox
    print(f'{char}:')
    for state, bbox in bboxes.items():
        if bbox:
            w = bbox[2] - bbox[0]
            h = bbox[3] - bbox[1]
            print(f'  {state}: bbox={bbox}, size={w}x{h}')
    heights = [bboxes[s][3]-bboxes[s][1] for s in ['idle','attack','bigwin'] if bboxes[s]]
    widths  = [bboxes[s][2]-bboxes[s][0] for s in ['idle','attack','bigwin'] if bboxes[s]]
    print(f'  height range: {min(heights)}-{max(heights)}px (diff={max(heights)-min(heights)}px)')
    print(f'  width  range: {min(widths)}-{max(widths)}px (diff={max(widths)-min(widths)}px)')
    print()
