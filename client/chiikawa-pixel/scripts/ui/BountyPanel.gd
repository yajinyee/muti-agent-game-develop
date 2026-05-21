## BountyPanel.gd
## 全服目標懸賞面板（DAY-137）
## 業界依據：strivecloud.io 2026「social streaks + tiered rewards」
## 玩家可對目標下懸賞，擊破者獲得額外金幣，增加社交互動

extends Control

# ---- 狀態 ----
var _active_bounties: Array = []  # 活躍懸賞列表
var _cooldown_left: int = 0

# ---- 節點 ----
var _banner: Control = null         # 懸賞發布橫幅
var _banner_label: Label = null     # 橫幅文字
var _claim_popup: Control = null    # 懸賞領取彈窗
var _claim_label: Label = null      # 領取文字
var _flash_overlay: ColorRect = null # 全螢幕閃光

# ---- 顏色 ----
const COLOR_BOUNTY = Color(1.0, 0.75, 0.0, 1.0)  # 懸賞金色
const COLOR_CLAIM = Color(0.2, 0.9, 0.3, 1.0)    # 領取綠色
const COLOR_BG = Color(0.0, 0.0, 0.0, 0.8)

func _ready() -> void:
	_build_ui()
	mouse_filter = Control.MOUSE_FILTER_IGNORE

func _build_ui() -> void:
	# 全螢幕閃光
	_flash_overlay = ColorRect.new()
	_flash_overlay.set_anchors_preset(Control.PRESET_FULL_RECT)
	_flash_overlay.color = Color(1.0, 0.75, 0.0, 0.0)
	_flash_overlay.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_flash_overlay)

	# 懸賞發布橫幅（頂部，在競速橫幅下方）
	_banner = Control.new()
	_banner.set_anchors_preset(Control.PRESET_TOP_WIDE)
	_banner.custom_minimum_size = Vector2(0, 44)
	_banner.position = Vector2(0, 60)  # 在競速橫幅下方
	_banner.modulate.a = 0.0  # 初始隱藏
	_banner.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_banner)

	var banner_bg = ColorRect.new()
	banner_bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	banner_bg.color = Color(0.2, 0.15, 0.0, 0.88)
	banner_bg.mouse_filter = Control.MOUSE_FILTER_IGNORE
	_banner.add_child(banner_bg)

	_banner_label = Label.new()
	_banner_label.set_anchors_preset(Control.PRESET_FULL_RECT)
	_banner_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_banner_label.add_theme_color_override("font_color", COLOR_BOUNTY)
	_banner_label.add_theme_font_size_override("font_size", 14)
	_banner_label.mouse_filter = Control.MOUSE_FILTER_IGNORE
	_banner.add_child(_banner_label)

	# 懸賞領取彈窗（中央偏下）
	_claim_popup = Control.new()
	_claim_popup.set_anchors_preset(Control.PRESET_CENTER_BOTTOM)
	_claim_popup.position = Vector2(-150, -120)
	_claim_popup.custom_minimum_size = Vector2(300, 80)
	_claim_popup.modulate.a = 0.0
	_claim_popup.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_claim_popup)

	var claim_bg = ColorRect.new()
	claim_bg.set_anchors_preset(Control.PRESET_FULL_RECT)
	claim_bg.color = COLOR_BG
	claim_bg.mouse_filter = Control.MOUSE_FILTER_IGNORE
	_claim_popup.add_child(claim_bg)

	_claim_label = Label.new()
	_claim_label.set_anchors_preset(Control.PRESET_FULL_RECT)
	_claim_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_claim_label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	_claim_label.add_theme_color_override("font_color", COLOR_CLAIM)
	_claim_label.add_theme_font_size_override("font_size", 16)
	_claim_label.autowrap_mode = TextServer.AUTOWRAP_WORD_SMART
	_claim_label.mouse_filter = Control.MOUSE_FILTER_IGNORE
	_claim_popup.add_child(_claim_label)

# ---- 外部呼叫 ----

