#!/usr/bin/env python3
"""在 HUD.gd 的 set_mission_reset_at 函數後插入 Session Stats 和 Jackpot 面板的初始化函數"""

path = r"d:\Kiro\client\chiikawa-pixel\scripts\ui\HUD.gd"

with open(path, "r", encoding="utf-8") as f:
    content = f.read()

# 找到 ESC 快捷鍵注釋的位置
esc_marker = "# ── ESC 快捷鍵（DAY-049，DAY-053 更新）──────────────────────────"

if esc_marker not in content:
    print("ERROR: ESC marker not found")
    exit(1)

# 要插入的新函數
new_functions = '''
# ── Session Stats 面板（DAY-046，DAY-053 拆分為 SessionStatsPanel.gd）──────

var _session_stats_node: SessionStatsPanelScript = null  # 獨立面板節點

## 初始化 Session Stats 面板（DAY-053）
func _init_session_stats() -> void:
\tvar panel = SessionStatsPanelScript.new()
\tadd_child(panel)
\tpanel.setup(_pixel_font)
\tvar top_bar = get_node_or_null("TopBar")
\tpanel.create_button(top_bar)
\t_session_stats_node = panel


# ── Progressive Jackpot 面板（DAY-048，DAY-053 拆分為 JackpotPanel.gd）──────

var _jackpot_panel_node: JackpotPanelScript = null  # 獨立面板節點

## 初始化 Jackpot 面板（DAY-053）
func _init_jackpot_panel() -> void:
\tvar panel = JackpotPanelScript.new()
\tpanel.position = Vector2(320, 42)  # TopBar 下方，畫面中央
\tpanel.size = Vector2(640, 36)
\tpanel.z_index = 10
\tadd_child(panel)
\tpanel.setup(_pixel_font)
\t_jackpot_panel_node = panel


'''

# 在 ESC 快捷鍵注釋前插入
new_content = content.replace(esc_marker, new_functions + esc_marker)

with open(path, "w", encoding="utf-8") as f:
    f.write(new_content)

# Verify
with open(path, "r", encoding="utf-8") as f:
    verify = f.read()

if "_init_session_stats" in verify and "_init_jackpot_panel" in verify:
    print("OK: Functions inserted successfully")
else:
    print("ERROR: Functions not found after insertion")

lines = verify.split('\n')
print(f"Total lines: {len(lines)}")
