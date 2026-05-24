## ReferralPanel.gd ???刻蝣潮?選?DAY-082嚗?
## 憿舐內?拙振??衣Ⅳ嚗??拙振頛詨隞犖?刻蝣?
extends Node2D

const PANEL_WIDTH  := 360
const PANEL_HEIGHT := 280

var _font: FontFile
var _bg: ColorRect
var _my_code_label: Label
var _referral_count_label: Label
var _total_reward_label: Label
var _input_field: LineEdit
var _use_btn: Button
var _close_btn: Button
var _status_label: Label
var _is_visible := false

signal referral_closed

func setup(font: FontFile) -> void:
	_font = font
	_build_ui()
	_connect_signals()
	hide()

func _build_ui() -> void:
	var px := (1600 - PANEL_WIDTH) / 2.0
	var py := (900 - PANEL_HEIGHT) / 2.0

	_bg = ColorRect.new()
	_bg.color = Color(0.05, 0.08, 0.15, 0.95)
	_bg.size = Vector2(PANEL_WIDTH, PANEL_HEIGHT)
	_bg.position = Vector2(px, py)
	add_child(_bg)

	# 璅?
	var title = Label.new()
	title.text = "?? ?刻憟賢?"
	title.position = _bg.position + Vector2(16, 12)
	title.add_theme_color_override("font_color", Color(1.0, 0.85, 0.2))
	if _font:
		title.add_theme_font_override("font", _font)
		title.add_theme_font_size_override("font_size", 18)
	add_child(title)

	# ????
	_close_btn = Button.new()
	_close_btn.text = "??"
	_close_btn.size = Vector2(28, 28)
	_close_btn.position = _bg.position + Vector2(PANEL_WIDTH - 36, 8)
	_close_btn.add_theme_color_override("font_color", Color(1, 0.4, 0.4))
	if _font:
		_close_btn.add_theme_font_override("font", _font)
	add_child(_close_btn)
	_close_btn.pressed.connect(_on_close_pressed)

	# ??蝺?
	var sep = ColorRect.new()
	sep.color = Color(0.3, 0.5, 0.8, 0.5)
	sep.size = Vector2(PANEL_WIDTH - 16, 2)
	sep.position = _bg.position + Vector2(8, 40)
	add_child(sep)

	# ???刻蝣?
	var my_code_title = Label.new()
	my_code_title.text = "???刻蝣潘?"
	my_code_title.position = _bg.position + Vector2(16, 52)
	my_code_title.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8))
	if _font:
		my_code_title.add_theme_font_override("font", _font)
		my_code_title.add_theme_font_size_override("font_size", 13)
	add_child(my_code_title)

	_my_code_label = Label.new()
	_my_code_label.text = "------"
	_my_code_label.position = _bg.position + Vector2(16, 72)
	_my_code_label.add_theme_color_override("font_color", Color(0.3, 1.0, 0.5))
	if _font:
		_my_code_label.add_theme_font_override("font", _font)
		_my_code_label.add_theme_font_size_override("font_size", 28)
	add_child(_my_code_label)

	# ?刻隤芣?
	var desc = Label.new()
	desc.text = "?刻鈭綽?+1000 ?馳  鋡急?虫犖嚗?500 ?馳"
	desc.position = _bg.position + Vector2(16, 108)
	desc.add_theme_color_override("font_color", Color(0.7, 0.9, 0.7))
	if _font:
		desc.add_theme_font_override("font", _font)
		desc.add_theme_font_size_override("font_size", 12)
	add_child(desc)

	# 蝯梯?
	_referral_count_label = Label.new()
	_referral_count_label.text = "撌脫?佗?0 鈭?"
	_referral_count_label.position = _bg.position + Vector2(16, 128)
	_referral_count_label.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8))
	if _font:
		_referral_count_label.add_theme_font_override("font", _font)
		_referral_count_label.add_theme_font_size_override("font_size", 13)
	add_child(_referral_count_label)

	_total_reward_label = Label.new()
	_total_reward_label.text = "蝝航??嚗? ?馳"
	_total_reward_label.position = _bg.position + Vector2(16, 148)
	_total_reward_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.2))
	if _font:
		_total_reward_label.add_theme_font_override("font", _font)
		_total_reward_label.add_theme_font_size_override("font_size", 13)
	add_child(_total_reward_label)

	# ??蝺?
	var sep2 = ColorRect.new()
	sep2.color = Color(0.3, 0.5, 0.8, 0.3)
	sep2.size = Vector2(PANEL_WIDTH - 16, 1)
	sep2.position = _bg.position + Vector2(8, 170)
	add_child(sep2)

	# 頛詨?刻蝣?
	var input_title = Label.new()
	input_title.text = "頛詨憟賢??刻蝣潘?"
	input_title.position = _bg.position + Vector2(16, 178)
	input_title.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8))
	if _font:
		input_title.add_theme_font_override("font", _font)
		input_title.add_theme_font_size_override("font_size", 13)
	add_child(input_title)

	_input_field = LineEdit.new()
	_input_field.placeholder_text = "頛詨6雿?衣Ⅳ"
	_input_field.max_length = 6
	_input_field.size = Vector2(180, 32)
	_input_field.position = _bg.position + Vector2(16, 198)
	if _font:
		_input_field.add_theme_font_override("font", _font)
		_input_field.add_theme_font_size_override("font_size", 16)
	add_child(_input_field)

	_use_btn = Button.new()
	_use_btn.text = "雿輻"
	_use_btn.size = Vector2(80, 32)
	_use_btn.position = _bg.position + Vector2(204, 198)
	if _font:
		_use_btn.add_theme_font_override("font", _font)
		_use_btn.add_theme_font_size_override("font_size", 14)
	add_child(_use_btn)
	_use_btn.pressed.connect(_on_use_btn_pressed)

	# ?????
	_status_label = Label.new()
	_status_label.text = ""
	_status_label.position = _bg.position + Vector2(16, 238)
	_status_label.size = Vector2(PANEL_WIDTH - 32, 24)
	if _font:
		_status_label.add_theme_font_override("font", _font)
		_status_label.add_theme_font_size_override("font_size", 12)
	add_child(_status_label)

