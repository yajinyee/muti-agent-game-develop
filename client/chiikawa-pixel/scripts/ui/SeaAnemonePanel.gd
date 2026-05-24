## SeaAnemonePanel.gd ??瘚瑁閫豢??餅??Ｘ嚗AY-174嚗?
## 璆剔?靘?嚗ILI Jackpot Fishing?ea Anemone introduces unique effects ??
## tentacle attacks that spread to nearby fish??
## 閬死閮剛?嚗?
##   - tentacle_start嚗??嚗絲?萎葉敹??+ ?璈怠???鈭箄孛?潭絲?蛛???
##   - tentacle_hit嚗??嚗?銝剖??璅?辣隡貊?蝝孛?? + ?賭葉??
##     - isKill嚗璅???+ 瘚桀????
##     - !isKill嚗璅????嚗?
##   - tentacle_miss嚗??嚗孛?辣隡詨?楠敺?憭?
##   - tentacle_result嚗??嚗?湔??亦???蝒????蝮賜??蛛?
##   - ?????湛??冽?撱?璈怠?嚗6??蝎?????
extends Node2D

# ---- 撣豢 ----
const SCREEN_W := 1280.0
const SCREEN_H := 720.0
const TENTACLE_COLOR := Color(1.0, 0.41, 0.71, 0.85)  # 蝎???
const TENTACLE_WIDTH := 3.0
const TENTACLE_DURATION := 0.4  # 閫豢?憿舐內??嚗?嚗?

# ---- ???----
var _pixel_font: Font = null
var _tentacle_nodes: Array = []   # 閫豢?蝺?暺?
var _result_panel: Node2D = null  # 蝯?敶?

# ---- ????----
func _ready() -> void:
	if ResourceLoader.exists("res://assets/fonts/pixel8.fnt"):
		_pixel_font = load("res://assets/fonts/pixel8.fnt")
	_connect_signals()

func _connect_signals() -> void:
	if GameManager.has_signal("sea_anemone"):
		GameManager.sea_anemone.connect(_on_sea_anemone)

# ---- 閮??? ----
func _on_sea_anemone(data: Dictionary) -> void:
	var phase = data.get("phase", "")
	match phase:
		"tentacle_start":
			_handle_tentacle_start(data)
		"tentacle_hit":
			_handle_tentacle_hit(data)
		"tentacle_miss":
			_handle_tentacle_miss(data)
		"tentacle_result":
			_handle_tentacle_result(data)

