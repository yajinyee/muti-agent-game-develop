## RainbowPhoenixPanel.gd
## 彩虹鳳凰 Power Up 面板（DAY-151）
## 業界依據：royal-fishing.co.uk 2026「Multicoloured phoenix (blue, pink, purple, orange)
## with magical aura. Power Up attack delivers 6x-10x boost for rewards up to 300 times bet.」
## 設計：彩虹漸層主題；activate 橫幅+倒數計時+全螢幕彩虹閃光；
## end 右側滑入彈窗含 Power Up 統計+獎勵

extends Control

# ---- 節點 ----
var _banner: Control = null
var _banner_label: Label = null
var _mult_label: Label = null        # Power Up 倍率顯示
var _countdown_label: Label = null   # 倒數計時
var _result_panel: Control = null
var _result_title: Label = null
var _result_stats: Label = null
var _result_reward: Label = null
var _flash_overlay: ColorRect = null
var _progress_bar: ColorRect = null
var _progress_fill: ColorRect = null

# ---- 彩虹顏色序列 ----
const RAINBOW_COLORS = [
	Color(1.0, 0.2, 0.2, 1.0),  # 紅
	Color(1.0, 0.6, 0.0, 1.0),  # 橙
	Color(1.0, 1.0, 0.0, 1.0),  # 黃
	Color(0.2, 1.0, 0.2, 1.0),  # 綠
	Color(0.2, 0.6, 1.0, 1.0),  # 藍
	Color(0.6, 0.2, 1.0, 1.0),  # 紫
	Color(1.0, 0.4, 1.0, 1.0),  # 粉
]
const COLOR_BG = Color(0.04, 0.0, 0.08, 0.95)

# ---- 狀態 ----
var _killer_id: String = ""
var _killer_name: String = ""
var _power_up_mult: float = 6.0
var _duration: int = 8
var _total_kills: int = 0
var _total_reward: int = 0
var _is_my_trigger: bool = false
var _countdown_timer: float = 0.0
var _is_active: bool = false
var _rainbow_index: int = 0
var _rainbow_timer: float = 0.0

func _ready() -> void:
	_build_ui()
	visible = false
	mouse_filter = Control.MOUSE_FILTER_IGNORE

func _process(delta: float) -> void:
	if not _is_active:
		return

	# 倒數計時
	_countdown_timer -= delta
	if _countdown_timer < 0.0:
		_countdown_timer = 0.0
	if _countdown_label:
		var secs = int(ceil(_countdown_timer))
		_countdown_label.text = "🌈 %d 秒" % secs
		if secs <= 3:
			_countdown_label.add_theme_color_override("font_color", Color(1.0, 0.3, 0.3, 1.0))

	# 彩虹顏色循環（橫幅文字）
	_rainbow_timer += delta
	if _rainbow_timer >= 0.15:
		_rainbow_timer = 0.0
		_rainbow_index = (_rainbow_index + 1) % RAINBOW_COLORS.size()
		if _banner_label:
			_banner_label.add_theme_color_override("font_color", RAINBOW_COLORS[_rainbow_index])
		if _mult_label:
			_mult_label.add_theme_color_override("font_color", RAINBOW_COLORS[(_rainbow_index + 2) % RAINBOW_COLORS.size()])

	# 進度條
	if _progress_fill and _duration > 0:
		var ratio = _countdown_timer / float(_duration)
		_progress_fill.size.x = 1240.0 * ratio
		# 進度條顏色也跟著彩虹變
		_progress_fill.color = RAINBOW_COLORS[_rainbow_index]

