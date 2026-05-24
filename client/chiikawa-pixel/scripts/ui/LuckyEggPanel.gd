## LuckyEggPanel.gd ??撟賊?敶抵?擳?選?DAY-172嚗?
## 璆剔?靘?嚗ILI Mega Fishing 2026?iant Prize Fish lets you easily win great prizes,
## with the chance for 5x multipliers?? Ocean King 2026?gg Fish drops golden eggs??
## 閬死閮剛?嚗?
##   - egg_start嚗??嚗??典?璈怠???鈭箄孛?澆兢?蔗??嚗? 敶抵???
##   - egg_open嚗犖嚗?敶抵?敺孛?潔?蝵桅???+ ??? + ?瘚桀???
##     - coins嚗??脣蔗??+ ?馳?典???+ "+XXX ?馳" 瘚桀???
##     - mult嚗?蝝蔗??+ "?2 ?? 5s" 瘚桀??? + ?喃?閫閮?
##     - weapon嚗予?蔗??+ "甇血? ?1" 瘚桀??? + 甇血?內??
##   - egg_result嚗犖嚗??喳皛蝯?敶?嚗蔗?/?馳/??/甇血嚗?
##   - mult_end嚗犖嚗??喃?閫閮?瘛∪
##   - ???蔗???冽?撱?璈怠?嚗5???????
extends Node2D

# ---- 撣豢 ----
const SCREEN_W := 1280.0
const SCREEN_H := 720.0
const EGG_SIZE := 32.0
const EGG_COLORS = {
	"coins":  Color(1.0, 0.85, 0.0),   # ?
	"mult":   Color(1.0, 0.41, 0.71),  # 蝎???
	"weapon": Color(0.0, 0.75, 1.0),   # 憭抵???
}
const EGG_ICONS = {
	"coins":  "??",
	"mult":   "??,"
	"weapon": "??,"
}

# ---- ???----
var _pixel_font: Font = null
var _mult_countdown_lbl: Label = null   # ???閮?璅惜
var _mult_elapsed: float = 0.0          # ??撌脤???
var _mult_duration: float = 5.0         # ??????
var _is_mult_active: bool = false       # ?臬??瞈瘣颱葉
var _mult_stack: int = 0                # ????甈⊥嚗??蔗???嚗?
var _egg_nodes: Array = []              # 敶抵?蝭暺?銵?
var _result_panel: Node2D = null        # 蝯?敶?蝭暺?

# ---- ????----
func _ready() -> void:
	if ResourceLoader.exists("res://assets/fonts/pixel8.fnt"):
		_pixel_font = load("res://assets/fonts/pixel8.fnt")
	_connect_signals()

func _connect_signals() -> void:
	if GameManager.has_signal("lucky_egg_fish"):
		GameManager.lucky_egg_fish.connect(_on_lucky_egg_fish)

# ---- 閮???----
func _process(delta: float) -> void:
	# ???閮?
	if _is_mult_active:
		_mult_elapsed += delta
		var remaining = _mult_duration - _mult_elapsed
		if remaining <= 0.0:
			_is_mult_active = false
			if is_instance_valid(_mult_countdown_lbl):
				_mult_countdown_lbl.queue_free()
				_mult_countdown_lbl = null
		elif is_instance_valid(_mult_countdown_lbl):
			_mult_countdown_lbl.text = "?2 %.1fs" % remaining

# ---- 閮??? ----
func _on_lucky_egg_fish(data: Dictionary) -> void:
	var phase = data.get("phase", "")
	match phase:
		"egg_start":
			_handle_egg_start(data)
		"egg_open":
			_handle_egg_open(data)
		"egg_result":
			_handle_egg_result(data)
		"egg_broadcast":
			_handle_egg_broadcast(data)
		"mult_end":
			_handle_mult_end()

# ---- egg_start嚗?誨?剜帖撟?----
func _handle_egg_start(data: Dictionary) -> void:
	var player_name = data.get("player_name", "?拙振")
	var egg_count = data.get("egg_count", 1)
	_show_broadcast_banner("?? %s 閫貊撟賊?敶抵?擳?? %d ?蔗??" % [player_name, egg_count])

