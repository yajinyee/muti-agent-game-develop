## LuckyBetFishPanel.gd — 幸運賭注魚系統面板（DAY-240）
## 業界原創「玩家主動風險決策+賭注翻倍」機制
##
## 數值設計（三選項期望值相同 2.0x，差異只在方差）：
##   - 選擇 A（保守）：×2.0，100% 成功 → 期望值 2.0x
##   - 選擇 B（激進）：×4.0，50% 成功，失敗歸零 → 期望值 2.0x
##   - 選擇 C（瘋狂）：×8.0，25% 成功，失敗歸零 → 期望值 2.0x
##
## 視覺設計：
##   - 紫金賭注主題（#9B59B6 + #F39C12 + #FFD700 + #F5EEF8）
##   - bet_start：紫色雙閃光 + 賭注選擇介面（三個按鈕 A/B/C）+ 倒數計時條
##   - bet_broadcast：頂部小橫幅（通知全服有人觸發）
##   - bet_decided：結果閃光 + 大字顯示結果（成功=金色/失敗=灰色）
##   - bet_timeout：橙色提示「超時自動選擇 A」
extends CanvasLayer

# 主題顏色
const COLOR_PURPLE   = Color("#9B59B6")  # 紫色（主題）
const COLOR_GOLD     = Color("#F39C12")  # 金橙（激進）
const COLOR_CRAZY    = Color("#E74C3C")  # 紅色（瘋狂）
const COLOR_SAFE     = Color("#27AE60")  # 綠色（保守）
const COLOR_FAIL     = Color("#7F8C8D")  # 灰色（失敗）
const COLOR_WHITE    = Color("#FFFFFF")  # 白色
const COLOR_PALE     = Color("#F5EEF8")  # 淡紫

# 賭注選擇介面節點
var _bet_panel: Control = null
var _countdown_bar: ColorRect = null
var _countdown_bar_bg: ColorRect = null
var _countdown_tween: Tween = null
var _current_instance_id: String = ""

# 是否為本玩家的賭注
var _is_my_bet: bool = false

func _ready() -> void:
	layer = 5  # 幸運賭注魚面板層級

