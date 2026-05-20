## VIPPanel.gd — VIP 等級系統面板（DAY-078）
## 顯示 VIP 等級、消費進度、週獎勵領取按鈕
## 位置：TopBar 上（可折疊）
extends Node2D

# ---- 常數 ----
const PANEL_WIDTH  := 300
const PANEL_HEIGHT := 220
const BTN_SIZE     := 26

# VIP 等級顏色
const TIER_COLORS := {
	0: Color(0.6, 0.6, 0.6),       # 無 VIP（灰）
	1: Color(0.80, 0.50, 0.20),    # 青銅
	2: Color(0.75, 0.75, 0.75),    # 白銀
	3: Color(1.00, 0.85, 0.10),    # 黃金
	4: Color(0.90, 0.90, 0.88),    # 白金
	5: Color(0.73, 0.95, 1.00),    # 鑽石
}

# ---- 節點引用 ----
var _pixel_font: Font = null
var _is_open: bool = false
var _toggle_btn: Button = null
var _panel_bg: ColorRect = null
var _tier_label: Label = null
var _spend_label: Label = null
var _progress_bar: ColorRect = null
var _progress_fill: ColorRect = null
var _progress_label: Label = null
var _cashback_label: Label = null
var _daily_mult_label: Label = null
var _weekly_btn: Button = null
var _weekly_label: Label = null
var _tier_rows: Array = []

# ---- VIP 資料 ----
var _vip_data: Dictionary = {
	"vip_level": 0,
	"tier_name": "一般玩家",
	"tier_icon": "👤",
	"tier_color": "#999999",
	"total_spend": 0,
	"cashback_rate": 0.0,
	"daily_bonus_mult": 1.0,
	"weekly_bonus": 0,
	"next_level": 1,
	"spend_to_next": 10000,
	"progress": 0.0,
	"can_claim_weekly": false
}

# ---- 初始化 ----
func _ready() -> void:
	if ResourceLoader.exists("res://assets/fonts/pixel8.fnt"):
		_pixel_font = load("res://assets/fonts/pixel8.fnt")
	_build_toggle_btn()
	_build_panel()
	_connect_signals()

func setup(font: Font) -> void:
	if font:
		_pixel_font = font

## 建立折疊按鈕（TopBar 上）
func _build_toggle_btn() -> void:
	_toggle_btn = Button.new()
	_toggle_btn.text = "💎"
	_toggle_btn.size = Vector2(32, 24)
	_toggle_btn.position = Vector2(0, 0)
	_toggle_btn.flat = true
	_toggle_btn.tooltip_text = "VIP 等級"
	if _pixel_font:
		_toggle_btn.add_theme_font_override("font", _pixel_font)
		_toggle_btn.add_theme_font_size_override("font_size", 14)
	add_child(_toggle_btn)

