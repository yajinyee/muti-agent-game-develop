# -*- coding: utf-8 -*-
"""
免費像素美術素材下載工具
從 OpenGameArt 和 itch.io 下載免費素材
"""
import os
import urllib.request
import json

DOWNLOAD_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\downloads"

# 可直接下載的免費 CC0/CC-BY 素材清單
# 來源：OpenGameArt.org
FREE_ASSETS = [
    {
        "name": "cute_sprites_pack",
        "url": "https://opengameart.org/sites/default/files/cute_sprites_pack_1.zip",
        "license": "CC0",
        "description": "可愛怪物 sprites，適合作為目標物"
    },
    {
        "name": "animated_coins",
        "url": "https://opengameart.org/sites/default/files/animated_coin_0.png",
        "license": "CC0",
        "description": "金幣動畫"
    },
]

def download_file(url: str, dest_path: str) -> bool:
    """下載檔案"""
    try:
        headers = {
            'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36'
        }
        req = urllib.request.Request(url, headers=headers)
        with urllib.request.urlopen(req, timeout=30) as response:
            with open(dest_path, 'wb') as f:
                f.write(response.read())
        return True
    except Exception as e:
        print(f"  NG 下載失敗: {e}")
        return False

def download_opengameart_asset(asset_id: str, dest_name: str) -> bool:
    """
    從 OpenGameArt 下載素材
    asset_id: OpenGameArt 的素材 ID 或直接 URL
    """
    os.makedirs(DOWNLOAD_DIR, exist_ok=True)
    dest_path = os.path.join(DOWNLOAD_DIR, dest_name)

    if os.path.exists(dest_path):
        print(f"  已存在: {dest_name}")
        return True

    print(f"  下載: {dest_name}...")
    return download_file(asset_id, dest_path)

def list_available_assets():
    """列出可用的免費素材"""
    print("=" * 60)
    print("可用的免費像素美術素材")
    print("=" * 60)
    print()
    print("【OpenGameArt.org - CC0 授權（可商用，不需標註）】")
    print("  1. Cute Sprites Pack 1")
    print("     https://opengameart.org/content/cute-sprites-pack-1")
    print("     適合：可愛怪物目標物")
    print()
    print("  2. 16-bit Cute Character")
    print("     https://opengameart.org/content/16-bit-cute-character")
    print("     適合：角色基底（CC-BY 3.0）")
    print()
    print("  3. Animated Coins")
    print("     https://opengameart.org/content/animated-coins-0")
    print("     適合：金幣動畫")
    print()
    print("【itch.io - 免費商用】")
    print("  4. Pixel Monsters Megapack")
    print("     https://blacis.itch.io/pixel-monsters-mega-pack")
    print("     適合：各種怪物目標物")
    print()
    print("  5. Fish Sprite Bundle")
    print("     https://llenpix.itch.io/fish-sprite")
    print("     適合：魚類目標物")
    print()
    print("【使用方式】")
    print("  1. 手動下載 PNG 到:")
    print(f"     {DOWNLOAD_DIR}")
    print("  2. 執行 tools/process_assets.py 自動裁切整理")
    print("  3. 執行 tools/generate_spritesheet.py 更新 Spritesheet")
    print("  4. 重新匯出 HTML5")

def process_downloaded_png(png_path: str, output_dir: str, cell_size: int = 32):
    """
    處理下載的 spritesheet，裁切成單個 sprite
    """
    from PIL import Image
    import math

    img = Image.open(png_path).convert("RGBA")
    w, h = img.size
    cols = w // cell_size
    rows = h // cell_size

    os.makedirs(output_dir, exist_ok=True)
    base_name = os.path.splitext(os.path.basename(png_path))[0]

    count = 0
    for row in range(rows):
        for col in range(cols):
            x = col * cell_size
            y = row * cell_size
            sprite = img.crop((x, y, x + cell_size, y + cell_size))

            # 跳過全透明的格子
            if sprite.getbbox() is None:
                continue

            out_path = os.path.join(output_dir, f"{base_name}_{row:02d}_{col:02d}.png")
            sprite.save(out_path)
            count += 1

    print(f"  裁切完成: {count} 個 sprite -> {output_dir}")
    return count

if __name__ == "__main__":
    print("像素美術素材下載工具")
    print()
    list_available_assets()
    print()
    print("提示：手動下載素材後，執行 process_assets.py 自動整理")
