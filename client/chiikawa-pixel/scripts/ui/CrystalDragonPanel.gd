## CrystalDragonPanel.gd
## 水晶龍收集大獎面板（DAY-153）
## 業界依據：jiligames.com JILI Flying Dragon 2026「collect crystals to get the grand prize!
## Kill the Underworld Dragon and win the prize!」
## 設計：紫水晶主題；crystal_dragon_drop 顯示水晶掉落動畫+進度條更新；
## crystal_dragon_reward 全螢幕紫色閃光+地獄龍大獎彈窗（含貢獻者列表）

extends Control

# ---- 節點 ----
var _progress_banner: Control = null
var _progress_bar: ColorRect = null
var _progress_fill: ColorRect = null
var _progress_label: Label = null
var _crystal_count_label: Label = null
var _reward_panel: Control = null
var _reward_title: Label = null
var _reward_contributors: Label = null
var _reward_total: Label = null
var _flash_overlay: ColorRect = null

# ---- 浮動水晶標籤池 ----
var _float_labels: Array = []
const MAX_FLOAT_LABELS = 8

# ---- 顏色（紫水晶主題）----
const COLOR_CRYSTAL  = Color(0.6, 0.2, 1.0, 1.0)   # 紫水晶
const COLOR_BRIGHT   = Color(0.8, 0.5, 1.0, 1.0)   # 亮紫
const COLOR_DARK     = Color(0.3, 0.0, 0.6, 1.0)   # 深紫
const COLOR_GOLD     = Color(1.0, 0.85, 0.0, 1.0)  # 金色（大獎）
const COLOR_BG       = Color(0.04, 0.0, 0.08, 0.95) # 深紫背景

# ---- 狀態 ----
var _current_crystals: int = 0
var _goal: int = 50
var _is_visible: bool = false

func _ready() -> void:
	_build_ui()
	visible = false
	mouse_filter = Control.MOUSE_FILTER_IGNORE

