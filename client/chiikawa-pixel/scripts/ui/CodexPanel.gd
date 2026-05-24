## CodexPanel.gd ??擳????Ｘ嚗AY-081嚗?
## 憿舐內??璅??脣漲嚗圾??憿舐內??
extends Node2D

# ---- 撣豢 ----
const PANEL_WIDTH  := 420
const PANEL_HEIGHT := 520
const RARITY_COLORS := {
	"common":    Color(0.8, 0.8, 0.8),
	"rare":      Color(0.3, 0.6, 1.0),
	"epic":      Color(0.7, 0.3, 1.0),
	"legendary": Color(1.0, 0.7, 0.0),
}
const RARITY_ICONS := {
	"common":    "漎?,"
	"rare":      "?",
	"epic":      "?",
	"legendary": "潃?,"
}

# ---- 蝭暺???----
var _font: FontFile
var _bg: ColorRect
var _title_label: Label
var _progress_label: Label
var _scroll_container: ScrollContainer
var _entry_container: VBoxContainer
var _close_btn: Button
var _is_visible := false

# ---- ??鞈? ----
var _entries: Array = []
var _unlocked_count: int = 0
var _total_count: int = 0

# ---- 閫??? ----
var _unlock_queue: Array = []
var _showing_unlock := false

signal codex_closed

func setup(font: FontFile) -> void:
	_font = font
	_build_ui()
	_connect_signals()
	hide()

func _build_ui() -> void:
	# ??Ｘ
	_bg = ColorRect.new()
	_bg.color = Color(0.05, 0.08, 0.15, 0.95)
	_bg.size = Vector2(PANEL_WIDTH, PANEL_HEIGHT)
	_bg.position = Vector2(
		(1600 - PANEL_WIDTH) / 2.0,
		(900 - PANEL_HEIGHT) / 2.0
	)
	add_child(_bg)

	# 璅?
	_title_label = Label.new()
	_title_label.text = "?? ??"
	_title_label.position = _bg.position + Vector2(16, 12)
	_title_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.2))
	if _font:
		_title_label.add_theme_font_override("font", _font)
		_title_label.add_theme_font_size_override("font_size", 20)
	add_child(_title_label)

	# ?脣漲璅惜
	_progress_label = Label.new()
	_progress_label.text = "0 / 12"
	_progress_label.position = _bg.position + Vector2(PANEL_WIDTH - 80, 14)
	_progress_label.add_theme_color_override("font_color", Color(0.7, 0.9, 0.7))
	if _font:
		_progress_label.add_theme_font_override("font", _font)
		_progress_label.add_theme_font_size_override("font_size", 16)
	add_child(_progress_label)

	# ??蝺?
	var sep = ColorRect.new()
	sep.color = Color(0.3, 0.5, 0.8, 0.5)
	sep.size = Vector2(PANEL_WIDTH - 16, 2)
	sep.position = _bg.position + Vector2(8, 40)
	add_child(sep)

	# ????
	_close_btn = Button.new()
	_close_btn.text = "??"
	_close_btn.size = Vector2(28, 28)
	_close_btn.position = _bg.position + Vector2(PANEL_WIDTH - 36, 8)
	_close_btn.add_theme_color_override("font_color", Color(1, 0.4, 0.4))
	if _font:
		_close_btn.add_theme_font_override("font", _font)
	add_child(_close_btn)
	_close_btn.pressed.connect(_on_close_pressed)

	# ?脣?摰孵
	_scroll_container = ScrollContainer.new()
	_scroll_container.position = _bg.position + Vector2(8, 48)
	_scroll_container.size = Vector2(PANEL_WIDTH - 16, PANEL_HEIGHT - 60)
	add_child(_scroll_container)

	_entry_container = VBoxContainer.new()
	_entry_container.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	_scroll_container.add_child(_entry_container)

func _connect_signals() -> void:
	if GameManager.has_signal("codex_updated"):
		GameManager.codex_updated.connect(_on_codex_updated)
	if GameManager.has_signal("codex_unlocked"):
		GameManager.codex_unlocked.connect(_on_codex_unlocked)
	if GameManager.has_signal("codex_complete"):
		GameManager.codex_complete.connect(_on_codex_complete)

func _on_codex_updated(data: Dictionary) -> void:
	_entries = data.get("entries", [])
	_unlocked_count = data.get("unlocked_count", 0)
	_total_count = data.get("total_count", 12)
	_refresh_entries()

