## Cannon.gd
## 玩家砲台控制（規格書 5章）
## 掛載在 Cannon 節點上

extends Node2D

const CANNON_POSITION = Vector2(640, 630)

const PROJECTILE_SPRITES = {
	"chiikawa": "res://assets/sprites/effects/projectile_chiikawa.png",
	"hachiware": "res://assets/sprites/effects/projectile_hachiware.png",
	"usagi":    "res://assets/sprites/effects/projectile_usagi.png",
}

const CHAR_COLORS = {
	"chiikawa": Color(1.0, 0.6, 0.8),
	"hachiware": Color(0.4, 0.6, 1.0),
	"usagi": Color(1.0, 0.9, 0.2),
}

const VOICE_TEXTS = {
	"chiikawa": "YaDa!",
	"hachiware": "Yagaina!",
	"usagi": "Yaha!"
}

@onready var cannon_sprite: Sprite2D = $CannonSprite
@onready var attack_label: Label = $AttackLabel

func _ready() -> void:
	GameManager.attack_result.connect(_on_attack_result)
	GameManager.reward_received.connect(_on_reward_received)
	GameManager.player_updated.connect(_on_player_updated)

func _process(_delta: float) -> void:
	pass

func _input(event: InputEvent) -> void:
	if not (event is InputEventMouseButton):
		return
	if not (event.button_index == MOUSE_BUTTON_LEFT and event.pressed):
		return
	# 只在可攻擊狀態下處理
	var state = GameManager.current_state
	if state not in ["normal_play", "special_target_event", "boss_battle"]:
		return
	_handle_click(event.position)

func _handle_click(click_pos: Vector2) -> void:
	# 確認父節點存在
	var parent = get_parent()
	if not is_instance_valid(parent):
		return

	var target_manager = parent.get_node_or_null("TargetManager")
	var target_id = ""
	if is_instance_valid(target_manager):
		target_id = target_manager.try_click_target(click_pos)

	if target_id != "":
		NetworkManager.send_lock(target_id)
		if is_instance_valid(target_manager):
			target_manager.show_lock_indicator(target_id)

	NetworkManager.send_attack(target_id, click_pos.x, click_pos.y)

	# 攻擊動畫（不阻塞）
	_fire_projectile(click_pos)

	var char_id = GameManager.player_data.get("character_id", "chiikawa")
	AudioManager.play_attack_by_character(char_id)

func _fire_projectile(target_pos: Vector2) -> void:
	var parent = get_parent()
	if not is_instance_valid(parent):
		return

	var char_id = GameManager.player_data.get("character_id", "chiikawa")

	# 依投注等級取得投射物速度（規格書 6章）
	var bet_level = GameManager.get_bet_level()
	var proj_speed = GameManager.player_data.get("projectile_speed", 700.0)
	if proj_speed <= 0:
		proj_speed = 700.0
	var dist = CANNON_POSITION.distance_to(target_pos)
	var flight_time = clamp(dist / proj_speed, 0.05, 0.25)

	# 建立投射物節點
	var proj := Node2D.new()
	proj.position = CANNON_POSITION

	# 嘗試載入 Sprite，失敗就用 ColorRect
	var sprite_path = PROJECTILE_SPRITES.get(char_id, "")
	if sprite_path != "" and ResourceLoader.exists(sprite_path):
		var s := Sprite2D.new()
		s.texture = load(sprite_path)
		s.texture_filter = CanvasItem.TEXTURE_FILTER_NEAREST
		s.scale = Vector2(1.0, 1.0)  # v2 投射物已是 32x16，不需要放大
		proj.add_child(s)
	else:
		var rect := ColorRect.new()
		rect.size = Vector2(10, 6)
		rect.position = Vector2(-5, -3)
		rect.color = CHAR_COLORS.get(char_id, Color.WHITE)
		proj.add_child(rect)

	# 計算方向（防止零向量）
	var diff = target_pos - CANNON_POSITION
	if diff.length() > 1.0:
		proj.rotation = diff.angle()

	parent.add_child(proj)

	# 飛行動畫（依實際速度計算時間）
	var tween = create_tween()
	tween.tween_property(proj, "position", target_pos, flight_time)
	tween.tween_callback(func():
		# 確認節點還存在才執行
		if is_instance_valid(proj):
			_spawn_hit_effect(target_pos, char_id)
			proj.queue_free()
	)

