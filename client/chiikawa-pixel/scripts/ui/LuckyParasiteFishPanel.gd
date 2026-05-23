## LuckyParasiteFishPanel.gd — 幸運寄生魚系統面板（DAY-229）
## 業界原創「寄生附著+跳躍」機制
##
## 視覺設計：
##   - 綠色寄生主題（#27AE60 + #2ECC71 + #A9DFBF + #FF4444）
##   - parasite_start：綠色雙閃光 + 頂部橫幅 + 「🦠 寄生釋放！」大字 + 寄生目標菱形標記（綠色閃爍）
##   - parasite_tick：目標閃爍 + HP 損失浮動文字（綠色）
##   - parasite_jump：跳躍軌跡動畫 + 「🦠 跳躍！」提示
##   - parasite_kill：消失閃光 + ×2.2 浮動文字（金色）
##   - parasite_end：標記淡出
extends CanvasLayer

# 寄生狀態
var _active: bool = false
var _parasite_markers: Dictionary = {}  # targetID → marker node

# 主題顏色
const COLOR_PRIMARY   = Color("#27AE60")  # 深綠
const COLOR_LIGHT     = Color("#2ECC71")  # 亮綠
const COLOR_PALE      = Color("#A9DFBF")  # 淡綠
const COLOR_DANGER    = Color("#FF4444")  # 紅色（HP 損失）
const COLOR_GOLD      = Color("#FFD700")  # 金黃

func _ready() -> void:
	layer = 16  # 幸運寄生魚面板層級

