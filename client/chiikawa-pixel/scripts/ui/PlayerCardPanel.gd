## PlayerCardPanel.gd — 玩家名片面板（DAY-106）
## 顯示其他玩家的名片（稱號/VIP/公會/統計亮點）
## 觸發方式：點擊排行榜玩家名稱 / 點擊公會成員 / 點擊好友列表
extends Control

# 面板節點
var _bg: ColorRect
var _title_bar: ColorRect
var _title_label: Label
var _close_btn: Button
var _content: VBoxContainer

# 當前顯示的玩家 ID
var _current_player_id: String = ""

func _ready() -> void:
	_build_ui()
	visible = false
	# 連接 GameManager 訊號
	if GameManager.has_signal("player_card_received"):
		GameManager.player_card_received.connect(_on_player_card_received)

func _build_ui() -> void:
	# 背景遮罩
	var overlay = ColorRect.new()
	overlay.color = Color(0, 0, 0, 0.5)
	overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	overlay.gui_input.connect(func(event):
		if event is InputEventMouseButton and event.pressed:
			hide_card()
	)
	add_child(overlay)

	# 名片容器
	var card = PanelContainer.new()
	card.set_anchors_preset(Control.PRESET_CENTER)
	card.custom_minimum_size = Vector2(380, 480)
	card.position = Vector2(-190, -240)
	add_child(card)

	var vbox = VBoxContainer.new()
	vbox.add_theme_constant_override("separation", 0)
	card.add_child(vbox)

	# 標題列
	_title_bar = ColorRect.new()
	_title_bar.color = Color(0.1, 0.1, 0.2, 1.0)
	_title_bar.custom_minimum_size = Vector2(380, 44)
	vbox.add_child(_title_bar)

	var title_hbox = HBoxContainer.new()
	title_hbox.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_title_bar.add_child(title_hbox)

	_title_label = Label.new()
	_title_label.text = "👤 玩家名片"
	_title_label.add_theme_color_override("font_color", Color.WHITE)
	_title_label.add_theme_font_size_override("font_size", 16)
	_title_label.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	_title_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_title_label.add_theme_constant_override("margin_left", 12)
	title_hbox.add_child(_title_label)

	_close_btn = Button.new()
	_close_btn.text = "✕"
	_close_btn.custom_minimum_size = Vector2(44, 44)
	_close_btn.flat = true
	_close_btn.pressed.connect(hide_card)
	title_hbox.add_child(_close_btn)

	# 內容區域
	var scroll = ScrollContainer.new()
	scroll.custom_minimum_size = Vector2(380, 436)
	vbox.add_child(scroll)

	_content = VBoxContainer.new()
	_content.add_theme_constant_override("separation", 8)
	_content.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	scroll.add_child(_content)

func show_card(player_id: String) -> void:
	_current_player_id = player_id
	visible = true
	# 顯示載入中
	for child in _content.get_children():
		child.queue_free()
	var loading = Label.new()
	loading.text = "載入中..."
	loading.add_theme_color_override("font_color", Color(0.6, 0.6, 0.6))
	loading.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_content.add_child(loading)
	# 發送查詢請求
	NetworkManager.send_get_player_card(player_id)

func hide_card() -> void:
	visible = false
	_current_player_id = ""

func _on_player_card_received(card_data: Dictionary) -> void:
	if card_data.get("player_id", "") != _current_player_id:
		return
	_render_card(card_data)

