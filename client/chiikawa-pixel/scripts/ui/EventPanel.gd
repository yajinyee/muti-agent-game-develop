## EventPanel.gd ????瘣餃??Ｘ嚗AY-079嚗?
## 憿舐內?嗅???瘣餃??閮????牧??
## 雿蔭嚗opBar 銝?撣賊?憿舐內嚗?瘣餃???鈭殷?
extends Node2D

# ---- 撣豢 ----
const PANEL_WIDTH  := 260
const PANEL_HEIGHT := 80

# 瘣餃?憿?憿
const EVENT_COLORS := {
	"golden_hour":  Color(1.00, 0.85, 0.10),  # 暺?
	"fish_frenzy":  Color(0.00, 0.75, 1.00),  # ?
	"lucky_moment": Color(0.00, 1.00, 0.50),  # 蝬
	"none":         Color(0.40, 0.40, 0.40),  # ?啗
}

# ---- 蝭暺???----
var _pixel_font: Font = null
var _is_open: bool = false
var _toggle_btn: Button = null
var _panel_bg: ColorRect = null
var _event_icon_label: Label = null
var _event_name_label: Label = null
var _event_desc_label: Label = null
var _timer_label: Label = null
var _effect_label: Label = null
var _glow_rect: ColorRect = null

# ---- 瘣餃?鞈? ----
var _event_data: Dictionary = {
	"type": "none",
	"name": "",
	"description": "",
	"icon": "??,"
	"color": "#666666",
	"is_active": false,
	"end_at": 0,
	"time_left": 0.0,
	"reward_mult": 1.0,
	"spawn_mult": 1.0,
	"kill_chance_add": 0.0
}

# ---- ????----
func _ready() -> void:
	if ResourceLoader.exists("res://assets/fonts/pixel8.fnt"):
		_pixel_font = load("res://assets/fonts/pixel8.fnt")
	_build_toggle_btn()
	_build_panel()
	_connect_signals()

func setup(font: Font) -> void:
	if font:
		_pixel_font = font

func _process(_delta: float) -> void:
	# 瘥??湔?閮?
	if _event_data.get("is_active", false) and _timer_label:
		var end_at_ms: int = _event_data.get("end_at", 0)
		if end_at_ms > 0:
			var now_ms := int(Time.get_unix_time_from_system() * 1000)
			var remaining_ms := end_at_ms - now_ms
			if remaining_ms > 0:
				var secs := remaining_ms / 1000
				var mins := secs / 60
				var s := secs % 60
				_timer_label.text = "??%02d:%02d" % [mins, s]
			else:
				_timer_label.text = "??蝯?銝?.."

## 撱箇?????嚗opBar 銝?
func _build_toggle_btn() -> void:
	_toggle_btn = Button.new()
	_toggle_btn.text = "??"
	_toggle_btn.size = Vector2(32, 24)
	_toggle_btn.position = Vector2(0, 0)
	_toggle_btn.flat = true
	_toggle_btn.tooltip_text = "??瘣餃?"
	if _pixel_font:
		_toggle_btn.add_theme_font_override("font", _pixel_font)
		_toggle_btn.add_theme_font_size_override("font_size", 14)
	add_child(_toggle_btn)

