# -*- coding: utf-8 -*-
"""
下載吉伊卡哇透明 PNG 素材
來源：pikakirakuzu.com, pixabay.com
"""
import urllib.request
import os
import re
import json

SAVE_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\reference"
os.makedirs(SAVE_DIR, exist_ok=True)

HEADERS = {
    "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/120.0.0.0 Safari/537.36",
    "Accept": "text/html,application/xhtml+xml,*/*;q=0.8",
    "Accept-Language": "en-US,en;q=0.5",
}

def fetch(url, timeout=15):
    req = urllib.request.Request(url, headers=HEADERS)
    with urllib.request.urlopen(req, timeout=timeout) as r:
        return r.read()

def download_img(url, save_path):
    try:
        data = fetch(url)
        with open(save_path, "wb") as f:
            f.write(data)
        return len(data)
    except Exception as e:
        return 0

def fetch_pikakirakuzu():
    """從 pikakirakuzu.com 取得吉伊卡哇資源"""
    print("[pikakirakuzu.com]")
    try:
        html = fetch("https://pikakirakuzu.com/2024/06/22/resource-chiikawa-fan-resources/").decode("utf-8", errors="replace")
        # 找圖片 URL
        imgs = re.findall(r'(https://[^\s"\'<>]+\.(?:png|gif|webp))', html)
        imgs = [u for u in imgs if "chiikawa" in u.lower() or "hachiware" in u.lower() or "usagi" in u.lower()]
        print(f"  Found {len(imgs)} chiikawa images")
        for i, url in enumerate(imgs[:5]):
            fname = f"pikakirakuzu_{i}.png"
            size = download_img(url, os.path.join(SAVE_DIR, fname))
            if size > 500:
                print(f"  OK {fname} ({size} bytes) <- {url[:60]}")
    except Exception as e:
        print(f"  Error: {e}")

def fetch_pixabay():
    """從 pixabay.com 搜尋吉伊卡哇 PNG"""
    print("[pixabay.com]")
    try:
        # Pixabay API（免費，不需要 key 的搜尋）
        url = "https://pixabay.com/api/?q=chiikawa&image_type=vector&per_page=5&safesearch=true"
        # 注意：pixabay API 需要 key，改用網頁搜尋
        html = fetch("https://pixabay.com/images/search/chiikawa%20png/").decode("utf-8", errors="replace")
        # 找圖片 URL
        imgs = re.findall(r'"previewURL":"(https://[^"]+\.(?:png|jpg|webp))"', html)
        print(f"  Found {len(imgs)} preview images")
        for i, url in enumerate(imgs[:3]):
            url = url.replace("\\u002F", "/")
            fname = f"pixabay_{i}.jpg"
            size = download_img(url, os.path.join(SAVE_DIR, fname))
            if size > 500:
                print(f"  OK {fname} ({size} bytes)")
    except Exception as e:
        print(f"  Error: {e}")

def fetch_deviantart_stickers():
    """從 DeviantArt 取得 Chiikawa Stickers PNG"""
    print("[DeviantArt - Chiikawa Stickers]")
    try:
        html = fetch("https://www.deviantart.com/maknaeae/art/PNG-Pack-Chiikawa-Stickers-1250986927").decode("utf-8", errors="replace")
        # 找 og:image
        og_imgs = re.findall(r'og:image.*?content="(https://[^"]+\.(?:png|jpg|webp))"', html)
        # 找其他圖片
        imgs = re.findall(r'"src":"(https://[^"]+\.(?:png|jpg|webp))"', html)
        all_imgs = og_imgs + imgs
        print(f"  Found {len(all_imgs)} images")
        for i, url in enumerate(all_imgs[:3]):
            fname = f"deviantart_sticker_{i}.png"
            size = download_img(url, os.path.join(SAVE_DIR, fname))
            if size > 1000:
                print(f"  OK {fname} ({size} bytes) <- {url[:60]}")
    except Exception as e:
        print(f"  Error: {e}")

def fetch_aigei():
    """從 aigei.com 取得吉伊卡哇免扣 PNG"""
    print("[aigei.com - chiikawa PNG]")
    try:
        html = fetch("https://www.aigei.com/element/element/chiikawa/").decode("utf-8", errors="replace")
        # 找圖片
        imgs = re.findall(r'"(https://[^"]+chiikawa[^"]+\.(?:png|webp))"', html)
        if not imgs:
            imgs = re.findall(r'src="(https://[^"]+\.(?:png|webp))"', html)
        print(f"  Found {len(imgs)} images")
        for i, url in enumerate(imgs[:5]):
            fname = f"aigei_{i}.png"
            size = download_img(url, os.path.join(SAVE_DIR, fname))
            if size > 1000:
                print(f"  OK {fname} ({size} bytes)")
    except Exception as e:
        print(f"  Error: {e}")

if __name__ == "__main__":
    print("Fetching Chiikawa transparent PNG resources...")
    fetch_pikakirakuzu()
    fetch_deviantart_stickers()
    fetch_aigei()
    
    print("\n=== Downloaded files ===")
    for f in sorted(os.listdir(SAVE_DIR)):
        path = os.path.join(SAVE_DIR, f)
        size = os.path.getsize(path)
        if size > 1000:
            print(f"  {f}: {size:,} bytes")
