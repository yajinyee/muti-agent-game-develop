## ChainLongKingPanel.gd — 長龍王雙環輪盤面板（DAY-194）
## 業界依據：Royal Fishing JILI「ChainLong King — dual-ring roulette activates when captured.
## You control when the pointer stops, multiplying inner and outer ring values together.
## Maximum combination delivers 350X, whilst the ChainLong King itself can award up to 1000X mega wins.」
##
## 視覺設計：
##   - 金龍主題（深金色 #B8860B + 亮金色 #FFD700）
##   - roulette_start：金色三次閃光 + 雙環輪盤出現（內環旋轉快，外環旋轉慢）
##   - inner_stop：內環停止動畫（彈跳）+ 顯示結果
##   - outer_stop：外環停止動畫（彈跳）+ 顯示結果
##   - result：結算彈窗（內環 × 外環 = 總倍率）+ 依倍率決定特效
##   - mega_win：全螢幕金色爆炸 + 1000x 大字 + 龍形粒子
##   - broadcast：全服廣播橫幅（其他玩家看到）
extends CanvasLayer

# 面板根節點
var _panel: Control
var _inner_ring_label: Label
var _outer_ring_label: Label
var _result_popup: Control
var _broadcast_banner: Label
var _spin_button: Button
var _mega_overlay: Control

# 輪盤狀態
var _instance_id: String = ""
var _phase: String = ""  # "inner_spin" / "outer_spin" / "result"
var _inner_result: int = 0
var _outer_result: int = 0
var _inner_ring: Array = [5, 10, 20, 50]
var _outer_ring: Array = [1, 2, 3, 5, 7]

# 旋轉動畫
var _inner_spin_angle: float = 0.0
var _outer_spin_angle: float = 0.0
var _inner_spin_speed: float = 720.0  # 度/秒（快）
var _outer_spin_speed: float = 360.0  # 度/秒（慢）
var _is_spinning: bool = false

func _ready() -> void:
	layer = 51
	_build_ui()

