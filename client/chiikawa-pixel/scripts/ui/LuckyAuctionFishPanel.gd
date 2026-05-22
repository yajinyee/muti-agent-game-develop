## LuckyAuctionFishPanel.gd — 幸運拍賣魚系統面板（DAY-217）
## 業界原創「全服競標」機制
##
## 視覺設計：
##   - 黃金競標主題（#FFD700 + #FF8C00 + #FFF8DC + #FF4500）
##   - auction_start：金色雙閃光 + 頂部橫幅 + 底部競標進度條 + 出價按鈕
##   - auction_bid：出價閃光 + 「💰 [玩家] 出價 N！最高 M！」浮動文字
##   - auction_result：金色三次強閃光 + 「🏆 [贏家] 獲得大獎控制權！」大字 + 結算彈窗
##   - auction_no_bid：灰色淡出 + 「😔 無人競標」提示
##   - auction_fish_killed：橙色閃光 + 「🎯 拍賣魚被擊破！競標結束！」提示
##   - control_shot：小閃光 + 浮動獎勵文字 + 計數器更新
##   - control_end：結算彈窗右側滑入
extends CanvasLayer

# 競標狀態
var _auction_active: bool = false
var _auction_timer: float = 0.0
var _auction_duration: float = 8.0
var _bid_base: int = 5
var _control_mult: float = 0.85
var _control_sec: int = 5

# 競標 UI 節點
var _auction_banner: Control = null
var _bid_bar: Control = null
var _bid_button: Control = null
var _timer_bar: Control = null
var _control_overlay: Control = null

# 控制權狀態
var _control_active: bool = false
var _control_shot_count: int = 0
var _control_total_reward: int = 0

func _ready() -> void:
	layer = 28  # 幸運拍賣魚面板層級

func _process(delta: float) -> void:
	if _auction_active:
		_auction_timer -= delta
		if _auction_timer <= 0.0:
			_auction_active = false
			_cleanup_auction_ui()
		else:
			_update_timer_bar()

