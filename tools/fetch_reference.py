# -*- coding: utf-8 -*-
"""
下載吉伊卡哇 Perler Bead Pattern 參考圖
從 kandipad.com 取得像素圖案作為美術參考
"""
import urllib.request
import re
import os

SAVE_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\reference"
os.makedirs(SAVE_DIR, exist_ok=True)

HEADERS = {
    "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/120.0.0.0 Safari/537.36",
    "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
    "Accept-Language": "en-US,en;q=0.5",
}

PAGES = [
    ("chiikawa", "https://kandipad.com/pattern/chiikawa-4603016"),
    ("hachiware", "https://kandipad.com/pattern/hachiware-from-chiikawa-8483016"),
    ("chiikawa2", "https://kandipad.com/pattern/chiikawa-15113137"),
    ("usagi", "https://kandipad.com/pattern/chiikawa-with-bag-11922333"),
    ("usagi2", "https://kandipad.com/pattern/usagi-from-chiikawa-10593053"),
]

def fetch_html(url):
    req = urllib.request.Request(url, headers=HEADERS)
    with urllib.request.urlopen(req, timeout=15) as r:
        return r.read().decode("utf-8", errors="replace")

def find_images(html):
    """找所有可能的圖片 URL"""
    patterns = [
        r'content="(https://[^"]+\.(?:png|jpg|webp))"',
        r'src="(https://[^"]+\.(?:png|jpg|webp))"',
        r"src='(https://[^']+\.(?:png|jpg|webp))'",
        r'"image":\s*"(https://[^"]+\.(?:png|jpg|webp))"',
        r'(https://[^\s"\'<>]+(?:pattern|preview|thumb|image)[^\s"\'<>]*\.(?:png|jpg|webp))',
    ]
    found = []
    for p in patterns:
        found.extend(re.findall(p, html))
    # 去重，過濾 favicon
    seen = set()
    result = []
    for url in found:
        if url not in seen and "favicon" not in url and "icon" not in url.lower():
            seen.add(url)
            result.append(url)
    return result

def download_image(url, save_path):
    req = urllib.request.Request(url, headers=HEADERS)
    try:
        with urllib.request.urlopen(req, timeout=15) as r:
            data = r.read()
        with open(save_path, "wb") as f:
            f.write(data)
        return len(data)
    except Exception as e:
        return 0

def main():
    print("Fetching Chiikawa reference images from kandipad.com...")
    
    for name, url in PAGES:
        print(f"\n[{name}] {url}")
        try:
            html = fetch_html(url)
            imgs = find_images(html)
            print(f"  Found {len(imgs)} images")
            
            # 下載前 5 張
            for i, img_url in enumerate(imgs[:5]):
                ext = img_url.split(".")[-1].split("?")[0]
                save_path = os.path.join(SAVE_DIR, f"{name}_{i}.{ext}")
                size = download_image(img_url, save_path)
                if size > 1000:
                    print(f"  OK {name}_{i}.{ext} ({size} bytes) <- {img_url[:60]}")
                else:
                    print(f"  SKIP {img_url[:60]} (too small: {size})")
        except Exception as e:
            print(f"  ERROR: {e}")

    # 列出下載結果
    print("\n=== Downloaded files ===")
    for f in os.listdir(SAVE_DIR):
        path = os.path.join(SAVE_DIR, f)
        size = os.path.getsize(path)
        print(f"  {f}: {size} bytes")

if __name__ == "__main__":
    main()
