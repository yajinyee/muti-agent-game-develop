## LuckyTimeWarpPanel.gd — T128 幸運時間扭曲魚 UI
## lucky-panel-agent 負責維護
## 業界依據：Fishing Fortune 2026「Time Warp — all fish slow to 30% speed, damage ×2.0 for 10s」
## 視覺主題：紫色 + 計時條 + 傷害加成指示器 + 時鐘圖案
extends CanvasLayer

const LAYER_Z = 28  # CanvasLayer layer 值

# 顏色主題
const COLOR_WARP     = Color(0.55, 0.2, 0.86)   # 紫色（扭曲主色）
const COLOR_DAMAGE   = Color(1.0, 0.4, 0.1)     # 橙紅（傷害加成）
const COLOR_COLLAPSE = Color(1.0, 0.85, 0.0)    # 金色（時間崩潰）
const COLOR_BG       = Color(0.05, 0.0, 0.1, 0.88)

var _banner: Control = null
var _warp_panel: Control = null
var _result_popup: Control = null
var _flash_overlay: ColorRect = null
var _timer_bar: ColorRect = null
var _timer_label: Label = null
var _kill_label: Label = null
var _dmg_label: Label = null
var _warp_active: bool = false
var _warp_duration: float = 10.0
var _warp_elapsed: float = 0.0

func _ready() -> void:
	layer = LAYER_Z
	_create_flash_overlay()
	_create_warp_panel()
	GameManager.lucky_time_warp.connect(_on_lucky_time_warp)

func _process(delta: float) -> void:
	if _warp_active and is_instance_valid(_timer_bar):
		_warp_elapsed += delta
		var pct = 1.0 - (_warp_elapsed / _warp_duration)
		pct = clamp(pct, 0.0, 1.0)
		_timer_bar.size.x = 284.0 * pct
		# 顏色隨時間變化：紫→橙→紅
		if pct > 0.5:
			_timer_bar.color = COLOR_WARP
		elif pct > 0.25:
			_timer_bar.color = Color(0.8, 0.3, 0.9)
		else:
			_timer_bar.color = Color(1.0, 0.2, 0.5)
		if is_instance_valid(_timer_label):
			_timer_label.text = "%.0fs" % max(0.0, _warp_duration - _warp_elapsed)

func _create_flash_overlay() -> void:
	_flash_overlay = ColorRect.new()
	_flash_overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_flash_overlay.color = Color(0.55, 0.2, 0.86, 0.0)
	_flash_overlay.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_flash_overlay)

func _create_warp_panel() -> void:
	# 右上角時間扭曲狀態顯示
	_warp_panel = Control.new()
	_warp_panel.position = Vector2(960, 10)
	_warp_panel.size = Vector2(300, 100)
	_warp_panel.visible = false
	add_child(_warp_panel)

	var bg = ColorRect.new()
	bg.size = _warp_panel.size
	bg.color = COLOR_BG
	_warp_panel.add_child(bg)

	var title = Label.new()
	title.text = "⏰ 時間扭曲進行中"
	title.position = Vector2(8, 4)
	title.add_theme_font_size_override("font_size", 13)
	title.modulate = COLOR_WARP
	_warp_panel.add_child(title)

	# 計時條背景
	var bar_bg = ColorRect.new()
	bar_bg.position = Vector2(8, 26)
	bar_bg.size = Vector2(284, 14)
	bar_bg.color = Color(0.1, 0.0, 0.15)
	_warp_panel.add_child(bar_bg)

	# 計時條填充
	_timer_bar = ColorRect.new()
	_timer_bar.position = Vector2(8, 26)
	_timer_bar.size = Vector2(284, 14)
	_timer_bar.color = COLOR_WARP
	_warp_panel.add_child(_timer_bar)

	# 計時器文字
	_timer_label = Label.new()
	_timer_label.name = "TimerLabel"
	_timer_label.text = "10s"
	_timer_label.position = Vector2(8, 44)
	_timer_label.add_theme_font_size_override("font_size", 14)
	_timer_label.modulate = Color(0.9, 0.7, 1.0)
	_warp_panel.add_child(_timer_label)

	# 傷害加成標籤
	_dmg_label = Label.new()
	_dmg_label.name = "DmgLabel"
	_dmg_label.text = "傷害 ×2.0"
	_dmg_label.position = Vector2(160, 44)
	_dmg_label.add_theme_font_size_override("font_size", 14)
	_dmg_label.modulate = COLOR_DAMAGE
	_warp_panel.add_child(_dmg_label)

	# 擊破計數
	_kill_label = Label.new()
	_kill_label.name = "KillLabel"
	_kill_label.text = "擊破 0 條"
	_kill_label.position = Vector2(8, 68)
	_kill_label.add_theme_font_size_override("font_size", 14)
	_kill_label.modulate = Color(0.8, 0.6, 1.0)
	_warp_panel.add_child(_kill_label)