## 處理幸運拍賣魚訊息
func handle_lucky_auction_fish(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"auction_start":
			_on_auction_start(payload)
		"auction_bid":
			_on_auction_bid(payload)
		"auction_result":
			_on_auction_result(payload)
		"auction_no_bid":
			_on_auction_no_bid()
		"auction_fish_killed":
			_on_auction_fish_killed(payload)
		"control_shot":
			_on_control_shot(payload)
		"control_end":
			_on_control_end(payload)

## 競標開始 — 金色雙閃光 + 頂部橫幅 + 計時條 + 出價按鈕
func _on_auction_start(payload: Dictionary) -> void:
	_auction_active = true
	_auction_duration = float(payload.get("duration_sec", 8))
	_auction_timer = _auction_duration
	_bid_base = payload.get("bid_base", 5)
	_control_mult = payload.get("control_mult", 0.85)
	_control_sec = payload.get("control_sec", 5)

	# 金色雙閃光
	_double_flash(Color("#FFD700"), 0.45)

	# 頂部橫幅
	var msg = "🏆 幸運拍賣魚出現！競標開始！出價最高者獲得 %d 秒大獎控制權（×%.2f 倍率）！" % [_control_sec, _control_mult]
	_auction_banner = _make_banner(msg, Color(0.1, 0.08, 0.0, 0.90), Color("#FFD700"))
	add_child(_auction_banner)

	# 底部計時條
	_timer_bar = _make_timer_bar()
	add_child(_timer_bar)

	# 出價按鈕（右下角）
	_bid_button = _make_bid_button()
	add_child(_bid_button)

## 有玩家出價 — 出價閃光 + 浮動文字
func _on_auction_bid(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "")
	var bid_amount: int = payload.get("bid_amount", 0)
	var top_bidder: String = payload.get("top_bidder", "")
	var top_bid: int = payload.get("top_bid_amount", 0)
	var total_bidders: int = payload.get("total_bidders", 0)

	# 小閃光
	_single_flash(Color("#FF8C00"), 0.25)

	# 浮動文字
	var text = "💰 %s 出價 %d！最高：%s（%d）" % [player_name, bid_amount, top_bidder, top_bid]
	_spawn_float_text(text, Color("#FFD700"), 28)

	# 更新橫幅文字
	if is_instance_valid(_auction_banner):
		var label = _auction_banner.get_node_or_null("Label")
		if label:
			label.text = "🏆 競標中！最高出價：%s（%d 金幣）| %d 人競標" % [top_bidder, top_bid, total_bidders]

## 競標結算 — 金色三次強閃光 + 大字 + 結算彈窗
func _on_auction_result(payload: Dictionary) -> void:
	_auction_active = false
	_cleanup_auction_ui()

	var winner_name: String = payload.get("winner_name", "")
	var winner_bid: int = payload.get("winner_bid", 0)
	var control_sec: int = payload.get("control_sec", 5)
	var control_mult: float = payload.get("control_mult", 0.85)
	var bid_results: Array = payload.get("bid_results", [])
	var refund_rate: float = payload.get("refund_rate", 0.5)

	# 金色三次強閃光
	_triple_flash(Color("#FFD700"))

	# 大字
	var big_label = _make_big_label(
		"🏆 %s 獲得大獎控制權！" % winner_name,
		Color("#FFD700"), 44
	)
	add_child(big_label)
	var tw = big_label.create_tween()
	tw.tween_interval(2.5)
	tw.tween_property(big_label, "modulate:a", 0.0, 0.4)
	tw.tween_callback(func(): if is_instance_valid(big_label): big_label.queue_free())

	# 結算彈窗（右側滑入）
	await get_tree().create_timer(0.5).timeout
	_show_result_popup(winner_name, winner_bid, control_sec, control_mult, bid_results, refund_rate)

	# 啟動控制權 UI
	_control_active = true
	_control_shot_count = 0
	_control_total_reward = 0
	_control_overlay = _make_control_overlay(winner_name, control_sec, control_mult)
	add_child(_control_overlay)

## 無人競標
func _on_auction_no_bid() -> void:
	_auction_active = false
	_cleanup_auction_ui()

	var label = _make_big_label("😔 無人競標，拍賣魚逃跑了...", Color("#AAAAAA"), 32)
	add_child(label)
	var tw = label.create_tween()
	tw.tween_interval(2.0)
	tw.tween_property(label, "modulate:a", 0.0, 0.5)
	tw.tween_callback(func(): if is_instance_valid(label): label.queue_free())

## 拍賣魚被擊破，競標提前結束
func _on_auction_fish_killed(payload: Dictionary) -> void:
	_auction_active = false
	_cleanup_auction_ui()

	var player_name: String = payload.get("player_name", "")
	_single_flash(Color("#FF8C00"), 0.3)

	var text = "🎯 %s 擊破拍賣魚！競標提前結束！" % player_name
	var label = _make_big_label(text, Color("#FF8C00"), 32)
	add_child(label)
	var tw = label.create_tween()
	tw.tween_interval(2.0)
	tw.tween_property(label, "modulate:a", 0.0, 0.4)
	tw.tween_callback(func(): if is_instance_valid(label): label.queue_free())

## 大獎控制權射擊 — 小閃光 + 浮動獎勵文字
func _on_control_shot(payload: Dictionary) -> void:
	var shot_reward: int = payload.get("shot_reward", 0)
	var shot_count: int = payload.get("shot_count", 0)
	var player_name: String = payload.get("player_name", "")

	_control_shot_count = shot_count
	_control_total_reward += shot_reward

	# 小閃光
	_single_flash(Color("#FFD700"), 0.15)

	# 浮動獎勵文字
	if shot_reward > 0:
		_spawn_float_text("🎯 +%d" % shot_reward, Color("#FFD700"), 32)

	# 更新控制權 overlay
	if is_instance_valid(_control_overlay):
		var label = _control_overlay.get_node_or_null("CountLabel")
		if label:
			label.text = "🎯 %s 控制中 | 擊破 %d 個 | +%d 金幣" % [player_name, shot_count, _control_total_reward]

## 大獎控制權結束 — 結算彈窗
func _on_control_end(payload: Dictionary) -> void:
	_control_active = false
	if is_instance_valid(_control_overlay):
		var tw = _control_overlay.create_tween()
		tw.tween_property(_control_overlay, "modulate:a", 0.0, 0.3)
		tw.tween_callback(func(): if is_instance_valid(_control_overlay): _control_overlay.queue_free())
		_control_overlay = null

	var player_name: String = payload.get("player_name", "")
	var shot_count: int = payload.get("shot_count", 0)
	var total_reward: int = payload.get("total_reward", 0)

	if shot_count > 0:
		_show_control_end_popup(player_name, shot_count, total_reward)

# ─── 內部 UI 工具函數 ───────────────────────────────────────────────────────

func _cleanup_auction_ui() -> void:
	if is_instance_valid(_auction_banner):
		_auction_banner.queue_free()
		_auction_banner = null
	if is_instance_valid(_timer_bar):
		_timer_bar.queue_free()
		_timer_bar = null
	if is_instance_valid(_bid_button):
		_bid_button.queue_free()
		_bid_button = null

func _update_timer_bar() -> void:
	if not is_instance_valid(_timer_bar):
		return
	var bar = _timer_bar.get_node_or_null("Bar")
	if not bar:
		return
	var pct = _auction_timer / _auction_duration
	bar.size.x = 480.0 * pct
	# 顏色漸變：金→橙→紅橙
	if pct > 0.6:
		bar.color = Color("#FFD700")
	elif pct > 0.3:
		bar.color = Color("#FF8C00")
	else:
		bar.color = Color("#FF4500")

func _make_banner(text: String, bg_color: Color, text_color: Color) -> Control:
	var panel = PanelContainer.new()
	panel.set_anchors_preset(Control.PRESET_TOP_WIDE)
	panel.position.y = 8
	panel.size.x = 600
	panel.position.x = (get_viewport().get_visible_rect().size.x - 600) / 2.0

	var style = StyleBoxFlat.new()
	style.bg_color = bg_color
	style.corner_radius_top_left = 8
	style.corner_radius_top_right = 8
	style.corner_radius_bottom_left = 8
	style.corner_radius_bottom_right = 8
	panel.add_theme_stylebox_override("panel", style)

	var label = Label.new()
	label.name = "Label"
	label.text = text
	label.add_theme_color_override("font_color", text_color)
	label.add_theme_font_size_override("font_size", 18)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.autowrap_mode = TextServer.AUTOWRAP_WORD_SMART
	panel.add_child(label)
	return panel

func _make_timer_bar() -> Control:
	var container = Control.new()
	container.set_anchors_preset(Control.PRESET_BOTTOM_WIDE)
	container.position.y = -28
	container.size = Vector2(get_viewport().get_visible_rect().size.x, 20)

	var bg = ColorRect.new()
	bg.color = Color(0.1, 0.1, 0.1, 0.7)
	bg.size = Vector2(container.size.x, 16)
	bg.position = Vector2(0, 2)
	container.add_child(bg)

	var bar = ColorRect.new()
	bar.name = "Bar"
	bar.color = Color("#FFD700")
	bar.size = Vector2(480.0, 12)
	bar.position = Vector2((container.size.x - 480.0) / 2.0, 4)
	container.add_child(bar)

	return container

func _make_bid_button() -> Control:
	var btn = Button.new()
	btn.text = "💰 出價競標"
	btn.set_anchors_preset(Control.PRESET_BOTTOM_RIGHT)
	btn.position = Vector2(-180, -80)
	btn.size = Vector2(160, 48)
	btn.add_theme_font_size_override("font_size", 20)
	btn.add_theme_color_override("font_color", Color("#FFD700"))
	btn.pressed.connect(func(): _on_bid_button_pressed())
	return btn

func _on_bid_button_pressed() -> void:
	# 發送出價訊息給 Server
	var gm = get_node_or_null("/root/GameManager")
	if gm and gm.has_method("send_message"):
		gm.send_message({"type": "lucky_auction_bid"})
	# 按鈕閃爍反饋
	if is_instance_valid(_bid_button):
		var tw = _bid_button.create_tween()
		tw.tween_property(_bid_button, "modulate", Color(2.0, 2.0, 0.5, 1.0), 0.08)
		tw.tween_property(_bid_button, "modulate", Color.WHITE, 0.12)

func _make_control_overlay(winner_name: String, control_sec: int, control_mult: float) -> Control:
	var container = Control.new()
	container.set_anchors_preset(Control.PRESET_TOP_WIDE)
	container.position.y = 60

	var label = Label.new()
	label.name = "CountLabel"
	label.text = "🎯 %s 大獎控制權（%d 秒 ×%.2f）" % [winner_name, control_sec, control_mult]
	label.add_theme_color_override("font_color", Color("#FFD700"))
	label.add_theme_font_size_override("font_size", 22)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.set_anchors_preset(Control.PRESET_TOP_WIDE)
	container.add_child(label)
	return container

func _show_result_popup(winner_name: String, winner_bid: int, control_sec: int,
		control_mult: float, bid_results: Array, refund_rate: float) -> void:
	var popup = PanelContainer.new()
	popup.set_anchors_preset(Control.PRESET_CENTER_RIGHT)
	popup.position.x = get_viewport().get_visible_rect().size.x + 10
	popup.size = Vector2(280, 200)

	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.08, 0.06, 0.0, 0.92)
	style.border_color = Color("#FFD700")
	style.border_width_left = 2
	style.border_width_right = 2
	style.border_width_top = 2
	style.border_width_bottom = 2
	style.corner_radius_top_left = 10
	style.corner_radius_top_right = 10
	style.corner_radius_bottom_left = 10
	style.corner_radius_bottom_right = 10
	popup.add_theme_stylebox_override("panel", style)

	var vbox = VBoxContainer.new()
	vbox.add_theme_constant_override("separation", 6)

	var title_label = Label.new()
	title_label.text = "🏆 競標結算"
	title_label.add_theme_color_override("font_color", Color("#FFD700"))
	title_label.add_theme_font_size_override("font_size", 22)
	title_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title_label)

	var winner_label = Label.new()
	winner_label.text = "🥇 %s\n出價：%d 金幣" % [winner_name, winner_bid]
	winner_label.add_theme_color_override("font_color", Color("#FFF8DC"))
	winner_label.add_theme_font_size_override("font_size", 18)
	winner_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(winner_label)

	var ctrl_label = Label.new()
	ctrl_label.text = "大獎控制權：%d 秒（×%.2f 倍率）" % [control_sec, control_mult]
	ctrl_label.add_theme_color_override("font_color", Color("#FF8C00"))
	ctrl_label.add_theme_font_size_override("font_size", 16)
	ctrl_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(ctrl_label)

	if bid_results.size() > 1:
		var refund_label = Label.new()
		refund_label.text = "失敗者退還 %.0f%% 出價" % (refund_rate * 100)
		refund_label.add_theme_color_override("font_color", Color("#AAAAAA"))
		refund_label.add_theme_font_size_override("font_size", 14)
		refund_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
		vbox.add_child(refund_label)

	popup.add_child(vbox)
	add_child(popup)

	# 右側滑入動畫
	var target_x = get_viewport().get_visible_rect().size.x - 300
	var tw = popup.create_tween()
	tw.tween_property(popup, "position:x", target_x, 0.3).set_ease(Tween.EASE_OUT)
	tw.tween_interval(4.0)
	tw.tween_property(popup, "position:x", get_viewport().get_visible_rect().size.x + 10, 0.3)
	tw.tween_callback(func(): if is_instance_valid(popup): popup.queue_free())

