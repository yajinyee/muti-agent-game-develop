## LuckyChainReactionPanel.gd — 幸運連鎖反應魚系統面板（DAY-241）
## 業界原創「多米諾骨牌效應」機制
##
## 視覺設計：
##   - 橙紅連鎖主題（#FF6B35 + #FF4500 + #FFD700 + #FFF0E6）
##   - chain_start：橙色雙閃光 + 頂部橫幅 + 「🔗 連鎖起點標記！」大字 + 目標標記閃爍
##   - chain_explode：每層引爆閃光 + 連鎖線動畫 + 層數計數 + 倍率浮動文字
##   - chain_broken：灰色提示「連鎖中斷」
##   - chain_complete：全螢幕三次強閃光 + 「🔗 連鎖完成！8層！」大字 + 結算彈窗
extends CanvasLayer

# 主題顏色
const COLOR_CHAIN    = Color("#FF6B35")  # 橙紅（主題）
const COLOR_FIRE     = Color("#FF4500")  # 深橙（爆炸）
const COLOR_GOLD     = Color("#FFD700")  # 金色（倍率）
const COLOR_PALE     = Color("#FFF0E6")  # 淡橙（文字）
const COLOR_BROKEN   = Color("#7F8C8D")  # 灰色（中斷）
const COLOR_WHITE    = Color("#FFFFFF")  # 白色

# 連鎖狀態
var _chain_layer: int = 0
var _chain_player: String = ""
var _total_reward: int = 0

# 層數計數器節點
var _layer_counter: Label = null
var _layer_counter_tween: Tween = null

func _ready() -> void:
	layer = 4  # 幸運連鎖反應魚面板層級

