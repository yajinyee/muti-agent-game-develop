## LuckyEventSystem.gd — 幸運特殊魚事件視覺系統
## lucky-panel-agent 負責維護
## 取代 HUD.gd 中的純文字橫幅，提供每個 Lucky 系統獨特的視覺演出
## 使用方式：LuckyEventSystem.show_event(config) 或 LuckyEventSystem.show_banner(text, color)
extends CanvasLayer

# ── 常數 ─────────────────────────────────────────────────────
const SCREEN_W = 1280.0
const SCREEN_H = 720.0

# Lucky 系統配置表（每個系統的視覺主題）
const LUCKY_CONFIGS = {
	"chain_lightning": {
		"icon": "⚡", "title": "連鎖閃電", "color": Color(0.0, 0.9, 1.0),
		"bg_color": Color(0.0, 0.1, 0.2, 0.92), "flash_color": Color(0.0, 0.9, 1.0),
		"shake": 0.4, "flash_times": 2
	},
	"crab_torpedo": {
		"icon": "🦀", "title": "螃蟹魚雷", "color": Color(1.0, 0.4, 0.1),
		"bg_color": Color(0.2, 0.05, 0.0, 0.92), "flash_color": Color(1.0, 0.4, 0.1),
		"shake": 0.35, "flash_times": 1
	},
	"vortex": {
		"icon": "🌀", "title": "渦旋海葵", "color": Color(0.7, 0.3, 1.0),
		"bg_color": Color(0.1, 0.0, 0.2, 0.92), "flash_color": Color(0.5, 0.2, 0.8),
		"shake": 0.3, "flash_times": 1
	},
	"golden_dragon": {
		"icon": "🐉", "title": "黃金龍魚", "color": Color(1.0, 0.85, 0.0),
		"bg_color": Color(0.15, 0.1, 0.0, 0.92), "flash_color": Color(1.0, 0.85, 0.0),
		"shake": 0.5, "flash_times": 3
	},
	"thunder_lobster": {
		"icon": "🦞⚡", "title": "雷霆龍蝦", "color": Color(1.0, 0.3, 0.0),
		"bg_color": Color(0.2, 0.05, 0.0, 0.92), "flash_color": Color(1.0, 0.5, 0.0),
		"shake": 0.45, "flash_times": 2
	},
	"awakened_phoenix": {
		"icon": "🔥", "title": "覺醒鳳凰", "color": Color(1.0, 0.42, 0.21),
		"bg_color": Color(0.2, 0.05, 0.0, 0.92), "flash_color": Color(1.0, 0.6, 0.2),
		"shake": 0.35, "flash_times": 2
	},
	"shockwave_bomb": {
		"icon": "💥", "title": "全場震盪", "color": Color(1.0, 0.27, 0.0),
		"bg_color": Color(0.2, 0.0, 0.0, 0.92), "flash_color": Color(1.0, 0.3, 0.0),
		"shake": 0.7, "flash_times": 3
	},
	"drill_torpedo": {
		"icon": "🚀", "title": "鑽頭魚雷", "color": Color(1.0, 0.55, 0.15),
		"bg_color": Color(0.15, 0.08, 0.0, 0.92), "flash_color": Color(1.0, 0.55, 0.15),
		"shake": 0.4, "flash_times": 1
	},
	"time_freeze": {
		"icon": "❄️", "title": "時間凍結", "color": Color(0.4, 0.85, 1.0),
		"bg_color": Color(0.0, 0.1, 0.2, 0.92), "flash_color": Color(0.4, 0.85, 1.0),
		"shake": 0.3, "flash_times": 2
	},
	"chain_explosion": {
		"icon": "💥🔗", "title": "連鎖爆炸", "color": Color(0.9, 0.2, 0.15),
		"bg_color": Color(0.2, 0.0, 0.0, 0.92), "flash_color": Color(1.0, 0.3, 0.1),
		"shake": 0.4, "flash_times": 2
	},
	"chain_long_king": {
		"icon": "👑", "title": "千龍王輪盤", "color": Color(1.0, 0.85, 0.0),
		"bg_color": Color(0.15, 0.1, 0.0, 0.92), "flash_color": Color(1.0, 0.85, 0.0),
		"shake": 0.5, "flash_times": 3
	},
	"dragon_shotgun": {
		"icon": "🐲", "title": "龍力散彈", "color": Color(0.8, 0.2, 0.9),
		"bg_color": Color(0.1, 0.0, 0.15, 0.92), "flash_color": Color(0.8, 0.2, 0.9),
		"shake": 0.45, "flash_times": 2
	},
	"rocket_cannon": {
		"icon": "🚀🔥", "title": "火箭砲", "color": Color(1.0, 0.3, 0.1),
		"bg_color": Color(0.2, 0.05, 0.0, 0.92), "flash_color": Color(1.0, 0.4, 0.1),
		"shake": 0.5, "flash_times": 2
	},
	"deep_whirlpool": {
		"icon": "🌊", "title": "深海漩渦", "color": Color(0.0, 0.6, 0.9),
		"bg_color": Color(0.0, 0.05, 0.15, 0.92), "flash_color": Color(0.0, 0.6, 0.9),
		"shake": 0.4, "flash_times": 2
	},
	"vampire_mult": {
		"icon": "🧛", "title": "吸血鬼", "color": Color(0.6, 0.0, 0.6),
		"bg_color": Color(0.1, 0.0, 0.1, 0.92), "flash_color": Color(0.7, 0.0, 0.7),
		"shake": 0.35, "flash_times": 2
	},
	"mirror_fish": {
		"icon": "🪞", "title": "鏡像魚", "color": Color(0.88, 0.67, 1.0),
		"bg_color": Color(0.1, 0.05, 0.15, 0.92), "flash_color": Color(0.88, 0.67, 1.0),
		"shake": 0.3, "flash_times": 1
	},
	"golden_rain": {
		"icon": "🌧️✨", "title": "黃金雨", "color": Color(1.0, 0.85, 0.0),
		"bg_color": Color(0.15, 0.1, 0.0, 0.92), "flash_color": Color(1.0, 0.85, 0.0),
		"shake": 0.4, "flash_times": 3
	},
	"freeze_bomb": {
		"icon": "💣❄️", "title": "冰凍炸彈", "color": Color(0.4, 0.85, 1.0),
		"bg_color": Color(0.0, 0.1, 0.2, 0.92), "flash_color": Color(0.4, 0.85, 1.0),
		"shake": 0.5, "flash_times": 2
	},
	"thunder_storm": {
		"icon": "⛈️", "title": "雷暴", "color": Color(0.5, 0.3, 1.0),
		"bg_color": Color(0.05, 0.0, 0.15, 0.92), "flash_color": Color(0.6, 0.4, 1.0),
		"shake": 0.6, "flash_times": 3
	},
	"lucky_wheel": {
		"icon": "🎡", "title": "大轉盤", "color": Color(1.0, 0.5, 0.0),
		"bg_color": Color(0.15, 0.05, 0.0, 0.92), "flash_color": Color(1.0, 0.6, 0.0),
		"shake": 0.4, "flash_times": 2
	},
}

