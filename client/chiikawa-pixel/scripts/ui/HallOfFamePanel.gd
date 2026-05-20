## HallOfFamePanel.gd — 全服名人堂面板（DAY-110）
## 展示全服最佳記錄，激勵玩家挑戰極限
extends CanvasLayer

# 記錄類型標籤（對應 Server 端 RecordType）
const RECORD_LABELS = {
	"best_streak": "最高連擊王",
	"best_multiplier": "最高倍率王",
	"best_bonus_reward": "Bonus 大師",
	"most_jackpots": "Jackpot 收集者",
	"grand_jackpot": "Grand Jackpot 傳說",
	"boss_kills": "BOSS 獵人",
	"max_coins": "金幣大亨",
	"best_rtp": "效率之王"
}

const RECORD_ICONS = {
	"best_streak": "🔥",
	"best_multiplier": "⚡",
	"best_bonus_reward": "🌾",
	"most_jackpots": "🎰",
	"grand_jackpot": "👑",
	"boss_kills": "⚔️",
	"max_coins": "💰",
	"best_rtp": "📊"
}

# 記錄顯示順序
const RECORD_ORDER = [
	"grand_jackpot", "best_multiplier", "best_streak",
	"max_coins", "boss_kills", "best_bonus_reward",
	"most_jackpots", "best_rtp"
]

var _panel: PanelContainer
var _records_container: VBoxContainer
var _new_record_overlay: Control
var _record_entries: Dictionary = {}

func _ready():
	layer = 88
	_build_ui()
	hide()

func _build_ui():
	# 半透明背景
	var bg = ColorRect.new()
	bg.color = Color(0, 0, 0, 0.6)
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	add_child(bg)
	bg.gui_input.connect(func(e): if e is InputEventMouseButton and e.pressed: hide())

	# 主面板
	_panel = PanelContainer.new()
	_panel.set_anchors_preset(Control.PRESET_CENTER)
	_panel.custom_minimum_size = Vector2(520, 580)
	_panel.offset_left = -260
	_panel.offset_top = -290
	_panel.offset_right = 260
	_panel.offset_bottom = 290
	add_child(_panel)

	var vbox = VBoxContainer.new()
	vbox.add_theme_constant_override("separation", 8)
	_panel.add_child(vbox)

	# 標題列
	var title_row = HBoxContainer.new()
	vbox.add_child(title_row)

	var title_lbl = Label.new()
	title_lbl.text = "🏆 全服名人堂"
	title_lbl.add_theme_font_size_override("font_size", 20)
	title_lbl.add_theme_color_override("font_color", Color(1.0, 0.85, 0.2))
	title_lbl.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	title_row.add_child(title_lbl)

	var close_btn = Button.new()
	close_btn.text = "✕"
	close_btn.custom_minimum_size = Vector2(32, 32)
	close_btn.pressed.connect(hide)
	title_row.add_child(close_btn)

	# 副標題
	var subtitle = Label.new()
	subtitle.text = "挑戰全服最佳記錄，留名青史！"
	subtitle.add_theme_font_size_override("font_size", 12)
	subtitle.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8))
	vbox.add_child(subtitle)

	var sep = HSeparator.new()
	vbox.add_child(sep)

	# 記錄列表（ScrollContainer）
	var scroll = ScrollContainer.new()
	scroll.size_flags_vertical = Control.SIZE_EXPAND_FILL
	scroll.custom_minimum_size = Vector2(0, 400)
	vbox.add_child(scroll)

	_records_container = VBoxContainer.new()
	_records_container.add_theme_constant_override("separation", 6)
	_records_container.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	scroll.add_child(_records_container)

	# 初始化所有記錄行
	for rt in RECORD_ORDER:
		var entry = _create_record_entry(rt)
		_records_container.add_child(entry)
		_record_entries[rt] = entry

	# 底部按鈕
	var refresh_btn = Button.new()
	refresh_btn.text = "🔄 重新整理"
	refresh_btn.pressed.connect(_on_refresh_pressed)
	vbox.add_child(refresh_btn)

	# 新記錄通知 Overlay（全畫面）
	_new_record_overlay = _create_new_record_overlay()
	add_child(_new_record_overlay)
	_new_record_overlay.hide()

func _create_record_entry(record_type: String) -> PanelContainer:
	var panel = PanelContainer.new()
	panel.name = "entry_" + record_type

	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.15, 0.15, 0.2, 0.9)
	style.corner_radius_top_left = 6
	style.corner_radius_top_right = 6
	style.corner_radius_bottom_left = 6
	style.corner_radius_bottom_right = 6
	panel.add_theme_stylebox_override("panel", style)

	var hbox = HBoxContainer.new()
	hbox.add_theme_constant_override("separation", 8)
	panel.add_child(hbox)

	# 圖示
	var icon_lbl = Label.new()
	icon_lbl.name = "icon"
	icon_lbl.text = RECORD_ICONS.get(record_type, "🏆")
	icon_lbl.add_theme_font_size_override("font_size", 24)
	icon_lbl.custom_minimum_size = Vector2(36, 36)
	hbox.add_child(icon_lbl)

	# 資訊區
	var info_vbox = VBoxContainer.new()
	info_vbox.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	hbox.add_child(info_vbox)

	var type_lbl = Label.new()
	type_lbl.name = "type_label"
	type_lbl.text = RECORD_LABELS.get(record_type, record_type)
	type_lbl.add_theme_font_size_override("font_size", 13)
	type_lbl.add_theme_color_override("font_color", Color(1.0, 0.85, 0.2))
	info_vbox.add_child(type_lbl)

	var holder_lbl = Label.new()
	holder_lbl.name = "holder"
	holder_lbl.text = "— 尚無記錄 —"
	holder_lbl.add_theme_font_size_override("font_size", 11)
	holder_lbl.add_theme_color_override("font_color", Color(0.6, 0.6, 0.6))
	info_vbox.add_child(holder_lbl)

	var desc_lbl = Label.new()
	desc_lbl.name = "description"
	desc_lbl.text = ""
	desc_lbl.add_theme_font_size_override("font_size", 10)
	desc_lbl.add_theme_color_override("font_color", Color(0.7, 0.7, 0.7))
	info_vbox.add_child(desc_lbl)

	# 數值
	var value_lbl = Label.new()
	value_lbl.name = "value"
	value_lbl.text = ""
	value_lbl.add_theme_font_size_override("font_size", 16)
	value_lbl.add_theme_color_override("font_color", Color(0.3, 1.0, 0.5))
	value_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_RIGHT
	value_lbl.custom_minimum_size = Vector2(80, 0)
	hbox.add_child(value_lbl)

	return panel

