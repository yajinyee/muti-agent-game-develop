## VortexFishPanel.gd — 漩渦魚群吸引面板（DAY-169）
## 業界依據：Ocean King（Google Play 2026）「Vortex Fish — catching a Vortex Fish will suck
## all fish of the same species in the area into a whirlpool, capturing them all at once.」
## 視覺設計：
##   - vortex_start：全螢幕深藍閃光 + 頂部橫幅滑入 + 漩渦旋轉動畫（中央）+ 目標數量提示
##   - 自己觸發時：中央大 🌀 標誌彈跳動畫 + 「漩渦吸引中！」提示
##   - vortex_suck（每個目標）：目標飛向漩渦中心的動畫 + 吸入計數器 + 小閃光
##   - vortex_end：全螢幕藍色爆炸閃光 + 右側滑入結果彈窗（吸入數/獎勵）
##   - ≥5個：雙閃光；≥8個：彩虹三閃光
extends Node2D

# ---- 常數 ----
const SCREEN_W := 1280.0
const SCREEN_H := 720.0

# ---- 狀態 ----
var _pixel_font: Font = null
var _banner: Node2D = null          # 頂部橫幅
var _vortex_center: Node2D = null   # 漩渦中心動畫節點
var _suck_counter_lbl: Label = null # 吸入計數器
var _is_my_vortex: bool = false     # 是否是自己觸發的漩渦
var _vortex_x: float = SCREEN_W / 2.0
var _vortex_y: float = SCREEN_H / 2.0
var _target_count: int = 0          # 預計吸入目標數
var _sucked_count: int = 0          # 已吸入數
var _is_active: bool = false        # 是否正在漩渦中
var _vortex_rotation: float = 0.0   # 漩渦旋轉角度

# ---- 初始化 ----
func _ready() -> void:
	if ResourceLoader.exists("res://assets/fonts/pixel8.fnt"):
		_pixel_font = load("res://assets/fonts/pixel8.fnt")
	_connect_signals()

func _connect_signals() -> void:
	if GameManager.has_signal("vortex_fish"):
		GameManager.vortex_fish.connect(_on_vortex_fish)

# ---- 漩渦旋轉動畫 ----
func _process(delta: float) -> void:
	if not _is_active:
		return
	_vortex_rotation += delta * 180.0  # 每秒旋轉 180 度
	if is_instance_valid(_vortex_center):
		_vortex_center.rotation_degrees = _vortex_rotation

# ---- 主要事件處理 ----
func _on_vortex_fish(data: Dictionary) -> void:
	var phase: String = data.get("phase", "")
	var trigger_id: String = data.get("trigger_id", "")
	var trigger_name: String = data.get("trigger_name", "漩渦魚")
	var vortex_x: float = data.get("vortex_x", SCREEN_W / 2.0)
	var vortex_y: float = data.get("vortex_y", SCREEN_H / 2.0)
	var group_name: String = data.get("group_name", "基礎目標群")
	var target_count: int = data.get("target_count", 0)

	match phase:
		"vortex_start":
			_start_vortex(trigger_id, trigger_name, vortex_x, vortex_y, group_name, target_count)
		"vortex_suck":
			var suck_entry = data.get("suck_entry", null)
			var suck_index: int = data.get("suck_index", 0)
			_on_suck(suck_entry, suck_index)
		"vortex_end":
			var killed_count: int = data.get("killed_count", 0)
			var total_reward: int = data.get("total_reward", 0)
			_end_vortex(killed_count, total_reward)

# ---- 漩渦開始 ----
func _start_vortex(trigger_id: String, trigger_name: String, vx: float, vy: float, group_name: String, target_count: int) -> void:
	_is_active = true
	_vortex_x = vx
	_vortex_y = vy
	_target_count = target_count
	_sucked_count = 0
	_vortex_rotation = 0.0

	# 判斷是否是自己觸發
	var my_id: String = ""
	if GameManager.has_method("get_player_id"):
		my_id = GameManager.get_player_id()
	_is_my_vortex = (trigger_id == my_id)

	# 全螢幕深藍閃光
	_flash_screen(Color(0.0, 0.5, 1.0, 0.6), 0.4)

	# 建立漩渦中心動畫
	_create_vortex_center(vx, vy)

	# 建立頂部橫幅
	_create_banner(trigger_name, group_name, target_count)

	# 建立吸入計數器
	_create_suck_counter()

	# 自己觸發時：中央大 🌀 標誌彈跳
	if _is_my_vortex:
		_show_my_trigger_anim()

