## DrillLobsterPanel.gd — 鑽頭龍蝦連帶效果面板（DAY-142）
## 業界依據：Royal Fishing JILI 2026「Drill Bit Lobster (80X) — penetrating drill through multiple fish, self-detonates at end of trajectory」
## 顯示鑽頭穿透動畫 + 爆炸效果 + 連帶擊破結果
extends CanvasLayer

# ---- 常數 ----
const PANEL_COLOR := Color(1.0, 0.42, 0.21)  # 橙紅色（龍蝦主題）
const DRILL_COLOR := Color(0.8, 0.6, 0.1)    # 金色（鑽頭）

# ---- 狀態 ----
var _pixel_font: Font = null
var _active_drill: bool = false

# ---- 初始化 ----
func _ready() -> void:
	if ResourceLoader.exists("res://assets/fonts/pixel8.fnt"):
		_pixel_font = load("res://assets/fonts/pixel8.fnt")
	_connect_signals()

func setup(font: Font) -> void:
	if font:
		_pixel_font = font

func _connect_signals() -> void:
	if GameManager.has_signal("drill_lobster_chain"):
		GameManager.drill_lobster_chain.connect(_on_drill_lobster_chain)

# ---- 事件處理 ----

func _on_drill_lobster_chain(data: Dictionary) -> void:
	var phase: String = data.get("phase", "")
	
	match phase:
		"drill_start":
			_show_drill_start(data)
		"explosion":
			_show_explosion(data)
		"result":
			_show_result(data)

func _show_drill_start(data: Dictionary) -> void:
	if _active_drill:
		return
	_active_drill = true
	
	# 頂部橫幅：鑽頭發射
	var banner := _create_banner("🦞 鑽頭龍蝦！穿透攻擊！", PANEL_COLOR)
	add_child(banner)
	
	# 全螢幕橙色閃光
	_flash_screen(Color(1.0, 0.42, 0.21, 0.3), 0.15)
	
	# 2 秒後移除橫幅
	var tween = banner.create_tween()
	tween.tween_interval(1.5)
	tween.tween_property(banner, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func():
		if is_instance_valid(banner): banner.queue_free()
	)

func _show_explosion(data: Dictionary) -> void:
	# 爆炸閃光（更強烈）
	_flash_screen(Color(1.0, 0.6, 0.1, 0.5), 0.2)
	
	# 爆炸橫幅
	var banner := _create_banner("💥 鑽頭爆炸！", Color(1.0, 0.6, 0.1))
	add_child(banner)
	
	var tween = banner.create_tween()
	tween.tween_interval(0.8)
	tween.tween_property(banner, "modulate:a", 0.0, 0.3)
	tween.tween_callback(func():
		if is_instance_valid(banner): banner.queue_free()
	)

func _show_result(data: Dictionary) -> void:
	_active_drill = false
	
	var killed_targets: Array = data.get("killed_targets", [])
	var total_reward: int = data.get("total_reward", 0)
	var killer_id: String = data.get("killer_id", "")
	var killer_name: String = data.get("killer_name", "")
	
	if total_reward <= 0:
		return
	
	var is_self: bool = (killer_id == GameManager.get_meta("player_id", ""))
	
	# 結果彈窗（右側滑入）
	var result_panel := _create_result_panel(
		killer_name, len(killed_targets), total_reward, is_self
	)
	add_child(result_panel)
	
	# 滑入動畫
	result_panel.position.x = 1280
	var tween = result_panel.create_tween()
	tween.tween_property(result_panel, "position:x", 1280 - 220, 0.3)
	tween.tween_interval(2.5)
	tween.tween_property(result_panel, "position:x", 1280, 0.3)
	tween.tween_callback(func():
		if is_instance_valid(result_panel): result_panel.queue_free()
	)
	
	# 自己觸發時：額外全螢幕金色閃光
	if is_self:
		_flash_screen(Color(1.0, 0.85, 0.0, 0.4), 0.25)

# ---- UI 建立輔助 ----

