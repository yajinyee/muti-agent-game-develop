## BuyBonusPanel.gd - DAY-114
## Buy Bonus 系統 UI：玩家可以花費金幣直接觸發 Bonus
## 參考 BGaming Fishing Club 2（2026-04）的 Buy Bonus 機制
## 標準 Bonus：BetCost × 100 / TNT Bonus：BetCost × 150（倍率加成 1.5x）
extends Node2D

const PANEL_W := 360
const PANEL_H := 280

var _font: FontFile
var _bg: ColorRect
var _overlay: ColorRect
var _title_label: Label
var _daily_label: Label

# 標準 Bonus 按鈕區
var _std_bg: ColorRect
var _std_title: Label
var _std_desc: Label
var _std_cost_label: Label
var _std_btn: Button

# TNT Bonus 按鈕區
var _tnt_bg: ColorRect
var _tnt_title: Label
var _tnt_desc: Label
var _tnt_cost_label: Label
var _tnt_btn: Button

var _close_btn: Button
var _status_label: Label

# 狀態
var _daily_left: int = 3
var _standard_cost: int = 0
var _tnt_cost: int = 0
var _can_buy: bool = true

func setup(font: FontFile) -> void:
	_font = font
	_build_ui()
	_connect_signals()
	hide()

func _build_ui() -> void:
	# 半透明遮罩
	_overlay = ColorRect.new()
	_overlay.color = Color(0.0, 0.0, 0.0, 0.65)
	_overlay.size = Vector2(1280, 720)
	_overlay.position = Vector2(-640, -180)
	add_child(_overlay)

	# 主面板
	_bg = ColorRect.new()
	_bg.color = Color(0.04, 0.08, 0.18, 0.97)
	_bg.size = Vector2(PANEL_W, PANEL_H)
	_bg.position = Vector2(-PANEL_W / 2, -PANEL_H / 2)
	add_child(_bg)

	# 頂部金色邊框
	var border = ColorRect.new()
	border.color = Color(1.0, 0.85, 0.1, 1.0)
	border.size = Vector2(PANEL_W, 4)
	border.position = Vector2(-PANEL_W / 2, -PANEL_H / 2)
	add_child(border)

	# 標題
	_title_label = Label.new()
	_title_label.text = "🎯 Buy Bonus"
	_title_label.position = Vector2(-PANEL_W / 2 + 12, -PANEL_H / 2 + 10)
	_title_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.1))
	if _font:
		_title_label.add_theme_font_override("font", _font)
		_title_label.add_theme_font_size_override("font_size", 20)
	add_child(_title_label)

	# 每日次數
	_daily_label = Label.new()
	_daily_label.text = "今日剩餘：3/3 次"
	_daily_label.position = Vector2(-PANEL_W / 2 + 12, -PANEL_H / 2 + 36)
	_daily_label.add_theme_color_override("font_color", Color(0.7, 0.9, 1.0))
	if _font:
		_daily_label.add_theme_font_override("font", _font)
		_daily_label.add_theme_font_size_override("font_size", 13)
	add_child(_daily_label)

	# ── 標準 Bonus 區塊 ──
	_std_bg = ColorRect.new()
	_std_bg.color = Color(0.08, 0.15, 0.08, 0.9)
	_std_bg.size = Vector2(PANEL_W - 24, 80)
	_std_bg.position = Vector2(-PANEL_W / 2 + 12, -PANEL_H / 2 + 60)
	add_child(_std_bg)

	_std_title = Label.new()
	_std_title.text = "🌿 標準 Bonus"
	_std_title.position = Vector2(-PANEL_W / 2 + 20, -PANEL_H / 2 + 65)
	_std_title.add_theme_color_override("font_color", Color(0.5, 1.0, 0.5))
	if _font:
		_std_title.add_theme_font_override("font", _font)
		_std_title.add_theme_font_size_override("font_size", 15)
	add_child(_std_title)

	_std_desc = Label.new()
	_std_desc.text = "直接觸發 Bonus Game（期望 ×60）"
	_std_desc.position = Vector2(-PANEL_W / 2 + 20, -PANEL_H / 2 + 84)
	_std_desc.add_theme_color_override("font_color", Color(0.8, 0.9, 0.8))
	if _font:
		_std_desc.add_theme_font_override("font", _font)
		_std_desc.add_theme_font_size_override("font_size", 11)
	add_child(_std_desc)

	_std_cost_label = Label.new()
	_std_cost_label.text = "費用：計算中..."
	_std_cost_label.position = Vector2(-PANEL_W / 2 + 20, -PANEL_H / 2 + 100)
	_std_cost_label.add_theme_color_override("font_color", Color(1.0, 0.9, 0.3))
	if _font:
		_std_cost_label.add_theme_font_override("font", _font)
		_std_cost_label.add_theme_font_size_override("font_size", 12)
	add_child(_std_cost_label)

	_std_btn = Button.new()
	_std_btn.text = "購買"
	_std_btn.size = Vector2(70, 28)
	_std_btn.position = Vector2(PANEL_W / 2 - 90, -PANEL_H / 2 + 80)
	_std_btn.pressed.connect(_on_buy_standard)
	if _font:
		_std_btn.add_theme_font_override("font", _font)
		_std_btn.add_theme_font_size_override("font_size", 13)
	add_child(_std_btn)

	# ── TNT Bonus 區塊 ──
	_tnt_bg = ColorRect.new()
	_tnt_bg.color = Color(0.18, 0.08, 0.04, 0.9)
	_tnt_bg.size = Vector2(PANEL_W - 24, 80)
	_tnt_bg.position = Vector2(-PANEL_W / 2 + 12, -PANEL_H / 2 + 152)
	add_child(_tnt_bg)

	_tnt_title = Label.new()
	_tnt_title.text = "💥 TNT Bonus"
	_tnt_title.position = Vector2(-PANEL_W / 2 + 20, -PANEL_H / 2 + 157)
	_tnt_title.add_theme_color_override("font_color", Color(1.0, 0.5, 0.2))
	if _font:
		_tnt_title.add_theme_font_override("font", _font)
		_tnt_title.add_theme_font_size_override("font_size", 15)
	add_child(_tnt_title)

	_tnt_desc = Label.new()
	_tnt_desc.text = "觸發 Bonus Game + 1.5x 倍率加成（期望 ×100）"
	_tnt_desc.position = Vector2(-PANEL_W / 2 + 20, -PANEL_H / 2 + 176)
	_tnt_desc.add_theme_color_override("font_color", Color(1.0, 0.8, 0.6))
	if _font:
		_tnt_desc.add_theme_font_override("font", _font)
		_tnt_desc.add_theme_font_size_override("font_size", 11)
	add_child(_tnt_desc)

	_tnt_cost_label = Label.new()
	_tnt_cost_label.text = "費用：計算中..."
	_tnt_cost_label.position = Vector2(-PANEL_W / 2 + 20, -PANEL_H / 2 + 192)
	_tnt_cost_label.add_theme_color_override("font_color", Color(1.0, 0.9, 0.3))
	if _font:
		_tnt_cost_label.add_theme_font_override("font", _font)
		_tnt_cost_label.add_theme_font_size_override("font_size", 12)
	add_child(_tnt_cost_label)

	_tnt_btn = Button.new()
	_tnt_btn.text = "購買"
	_tnt_btn.size = Vector2(70, 28)
	_tnt_btn.position = Vector2(PANEL_W / 2 - 90, -PANEL_H / 2 + 172)
	_tnt_btn.pressed.connect(_on_buy_tnt)
	if _font:
		_tnt_btn.add_theme_font_override("font", _font)
		_tnt_btn.add_theme_font_size_override("font_size", 13)
	add_child(_tnt_btn)

	# 狀態標籤
	_status_label = Label.new()
	_status_label.text = ""
	_status_label.position = Vector2(-PANEL_W / 2 + 12, PANEL_H / 2 - 44)
	_status_label.add_theme_color_override("font_color", Color(0.5, 1.0, 0.5))
	if _font:
		_status_label.add_theme_font_override("font", _font)
		_status_label.add_theme_font_size_override("font_size", 12)
	add_child(_status_label)

	# 關閉按鈕
	_close_btn = Button.new()
	_close_btn.text = "✕ 關閉"
	_close_btn.size = Vector2(80, 28)
	_close_btn.position = Vector2(PANEL_W / 2 - 92, PANEL_H / 2 - 36)
	_close_btn.pressed.connect(_on_close)
	if _font:
		_close_btn.add_theme_font_override("font", _font)
		_close_btn.add_theme_font_size_override("font_size", 12)
	add_child(_close_btn)

