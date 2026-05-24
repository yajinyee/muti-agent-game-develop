## ShopPanel.gd ?????Ｘ嚗AY-094嚗?
## 憿舐內???”???鞈?頃鞎瑕???
## ??? Server ?券?摨???
extends Node2D

# ---- 撣豢 ----
const PANEL_WIDTH  := 280
const PANEL_HEIGHT := 380
const PANEL_X      := 300
const PANEL_Y      := 80

# ---- 蝭暺???----
var _panel_bg: ColorRect
var _title_label: Label
var _toggle_btn: Button
var _flash_label: Label
var _flash_time_label: Label
var _items_container: VBoxContainer
var _scroll: ScrollContainer

# ---- ???----
var _is_expanded: bool = false
var _pixel_font: Font = null
var _items: Array = []
var _flash_seconds_left: int = 0

# ---- ????----
func _ready() -> void:
	if ResourceLoader.exists("res://assets/fonts/pixel8.fnt"):
		_pixel_font = load("res://assets/fonts/pixel8.fnt")
	_build_ui()
	_connect_signals()

func setup(font: Font) -> void:
	if font:
		_pixel_font = font

func _build_ui() -> void:
	# ??Ｘ
	_panel_bg = ColorRect.new()
	_panel_bg.position = Vector2(PANEL_X, PANEL_Y)
	_panel_bg.size = Vector2(PANEL_WIDTH, 36)
	_panel_bg.color = Color(0.08, 0.05, 0.15, 0.90)
	add_child(_panel_bg)

	# 璅???
	var title_bar := ColorRect.new()
	title_bar.position = Vector2(0, 0)
	title_bar.size = Vector2(PANEL_WIDTH, 36)
	title_bar.color = Color(0.3, 0.1, 0.5, 0.95)
	_panel_bg.add_child(title_bar)

	# 璅???
	var title := Label.new()
	title.position = Vector2(8, 6)
	title.text = "?? ??"
	title.add_theme_color_override("font_color", Color(1.0, 0.85, 1.0))
	if _pixel_font:
		title.add_theme_font_override("font", _pixel_font)
		title.add_theme_font_size_override("font_size", 14)
	title_bar.add_child(title)

	# 撅?/????
	_toggle_btn = Button.new()
	_toggle_btn.position = Vector2(PANEL_WIDTH - 32, 4)
	_toggle_btn.size = Vector2(28, 28)
	_toggle_btn.text = "??"
	_toggle_btn.flat = true
	_toggle_btn.add_theme_color_override("font_color", Color(0.9, 0.8, 1.0))
	title_bar.add_child(_toggle_btn)

	# ???寡都?嚗???銋＊蝷綽?
	_flash_label = Label.new()
	_flash_label.position = Vector2(8, 38)
	_flash_label.text = "?????寡都"
	_flash_label.add_theme_color_override("font_color", Color(1.0, 0.7, 0.0))
	if _pixel_font:
		_flash_label.add_theme_font_override("font", _pixel_font)
		_flash_label.add_theme_font_size_override("font_size", 11)
	_panel_bg.add_child(_flash_label)

	_flash_time_label = Label.new()
	_flash_time_label.position = Vector2(PANEL_WIDTH - 80, 38)
	_flash_time_label.text = ""
	_flash_time_label.add_theme_color_override("font_color", Color(1.0, 0.5, 0.0))
	if _pixel_font:
		_flash_time_label.add_theme_font_override("font", _pixel_font)
		_flash_time_label.add_theme_font_size_override("font_size", 10)
	_panel_bg.add_child(_flash_time_label)

	# ???”嚗???憿舐內嚗?
	_scroll = ScrollContainer.new()
	_scroll.position = Vector2(0, 58)
	_scroll.size = Vector2(PANEL_WIDTH, PANEL_HEIGHT - 60)
	_scroll.visible = false
	_panel_bg.add_child(_scroll)

	_items_container = VBoxContainer.new()
	_items_container.size = Vector2(PANEL_WIDTH - 8, 0)
	_scroll.add_child(_items_container)

func _connect_signals() -> void:
	_toggle_btn.pressed.connect(_on_toggle_pressed)
	if GameManager.has_signal("shop_updated"):
		GameManager.shop_updated.connect(_on_shop_updated)
	if GameManager.has_signal("shop_purchased"):
		GameManager.shop_purchased.connect(_on_shop_purchased)
	if GameManager.has_signal("shop_error"):
		GameManager.shop_error.connect(_on_shop_error)