func _build_ui() -> void:
	# 全螢幕閃光層
	_flash_overlay = ColorRect.new()
	_flash_overlay.set_anchors_preset(Control.PRESET_FULL_RECT)
	_flash_overlay.color = Color(0.5, 0.0, 1.0, 0.0)
	_flash_overlay.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_flash_overlay)

	# 頂部進度橫幅（常駐顯示）
	_progress_banner = Control.new()
	_progress_banner.position = Vector2(0, 0)
	_progress_banner.size = Vector2(1280, 48)
	_progress_banner.visible = false
	add_child(_progress_banner)

	var banner_bg = ColorRect.new()
	banner_bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	banner_bg.color = Color(0.04, 0.0, 0.08, 0.88)
	_progress_banner.add_child(banner_bg)

	var banner_border = ColorRect.new()
	banner_border.color = COLOR_CRYSTAL
	banner_border.position = Vector2(0, 44)
	banner_border.size = Vector2(1280, 4)
	_progress_banner.add_child(banner_border)

	# 水晶圖示 + 標題
	var title_label = Label.new()
	title_label.text = "💎 水晶龍"
	title_label.position = Vector2(12, 8)
	title_label.size = Vector2(120, 32)
	title_label.add_theme_color_override("font_color", COLOR_BRIGHT)
	title_label.add_theme_font_size_override("font_size", 18)
	_progress_banner.add_child(title_label)

	# 進度條背景
	_progress_bar = ColorRect.new()
	_progress_bar.position = Vector2(140, 14)
	_progress_bar.size = Vector2(800, 20)
	_progress_bar.color = Color(0.15, 0.05, 0.25, 1.0)
	_progress_banner.add_child(_progress_bar)

	# 進度條填充
	_progress_fill = ColorRect.new()
	_progress_fill.position = Vector2(140, 14)
	_progress_fill.size = Vector2(0, 20)
	_progress_fill.color = COLOR_CRYSTAL
	_progress_banner.add_child(_progress_fill)

	# 進度文字
	_progress_label = Label.new()
	_progress_label.text = "0/50"
	_progress_label.position = Vector2(950, 8)
	_progress_label.size = Vector2(120, 32)
	_progress_label.add_theme_color_override("font_color", COLOR_BRIGHT)
	_progress_label.add_theme_font_size_override("font_size", 16)
	_progress_banner.add_child(_progress_label)

	# 水晶數量大字
	_crystal_count_label = Label.new()
	_crystal_count_label.text = "💎 0"
	_crystal_count_label.position = Vector2(1080, 6)
	_crystal_count_label.size = Vector2(180, 36)
	_crystal_count_label.add_theme_color_override("font_color", COLOR_GOLD)
	_crystal_count_label.add_theme_font_size_override("font_size", 20)
	_progress_banner.add_child(_crystal_count_label)

	# 地獄龍大獎彈窗（初始隱藏）
	_reward_panel = Control.new()
	_reward_panel.position = Vector2(1280, 100)
	_reward_panel.size = Vector2(380, 300)
	_reward_panel.visible = false
	add_child(_reward_panel)

	var reward_bg = ColorRect.new()
	reward_bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	reward_bg.color = COLOR_BG
	_reward_panel.add_child(reward_bg)

	var reward_border = ColorRect.new()
	reward_border.color = COLOR_CRYSTAL
	reward_border.position = Vector2(0, 0)
	reward_border.size = Vector2(380, 4)
	_reward_panel.add_child(reward_border)

	var reward_border_b = ColorRect.new()
	reward_border_b.color = COLOR_CRYSTAL
	reward_border_b.position = Vector2(0, 296)
	reward_border_b.size = Vector2(380, 4)
	_reward_panel.add_child(reward_border_b)

	_reward_title = Label.new()
	_reward_title.text = "🐉 地獄龍大獎！"
	_reward_title.position = Vector2(0, 12)
	_reward_title.size = Vector2(380, 40)
	_reward_title.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_reward_title.add_theme_color_override("font_color", COLOR_GOLD)
	_reward_title.add_theme_font_size_override("font_size", 24)
	_reward_panel.add_child(_reward_title)

	_reward_contributors = Label.new()
	_reward_contributors.text = ""
	_reward_contributors.position = Vector2(12, 60)
	_reward_contributors.size = Vector2(356, 180)
	_reward_contributors.add_theme_color_override("font_color", COLOR_BRIGHT)
	_reward_contributors.add_theme_font_size_override("font_size", 14)
	_reward_contributors.autowrap_mode = TextServer.AUTOWRAP_WORD_SMART
	_reward_panel.add_child(_reward_contributors)

	_reward_total = Label.new()
	_reward_total.text = "總獎勵：0 金幣"
	_reward_total.position = Vector2(0, 252)
	_reward_total.size = Vector2(380, 36)
	_reward_total.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_reward_total.add_theme_color_override("font_color", COLOR_GOLD)
	_reward_total.add_theme_font_size_override("font_size", 18)
	_reward_panel.add_child(_reward_total)

# ---- 外部呼叫 API ----

## 更新水晶進度（登入時或衰減時）
func update_status(total_crystals: int, goal: int, progress: float) -> void:
	_current_crystals = total_crystals
	_goal = goal
	_update_progress_bar(progress)
	if total_crystals > 0:
		_show_progress_banner()

## 處理水晶掉落事件
func on_crystal_drop(killer_name: String, crystals_gain: int, total_crystals: int, goal: int, progress: float) -> void:
	_current_crystals = total_crystals
	_goal = goal
	_update_progress_bar(progress)
	_show_progress_banner()

	# 浮動水晶文字
	var msg = "+%d 💎 %s" % [crystals_gain, killer_name]
	_spawn_float_label(msg, COLOR_CRYSTAL)

	# 進度接近目標時閃爍
	if progress >= 0.8:
		_flash_progress_bar()

## 處理地獄龍大獎事件
func on_crystal_reward(contributors: Array, total_reward: int, message: String) -> void:
	# 全螢幕紫色閃光
	_do_flash(Color(0.5, 0.0, 1.0, 0.7), 0.15)
	await get_tree().create_timer(0.15).timeout
	_do_flash(Color(0.5, 0.0, 1.0, 0.5), 0.1)

	# 重置進度條
	_current_crystals = 0
	_update_progress_bar(0.0)

	# 顯示大獎彈窗
	_show_reward_panel(contributors, total_reward)

	# 大獎後隱藏進度條（冷卻中）
	await get_tree().create_timer(5.0).timeout
	_hide_progress_banner()

