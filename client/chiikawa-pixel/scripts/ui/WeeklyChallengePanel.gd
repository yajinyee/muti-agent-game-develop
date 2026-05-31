## WeeklyChallengePanel.gd — 每週挑戰面板（DAY-346）
## 設計：比每日任務更難，獎勵更豐厚，週一重置
## 功能：顯示每週挑戰進度、完成通知、領取獎勵
extends CanvasLayer

# ── 節點引用 ──────────────────────────────────────────────────
var _panel: Panel
var _title_label: Label
var _challenge_items: Array = []
var _coins_label: Label
var _week_label: Label
var _reset_label: Label
var _close_btn: Button

# ── 狀態 ──────────────────────────────────────────────────────
var _is_visible: bool = false
var _challenge_data: Array = []
var _weekly_coins: int = 0
var _week_key: String = ""
var _reset_at: int = 0

# ── 顏色常數 ──────────────────────────────────────────────────
const COLOR_BG = Color(0.08, 0.05, 0.15, 0.95)
const COLOR_PANEL = Color(0.12, 0.08, 0.22, 1.0)
const COLOR_CHALLENGE_BG = Color(0.15, 0.10, 0.28, 1.0)
const COLOR_CHALLENGE_DONE = Color(0.15, 0.25, 0.1, 1.0)
const COLOR_GOLD = Color(1.0, 0.85, 0.2, 1.0)
const COLOR_ORANGE = Color(1.0, 0.55, 0.1, 1.0)
const COLOR_GREEN = Color(0.3, 0.9, 0.3, 1.0)
const COLOR_GRAY = Color(0.6, 0.6, 0.6, 1.0)
const COLOR_WHITE = Color(1.0, 1.0, 1.0, 1.0)
const COLOR_PURPLE = Color(0.7, 0.4, 1.0, 1.0)

# 難度等級顏色
const TIER_COLORS = [
	Color(0.7, 0.7, 0.7, 1.0),  # Tier 0（未使用）
	Color(0.8, 0.6, 0.2, 1.0),  # Tier 1 — 銅色
	Color(0.7, 0.7, 0.8, 1.0),  # Tier 2 — 銀色
	Color(1.0, 0.85, 0.2, 1.0), # Tier 3 — 金色
]

const TIER_ICONS = ["", "🥉", "🥈", "🥇"]

func _ready() -> void:
	layer = 21  # 比每日任務面板高一層
	_build_ui()
	_connect_signals()
	# 連線後自動請求挑戰狀態
	if GameManager.has_signal("weekly_challenge_update"):
		GameManager.weekly_challenge_update.connect(_on_challenge_update)
	if GameManager.has_signal("weekly_challenge_complete"):
		GameManager.weekly_challenge_complete.connect(_on_challenge_complete)
	# 延遲請求，確保連線完成
	get_tree().create_timer(2.0).timeout.connect(func(): _request_challenges())

