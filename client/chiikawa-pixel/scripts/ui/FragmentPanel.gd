## FragmentPanel.gd
## 碎片收集大獎系統 UI（DAY-116）
## 顯示碎片收集槽（左下角），掉落動畫，集齊大獎特效
## 業界依據：Hidden Treasure Unlocks — 碎片收集讓玩家留存率提升 28%

extends Control

# 碎片類型顏色
const BRONZE_COLOR = Color(0.804, 0.498, 0.196)  # #CD7F32
const SILVER_COLOR = Color(0.753, 0.753, 0.753)  # #C0C0C0
const GOLD_COLOR   = Color(1.0, 0.843, 0.0)      # #FFD700

const REQUIRED = 5  # 集齊需要數量

# 碎片狀態
var _bronze: int = 0
var _silver: int = 0
var _gold: int = 0

# UI 節點
var _panel: Control = null
var _bronze_bar: Control = null
var _silver_bar: Control = null
var _gold_bar: Control = null
var _complete_overlay: Control = null

func _ready() -> void:
	_build_ui()

func _build_ui() -> void:
	# 主面板（左下角，z_index=62）
	_panel = Control.new()
	_panel.name = "FragmentPanel"
	_panel.position = Vector2(8, 560)
	_panel.size = Vector2(180, 110)
	_panel.z_index = 62
	add_child(_panel)

	# 背景
	var bg = ColorRect.new()
	bg.size = Vector2(180, 110)
	bg.color = Color(0.03, 0.06, 0.18, 0.82)
	_panel.add_child(bg)

	# 邊框
	var border = ColorRect.new()
	border.size = Vector2(180, 2)
	border.color = Color(0.9, 0.75, 0.2, 0.7)
	_panel.add_child(border)

	# 標題
	var title = Label.new()
	title.text = "🧩 碎片收集"
	title.position = Vector2(8, 4)
	title.add_theme_font_size_override("font_size", 11)
	title.modulate = Color(0.9, 0.85, 0.6)
	_panel.add_child(title)

	# 三種碎片進度條
	_bronze_bar = _create_fragment_bar("銅", BRONZE_COLOR, 0)
	_silver_bar = _create_fragment_bar("銀", SILVER_COLOR, 1)
	_gold_bar   = _create_fragment_bar("金", GOLD_COLOR,   2)

	# 集齊大獎覆蓋層（預設隱藏）
	_complete_overlay = Control.new()
	_complete_overlay.name = "CompleteOverlay"
	_complete_overlay.visible = false
	_complete_overlay.z_index = 100
	add_child(_complete_overlay)

func _create_fragment_bar(label_text: String, color: Color, row: int) -> Control:
	var container = Control.new()
	container.position = Vector2(8, 24 + row * 28)
	container.size = Vector2(164, 24)
	_panel.add_child(container)

	# 標籤
	var lbl = Label.new()
	lbl.text = label_text
	lbl.position = Vector2(0, 2)
	lbl.add_theme_font_size_override("font_size", 10)
	lbl.modulate = color
	container.add_child(lbl)

	# 進度背景
	var bar_bg = ColorRect.new()
	bar_bg.name = "BarBG"
	bar_bg.position = Vector2(20, 4)
	bar_bg.size = Vector2(100, 14)
	bar_bg.color = Color(0.1, 0.1, 0.1, 0.8)
	container.add_child(bar_bg)

	# 進度填充
	var bar_fill = ColorRect.new()
	bar_fill.name = "BarFill"
	bar_fill.position = Vector2(20, 4)
	bar_fill.size = Vector2(0, 14)
	bar_fill.color = color
	container.add_child(bar_fill)

	# 計數標籤
	var count_lbl = Label.new()
	count_lbl.name = "CountLabel"
	count_lbl.text = "0/5"
	count_lbl.position = Vector2(124, 2)
	count_lbl.add_theme_font_size_override("font_size", 10)
	count_lbl.modulate = Color(0.8, 0.8, 0.8)
	container.add_child(count_lbl)

	return container

