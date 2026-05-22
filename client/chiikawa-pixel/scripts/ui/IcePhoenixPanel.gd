## IcePhoenixPanel.gd — 冰鳳凰覺醒 BOSS 面板（DAY-200）
## 業界依據：Royal Fishing JILI「Ice Phoenix Awaken Feature — awards up to 300x.
## Multicoloured phoenix (blue, pink, purple, orange) with magical aura.
## Awaken Boss with 30x basic multiplier. Power Up attack delivers 6x-10x boost.」
##
## 視覺設計：
##   - 冰藍紫主題（#00BFFF + #9400D3 + #FF69B4 + #FFFFFF）
##   - awaken_start：全螢幕冰藍三次閃光 + 頂部橫幅「❄️ 冰鳳凰覺醒！」+ 基礎獎勵顯示
##   - power_up_shot：冰藍閃光 + 目標位置「❄️ ×N倍」浮動文字 + 進度條更新
##   - frost_burst_start：全螢幕冰白強閃光 + 「❄️💥 冰霜爆發！」大字
##   - frost_burst_result：冰霜擊破數/獎勵顯示
##   - awaken_result：右側滑入結算彈窗（基礎/Power Up/冰霜/總計）
extends CanvasLayer

var _panel: Control
var _banner: Label
var _power_up_bar: ColorRect     # Power Up 進度條
var _power_up_bar_bg: ColorRect
var _power_up_counter: Label     # Power Up 計數器
var _result_popup: Control

var _total_shots: int = 0
var _current_shot: int = 0

func _ready() -> void:
	layer = 45
	_build_ui()

func _build_ui() -> void:
	_panel = Control.new()
	_panel.set_anchors_preset(Control.PRESET_FULL_RECT)
	_panel.visible = false
	_panel.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_panel)

	# 頂部橫幅
	_banner = Label.new()
	_banner.text = "❄️ 冰鳳凰覺醒！"
	_banner.add_theme_font_size_override("font_size", 26)
	_banner.add_theme_color_override("font_color", Color("#00BFFF"))
	_banner.set_anchors_preset(Control.PRESET_TOP_WIDE)
	_banner.position = Vector2(0, 12)
	_banner.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner.visible = false
	add_child(_banner)

	# Power Up 計數器
	_power_up_counter = Label.new()
	_power_up_counter.text = "Power Up: 0 / 0"
	_power_up_counter.add_theme_font_size_override("font_size", 18)
	_power_up_counter.add_theme_color_override("font_color", Color("#9400D3"))
	_power_up_counter.set_anchors_preset(Control.PRESET_TOP_WIDE)
	_power_up_counter.position = Vector2(0, 44)
	_power_up_counter.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_power_up_counter.visible = false
	add_child(_power_up_counter)

	# Power Up 進度條背景
	_power_up_bar_bg = ColorRect.new()
	_power_up_bar_bg.color = Color(0.05, 0.05, 0.15, 0.7)
	_power_up_bar_bg.set_anchors_preset(Control.PRESET_TOP_WIDE)
	_power_up_bar_bg.size = Vector2(1280, 10)
	_power_up_bar_bg.position = Vector2(0, 70)
	_power_up_bar_bg.visible = false
	add_child(_power_up_bar_bg)

	# Power Up 進度條
	_power_up_bar = ColorRect.new()
	_power_up_bar.color = Color("#00BFFF")
	_power_up_bar.set_anchors_preset(Control.PRESET_TOP_LEFT)
	_power_up_bar.size = Vector2(0, 10)
	_power_up_bar.position = Vector2(0, 70)
	_power_up_bar.visible = false
	add_child(_power_up_bar)

	# 結算彈窗
	_result_popup = Control.new()
	_result_popup.visible = false
	_result_popup.set_anchors_preset(Control.PRESET_CENTER_RIGHT)
	_result_popup.size = Vector2(320, 240)
	_result_popup.position = Vector2(-340, -120)
	add_child(_result_popup)

	var popup_bg = ColorRect.new()
	popup_bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	popup_bg.color = Color(0.02, 0.05, 0.15, 0.95)
	_result_popup.add_child(popup_bg)

	var popup_label = Label.new()
	popup_label.name = "ResultLabel"
	popup_label.add_theme_font_size_override("font_size", 15)
	popup_label.add_theme_color_override("font_color", Color("#00BFFF"))
	popup_label.set_anchors_preset(Control.PRESET_FULL_RECT)
	popup_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	popup_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_result_popup.add_child(popup_label)

## 處理冰鳳凰訊息
func handle_ice_phoenix(payload: Dictionary) -> void:
	var event = payload.get("event", "")
	match event:
		"awaken_start":
			_on_awaken_start(payload)
		"power_up_shot":
			_on_power_up_shot(payload)
		"frost_burst_start":
			_on_frost_burst_start()
		"frost_burst_result":
			_on_frost_burst_result(payload)
		"awaken_result":
			_on_awaken_result(payload)

## 覺醒開始
func _on_awaken_start(payload: Dictionary) -> void:
	var killer_name = payload.get("killer_name", "玩家")
	var base_reward = payload.get("base_reward", 0)

	_panel.visible = true

	# 全螢幕冰藍三次閃光
	_flash_screen(Color("#00BFFF"), 0.6)
	await get_tree().create_timer(0.15).timeout
	_flash_screen(Color("#9400D3"), 0.4)
	await get_tree().create_timer(0.15).timeout
	_flash_screen(Color("#00BFFF"), 0.35)

	# 顯示橫幅
	_banner.text = "❄️ " + killer_name + " 觸發冰鳳凰覺醒！"
	_banner.visible = true

	# 顯示基礎獎勵
	_spawn_float_text("❄️ 基礎獎勵 +" + str(base_reward) + " 金幣", Color("#00BFFF"), Vector2(640, 300))