# ---- 閮??? ----
func _on_toggle_pressed() -> void:
	_is_expanded = !_is_expanded
	_toggle_btn.text = "?? if _is_expanded else "??
	_scroll.visible = _is_expanded

	if _is_expanded:
		_panel_bg.size = Vector2(PANEL_WIDTH, PANEL_HEIGHT)
		_rebuild_items()
	else:
		_panel_bg.size = Vector2(PANEL_WIDTH, 58)

func _on_shop_updated(data: Dictionary) -> void:
	_items = data.get("items", [])
	_flash_seconds_left = data.get("seconds_left", 0)
	_update_flash_time_label()
	if _is_expanded:
		_rebuild_items()

func _on_shop_purchased(data: Dictionary) -> void:
	var item_name: String = data.get("item_name", "")
	var price: int = data.get("price", 0)
	# 憿舐內鞈潸眺???內
	_show_purchase_toast("??鞈潸眺??嚗?s" % item_name, Color(0.3, 1.0, 0.3))
	# ?瑟???”
	if _is_expanded:
		_rebuild_items()

func _on_shop_error(data: Dictionary) -> void:
	var reason: String = data.get("reason", "")
	var msg := "??鞈潸眺憭望?"
	match reason:
		"insufficient_coins": msg = "???馳銝雲"
		"daily_limit_reached": msg = "??隞撌脤?鞈潸眺銝?"
		"out_of_stock": msg = "????撌脣摰?"
		"item_not_found": msg = "????銝???"
	_show_purchase_toast(msg, Color(1.0, 0.3, 0.3))

# ---- UI ?湔 ----
func _update_flash_time_label() -> void:
	if _flash_seconds_left <= 0:
		_flash_time_label.text = "撌脩???"
		return
	var hours := _flash_seconds_left / 3600
	var mins := (_flash_seconds_left % 3600) / 60
	if hours > 0:
		_flash_time_label.text = "%dh%dm" % [hours, mins]
	else:
		_flash_time_label.text = "%dm" % mins

func _rebuild_items() -> void:
	for child in _items_container.get_children():
		child.queue_free()

	if _items.is_empty():
		var empty := Label.new()
		empty.text = "  ???怎??"
		empty.add_theme_color_override("font_color", Color(0.5, 0.5, 0.5))
		if _pixel_font:
			empty.add_theme_font_override("font", _pixel_font)
			empty.add_theme_font_size_override("font_size", 11)
		_items_container.add_child(empty)
		return

	for item in _items:
		_add_item_row(item)

