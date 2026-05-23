## LuckyChainBombPanel.gd — 幸運鏈鎖爆炸魚系統面板（DAY-226）
## 業界原創「連鎖爆炸」機制
##
## 視覺設計：
##   - 橙紅爆炸主題（#FF4500 + #FF6B35 + #FFD700 + #FFF0E6）
##   - chain_bomb_start：橙紅三次強閃光 + 頂部橫幅 + 引爆標記菱形（橙紅色）
##   - chain_bomb_trigger：全螢幕橙色閃光 + 「💣 引爆！」大字 + 爆炸圓圈
##   - chain_bomb_blast：爆炸圓圈擴散 + 浮動獎勵文字 + 連鎖層數標記
##   - chain_bomb_expire：橙色淡出
extends CanvasLayer

# 引爆標記節點（targetID → Control）
var _mark_nodes: Dictionary = {}
var _active_instance: String = ""

# 主題顏色
const COLOR_PRIMARY   = Color("#FF4500")  # 橙紅
const COLOR_ORANGE    = Color("#FF6B35")  # 橙色
const COLOR_GOLD      = Color("#FFD700")  # 金黃
const COLOR_PALE      = Color("#FFF0E6")  # 極淡橙
const COLOR_BG        = Color(0.12, 0.04, 0.0, 0.88)

func _ready() -> void:
	layer = 19  # 幸運鏈鎖爆炸魚面板層級

