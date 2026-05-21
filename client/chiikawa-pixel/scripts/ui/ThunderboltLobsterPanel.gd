## ThunderboltLobsterPanel.gd
## 雷霆龍蝦免費射擊面板（DAY-150）
## 業界依據：royalfishingsite.com 2026「Thunderbolt Lobster feature —
## 15 seconds of free play followed by automatic shooting」
## 設計：橙紅電流主題；activate 橫幅+倒數計時+全螢幕橙色閃光；
## shot 自動射擊動畫+計數器；end 右側滑入彈窗含射擊統計+獎勵

extends Control

# ---- 節點 ----
var _banner: Control = null
var _banner_label: Label = null
var _countdown_label: Label = null   # 倒數計時
var _shot_counter: Label = null      # 射擊計數器
var _result_panel: Control = null
var _result_title: Label = null
var _result_stats: Label = null
var _result_reward: Label = null
var _flash_overlay: ColorRect = null
var _progress_bar: ColorRect = null  # 免費射擊進度條
var _progress_fill: ColorRect = null

# ---- 顏色（橙紅電流主題）----
const COLOR_ORANGE   = Color(1.0, 0.5, 0.0, 1.0)
const COLOR_ELECTRIC = Color(1.0, 0.7, 0.1, 1.0)
const COLOR_RED      = Color(0.9, 0.2, 0.1, 1.0)
const COLOR_BG       = Color(0.06, 0.02, 0.0, 0.95)

# ---- 狀態 ----
var _killer_id: String = ""
var _killer_name: String = ""
var _duration: int = 15
var _shot_interval: int = 500
var _total_shots: int = 0
var _total_kills: int = 0
var _total_reward: int = 0
var _shots_left: int = 30
var _is_my_trigger: bool = false
var _countdown_timer: float = 0.0
var _is_active: bool = false

func _ready() -> void:
	_build_ui()
	visible = false
	mouse_filter = Control.MOUSE_FILTER_IGNORE

func _process(delta: float) -> void:
	if not _is_active:
		return
	_countdown_timer -= delta
	if _countdown_timer < 0.0:
		_countdown_timer = 0.0
	# 更新倒數計時
	if _countdown_label:
		var secs = int(ceil(_countdown_timer))
		_countdown_label.text = "⚡ %d 秒" % secs
		if secs <= 5:
			_countdown_label.add_theme_color_override("font_color", COLOR_RED)
		else:
			_countdown_label.add_theme_color_override("font_color", COLOR_ELECTRIC)
	# 更新進度條
	if _progress_fill and _duration > 0:
		var ratio = _countdown_timer / float(_duration)
		_progress_fill.size.x = 1240.0 * ratio