## 處理幸運賭注魚訊息
func handle_lucky_bet_fish(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"bet_start":
			_on_bet_start(payload)
		"bet_broadcast":
			_on_bet_broadcast(payload)
		"bet_decided":
			_on_bet_decided(payload)
		"bet_timeout":
			_on_bet_timeout(payload)

## bet_start — 觸發賭注選擇（個人訊息）
func _on_bet_start(payload: Dictionary) -> void:
	_is_my_bet = true
	_current_instance_id = payload.get("instance_id", "")
	var decision_sec: int = payload.get("decision_sec", 10)
	var option_a: Dictionary = payload.get("option_a", {})
	var option_b: Dictionary = payload.get("option_b", {})
	var option_c: Dictionary = payload.get("option_c", {})

	var vp_size = get_viewport().size

	# 紫色雙閃光
	_flash_screen(COLOR_PURPLE, 0.13)
	await get_tree().create_timer(0.09).timeout
	_flash_screen(COLOR_PALE, 0.10)

	# 建立賭注選擇面板
	_bet_panel = Control.new()
	_bet_panel.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	add_child(_bet_panel)

	# 半透明背景
	var bg = ColorRect.new()
	bg.color = Color(0.05, 0.0, 0.1, 0.72)
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_bet_panel.add_child(bg)

	# 標題
	var title = Label.new()
	title.text = "🎲 幸運賭注魚！選擇你的賭注"
	title.add_theme_font_size_override("font_size", 28)
	title.add_theme_color_override("font_color", COLOR_PALE)
	title.position = Vector2(vp_size.x / 2 - 160, vp_size.y / 2 - 130)
	_bet_panel.add_child(title)

	# 副標題
	var subtitle = Label.new()
	subtitle.text = "下一次擊破目標時套用選擇的倍率"
	subtitle.add_theme_font_size_override("font_size", 14)
	subtitle.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8))
	subtitle.position = Vector2(vp_size.x / 2 - 130, vp_size.y / 2 - 95)
	_bet_panel.add_child(subtitle)

	# 三個選項按鈕
	var btn_y = vp_size.y / 2 - 60
	var btn_x_start = vp_size.x / 2 - 210

	_create_bet_button(_bet_panel, "A", option_a, Vector2(btn_x_start, btn_y), COLOR_SAFE)
	_create_bet_button(_bet_panel, "B", option_b, Vector2(btn_x_start + 140, btn_y), COLOR_GOLD)
	_create_bet_button(_bet_panel, "C", option_c, Vector2(btn_x_start + 280, btn_y), COLOR_CRAZY)

	# 倒數計時條背景
	_countdown_bar_bg = ColorRect.new()
	_countdown_bar_bg.color = Color(0.2, 0.1, 0.3, 0.8)
	_countdown_bar_bg.size = Vector2(400, 12)
	_countdown_bar_bg.position = Vector2(vp_size.x / 2 - 200, vp_size.y / 2 + 80)
	_bet_panel.add_child(_countdown_bar_bg)

	# 倒數計時條
	_countdown_bar = ColorRect.new()
	_countdown_bar.color = COLOR_PURPLE
	_countdown_bar.size = Vector2(400, 12)
	_countdown_bar.position = Vector2(vp_size.x / 2 - 200, vp_size.y / 2 + 80)
	_bet_panel.add_child(_countdown_bar)

	# 倒數計時條動畫
	_countdown_tween = _countdown_bar.create_tween()
	_countdown_tween.tween_property(_countdown_bar, "size:x", 0.0, float(decision_sec))

	# 倒數文字
	var countdown_label = Label.new()
	countdown_label.name = "CountdownLabel"
	countdown_label.text = "%d 秒後自動選擇 A" % decision_sec
	countdown_label.add_theme_font_size_override("font_size", 12)
	countdown_label.add_theme_color_override("font_color", Color(0.7, 0.7, 0.7))
	countdown_label.position = Vector2(vp_size.x / 2 - 80, vp_size.y / 2 + 96)
	_bet_panel.add_child(countdown_label)

	# 倒數計時更新
	var elapsed = 0.0
	while elapsed < float(decision_sec) and is_instance_valid(_bet_panel):
		await get_tree().create_timer(1.0).timeout
		elapsed += 1.0
		var remaining = decision_sec - int(elapsed)
		if is_instance_valid(countdown_label):
			countdown_label.text = "%d 秒後自動選擇 A" % remaining

