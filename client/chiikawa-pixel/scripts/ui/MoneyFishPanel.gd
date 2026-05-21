## MoneyFishPanel.gd — 金幣魚王即時獎勵面板（DAY-162）
## 擊破金幣魚王後立即給予玩家即時獎勵（betLevel × 20-50）
## 視覺：金幣從目標位置噴出 + 右側滑入獎勵彈窗 + 金色閃光
## 業界依據：King of Ocean 2026「Money Fish trigger instant payouts」
extends Node2D

var _pixel_font: Font = null
var _my_player_id: String = ""

func _ready() -> void:
	if ResourceLoader.exists("res://assets/fonts/pixel8.fnt"):
		_pixel_font = load("res://assets/fonts/pixel8.fnt")
	_connect_signals()

func setup(font: Font) -> void:
	if font:
		_pixel_font = font

func _connect_signals() -> void:
	if GameManager.has_signal("money_fish_reward"):
		GameManager.money_fish_reward.connect(_on_money_fish_reward)

## 處理金幣魚王即時獎勵事件
func _on_money_fish_reward(data: Dictionary) -> void:
	var killer_id: String = data.get("killer_id", "")
	var killer_name: String = data.get("killer_name", "")
	var instant_reward: int = data.get("instant_reward", 0)
	var mult_used: int = data.get("mult_used", 0)
	var trigger_x: float = data.get("trigger_x", 640.0)
	var trigger_y: float = data.get("trigger_y", 360.0)

	_my_player_id = NetworkManager.get_player_id() if NetworkManager.has_method("get_player_id") else ""
	var is_me = (killer_id == _my_player_id)

	# 金幣噴射動畫（從目標位置噴出）
	_spawn_coin_burst(trigger_x, trigger_y, instant_reward, is_me)

	# 右側滑入獎勵彈窗
	_show_reward_popup(killer_name, instant_reward, mult_used, is_me)

## 金幣噴射動畫
func _spawn_coin_burst(tx: float, ty: float, reward: int, is_me: bool) -> void:
	# 轉換到本地座標（面板在畫面中心 640,360）
	var local_x = tx - 640.0
	var local_y = ty - 360.0

	# 金色閃光（在目標位置）
	var flash := ColorRect.new()
	flash.position = Vector2(local_x - 30, local_y - 30)
	flash.size = Vector2(60, 60)
	flash.color = Color(1.0, 0.9, 0.0, 0.8)
	add_child(flash)
	var flash_tween = flash.create_tween()
	flash_tween.tween_property(flash, "scale", Vector2(2.0, 2.0), 0.15)
	flash_tween.tween_property(flash, "modulate:a", 0.0, 0.2)
	flash_tween.tween_callback(func():
		if is_instance_valid(flash): flash.queue_free()
	)

	# 噴射 8-12 個金幣粒子
	var coin_count = 8 if not is_me else 12
	for i in range(coin_count):
		var coin := Label.new()
		coin.text = "💰"
		coin.position = Vector2(local_x - 8, local_y - 8)
		coin.add_theme_font_size_override("font_size", 16)
		add_child(coin)

		# 隨機方向噴射
		var angle = (float(i) / float(coin_count)) * TAU + randf() * 0.3
		var speed = 80.0 + randf() * 60.0
		var target_x = local_x + cos(angle) * speed
		var target_y = local_y + sin(angle) * speed

		var coin_tween = coin.create_tween()
		coin_tween.tween_property(coin, "position", Vector2(target_x - 8, target_y - 8), 0.4)
		coin_tween.parallel().tween_property(coin, "modulate:a", 0.0, 0.4)
		coin_tween.tween_callback(func():
			if is_instance_valid(coin): coin.queue_free()
		)

	# 浮動獎勵文字（從目標位置向上飄）
	var reward_lbl := Label.new()
	reward_lbl.text = "+%d 💰" % reward
	reward_lbl.position = Vector2(local_x - 40, local_y - 20)
	reward_lbl.size = Vector2(80, 30)
	reward_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	reward_lbl.add_theme_color_override("font_color", Color(1.0, 0.95, 0.1))
	reward_lbl.add_theme_font_size_override("font_size", 18 if is_me else 14)
	if _pixel_font:
		reward_lbl.add_theme_font_override("font", _pixel_font)
	add_child(reward_lbl)

	var float_tween = reward_lbl.create_tween()
	float_tween.tween_property(reward_lbl, "position:y", local_y - 80, 0.8)
	float_tween.parallel().tween_property(reward_lbl, "modulate:a", 0.0, 0.8)
	float_tween.tween_callback(func():
		if is_instance_valid(reward_lbl): reward_lbl.queue_free()
	)