## 建立主面板（預設隱藏）
func _build_panel() -> void:
	_panel_bg = ColorRect.new()
	_panel_bg.position = Vector2(-PANEL_WIDTH + 32, 28)
	_panel_bg.size = Vector2(PANEL_WIDTH, PANEL_HEIGHT)
	_panel_bg.color = Color(0.03, 0.05, 0.15, 0.93)
	_panel_bg.visible = false
	add_child(_panel_bg)

	# 標題
	var title := Label.new()
	title.position = Vector2(8, 4)
	title.text = "💎 VIP 會員系統"
	title.add_theme_color_override("font_color", Color(0.73, 0.95, 1.0))
	if _pixel_font:
		title.add_theme_font_override("font", _pixel_font)
		title.add_theme_font_size_override("font_size", 12)
	_panel_bg.add_child(title)

	# 等級標籤
	_tier_label = Label.new()
	_tier_label.position = Vector2(8, 22)
	_tier_label.text = "👤 一般玩家"
	_tier_label.add_theme_color_override("font_color", Color(0.6, 0.6, 0.6))
	if _pixel_font:
		_tier_label.add_theme_font_override("font", _pixel_font)
		_tier_label.add_theme_font_size_override("font_size", 11)
	_panel_bg.add_child(_tier_label)

	# 消費標籤
	_spend_label = Label.new()
	_spend_label.position = Vector2(8, 38)
	_spend_label.text = "累積消費：0 金幣"
	_spend_label.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8))
	if _pixel_font:
		_spend_label.add_theme_font_override("font", _pixel_font)
		_spend_label.add_theme_font_size_override("font_size", 10)
	_panel_bg.add_child(_spend_label)

	# 進度條背景
	_progress_bar = ColorRect.new()
	_progress_bar.position = Vector2(8, 54)
	_progress_bar.size = Vector2(PANEL_WIDTH - 16, 10)
	_progress_bar.color = Color(0.15, 0.15, 0.25)
	_panel_bg.add_child(_progress_bar)

	# 進度條填充
	_progress_fill = ColorRect.new()
	_progress_fill.position = Vector2(0, 0)
	_progress_fill.size = Vector2(0, 10)
	_progress_fill.color = Color(0.73, 0.95, 1.0)
	_progress_bar.add_child(_progress_fill)

	# 進度標籤
	_progress_label = Label.new()
	_progress_label.position = Vector2(8, 66)
	_progress_label.text = "距下一等級：10,000 金幣"
	_progress_label.add_theme_color_override("font_color", Color(0.7, 0.7, 0.9))
	if _pixel_font:
		_progress_label.add_theme_font_override("font", _pixel_font)
		_progress_label.add_theme_font_size_override("font_size", 9)
	_panel_bg.add_child(_progress_label)

	# 返還率標籤
	_cashback_label = Label.new()
	_cashback_label.position = Vector2(8, 82)
	_cashback_label.text = "💰 金幣返還：0%"
	_cashback_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.2))
	if _pixel_font:
		_cashback_label.add_theme_font_override("font", _pixel_font)
		_cashback_label.add_theme_font_size_override("font_size", 10)
	_panel_bg.add_child(_cashback_label)

	# 每日獎勵倍率標籤
	_daily_mult_label = Label.new()
	_daily_mult_label.position = Vector2(8, 96)
	_daily_mult_label.text = "📅 每日獎勵倍率：×1.0"
	_daily_mult_label.add_theme_color_override("font_color", Color(0.6, 1.0, 0.6))
	if _pixel_font:
		_daily_mult_label.add_theme_font_override("font", _pixel_font)
		_daily_mult_label.add_theme_font_size_override("font_size", 10)
	_panel_bg.add_child(_daily_mult_label)

	# 週獎勵標籤
	_weekly_label = Label.new()
	_weekly_label.position = Vector2(8, 112)
	_weekly_label.text = "🎁 週獎勵：--"
	_weekly_label.add_theme_color_override("font_color", Color(0.9, 0.7, 1.0))
	if _pixel_font:
		_weekly_label.add_theme_font_override("font", _pixel_font)
		_weekly_label.add_theme_font_size_override("font_size", 10)
	_panel_bg.add_child(_weekly_label)

	# 週獎勵領取按鈕
	_weekly_btn = Button.new()
	_weekly_btn.position = Vector2(8, 128)
	_weekly_btn.size = Vector2(PANEL_WIDTH - 16, 24)
	_weekly_btn.text = "領取週獎勵"
	_weekly_btn.disabled = true
	if _pixel_font:
		_weekly_btn.add_theme_font_override("font", _pixel_font)
		_weekly_btn.add_theme_font_size_override("font_size", 10)
	_panel_bg.add_child(_weekly_btn)

	# VIP 等級列表（5 個等級）
	var tier_defs := [
		{"level": 1, "name": "青銅", "icon": "🥉", "spend": "10,000"},
		{"level": 2, "name": "白銀", "icon": "🥈", "spend": "50,000"},
		{"level": 3, "name": "黃金", "icon": "🥇", "spend": "200,000"},
		{"level": 4, "name": "白金", "icon": "💎", "spend": "500,000"},
		{"level": 5, "name": "鑽石", "icon": "👑", "spend": "2,000,000"},
	]
	for i in range(tier_defs.size()):
		var td = tier_defs[i]
		var row := Label.new()
		row.position = Vector2(8, 158 + i * 12)
		row.text = "%s %s  (%s)" % [td["icon"], td["name"], td["spend"]]
		row.add_theme_color_override("font_color", Color(0.5, 0.5, 0.5))
		if _pixel_font:
			row.add_theme_font_override("font", _pixel_font)
			row.add_theme_font_size_override("font_size", 9)
		_panel_bg.add_child(row)
		_tier_rows.append(row)

## 連接訊號
func _connect_signals() -> void:
	if _toggle_btn:
		_toggle_btn.pressed.connect(_on_toggle_pressed)
	if _weekly_btn:
		_weekly_btn.pressed.connect(_on_weekly_btn_pressed)
	# 連接 GameManager 訊號
	if GameManager:
		if GameManager.has_signal("vip_updated"):
			GameManager.vip_updated.connect(_on_vip_updated)
		if GameManager.has_signal("vip_level_up"):
			GameManager.vip_level_up.connect(_on_vip_level_up)
		if GameManager.has_signal("vip_weekly_claimed"):
			GameManager.vip_weekly_claimed.connect(_on_vip_weekly_claimed)

## 折疊/展開面板
func _on_toggle_pressed() -> void:
	_is_open = !_is_open
	if _panel_bg:
		_panel_bg.visible = _is_open

## 更新 VIP 資料
func _on_vip_updated(data: Dictionary) -> void:
	_vip_data = data
	_refresh_ui()