func _refresh_entries() -> void:
	# 皜????
	for child in _entry_container.get_children():
		child.queue_free()

	# ?湔?脣漲
	_progress_label.text = "%d / %d" % [_unlocked_count, _total_count]

	# 靘??漲??憿舐內
	var groups := {
		"legendary": [],
		"epic":      [],
		"rare":      [],
		"common":    [],
	}
	for entry in _entries:
		var rarity = entry.get("rarity", "common")
		if groups.has(rarity):
			groups[rarity].append(entry)

	var group_order := ["legendary", "epic", "rare", "common"]
	var group_names := {
		"legendary": "潃??唾牧",
		"epic":      "? ?脰帘",
		"rare":      "? 蝔??,"
		"common":    "漎??桅?,"
	}

	for rarity in group_order:
		var group_entries = groups[rarity]
		if group_entries.is_empty():
			continue

		# ??璅?
		var group_label = Label.new()
		group_label.text = group_names[rarity]
		group_label.add_theme_color_override("font_color", RARITY_COLORS[rarity])
		if _font:
			group_label.add_theme_font_override("font", _font)
			group_label.add_theme_font_size_override("font_size", 14)
		group_label.custom_minimum_size = Vector2(0, 24)
		_entry_container.add_child(group_label)

		# 璇?”
		for entry in group_entries:
			_add_entry_row(entry)

		# ??
		var spacer = Control.new()
		spacer.custom_minimum_size = Vector2(0, 8)
		_entry_container.add_child(spacer)

func _add_entry_row(entry: Dictionary) -> void:
	var unlocked: bool = entry.get("unlocked", false)
	var rarity: String = entry.get("rarity", "common")
	var name_text: String = entry.get("target_name", "???")
	var kill_count: int = entry.get("kill_count", 0)
	var max_mult: float = entry.get("max_multiplier", 0.0)

	var row = HBoxContainer.new()
	row.custom_minimum_size = Vector2(0, 36)

	# ?內嚗圾??蝔?漲?內嚗閫??=??
	var icon_label = Label.new()
	icon_label.text = RARITY_ICONS[rarity] if unlocked else "??"
	icon_label.custom_minimum_size = Vector2(28, 0)
	if _font:
		icon_label.add_theme_font_override("font", _font)
		icon_label.add_theme_font_size_override("font_size", 16)
	row.add_child(icon_label)

	# ?迂
	var name_label = Label.new()
	name_label.text = name_text if unlocked else "???"
	name_label.custom_minimum_size = Vector2(120, 0)
	name_label.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	var name_color = RARITY_COLORS[rarity] if unlocked else Color(0.4, 0.4, 0.4)
	name_label.add_theme_color_override("font_color", name_color)
	if _font:
		name_label.add_theme_font_override("font", _font)
		name_label.add_theme_font_size_override("font_size", 14)
	row.add_child(name_label)

	if unlocked:
		# ?甈⊥
		var kill_label = Label.new()
		kill_label.text = "?%d" % kill_count
		kill_label.custom_minimum_size = Vector2(50, 0)
		kill_label.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8))
		if _font:
			kill_label.add_theme_font_override("font", _font)
			kill_label.add_theme_font_size_override("font_size", 12)
		row.add_child(kill_label)

		# ?擃?
		var mult_label = Label.new()
		mult_label.text = "?擃?%.0fx" % max_mult
		mult_label.custom_minimum_size = Vector2(80, 0)
		mult_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.2))
		if _font:
			mult_label.add_theme_font_override("font", _font)
			mult_label.add_theme_font_size_override("font_size", 12)
		row.add_child(mult_label)
	else:
		# ?芾圾??蝷?
		var hint_label = Label.new()
		hint_label.text = "?閫??"
		hint_label.add_theme_color_override("font_color", Color(0.5, 0.5, 0.5))
		if _font:
			hint_label.add_theme_font_override("font", _font)
			hint_label.add_theme_font_size_override("font_size", 12)
		row.add_child(hint_label)

	_entry_container.add_child(row)

	# ??蝺?
	var sep = ColorRect.new()
	sep.color = Color(0.2, 0.3, 0.5, 0.3)
	sep.custom_minimum_size = Vector2(0, 1)
	_entry_container.add_child(sep)

# ---- 閫??? ----

func _on_codex_unlocked(data: Dictionary) -> void:
	_unlock_queue.append(data)
	if not _showing_unlock:
		_show_next_unlock()

func _show_next_unlock() -> void:
	if _unlock_queue.is_empty():
		_showing_unlock = false
		return
	_showing_unlock = true
	var data = _unlock_queue.pop_front()
	_spawn_unlock_popup(data)

