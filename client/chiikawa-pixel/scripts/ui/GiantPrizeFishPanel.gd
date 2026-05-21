## GiantPrizeFishPanel.gd — 夢幻巨型獎勵魚面板（DAY-147）
## 業界依據：jiligames.com 2026「The dreamy Giant Prize Fish lets you easily win great prizes,
## with the chance for 5x multipliers」
## 觸發玩家在 10 秒內所有擊破獎勵 ×5，全服廣播夢幻模式開始/結束
## 視覺設計：粉紅夢幻主題；頂部橫幅滑入；倒數計時；×5 倍率顯示；全螢幕粉紅閃光；結束彈窗
extends Control

# ---- 常數 ----
const PANEL_COLOR_BG    := Color(0.95, 0.4, 0.7, 0.92)   # 粉紅主題背景
const PANEL_COLOR_GOLD  := Color(1.0, 0.85, 0.0, 1.0)    # 金色文字
const PANEL_COLOR_WHITE := Color(1.0, 1.0, 1.0, 1.0)
const PANEL_COLOR_PINK  := Color(1.0, 0.4, 0.7, 1.0)     # 粉紅色
const PANEL_COLOR_DREAM := Color(0.98, 0.75, 0.9, 1.0)   # 夢幻淡粉

# ---- 狀態 ----
var _pixel_font: Font = null
var _is_active: bool = false
var _is_my_session: bool = false   # 是否是自己觸發的
var _duration: float = 10.0
var _elapsed: float = 0.0
var _mult_bonus: float = 5.0
var _killer_name: String = ""

# ---- 動態節點 ----
var _banner: Control = null
var _timer_label: Label = null
var _progress_bar: ColorRect = null

## setup — 由 HUD.gd 呼叫，連接 GameManager 訊號
func setup(font: Font) -> void:
	_pixel_font = font
	GameManager.giant_prize_fish.connect(_on_giant_prize_fish)

## _on_giant_prize_fish — 處理 Server 廣播的夢幻獎勵魚事件
func _on_giant_prize_fish(data: Dictionary) -> void:
	var phase: String = data.get("phase", "")
	var killer_id: String = data.get("killer_id", "")
	_killer_name = data.get("killer_name", "玩家")
	_mult_bonus = data.get("mult_bonus", 5.0)
	_duration = float(data.get("duration", 10))
	_is_my_session = (killer_id == NetworkManager.get_player_id())

	match phase:
		"activate":
			_on_activate(data)
		"end":
			_on_end(data)

func _process(delta: float) -> void:
	if not _is_active:
		return

	_elapsed += delta
	var remaining = max(0.0, _duration - _elapsed)
	var pct = remaining / _duration

	# 更新倒數計時
	if is_instance_valid(_timer_label):
		_timer_label.text = "✨ %.1fs" % remaining

	# 更新進度條
	if is_instance_valid(_progress_bar):
		_progress_bar.size.x = 180.0 * pct

	# 最後 3 秒：紅色閃爍
	if is_instance_valid(_timer_label) and remaining <= 3.0:
		var blink = sin(_elapsed * 8.0) > 0.0
		_timer_label.add_theme_color_override("font_color",
			Color(1.0, 0.2, 0.2) if blink else PANEL_COLOR_PINK)

	# 時間到：備用清理
	if _elapsed >= _duration + 0.5:
		_is_active = false

