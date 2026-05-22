## TimeBombFishPanel.gd — 時間炸彈魚 UI 面板（DAY-189）
## 業界靈感：Ocean King 炸彈魚概念 + 倒數計時緊張感設計
## 顯示倒數計時、拆彈成功/爆炸動畫、結果彈窗
## Phase: bomb_appear → bomb_tick(×N) → bomb_defused / bomb_explode → bomb_result / defuse_end
extends CanvasLayer

# ---- 常數 ----
const COLOR_SAFE    := Color(0.1, 0.85, 0.3, 1.0)   # 安全綠（倒數充裕）
const COLOR_WARNING := Color(1.0, 0.75, 0.0, 1.0)   # 警告黃（倒數 5 秒）
const COLOR_DANGER  := Color(1.0, 0.2, 0.1, 1.0)    # 危險紅（倒數 3 秒）
const COLOR_DEFUSED := Color(0.1, 0.9, 0.4, 1.0)    # 拆彈成功綠
const COLOR_EXPLODE := Color(1.0, 0.4, 0.0, 1.0)    # 爆炸橙
const COLOR_GOLD    := Color(1.0, 0.85, 0.0, 1.0)   # 金色
const COLOR_WHITE   := Color(1.0, 1.0, 1.0, 1.0)

# ---- 節點引用 ----
var _banner_container : Control
var _banner_label     : Label
var _countdown_label  : Label   # 大型倒數數字
var _status_label     : Label   # 狀態提示
var _result_panel     : Control
var _result_label     : Label
var _flash_overlay    : ColorRect
var _defuse_bar       : ColorRect  # 拆彈加成進度條

# ---- 狀態 ----
var _countdown        : int = 10
var _instance_id      : String = ""
var _defuse_tween     : Tween = null

func _ready() -> void:
	layer = 56  # 比 CrocodileHunterPanel(57) 低一層
	_build_ui()
	hide()

func _build_ui() -> void:
	# 全螢幕閃光 overlay
	_flash_overlay = ColorRect.new()
	_flash_overlay.color = Color(1.0, 0.2, 0.1, 0.0)
	_flash_overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_flash_overlay.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_flash_overlay)

	# 頂部橫幅
	_banner_container = PanelContainer.new()
	_banner_container.set_anchors_preset(Control.PRESET_TOP_WIDE)
	_banner_container.offset_top = 8
	_banner_container.offset_bottom = 80
	_banner_container.offset_left = 80
	_banner_container.offset_right = -80
	var banner_style := StyleBoxFlat.new()
	banner_style.bg_color = Color(0.12, 0.02, 0.0, 0.92)
	banner_style.corner_radius_top_left = 8
	banner_style.corner_radius_top_right = 8
	banner_style.corner_radius_bottom_left = 8
	banner_style.corner_radius_bottom_right = 8
	banner_style.border_width_bottom = 2
	banner_style.border_color = COLOR_DANGER
	_banner_container.add_theme_stylebox_override("panel", banner_style)
	add_child(_banner_container)

	var banner_vbox := VBoxContainer.new()
	banner_vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	_banner_container.add_child(banner_vbox)

	_banner_label = Label.new()
	_banner_label.text = "💣 時間炸彈魚！"
	_banner_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner_label.add_theme_color_override("font_color", COLOR_DANGER)
	_banner_label.add_theme_font_size_override("font_size", 18)
	banner_vbox.add_child(_banner_label)

	# 大型倒數數字（中央顯示）
	_countdown_label = Label.new()
	_countdown_label.text = "10"
	_countdown_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_countdown_label.add_theme_color_override("font_color", COLOR_SAFE)
	_countdown_label.add_theme_font_size_override("font_size", 48)
	banner_vbox.add_child(_countdown_label)

	_status_label = Label.new()
	_status_label.text = "擊破可拆彈！否則全場爆炸！"
	_status_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_status_label.add_theme_color_override("font_color", COLOR_WHITE)
	_status_label.add_theme_font_size_override("font_size", 13)
	banner_vbox.add_child(_status_label)

	# 拆彈加成進度條（底部）
	_defuse_bar = ColorRect.new()
	_defuse_bar.color = COLOR_DEFUSED
	_defuse_bar.set_anchors_preset(Control.PRESET_BOTTOM_WIDE)
	_defuse_bar.offset_top = -8
	_defuse_bar.offset_bottom = 0
	_defuse_bar.modulate.a = 0.0
	add_child(_defuse_bar)

	# 結果彈窗（右側滑入）
	_result_panel = PanelContainer.new()
	_result_panel.set_anchors_preset(Control.PRESET_CENTER_RIGHT)
	_result_panel.offset_right = -16
	_result_panel.offset_left = -320
	_result_panel.offset_top = -110
	_result_panel.offset_bottom = 110
	var result_style := StyleBoxFlat.new()
	result_style.bg_color = Color(0.1, 0.05, 0.0, 0.95)
	result_style.corner_radius_top_left = 12
	result_style.corner_radius_top_right = 12
	result_style.corner_radius_bottom_left = 12
	result_style.corner_radius_bottom_right = 12
	result_style.border_width_left = 3
	result_style.border_color = COLOR_GOLD
	_result_panel.add_theme_stylebox_override("panel", result_style)
	_result_panel.modulate.a = 0.0
	add_child(_result_panel)

	_result_label = Label.new()
	_result_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_result_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_result_label.add_theme_color_override("font_color", COLOR_GOLD)
	_result_label.add_theme_font_size_override("font_size", 15)
	_result_label.autowrap_mode = TextServer.AUTOWRAP_WORD_SMART
	_result_panel.add_child(_result_label)