func _show_control_end_popup(player_name: String, shot_count: int, total_reward: int) -> void:
	var popup = PanelContainer.new()
	popup.set_anchors_preset(Control.PRESET_CENTER_RIGHT)
	popup.position.x = get_viewport().get_visible_rect().size.x + 10
	popup.size = Vector2(260, 140)

	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.08, 0.06, 0.0, 0.92)
	style.border_color = Color("#FF8C00")
	style.border_width_left = 2
	style.border_width_right = 2
	style.border_width_top = 2
	style.border_width_bottom = 2
	style.corner_radius_top_left = 10
	style.corner_radius_top_right = 10
	style.corner_radius_bottom_left = 10
	style.corner_radius_bottom_right = 10
	popup.add_theme_stylebox_override("panel", style)

	var vbox = VBoxContainer.new()
	vbox.add_theme_constant_override("separation", 8)

	var title = Label.new()
	title.text = "🎯 控制權結束"
	title.add_theme_color_override("font_color", Color("#FF8C00"))
	title.add_theme_font_size_override("font_size", 20)
	title.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title)

	var result = Label.new()
	result.text = "%s\n擊破 %d 個目標\n獲得 %d 金幣" % [player_name, shot_count, total_reward]
	result.add_theme_color_override("font_color", Color("#FFF8DC"))
	result.add_theme_font_size_override("font_size", 18)
	result.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(result)

	popup.add_child(vbox)
	add_child(popup)

	var target_x = get_viewport().get_visible_rect().size.x - 280
	var tw = popup.create_tween()
	tw.tween_property(popup, "position:x", target_x, 0.3).set_ease(Tween.EASE_OUT)
	tw.tween_interval(3.5)
	tw.tween_property(popup, "position:x", get_viewport().get_visible_rect().size.x + 10, 0.3)
	tw.tween_callback(func(): if is_instance_valid(popup): popup.queue_free())

