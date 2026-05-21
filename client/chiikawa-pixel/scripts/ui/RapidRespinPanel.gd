# RapidRespinPanel.gd — Rapid Respin 觸發通知面板（DAY-121）
# 業界依據：Reflex Gaming Big Game Fishing Rapid Riches（2026-05-14）
# Rapid Respin 觸發時全螢幕閃光 + 頂部橫幅通知，連鎖時顯示倍率遞增
extends Control

# 連鎖倍率顏色（依連鎖次數）
const CHAIN_COLORS = [
	Color(0.3, 0.8, 1.0),   # 第1次：天藍（1.0x）
	Color(0.2, 1.0, 0.4),   # 第2次：綠色（1.5x）
	Color(1.0, 0.8, 0.0),   # 第3次：金色（2.0x）
	Color(1.0, 0.4, 0.0),   # 第4次：橙紅（3.0x）
	Color(1.0, 0.2, 0.8),   # 第5次：粉紫（5.0x）
]

# 連鎖倍率標籤
const CHAIN_LABELS = ["⚡ RAPID RESPIN", "🔥 CHAIN x2", "💥 CHAIN x3", "🌟 CHAIN x4", "🔥 MAX CHAIN x5"]

var _banner: Control = null
var _flash_overlay: ColorRect = null

func _ready():
	# 建立全螢幕閃光遮罩（預設隱藏）
	_flash_overlay = ColorRect.new()
	_flash_overlay.size = Vector2(1280, 720)
	_flash_overlay.position = Vector2.ZERO
	_flash_overlay.color = Color(0.3, 0.8, 1.0, 0.0)
	_flash_overlay.z_index = 70
	_flash_overlay.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_flash_overlay)

	# 連接 GameManager 訊號
	if GameManager.has_signal("rapid_respin"):
		GameManager.rapid_respin.connect(_on_rapid_respin)
	if GameManager.has_signal("rapid_respin_end"):
		GameManager.rapid_respin_end.connect(_on_rapid_respin_end)

func _on_rapid_respin(data: Dictionary) -> void:
	var player_name = data.get("player_name", "玩家")
	var chain_count = data.get("chain_count", 0)
	var chain_mult = data.get("chain_mult", 1.0)
	var is_chain = data.get("is_chain", false)
	var icon = data.get("icon", "⚡🔄")
	var player_id = data.get("player_id", "")

	var is_self = (player_id == GameManager.player_data.get("player_id", ""))

	_show_respin_effect(player_name, chain_count, chain_mult, is_chain, icon, is_self)

func _on_rapid_respin_end(data: Dictionary) -> void:
	var player_name = data.get("player_name", "玩家")
	var total_chain = data.get("total_chain", 1)
	var player_id = data.get("player_id", "")
	var is_self = (player_id == GameManager.player_data.get("player_id", ""))

	if is_self and total_chain >= 2:
		_show_chain_end_banner(total_chain)