func _on_activate(data: Dictionary) -> void:
	_is_active = true
	_elapsed = 0.0

	var canvas_layer = get_parent()
	if not is_instance_valid(canvas_layer):
		return

	# 頂部橫幅（粉紅夢幻主題）
	_banner = Control.new()
	_banner.name = "GiantPrizeFishBanner"
	_banner.position = Vector2(0, -60)
	_banner.size = Vector2(1280, 56)
	_banner.z_index = 88
	canvas_layer.add_child(_banner)

	var bg = ColorRect.new()
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	bg.color = Color(0.6, 0.1, 0.35, 0.92)
	_banner.add_child(bg)

	# 標題
	var title_lbl = Label.new()
	if _is_my_session:
		title_lbl.text = "✨ 夢幻獎勵模式！×%.0f 加成！" % _mult_bonus
	else:
		title_lbl.text = "✨ %s 觸發夢幻獎勵！×%.0f 加成！" % [_killer_name, _mult_bonus]
	title_lbl.position = Vector2(0, 4)
	title_lbl.size = Vector2(1280, 28)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title_lbl.add_theme_font_size_override("font_size", 20)
	title_lbl.add_theme_color_override("font_color", PANEL_COLOR_GOLD)
	if is_instance_valid(_pixel_font):
		title_lbl.add_theme_font_override("font", _pixel_font)
	_banner.add_child(title_lbl)

	# 倒數計時標籤
	_timer_label = Label.new()
	_timer_label.name = "TimerLabel"
	_timer_label.text = "✨ %.1fs" % _duration
	_timer_label.position = Vector2(0, 32)
	_timer_label.size = Vector2(1280, 20)
	_timer_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_timer_label.add_theme_font_size_override("font_size", 12)
	_timer_label.add_theme_color_override("font_color", PANEL_COLOR_PINK)
	if is_instance_valid(_pixel_font):
		_timer_label.add_theme_font_override("font", _pixel_font)
	_banner.add_child(_timer_label)

	# 進度條背景
	var bar_bg = ColorRect.new()
	bar_bg.position = Vector2(550, 50)
	bar_bg.size = Vector2(180, 6)
	bar_bg.color = Color(0.1, 0.1, 0.1, 0.8)
	_banner.add_child(bar_bg)

	# 進度條填充
	_progress_bar = ColorRect.new()
	_progress_bar.name = "ProgressBar"
	_progress_bar.position = Vector2(550, 50)
	_progress_bar.size = Vector2(180, 6)
	_progress_bar.color = PANEL_COLOR_PINK
	_banner.add_child(_progress_bar)

	# 橫幅滑入動畫
	var tween = _banner.create_tween()
	tween.tween_property(_banner, "position:y", 0.0, 0.3).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)

	# 全螢幕粉紅閃光（自己觸發時更強烈）
	var flash_alpha = 0.4 if _is_my_session else 0.2
	var flash = ColorRect.new()
	flash.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	flash.color = Color(1.0, 0.4, 0.7, 0.0)
	flash.z_index = 87
	canvas_layer.add_child(flash)
	var flash_tween = flash.create_tween()
	flash_tween.tween_property(flash, "color:a", flash_alpha, 0.1)
	flash_tween.tween_property(flash, "color:a", 0.0, 0.4)
	flash_tween.tween_callback(flash.queue_free)

	# 自己觸發時：生成夢幻星星粒子
	if _is_my_session:
		_spawn_dream_particles(canvas_layer)