func _make_big_label(text: String, color: Color, font_size: int) -> Label:
	var label = Label.new()
	label.text = text
	label.add_theme_color_override("font_color", color)
	label.add_theme_font_size_override("font_size", font_size)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.set_anchors_preset(Control.PRESET_CENTER)
	label.position.y = get_viewport().get_visible_rect().size.y * 0.35
	label.position.x = -300
	label.size.x = 600
	return label

func _spawn_float_text(text: String, color: Color, font_size: int) -> void:
	var label = Label.new()
	label.text = text
	label.add_theme_color_override("font_color", color)
	label.add_theme_font_size_override("font_size", font_size)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	var vp_size = get_viewport().get_visible_rect().size
	label.position = Vector2(vp_size.x / 2.0 - 200, vp_size.y * 0.45)
	label.size.x = 400
	add_child(label)
	var tw = label.create_tween()
	tw.tween_property(label, "position:y", label.position.y - 60, 1.2)
	tw.parallel().tween_property(label, "modulate:a", 0.0, 1.2)
	tw.tween_callback(func(): if is_instance_valid(label): label.queue_free())

func _double_flash(color: Color, duration: float) -> void:
	var overlay = ColorRect.new()
	overlay.color = Color(color.r, color.g, color.b, 0.0)
	overlay.set_anchors_preset(Control.PRESET_FULL_RECT)
	add_child(overlay)
	var tw = overlay.create_tween()
	tw.tween_property(overlay, "color:a", 0.55, duration * 0.3)
	tw.tween_property(overlay, "color:a", 0.0, duration * 0.3)
	tw.tween_property(overlay, "color:a", 0.45, duration * 0.2)
	tw.tween_property(overlay, "color:a", 0.0, duration * 0.2)
	tw.tween_callback(func(): if is_instance_valid(overlay): overlay.queue_free())