func _show_respin_effect(player_name: String, chain_count: int, chain_mult: float,
		is_chain: bool, icon: String, is_self: bool) -> void:

	var color_idx = clamp(chain_count, 0, CHAIN_COLORS.size() - 1)
	var color = CHAIN_COLORS[color_idx]
	var label_text = CHAIN_LABELS[color_idx] if chain_count < CHAIN_LABELS.size() else CHAIN_LABELS[-1]

	# 全螢幕閃光效果
	_flash_overlay.color = Color(color.r, color.g, color.b, 0.0)
	var flash_tween = create_tween()
	flash_tween.tween_property(_flash_overlay, "color:a", 0.35, 0.08)
	flash_tween.tween_property(_flash_overlay, "color:a", 0.0, 0.25)

	# 移除舊橫幅
	if is_instance_valid(_banner):
		_banner.queue_free()

	# 建立頂部橫幅
	_banner = Control.new()
	_banner.z_index = 71
	add_child(_banner)

	# 橫幅背景
	var bg = ColorRect.new()
	bg.size = Vector2(1280, 64)
	bg.position = Vector2(0, -64)  # 從頂部外開始
	bg.color = Color(color.r * 0.15, color.g * 0.15, color.b * 0.15, 0.95)
	_banner.add_child(bg)

	# 頂部彩色邊條
	var top_bar = ColorRect.new()
	top_bar.size = Vector2(1280, 4)
	top_bar.position = Vector2(0, 0)
	top_bar.color = color
	bg.add_child(top_bar)

	# 底部彩色邊條
	var bot_bar = ColorRect.new()
	bot_bar.size = Vector2(1280, 4)
	bot_bar.position = Vector2(0, 60)
	bot_bar.color = color
	bg.add_child(bot_bar)

	# 主標題（連鎖類型）
	var title_lbl = Label.new()
	title_lbl.text = label_text
	title_lbl.position = Vector2(40, 8)
	title_lbl.add_theme_font_size_override("font_size", 28)
	title_lbl.add_theme_color_override("font_color", color)
	bg.add_child(title_lbl)

	# 倍率顯示
	var mult_lbl = Label.new()
	mult_lbl.text = "×%.1f" % chain_mult
	mult_lbl.position = Vector2(400, 8)
	mult_lbl.add_theme_font_size_override("font_size", 28)
	mult_lbl.add_theme_color_override("font_color", Color.WHITE)
	bg.add_child(mult_lbl)

	# 玩家名稱
	var name_lbl = Label.new()
	if is_self:
		name_lbl.text = "你觸發了！"
		name_lbl.add_theme_color_override("font_color", Color(1.0, 1.0, 0.3))
	else:
		name_lbl.text = player_name + " 觸發"
		name_lbl.add_theme_color_override("font_color", Color(0.8, 0.8, 0.8))
	name_lbl.position = Vector2(700, 8)
	name_lbl.add_theme_font_size_override("font_size", 22)
	bg.add_child(name_lbl)

	# 連鎖進度指示器（小圓點）
	for i in range(5):
		var dot = ColorRect.new()
		dot.size = Vector2(12, 12)
		dot.position = Vector2(1100 + i * 20, 26)
		if i <= chain_count:
			dot.color = color
		else:
			dot.color = Color(0.3, 0.3, 0.3, 0.8)
		bg.add_child(dot)

	# 滑入動畫
	var slide_tween = create_tween()
	slide_tween.tween_property(bg, "position:y", 0.0, 0.15).set_ease(Tween.EASE_OUT)

	# 自動滑出（3 秒後）
	var timer = get_tree().create_timer(3.0)
	timer.timeout.connect(func():
		if is_instance_valid(bg):
			var out_tween = create_tween()
			out_tween.tween_property(bg, "position:y", -64.0, 0.2).set_ease(Tween.EASE_IN)
			out_tween.tween_callback(func():
				if is_instance_valid(_banner):
					_banner.queue_free()
					_banner = null
			)
	)

	# 自己觸發時額外顯示金色粒子效果
	if is_self:
		_spawn_star_particles(color)

func _show_chain_end_banner(total_chain: int) -> void:
	# 連鎖結束時顯示總結橫幅
	var end_lbl = Label.new()
	end_lbl.text = "🔄 Rapid Respin 連鎖結束！共 %d 次" % total_chain
	end_lbl.position = Vector2(400, 680)
	end_lbl.add_theme_font_size_override("font_size", 20)
	end_lbl.add_theme_color_override("font_color", Color(0.7, 0.7, 0.7))
	end_lbl.z_index = 71
	add_child(end_lbl)

	var tween = create_tween()
	tween.tween_property(end_lbl, "modulate:a", 0.0, 2.0).set_delay(1.5)
	tween.tween_callback(end_lbl.queue_free)

func _spawn_star_particles(color: Color) -> void:
	# 在畫面中央生成 8 個星形粒子
	for i in range(8):
		var star = Label.new()
		star.text = "✦"
		star.add_theme_font_size_override("font_size", 24)
		star.add_theme_color_override("font_color", color)
		star.z_index = 72

		var angle = i * PI / 4.0
		var start_x = 640.0
		var start_y = 360.0
		star.position = Vector2(start_x, start_y)
		add_child(star)

		var end_x = start_x + cos(angle) * 200.0
		var end_y = start_y + sin(angle) * 200.0

		var tween = create_tween()
		tween.set_parallel(true)
		tween.tween_property(star, "position", Vector2(end_x, end_y), 0.6).set_ease(Tween.EASE_OUT)
		tween.tween_property(star, "modulate:a", 0.0, 0.6).set_delay(0.2)
		tween.tween_callback(star.queue_free).set_delay(0.6)