func _spawn_unlock_popup(data: Dictionary) -> void:
	var target_name: String = data.get("target_name", "???")
	var rarity: String = data.get("rarity", "common")
	var reward: int = data.get("reward", 200)
	var unlocked: int = data.get("unlocked_count", 0)
	var total: int = data.get("total_count", 12)

	var popup = Node2D.new()
	popup.position = Vector2(800, 400)
	add_child(popup)

	# ?
	var bg = ColorRect.new()
	bg.color = Color(0.05, 0.08, 0.15, 0.92)
	bg.size = Vector2(280, 80)
	bg.position = Vector2(-140, -40)
	popup.add_child(bg)

	# ??
	var border = ColorRect.new()
	border.color = RARITY_COLORS.get(rarity, Color.WHITE)
	border.size = Vector2(280, 3)
	border.position = Vector2(-140, -40)
	popup.add_child(border)

	# 璅?
	var title_lbl = Label.new()
	title_lbl.text = "?? ??閫??嚗?"
	title_lbl.position = Vector2(-130, -34)
	title_lbl.add_theme_color_override("font_color", Color(1.0, 0.85, 0.2))
	if _font:
		title_lbl.add_theme_font_override("font", _font)
		title_lbl.add_theme_font_size_override("font_size", 14)
	popup.add_child(title_lbl)

	# ?格??迂
	var name_lbl = Label.new()
	name_lbl.text = "%s %s" % [RARITY_ICONS.get(rarity, ""), target_name]
	name_lbl.position = Vector2(-130, -14)
	name_lbl.add_theme_color_override("font_color", RARITY_COLORS.get(rarity, Color.WHITE))
	if _font:
		name_lbl.add_theme_font_override("font", _font)
		name_lbl.add_theme_font_size_override("font_size", 16)
	popup.add_child(name_lbl)

	# ? + ?脣漲
	var reward_lbl = Label.new()
	reward_lbl.text = "+%d ?馳  (%d/%d)" % [reward, unlocked, total]
	reward_lbl.position = Vector2(-130, 10)
	reward_lbl.add_theme_color_override("font_color", Color(0.7, 1.0, 0.7))
	if _font:
		reward_lbl.add_theme_font_override("font", _font)
		reward_lbl.add_theme_font_size_override("font_size", 13)
	popup.add_child(reward_lbl)

	# ?嚗楚?????? ??瘛∪
	popup.modulate.a = 0.0
	var tween = create_tween()
	tween.tween_property(popup, "modulate:a", 1.0, 0.3)
	tween.tween_interval(2.0)
	tween.tween_property(popup, "modulate:a", 0.0, 0.4)
	tween.tween_callback(func():
		popup.queue_free()
		_show_next_unlock()
	)

func _on_codex_complete(data: Dictionary) -> void:
	var reward: int = data.get("reward", 5000)
	var title_name: String = data.get("title_name", "??摰???)"

	# ?典????之敶?
	var popup = Node2D.new()
	popup.position = Vector2(800, 450)
	add_child(popup)

	var bg = ColorRect.new()
	bg.color = Color(0.05, 0.05, 0.1, 0.95)
	bg.size = Vector2(360, 120)
	bg.position = Vector2(-180, -60)
	popup.add_child(bg)

	var border = ColorRect.new()
	border.color = Color(1.0, 0.7, 0.0)
	border.size = Vector2(360, 4)
	border.position = Vector2(-180, -60)
	popup.add_child(border)

	var title_lbl = Label.new()
	title_lbl.text = "?? ???冽????"
	title_lbl.position = Vector2(-170, -52)
	title_lbl.add_theme_color_override("font_color", Color(1.0, 0.85, 0.2))
	if _font:
		title_lbl.add_theme_font_override("font", _font)
		title_lbl.add_theme_font_size_override("font_size", 18)
	popup.add_child(title_lbl)

	var reward_lbl = Label.new()
	reward_lbl.text = "+%d ?馳" % reward
	reward_lbl.position = Vector2(-170, -20)
	reward_lbl.add_theme_color_override("font_color", Color(0.7, 1.0, 0.7))
	if _font:
		reward_lbl.add_theme_font_override("font", _font)
		reward_lbl.add_theme_font_size_override("font_size", 20)
	popup.add_child(reward_lbl)

	var title_unlock_lbl = Label.new()
	title_unlock_lbl.text = "蝔梯?閫??嚗?s?? % title_name"
	title_unlock_lbl.position = Vector2(-170, 14)
	title_unlock_lbl.add_theme_color_override("font_color", Color(1.0, 0.7, 0.0))
	if _font:
		title_unlock_lbl.add_theme_font_override("font", _font)
		title_unlock_lbl.add_theme_font_size_override("font_size", 14)
	popup.add_child(title_unlock_lbl)

	popup.modulate.a = 0.0
	var tween = create_tween()
	tween.tween_property(popup, "modulate:a", 1.0, 0.5)
	tween.tween_interval(4.0)
	tween.tween_property(popup, "modulate:a", 0.0, 0.5)
	tween.tween_callback(popup.queue_free)

# ---- 憿舐內/?梯? ----

func toggle() -> void:
	if _is_visible:
		_hide_panel()
	else:
		_show_panel()

func _show_panel() -> void:
	_is_visible = true
	show()
	# 隢???啣?????
	if GameManager.has_method("request_codex"):
		GameManager.request_codex()

func _hide_panel() -> void:
	_is_visible = false
	hide()
	emit_signal("codex_closed")

func _on_close_pressed() -> void:
	_hide_panel()