## 懸賞發布廣播（全服）
func on_bounty_posted(data: Dictionary) -> void:
	var poster_name: String = data.get("poster_name", "玩家")
	var target_name: String = data.get("target_name", "目標")
	var target_mult: float = data.get("target_mult", 0.0)
	var amount: int = data.get("amount", 0)

	if is_instance_valid(_banner_label):
		_banner_label.text = "💰 %s 對【%s】(×%.0f) 懸賞 %d 金幣！" % [poster_name, target_name, target_mult, amount]

	# 橫幅淡入
	_show_banner()
	# 輕微金色閃光
	_flash_screen(Color(1.0, 0.75, 0.0, 0.2), 0.3)

	# 3 秒後淡出
	var tween = create_tween()
	tween.tween_interval(3.0)
	tween.tween_callback(_hide_banner)

## 懸賞個人領取通知
func on_bounty_claimed(data: Dictionary) -> void:
	var total_amount: int = data.get("total_amount", 0)
	var bounty_count: int = data.get("bounty_count", 1)

	if is_instance_valid(_claim_label):
		if bounty_count > 1:
			_claim_label.text = "💰 你領取了 %d 筆懸賞！\n共獲得 %d 金幣！" % [bounty_count, total_amount]
		else:
			_claim_label.text = "💰 你領取了懸賞！\n獲得 %d 金幣！" % total_amount

	# 彈窗淡入
	_show_claim_popup()
	# 綠色閃光
	_flash_screen(Color(0.2, 0.9, 0.3, 0.35), 0.4)

	# 3 秒後淡出
	var tween = create_tween()
	tween.tween_interval(3.0)
	tween.tween_callback(_hide_claim_popup)

## 懸賞目標擊破廣播（全服）
func on_bounty_killed(data: Dictionary) -> void:
	var killer_name: String = data.get("killer_name", "玩家")
	var target_name: String = data.get("target_name", "目標")
	var total_amount: int = data.get("total_amount", 0)

	if is_instance_valid(_banner_label):
		_banner_label.text = "💰 %s 擊破懸賞目標【%s】！獲得 %d 金幣！" % [killer_name, target_name, total_amount]
		_banner_label.add_theme_color_override("font_color", COLOR_CLAIM)

	_show_banner()
	_flash_screen(Color(0.2, 0.9, 0.3, 0.25), 0.4)

	var tween = create_tween()
	tween.tween_interval(3.0)
	tween.tween_callback(func():
		if is_instance_valid(_banner_label):
			_banner_label.add_theme_color_override("font_color", COLOR_BOUNTY)
		_hide_banner()
	)

## 懸賞過期通知（個人）
func on_bounty_expired(data: Dictionary) -> void:
	var target_name: String = data.get("target_name", "目標")
	var amount: int = data.get("amount", 0)
	var message: String = data.get("message", "")

	# 只顯示退款通知（有 amount 才是退款給自己）
	if amount > 0 and message.contains("退還"):
		if is_instance_valid(_claim_label):
			_claim_label.text = "⏰ 懸賞超時\n【%s】的 %d 金幣已退還" % [target_name, amount]
			_claim_label.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8, 1.0))
		_show_claim_popup()
		var tween = create_tween()
		tween.tween_interval(2.5)
		tween.tween_callback(func():
			if is_instance_valid(_claim_label):
				_claim_label.add_theme_color_override("font_color", COLOR_CLAIM)
			_hide_claim_popup()
		)

# ---- 私有方法 ----

func _show_banner() -> void:
	if not is_instance_valid(_banner):
		return
	var tween = create_tween()
	tween.tween_property(_banner, "modulate:a", 1.0, 0.2)

func _hide_banner() -> void:
	if not is_instance_valid(_banner):
		return
	var tween = create_tween()
	tween.tween_property(_banner, "modulate:a", 0.0, 0.3)

func _show_claim_popup() -> void:
	if not is_instance_valid(_claim_popup):
		return
	_claim_popup.modulate.a = 0.0
	var tween = create_tween()
	tween.tween_property(_claim_popup, "modulate:a", 1.0, 0.2)

func _hide_claim_popup() -> void:
	if not is_instance_valid(_claim_popup):
		return
	var tween = create_tween()
	tween.tween_property(_claim_popup, "modulate:a", 0.0, 0.3)

func _flash_screen(color: Color, duration: float) -> void:
	if not is_instance_valid(_flash_overlay):
		return
	_flash_overlay.color = color
	var tween = create_tween()
	tween.tween_property(_flash_overlay, "color:a", 0.0, duration)
