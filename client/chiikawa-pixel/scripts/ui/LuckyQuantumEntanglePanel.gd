## LuckyQuantumEntanglePanel.gd — 幸運量子糾纏魚系統面板（DAY-251）
## 業界原創「量子糾纏+同步爆炸+量子共鳴」機制
##
## 視覺設計：
##   - 深藍量子主題（#1A5276 + #2471A3 + #AED6F1 + #EBF5FB）
##   - entangle_start：深藍三次強閃光 + 頂部橫幅 + 「⚛️ 量子糾纏！」大字 + 糾纏連線 + 計時條
##   - entangle_broadcast：頂部小橫幅（全服廣播）
##   - entangle_sync：深藍閃光 + 「⚛️ 同步爆炸！×1.8」浮動文字 + 全服獎勵
##   - entangle_resonance：全螢幕三次強閃光 + 「⚛️ 量子共鳴！×3.5」大字 + 結算彈窗
##   - entangle_decay：灰色閃光 + 「⚛️ 量子衰變！HP-60%」提示
##   - entangle_end：計時條淡出
extends CanvasLayer

# 主題顏色
const COLOR_QUANTUM  = Color("#1A5276")  # 深藍（主題）
const COLOR_BRIGHT   = Color("#2471A3")  # 亮藍（強調）
const COLOR_LIGHT    = Color("#AED6F1")  # 淺藍（背景）
const COLOR_RESON    = Color("#F39C12")  # 金色（量子共鳴）
const COLOR_DECAY    = Color("#7F8C8D")  # 灰色（衰變）
const COLOR_WHITE    = Color("#FFFFFF")  # 白色

# 計時條
var _timer_bar: ColorRect = null
var _timer_bar_bg: ColorRect = null
var _timer_tween: Tween = null

# 糾纏連線（兩個目標之間的連線）
var _entangle_lines: Array = []

func _ready() -> void:
	layer = 24  # 幸運量子糾纏魚面板層級（DAY-251）

