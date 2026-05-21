## RoyalChainLightningPanel.gd — 皇家閃電鰻持續連鎖電擊面板（DAY-156）
## 業界依據：royal-fishing.co.uk 2026「Creates chain lightning that shocks nearby fish
##   consecutively until targeting turns off. Devastating against clustered schools.」
## 視覺設計：
##   - chain_start：頂部橫幅滑入（電藍色）+ 全螢幕電藍閃光
##   - jump_N：每跳顯示電擊線條動畫 + 跳數計數器 + 浮動獎勵文字
##   - result：右側滑入彈窗（含跳數/擊破數/總獎勵）
##   - ≥8 跳：雙閃光（傳說連鎖）
extends Control

# ---- 常數 ----
const PANEL_WIDTH := 1280.0
const PANEL_HEIGHT := 720.0

# ---- 節點引用 ----
var _banner: Control = null
var _banner_label: Label = null
var _jump_counter: Label = null
var _result_panel: Control = null
var _pixel_font: FontFile = null

# ---- 狀態 ----
var _is_active: bool = false
var _jump_count: int = 0
var _killed_count: int = 0
var _total_reward: int = 0
var _killer_name: String = ""

func setup(font: FontFile) -> void:
	_pixel_font = font
	_build_ui()
	_connect_signals()

func _build_ui() -> void:
	size = Vector2(PANEL_WIDTH, PANEL_HEIGHT)
	mouse_filter = Control.MOUSE_FILTER_IGNORE

	# 頂部橫幅（chain_start 時滑入）
	_banner = Control.new()
	_banner.position = Vector2(0, -60)
	_banner.size = Vector2(PANEL_WIDTH, 52)
	_banner.mouse_filter = Control.MOUSE_FILTER_IGNORE

	var banner_bg = ColorRect.new()
	banner_bg.size = Vector2(PANEL_WIDTH, 52)
	banner_bg.color = Color(0.0, 0.1, 0.3, 0.88)
	_banner.add_child(banner_bg)

	_banner_label = Label.new()
	_banner_label.position = Vector2(0, 8)
	_banner_label.size = Vector2(PANEL_WIDTH, 36)
	_banner_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner_label.text = "⚡ 皇家閃電鰻 — 持續連鎖電擊！"
	_banner_label.add_theme_color_override("font_color", Color(0.0, 0.75, 1.0))
	if _pixel_font:
		_banner_label.add_theme_font_override("font", _pixel_font)
		_banner_label.add_theme_font_size_override("font_size", 20)
	_banner.add_child(_banner_label)
	add_child(_banner)

	# 跳數計數器（左上角）
	_jump_counter = Label.new()
	_jump_counter.position = Vector2(12, 60)
	_jump_counter.size = Vector2(120, 28)
	_jump_counter.text = ""
	_jump_counter.add_theme_color_override("font_color", Color(0.0, 0.9, 1.0))
	if _pixel_font:
		_jump_counter.add_theme_font_override("font", _pixel_font)
		_jump_counter.add_theme_font_size_override("font_size", 14)
	add_child(_jump_counter)

func _connect_signals() -> void:
	if GameManager.has_signal("royal_chain_lightning"):
		GameManager.royal_chain_lightning.connect(_on_royal_chain_lightning)

# ---- 事件處理 ----

func _on_royal_chain_lightning(data: Dictionary) -> void:
	var phase: String = data.get("phase", "")
	var killer_name: String = data.get("killer_name", "")
	var jump_index: int = data.get("jump_index", 0)
	var jump_entry: Dictionary = data.get("jump_entry", {})
	var total_reward: int = data.get("total_reward", 0)
	var total_jumps: int = data.get("total_jumps", 0)

	match phase:
		"chain_start":
			_killer_name = killer_name
			_jump_count = 0
			_killed_count = 0
			_total_reward = 0
			_is_active = true
			_show_banner()
			_flash_screen(Color(0.0, 0.5, 1.0, 0.35))

		"jump":
			_jump_count = jump_index
			var killed: bool = jump_entry.get("killed", false)
			var reward: int = jump_entry.get("reward", 0)
			var to_x: float = jump_entry.get("to_x", 640.0)
			var to_y: float = jump_entry.get("to_y", 360.0)
			if killed:
				_killed_count += 1
				_total_reward += reward
			_update_jump_counter()
			_show_jump_effect(to_x, to_y, killed, reward)

		"result":
			_is_active = false
			_hide_banner()
			_jump_counter.text = ""
			if total_reward > 0:
				_show_result_panel(total_jumps, _killed_count, total_reward)
				# 傳說連鎖（≥8 跳）雙閃光
				if total_jumps >= 8:
					_flash_screen(Color(0.0, 0.75, 1.0, 0.5))
					await get_tree().create_timer(0.15).timeout
					_flash_screen(Color(0.0, 0.75, 1.0, 0.4))