func _render_card(d: Dictionary) -> void:
	# 清空舊內容
	for child in _content.get_children():
		child.queue_free()

	var margin = MarginContainer.new()
	margin.add_theme_constant_override("margin_left", 12)
	margin.add_theme_constant_override("margin_right", 12)
	margin.add_theme_constant_override("margin_top", 8)
	margin.add_theme_constant_override("margin_bottom", 8)
	_content.add_child(margin)

	var inner = VBoxContainer.new()
	inner.add_theme_constant_override("separation", 10)
	margin.add_child(inner)

	# 玩家名稱 + 稱號
	var name_row = HBoxContainer.new()
	inner.add_child(name_row)

	var name_label = Label.new()
	name_label.text = d.get("display_name", "未知玩家")
	name_label.add_theme_font_size_override("font_size", 20)
	name_label.add_theme_color_override("font_color", Color.WHITE)
	name_row.add_child(name_label)

	if d.get("is_online", false):
		var online_dot = Label.new()
		online_dot.text = " 🟢"
		name_row.add_child(online_dot)

	# 稱號
	var title_name = d.get("title_name", "")
	if title_name != "":
		var title_label = Label.new()
		var title_icon = d.get("title_icon", "")
		title_label.text = title_icon + " " + title_name
		var title_color_str = d.get("title_color", "#FFD700")
		title_label.add_theme_color_override("font_color", Color.from_string(title_color_str, Color.GOLD))
		title_label.add_theme_font_size_override("font_size", 14)
		inner.add_child(title_label)

	# 分隔線
	inner.add_child(_make_separator())

	# VIP + 公會
	var info_grid = GridContainer.new()
	info_grid.columns = 2
	info_grid.add_theme_constant_override("h_separation", 16)
	info_grid.add_theme_constant_override("v_separation", 6)
	inner.add_child(info_grid)

	var vip_level = d.get("vip_level", 0)
	var vip_name = d.get("vip_name", "一般")
	_add_info_row(info_grid, "💎 VIP", "Lv.%d %s" % [vip_level, vip_name])

	var guild_name = d.get("guild_name", "")
	var guild_role = d.get("guild_role", "")
	if guild_name != "":
		var role_map = {"leader": "會長", "officer": "副會長", "member": "成員"}
		var role_str = role_map.get(guild_role, guild_role)
		_add_info_row(info_grid, "⚔️ 公會", "%s（%s）" % [guild_name, role_str])
	else:
		_add_info_row(info_grid, "⚔️ 公會", "未加入")

	_add_info_row(info_grid, "🔥 登入連續", "%d 天" % d.get("login_streak", 0))
	_add_info_row(info_grid, "🏆 成就數", "%d 個" % d.get("achievement_count", 0))

	# 分隔線
	inner.add_child(_make_separator())

	# 統計亮點
	var stats_label = Label.new()
	stats_label.text = "📊 統計亮點"
	stats_label.add_theme_font_size_override("font_size", 14)
	stats_label.add_theme_color_override("font_color", Color(0.7, 0.9, 1.0))
	inner.add_child(stats_label)

	var stats_grid = GridContainer.new()
	stats_grid.columns = 2
	stats_grid.add_theme_constant_override("h_separation", 16)
	stats_grid.add_theme_constant_override("v_separation", 6)
	inner.add_child(stats_grid)

	_add_info_row(stats_grid, "💀 擊破數", "%d" % d.get("kill_count", 0))
	_add_info_row(stats_grid, "💰 最高金幣", _fmt_coins(d.get("max_coins", 0)))
	_add_info_row(stats_grid, "⚡ 最高連擊", "%d 連" % d.get("best_streak", 0))
	_add_info_row(stats_grid, "🎯 最高倍率", "%.1fx" % d.get("best_mult", 0.0))
	_add_info_row(stats_grid, "🎰 Jackpot", "%d 次" % d.get("jackpot_wins", 0))
	var rtp = d.get("rtp", 0.0)
	var rtp_str = "%.1f%%" % (rtp * 100) if rtp > 0 else "—"
	_add_info_row(stats_grid, "📈 RTP", rtp_str)

func _add_info_row(parent: GridContainer, key: String, value: String) -> void:
	var key_label = Label.new()
	key_label.text = key
	key_label.add_theme_color_override("font_color", Color(0.6, 0.6, 0.6))
	key_label.add_theme_font_size_override("font_size", 13)
	parent.add_child(key_label)

	var val_label = Label.new()
	val_label.text = value
	val_label.add_theme_color_override("font_color", Color.WHITE)
	val_label.add_theme_font_size_override("font_size", 13)
	parent.add_child(val_label)

func _make_separator() -> HSeparator:
	var sep = HSeparator.new()
	sep.add_theme_color_override("color", Color(0.3, 0.3, 0.3))
	return sep

func _fmt_coins(n: int) -> String:
	if n >= 1000000:
		return "%.1fM" % (n / 1000000.0)
	elif n >= 1000:
		return "%.1fK" % (n / 1000.0)
	return str(n)