## 撱箇?銝駁?選??身?梯?嚗?
func _build_panel() -> void:
	_panel_bg = ColorRect.new()
	_panel_bg.position = Vector2(-PANEL_WIDTH + 32, 28)
	_panel_bg.size = Vector2(PANEL_WIDTH, PANEL_HEIGHT)
	_panel_bg.color = Color(0.03, 0.05, 0.15, 0.93)
	_panel_bg.visible = false
	add_child(_panel_bg)

	# ????嚗暑??憿舐內嚗?
	_glow_rect = ColorRect.new()
	_glow_rect.position = Vector2(-1, -1)
	_glow_rect.size = Vector2(PANEL_WIDTH + 2, PANEL_HEIGHT + 2)
	_glow_rect.color = Color(1.0, 0.85, 0.1, 0.0)
	_glow_rect.z_index = -1
	_panel_bg.add_child(_glow_rect)

	# 瘣餃??內 + ?迂
	_event_icon_label = Label.new()
	_event_icon_label.position = Vector2(8, 4)
	_event_icon_label.text = "??"
	_event_icon_label.add_theme_color_override("font_color", Color(0.6, 0.6, 0.6))
	if _pixel_font:
		_event_icon_label.add_theme_font_override("font", _pixel_font)
		_event_icon_label.add_theme_font_size_override("font_size", 16)
	_panel_bg.add_child(_event_icon_label)

	_event_name_label = Label.new()
	_event_name_label.position = Vector2(32, 6)
	_event_name_label.text = "?桀??⊥暑??"
	_event_name_label.add_theme_color_override("font_color", Color(0.6, 0.6, 0.6))
	if _pixel_font:
		_event_name_label.add_theme_font_override("font", _pixel_font)
		_event_name_label.add_theme_font_size_override("font_size", 12)
	_panel_bg.add_child(_event_name_label)

	# ?閮?
	_timer_label = Label.new()
	_timer_label.position = Vector2(PANEL_WIDTH - 80, 6)
	_timer_label.text = ""
	_timer_label.add_theme_color_override("font_color", Color(1.0, 0.9, 0.5))
	if _pixel_font:
		_timer_label.add_theme_font_override("font", _pixel_font)
		_timer_label.add_theme_font_size_override("font_size", 11)
	_panel_bg.add_child(_timer_label)

	# 瘣餃??膩
	_event_desc_label = Label.new()
	_event_desc_label.position = Vector2(8, 26)
	_event_desc_label.text = "蝑?銝??暑??.."
	_event_desc_label.add_theme_color_override("font_color", Color(0.7, 0.7, 0.7))
	if _pixel_font:
		_event_desc_label.add_theme_font_override("font", _pixel_font)
		_event_desc_label.add_theme_font_size_override("font_size", 10)
	_panel_bg.add_child(_event_desc_label)

	# ??隤芣?
	_effect_label = Label.new()
	_effect_label.position = Vector2(8, 44)
	_effect_label.text = ""
	_effect_label.add_theme_color_override("font_color", Color(0.9, 0.9, 0.5))
	if _pixel_font:
		_effect_label.add_theme_font_override("font", _pixel_font)
		_effect_label.add_theme_font_size_override("font_size", 10)
	_panel_bg.add_child(_effect_label)

## ??閮?
func _connect_signals() -> void:
	if _toggle_btn:
		_toggle_btn.pressed.connect(_on_toggle_pressed)
	if GameManager:
		if GameManager.has_signal("event_updated"):
			GameManager.event_updated.connect(_on_event_updated)

## ??/撅??Ｘ
func _on_toggle_pressed() -> void:
	_is_open = !_is_open
	if _panel_bg:
		_panel_bg.visible = _is_open

## ?湔瘣餃?鞈?
func _on_event_updated(data: Dictionary) -> void:
	var was_active: bool = _event_data.get("is_active", false)
	_event_data = data
	_refresh_ui()
	# ?交暑????嚗＊蝷粹
	var is_active: bool = data.get("is_active", false)
	if is_active and not was_active:
		_show_event_start_popup(data)

