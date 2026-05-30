#!/usr/bin/env python3
"""為升級的基礎目標物生成 .import 檔案"""
import os

IMPORT_DIR = "client/chiikawa-pixel/assets/sprites/targets"

IMPORT_TEMPLATE = """[remap]

importer="texture"
type="CompressedTexture2D"
uid="uid://auto_{uid_name}"
path="res://.godot/imported/{filename}-auto.ctex"
metadata={{
"vram_texture": false
}}

[deps]

source_file="res://assets/sprites/targets/{filename}"
dest_files=["res://.godot/imported/{filename}-auto.ctex"]

[params]

compress/mode=0
compress/high_quality=false
compress/lossy_quality=0.7
compress/hdr_compression=1
compress/normal_map=0
compress/channel_pack=0
mipmaps/generate=false
mipmaps/limit=-1
roughness/mode=0
roughness/src_normal=""
process/fix_alpha_border=true
process/premult_alpha=false
process/normal_map_invert_y=false
process/hdr_as_srgb=false
process/hdr_clamp_exposure=false
process/size_limit=0
detect_3d/compress_to=1
svg/scale=1.0
editor/scale_with_editor_scale=false
editor/convert_colors_with_editor_theme=false
"""

targets = [
    "T001_grass",
    "T002_bug_g",
    "T003_bug_r",
    "T004_bug_b",
    "T005_pudding",
    "T006_mushroom",
]

for t in targets:
    filename = f"{t}.png"
    uid_name = t.replace("_", "_") + "_png"
    import_content = IMPORT_TEMPLATE.format(
        uid_name=uid_name,
        filename=filename,
    )
    import_path = os.path.join(IMPORT_DIR, filename + ".import")
    with open(import_path, "w", encoding="utf-8") as f:
        f.write(import_content)
    print(f"✅ 建立 {filename}.import")

print("\n完成！Godot 重新開啟專案後會自動重新匯入這些 PNG。")