# ---- tentacle_start嚗絲?菔孛???憪?----
func _handle_tentacle_start(data: Dictionary) -> void:
	var player_name = data.get("killer_name", "?拙振")
	var trigger_x = data.get("trigger_x", SCREEN_W / 2.0)
	var trigger_y = data.get("trigger_y", SCREEN_H / 2.0)

	# 銝剖????
	_flash_at(trigger_x, trigger_y, TENTACLE_COLOR, 40.0)

	# ?璈怠?
	_show_broadcast_banner("?疙 %s ?絲?菔孛???撱嗡撓嚗? % player_name)"

# ---- tentacle_hit嚗孛?銝剔璅?----
func _handle_tentacle_hit(data: Dictionary) -> void:
	var trigger_x = data.get("trigger_x", SCREEN_W / 2.0)
	var trigger_y = data.get("trigger_y", SCREEN_H / 2.0)
	var hit_x = data.get("hit_x", trigger_x)
	var hit_y = data.get("hit_y", trigger_y)
	var is_kill = data.get("is_kill", false)
	var reward = data.get("reward", 0)
	var multiplier = data.get("multiplier", 1.0)

	# 蝜芾ˊ閫豢?蝺?敺葉敹?格?嚗?
	_draw_tentacle_line(trigger_x, trigger_y, hit_x, hit_y)

	# ?賭葉??
	var flash_color = Color(1.0, 0.2, 0.2, 0.7) if is_kill else Color(1.0, 0.8, 0.0, 0.5)
	_flash_at(hit_x, hit_y, flash_color, 24.0)

	# ??＊蝷箇??菜筑??摮?
	if is_kill and reward > 0:
		_show_reward_float(hit_x, hit_y - 20.0, "+%d" % reward, Color(1.0, 0.85, 0.0))

# ---- tentacle_miss嚗孛??賭葉 ----
func _handle_tentacle_miss(data: Dictionary) -> void:
	var trigger_x = data.get("trigger_x", SCREEN_W / 2.0)
	var trigger_y = data.get("trigger_y", SCREEN_H / 2.0)
	var angle = data.get("angle", 0.0)

	# 蝜芾ˊ閫豢?蝺?楠嚗??哨?銵函內?芸銝哨?
	var angle_rad = deg_to_rad(angle)
	var end_x = trigger_x + cos(angle_rad) * 150.0
	var end_y = trigger_y + sin(angle_rad) * 150.0
	_draw_tentacle_line(trigger_x, trigger_y, end_x, end_y, 0.5)  # ??

# ---- tentacle_result嚗孛?????----
func _handle_tentacle_result(data: Dictionary) -> void:
	var kill_count = data.get("kill_count", 0)
	var total_reward = data.get("total_reward", 0)
	var player_name = data.get("killer_name", "?拙振")

	if kill_count == 0:
		return

	# 撱箇?蝯?敶?嚗?湔??伐?
	if is_instance_valid(_result_panel):
		_result_panel.queue_free()

	_result_panel = Node2D.new()
	add_child(_result_panel)

	# ?
	var bg = ColorRect.new()
	bg.size = Vector2(200, 100)
	bg.position = Vector2(0, -50)
	bg.color = Color(0.1, 0.0, 0.1, 0.85)
	_result_panel.add_child(bg)

	# 璅?
	var title_lbl = Label.new()
	title_lbl.text = "?疙 瘚瑁閫豢??餅?"
	title_lbl.position = Vector2(10, -45)
	if _pixel_font:
		title_lbl.add_theme_font_override("font", _pixel_font)
		title_lbl.add_theme_font_size_override("font_size", 13)
	title_lbl.add_theme_color_override("font_color", TENTACLE_COLOR)
	_result_panel.add_child(title_lbl)

	# ???
	var kill_lbl = Label.new()
	kill_lbl.text = "?嚗?d ?璅? % kill_count"
	kill_lbl.position = Vector2(10, -25)
	if _pixel_font:
		kill_lbl.add_theme_font_override("font", _pixel_font)
		kill_lbl.add_theme_font_size_override("font_size", 12)
	kill_lbl.add_theme_color_override("font_color", Color.WHITE)
	_result_panel.add_child(kill_lbl)

	# 蝮賜???
	var reward_lbl = Label.new()
	reward_lbl.text = "?? ?嚗?%d" % total_reward
	reward_lbl.position = Vector2(10, -5)
	if _pixel_font:
		reward_lbl.add_theme_font_override("font", _pixel_font)
		reward_lbl.add_theme_font_size_override("font_size", 12)
	reward_lbl.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	_result_panel.add_child(reward_lbl)

	# 敺?湔???
	_result_panel.position = Vector2(SCREEN_W + 50, SCREEN_H / 2.0)
	var tween = _result_panel.create_tween()
	tween.tween_property(_result_panel, "position:x", SCREEN_W - 220.0, 0.4).set_ease(Tween.EASE_OUT)

	# ?????湛?蝎?????
	if kill_count >= 6:
		_flash_screen(Color(1.0, 0.41, 0.71, 0.4))
		tween.tween_interval(0.15)
		tween.tween_callback(func(): _flash_screen(Color(1.0, 0.41, 0.71, 0.4)))

	# 3 蝘?瘛∪
	tween.tween_interval(3.0)
	tween.tween_property(_result_panel, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func():
		if is_instance_valid(_result_panel):
			_result_panel.queue_free()
			_result_panel = null
	)

# ---- 頛嚗鼓鋆質孛?? ----
func _draw_tentacle_line(from_x: float, from_y: float, to_x: float, to_y: float, alpha: float = 1.0) -> void:
	var line = Line2D.new()
	line.add_point(Vector2(from_x, from_y))
	line.add_point(Vector2(to_x, to_y))
	line.width = TENTACLE_WIDTH
	var color = TENTACLE_COLOR
	color.a = alpha
	line.default_color = color
	add_child(line)
	_tentacle_nodes.append(line)

	# 閫豢?瘛∪
	var tween = line.create_tween()
	tween.tween_interval(TENTACLE_DURATION * 0.5)
	tween.tween_property(line, "modulate:a", 0.0, TENTACLE_DURATION * 0.5)
	tween.tween_callback(func():
		if is_instance_valid(line):
			line.queue_free()
		_tentacle_nodes.erase(line)
	)

# ---- 頛嚗?蝵桅???----
func _flash_at(x: float, y: float, color: Color, size: float) -> void:
	var flash = ColorRect.new()
	flash.size = Vector2(size, size)
	flash.position = Vector2(x - size / 2.0, y - size / 2.0)
	flash.color = color
	add_child(flash)

	var tween = flash.create_tween()
	tween.tween_property(flash, "modulate:a", 0.0, 0.2)
	tween.tween_callback(func():
		if is_instance_valid(flash):
			flash.queue_free()
	)

# ---- 頛嚗??菜筑??摮?----
func _show_reward_float(x: float, y: float, text: String, color: Color) -> void:
	var lbl = Label.new()
	lbl.text = text
	lbl.position = Vector2(x - 20, y)
	if _pixel_font:
		lbl.add_theme_font_override("font", _pixel_font)
		lbl.add_theme_font_size_override("font_size", 13)
	lbl.add_theme_color_override("font_color", color)
	add_child(lbl)

	var tween = lbl.create_tween()
	tween.tween_property(lbl, "position:y", y - 35.0, 0.7).set_ease(Tween.EASE_OUT)
	tween.parallel().tween_property(lbl, "modulate:a", 0.0, 0.7)
	tween.tween_callback(func():
		if is_instance_valid(lbl):
			lbl.queue_free()
	)

# ---- 頛嚗?誨?剜帖撟?----
func _show_broadcast_banner(text: String) -> void:
	var banner = ColorRect.new()
	banner.size = Vector2(SCREEN_W, 34)
	banner.position = Vector2(0, 0)
	banner.color = Color(0.1, 0.0, 0.1, 0.85)
	add_child(banner)

	var lbl = Label.new()
	lbl.text = text
	lbl.position = Vector2(10, 6)
	if _pixel_font:
		lbl.add_theme_font_override("font", _pixel_font)
		lbl.add_theme_font_size_override("font_size", 13)
	lbl.add_theme_color_override("font_color", TENTACLE_COLOR)
	banner.add_child(lbl)

	var tween = banner.create_tween()
	tween.tween_interval(2.5)
	tween.tween_property(banner, "modulate:a", 0.0, 0.3)
	tween.tween_callback(func():
		if is_instance_valid(banner):
			banner.queue_free()
	)

# ---- 頛嚗?Ｗ??? ----
func _flash_screen(color: Color) -> void:
	var flash = ColorRect.new()
	flash.size = Vector2(SCREEN_W, SCREEN_H)
	flash.position = Vector2(0, 0)
	flash.color = color
	add_child(flash)

	var tween = flash.create_tween()
	tween.tween_property(flash, "modulate:a", 0.0, 0.25)
	tween.tween_callback(func():
		if is_instance_valid(flash):
			flash.queue_free()
	)