func _build_ui() -> void:
	# 主面板（左側滑出，與每日任務面板錯開）
	_panel = Panel.new()
	_panel.size = Vector2(300, 520)
	_panel.position = Vector2(10, 80)
	var style = StyleBoxFlat.new()
	style.bg_color = COLOR_PANEL
	style.corner_radius_top_right = 12
	style.corner_radius_bottom_right = 12
	style.border_color = COLOR_ORANGE
	style.border_width_right = 2
	style.border_width_top = 2
	style.border_width_bottom = 2
	_panel.add_theme_stylebox_override("panel", style)
	add_child(_panel)

	# 標題
	_title_label = Label.new()
	_title_label.text = "🏆 每週挑戰"
	_title_label.position = Vector2(16, 12)
	_title_label.size = Vector2(220, 28)
	_title_label.add_theme_color_override("font_color", COLOR_ORANGE)
	_title_label.add_theme_font_size_override("font_size", 18)
	_panel.add_child(_title_label)

	# 關閉按鈕
	_close_btn = Button.new()
	_close_btn.text = "✕"
	_close_btn.position = Vector2(268, 8)
	_close_btn.size = Vector2(28, 28)
	_close_btn.add_theme_color_override("font_color", COLOR_GRAY)
	_close_btn.pressed.connect(_toggle_panel)
	_panel.add_child(_close_btn)

	# 週次顯示
	_week_label = Label.new()
	_week_label.text = "本週挑戰"
	_week_label.position = Vector2(16, 44)
	_week_label.size = Vector2(268, 18)
	_week_label.add_theme_color_override("font_color", COLOR_PURPLE)
	_week_label.add_theme_font_size_override("font_size", 12)
	_panel.add_child(_week_label)

	# 任務幣顯示
	_coins_label = Label.new()
	_coins_label.text = "🏅 本週任務幣：0"
	_coins_label.position = Vector2(16, 62)
	_coins_label.size = Vector2(268, 20)
	_coins_label.add_theme_color_override("font_color", COLOR_ORANGE)
	_coins_label.add_theme_font_size_override("font_size", 13)
	_panel.add_child(_coins_label)

	# 重置時間
	_reset_label = Label.new()
	_reset_label.text = "週一重置：--天--:--:--"
	_reset_label.position = Vector2(16, 82)
	_reset_label.size = Vector2(268, 16)
	_reset_label.add_theme_color_override("font_color", COLOR_GRAY)
	_reset_label.add_theme_font_size_override("font_size", 10)
	_panel.add_child(_reset_label)

	# 分隔線
	var sep = ColorRect.new()
	sep.position = Vector2(8, 102)
	sep.size = Vector2(284, 1)
	sep.color = COLOR_ORANGE * Color(1, 1, 1, 0.3)
	_panel.add_child(sep)

	# 挑戰項目（5個）
	for i in range(5):
		var item = _create_challenge_item(i)
		_challenge_items.append(item)
		_panel.add_child(item)

	# 預設隱藏
	_panel.visible = false

func _create_challenge_item(index: int) -> Control:
	var container = Control.new()
	container.position = Vector2(8, 108 + index * 80)
	container.size = Vector2(284, 74)

	# 背景
	var bg = ColorRect.new()
	bg.size = Vector2(284, 70)
	bg.color = COLOR_CHALLENGE_BG
	bg.name = "bg"
	container.add_child(bg)

	# 難度圖示
	var tier_label = Label.new()
	tier_label.position = Vector2(6, 6)
	tier_label.size = Vector2(24, 20)
	tier_label.add_theme_font_size_override("font_size", 14)
	tier_label.name = "tier_label"
	container.add_child(tier_label)

	# 挑戰名稱
	var name_label = Label.new()
	name_label.position = Vector2(30, 6)
	name_label.size = Vector2(180, 20)
	name_label.add_theme_font_size_override("font_size", 13)
	name_label.add_theme_color_override("font_color", COLOR_WHITE)
	name_label.name = "name_label"
	container.add_child(name_label)

	# 獎勵標籤
	var reward_label = Label.new()
	reward_label.position = Vector2(210, 6)
	reward_label.size = Vector2(70, 20)
	reward_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_RIGHT
	reward_label.add_theme_font_size_override("font_size", 12)
	reward_label.add_theme_color_override("font_color", COLOR_ORANGE)
	reward_label.name = "reward_label"
	container.add_child(reward_label)

	# 描述
	var desc_label = Label.new()
	desc_label.position = Vector2(8, 26)
	desc_label.size = Vector2(268, 16)
	desc_label.add_theme_font_size_override("font_size", 10)
	desc_label.add_theme_color_override("font_color", COLOR_GRAY)
	desc_label.name = "desc_label"
	container.add_child(desc_label)

	# 進度條背景
	var bar_bg = ColorRect.new()
	bar_bg.position = Vector2(8, 44)
	bar_bg.size = Vector2(268, 8)
	bar_bg.color = Color(0.2, 0.15, 0.3, 1.0)
	bar_bg.name = "bar_bg"
	container.add_child(bar_bg)

	# 進度條
	var bar = ColorRect.new()
	bar.position = Vector2(8, 44)
	bar.size = Vector2(0, 8)
	bar.color = COLOR_PURPLE
	bar.name = "bar"
	container.add_child(bar)

	# 進度文字
	var progress_label = Label.new()
	progress_label.position = Vector2(8, 54)
	progress_label.size = Vector2(200, 14)
	progress_label.add_theme_font_size_override("font_size", 10)
	progress_label.add_theme_color_override("font_color", COLOR_GRAY)
	progress_label.name = "progress_label"
	container.add_child(progress_label)

	# 領取按鈕
	var claim_btn = Button.new()
	claim_btn.position = Vector2(208, 50)
	claim_btn.size = Vector2(68, 18)
	claim_btn.text = "領取"
	claim_btn.visible = false
	claim_btn.name = "claim_btn"
	container.add_child(claim_btn)

	return container

