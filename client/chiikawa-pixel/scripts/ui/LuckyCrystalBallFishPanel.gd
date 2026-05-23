## LuckyCrystalBallFishPanel.gd — 幸運水晶球魚系統面板（DAY-246）
## 業界原創「預測未來+命中率提升」機制
##
## 視覺設計：
##   - 青綠水晶主題（#1ABC9C + #16A085 + #A3E4D7 + #E8F8F5）
##   - crystal_start：青綠三次強閃光 + 頂部橫幅 + 「🔮 水晶預言！」大字 + 目標標記（水晶圓圈）+ 計時條 + 倍率說明
##   - crystal_broadcast：頂部小橫幅（全服廣播）
##   - crystal_hit：水晶擊破閃光 + ×2.5 浮動文字（金色）
##   - crystal_blast：水晶爆炸閃光 + ×1.8 浮動文字（青綠）
##   - crystal_end：三次強閃光 + 「🔮 預言結束！」大字 + 結算彈窗
extends CanvasLayer

# 主題顏色
const COLOR_CRYSTAL  = Color("#1ABC9C")  # 青綠（主題）
const COLOR_DARK     = Color("#16A085")  # 深青綠（強調）
const COLOR_LIGHT    = Color("#A3E4D7")  # 淺青綠（背景）
const COLOR_GOLD     = Color("#F39C12")  # 金色（必中獎勵）
const COLOR_WHITE    = Color("#FFFFFF")  # 白色

# 計時條節點
var _timer_bar: ColorRect = null
var _timer_bar_bg: ColorRect = null
var _timer_tween: Tween = null
var _duration_sec: int = 8

# 水晶目標標記（instanceID → Control）
var _crystal_markers: Dictionary = {}

func _ready() -> void:
	layer = 19  # 幸運水晶球魚面板層級（DAY-246）

## 處理幸運水晶球魚訊息
func handle_lucky_crystal_ball_fish(payload: Dictionary) -> void:
	var event: String = payload.get("event", "")
	match event:
		"crystal_start":
			_on_crystal_start(payload)
		"crystal_broadcast":
			_on_crystal_broadcast(payload)
		"crystal_hit":
			_on_crystal_hit(payload)
		"crystal_blast":
			_on_crystal_blast(payload)
		"crystal_end":
			_on_crystal_end(payload)

## crystal_start — 水晶預言啟動（個人訊息）
func _on_crystal_start(payload: Dictionary) -> void:
	_duration_sec = payload.get("duration_sec", 8)
	var hit_mult: float = payload.get("hit_mult", 2.5)
	var blast_mult: float = payload.get("blast_mult", 1.8)
	var target_ids: Array = payload.get("target_ids", [])

	# 青綠三次強閃光
	_flash_screen(COLOR_CRYSTAL, 0.5, 3)

	# 頂部橫幅
	_show_banner("🔮 水晶預言！%d 個目標必中！" % target_ids.size(), COLOR_CRYSTAL, 3.5)

	# 中央大字
	_show_big_text("🔮 水晶預言！", COLOR_CRYSTAL, 52, 2.5)

	# 倍率說明
	_show_sub_text("必中擊破 ×%.1f  水晶爆炸 ×%.1f" % [hit_mult, blast_mult], COLOR_GOLD, 2.0)

	# 右側豎向計時條
	_start_timer_bar(_duration_sec)

## crystal_broadcast — 全服廣播水晶預言啟動
func _on_crystal_broadcast(payload: Dictionary) -> void:
	var player_name: String = payload.get("player_name", "某玩家")
	_show_top_banner("🔮 %s 觸發水晶預言！" % player_name, COLOR_CRYSTAL, 2.5)

## crystal_hit — 水晶預言目標被必中擊破
func _on_crystal_hit(payload: Dictionary) -> void:
	var target_id: String = payload.get("target_id", "")
	var reward: int = payload.get("reward", 0)
	var hit_mult: float = payload.get("hit_mult", 2.5)

	# 移除水晶標記
	if _crystal_markers.has(target_id):
		var marker = _crystal_markers[target_id]
		if is_instance_valid(marker):
			# 擊破閃光
			var tween = create_tween()
			tween.tween_property(marker, "modulate:a", 0.0, 0.25)
			tween.tween_callback(marker.queue_free)
		_crystal_markers.erase(target_id)

	# 金色浮動獎勵文字
	_show_float_text("🔮 必中！×%.1f  +%d" % [hit_mult, reward], COLOR_GOLD, 1.8)