## VIP 升級通知
func _on_vip_level_up(data: Dictionary) -> void:
	_show_level_up_popup(data)

## 週獎勵領取通知
func _on_vip_weekly_claimed(data: Dictionary) -> void:
	_show_weekly_claimed_popup(data)

## 領取週獎勵按鈕
func _on_weekly_btn_pressed() -> void:
	if GameManager:
		GameManager.claim_vip_weekly()

## 刷新 UI
func _refresh_ui() -> void:
	if not _panel_bg:
		return

	var level: int = _vip_data.get("vip_level", 0)
	var tier_name: String = _vip_data.get("tier_name", "一般玩家")
	var tier_icon: String = _vip_data.get("tier_icon", "👤")
	var total_spend: int = _vip_data.get("total_spend", 0)
	var cashback_rate: float = _vip_data.get("cashback_rate", 0.0)
	var daily_mult: float = _vip_data.get("daily_bonus_mult", 1.0)
	var weekly_bonus: int = _vip_data.get("weekly_bonus", 0)
	var spend_to_next: int = _vip_data.get("spend_to_next", 10000)
	var progress: float = _vip_data.get("progress", 0.0)
	var can_claim: bool = _vip_data.get("can_claim_weekly", false)

	# 等級顏色
	var tier_color: Color = TIER_COLORS.get(level, Color(0.6, 0.6, 0.6))

	# 更新按鈕圖示顏色
	if _toggle_btn:
		_toggle_btn.modulate = tier_color

	# 等級標籤
	if _tier_label:
		_tier_label.text = "%s %s" % [tier_icon, tier_name]
		_tier_label.add_theme_color_override("font_color", tier_color)

	# 消費標籤
	if _spend_label:
		_spend_label.text = "累積消費：%s 金幣" % _format_number(total_spend)

	# 進度條
	if _progress_fill and _progress_bar:
		var bar_width := _progress_bar.size.x
		_progress_fill.size.x = bar_width * clamp(progress, 0.0, 1.0)
		_progress_fill.color = tier_color

	# 進度標籤
	if _progress_label:
		if level >= 5:
			_progress_label.text = "✨ 已達最高等級"
			_progress_label.add_theme_color_override("font_color", TIER_COLORS[5])
		else:
			_progress_label.text = "距下一等級：%s 金幣" % _format_number(spend_to_next)

	# 返還率
	if _cashback_label:
		if cashback_rate > 0:
			_cashback_label.text = "💰 金幣返還：%.0f%%" % (cashback_rate * 100)
			_cashback_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.2))
		else:
			_cashback_label.text = "💰 金幣返還：無"
			_cashback_label.add_theme_color_override("font_color", Color(0.5, 0.5, 0.5))

	# 每日獎勵倍率
	if _daily_mult_label:
		_daily_mult_label.text = "📅 每日獎勵倍率：×%.1f" % daily_mult
		if daily_mult > 1.0:
			_daily_mult_label.add_theme_color_override("font_color", Color(0.6, 1.0, 0.6))
		else:
			_daily_mult_label.add_theme_color_override("font_color", Color(0.5, 0.5, 0.5))

	# 週獎勵
	if _weekly_label:
		if weekly_bonus > 0:
			_weekly_label.text = "🎁 週獎勵：%s 金幣" % _format_number(weekly_bonus)
			_weekly_label.add_theme_color_override("font_color", Color(0.9, 0.7, 1.0))
		else:
			_weekly_label.text = "🎁 週獎勵：升級 VIP 解鎖"
			_weekly_label.add_theme_color_override("font_color", Color(0.5, 0.5, 0.5))

	# 週獎勵按鈕
	if _weekly_btn:
		_weekly_btn.disabled = not can_claim
		if can_claim:
			_weekly_btn.text = "🎁 領取週獎勵 (%s 金幣)" % _format_number(weekly_bonus)
			_weekly_btn.modulate = Color(1.0, 1.0, 1.0)
		else:
			_weekly_btn.text = "週獎勵（7天後可領）"
			_weekly_btn.modulate = Color(0.6, 0.6, 0.6)

	# 更新等級列表顏色
	for i in range(_tier_rows.size()):
		var row: Label = _tier_rows[i]
		var row_level := i + 1
		if row_level <= level:
			row.add_theme_color_override("font_color", TIER_COLORS.get(row_level, Color(0.6, 0.6, 0.6)))
		elif row_level == level + 1:
			row.add_theme_color_override("font_color", Color(0.9, 0.9, 0.9))
		else:
			row.add_theme_color_override("font_color", Color(0.4, 0.4, 0.4))

