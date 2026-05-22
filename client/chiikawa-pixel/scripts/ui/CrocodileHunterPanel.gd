## CrocodileHunterPanel.gd — 巨型鱷魚獵食 UI 面板（DAY-188）
## 業界依據：JILI Mega Fishing「giant crocodiles awaken to hunt fish on the fish farm
## to accumulate big prizes」
## 顯示鱷魚出現橫幅、每次獵食動畫、累積獎池、擊破結算彈窗
## Phase: croc_appear → croc_hunt(×N) / croc_miss → croc_killed / croc_leave
extends CanvasLayer

# ---- 常數 ----
const PANEL_COLOR_BG      := Color(0.0, 0.12, 0.05, 0.92)  # 深綠（鱷魚感）
const PANEL_COLOR_GREEN   := Color(0.1, 0.85, 0.3, 1.0)    # 亮綠（鱷魚出現）
const PANEL_COLOR_ORANGE  := Color(1.0, 0.55, 0.0, 1.0)    # 橙色（獵食中）
const PANEL_COLOR_GOLD    := Color(1.0, 0.85, 0.0, 1.0)    # 金色（擊破大獎）
const PANEL_COLOR_RED     := Color(1.0, 0.2, 0.1, 1.0)     # 紅色（鱷魚離去）
const PANEL_COLOR_WHITE   := Color(1.0, 1.0, 1.0, 1.0)

# ---- 節點引用 ----
var _banner_container : Control
var _banner_label     : Label
var _hunt_counter     : Label   # 獵食次數計數器
var _pool_label       : Label   # 累積獎池顯示
var _result_panel     : Control
var _result_label     : Label
var _flash_overlay    : ColorRect

# ---- 狀態 ----
var _hunt_count       : int = 0
var _max_hunts        : int = 8
var _total_pool       : int = 0
var _instance_id      : String = ""

func _ready() -> void:
	layer = 57  # 比 ChainBombPanel(58) 低一層
	_build_ui()
	hide()

func _build_ui() -> void:
	# 全螢幕閃光 overlay（深綠色）
	_flash_overlay = ColorRect.new()
	_flash_overlay.color = Color(0.0, 0.6, 0.2, 0.0)
	_flash_overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_flash_overlay.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_flash_overlay)

	# 頂部橫幅
	_banner_container = PanelContainer.new()
	_banner_container.set_anchors_preset(Control.PRESET_TOP_WIDE)
	_banner_container.offset_top = 8
	_banner_container.offset_bottom = 60
	_banner_container.offset_left = 80
	_banner_container.offset_right = -80
	var banner_style := StyleBoxFlat.new()
	banner_style.bg_color = Color(0.0, 0.15, 0.06, 0.92)
	banner_style.corner_radius_top_left = 8
	banner_style.corner_radius_top_right = 8
	banner_style.corner_radius_bottom_left = 8
	banner_style.corner_radius_bottom_right = 8
	banner_style.border_width_bottom = 2
	banner_style.border_color = PANEL_COLOR_GREEN
	_banner_container.add_theme_stylebox_override("panel", banner_style)
	add_child(_banner_container)

	var banner_vbox := VBoxContainer.new()
	banner_vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	_banner_container.add_child(banner_vbox)

	_banner_label = Label.new()
	_banner_label.text = "🐊 巨型鱷魚出現！"
	_banner_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner_label.add_theme_color_override("font_color", PANEL_COLOR_GREEN)
	_banner_label.add_theme_font_size_override("font_size", 20)
	banner_vbox.add_child(_banner_label)

	# 獵食計數器 + 獎池（橫排）
	var counter_hbox := HBoxContainer.new()
	counter_hbox.alignment = BoxContainer.ALIGNMENT_CENTER
	banner_vbox.add_child(counter_hbox)

	_hunt_counter = Label.new()
	_hunt_counter.text = "獵食：0/8"
	_hunt_counter.add_theme_color_override("font_color", PANEL_COLOR_ORANGE)
	_hunt_counter.add_theme_font_size_override("font_size", 14)
	counter_hbox.add_child(_hunt_counter)

	var spacer := Label.new()
	spacer.text = "  |  "
	spacer.add_theme_color_override("font_color", PANEL_COLOR_WHITE)
	counter_hbox.add_child(spacer)

	_pool_label = Label.new()
	_pool_label.text = "獎池：0"
	_pool_label.add_theme_color_override("font_color", PANEL_COLOR_GOLD)
	_pool_label.add_theme_font_size_override("font_size", 14)
	counter_hbox.add_child(_pool_label)

	# 結果彈窗（右側滑入）
	_result_panel = PanelContainer.new()
	_result_panel.set_anchors_preset(Control.PRESET_CENTER_RIGHT)
	_result_panel.offset_right = -16
	_result_panel.offset_left = -320
	_result_panel.offset_top = -100
	_result_panel.offset_bottom = 100
	var result_style := StyleBoxFlat.new()
	result_style.bg_color = Color(0.0, 0.15, 0.06, 0.95)
	result_style.corner_radius_top_left = 12
	result_style.corner_radius_top_right = 12
	result_style.corner_radius_bottom_left = 12
	result_style.corner_radius_bottom_right = 12
	result_style.border_width_left = 3
	result_style.border_color = PANEL_COLOR_GOLD
	_result_panel.add_theme_stylebox_override("panel", result_style)
	_result_panel.modulate.a = 0.0
	add_child(_result_panel)

	_result_label = Label.new()
	_result_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_result_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_result_label.add_theme_color_override("font_color", PANEL_COLOR_GOLD)
	_result_label.add_theme_font_size_override("font_size", 16)
	_result_label.autowrap_mode = TextServer.AUTOWRAP_WORD_SMART
	_result_panel.add_child(_result_label)