## 處理幸運量子糾纏魚訊息
func handle_lucky_quantum_entangle(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"entangle_start":
			_on_entangle_start(payload)
		"entangle_broadcast":
			_on_entangle_broadcast(payload)
		"entangle_sync":
			_on_entangle_sync(payload)
		"entangle_resonance":
			_on_entangle_resonance(payload)
		"entangle_decay":
			_on_entangle_decay(payload)
		"entangle_end":
			_on_entangle_end(payload)

## entangle_start — 量子糾纏啟動（個人訊息）
func _on_entangle_start(payload: Dictionary) -> void:
	var duration_sec: int = payload.get("duration_sec", 20)
	var sync_mult: float = payload.get("sync_mult", 1.8)
	var reson_mult: float = payload.get("reson_mult", 3.5)
	var targets = payload.get("targets", [])

	# 深藍三次強閃光
	_flash_screen(COLOR_QUANTUM, 0.55, 3)

	# 頂部橫幅
	_show_banner("⚛️ 量子糾纏！2 個目標被量子連結！", COLOR_QUANTUM, 4.0)

	# 中央大字
	_show_big_text("⚛️ 量子糾纏！", COLOR_BRIGHT, 52, 2.5)

	# 倍率說明
	_show_sub_text("同步爆炸 ×%.1f（全服）  量子共鳴 ×%.1f（全服大獎）" % [sync_mult, reson_mult], COLOR_RESON, 2.5)

	# 糾纏連線（如果有目標位置資訊）
	if targets is Array and targets.size() >= 2:
		_show_entangle_connection(targets)

	# 右側豎向計時條
	_start_timer_bar(duration_sec)

## entangle_broadcast — 全服廣播糾纏
func _on_entangle_broadcast(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "某玩家")
	var reson_mult: float = payload.get("reson_mult", 3.5)
	_show_top_banner("⚛️ %s 觸發量子糾纏！同時擊破可觸發 ×%.1f 共鳴！" % [player_name, reson_mult], COLOR_QUANTUM, 3.0)

## entangle_sync — 同步爆炸（全服廣播）
func _on_entangle_sync(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "某玩家")
	var sync_mult: float = payload.get("sync_mult", 1.8)
	var total_reward: int = payload.get("total_reward", 0)

	# 深藍閃光
	_flash_screen(COLOR_BRIGHT, 0.4, 2)

	# 浮動文字
	_show_float_text("⚛️ %s 觸發同步爆炸！×%.1f  全服+%d" % [player_name, sync_mult, total_reward], COLOR_LIGHT, 2.0)

	# 清除糾纏連線（一個目標已爆炸）
	_clear_entangle_lines()

## entangle_resonance — 量子共鳴（全服廣播，大獎）
func _on_entangle_resonance(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "某玩家")
	var first_killer_name: String = payload.get("first_killer_name", "某玩家")
	var reson_mult: float = payload.get("reson_mult", 3.5)
	var total_reward: int = payload.get("total_reward", 0)
	var time_diff_ms: int = payload.get("time_diff_ms", 0)

	# 停止計時條
	_stop_timer_bar()
	_clear_entangle_lines()

	# 全螢幕三次強閃光（金色，量子共鳴感）
	_flash_screen(COLOR_RESON, 0.7, 3)

	# 大字
	_show_big_text("⚛️ 量子共鳴！", COLOR_RESON, 56, 3.0)
	_show_sub_text("%s + %s 同時擊破！×%.1f 全服大獎！" % [first_killer_name, player_name, reson_mult], COLOR_WHITE, 3.0)

	# 結算彈窗
	_show_resonance_popup(first_killer_name, player_name, reson_mult, total_reward, time_diff_ms)

## entangle_decay — 量子衰變（全服廣播）
func _on_entangle_decay(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "某玩家")
	var decay_hp: int = payload.get("decay_hp", 60)

	# 停止計時條
	_stop_timer_bar()
	_clear_entangle_lines()

	# 灰色閃光（衰變感）
	_flash_screen(COLOR_DECAY, 0.35, 2)

	# 衰變提示
	_show_big_text("⚛️ 量子衰變！", COLOR_DECAY, 44, 2.0)
	_show_sub_text("%s 的糾纏目標衰變！HP -%d%%！趕快擊破！" % [player_name, decay_hp], COLOR_LIGHT, 2.0)

## entangle_end — 糾纏結束（全服廣播）
func _on_entangle_end(_payload: Dictionary) -> void:
	_stop_timer_bar()
	_clear_entangle_lines()

# ─── 糾纏連線 ────────────────────────────────────────────────────────────────

func _show_entangle_connection(targets: Array) -> void:
	# 在兩個目標之間畫一條閃爍的量子連線
	# targets[0] 和 targets[1] 各有 x, y 座標
	if targets.size() < 2:
		return

	var vp_size = get_viewport().size
	var t0 = targets[0]
	var t1 = targets[1]

	# 遊戲座標 → 螢幕座標（假設遊戲寬 1000，高 600）
	var game_w: float = 1000.0
	var game_h: float = 600.0
	var sx0: float = float(t0.get("x", 250)) / game_w * vp_size.x
	var sy0: float = float(t0.get("y", 300)) / game_h * vp_size.y
	var sx1: float = float(t1.get("x", 750)) / game_w * vp_size.x
	var sy1: float = float(t1.get("y", 300)) / game_h * vp_size.y

	# 用多個小圓點模擬連線（每 30px 一個點）
	var dx: float = sx1 - sx0
	var dy: float = sy1 - sy0
	var dist: float = sqrt(dx * dx + dy * dy)
	var steps: int = max(1, int(dist / 30.0))

	for i in range(steps + 1):
		var t: float = float(i) / float(steps)
		var px: float = sx0 + dx * t
		var py: float = sy0 + dy * t

		var dot = ColorRect.new()
		dot.color = Color(COLOR_BRIGHT.r, COLOR_BRIGHT.g, COLOR_BRIGHT.b, 0.8)
		dot.size = Vector2(6, 6)
		dot.position = Vector2(px - 3, py - 3)
		add_child(dot)
		_entangle_lines.append(dot)

		# 閃爍動畫（每個點延遲不同，製造波動感）
		var tween = dot.create_tween().set_loops()
		tween.tween_interval(float(i) * 0.05)
		tween.tween_property(dot, "modulate:a", 0.2, 0.4)
		tween.tween_property(dot, "modulate:a", 1.0, 0.4)

	# 在兩個目標位置加上 ⚛️ 標記
	for pos in [[sx0, sy0], [sx1, sy1]]:
		var marker = Label.new()
		marker.text = "⚛️"
		marker.add_theme_font_size_override("font_size", 24)
		marker.position = Vector2(pos[0] - 16, pos[1] - 16)
		add_child(marker)
		_entangle_lines.append(marker)

		# 脈衝縮放動畫
		var tween2 = marker.create_tween().set_loops()
		tween2.tween_property(marker, "scale", Vector2(1.3, 1.3), 0.5)
		tween2.tween_property(marker, "scale", Vector2(1.0, 1.0), 0.5)

func _clear_entangle_lines() -> void:
	for node in _entangle_lines:
		if is_instance_valid(node):
			node.queue_free()
	_entangle_lines.clear()

# ─── 量子共鳴結算彈窗 ────────────────────────────────────────────────────────

func _show_resonance_popup(first_killer: String, second_killer: String, reson_mult: float, total_reward: int, time_diff_ms: int) -> void:
	var vp_size = get_viewport().size
	var popup = PanelContainer.new()
	popup.size = Vector2(340, 180)
	popup.position = Vector2(vp_size.x + 10, vp_size.y / 2.0 - 90)
	add_child(popup)

	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.02, 0.05, 0.12, 0.94)
	style.border_color = COLOR_RESON
	style.set_border_width_all(3)
	style.corner_radius_top_left = 12
	style.corner_radius_top_right = 12
	style.corner_radius_bottom_left = 12
	style.corner_radius_bottom_right = 12
	popup.add_theme_stylebox_override("panel", style)

	var vbox = VBoxContainer.new()
	vbox.add_theme_constant_override("separation", 8)
	popup.add_child(vbox)

	var title_lbl = Label.new()
	title_lbl.text = "⚛️ 量子共鳴！"
	title_lbl.add_theme_font_size_override("font_size", 22)
	title_lbl.add_theme_color_override("font_color", COLOR_RESON)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title_lbl)

	var killers_lbl = Label.new()
	killers_lbl.text = "%s + %s" % [first_killer, second_killer]
	killers_lbl.add_theme_font_size_override("font_size", 14)
	killers_lbl.add_theme_color_override("font_color", COLOR_LIGHT)
	killers_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(killers_lbl)

	var time_lbl = Label.new()
	time_lbl.text = "間隔 %.2f 秒" % (float(time_diff_ms) / 1000.0)
	time_lbl.add_theme_font_size_override("font_size", 13)
	time_lbl.add_theme_color_override("font_color", COLOR_DECAY)
	time_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(time_lbl)

	var mult_lbl = Label.new()
	mult_lbl.text = "×%.1f 量子共鳴倍率" % reson_mult
	mult_lbl.add_theme_font_size_override("font_size", 18)
	mult_lbl.add_theme_color_override("font_color", COLOR_RESON)
	mult_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(mult_lbl)

	var reward_lbl = Label.new()
	reward_lbl.text = "全服獎勵：+%d 籌碼" % total_reward
	reward_lbl.add_theme_font_size_override("font_size", 18)
	reward_lbl.add_theme_color_override("font_color", COLOR_WHITE)
	reward_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(reward_lbl)

	# 右側滑入動畫
	var tween = create_tween()
	tween.tween_property(popup, "position:x", vp_size.x - 360.0, 0.4).set_ease(Tween.EASE_OUT)
	tween.tween_interval(4.0)
	tween.tween_property(popup, "position:x", vp_size.x + 10.0, 0.4).set_ease(Tween.EASE_IN)
	tween.tween_callback(popup.queue_free)