func _show_start_banner(name: String, duration: float, dmg: float) -> void:
	if is_instance_valid(_banner):
		_banner.queue_free()
	_banner = Control.new()
	_banner.position = Vector2(0, 100)
	_banner.size = Vector2(1280, 90)
	add_child(_banner)

	var bg = ColorRect.new()
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	bg.color = COLOR_BG
	_banner.add_child(bg)

	# 頂部紫色線
	var top_line = ColorRect.new()
	top_line.size = Vector2(1280, 3)
	top_line.color = COLOR_WARP
	_banner.add_child(top_line)

	var lbl = Label.new()
	lbl.text = "⏰ %s 觸發時間扭曲！全場慢速 %.0f 秒！" % [name, duration]
	lbl.position = Vector2(0, 10)
	lbl.size = Vector2(1280, 40)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.add_theme_font_size_override("font_size", 28)
	lbl.modulate = COLOR_WARP
	_banner.add_child(lbl)

	var sub = Label.new()
	sub.text = "傷害 ×%.0f！趁慢速瘋狂打魚！擊破 ≥6 條觸發時間崩潰！" % dmg
	sub.position = Vector2(0, 50)
	sub.size = Vector2(1280, 30)
	sub.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	sub.add_theme_font_size_override("font_size", 16)
	sub.modulate = Color(0.9, 0.7, 1.0)
	_banner.add_child(sub)

	# 底部紫色線
	var bot_line = ColorRect.new()
	bot_line.size = Vector2(1280, 3)
	bot_line.position = Vector2(0, 87)
	bot_line.color = COLOR_WARP
	_banner.add_child(bot_line)

	# 滑入動畫
	_banner.position.y = -90
	var tween = _banner.create_tween()
	tween.tween_property(_banner, "position:y", 100.0, 0.3).set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_BACK)
	tween.tween_interval(3.0)
	tween.tween_property(_banner, "modulate:a", 0.0, 0.4)
	tween.tween_callback(func():
		if is_instance_valid(_banner):
			_banner.queue_free()
			_banner = null
	)

func _show_collapse_popup(name: String, kills: int, boost: float, secs: int) -> void:
	if is_instance_valid(_result_popup):
		_result_popup.queue_free()

	_do_flash(COLOR_COLLAPSE, 5)

	_result_popup = Control.new()
	_result_popup.position = Vector2(340, 200)
	_result_popup.size = Vector2(600, 180)
	add_child(_result_popup)

	var bg = ColorRect.new()
	bg.size = _result_popup.size
	bg.color = COLOR_BG
	_result_popup.add_child(bg)

	var border = ColorRect.new()
	border.size = _result_popup.size
	border.color = Color(1.0, 0.85, 0.0, 0.2)
	_result_popup.add_child(border)

	var title_lbl = Label.new()
	title_lbl.text = "⏰💥 時間崩潰！"
	title_lbl.position = Vector2(0, 15)
	title_lbl.size = Vector2(600, 50)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title_lbl.add_theme_font_size_override("font_size", 32)
	title_lbl.modulate = COLOR_COLLAPSE
	_result_popup.add_child(title_lbl)

	var detail_lbl = Label.new()
	detail_lbl.text = "%s 擊破 %d 條！" % [name, kills]
	detail_lbl.position = Vector2(0, 70)
	detail_lbl.size = Vector2(600, 30)
	detail_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	detail_lbl.add_theme_font_size_override("font_size", 18)
	detail_lbl.modulate = Color(0.9, 0.9, 0.9)
	_result_popup.add_child(detail_lbl)

	var boost_lbl = Label.new()
	boost_lbl.text = "全服 ×%.0f 加成 %d 秒！" % [boost, secs]
	boost_lbl.position = Vector2(0, 110)
	boost_lbl.size = Vector2(600, 50)
	boost_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	boost_lbl.add_theme_font_size_override("font_size", 24)
	boost_lbl.modulate = COLOR_WARP
	_result_popup.add_child(boost_lbl)

	# 彈出動畫
	_result_popup.scale = Vector2.ZERO
	_result_popup.pivot_offset = Vector2(300, 90)
	var tween = _result_popup.create_tween()
	tween.tween_property(_result_popup, "scale", Vector2(1.0, 1.0), 0.3).set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_ELASTIC)
	tween.tween_interval(3.0)
	tween.tween_property(_result_popup, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func():
		if is_instance_valid(_result_popup):
			_result_popup.queue_free()
			_result_popup = null
	)

func _do_flash(color: Color, count: int) -> void:
	if not is_instance_valid(_flash_overlay):
		return
	var tween = _flash_overlay.create_tween()
	for i in count:
		tween.tween_property(_flash_overlay, "color", Color(color.r, color.g, color.b, 0.45), 0.08)
		tween.tween_property(_flash_overlay, "color", Color(color.r, color.g, color.b, 0.0), 0.12)

func _on_lucky_time_warp(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"warp_start":
			var duration = data.get("duration", 10.0)
			var dmg = data.get("damage_mult", 2.0)
			_show_start_banner(name, duration, dmg)
			# 啟動計時器
			_warp_active = true
			_warp_duration = duration
			_warp_elapsed = 0.0
			if is_instance_valid(_warp_panel):
				_warp_panel.visible = true
			if is_instance_valid(_dmg_label):
				_dmg_label.text = "傷害 ×%.0f" % dmg
			if is_instance_valid(_kill_label):
				_kill_label.text = "擊破 0 條"
			_do_flash(COLOR_WARP, 3)
			ScreenShake.add_trauma(0.35)
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
		"warp_end":
			_warp_active = false
			if is_instance_valid(_warp_panel):
				_warp_panel.visible = false
			var kills = data.get("kill_count", 0)
			if is_instance_valid(_kill_label):
				_kill_label.text = "擊破 %d 條" % kills
			_do_flash(COLOR_WARP, 2)
			ScreenShake.add_trauma(0.5)
		"time_collapse":
			_warp_active = false
			if is_instance_valid(_warp_panel):
				_warp_panel.visible = false
			var kills = data.get("kill_count", 0)
			var boost = data.get("boost_mult", 2.5)
			var secs = data.get("boost_secs", 6)
			_show_collapse_popup(name, kills, boost, secs)
			ScreenShake.add_trauma(0.7)
		"collapse_end":
			pass  # HUD 已處理橫幅