# ---- 公開 API ----

## handle 處理 crocodile_hunter 訊息
func handle(payload: Dictionary) -> void:
	var phase : String = payload.get("phase", "")
	match phase:
		"croc_appear":
			_on_croc_appear(payload)
		"croc_hunt":
			_on_croc_hunt(payload)
		"croc_miss":
			_on_croc_miss(payload)
		"croc_killed":
			_on_croc_killed(payload)
		"croc_leave":
			_on_croc_leave(payload)

# ---- 私有處理函數 ----

func _on_croc_appear(payload: Dictionary) -> void:
	_instance_id = payload.get("instance_id", "")
	_max_hunts   = payload.get("max_hunts", 8)
	_hunt_count  = 0
	_total_pool  = 0

	show()
	_banner_label.text = "🐊 巨型鱷魚出現！"
	_banner_label.add_theme_color_override("font_color", PANEL_COLOR_GREEN)
	_hunt_counter.text = "獵食：0/%d" % _max_hunts
	_pool_label.text = "獎池：0"

	# 橫幅從上方滑入
	_banner_container.offset_top = -60
	_banner_container.offset_bottom = -12
	var tween := create_tween()
	tween.tween_property(_banner_container, "offset_top", 8.0, 0.35).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)
	tween.parallel().tween_property(_banner_container, "offset_bottom", 60.0, 0.35).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)

	# 綠色閃光
	_flash_green(0.45)

func _on_croc_hunt(payload: Dictionary) -> void:
	_hunt_count = payload.get("hunt_index", _hunt_count + 1)
	_total_pool = payload.get("total_pool", _total_pool)
	var target_name : String = payload.get("target_name", "目標")
	var target_mult : float = payload.get("target_mult", 1.0)
	var hunt_reward : int   = payload.get("hunt_reward", 0)
	var target_x    : float = payload.get("target_x", 640.0)
	var target_y    : float = payload.get("target_y", 360.0)

	# 更新計數器和獎池
	_hunt_counter.text = "獵食：%d/%d" % [_hunt_count, _max_hunts]
	_pool_label.text = "獎池：%d" % _total_pool

	# 計數器彈跳動畫
	var tween := create_tween()
	tween.tween_property(_hunt_counter, "scale", Vector2(1.3, 1.3), 0.1)
	tween.tween_property(_hunt_counter, "scale", Vector2(1.0, 1.0), 0.15)

	# 在目標位置顯示「🐊 +獎勵」浮動文字
	_spawn_hunt_text(target_x, target_y, "+%d" % hunt_reward, PANEL_COLOR_ORANGE)

	# 橙色小閃光
	_flash_color(Color(1.0, 0.55, 0.0, 0.3), 0.2)

	# ≥4 次獵食：橙色雙閃光
	if _hunt_count >= 4:
		_flash_color(Color(1.0, 0.55, 0.0, 0.5), 0.3)

	# ≥6 次獵食：金色強閃光
	if _hunt_count >= 6:
		_flash_color(Color(1.0, 0.85, 0.0, 0.6), 0.4)

