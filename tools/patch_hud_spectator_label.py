#!/usr/bin/env python3
"""在 HUD.gd TopBar 初始化後加入觀戰者計數標籤"""
import re

path = r"d:\Kiro\client\chiikawa-pixel\scripts\ui\HUD.gd"

with open(path, "r", encoding="utf-8") as f:
    content = f.read()

# 找到 top_bar.add_child(top_line) 後面緊接 if is_instance_valid(bottom_bar)
# 在中間插入觀戰者計數標籤
insert_code = """		# 觀戰者計數標籤（DAY-055）
		var spectator_lbl = Label.new()
		spectator_lbl.name = "SpectatorCountLabel"
		spectator_lbl.text = ""
		spectator_lbl.position = Vector2(1180, 8)
		spectator_lbl.size = Vector2(90, 24)
		spectator_lbl.add_theme_font_size_override("font_size", 12)
		spectator_lbl.modulate = Color(0.7, 0.85, 1.0, 0.8)
		spectator_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_RIGHT
		top_bar.add_child(spectator_lbl)
"""

# 找到插入點：top_bar.add_child(top_line) 後面
pattern = r'(		top_bar\.add_child\(top_line\)\n)(	if is_instance_valid\(bottom_bar\):)'
replacement = r'\1' + insert_code + r'\2'

new_content = re.sub(pattern, replacement, content)

if new_content == content:
    print("ERROR: Pattern not found!")
    # 嘗試找到 top_line 的位置
    idx = content.find("top_bar.add_child(top_line)")
    if idx >= 0:
        print(f"Found at index {idx}, context:")
        print(repr(content[idx-5:idx+60]))
    else:
        print("top_bar.add_child(top_line) not found at all")
else:
    with open(path, "w", encoding="utf-8") as f:
        f.write(new_content)
    print("SUCCESS: SpectatorCountLabel added to TopBar")
