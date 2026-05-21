# LuckyCatchPanel.gd — 幸運捕獲通知面板（DAY-119）
# 業界依據：betway.com Lucky Catch Pick and Win（2026-04）即時獎勵機制
# 幸運捕獲時在畫面右側顯示滑入通知，全服玩家都能看到
extends Control

# 觸發類型顏色
const TRIGGER_COLORS = {
	"streak":   Color(1.0, 0.6, 0.0),   # 連擊觸發：橙色
	"weather":  Color(0.3, 0.8, 1.0),   # 天氣觸發：天藍
	"festival": Color(1.0, 0.3, 0.8),   # 節日觸發：粉紅
}

# 觸發類型名稱
const TRIGGER_NAMES = {
	"streak":   "連擊幸運",
	"weather":  "天氣幸運",
	"festival": "節日幸運",
}

# 最多同時顯示 3 條通知
const MAX_VISIBLE = 3
var _active_notifies: Array = []

func _ready():
	# 連接 GameManager 訊號
	if GameManager.has_signal("lucky_catch"):
		GameManager.lucky_catch.connect(_on_lucky_catch)

func _on_lucky_catch(data: Dictionary) -> void:
	var player_name = data.get("player_name", "玩家")
	var target_name = data.get("target_name", "目標")
	var multiplier = data.get("multiplier", 1.0)
	var bonus_mult = data.get("bonus_mult", 2.0)
	var reward = data.get("reward", 0)
	var trigger_type = data.get("trigger_type", "streak")
	var icon = data.get("icon", "🍀")
	var player_id = data.get("player_id", "")

	var is_self = (player_id == GameManager.player_data.get("player_id", ""))

	_show_lucky_notify(player_name, target_name, multiplier, bonus_mult, reward, trigger_type, icon, is_self)

func _show_lucky_notify(player_name: String, target_name: String, multiplier: float64,
		bonus_mult: float64, reward: int, trigger_type: String, icon: String, is_self: bool) -> void:

	# 超過最大數量時移除最舊的
	if _active_notifies.size() >= MAX_VISIBLE:
		var oldest = _active_notifies.pop_front()
		if is_instance_valid(oldest):
			oldest.queue_free()

	var color = TRIGGER_COLORS.get(trigger_type, Color.WHITE)
	var trigger_name = TRIGGER_NAMES.get(trigger_type, "幸運")

	# 建立通知容器
	var notify = Control.new()
	notify.z_index = 75
	add_child(notify)
	_active_notifies.append(notify)

	# 計算垂直位置（從下往上堆疊）
	var idx = _active_notifies.size() - 1
	var base_y = 580.0 - idx * 90.0

	# 背景面板（圓角矩形效果）
	var bg = ColorRect.new()
	bg.size = Vector2(320, 80)
	bg.position = Vector2(1300, base_y)  # 從右側畫面外開始
	bg.color = Color(0.05, 0.08, 0.18, 0.92)
	notify.add_child(bg)

	# 左側彩色邊條（觸發類型顏色）
	var side_bar = ColorRect.new()
	side_bar.size = Vector2(5, 80)
	side_bar.position = Vector2(0, 0)
	side_bar.color = color
	bg.add_child(side_bar)

	# 圖示標籤
	var icon_lbl = Label.new()
	icon_lbl.text = icon
	icon_lbl.position = Vector2(12, 8)
	icon_lbl.add_theme_font_size_override("font_size", 28)
	bg.add_child(icon_lbl)

	# 觸發類型標籤
	var type_lbl = Label.new()
	type_lbl.text = trigger_name
	type_lbl.position = Vector2(50, 6)
	type_lbl.add_theme_font_size_override("font_size", 11)
	type_lbl.modulate = color
	bg.add_child(type_lbl)

	# 玩家名稱 + 目標名稱
	var main_lbl = Label.new()
	var total_mult = multiplier * bonus_mult
	main_lbl.text = "%s 捕獲了 %s" % [player_name, target_name]
	main_lbl.position = Vector2(50, 22)
	main_lbl.size = Vector2(260, 20)
	main_lbl.add_theme_font_size_override("font_size", 13)
	main_lbl.modulate = Color.WHITE
	bg.add_child(main_lbl)

	# 獎勵標籤
	var reward_lbl = Label.new()
	reward_lbl.text = "%.1fx × %.1fx = +%d 🪙" % [multiplier, bonus_mult, reward]
	reward_lbl.position = Vector2(50, 44)
	reward_lbl.size = Vector2(260, 20)
	reward_lbl.add_theme_font_size_override("font_size", 12)
	reward_lbl.modulate = Color(1.0, 0.9, 0.3)
	bg.add_child(reward_lbl)

	# 自己觸發時加金色邊框閃爍
	if is_self:
		var border = ColorRect.new()
		border.size = Vector2(320, 80)
		border.position = Vector2(0, 0)
		border.color = Color(1.0, 0.85, 0.0, 0.0)
		bg.add_child(border)
		# 閃爍動畫
		var flash_tween = notify.create_tween().set_loops(3)
		flash_tween.tween_property(border, "color:a", 0.4, 0.15)
		flash_tween.tween_property(border, "color:a", 0.0, 0.15)

	# 滑入動畫（從右側畫面外滑入）
	var slide_tween = notify.create_tween()
	slide_tween.tween_property(bg, "position:x", 950.0, 0.3).set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_BACK)

	# 8 秒後淡出並移除
	var timer = notify.create_tween()
	timer.tween_interval(8.0)
	timer.tween_property(bg, "modulate:a", 0.0, 0.5)
	timer.tween_callback(func():
		_active_notifies.erase(notify)
		if is_instance_valid(notify):
			notify.queue_free()
	)