func _add_item_row(item: Dictionary) -> void:
	var item_id: String = item.get("id", "")
	var name: String = item.get("name", "")
	var desc: String = item.get("description", "")
	var price: int = item.get("price", 0)
	var orig_price: int = item.get("orig_price", 0)
	var is_flash: bool = item.get("is_flash_sale", false)
	var limit: int = item.get("limit_per_day", 0)
	var purchased: int = item.get("purchased_today", 0)
	var can_buy: bool = (limit == 0 or purchased < limit)

	var row := VBoxContainer.new()
	row.custom_minimum_size = Vector2(PANEL_WIDTH - 8, 0)

	# ?
	var row_bg := ColorRect.new()
	if is_flash:
		row_bg.color = Color(0.3, 0.15, 0.0, 0.5)  # 璈?嚗鞈??
	else:
		row_bg.color = Color(0.1, 0.05, 0.2, 0.4)
	row_bg.size = Vector2(PANEL_WIDTH - 8, 52)
	row.add_child(row_bg)

	# ???迂銵?
	var name_row := HBoxContainer.new()
	name_row.position = Vector2(4, 2)
	row_bg.add_child(name_row)

	var name_label := Label.new()
	name_label.text = name
	name_label.custom_minimum_size = Vector2(160, 18)
	var name_color := Color(1.0, 0.85, 0.0) if is_flash else Color(0.9, 0.9, 1.0)
	name_label.add_theme_color_override("font_color", name_color)
	if _pixel_font:
		name_label.add_theme_font_override("font", _pixel_font)
		name_label.add_theme_font_size_override("font_size", 11)
	name_row.add_child(name_label)

	# ?寡都璅惜
	if is_flash:
		var flash_tag := Label.new()
		flash_tag.text = "?∠鞈?"
		flash_tag.add_theme_color_override("font_color", Color(1.0, 0.5, 0.0))
		if _pixel_font:
			flash_tag.add_theme_font_override("font", _pixel_font)
			flash_tag.add_theme_font_size_override("font_size", 9)
		name_row.add_child(flash_tag)

	# ?膩銵?
	var desc_label := Label.new()
	desc_label.position = Vector2(4, 20)
	desc_label.text = desc
	desc_label.add_theme_color_override("font_color", Color(0.7, 0.7, 0.8))
	if _pixel_font:
		desc_label.add_theme_font_override("font", _pixel_font)
		desc_label.add_theme_font_size_override("font_size", 9)
	row_bg.add_child(desc_label)

	# ?寞 + 鞈潸眺??銵?
	var price_row := HBoxContainer.new()
	price_row.position = Vector2(4, 32)
	row_bg.add_child(price_row)

	# ?寞憿舐內
	var price_text := ""
	if price == 0:
		price_text = "?祥"
	elif is_flash and orig_price > price:
		price_text = "??%d (??d)" % [price, orig_price]
	else:
		price_text = "??%d" % price

	var price_label := Label.new()
	price_label.text = price_text
	price_label.custom_minimum_size = Vector2(140, 16)
	price_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	if _pixel_font:
		price_label.add_theme_font_override("font", _pixel_font)
		price_label.add_theme_font_size_override("font_size", 10)
	price_row.add_child(price_label)

	# 鞈潸眺甈⊥憿舐內
	if limit > 0:
		var limit_label := Label.new()
		limit_label.text = "%d/%d" % [purchased, limit]
		limit_label.custom_minimum_size = Vector2(40, 16)
		limit_label.add_theme_color_override("font_color", Color(0.6, 0.6, 0.6))
		if _pixel_font:
			limit_label.add_theme_font_override("font", _pixel_font)
			limit_label.add_theme_font_size_override("font_size", 9)
		price_row.add_child(limit_label)

	# 鞈潸眺??
	var buy_btn := Button.new()
	buy_btn.text = "鞈潸眺" if can_buy else "撌脤?銝?"
	buy_btn.disabled = !can_buy
	buy_btn.custom_minimum_size = Vector2(60, 16)
	if can_buy:
		buy_btn.add_theme_color_override("font_color", Color(0.2, 1.0, 0.4))
	else:
		buy_btn.add_theme_color_override("font_color", Color(0.5, 0.5, 0.5))
	if _pixel_font:
		buy_btn.add_theme_font_override("font", _pixel_font)
		buy_btn.add_theme_font_size_override("font_size", 9)
	buy_btn.pressed.connect(func(): _on_buy_pressed(item_id))
	price_row.add_child(buy_btn)

	# ??蝺?
	var sep := ColorRect.new()
	sep.color = Color(0.3, 0.2, 0.4, 0.5)
	sep.custom_minimum_size = Vector2(PANEL_WIDTH - 8, 1)
	row.add_child(sep)

	_items_container.add_child(row)

func _on_buy_pressed(item_id: String) -> void:
	NetworkManager.send_buy_shop_item(item_id)

func _show_purchase_toast(msg: String, color: Color) -> void:
	var toast := Label.new()
	toast.text = msg
	toast.position = Vector2(PANEL_X + 10, PANEL_Y + PANEL_HEIGHT + 5)
	toast.add_theme_color_override("font_color", color)
	if _pixel_font:
		toast.add_theme_font_override("font", _pixel_font)
		toast.add_theme_font_size_override("font_size", 11)
	get_parent().add_child(toast)

	var tween := toast.create_tween()
	tween.tween_property(toast, "position:y", toast.position.y - 30, 1.5)
	tween.parallel().tween_property(toast, "modulate:a", 0.0, 1.5)
	tween.tween_callback(func(): if is_instance_valid(toast): toast.queue_free())

# ---- 瘥??湔? ----
func _process(delta: float) -> void:
	if _flash_seconds_left > 0:
		_flash_seconds_left -= int(delta)
		if _flash_seconds_left < 0:
			_flash_seconds_left = 0
		_update_flash_time_label()