# ---- 私有方法 ----

func _show_progress_banner() -> void:
	if _progress_banner.visible:
		return
	_progress_banner.visible = true
	visible = true
	_is_visible = true
	var tween = create_tween()
	_progress_banner.position.y = -48
	tween.tween_property(_progress_banner, "position:y", 0.0, 0.3).set_ease(Tween.EASE_OUT)

func _hide_progress_banner() -> void:
	if not _progress_banner.visible:
		return
	var tween = create_tween()
	tween.tween_property(_progress_banner, "position:y", -48.0, 0.3).set_ease(Tween.EASE_IN)
	tween.tween_callback(func():
		_progress_banner.visible = false
		if not _reward_panel.visible:
			visible = false
	)

func _update_progress_bar(progress: float) -> void:
	var fill_width = 800.0 * clamp(progress, 0.0, 1.0)
	if _progress_fill:
		_progress_fill.size.x = fill_width
	if _progress_label:
		_progress_label.text = "%d/%d" % [_current_crystals, _goal]
	if _crystal_count_label:
		_crystal_count_label.text = "💎 %d" % _current_crystals

	# 進度條顏色：低→紫，高→金
	if _progress_fill:
		if progress >= 0.9:
			_progress_fill.color = COLOR_GOLD
		elif progress >= 0.6:
			_progress_fill.color = COLOR_BRIGHT
		else:
			_progress_fill.color = COLOR_CRYSTAL

func _flash_progress_bar() -> void:
	if not _progress_fill:
		return
	var tween = create_tween()
	tween.tween_property(_progress_fill, "color", Color(1.0, 1.0, 1.0, 1.0), 0.1)
	tween.tween_property(_progress_fill, "color", COLOR_GOLD, 0.15)

func _show_reward_panel(contributors: Array, total_reward: int) -> void:
	# 建立貢獻者文字
	var text = ""
	var rank = 1
	for c in contributors:
		var medal = "🥇" if rank == 1 else ("🥈" if rank == 2 else ("🥉" if rank == 3 else "  "))
		text += "%s %s：%d 💎 → +%d 金幣\n" % [medal, c.get("player_name", "?"), c.get("crystals", 0), c.get("reward", 0)]
		rank += 1
		if rank > 5:
			break

	_reward_contributors.text = text
	_reward_total.text = "🐉 全服總獎勵：%d 金幣" % total_reward

	# 滑入動畫
	_reward_panel.visible = true
	visible = true
	var tween = create_tween()
	_reward_panel.position.x = 1280
	tween.tween_property(_reward_panel, "position:x", 900.0, 0.4).set_ease(Tween.EASE_OUT)

	# 5 秒後滑出
	await get_tree().create_timer(5.0).timeout
	var tween2 = create_tween()
	tween2.tween_property(_reward_panel, "position:x", 1280.0, 0.3).set_ease(Tween.EASE_IN)
	tween2.tween_callback(func():
		_reward_panel.visible = false
		if not _progress_banner.visible:
			visible = false
	)

func _do_flash(color: Color, duration: float) -> void:
	if not _flash_overlay:
		return
	_flash_overlay.color = color
	var tween = create_tween()
	tween.tween_property(_flash_overlay, "color:a", 0.0, duration)

func _spawn_float_label(text: String, color: Color) -> void:
	# 回收舊標籤
	if _float_labels.size() >= MAX_FLOAT_LABELS:
		var old = _float_labels.pop_front()
		if is_instance_valid(old):
			old.queue_free()

	var lbl = Label.new()
	lbl.text = text
	lbl.add_theme_color_override("font_color", color)
	lbl.add_theme_font_size_override("font_size", 18)
	# 隨機位置（畫面中央偏上）
	var rx = randf_range(300, 900)
	var ry = randf_range(60, 120)
	lbl.position = Vector2(rx, ry)
	add_child(lbl)
	_float_labels.append(lbl)

	# 上浮 + 淡出
	var tween = create_tween()
	tween.set_parallel(true)
	tween.tween_property(lbl, "position:y", ry - 50, 1.2)
	tween.tween_property(lbl, "modulate:a", 0.0, 1.2)
	tween.tween_callback(func():
		if is_instance_valid(lbl):
			lbl.queue_free()
			_float_labels.erase(lbl)
	)
