# -*- coding: utf-8 -*-
"""
generate_missing_imports.py — 為缺少 .import 檔案的 PNG 批次生成 .import
Godot 4 格式，使用 MD5 hash 作為 ctex 路徑，隨機 UID

用法：py tools/generate_missing_imports.py
"""
import os
import hashlib
import random
import string

TARGETS_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"

def md5_of_file(path):
    """計算檔案的 MD5 hash"""
    h = hashlib.md5()
    with open(path, "rb") as f:
        h.update(f.read())
    return h.hexdigest()

def random_uid():
    """生成 Godot 風格的 UID（uid://[a-z0-9]{13}）"""
    chars = string.ascii_lowercase + string.digits
    return "uid://" + "".join(random.choices(chars, k=13))

def generate_import(png_path):
    """為單個 PNG 生成 .import 檔案"""
    import_path = png_path + ".import"
    if os.path.exists(import_path):
        return False  # 已存在，跳過

    filename = os.path.basename(png_path)
    md5 = md5_of_file(png_path)
    uid = random_uid()
    ctex_name = f"{filename}-{md5}.ctex"
    res_path = f"res://assets/sprites/targets/{filename}"
    ctex_path = f"res://.godot/imported/{ctex_name}"

    content = f"""[remap]

importer="texture"
type="CompressedTexture2D"
uid="{uid}"
path="{ctex_path}"
metadata={{
"vram_texture": false
}}

[deps]

source_file="{res_path}"
dest_files=["{ctex_path}"]

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

    with open(import_path, "w", encoding="utf-8") as f:
        f.write(content)
    return True

def main():
    png_files = [f for f in os.listdir(TARGETS_DIR) if f.endswith(".png")]
    created = 0
    skipped = 0

    for fname in sorted(png_files):
        full_path = os.path.join(TARGETS_DIR, fname)
        if generate_import(full_path):
            print(f"  ✅ Created: {fname}.import")
            created += 1
        else:
            skipped += 1

    print(f"\n完成！新建 {created} 個 .import，跳過 {skipped} 個（已存在）")

if __name__ == "__main__":
    main()