## Power Up 攻擊
func _on_power_up_shot(payload: Dictionary) -> void:
	var shot_index = payload.get("shot_index", 1)
	var total_shots = payload.get("total_shots", 5)
	var result = payload.get("power_up_result", {})
	var killed = result.get("killed", false)
	var mult = result.get("mult", 6.0)
	var reward = result.get("reward", 0)

	_total_shots = total_shots
	_current_shot = shot_index

	# 更新計數器
	_power_up_counter.text = "Power Up: %d / %d" % [shot_index, total_shots]
	_power_up_counter.visible = true

	# 更新進度條
	_power_up_bar_bg.visible = true
	_power_up_bar.visible = true
	var ratio = float(shot_index) / float(total_shots)
	var tween = create_tween()
	tween.tween_property(_power_up_bar, "size:x", 1280.0 * ratio, 0.3)

	if killed:
		# 冰藍閃光
		_flash_screen(Color("#00BFFF"), 0.3)
		# 浮動文字
		_spawn_float_text("❄️ ×%.1f +%d" % [mult, reward], Color("#FF69B4"))
	else:
		# 未擊破，淡藍色
		_spawn_float_text("❄️ 閃避！", Color("#87CEEB"))

## 冰霜爆發開始
func _on_frost_burst_start() -> void:
	# 全螢幕冰白強閃光
	_flash_screen(Color("#FFFFFF"), 0.9)
	await get_tree().create_timer(0.1).timeout
	_flash_screen(Color("#00BFFF"), 0.7)

	# 大字提示
	var frost_label = Label.new()
	frost_label.text = "❄️💥 冰霜爆發！"
	frost_label.add_theme_font_size_override("font_size", 42)
	frost_label.add_theme_color_override("font_color", Color("#00BFFF"))
	frost_label.set_anchors_preset(Control.PRESET_CENTER)
	frost_label.position = Vector2(-200, -30)
	add_child(frost_label)
	var tween = create_tween()
	tween.tween_property(frost_label, "modulate:a", 0.0, 1.5)
	tween.tween_callback(frost_label.queue_free)

## 冰霜爆發結果
func _on_frost_burst_result(payload: Dictionary) -> void:
	var frost_kills = payload.get("frost_kills", 0)
	var frost_reward = payload.get("frost_reward", 0)
	if frost_kills > 0:
		_spawn_float_text("❄️ 冰霜擊破 %d 個！+%d" % [frost_kills, frost_reward], Color("#00BFFF"), Vector2(640, 250))

## 覺醒結算
func _on_awaken_result(payload: Dictionary) -> void:
	var killer_name = payload.get("killer_name", "玩家")
	var base_reward = payload.get("base_reward", 0)
	var power_up_kills = payload.get("power_up_kills", 0)
	var power_up_reward = payload.get("power_up_reward", 0)
	var frost_kills = payload.get("frost_kills", 0)
	var frost_reward = payload.get("frost_reward", 0)
	var total_reward = payload.get("total_reward", 0)
	var has_frost = payload.get("has_frost", false)

	# 隱藏橫幅和進度條
	_banner.visible = false
	_power_up_counter.visible = false
	_power_up_bar.visible = false
	_power_up_bar_bg.visible = false

	# 依總獎勵決定閃光
	if total_reward >= 200 * 10:  # 假設 betLevel=10
		_flash_screen(Color("#FFD700"), 0.8)
	elif total_reward >= 100 * 10:
		_flash_screen(Color("#00BFFF"), 0.6)

	# 顯示結算彈窗
	var result_label = _result_popup.get_node("ResultLabel")
	var frost_text = ""
	if has_frost:
		frost_text = "\n❄️ 冰霜爆發: %d 個 +%d" % [frost_kills, frost_reward]
	result_label.text = (
		"❄️ 冰鳳凰覺醒結算\n\n"
		+ "基礎獎勵: +" + str(base_reward) + "\n"
		+ "Power Up: %d 個 +%d" % [power_up_kills, power_up_reward]
		+ frost_text
		+ "\n\n總獎勵: " + str(total_reward) + " 金幣"
	)
	_result_popup.visible = true

	# 5 秒後淡出
	await get_tree().create_timer(5.0).timeout
	var tween = create_tween()
	tween.tween_property(_result_popup, "modulate:a", 0.0, 0.5)
	await tween.finished
	_result_popup.visible = false
	_result_popup.modulate.a = 1.0
	_panel.visible = false

## 全螢幕閃光
func _flash_screen(color: Color, alpha: float) -> void:
	var flash = ColorRect.new()
	flash.set_anchors_preset(Control.PRESET_FULL_RECT)
	flash.color = Color(color.r, color.g, color.b, alpha)
	flash.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(flash)
	var tween = create_tween()
	tween.tween_property(flash, "modulate:a", 0.0, 0.3)
	tween.tween_callback(flash.queue_free)

## 浮動文字
func _spawn_float_text(text: String, color: Color, pos: Vector2 = Vector2(-1, -1)) -> void:
	var label = Label.new()
	label.text = text
	label.add_theme_font_size_override("font_size", 20)
	label.add_theme_color_override("font_color", color)
	if pos.x < 0:
		label.position = Vector2(randf_range(300, 900), randf_range(200, 500))
	else:
		label.position = pos + Vector2(-100, 0)
	add_child(label)
	var tween = create_tween()
	tween.tween_property(label, "position:y", label.position.y - 70, 1.2)
	tween.parallel().tween_property(label, "modulate:a", 0.0, 1.2)
	tween.tween_callback(label.queue_free)