func _build_ui() -> void:
	# 全螢幕閃光層
	_flash_overlay = ColorRect.new()
	_flash_overlay.set_anchors_preset(Control.PRESET_FULL_RECT)
	_flash_overlay.color = Color(0.8, 0.2, 1.0, 0.0)
	_flash_overlay.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_flash_overlay)

	# 頂部橫幅
	_banner = Control.new()
	_banner.position = Vector2(0, -80)
	_banner.size = Vector2(1280, 72)
	add_child(_banner)

	var banner_bg = ColorRect.new()
	banner_bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	banner_bg.color = Color(0.06, 0.0, 0.12, 0.92)
	_banner.add_child(banner_bg)

	var banner_border = ColorRect.new()
	banner_border.color = Color(0.8, 0.2, 1.0, 1.0)
	banner_border.position = Vector2(0, 68)
	banner_border.size = Vector2(1280, 4)
	_banner.add_child(banner_border)

	_banner_label = Label.new()
	_banner_label.text = "🌈 彩虹鳳凰 Power Up！"
	_banner_label.set_anchors_preset(Control.PRESET_FULL_RECT)
	_banner_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_banner_label.add_theme_color_override("font_color", Color(1.0, 0.4, 1.0, 1.0))
	_banner_label.add_theme_font_size_override("font_size", 26)
	_banner.add_child(_banner_label)

	# Power Up 倍率顯示（橫幅左側）
	_mult_label = Label.new()
	_mult_label.text = "×6"
	_mult_label.position = Vector2(20, 0)
	_mult_label.size = Vector2(120, 72)
	_mult_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_mult_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_mult_label.add_theme_color_override("font_color", Color(1.0, 1.0, 0.0, 1.0))
	_mult_label.add_theme_font_size_override("font_size", 32)
	_banner.add_child(_mult_label)

	# 倒數計時（橫幅右側）
	_countdown_label = Label.new()
	_countdown_label.text = "🌈 8 秒"
	_countdown_label.position = Vector2(1100, 0)
	_countdown_label.size = Vector2(160, 72)
	_countdown_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_countdown_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_countdown_label.add_theme_color_override("font_color", Color(0.8, 0.6, 1.0, 1.0))
	_countdown_label.add_theme_font_size_override("font_size", 20)
	_banner.add_child(_countdown_label)

	# 進度條
	_progress_bar = ColorRect.new()
	_progress_bar.position = Vector2(20, 76)
	_progress_bar.size = Vector2(1240, 8)
	_progress_bar.color = Color(0.15, 0.0, 0.2, 0.8)
	add_child(_progress_bar)

	_progress_fill = ColorRect.new()
	_progress_fill.position = Vector2(20, 76)
	_progress_fill.size = Vector2(1240, 8)
	_progress_fill.color = Color(0.8, 0.2, 1.0, 1.0)
	add_child(_progress_fill)

	# 右側結果彈窗
	_result_panel = Control.new()
	_result_panel.position = Vector2(1280, 120)
	_result_panel.size = Vector2(320, 260)
	_result_panel.visible = false
	add_child(_result_panel)

	var result_bg = ColorRect.new()
	result_bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	result_bg.color = COLOR_BG
	_result_panel.add_child(result_bg)

	var result_border = ColorRect.new()
	result_border.set_anchors_preset(Control.PRESET_FULL_RECT)
	result_border.color = Color(0.8, 0.2, 1.0, 1.0)
	result_border.custom_minimum_size = Vector2(320, 260)
	_result_panel.add_child(result_border)

	var result_inner = ColorRect.new()
	result_inner.color = COLOR_BG
	result_inner.position = Vector2(2, 2)
	result_inner.size = Vector2(316, 256)
	_result_panel.add_child(result_inner)

	_result_title = Label.new()
	_result_title.text = "🌈 Power Up 結果"
	_result_title.position = Vector2(0, 10)
	_result_title.size = Vector2(320, 36)
	_result_title.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_result_title.add_theme_color_override("font_color", Color(0.8, 0.2, 1.0, 1.0))
	_result_title.add_theme_font_size_override("font_size", 18)
	_result_panel.add_child(_result_title)

	_result_stats = Label.new()
	_result_stats.text = ""
	_result_stats.position = Vector2(12, 50)
	_result_stats.size = Vector2(296, 140)
	_result_stats.add_theme_color_override("font_color", Color(0.9, 0.9, 0.9, 1.0))
	_result_stats.add_theme_font_size_override("font_size", 14)
	_result_stats.autowrap_mode = TextServer.AUTOWRAP_WORD_SMART
	_result_panel.add_child(_result_stats)

	_result_reward = Label.new()
	_result_reward.text = "+0 金幣"
	_result_reward.position = Vector2(0, 200)
	_result_reward.size = Vector2(320, 48)
	_result_reward.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_result_reward.add_theme_color_override("font_color", Color(0.2, 1.0, 0.4, 1.0))
	_result_reward.add_theme_font_size_override("font_size", 24)
	_result_panel.add_child(_result_reward)

# ---- 公開 API ----

