extends Control
# DAY-347 賽季通行證面板
# BGaming Quests 啟發：每月賽季，10個等級，XP 累積升級

const PANEL_WIDTH = 700
const PANEL_HEIGHT = 480

var _font = null
var _panel_bg: ColorRect
var _title_label: Label
var _season_info_label: Label
var _xp_bar_bg: ColorRect
var _xp_bar_fill: ColorRect
var _xp_label: Label
var _level_label: Label
var _tiers_container: Control
var _close_btn: Button
var _xp_popup: Label  # XP 獲得提示

# 賽季狀態
var _current_xp: int = 0
var _current_level: int = 1
var _next_level_xp: int = 100
var _season_id: String = ""
var _days_left: int = 30
var _is_premium: bool = false

# 等級名稱和徽章
const TIER_BADGES = ["🌱", "🗡️", "⚔️", "🛡️", "🏆", "⭐", "🌟", "💫", "🌌", "👑"]
const TIER_NAMES = [
	"新手探索者", "初級獵人", "中級戰士", "高級勇者", "精英鬥士",
	"傳說英雄", "神話戰神", "宇宙霸主", "時空征服者", "終極大師"
]
const TIER_COLORS = [
	Color(0.5, 0.8, 0.5),   # 綠（LV1）
	Color(0.6, 0.7, 0.9),   # 藍（LV2）
	Color(0.7, 0.6, 0.9),   # 紫（LV3）
	Color(0.9, 0.7, 0.4),   # 橙（LV4）
	Color(1.0, 0.85, 0.2),  # 金（LV5）
	Color(1.0, 0.9, 0.3),   # 亮金（LV6）
	Color(0.9, 0.5, 1.0),   # 洋紅（LV7）
	Color(0.5, 0.9, 1.0),   # 青（LV8）
	Color(0.8, 0.6, 1.0),   # 淡紫（LV9）
	Color(1.0, 0.8, 0.0),   # 皇金（LV10）
]

func setup(font):
	_font = font
	_build_ui()
	_connect_signals()

func _build_ui():
	# 主背景
	_panel_bg = ColorRect.new()
	_panel_bg.color = Color(0.05, 0.05, 0.15, 0.95)
	_panel_bg.size = Vector2(PANEL_WIDTH, PANEL_HEIGHT)
	_panel_bg.position = Vector2(
		(1280 - PANEL_WIDTH) / 2,
		(720 - PANEL_HEIGHT) / 2
	)
	add_child(_panel_bg)
	
	# 邊框
	var border = ColorRect.new()
	border.color = Color(0.6, 0.4, 1.0, 0.8)
	border.size = Vector2(PANEL_WIDTH, 3)
	border.position = Vector2(0, 0)
	_panel_bg.add_child(border)
	
	var border_b = ColorRect.new()
	border_b.color = Color(0.6, 0.4, 1.0, 0.8)
	border_b.size = Vector2(PANEL_WIDTH, 3)
	border_b.position = Vector2(0, PANEL_HEIGHT - 3)
	_panel_bg.add_child(border_b)
	
	# 標題
	_title_label = Label.new()
	_title_label.text = "🌟 賽季通行證"
	_title_label.position = Vector2(20, 12)
	_title_label.add_theme_color_override("font_color", Color(0.9, 0.7, 1.0))
	_title_label.add_theme_font_size_override("font_size", 20)
	if _font:
		_title_label.add_theme_font_override("font", _font)
	_panel_bg.add_child(_title_label)
	
	# 賽季資訊
	_season_info_label = Label.new()
	_season_info_label.text = "賽季 2026-05 | 剩餘 30 天"
	_season_info_label.position = Vector2(20, 40)
	_season_info_label.add_theme_color_override("font_color", Color(0.7, 0.7, 0.9))
	_season_info_label.add_theme_font_size_override("font_size", 13)
	if _font:
		_season_info_label.add_theme_font_override("font", _font)
	_panel_bg.add_child(_season_info_label)
	
	# 關閉按鈕
	_close_btn = Button.new()
	_close_btn.text = "✕"
	_close_btn.position = Vector2(PANEL_WIDTH - 40, 8)
	_close_btn.size = Vector2(30, 30)
	_close_btn.add_theme_color_override("font_color", Color(0.8, 0.5, 0.5))
	_panel_bg.add_child(_close_btn)
	_close_btn.pressed.connect(_on_close)
	
	# 等級顯示
	_level_label = Label.new()
	_level_label.text = "LV.1 🌱 新手探索者"
	_level_label.position = Vector2(20, 65)
	_level_label.add_theme_color_override("font_color", Color(0.5, 0.8, 0.5))
	_level_label.add_theme_font_size_override("font_size", 16)
	if _font:
		_level_label.add_theme_font_override("font", _font)
	_panel_bg.add_child(_level_label)
	
	# XP 進度條背景
	_xp_bar_bg = ColorRect.new()
	_xp_bar_bg.color = Color(0.1, 0.1, 0.2, 0.9)
	_xp_bar_bg.size = Vector2(PANEL_WIDTH - 40, 20)
	_xp_bar_bg.position = Vector2(20, 92)
	_panel_bg.add_child(_xp_bar_bg)
	
	# XP 進度條填充
	_xp_bar_fill = ColorRect.new()
	_xp_bar_fill.color = Color(0.5, 0.3, 0.9)
	_xp_bar_fill.size = Vector2(0, 20)
	_xp_bar_fill.position = Vector2(0, 0)
	_xp_bar_bg.add_child(_xp_bar_fill)
	
	# XP 標籤
	_xp_label = Label.new()
	_xp_label.text = "0 / 100 XP"
	_xp_label.position = Vector2(20, 115)
	_xp_label.add_theme_color_override("font_color", Color(0.7, 0.7, 0.9))
	_xp_label.add_theme_font_size_override("font_size", 12)
	if _font:
		_xp_label.add_theme_font_override("font", _font)
	_panel_bg.add_child(_xp_label)
	
	# 等級格子容器
	_tiers_container = Control.new()
	_tiers_container.position = Vector2(10, 140)
	_tiers_container.size = Vector2(PANEL_WIDTH - 20, 320)
	_panel_bg.add_child(_tiers_container)
	
	_build_tier_grid()
	
	# XP 獲得提示（浮動文字）
	_xp_popup = Label.new()
	_xp_popup.text = "+1 XP"
	_xp_popup.position = Vector2(PANEL_WIDTH / 2, 80)
	_xp_popup.add_theme_color_override("font_color", Color(0.8, 1.0, 0.4))
	_xp_popup.add_theme_font_size_override("font_size", 18)
	_xp_popup.visible = false
	if _font:
		_xp_popup.add_theme_font_override("font", _font)
	_panel_bg.add_child(_xp_popup)