func _build_ui() -> void:
	# 全螢幕閃光層
	_flash_overlay = ColorRect.new()
	_flash_overlay.set_anchors_preset(Control.PRESET_FULL_RECT)
	_flash_overlay.color = Color(1.0, 0.5, 0.0, 0.0)
	_flash_overlay.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_flash_overlay)

	# 頂部橫幅
	_banner = Control.new()
	_banner.position = Vector2(0, -80)
	_banner.size = Vector2(1280, 72)
	add_child(_banner)

	var banner_bg = ColorRect.new()
	banner_bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	banner_bg.color = Color(0.12, 0.04, 0.0, 0.92)
	_banner.add_child(banner_bg)

	var banner_border = ColorRect.new()
	banner_border.color = COLOR_ORANGE
	banner_border.position = Vector2(0, 68)
	banner_border.size = Vector2(1280, 4)
	_banner.add_child(banner_border)

	_banner_label = Label.new()
	_banner_label.text = "⚡ 雷霆龍蝦！免費射擊 15 秒！"
	_banner_label.set_anchors_preset(Control.PRESET_FULL_RECT)
	_banner_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_banner_label.add_theme_color_override("font_color", COLOR_ORANGE)
	_banner_label.add_theme_font_size_override("font_size", 26)
	_banner.add_child(_banner_label)

	# 倒數計時（橫幅右側）
	_countdown_label = Label.new()
	_countdown_label.text = "⚡ 15 秒"
	_countdown_label.position = Vector2(1080, 0)
	_countdown_label.size = Vector2(180, 72)
	_countdown_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_countdown_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_countdown_label.add_theme_color_override("font_color", COLOR_ELECTRIC)
	_countdown_label.add_theme_font_size_override("font_size", 22)
	_banner.add_child(_countdown_label)

	# 進度條（橫幅下方）
	_progress_bar = ColorRect.new()
	_progress_bar.position = Vector2(20, 76)
	_progress_bar.size = Vector2(1240, 8)
	_progress_bar.color = Color(0.2, 0.1, 0.0, 0.8)
	add_child(_progress_bar)

	_progress_fill = ColorRect.new()
	_progress_fill.position = Vector2(20, 76)
	_progress_fill.size = Vector2(1240, 8)
	_progress_fill.color = COLOR_ORANGE
	add_child(_progress_fill)

	# 射擊計數器（左上角）
	_shot_counter = Label.new()
	_shot_counter.text = ""
	_shot_counter.position = Vector2(20, 90)
	_shot_counter.size = Vector2(300, 40)
	_shot_counter.add_theme_color_override("font_color", COLOR_ELECTRIC)
	_shot_counter.add_theme_font_size_override("font_size", 16)
	_shot_counter.visible = false
	add_child(_shot_counter)

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
	result_border.color = COLOR_ORANGE
	result_border.custom_minimum_size = Vector2(320, 260)
	_result_panel.add_child(result_border)

	var result_inner = ColorRect.new()
	result_inner.color = COLOR_BG
	result_inner.position = Vector2(2, 2)
	result_inner.size = Vector2(316, 256)
	_result_panel.add_child(result_inner)

	_result_title = Label.new()
	_result_title.text = "⚡ 免費射擊結果"
	_result_title.position = Vector2(0, 10)
	_result_title.size = Vector2(320, 36)
	_result_title.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_result_title.add_theme_color_override("font_color", COLOR_ORANGE)
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

## handle_activate 處理免費射擊開始（收到 thunderbolt_lobster_activate 時呼叫）
func handle_activate(data: Dictionary) -> void:
	_killer_id = data.get("killer_id", "")
	_killer_name = data.get("killer_name", "")
	_duration = data.get("duration", 15)
	_shot_interval = data.get("shot_interval", 500)
	_total_shots = 0
	_total_kills = 0
	_total_reward = 0
	_shots_left = 30
	_is_active = true
	_countdown_timer = float(_duration)

	# 判斷是否為自己觸發
	var my_id = ""
	if has_node("/root/GameManager"):
		my_id = get_node("/root/GameManager").get_player_id()
	_is_my_trigger = (_killer_id == my_id)

	# 更新橫幅
	_banner_label.text = "⚡ %s 觸發雷霆龍蝦！免費射擊 %d 秒！" % [_killer_name, _duration]

	# 顯示橫幅（從上方滑入）
	visible = true
	_banner.position = Vector2(0, -80)
	var tween = create_tween()
	tween.tween_property(_banner, "position", Vector2(0, 0), 0.3).set_ease(Tween.EASE_OUT)

	# 全螢幕橙色閃光
	var flash_alpha = 0.7 if _is_my_trigger else 0.4
	_flash_overlay.color = Color(1.0, 0.5, 0.0, flash_alpha)
	var flash_tween = create_tween()
	flash_tween.tween_property(_flash_overlay, "color", Color(1.0, 0.5, 0.0, 0.0), 0.5)

	# 顯示射擊計數器
	_shot_counter.text = "⚡ 免費射擊中..."
	_shot_counter.visible = true

