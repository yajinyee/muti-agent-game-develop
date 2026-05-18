from PIL import Image
import os

sheets_dir = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\sheets"
chars = ["chiikawa", "hachiware", "usagi"]

for char in chars:
    path = os.path.join(sheets_dir, f"{char}_animated.png")
    if os.path.exists(path):
        img = Image.open(path)
        w, h = img.size
        cols = w // 96
        rows = h // 96
        print(f"{char}: {w}x{h} = {cols}cols x {rows}rows")
    else:
        print(f"{char}: NOT FOUND")