func _connect_signals() -> void:
	pass

func _request_challenges() -> void:
	if GameManager.has_method("request_weekly_challenges"):
		GameManager.request_weekly_challenges()

func _toggle_panel() -> void:
	_is_visible = !_is_visible
	_panel.visible = _is_visible
	if _is_visible:
		_request_challenges()
		var tween = create_tween()
		_panel.modulate.a = 0
		tween.tween_property(_panel, "modulate:a", 1.0, 0.2)

func show_panel() -> void:
	_is_visible = true
	_panel.visible = true
	_request_challenges()

func hide_panel() -> void:
	_is_visible = false
	_panel.visible = false

func _on_challenge_update(data: Dictionary) -> void:
	_challenge_data = data.get("challenges", [])
	_weekly_coins = data.get("weekly_coins", 0)
	_week_key = data.get("week_key", "")
	_reset_at = data.get("reset_at", 0)
	_refresh_ui()

func _refresh_ui() -> void:
	# 更新週次
	if _week_key != "":
		_week_label.text = "📅 " + _week_key + " 週間挑戰"

	# 更新任務幣
	_coins_label.text = "🏅 本週任務幣：%d" % _weekly_coins

	# 更新重置時間
	_update_reset_label()

	# 更新挑戰項目
	for i in range(min(_challenge_data.size(), _challenge_items.size())):
		var challenge = _challenge_data[i]
		var item = _challenge_items[i]
		_update_challenge_item(item, challenge)

func _update_reset_label() -> void:
	if _reset_at > 0:
		var now_ms = Time.get_unix_time_from_system() * 1000
		var diff_sec = int((_reset_at - now_ms) / 1000)
		if diff_sec > 0:
			var days = diff_sec / 86400
			var h = (diff_sec % 86400) / 3600
			var m = (diff_sec % 3600) / 60
			var s = diff_sec % 60
			if days > 0:
				_reset_label.text = "週一重置：%d天%02d:%02d:%02d" % [days, h, m, s]
			else:
				_reset_label.text = "週一重置：%02d:%02d:%02d" % [h, m, s]
		else:
			_reset_label.text = "週一重置：即將重置"

