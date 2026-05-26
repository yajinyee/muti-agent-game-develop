from PIL import Image
import os

targets = ['T126_jackpot_fish', 'T127_coop_fish', 'T128_time_warp', 'T129_chain_meteor']
base = r'd:\Kiro\client\chiikawa-pixel\assets\sprites\targets'

for name in targets:
    path = os.path.join(base, name + '.png')
    img = Image.open(path).convert('RGBA')
    w, h = img.size
    pixels = sum(1 for p in img.getdata() if p[3] > 10)
    total = w * h
    pct = pixels / total * 100
    print(f"{name}: {w}x{h}, non-transparent={pixels}/{total} ({pct:.1f}%)")