func _create_new_record_overlay() -> Control:
	var overlay = Control.new()
	overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)

	var bg = ColorRect.new()
	bg.color = Color(0, 0, 0, 0.7)
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	overlay.add_child(bg)

	var center = CenterContainer.new()
	center.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	overlay.add_child(center)

	var vbox = VBoxContainer.new()
	vbox.add_theme_constant_override("separation", 12)
	center.add_child(vbox)

	var crown = Label.new()
	crown.name = "crown"
	crown.text = "👑"
	crown.add_theme_font_size_override("font_size", 64)
	crown.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(crown)

	var title = Label.new()
	title.name = "title"
	title.text = "🏆 新記錄誕生！"
	title.add_theme_font_size_override("font_size", 28)
	title.add_theme_color_override("font_color", Color(1.0, 0.85, 0.2))
	title.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title)

	var record_type_lbl = Label.new()
	record_type_lbl.name = "record_type"
	record_type_lbl.text = ""
	record_type_lbl.add_theme_font_size_override("font_size", 18)
	record_type_lbl.add_theme_color_override("font_color", Color(0.9, 0.9, 0.9))
	record_type_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(record_type_lbl)

	var holder_lbl = Label.new()
	holder_lbl.name = "holder"
	holder_lbl.text = ""
	holder_lbl.add_theme_font_size_override("font_size", 22)
	holder_lbl.add_theme_color_override("font_color", Color(0.3, 1.0, 0.5))
	holder_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(holder_lbl)

	var desc_lbl = Label.new()
	desc_lbl.name = "description"
	desc_lbl.text = ""
	desc_lbl.add_theme_font_size_override("font_size", 14)
	desc_lbl.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8))
	desc_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(desc_lbl)

	return overlay

# ---- 公開方法 ----

func show_panel():
	show()
	GameManager.request_hall_of_fame()

func update_records(data: Dictionary):
	var records = data.get("records", [])
	# 先清空所有記錄
	for rt in _record_entries:
		_set_entry_empty(_record_entries[rt])

	# 填入資料
	for entry_data in records:
		var rt = entry_data.get("record_type", "")
		if rt in _record_entries:
			_update_entry(_record_entries[rt], entry_data)

func show_new_record(data: Dictionary):
	var entry = data.get("entry", {})
	var rt = entry.get("record_type", "")
	var holder = entry.get("display_name", "")
	var desc = entry.get("description", "")
	var record_label = entry.get("record_label", RECORD_LABELS.get(rt, rt))
	var record_icon = entry.get("record_icon", "🏆")

	# 更新對應記錄行
	if rt in _record_entries:
		_update_entry(_record_entries[rt], entry)

	# 顯示全畫面通知（3秒後自動消失）
	var overlay = _new_record_overlay
	overlay.get_node("CenterContainer/VBoxContainer/record_type").text = record_icon + " " + record_label
	overlay.get_node("CenterContainer/VBoxContainer/holder").text = holder + " 創下新記錄！"
	overlay.get_node("CenterContainer/VBoxContainer/description").text = desc
	overlay.show()

	# 金色閃光動畫
	var tween = create_tween()
	tween.tween_property(overlay, "modulate:a", 1.0, 0.3)
	tween.tween_interval(2.5)
	tween.tween_property(overlay, "modulate:a", 0.0, 0.5)
	tween.tween_callback(overlay.hide)

# ---- 私有方法 ----

func _set_entry_empty(panel: PanelContainer):
	panel.get_node("HBoxContainer/VBoxContainer/holder").text = "— 尚無記錄 —"
	panel.get_node("HBoxContainer/VBoxContainer/holder").add_theme_color_override("font_color", Color(0.6, 0.6, 0.6))
	panel.get_node("HBoxContainer/VBoxContainer/description").text = ""
	panel.get_node("HBoxContainer/value").text = ""

func _update_entry(panel: PanelContainer, data: Dictionary):
	var holder = data.get("display_name", "")
	var desc = data.get("description", "")
	var value = data.get("value", 0.0)
	var rt = data.get("record_type", "")

	panel.get_node("HBoxContainer/VBoxContainer/holder").text = "👤 " + holder
	panel.get_node("HBoxContainer/VBoxContainer/holder").add_theme_color_override("font_color", Color(0.3, 1.0, 0.5))
	panel.get_node("HBoxContainer/VBoxContainer/description").text = desc

	# 格式化數值
	var value_str = ""
	match rt:
		"best_multiplier":
			value_str = "%.0fx" % value
		"max_coins", "best_bonus_reward", "grand_jackpot":
			value_str = "%d" % int(value)
		"best_rtp":
			value_str = "%.0f%%" % (value * 100)
		_:
			value_str = "%d" % int(value)

	panel.get_node("HBoxContainer/value").text = value_str

func _on_refresh_pressed():
	GameManager.request_hall_of_fame()