## 處理幸運鏈鎖爆炸魚訊息
func handle_lucky_chain_bomb(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"chain_bomb_start":
			_on_chain_bomb_start(payload)
		"chain_bomb_trigger":
			_on_chain_bomb_trigger(payload)
		"chain_bomb_blast":
			_on_chain_bomb_blast(payload)
		"chain_bomb_expire":
			_on_chain_bomb_expire(payload)

## chain_bomb_start — 引爆標記建立（全服廣播）
func _on_chain_bomb_start(payload: Dictionary) -> void:
	var instance_id: String = payload.get("instance_id", "")
	var player_name: String = payload.get("player_name", "")
	var marked: Array = payload.get("marked", [])
	var duration_sec: int = payload.get("duration_sec", 15)

	_active_instance = instance_id

	# 橙紅三次強閃光
	_flash_screen(COLOR_PRIMARY, 0.15)
	await get_tree().create_timer(0.1).timeout
	_flash_screen(COLOR_ORANGE, 0.12)
	await get_tree().create_timer(0.1).timeout
	_flash_screen(COLOR_GOLD, 0.1)

	# 頂部橫幅
	var vp_size = get_viewport().size
	var banner = Label.new()
	banner.text = "💣 %s 觸發鏈鎖爆炸！%d 個目標被引爆！" % [player_name, marked.size()]
	banner.add_theme_font_size_override("font_size", 15)
	banner.add_theme_color_override("font_color", COLOR_GOLD)
	banner.position = Vector2(vp_size.x / 2 - 160, 6)
	add_child(banner)

	var tween_banner = banner.create_tween()
	tween_banner.tween_interval(3.0)
	tween_banner.tween_property(banner, "modulate:a", 0.0, 0.5)
	tween_banner.tween_callback(banner.queue_free)

	# 為每個引爆標記目標建立菱形標記
	for info in marked:
		var target_id: String = info.get("target_id", "")
		var tx: float = info.get("x", 0.0)
		var ty: float = info.get("y", 0.0)
		_create_bomb_mark(target_id, tx, ty, duration_sec)

## 建立引爆標記（菱形，橙紅色）
func _create_bomb_mark(target_id: String, tx: float, ty: float, duration_sec: int) -> void:
	var container = Control.new()
	container.position = Vector2(tx - 20, ty - 20)
	container.size = Vector2(40, 40)
	add_child(container)
	_mark_nodes[target_id] = container

	# 菱形輪廓（4 個 ColorRect 組成菱形）
	var mark_size = 6
	var offsets = [
		Vector2(17, 0),   # 上
		Vector2(17, 34),  # 下
		Vector2(0, 17),   # 左
		Vector2(34, 17),  # 右
	]
	for offset in offsets:
		var rect = ColorRect.new()
		rect.color = COLOR_PRIMARY
		rect.size = Vector2(mark_size, mark_size)
		rect.position = offset
		container.add_child(rect)

	# 中心 💣 符號
	var icon = Label.new()
	icon.text = "💣"
	icon.add_theme_font_size_override("font_size", 14)
	icon.position = Vector2(10, 10)
	container.add_child(icon)

	# 閃爍動畫
	var tween = container.create_tween().set_loops()
	tween.tween_property(container, "modulate", Color(2.0, 1.0, 0.5, 1.0), 0.4)
	tween.tween_property(container, "modulate", Color.WHITE, 0.4)

	# 計時條（底部）
	var timer_bar = ColorRect.new()
	timer_bar.color = COLOR_ORANGE
	timer_bar.size = Vector2(40, 3)
	timer_bar.position = Vector2(0, 38)
	container.add_child(timer_bar)

	var tween_timer = timer_bar.create_tween()
	tween_timer.tween_property(timer_bar, "size:x", 0.0, float(duration_sec))

## chain_bomb_trigger — 引爆標記目標被擊破，連鎖開始
func _on_chain_bomb_trigger(payload: Dictionary) -> void:
	var trigger_id: String = payload.get("trigger_id", "")
	var trigger_x: float = payload.get("trigger_x", 0.0)
	var trigger_y: float = payload.get("trigger_y", 0.0)
	var player_name: String = payload.get("player_name", "")

	# 移除引爆標記節點
	_remove_mark(trigger_id)

	# 全螢幕橙色閃光
	_flash_screen(COLOR_PRIMARY, 0.2)

	# 「💣 引爆！」大字
	var vp_size = get_viewport().size
	var big_label = Label.new()
	big_label.text = "💣 引爆！"
	big_label.add_theme_font_size_override("font_size", 48)
	big_label.add_theme_color_override("font_color", COLOR_GOLD)
	big_label.position = vp_size / 2 - Vector2(80, 28)
	add_child(big_label)

	var tween_label = big_label.create_tween()
	tween_label.tween_property(big_label, "scale", Vector2(1.3, 1.3), 0.1)
	tween_label.tween_property(big_label, "scale", Vector2(1.0, 1.0), 0.08)
	tween_label.tween_interval(0.3)
	tween_label.tween_property(big_label, "modulate:a", 0.0, 0.3)
	tween_label.tween_callback(big_label.queue_free)

	# 爆炸圓圈（在觸發位置）
	_spawn_explosion_circle(trigger_x, trigger_y, 200.0, COLOR_PRIMARY)

	# 頂部橫幅
	var banner = Label.new()
	banner.text = "💥 %s 觸發連鎖！" % player_name
	banner.add_theme_font_size_override("font_size", 13)
	banner.add_theme_color_override("font_color", COLOR_ORANGE)
	banner.position = Vector2(vp_size.x / 2 - 80, 6)
	add_child(banner)

	var tween_banner = banner.create_tween()
	tween_banner.tween_interval(2.0)
	tween_banner.tween_property(banner, "modulate:a", 0.0, 0.4)
	tween_banner.tween_callback(banner.queue_free)

## chain_bomb_blast — 連鎖爆炸結果
func _on_chain_bomb_blast(payload: Dictionary) -> void:
	var trigger_x: float = payload.get("trigger_x", 0.0)
	var trigger_y: float = payload.get("trigger_y", 0.0)
	var chain_layer: int = payload.get("chain_layer", 1)
	var blast_radius: float = payload.get("blast_radius", 200.0)
	var results: Array = payload.get("results", [])
	var total_reward: int = payload.get("total_reward", 0)
	var player_name: String = payload.get("player_name", "")

	# 爆炸圓圈（依連鎖層數改變顏色）
	var blast_color = COLOR_PRIMARY
	if chain_layer == 2:
		blast_color = COLOR_ORANGE
	elif chain_layer >= 3:
		blast_color = COLOR_GOLD
	_spawn_explosion_circle(trigger_x, trigger_y, blast_radius, blast_color)

	# 連鎖層數標記
	if chain_layer >= 2:
		var vp_size = get_viewport().size
		var chain_label = Label.new()
		chain_label.text = "🔗 第 %d 層連鎖！" % chain_layer
		chain_label.add_theme_font_size_override("font_size", 22)
		chain_label.add_theme_color_override("font_color", blast_color)
		chain_label.position = Vector2(trigger_x - 60, trigger_y - 40)
		add_child(chain_label)

		var tween_chain = chain_label.create_tween()
		tween_chain.tween_property(chain_label, "position:y", chain_label.position.y - 30, 0.5)
		tween_chain.parallel().tween_property(chain_label, "modulate:a", 0.0, 0.5)
		tween_chain.tween_callback(chain_label.queue_free)

	# 每個被擊破目標的浮動獎勵
	for result in results:
		if not result.get("killed", false):
			continue
		var rx: float = result.get("x", 0.0)
		var ry: float = result.get("y", 0.0)
		var reward: int = result.get("reward", 0)
		var is_chain: bool = result.get("is_chain", false)

		# 移除引爆標記（如果有）
		var target_id: String = result.get("target_id", "")
		if is_chain:
			_remove_mark(target_id)

		# 浮動獎勵文字
		if reward > 0:
			var reward_label = Label.new()
			reward_label.text = "+%d" % reward
			if is_chain:
				reward_label.text += " 💣"
			reward_label.add_theme_font_size_override("font_size", 14)
			reward_label.add_theme_color_override("font_color", COLOR_PALE)
			reward_label.position = Vector2(rx - 20, ry - 10)
			add_child(reward_label)

			var tween_r = reward_label.create_tween()
			tween_r.tween_property(reward_label, "position:y", reward_label.position.y - 35, 0.5)
			tween_r.parallel().tween_property(reward_label, "modulate:a", 0.0, 0.5)
			tween_r.tween_callback(reward_label.queue_free)

	# 總獎勵（連鎖層數 ≥ 2 時顯示）
	if chain_layer >= 2 and total_reward > 0:
		var vp_size = get_viewport().size
		var total_label = Label.new()
		total_label.text = "💥 連鎖 +%d！" % total_reward
		total_label.add_theme_font_size_override("font_size", 26)
		total_label.add_theme_color_override("font_color", COLOR_GOLD)
		total_label.position = vp_size / 2 - Vector2(70, 0)
		add_child(total_label)

		var tween_total = total_label.create_tween()
		tween_total.tween_property(total_label, "scale", Vector2(1.2, 1.2), 0.1)
		tween_total.tween_property(total_label, "scale", Vector2(1.0, 1.0), 0.08)
		tween_total.tween_interval(0.5)
		tween_total.tween_property(total_label, "modulate:a", 0.0, 0.4)
		tween_total.tween_callback(total_label.queue_free)

## chain_bomb_expire — 引爆標記過期
func _on_chain_bomb_expire(payload: Dictionary) -> void:
	var instance_id: String = payload.get("instance_id", "")
	if instance_id != _active_instance:
		return

	# 淡出所有引爆標記
	for target_id in _mark_nodes.keys():
		var node = _mark_nodes[target_id]
		if is_instance_valid(node):
			var tween = node.create_tween()
			tween.tween_property(node, "modulate:a", 0.0, 0.4)
			tween.tween_callback(node.queue_free)
	_mark_nodes.clear()
	_active_instance = ""

# ---- 輔助函數 ----

## 移除引爆標記節點
func _remove_mark(target_id: String) -> void:
	if _mark_nodes.has(target_id):
		var node = _mark_nodes[target_id]
		if is_instance_valid(node):
			# 消失閃光
			var tween = node.create_tween()
			tween.tween_property(node, "modulate", Color(3.0, 1.5, 0.5, 1.0), 0.06)
			tween.tween_property(node, "modulate:a", 0.0, 0.15)
			tween.tween_callback(node.queue_free)
		_mark_nodes.erase(target_id)

## 爆炸圓圈特效（在指定位置擴散）
func _spawn_explosion_circle(cx: float, cy: float, radius: float, color: Color) -> void:
	# 外圈（擴散）
	var outer = ColorRect.new()
	outer.color = Color(color.r, color.g, color.b, 0.0)
	var outer_size = radius * 2
	outer.size = Vector2(outer_size, outer_size)
	outer.position = Vector2(cx - radius, cy - radius)
	add_child(outer)

	# 用 Label 模擬圓圈（用 ○ 字符）
	var circle_label = Label.new()
	circle_label.text = "○"
	circle_label.add_theme_font_size_override("font_size", int(radius * 1.5))
	circle_label.add_theme_color_override("font_color", color)
	circle_label.position = Vector2(cx - radius * 0.75, cy - radius * 0.75)
	add_child(circle_label)

	var tween = circle_label.create_tween()
	tween.tween_property(circle_label, "scale", Vector2(1.3, 1.3), 0.25)
	tween.parallel().tween_property(circle_label, "modulate:a", 0.0, 0.25)
	tween.tween_callback(circle_label.queue_free)
	outer.queue_free()

## 全螢幕閃光效果
func _flash_screen(color: Color, duration: float) -> void:
	var flash = ColorRect.new()
	flash.color = Color(color.r, color.g, color.b, 0.35)
	flash.set_anchors_preset(Control.PRESET_FULL_RECT)
	add_child(flash)

	var tween = flash.create_tween()
	tween.tween_property(flash, "modulate:a", 0.0, duration)
	tween.tween_callback(flash.queue_free)
