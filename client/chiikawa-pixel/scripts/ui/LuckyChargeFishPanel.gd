## LuckyChargeFishPanel.gd — 幸運充能魚系統面板（DAY-225）
## 業界原創「射擊充能→爆發」機制
##
## 視覺設計：
##   - 黃橙充能主題（#F39C12 + #E67E22 + #F1C40F + #FEF9E7）
##   - charge_start：黃橙雙閃光 + 右側充能進度條（豎向）+ 計時條
##   - charge_progress：進度條更新 + 閃光 + 「還差 N 槍」提示
##   - charge_ready：全螢幕黃色強閃光 + 「⚡ 充能就緒！」大字 + 閃爍提示
##   - charge_burst：全螢幕三次黃橙強閃光 + 「⚡ ×5.0 爆發！」大字
##   - charge_end：進度條淡出
extends CanvasLayer

# 充能狀態
var _active: bool = false
var _charge_count: int = 0
var _charge_target: int = 10
var _burst_ready: bool = false
var _charge_bar: Control = null
var _timer_bar: Control = null
var _ready_label: Control = null

# 主題顏色
const COLOR_PRIMARY   = Color("#F39C12")  # 黃橙
const COLOR_DARK      = Color("#E67E22")  # 深橙
const COLOR_GOLD      = Color("#F1C40F")  # 金黃
const COLOR_PALE      = Color("#FEF9E7")  # 極淡黃
const COLOR_BG        = Color(0.1, 0.06, 0.0, 0.88)

func _ready() -> void:
	layer = 20  # 幸運充能魚面板層級