# ---- egg_open嚗犖敶抵???? ----
func _handle_egg_open(data: Dictionary) -> void:
	var egg_result = data.get("egg_result", {})
	var reward_type = egg_result.get("reward_type", "coins")
	var label_text = egg_result.get("label", "")
	var egg_index = data.get("egg_index", 0)
	var trigger_x = data.get("trigger_x", SCREEN_W / 2.0)
	var trigger_y = data.get("trigger_y", SCREEN_H / 2.0)

	# 敶抵?憿
	var egg_color = EGG_COLORS.get(reward_type, Color.WHITE)
	var egg_icon = EGG_ICONS.get(reward_type, "??")

	# 撱箇?敶抵?蝭暺?敺孛?潔?蝵桅??綽?
	var egg_node = Node2D.new()
	add_child(egg_node)
	_egg_nodes.append(egg_node)

	# 敶抵??耦
	var egg_circle = ColorRect.new()
	egg_circle.size = Vector2(EGG_SIZE, EGG_SIZE)
	egg_circle.position = Vector2(-EGG_SIZE / 2.0, -EGG_SIZE / 2.0)
	egg_circle.color = egg_color
	egg_node.add_child(egg_circle)

	# 敶抵??內
	var icon_lbl = Label.new()
	icon_lbl.text = egg_icon
	icon_lbl.position = Vector2(-12, -14)
	if _pixel_font:
		icon_lbl.add_theme_font_override("font", _pixel_font)
		icon_lbl.add_theme_font_size_override("font_size", 20)
	egg_node.add_child(icon_lbl)

	# 韏瑕?雿蔭嚗孛?潔?蝵殷?
	egg_node.position = Vector2(trigger_x, trigger_y)

	# 憌?格?雿蔭嚗????恍銝剖亢嚗?
	var spread_x = SCREEN_W / 2.0 + (egg_index - 2) * 80.0
	var spread_y = SCREEN_H / 2.0 - 50.0

	# 憌?
	var tween = egg_node.create_tween()
	tween.tween_property(egg_node, "position",
		Vector2(spread_x, spread_y), 0.3).set_ease(Tween.EASE_OUT)
	tween.tween_interval(0.1)

	# ???嚗葬?曄??賂?
	tween.tween_property(egg_node, "scale", Vector2(1.5, 1.5), 0.1)
	tween.tween_property(egg_node, "scale", Vector2(0.8, 0.8), 0.1)
	tween.tween_property(egg_node, "scale", Vector2(1.0, 1.0), 0.1)

	# 憿舐內?瘚桀???
	tween.tween_callback(func():
		_show_reward_float(spread_x, spread_y - 40.0, label_text, egg_color)
	)

	# ?寞???嚗???
	if reward_type == "mult":
		tween.tween_callback(func():
			_activate_mult_display()
		)

	# 瘛∪敶抵?
	tween.tween_interval(0.5)
	tween.tween_property(egg_node, "modulate:a", 0.0, 0.3)
	tween.tween_callback(func():
		if is_instance_valid(egg_node):
			egg_node.queue_free()
		_egg_nodes.erase(egg_node)
	)

# ---- egg_result嚗犖蝯?敶? ----
func _handle_egg_result(data: Dictionary) -> void:
	var egg_count = data.get("egg_count", 1)
	var total_coins = data.get("total_coins", 0)
	var mult_count = data.get("mult_count", 0)
	var weapon_count = data.get("weapon_count", 0)

	# 撱箇?蝯?敶?嚗?湔??伐?
	if is_instance_valid(_result_panel):
		_result_panel.queue_free()

	_result_panel = Node2D.new()
	add_child(_result_panel)

	# ?
	var bg = ColorRect.new()
	bg.size = Vector2(220, 140)
	bg.position = Vector2(0, -70)
	bg.color = Color(0.1, 0.1, 0.1, 0.85)
	_result_panel.add_child(bg)

	# 璅?
	var title_lbl = Label.new()
	title_lbl.text = "?? 敶抵?蝯?"
	title_lbl.position = Vector2(10, -65)
	if _pixel_font:
		title_lbl.add_theme_font_override("font", _pixel_font)
		title_lbl.add_theme_font_size_override("font_size", 14)
	title_lbl.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	_result_panel.add_child(title_lbl)

	# 敶抵???
	var count_lbl = Label.new()
	count_lbl.text = "敶抵??賂?%d ?? % egg_count"
	count_lbl.position = Vector2(10, -45)
	if _pixel_font:
		count_lbl.add_theme_font_override("font", _pixel_font)
		count_lbl.add_theme_font_size_override("font_size", 12)
	count_lbl.add_theme_color_override("font_color", Color.WHITE)
	_result_panel.add_child(count_lbl)

	# ?馳?
	if total_coins > 0:
		var coins_lbl = Label.new()
		coins_lbl.text = "?? ?馳嚗?%d" % total_coins
		coins_lbl.position = Vector2(10, -25)
		if _pixel_font:
			coins_lbl.add_theme_font_override("font", _pixel_font)
			coins_lbl.add_theme_font_size_override("font_size", 12)
		coins_lbl.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
		_result_panel.add_child(coins_lbl)

	# ????
	if mult_count > 0:
		var mult_lbl = Label.new()
		mult_lbl.text = "??????嚗?d" % mult_count
		mult_lbl.position = Vector2(10, -5)
		if _pixel_font:
			mult_lbl.add_theme_font_override("font", _pixel_font)
			mult_lbl.add_theme_font_size_override("font_size", 12)
		mult_lbl.add_theme_color_override("font_color", Color(1.0, 0.41, 0.71))
		_result_panel.add_child(mult_lbl)

	# 甇血?
	if weapon_count > 0:
		var weapon_lbl = Label.new()
		weapon_lbl.text = "??甇血?嚗?d" % weapon_count
		weapon_lbl.position = Vector2(10, 15)
		if _pixel_font:
			weapon_lbl.add_theme_font_override("font", _pixel_font)
			weapon_lbl.add_theme_font_size_override("font_size", 12)
		weapon_lbl.add_theme_color_override("font_color", Color(0.0, 0.75, 1.0))
		_result_panel.add_child(weapon_lbl)

	# 敺?湔???
	_result_panel.position = Vector2(SCREEN_W + 50, SCREEN_H / 2.0)
	var tween = _result_panel.create_tween()
	tween.tween_property(_result_panel, "position:x", SCREEN_W - 240.0, 0.4).set_ease(Tween.EASE_OUT)

	# ???蔗???????
	if egg_count >= 5:
		_flash_screen(Color(1.0, 0.85, 0.0, 0.5))
		tween.tween_interval(0.15)
		tween.tween_callback(func(): _flash_screen(Color(1.0, 0.85, 0.0, 0.5)))

	# 3 蝘?瘛∪
	tween.tween_interval(3.0)
	tween.tween_property(_result_panel, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func():
		if is_instance_valid(_result_panel):
			_result_panel.queue_free()
			_result_panel = null
	)

