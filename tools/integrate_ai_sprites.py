"""
integrate_ai_sprites.py — 把 AI 生成圖整合進 characters/ 目錄（DAY-115）
用 process_sprites.py 的核心邏輯處理 AI 生成圖，確保一致性後替換
"""
from PIL import Image, ImageEnhance
import os
import shutil

AI_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\ai_generated"
CHARS_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\characters"
BACKUP_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\characters_backup"

CHARS = ["chiikawa", "hachiware", "usagi"]
STATES = ["idle", "attack", "bigwin"]
TARGET_SIZE = 96  # 輸出尺寸

def remove_background(img: Image.Image) -> Image.Image:
    """去除背景（洋紅色或白色），保留角色"""
    img = img.convert("RGBA")
    pixels = img.load()
    w, h = img.size
    
    for y in range(h):
        for x in range(w):
            r, g, b, a = pixels[x, y]
            # 去除洋紅色背景
            if r > 200 and g < 80 and b > 200:
                pixels[x, y] = (0, 0, 0, 0)
            # 去除純白背景（邊緣 flood fill 更安全，這裡用簡單閾值）
            elif r > 240 and g > 240 and b > 240 and a > 200:
                # 只去除邊緣的白色（簡單判斷：距離邊緣 5px 內）
                if x < 5 or x >= w-5 or y < 5 or y >= h-5:
                    pixels[x, y] = (0, 0, 0, 0)
    return img

def get_bbox(img: Image.Image):
    """取得非透明像素的 bounding box"""
    pixels = img.load()
    w, h = img.size
    min_x, min_y = w, h
    max_x, max_y = 0, 0
    for y in range(h):
        for x in range(w):
            if pixels[x, y][3] > 10:
                min_x = min(min_x, x)
                min_y = min(min_y, y)
                max_x = max(max_x, x)
                max_y = max(max_y, y)
    if min_x > max_x:
        return None
    return (min_x, min_y, max_x+1, max_y+1)

def process_ai_sprite(src_path: str, ref_bbox=None) -> Image.Image:
    """處理 AI 生成圖：去背 → 裁切 → 縮放 → 置中到 96x96"""
    img = Image.open(src_path).convert("RGBA")
    
    # 去背
    img = remove_background(img)
    
    # 取 bbox
    bbox = get_bbox(img)
    if bbox is None:
        print(f"  WARNING: no non-transparent pixels in {src_path}")
        return Image.new("RGBA", (TARGET_SIZE, TARGET_SIZE), (0, 0, 0, 0))
    
    # 裁切到 bbox
    cropped = img.crop(bbox)
    cw, ch = cropped.size
    
    # 如果有參考 bbox（idle 幀），用相同的縮放比例
    if ref_bbox:
        ref_w = ref_bbox[2] - ref_bbox[0]
        ref_h = ref_bbox[3] - ref_bbox[1]
        scale = min((TARGET_SIZE - 16) / ref_w, (TARGET_SIZE - 16) / ref_h)
    else:
        scale = min((TARGET_SIZE - 16) / cw, (TARGET_SIZE - 16) / ch)
    
    new_w = max(1, int(cw * scale))
    new_h = max(1, int(ch * scale))
    
    # 縮放（NEAREST 保持像素感）
    resized = cropped.resize((new_w, new_h), Image.NEAREST)
    
    # 置中到 TARGET_SIZE x TARGET_SIZE
    result = Image.new("RGBA", (TARGET_SIZE, TARGET_SIZE), (0, 0, 0, 0))
    paste_x = (TARGET_SIZE - new_w) // 2
    paste_y = (TARGET_SIZE - new_h) // 2
    result.paste(resized, (paste_x, paste_y))
    
    return result

def main():
    print("=== 整合 AI 生成圖到 characters/ 目錄 ===\n")
    
    # 備份現有 characters/
    if not os.path.exists(BACKUP_DIR):
        shutil.copytree(CHARS_DIR, BACKUP_DIR)
        print(f"✅ 備份現有 characters/ 到 {BACKUP_DIR}\n")
    else:
        print(f"ℹ️  備份已存在，跳過備份\n")
    
    results = {}
    
    for char in CHARS:
        print(f"[{char}]")
        
        # 先處理 idle 幀，取得參考 bbox
        idle_path = os.path.join(AI_DIR, f"{char}_idle.png")
        if not os.path.exists(idle_path):
            print(f"  ❌ {char}_idle.png not found in ai_generated/")
            continue
        
        idle_img = Image.open(idle_path).convert("RGBA")
        idle_img = remove_background(idle_img)
        ref_bbox = get_bbox(idle_img)
        
        for state in STATES:
            ai_path = os.path.join(AI_DIR, f"{char}_{state}.png")
            if not os.path.exists(ai_path):
                print(f"  ❌ {char}_{state}.png not found")
                continue
            
            # 處理 AI 生成圖
            processed = process_ai_sprite(ai_path, ref_bbox if state != "idle" else None)
            
            # 統計
            pixels = processed.load()
            non_trans = sum(1 for y in range(TARGET_SIZE) for x in range(TARGET_SIZE) if pixels[x, y][3] > 10)
            
            # 儲存到 characters/
            out_path = os.path.join(CHARS_DIR, f"{char}_{state}.png")
            processed.save(out_path)
            
            print(f"  ✅ {char}_{state}: {non_trans}/{TARGET_SIZE*TARGET_SIZE} ({100*non_trans//(TARGET_SIZE*TARGET_SIZE)}%) → {out_path}")
            results[f"{char}_{state}"] = non_trans
        
        print()
    
    print("=== 整合完成 ===")
    print(f"共處理 {len(results)} 個 sprites")
    avg_density = sum(results.values()) / len(results) / (TARGET_SIZE * TARGET_SIZE) * 100 if results else 0
    print(f"平均像素密度: {avg_density:.1f}%")
    print("\n下一步：執行 py tools/generate_spritesheet.py 重建 Spritesheet")

if __name__ == "__main__":
    main()