func _show_banner() -> void:
	_banner.position = Vector2(0, -60)
	var tween = _banner.create_tween()
	tween.tween_property(_banner, "position:y", 0.0, 0.3).set_ease(Tween.EASE_OUT)

func _hide_banner() -> void:
	var tween = _banner.create_tween()
	tween.tween_property(_banner, "position:y", -60.0, 0.25).set_ease(Tween.EASE_IN)

func _update_jump_counter() -> void:
	_jump_counter.text = "⚡ 第 %d 跳" % _jump_count

func _show_jump_effect(x: float, y: float, killed: bool, reward: int) -> void:
	# 電擊閃光（在目標位置）
	var flash = ColorRect.new()
	flash.size = Vector2(40, 40)
	flash.position = Vector2(x - 20, y - 20)
	flash.color = Color(0.0, 0.8, 1.0, 0.9)
	flash.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(flash)

	var tween = flash.create_tween()
	tween.tween_property(flash, "modulate:a", 0.0, 0.25)
	tween.tween_callback(func(): if is_instance_valid(flash): flash.queue_free())

	# 浮動獎勵文字（擊破時）
	if killed and reward > 0:
		var lbl = Label.new()
		lbl.text = "+%d" % reward
		lbl.position = Vector2(x - 20, y - 30)
		lbl.size = Vector2(60, 20)
		lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
		lbl.add_theme_color_override("font_color", Color(0.0, 1.0, 1.0))
		if _pixel_font:
			lbl.add_theme_font_override("font", _pixel_font)
			lbl.add_theme_font_size_override("font_size", 12)
		lbl.mouse_filter = Control.MOUSE_FILTER_IGNORE
		add_child(lbl)

		var lbl_tween = lbl.create_tween()
		lbl_tween.tween_property(lbl, "position:y", lbl.position.y - 24, 0.8)
		lbl_tween.parallel().tween_property(lbl, "modulate:a", 0.0, 0.8)
		lbl_tween.tween_callback(func(): if is_instance_valid(lbl): lbl.queue_free())

func _show_result_panel(jumps: int, killed: int, reward: int) -> void:
	if is_instance_valid(_result_panel):
		_result_panel.queue_free()

	_result_panel = Control.new()
	_result_panel.position = Vector2(PANEL_WIDTH + 10, 200)
	_result_panel.size = Vector2(220, 140)
	_result_panel.mouse_filter = Control.MOUSE_FILTER_IGNORE

	var bg = ColorRect.new()
	bg.size = Vector2(220, 140)
	bg.color = Color(0.0, 0.05, 0.15, 0.92)
	_result_panel.add_child(bg)

	var border = ColorRect.new()
	border.size = Vector2(220, 3)
	border.color = Color(0.0, 0.75, 1.0, 1.0)
	_result_panel.add_child(border)

	var title = Label.new()
	title.position = Vector2(0, 8)
	title.size = Vector2(220, 24)
	title.text = "⚡ 皇家閃電鰻連鎖"
	title.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title.add_theme_color_override("font_color", Color(0.0, 0.9, 1.0))
	if _pixel_font:
		title.add_theme_font_override("font", _pixel_font)
		title.add_theme_font_size_override("font_size", 13)
	_result_panel.add_child(title)

	var lines = [
		"連鎖跳數：%d 跳" % jumps,
		"擊破目標：%d 個" % killed,
		"獲得獎勵：%d 金幣" % reward,
	]
	for i in range(lines.size()):
		var lbl = Label.new()
		lbl.position = Vector2(12, 38 + i * 26)
		lbl.size = Vector2(196, 22)
		lbl.text = lines[i]
		lbl.add_theme_color_override("font_color", Color(0.85, 0.95, 1.0))
		if _pixel_font:
			lbl.add_theme_font_override("font", _pixel_font)
			lbl.add_theme_font_size_override("font_size", 12)
		_result_panel.add_child(lbl)

	add_child(_result_panel)

	# 從右側滑入
	var tween = _result_panel.create_tween()
	tween.tween_property(_result_panel, "position:x", PANEL_WIDTH - 230.0, 0.35).set_ease(Tween.EASE_OUT)
	# 3 秒後淡出
	tween.tween_interval(3.0)
	tween.tween_property(_result_panel, "modulate:a", 0.0, 0.4)
	tween.tween_callback(func():
		if is_instance_valid(_result_panel):
			_result_panel.queue_free()
			_result_panel = null
	)

func _flash_screen(color: Color) -> void:
	var flash = ColorRect.new()
	flash.size = Vector2(PANEL_WIDTH, PANEL_HEIGHT)
	flash.color = color
	flash.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(flash)

	var tween = flash.create_tween()
	tween.tween_property(flash, "modulate:a", 0.0, 0.3)
	tween.tween_callback(func(): if is_instance_valid(flash): flash.queue_free())