## crystal_blast — 水晶爆炸
func _on_crystal_blast(payload: Dictionary) -> void:
	var target_id: String = payload.get("target_id", "")
	var reward: int = payload.get("reward", 0)
	var blast_mult: float = payload.get("blast_mult", 1.8)

	# 移除水晶標記
	if _crystal_markers.has(target_id):
		var marker = _crystal_markers[target_id]
		if is_instance_valid(marker):
			var tween = create_tween()
			tween.tween_property(marker, "scale", Vector2(2.0, 2.0), 0.2)
			tween.parallel().tween_property(marker, "modulate:a", 0.0, 0.2)
			tween.tween_callback(marker.queue_free)
		_crystal_markers.erase(target_id)

	# 青綠浮動文字
	_show_float_text("💥 爆炸！×%.1f  +%d" % [blast_mult, reward], COLOR_CRYSTAL, 1.5)

## crystal_end — 水晶預言結束
func _on_crystal_end(payload: Dictionary) -> void:
	var hit_count: int = payload.get("hit_count", 0)
	var blast_count: int = payload.get("blast_count", 0)
	var total_reward: int = payload.get("total_reward", 0)

	# 清除所有水晶標記
	for target_id in _crystal_markers.keys():
		var marker = _crystal_markers[target_id]
		if is_instance_valid(marker):
			marker.queue_free()
	_crystal_markers.clear()

	# 停止計時條
	_stop_timer_bar()

	if hit_count + blast_count > 0:
		# 三次強閃光
		_flash_screen(COLOR_DARK, 0.55, 3)

		# 中央大字
		_show_big_text("🔮 預言結束！", COLOR_CRYSTAL, 48, 2.0)

		# 結算彈窗
		_show_result_popup(hit_count, blast_count, total_reward)

# ─── 結算彈窗 ────────────────────────────────────────────────────────────────

func _show_result_popup(hit_count: int, blast_count: int, total_reward: int) -> void:
	var vp_size = get_viewport().size
	var popup = PanelContainer.new()
	popup.size = Vector2(320, 170)
	popup.position = Vector2(vp_size.x + 10, vp_size.y / 2.0 - 85)
	add_child(popup)

	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.02, 0.1, 0.1, 0.92)
	style.border_color = COLOR_CRYSTAL
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
	title_lbl.text = "🔮 水晶預言結算"
	title_lbl.add_theme_font_size_override("font_size", 20)
	title_lbl.add_theme_color_override("font_color", COLOR_CRYSTAL)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(title_lbl)

	var hit_lbl = Label.new()
	hit_lbl.text = "必中擊破：%d 個" % hit_count
	hit_lbl.add_theme_font_size_override("font_size", 16)
	hit_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	hit_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(hit_lbl)

	var blast_lbl = Label.new()
	blast_lbl.text = "水晶爆炸：%d 個" % blast_count
	blast_lbl.add_theme_font_size_override("font_size", 16)
	blast_lbl.add_theme_color_override("font_color", COLOR_LIGHT)
	blast_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(blast_lbl)

	var reward_lbl = Label.new()
	reward_lbl.text = "總獎勵：+%d 籌碼" % total_reward
	reward_lbl.add_theme_font_size_override("font_size", 18)
	reward_lbl.add_theme_color_override("font_color", COLOR_GOLD)
	reward_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(reward_lbl)

	# 右側滑入動畫
	var tween = create_tween()
	tween.tween_property(popup, "position:x", vp_size.x - 340.0, 0.4).set_ease(Tween.EASE_OUT)
	tween.tween_interval(3.5)
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
	banner.size = Vector2(vp_size.x * 0.7, 36)
	banner.position = Vector2(vp_size.x * 0.15, 52)
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
	lbl.add_theme_font_size_override("font_size", 14)
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
	lbl.add_theme_font_size_override("font_size", 18)
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
	lbl.add_theme_font_size_override("font_size", 20)
	lbl.add_theme_color_override("font_color", color)
	lbl.position = Vector2(
		randf_range(vp_size.x * 0.3, vp_size.x * 0.7),
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
	var bar_x: float = vp_size.x - 44.0  # 與幽靈魚計時條錯開
	var bar_y: float = vp_size.y * 0.25

	_timer_bar_bg = ColorRect.new()
	_timer_bar_bg.color = Color(0.0, 0.08, 0.08, 0.7)
	_timer_bar_bg.size = Vector2(bar_w, bar_h)
	_timer_bar_bg.position = Vector2(bar_x, bar_y)
	add_child(_timer_bar_bg)

	_timer_bar = ColorRect.new()
	_timer_bar.color = COLOR_CRYSTAL
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
