# TreasureMapPanel.gd ??撖嗉??啣??Ｘ嚗AY-122嚗?
# 璆剔?靘?嚗su.edu嚗?026嚗Ⅱ隤idden Treasure Unlocks? 2026 撟湔?擳???啗隅??
# 3?3 鞈?撘????孵??格??拙‵皛踵摮??遛銝銵???撠?蝺孛?澆窄????
extends Control

# ?澆?憭批???頝?
const CELL_SIZE = 72
const CELL_GAP = 6
const GRID_OFFSET_X = 20
const GRID_OFFSET_Y = 60

# ?澆?憿
const COLOR_EMPTY   = Color(0.12, 0.15, 0.25, 0.9)   # ?芸‵皛選?瘛梯???
const COLOR_FILLED  = Color(0.15, 0.55, 0.25, 0.95)  # 撌脣‵皛選?瘛梁?
const COLOR_LINE    = Color(0.8, 0.65, 0.0, 1.0)     # 摰?銵??

var _panel_bg: Control = null
var _cell_nodes: Array = []  # 9 ?摮?暺?
var _is_open: bool = false
var _current_data: Dictionary = {}

func _ready():
	# ?身?梯?
	visible = false

	# ?? GameManager 閮?
	if GameManager.has_signal("treasure_map_updated"):
		GameManager.treasure_map_updated.connect(_on_treasure_map_updated)
	if GameManager.has_signal("treasure_map_line"):
		GameManager.treasure_map_line.connect(_on_treasure_map_line)
	if GameManager.has_signal("treasure_map_full"):
		GameManager.treasure_map_full.connect(_on_treasure_map_full)

	_build_panel()

func _build_panel() -> void:
	# 銝駁?輯???
	_panel_bg = Control.new()
	_panel_bg.z_index = 80
	add_child(_panel_bg)

	var bg = ColorRect.new()
	bg.size = Vector2(280, 340)
	bg.position = Vector2(490, 190)  # ?恍銝剖亢?椰
	bg.color = Color(0.06, 0.08, 0.16, 0.96)
	_panel_bg.add_child(bg)

	# 璅?
	var title = Label.new()
	title.text = "?儭?撖嗉??啣?"
	title.position = Vector2(500, 198)
	title.add_theme_font_size_override("font_size", 20)
	title.add_theme_color_override("font_color", Color(1.0, 0.85, 0.2))
	_panel_bg.add_child(title)

	# ?舀?憿?隞?脣漲嚗?
	var subtitle = Label.new()
	subtitle.name = "Subtitle"
	subtitle.text = "??格??拙‵皛踵摮?"
	subtitle.position = Vector2(500, 222)
	subtitle.add_theme_font_size_override("font_size", 13)
	subtitle.add_theme_color_override("font_color", Color(0.6, 0.6, 0.7))
	_panel_bg.add_child(subtitle)

	# 撱箇? 3?3 ?澆?
	_cell_nodes.clear()
	for r in range(3):
		for c in range(3):
			var cell = _create_cell(r, c)
			_panel_bg.add_child(cell)
			_cell_nodes.append(cell)

	# ?隤芣?
	var reward_lbl = Label.new()
	reward_lbl.text = "銝銵???撠?蝺???50??釣\n?遛?典? ??500??釣 ??"
	reward_lbl.position = Vector2(500, 318)
	reward_lbl.add_theme_font_size_override("font_size", 12)
	reward_lbl.add_theme_color_override("font_color", Color(0.7, 0.7, 0.5))
	_panel_bg.add_child(reward_lbl)

	# ????
	var close_btn = Button.new()
	close_btn.text = "??"
	close_btn.position = Vector2(748, 196)
	close_btn.size = Vector2(24, 24)
	close_btn.add_theme_font_size_override("font_size", 14)
	close_btn.pressed.connect(hide_panel)
	_panel_bg.add_child(close_btn)

func _create_cell(row: int, col: int) -> Control:
	var cell = Control.new()
	cell.name = "Cell_%d_%d" % [row, col]

	var x = 500 + GRID_OFFSET_X + col * (CELL_SIZE + CELL_GAP)
	var y = 240 + GRID_OFFSET_Y + row * (CELL_SIZE + CELL_GAP)

	# ?澆??
	var bg = ColorRect.new()
	bg.name = "BG"
	bg.size = Vector2(CELL_SIZE, CELL_SIZE)
	bg.position = Vector2(x, y)
	bg.color = COLOR_EMPTY
	cell.add_child(bg)

	# ?澆??內
	var icon_lbl = Label.new()
	icon_lbl.name = "Icon"
	icon_lbl.text = "?"
	icon_lbl.position = Vector2(x + 20, y + 10)
	icon_lbl.add_theme_font_size_override("font_size", 24)
	icon_lbl.add_theme_color_override("font_color", Color(0.5, 0.5, 0.6))
	cell.add_child(icon_lbl)

	# ?澆??迂
	var name_lbl = Label.new()
	name_lbl.name = "Name"
	name_lbl.text = ""
	name_lbl.position = Vector2(x + 4, y + 48)
	name_lbl.add_theme_font_size_override("font_size", 10)
	name_lbl.add_theme_color_override("font_color", Color(0.6, 0.6, 0.7))
	cell.add_child(name_lbl)

	# 憛急遛?暸璅?嚗?閮剝??
	var check = Label.new()
	check.name = "Check"
	check.text = "??"
	check.position = Vector2(x + 48, y + 4)
	check.add_theme_font_size_override("font_size", 18)
	check.add_theme_color_override("font_color", Color(0.3, 1.0, 0.4))
	check.visible = false
	cell.add_child(check)

	return cell