## 處理幸運寄生魚訊息
func handle_lucky_parasite_fish(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"parasite_start":
			_on_parasite_start(payload)
		"parasite_tick":
			_on_parasite_tick(payload)
		"parasite_jump":
			_on_parasite_jump(payload)
		"parasite_kill":
			_on_parasite_kill(payload)
		"parasite_end":
			_on_parasite_end(payload)

## parasite_start — 寄生釋放開始
func _on_parasite_start(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "")
	var parasite_count: int = payload.get("parasite_count", 0)
	var kill_mult: float = payload.get("kill_mult", 2.2)

	_active = true

	# 綠色雙閃光
	_flash_screen(COLOR_PRIMARY, 0.16)
	await get_tree().create_timer(0.1).timeout
	_flash_screen(COLOR_LIGHT, 0.13)

	var vp_size = get_viewport().size

	# 頂部橫幅
	var banner = Label.new()
	banner.text = "🦠 %s 觸發寄生釋放！%d 個目標被寄生！" % [player_name, parasite_count]
	banner.add_theme_font_size_override("font_size", 14)
	banner.add_theme_color_override("font_color", COLOR_LIGHT)
	banner.position = Vector2(vp_size.x / 2 - 160, 6)
	add_child(banner)

	var tween_banner = banner.create_tween()
	tween_banner.tween_interval(3.5)
	tween_banner.tween_property(banner, "modulate:a", 0.0, 0.5)
	tween_banner.tween_callback(banner.queue_free)

	# 「🦠 寄生釋放！」大字
	var big_label = Label.new()
	big_label.text = "🦠 寄生釋放！"
	big_label.add_theme_font_size_override("font_size", 48)
	big_label.add_theme_color_override("font_color", COLOR_PRIMARY)
	big_label.position = vp_size / 2 - Vector2(110, 28)
	add_child(big_label)

	var tween_big = big_label.create_tween()
	tween_big.tween_property(big_label, "scale", Vector2(1.2, 1.2), 0.1)
	tween_big.tween_property(big_label, "scale", Vector2(1.0, 1.0), 0.08)
	tween_big.tween_interval(0.4)
	tween_big.tween_property(big_label, "modulate:a", 0.0, 0.4)
	tween_big.tween_callback(big_label.queue_free)

	# 倍率提示
	var mult_label = Label.new()
	mult_label.text = "擊破寄生目標 ×%.1f！" % kill_mult
	mult_label.add_theme_font_size_override("font_size", 14)
	mult_label.add_theme_color_override("font_color", COLOR_GOLD)
	mult_label.position = Vector2(vp_size.x / 2 - 70, vp_size.y / 2 + 28)
	add_child(mult_label)

	var tween_mult = mult_label.create_tween()
	tween_mult.tween_interval(2.0)
	tween_mult.tween_property(mult_label, "modulate:a", 0.0, 0.5)
	tween_mult.tween_callback(mult_label.queue_free)

## parasite_tick — 寄生 HP 損失
func _on_parasite_tick(payload: Dictionary) -> void:
	var target_id: String = payload.get("target_id", "")
	var hp_loss: float = payload.get("hp_loss", 0.08)
	var tick_count: int = payload.get("tick_count", 1)

	var vp_size = get_viewport().size

	# HP 損失浮動文字
	var loss_label = Label.new()
	loss_label.text = "-%.0f%% HP" % (hp_loss * 100)
	loss_label.add_theme_font_size_override("font_size", 13)
	loss_label.add_theme_color_override("font_color", COLOR_DANGER)
	# 隨機位置（模擬在目標附近）
	loss_label.position = Vector2(
		vp_size.x * 0.3 + randf() * vp_size.x * 0.4,
		vp_size.y * 0.3 + randf() * vp_size.y * 0.3
	)
	add_child(loss_label)

	var tween_loss = loss_label.create_tween()
	tween_loss.tween_property(loss_label, "position:y", loss_label.position.y - 20, 0.5)
	tween_loss.parallel().tween_property(loss_label, "modulate:a", 0.0, 0.5)
	tween_loss.tween_callback(loss_label.queue_free)

## parasite_jump — 寄生蟲跳躍
func _on_parasite_jump(payload: Dictionary) -> void:
	var jump_layer: int = payload.get("jump_layer", 1)
	var kill_mult: float = payload.get("kill_mult", 2.2)

	var vp_size = get_viewport().size

	# 「🦠 跳躍！」提示
	var jump_label = Label.new()
	jump_label.text = "🦠 寄生蟲跳躍！（第%d跳）" % jump_layer
	jump_label.add_theme_font_size_override("font_size", 16)
	jump_label.add_theme_color_override("font_color", COLOR_LIGHT)
	jump_label.position = Vector2(vp_size.x / 2 - 80, vp_size.y / 2 - 10)
	add_child(jump_label)

	var tween_jump = jump_label.create_tween()
	tween_jump.tween_property(jump_label, "scale", Vector2(1.15, 1.15), 0.08)
	tween_jump.tween_property(jump_label, "scale", Vector2(1.0, 1.0), 0.06)
	tween_jump.tween_interval(0.5)
	tween_jump.tween_property(jump_label, "modulate:a", 0.0, 0.4)
	tween_jump.tween_callback(jump_label.queue_free)

	# 綠色閃光
	_flash_screen(COLOR_LIGHT, 0.12)

## parasite_kill — 寄生目標被玩家擊破
func _on_parasite_kill(payload: Dictionary) -> void:
	var target_id: String = payload.get("target_id", "")
	var kill_mult: float = payload.get("kill_mult", 2.2)
	var kill_reward: int = payload.get("kill_reward", 0)
	var jump_layer: int = payload.get("jump_layer", 0)

	# 移除寄生標記
	if _parasite_markers.has(target_id):
		var marker = _parasite_markers[target_id]
		if is_instance_valid(marker):
			var tween_fade = marker.create_tween()
			tween_fade.tween_property(marker, "modulate:a", 0.0, 0.2)
			tween_fade.tween_callback(marker.queue_free)
		_parasite_markers.erase(target_id)

	var vp_size = get_viewport().size

	# ×2.2 浮動文字（金色）
	var reward_label = Label.new()
	reward_label.text = "🦠 ×%.1f 擊破！" % kill_mult
	reward_label.add_theme_font_size_override("font_size", 22)
	reward_label.add_theme_color_override("font_color", COLOR_GOLD)
	reward_label.position = Vector2(vp_size.x / 2 - 70, vp_size.y / 2 - 15)
	add_child(reward_label)

	var tween_r = reward_label.create_tween()
	tween_r.tween_property(reward_label, "position:y", reward_label.position.y - 35, 0.6)
	tween_r.parallel().tween_property(reward_label, "modulate:a", 0.0, 0.6)
	tween_r.tween_callback(reward_label.queue_free)

	# 綠色消失閃光
	_flash_screen(COLOR_PRIMARY, 0.12)

## parasite_end — 寄生消散
func _on_parasite_end(payload: Dictionary) -> void:
	var target_id: String = payload.get("target_id", "")

	# 移除寄生標記
	if _parasite_markers.has(target_id):
		var marker = _parasite_markers[target_id]
		if is_instance_valid(marker):
			var tween_fade = marker.create_tween()
			tween_fade.tween_property(marker, "modulate:a", 0.0, 0.3)
			tween_fade.tween_callback(marker.queue_free)
		_parasite_markers.erase(target_id)

# ---- 輔助函數 ----

## 全螢幕閃光效果
func _flash_screen(color: Color, duration: float) -> void:
	var flash = ColorRect.new()
	flash.color = Color(color.r, color.g, color.b, 0.28)
	flash.set_anchors_preset(Control.PRESET_FULL_RECT)
	add_child(flash)

	var tween = flash.create_tween()
	tween.tween_property(flash, "modulate:a", 0.0, duration)
	tween.tween_callback(flash.queue_free)
