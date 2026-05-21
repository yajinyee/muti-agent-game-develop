# TreasureMapPanel.gd — 寶藏地圖面板（DAY-122）
# 業界依據：bsu.edu（2026）確認「Hidden Treasure Unlocks」是 2026 年捕魚機最新趨勢
# 3×3 賓果式地圖，擊破特定目標物填滿格子，集滿一行/列/對角線觸發寶藏獎勵
extends Control

# 格子大小和間距
const CELL_SIZE = 72
const CELL_GAP = 6
const GRID_OFFSET_X = 20
const GRID_OFFSET_Y = 60

# 格子顏色
const COLOR_EMPTY   = Color(0.12, 0.15, 0.25, 0.9)   # 未填滿：深藍灰
const COLOR_FILLED  = Color(0.15, 0.55, 0.25, 0.95)  # 已填滿：深綠
const COLOR_LINE    = Color(0.8, 0.65, 0.0, 1.0)     # 完成行：金色

var _panel_bg: Control = null
var _cell_nodes: Array = []  # 9 個格子節點
var _is_open: bool = false
var _current_data: Dictionary = {}

func _ready():
	# 預設隱藏
	visible = false

	# 連接 GameManager 訊號
	if GameManager.has_signal("treasure_map_updated"):
		GameManager.treasure_map_updated.connect(_on_treasure_map_updated)
	if GameManager.has_signal("treasure_map_line"):
		GameManager.treasure_map_line.connect(_on_treasure_map_line)
	if GameManager.has_signal("treasure_map_full"):
		GameManager.treasure_map_full.connect(_on_treasure_map_full)

	_build_panel()

func _build_panel() -> void:
	# 主面板背景
	_panel_bg = Control.new()
	_panel_bg.z_index = 80
	add_child(_panel_bg)

	var bg = ColorRect.new()
	bg.size = Vector2(280, 340)
	bg.position = Vector2(490, 190)  # 畫面中央偏左
	bg.color = Color(0.06, 0.08, 0.16, 0.96)
	_panel_bg.add_child(bg)

	# 標題
	var title = Label.new()
	title.text = "🗺️ 寶藏地圖"
	title.position = Vector2(500, 198)
	title.add_theme_font_size_override("font_size", 20)
	title.add_theme_color_override("font_color", Color(1.0, 0.85, 0.2))
	_panel_bg.add_child(title)

	# 副標題（今日進度）
	var subtitle = Label.new()
	subtitle.name = "Subtitle"
	subtitle.text = "擊破目標物填滿格子"
	subtitle.position = Vector2(500, 222)
	subtitle.add_theme_font_size_override("font_size", 13)
	subtitle.add_theme_color_override("font_color", Color(0.6, 0.6, 0.7))
	_panel_bg.add_child(subtitle)

	# 建立 3×3 格子
	_cell_nodes.clear()
	for r in range(3):
		for c in range(3):
			var cell = _create_cell(r, c)
			_panel_bg.add_child(cell)
			_cell_nodes.append(cell)

	# 獎勵說明
	var reward_lbl = Label.new()
	reward_lbl.text = "一行/列/對角線 → 50×投注\n集滿全圖 → 500×投注 🏆"
	reward_lbl.position = Vector2(500, 318)
	reward_lbl.add_theme_font_size_override("font_size", 12)
	reward_lbl.add_theme_color_override("font_color", Color(0.7, 0.7, 0.5))
	_panel_bg.add_child(reward_lbl)

	# 關閉按鈕
	var close_btn = Button.new()
	close_btn.text = "✕"
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

	# 格子背景
	var bg = ColorRect.new()
	bg.name = "BG"
	bg.size = Vector2(CELL_SIZE, CELL_SIZE)
	bg.position = Vector2(x, y)
	bg.color = COLOR_EMPTY
	cell.add_child(bg)

	# 格子圖示
	var icon_lbl = Label.new()
	icon_lbl.name = "Icon"
	icon_lbl.text = "?"
	icon_lbl.position = Vector2(x + 20, y + 10)
	icon_lbl.add_theme_font_size_override("font_size", 24)
	icon_lbl.add_theme_color_override("font_color", Color(0.5, 0.5, 0.6))
	cell.add_child(icon_lbl)

	# 格子名稱
	var name_lbl = Label.new()
	name_lbl.name = "Name"
	name_lbl.text = ""
	name_lbl.position = Vector2(x + 4, y + 48)
	name_lbl.add_theme_font_size_override("font_size", 10)
	name_lbl.add_theme_color_override("font_color", Color(0.6, 0.6, 0.7))
	cell.add_child(name_lbl)

	# 填滿勾選標記（預設隱藏）
	var check = Label.new()
	check.name = "Check"
	check.text = "✓"
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

	# 更新副標題
	var subtitle = _panel_bg.get_node_or_null("Subtitle")
	if subtitle:
		subtitle.text = "今日進度：%d/9 格" % filled_count

	# 更新每個格子
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

		# 更新圖示
		var icon_lbl = cell.get_node_or_null("Icon")
		if icon_lbl:
			icon_lbl.text = icon
			icon_lbl.add_theme_color_override("font_color",
				Color(1.0, 1.0, 1.0) if filled else Color(0.5, 0.5, 0.6))

		# 更新名稱
		var name_lbl = cell.get_node_or_null("Name")
		if name_lbl:
			name_lbl.text = name_text

		# 更新背景顏色
		var bg = cell.get_node_or_null("BG")
		if bg:
			bg.color = COLOR_FILLED if filled else COLOR_EMPTY

		# 更新勾選標記
		var check = cell.get_node_or_null("Check")
		if check:
			check.visible = filled

func _on_treasure_map_line(data: Dictionary) -> void:
	var line_type = data.get("line_type", "")
	var reward = data.get("reward", 0)
	var message = data.get("message", "完成一條線！")

	# 顯示完成通知
	_show_reward_popup(message, reward, Color(1.0, 0.85, 0.2))

	# 高亮完成的行/列/對角線
	_highlight_line(line_type)

func _on_treasure_map_full(data: Dictionary) -> void:
	var reward = data.get("reward", 0)
	var message = data.get("message", "傳說寶藏！")

	# 顯示大獎通知（金色閃光）
	_show_reward_popup(message, reward, Color(1.0, 0.7, 0.0))
	_show_full_map_effect()

func _highlight_line(line_type: String) -> void:
	# 依行/列/對角線類型高亮對應格子
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
	popup.text = "%s\n+%d 金幣" % [message, reward]
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
	# 整張地圖金色閃光
	for cell in _cell_nodes:
		var bg = cell.get_node_or_null("BG")
		if bg:
			var tween = create_tween().set_loops(5)
			tween.tween_property(bg, "color", Color(1.0, 0.85, 0.0, 1.0), 0.1)
			tween.tween_property(bg, "color", COLOR_FILLED, 0.1)

func show_panel() -> void:
	visible = true
	_is_open = true
	# 請求最新地圖狀態
	if NetworkManager.has_method("send_get_treasure_map"):
		NetworkManager.send_get_treasure_map()

func hide_panel() -> void:
	visible = false
	_is_open = false
