# -*- coding: utf-8 -*-
"""診斷 usagi_attack 的連通區域問題"""
from PIL import Image
import math
from collections import deque

def remove_bg_white(img, threshold=200):
    img = img.convert("RGBA")
    pixels = img.load()
    w, h = img.size
    def is_white(px):
        r, g, b, a = px
        return r > threshold and g > threshold and b > threshold and a > 10
    queue = deque()
    visited = [[False]*h for _ in range(w)]
    for x in range(w):
        for y in [0, h-1]:
            if not visited[x][y] and is_white(pixels[x, y]):
                queue.append((x, y)); visited[x][y] = True
    for y in range(h):
        for x in [0, w-1]:
            if not visited[x][y] and is_white(pixels[x, y]):
                queue.append((x, y)); visited[x][y] = True
    while queue:
        x, y = queue.popleft()
        pixels[x, y] = (0, 0, 0, 0)
        for dx, dy in [(0,1),(0,-1),(1,0),(-1,0)]:
            nx, ny = x+dx, y+dy
            if 0<=nx<w and 0<=ny<h and not visited[nx][ny] and is_white(pixels[nx,ny]):
                visited[nx][ny] = True; queue.append((nx, ny))
    return img

def connected_components(img, min_area=1):
    alpha = img.getchannel("A")
    pixels = alpha.load()
    width, height = img.size
    visited = [[False]*width for _ in range(height)]
    components = []
    for y in range(height):
        for x in range(width):
            if pixels[x,y] == 0 or visited[y][x]:
                continue
            queue = deque([(x,y)])
            visited[y][x] = True
            area = 0
            min_x = max_x = x; min_y = max_y = y
            while queue:
                cx, cy = queue.popleft()
                area += 1
                min_x = min(min_x, cx); min_y = min(min_y, cy)
                max_x = max(max_x, cx); max_y = max(max_y, cy)
                for dx, dy in ((1,0),(-1,0),(0,1),(0,-1)):
                    nx, ny = cx+dx, cy+dy
                    if 0<=nx<width and 0<=ny<height and pixels[nx,ny]>0 and not visited[ny][nx]:
                        visited[ny][nx] = True; queue.append((nx,ny))
            if area >= min_area:
                components.append({"area": area, "bbox": (min_x, min_y, max_x+1, max_y+1)})
    components.sort(key=lambda c: c["area"], reverse=True)
    return components

# 載入並去背
img = Image.open("D:/Kiro/client/chiikawa-pixel/assets/sprites/characters/usagi_attack.png").convert("RGBA")
print(f"原始: {img.size}")
print(f"原始 bbox: {img.getbbox()}")

# 去白色背景
img_clean = remove_bg_white(img)
print(f"去背後 bbox: {img_clean.getbbox()}")

# 找連通區域
comps = connected_components(img_clean, min_area=1)
print(f"\n連通區域數量: {len(comps)}")
for i, c in enumerate(comps[:10]):
    bbox = c["bbox"]
    w = bbox[2]-bbox[0]; h = bbox[3]-bbox[1]
    print(f"  [{i}] area={c['area']}, bbox={bbox}, size={w}x{h}")

# 儲存去背後的圖供檢查
img_clean.save("D:/Kiro/tools/debug_usagi_clean.png")
print("\n儲存: tools/debug_usagi_clean.png")
