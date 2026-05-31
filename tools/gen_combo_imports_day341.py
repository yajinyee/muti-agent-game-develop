"""
gen_combo_imports_day341.py — 為 DAY-341 Combo 音效生成 .import 檔案
"""
import os

SFX_DIR = r"d:\Kiro\client\chiikawa-pixel\assets\audio\sfx"

IMPORT_TEMPLATE = """[remap]

importer="wav"
type="AudioStreamWAV"
uid="uid://PLACEHOLDER_{name}"
path="res://.godot/imported/{name}.wav-PLACEHOLDER_{name}.sample"

[deps]

source_file="res://assets/audio/sfx/{name}.wav"
dest_files=["res://.godot/imported/{name}.wav-PLACEHOLDER_{name}.sample"]

[params]

force/max_rate=false
force/max_rate_hz=44100
edit/trim=false
edit/normalize=false
edit/loop_mode=0
edit/loop_begin=0
edit/loop_end=-1
compress/mode=0
"""

combo_files = ["combo_5", "combo_10", "combo_20", "combo_30"]

for name in combo_files:
    import_path = os.path.join(SFX_DIR, f"{name}.wav.import")
    content = IMPORT_TEMPLATE.format(name=name)
    with open(import_path, 'w', encoding='utf-8') as f:
        f.write(content)
    print(f"✅ {name}.wav.import")

print("完成！")