## 建立賭注選項按鈕
func _create_bet_button(parent: Control, choice: String, option: Dictionary, pos: Vector2, color: Color) -> void:
	var label_text = option.get("label", choice)
	var mult = option.get("mult", 1.0)
	var chance = option.get("success_chance", 1.0)
	var fail_mult = option.get("fail_mult", 1.0)

	# 按鈕容器
	var container = Control.new()
	container.position = pos
	container.size = Vector2(120, 120)
	parent.add_child(container)

	# 按鈕背景
	var btn_bg = ColorRect.new()
	btn_bg.color = Color(color.r * 0.3, color.g * 0.3, color.b * 0.3, 0.9)
	btn_bg.size = Vector2(120, 120)
	btn_bg.position = Vector2.ZERO
	container.add_child(btn_bg)

	# 按鈕邊框（用 Label 模擬）
	var border = ColorRect.new()
	border.color = color
	border.size = Vector2(120, 3)
	border.position = Vector2(0, 0)
	container.add_child(border)

	var border_b = ColorRect.new()
	border_b.color = color
	border_b.size = Vector2(120, 3)
	border_b.position = Vector2(0, 117)
	container.add_child(border_b)

	# 選項標籤
	var choice_label = Label.new()
	choice_label.text = choice
	choice_label.add_theme_font_size_override("font_size", 32)
	choice_label.add_theme_color_override("font_color", color)
	choice_label.position = Vector2(48, 8)
	container.add_child(choice_label)

	# 選項名稱
	var name_label = Label.new()
	name_label.text = label_text
	name_label.add_theme_font_size_override("font_size", 14)
	name_label.add_theme_color_override("font_color", COLOR_WHITE)
	name_label.position = Vector2(10, 48)
	container.add_child(name_label)

	# 倍率
	var mult_label = Label.new()
	mult_label.text = "×%.1f" % mult
	mult_label.add_theme_font_size_override("font_size", 18)
	mult_label.add_theme_color_override("font_color", Color("#FFD700"))
	mult_label.position = Vector2(10, 66)
	container.add_child(mult_label)

	# 成功機率
	var chance_label = Label.new()
	chance_label.text = "%.0f%% 成功" % (chance * 100)
	chance_label.add_theme_font_size_override("font_size", 11)
	chance_label.add_theme_color_override("font_color", Color(0.8, 0.9, 0.8))
	chance_label.position = Vector2(10, 90)
	container.add_child(chance_label)

	# 失敗倍率（只有 B/C 顯示）
	if choice != "A":
		var fail_label = Label.new()
		fail_label.text = "失敗 ×0（歸零）" if fail_mult == 0.0 else "失敗 ×%.1f" % fail_mult
		fail_label.add_theme_font_size_override("font_size", 10)
		fail_label.add_theme_color_override("font_color", Color(0.8, 0.3, 0.3))
		fail_label.position = Vector2(10, 104)
		container.add_child(fail_label)

	# 點擊區域（Button）
	var btn = Button.new()
	btn.flat = true
	btn.size = Vector2(120, 120)
	btn.position = Vector2.ZERO
	btn.modulate = Color(1, 1, 1, 0)  # 透明，只用於點擊
	container.add_child(btn)

	# 懸停效果
	btn.mouse_entered.connect(func():
		var tw = container.create_tween()
		tw.tween_property(btn_bg, "color", Color(color.r * 0.5, color.g * 0.5, color.b * 0.5, 0.95), 0.1)
	)
	btn.mouse_exited.connect(func():
		var tw = container.create_tween()
		tw.tween_property(btn_bg, "color", Color(color.r * 0.3, color.g * 0.3, color.b * 0.3, 0.9), 0.1)
	)

	# 點擊事件
	btn.pressed.connect(func():
		_on_choice_selected(choice)
	)

## 玩家選擇賭注
func _on_choice_selected(choice: String) -> void:
	if not _is_my_bet or _current_instance_id.is_empty():
		return

	# 關閉選擇介面
	_close_bet_panel()

	# 發送選擇到 Server
	var gm = get_node_or_null("/root/Main/GameManager")
	if gm and gm.has_method("send_bet_choice"):
		gm.send_bet_choice(choice, _current_instance_id)
	else:
		# 直接透過 NetworkManager 發送
		var nm = get_node_or_null("/root/Main/NetworkManager")
		if nm and nm.has_method("send"):
			nm.send({
				"type": "lucky_bet_choice",
				"payload": {
					"choice": choice,
					"instance_id": _current_instance_id
				}
			})

	_is_my_bet = false
	_current_instance_id = ""

## 關閉賭注選擇介面
func _close_bet_panel() -> void:
	if _countdown_tween and _countdown_tween.is_valid():
		_countdown_tween.kill()
	if is_instance_valid(_bet_panel):
		var tw = _bet_panel.create_tween()
		tw.tween_property(_bet_panel, "modulate:a", 0.0, 0.2)
		tw.tween_callback(_bet_panel.queue_free)
		_bet_panel = null

## bet_broadcast — 通知全服有人觸發賭注魚
func _on_bet_broadcast(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "")
	var vp_size = get_viewport().size

	# 頂部小橫幅
	var banner = Label.new()
	banner.text = "🎲 %s 正在選擇賭注..." % player_name
	banner.add_theme_font_size_override("font_size", 13)
	banner.add_theme_color_override("font_color", COLOR_PALE)
	banner.position = Vector2(vp_size.x / 2 - 100, 8)
	add_child(banner)

	var tw = banner.create_tween()
	tw.tween_interval(3.0)
	tw.tween_property(banner, "modulate:a", 0.0, 0.5)
	tw.tween_callback(banner.queue_free)