func _on_croc_miss(payload: Dictionary) -> void:
	var target_x : float = payload.get("target_x", 640.0)
	var target_y : float = payload.get("target_y", 360.0)
	# 顯示「💨 逃脫！」浮動文字
	_spawn_hunt_text(target_x, target_y, "💨 逃脫！", PANEL_COLOR_WHITE)

func _on_croc_killed(payload: Dictionary) -> void:
	var killer_name : String = payload.get("killer_name", "玩家")
	var hunt_count  : int    = payload.get("hunt_count", _hunt_count)
	var total_pool  : int    = payload.get("total_pool", _total_pool)
	var pool_bonus  : int    = payload.get("pool_bonus", 0)
	var base_reward : int    = payload.get("base_reward", 0)
	var total_reward: int    = payload.get("total_reward", 0)

	# 金色強閃光
	_flash_color(Color(1.0, 0.85, 0.0, 0.7), 0.5)

	# 結果彈窗
	_result_label.text = "🐊 %s 擊破鱷魚！\n\n獵食次數：%d\n累積獎池：%d\n獎池加成：+%d\n基礎獎勵：%d\n\n🏆 總獎勵：%d" % [
		killer_name, hunt_count, total_pool, pool_bonus, base_reward, total_reward
	]
	_result_label.add_theme_color_override("font_color", PANEL_COLOR_GOLD)

	# 右側滑入
	_result_panel.offset_right = 400
	_result_panel.offset_left = 80
	var tween := create_tween()
	tween.tween_property(_result_panel, "modulate:a", 1.0, 0.3)
	tween.parallel().tween_property(_result_panel, "offset_right", -16.0, 0.4).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)
	tween.parallel().tween_property(_result_panel, "offset_left", -320.0, 0.4).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)

	# 5 秒後淡出
	tween.tween_interval(5.0)
	tween.tween_callback(_hide_all)

func _on_croc_leave(payload: Dictionary) -> void:
	var hunt_count : int = payload.get("hunt_count", _hunt_count)
	var total_pool : int = payload.get("total_pool", _total_pool)

	# 紅色閃光（鱷魚離去，獎池未被領取）
	_flash_color(Color(1.0, 0.2, 0.1, 0.5), 0.4)

	# 結果彈窗（紅色主題）
	_result_label.text = "🐊 巨型鱷魚離去！\n\n獵食次數：%d\n未領取獎池：%d\n\n下次要快點擊破牠！" % [hunt_count, total_pool]
	_result_label.add_theme_color_override("font_color", PANEL_COLOR_RED)

	var tween := create_tween()
	tween.tween_property(_result_panel, "modulate:a", 1.0, 0.3)
	tween.tween_interval(4.0)
	tween.tween_callback(_hide_all)

# ---- 輔助函數 ----

func _flash_green(duration: float) -> void:
	_flash_color(Color(0.0, 0.6, 0.2, 0.5), duration)

func _flash_color(color: Color, duration: float) -> void:
	_flash_overlay.color = color
	var tween := create_tween()
	tween.tween_property(_flash_overlay, "color:a", 0.0, duration)

func _spawn_hunt_text(x: float, y: float, text: String, color: Color) -> void:
	var label := Label.new()
	label.text = text
	label.add_theme_color_override("font_color", color)
	label.add_theme_font_size_override("font_size", 18)
	label.position = Vector2(x - 40, y - 20)
	add_child(label)

	var tween := label.create_tween()
	tween.tween_property(label, "position:y", y - 70, 0.8).set_trans(Tween.TRANS_QUAD).set_ease(Tween.EASE_OUT)
	tween.parallel().tween_property(label, "modulate:a", 0.0, 0.8)
	tween.tween_callback(label.queue_free)

func _hide_all() -> void:
	var tween := create_tween()
	tween.tween_property(_banner_container, "modulate:a", 0.0, 0.4)
	tween.parallel().tween_property(_result_panel, "modulate:a", 0.0, 0.4)
	tween.tween_callback(func():
		hide()
		_banner_container.modulate.a = 1.0
		_result_panel.modulate.a = 0.0
		_hunt_count = 0
		_total_pool = 0
	)