## 處理幸運連鎖反應魚訊息
func handle_lucky_chain_reaction(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"chain_start":
			_on_chain_start(payload)
		"chain_explode":
			_on_chain_explode(payload)
		"chain_broken":
			_on_chain_broken(payload)
		"chain_complete":
			_on_chain_complete(payload)

## chain_start — 連鎖起點標記
func _on_chain_start(payload: Dictionary) -> void:
	_chain_player = payload.get("player_name", "")
	_chain_layer = 0
	_total_reward = 0
	var max_layers: int = payload.get("max_layers", 8)
	var base_mult: float = payload.get("base_mult", 1.4)
	var vp_size = get_viewport().size

	# 橙色雙閃光
	_flash_screen(COLOR_CHAIN, 0.12)
	await get_tree().create_timer(0.08).timeout
	_flash_screen(COLOR_PALE, 0.09)

	# 頂部橫幅
	var banner = Label.new()
	banner.text = "🔗 %s 觸發連鎖反應！最多 %d 層！" % [_chain_player, max_layers]
	banner.add_theme_font_size_override("font_size", 14)
	banner.add_theme_color_override("font_color", COLOR_PALE)
	banner.position = Vector2(vp_size.x / 2 - 140, 6)
	add_child(banner)

	var tw_banner = banner.create_tween()
	tw_banner.tween_interval(4.0)
	tw_banner.tween_property(banner, "modulate:a", 0.0, 0.4)
	tw_banner.tween_callback(banner.queue_free)

	# 大字提示
	var big_label = Label.new()
	big_label.text = "🔗 連鎖起點標記！"
	big_label.add_theme_font_size_override("font_size", 36)
	big_label.add_theme_color_override("font_color", COLOR_CHAIN)
	big_label.position = Vector2(vp_size.x / 2 - 120, vp_size.y / 2 - 80)
	add_child(big_label)

	big_label.scale = Vector2(0.7, 0.7)
	var tw_big = big_label.create_tween()
	tw_big.tween_property(big_label, "scale", Vector2(1.1, 1.1), 0.15)
	tw_big.tween_property(big_label, "scale", Vector2(1.0, 1.0), 0.08)
	tw_big.tween_interval(1.5)
	tw_big.tween_property(big_label, "modulate:a", 0.0, 0.4)
	tw_big.tween_callback(big_label.queue_free)

	# 倍率說明
	var info_label = Label.new()
	info_label.text = "擊破標記目標引發連鎖！×%.1f → ×%.1f（每層遞減）" % [base_mult, 0.7]
	info_label.add_theme_font_size_override("font_size", 13)
	info_label.add_theme_color_override("font_color", Color(0.9, 0.8, 0.7))
	info_label.position = Vector2(vp_size.x / 2 - 170, vp_size.y / 2 - 38)
	add_child(info_label)

	var tw_info = info_label.create_tween()
	tw_info.tween_interval(3.0)
	tw_info.tween_property(info_label, "modulate:a", 0.0, 0.4)
	tw_info.tween_callback(info_label.queue_free)

	# 建立層數計數器（右上角）
	_create_layer_counter(max_layers)

## 建立層數計數器
func _create_layer_counter(max_layers: int) -> void:
	if is_instance_valid(_layer_counter):
		_layer_counter.queue_free()

	var vp_size = get_viewport().size

	var counter_bg = ColorRect.new()
	counter_bg.color = Color(0.1, 0.05, 0.0, 0.75)
	counter_bg.size = Vector2(110, 44)
	counter_bg.position = Vector2(vp_size.x - 120, 50)
	add_child(counter_bg)

	_layer_counter = Label.new()
	_layer_counter.text = "🔗 0 / %d 層" % max_layers
	_layer_counter.add_theme_font_size_override("font_size", 16)
	_layer_counter.add_theme_color_override("font_color", COLOR_CHAIN)
	_layer_counter.position = Vector2(vp_size.x - 116, 58)
	add_child(_layer_counter)

	# 計數器自動消失（15 秒後）
	var tw = counter_bg.create_tween()
	tw.tween_interval(15.0)
	tw.tween_property(counter_bg, "modulate:a", 0.0, 0.5)
	tw.tween_callback(counter_bg.queue_free)

	var tw2 = _layer_counter.create_tween()
	tw2.tween_interval(15.0)
	tw2.tween_property(_layer_counter, "modulate:a", 0.0, 0.5)
	tw2.tween_callback(func():
		if is_instance_valid(_layer_counter):
			_layer_counter.queue_free()
			_layer_counter = null
	)

## chain_explode — 本層連鎖引爆
func _on_chain_explode(payload: Dictionary) -> void:
	var layer_num: int = payload.get("layer", 1)
	var max_layers: int = payload.get("max_layers", 8)
	var mult: float = payload.get("mult", 1.0)
	var reward: int = payload.get("reward", 0)
	var from_x: float = payload.get("from_x", 0.0)
	var from_y: float = payload.get("from_y", 0.0)
	var to_x: float = payload.get("to_x", 0.0)
	var to_y: float = payload.get("to_y", 0.0)

	_chain_layer = layer_num
	_total_reward += reward

	var vp_size = get_viewport().size

	# 更新層數計數器
	if is_instance_valid(_layer_counter):
		_layer_counter.text = "🔗 %d / %d 層" % [layer_num, max_layers]
		# 閃爍動畫
		if is_instance_valid(_layer_counter_tween):
			_layer_counter_tween.kill()
		_layer_counter_tween = _layer_counter.create_tween()
		_layer_counter_tween.tween_property(_layer_counter, "modulate", Color(1.5, 1.5, 0.5), 0.08)
		_layer_counter_tween.tween_property(_layer_counter, "modulate", Color.WHITE, 0.12)

	# 連鎖閃光（強度隨層數遞增）
	var flash_alpha = 0.08 + float(layer_num) * 0.02
	_flash_screen(COLOR_FIRE, flash_alpha)

	# 連鎖線動畫（從 from 到 to）
	_draw_chain_line(from_x, from_y, to_x, to_y)

	# 倍率浮動文字（在目標位置）
	var mult_label = Label.new()
	mult_label.text = "🔗 ×%.1f" % mult
	mult_label.add_theme_font_size_override("font_size", 20)
	var mult_color = COLOR_GOLD if mult >= 1.0 else Color(0.8, 0.5, 0.3)
	mult_label.add_theme_color_override("font_color", mult_color)
	mult_label.position = Vector2(to_x - 30, to_y - 40)
	add_child(mult_label)

	var tw_mult = mult_label.create_tween()
	tw_mult.tween_property(mult_label, "position:y", to_y - 80, 0.6)
	tw_mult.parallel().tween_property(mult_label, "modulate:a", 0.0, 0.6)
	tw_mult.tween_callback(mult_label.queue_free)

	# 層數標記（左側）
	var layer_label = Label.new()
	layer_label.text = "第 %d 層" % layer_num
	layer_label.add_theme_font_size_override("font_size", 13)
	layer_label.add_theme_color_override("font_color", COLOR_PALE)
	layer_label.position = Vector2(8, vp_size.y / 2 - 20 + layer_num * 18)
	add_child(layer_label)

	var tw_layer = layer_label.create_tween()
	tw_layer.tween_interval(1.5)
	tw_layer.tween_property(layer_label, "modulate:a", 0.0, 0.4)
	tw_layer.tween_callback(layer_label.queue_free)

## 繪製連鎖線（從起點到目標）
func _draw_chain_line(from_x: float, from_y: float, to_x: float, to_y: float) -> void:
	# 用多個小 ColorRect 模擬連鎖線
	var dx = to_x - from_x
	var dy = to_y - from_y
	var dist = sqrt(dx * dx + dy * dy)
	if dist < 1.0:
		return

	var steps = int(dist / 20.0)
	if steps < 2:
		steps = 2

	for i in range(steps):
		var t = float(i) / float(steps)
		var px = from_x + dx * t
		var py = from_y + dy * t

		var dot = ColorRect.new()
		dot.color = COLOR_CHAIN
		dot.size = Vector2(6, 6)
		dot.position = Vector2(px - 3, py - 3)
		add_child(dot)

		var delay = float(i) * 0.015
		var tw = dot.create_tween()
		tw.tween_interval(delay)
		tw.tween_property(dot, "modulate:a", 1.0, 0.0)  # 立即顯示
		tw.tween_interval(0.3)
		tw.tween_property(dot, "modulate:a", 0.0, 0.2)
		tw.tween_callback(dot.queue_free)

## chain_broken — 連鎖中斷
func _on_chain_broken(payload: Dictionary) -> void:
	var layer_num: int = payload.get("layer", 1)
	var vp_size = get_viewport().size

	# 灰色閃光
	_flash_screen(COLOR_BROKEN, 0.08)

	# 中斷提示
	var broken_label = Label.new()
	broken_label.text = "🔗 連鎖中斷（第 %d 層，範圍內無目標）" % layer_num
	broken_label.add_theme_font_size_override("font_size", 15)
	broken_label.add_theme_color_override("font_color", COLOR_BROKEN)
	broken_label.position = Vector2(vp_size.x / 2 - 160, vp_size.y / 2 - 10)
	add_child(broken_label)

	var tw = broken_label.create_tween()
	tw.tween_interval(2.0)
	tw.tween_property(broken_label, "modulate:a", 0.0, 0.4)
	tw.tween_callback(broken_label.queue_free)

	# 清除層數計數器
	if is_instance_valid(_layer_counter):
		var tw2 = _layer_counter.create_tween()
		tw2.tween_property(_layer_counter, "modulate:a", 0.0, 0.3)
		tw2.tween_callback(func():
			if is_instance_valid(_layer_counter):
				_layer_counter.queue_free()
				_layer_counter = null
		)

## chain_complete — 連鎖完成（8層全部引爆）
func _on_chain_complete(payload: Dictionary) -> void:
	var layer_num: int = payload.get("layer", 8)
	var vp_size = get_viewport().size

	# 三次強閃光
	_flash_screen(COLOR_CHAIN, 0.20)
	await get_tree().create_timer(0.12).timeout
	_flash_screen(COLOR_FIRE, 0.18)
	await get_tree().create_timer(0.10).timeout
	_flash_screen(COLOR_GOLD, 0.15)

	# 大字
	var big_label = Label.new()
	big_label.text = "🔗 連鎖完成！%d 層！" % layer_num
	big_label.add_theme_font_size_override("font_size", 48)
	big_label.add_theme_color_override("font_color", COLOR_GOLD)
	big_label.position = Vector2(vp_size.x / 2 - 160, vp_size.y / 2 - 60)
	add_child(big_label)

	big_label.scale = Vector2(0.6, 0.6)
	var tw_big = big_label.create_tween()
	tw_big.tween_property(big_label, "scale", Vector2(1.15, 1.15), 0.18)
	tw_big.tween_property(big_label, "scale", Vector2(1.0, 1.0), 0.10)
	tw_big.tween_interval(2.5)
	tw_big.tween_property(big_label, "modulate:a", 0.0, 0.5)
	tw_big.tween_callback(big_label.queue_free)

	# 結算彈窗（右側滑入）
	if _total_reward > 0:
		var popup = Control.new()
		popup.position = Vector2(vp_size.x + 10, vp_size.y / 2 - 50)
		popup.size = Vector2(180, 100)
		add_child(popup)

		var popup_bg = ColorRect.new()
		popup_bg.color = Color(0.12, 0.06, 0.0, 0.88)
		popup_bg.size = Vector2(180, 100)
		popup_bg.position = Vector2.ZERO
		popup.add_child(popup_bg)

		var popup_border = ColorRect.new()
		popup_border.color = COLOR_CHAIN
		popup_border.size = Vector2(180, 3)
		popup_border.position = Vector2(0, 0)
		popup.add_child(popup_border)

		var popup_title = Label.new()
		popup_title.text = "🔗 連鎖結算"
		popup_title.add_theme_font_size_override("font_size", 14)
		popup_title.add_theme_color_override("font_color", COLOR_CHAIN)
		popup_title.position = Vector2(10, 10)
		popup.add_child(popup_title)

		var popup_layers = Label.new()
		popup_layers.text = "連鎖層數：%d 層" % layer_num
		popup_layers.add_theme_font_size_override("font_size", 13)
		popup_layers.add_theme_color_override("font_color", COLOR_PALE)
		popup_layers.position = Vector2(10, 34)
		popup.add_child(popup_layers)

		var popup_reward = Label.new()
		popup_reward.text = "連鎖獎勵：+%d" % _total_reward
		popup_reward.add_theme_font_size_override("font_size", 16)
		popup_reward.add_theme_color_override("font_color", COLOR_GOLD)
		popup_reward.position = Vector2(10, 56)
		popup.add_child(popup_reward)

		var popup_player = Label.new()
		popup_player.text = _chain_player
		popup_player.add_theme_font_size_override("font_size", 11)
		popup_player.add_theme_color_override("font_color", Color(0.7, 0.7, 0.7))
		popup_player.position = Vector2(10, 80)
		popup.add_child(popup_player)

		# 滑入動畫
		var tw_popup = popup.create_tween()
		tw_popup.tween_property(popup, "position:x", vp_size.x - 190, 0.3)
		tw_popup.tween_interval(3.5)
		tw_popup.tween_property(popup, "position:x", vp_size.x + 10, 0.3)
		tw_popup.tween_callback(popup.queue_free)

	# 清除層數計數器
	if is_instance_valid(_layer_counter):
		var tw3 = _layer_counter.create_tween()
		tw3.tween_property(_layer_counter, "modulate:a", 0.0, 0.4)
		tw3.tween_callback(func():
			if is_instance_valid(_layer_counter):
				_layer_counter.queue_free()
				_layer_counter = null
		)

	# 重置狀態
	_chain_layer = 0
	_total_reward = 0

## 全螢幕閃光效果
func _flash_screen(color: Color, duration: float) -> void:
	var vp_size = get_viewport().size
	var flash = ColorRect.new()
	flash.color = Color(color.r, color.g, color.b, 0.42)
	flash.size = vp_size
	flash.position = Vector2.ZERO
	add_child(flash)

	var tw = flash.create_tween()
	tw.tween_property(flash, "modulate:a", 0.0, duration)
	tw.tween_callback(flash.queue_free)