## bet_decided — 玩家決策結果
func _on_bet_decided(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "")
	var choice: String = payload.get("choice", "A")
	var choice_label: String = payload.get("choice_label", "")
	var success: bool = payload.get("success", true)
	var result_mult: float = payload.get("result_mult", 1.0)

	# 關閉選擇介面（如果還開著）
	if _is_my_bet:
		_close_bet_panel()
		_is_my_bet = false

	var vp_size = get_viewport().size

	# 結果閃光
	if success:
		var flash_color = {
			"A": COLOR_SAFE,
			"B": COLOR_GOLD,
			"C": COLOR_CRAZY
		}.get(choice, COLOR_PURPLE)
		_flash_screen(flash_color, 0.15)
		await get_tree().create_timer(0.1).timeout
		_flash_screen(COLOR_WHITE, 0.10)
	else:
		_flash_screen(COLOR_FAIL, 0.12)

	# 結果大字
	var result_text: String
	var result_color: Color
	if success:
		result_text = "🎲 %s 賭注成功！×%.1f" % [player_name, result_mult]
		result_color = {
			"A": COLOR_SAFE,
			"B": COLOR_GOLD,
			"C": COLOR_CRAZY
		}.get(choice, COLOR_PURPLE)
	else:
		result_text = "🎲 %s 賭注失敗... ×%.1f" % [player_name, result_mult]
		result_color = COLOR_FAIL

	var big_label = Label.new()
	big_label.text = result_text
	big_label.add_theme_font_size_override("font_size", 32)
	big_label.add_theme_color_override("font_color", result_color)
	big_label.position = Vector2(vp_size.x / 2 - 160, vp_size.y / 2 - 20)
	add_child(big_label)

	# 縮放動畫
	big_label.scale = Vector2(0.8, 0.8)
	var tw = big_label.create_tween()
	tw.tween_property(big_label, "scale", Vector2(1.1, 1.1), 0.12)
	tw.tween_property(big_label, "scale", Vector2(1.0, 1.0), 0.08)
	tw.tween_interval(2.0)
	tw.tween_property(big_label, "modulate:a", 0.0, 0.5)
	tw.tween_callback(big_label.queue_free)

	# 選擇標籤
	var choice_info = Label.new()
	choice_info.text = "選擇：%s" % choice_label
	choice_info.add_theme_font_size_override("font_size", 14)
	choice_info.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8))
	choice_info.position = Vector2(vp_size.x / 2 - 80, vp_size.y / 2 + 20)
	add_child(choice_info)

	var tw2 = choice_info.create_tween()
	tw2.tween_interval(2.5)
	tw2.tween_property(choice_info, "modulate:a", 0.0, 0.4)
	tw2.tween_callback(choice_info.queue_free)

## bet_timeout — 超時自動選擇 A
func _on_bet_timeout(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "")

	# 關閉選擇介面
	if _is_my_bet:
		_close_bet_panel()
		_is_my_bet = false

	var vp_size = get_viewport().size

	# 橙色提示
	_flash_screen(Color("#E67E22"), 0.10)

	var timeout_label = Label.new()
	timeout_label.text = "⏰ %s 超時，自動選擇保守 ×2.0" % player_name
	timeout_label.add_theme_font_size_override("font_size", 16)
	timeout_label.add_theme_color_override("font_color", Color("#FAD7A0"))
	timeout_label.position = Vector2(vp_size.x / 2 - 160, vp_size.y / 2 - 10)
	add_child(timeout_label)

	var tw = timeout_label.create_tween()
	tw.tween_interval(2.0)
	tw.tween_property(timeout_label, "modulate:a", 0.0, 0.4)
	tw.tween_callback(timeout_label.queue_free)

## 全螢幕閃光效果
func _flash_screen(color: Color, duration: float) -> void:
	var vp_size = get_viewport().size
	var flash = ColorRect.new()
	flash.color = Color(color.r, color.g, color.b, 0.45)
	flash.size = vp_size
	flash.position = Vector2.ZERO
	add_child(flash)

	var tw = flash.create_tween()
	tw.tween_property(flash, "modulate:a", 0.0, duration)
	tw.tween_callback(flash.queue_free)
