## ActivityFeedPanel.gd — 成就動態牆面板（DAY-112）
## 顯示全服玩家的重要成就事件，製造社交證明和 FOMO 效應
## 右下角滑入，最多顯示 5 條，自動淡出
extends Node2D

# ---- 常數 ----
const PANEL_WIDTH  := 240
const PANEL_X      := 1040  # 右側，公告面板左邊
const PANEL_Y      := 520   # 底部
const MAX_VISIBLE  := 5     # 最多同時顯示 5 條
const ENTRY_HEIGHT := 44    # 每條高度
const AUTO_HIDE_SEC := 8.0  # 8 秒後自動淡出

# 稀有度顏色
const RARITY_COLORS = {
	"common":    Color(0.7, 0.7, 0.7, 1.0),
	"uncommon":  Color(0.3, 0.9, 0.3, 1.0),
	"rare":      Color(0.3, 0.6, 1.0, 1.0),
	"epic":      Color(0.8, 0.3, 1.0, 1.0),
	"legendary": Color(1.0, 0.85, 0.0, 1.0),
}

# ---- 節點引用 ----
var _container: VBoxContainer
var _pixel_font: Font = null

# ---- 狀態 ----
var _entries: Array = []  # Array of {node, timer}

# ---- 初始化 ----
func _ready() -> void:
	if ResourceLoader.exists("res://assets/fonts/pixel8.fnt"):
		_pixel_font = load("res://assets/fonts/pixel8.fnt")
	_build_ui()
	_connect_signals()

func _build_ui() -> void:
	_container = VBoxContainer.new()
	_container.position = Vector2(PANEL_X, PANEL_Y - MAX_VISIBLE * ENTRY_HEIGHT)
	_container.size = Vector2(PANEL_WIDTH, MAX_VISIBLE * ENTRY_HEIGHT)
	_container.alignment = BoxContainer.ALIGNMENT_END
	add_child(_container)

func _connect_signals() -> void:
	if GameManager.has_signal("activity_feed_event"):
		GameManager.activity_feed_event.connect(_on_feed_event)
	if GameManager.has_signal("activity_feed_history"):
		GameManager.activity_feed_history.connect(_on_feed_history)

# ---- 訊號處理 ----
func _on_feed_event(data: Dictionary) -> void:
	_add_entry(data)

func _on_feed_history(data: Dictionary) -> void:
	# 上線時收到最近 10 條，只顯示最新 3 條（避免刷屏）
	var events: Array = data.get("events", [])
	var show_count := mini(3, events.size())
	for i in range(show_count - 1, -1, -1):
		_add_entry(events[i])

# ---- 新增動態條目 ----
func _add_entry(data: Dictionary) -> void:
	# 超過上限時移除最舊的
	while _entries.size() >= MAX_VISIBLE:
		var oldest = _entries.pop_front()
		if is_instance_valid(oldest.node):
			oldest.node.queue_free()

	var rarity: String = data.get("rarity", "common")
	var icon: String = data.get("icon", "⭐")
	var title: String = data.get("title", "")
	var detail: String = data.get("detail", "")
	var rarity_color: Color = RARITY_COLORS.get(rarity, RARITY_COLORS["common"])

	# 建立條目節點
	var entry_bg := ColorRect.new()
	entry_bg.custom_minimum_size = Vector2(PANEL_WIDTH, ENTRY_HEIGHT - 2)
	entry_bg.color = Color(0.05, 0.05, 0.1, 0.85)
	entry_bg.modulate.a = 0.0  # 初始透明，滑入動畫

	# 左側稀有度邊條
	var rarity_bar := ColorRect.new()
	rarity_bar.position = Vector2(0, 0)
	rarity_bar.size = Vector2(3, ENTRY_HEIGHT - 2)
	rarity_bar.color = rarity_color
	entry_bg.add_child(rarity_bar)

	# 圖示
	var icon_label := Label.new()
	icon_label.position = Vector2(6, 4)
	icon_label.text = icon
	if _pixel_font:
		icon_label.add_theme_font_override("font", _pixel_font)
		icon_label.add_theme_font_size_override("font_size", 16)
	entry_bg.add_child(icon_label)

	# 標題（玩家名稱 + 動作）
	var title_label := Label.new()
	title_label.position = Vector2(28, 3)
	title_label.size = Vector2(PANEL_WIDTH - 32, 18)
	# 截斷過長的標題
	var short_title := title if title.length() <= 22 else title.substr(0, 21) + "…"
	title_label.text = short_title
	title_label.add_theme_color_override("font_color", rarity_color)
	if _pixel_font:
		title_label.add_theme_font_override("font", _pixel_font)
		title_label.add_theme_font_size_override("font_size", 10)
	entry_bg.add_child(title_label)

	# 詳情
	var detail_label := Label.new()
	detail_label.position = Vector2(28, 22)
	detail_label.size = Vector2(PANEL_WIDTH - 32, 16)
	var short_detail := detail if detail.length() <= 24 else detail.substr(0, 23) + "…"
	detail_label.text = short_detail
	detail_label.add_theme_color_override("font_color", Color(0.85, 0.85, 0.85))
	if _pixel_font:
		detail_label.add_theme_font_override("font", _pixel_font)
		detail_label.add_theme_font_size_override("font_size", 9)
	entry_bg.add_child(detail_label)

	_container.add_child(entry_bg)

	# 滑入動畫
	var tween := entry_bg.create_tween()
	tween.tween_property(entry_bg, "modulate:a", 1.0, 0.3)

	# 記錄條目
	_entries.append({"node": entry_bg, "timer": AUTO_HIDE_SEC})

	# 傳說/史詩稀有度加閃光效果
	if rarity == "legendary":
		_add_legendary_glow(entry_bg, rarity_color)
	elif rarity == "epic":
		_add_epic_pulse(entry_bg, rarity_color)

func _add_legendary_glow(node: ColorRect, color: Color) -> void:
	# 金色閃光邊框
	var glow := ColorRect.new()
	glow.position = Vector2(0, 0)
	glow.size = Vector2(PANEL_WIDTH, ENTRY_HEIGHT - 2)
	glow.color = Color(color.r, color.g, color.b, 0.0)
	node.add_child(glow)
	# 閃爍 3 次
	var tween := glow.create_tween().set_loops(3)
	tween.tween_property(glow, "color:a", 0.3, 0.2)
	tween.tween_property(glow, "color:a", 0.0, 0.2)

func _add_epic_pulse(node: ColorRect, color: Color) -> void:
	# 紫色脈衝
	var tween := node.create_tween().set_loops(2)
	tween.tween_property(node, "modulate", Color(1.2, 1.0, 1.2, 1.0), 0.15)
	tween.tween_property(node, "modulate", Color(1.0, 1.0, 1.0, 1.0), 0.15)

# ---- 每幀更新（自動淡出）----
func _process(delta: float) -> void:
	var to_remove: Array = []
	for entry in _entries:
		entry.timer -= delta
		if entry.timer <= 0 and is_instance_valid(entry.node):
			# 淡出動畫
			var tween := entry.node.create_tween()
			tween.tween_property(entry.node, "modulate:a", 0.0, 0.5)
			tween.tween_callback(entry.node.queue_free)
			to_remove.append(entry)
		elif entry.timer <= 1.5 and is_instance_valid(entry.node):
			# 最後 1.5 秒開始半透明
			entry.node.modulate.a = entry.timer / 1.5

	for entry in to_remove:
		_entries.erase(entry)