func _on_end(data: Dictionary) -> void:
	_is_active = false
	var total_reward: int = data.get("total_reward", 0)
	var kill_count: int = data.get("kill_count", 0)

	# 橫幅滑出
	if is_instance_valid(_banner):
		var tween = _banner.create_tween()
		tween.tween_property(_banner, "position:y", -60.0, 0.3).set_trans(Tween.TRANS_QUAD).set_ease(Tween.EASE_IN)
		tween.tween_callback(_banner.queue_free)
	_banner = null
	_timer_label = null
	_progress_bar = null

	# 只有自己觸發且有獎勵時顯示結果彈窗
	if not _is_my_session or total_reward <= 0:
		return

	var canvas_layer = get_parent()
	if not is_instance_valid(canvas_layer):
		return

	# 結果彈窗（右側滑入）
	var panel = Control.new()
	panel.name = "GiantPrizeFishResult"
	panel.position = Vector2(1280, 80)
	panel.size = Vector2(280, 140)
	panel.z_index = 89
	canvas_layer.add_child(panel)

	var bg = ColorRect.new()
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	bg.color = Color(0.1, 0.03, 0.08, 0.95)
	panel.add_child(bg)

	# 粉紅左側邊框
	var left_border = ColorRect.new()
	left_border.size = Vector2(4, 140)
	left_border.color = PANEL_COLOR_PINK
	panel.add_child(left_border)

	# 標題
	var title_lbl = Label.new()
	title_lbl.text = "✨ 夢幻模式結束！"
	title_lbl.position = Vector2(12, 8)
	title_lbl.size = Vector2(260, 22)
	title_lbl.add_theme_font_size_override("font_size", 14)
	title_lbl.add_theme_color_override("font_color", PANEL_COLOR_GOLD)
	if is_instance_valid(_pixel_font):
		title_lbl.add_theme_font_override("font", _pixel_font)
	panel.add_child(title_lbl)

	# 擊破數
	var kill_lbl = Label.new()
	kill_lbl.text = "擊破目標：%d 個" % kill_count
	kill_lbl.position = Vector2(12, 34)
	kill_lbl.size = Vector2(260, 18)
	kill_lbl.add_theme_font_size_override("font_size", 12)
	kill_lbl.add_theme_color_override("font_color", PANEL_COLOR_WHITE)
	if is_instance_valid(_pixel_font):
		kill_lbl.add_theme_font_override("font", _pixel_font)
	panel.add_child(kill_lbl)

	# 倍率說明
	var mult_lbl = Label.new()
	mult_lbl.text = "×%.0f 加成已套用" % _mult_bonus
	mult_lbl.position = Vector2(12, 56)
	mult_lbl.size = Vector2(260, 18)
	mult_lbl.add_theme_font_size_override("font_size", 12)
	mult_lbl.add_theme_color_override("font_color", PANEL_COLOR_PINK)
	if is_instance_valid(_pixel_font):
		mult_lbl.add_theme_font_override("font", _pixel_font)
	panel.add_child(mult_lbl)

	# 總獎勵
	var reward_lbl = Label.new()
	reward_lbl.text = "🪙 +%d" % total_reward
	reward_lbl.position = Vector2(12, 82)
	reward_lbl.size = Vector2(260, 40)
	reward_lbl.add_theme_font_size_override("font_size", 26)
	reward_lbl.add_theme_color_override("font_color", PANEL_COLOR_GOLD)
	if is_instance_valid(_pixel_font):
		reward_lbl.add_theme_font_override("font", _pixel_font)
	panel.add_child(reward_lbl)

	# 粉紅閃光
	var flash = ColorRect.new()
	flash.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	flash.color = Color(1.0, 0.4, 0.7, 0.0)
	flash.z_index = 1
	panel.add_child(flash)
	var flash_tween = flash.create_tween()
	flash_tween.tween_property(flash, "color:a", 0.5, 0.1)
	flash_tween.tween_property(flash, "color:a", 0.0, 0.4)
	flash_tween.tween_callback(flash.queue_free)

	# 滑入動畫
	var tween = panel.create_tween()
	tween.tween_property(panel, "position:x", 990.0, 0.35).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)
	tween.tween_interval(3.5)
	tween.tween_property(panel, "position:x", 1280.0, 0.3).set_trans(Tween.TRANS_QUAD).set_ease(Tween.EASE_IN)
	tween.tween_callback(panel.queue_free)

## _spawn_dream_particles — 生成夢幻星星粒子（浮動文字模擬）
func _spawn_dream_particles(canvas_layer: Node) -> void:
	var stars = ["✨", "⭐", "💫", "🌟", "✨"]
	for i in range(5):
		var star_label = Label.new()
		star_label.text = stars[i % stars.size()]
		star_label.add_theme_font_size_override("font_size", 24 + randi() % 16)
		star_label.position = Vector2(
			randf_range(100, 1180),
			randf_range(100, 620)
		)
		star_label.z_index = 90
		canvas_layer.add_child(star_label)

		var star_tween = star_label.create_tween()
		star_tween.tween_property(star_label, "position:y", star_label.position.y - 80, 1.2)
		star_tween.parallel().tween_property(star_label, "modulate:a", 0.0, 1.2)
		star_tween.tween_callback(star_label.queue_free)

