## LuckyEchoFishPanel.gd — 幸運回聲魚系統面板（DAY-233）
## 業界原創「回聲分身+層疊倍率」機制
##
## 視覺設計：
##   - 紫色回聲主題（#9B59B6 + #8E44AD + #D7BDE2 + #F5EEF8）
##   - echo_ready：紫色雙閃光 + 「🔊 回聲模式！」大字 + 「下次擊破產生分身」提示
##   - echo_broadcast：頂部小橫幅（全服廣播）
##   - echo_spawn：紫色閃光 + 分身標記 + 層數倍率浮動文字
##   - echo_spawn_personal：個人提示（「🔊 第N層回聲！×X.X」）
##   - echo_expire：分身標記淡出
extends CanvasLayer

# 主題顏色
const COLOR_PRIMARY  = Color("#9B59B6")  # 紫色
const COLOR_DARK     = Color("#8E44AD")  # 深紫
const COLOR_PALE     = Color("#D7BDE2")  # 淡紫
const COLOR_LIGHT_BG = Color("#F5EEF8")  # 極淡紫
const COLOR_GOLD     = Color("#FFD700")  # 金黃
const COLOR_LAYER1   = Color("#CE93D8")  # 第1層（淡紫）
const COLOR_LAYER2   = Color("#AB47BC")  # 第2層（中紫）
const COLOR_LAYER3   = Color("#7B1FA2")  # 第3層（深紫）

# 回聲分身標記（instanceID → 節點）
var _echo_markers: Dictionary = {}

func _ready() -> void:
	layer = 12  # 幸運回聲魚面板層級