## handle_shot 處理自動射擊一次（收到 thunderbolt_lobster_shot 時呼叫）
func handle_shot(data: Dictionary) -> void:
	if not _is_active:
		return

	var is_kill: bool = data.get("is_kill", false)
	var reward: int = data.get("reward", 0)
	var target_name: String = data.get("target_name", "")
	var multiplier: float = data.get("multiplier", 1.0)
	var shot_index: int = data.get("shot_index", 0)
	_shots_left = data.get("shots_left", 0)

	_total_shots = shot_index + 1
	if is_kill:
		_total_kills += 1
		_total_reward += reward

	# 更新計數器
	var kill_icon = "💥" if is_kill else "•"
	_shot_counter.text = "⚡ 射擊 #%d %s %s (%.0fx)" % [_total_shots, kill_icon, target_name, multiplier]

	# 每次射擊小閃光
	_flash_overlay.color = Color(1.0, 0.6, 0.0, 0.15)
	var flash_tween = create_tween()
	flash_tween.tween_property(_flash_overlay, "color", Color(1.0, 0.6, 0.0, 0.0), 0.1)

## handle_end 處理免費射擊結束（收到 thunderbolt_lobster_end 時呼叫）
func handle_end(data: Dictionary) -> void:
	_is_active = false
	_total_shots = data.get("total_shots", _total_shots)
	_total_kills = data.get("total_kills", _total_kills)
	_total_reward = data.get("total_reward", _total_reward)

	# 隱藏計數器
	_shot_counter.visible = false

	# 建立統計文字
	var stats_text = ""
	stats_text += "🎯 總射擊：%d 次\n" % _total_shots
	stats_text += "💥 擊破：%d 個目標\n" % _total_kills
	if _total_shots > 0:
		var hit_rate = int(float(_total_kills) / float(_total_shots) * 100.0)
		stats_text += "📊 擊破率：%d%%\n" % hit_rate
	stats_text += "\n⚡ 全部免費！零消耗！"

	_result_stats.text = stats_text
	_result_reward.text = "+%d 金幣" % _total_reward

	# 根據擊破數決定顏色和閃光
	if _total_kills >= 10:
		_result_reward.add_theme_color_override("font_color", COLOR_ORANGE)
		_do_mega_flash()
	elif _total_kills >= 5:
		_result_reward.add_theme_color_override("font_color", COLOR_ELECTRIC)
		_do_gold_flash()
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
	# 橫幅滑出
	var tween = create_tween()
	tween.tween_property(_banner, "position", Vector2(0, -80), 0.3).set_ease(Tween.EASE_IN)
	# 進度條淡出
	if _progress_fill:
		var tween2 = create_tween()
		tween2.tween_property(_progress_fill, "color", Color(1.0, 0.5, 0.0, 0.0), 0.3)
	# 結果面板滑出
	var tween3 = create_tween()
	tween3.tween_property(_result_panel, "position", Vector2(1280, 120), 0.3).set_ease(Tween.EASE_IN)
	tween3.tween_callback(func():
		visible = false
		_result_panel.visible = false
		_progress_fill.color = COLOR_ORANGE  # 重置進度條顏色
	)

func _do_gold_flash() -> void:
	_flash_overlay.color = Color(1.0, 0.6, 0.0, 0.5)
	var tween = create_tween()
	tween.tween_property(_flash_overlay, "color", Color(1.0, 0.6, 0.0, 0.0), 0.4)

func _do_mega_flash() -> void:
	# 雙閃光（≥10 個擊破）
	_flash_overlay.color = Color(1.0, 0.5, 0.0, 0.7)
	var tween = create_tween()
	tween.tween_property(_flash_overlay, "color", Color(1.0, 0.5, 0.0, 0.0), 0.25)
	tween.tween_property(_flash_overlay, "color", Color(1.0, 0.5, 0.0, 0.5), 0.15)
	tween.tween_property(_flash_overlay, "color", Color(1.0, 0.5, 0.0, 0.0), 0.3)
