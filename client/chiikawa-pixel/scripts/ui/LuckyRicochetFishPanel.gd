## LuckyRicochetFishPanel.gd — 幸運反彈魚系統面板（DAY-220）
## 業界原創「子彈反彈」機制
##
## 視覺設計：
##   - 橙色反彈主題（#FF8C00 + #FFD700 + #FF4500 + #FFF8DC）
##   - ricochet_start：橙色雙閃光 + 頂部橫幅 + 底部計時條
##   - ricochet_bounce：反彈軌跡線 + 爆炸圓圈 + 浮動獎勵文字
##   - ricochet_end：橙色淡出 + 「反彈模式結束」提示
extends CanvasLayer

# 反彈狀態
var _ricochet_active: bool = false
var _my_player_id: String = ""
var _timer_bar: Control = null
var _banner: Control = null

# 主題顏色
const COLOR_PRIMARY   = Color("#FF8C00")  # 橙色
const COLOR_GOLD      = Color("#FFD700")  # 金色
const COLOR_ACCENT    = Color("#FF4500")  # 深橙
const COLOR_LIGHT     = Color("#FFF8DC")  # 淡黃
const COLOR_BG        = Color(0.15, 0.08, 0.0, 0.85)

func _ready() -> void:
	layer = 25  # 幸運反彈魚面板層級
	# 取得本地玩家 ID
	if GameManager.has_method("get_player_id"):
		_my_player_id = GameManager.get_player_id()

