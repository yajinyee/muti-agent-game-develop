## LuckyTimeBombFishPanel.gd — 幸運時間炸彈魚系統面板（DAY-235）
## 業界原創「倒數計時+提前引爆+連鎖爆炸」機制
##
## 視覺設計：
##   - 紅橙炸彈主題（#E74C3C + #C0392B + #F1948A + #FDEDEC）
##   - bomb_placed：紅色雙閃光 + 頂部橫幅 + 「💣 時間炸彈！」大字 + 炸彈標記 + 倒數提示
##   - bomb_countdown：倒數數字更新（每秒）
##   - bomb_early_detonate：橙色強閃光 + 「💥 提前引爆！×2.0」大字
##   - bomb_chain_blast：全螢幕三次強閃光 + 「💥 連鎖爆炸！」大字 + 結算彈窗
##   - bomb_auto_explode：紅色閃光 + 「💣 自動爆炸！」提示
extends CanvasLayer

# 炸彈狀態
var _active: bool = false
var _fuse_sec: int = 8
var _bomb_count: int = 0

# 倒數標籤（targetID → Label）
var _countdown_labels: Dictionary = {}

# 計時條節點
var _timer_bar: ColorRect = null
var _timer_bar_bg: ColorRect = null

# 主題顏色
const COLOR_PRIMARY  = Color("#E74C3C")  # 紅色
const COLOR_DARK     = Color("#C0392B")  # 深紅
const COLOR_PALE     = Color("#F1948A")  # 淡紅
const COLOR_LIGHT_BG = Color("#FDEDEC")  # 極淡紅
const COLOR_ORANGE   = Color("#E67E22")  # 橙色
const COLOR_GOLD     = Color("#FFD700")  # 金黃

func _ready() -> void:
	layer = 10  # 幸運時間炸彈魚面板層級