## ?瑟 UI
func _refresh_ui() -> void:
	if not _panel_bg:
		return

	var event_type: String = _event_data.get("type", "none")
	var event_name: String = _event_data.get("name", "")
	var event_icon: String = _event_data.get("icon", "??)"
	var event_desc: String = _event_data.get("description", "")
	var is_active: bool = _event_data.get("is_active", false)
	var reward_mult: float = _event_data.get("reward_mult", 1.0)
	var spawn_mult: float = _event_data.get("spawn_mult", 1.0)
	var kill_add: float = _event_data.get("kill_chance_add", 0.0)

	var event_color: Color = EVENT_COLORS.get(event_type, Color(0.4, 0.4, 0.4))

	# ?湔??
	if _toggle_btn:
		if is_active:
			_toggle_btn.text = event_icon
			_toggle_btn.modulate = event_color
		else:
			_toggle_btn.text = "??"
			_toggle_btn.modulate = Color(0.5, 0.5, 0.5)

	# ?湔?內
	if _event_icon_label:
		_event_icon_label.text = event_icon
		_event_icon_label.add_theme_color_override("font_color", event_color)

	# ?湔?迂
	if _event_name_label:
		if is_active:
			_event_name_label.text = event_name
			_event_name_label.add_theme_color_override("font_color", event_color)
		else:
			_event_name_label.text = "?桀??⊥暑??"
			_event_name_label.add_theme_color_override("font_color", Color(0.5, 0.5, 0.5))

	# ?湔?膩
	if _event_desc_label:
		if is_active:
			_event_desc_label.text = event_desc
			_event_desc_label.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8))
		else:
			_event_desc_label.text = "蝑?銝??暑??.."
			_event_desc_label.add_theme_color_override("font_color", Color(0.5, 0.5, 0.5))

	# ?湔??隤芣?
	if _effect_label:
		var effects := []
		if reward_mult > 1.0:
			effects.append("? ?%.1f" % reward_mult)
		if spawn_mult > 1.0:
			effects.append("?格? ?%.1f" % spawn_mult)
		if kill_add > 0:
			effects.append("???+%.0f%%" % (kill_add * 100))
		if effects.size() > 0:
			_effect_label.text = "??嚗? + "  ".join(effects)"
			_effect_label.add_theme_color_override("font_color", event_color)
		else:
			_effect_label.text = ""

	# ????
	if _glow_rect:
		if is_active:
			_glow_rect.color = Color(event_color.r, event_color.g, event_color.b, 0.15)
		else:
			_glow_rect.color = Color(0, 0, 0, 0)

	# 閮??剁??⊥暑??皜征嚗?
	if _timer_label and not is_active:
		_timer_label.text = ""

## 憿舐內瘣餃???敶?
func _show_event_start_popup(data: Dictionary) -> void:
	var event_name: String = data.get("name", "")
	var event_icon: String = data.get("icon", "??)"
	var event_desc: String = data.get("description", "")
	var event_type: String = data.get("type", "none")
	var event_color: Color = EVENT_COLORS.get(event_type, Color(1.0, 1.0, 1.0))

	var canvas := get_viewport().get_canvas_layer_node(1) if get_viewport() else null
	var popup_parent := canvas if canvas else get_parent()
	if not is_instance_valid(popup_parent):
		return

	var popup := ColorRect.new()
	popup.size = Vector2(300, 70)
	popup.position = Vector2(
		(get_viewport().get_visible_rect().size.x - 300) / 2.0,
		get_viewport().get_visible_rect().size.y * 0.3
	)
	popup.color = Color(0.03, 0.05, 0.15, 0.95)
	popup_parent.add_child(popup)

	var lbl := Label.new()
	lbl.position = Vector2(8, 6)
	lbl.text = "%s ??瘣餃???嚗? % event_icon"
	lbl.add_theme_color_override("font_color", event_color)
	if _pixel_font:
		lbl.add_theme_font_override("font", _pixel_font)
		lbl.add_theme_font_size_override("font_size", 14)
	popup.add_child(lbl)

	var lbl2 := Label.new()
	lbl2.position = Vector2(8, 26)
	lbl2.text = event_name
	lbl2.add_theme_color_override("font_color", Color(1.0, 1.0, 1.0))
	if _pixel_font:
		lbl2.add_theme_font_override("font", _pixel_font)
		lbl2.add_theme_font_size_override("font_size", 12)
	popup.add_child(lbl2)

	var lbl3 := Label.new()
	lbl3.position = Vector2(8, 44)
	lbl3.text = event_desc
	lbl3.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8))
	if _pixel_font:
		lbl3.add_theme_font_override("font", _pixel_font)
		lbl3.add_theme_font_size_override("font_size", 10)
	popup.add_child(lbl3)

	# ?嚗??????? ??瘛∪
	var tween := popup.create_tween()
	popup.modulate.a = 0.0
	tween.tween_property(popup, "modulate:a", 1.0, 0.4)
	tween.tween_interval(3.5)
	tween.tween_property(popup, "modulate:a", 0.0, 0.5)
	tween.tween_callback(popup.queue_free)