# ── 節點 ─────────────────────────────────────────────────────
var _banner_panel: Control = null
var _banner_icon: Label = null
var _banner_title: Label = null
var _banner_msg: Label = null
var _banner_bg: ColorRect = null
var _banner_accent: ColorRect = null

var _indicator_panel: Control = null
var _indicator_title: Label = null
var _indicator_value: Label = null
var _indicator_bar_bg: ColorRect = null
var _indicator_bar: ColorRect = null

var _announce_queue: Array = []
var _announce_showing: bool = false

func _ready() -> void:
	layer = 62
	_build_banner()
	_build_indicator()

# ── 建立橫幅 ─────────────────────────────────────────────────
func _build_banner() -> void:
	_banner_panel = Control.new()
	_banner_panel.position = Vector2(0, 100)
	_banner_panel.size = Vector2(SCREEN_W, 90)
	_banner_panel.visible = false
	_banner_panel.z_index = 62
	add_child(_banner_panel)

	_banner_bg = ColorRect.new()
	_banner_bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_banner_bg.color = Color(0, 0, 0, 0.88)
	_banner_panel.add_child(_banner_bg)

	# 頂部強調線
	_banner_accent = ColorRect.new()
	_banner_accent.size = Vector2(SCREEN_W, 4)
	_banner_accent.position = Vector2(0, 0)
	_banner_panel.add_child(_banner_accent)

	# 底部強調線
	var accent_bottom = ColorRect.new()
	accent_bottom.name = "AccentBottom"
	accent_bottom.size = Vector2(SCREEN_W, 4)
	accent_bottom.position = Vector2(0, 86)
	_banner_panel.add_child(accent_bottom)

	# 圖示
	_banner_icon = Label.new()
	_banner_icon.position = Vector2(20, 10)
	_banner_icon.size = Vector2(70, 70)
	_banner_icon.add_theme_font_size_override("font_size", 40)
	_banner_icon.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner_icon.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_banner_panel.add_child(_banner_icon)

	# 標題（系統名稱）
	_banner_title = Label.new()
	_banner_title.position = Vector2(100, 8)
	_banner_title.size = Vector2(400, 32)
	_banner_title.add_theme_font_size_override("font_size", 14)
	_banner_title.modulate = Color(0.8, 0.8, 0.8)
	_banner_panel.add_child(_banner_title)

	# 主訊息
	_banner_msg = Label.new()
	_banner_msg.position = Vector2(100, 36)
	_banner_msg.size = Vector2(1100, 44)
	_banner_msg.add_theme_font_size_override("font_size", 26)
	_banner_msg.modulate = Color.WHITE
	_banner_panel.add_child(_banner_msg)

