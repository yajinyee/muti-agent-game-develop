# -*- coding: utf-8 -*-
import urllib.request
import os
import re

headers = {"User-Agent": "Mozilla/5.0 Chrome/120.0.0.0"}
ref_dir = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\reference"

urls = [
    ("usagi_ref2", "https://kandipad.com/pattern/usagi-from-chiikawa-10593053"),
    ("chiikawa_ref3", "https://kandipad.com/pattern/chiikawa-15113137"),
]

for name, url in urls:
    req = urllib.request.Request(url, headers=headers)
    with urllib.request.urlopen(req, timeout=15) as r:
        html = r.read().decode("utf-8", errors="replace")
    imgs = re.findall(r"(https://kandipad\.com/assets/images/projects/pp/full/[^\s\"'<>]+\.png)", html)
    print(f"{name}: {len(imgs)} pattern imgs")
    for i, img_url in enumerate(imgs[:1]):
        req2 = urllib.request.Request(img_url, headers=headers)
        with urllib.request.urlopen(req2, timeout=15) as r2:
            data = r2.read()
        fname = f"{name}_{i}.png"
        with open(os.path.join(ref_dir, fname), "wb") as f:
            f.write(data)
        print(f"  Saved {fname} ({len(data)} bytes)")
