## RainbowPrismPanel.gd — 彩虹稜鏡魚系統面板（DAY-213）
## 業界依據：Dive Down 2026「Rainbow is the strongest mutation with 3.0x multiplier」
## 業界原創「稜鏡折射染色」機制
##
## 視覺設計：
##   - 彩虹主題（#FF4444 紅 / #FF8C00 橙 / #FFD700 黃 / #00CC44 綠 / #0088FF 藍）
##   - prism_start：全螢幕彩虹三次強閃光 + 頂部橫幅 + 染色目標標記（彩色光環）
##   - prism_blast：全螢幕彩虹爆炸閃光 + 「🌈 彩虹爆炸！」大字 + 結算彈窗右側滑入
extends CanvasLayer

# 顏色對應表（color_name → Color）
const COLOR_MAP: Dictionary = {
	"red":    Color("#FF4444"),
	"orange": Color("#FF8C00"),
	"yellow": Color("#FFD700"),
	"green":  Color("#00CC44"),
	"blue":   Color("#0088FF"),
}

# 顏色中文名稱
const COLOR_NAMES: Dictionary = {
	"red":    "紅色 ×1.5",
	"orange": "橙色 ×2.0",
	"yellow": "黃色 ×2.5",
	"green":  "綠色 ×3.0",
	"blue":   "藍色 ×5.0",
}

# 染色目標標記（target_id → Control）
var _color_markers: Dictionary = {}

func _ready() -> void:
	layer = 32  # 彩虹稜鏡面板層級

## 處理彩虹稜鏡魚訊息
func handle_rainbow_prism(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"prism_start":
			_on_prism_start(payload)
		"prism_blast":
			_on_prism_blast(payload)

## 稜鏡折射開始 — 全螢幕彩虹三次強閃光 + 頂部橫幅 + 染色目標標記
func _on_prism_start(payload: Dictionary) -> void:
	var trigger_player: String = payload.get("trigger_player", "玩家")
	var colored_targets: Array = payload.get("colored_targets", [])
	var duration: int = payload.get("duration", 10)

	# 全螢幕彩虹三次強閃光（紅→綠→藍）
	_rainbow_triple_flash()

	# 頂部橫幅
	var count = colored_targets.size()
	var banner = _make_banner(
		"🌈 %s 觸發彩虹稜鏡！%d 個目標被染色！快打高倍率顏色！" % [trigger_player, count],
		Color(0.05, 0.0, 0.1, 0.88),
		Color("#FF69B4")
	)
	add_child(banner)
	var tween_b = create_tween()
	tween_b.tween_property(banner, "position:y", 0.0, 0.25)
	tween_b.tween_interval(3.5)
	tween_b.tween_property(banner, "modulate:a", 0.0, 0.4)
	tween_b.tween_callback(banner.queue_free)

	# 顯示顏色圖例（右側小面板）
	_show_color_legend(colored_targets, duration)

	# 10 秒後自動清除標記
	var timer = get_tree().create_timer(float(duration) + 0.5)
	timer.timeout.connect(_clear_color_markers)

## 彩虹爆炸結算 — 全螢幕彩虹爆炸閃光 + 大字 + 結算彈窗
func _on_prism_blast(payload: Dictionary) -> void:
	var blast_kills: int = payload.get("blast_kills", 0)
	var blast_reward: int = payload.get("blast_reward", 0)

	# 清除所有染色標記
	_clear_color_markers()

	# 全螢幕彩虹爆炸閃光
	_rainbow_triple_flash()

	# 「🌈 彩虹爆炸！」大字
	var big_label = Label.new()
	big_label.text = "🌈 彩虹爆炸！"
	big_label.add_theme_font_size_override("font_size", 52)
	big_label.add_theme_color_override("font_color", Color("#FF69B4"))
	big_label.position = Vector2(640 - 180, 280)
	big_label.size = Vector2(360, 70)
	big_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	add_child(big_label)
	var tween_big = create_tween()
	tween_big.tween_property(big_label, "scale", Vector2(1.2, 1.2), 0.1)
	tween_big.tween_property(big_label, "scale", Vector2(1.0, 1.0), 0.1)
	tween_big.tween_interval(1.2)
	tween_big.tween_property(big_label, "modulate:a", 0.0, 0.4)
	tween_big.tween_callback(big_label.queue_free)

	# 結算彈窗（右側滑入）
	if blast_kills > 0 or blast_reward > 0:
		_show_result_popup(blast_kills, blast_reward)

## 顯示顏色圖例（右側小面板，顯示各顏色倍率）
func _show_color_legend(colored_targets: Array, duration: int) -> void:
	if colored_targets.is_empty():
		return

	var panel = PanelContainer.new()
	panel.position = Vector2(1280, 120)
	panel.size = Vector2(200, 180)

	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.05, 0.0, 0.1, 0.88)
	style.border_color = Color("#FF69B4")
	style.border_width_left = 2
	style.border_width_right = 2
	style.border_width_top = 2
	style.border_width_bottom = 2
	style.corner_radius_top_left = 8
	style.corner_radius_top_right = 8
	style.corner_radius_bottom_left = 8
	style.corner_radius_bottom_right = 8
	panel.add_theme_stylebox_override("panel", style)

	var vbox = VBoxContainer.new()
	vbox.add_theme_constant_override("separation", 4)
	panel.add_child(vbox)

	var title = Label.new()
	title.text = "🌈 稜鏡染色"
	title.add_theme_font_size_override("font_size", 14)
	title.add_theme_color_override("font_color", Color("#FF69B4"))
	title.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title)

	# 顯示每個染色目標的顏色和倍率
	for info in colored_targets:
		var color_name: String = info.get("color_name", "red")
		var mult: float = info.get("mult_bonus", 1.5)
		var color_hex: String = info.get("color_hex", "#FF4444")

		var row = Label.new()
		row.text = "● %s ×%.1f" % [_get_color_cn(color_name), mult]
		row.add_theme_font_size_override("font_size", 13)
		row.add_theme_color_override("font_color", Color(color_hex))
		row.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
		vbox.add_child(row)

	add_child(panel)

	# 右側滑入動畫
	var tween = create_tween()
	tween.tween_property(panel, "position:x", 1060.0, 0.3).set_ease(Tween.EASE_OUT)
	tween.tween_interval(float(duration) - 0.5)
	tween.tween_property(panel, "modulate:a", 0.0, 0.5)
	tween.tween_callback(panel.queue_free)