## 右側滑入獎勵彈窗
func _show_reward_popup(killer_name: String, reward: int, mult: int, is_me: bool) -> void:
	# 彈窗背景
	var popup_bg := ColorRect.new()
	popup_bg.name = "MoneyFishPopup"
	popup_bg.size = Vector2(220, 80)
	popup_bg.color = Color(0.12, 0.10, 0.0, 0.95)
	popup_bg.position = Vector2(660, -40)  # 從右側畫面外開始
	add_child(popup_bg)

	# 標題
	var title_lbl := Label.new()
	title_lbl.text = "💰 金幣魚王！"
	title_lbl.position = Vector2(660, -38)
	title_lbl.size = Vector2(220, 24)
	title_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title_lbl.add_theme_color_override("font_color", Color(1.0, 0.95, 0.1))
	title_lbl.add_theme_font_size_override("font_size", 14)
	if _pixel_font:
		title_lbl.add_theme_font_override("font", _pixel_font)
	add_child(title_lbl)

	# 玩家名稱
	var name_text = "你" if is_me else killer_name
	var name_lbl := Label.new()
	name_lbl.text = "%s 即時獲得" % name_text
	name_lbl.position = Vector2(660, -16)
	name_lbl.size = Vector2(220, 20)
	name_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	name_lbl.add_theme_color_override("font_color", Color(0.9, 0.85, 0.6))
	name_lbl.add_theme_font_size_override("font_size", 12)
	if _pixel_font:
		name_lbl.add_theme_font_override("font", _pixel_font)
	add_child(name_lbl)

	# 獎勵金額
	var reward_lbl := Label.new()
	reward_lbl.text = "+%d 金幣 (×%d)" % [reward, mult]
	reward_lbl.position = Vector2(660, 4)
	reward_lbl.size = Vector2(220, 28)
	reward_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	reward_lbl.add_theme_color_override("font_color", Color(1.0, 0.9, 0.0))
	reward_lbl.add_theme_font_size_override("font_size", 16 if is_me else 14)
	if _pixel_font:
		reward_lbl.add_theme_font_override("font", _pixel_font)
	add_child(reward_lbl)

	# 滑入動畫（從右側滑入）
	var slide_in = popup_bg.create_tween()
	slide_in.tween_property(popup_bg, "position:x", 420.0, 0.25)
	var slide_in2 = title_lbl.create_tween()
	slide_in2.tween_property(title_lbl, "position:x", 420.0, 0.25)
	var slide_in3 = name_lbl.create_tween()
	slide_in3.tween_property(name_lbl, "position:x", 420.0, 0.25)
	var slide_in4 = reward_lbl.create_tween()
	slide_in4.tween_property(reward_lbl, "position:x", 420.0, 0.25)

	# 2.5 秒後淡出
	var fade_tween = popup_bg.create_tween()
	fade_tween.tween_interval(2.5)
	fade_tween.tween_property(popup_bg, "modulate:a", 0.0, 0.4)
	fade_tween.tween_callback(func():
		if is_instance_valid(popup_bg): popup_bg.queue_free()
	)
	for lbl in [title_lbl, name_lbl, reward_lbl]:
		var lbl_fade = lbl.create_tween()
		lbl_fade.tween_interval(2.5)
		lbl_fade.tween_property(lbl, "modulate:a", 0.0, 0.4)
		lbl_fade.tween_callback(func():
			if is_instance_valid(lbl): lbl.queue_free()
		)