func _connect_signals() -> void:
	if GameManager.has_signal("referral_info_received"):
		GameManager.referral_info_received.connect(_on_referral_info)
	if GameManager.has_signal("referral_success"):
		GameManager.referral_success.connect(_on_referral_success)
	if GameManager.has_signal("referral_error"):
		GameManager.referral_error.connect(_on_referral_error)

func _on_referral_info(data: Dictionary) -> void:
	var my_code: String = data.get("my_code", "------")
	var count: int = data.get("referral_count", 0)
	var total_reward: int = data.get("total_reward", 0)
	var used_code: String = data.get("used_code", "")

	_my_code_label.text = my_code
	_referral_count_label.text = "撌脫?佗?%d 鈭綽??憭?20 鈭綽?" % count
	_total_reward_label.text = "蝝航??嚗?d ?馳" % total_reward

	# 撌脖蝙?冽?衣Ⅳ???刻撓??
	if used_code != "":
		_input_field.editable = false
		_use_btn.disabled = true
		_input_field.text = used_code
		_status_label.text = "??撌脖蝙?冽?衣Ⅳ嚗?s" % used_code
		_status_label.add_theme_color_override("font_color", Color(0.5, 1.0, 0.5))

func _on_use_btn_pressed() -> void:
	var code = _input_field.text.strip_edges().to_upper()
	if code.length() != 6:
		_show_status("隢撓??雿?衣Ⅳ", Color(1.0, 0.4, 0.4))
		return
	if GameManager.has_method("use_referral_code"):
		GameManager.use_referral_code(code)
	_use_btn.disabled = true

func _on_referral_success(data: Dictionary) -> void:
	var reward: int = data.get("reward", 0)
	var msg: String = data.get("message", "")
	_show_status("??%s +%d ?馳" % [msg, reward], Color(0.5, 1.0, 0.5))
	_input_field.editable = false
	_use_btn.disabled = true
	# ?瑟鞈?
	if GameManager.has_method("request_referral_info"):
		GameManager.request_referral_info()

func _on_referral_error(data: Dictionary) -> void:
	var reason: String = data.get("reason", "?芰?航炊")
	_show_status("??%s" % reason, Color(1.0, 0.4, 0.4))
	_use_btn.disabled = false

func _show_status(text: String, color: Color) -> void:
	_status_label.text = text
	_status_label.add_theme_color_override("font_color", color)

func toggle() -> void:
	if _is_visible:
		_hide_panel()
	else:
		_show_panel()

func _show_panel() -> void:
	_is_visible = true
	show()
	if GameManager.has_method("request_referral_info"):
		GameManager.request_referral_info()

func _hide_panel() -> void:
	_is_visible = false
	hide()
	emit_signal("referral_closed")

func _on_close_pressed() -> void:
	_hide_panel()
