## LuckyProphecyFishPanel.gd — 幸運預言魚系統面板（DAY-243）
## 業界原創「預言指定目標」機制
##
## 視覺設計：
##   - 紫金預言主題（#9B59B6 + #F39C12 + #D7BDE2 + #FFF9E6）
##   - prophecy_start：紫色三次強閃光 + 頂部橫幅 + 「🔮 預言降臨！」大字 + 目標標記 + 計時條
##   - prophecy_broadcast：頂部小橫幅（通知全服有人觸發預言）
##   - prophecy_fulfilled：金色三次強閃光 + 「🔮 預言成真！×3.5」大字 + 結算彈窗
##   - prophecy_broadcast_fulfilled：頂部小橫幅（通知全服預言成真）
##   - prophecy_transfer：橙色閃光 + 「🔮 預言轉移！」提示 + 新目標標記
##   - prophecy_broadcast_transfer：頂部小橫幅（通知全服預言轉移）
##   - prophecy_fail：灰色閃光 + 「🔮 預言失敗！HP-20%」提示
extends CanvasLayer

# 主題顏色
const COLOR_PROPHECY = Color("#9B59B6")  # 紫色（主題）
const COLOR_GOLD     = Color("#F39C12")  # 金色（成真）
const COLOR_ORANGE   = Color("#E67E22")  # 橙色（轉移）
const COLOR_PALE     = Color("#FFF9E6")  # 極淡黃
const COLOR_FAIL     = Color("#7F8C8D")  # 灰色（失敗）
const COLOR_WHITE    = Color("#FFFFFF")  # 白色
const COLOR_MARK     = Color("#E74C3C")  # 紅色（目標標記）

# 計時條節點
var _timer_bar: ColorRect = null
var _timer_bar_bg: ColorRect = null
var _timer_tween: Tween = null
var _duration_sec: int = 12

# 目標標記節點（顯示在預言目標上方）
var _target_marker: Control = null
var _current_target_id: String = ""

func _ready() -> void:
	layer = 2  # 幸運預言魚面板層級

