## LuckyStormFishPanel.gd — 幸運風暴魚系統面板（DAY-230）
## 業界原創「風暴旋轉+位置混亂」機制
##
## 視覺設計：
##   - 青綠風暴主題（#1ABC9C + #16A085 + #A3E4D7 + #F0FFF4）
##   - storm_start：青綠三次強閃光 + 頂部橫幅 + 「🌪️ 風暴來襲！」大字 + 風暴圓圈 + 計時條
##   - storm_rotate：全螢幕青綠閃光 + 「🌪️ 旋轉！」提示 + 移動目標數
##   - storm_blast：全螢幕三次強閃光 + 「🌪️ 風暴爆發！」大字 + 結算彈窗
extends CanvasLayer

# 風暴狀態
var _active: bool = false
var _instance_id: String = ""
var _storm_x: float = 500.0
var _storm_y: float = 300.0
var _radius: float = 320.0
var _duration_sec: int = 10
var _kill_mult: float = 2.5

# 計時條節點
var _timer_bar: ColorRect = null
var _timer_bar_bg: ColorRect = null
var _storm_circle: Control = null

# 主題顏色
const COLOR_PRIMARY   = Color("#1ABC9C")  # 青綠
const COLOR_DARK      = Color("#16A085")  # 深青綠
const COLOR_PALE      = Color("#A3E4D7")  # 淡青綠
const COLOR_LIGHT_BG  = Color("#F0FFF4")  # 極淡綠
const COLOR_GOLD      = Color("#FFD700")  # 金黃
const COLOR_BLAST     = Color("#00FF7F")  # 爆發亮綠

func _ready() -> void:
	layer = 15  # 幸運風暴魚面板層級

