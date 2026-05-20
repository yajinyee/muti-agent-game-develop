# ChainExplosionPanel.gd — 連鎖爆炸通知面板（DAY-088）
# 顯示連鎖爆炸等級、目標數量、總獎勵
# 連鎖爆炸時在畫面中央顯示大特效通知
extends Control

# 連鎖等級顏色
const CHAIN_COLORS = {
	1: Color(1.0, 1.0, 1.0),    # 小連鎖：白色
	2: Color(0.0, 0.75, 1.0),   # 中連鎖：天藍
	3: Color(1.0, 0.84, 0.0),   # 大連鎖：金色
	4: Color(1.0, 0.27, 0.0),   # 超級連鎖：橙紅
}

const CHAIN_SIZES = {
	1: 24,
	2: 28,
	3: 34,
	4: 42,
}

func _ready():
	# 連接 GameManager 訊號
	if GameManager.has_signal("chain_explosion"):
		GameManager.chain_explosion.connect(_on_chain_explosion)

func _on_chain_explosion(data: Dictionary) -> void:
	var level = data.get("level", 1)
	var level_name = data.get("level_name", "連鎖！")
	var chains = data.get("chains", [])
	var total_reward = data.get("total_reward", 0)
	var bonus_mult = data.get("bonus_mult", 1.0)
	var player_id = data.get("player_id", "")

	# 只對觸發玩家顯示詳細通知，其他玩家顯示簡化版
	var is_self = (player_id == GameManager.player_data.get("player_id", ""))

	_show_chain_notify(level, level_name, len(chains), total_reward, bonus_mult, is_self)

	# 對每個被連鎖的目標播放爆炸特效
	for entry in chains:
		var instance_id = entry.get("instance_id", "")
		var mult = entry.get("multiplier", 1.0)
		# 通知 TargetManager 播放連鎖爆炸特效
		if GameManager.has_signal("chain_target_killed"):
			GameManager.emit_signal("chain_target_killed", instance_id, mult)

func _show_chain_notify(level: int, level_name: String, count: int, reward: int, bonus_mult: float64, is_self: bool) -> void:
	var color = CHAIN_COLORS.get(level, Color.WHITE)
	var font_size = CHAIN_SIZES.get(level, 24)

	# 建立通知節點
	var notify = Control.new()
	notify.z_index = 70
	add_child(notify)

	# 連鎖等級標籤（畫面中央偏上）
	var name_lbl = Label.new()
	name_lbl.text = level_name
	name_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	name_lbl.add_theme_font_size_override("font_size", font_size)
	name_lbl.add_theme_color_override("font_color", color)
	name_lbl.add_theme_color_override("font_shadow_color", Color(0, 0, 0, 0.8))
	name_lbl.add_theme_constant_override("shadow_offset_x", 2)
	name_lbl.add_theme_constant_override("shadow_offset_y", 2)
	name_lbl.position = Vector2(540, 280)
	name_lbl.size = Vector2(200, 50)
	notify.add_child(name_lbl)

	# 目標數量 + 獎勵（只對自己顯示）
	if is_self and count > 0:
		var detail_lbl = Label.new()
		detail_lbl.text = "×%d 目標  +%d" % [count, reward]
		detail_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
		detail_lbl.add_theme_font_size_override("font_size", int(font_size * 0.6))
		detail_lbl.add_theme_color_override("font_color", color * 0.9)
		detail_lbl.position = Vector2(520, 280 + font_size + 4)
		detail_lbl.size = Vector2(240, 30)
		notify.add_child(detail_lbl)

		# 倍率加成（中連鎖以上才顯示）
		if bonus_mult > 1.0:
			var mult_lbl = Label.new()
			mult_lbl.text = "連鎖加成 ×%.1f" % bonus_mult
			mult_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
			mult_lbl.add_theme_font_size_override("font_size", int(font_size * 0.5))
			mult_lbl.add_theme_color_override("font_color", Color(1.0, 0.9, 0.3))
			mult_lbl.position = Vector2(520, 280 + font_size + 36)
			mult_lbl.size = Vector2(240, 24)
			notify.add_child(mult_lbl)

	# 動畫：縮放彈入 → 停留 → 上移淡出
	notify.scale = Vector2(0.3, 0.3)
	notify.modulate.a = 0.0

	var tween = notify.create_tween()
	tween.tween_property(notify, "scale", Vector2(1.0, 1.0), 0.15).set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)
	tween.parallel().tween_property(notify, "modulate:a", 1.0, 0.15)
	tween.tween_interval(0.8 if level < 3 else 1.2)
	tween.tween_property(notify, "position:y", notify.position.y - 60, 0.4)
	tween.parallel().tween_property(notify, "modulate:a", 0.0, 0.4)
	tween.tween_callback(func(): if is_instance_valid(notify): notify.queue_free())

	# 超級連鎖：額外全畫面閃光
	if level >= 4:
		_spawn_mega_flash(color)

func _spawn_mega_flash(color: Color) -> void:
	var flash = ColorRect.new()
	flash.set_anchors_preset(Control.PRESET_FULL_RECT)
	flash.color = Color(color.r, color.g, color.b, 0.3)
	flash.z_index = 68
	add_child(flash)

	var tween = flash.create_tween()
	tween.tween_property(flash, "modulate:a", 0.0, 0.3)
	tween.tween_callback(func(): if is_instance_valid(flash): flash.queue_free())