func _build_ui() -> void:
	# 主面板（預設隱藏）
	_panel = Control.new()
	_panel.set_anchors_preset(Control.PRESET_FULL_RECT)
	_panel.visible = false
	add_child(_panel)

	# 半透明背景
	var bg = ColorRect.new()
	bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	bg.color = Color(0, 0, 0, 0.6)
	_panel.add_child(bg)

	# 標題橫幅
	var title = Label.new()
	title.text = "🐉 長龍王雙環輪盤"
	title.add_theme_font_size_override("font_size", 28)
	title.add_theme_color_override("font_color", Color("#FFD700"))
	title.set_anchors_preset(Control.PRESET_TOP_WIDE)
	title.position = Vector2(0, 40)
	title.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_panel.add_child(title)

	# 內環顯示
	_inner_ring_label = Label.new()
	_inner_ring_label.text = "內環: ?"
	_inner_ring_label.add_theme_font_size_override("font_size", 36)
	_inner_ring_label.add_theme_color_override("font_color", Color("#FFD700"))
	_inner_ring_label.position = Vector2(0, 160)
	_inner_ring_label.set_anchors_preset(Control.PRESET_TOP_WIDE)
	_inner_ring_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_panel.add_child(_inner_ring_label)

	# 外環顯示
	_outer_ring_label = Label.new()
	_outer_ring_label.text = "外環: ?"
	_outer_ring_label.add_theme_font_size_override("font_size", 36)
	_outer_ring_label.add_theme_color_override("font_color", Color("#FFA500"))
	_outer_ring_label.position = Vector2(0, 220)
	_outer_ring_label.set_anchors_preset(Control.PRESET_TOP_WIDE)
	_outer_ring_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_panel.add_child(_outer_ring_label)

	# 停止按鈕
	_spin_button = Button.new()
	_spin_button.text = "🛑 停止內環"
	_spin_button.add_theme_font_size_override("font_size", 22)
	_spin_button.position = Vector2(0, 300)
	_spin_button.size = Vector2(200, 50)
	_spin_button.set_anchors_preset(Control.PRESET_CENTER)
	_spin_button.position.y = 300
	_spin_button.pressed.connect(_on_stop_pressed)
	_panel.add_child(_spin_button)

	# 結算彈窗（預設隱藏）
	_result_popup = Control.new()
	_result_popup.visible = false
	_result_popup.set_anchors_preset(Control.PRESET_CENTER)
	_result_popup.size = Vector2(320, 200)
	_result_popup.position = Vector2(-160, -100)
	_panel.add_child(_result_popup)

	var popup_bg = ColorRect.new()
	popup_bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	popup_bg.color = Color(0.1, 0.08, 0.0, 0.95)
	_result_popup.add_child(popup_bg)

	var popup_label = Label.new()
	popup_label.name = "ResultLabel"
	popup_label.add_theme_font_size_override("font_size", 20)
	popup_label.add_theme_color_override("font_color", Color("#FFD700"))
	popup_label.set_anchors_preset(Control.PRESET_FULL_RECT)
	popup_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	popup_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_result_popup.add_child(popup_label)

	# 全服廣播橫幅
	_broadcast_banner = Label.new()
	_broadcast_banner.text = ""
	_broadcast_banner.add_theme_font_size_override("font_size", 18)
	_broadcast_banner.add_theme_color_override("font_color", Color("#FFD700"))
	_broadcast_banner.set_anchors_preset(Control.PRESET_TOP_WIDE)
	_broadcast_banner.position = Vector2(0, 10)
	_broadcast_banner.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_broadcast_banner.visible = false
	add_child(_broadcast_banner)

	# 千倍大獎覆蓋層（預設隱藏）
	_mega_overlay = Control.new()
	_mega_overlay.set_anchors_preset(Control.PRESET_FULL_RECT)
	_mega_overlay.visible = false
	add_child(_mega_overlay)

	var mega_bg = ColorRect.new()
	mega_bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	mega_bg.color = Color(0.8, 0.6, 0.0, 0.85)
	_mega_overlay.add_child(mega_bg)

	var mega_label = Label.new()
	mega_label.text = "🐉 千倍大獎！\n1000x"
	mega_label.add_theme_font_size_override("font_size", 56)
	mega_label.add_theme_color_override("font_color", Color.WHITE)
	mega_label.set_anchors_preset(Control.PRESET_FULL_RECT)
	mega_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	mega_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_mega_overlay.add_child(mega_label)

func _process(delta: float) -> void:
	if not _is_spinning:
		return
	_inner_spin_angle += _inner_spin_speed * delta
	_outer_spin_angle += _outer_spin_speed * delta
	# 更新旋轉顯示（用角度模擬輪盤位置）
	var inner_idx = int(_inner_spin_angle / 90.0) % _inner_ring.size()
	var outer_idx = int(_outer_spin_angle / 72.0) % _outer_ring.size()
	if _phase == "inner_spin":
		_inner_ring_label.text = "內環: %dx ←旋轉中" % _inner_ring[inner_idx]
		_outer_ring_label.text = "外環: %dx ←等待" % _outer_ring[outer_idx]
	elif _phase == "outer_spin":
		_outer_ring_label.text = "外環: %dx ←旋轉中" % _outer_ring[outer_idx]

## 處理長龍王訊息
func handle_chainlong_king(payload: Dictionary) -> void:
	var phase = payload.get("phase", "")
	match phase:
		"roulette_start":
			_on_roulette_start(payload)
		"inner_stop":
			_on_inner_stop(payload)
		"outer_stop":
			_on_outer_stop(payload)
		"result":
			_on_result(payload)
		"mega_win":
			_on_mega_win(payload)
		"broadcast", "mega_broadcast":
			_on_broadcast(payload)