## 顯示結算彈窗（右側滑入）
func _show_result_popup(blast_kills: int, blast_reward: int) -> void:
	var popup = PanelContainer.new()
	popup.position = Vector2(1280, 200)
	popup.size = Vector2(280, 140)

	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.05, 0.0, 0.1, 0.92)
	style.border_color = Color("#FF69B4")
	style.border_width_left = 2
	style.border_width_right = 2
	style.border_width_top = 2
	style.border_width_bottom = 2
	style.corner_radius_top_left = 8
	style.corner_radius_top_right = 8
	style.corner_radius_bottom_left = 8
	style.corner_radius_bottom_right = 8
	popup.add_theme_stylebox_override("panel", style)

	var vbox = VBoxContainer.new()
	vbox.add_theme_constant_override("separation", 6)
	popup.add_child(vbox)

	var title = Label.new()
	title.text = "🌈 彩虹爆炸結算"
	title.add_theme_font_size_override("font_size", 16)
	title.add_theme_color_override("font_color", Color("#FF69B4"))
	title.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title)

	var sep = HSeparator.new()
	vbox.add_child(sep)

	var kills_label = Label.new()
	kills_label.text = "爆炸擊破：%d 個目標" % blast_kills
	kills_label.add_theme_font_size_override("font_size", 14)
	kills_label.add_theme_color_override("font_color", Color("#FFFFFF"))
	kills_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(kills_label)

	if blast_reward > 0:
		var reward_label = Label.new()
		reward_label.text = "💰 獲得 %d 金幣" % blast_reward
		reward_label.add_theme_font_size_override("font_size", 16)
		reward_label.add_theme_color_override("font_color", Color("#FFD700"))
		reward_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
		vbox.add_child(reward_label)

	add_child(popup)

	# 右側滑入動畫
	var tween = create_tween()
	tween.tween_property(popup, "position:x", 980.0, 0.3).set_ease(Tween.EASE_OUT)
	tween.tween_interval(4.0)
	tween.tween_property(popup, "modulate:a", 0.0, 0.5)
	tween.tween_callback(popup.queue_free)

## 清除所有染色標記
func _clear_color_markers() -> void:
	for tid in _color_markers:
		var marker = _color_markers[tid]
		if is_instance_valid(marker):
			marker.queue_free()
	_color_markers.clear()

## 彩虹三次強閃光（紅→綠→藍）
func _rainbow_triple_flash() -> void:
	var colors = [Color("#FF4444"), Color("#00CC44"), Color("#0088FF")]
	for i in range(3):
		var flash = ColorRect.new()
		flash.color = Color(colors[i].r, colors[i].g, colors[i].b, 0.0)
		flash.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
		add_child(flash)
		var delay = i * 0.18
		var tween = create_tween()
		tween.tween_interval(delay)
		tween.tween_property(flash, "color:a", 0.65, 0.08)
		tween.tween_property(flash, "color:a", 0.0, 0.14)
		tween.tween_callback(flash.queue_free)

## 建立頂部橫幅
func _make_banner(text: String, bg_color: Color, border_color: Color) -> PanelContainer:
	var panel = PanelContainer.new()
	panel.position = Vector2(0, -60)
	panel.size = Vector2(1280, 52)

	var style = StyleBoxFlat.new()
	style.bg_color = bg_color
	style.border_color = border_color
	style.border_width_bottom = 2
	panel.add_theme_stylebox_override("panel", style)

	var label = Label.new()
	label.text = text
	label.add_theme_font_size_override("font_size", 18)
	label.add_theme_color_override("font_color", Color("#FFFFFF"))
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	panel.add_child(label)

	return panel

## 取得顏色中文名稱
func _get_color_cn(color_name: String) -> String:
	match color_name:
		"red":    return "紅"
		"orange": return "橙"
		"yellow": return "黃"
		"green":  return "綠"
		"blue":   return "藍"
	return color_name