func _create_banner(text: String, color: Color) -> Control:
	var root := Control.new()
	root.set_anchors_preset(Control.PRESET_TOP_WIDE)
	root.position = Vector2(0, 8)
	root.size = Vector2(1280, 32)
	
	var bg := ColorRect.new()
	bg.size = Vector2(1280, 32)
	bg.color = Color(color.r * 0.2, color.g * 0.2, color.b * 0.2, 0.88)
	root.add_child(bg)
	
	var lbl := Label.new()
	lbl.text = text
	lbl.size = Vector2(1280, 32)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	lbl.add_theme_color_override("font_color", color)
	lbl.add_theme_font_size_override("font_size", 14)
	if _pixel_font:
		lbl.add_theme_font_override("font", _pixel_font)
	root.add_child(lbl)
	
	return root

func _create_result_panel(killer_name: String, kill_count: int, reward: int, is_self: bool) -> Control:
	var panel := Control.new()
	panel.position = Vector2(1060, 200)
	panel.size = Vector2(210, 90)
	
	var bg := ColorRect.new()
	bg.size = Vector2(210, 90)
	bg.color = Color(0.08, 0.04, 0.02, 0.92)
	panel.add_child(bg)
	
	# 橙色邊框
	var border := ColorRect.new()
	border.position = Vector2(-1, -1)
	border.size = Vector2(212, 92)
	border.color = Color(1.0, 0.42, 0.21, 0.7)
	border.z_index = -1
	panel.add_child(border)
	
	# 標題
	var title := Label.new()
	title.text = "🦞 鑽頭連帶！"
	title.position = Vector2(4, 4)
	title.size = Vector2(202, 20)
	title.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title.add_theme_color_override("font_color", Color(1.0, 0.42, 0.21))
	title.add_theme_font_size_override("font_size", 11)
	if _pixel_font:
		title.add_theme_font_override("font", _pixel_font)
	panel.add_child(title)
	
	# 玩家名稱
	var name_lbl := Label.new()
	name_lbl.text = killer_name if killer_name != "" else "玩家"
	name_lbl.position = Vector2(4, 24)
	name_lbl.size = Vector2(202, 16)
	name_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	name_lbl.add_theme_color_override("font_color", Color(1.0, 1.0, 1.0))
	name_lbl.add_theme_font_size_override("font_size", 10)
	if _pixel_font:
		name_lbl.add_theme_font_override("font", _pixel_font)
	panel.add_child(name_lbl)
	
	# 擊破數
	var kill_lbl := Label.new()
	kill_lbl.text = "連帶擊破 %d 個目標" % kill_count
	kill_lbl.position = Vector2(4, 42)
	kill_lbl.size = Vector2(202, 16)
	kill_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	kill_lbl.add_theme_color_override("font_color", Color(1.0, 0.8, 0.4))
	kill_lbl.add_theme_font_size_override("font_size", 10)
	if _pixel_font:
		kill_lbl.add_theme_font_override("font", _pixel_font)
	panel.add_child(kill_lbl)
	
	# 獎勵
	var reward_lbl := Label.new()
	reward_lbl.text = "🪙 +%d" % reward
	reward_lbl.position = Vector2(4, 60)
	reward_lbl.size = Vector2(202, 22)
	reward_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	reward_lbl.add_theme_color_override("font_color", Color(1.0, 0.9, 0.2))
	reward_lbl.add_theme_font_size_override("font_size", 14)
	if _pixel_font:
		reward_lbl.add_theme_font_override("font", _pixel_font)
	panel.add_child(reward_lbl)
	
	return panel

func _flash_screen(color: Color, duration: float) -> void:
	var flash := ColorRect.new()
	flash.set_anchors_preset(Control.PRESET_FULL_RECT)
	flash.color = color
	flash.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(flash)
	
	var tween = flash.create_tween()
	tween.tween_property(flash, "color:a", 0.0, duration)
	tween.tween_callback(func():
		if is_instance_valid(flash): flash.queue_free()
	)