func _connect_signals() -> void:
	if GameManager.has_signal("buy_bonus_status"):
		GameManager.buy_bonus_status.connect(_on_buy_bonus_status)
	if GameManager.has_signal("buy_bonus_success"):
		GameManager.buy_bonus_success.connect(_on_buy_bonus_success)
	if GameManager.has_signal("buy_bonus_error"):
		GameManager.buy_bonus_error.connect(_on_buy_bonus_error)

func show_panel() -> void:
	# 請求最新狀態
	GameManager.send_message("get_buy_bonus_status", {})
	show()
	scale = Vector2(0.8, 0.8)
	var tween = create_tween()
	tween.tween_property(self, "scale", Vector2(1.05, 1.05), 0.12)
	tween.tween_property(self, "scale", Vector2(1.0, 1.0), 0.08)

func _on_buy_bonus_status(data: Dictionary) -> void:
	_daily_left = data.get("daily_left", 0)
	_standard_cost = data.get("standard_cost", 0)
	_tnt_cost = data.get("tnt_cost", 0)
	_can_buy = data.get("can_buy", false)

	var daily_used: int = data.get("daily_used", 0)
	var daily_limit: int = data.get("daily_limit", 3)
	_daily_label.text = "今日剩餘：%d/%d 次" % [_daily_left, daily_limit]

	_std_cost_label.text = "費用：🪙 %d 金幣" % _standard_cost
	_tnt_cost_label.text = "費用：🪙 %d 金幣" % _tnt_cost

	# 更新按鈕狀態
	_std_btn.disabled = not _can_buy
	_tnt_btn.disabled = not _can_buy

	if not _can_buy:
		if _daily_left <= 0:
			_status_label.text = "⚠️ 今日購買次數已達上限"
			_status_label.add_theme_color_override("font_color", Color(1.0, 0.5, 0.2))
		else:
			_status_label.text = "⚠️ 遊戲進行中，請稍後"
			_status_label.add_theme_color_override("font_color", Color(1.0, 0.5, 0.2))
	else:
		_status_label.text = ""