func _on_treasure_map_updated(data: Dictionary) -> void:
	_current_data = data
	_update_grid(data)

func _update_grid(data: Dictionary) -> void:
	var cells = data.get("cells", [])
	var filled_count = data.get("filled_count", 0)

	# ?湔?舀?憿?
	var subtitle = _panel_bg.get_node_or_null("Subtitle")
	if subtitle:
		subtitle.text = "隞?脣漲嚗?d/9 ?? % filled_count"

	# ?湔瘥摮?
	for cell_data in cells:
		var r = cell_data.get("row", 0)
		var c = cell_data.get("col", 0)
		var idx = r * 3 + c
		if idx >= _cell_nodes.size():
			continue

		var cell = _cell_nodes[idx]
		var filled = cell_data.get("filled", false)
		var icon = cell_data.get("icon", "?")
		var name_text = cell_data.get("name", "")

		# ?湔?內
		var icon_lbl = cell.get_node_or_null("Icon")
		if icon_lbl:
			icon_lbl.text = icon
			icon_lbl.add_theme_color_override("font_color",
				Color(1.0, 1.0, 1.0) if filled else Color(0.5, 0.5, 0.6))

		# ?湔?迂
		var name_lbl = cell.get_node_or_null("Name")
		if name_lbl:
			name_lbl.text = name_text

		# ?湔?憿
		var bg = cell.get_node_or_null("BG")
		if bg:
			bg.color = COLOR_FILLED if filled else COLOR_EMPTY

		# ?湔?暸璅?
		var check = cell.get_node_or_null("Check")
		if check:
			check.visible = filled

func _on_treasure_map_line(data: Dictionary) -> void:
	var line_type = data.get("line_type", "")
	var reward = data.get("reward", 0)
	var message = data.get("message", "")"

	# 憿舐內摰??
	_show_reward_popup(message, reward, Color(1.0, 0.85, 0.2))

	# 擃漁摰???/??撠?蝺?
	_highlight_line(line_type)

func _on_treasure_map_full(data: Dictionary) -> void:
	var reward = data.get("reward", 0)
	var message = data.get("message", "?唾牧撖嗉?嚗?)"

	# 憿舐內憭抒??嚗??脤???
	_show_reward_popup(message, reward, Color(1.0, 0.7, 0.0))
	_show_full_map_effect()

func _highlight_line(line_type: String) -> void:
	# 靘?/??撠?蝺???鈭桀??摮?
	var indices: Array = []
	match line_type:
		"row0": indices = [0, 1, 2]
		"row1": indices = [3, 4, 5]
		"row2": indices = [6, 7, 8]
		"col0": indices = [0, 3, 6]
		"col1": indices = [1, 4, 7]
		"col2": indices = [2, 5, 8]
		"diag0": indices = [0, 4, 8]
		"diag1": indices = [2, 4, 6]

	for idx in indices:
		if idx >= _cell_nodes.size():
			continue
		var cell = _cell_nodes[idx]
		var bg = cell.get_node_or_null("BG")
		if bg:
			var tween = create_tween().set_loops(3)
			tween.tween_property(bg, "color", COLOR_LINE, 0.15)
			tween.tween_property(bg, "color", COLOR_FILLED, 0.15)

func _show_reward_popup(message: String, reward: int, color: Color) -> void:
	var popup = Label.new()
	popup.text = "%s\n+%d ?馳" % [message, reward]
	popup.position = Vector2(490, 160)
	popup.add_theme_font_size_override("font_size", 18)
	popup.add_theme_color_override("font_color", color)
	popup.z_index = 82
	add_child(popup)

	var tween = create_tween()
	tween.tween_property(popup, "position:y", 120.0, 1.0).set_ease(Tween.EASE_OUT)
	tween.tween_property(popup, "modulate:a", 0.0, 0.5).set_delay(0.8)
	tween.tween_callback(popup.queue_free)

func _show_full_map_effect() -> void:
	# ?游撐?啣????
	for cell in _cell_nodes:
		var bg = cell.get_node_or_null("BG")
		if bg:
			var tween = create_tween().set_loops(5)
			tween.tween_property(bg, "color", Color(1.0, 0.85, 0.0, 1.0), 0.1)
			tween.tween_property(bg, "color", COLOR_FILLED, 0.1)

func show_panel() -> void:
	visible = true
	_is_open = true
	# 隢???啣????
	if NetworkManager.has_method("send_get_treasure_map"):
		NetworkManager.send_get_treasure_map()

func hide_panel() -> void:
	visible = false
	_is_open = false
