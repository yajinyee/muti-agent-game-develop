"""
Spritesheet 生成器
把所有角色 Sprite 合併成 Spritesheet，提升 Godot 載入效能
參考：sethmlarson.dev 的 spritesheet 最佳實踐
"""
from PIL import Image
import os
import json

SPRITES_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites"
OUTPUT_DIR = r"D:\Kiro\client\chiikawa-pixel\assets\sprites\sheets"

def create_character_sheet():
    """建立角色 Spritesheet（3角色 × 3狀態 = 9格）"""
    chars = ["chiikawa", "hachiware", "usagi"]
    states = ["idle", "attack", "bigwin"]
    cell_size = 96  # v5 生成的是 96x96
    cols = len(states)
    rows = len(chars)

    sheet = Image.new("RGBA", (cell_size * cols, cell_size * rows), (0, 0, 0, 0))
    metadata = {"cell_size": cell_size, "cols": cols, "rows": rows, "sprites": {}}

    for row, char in enumerate(chars):
        for col, state in enumerate(states):
            path = os.path.join(SPRITES_DIR, "characters", f"{char}_{state}.png")
            if os.path.exists(path):
                sprite = Image.open(path).convert("RGBA")
                sheet.paste(sprite, (col * cell_size, row * cell_size))
                metadata["sprites"][f"{char}_{state}"] = {
                    "x": col * cell_size,
                    "y": row * cell_size,
                    "w": cell_size,
                    "h": cell_size
                }

    os.makedirs(OUTPUT_DIR, exist_ok=True)
    sheet_path = os.path.join(OUTPUT_DIR, "characters_sheet.png")
    sheet.save(sheet_path)
    print(f"  ✅ characters_sheet.png ({sheet.width}x{sheet.height})")

    # 儲存 metadata
    meta_path = os.path.join(OUTPUT_DIR, "characters_sheet.json")
    with open(meta_path, "w", encoding="utf-8") as f:
        json.dump(metadata, f, indent=2, ensure_ascii=False)
    print(f"  ✅ characters_sheet.json")
    return sheet_path

def create_targets_sheet():
    """建立目標物 Spritesheet"""
    target_files = [f for f in os.listdir(os.path.join(SPRITES_DIR, "targets"))
                    if f.endswith(".png") and "B001" not in f]  # BOSS 單獨處理
    target_files.sort()

    cell_size = 64  # v3 升級到 64x64
    cols = 4
    rows = (len(target_files) + cols - 1) // cols

    sheet = Image.new("RGBA", (cell_size * cols, cell_size * rows), (0, 0, 0, 0))
    metadata = {"cell_size": cell_size, "cols": cols, "rows": rows, "sprites": {}}

    for i, filename in enumerate(target_files):
        row, col = divmod(i, cols)
        path = os.path.join(SPRITES_DIR, "targets", filename)
        sprite = Image.open(path).convert("RGBA")
        # 縮放到 cell_size
        sprite = sprite.resize((cell_size, cell_size), Image.NEAREST)
        sheet.paste(sprite, (col * cell_size, row * cell_size))
        name = filename.replace(".png", "")
        metadata["sprites"][name] = {
            "x": col * cell_size,
            "y": row * cell_size,
            "w": cell_size,
            "h": cell_size
        }

    os.makedirs(OUTPUT_DIR, exist_ok=True)
    sheet_path = os.path.join(OUTPUT_DIR, "targets_sheet.png")
    sheet.save(sheet_path)
    print(f"  ✅ targets_sheet.png ({sheet.width}x{sheet.height})")

    meta_path = os.path.join(OUTPUT_DIR, "targets_sheet.json")
    with open(meta_path, "w", encoding="utf-8") as f:
        json.dump(metadata, f, indent=2, ensure_ascii=False)
    print(f"  ✅ targets_sheet.json")
    return sheet_path

def create_effects_sheet():
    """建立特效 Spritesheet"""
    effect_files = [f for f in os.listdir(os.path.join(SPRITES_DIR, "effects"))
                    if f.endswith(".png") and "warning" not in f]  # warning 單獨處理
    effect_files.sort()

    cell_size = 48  # v2 升級到 48x48
    cols = 4
    rows = (len(effect_files) + cols - 1) // cols

    sheet = Image.new("RGBA", (cell_size * cols, cell_size * rows), (0, 0, 0, 0))
    metadata = {"cell_size": cell_size, "cols": cols, "rows": rows, "sprites": {}}

    for i, filename in enumerate(effect_files):
        row, col = divmod(i, cols)
        path = os.path.join(SPRITES_DIR, "effects", filename)
        sprite = Image.open(path).convert("RGBA")
        sprite = sprite.resize((cell_size, cell_size), Image.NEAREST)
        sheet.paste(sprite, (col * cell_size, row * cell_size))
        name = filename.replace(".png", "")
        metadata["sprites"][name] = {
            "x": col * cell_size,
            "y": row * cell_size,
            "w": cell_size,
            "h": cell_size
        }

    os.makedirs(OUTPUT_DIR, exist_ok=True)
    sheet_path = os.path.join(OUTPUT_DIR, "effects_sheet.png")
    sheet.save(sheet_path)
    print(f"  ✅ effects_sheet.png ({sheet.width}x{sheet.height})")

    meta_path = os.path.join(OUTPUT_DIR, "effects_sheet.json")
    with open(meta_path, "w", encoding="utf-8") as f:
        json.dump(metadata, f, indent=2, ensure_ascii=False)
    print(f"  ✅ effects_sheet.json")

def verify_sheets():
    """驗證 Spritesheet 品質"""
    print("\n[驗證]")
    for name in ["characters_sheet.png", "targets_sheet.png", "effects_sheet.png"]:
        path = os.path.join(OUTPUT_DIR, name)
        if os.path.exists(path):
            img = Image.open(path)
            print(f"  ✅ {name}: {img.width}x{img.height}, mode={img.mode}")
        else:
            print(f"  ❌ {name}: 不存在")

if __name__ == "__main__":
    print("Spritesheet generating...")
    print("\n[Character Spritesheet]")
    create_character_sheet()
    print("\n[目標物 Spritesheet]")
    create_targets_sheet()
    print("\n[特效 Spritesheet]")
    create_effects_sheet()
    verify_sheets()
    print("\n✅ Spritesheet 生成完畢！")
    print(f"輸出目錄: {OUTPUT_DIR}")