# ---- 公開 API ----

func handle(payload: Dictionary) -> void:
	var phase : String = payload.get("phase", "")
	match phase:
		"bomb_appear":
			_on_bomb_appear(payload)
		"bomb_tick":
			_on_bomb_tick(payload)
		"bomb_defused":
			_on_bomb_defused(payload)
		"bomb_explode":
			_on_bomb_explode(payload)
		"bomb_result":
			_on_bomb_result(payload)
		"defuse_end":
			_on_defuse_end(payload)

# ---- 私有處理函數 ----

func _on_bomb_appear(payload: Dictionary) -> void:
	_instance_id = payload.get("instance_id", "")
	_countdown    = payload.get("countdown", 10)

	show()
	_banner_label.text = "💣 時間炸彈魚！"
	_banner_label.add_theme_color_override("font_color", COLOR_DANGER)
	_countdown_label.text = str(_countdown)
	_countdown_label.add_theme_color_override("font_color", COLOR_SAFE)
	_status_label.text = "擊破可拆彈！否則全場爆炸！"
	_defuse_bar.modulate.a = 0.0

	# 橫幅從上方滑入
	_banner_container.offset_top = -90
	_banner_container.offset_bottom = -2
	var tween := create_tween()
	tween.tween_property(_banner_container, "offset_top", 8.0, 0.4).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)
	tween.parallel().tween_property(_banner_container, "offset_bottom", 80.0, 0.4).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)

	# 紅色閃光
	_flash_color(Color(1.0, 0.2, 0.1, 0.5), 0.5)

func _on_bomb_tick(payload: Dictionary) -> void:
	_countdown = payload.get("countdown", _countdown - 1)
	_countdown_label.text = str(_countdown)

	# 依剩餘時間改變顏色
	var color : Color
	if _countdown > 5:
		color = COLOR_SAFE
	elif _countdown > 3:
		color = COLOR_WARNING
	else:
		color = COLOR_DANGER

	_countdown_label.add_theme_color_override("font_color", color)

	# 倒數數字彈跳動畫
	var tween := create_tween()
	tween.tween_property(_countdown_label, "scale", Vector2(1.4, 1.4), 0.1)
	tween.tween_property(_countdown_label, "scale", Vector2(1.0, 1.0), 0.2)

	# ≤3 秒：紅色閃光 + 震動
	if _countdown <= 3:
		_flash_color(Color(1.0, 0.2, 0.1, 0.3), 0.3)
		_shake_banner()

