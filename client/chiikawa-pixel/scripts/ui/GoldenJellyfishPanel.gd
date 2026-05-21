## GoldenJellyfishPanel.gd
## 黃金水母全場電擊面板（DAY-149）
## 業界依據：Ocean King 3 2026「Electric Jellyfish chain shocks across multiple targets.
## Devastating against clustered schools.」
## 設計：黃金電流主題；shock_start 橫幅+全螢幕黃色閃光；
## shock 逐一電擊動畫；result 右側滑入彈窗含擊破列表+獎勵

extends Control

# ---- 節點 ----
var _banner: Control = null
var _banner_label: Label = null
var _result_panel: Control = null
var _result_title: Label = null
var _result_list: Label = null
var _result_reward: Label = null
var _flash_overlay: ColorRect = null
var _shock_counter: Label = null  # 電擊計數器

# ---- 顏色（黃金電流主題）----
const COLOR_GOLD    = Color(1.0, 0.84, 0.0, 1.0)
const COLOR_ELECTRIC = Color(0.9, 0.95, 0.2, 1.0)  # 黃綠電流色
const COLOR_BG      = Color(0.04, 0.04, 0.08, 0.95)

# ---- 狀態 ----
var _shock_entries: Array = []
var _total_kills: int = 0
var _total_reward: int = 0
var _killer_name: String = ""

func _ready() -> void:
	_build_ui()
	visible = false
	mouse_filter = Control.MOUSE_FILTER_IGNORE

func _build_ui() -> void:
	# 全螢幕閃光層
	_flash_overlay = ColorRect.new()
	_flash_overlay.set_anchors_preset(Control.PRESET_FULL_RECT)
	_flash_overlay.color = Color(1.0, 0.84, 0.0, 0.0)
	_flash_overlay.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_flash_overlay)

	# 頂部橫幅
	_banner = Control.new()
	_banner.position = Vector2(0, -80)
	_banner.size = Vector2(1280, 72)
	add_child(_banner)

	var banner_bg = ColorRect.new()
	banner_bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	banner_bg.color = Color(0.1, 0.1, 0.0, 0.92)
	_banner.add_child(banner_bg)

	var banner_border = ColorRect.new()
	banner_border.color = COLOR_GOLD
	banner_border.position = Vector2(0, 68)
	banner_border.size = Vector2(1280, 4)
	_banner.add_child(banner_border)

	_banner_label = Label.new()
	_banner_label.text = "⚡ 黃金水母全場電擊！"
	_banner_label.set_anchors_preset(Control.PRESET_FULL_RECT)
	_banner_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_banner_label.add_theme_color_override("font_color", COLOR_GOLD)
	_banner_label.add_theme_font_size_override("font_size", 26)
	_banner.add_child(_banner_label)

	# 電擊計數器（左上角）
	_shock_counter = Label.new()
	_shock_counter.text = ""
	_shock_counter.position = Vector2(20, 80)
	_shock_counter.size = Vector2(200, 40)
	_shock_counter.add_theme_color_override("font_color", COLOR_ELECTRIC)
	_shock_counter.add_theme_font_size_override("font_size", 18)
	_shock_counter.visible = false
	add_child(_shock_counter)

	# 右側結果彈窗
	_result_panel = Control.new()
	_result_panel.position = Vector2(1280, 120)  # 初始在螢幕右側外
	_result_panel.size = Vector2(320, 280)
	_result_panel.visible = false
	add_child(_result_panel)

	var result_bg = ColorRect.new()
	result_bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	result_bg.color = COLOR_BG
	_result_panel.add_child(result_bg)

	var result_border = ColorRect.new()
	result_border.set_anchors_preset(Control.PRESET_FULL_RECT)
	result_border.color = COLOR_GOLD
	result_border.custom_minimum_size = Vector2(320, 280)
	_result_panel.add_child(result_border)

	var result_inner = ColorRect.new()
	result_inner.color = COLOR_BG
	result_inner.position = Vector2(2, 2)
	result_inner.size = Vector2(316, 276)
	_result_panel.add_child(result_inner)

	_result_title = Label.new()
	_result_title.text = "⚡ 電擊結果"
	_result_title.position = Vector2(0, 10)
	_result_title.size = Vector2(320, 36)
	_result_title.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_result_title.add_theme_color_override("font_color", COLOR_GOLD)
	_result_title.add_theme_font_size_override("font_size", 18)
	_result_panel.add_child(_result_title)

	_result_list = Label.new()
	_result_list.text = ""
	_result_list.position = Vector2(12, 50)
	_result_list.size = Vector2(296, 160)
	_result_list.add_theme_color_override("font_color", Color(0.9, 0.9, 0.9, 1.0))
	_result_list.add_theme_font_size_override("font_size", 13)
	_result_list.autowrap_mode = TextServer.AUTOWRAP_WORD_SMART
	_result_panel.add_child(_result_list)

	_result_reward = Label.new()
	_result_reward.text = "+0 金幣"
	_result_reward.position = Vector2(0, 220)
	_result_reward.size = Vector2(320, 48)
	_result_reward.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_result_reward.add_theme_color_override("font_color", Color(0.2, 1.0, 0.4, 1.0))
	_result_reward.add_theme_font_size_override("font_size", 24)
	_result_panel.add_child(_result_reward)

# ---- 公開 API ----

