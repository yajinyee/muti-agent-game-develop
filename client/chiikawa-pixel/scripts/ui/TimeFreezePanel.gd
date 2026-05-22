## TimeFreezePanel.gd — 時間凍結魚系統面板（DAY-212）
## 業界依據：Evolution Ice Fishing Live 2026「frozen paradise」概念
## 業界原創設計：「時間停止」讓玩家在凍結期間免費射擊所有靜止目標
##
## 視覺設計：
##   - 冰藍白主題（#00BFFF + #87CEEB + #FFFFFF + #1E90FF）
##   - freeze_start：全螢幕冰藍三次強閃光 + 頂部橫幅 + 底部計時條 + 畫面邊框冰晶效果
##   - freeze_end：冰藍淡出 + 「解凍中...」提示
##   - thaw_blast：全螢幕白色爆炸閃光 + 「💥 解凍爆炸！」大字 + 結算彈窗右側滑入
extends CanvasLayer

var _timer_bar: Control = null      # 底部計時條
var _ice_border: Control = null     # 畫面邊框冰晶效果
var _duration_sec: int = 5

func _ready() -> void:
	layer = 33  # 時間凍結面板層級

## 處理時間凍結魚訊息
func handle_time_freeze(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"freeze_start":
			_on_freeze_start(payload)
		"freeze_end":
			_on_freeze_end(payload)
		"thaw_blast":
			_on_thaw_blast(payload)

## 凍結開始 — 全螢幕冰藍三次強閃光 + 頂部橫幅 + 底部計時條 + 邊框冰晶
func _on_freeze_start(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "玩家")
	_duration_sec = payload.get("duration_sec", 5)
	var hit_bonus: float = payload.get("hit_bonus", 0.30)

	# 全螢幕冰藍三次強閃光
	_triple_flash(Color("#00BFFF"), 0.75)

	# 頂部橫幅
	var banner = _make_banner(
		"❄️ %s 觸發時間凍結！全場靜止 %d 秒！命中率 +%.0f%%！" % [player_name, _duration_sec, hit_bonus * 100],
		Color(0.0, 0.1, 0.2, 0.88),
		Color("#00BFFF")
	)
	add_child(banner)
	var tween_b = create_tween()
	tween_b.tween_property(banner, "position:y", 0.0, 0.25)
	tween_b.tween_interval(3.5)
	tween_b.tween_property(banner, "modulate:a", 0.0, 0.4)
	tween_b.tween_callback(banner.queue_free)

	# 畫面邊框冰晶效果（四邊冰藍半透明邊框）
	_show_ice_border()

	# 底部計時條
	_show_timer_bar(_duration_sec)

## 凍結結束 — 冰藍淡出 + 「解凍中...」提示
func _on_freeze_end(payload: Dictionary) -> void:
	var hit_count: int = payload.get("hit_count", 0)

	# 移除邊框冰晶
	if is_instance_valid(_ice_border):
		var tween = create_tween()
		tween.tween_property(_ice_border, "modulate:a", 0.0, 0.3)
		tween.tween_callback(_ice_border.queue_free)
		_ice_border = null

	# 移除計時條
	if is_instance_valid(_timer_bar):
		var tween2 = create_tween()
		tween2.tween_property(_timer_bar, "modulate:a", 0.0, 0.2)
		tween2.tween_callback(_timer_bar.queue_free)
		_timer_bar = null

	# 「解凍中...」提示
	if hit_count > 0:
		var label = Label.new()
		label.text = "🧊 解凍中... %d 個目標即將爆炸！" % hit_count
		label.add_theme_font_size_override("font_size", 22)
		label.add_theme_color_override("font_color", Color("#87CEEB"))
		label.position = Vector2(640 - 200, 300)
		label.size = Vector2(400, 40)
		label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
		add_child(label)
		var tween = create_tween()
		tween.tween_interval(0.25)
		tween.tween_property(label, "modulate:a", 0.0, 0.3)
		tween.tween_callback(label.queue_free)

## 解凍爆炸 — 全螢幕白色爆炸閃光 + 大字 + 結算彈窗
func _on_thaw_blast(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "玩家")
	var hit_count: int = payload.get("hit_count", 0)
	var killed_count: int = payload.get("killed_count", 0)
	var reward_per_player: int = payload.get("reward_per_player", 0)

	# 全螢幕白色爆炸閃光（冰→白→藍）
	_triple_flash(Color(1, 1, 1, 0.9), 0.8)

	# 「💥 解凍爆炸！」大字
	var big_label = Label.new()
	big_label.text = "💥 解凍爆炸！"
	big_label.add_theme_font_size_override("font_size", 52)
	big_label.add_theme_color_override("font_color", Color("#00BFFF"))
	big_label.position = Vector2(640 - 180, 280)
	big_label.size = Vector2(360, 70)
	big_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	add_child(big_label)
	var tween_big = create_tween()
	tween_big.tween_property(big_label, "scale", Vector2(1.2, 1.2), 0.1)
	tween_big.tween_property(big_label, "scale", Vector2(1.0, 1.0), 0.1)
	tween_big.tween_interval(1.2)
	tween_big.tween_property(big_label, "modulate:a", 0.0, 0.4)
	tween_big.tween_callback(big_label.queue_free)

	# 結算彈窗（右側滑入）
	if killed_count > 0 or reward_per_player > 0:
		_show_result_popup(player_name, hit_count, killed_count, reward_per_player)

## 顯示結算彈窗（右側滑入）
func _show_result_popup(player_name: String, hit_count: int, killed_count: int, reward: int) -> void:
	var popup = PanelContainer.new()
	popup.position = Vector2(1280, 200)
	popup.size = Vector2(280, 160)

	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.0, 0.08, 0.18, 0.92)
	style.border_color = Color("#00BFFF")
	style.border_width_left = 2
	style.border_width_right = 2
	style.border_width_top = 2
	style.border_width_bottom = 2
	style.corner_radius_top_left = 8
	style.corner_radius_top_right = 8
	style.corner_radius_bottom_left = 8
	style.corner_radius_bottom_right = 8
	popup.add_theme_stylebox_override("panel", style)

	var vbox = VBoxContainer.new()
	vbox.add_theme_constant_override("separation", 6)
	popup.add_child(vbox)

	var title = Label.new()
	title.text = "❄️ 時間凍結結算"
	title.add_theme_font_size_override("font_size", 16)
	title.add_theme_color_override("font_color", Color("#00BFFF"))
	title.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title)

	var sep = HSeparator.new()
	vbox.add_child(sep)

	var info1 = Label.new()
	info1.text = "命中目標：%d 個" % hit_count
	info1.add_theme_font_size_override("font_size", 14)
	info1.add_theme_color_override("font_color", Color("#87CEEB"))
	info1.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(info1)

	var info2 = Label.new()
	info2.text = "解凍擊破：%d 個" % killed_count
	info2.add_theme_font_size_override("font_size", 14)
	info2.add_theme_color_override("font_color", Color("#FFFFFF"))
	info2.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(info2)

	if reward > 0:
		var reward_label = Label.new()
		reward_label.text = "💰 每人獲得 %d 金幣" % reward
		reward_label.add_theme_font_size_override("font_size", 16)
		reward_label.add_theme_color_override("font_color", Color("#FFD700"))
		reward_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
		vbox.add_child(reward_label)

	add_child(popup)

	# 右側滑入動畫
	var tween = create_tween()
	tween.tween_property(popup, "position:x", 980.0, 0.3).set_ease(Tween.EASE_OUT)
	tween.tween_interval(4.0)
	tween.tween_property(popup, "modulate:a", 0.0, 0.5)
	tween.tween_callback(popup.queue_free)

