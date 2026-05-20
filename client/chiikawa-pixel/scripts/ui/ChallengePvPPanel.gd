## ChallengePvPPanel.gd — 好友挑戰面板（DAY-102）
## 顯示挑戰邀請、進行中的挑戰分數、挑戰結果
## 整合到 FriendPanel 的好友行（挑戰按鈕）
extends Node2D

# ---- 常數 ----
const PANEL_WIDTH  := 320
const PANEL_HEIGHT := 160

# ---- 節點引用 ----
var _pixel_font: Font = null
var _panel_bg: ColorRect = null
var _my_score_label: Label = null
var _opponent_score_label: Label = null
var _timer_label: Label = null
var _status_label: Label = null

# ---- 挑戰狀態 ----
var _challenge_id: String = ""
var _opponent_id: String = ""
var _opponent_name: String = ""
var _my_score: int = 0
var _opponent_score: int = 0
var _time_remaining: int = 0
var _is_active: bool = false

# ---- 初始化 ----
func _ready() -> void:
	if ResourceLoader.exists("res://assets/fonts/pixel8.fnt"):
		_pixel_font = load("res://assets/fonts/pixel8.fnt")
	_build_panel()
	_connect_signals()
	visible = false

func setup(font: Font) -> void:
	if font:
		_pixel_font = font

func _process(_delta: float) -> void:
	if _is_active and _time_remaining > 0:
		_time_remaining -= 1
		_update_timer_display()

## 建立面板
func _build_panel() -> void:
	_panel_bg = ColorRect.new()
	_panel_bg.position = Vector2(0, 0)
	_panel_bg.size = Vector2(PANEL_WIDTH, PANEL_HEIGHT)
	_panel_bg.color = Color(0.05, 0.02, 0.18, 0.95)
	add_child(_panel_bg)

	# 標題
	var title := Label.new()
	title.position = Vector2(8, 4)
	title.text = "⚔️ 好友挑戰進行中！"
	title.add_theme_color_override("font_color", Color(1.0, 0.85, 0.2))
	if _pixel_font:
		title.add_theme_font_override("font", _pixel_font)
		title.add_theme_font_size_override("font_size", 13)
	_panel_bg.add_child(title)

	# 計時器
	_timer_label = Label.new()
	_timer_label.position = Vector2(PANEL_WIDTH - 70, 4)
	_timer_label.text = "⏱ 3:00"
	_timer_label.add_theme_color_override("font_color", Color(0.4, 1.0, 0.6))
	if _pixel_font:
		_timer_label.add_theme_font_override("font", _pixel_font)
		_timer_label.add_theme_font_size_override("font_size", 12)
	_panel_bg.add_child(_timer_label)

	# 分隔線
	var sep := ColorRect.new()
	sep.position = Vector2(4, 22)
	sep.size = Vector2(PANEL_WIDTH - 8, 1)
	sep.color = Color(0.4, 0.3, 0.6, 0.6)
	_panel_bg.add_child(sep)

	# 我的分數
	var my_label := Label.new()
	my_label.position = Vector2(8, 28)
	my_label.text = "我："
	my_label.add_theme_color_override("font_color", Color(0.6, 0.9, 1.0))
	if _pixel_font:
		my_label.add_theme_font_override("font", _pixel_font)
		my_label.add_theme_font_size_override("font_size", 10)
	_panel_bg.add_child(my_label)

	_my_score_label = Label.new()
	_my_score_label.position = Vector2(40, 28)
	_my_score_label.text = "0"
	_my_score_label.add_theme_color_override("font_color", Color(1.0, 1.0, 1.0))
	if _pixel_font:
		_my_score_label.add_theme_font_override("font", _pixel_font)
		_my_score_label.add_theme_font_size_override("font_size", 14)
	_panel_bg.add_child(_my_score_label)

	# VS
	var vs_label := Label.new()
	vs_label.position = Vector2(PANEL_WIDTH / 2 - 12, 28)
	vs_label.text = "VS"
	vs_label.add_theme_color_override("font_color", Color(1.0, 0.5, 0.2))
	if _pixel_font:
		vs_label.add_theme_font_override("font", _pixel_font)
		vs_label.add_theme_font_size_override("font_size", 12)
	_panel_bg.add_child(vs_label)

	# 對手分數
	var opp_label := Label.new()
	opp_label.position = Vector2(PANEL_WIDTH / 2 + 20, 28)
	opp_label.text = "對手："
	opp_label.add_theme_color_override("font_color", Color(1.0, 0.6, 0.6))
	if _pixel_font:
		opp_label.add_theme_font_override("font", _pixel_font)
		opp_label.add_theme_font_size_override("font_size", 10)
	_panel_bg.add_child(opp_label)

	_opponent_score_label = Label.new()
	_opponent_score_label.position = Vector2(PANEL_WIDTH / 2 + 65, 28)
	_opponent_score_label.text = "0"
	_opponent_score_label.add_theme_color_override("font_color", Color(1.0, 0.7, 0.7))
	if _pixel_font:
		_opponent_score_label.add_theme_font_override("font", _pixel_font)
		_opponent_score_label.add_theme_font_size_override("font_size", 14)
	_panel_bg.add_child(_opponent_score_label)

	# 狀態文字
	_status_label = Label.new()
	_status_label.position = Vector2(8, 55)
	_status_label.text = "賭注：各 1000🪙，勝者獲得全部！"
	_status_label.add_theme_color_override("font_color", Color(0.8, 0.8, 0.6))
	if _pixel_font:
		_status_label.add_theme_font_override("font", _pixel_font)
		_status_label.add_theme_font_size_override("font_size", 9)
	_panel_bg.add_child(_status_label)

