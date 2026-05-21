## LightningEelPanel.gd — 閃電鰻連鎖攻擊面板（DAY-132）
## 業界依據：JILI Royal Fishing 2026 — 閃電鰻連鎖攻擊，最多跳 5 次
## 顯示連鎖攻擊結果：跳躍路徑、擊破數、總獎勵
## 全服廣播：所有玩家都能看到閃電連鎖效果
extends Control

# ---- 節點引用 ----
var _chain_banner: PanelContainer   # 頂部橫幅（全服廣播）
var _banner_label: Label
var _banner_sub_label: Label
var _chain_result_popup: PanelContainer # 個人結果彈窗
var _result_title_label: Label
var _result_kills_label: Label
var _result_reward_label: Label
var _jump_container: VBoxContainer  # 跳躍列表
var _flash_overlay: ColorRect       # 全螢幕閃電閃光

# ---- 狀態 ----
var _banner_tween: Tween
var _popup_tween: Tween
var _flash_tween: Tween
var _my_player_id: String = ""

const BANNER_DURATION := 4.0
const POPUP_DURATION := 3.5
const FLASH_DURATION := 0.25

func _ready() -> void:
	_build_ui()
	_chain_banner.visible = false
	_chain_result_popup.visible = false
	_flash_overlay.visible = false

func _build_ui() -> void:
	# 全螢幕閃電閃光（最底層）
	_flash_overlay = ColorRect.new()
	_flash_overlay.set_anchors_preset(Control.PRESET_FULL_RECT)
	_flash_overlay.color = Color(1.0, 1.0, 0.3, 0.0)
	_flash_overlay.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_flash_overlay)

	# 頂部橫幅（全服廣播）
	_chain_banner = PanelContainer.new()
	_chain_banner.position = Vector2(0, -80)
	_chain_banner.custom_minimum_size = Vector2(1280, 70)
	add_child(_chain_banner)

	var banner_style = StyleBoxFlat.new()
	banner_style.bg_color = Color(0.05, 0.05, 0.15, 0.92)
	banner_style.border_color = Color(1.0, 0.9, 0.1, 1.0)
	banner_style.set_border_width_all(2)
	_chain_banner.add_theme_stylebox_override("panel", banner_style)

	var banner_vbox = VBoxContainer.new()
	banner_vbox.alignment = BoxContainer.ALIGNMENT_CENTER
	banner_vbox.add_theme_constant_override("separation", 2)
	_chain_banner.add_child(banner_vbox)

	_banner_label = Label.new()
	_banner_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner_label.add_theme_font_size_override("font_size", 22)
	_banner_label.add_theme_color_override("font_color", Color(1.0, 0.95, 0.2))
	_banner_label.text = "⚡ 閃電連鎖！"
	banner_vbox.add_child(_banner_label)

	_banner_sub_label = Label.new()
	_banner_sub_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_banner_sub_label.add_theme_font_size_override("font_size", 14)
	_banner_sub_label.add_theme_color_override("font_color", Color(0.9, 0.9, 1.0))
	_banner_sub_label.text = ""
	banner_vbox.add_child(_banner_sub_label)

	# 個人結果彈窗（右側滑入）
	_chain_result_popup = PanelContainer.new()
	_chain_result_popup.position = Vector2(1280, 200)
	_chain_result_popup.custom_minimum_size = Vector2(220, 160)
	add_child(_chain_result_popup)

	var popup_style = StyleBoxFlat.new()
	popup_style.bg_color = Color(0.05, 0.05, 0.15, 0.92)
	popup_style.border_color = Color(1.0, 0.9, 0.1, 1.0)
	popup_style.set_border_width_all(2)
	popup_style.set_corner_radius_all(8)
	_chain_result_popup.add_theme_stylebox_override("panel", popup_style)

	var popup_vbox = VBoxContainer.new()
	popup_vbox.add_theme_constant_override("separation", 6)
	_chain_result_popup.add_child(popup_vbox)

	_result_title_label = Label.new()
	_result_title_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_result_title_label.add_theme_font_size_override("font_size", 18)
	_result_title_label.add_theme_color_override("font_color", Color(1.0, 0.95, 0.2))
	_result_title_label.text = "⚡ 閃電連鎖！"
	popup_vbox.add_child(_result_title_label)

	_result_kills_label = Label.new()
	_result_kills_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_result_kills_label.add_theme_font_size_override("font_size", 14)
	_result_kills_label.add_theme_color_override("font_color", Color(0.9, 0.9, 1.0))
	_result_kills_label.text = "擊破 0 個目標"
	popup_vbox.add_child(_result_kills_label)

	_result_reward_label = Label.new()
	_result_reward_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_result_reward_label.add_theme_font_size_override("font_size", 20)
	_result_reward_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.1))
	_result_reward_label.text = "🪙 0"
	popup_vbox.add_child(_result_reward_label)

	# 跳躍列表（顯示每次跳躍結果）
	_jump_container = VBoxContainer.new()
	_jump_container.add_theme_constant_override("separation", 2)
	popup_vbox.add_child(_jump_container)