## 顯示 VIP 升級彈窗
func _show_level_up_popup(data: Dictionary) -> void:
	var new_level: int = data.get("new_level", 1)
	var tier_name: String = data.get("tier_name", "")
	var tier_icon: String = data.get("tier_icon", "💎")
	var tier_color_hex: String = data.get("tier_color", "#FFFFFF")
	var weekly_bonus: int = data.get("weekly_bonus", 0)

	var tier_color := Color(TIER_COLORS.get(new_level, Color(1.0, 1.0, 1.0)))

	# 建立彈窗
	var canvas := get_viewport().get_canvas_layer_node(1) if get_viewport() else null
	var popup_parent := canvas if canvas else get_parent()
	if not is_instance_valid(popup_parent):
		return

	var popup := ColorRect.new()
	popup.size = Vector2(280, 80)
	popup.position = Vector2(
		(get_viewport().get_visible_rect().size.x - 280) / 2.0,
		get_viewport().get_visible_rect().size.y * 0.35
	)
	popup.color = Color(0.03, 0.05, 0.15, 0.95)
	popup_parent.add_child(popup)

	var lbl := Label.new()
	lbl.position = Vector2(8, 8)
	lbl.text = "%s VIP 升級！" % tier_icon
	lbl.add_theme_color_override("font_color", tier_color)
	if _pixel_font:
		lbl.add_theme_font_override("font", _pixel_font)
		lbl.add_theme_font_size_override("font_size", 14)
	popup.add_child(lbl)

	var lbl2 := Label.new()
	lbl2.position = Vector2(8, 28)
	lbl2.text = "恭喜成為 %s！" % tier_name
	lbl2.add_theme_color_override("font_color", Color(1.0, 1.0, 1.0))
	if _pixel_font:
		lbl2.add_theme_font_override("font", _pixel_font)
		lbl2.add_theme_font_size_override("font_size", 11)
	popup.add_child(lbl2)

	var lbl3 := Label.new()
	lbl3.position = Vector2(8, 46)
	lbl3.text = "週獎勵：%s 金幣 | 立即可領取！" % _format_number(weekly_bonus)
	lbl3.add_theme_color_override("font_color", Color(0.9, 0.7, 1.0))
	if _pixel_font:
		lbl3.add_theme_font_override("font", _pixel_font)
		lbl3.add_theme_font_size_override("font_size", 10)
	popup.add_child(lbl3)

	# 動畫：彈入 → 停留 → 淡出
	var tween := popup.create_tween()
	popup.modulate.a = 0.0
	tween.tween_property(popup, "modulate:a", 1.0, 0.3)
	tween.tween_interval(3.0)
	tween.tween_property(popup, "modulate:a", 0.0, 0.5)
	tween.tween_callback(popup.queue_free)

## 顯示週獎勵領取彈窗
func _show_weekly_claimed_popup(data: Dictionary) -> void:
	var coins: int = data.get("coins", 0)
	var tier_name: String = data.get("tier_name", "")

	var canvas := get_viewport().get_canvas_layer_node(1) if get_viewport() else null
	var popup_parent := canvas if canvas else get_parent()
	if not is_instance_valid(popup_parent):
		return

	var popup := ColorRect.new()
	popup.size = Vector2(240, 60)
	popup.position = Vector2(
		(get_viewport().get_visible_rect().size.x - 240) / 2.0,
		get_viewport().get_visible_rect().size.y * 0.4
	)
	popup.color = Color(0.05, 0.15, 0.05, 0.95)
	popup_parent.add_child(popup)

	var lbl := Label.new()
	lbl.position = Vector2(8, 8)
	lbl.text = "🎁 VIP 週獎勵領取成功！"
	lbl.add_theme_color_override("font_color", Color(0.9, 0.7, 1.0))
	if _pixel_font:
		lbl.add_theme_font_override("font", _pixel_font)
		lbl.add_theme_font_size_override("font_size", 11)
	popup.add_child(lbl)

	var lbl2 := Label.new()
	lbl2.position = Vector2(8, 28)
	lbl2.text = "+%s 金幣（%s 福利）" % [_format_number(coins), tier_name]
	lbl2.add_theme_color_override("font_color", Color(1.0, 0.85, 0.2))
	if _pixel_font:
		lbl2.add_theme_font_override("font", _pixel_font)
		lbl2.add_theme_font_size_override("font_size", 10)
	popup.add_child(lbl2)

	var tween := popup.create_tween()
	popup.modulate.a = 0.0
	tween.tween_property(popup, "modulate:a", 1.0, 0.3)
	tween.tween_interval(2.5)
	tween.tween_property(popup, "modulate:a", 0.0, 0.4)
	tween.tween_callback(popup.queue_free)

## 格式化數字（加千分位）
func _format_number(n: int) -> String:
	var s := str(n)
	var result := ""
	var count := 0
	for i in range(s.length() - 1, -1, -1):
		if count > 0 and count % 3 == 0:
			result = "," + result
		result = s[i] + result
		count += 1
	return result