## 連接訊號
func _connect_signals() -> void:
	if GameManager.has_signal("challenge_request"):
		GameManager.challenge_request.connect(_on_challenge_request)
	if GameManager.has_signal("challenge_updated"):
		GameManager.challenge_updated.connect(_on_challenge_updated)
	if GameManager.has_signal("challenge_result"):
		GameManager.challenge_result.connect(_on_challenge_result)
	if GameManager.has_signal("challenge_error"):
		GameManager.challenge_error.connect(_on_challenge_error)

## 收到挑戰邀請
func _on_challenge_request(data: Dictionary) -> void:
	var challenger_name = data.get("challenger_name", "好友")
	var challenge_id = data.get("challenge_id", "")
	var stake = data.get("stake", 1000)
	var expires = data.get("expires_in_sec", 30)

	# 顯示挑戰邀請彈窗
	_show_challenge_invite(challenge_id, challenger_name, stake, expires)

## 挑戰狀態更新
func _on_challenge_updated(data: Dictionary) -> void:
	_challenge_id = data.get("challenge_id", "")
	_opponent_id = data.get("opponent_id", "")
	_opponent_name = data.get("opponent_name", "對手")
	_my_score = data.get("my_score", 0)
	_opponent_score = data.get("opponent_score", 0)
	_time_remaining = data.get("time_remaining", 0)

	var status = data.get("status", "")
	if status == "active":
		_is_active = true
		visible = true
		_update_score_display()
		_update_timer_display()
	elif status == "pending":
		_status_label.text = "等待 %s 接受挑戰..." % _opponent_name

## 挑戰結果
func _on_challenge_result(data: Dictionary) -> void:
	_is_active = false
	visible = false

	var is_winner = data.get("is_winner", false)
	var is_draw = data.get("is_draw", false)
	var prize = data.get("prize", 0)
	var my_score = data.get("my_score", 0)
	var opp_score = data.get("opponent_score", 0)
	var opp_name = data.get("opponent_name", "對手")

	var result_text: String
	var result_color: Color
	if is_draw:
		result_text = "⚔️ 平局！各退回 %d🪙\n你：%d vs %s：%d" % [prize, my_score, opp_name, opp_score]
		result_color = Color(0.8, 0.8, 0.4)
	elif is_winner:
		result_text = "🏆 你贏了！獲得 %d🪙\n你：%d vs %s：%d" % [prize, my_score, opp_name, opp_score]
		result_color = Color(1.0, 0.85, 0.2)
	else:
		result_text = "💔 你輸了...\n你：%d vs %s：%d" % [my_score, opp_name, opp_score]
		result_color = Color(1.0, 0.4, 0.4)

	_show_result_notification(result_text, result_color)

## 挑戰錯誤
func _on_challenge_error(data: Dictionary) -> void:
	var msg = data.get("message", "挑戰失敗")
	_show_notification("❌ %s" % msg, Color(1.0, 0.4, 0.4))

## 更新分數顯示
func _update_score_display() -> void:
	if is_instance_valid(_my_score_label):
		_my_score_label.text = str(_my_score)
		# 領先時顯示金色
		if _my_score > _opponent_score:
			_my_score_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.2))
		else:
			_my_score_label.add_theme_color_override("font_color", Color(1.0, 1.0, 1.0))

	if is_instance_valid(_opponent_score_label):
		_opponent_score_label.text = str(_opponent_score)
		if _opponent_score > _my_score:
			_opponent_score_label.add_theme_color_override("font_color", Color(1.0, 0.5, 0.5))
		else:
			_opponent_score_label.add_theme_color_override("font_color", Color(0.8, 0.6, 0.6))

## 更新計時器顯示
func _update_timer_display() -> void:
	if not is_instance_valid(_timer_label):
		return
	var mins = _time_remaining / 60
	var secs = _time_remaining % 60
	_timer_label.text = "⏱ %d:%02d" % [mins, secs]
	# 最後 30 秒變紅色閃爍
	if _time_remaining <= 30:
		_timer_label.add_theme_color_override("font_color", Color(1.0, 0.3, 0.3))
	else:
		_timer_label.add_theme_color_override("font_color", Color(0.4, 1.0, 0.6))