## 處理幸運反彈魚訊息
func handle_lucky_ricochet_fish(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"ricochet_start":
			_on_ricochet_start(payload)
		"ricochet_bounce":
			_on_ricochet_bounce(payload)
		"ricochet_end":
			_on_ricochet_end(payload)

## ricochet_start — 反彈模式開始
func _on_ricochet_start(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "")
	var duration_sec: int = payload.get("duration_sec", 8)
	var player_id: String = payload.get("player_id", "")

	_ricochet_active = true

	# 橙色雙閃光
	_flash_screen(COLOR_PRIMARY, 0.25)
	await get_tree().create_timer(0.15).timeout
	_flash_screen(COLOR_GOLD, 0.2)

	# 頂部橫幅
	_show_banner("🎯 反彈模式！", "%s 觸發反彈魚！每槍最多反彈 3 次" % player_name, duration_sec)

## ricochet_bounce — 子彈反彈命中
func _on_ricochet_bounce(payload: Dictionary) -> void:
	var bounce_num: int = payload.get("bounce_num", 1)
	var killed: bool = payload.get("killed", false)
	var reward: int = payload.get("reward", 0)
	var x: float = payload.get("x", 0.0)
	var y: float = payload.get("y", 0.0)

	# 爆炸圓圈（依反彈次數決定大小）
	var radius: float = 30.0 - float(bounce_num) * 5.0
	_show_bounce_effect(Vector2(x, y), radius, bounce_num)

	# 浮動文字（只有擊破才顯示獎勵）
	if killed and reward > 0:
		var color = COLOR_GOLD if bounce_num == 1 else COLOR_PRIMARY
		_show_float_text("↩ +%d" % reward, color, Vector2(x, y))

## ricochet_end — 反彈模式結束
func _on_ricochet_end(payload: Dictionary) -> void:
	_ricochet_active = false
	_hide_banner()

	# 橙色淡出提示
	var label = Label.new()
	label.text = "🎯 反彈結束"
	label.add_theme_font_size_override("font_size", 16)
	label.add_theme_color_override("font_color", COLOR_PRIMARY)
	label.set_anchors_preset(Control.PRESET_CENTER)
	label.position = get_viewport().size / 2 - Vector2(60, 20)
	add_child(label)

	var tween = label.create_tween()
	tween.tween_property(label, "modulate:a", 0.0, 0.8)
	tween.tween_callback(label.queue_free)

# ---- 輔助函數 ----

## 顯示頂部橫幅
func _show_banner(title: String, subtitle: String, duration_sec: int) -> void:
	_hide_banner()

	var banner = Control.new()
	banner.set_anchors_preset(Control.PRESET_TOP_WIDE)
	banner.position = Vector2(0, 8)
	banner.size = Vector2(get_viewport().size.x, 52)
	add_child(banner)

	var bg = ColorRect.new()
	bg.color = COLOR_BG
	bg.size = banner.size
	banner.add_child(bg)

	var title_label = Label.new()
	title_label.text = title
	title_label.add_theme_font_size_override("font_size", 20)
	title_label.add_theme_color_override("font_color", COLOR_PRIMARY)
	title_label.position = Vector2(12, 4)
	banner.add_child(title_label)

	var sub_label = Label.new()
	sub_label.text = subtitle
	sub_label.add_theme_font_size_override("font_size", 13)
	sub_label.add_theme_color_override("font_color", COLOR_LIGHT)
	sub_label.position = Vector2(12, 28)
	banner.add_child(sub_label)

	# 計時條（底部，橙→紅橙漸變）
	var timer_bar = ColorRect.new()
	timer_bar.name = "TimerBar"
	timer_bar.color = COLOR_PRIMARY
	timer_bar.position = Vector2(0, 48)
	timer_bar.size = Vector2(get_viewport().size.x, 4)
	banner.add_child(timer_bar)

	var tween = banner.create_tween()
	tween.tween_property(timer_bar, "size:x", 0.0, float(duration_sec))
	tween.parallel().tween_property(timer_bar, "color", COLOR_ACCENT, float(duration_sec))

	_banner = banner

## 隱藏橫幅
func _hide_banner() -> void:
	if _banner != null and is_instance_valid(_banner):
		_banner.queue_free()
	_banner = null

## 反彈爆炸效果（擴散圓圈）
func _show_bounce_effect(pos: Vector2, radius: float, bounce_num: int) -> void:
	# 外圈（橙色）
	var outer = ColorRect.new()
	outer.color = Color(COLOR_PRIMARY.r, COLOR_PRIMARY.g, COLOR_PRIMARY.b, 0.6)
	outer.size = Vector2(radius * 2, radius * 2)
	outer.position = pos - Vector2(radius, radius)
	add_child(outer)

	# 反彈次數標記
	var num_label = Label.new()
	num_label.text = "↩%d" % bounce_num
	num_label.add_theme_font_size_override("font_size", 12)
	num_label.add_theme_color_override("font_color", COLOR_GOLD)
	num_label.position = pos - Vector2(10, 8)
	add_child(num_label)

	# 擴散動畫
	var tween = outer.create_tween()
	tween.tween_property(outer, "size", Vector2(radius * 4, radius * 4), 0.3)
	tween.parallel().tween_property(outer, "position", pos - Vector2(radius * 2, radius * 2), 0.3)
	tween.parallel().tween_property(outer, "modulate:a", 0.0, 0.3)
	tween.tween_callback(outer.queue_free)

	var tween2 = num_label.create_tween()
	tween2.tween_property(num_label, "modulate:a", 0.0, 0.4)
	tween2.tween_callback(num_label.queue_free)

## 全螢幕閃光效果
func _flash_screen(color: Color, duration: float) -> void:
	var flash = ColorRect.new()
	flash.color = Color(color.r, color.g, color.b, 0.4)
	flash.set_anchors_preset(Control.PRESET_FULL_RECT)
	add_child(flash)

	var tween = flash.create_tween()
	tween.tween_property(flash, "modulate:a", 0.0, duration)
	tween.tween_callback(flash.queue_free)

## 浮動文字
func _show_float_text(text: String, color: Color, pos: Vector2) -> void:
	var label = Label.new()
	label.text = text
	label.add_theme_font_size_override("font_size", 16)
	label.add_theme_color_override("font_color", color)
	label.position = pos - Vector2(30, 15)
	add_child(label)

	var tween = label.create_tween()
	tween.tween_property(label, "position:y", label.position.y - 35, 0.7)
	tween.parallel().tween_property(label, "modulate:a", 0.0, 0.7)
	tween.tween_callback(label.queue_free)
