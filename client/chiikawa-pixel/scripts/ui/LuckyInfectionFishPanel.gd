## LuckyInfectionFishPanel.gd — 幸運連鎖感染魚系統面板（DAY-219）
## 業界原創「病毒式蔓延」機制
##
## 視覺設計：
##   - 感染綠色主題（#00FF88 + #00CC66 + #88FFCC + #FF4444）
##   - infection_start：綠色雙閃光 + 頂部橫幅 + 感染目標標記（菱形輪廓）
##   - infection_spread：蔓延閃光 + 新感染目標標記 + 「🦠 感染蔓延！+N 個」浮動文字
##   - infection_kill：消失閃光 + 浮動獎勵文字（×2.0）
##   - infection_blast：全螢幕三次綠色強閃光 + 「🦠 感染爆發！」52px大字 + 結算彈窗
extends CanvasLayer

# 感染狀態
var _infection_active: bool = false
var _session_id: String = ""
var _infected_markers: Dictionary = {}  # instanceID -> Control（標記節點）
var _total_infected: int = 0

# 感染主題顏色
const COLOR_PRIMARY   = Color("#00FF88")  # 主色：感染綠
const COLOR_SECONDARY = Color("#00CC66")  # 次色：深綠
const COLOR_ACCENT    = Color("#88FFCC")  # 強調：淡綠
const COLOR_DANGER    = Color("#FF4444")  # 危險：紅色（爆發）
const COLOR_BG        = Color(0, 0.2, 0.1, 0.85)  # 背景

# UI 節點
var _banner: Control = null
var _blast_label: Label = null
var _result_panel: Control = null

func _ready() -> void:
	layer = 26  # 幸運連鎖感染魚面板層級