## handle_shock_event 處理電擊事件（收到 golden_jellyfish_shock 時呼叫）
func handle_shock_event(data: Dictionary) -> void:
	var phase: String = data.get("phase", "")
	_killer_name = data.get("killer_name", "")

	match phase:
		"shock_start":
			_on_shock_start(data)
		"shock":
			_on_shock(data)
		"result":
			_on_result(data)

func _on_shock_start(data: Dictionary) -> void:
	_shock_entries = []
	_total_kills = 0
	_total_reward = 0

	# 更新橫幅
	_banner_label.text = "⚡ %s 的黃金水母全場電擊！" % _killer_name

	# 顯示橫幅（從上方滑入）
	visible = true
	_banner.position = Vector2(0, -80)
	var tween = create_tween()
	tween.tween_property(_banner, "position", Vector2(0, 0), 0.3).set_ease(Tween.EASE_OUT)

	# 全螢幕黃金電流閃光
	_flash_overlay.color = Color(1.0, 0.9, 0.0, 0.6)
	var flash_tween = create_tween()
	flash_tween.tween_property(_flash_overlay, "color", Color(1.0, 0.9, 0.0, 0.0), 0.4)

	# 顯示電擊計數器
	_shock_counter.text = "⚡ 電擊中..."
	_shock_counter.visible = true

func _on_shock(data: Dictionary) -> void:
	var targets: Array = data.get("targets", [])
	if targets.is_empty():
		return

	var entry = targets[0]
	var killed: bool = entry.get("killed", false)
	var shock_index: int = entry.get("shock_index", 0)
	var target_name: String = entry.get("target_name", "")
	var reward: int = entry.get("reward", 0)

	# 更新計數器
	_shock_counter.text = "⚡ 電擊 #%d：%s %s" % [shock_index + 1, target_name, "💥" if killed else "✗"]

	# 電擊閃光（每次電擊一個小閃光）
	_flash_overlay.color = Color(1.0, 0.95, 0.0, 0.25)
	var flash_tween = create_tween()
	flash_tween.tween_property(_flash_overlay, "color", Color(1.0, 0.95, 0.0, 0.0), 0.12)

	# 記錄
	_shock_entries.append(entry)
	if killed:
		_total_kills += 1
		_total_reward += reward

func _on_result(data: Dictionary) -> void:
	_total_kills = data.get("total_kills", _total_kills)
	_total_reward = data.get("total_reward", _total_reward)

	# 隱藏計數器
	_shock_counter.visible = false

	# 建立結果列表文字
	var list_text = ""
	var entries: Array = data.get("targets", _shock_entries)
	var shown = 0
	for entry in entries:
		if shown >= 6:
			list_text += "...\n"
			break
		var killed: bool = entry.get("killed", false)
		var name: String = entry.get("target_name", "")
		var reward: int = entry.get("reward", 0)
		var mult: float = entry.get("multiplier", 1.0)
		if killed:
			list_text += "💥 %s (%.0fx) +%d\n" % [name, mult, reward]
		else:
			list_text += "✗ %s (%.0fx)\n" % [name, mult]
		shown += 1

	_result_list.text = list_text
	_result_reward.text = "+%d 金幣（%d 個擊破）" % [_total_reward, _total_kills]

	# 根據擊破數決定顏色
	if _total_kills >= 5:
		_result_reward.add_theme_color_override("font_color", Color(1.0, 0.84, 0.0, 1.0))
		_do_mega_flash()
	elif _total_kills >= 3:
		_result_reward.add_theme_color_override("font_color", Color(0.2, 1.0, 0.4, 1.0))
		_do_gold_flash()

	# 顯示結果面板（從右側滑入）
	_result_panel.visible = true
	_result_panel.position = Vector2(1280, 120)
	var tween = create_tween()
	tween.tween_property(_result_panel, "position", Vector2(950, 120), 0.35).set_ease(Tween.EASE_OUT)

	# 3 秒後關閉
	var close_timer = get_tree().create_timer(3.0)
	close_timer.timeout.connect(_close_panel)

func _close_panel() -> void:
	# 橫幅滑出
	var tween = create_tween()
	tween.tween_property(_banner, "position", Vector2(0, -80), 0.3).set_ease(Tween.EASE_IN)
	# 結果面板滑出
	var tween2 = create_tween()
	tween2.tween_property(_result_panel, "position", Vector2(1280, 120), 0.3).set_ease(Tween.EASE_IN)
	tween2.tween_callback(func():
		visible = false
		_result_panel.visible = false
	)

func _do_gold_flash() -> void:
	_flash_overlay.color = Color(1.0, 0.84, 0.0, 0.5)
	var tween = create_tween()
	tween.tween_property(_flash_overlay, "color", Color(1.0, 0.84, 0.0, 0.0), 0.5)

func _do_mega_flash() -> void:
	# 雙閃光（5個以上擊破）
	_flash_overlay.color = Color(1.0, 0.9, 0.0, 0.7)
	var tween = create_tween()
	tween.tween_property(_flash_overlay, "color", Color(1.0, 0.9, 0.0, 0.0), 0.25)
	tween.tween_property(_flash_overlay, "color", Color(0.9, 1.0, 0.2, 0.5), 0.1)
	tween.tween_property(_flash_overlay, "color", Color(0.9, 1.0, 0.2, 0.0), 0.35)