## 更新碎片狀態（玩家上線時或查詢回應時呼叫）
func update_status(data: Dictionary) -> void:
	_bronze = data.get("bronze", 0)
	_silver = data.get("silver", 0)
	_gold   = data.get("gold", 0)
	_refresh_bars()

## 處理碎片掉落通知
func on_fragment_drop(data: Dictionary) -> void:
	var ftype = data.get("fragment_type", "bronze")
	var new_count = data.get("new_count", 0)
	var drop_x = data.get("drop_x", 640.0)
	var drop_y = data.get("drop_y", 360.0)

	# 更新本地計數
	match ftype:
		"bronze": _bronze = new_count
		"silver": _silver = new_count
		"gold":   _gold   = new_count

	_refresh_bars()

	# 掉落動畫：從目標位置飛向收集槽
	var color = _get_fragment_color(ftype)
	_spawn_drop_animation(drop_x, drop_y, ftype, color)

## 處理集齊大獎通知（廣播給所有玩家）
func on_fragment_complete(data: Dictionary) -> void:
	var ftype = data.get("fragment_type", "gold")
	var label = data.get("label", "碎片大獎")
	var reward = data.get("reward", 0)
	var is_self = data.get("is_self", false)
	var display_name = data.get("display_name", "玩家")
	var color = _get_fragment_color(ftype)

	# 重置本地計數
	if is_self:
		match ftype:
			"bronze": _bronze = 0
			"silver": _silver = 0
			"gold":   _gold   = 0
		_refresh_bars()

	# 顯示大獎特效
	_show_complete_effect(display_name, label, reward, color, is_self)

## 刷新進度條顯示
func _refresh_bars() -> void:
	_update_bar(_bronze_bar, _bronze)
	_update_bar(_silver_bar, _silver)
	_update_bar(_gold_bar, _gold)

func _update_bar(container: Control, count: int) -> void:
	if not is_instance_valid(container):
		return
	var bar_fill = container.get_node_or_null("BarFill")
	var count_lbl = container.get_node_or_null("CountLabel")
	if is_instance_valid(bar_fill):
		var pct = float(count) / float(REQUIRED)
		bar_fill.size.x = 100.0 * pct
		# 快滿時閃爍
		if count >= REQUIRED - 1 and count > 0:
			var tween = create_tween().set_loops(3)
			tween.tween_property(bar_fill, "modulate:a", 0.4, 0.15)
			tween.tween_property(bar_fill, "modulate:a", 1.0, 0.15)
	if is_instance_valid(count_lbl):
		count_lbl.text = "%d/%d" % [count, REQUIRED]

## 掉落動畫：碎片從目標位置飛向收集槽
func _spawn_drop_animation(from_x: float, from_y: float, ftype: String, color: Color) -> void:
	# 目標位置（收集槽中心）
	var row = 0
	match ftype:
		"silver": row = 1
		"gold":   row = 2
	var target_pos = Vector2(98, 570 + row * 28)  # 收集槽位置

	# 建立碎片粒子
	var gem = ColorRect.new()
	gem.size = Vector2(10, 10)
	gem.position = Vector2(from_x - 5, from_y - 5)
	gem.color = color
	gem.z_index = 80
	add_child(gem)

	# 飛行動畫
	var tween = create_tween()
	tween.tween_property(gem, "position", target_pos, 0.6).set_trans(Tween.TRANS_QUAD).set_ease(Tween.EASE_IN)
	tween.tween_callback(func():
		if is_instance_valid(gem):
			gem.queue_free()
		# 收集槽閃光
		_flash_bar_row(row, color)
	)

## 收集槽閃光效果
func _flash_bar_row(row: int, color: Color) -> void:
	var container: Control
	match row:
		0: container = _bronze_bar
		1: container = _silver_bar
		2: container = _gold_bar
	if not is_instance_valid(container):
		return
	var tween = create_tween()
	tween.tween_property(container, "modulate", Color(2.0, 2.0, 2.0), 0.08)
	tween.tween_property(container, "modulate", Color.WHITE, 0.15)