# ---- 公開 API ----

## 設定玩家 ID（用於判斷是否為自己觸發）
func set_player_id(pid: String) -> void:
	_my_player_id = pid

## 顯示連鎖攻擊結果（由 GameManager 呼叫）
func show_chain_result(data: Dictionary) -> void:
	var player_id: String = data.get("player_id", "")
	var player_name: String = data.get("player_name", "???")
	var jumps: Array = data.get("jumps", [])
	var total_kills: int = data.get("total_kills", 0)
	var total_reward: int = data.get("total_reward", 0)
	var is_self: bool = (player_id == _my_player_id)

	# 全螢幕閃電閃光
	_play_flash(is_self)

	# 頂部橫幅（全服可見）
	var icon = "⚡" if total_kills < 5 else "🌩️"
	_banner_label.text = "%s %s 閃電連鎖！" % [icon, player_name]
	_banner_sub_label.text = "連鎖擊破 %d 個目標！獲得 🪙%d！" % [total_kills, total_reward]

	# 橫幅顏色：自己觸發時金色，他人觸發時白色
	if is_self:
		_banner_label.add_theme_color_override("font_color", Color(1.0, 0.95, 0.2))
	else:
		_banner_label.add_theme_color_override("font_color", Color(0.9, 0.9, 1.0))

	_show_banner()

	# 個人結果彈窗（只有自己觸發時顯示）
	if is_self and total_kills > 0:
		_show_personal_result(total_kills, total_reward, jumps)

## 顯示閃電鰻冷卻狀態（登入時呼叫）
func update_status(data: Dictionary) -> void:
	# 目前不顯示冷卻 UI，只記錄狀態
	pass

# ---- 私有方法 ----

func _play_flash(is_self: bool) -> void:
	if _flash_tween:
		_flash_tween.kill()
	_flash_overlay.visible = true
	var alpha = 0.5 if is_self else 0.25
	_flash_overlay.color = Color(1.0, 1.0, 0.3, alpha)
	_flash_tween = create_tween()
	_flash_tween.tween_property(_flash_overlay, "color:a", 0.0, FLASH_DURATION)
	_flash_tween.tween_callback(func(): _flash_overlay.visible = false)

func _show_banner() -> void:
	if _banner_tween:
		_banner_tween.kill()
	_chain_banner.visible = true
	_chain_banner.position.y = -80
	_banner_tween = create_tween()
	# 滑入
	_banner_tween.tween_property(_chain_banner, "position:y", 0.0, 0.3).set_ease(Tween.EASE_OUT)
	# 停留
	_banner_tween.tween_interval(BANNER_DURATION)
	# 滑出
	_banner_tween.tween_property(_chain_banner, "position:y", -80.0, 0.3).set_ease(Tween.EASE_IN)
	_banner_tween.tween_callback(func(): _chain_banner.visible = false)

func _show_personal_result(kills: int, reward: int, jumps: Array) -> void:
	# 清空跳躍列表
	for child in _jump_container.get_children():
		child.queue_free()

	_result_kills_label.text = "擊破 %d 個目標" % kills
	_result_reward_label.text = "🪙 %d" % reward

	# 顯示每次跳躍結果（最多 5 條）
	var show_count = min(jumps.size(), 5)
	for i in range(show_count):
		var jump = jumps[i]
		var jump_label = Label.new()
		jump_label.add_theme_font_size_override("font_size", 11)
		var killed = jump.get("killed", false)
		var target_name = jump.get("target_name", "目標")
		var jump_reward = jump.get("reward", 0)
		if killed:
			jump_label.text = "  ⚡ %s → 🪙%d" % [target_name, jump_reward]
			jump_label.add_theme_color_override("font_color", Color(1.0, 0.9, 0.2))
		else:
			jump_label.text = "  ⚡ %s → 未擊破" % target_name
			jump_label.add_theme_color_override("font_color", Color(0.6, 0.6, 0.8))
		_jump_container.add_child(jump_label)

	# 滑入
	if _popup_tween:
		_popup_tween.kill()
	_chain_result_popup.visible = true
	_chain_result_popup.position.x = 1280
	_popup_tween = create_tween()
	_popup_tween.tween_property(_chain_result_popup, "position:x", 1060.0, 0.3).set_ease(Tween.EASE_OUT)
	_popup_tween.tween_interval(POPUP_DURATION)
	_popup_tween.tween_property(_chain_result_popup, "position:x", 1280.0, 0.3).set_ease(Tween.EASE_IN)
	_popup_tween.tween_callback(func(): _chain_result_popup.visible = false)