# ---- egg_broadcast嚗?誨?剜帖撟?----
func _handle_egg_broadcast(data: Dictionary) -> void:
	var player_name = data.get("player_name", "?拙振")
	var egg_count = data.get("egg_count", 1)
	var total_coins = data.get("total_coins", 0)
	var mult_count = data.get("mult_count", 0)

	var msg = "?? %s 撟賊?敶抵?擳???%d ?蔗??" % [player_name, egg_count]
	if mult_count > 0:
		msg += " %d 甈∪???嚗? % mult_count"
	if total_coins > 0:
		msg += " +%d ?馳嚗? % total_coins"
	_show_broadcast_banner(msg)

# ---- mult_end嚗?蝯? ----
func _handle_mult_end() -> void:
	_is_mult_active = false
	_mult_stack = 0
	if is_instance_valid(_mult_countdown_lbl):
		var tween = _mult_countdown_lbl.create_tween()
		tween.tween_property(_mult_countdown_lbl, "modulate:a", 0.0, 0.3)
		tween.tween_callback(func():
			if is_instance_valid(_mult_countdown_lbl):
				_mult_countdown_lbl.queue_free()
				_mult_countdown_lbl = null
		)

# ---- 頛嚗?瘣餃?憿舐內 ----
func _activate_mult_display() -> void:
	_is_mult_active = true
	_mult_elapsed = 0.0
	_mult_duration = 5.0
	_mult_stack += 1

	# 撱箇??喃?閫閮?璅惜
	if not is_instance_valid(_mult_countdown_lbl):
		_mult_countdown_lbl = Label.new()
		add_child(_mult_countdown_lbl)
		_mult_countdown_lbl.position = Vector2(SCREEN_W - 120, 60)
		if _pixel_font:
			_mult_countdown_lbl.add_theme_font_override("font", _pixel_font)
			_mult_countdown_lbl.add_theme_font_size_override("font_size", 16)
		_mult_countdown_lbl.add_theme_color_override("font_color", Color(1.0, 0.41, 0.71))

	_mult_countdown_lbl.text = "?2 5.0s"

	# 敶歲?
	var tween = _mult_countdown_lbl.create_tween()
	tween.tween_property(_mult_countdown_lbl, "scale", Vector2(1.4, 1.4), 0.1)
	tween.tween_property(_mult_countdown_lbl, "scale", Vector2(1.0, 1.0), 0.15)

# ---- 頛嚗??菜筑??摮?----
func _show_reward_float(x: float, y: float, text: String, color: Color) -> void:
	var lbl = Label.new()
	lbl.text = text
	lbl.position = Vector2(x - 40, y)
	if _pixel_font:
		lbl.add_theme_font_override("font", _pixel_font)
		lbl.add_theme_font_size_override("font_size", 14)
	lbl.add_theme_color_override("font_color", color)
	add_child(lbl)

	var tween = lbl.create_tween()
	tween.tween_property(lbl, "position:y", y - 40.0, 0.8).set_ease(Tween.EASE_OUT)
	tween.parallel().tween_property(lbl, "modulate:a", 0.0, 0.8)
	tween.tween_callback(func():
		if is_instance_valid(lbl):
			lbl.queue_free()
	)

# ---- 頛嚗?誨?剜帖撟?----
func _show_broadcast_banner(text: String) -> void:
	var banner = ColorRect.new()
	banner.size = Vector2(SCREEN_W, 36)
	banner.position = Vector2(0, -36)
	banner.color = Color(0.1, 0.1, 0.1, 0.85)
	add_child(banner)

	var lbl = Label.new()
	lbl.text = text
	lbl.position = Vector2(10, 6)
	if _pixel_font:
		lbl.add_theme_font_override("font", _pixel_font)
		lbl.add_theme_font_size_override("font_size", 14)
	lbl.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	banner.add_child(lbl)

	# 敺??冽???
	banner.position = Vector2(0, 0)
	var tween = banner.create_tween()
	tween.tween_interval(2.5)
	tween.tween_property(banner, "position:y", -36.0, 0.3)
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