## 處理幸運連鎖感染魚訊息
func handle_lucky_infection_fish(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"infection_start":
			_on_infection_start(payload)
		"infection_spread":
			_on_infection_spread(payload)
		"infection_kill":
			_on_infection_kill(payload)
		"infection_blast":
			_on_infection_blast(payload)

## infection_start — 感染開始
func _on_infection_start(payload: Dictionary) -> void:
	_infection_active = true
	_session_id = payload.get("session_id", "")
	_total_infected = payload.get("total_infected", 0)
	var trigger_player: String = payload.get("trigger_player", "")
	var duration_sec: int = payload.get("duration_sec", 12)

	# 清除舊標記
	_clear_all_markers()

	# 綠色雙閃光
	_flash_screen(COLOR_PRIMARY, 0.25)
	await get_tree().create_timer(0.15).timeout
	_flash_screen(COLOR_PRIMARY, 0.2)

	# 頂部橫幅
	_show_banner("🦠 感染蔓延！", "%s 觸發感染魚！%d 個目標被感染" % [trigger_player, _total_infected], duration_sec)

	# 建立感染目標標記
	var infected_targets: Array = payload.get("infected_targets", [])
	for target_info in infected_targets:
		_add_infection_marker(target_info)

## infection_spread — 感染蔓延
func _on_infection_spread(payload: Dictionary) -> void:
	_total_infected = payload.get("total_infected", _total_infected)
	var new_targets: Array = payload.get("infected_targets", [])

	if new_targets.size() == 0:
		return

	# 蔓延閃光（較小）
	_flash_screen(COLOR_SECONDARY, 0.15)

	# 建立新感染目標標記
	for target_info in new_targets:
		_add_infection_marker(target_info)

	# 浮動文字：感染蔓延
	_show_float_text("🦠 感染蔓延！+%d 個" % new_targets.size(), COLOR_PRIMARY, Vector2(get_viewport().size / 2))

	# 更新橫幅計數
	if _banner != null:
		var count_label = _banner.get_node_or_null("CountLabel")
		if count_label:
			count_label.text = "感染目標：%d / %d" % [_total_infected, payload.get("max_infected", 8)]

## infection_kill — 感染目標被玩家擊破
func _on_infection_kill(payload: Dictionary) -> void:
	var killed_id: String = payload.get("killed_target", "")
	_total_infected = payload.get("total_infected", _total_infected)

	# 移除標記（消失閃光）
	if _infected_markers.has(killed_id):
		var marker = _infected_markers[killed_id]
		if is_instance_valid(marker):
			var tween = marker.create_tween()
			tween.tween_property(marker, "modulate:a", 0.0, 0.2)
			tween.tween_callback(marker.queue_free)
		_infected_markers.erase(killed_id)

	# 浮動文字：×2.0 倍率
	_show_float_text("×2.0 感染加成！", COLOR_ACCENT, Vector2(get_viewport().size / 2))

## infection_blast — 感染爆發
func _on_infection_blast(payload: Dictionary) -> void:
	_infection_active = false
	var total_killed: int = payload.get("total_killed", 0)
	var total_reward: int = payload.get("total_reward", 0)

	# 清除所有感染標記
	_clear_all_markers()
	_hide_banner()

	# 全螢幕三次強閃光（綠→白→綠）
	_flash_screen(COLOR_PRIMARY, 0.35)
	await get_tree().create_timer(0.12).timeout
	_flash_screen(Color.WHITE, 0.3)
	await get_tree().create_timer(0.12).timeout
	_flash_screen(COLOR_PRIMARY, 0.25)

	# 大字顯示
	_show_blast_label("🦠 感染爆發！")
	await get_tree().create_timer(1.5).timeout
	_hide_blast_label()

	# 結算彈窗（≥3 個擊破才顯示）
	if total_killed >= 3:
		_show_result_panel(total_killed, total_reward)

# ---- 輔助函數 ----

## 建立感染目標標記（菱形輪廓）
func _add_infection_marker(target_info: Dictionary) -> void:
	var iid: String = target_info.get("instance_id", "")
	if iid.is_empty() or _infected_markers.has(iid):
		return

	var x: float = target_info.get("x", 0.0)
	var y: float = target_info.get("y", 0.0)
	var layer_num: int = target_info.get("layer", 0)

	# 依層數決定顏色（層數越高越深）
	var marker_color: Color = COLOR_PRIMARY
	if layer_num == 1:
		marker_color = COLOR_SECONDARY
	elif layer_num >= 2:
		marker_color = COLOR_ACCENT

	# 建立菱形標記（4 個 ColorRect 組成）
	var container = Control.new()
	container.position = Vector2(x - 20, y - 20)
	container.size = Vector2(40, 40)
	add_child(container)

	# 四個角的 L 形線段（模擬菱形輪廓）
	var corners = [
		[Vector2(0, 15), Vector2(15, 0)],   # 左上
		[Vector2(25, 0), Vector2(40, 15)],  # 右上
		[Vector2(0, 25), Vector2(15, 40)],  # 左下
		[Vector2(25, 40), Vector2(40, 25)], # 右下
	]
	for corner in corners:
		var line = ColorRect.new()
		line.color = marker_color
		line.size = Vector2(2, 2)
		line.position = corner[0]
		container.add_child(line)

	# 感染符號（🦠）
	var label = Label.new()
	label.text = "🦠"
	label.add_theme_font_size_override("font_size", 14)
	label.position = Vector2(10, 10)
	container.add_child(label)

	# 閃爍動畫
	var tween = container.create_tween().set_loops()
	tween.tween_property(container, "modulate:a", 0.4, 0.5)
	tween.tween_property(container, "modulate:a", 1.0, 0.5)

	_infected_markers[iid] = container

## 清除所有感染標記
func _clear_all_markers() -> void:
	for iid in _infected_markers.keys():
		var marker = _infected_markers[iid]
		if is_instance_valid(marker):
			marker.queue_free()
	_infected_markers.clear()

## 顯示頂部橫幅
func _show_banner(title: String, subtitle: String, duration_sec: int) -> void:
	_hide_banner()

	var banner = Control.new()
	banner.set_anchors_preset(Control.PRESET_TOP_WIDE)
	banner.position = Vector2(0, 8)
	banner.size = Vector2(get_viewport().size.x, 52)
	add_child(banner)

	var bg = ColorRect.new()
	bg.color = COLOR_BG
	bg.size = banner.size
	banner.add_child(bg)

	var title_label = Label.new()
	title_label.text = title
	title_label.add_theme_font_size_override("font_size", 20)
	title_label.add_theme_color_override("font_color", COLOR_PRIMARY)
	title_label.position = Vector2(12, 4)
	banner.add_child(title_label)

	var sub_label = Label.new()
	sub_label.name = "CountLabel"
	sub_label.text = subtitle
	sub_label.add_theme_font_size_override("font_size", 13)
	sub_label.add_theme_color_override("font_color", COLOR_ACCENT)
	sub_label.position = Vector2(12, 28)
	banner.add_child(sub_label)

	# 計時條（底部）
	var timer_bar = ColorRect.new()
	timer_bar.name = "TimerBar"
	timer_bar.color = COLOR_PRIMARY
	timer_bar.position = Vector2(0, 48)
	timer_bar.size = Vector2(get_viewport().size.x, 4)
	banner.add_child(timer_bar)

	# 計時條縮短動畫
	var tween = banner.create_tween()
	tween.tween_property(timer_bar, "size:x", 0.0, float(duration_sec))

	_banner = banner

## 隱藏橫幅
func _hide_banner() -> void:
	if _banner != null and is_instance_valid(_banner):
		_banner.queue_free()
	_banner = null

## 顯示爆發大字
func _show_blast_label(text: String) -> void:
	_hide_blast_label()

	var label = Label.new()
	label.text = text
	label.add_theme_font_size_override("font_size", 52)
	label.add_theme_color_override("font_color", COLOR_PRIMARY)
	label.add_theme_color_override("font_shadow_color", Color.BLACK)
	label.add_theme_constant_override("shadow_offset_x", 3)
	label.add_theme_constant_override("shadow_offset_y", 3)
	label.set_anchors_preset(Control.PRESET_CENTER)
	label.position = get_viewport().size / 2 - Vector2(200, 40)
	add_child(label)

	# 彈跳縮放動畫
	label.scale = Vector2(0.5, 0.5)
	var tween = label.create_tween()
	tween.tween_property(label, "scale", Vector2(1.1, 1.1), 0.2)
	tween.tween_property(label, "scale", Vector2(1.0, 1.0), 0.1)

	_blast_label = label

## 隱藏爆發大字
func _hide_blast_label() -> void:
	if _blast_label != null and is_instance_valid(_blast_label):
		var tween = _blast_label.create_tween()
		tween.tween_property(_blast_label, "modulate:a", 0.0, 0.3)
		tween.tween_callback(_blast_label.queue_free)
	_blast_label = null

## 顯示結算彈窗
func _show_result_panel(killed: int, reward: int) -> void:
	if _result_panel != null and is_instance_valid(_result_panel):
		_result_panel.queue_free()

	var panel = Control.new()
	panel.size = Vector2(260, 120)
	panel.position = Vector2(get_viewport().size.x + 10, get_viewport().size.y / 2 - 60)
	add_child(panel)

	var bg = ColorRect.new()
	bg.color = COLOR_BG
	bg.size = panel.size
	panel.add_child(bg)

	var border = ColorRect.new()
	border.color = COLOR_PRIMARY
	border.size = Vector2(panel.size.x, 3)
	panel.add_child(border)

	var title_label = Label.new()
	title_label.text = "🦠 感染爆發結算"
	title_label.add_theme_font_size_override("font_size", 16)
	title_label.add_theme_color_override("font_color", COLOR_PRIMARY)
	title_label.position = Vector2(12, 12)
	panel.add_child(title_label)

	var killed_label = Label.new()
	killed_label.text = "擊破目標：%d 個" % killed
	killed_label.add_theme_font_size_override("font_size", 14)
	killed_label.add_theme_color_override("font_color", Color.WHITE)
	killed_label.position = Vector2(12, 42)
	panel.add_child(killed_label)

	var reward_label = Label.new()
	reward_label.text = "全服共享：%d 金幣" % reward
	reward_label.add_theme_font_size_override("font_size", 14)
	reward_label.add_theme_color_override("font_color", COLOR_ACCENT)
	reward_label.position = Vector2(12, 68)
	panel.add_child(reward_label)

	# 右側滑入動畫
	var tween = panel.create_tween()
	tween.tween_property(panel, "position:x", get_viewport().size.x - 280.0, 0.3)

	# 3 秒後淡出
	await get_tree().create_timer(3.0).timeout
	if is_instance_valid(panel):
		var fade_tween = panel.create_tween()
		fade_tween.tween_property(panel, "modulate:a", 0.0, 0.4)
		fade_tween.tween_callback(panel.queue_free)

	_result_panel = null

## 全螢幕閃光效果
func _flash_screen(color: Color, duration: float) -> void:
	var flash = ColorRect.new()
	flash.color = Color(color.r, color.g, color.b, 0.45)
	flash.set_anchors_preset(Control.PRESET_FULL_RECT)
	add_child(flash)

	var tween = flash.create_tween()
	tween.tween_property(flash, "modulate:a", 0.0, duration)
	tween.tween_callback(flash.queue_free)

## 浮動文字
func _show_float_text(text: String, color: Color, pos: Vector2) -> void:
	var label = Label.new()
	label.text = text
	label.add_theme_font_size_override("font_size", 18)
	label.add_theme_color_override("font_color", color)
	label.position = pos - Vector2(80, 20)
	add_child(label)

	var tween = label.create_tween()
	tween.tween_property(label, "position:y", label.position.y - 40, 0.8)
	tween.parallel().tween_property(label, "modulate:a", 0.0, 0.8)
	tween.tween_callback(label.queue_free)