## 集齊大獎特效
func _show_complete_effect(display_name: String, label: String, reward: int, color: Color, is_self: bool) -> void:
	# 全螢幕覆蓋層
	var overlay = ColorRect.new()
	overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	overlay.color = Color(color.r * 0.3, color.g * 0.3, color.b * 0.3, 0.0)
	overlay.z_index = 95
	add_child(overlay)

	# 淡入
	var tween = create_tween()
	tween.tween_property(overlay, "color:a", 0.6, 0.2)

	# 大獎面板
	var panel = Control.new()
	panel.position = Vector2(290, 200)
	panel.size = Vector2(700, 200)
	panel.z_index = 96
	panel.modulate.a = 0.0
	add_child(panel)

	# 面板背景
	var bg = ColorRect.new()
	bg.size = Vector2(700, 200)
	bg.color = Color(0.02, 0.04, 0.12, 0.95)
	panel.add_child(bg)

	# 頂部彩色邊條
	var top_bar = ColorRect.new()
	top_bar.size = Vector2(700, 5)
	top_bar.color = color
	panel.add_child(top_bar)

	# 大獎標題
	var title_lbl = Label.new()
	title_lbl.text = "🧩 " + label + " 🧩"
	title_lbl.position = Vector2(0, 20)
	title_lbl.size = Vector2(700, 50)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title_lbl.add_theme_font_size_override("font_size", 36)
	title_lbl.modulate = color
	panel.add_child(title_lbl)

	# 玩家名稱
	var name_lbl = Label.new()
	name_lbl.text = display_name + " 集齊碎片！"
	name_lbl.position = Vector2(0, 75)
	name_lbl.size = Vector2(700, 40)
	name_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	name_lbl.add_theme_font_size_override("font_size", 20)
	name_lbl.modulate = Color(0.9, 0.9, 0.9)
	panel.add_child(name_lbl)

	# 獎勵金額
	var reward_lbl = Label.new()
	reward_lbl.text = "🪙 +" + str(reward) + " 金幣！"
	reward_lbl.position = Vector2(0, 120)
	reward_lbl.size = Vector2(700, 50)
	reward_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	reward_lbl.add_theme_font_size_override("font_size", 28)
	reward_lbl.modulate = Color(1.0, 0.9, 0.2)
	panel.add_child(reward_lbl)

	# 面板淡入
	var tween2 = create_tween()
	tween2.tween_property(panel, "modulate:a", 1.0, 0.25)

	# 金色粒子效果（is_self 時更強烈）
	if is_self or ftype_is_gold(label):
		_spawn_complete_particles(color)

	# 5 秒後淡出
	await get_tree().create_timer(5.0).timeout
	var tween3 = create_tween()
	tween3.tween_property(panel, "modulate:a", 0.0, 0.4)
	tween3.tween_property(overlay, "color:a", 0.0, 0.4)
	tween3.tween_callback(func():
		if is_instance_valid(panel): panel.queue_free()
		if is_instance_valid(overlay): overlay.queue_free()
	)

func ftype_is_gold(label: String) -> bool:
	return "金" in label

## 集齊粒子效果
func _spawn_complete_particles(color: Color) -> void:
	for i in range(20):
		var p = ColorRect.new()
		p.size = Vector2(8, 8)
		p.color = color
		p.position = Vector2(randf_range(200, 1080), randf_range(150, 500))
		p.z_index = 97
		add_child(p)
		var tween = create_tween()
		tween.tween_property(p, "position:y", p.position.y - randf_range(50, 150), randf_range(0.5, 1.2))
		tween.parallel().tween_property(p, "modulate:a", 0.0, randf_range(0.5, 1.2))
		tween.tween_callback(func():
			if is_instance_valid(p): p.queue_free()
		)

## 取得碎片顏色
func _get_fragment_color(ftype: String) -> Color:
	match ftype:
		"bronze": return BRONZE_COLOR
		"silver": return SILVER_COLOR
		"gold":   return GOLD_COLOR
	return Color.WHITE