func _build_tier_grid():
	# 清除舊格子
	for child in _tiers_container.get_children():
		child.queue_free()
	
	# 建立 10 個等級格子（2行×5列）
	var cell_w = (PANEL_WIDTH - 20) / 5
	var cell_h = 150
	
	for i in range(10):
		var row = i / 5
		var col = i % 5
		var level = i + 1
		
		var cell = ColorRect.new()
		var is_unlocked = level <= _current_level
		var is_current = level == _current_level
		
		if is_current:
			cell.color = Color(0.15, 0.1, 0.3, 0.95)
		elif is_unlocked:
			cell.color = Color(0.08, 0.08, 0.18, 0.9)
		else:
			cell.color = Color(0.05, 0.05, 0.1, 0.8)
		
		cell.size = Vector2(cell_w - 4, cell_h - 4)
		cell.position = Vector2(col * cell_w + 2, row * cell_h + 2)
		_tiers_container.add_child(cell)
		
		# 等級邊框（當前等級高亮）
		if is_current:
			var border = ColorRect.new()
			border.color = TIER_COLORS[i]
			border.size = Vector2(cell_w - 4, 2)
			border.position = Vector2(0, 0)
			cell.add_child(border)
		
		# 徽章
		var badge = Label.new()
		badge.text = TIER_BADGES[i]
		badge.position = Vector2(5, 5)
		badge.add_theme_font_size_override("font_size", 24)
		if not is_unlocked:
			badge.modulate = Color(0.4, 0.4, 0.4)
		cell.add_child(badge)
		
		# 等級數字
		var lv_label = Label.new()
		lv_label.text = "LV.%d" % level
		lv_label.position = Vector2(5, 38)
		lv_label.add_theme_font_size_override("font_size", 11)
		lv_label.add_theme_color_override("font_color", TIER_COLORS[i] if is_unlocked else Color(0.4, 0.4, 0.4))
		if _font:
			lv_label.add_theme_font_override("font", _font)
		cell.add_child(lv_label)
		
		# 等級名稱
		var name_label = Label.new()
		name_label.text = TIER_NAMES[i]
		name_label.position = Vector2(5, 55)
		name_label.add_theme_font_size_override("font_size", 10)
		name_label.add_theme_color_override("font_color", Color(0.8, 0.8, 0.9) if is_unlocked else Color(0.4, 0.4, 0.4))
		if _font:
			name_label.add_theme_font_override("font", _font)
		cell.add_child(name_label)
		
		# 免費獎勵
		var free_rewards = [100, 200, 300, 500, 800, 1200, 1800, 2500, 3500, 5000]
		var free_label = Label.new()
		free_label.text = "🆓 %d 金" % free_rewards[i]
		free_label.position = Vector2(5, 80)
		free_label.add_theme_font_size_override("font_size", 10)
		free_label.add_theme_color_override("font_color", Color(0.9, 0.8, 0.3) if is_unlocked else Color(0.4, 0.4, 0.4))
		if _font:
			free_label.add_theme_font_override("font", _font)
		cell.add_child(free_label)
		
		# 鎖定圖示
		if not is_unlocked:
			var lock_label = Label.new()
			lock_label.text = "🔒"
			lock_label.position = Vector2(cell_w / 2 - 15, cell_h / 2 - 15)
			lock_label.add_theme_font_size_override("font_size", 20)
			lock_label.modulate = Color(0.5, 0.5, 0.5, 0.7)
			cell.add_child(lock_label)