# ---- 建立漩渦中心動畫 ----
func _create_vortex_center(vx: float, vy: float) -> void:
	if is_instance_valid(_vortex_center):
		_vortex_center.queue_free()

	_vortex_center = Node2D.new()
	_vortex_center.position = Vector2(vx, vy)
	add_child(_vortex_center)

	# 漩渦外圈（深藍色）
	for i in range(3):
		var ring = ColorRect.new()
		var r := 40.0 + i * 20.0
		ring.size = Vector2(r * 2, r * 2)
		ring.position = Vector2(-r, -r)
		ring.color = Color(0.0, 0.3 + i * 0.2, 1.0, 0.3 - i * 0.08)
		_vortex_center.add_child(ring)

	# 漩渦中心圓
	var center_dot = ColorRect.new()
	center_dot.size = Vector2(20, 20)
	center_dot.position = Vector2(-10, -10)
	center_dot.color = Color(0.0, 0.8, 1.0, 0.9)
	_vortex_center.add_child(center_dot)

	# 漩渦圖示
	var icon_lbl = Label.new()
	icon_lbl.text = "🌀"
	icon_lbl.position = Vector2(-16, -16)
	if _pixel_font:
		icon_lbl.add_theme_font_override("font", _pixel_font)
	icon_lbl.add_theme_font_size_override("font_size", 28)
	_vortex_center.add_child(icon_lbl)

# ---- 建立頂部橫幅 ----
func _create_banner(trigger_name: String, group_name: String, target_count: int) -> void:
	if is_instance_valid(_banner):
		_banner.queue_free()

	_banner = Node2D.new()
	_banner.position = Vector2(SCREEN_W / 2.0, -60)
	add_child(_banner)

	# 橫幅背景
	var bg = ColorRect.new()
	bg.size = Vector2(600, 52)
	bg.position = Vector2(-300, -26)
	bg.color = Color(0.0, 0.2, 0.6, 0.88)
	_banner.add_child(bg)

	# 橫幅文字
	var lbl = Label.new()
	lbl.text = "🌀 %s 觸發漩渦魚！吸入 %d 個%s！" % [trigger_name, target_count, group_name]
	lbl.position = Vector2(-290, -18)
	lbl.add_theme_color_override("font_color", Color(0.4, 0.9, 1.0))
	if _pixel_font:
		lbl.add_theme_font_override("font", _pixel_font)
	lbl.add_theme_font_size_override("font_size", 16)
	_banner.add_child(lbl)

	# 橫幅滑入動畫
	var tween = create_tween()
	tween.tween_property(_banner, "position:y", 36.0, 0.3).set_ease(Tween.EASE_OUT)

# ---- 建立吸入計數器 ----
func _create_suck_counter() -> void:
	if is_instance_valid(_suck_counter_lbl):
		_suck_counter_lbl.queue_free()

	_suck_counter_lbl = Label.new()
	_suck_counter_lbl.text = "吸入：0 / %d" % _target_count
	_suck_counter_lbl.position = Vector2(SCREEN_W / 2.0 - 80, SCREEN_H - 80)
	_suck_counter_lbl.add_theme_color_override("font_color", Color(0.4, 0.9, 1.0))
	if _pixel_font:
		_suck_counter_lbl.add_theme_font_override("font", _pixel_font)
	_suck_counter_lbl.add_theme_font_size_override("font_size", 18)
	add_child(_suck_counter_lbl)

# ---- 自己觸發動畫 ----
func _show_my_trigger_anim() -> void:
	var anim_node = Node2D.new()
	anim_node.position = Vector2(SCREEN_W / 2.0, SCREEN_H / 2.0)
	add_child(anim_node)

	var lbl = Label.new()
	lbl.text = "🌀"
	lbl.position = Vector2(-24, -24)
	if _pixel_font:
		lbl.add_theme_font_override("font", _pixel_font)
	lbl.add_theme_font_size_override("font_size", 48)
	anim_node.add_child(lbl)

	var sub_lbl = Label.new()
	sub_lbl.text = "漩渦吸引中！"
	sub_lbl.position = Vector2(-60, 30)
	sub_lbl.add_theme_color_override("font_color", Color(0.4, 0.9, 1.0))
	if _pixel_font:
		sub_lbl.add_theme_font_override("font", _pixel_font)
	sub_lbl.add_theme_font_size_override("font_size", 16)
	anim_node.add_child(sub_lbl)

	# 彈跳動畫
	var tween = create_tween()
	tween.tween_property(anim_node, "scale", Vector2(1.4, 1.4), 0.15).set_ease(Tween.EASE_OUT)
	tween.tween_property(anim_node, "scale", Vector2(1.0, 1.0), 0.1)
	tween.tween_interval(1.0)
	tween.tween_property(anim_node, "modulate:a", 0.0, 0.3)
	tween.tween_callback(func(): if is_instance_valid(anim_node): anim_node.queue_free())