func _on_buy_bonus_success(data: Dictionary) -> void:
	var bonus_type: String = data.get("bonus_type", "standard")
	var cost: int = data.get("cost", 0)
	var daily_left: int = data.get("daily_left", 0)
	var mult_bonus: float = data.get("mult_bonus", 1.0)

	_daily_left = daily_left
	_daily_label.text = "今日剩餘：%d/3 次" % daily_left

	if bonus_type == "tnt":
		_status_label.text = "💥 TNT Bonus 已購買！倍率 ×%.1f" % mult_bonus
	else:
		_status_label.text = "✅ 標準 Bonus 已購買！-🪙%d" % cost
	_status_label.add_theme_color_override("font_color", Color(0.5, 1.0, 0.5))

	# 購買成功後關閉面板
	var tween = create_tween()
	tween.tween_interval(1.5)
	tween.tween_callback(_on_close)

func _on_buy_bonus_error(data: Dictionary) -> void:
	var reason: String = data.get("reason", "")
	var message: String = data.get("message", "購買失敗")
	var cost: int = data.get("cost", 0)
	var balance: int = data.get("balance", 0)

	if reason == "insufficient_coins":
		_status_label.text = "❌ 金幣不足（需要 %d，現有 %d）" % [cost, balance]
	elif reason == "daily_limit":
		_status_label.text = "❌ 今日購買次數已達上限"
	else:
		_status_label.text = "❌ %s" % message
	_status_label.add_theme_color_override("font_color", Color(1.0, 0.3, 0.3))

func _on_buy_standard() -> void:
	_status_label.text = "購買中..."
	_status_label.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8))
	GameManager.send_message("buy_bonus", {"bonus_type": "standard"})

func _on_buy_tnt() -> void:
	_status_label.text = "購買中..."
	_status_label.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8))
	GameManager.send_message("buy_bonus", {"bonus_type": "tnt"})

func _on_close() -> void:
	var tween = create_tween()
	tween.tween_property(self, "scale", Vector2(0.0, 0.0), 0.12)
	tween.tween_callback(func():
		scale = Vector2(1.0, 1.0)
		hide()
	)
