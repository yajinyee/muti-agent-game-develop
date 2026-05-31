## QuestShopPanel.gd — 任務幣兌換商店（DAY-348）
## 玩家用任務幣兌換 BET 加成、金幣、特殊道具
## 靈感來源：BGaming Quests 2026 核心機制
extends CanvasLayer

# ── 節點引用 ──────────────────────────────────────────────────
var _panel: PanelContainer
var _title_label: Label
var _coins_label: Label
var _items_container: VBoxContainer
var _close_btn: Button
var _effects_label: Label
var _status_label: Label

# ── 狀態 ──────────────────────────────────────────────────────
var _quest_coins: int = 0
var _shop_items: Array = []
var _active_effects: Array = []

# ── 常數 ──────────────────────────────────────────────────────
const PANEL_WIDTH = 520
const PANEL_HEIGHT = 580
const ITEM_HEIGHT = 72

func _ready() -> void:
	layer = 22  # 在賽季通行證（layer=23）下方
	_build_ui()
	_connect_signals()
	visible = false

func _build_ui() -> void:
	# 半透明背景遮罩
	var overlay = ColorRect.new()
	overlay.color = Color(0, 0, 0, 0.6)
	overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	add_child(overlay)
	
	# 主面板
	_panel = PanelContainer.new()
	_panel.set_anchors_preset(Control.PRESET_CENTER)
	_panel.custom_minimum_size = Vector2(PANEL_WIDTH, PANEL_HEIGHT)
	_panel.position = Vector2(-PANEL_WIDTH / 2, -PANEL_HEIGHT / 2)
	add_child(_panel)
	
	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.08, 0.06, 0.15, 0.97)
	style.border_color = Color(0.8, 0.6, 0.1, 1.0)
	style.border_width_left = 2
	style.border_width_right = 2
	style.border_width_top = 2
	style.border_width_bottom = 2
	style.corner_radius_top_left = 12
	style.corner_radius_top_right = 12
	style.corner_radius_bottom_left = 12
	style.corner_radius_bottom_right = 12
	_panel.add_theme_stylebox_override("panel", style)
	
	var vbox = VBoxContainer.new()
	vbox.add_theme_constant_override("separation", 8)
	_panel.add_child(vbox)
	
	# ── 標題列 ──────────────────────────────────────────────────
	var header = HBoxContainer.new()
	vbox.add_child(header)
	
	_title_label = Label.new()
	_title_label.text = "🛒 任務幣兌換商店"
	_title_label.add_theme_font_size_override("font_size", 20)
	_title_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.2))
	_title_label.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	header.add_child(_title_label)
	
	_close_btn = Button.new()
	_close_btn.text = "✕"
	_close_btn.custom_minimum_size = Vector2(36, 36)
	_close_btn.add_theme_font_size_override("font_size", 18)
	header.add_child(_close_btn)
	
	# ── 任務幣顯示 ──────────────────────────────────────────────
	var coins_row = HBoxContainer.new()
	vbox.add_child(coins_row)
	
	_coins_label = Label.new()
	_coins_label.text = "🪙 任務幣：0"
	_coins_label.add_theme_font_size_override("font_size", 16)
	_coins_label.add_theme_color_override("font_color", Color(0.9, 0.8, 0.2))
	_coins_label.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	coins_row.add_child(_coins_label)
	
	# ── 有效效果顯示 ────────────────────────────────────────────
	_effects_label = Label.new()
	_effects_label.text = ""
	_effects_label.add_theme_font_size_override("font_size", 12)
	_effects_label.add_theme_color_override("font_color", Color(0.5, 1.0, 0.5))
	_effects_label.autowrap_mode = TextServer.AUTOWRAP_WORD_SMART
	vbox.add_child(_effects_label)
	
	# 分隔線
	var sep = HSeparator.new()
	vbox.add_child(sep)
	
	# ── 商品列表（可滾動）──────────────────────────────────────
	var scroll = ScrollContainer.new()
	scroll.custom_minimum_size = Vector2(PANEL_WIDTH - 20, PANEL_HEIGHT - 200)
	scroll.size_flags_vertical = Control.SIZE_EXPAND_FILL
	vbox.add_child(scroll)
	
	_items_container = VBoxContainer.new()
	_items_container.add_theme_constant_override("separation", 6)
	_items_container.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	scroll.add_child(_items_container)
	
	# ── 狀態訊息 ────────────────────────────────────────────────
	_status_label = Label.new()
	_status_label.text = ""
	_status_label.add_theme_font_size_override("font_size", 13)
	_status_label.add_theme_color_override("font_color", Color(0.5, 1.0, 0.5))
	_status_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	vbox.add_child(_status_label)

func _connect_signals() -> void:
	_close_btn.pressed.connect(_on_close)
	GameManager.shop_items_received.connect(_on_shop_items_received)
	GameManager.shop_purchase_result.connect(_on_purchase_result)
	GameManager.shop_effect_update.connect(_on_effect_update)

func show_panel() -> void:
	visible = true
	# 請求商店資訊
	GameManager.send_message("shop_request", {})

func _on_close() -> void:
	visible = false

# ── 商品列表渲染 ──────────────────────────────────────────────

func _on_shop_items_received(data: Dictionary) -> void:
	_quest_coins = data.get("quest_coins", 0)
	_shop_items = data.get("items", [])
	_active_effects = data.get("active_effects", [])
	_refresh_ui()