# ---- 常數 ----
const PANEL_COLOR_BG    := Color(0.95, 0.4, 0.7, 0.92)   # 粉紅主題背景
const PANEL_COLOR_GOLD  := Color(1.0, 0.85, 0.0, 1.0)    # 金色文字
const PANEL_COLOR_WHITE := Color(1.0, 1.0, 1.0, 1.0)
const PANEL_COLOR_PINK  := Color(1.0, 0.4, 0.7, 1.0)     # 粉紅色
const PANEL_COLOR_DREAM := Color(0.98, 0.75, 0.9, 1.0)   # 夢幻淡粉

# ---- 節點引用 ----
var _banner: PanelContainer       # 頂部橫幅
var _banner_label: Label          # 橫幅文字
var _mult_label: Label            # ×5 倍率顯示
var _timer_label: Label           # 倒數計時
var _progress_bar: ProgressBar    # 剩餘時間進度條
var _flash_overlay: ColorRect     # 全螢幕閃光
var _result_panel: PanelContainer # 結束結果彈窗
var _result_label: Label          # 結果文字

# ---- 狀態 ----
var _is_active: bool = false
var _is_my_session: bool = false   # 是否是自己觸發的
var _duration: float = 10.0
var _elapsed: float = 0.0
var _mult_bonus: float = 5.0
var _killer_name: String = ""

func _ready() -> void:
	layer = 88  # 在 CrocodilePanel(87) 之上
	_build_ui()
	visible = true

