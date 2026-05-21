## UnluckyBonusPanel.gd — 失敗補償系統面板（DAY-135）
## 業界依據：Funrize 2026 的「Unlucky Bonus」
## 連續花費超過一定金額但獲得低回報時，自動給予補償獎勵
## 防止玩家因為「運氣太差」而離開，是 2026 年業界最新的留存機制
extends Node2D

# ---- 常數 ----
const PANEL_WIDTH  := 200
const PANEL_HEIGHT := 50

# ---- 狀態 ----
var _pixel_font: Font = null
var _bonus_queue: Array = []  # 待顯示的補償通知佇列
var _is_showing: bool = false

# ---- 初始化 ----
func _ready() -> void:
	if ResourceLoader.exists("res://assets/fonts/pixel8.fnt"):
		_pixel_font = load("res://assets/fonts/pixel8.fnt")
	_connect_signals()

func setup(font: Font) -> void:
	if font:
		_pixel_font = font

func _connect_signals() -> void:
	if GameManager.has_signal("unlucky_bonus"):
		GameManager.unlucky_bonus.connect(_on_unlucky_bonus)

# ---- 事件處理 ----

func _on_unlucky_bonus(data: Dictionary) -> void:
	var bonus_amount = data.get("bonus_amount", 0)
	var message = data.get("message", "🍀 運氣補償！")
	_bonus_queue.append({"amount": bonus_amount, "message": message})
	if not _is_showing:
		_show_next_bonus()

func _show_next_bonus() -> void:
	if _bonus_queue.is_empty():
		_is_showing = false
		return

	_is_showing = true
	var bonus_data = _bonus_queue.pop_front()
	_show_bonus_popup(bonus_data["amount"], bonus_data["message"])

func _show_bonus_popup(amount: int, message: String) -> void:
	# 建立補償通知彈窗（螢幕中央偏下）
	var popup := Node2D.new()
	add_child(popup)

	# 背景（綠色主題）
	var bg := ColorRect.new()
	bg.size = Vector2(PANEL_WIDTH, PANEL_HEIGHT)
	bg.position = Vector2(-PANEL_WIDTH / 2.0, 0)
	bg.color = Color(0.05, 0.25, 0.08, 0.92)
	popup.add_child(bg)

	# 邊框（亮綠色）
	var border := ColorRect.new()
	border.size = Vector2(PANEL_WIDTH + 2, PANEL_HEIGHT + 2)
	border.position = Vector2(-PANEL_WIDTH / 2.0 - 1, -1)
	border.color = Color(0.2, 0.9, 0.3, 0.8)
	border.z_index = -1
	popup.add_child(border)

	# 四葉草圖示（大）
	var clover_lbl := Label.new()
	clover_lbl.text = "🍀"
	clover_lbl.position = Vector2(-PANEL_WIDTH / 2.0 + 6, 4)
	clover_lbl.add_theme_font_size_override("font_size", 28)
	popup.add_child(clover_lbl)

	# 補償文字
	var title_lbl := Label.new()
	title_lbl.text = "運氣補償！"
	title_lbl.position = Vector2(-PANEL_WIDTH / 2.0 + 44, 4)
	title_lbl.add_theme_color_override("font_color", Color(0.4, 1.0, 0.5))
	title_lbl.add_theme_font_size_override("font_size", 13)
	if _pixel_font:
		title_lbl.add_theme_font_override("font", _pixel_font)
	popup.add_child(title_lbl)

	# 金額
	var amount_lbl := Label.new()
	amount_lbl.text = "+%d 🪙" % amount
	amount_lbl.position = Vector2(-PANEL_WIDTH / 2.0 + 44, 24)
	amount_lbl.add_theme_color_override("font_color", Color(1.0, 0.95, 0.3))
	amount_lbl.add_theme_font_size_override("font_size", 16)
	if _pixel_font:
		amount_lbl.add_theme_font_override("font", _pixel_font)
	popup.add_child(amount_lbl)

	# 動畫：從下方滑入 → 停留 → 淡出
	popup.position = Vector2(640, 500)  # 螢幕中央偏下
	popup.modulate.a = 0.0

	var tween = popup.create_tween()
	# 滑入 + 淡入
	tween.set_parallel(true)
	tween.tween_property(popup, "position:y", 460.0, 0.3).set_ease(Tween.EASE_OUT)
	tween.tween_property(popup, "modulate:a", 1.0, 0.3)
	tween.set_parallel(false)

	# 停留 2.5 秒
	tween.tween_interval(2.5)

	# 縮放彈跳（強調感）
	tween.tween_property(popup, "scale", Vector2(1.1, 1.1), 0.1)
	tween.tween_property(popup, "scale", Vector2(1.0, 1.0), 0.1)

	# 淡出
	tween.tween_property(popup, "modulate:a", 0.0, 0.4)
	tween.tween_callback(func():
		if is_instance_valid(popup):
			popup.queue_free()
		_show_next_bonus()
	)

	# 全螢幕綠色閃光
	_show_lucky_flash()

func _show_lucky_flash() -> void:
	# 全螢幕淡綠色閃光
	var flash := ColorRect.new()
	flash.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	flash.color = Color(0.1, 0.8, 0.2, 0.0)
	flash.z_index = 200
	add_child(flash)

	var tween = flash.create_tween()
	tween.tween_property(flash, "color:a", 0.25, 0.1)
	tween.tween_property(flash, "color:a", 0.0, 0.3)
	tween.tween_callback(func():
		if is_instance_valid(flash):
			flash.queue_free()
	)
