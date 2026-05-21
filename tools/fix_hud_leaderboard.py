#!/usr/bin/env python3
"""
fix_hud_leaderboard.py
在 HUD.gd 的 _init_jackpot_panel 函數後插入 _create_leaderboard_panel 和 _on_leaderboard_updated 函數定義
"""
import sys

HUD_PATH = r"d:\Kiro\client\chiikawa-pixel\scripts\ui\HUD.gd"

NEW_FUNCTIONS = '''

# ── 排行榜面板（DAY-058 拆分到 LeaderboardPanel.gd）──────────────────────────

var _leaderboard_node = null  # LeaderboardPanelScript 實例

## 初始化排行榜面板（委派給 LeaderboardPanel.gd）
func _create_leaderboard_panel() -> void:
\t_leaderboard_node = LeaderboardPanelScript.new()
\t_leaderboard_node.setup(self, _pixel_font)

## 排行榜更新事件（委派給 LeaderboardPanel.gd）
func _on_leaderboard_updated(entries: Array) -> void:
\tif _leaderboard_node:
\t\t_leaderboard_node.update(entries, GameManager.get_player_id())

'''

def main():
    with open(HUD_PATH, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 找到 _init_jackpot_panel 函數結尾後的位置
    # 在 "# ── ESC 快捷鍵" 之前插入
    marker = "# ── ESC 快捷鍵"
    if marker not in content:
        print(f"ERROR: marker not found: {marker}")
        sys.exit(1)
    
    # 確認 _create_leaderboard_panel 是否已存在
    if "func _create_leaderboard_panel" in content:
        print("INFO: _create_leaderboard_panel already exists, skipping")
        sys.exit(0)
    
    # 插入新函數
    new_content = content.replace(marker, NEW_FUNCTIONS + marker, 1)
    
    with open(HUD_PATH, 'w', encoding='utf-8') as f:
        f.write(new_content)
    
    # 確認行數
    lines = new_content.split('\n')
    print(f"Done. New line count: {len(lines)}")
    print(f"Inserted _create_leaderboard_panel and _on_leaderboard_updated")

if __name__ == "__main__":
    main()