## 顯示挑戰邀請彈窗
func _show_challenge_invite(challenge_id: String, challenger_name: String, stake: int, expires: int) -> void:
	var invite_bg := ColorRect.new()
	invite_bg.position = Vector2(640 - 160, 360 - 60)
	invite_bg.size = Vector2(320, 120)
	invite_bg.color = Color(0.05, 0.02, 0.18, 0.97)
	invite_bg.z_index = 90
	get_tree().root.add_child(invite_bg)

	var title := Label.new()
	title.position = Vector2(8, 8)
	title.text = "⚔️ %s 向你發起挑戰！" % challenger_name
	title.add_theme_color_override("font_color", Color(1.0, 0.85, 0.2))
	if _pixel_font:
		title.add_theme_font_override("font", _pixel_font)
		title.add_theme_font_size_override("font_size", 12)
	invite_bg.add_child(title)

	var desc := Label.new()
	desc.position = Vector2(8, 28)
	desc.text = "賭注：各 %d🪙，3分鐘比分數，勝者獲得全部！" % stake
	desc.add_theme_color_override("font_color", Color(0.8, 0.8, 0.6))
	if _pixel_font:
		desc.add_theme_font_override("font", _pixel_font)
		desc.add_theme_font_size_override("font_size", 9)
	invite_bg.add_child(desc)

	var timer_lbl := Label.new()
	timer_lbl.position = Vector2(8, 44)
	timer_lbl.text = "⏱ %d 秒後過期" % expires
	timer_lbl.add_theme_color_override("font_color", Color(0.7, 0.7, 0.7))
	if _pixel_font:
		timer_lbl.add_theme_font_override("font", _pixel_font)
		timer_lbl.add_theme_font_size_override("font_size", 9)
	invite_bg.add_child(timer_lbl)

	# 接受按鈕
	var accept_btn := Button.new()
	accept_btn.text = "✓ 接受"
	accept_btn.position = Vector2(60, 70)
	accept_btn.size = Vector2(80, 28)
	if _pixel_font:
		accept_btn.add_theme_font_override("font", _pixel_font)
		accept_btn.add_theme_font_size_override("font_size", 11)
	accept_btn.pressed.connect(func():
		NetworkManager.send_message({
			"type": "accept_challenge",
			"payload": {"challenge_id": challenge_id}
		})
		if is_instance_valid(invite_bg): invite_bg.queue_free()
	)
	invite_bg.add_child(accept_btn)

	# 拒絕按鈕
	var decline_btn := Button.new()
	decline_btn.text = "✕ 拒絕"
	decline_btn.position = Vector2(180, 70)
	decline_btn.size = Vector2(80, 28)
	if _pixel_font:
		decline_btn.add_theme_font_override("font", _pixel_font)
		decline_btn.add_theme_font_size_override("font_size", 11)
	decline_btn.pressed.connect(func():
		NetworkManager.send_message({
			"type": "decline_challenge",
			"payload": {"challenge_id": challenge_id}
		})
		if is_instance_valid(invite_bg): invite_bg.queue_free()
	)
	invite_bg.add_child(decline_btn)

	# 自動過期
	var tween = create_tween()
	tween.tween_interval(float(expires))
	tween.tween_callback(func():
		if is_instance_valid(invite_bg): invite_bg.queue_free()
	)

## 顯示結果通知
func _show_result_notification(text: String, color: Color) -> void:
	var notify_bg := ColorRect.new()
	notify_bg.position = Vector2(640 - 180, 360 - 50)
	notify_bg.size = Vector2(360, 100)
	notify_bg.color = Color(0.05, 0.02, 0.18, 0.97)
	notify_bg.z_index = 90
	get_tree().root.add_child(notify_bg)

	var label := Label.new()
	label.position = Vector2(8, 8)
	label.text = text
	label.add_theme_color_override("font_color", color)
	if _pixel_font:
		label.add_theme_font_override("font", _pixel_font)
		label.add_theme_font_size_override("font_size", 12)
	notify_bg.add_child(label)

	var tween = create_tween()
	tween.tween_interval(3.0)
	tween.tween_property(notify_bg, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func():
		if is_instance_valid(notify_bg): notify_bg.queue_free()
	)

## 顯示通知
func _show_notification(text: String, color: Color) -> void:
	var notify := Label.new()
	notify.text = text
	notify.position = Vector2(640 - 100, 300)
	notify.add_theme_color_override("font_color", color)
	if _pixel_font:
		notify.add_theme_font_override("font", _pixel_font)
		notify.add_theme_font_size_override("font_size", 11)
	get_tree().root.add_child(notify)

	var tween = create_tween()
	tween.tween_interval(2.5)
	tween.tween_property(notify, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func():
		if is_instance_valid(notify): notify.queue_free()
	)