func _on_roulette_start(payload: Dictionary) -> void:
	_instance_id = payload.get("instance_id", "")
	_phase = "inner_spin"
	_inner_result = 0
	_outer_result = 0
	_is_spinning = true

	if payload.has("inner_ring"):
		_inner_ring = payload["inner_ring"]
	if payload.has("outer_ring"):
		_outer_ring = payload["outer_ring"]

	_panel.visible = true
	_result_popup.visible = false
	_spin_button.text = "🛑 停止內環"
	_spin_button.visible = true

	# 金色三次閃光
	_flash_screen(Color("#FFD700"), 0.3)
	await get_tree().create_timer(0.35).timeout
	_flash_screen(Color("#FFD700"), 0.25)
	await get_tree().create_timer(0.3).timeout
	_flash_screen(Color("#FFD700"), 0.2)

func _on_inner_stop(payload: Dictionary) -> void:
	_inner_result = payload.get("inner_result", 5)
	_phase = "outer_spin"
	var is_timeout = payload.get("is_timeout", false)

	# 內環停止動畫（彈跳）
	_inner_ring_label.text = "內環: %dx ✓" % _inner_result
	_inner_ring_label.add_theme_color_override("font_color", Color("#00FF88"))
	var tween = create_tween()
	tween.tween_property(_inner_ring_label, "scale", Vector2(1.3, 1.3), 0.1)
	tween.tween_property(_inner_ring_label, "scale", Vector2(1.0, 1.0), 0.15)

	_spin_button.text = "🛑 停止外環"
	if is_timeout:
		_spin_button.text = "⏱ 超時自動停止"

func _on_outer_stop(payload: Dictionary) -> void:
	_outer_result = payload.get("outer_result", 1)
	_phase = "result"
	_is_spinning = false
	var is_timeout = payload.get("is_timeout", false)

	# 外環停止動畫（彈跳）
	_outer_ring_label.text = "外環: %dx ✓" % _outer_result
	_outer_ring_label.add_theme_color_override("font_color", Color("#00FF88"))
	var tween = create_tween()
	tween.tween_property(_outer_ring_label, "scale", Vector2(1.3, 1.3), 0.1)
	tween.tween_property(_outer_ring_label, "scale", Vector2(1.0, 1.0), 0.15)

	_spin_button.visible = false
	if is_timeout:
		_spin_button.text = "⏱ 超時自動停止"

func _on_result(payload: Dictionary) -> void:
	var inner = payload.get("inner_result", _inner_result)
	var outer = payload.get("outer_result", _outer_result)
	var total_mult = payload.get("total_mult", inner * outer)
	var reward = payload.get("reward", 0)
	var is_big_win = payload.get("is_big_win", false)

	# 顯示結算彈窗
	_result_popup.visible = true
	var label = _result_popup.get_node("ResultLabel")
	label.text = "🐉 長龍王輪盤結算\n內環 %dx × 外環 %dx\n= %dx\n獎勵: %d 金幣" % [inner, outer, total_mult, reward]

	# 依倍率決定特效
	if total_mult >= 350:
		# 最高倍率：橙紅色強閃光 × 3
		_flash_screen(Color("#FF4500"), 0.5)
		await get_tree().create_timer(0.55).timeout
		_flash_screen(Color("#FF4500"), 0.4)
		await get_tree().create_timer(0.45).timeout
		_flash_screen(Color("#FFD700"), 0.6)
		label.add_theme_color_override("font_color", Color("#FF4500"))
	elif total_mult >= 150:
		# 大獎：金色雙閃光
		_flash_screen(Color("#FFD700"), 0.4)
		await get_tree().create_timer(0.45).timeout
		_flash_screen(Color("#FFD700"), 0.3)
		label.add_theme_color_override("font_color", Color("#FFD700"))
	elif is_big_win:
		# 100x+：橙色閃光
		_flash_screen(Color("#FFA500"), 0.3)

	# 4 秒後淡出
	await get_tree().create_timer(4.0).timeout
	_fade_out()

