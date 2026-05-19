"""
create_swim_imports.py
為游泳動畫 PNG 建立 Godot .import 檔案

Godot 4 需要每個 PNG 都有對應的 .import 檔案才能正確載入
"""

import os
import hashlib

TARGETS = [
    "T001_grass",
    "T002_bug_g",
    "T003_bug_r",
    "T004_bug_b",
    "T005_pudding",
    "T006_mushroom",
    "T101_mimic",
    "T102_chest",
    "T103_meteor",
    "T104_gold_grass",
    "T105_coin_fish",
]

SPRITE_DIR = "client/chiikawa-pixel/assets/sprites/targets"

# UID 前綴（Godot 4 的 uid 格式）
# 用 MD5 hash 生成唯一 UID
def make_uid(filename: str) -> str:
    h = hashlib.md5(filename.encode()).hexdigest()[:16]
    return f"uid://{h}"

def make_ctex_hash(filename: str) -> str:
    return hashlib.md5(filename.encode()).hexdigest()

def create_import(target_id: str) -> bool:
    png_name = f"{target_id}_swim.png"
    png_path = os.path.join(SPRITE_DIR, png_name)
    import_path = png_path + ".import"

    if not os.path.exists(png_path):
        print(f"  ⚠️  {png_name} 不存在，跳過")
        return False

    uid = make_uid(png_name)
    ctex_hash = make_ctex_hash(png_name)
    ctex_name = f"{png_name}-{ctex_hash}.ctex"

    content = f"""[remap]

importer="texture"
type="CompressedTexture2D"
uid="{uid}"
path="res://.godot/imported/{ctex_name}"
metadata={{
"vram_texture": false
}}

[deps]

source_file="res://assets/sprites/targets/{png_name}"
dest_files=["res://.godot/imported/{ctex_name}"]

[params]

compress/mode=0
compress/high_quality=false
compress/lossy_quality=0.7
compress/uastc_level=0
compress/rdo_quality_loss=0.0
compress/hdr_compression=1
compress/normal_map=0
compress/channel_pack=0
mipmaps/generate=false
mipmaps/limit=-1
roughness/mode=0
roughness/src_normal=""
process/channel_remap/red=0
process/channel_remap/green=1
process/channel_remap/blue=2
process/channel_remap/alpha=3
process/fix_alpha_border=true
process/premult_alpha=false
process/normal_map_invert_y=false
process/hdr_as_srgb=false
process/hdr_clamp_exposure=false
process/size_limit=0
detect_3d/compress_to=1
"""

    with open(import_path, 'w', encoding='utf-8') as f:
        f.write(content)

    print(f"  ✅ {png_name}.import")
    return True


def main():
    print("=" * 60)
    print("建立游泳動畫 .import 檔案")
    print("=" * 60)

    success = 0
    for target_id in TARGETS:
        if create_import(target_id):
            success += 1

    print(f"\n完成！建立 {success}/{len(TARGETS)} 個 .import 檔案")


if __name__ == "__main__":
    main()