## 顯示畫面邊框冰晶效果（四邊冰藍半透明邊框）
func _show_ice_border() -> void:
	if is_instance_valid(_ice_border):
		_ice_border.queue_free()

	var border = Control.new()
	border.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_ice_border = border
	add_child(border)

	# 四邊冰藍邊框（ColorRect）
	var border_width = 12
	var ice_color = Color(0.0, 0.75, 1.0, 0.55)

	for i in range(4):
		var rect = ColorRect.new()
		rect.color = ice_color
		match i:
			0: rect.position = Vector2(0, 0); rect.size = Vector2(1280, border_width)       # 上
			1: rect.position = Vector2(0, 720 - border_width); rect.size = Vector2(1280, border_width) # 下
			2: rect.position = Vector2(0, 0); rect.size = Vector2(border_width, 720)        # 左
			3: rect.position = Vector2(1280 - border_width, 0); rect.size = Vector2(border_width, 720) # 右
		border.add_child(rect)

	# 邊框脈衝動畫（sin 波動透明度）
	var tween = border.create_tween().set_loops()
	tween.tween_property(border, "modulate:a", 0.5, 0.4)
	tween.tween_property(border, "modulate:a", 1.0, 0.4)

## 顯示底部計時條
func _show_timer_bar(duration: int) -> void:
	if is_instance_valid(_timer_bar):
		_timer_bar.queue_free()

	var bar_container = Control.new()
	bar_container.position = Vector2(0, 700)
	bar_container.size = Vector2(1280, 20)
	_timer_bar = bar_container
	add_child(bar_container)

	# 背景
	var bg = ColorRect.new()
	bg.color = Color(0.0, 0.1, 0.2, 0.7)
	bg.position = Vector2(0, 0)
	bg.size = Vector2(1280, 20)
	bar_container.add_child(bg)

	# 進度條
	var bar = ColorRect.new()
	bar.color = Color("#00BFFF")
	bar.position = Vector2(0, 0)
	bar.size = Vector2(1280, 20)
	bar_container.add_child(bar)

	# 計時條縮短動畫（顏色漸變：冰藍→淡藍→白）
	var tween = create_tween()
	tween.tween_property(bar, "size:x", 0.0, float(duration))
	tween.parallel().tween_method(
		func(t: float) -> void:
			var c = Color("#00BFFF").lerp(Color("#FFFFFF"), t)
			bar.color = c,
		0.0, 1.0, float(duration)
	)

## 三次強閃光效果
func _triple_flash(color: Color, alpha: float) -> void:
	for i in range(3):
		var flash = ColorRect.new()
		flash.color = Color(color.r, color.g, color.b, 0.0)
		flash.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
		add_child(flash)
		var delay = i * 0.18
		var tween = create_tween()
		tween.tween_interval(delay)
		tween.tween_property(flash, "color:a", alpha, 0.08)
		tween.tween_property(flash, "color:a", 0.0, 0.14)
		tween.tween_callback(flash.queue_free)

## 建立頂部橫幅
func _make_banner(text: String, bg_color: Color, border_color: Color) -> PanelContainer:
	var panel = PanelContainer.new()
	panel.position = Vector2(0, -60)
	panel.size = Vector2(1280, 52)

	var style = StyleBoxFlat.new()
	style.bg_color = bg_color
	style.border_color = border_color
	style.border_width_bottom = 2
	panel.add_theme_stylebox_override("panel", style)

	var label = Label.new()
	label.text = text
	label.add_theme_font_size_override("font_size", 18)
	label.add_theme_color_override("font_color", Color("#FFFFFF"))
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	panel.add_child(label)

	return panel