# ── 建立右側指示器 ────────────────────────────────────────────
func _build_indicator() -> void:
	_indicator_panel = Control.new()
	_indicator_panel.position = Vector2(1050, 200)
	_indicator_panel.size = Vector2(220, 80)
	_indicator_panel.visible = false
	_indicator_panel.z_index = 65
	add_child(_indicator_panel)

	var bg = ColorRect.new()
	bg.size = Vector2(220, 80)
	bg.color = Color(0, 0, 0, 0.82)
	_indicator_panel.add_child(bg)

	var border = ColorRect.new()
	border.size = Vector2(220, 3)
	border.color = Color(1.0, 0.85, 0.0)
	_indicator_panel.add_child(border)

	_indicator_title = Label.new()
	_indicator_title.position = Vector2(8, 6)
	_indicator_title.size = Vector2(204, 22)
	_indicator_title.add_theme_font_size_override("font_size", 12)
	_indicator_title.modulate = Color(1.0, 0.85, 0.0)
	_indicator_panel.add_child(_indicator_title)

	_indicator_value = Label.new()
	_indicator_value.position = Vector2(8, 28)
	_indicator_value.size = Vector2(204, 30)
	_indicator_value.add_theme_font_size_override("font_size", 22)
	_indicator_value.modulate = Color.WHITE
	_indicator_panel.add_child(_indicator_value)

	_indicator_bar_bg = ColorRect.new()
	_indicator_bar_bg.position = Vector2(8, 62)
	_indicator_bar_bg.size = Vector2(204, 8)
	_indicator_bar_bg.color = Color(0.2, 0.2, 0.2, 0.8)
	_indicator_panel.add_child(_indicator_bar_bg)

	_indicator_bar = ColorRect.new()
	_indicator_bar.position = Vector2(8, 62)
	_indicator_bar.size = Vector2(204, 8)
	_indicator_bar.color = Color(1.0, 0.85, 0.0)
	_indicator_panel.add_child(_indicator_bar)

# ── 公開 API ─────────────────────────────────────────────────

## 顯示 Lucky 系統觸發橫幅（帶主題色彩和動畫）
## lucky_key: LUCKY_CONFIGS 的 key（如 "chain_lightning"）
## msg: 主要訊息文字
## duration: 顯示時長（秒）
func show_lucky_banner(lucky_key: String, msg: String, duration: float = 2.5) -> void:
	var cfg = LUCKY_CONFIGS.get(lucky_key, {})
	var color = cfg.get("color", Color.WHITE)
	var bg_color = cfg.get("bg_color", Color(0, 0, 0, 0.88))
	var icon = cfg.get("icon", "✨")
	var title = cfg.get("title", lucky_key)
	var flash_color = cfg.get("flash_color", color)
	var flash_times = cfg.get("flash_times", 1)

	_enqueue_banner({
		"icon": icon,
		"title": "【%s】" % title,
		"msg": msg,
		"color": color,
		"bg_color": bg_color,
		"flash_color": flash_color,
		"flash_times": flash_times,
		"duration": duration,
	})

## 顯示一般公告橫幅（不帶 Lucky 主題）
func show_banner(msg: String, color: Color = Color.WHITE, duration: float = 2.5) -> void:
	_enqueue_banner({
		"icon": "📢",
		"title": "公告",
		"msg": msg,
		"color": color,
		"bg_color": Color(0, 0, 0, 0.88),
		"flash_color": color,
		"flash_times": 0,
		"duration": duration,
	})

## 更新右側指示器（顯示進度/計時/倍率等）
func update_indicator(title: String, value: String, bar_pct: float = -1.0, color: Color = Color(1.0, 0.85, 0.0)) -> void:
	if not is_instance_valid(_indicator_panel):
		return
	_indicator_title.text = title
	_indicator_value.text = value
	_indicator_value.modulate = color

	if bar_pct >= 0.0:
		_indicator_bar_bg.visible = true
		_indicator_bar.visible = true
		var w = 204.0 * clamp(bar_pct, 0.0, 1.0)
		_indicator_bar.size.x = w
		# 顏色隨進度：金→橙→紅
		if bar_pct > 0.6:
			_indicator_bar.color = Color(1.0, 0.85, 0.0)
		elif bar_pct > 0.3:
			_indicator_bar.color = Color(1.0, 0.5, 0.1)
		else:
			_indicator_bar.color = Color(1.0, 0.2, 0.2)
	else:
		_indicator_bar_bg.visible = false
		_indicator_bar.visible = false

	if not _indicator_panel.visible:
		_indicator_panel.visible = true
		_indicator_panel.modulate.a = 0.0
		var tween = _indicator_panel.create_tween()
		tween.tween_property(_indicator_panel, "modulate:a", 1.0, 0.2)