func _on_mega_win(payload: Dictionary) -> void:
	var reward = payload.get("reward", 0)
	_is_spinning = false
	_panel.visible = false

	# 千倍大獎：全螢幕金色爆炸
	_mega_overlay.visible = true
	_mega_overlay.modulate.a = 0.0
	var tween = create_tween()
	tween.tween_property(_mega_overlay, "modulate:a", 1.0, 0.3)

	# 三次強閃光
	for i in range(3):
		_flash_screen(Color("#FFD700"), 0.6)
		await get_tree().create_timer(0.65).timeout

	# 顯示獎勵
	var reward_label = Label.new()
	reward_label.text = "獎勵: %d 金幣！" % reward
	reward_label.add_theme_font_size_override("font_size", 32)
	reward_label.add_theme_color_override("font_color", Color.WHITE)
	reward_label.set_anchors_preset(Control.PRESET_BOTTOM_WIDE)
	reward_label.position.y = -80
	reward_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_mega_overlay.add_child(reward_label)

	# 5 秒後淡出
	await get_tree().create_timer(5.0).timeout
	var fade_tween = create_tween()
	fade_tween.tween_property(_mega_overlay, "modulate:a", 0.0, 0.5)
	await fade_tween.finished
	_mega_overlay.visible = false
	reward_label.queue_free()

func _on_broadcast(payload: Dictionary) -> void:
	var player_name = payload.get("player_name", "玩家")
	var total_mult = payload.get("total_mult", 0)
	var is_mega = payload.get("is_mega", false)

	var text: String
	if is_mega:
		text = "🐉 %s 觸發長龍王千倍大獎！1000x！" % player_name
		_broadcast_banner.add_theme_color_override("font_color", Color("#FFD700"))
	elif total_mult >= 350:
		text = "🐉 %s 長龍王最高倍率 %dx！" % [player_name, total_mult]
		_broadcast_banner.add_theme_color_override("font_color", Color("#FF4500"))
	else:
		text = "🐉 %s 長龍王輪盤獲得 %dx！" % [player_name, total_mult]
		_broadcast_banner.add_theme_color_override("font_color", Color("#FFA500"))

	_broadcast_banner.text = text
	_broadcast_banner.visible = true
	_broadcast_banner.modulate.a = 1.0

	# 3 秒後淡出
	await get_tree().create_timer(3.0).timeout
	var tween = create_tween()
	tween.tween_property(_broadcast_banner, "modulate:a", 0.0, 0.5)
	await tween.finished
	_broadcast_banner.visible = false

func _on_stop_pressed() -> void:
	if _instance_id.is_empty():
		return
	# 發送停止訊號給 GameManager
	if Engine.has_singleton("GameManager"):
		pass
	# 透過訊號通知 GameManager 發送 WebSocket 訊息
	emit_signal("chainlong_king_stop_pressed", _instance_id)

func _fade_out() -> void:
	_is_spinning = false
	var tween = create_tween()
	tween.tween_property(_panel, "modulate:a", 0.0, 0.4)
	await tween.finished
	_panel.visible = false
	_panel.modulate.a = 1.0
	_result_popup.visible = false
	_spin_button.visible = true
	_inner_ring_label.add_theme_color_override("font_color", Color("#FFD700"))
	_outer_ring_label.add_theme_color_override("font_color", Color("#FFA500"))
	_inner_ring_label.text = "內環: ?"
	_outer_ring_label.text = "外環: ?"
	_instance_id = ""
	_phase = ""

func _flash_screen(color: Color, intensity: float) -> void:
	var flash = ColorRect.new()
	flash.set_anchors_preset(Control.PRESET_FULL_RECT)
	flash.color = Color(color.r, color.g, color.b, intensity)
	flash.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(flash)
	var tween = create_tween()
	tween.tween_property(flash, "modulate:a", 0.0, 0.25)
	await tween.finished
	flash.queue_free()

signal chainlong_king_stop_pressed(instance_id: String)