func _on_attack_result(result: Dictionary) -> void:
	if result.get("is_hit", false):
		_show_hit_flash()
		AudioManager.play_sfx(AudioManager.SFX.HIT)
	if result.get("is_kill", false):
		AudioManager.play_sfx(AudioManager.SFX.KILL)

func _on_reward_received(reward: Dictionary) -> void:
	var multiplier = reward.get("multiplier", 1.0)
	if multiplier < 20:
		return

	var parent = get_parent()
	if not is_instance_valid(parent):
		return

	var char_id = GameManager.player_data.get("character_id", "chiikawa")
	var text = VOICE_TEXTS.get(char_id, "!")
	var color = CHAR_COLORS.get(char_id, Color.WHITE)

	# 語音字卡
	var label := Label.new()
	label.text = text
	label.position = Vector2(580, 520)
	label.add_theme_font_size_override("font_size", 32)
	label.modulate = color
	parent.add_child(label)

	var tween = create_tween()
	tween.tween_property(label, "scale", Vector2(1.4, 1.4), 0.15)
	tween.tween_property(label, "scale", Vector2(1.0, 1.0), 0.1)
	tween.tween_interval(0.4)
	tween.tween_property(label, "modulate:a", 0.0, 0.3)
	tween.tween_callback(func():
		if is_instance_valid(label):
			label.queue_free()
	)

	# 角色跳起
	var tween2 = create_tween()
	tween2.tween_property(self, "position:y", position.y - 18, 0.12)
	tween2.tween_property(self, "position:y", position.y, 0.12)

	AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)

func _show_hit_flash() -> void:
	if not is_instance_valid(cannon_sprite):
		return
	var tween = create_tween()
	tween.tween_property(cannon_sprite, "modulate", Color(2.0, 2.0, 2.0, 1.0), 0.04)
	tween.tween_property(cannon_sprite, "modulate", Color.WHITE, 0.06)

func _spawn_hit_effect(pos: Vector2, char_id: String) -> void:
	var parent = get_parent()
	if not is_instance_valid(parent):
		return

	var path = "res://assets/sprites/effects/hit_" + char_id + ".png"
	var effect := Node2D.new()

	if ResourceLoader.exists(path):
		var s := Sprite2D.new()
		s.texture = load(path)
		s.texture_filter = CanvasItem.TEXTURE_FILTER_NEAREST
		s.scale = Vector2(1.0, 1.0)  # v2 特效已是 48x48，不需要放大
		effect.add_child(s)
	else:
		# 備用：簡單閃光圓圈
		var rect := ColorRect.new()
		rect.size = Vector2(20, 20)
		rect.position = Vector2(-10, -10)
		rect.color = CHAR_COLORS.get(char_id, Color.WHITE)
		effect.add_child(rect)

	effect.position = pos
	parent.add_child(effect)

	var tween = create_tween()
	tween.tween_property(effect, "scale", Vector2(2.5, 2.5), 0.10)
	tween.parallel().tween_property(effect, "modulate:a", 0.0, 0.10)
	tween.tween_callback(func():
		if is_instance_valid(effect):
			effect.queue_free()
	)

func _on_player_updated(data: Dictionary) -> void:
	var char_name = GameManager.get_character_name()
	var color = CHAR_COLORS.get(data.get("character_id", "chiikawa"), Color.WHITE)
	if is_instance_valid(attack_label):
		attack_label.text = char_name
		attack_label.modulate = color