func _build_ui() -> void:
	# 全螢幕閃光遮罩
	_flash_overlay = ColorRect.new()
	_flash_overlay.color = Color(1.0, 0.4, 0.7, 0.0)
	_flash_overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_flash_overlay.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_flash_overlay)

	# 頂部橫幅（從上方滑入）
	_banner = PanelContainer.new()
	var banner_style = StyleBoxFlat.new()
	banner_style.bg_color = PANEL_COLOR_BG
	banner_style.corner_radius_bottom_left = 12
	banner_style.corner_radius_bottom_right = 12
	banner_style.border_width_bottom = 3
	banner_style.border_color = PANEL_COLOR_GOLD
	_banner.add_theme_stylebox_override("panel", banner_style)
	_banner.set_anchors_preset(Control.PRESET_TOP_WIDE)
	_banner.offset_top = -120
	_banner.offset_bottom = -120 + 110
	_banner.offset_left = 200
	_banner.offset_right = -200
	add_child(_banner)

	var banner_vbox = VBoxContainer.new()
	banner_vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	_banner.add_child(banner_vbox)

	_banner_label = Label.new()
	_banner_label.text = "✨ 夢幻獎勵模式！"
	_banner_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner_label.add_theme_color_override("font_color", PANEL_COLOR_GOLD)
	_banner_label.add_theme_font_size_override("font_size", 22)
	banner_vbox.add_child(_banner_label)

	_mult_label = Label.new()
	_mult_label.text = "×5 獎勵加成！"
	_mult_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_mult_label.add_theme_color_override("font_color", PANEL_COLOR_WHITE)
	_mult_label.add_theme_font_size_override("font_size", 18)
	banner_vbox.add_child(_mult_label)

	# 倒數計時 + 進度條（左下角）
	var timer_container = VBoxContainer.new()
	timer_container.set_anchors_preset(Control.PRESET_BOTTOM_LEFT)
	timer_container.offset_left = 20
	timer_container.offset_bottom = -20
	timer_container.offset_top = -80
	timer_container.offset_right = 200
	add_child(timer_container)

	_timer_label = Label.new()
	_timer_label.text = "✨ 10.0s"
	_timer_label.add_theme_color_override("font_color", PANEL_COLOR_PINK)
	_timer_label.add_theme_font_size_override("font_size", 20)
	timer_container.add_child(_timer_label)

	_progress_bar = ProgressBar.new()
	_progress_bar.min_value = 0.0
	_progress_bar.max_value = 1.0
	_progress_bar.value = 1.0
	_progress_bar.custom_minimum_size = Vector2(180, 14)
	var pb_style = StyleBoxFlat.new()
	pb_style.bg_color = PANEL_COLOR_PINK
	pb_style.corner_radius_top_left = 7
	pb_style.corner_radius_top_right = 7
	pb_style.corner_radius_bottom_left = 7
	pb_style.corner_radius_bottom_right = 7
	_progress_bar.add_theme_stylebox_override("fill", pb_style)
	timer_container.add_child(_progress_bar)

	# 結束結果彈窗（右側滑入）
	_result_panel = PanelContainer.new()
	var result_style = StyleBoxFlat.new()
	result_style.bg_color = Color(0.1, 0.05, 0.15, 0.95)
	result_style.corner_radius_top_left = 12
	result_style.corner_radius_bottom_left = 12
	result_style.border_width_left = 4
	result_style.border_color = PANEL_COLOR_PINK
	_result_panel.add_theme_stylebox_override("panel", result_style)
	_result_panel.set_anchors_preset(Control.PRESET_CENTER_RIGHT)
	_result_panel.offset_right = 0
	_result_panel.offset_left = -280
	_result_panel.offset_top = -80
	_result_panel.offset_bottom = 80
	_result_panel.visible = false
	add_child(_result_panel)

	_result_label = Label.new()
	_result_label.add_theme_color_override("font_color", PANEL_COLOR_WHITE)
	_result_label.add_theme_font_size_override("font_size", 16)
	_result_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_result_panel.add_child(_result_label)

	# 初始隱藏
	_banner.modulate.a = 0.0
	timer_container.visible = false
	_timer_container_ref = timer_container

var _timer_container_ref: VBoxContainer = null

func _process(delta: float) -> void:
	if not _is_active:
		return

	_elapsed += delta
	var remaining = max(0.0, _duration - _elapsed)
	var pct = remaining / _duration

	# 更新倒數計時
	_timer_label.text = "✨ %.1fs" % remaining
	_progress_bar.value = pct

	# 最後 3 秒：紅色閃爍
	if remaining <= 3.0:
		var blink = sin(_elapsed * 8.0) > 0.0
		_timer_label.add_theme_color_override("font_color",
			Color(1.0, 0.2, 0.2) if blink else PANEL_COLOR_PINK)

	# 時間到：結束（由 Server 的 "end" 訊息觸發，這裡只是備用）
	if _elapsed >= _duration + 0.5:
		_is_active = false

# ---- 公開 API ----

## on_giant_prize_fish_event — 處理 Server 廣播的夢幻獎勵魚事件
func on_giant_prize_fish_event(data: Dictionary, my_player_id: String) -> void:
	var phase: String = data.get("phase", "")
	var killer_id: String = data.get("killer_id", "")
	_killer_name = data.get("killer_name", "玩家")
	_mult_bonus = data.get("mult_bonus", 5.0)
	_duration = float(data.get("duration", 10))
	_is_my_session = (killer_id == my_player_id)

	match phase:
		"activate":
			_on_activate(data)
		"end":
			_on_end(data)