# ---- 目標被吸入 ----
func _on_suck(suck_entry, suck_index: int) -> void:
	if suck_entry == null:
		return

	_sucked_count += 1

	# 更新計數器
	if is_instance_valid(_suck_counter_lbl):
		_suck_counter_lbl.text = "吸入：%d / %d" % [_sucked_count, _target_count]

	# 目標飛向漩渦中心的動畫
	var entry_x: float = suck_entry.get("x", _vortex_x)
	var entry_y: float = suck_entry.get("y", _vortex_y)
	var reward: int = suck_entry.get("reward", 0)

	_spawn_suck_particle(entry_x, entry_y, reward)

	# 小閃光
	_flash_screen(Color(0.0, 0.6, 1.0, 0.15), 0.1)

# ---- 生成吸入粒子動畫 ----
func _spawn_suck_particle(from_x: float, from_y: float, reward: int) -> void:
	var particle = Node2D.new()
	particle.position = Vector2(from_x, from_y)
	add_child(particle)

	# 粒子圖示
	var lbl = Label.new()
	lbl.text = "💧"
	lbl.position = Vector2(-8, -8)
	if _pixel_font:
		lbl.add_theme_font_override("font", _pixel_font)
	lbl.add_theme_font_size_override("font_size", 16)
	particle.add_child(lbl)

	# 獎勵文字
	if reward > 0:
		var reward_lbl = Label.new()
		reward_lbl.text = "+%d" % reward
		reward_lbl.position = Vector2(10, -8)
		reward_lbl.add_theme_color_override("font_color", Color(0.4, 1.0, 0.8))
		if _pixel_font:
			reward_lbl.add_theme_font_override("font", _pixel_font)
		reward_lbl.add_theme_font_size_override("font_size", 12)
		particle.add_child(reward_lbl)

	# 飛向漩渦中心
	var tween = create_tween()
	tween.tween_property(particle, "position", Vector2(_vortex_x, _vortex_y), 0.3).set_ease(Tween.EASE_IN)
	tween.tween_property(particle, "modulate:a", 0.0, 0.1)
	tween.tween_callback(func(): if is_instance_valid(particle): particle.queue_free())

# ---- 漩渦結束 ----
func _end_vortex(killed_count: int, total_reward: int) -> void:
	_is_active = false

	# 清理漩渦中心
	if is_instance_valid(_vortex_center):
		var tween = create_tween()
		tween.tween_property(_vortex_center, "scale", Vector2(2.0, 2.0), 0.2)
		tween.tween_property(_vortex_center, "modulate:a", 0.0, 0.2)
		tween.tween_callback(func(): if is_instance_valid(_vortex_center): _vortex_center.queue_free())

	# 清理橫幅
	if is_instance_valid(_banner):
		var tween2 = create_tween()
		tween2.tween_property(_banner, "modulate:a", 0.0, 0.3)
		tween2.tween_callback(func(): if is_instance_valid(_banner): _banner.queue_free())

	# 清理計數器
	if is_instance_valid(_suck_counter_lbl):
		var tween3 = create_tween()
		tween3.tween_property(_suck_counter_lbl, "modulate:a", 0.0, 0.3)
		tween3.tween_callback(func(): if is_instance_valid(_suck_counter_lbl): _suck_counter_lbl.queue_free())

	# 爆炸閃光
	if killed_count >= 8:
		# 彩虹三閃光
		_flash_screen(Color(0.0, 0.8, 1.0, 0.7), 0.15)
		await get_tree().create_timer(0.2).timeout
		_flash_screen(Color(0.4, 0.0, 1.0, 0.6), 0.15)
		await get_tree().create_timer(0.2).timeout
		_flash_screen(Color(0.0, 1.0, 0.5, 0.5), 0.15)
	elif killed_count >= 5:
		# 雙閃光
		_flash_screen(Color(0.0, 0.6, 1.0, 0.6), 0.15)
		await get_tree().create_timer(0.2).timeout
		_flash_screen(Color(0.0, 0.6, 1.0, 0.4), 0.15)
	else:
		_flash_screen(Color(0.0, 0.5, 1.0, 0.5), 0.2)

	# 顯示結果彈窗
	await get_tree().create_timer(0.3).timeout
	_show_result_popup(killed_count, total_reward)