## 處理幸運時間炸彈魚訊息
func handle_lucky_time_bomb_fish(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"bomb_placed":
			_on_bomb_placed(payload)
		"bomb_countdown":
			_on_bomb_countdown(payload)
		"bomb_early_detonate":
			_on_bomb_early_detonate(payload)
		"bomb_chain_blast":
			_on_bomb_chain_blast(payload)
		"bomb_auto_explode":
			_on_bomb_auto_explode(payload)

## bomb_placed — 炸彈標記放置
func _on_bomb_placed(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "")
	_bomb_count = payload.get("bomb_count", 4)
	_fuse_sec = payload.get("fuse_sec", 8)
	var early_mult: float = payload.get("early_mult", 2.0)
	var auto_mult: float = payload.get("auto_mult", 1.6)
	_active = true

	var vp_size = get_viewport().size

	# 紅色雙閃光
	_flash_screen(COLOR_PRIMARY, 0.12)
	await get_tree().create_timer(0.08).timeout
	_flash_screen(COLOR_DARK, 0.10)

	# 頂部橫幅
	var banner = Label.new()
	banner.text = "💣 %s 放置了 %d 個時間炸彈！" % [player_name, _bomb_count]
	banner.add_theme_font_size_override("font_size", 14)
	banner.add_theme_color_override("font_color", COLOR_PALE)
	banner.position = Vector2(vp_size.x / 2 - 140, 6)
	add_child(banner)

	var tween_banner = banner.create_tween()
	tween_banner.tween_interval(float(_fuse_sec) - 0.5)
	tween_banner.tween_property(banner, "modulate:a", 0.0, 0.5)
	tween_banner.tween_callback(banner.queue_free)

	# 「💣 時間炸彈！」大字
	var big_label = Label.new()
	big_label.text = "💣 時間炸彈！"
	big_label.add_theme_font_size_override("font_size", 48)
	big_label.add_theme_color_override("font_color", COLOR_PRIMARY)
	big_label.position = vp_size / 2 - Vector2(100, 28)
	add_child(big_label)

	var tween_big = big_label.create_tween()
	tween_big.tween_property(big_label, "scale", Vector2(1.2, 1.2), 0.10)
	tween_big.tween_property(big_label, "scale", Vector2(1.0, 1.0), 0.08)
	tween_big.tween_interval(0.5)
	tween_big.tween_property(big_label, "modulate:a", 0.0, 0.4)
	tween_big.tween_callback(big_label.queue_free)

	# 倍率說明
	var mult_label = Label.new()
	mult_label.text = "提前引爆 ×%.1f | 自動爆炸 ×%.1f" % [early_mult, auto_mult]
	mult_label.add_theme_font_size_override("font_size", 12)
	mult_label.add_theme_color_override("font_color", COLOR_GOLD)
	mult_label.position = Vector2(vp_size.x / 2 - 90, vp_size.y / 2 + 28)
	add_child(mult_label)

	var tween_mult = mult_label.create_tween()
	tween_mult.tween_interval(2.5)
	tween_mult.tween_property(mult_label, "modulate:a", 0.0, 0.5)
	tween_mult.tween_callback(mult_label.queue_free)

	# 底部計時條（紅→深紅漸變）
	_spawn_timer_bar(float(_fuse_sec))

## bomb_countdown — 倒數更新（每秒）
func _on_bomb_countdown(payload: Dictionary) -> void:
	var target_id: String = payload.get("target_id", "")
	var remaining: int = payload.get("remaining", 0)

	# 更新或建立倒數標籤
	if target_id in _countdown_labels:
		var lbl = _countdown_labels[target_id]
		if is_instance_valid(lbl):
			lbl.text = "💣 %d" % remaining
			# 最後 3 秒變橙色閃爍
			if remaining <= 3:
				lbl.add_theme_color_override("font_color", COLOR_ORANGE)
				var tween_blink = lbl.create_tween()
				tween_blink.tween_property(lbl, "modulate:a", 0.3, 0.15)
				tween_blink.tween_property(lbl, "modulate:a", 1.0, 0.15)

## bomb_early_detonate — 提前引爆
func _on_bomb_early_detonate(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "")
	var target_id: String = payload.get("target_id", "")
	var mult: float = payload.get("mult", 2.0)
	var x: float = payload.get("x", 0.0)
	var y: float = payload.get("y", 0.0)

	# 移除倒數標籤
	if target_id in _countdown_labels:
		var lbl = _countdown_labels[target_id]
		if is_instance_valid(lbl):
			lbl.queue_free()
		_countdown_labels.erase(target_id)

	var vp_size = get_viewport().size

	# 橙色強閃光
	_flash_screen(COLOR_ORANGE, 0.10)

	# 「💥 提前引爆！×2.0」大字
	var detonate_label = Label.new()
	detonate_label.text = "💥 提前引爆！×%.1f" % mult
	detonate_label.add_theme_font_size_override("font_size", 36)
	detonate_label.add_theme_color_override("font_color", COLOR_ORANGE)
	detonate_label.position = vp_size / 2 - Vector2(90, 22)
	add_child(detonate_label)

	var tween_det = detonate_label.create_tween()
	tween_det.tween_property(detonate_label, "scale", Vector2(1.15, 1.15), 0.08)
	tween_det.tween_property(detonate_label, "scale", Vector2(1.0, 1.0), 0.07)
	tween_det.tween_interval(0.5)
	tween_det.tween_property(detonate_label, "modulate:a", 0.0, 0.4)
	tween_det.tween_callback(detonate_label.queue_free)

	# 爆炸圓圈（在目標位置）
	_spawn_explosion_ring(vp_size, x, y, COLOR_ORANGE)

	# 玩家名稱提示
	if player_name != "":
		var name_label = Label.new()
		name_label.text = player_name
		name_label.add_theme_font_size_override("font_size", 11)
		name_label.add_theme_color_override("font_color", COLOR_PALE)
		name_label.position = Vector2(vp_size.x / 2 - 30, vp_size.y / 2 + 20)
		add_child(name_label)

		var tween_name = name_label.create_tween()
		tween_name.tween_interval(0.8)
		tween_name.tween_property(name_label, "modulate:a", 0.0, 0.3)
		tween_name.tween_callback(name_label.queue_free)

## bomb_chain_blast — 連鎖爆炸
func _on_bomb_chain_blast(payload: Dictionary) -> void:
	var killed_count: int = payload.get("killed_count", 0)
	var total_reward: int = payload.get("total_reward", 0)
	var source_x: float = payload.get("source_x", 0.0)
	var source_y: float = payload.get("source_y", 0.0)

	if killed_count == 0:
		return

	var vp_size = get_viewport().size

	# 全螢幕三次強閃光（橙→紅→深紅）
	_flash_screen(COLOR_ORANGE, 0.12)
	await get_tree().create_timer(0.08).timeout
	_flash_screen(COLOR_PRIMARY, 0.10)
	await get_tree().create_timer(0.07).timeout
	_flash_screen(COLOR_DARK, 0.12)

	# 「💥 連鎖爆炸！」大字
	var chain_label = Label.new()
	chain_label.text = "💥 連鎖爆炸！"
	chain_label.add_theme_font_size_override("font_size", 44)
	chain_label.add_theme_color_override("font_color", COLOR_ORANGE)
	chain_label.position = vp_size / 2 - Vector2(95, 28)
	add_child(chain_label)

	var tween_chain = chain_label.create_tween()
	tween_chain.tween_property(chain_label, "scale", Vector2(1.25, 1.25), 0.12)
	tween_chain.tween_property(chain_label, "scale", Vector2(1.0, 1.0), 0.10)
	tween_chain.tween_interval(0.6)
	tween_chain.tween_property(chain_label, "modulate:a", 0.0, 0.5)
	tween_chain.tween_callback(chain_label.queue_free)

	# 爆炸圓圈（在連鎖源位置）
	_spawn_explosion_ring(vp_size, source_x, source_y, COLOR_PRIMARY)

	# 結算彈窗（右側滑入）
	_spawn_chain_result_panel(vp_size, killed_count, total_reward)

## bomb_auto_explode — 自動爆炸（倒數結束）
func _on_bomb_auto_explode(payload: Dictionary) -> void:
	var target_id: String = payload.get("target_id", "")
	var killed: bool = payload.get("killed", false)
	var reward: int = payload.get("reward", 0)
	var x: float = payload.get("x", 0.0)
	var y: float = payload.get("y", 0.0)

	# 移除倒數標籤
	if target_id in _countdown_labels:
		var lbl = _countdown_labels[target_id]
		if is_instance_valid(lbl):
			lbl.queue_free()
		_countdown_labels.erase(target_id)

	var vp_size = get_viewport().size

	# 紅色閃光
	_flash_screen(COLOR_PRIMARY, 0.08)

	# 「💣 自動爆炸！」提示
	var auto_label = Label.new()
	if killed:
		auto_label.text = "💣 自動爆炸！+%d" % reward
		auto_label.add_theme_color_override("font_color", COLOR_GOLD)
	else:
		auto_label.text = "💣 自動爆炸！未命中"
		auto_label.add_theme_color_override("font_color", COLOR_PALE)
	auto_label.add_theme_font_size_override("font_size", 18)
	auto_label.position = Vector2(vp_size.x / 2 - 70, vp_size.y / 2 - 50)
	add_child(auto_label)

	var tween_auto = auto_label.create_tween()
	tween_auto.tween_property(auto_label, "position:y", auto_label.position.y - 20, 0.4)
	tween_auto.parallel().tween_property(auto_label, "modulate:a", 0.0, 0.4)
	tween_auto.tween_callback(auto_label.queue_free)

	# 爆炸圓圈
	_spawn_explosion_ring(vp_size, x, y, COLOR_PRIMARY)

	# 檢查是否所有炸彈都已爆炸
	if _countdown_labels.is_empty():
		_active = false
		# 清除計時條
		if is_instance_valid(_timer_bar):
			_timer_bar.queue_free()
			_timer_bar = null
		if is_instance_valid(_timer_bar_bg):
			_timer_bar_bg.queue_free()
			_timer_bar_bg = null

## 建立爆炸圓圈（在目標位置）
func _spawn_explosion_ring(vp_size: Vector2, world_x: float, world_y: float, color: Color) -> void:
	var scale_x = vp_size.x / 1000.0
	var scale_y = vp_size.y / 600.0
	var screen_x = world_x * scale_x
	var screen_y = world_y * scale_y
	var radius = 40.0

	var ring = ColorRect.new()
	ring.color = Color(color.r, color.g, color.b, 0.5)
	ring.size = Vector2(radius * 2, radius * 2)
	ring.position = Vector2(screen_x - radius, screen_y - radius)
	add_child(ring)

	var tween_ring = ring.create_tween()
	tween_ring.tween_property(ring, "scale", Vector2(2.5, 2.5), 0.25)
	tween_ring.parallel().tween_property(ring, "modulate:a", 0.0, 0.25)
	tween_ring.tween_callback(ring.queue_free)

## 建立連鎖爆炸結算彈窗
func _spawn_chain_result_panel(vp_size: Vector2, killed_count: int, total_reward: int) -> void:
	var panel = ColorRect.new()
	panel.color = Color(0.25, 0.05, 0.05, 0.88)
	panel.size = Vector2(200, 90)
	panel.position = Vector2(vp_size.x + 10, vp_size.y / 2 - 45)
	add_child(panel)

	var title = Label.new()
	title.text = "💥 連鎖爆炸結算"
	title.add_theme_font_size_override("font_size", 13)
	title.add_theme_color_override("font_color", COLOR_PALE)
	title.position = Vector2(10, 8)
	panel.add_child(title)

	var kill_label = Label.new()
	kill_label.text = "連鎖擊破：%d 個目標" % killed_count
	kill_label.add_theme_font_size_override("font_size", 12)
	kill_label.add_theme_color_override("font_color", Color.WHITE)
	kill_label.position = Vector2(10, 30)
	panel.add_child(kill_label)

	var reward_label = Label.new()
	reward_label.text = "個人獎勵：+%d" % total_reward
	reward_label.add_theme_font_size_override("font_size", 12)
	reward_label.add_theme_color_override("font_color", COLOR_GOLD)
	reward_label.position = Vector2(10, 52)
	panel.add_child(reward_label)

	var tween_panel = panel.create_tween()
	tween_panel.tween_property(panel, "position:x", vp_size.x - 210, 0.3)
	tween_panel.tween_interval(2.5)
	tween_panel.tween_property(panel, "modulate:a", 0.0, 0.4)
	tween_panel.tween_callback(panel.queue_free)

## 建立底部計時條（紅→深紅漸變）
func _spawn_timer_bar(duration: float) -> void:
	var vp_size = get_viewport().size

	var bg = ColorRect.new()
	bg.color = Color(0.1, 0.1, 0.1, 0.6)
	bg.size = Vector2(vp_size.x, 8)
	bg.position = Vector2(0, vp_size.y - 8)
	add_child(bg)
	_timer_bar_bg = bg

	var bar = ColorRect.new()
	bar.color = COLOR_PRIMARY
	bar.size = Vector2(vp_size.x, 8)
	bar.position = Vector2(0, vp_size.y - 8)
	add_child(bar)
	_timer_bar = bar

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
	flash.color = Color(color.r, color.g, color.b, 0.26)
	flash.set_anchors_preset(Control.PRESET_FULL_RECT)
	add_child(flash)

	var tween = flash.create_tween()
	tween.tween_property(flash, "modulate:a", 0.0, duration)
	tween.tween_callback(flash.queue_free)