## 處理幸運回聲魚訊息
func handle_lucky_echo_fish(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"echo_ready":
			_on_echo_ready(payload)
		"echo_broadcast":
			_on_echo_broadcast(payload)
		"echo_spawn":
			_on_echo_spawn(payload)
		"echo_spawn_personal":
			_on_echo_spawn_personal(payload)
		"echo_expire":
			_on_echo_expire(payload)

## echo_ready — 回聲模式啟動（個人訊息）
func _on_echo_ready(_payload: Dictionary) -> void:
	var vp_size = get_viewport().size

	# 紫色雙閃光
	_flash_screen(COLOR_PRIMARY, 0.14)
	await get_tree().create_timer(0.09).timeout
	_flash_screen(COLOR_DARK, 0.11)

	# 「🔊 回聲模式！」大字
	var big_label = Label.new()
	big_label.text = "🔊 回聲模式！"
	big_label.add_theme_font_size_override("font_size", 46)
	big_label.add_theme_color_override("font_color", COLOR_PRIMARY)
	big_label.position = vp_size / 2 - Vector2(100, 28)
	add_child(big_label)

	var tween_big = big_label.create_tween()
	tween_big.tween_property(big_label, "scale", Vector2(1.2, 1.2), 0.10)
	tween_big.tween_property(big_label, "scale", Vector2(1.0, 1.0), 0.08)
	tween_big.tween_interval(0.5)
	tween_big.tween_property(big_label, "modulate:a", 0.0, 0.4)
	tween_big.tween_callback(big_label.queue_free)

	# 提示文字
	var hint_label = Label.new()
	hint_label.text = "下次擊破將產生回聲分身！"
	hint_label.add_theme_font_size_override("font_size", 13)
	hint_label.add_theme_color_override("font_color", COLOR_PALE)
	hint_label.position = Vector2(vp_size.x / 2 - 80, vp_size.y / 2 + 28)
	add_child(hint_label)

	var tween_hint = hint_label.create_tween()
	tween_hint.tween_interval(1.5)
	tween_hint.tween_property(hint_label, "modulate:a", 0.0, 0.5)
	tween_hint.tween_callback(hint_label.queue_free)

## echo_broadcast — 全服廣播（頂部小橫幅）
func _on_echo_broadcast(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "")
	var vp_size = get_viewport().size

	var banner = Label.new()
	banner.text = "🔊 %s 觸發回聲模式！" % player_name
	banner.add_theme_font_size_override("font_size", 13)
	banner.add_theme_color_override("font_color", COLOR_PALE)
	banner.position = Vector2(vp_size.x / 2 - 100, 6)
	add_child(banner)

	var tween_banner = banner.create_tween()
	tween_banner.tween_interval(2.5)
	tween_banner.tween_property(banner, "modulate:a", 0.0, 0.4)
	tween_banner.tween_callback(banner.queue_free)

## echo_spawn — 回聲分身生成（全服廣播）
func _on_echo_spawn(payload: Dictionary) -> void:
	var layer_num: int = payload.get("layer", 1)
	var echo_instance_id: String = payload.get("echo_instance_id", "")
	var echo_x: float = payload.get("echo_x", 400.0)
	var echo_y: float = payload.get("echo_y", 300.0)
	var mult_label: String = payload.get("mult_label", "×1.5")
	var player_name: String = payload.get("player_name", "")

	var vp_size = get_viewport().size

	# 依層數選擇顏色
	var layer_color: Color
	match layer_num:
		1: layer_color = COLOR_LAYER1
		2: layer_color = COLOR_LAYER2
		3: layer_color = COLOR_LAYER3
		_: layer_color = COLOR_PRIMARY

	# 紫色閃光（輕微）
	_flash_screen(layer_color, 0.07)

	# 在分身位置顯示標記
	var scale_x = vp_size.x / 1000.0
	var scale_y = vp_size.y / 600.0
	var screen_x = echo_x * scale_x
	var screen_y = echo_y * scale_y

	# 分身標記（菱形輪廓 + 層數）
	var marker = Label.new()
	marker.text = "◆ L%d" % layer_num
	marker.add_theme_font_size_override("font_size", 12)
	marker.add_theme_color_override("font_color", layer_color)
	marker.position = Vector2(screen_x - 16, screen_y - 20)
	add_child(marker)

	# 閃爍動畫
	var tween_marker = marker.create_tween().set_loops(5)
	tween_marker.tween_property(marker, "modulate:a", 0.3, 0.3)
	tween_marker.tween_property(marker, "modulate:a", 1.0, 0.3)

	# 記錄標記（供 echo_expire 使用）
	if echo_instance_id != "":
		_echo_markers[echo_instance_id] = marker

	# 3 秒後自動清除（若未被 echo_expire 清除）
	get_tree().create_timer(3.0).timeout.connect(func():
		if is_instance_valid(marker):
			marker.queue_free()
		if _echo_markers.has(echo_instance_id):
			_echo_markers.erase(echo_instance_id)
	)

	# 倍率浮動文字
	var mult_text = Label.new()
	mult_text.text = "🔊 %s" % mult_label
	mult_text.add_theme_font_size_override("font_size", 18)
	mult_text.add_theme_color_override("font_color", COLOR_GOLD)
	mult_text.position = Vector2(screen_x - 20, screen_y - 40)
	add_child(mult_text)

	var tween_mult = mult_text.create_tween()
	tween_mult.tween_property(mult_text, "position:y", mult_text.position.y - 25, 0.5)
	tween_mult.parallel().tween_property(mult_text, "scale", Vector2(1.3, 1.3), 0.2)
	tween_mult.tween_property(mult_text, "modulate:a", 0.0, 0.4)
	tween_mult.tween_callback(mult_text.queue_free)

	# 廣播提示（非本人）
	if player_name != "":
		var spawn_label = Label.new()
		spawn_label.text = "🔊 %s 的第%d層回聲！" % [player_name, layer_num]
		spawn_label.add_theme_font_size_override("font_size", 12)
		spawn_label.add_theme_color_override("font_color", COLOR_PALE)
		spawn_label.position = Vector2(vp_size.x / 2 - 80, vp_size.y / 2 - 50)
		add_child(spawn_label)

		var tween_spawn = spawn_label.create_tween()
		tween_spawn.tween_property(spawn_label, "position:y", spawn_label.position.y - 15, 0.4)
		tween_spawn.parallel().tween_property(spawn_label, "modulate:a", 0.0, 0.4)
		tween_spawn.tween_callback(spawn_label.queue_free)

## echo_spawn_personal — 個人回聲分身提示
func _on_echo_spawn_personal(payload: Dictionary) -> void:
	var layer_num: int = payload.get("layer", 1)
	var mult_label: String = payload.get("mult_label", "×1.5")
	var vp_size = get_viewport().size

	# 依層數選擇顏色
	var layer_color: Color
	match layer_num:
		1: layer_color = COLOR_LAYER1
		2: layer_color = COLOR_LAYER2
		3: layer_color = COLOR_LAYER3
		_: layer_color = COLOR_PRIMARY

	# 個人提示（右側浮現）
	var personal_label = Label.new()
	personal_label.text = "🔊 第%d層回聲！%s" % [layer_num, mult_label]
	personal_label.add_theme_font_size_override("font_size", 20)
	personal_label.add_theme_color_override("font_color", layer_color)
	personal_label.position = Vector2(vp_size.x - 200, vp_size.y / 2 - 20 + (layer_num - 1) * 30)
	add_child(personal_label)

	var tween_p = personal_label.create_tween()
	tween_p.tween_property(personal_label, "scale", Vector2(1.2, 1.2), 0.12)
	tween_p.tween_property(personal_label, "scale", Vector2(1.0, 1.0), 0.10)
	tween_p.tween_interval(0.8)
	tween_p.tween_property(personal_label, "modulate:a", 0.0, 0.4)
	tween_p.tween_callback(personal_label.queue_free)

## echo_expire — 回聲分身消失
func _on_echo_expire(payload: Dictionary) -> void:
	var echo_instance_id: String = payload.get("echo_instance_id", "")
	if echo_instance_id == "":
		return

	if _echo_markers.has(echo_instance_id):
		var marker = _echo_markers[echo_instance_id]
		if is_instance_valid(marker):
			var tween_expire = marker.create_tween()
			tween_expire.tween_property(marker, "modulate:a", 0.0, 0.3)
			tween_expire.tween_callback(marker.queue_free)
		_echo_markers.erase(echo_instance_id)

# ---- 輔助函數 ----

## 全螢幕閃光效果
func _flash_screen(color: Color, duration: float) -> void:
	var flash = ColorRect.new()
	flash.color = Color(color.r, color.g, color.b, 0.25)
	flash.set_anchors_preset(Control.PRESET_FULL_RECT)
	add_child(flash)

	var tween = flash.create_tween()
	tween.tween_property(flash, "modulate:a", 0.0, duration)
	tween.tween_callback(flash.queue_free)