# ---- 結果彈窗 ----
func _show_result_popup(killed_count: int, total_reward: int) -> void:
	var popup = Node2D.new()
	popup.position = Vector2(SCREEN_W + 200, SCREEN_H / 2.0 - 80)
	add_child(popup)

	# 彈窗背景
	var bg = ColorRect.new()
	bg.size = Vector2(260, 160)
	bg.position = Vector2(-130, -80)
	bg.color = Color(0.0, 0.1, 0.3, 0.92)
	popup.add_child(bg)

	# 邊框
	var border = ColorRect.new()
	border.size = Vector2(264, 164)
	border.position = Vector2(-132, -82)
	border.color = Color(0.0, 0.6, 1.0, 0.8)
	popup.add_child(border)
	popup.move_child(border, 0)

	# 標題
	var title_lbl = Label.new()
	title_lbl.text = "🌀 漩渦魚大豐收！"
	title_lbl.position = Vector2(-120, -70)
	title_lbl.add_theme_color_override("font_color", Color(0.4, 0.9, 1.0))
	if _pixel_font:
		title_lbl.add_theme_font_override("font", _pixel_font)
	title_lbl.add_theme_font_size_override("font_size", 16)
	popup.add_child(title_lbl)

	# 吸入數
	var killed_lbl = Label.new()
	killed_lbl.text = "吸入目標：%d 個" % killed_count
	killed_lbl.position = Vector2(-110, -30)
	killed_lbl.add_theme_color_override("font_color", Color(0.8, 0.95, 1.0))
	if _pixel_font:
		killed_lbl.add_theme_font_override("font", _pixel_font)
	killed_lbl.add_theme_font_size_override("font_size", 14)
	popup.add_child(killed_lbl)

	# 獎勵
	var reward_lbl = Label.new()
	reward_lbl.text = "獲得獎勵：🪙%d" % total_reward
	reward_lbl.position = Vector2(-110, 0)
	reward_lbl.add_theme_color_override("font_color", Color(1.0, 0.9, 0.3))
	if _pixel_font:
		reward_lbl.add_theme_font_override("font", _pixel_font)
	reward_lbl.add_theme_font_size_override("font_size", 14)
	popup.add_child(reward_lbl)

	# 評語
	var comment_lbl = Label.new()
	if killed_count >= 8:
		comment_lbl.text = "🌊 傳說漩渦！"
		comment_lbl.add_theme_color_override("font_color", Color(0.0, 1.0, 1.0))
	elif killed_count >= 5:
		comment_lbl.text = "💧 大豐收！"
		comment_lbl.add_theme_color_override("font_color", Color(0.4, 0.9, 1.0))
	else:
		comment_lbl.text = "🌀 漩渦完成"
		comment_lbl.add_theme_color_override("font_color", Color(0.7, 0.85, 1.0))
	comment_lbl.position = Vector2(-110, 35)
	if _pixel_font:
		comment_lbl.add_theme_font_override("font", _pixel_font)
	comment_lbl.add_theme_font_size_override("font_size", 14)
	popup.add_child(comment_lbl)

	# 滑入動畫
	var tween = create_tween()
	tween.tween_property(popup, "position:x", SCREEN_W - 160.0, 0.35).set_ease(Tween.EASE_OUT)
	tween.tween_interval(2.5)
	tween.tween_property(popup, "modulate:a", 0.0, 0.4)
	tween.tween_callback(func(): if is_instance_valid(popup): popup.queue_free())

# ---- 全螢幕閃光 ----
func _flash_screen(color: Color, duration: float) -> void:
	var flash = ColorRect.new()
	flash.size = Vector2(SCREEN_W, SCREEN_H)
	flash.color = color
	add_child(flash)
	var tween = create_tween()
	tween.tween_property(flash, "modulate:a", 0.0, duration)
	tween.tween_callback(func(): if is_instance_valid(flash): flash.queue_free())