func _on_bomb_defused(payload: Dictionary) -> void:
	var killer_name   : String = payload.get("killer_name", "玩家")
	var base_reward   : int    = payload.get("base_reward", 0)
	var bonus_pct     : int    = payload.get("bonus_pct", 25)
	var bonus_duration: int    = payload.get("bonus_duration", 15)

	# 停止倒數動畫
	_countdown_label.text = "✓"
	_countdown_label.add_theme_color_override("font_color", COLOR_DEFUSED)
	_banner_label.text = "💚 拆彈成功！"
	_banner_label.add_theme_color_override("font_color", COLOR_DEFUSED)
	_status_label.text = "全服 +%d%% 加成持續 %d 秒！" % [bonus_pct, bonus_duration]

	# 綠色強閃光
	_flash_color(Color(0.1, 0.9, 0.4, 0.7), 0.6)

	# 拆彈加成進度條（從左到右，持續 bonus_duration 秒）
	_defuse_bar.color = COLOR_DEFUSED
	_defuse_bar.modulate.a = 0.8
	_defuse_bar.set_anchors_preset(Control.PRESET_BOTTOM_WIDE)
	_defuse_bar.offset_top = -8
	_defuse_bar.offset_bottom = 0
	_defuse_bar.offset_right = 0  # 從全寬開始
	if _defuse_tween:
		_defuse_tween.kill()
	_defuse_tween = create_tween()
	_defuse_tween.tween_property(_defuse_bar, "offset_right", -get_viewport().get_visible_rect().size.x, float(bonus_duration))
	_defuse_tween.tween_callback(func(): _defuse_bar.modulate.a = 0.0)

	# 結果彈窗
	_result_label.text = "💚 %s 拆彈成功！\n\n擊破獎勵：%d\n全服加成：+%d%%\n持續時間：%d 秒" % [
		killer_name, base_reward, bonus_pct, bonus_duration
	]
	_result_label.add_theme_color_override("font_color", COLOR_DEFUSED)

	var tween := create_tween()
	tween.tween_property(_result_panel, "modulate:a", 1.0, 0.3)
	tween.parallel().tween_property(_result_panel, "offset_right", -16.0, 0.4).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)
	tween.parallel().tween_property(_result_panel, "offset_left", -320.0, 0.4).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)
	tween.tween_interval(5.0)
	tween.tween_property(_result_panel, "modulate:a", 0.0, 0.4)

func _on_bomb_explode(_payload: Dictionary) -> void:
	# 爆炸動畫：多次強閃光
	_countdown_label.text = "💥"
	_countdown_label.add_theme_color_override("font_color", COLOR_EXPLODE)
	_banner_label.text = "💥 炸彈爆炸！"
	_banner_label.add_theme_color_override("font_color", COLOR_EXPLODE)
	_status_label.text = "全場目標受到爆炸傷害！"

	# 橙紅色強閃光（3次）
	var tween := create_tween()
	tween.tween_property(_flash_overlay, "color:a", 0.8, 0.1)
	tween.tween_property(_flash_overlay, "color:a", 0.0, 0.15)
	tween.tween_property(_flash_overlay, "color:a", 0.6, 0.1)
	tween.tween_property(_flash_overlay, "color:a", 0.0, 0.15)
	tween.tween_property(_flash_overlay, "color:a", 0.4, 0.1)
	tween.tween_property(_flash_overlay, "color:a", 0.0, 0.2)
	_flash_overlay.color = COLOR_EXPLODE

func _on_bomb_result(payload: Dictionary) -> void:
	var kill_count   : int = payload.get("kill_count", 0)
	var total_reward : int = payload.get("total_reward", 0)

	# 結果彈窗
	_result_label.text = "💥 炸彈爆炸！\n\n擊破目標：%d\n全服共享：%d 金幣\n\n（每位玩家均分）" % [kill_count, total_reward]
	_result_label.add_theme_color_override("font_color", COLOR_EXPLODE)

	var tween := create_tween()
	tween.tween_property(_result_panel, "modulate:a", 1.0, 0.3)
	tween.parallel().tween_property(_result_panel, "offset_right", -16.0, 0.4).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)
	tween.parallel().tween_property(_result_panel, "offset_left", -320.0, 0.4).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)
	tween.tween_interval(4.0)
	tween.tween_callback(_hide_all)

func _on_defuse_end(_payload: Dictionary) -> void:
	# 拆彈加成結束
	if _defuse_tween:
		_defuse_tween.kill()
	var tween := create_tween()
	tween.tween_property(_defuse_bar, "modulate:a", 0.0, 0.5)
	tween.tween_callback(_hide_all)

# ---- 輔助函數 ----

func _flash_color(color: Color, duration: float) -> void:
	_flash_overlay.color = color
	var tween := create_tween()
	tween.tween_property(_flash_overlay, "color:a", 0.0, duration)

func _shake_banner() -> void:
	var original_x := _banner_container.offset_left
	var tween := create_tween()
	tween.tween_property(_banner_container, "offset_left", original_x + 6, 0.05)
	tween.tween_property(_banner_container, "offset_left", original_x - 6, 0.05)
	tween.tween_property(_banner_container, "offset_left", original_x + 4, 0.04)
	tween.tween_property(_banner_container, "offset_left", original_x, 0.04)

func _hide_all() -> void:
	var tween := create_tween()
	tween.tween_property(_banner_container, "modulate:a", 0.0, 0.4)
	tween.parallel().tween_property(_result_panel, "modulate:a", 0.0, 0.4)
	tween.tween_callback(func():
		hide()
		_banner_container.modulate.a = 1.0
		_result_panel.modulate.a = 0.0
		_countdown = 10
	)