func _single_flash(color: Color, duration: float) -> void:
	var overlay = ColorRect.new()
	overlay.color = Color(color.r, color.g, color.b, 0.0)
	overlay.set_anchors_preset(Control.PRESET_FULL_RECT)
	add_child(overlay)
	var tw = overlay.create_tween()
	tw.tween_property(overlay, "color:a", 0.4, duration * 0.4)
	tw.tween_property(overlay, "color:a", 0.0, duration * 0.6)
	tw.tween_callback(func(): if is_instance_valid(overlay): overlay.queue_free())

func _triple_flash(color: Color) -> void:
	var overlay = ColorRect.new()
	overlay.color = Color(color.r, color.g, color.b, 0.0)
	overlay.set_anchors_preset(Control.PRESET_FULL_RECT)
	add_child(overlay)
	var tw = overlay.create_tween()
	tw.tween_property(overlay, "color:a", 0.7, 0.10)
	tw.tween_property(overlay, "color:a", 0.0, 0.12)
	tw.tween_property(overlay, "color:a", 0.6, 0.10)
	tw.tween_property(overlay, "color:a", 0.0, 0.12)
	tw.tween_property(overlay, "color:a", 0.5, 0.10)
	tw.tween_property(overlay, "color:a", 0.0, 0.15)
	tw.tween_callback(func(): if is_instance_valid(overlay): overlay.queue_free())