func _refresh_ui() -> void:
	_coins_label.text = "🪙 任務幣：%d" % _quest_coins
	
	# 更新有效效果
	if _active_effects.size() > 0:
		var effect_texts = []
		for effect in _active_effects:
			var item_type = effect.get("item_type", "")
			var value = effect.get("value", 0)
			var expires_in = effect.get("expires_in", -1)
			match item_type:
				"bet_boost":
					effect_texts.append("🎯 BET ×%d 加成（待使用）" % value)
				"xp_boost":
					effect_texts.append("⭐ XP ×2（剩 %ds）" % expires_in)
				"lucky_charm":
					effect_texts.append("🍀 幸運符（剩 %ds）" % expires_in)
				"auto_ammo":
					effect_texts.append("🔫 AUTO 免費（剩 %ds）" % expires_in)
		_effects_label.text = "✨ 有效效果：" + "  ".join(effect_texts)
	else:
		_effects_label.text = ""
	
	# 清空並重建商品列表
	for child in _items_container.get_children():
		child.queue_free()
	
	# 按類型分組顯示
	var categories = {
		"bet_boost": {"name": "🎯 BET 加成", "items": []},
		"coin_bonus": {"name": "🪙 金幣獎勵", "items": []},
		"xp_boost": {"name": "⭐ XP 加速", "items": []},
		"lucky_charm": {"name": "🍀 幸運符", "items": []},
		"auto_ammo": {"name": "🔫 AUTO 彈藥", "items": []},
	}
	
	for item in _shop_items:
		var item_type = item.get("type", "")
		if item_type in categories:
			categories[item_type]["items"].append(item)
	
	for cat_key in ["bet_boost", "coin_bonus", "xp_boost", "lucky_charm", "auto_ammo"]:
		var cat = categories[cat_key]
		if cat["items"].size() == 0:
			continue
		
		# 分類標題
		var cat_label = Label.new()
		cat_label.text = cat["name"]
		cat_label.add_theme_font_size_override("font_size", 14)
		cat_label.add_theme_color_override("font_color", Color(0.7, 0.7, 1.0))
		_items_container.add_child(cat_label)
		
		for item in cat["items"]:
			_items_container.add_child(_create_item_row(item))

func _create_item_row(item: Dictionary) -> Control:
	var row = HBoxContainer.new()
	row.custom_minimum_size = Vector2(0, ITEM_HEIGHT)
	
	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.12, 0.10, 0.22, 0.9)
	style.corner_radius_top_left = 6
	style.corner_radius_top_right = 6
	style.corner_radius_bottom_left = 6
	style.corner_radius_bottom_right = 6
	
	var panel = PanelContainer.new()
	panel.add_theme_stylebox_override("panel", style)
	panel.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	row.add_child(panel)
	
	var inner = HBoxContainer.new()
	inner.add_theme_constant_override("separation", 10)
	panel.add_child(inner)
	
	# 圖示
	var icon_label = Label.new()
	icon_label.text = item.get("icon", "🎁")
	icon_label.add_theme_font_size_override("font_size", 28)
	icon_label.custom_minimum_size = Vector2(40, 0)
	icon_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	inner.add_child(icon_label)
	
	# 名稱和描述
	var info_vbox = VBoxContainer.new()
	info_vbox.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	inner.add_child(info_vbox)
	
	var name_label = Label.new()
	name_label.text = item.get("name", "")
	name_label.add_theme_font_size_override("font_size", 14)
	name_label.add_theme_color_override("font_color", Color(1.0, 0.95, 0.8))
	info_vbox.add_child(name_label)
	
	var desc_label = Label.new()
	desc_label.text = item.get("description", "")
	desc_label.add_theme_font_size_override("font_size", 11)
	desc_label.add_theme_color_override("font_color", Color(0.7, 0.7, 0.7))
	info_vbox.add_child(desc_label)
	
	# 價格和購買按鈕
	var right_vbox = VBoxContainer.new()
	right_vbox.custom_minimum_size = Vector2(90, 0)
	inner.add_child(right_vbox)
	
	var cost = item.get("cost", 0)
	var cost_label = Label.new()
	cost_label.text = "🪙 %d" % cost
	cost_label.add_theme_font_size_override("font_size", 14)
	var can_afford = _quest_coins >= cost
	cost_label.add_theme_color_override("font_color", Color(0.9, 0.8, 0.2) if can_afford else Color(0.5, 0.5, 0.5))
	cost_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	right_vbox.add_child(cost_label)
	
	var buy_btn = Button.new()
	buy_btn.text = "購買" if can_afford else "不足"
	buy_btn.disabled = not can_afford
	buy_btn.custom_minimum_size = Vector2(80, 30)
	buy_btn.add_theme_font_size_override("font_size", 13)
	if can_afford:
		buy_btn.pressed.connect(_on_buy_pressed.bind(item.get("id", "")))
	right_vbox.add_child(buy_btn)
	
	return row

func _on_buy_pressed(item_id: String) -> void:
	GameManager.send_message("shop_purchase", {"item_id": item_id})
	_status_label.text = "購買中..."

func _on_purchase_result(data: Dictionary) -> void:
	var success = data.get("success", false)
	var message = data.get("message", "")
	_quest_coins = data.get("quest_coins", _quest_coins)
	
	if success:
		_status_label.add_theme_color_override("font_color", Color(0.3, 1.0, 0.3))
		_status_label.text = "✅ " + message
		var coin_reward = data.get("coin_reward", 0)
		if coin_reward > 0:
			_status_label.text += "（+%d 金幣）" % coin_reward
		# 重新請求商店資訊（更新任務幣顯示）
		GameManager.send_message("shop_request", {})
	else:
		_status_label.add_theme_color_override("font_color", Color(1.0, 0.3, 0.3))
		_status_label.text = "❌ " + message
	
	# 3秒後清除狀態訊息
	var timer = get_tree().create_timer(3.0)
	timer.timeout.connect(func(): _status_label.text = "")

func _on_effect_update(data: Dictionary) -> void:
	_active_effects = data.get("active_effects", [])
	_refresh_ui()