func _update_challenge_item(item: Control, challenge: Dictionary) -> void:
	var tier_label = item.get_node("tier_label")
	var name_label = item.get_node("name_label")
	var desc_label = item.get_node("desc_label")
	var reward_label = item.get_node("reward_label")
	var bar = item.get_node("bar")
	var progress_label = item.get_node("progress_label")
	var claim_btn = item.get_node("claim_btn")
	var bg = item.get_node("bg")

	var challenge_name = challenge.get("name", "")
	var description = challenge.get("description", "")
	var target = challenge.get("target", 1)
	var progress = challenge.get("progress", 0)
	var completed = challenge.get("completed", false)
	var claimed = challenge.get("claimed", false)
	var reward = challenge.get("reward", 0)
	var tier = challenge.get("tier", 1)
	var challenge_id = challenge.get("id", "")

	# 難度圖示
	tier_label.text = TIER_ICONS[clamp(tier, 0, 3)]
	var tier_color = TIER_COLORS[clamp(tier, 0, 3)]

	name_label.text = challenge_name
	name_label.add_theme_color_override("font_color", tier_color)
	desc_label.text = description
	reward_label.text = "🏅%d" % reward

	# 進度條
	var pct = clamp(float(progress) / float(max(target, 1)), 0.0, 1.0)
	bar.size.x = 268 * pct

	if completed and not claimed:
		# 可領取
		bg.color = COLOR_CHALLENGE_DONE
		bar.color = COLOR_GREEN
		progress_label.text = "✅ 完成！點擊領取"
		claim_btn.visible = true
		claim_btn.text = "領取"
		# 連接按鈕（清除舊連接）
		for conn in claim_btn.pressed.get_connections():
			claim_btn.pressed.disconnect(conn["callable"])
		claim_btn.pressed.connect(func(): _claim_challenge(challenge_id))
	elif claimed:
		# 已領取
		bg.color = COLOR_CHALLENGE_BG * Color(0.8, 0.8, 0.8, 1.0)
		bar.color = COLOR_GRAY
		progress_label.text = "✅ 已領取"
		claim_btn.visible = false
	else:
		# 進行中
		bg.color = COLOR_CHALLENGE_BG
		bar.color = COLOR_PURPLE
		progress_label.text = "%d / %d" % [progress, target]
		claim_btn.visible = false

func _claim_challenge(challenge_id: String) -> void:
	if GameManager.has_method("claim_weekly_challenge"):
		GameManager.claim_weekly_challenge(challenge_id)
	if AudioManager:
		AudioManager.play_sfx(AudioManager.SFX.COIN_DROP)

func _on_challenge_complete(data: Dictionary) -> void:
	var challenge_name = data.get("challenge_name", "挑戰")
	var reward = data.get("reward", 0)
	var tier = data.get("tier", 1)
	_show_complete_notification(challenge_name, reward, tier)

func _show_complete_notification(challenge_name: String, reward: int, tier: int) -> void:
	var tier_icon = TIER_ICONS[clamp(tier, 0, 3)]
	var tier_color = TIER_COLORS[clamp(tier, 0, 3)]

	var notif = Panel.new()
	notif.size = Vector2(300, 64)
	notif.position = Vector2(
		get_viewport().size.x / 2 - 150,
		get_viewport().size.y - 180  # 比每日任務通知高一點
	)
	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.1, 0.08, 0.25, 0.95)
	style.corner_radius_top_left = 8
	style.corner_radius_top_right = 8
	style.corner_radius_bottom_left = 8
	style.corner_radius_bottom_right = 8
	style.border_color = tier_color
	style.border_width_left = 2
	style.border_width_right = 2
	style.border_width_top = 2
	style.border_width_bottom = 2
	notif.add_theme_stylebox_override("panel", style)
	add_child(notif)

	var label = Label.new()
	label.text = "%s 週間挑戰完成：%s\n點擊挑戰面板領取 🏅%d 任務幣" % [tier_icon, challenge_name, reward]
	label.position = Vector2(8, 8)
	label.size = Vector2(284, 48)
	label.add_theme_font_size_override("font_size", 12)
	label.add_theme_color_override("font_color", tier_color)
	notif.add_child(label)

	# 動畫：上升 + 淡出
	var tween = create_tween()
	tween.tween_property(notif, "position:y", notif.position.y - 50, 2.5)
	tween.parallel().tween_property(notif, "modulate:a", 0.0, 2.5).set_delay(1.5)
	tween.tween_callback(notif.queue_free)

func _process(_delta: float) -> void:
	# 每秒更新重置倒數
	if _is_visible and _reset_at > 0:
		_update_reset_label()