# ─── 通用 UI 工具 ─────────────────────────────────────────────────────────────

func _flash_screen(color: Color, alpha: float, times: int) -> void:
	var vp_size = get_viewport().size
	var flash = ColorRect.new()
	flash.color = Color(color.r, color.g, color.b, 0.0)
	flash.size = vp_size
	add_child(flash)
	var tween = create_tween()
	for i in range(times):
		tween.tween_property(flash, "color:a", alpha, 0.08)
		tween.tween_property(flash, "color:a", 0.0, 0.12)
	tween.tween_callback(flash.queue_free)

func _show_banner(text: String, color: Color, duration: float) -> void:
	var vp_size = get_viewport().size
	var banner = PanelContainer.new()
	banner.size = Vector2(vp_size.x, 48)
	banner.position = Vector2(0, 0)
	add_child(banner)
	var style = StyleBoxFlat.new()
	style.bg_color = Color(color.r, color.g, color.b, 0.88)
	banner.add_theme_stylebox_override("panel", style)
	var lbl = Label.new()
	lbl.text = text
	lbl.add_theme_font_size_override("font_size", 18)
	lbl.add_theme_color_override("font_color", COLOR_WHITE)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	lbl.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	banner.add_child(lbl)
	var tween = create_tween()
	tween.tween_interval(duration)
	tween.tween_property(banner, "modulate:a", 0.0, 0.5)
	tween.tween_callback(banner.queue_free)