## 處理幸運充能魚訊息
func handle_lucky_charge_fish(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"charge_start":
			_on_charge_start(payload)
		"charge_progress":
			_on_charge_progress(payload)
		"charge_ready":
			_on_charge_ready(payload)
		"charge_burst":
			_on_charge_burst(payload)
		"charge_end":
			_on_charge_end(payload)
		"charge_broadcast":
			_on_charge_broadcast(payload)
		"charge_burst_broadcast":
			_on_charge_burst_broadcast(payload)

## charge_start — 充能模式開始（個人）
func _on_charge_start(payload: Dictionary) -> void:
	var duration_sec: int = payload.get("duration_sec", 12)
	var target: int = payload.get("target", 10)
	var burst_mult: float = payload.get("burst_mult", 5.0)

	_active = true
	_charge_count = 0
	_charge_target = target
	_burst_ready = false

	# 黃橙雙閃光
	_flash_screen(COLOR_PRIMARY, 0.18)
	await get_tree().create_timer(0.1).timeout
	_flash_screen(COLOR_GOLD, 0.15)

	# 右側充能進度條（豎向）
	_show_charge_bar(target, burst_mult)

	# 底部計時條
	_show_timer_bar(duration_sec)

## charge_progress — 充能進度（個人）
func _on_charge_progress(payload: Dictionary) -> void:
	var count: int = payload.get("count", 0)
	var target: int = payload.get("target", 10)

	_charge_count = count

	# 更新充能進度條
	if _charge_bar != null and is_instance_valid(_charge_bar):
		var pct: float = float(count) / float(target)
		var bar_fill = _charge_bar.get_node_or_null("BarFill")
		if bar_fill != null:
			var tween = bar_fill.create_tween()
			var vp_h = get_viewport().size.y * 0.5
			tween.tween_property(bar_fill, "size:y", vp_h * pct, 0.15)
			# 顏色漸變（黃橙→金黃）
			tween.parallel().tween_property(bar_fill, "color", COLOR_GOLD.lerp(COLOR_PRIMARY, pct), 0.15)

	# 閃光
	_flash_screen(COLOR_PRIMARY, 0.08)

	# 「還差 N 槍」提示
	var remaining: int = target - count
	if remaining > 0:
		var vp_size = get_viewport().size
		var hint = Label.new()
		hint.text = "⚡ 還差 %d 槍！" % remaining
		hint.add_theme_font_size_override("font_size", 13)
		hint.add_theme_color_override("font_color", COLOR_GOLD)
		hint.position = Vector2(vp_size.x - 120, vp_size.y / 2 - 10)
		add_child(hint)

		var tween = hint.create_tween()
		tween.tween_property(hint, "modulate:a", 0.0, 0.6)
		tween.tween_callback(hint.queue_free)

## charge_ready — 充能爆發就緒（個人）
func _on_charge_ready(payload: Dictionary) -> void:
	var burst_mult: float = payload.get("burst_mult", 5.0)
	_burst_ready = true

	# 全螢幕黃色強閃光
	_flash_screen(COLOR_GOLD, 0.25)

	# 「⚡ 充能就緒！」大字
	var vp_size = get_viewport().size
	var big_label = Label.new()
	big_label.text = "⚡ 充能就緒！×%.1f" % burst_mult
	big_label.add_theme_font_size_override("font_size", 46)
	big_label.add_theme_color_override("font_color", COLOR_GOLD)
	big_label.position = vp_size / 2 - Vector2(150, 28)
	add_child(big_label)

	var tween_label = big_label.create_tween()
	tween_label.tween_property(big_label, "scale", Vector2(1.2, 1.2), 0.12)
	tween_label.tween_property(big_label, "scale", Vector2(1.0, 1.0), 0.08)
	tween_label.tween_interval(0.4)
	tween_label.tween_property(big_label, "modulate:a", 0.0, 0.4)
	tween_label.tween_callback(big_label.queue_free)

	# 右側充能條閃爍提示
	if _charge_bar != null and is_instance_valid(_charge_bar):
		var tween_bar = _charge_bar.create_tween().set_loops(6)
		tween_bar.tween_property(_charge_bar, "modulate", Color(2.0, 2.0, 0.5, 1.0), 0.15)
		tween_bar.tween_property(_charge_bar, "modulate", Color.WHITE, 0.15)

## charge_burst — 充能爆發觸發（個人）
func _on_charge_burst(payload: Dictionary) -> void:
	var burst_mult: float = payload.get("burst_mult", 5.0)
	var reward: int = payload.get("reward", 0)
	_burst_ready = false
	_charge_count = 0

	# 全螢幕三次黃橙強閃光
	_flash_screen(COLOR_GOLD, 0.15)
	await get_tree().create_timer(0.1).timeout
	_flash_screen(COLOR_PRIMARY, 0.12)
	await get_tree().create_timer(0.1).timeout
	_flash_screen(Color.WHITE, 0.1)

	# 「⚡ ×5.0 爆發！」大字
	var vp_size = get_viewport().size
	var big_label = Label.new()
	big_label.text = "⚡ ×%.1f 爆發！" % burst_mult
	big_label.add_theme_font_size_override("font_size", 52)
	big_label.add_theme_color_override("font_color", COLOR_GOLD)
	big_label.position = vp_size / 2 - Vector2(140, 30)
	add_child(big_label)

	var tween_label = big_label.create_tween()
	tween_label.tween_property(big_label, "scale", Vector2(1.3, 1.3), 0.12)
	tween_label.tween_property(big_label, "scale", Vector2(1.0, 1.0), 0.08)
	tween_label.tween_interval(0.5)
	tween_label.tween_property(big_label, "modulate:a", 0.0, 0.4)
	tween_label.tween_callback(big_label.queue_free)

	# 獎勵浮動文字
	if reward > 0:
		var reward_label = Label.new()
		reward_label.text = "+%d 💰" % reward
		reward_label.add_theme_font_size_override("font_size", 20)
		reward_label.add_theme_color_override("font_color", COLOR_PALE)
		reward_label.position = vp_size / 2 - Vector2(40, -20)
		add_child(reward_label)

		var tween_r = reward_label.create_tween()
		tween_r.tween_property(reward_label, "position:y", reward_label.position.y - 40, 0.6)
		tween_r.parallel().tween_property(reward_label, "modulate:a", 0.0, 0.6)
		tween_r.tween_callback(reward_label.queue_free)

	# 重置充能條
	if _charge_bar != null and is_instance_valid(_charge_bar):
		var bar_fill = _charge_bar.get_node_or_null("BarFill")
		if bar_fill != null:
			var tween_reset = bar_fill.create_tween()
			tween_reset.tween_property(bar_fill, "size:y", 0.0, 0.2)
			tween_reset.tween_property(bar_fill, "color", COLOR_PRIMARY, 0.1)

## charge_end — 充能模式結束（個人）
func _on_charge_end(_payload: Dictionary) -> void:
	_active = false
	_burst_ready = false
	_charge_count = 0

	# 淡出充能條
	if _charge_bar != null and is_instance_valid(_charge_bar):
		var tween = _charge_bar.create_tween()
		tween.tween_property(_charge_bar, "modulate:a", 0.0, 0.4)
		tween.tween_callback(_charge_bar.queue_free)
		_charge_bar = null

	# 淡出計時條
	if _timer_bar != null and is_instance_valid(_timer_bar):
		var tween = _timer_bar.create_tween()
		tween.tween_property(_timer_bar, "modulate:a", 0.0, 0.3)
		tween.tween_callback(_timer_bar.queue_free)
		_timer_bar = null

## charge_broadcast — 充能模式開始廣播（全服）
func _on_charge_broadcast(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "")
	var vp_size = get_viewport().size

	# 頂部小橫幅（其他玩家看到）
	var banner = Label.new()
	banner.text = "⚡ %s 進入充能模式！" % player_name
	banner.add_theme_font_size_override("font_size", 13)
	banner.add_theme_color_override("font_color", COLOR_PRIMARY)
	banner.position = Vector2(vp_size.x / 2 - 100, 4)
	add_child(banner)

	var tween = banner.create_tween()
	tween.tween_interval(2.0)
	tween.tween_property(banner, "modulate:a", 0.0, 0.4)
	tween.tween_callback(banner.queue_free)

## charge_burst_broadcast — 充能爆發廣播（全服）
func _on_charge_burst_broadcast(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "")
	var burst_mult: float = payload.get("burst_mult", 5.0)
	var vp_size = get_viewport().size

	# 頂部小橫幅（其他玩家看到）
	var banner = Label.new()
	banner.text = "⚡ %s 充能爆發 ×%.1f！" % [player_name, burst_mult]
	banner.add_theme_font_size_override("font_size", 14)
	banner.add_theme_color_override("font_color", COLOR_GOLD)
	banner.position = Vector2(vp_size.x / 2 - 110, 4)
	add_child(banner)

	var tween = banner.create_tween()
	tween.tween_interval(2.5)
	tween.tween_property(banner, "modulate:a", 0.0, 0.4)
	tween.tween_callback(banner.queue_free)

# ---- 輔助函數 ----

## 顯示右側充能進度條（豎向）
func _show_charge_bar(target: int, burst_mult: float) -> void:
	if _charge_bar != null and is_instance_valid(_charge_bar):
		_charge_bar.queue_free()

	var vp_size = get_viewport().size
	var bar_h = vp_size.y * 0.5
	var bar_container = Control.new()
	bar_container.position = Vector2(vp_size.x - 28, vp_size.y / 2 - bar_h / 2)
	bar_container.size = Vector2(24, bar_h)
	add_child(bar_container)
	_charge_bar = bar_container

	# 背景
	var bar_bg = ColorRect.new()
	bar_bg.color = Color(0.1, 0.06, 0.0, 0.8)
	bar_bg.size = bar_container.size
	bar_container.add_child(bar_bg)

	# 填充（從底部往上）
	var bar_fill = ColorRect.new()
	bar_fill.name = "BarFill"
	bar_fill.color = COLOR_PRIMARY
	bar_fill.size = Vector2(24, 0)
	bar_fill.position = Vector2(0, bar_h)  # 從底部開始
	bar_container.add_child(bar_fill)

	# 標籤（頂部）
	var label = Label.new()
	label.text = "⚡\n×%.0f" % burst_mult
	label.add_theme_font_size_override("font_size", 10)
	label.add_theme_color_override("font_color", COLOR_GOLD)
	label.position = Vector2(-2, -28)
	bar_container.add_child(label)

## 顯示底部計時條
func _show_timer_bar(duration_sec: int) -> void:
	if _timer_bar != null and is_instance_valid(_timer_bar):
		_timer_bar.queue_free()

	var vp_size = get_viewport().size
	var timer_container = Control.new()
	timer_container.position = Vector2(0, vp_size.y - 6)
	timer_container.size = Vector2(vp_size.x, 6)
	add_child(timer_container)
	_timer_bar = timer_container

	var bar = ColorRect.new()
	bar.color = COLOR_PRIMARY
	bar.size = Vector2(vp_size.x, 6)
	timer_container.add_child(bar)

	var tween = bar.create_tween()
	tween.tween_property(bar, "size:x", 0.0, float(duration_sec))
	tween.parallel().tween_property(bar, "color", COLOR_DARK, float(duration_sec))

## 全螢幕閃光效果
func _flash_screen(color: Color, duration: float) -> void:
	var flash = ColorRect.new()
	flash.color = Color(color.r, color.g, color.b, 0.35)
	flash.set_anchors_preset(Control.PRESET_FULL_RECT)
	add_child(flash)

	var tween = flash.create_tween()
	tween.tween_property(flash, "modulate:a", 0.0, duration)
	tween.tween_callback(flash.queue_free)
