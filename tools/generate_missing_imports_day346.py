"""
DAY-346: 批次補齊缺少 .import 的目標物精靈圖
"""
import os
import glob

TARGETS_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"

IMPORT_TEMPLATE = """\
[remap]

importer="texture"
type="CompressedTexture2D"
uid="uid://auto_{safe_name}"
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

def main():
    png_files = glob.glob(os.path.join(TARGETS_DIR, "*.png"))
    png_files = [f for f in png_files if not f.endswith(".import")]
    
    created = 0
    skipped = 0
    
    for png_path in sorted(png_files):
        import_path = png_path + ".import"
        if os.path.exists(import_path):
            skipped += 1
            continue
        
        filename = os.path.basename(png_path)
        # safe_name: 把非字母數字字元換成底線
        safe_name = filename.replace(".", "_").replace("-", "_")
        
        content = IMPORT_TEMPLATE.format(
            safe_name=safe_name,
            filename=filename
        )
        
        with open(import_path, "w", encoding="utf-8") as f:
            f.write(content)
        
        print(f"  ✅ 建立: {filename}.import")
        created += 1
    
    print(f"\n完成！建立 {created} 個 .import，跳過 {skipped} 個（已存在）")

if __name__ == "__main__":
    main()