## 處理幸運預言魚訊息
func handle_lucky_prophecy_fish(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"prophecy_start":
			_on_prophecy_start(payload)
		"prophecy_broadcast":
			_on_prophecy_broadcast(payload)
		"prophecy_fulfilled":
			_on_prophecy_fulfilled(payload)
		"prophecy_broadcast_fulfilled":
			_on_prophecy_broadcast_fulfilled(payload)
		"prophecy_transfer":
			_on_prophecy_transfer(payload)
		"prophecy_broadcast_transfer":
			_on_prophecy_broadcast_transfer(payload)
		"prophecy_fail":
			_on_prophecy_fail(payload)

## prophecy_start — 預言開始（個人訊息）
func _on_prophecy_start(payload: Dictionary) -> void:
	_duration_sec = payload.get("duration_sec", 12)
	var kill_mult: float = payload.get("kill_mult", 3.5)
	_current_target_id = payload.get("target_id", "")
	var target_x: float = payload.get("x", 0.0)
	var target_y: float = payload.get("y", 0.0)
	var vp_size = get_viewport().size

	# 紫色三次強閃光
	_flash_screen(COLOR_PROPHECY, 0.5, 3)

	# 頂部橫幅
	_show_banner("🔮 預言降臨！", COLOR_PROPHECY, 3.5)

	# 中央大字
	_show_big_text("🔮 預言降臨！", COLOR_PROPHECY, 52, 2.5)

	# 倍率說明
	_show_sub_text("擊破指定目標獲得 ×%.1f 倍率！" % kill_mult, COLOR_GOLD, 2.0)

	# 右側豎向計時條
	_create_timer_bar(_duration_sec)

	# 目標標記（在目標位置顯示閃爍標記）
	_create_target_marker(target_x, target_y)

## prophecy_broadcast — 全服廣播有人觸發預言
func _on_prophecy_broadcast(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "某玩家")
	var duration_sec: int = payload.get("duration_sec", 12)
	_show_top_banner("🔮 %s 觸發預言！%d 秒內追蹤指定目標！" % [player_name, duration_sec], COLOR_PROPHECY, 3.0)

## prophecy_fulfilled — 預言成真（個人訊息）
func _on_prophecy_fulfilled(payload: Dictionary) -> void:
	var kill_mult: float = payload.get("kill_mult", 3.5)
	var reward: int = payload.get("reward", 0)

	# 清除計時條和目標標記
	_clear_timer_bar()
	_clear_target_marker()

	# 金色三次強閃光
	_flash_screen(COLOR_GOLD, 0.6, 3)

	# 中央大字
	_show_big_text("🔮 預言成真！×%.1f" % kill_mult, COLOR_GOLD, 3.0)

	# 結算彈窗
	_show_reward_popup(reward, kill_mult)

## prophecy_broadcast_fulfilled — 全服廣播預言成真
func _on_prophecy_broadcast_fulfilled(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "某玩家")
	var kill_mult: float = payload.get("kill_mult", 3.5)
	var reward: int = payload.get("reward", 0)
	_show_top_banner("🔮 %s 預言成真！×%.1f 倍率！獲得 %d 金幣！" % [player_name, kill_mult, reward], COLOR_GOLD, 3.5)

## prophecy_transfer — 預言轉移（個人訊息）
func _on_prophecy_transfer(payload: Dictionary) -> void:
	var new_target_id: String = payload.get("target_id", "")
	var new_x: float = payload.get("x", 0.0)
	var new_y: float = payload.get("y", 0.0)
	var transfer_count: int = payload.get("transfer_count", 1)
	var kill_mult: float = payload.get("kill_mult", 3.5)

	_current_target_id = new_target_id

	# 橙色閃光
	_flash_screen(COLOR_ORANGE, 0.3, 1)

	# 提示文字
	_show_big_text("🔮 預言轉移！（第 %d 次）" % transfer_count, COLOR_ORANGE, 2.0)

	# 更新目標標記位置
	_clear_target_marker()
	_create_target_marker(new_x, new_y)

## prophecy_broadcast_transfer — 全服廣播預言轉移
func _on_prophecy_broadcast_transfer(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "某玩家")
	var transfer_count: int = payload.get("transfer_count", 1)
	_show_top_banner("🔮 %s 的預言轉移！（第 %d 次）" % [player_name, transfer_count], COLOR_ORANGE, 2.5)

## prophecy_fail — 預言失敗（全服廣播）
func _on_prophecy_fail(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "某玩家")
	var affected_count: int = payload.get("affected_count", 0)

	# 清除計時條和目標標記
	_clear_timer_bar()
	_clear_target_marker()

	# 灰色閃光
	_flash_screen(COLOR_FAIL, 0.4, 2)

	# 提示文字
	_show_big_text("🔮 預言失敗！HP -20%%", COLOR_FAIL, 2.5)
	_show_sub_text("%s 的預言未能成真，%d 個目標受到傷害！" % [player_name, affected_count], COLOR_FAIL, 2.0)

# ─── 內部 UI 工具函數 ───────────────────────────────────────────────────────

## 全螢幕閃光
func _flash_screen(color: Color, alpha: float, times: int) -> void:
	var vp_size = get_viewport().size
	var flash = ColorRect.new()
	flash.color = Color(color.r, color.g, color.b, alpha)
	flash.size = vp_size
	flash.position = Vector2.ZERO
	add_child(flash)

	var tween = create_tween()
	for i in range(times):
		tween.tween_property(flash, "modulate:a", 1.0, 0.05)
		tween.tween_property(flash, "modulate:a", 0.0, 0.1)
	tween.tween_callback(flash.queue_free)

## 頂部橫幅（持續顯示）
func _show_banner(text: String, color: Color, duration: float) -> void:
	var vp_size = get_viewport().size
	var banner = ColorRect.new()
	banner.color = Color(color.r, color.g, color.b, 0.85)
	banner.size = Vector2(vp_size.x, 48)
	banner.position = Vector2(0, 0)
	add_child(banner)

	var label = Label.new()
	label.text = text
	label.add_theme_font_size_override("font_size", 20)
	label.add_theme_color_override("font_color", Color.WHITE)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	label.size = banner.size
	label.position = Vector2.ZERO
	banner.add_child(label)

	var tween = create_tween()
	tween.tween_interval(duration)
	tween.tween_property(banner, "modulate:a", 0.0, 0.5)
	tween.tween_callback(banner.queue_free)

## 頂部小橫幅（全服廣播用）
func _show_top_banner(text: String, color: Color, duration: float) -> void:
	var vp_size = get_viewport().size
	var banner = ColorRect.new()
	banner.color = Color(color.r, color.g, color.b, 0.75)
	banner.size = Vector2(vp_size.x * 0.7, 36)
	banner.position = Vector2(vp_size.x * 0.15, 56)
	add_child(banner)

	var label = Label.new()
	label.text = text
	label.add_theme_font_size_override("font_size", 16)
	label.add_theme_color_override("font_color", Color.WHITE)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	label.size = banner.size
	label.position = Vector2.ZERO
	banner.add_child(label)

	var tween = create_tween()
	tween.tween_interval(duration)
	tween.tween_property(banner, "modulate:a", 0.0, 0.4)
	tween.tween_callback(banner.queue_free)

## 中央大字
func _show_big_text(text: String, color: Color, duration: float) -> void:
	var vp_size = get_viewport().size
	var label = Label.new()
	label.text = text
	label.add_theme_font_size_override("font_size", 48)
	label.add_theme_color_override("font_color", color)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.size = Vector2(vp_size.x, 80)
	label.position = Vector2(0, vp_size.y * 0.35)
	add_child(label)

	var tween = create_tween()
	tween.tween_property(label, "position:y", vp_size.y * 0.30, 0.3)
	tween.tween_interval(duration - 0.8)
	tween.tween_property(label, "modulate:a", 0.0, 0.5)
	tween.tween_callback(label.queue_free)

## 副標題文字
func _show_sub_text(text: String, color: Color, duration: float) -> void:
	var vp_size = get_viewport().size
	var label = Label.new()
	label.text = text
	label.add_theme_font_size_override("font_size", 22)
	label.add_theme_color_override("font_color", color)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.size = Vector2(vp_size.x, 40)
	label.position = Vector2(0, vp_size.y * 0.45)
	add_child(label)

	var tween = create_tween()
	tween.tween_interval(duration)
	tween.tween_property(label, "modulate:a", 0.0, 0.4)
	tween.tween_callback(label.queue_free)

## 建立右側豎向計時條
func _create_timer_bar(duration: int) -> void:
	_clear_timer_bar()
	var vp_size = get_viewport().size

	# 背景條
	_timer_bar_bg = ColorRect.new()
	_timer_bar_bg.color = Color(0.2, 0.1, 0.3, 0.6)
	_timer_bar_bg.size = Vector2(16, 200)
	_timer_bar_bg.position = Vector2(vp_size.x - 28, vp_size.y * 0.3)
	add_child(_timer_bar_bg)

	# 計時條（紫→金漸變，從上往下縮短）
	_timer_bar = ColorRect.new()
	_timer_bar.color = COLOR_PROPHECY
	_timer_bar.size = Vector2(16, 200)
	_timer_bar.position = Vector2(vp_size.x - 28, vp_size.y * 0.3)
	add_child(_timer_bar)

	# 計時條動畫（從 200 縮到 0）
	_timer_tween = create_tween()
	_timer_tween.tween_property(_timer_bar, "size:y", 0.0, float(duration))

## 清除計時條
func _clear_timer_bar() -> void:
	if is_instance_valid(_timer_tween):
		_timer_tween.kill()
		_timer_tween = null
	if is_instance_valid(_timer_bar):
		_timer_bar.queue_free()
		_timer_bar = null
	if is_instance_valid(_timer_bar_bg):
		_timer_bar_bg.queue_free()
		_timer_bar_bg = null

## 建立目標標記（在目標位置顯示閃爍的 🔮 標記）
func _create_target_marker(target_x: float, target_y: float) -> void:
	_clear_target_marker()
	if target_x <= 0 and target_y <= 0:
		return

	_target_marker = Control.new()
	_target_marker.position = Vector2(target_x - 20, target_y - 50)
	_target_marker.size = Vector2(40, 40)
	add_child(_target_marker)

	# 閃爍的標記文字
	var label = Label.new()
	label.text = "🔮"
	label.add_theme_font_size_override("font_size", 28)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.size = Vector2(40, 40)
	label.position = Vector2.ZERO
	_target_marker.add_child(label)

	# 閃爍動畫
	var tween = _target_marker.create_tween().set_loops()
	tween.tween_property(label, "modulate:a", 0.3, 0.4)
	tween.tween_property(label, "modulate:a", 1.0, 0.4)

## 清除目標標記
func _clear_target_marker() -> void:
	if is_instance_valid(_target_marker):
		_target_marker.queue_free()
		_target_marker = null
	_current_target_id = ""

## 結算彈窗
func _show_reward_popup(reward: int, kill_mult: float) -> void:
	var vp_size = get_viewport().size
	var popup = ColorRect.new()
	popup.color = Color(0.1, 0.05, 0.2, 0.92)
	popup.size = Vector2(280, 120)
	popup.position = Vector2(vp_size.x, vp_size.y * 0.4)  # 從右側滑入
	add_child(popup)

	var title_label = Label.new()
	title_label.text = "🔮 預言成真！"
	title_label.add_theme_font_size_override("font_size", 22)
	title_label.add_theme_color_override("font_color", COLOR_GOLD)
	title_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title_label.size = Vector2(280, 40)
	title_label.position = Vector2(0, 8)
	popup.add_child(title_label)

	var mult_label = Label.new()
	mult_label.text = "×%.1f 倍率加成" % kill_mult
	mult_label.add_theme_font_size_override("font_size", 18)
	mult_label.add_theme_color_override("font_color", COLOR_PROPHECY)
	mult_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	mult_label.size = Vector2(280, 32)
	mult_label.position = Vector2(0, 44)
	popup.add_child(mult_label)

	var reward_label = Label.new()
	reward_label.text = "+%d 金幣" % reward
	reward_label.add_theme_font_size_override("font_size", 20)
	reward_label.add_theme_color_override("font_color", COLOR_GOLD)
	reward_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	reward_label.size = Vector2(280, 32)
	reward_label.position = Vector2(0, 76)
	popup.add_child(reward_label)

	# 從右側滑入
	var tween = create_tween()
	tween.tween_property(popup, "position:x", vp_size.x - 300, 0.3)
	tween.tween_interval(2.5)
	tween.tween_property(popup, "modulate:a", 0.0, 0.5)
	tween.tween_callback(popup.queue_free)