## 隱藏指示器
func hide_indicator() -> void:
	if not is_instance_valid(_indicator_panel) or not _indicator_panel.visible:
		return
	var tween = _indicator_panel.create_tween()
	tween.tween_property(_indicator_panel, "modulate:a", 0.0, 0.3)
	tween.tween_callback(func():
		if is_instance_valid(_indicator_panel):
			_indicator_panel.visible = false
	)

## 顯示結算彈窗（從右側滑入）
func show_settle(lines: Array, duration: float = 3.5) -> void:
	var panel = Control.new()
	panel.position = Vector2(SCREEN_W + 10, 280)
	panel.size = Vector2(300, 160)
	panel.z_index = 70
	add_child(panel)

	var bg = ColorRect.new()
	bg.size = panel.size
	bg.color = Color(0.05, 0.05, 0.1, 0.93)
	panel.add_child(bg)

	var border_top = ColorRect.new()
	border_top.size = Vector2(300, 3)
	border_top.color = Color(1.0, 0.85, 0.0)
	panel.add_child(border_top)

	for i in range(lines.size()):
		var line = lines[i]
		var lbl = Label.new()
		lbl.text = line.get("text", "")
		lbl.position = Vector2(12, 14 + i * 34)
		lbl.add_theme_font_size_override("font_size", line.get("size", 16))
		lbl.modulate = line.get("color", Color.WHITE)
		panel.add_child(lbl)

	var tween = panel.create_tween()
	tween.tween_property(panel, "position:x", SCREEN_W - 310.0, 0.3).set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_BACK)
	tween.tween_interval(duration - 0.6)
	tween.tween_property(panel, "position:x", SCREEN_W + 10.0, 0.3).set_ease(Tween.EASE_IN)
	tween.tween_callback(func():
		if is_instance_valid(panel):
			panel.queue_free()
	)

## 全螢幕閃光（觸發 Lucky 系統時）
func fullscreen_flash(color: Color, times: int = 2) -> void:
	var flash = ColorRect.new()
	flash.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	flash.color = Color(color.r, color.g, color.b, 0.0)
	flash.z_index = 90
	add_child(flash)

	var tween = flash.create_tween()
	for i in range(times):
		tween.tween_property(flash, "modulate:a", 0.65, 0.05)
		tween.tween_property(flash, "modulate:a", 0.0, 0.12)
	tween.tween_callback(func():
		if is_instance_valid(flash):
			flash.queue_free()
	)

# ── 內部：橫幅佇列 ───────────────────────────────────────────
func _enqueue_banner(item: Dictionary) -> void:
	_announce_queue.append(item)
	if not _announce_showing:
		_process_queue()

func _process_queue() -> void:
	if _announce_queue.is_empty():
		_announce_showing = false
		return
	_announce_showing = true
	var item = _announce_queue.pop_front()
	_display_banner(item)
	var tween = create_tween()
	tween.tween_interval(item.get("duration", 2.5) + 0.25)
	tween.tween_callback(_process_queue)

func _display_banner(item: Dictionary) -> void:
	if not is_instance_valid(_banner_panel):
		return

	# 設定內容
	_banner_icon.text = item.get("icon", "✨")
	_banner_title.text = item.get("title", "")
	_banner_msg.text = item.get("msg", "")
	_banner_msg.modulate = item.get("color", Color.WHITE)

	var bg_color = item.get("bg_color", Color(0, 0, 0, 0.88))
	_banner_bg.color = bg_color

	var accent_color = item.get("color", Color(1.0, 0.85, 0.0))
	_banner_accent.color = accent_color
	var accent_bottom = _banner_panel.get_node_or_null("AccentBottom")
	if is_instance_valid(accent_bottom):
		accent_bottom.color = accent_color

	_banner_icon.modulate = accent_color

	# 進場動畫（從上方滑入）
	_banner_panel.position.y = 60.0
	_banner_panel.modulate.a = 0.0
	_banner_panel.visible = true

	var duration = item.get("duration", 2.5)
	var tween = _banner_panel.create_tween()
	tween.tween_property(_banner_panel, "position:y", 100.0, 0.18).set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_BACK)
	tween.parallel().tween_property(_banner_panel, "modulate:a", 1.0, 0.18)
	tween.tween_interval(max(0.1, duration - 0.5))
	tween.tween_property(_banner_panel, "modulate:a", 0.0, 0.35)
	tween.tween_callback(func():
		if is_instance_valid(_banner_panel):
			_banner_panel.visible = false
	)

	# 閃光效果
	var flash_times = item.get("flash_times", 0)
	var flash_color = item.get("flash_color", accent_color)
	if flash_times > 0:
		fullscreen_flash(flash_color, flash_times)