func _connect_signals():
	if GameManager:
		if not GameManager.is_connected("season_pass_updated", _on_season_pass_updated):
			GameManager.season_pass_updated.connect(_on_season_pass_updated)
		if not GameManager.is_connected("season_pass_level_up", _on_season_pass_level_up):
			GameManager.season_pass_level_up.connect(_on_season_pass_level_up)

func _on_season_pass_updated(data: Dictionary):
	_current_xp = data.get("current_xp", 0)
	_current_level = data.get("current_level", 1)
	_next_level_xp = data.get("next_level_xp", 100)
	_season_id = data.get("season_id", "")
	_days_left = data.get("days_left", 30)
	
	_update_display()
	
	# 顯示 XP 獲得提示
	var xp_gained = data.get("xp_gained", 0)
	var xp_source = data.get("xp_source", "")
	if xp_gained > 0 and visible:
		_show_xp_popup(xp_gained, xp_source)

func _on_season_pass_level_up(data: Dictionary):
	var new_level = data.get("new_level", 1)
	var level_name = data.get("level_name", "")
	var badge = data.get("badge_name", "⭐")
	var free_reward = data.get("free_reward", 0)
	
	# 顯示升級慶祝
	_show_level_up_celebration(new_level, level_name, badge, free_reward)

func _update_display():
	# 更新賽季資訊
	_season_info_label.text = "賽季 %s | 剩餘 %d 天" % [_season_id, _days_left]
	
	# 更新等級顯示
	var level_idx = clamp(_current_level - 1, 0, 9)
	_level_label.text = "LV.%d %s %s" % [_current_level, TIER_BADGES[level_idx], TIER_NAMES[level_idx]]
	_level_label.add_theme_color_override("font_color", TIER_COLORS[level_idx])
	
	# 更新 XP 進度條
	var xp_pct = 0.0
	if _next_level_xp > 0:
		# 計算當前等級的起始 XP
		var level_start_xp = 0
		var xp_thresholds = [0, 100, 300, 600, 1000, 1500, 2100, 2800, 3600, 4500]
		if _current_level >= 1 and _current_level <= 10:
			level_start_xp = xp_thresholds[_current_level - 1]
		var level_range = _next_level_xp - level_start_xp
		if level_range > 0:
			xp_pct = float(_current_xp - level_start_xp) / float(level_range)
	
	xp_pct = clamp(xp_pct, 0.0, 1.0)
	var bar_width = (_xp_bar_bg.size.x) * xp_pct
	
	var tween = create_tween()
	tween.tween_property(_xp_bar_fill, "size:x", bar_width, 0.3)
	
	# XP 標籤
	if _current_level >= 10:
		_xp_label.text = "MAX LEVEL 👑 總 XP: %d" % _current_xp
	else:
		_xp_label.text = "%d / %d XP（%.0f%%）" % [_current_xp, _next_level_xp, xp_pct * 100]
	
	# 重建等級格子
	_build_tier_grid()

func _show_xp_popup(xp: int, source: String):
	var source_emoji = "⚔️"
	match source:
		"boss": source_emoji = "👹"
		"bonus": source_emoji = "🌿"
		"combo": source_emoji = "🔥"
		"daily_quest": source_emoji = "🎯"
		"weekly_challenge": source_emoji = "🏆"
	
	_xp_popup.text = "%s +%d XP" % [source_emoji, xp]
	_xp_popup.visible = true
	_xp_popup.modulate = Color(1, 1, 1, 1)
	_xp_popup.position = Vector2(PANEL_WIDTH / 2 - 40, 75)
	
	var tween = create_tween()
	tween.tween_property(_xp_popup, "position:y", 50, 0.8)
	tween.parallel().tween_property(_xp_popup, "modulate:a", 0.0, 0.8)
	tween.tween_callback(func(): _xp_popup.visible = false)

func _show_level_up_celebration(level: int, name: String, badge: String, reward: int):
	# 建立升級慶祝 overlay
	var overlay = ColorRect.new()
	overlay.color = Color(0.0, 0.0, 0.0, 0.0)
	overlay.size = Vector2(1280, 720)
	overlay.position = Vector2(-_panel_bg.position.x, -_panel_bg.position.y)
	_panel_bg.add_child(overlay)
	
	var msg = Label.new()
	msg.text = "🎉 賽季等級提升！\n%s LV.%d %s\n免費獎勵：%d 金幣" % [badge, level, name, reward]
	msg.position = Vector2(640 - 150, 300)
	msg.add_theme_font_size_override("font_size", 22)
	msg.add_theme_color_override("font_color", Color(1.0, 0.9, 0.3))
	if _font:
		msg.add_theme_font_override("font", _font)
	overlay.add_child(msg)
	
	var tween = create_tween()
	tween.tween_property(overlay, "color:a", 0.7, 0.3)
	tween.tween_interval(2.0)
	tween.tween_property(overlay, "color:a", 0.0, 0.5)
	tween.tween_callback(func(): overlay.queue_free())

func _on_close():
	visible = false