## 處理幸運風暴魚訊息
func handle_lucky_storm_fish(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"storm_start":
			_on_storm_start(payload)
		"storm_rotate":
			_on_storm_rotate(payload)
		"storm_blast":
			_on_storm_blast(payload)

## storm_start — 風暴建立
func _on_storm_start(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "")
	_instance_id = payload.get("instance_id", "")
	_storm_x = payload.get("storm_x", 500.0)
	_storm_y = payload.get("storm_y", 300.0)
	_radius = payload.get("radius", 320.0)
	_duration_sec = payload.get("duration_sec", 10)
	_kill_mult = payload.get("kill_mult", 2.5)
	_active = true

	var vp_size = get_viewport().size

	# 青綠三次強閃光
	_flash_screen(COLOR_PRIMARY, 0.18)
	await get_tree().create_timer(0.12).timeout
	_flash_screen(COLOR_DARK, 0.15)
	await get_tree().create_timer(0.10).timeout
	_flash_screen(COLOR_PRIMARY, 0.12)

	# 頂部橫幅
	var banner = Label.new()
	banner.text = "🌪️ %s 觸發風暴！範圍內目標擊破 ×%.1f！" % [player_name, _kill_mult]
	banner.add_theme_font_size_override("font_size", 14)
	banner.add_theme_color_override("font_color", COLOR_PALE)
	banner.position = Vector2(vp_size.x / 2 - 170, 6)
	add_child(banner)

	var tween_banner = banner.create_tween()
	tween_banner.tween_interval(float(_duration_sec) - 0.5)
	tween_banner.tween_property(banner, "modulate:a", 0.0, 0.5)
	tween_banner.tween_callback(banner.queue_free)

	# 「🌪️ 風暴來襲！」大字
	var big_label = Label.new()
	big_label.text = "🌪️ 風暴來襲！"
	big_label.add_theme_font_size_override("font_size", 52)
	big_label.add_theme_color_override("font_color", COLOR_PRIMARY)
	big_label.position = vp_size / 2 - Vector2(120, 30)
	add_child(big_label)

	var tween_big = big_label.create_tween()
	tween_big.tween_property(big_label, "scale", Vector2(1.25, 1.25), 0.12)
	tween_big.tween_property(big_label, "scale", Vector2(1.0, 1.0), 0.10)
	tween_big.tween_interval(0.5)
	tween_big.tween_property(big_label, "modulate:a", 0.0, 0.5)
	tween_big.tween_callback(big_label.queue_free)

	# 倍率提示
	var mult_label = Label.new()
	mult_label.text = "風暴範圍內擊破 ×%.1f！" % _kill_mult
	mult_label.add_theme_font_size_override("font_size", 14)
	mult_label.add_theme_color_override("font_color", COLOR_GOLD)
	mult_label.position = Vector2(vp_size.x / 2 - 80, vp_size.y / 2 + 32)
	add_child(mult_label)

	var tween_mult = mult_label.create_tween()
	tween_mult.tween_interval(2.5)
	tween_mult.tween_property(mult_label, "modulate:a", 0.0, 0.5)
	tween_mult.tween_callback(mult_label.queue_free)

	# 風暴圓圈（在遊戲座標系上顯示風暴範圍）
	_spawn_storm_circle()

	# 底部計時條
	_spawn_timer_bar(float(_duration_sec))

## storm_rotate — 風暴旋轉
func _on_storm_rotate(payload: Dictionary) -> void:
	var rotate_count: int = payload.get("rotate_count", 1)
	var moved_count: int = payload.get("moved_count", 0)

	if moved_count == 0:
		return

	var vp_size = get_viewport().size

	# 青綠閃光
	_flash_screen(COLOR_PRIMARY, 0.10)

	# 「🌪️ 旋轉！」提示
	var rotate_label = Label.new()
	rotate_label.text = "🌪️ 旋轉！（%d 個目標移動）" % moved_count
	rotate_label.add_theme_font_size_override("font_size", 16)
	rotate_label.add_theme_color_override("font_color", COLOR_PALE)
	rotate_label.position = Vector2(vp_size.x / 2 - 90, vp_size.y / 2 - 12)
	add_child(rotate_label)

	var tween_r = rotate_label.create_tween()
	tween_r.tween_property(rotate_label, "scale", Vector2(1.1, 1.1), 0.07)
	tween_r.tween_property(rotate_label, "scale", Vector2(1.0, 1.0), 0.06)
	tween_r.tween_interval(0.4)
	tween_r.tween_property(rotate_label, "modulate:a", 0.0, 0.35)
	tween_r.tween_callback(rotate_label.queue_free)

	# 風暴圓圈旋轉動畫
	if is_instance_valid(_storm_circle):
		var tween_circle = _storm_circle.create_tween()
		tween_circle.tween_property(_storm_circle, "rotation", _storm_circle.rotation + PI * 0.5, 0.3)

## storm_blast — 風暴爆發
func _on_storm_blast(payload: Dictionary) -> void:
	var killed_count: int = payload.get("killed_count", 0)
	var total_reward: int = payload.get("total_reward", 0)
	var blast_mult: float = payload.get("blast_mult", 0.75)

	_active = false

	# 清除計時條和風暴圓圈
	if is_instance_valid(_timer_bar):
		_timer_bar.queue_free()
		_timer_bar = null
	if is_instance_valid(_timer_bar_bg):
		_timer_bar_bg.queue_free()
		_timer_bar_bg = null
	if is_instance_valid(_storm_circle):
		_storm_circle.queue_free()
		_storm_circle = null

	var vp_size = get_viewport().size

	# 全螢幕三次強閃光（青綠→亮綠→青綠）
	_flash_screen(COLOR_PRIMARY, 0.20)
	await get_tree().create_timer(0.14).timeout
	_flash_screen(COLOR_BLAST, 0.18)
	await get_tree().create_timer(0.12).timeout
	_flash_screen(COLOR_PRIMARY, 0.15)

	# 「🌪️ 風暴爆發！」大字
	var blast_label = Label.new()
	blast_label.text = "🌪️ 風暴爆發！"
	blast_label.add_theme_font_size_override("font_size", 52)
	blast_label.add_theme_color_override("font_color", COLOR_BLAST)
	blast_label.position = vp_size / 2 - Vector2(120, 30)
	add_child(blast_label)

	var tween_blast = blast_label.create_tween()
	tween_blast.tween_property(blast_label, "scale", Vector2(1.3, 1.3), 0.14)
	tween_blast.tween_property(blast_label, "scale", Vector2(1.0, 1.0), 0.10)
	tween_blast.tween_interval(0.6)
	tween_blast.tween_property(blast_label, "modulate:a", 0.0, 0.5)
	tween_blast.tween_callback(blast_label.queue_free)

	# 結算彈窗（右側滑入）
	if killed_count > 0:
		await get_tree().create_timer(0.5).timeout
		_show_blast_result(killed_count, total_reward, blast_mult)

## 顯示風暴爆發結算彈窗
func _show_blast_result(killed_count: int, total_reward: int, blast_mult: float) -> void:
	var vp_size = get_viewport().size

	var panel = PanelContainer.new()
	panel.position = Vector2(vp_size.x + 10, vp_size.y / 2 - 60)
	panel.custom_minimum_size = Vector2(200, 120)
	add_child(panel)

	var vbox = VBoxContainer.new()
	panel.add_child(vbox)

	var title_lbl = Label.new()
	title_lbl.text = "🌪️ 風暴爆發結算"
	title_lbl.add_theme_font_size_override("font_size", 14)
	title_lbl.add_theme_color_override("font_color", COLOR_BLAST)
	vbox.add_child(title_lbl)

	var killed_lbl = Label.new()
	killed_lbl.text = "擊破目標：%d 個" % killed_count
	killed_lbl.add_theme_font_size_override("font_size", 13)
	killed_lbl.add_theme_color_override("font_color", COLOR_PALE)
	vbox.add_child(killed_lbl)

	var reward_lbl = Label.new()
	reward_lbl.text = "總獎勵：%d" % total_reward
	reward_lbl.add_theme_font_size_override("font_size", 13)
	reward_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	vbox.add_child(reward_lbl)

	var mult_lbl = Label.new()
	mult_lbl.text = "爆發倍率：×%.2f" % blast_mult
	mult_lbl.add_theme_font_size_override("font_size", 12)
	mult_lbl.add_theme_color_override("font_color", COLOR_PRIMARY)
	vbox.add_child(mult_lbl)

	# 右側滑入
	var tween_slide = panel.create_tween()
	tween_slide.tween_property(panel, "position:x", vp_size.x - 220.0, 0.3)
	tween_slide.tween_interval(2.5)
	tween_slide.tween_property(panel, "modulate:a", 0.0, 0.5)
	tween_slide.tween_callback(panel.queue_free)

## 建立風暴圓圈視覺
func _spawn_storm_circle() -> void:
	var vp_size = get_viewport().size

	# 用 Control 節點繪製圓圈（三層同心圓）
	var circle_container = Control.new()
	circle_container.set_anchors_preset(Control.PRESET_FULL_RECT)
	add_child(circle_container)
	_storm_circle = circle_container

	# 外圈（淡青綠，半透明）
	var outer = ColorRect.new()
	outer.color = Color(COLOR_PRIMARY.r, COLOR_PRIMARY.g, COLOR_PRIMARY.b, 0.08)
	var outer_size = _radius * 2.0
	outer.size = Vector2(outer_size, outer_size)
	outer.position = Vector2(_storm_x - _radius, _storm_y - _radius)
	circle_container.add_child(outer)

	# 中圈（青綠，半透明）
	var mid = ColorRect.new()
	mid.color = Color(COLOR_PRIMARY.r, COLOR_PRIMARY.g, COLOR_PRIMARY.b, 0.12)
	var mid_size = _radius * 1.4
	mid.size = Vector2(mid_size, mid_size)
	mid.position = Vector2(_storm_x - mid_size / 2, _storm_y - mid_size / 2)
	circle_container.add_child(mid)

	# 旋轉動畫（持續旋轉）
	var tween_rot = circle_container.create_tween().set_loops()
	tween_rot.tween_property(circle_container, "rotation", TAU, float(_duration_sec))

	# 持續時間後淡出
	var tween_fade = circle_container.create_tween()
	tween_fade.tween_interval(float(_duration_sec) - 0.5)
	tween_fade.tween_property(circle_container, "modulate:a", 0.0, 0.5)
	tween_fade.tween_callback(circle_container.queue_free)

## 建立底部計時條
func _spawn_timer_bar(duration: float) -> void:
	var vp_size = get_viewport().size

	# 背景條
	var bg = ColorRect.new()
	bg.color = Color(0.1, 0.1, 0.1, 0.6)
	bg.size = Vector2(vp_size.x, 8)
	bg.position = Vector2(0, vp_size.y - 8)
	add_child(bg)
	_timer_bar_bg = bg

	# 計時條（青綠→深青綠漸變）
	var bar = ColorRect.new()
	bar.color = COLOR_PRIMARY
	bar.size = Vector2(vp_size.x, 8)
	bar.position = Vector2(0, vp_size.y - 8)
	add_child(bar)
	_timer_bar = bar

	# 計時條縮短動畫
	var tween_bar = bar.create_tween()
	tween_bar.tween_property(bar, "size:x", 0.0, duration)
	tween_bar.tween_callback(func():
		if is_instance_valid(bar):
			bar.queue_free()
		if is_instance_valid(bg):
			bg.queue_free()
	)

# ---- 輔助函數 ----

## 全螢幕閃光效果
func _flash_screen(color: Color, duration: float) -> void:
	var flash = ColorRect.new()
	flash.color = Color(color.r, color.g, color.b, 0.30)
	flash.set_anchors_preset(Control.PRESET_FULL_RECT)
	add_child(flash)

	var tween = flash.create_tween()
	tween.tween_property(flash, "modulate:a", 0.0, duration)
	tween.tween_callback(flash.queue_free)