func _on_activate(data: Dictionary) -> void:
	_is_active = true
	_elapsed = 0.0

	# 更新橫幅文字
	if _is_my_session:
		_banner_label.text = "✨ 夢幻獎勵模式啟動！"
		_mult_label.text = "×%.0f 獎勵加成！10 秒！" % _mult_bonus
	else:
		_banner_label.text = "✨ %s 觸發夢幻獎勵！" % _killer_name
		_mult_label.text = "×%.0f 加成中..." % _mult_bonus

	# 橫幅滑入
	_banner.modulate.a = 0.0
	_banner.offset_top = -120
	_banner.offset_bottom = -120 + 110
	var tween = create_tween()
	tween.tween_property(_banner, "modulate:a", 1.0, 0.3)
	tween.parallel().tween_property(_banner, "offset_top", 0.0, 0.4).set_ease(Tween.EASE_OUT)
	tween.parallel().tween_property(_banner, "offset_bottom", 110.0, 0.4).set_ease(Tween.EASE_OUT)

	# 顯示計時器
	if _timer_container_ref:
		_timer_container_ref.visible = true

	# 全螢幕粉紅閃光（自己觸發時更強烈）
	var flash_alpha = 0.5 if _is_my_session else 0.25
	var flash_tween = create_tween()
	flash_tween.tween_property(_flash_overlay, "color:a", flash_alpha, 0.1)
	flash_tween.tween_property(_flash_overlay, "color:a", 0.0, 0.4)

	# 自己觸發時：額外的星星粒子效果（用多個浮動文字模擬）
	if _is_my_session:
		_spawn_dream_particles()

func _on_end(data: Dictionary) -> void:
	_is_active = false
	var total_reward: int = data.get("total_reward", 0)
	var kill_count: int = data.get("kill_count", 0)

	# 橫幅滑出
	var tween = create_tween()
	tween.tween_property(_banner, "offset_top", -120.0, 0.3).set_ease(Tween.EASE_IN)
	tween.parallel().tween_property(_banner, "offset_bottom", -10.0, 0.3).set_ease(Tween.EASE_IN)
	tween.parallel().tween_property(_banner, "modulate:a", 0.0, 0.3)

	# 隱藏計時器
	if _timer_container_ref:
		_timer_container_ref.visible = false

	# 只有自己觸發時顯示結果彈窗
	if not _is_my_session or total_reward <= 0:
		return

	# 結果彈窗（右側滑入）
	_result_label.text = "✨ 夢幻模式結束！\n擊破 %d 個目標\n獲得 %d 金幣" % [kill_count, total_reward]
	_result_panel.visible = true
	_result_panel.offset_right = 300  # 先在畫面外
	_result_panel.offset_left = 20

	var result_tween = create_tween()
	result_tween.tween_property(_result_panel, "offset_right", 0.0, 0.4).set_ease(Tween.EASE_OUT)
	result_tween.parallel().tween_property(_result_panel, "offset_left", -280.0, 0.4).set_ease(Tween.EASE_OUT)

	# 3 秒後淡出
	await get_tree().create_timer(3.0).timeout
	var fade_tween = create_tween()
	fade_tween.tween_property(_result_panel, "modulate:a", 0.0, 0.5)
	await fade_tween.finished
	_result_panel.visible = false
	_result_panel.modulate.a = 1.0

## _spawn_dream_particles — 生成夢幻星星粒子（浮動文字模擬）
func _spawn_dream_particles() -> void:
	var stars = ["✨", "⭐", "💫", "🌟", "✨"]
	var viewport_size = get_viewport().get_visible_rect().size
	for i in range(5):
		var star_label = Label.new()
		star_label.text = stars[i % stars.size()]
		star_label.add_theme_font_size_override("font_size", 24 + randi() % 16)
		star_label.position = Vector2(
			randf_range(100, viewport_size.x - 100),
			randf_range(100, viewport_size.y - 100)
		)
		add_child(star_label)

		var star_tween = create_tween()
		star_tween.tween_property(star_label, "position:y", star_label.position.y - 80, 1.2)
		star_tween.parallel().tween_property(star_label, "modulate:a", 0.0, 1.2)
		star_tween.tween_callback(star_label.queue_free)