func _show_top_banner(text: String, color: Color, duration: float) -> void:
	var vp_size = get_viewport().size
	var banner = PanelContainer.new()
	banner.size = Vector2(vp_size.x * 0.75, 36)
	banner.position = Vector2(vp_size.x * 0.125, 52)
	add_child(banner)
	var style = StyleBoxFlat.new()
	style.bg_color = Color(color.r, color.g, color.b, 0.82)
	style.corner_radius_top_left = 8
	style.corner_radius_top_right = 8
	style.corner_radius_bottom_left = 8
	style.corner_radius_bottom_right = 8
	banner.add_theme_stylebox_override("panel", style)
	var lbl = Label.new()
	lbl.text = text
	lbl.add_theme_font_size_override("font_size", 13)
	lbl.add_theme_color_override("font_color", COLOR_WHITE)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	lbl.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	banner.add_child(lbl)
	var tween = create_tween()
	tween.tween_interval(duration)
	tween.tween_property(banner, "modulate:a", 0.0, 0.4)
	tween.tween_callback(banner.queue_free)

func _show_big_text(text: String, color: Color, font_size: int, duration: float) -> void:
	var vp_size = get_viewport().size
	var lbl = Label.new()
	lbl.text = text
	lbl.add_theme_font_size_override("font_size", font_size)
	lbl.add_theme_color_override("font_color", color)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.size = Vector2(vp_size.x, 80)
	lbl.position = Vector2(0, vp_size.y * 0.35)
	add_child(lbl)
	var tween = create_tween()
	tween.tween_property(lbl, "scale", Vector2(1.15, 1.15), 0.15).set_ease(Tween.EASE_OUT)
	tween.tween_property(lbl, "scale", Vector2(1.0, 1.0), 0.1)
	tween.tween_interval(duration - 0.5)
	tween.tween_property(lbl, "modulate:a", 0.0, 0.5)
	tween.tween_callback(lbl.queue_free)

func _show_sub_text(text: String, color: Color, duration: float) -> void:
	var vp_size = get_viewport().size
	var lbl = Label.new()
	lbl.text = text
	lbl.add_theme_font_size_override("font_size", 14)
	lbl.add_theme_color_override("font_color", color)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.size = Vector2(vp_size.x, 40)
	lbl.position = Vector2(0, vp_size.y * 0.35 + 80)
	add_child(lbl)
	var tween = create_tween()
	tween.tween_interval(duration)
	tween.tween_property(lbl, "modulate:a", 0.0, 0.4)
	tween.tween_callback(lbl.queue_free)

func _show_float_text(text: String, color: Color, duration: float) -> void:
	var vp_size = get_viewport().size
	var lbl = Label.new()
	lbl.text = text
	lbl.add_theme_font_size_override("font_size", 17)
	lbl.add_theme_color_override("font_color", color)
	lbl.position = Vector2(
		randf_range(vp_size.x * 0.25, vp_size.x * 0.65),
		randf_range(vp_size.y * 0.3, vp_size.y * 0.6)
	)
	add_child(lbl)
	var tween = create_tween()
	tween.tween_property(lbl, "position:y", lbl.position.y - 60, duration)
	tween.parallel().tween_property(lbl, "modulate:a", 0.0, duration)
	tween.tween_callback(lbl.queue_free)

func _start_timer_bar(duration_sec: int) -> void:
	var vp_size = get_viewport().size
	var bar_h: float = vp_size.y * 0.5
	var bar_w: float = 12.0
	# x=-100 與其他計時條錯開（龍捲風-58，黑洞-72，鏡像分裂-86，量子糾纏-100）
	var bar_x: float = vp_size.x - 100.0
	var bar_y: float = vp_size.y * 0.25

	_timer_bar_bg = ColorRect.new()
	_timer_bar_bg.color = Color(0.02, 0.05, 0.12, 0.7)
	_timer_bar_bg.size = Vector2(bar_w, bar_h)
	_timer_bar_bg.position = Vector2(bar_x, bar_y)
	add_child(_timer_bar_bg)

	_timer_bar = ColorRect.new()
	_timer_bar.color = COLOR_BRIGHT
	_timer_bar.size = Vector2(bar_w, bar_h)
	_timer_bar.position = Vector2(bar_x, bar_y)
	add_child(_timer_bar)

	_timer_tween = create_tween()
	_timer_tween.tween_property(_timer_bar, "size:y", 0.0, float(duration_sec)).set_ease(Tween.EASE_IN_OUT)

func _stop_timer_bar() -> void:
	if is_instance_valid(_timer_tween):
		_timer_tween.kill()
	if is_instance_valid(_timer_bar):
		var tween = create_tween()
		tween.tween_property(_timer_bar, "modulate:a", 0.0, 0.5)
		tween.tween_callback(_timer_bar.queue_free)
		_timer_bar = null
	if is_instance_valid(_timer_bar_bg):
		var tween2 = create_tween()
		tween2.tween_property(_timer_bar_bg, "modulate:a", 0.0, 0.5)
		tween2.tween_callback(_timer_bar_bg.queue_free)
		_timer_bar_bg = null