## handle_activate 處理 Power Up 開始（收到 rainbow_phoenix_activate 時呼叫）
func handle_activate(data: Dictionary) -> void:
	_killer_id = data.get("killer_id", "")
	_killer_name = data.get("killer_name", "")
	_power_up_mult = data.get("power_up_mult", 6.0)
	_duration = data.get("duration", 8)
	_total_kills = 0
	_total_reward = 0
	_is_active = true
	_countdown_timer = float(_duration)
	_rainbow_index = 0
	_rainbow_timer = 0.0

	# 判斷是否為自己觸發
	var my_id = ""
	if has_node("/root/GameManager"):
		my_id = get_node("/root/GameManager").get_player_id()
	_is_my_trigger = (_killer_id == my_id)

	# 更新橫幅
	_banner_label.text = "🌈 %s 觸發彩虹鳳凰！Power Up ×%.0f！" % [_killer_name, _power_up_mult]
	_mult_label.text = "×%.0f" % _power_up_mult

	# 顯示橫幅（從上方滑入）
	visible = true
	_banner.position = Vector2(0, -80)
	var tween = create_tween()
	tween.tween_property(_banner, "position", Vector2(0, 0), 0.3).set_ease(Tween.EASE_OUT)

	# 全螢幕彩虹閃光
	var flash_alpha = 0.65 if _is_my_trigger else 0.35
	_flash_overlay.color = Color(0.8, 0.2, 1.0, flash_alpha)
	var flash_tween = create_tween()
	flash_tween.tween_property(_flash_overlay, "color", Color(0.8, 0.2, 1.0, 0.0), 0.5)

## handle_end 處理 Power Up 結束（收到 rainbow_phoenix_end 時呼叫）
func handle_end(data: Dictionary) -> void:
	_is_active = false
	_total_kills = data.get("total_kills", _total_kills)
	_total_reward = data.get("total_reward", _total_reward)
	_power_up_mult = data.get("power_up_mult", _power_up_mult)

	# 建立統計文字
	var stats_text = ""
	stats_text += "🌈 Power Up 倍率：×%.0f\n" % _power_up_mult
	stats_text += "💥 擊破目標：%d 個\n" % _total_kills
	stats_text += "💰 Power Up 獎勵：%d 金幣\n" % _total_reward
	stats_text += "\n✨ 所有獎勵已套用 ×%.0f 倍率！" % _power_up_mult

	_result_stats.text = stats_text
	_result_reward.text = "+%d 金幣" % _total_reward

	# 根據擊破數決定顏色和閃光
	if _total_kills >= 5:
		_result_reward.add_theme_color_override("font_color", Color(1.0, 0.4, 1.0, 1.0))
		_do_rainbow_flash()
	elif _total_kills >= 3:
		_result_reward.add_theme_color_override("font_color", Color(0.8, 0.2, 1.0, 1.0))
		_do_purple_flash()
	else:
		_result_reward.add_theme_color_override("font_color", Color(0.2, 1.0, 0.4, 1.0))

	# 顯示結果面板（從右側滑入）
	_result_panel.visible = true
	_result_panel.position = Vector2(1280, 120)
	var tween = create_tween()
	tween.tween_property(_result_panel, "position", Vector2(950, 120), 0.35).set_ease(Tween.EASE_OUT)

	# 4 秒後關閉
	var close_timer = get_tree().create_timer(4.0)
	close_timer.timeout.connect(_close_panel)

func _close_panel() -> void:
	_is_active = false
	var tween = create_tween()
	tween.tween_property(_banner, "position", Vector2(0, -80), 0.3).set_ease(Tween.EASE_IN)
	if _progress_fill:
		var tween2 = create_tween()
		tween2.tween_property(_progress_fill, "color", Color(0.8, 0.2, 1.0, 0.0), 0.3)
	var tween3 = create_tween()
	tween3.tween_property(_result_panel, "position", Vector2(1280, 120), 0.3).set_ease(Tween.EASE_IN)
	tween3.tween_callback(func():
		visible = false
		_result_panel.visible = false
		_progress_fill.color = Color(0.8, 0.2, 1.0, 1.0)
	)

func _do_purple_flash() -> void:
	_flash_overlay.color = Color(0.6, 0.0, 1.0, 0.5)
	var tween = create_tween()
	tween.tween_property(_flash_overlay, "color", Color(0.6, 0.0, 1.0, 0.0), 0.4)

func _do_rainbow_flash() -> void:
	# 彩虹雙閃光（≥5 個擊破）
	_flash_overlay.color = Color(1.0, 0.4, 1.0, 0.7)
	var tween = create_tween()
	tween.tween_property(_flash_overlay, "color", Color(1.0, 0.4, 1.0, 0.0), 0.2)
	tween.tween_property(_flash_overlay, "color", Color(0.4, 0.8, 1.0, 0.5), 0.15)
	tween.tween_property(_flash_overlay, "color", Color(0.4, 0.8, 1.0, 0.0), 0.3)
