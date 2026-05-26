## Cannon.gd — 射擊系統
## cannon-agent 負責維護
extends Node2D

const CANNON_POS = Vector2(640, 630)

const CHAR_COLORS = {
	"chiikawa": Color(1.0, 0.6, 0.8),
	"hachiware": Color(0.4, 0.6, 1.0),
	"usagi": Color(1.0, 0.9, 0.2),
}

const VOICE_TEXTS = {
	"chiikawa": "YaDa!",
	"hachiware": "尖尖哇嘎乃！",
	"usagi": "Yaha!",
}

@onready var cannon_sprite: Sprite2D = $CannonSprite
@onready var char_label: Label = $CharLabel

var _auto_timer: float = 0.0

func _ready() -> void:
	GameManager.attack_result.connect(_on_attack_result)
	GameManager.reward_received.connect(_on_reward_received)
	GameManager.player_updated.connect(_on_player_updated)

func _process(delta: float) -> void:
	# AUTO 自動射擊
	if not GameManager.is_auto():
		return
	var state = GameManager.current_state
	if state not in ["normal_play", "boss_battle"]:
		return

	var fire_rate = GameManager.get_fire_rate()
	if fire_rate <= 0:
		fire_rate = 2.0
	_auto_timer += delta
	if _auto_timer < 1.0 / fire_rate:
		return
	_auto_timer = 0.0

	# 找最高價值目標
	var parent = get_parent()
	if not is_instance_valid(parent):
		return
	var tm = parent.get_node_or_null("TargetManager")
	if not is_instance_valid(tm):
		return

	var best_id = ""
	var best_score = -1.0
	for iid in tm._target_nodes:
		var node = tm._target_nodes[iid]
		if not is_instance_valid(node):
			continue
		var mult = node.get_meta("multiplier", 2.0)
		# 評分系統（knowhow-log #71）
		var score = mult * 2.0
		# HP 低的優先（快要擊破）
		var hp_bar = node.get_node_or_null("HPBar")
		var hp_bg = node.get_node_or_null("HPBarBG")
		if is_instance_valid(hp_bar) and is_instance_valid(hp_bg) and hp_bg.size.x > 0:
			var hp_pct = hp_bar.size.x / hp_bg.size.x
			score += (1.0 - hp_pct) * 30.0
		# 快要離開畫面的優先
		if node.position.x < 400:
			score += 20.0
		# BOSS 最優先
		if node.get_meta("target_type", "") == "boss":
			score += 500.0
		if score > best_score:
			best_score = score
			best_id = iid

	if best_id == "":
		return

	var best_node = tm._target_nodes.get(best_id, null)
	if not is_instance_valid(best_node):
		return

	_fire(best_id, best_node.position)

func _input(event: InputEvent) -> void:
	if not (event is InputEventMouseButton):
		return
	if not (event.button_index == MOUSE_BUTTON_LEFT and event.pressed):
		return
	var state = GameManager.current_state
	if state not in ["normal_play", "boss_battle"]:
		return
	_handle_click(event.position)

func _handle_click(click_pos: Vector2) -> void:
	var parent = get_parent()
	if not is_instance_valid(parent):
		return
	var tm = parent.get_node_or_null("TargetManager")
	var target_id = ""
	if is_instance_valid(tm):
		target_id = tm.try_click_target(click_pos)
	_fire(target_id, click_pos)

func _fire(target_id: String, target_pos: Vector2) -> void:
	NetworkManager.send_attack(target_id, target_pos.x, target_pos.y)
	_spawn_projectile(target_pos)
	var char_id = GameManager.get_character_id()
	AudioManager.play_attack_by_character(char_id)
	# 觸發角色攻擊動畫
	var animator = get_node_or_null("CharacterAnimator")
	if is_instance_valid(animator):
		animator.play_attack()

func _spawn_projectile(target_pos: Vector2) -> void:
	var parent = get_parent()
	if not is_instance_valid(parent):
		return

	var char_id = GameManager.get_character_id()
	var color = CHAR_COLORS.get(char_id, Color.WHITE)
	var speed = GameManager.get_projectile_speed()
	if speed <= 0:
		speed = 700.0

	var dist = CANNON_POS.distance_to(target_pos)
	var flight_time = clamp(dist / speed, 0.05, 0.3)

	# 投射物（ColorRect 備用）
	var proj = ColorRect.new()
	proj.size = Vector2(16, 10)
	proj.position = CANNON_POS - proj.size / 2
	proj.color = color
	proj.z_index = 20
	parent.add_child(proj)

	# 方向
	var diff = target_pos - CANNON_POS
	if diff.length() > 1.0:
		proj.rotation = diff.angle()

	# 飛行動畫
	var tween = proj.create_tween()
	tween.tween_property(proj, "position", target_pos - proj.size / 2, flight_time)
	# 烏薩奇旋轉殘影
	if char_id == "usagi":
		tween.parallel().tween_property(proj, "rotation_degrees", 720.0, flight_time)
	tween.tween_callback(func():
		if is_instance_valid(proj):
			HitEffect.spawn_hit(target_pos, char_id)
			proj.queue_free()
	)

	# 拖尾
	_spawn_trail(parent, CANNON_POS, target_pos, flight_time, color)

func _spawn_trail(parent: Node, from: Vector2, to: Vector2, duration: float, color: Color) -> void:
	if not is_instance_valid(parent):
		return
	var steps = clamp(int(duration / 0.03), 3, 8)
	var seq = create_tween()
	for i in steps:
		var t = float(i) / float(steps)
		var pos = from.lerp(to, t)
		var size = lerp(8.0, 3.0, t)
		seq.tween_interval(duration * 0.7 / float(steps))
		seq.tween_callback(func():
			if not is_instance_valid(parent):
				return
			var dot = ColorRect.new()
			dot.size = Vector2(size, size)
			dot.position = pos - Vector2(size/2, size/2)
			dot.color = Color(color.r, color.g, color.b, 0.5 * (1.0 - t))
			dot.z_index = 18
			parent.add_child(dot)
			var tw = dot.create_tween()
			tw.tween_property(dot, "modulate:a", 0.0, 0.1)
			tw.tween_callback(func(): if is_instance_valid(dot): dot.queue_free())
		)

func _on_attack_result(result: Dictionary) -> void:
	if result.get("is_hit", false):
		_show_hit_flash()
		AudioManager.play_sfx(AudioManager.SFX.HIT)
		ScreenShake.add_trauma(0.18)
		HitEffect.hit_stop(0.04)

func _on_reward_received(reward: Dictionary) -> void:
	var mult = reward.get("multiplier", 1.0)
	if mult < 20:
		return
	var char_id = GameManager.get_character_id()
	var text = VOICE_TEXTS.get(char_id, "!")
	var color = CHAR_COLORS.get(char_id, Color.WHITE)

	# 觸發大獎動畫
	var animator = get_node_or_null("CharacterAnimator")
	if is_instance_valid(animator):
		animator.play_bigwin()

	# 語音字卡
	var parent = get_parent()
	if is_instance_valid(parent):
		var label = Label.new()
		label.text = text
		label.position = Vector2(560, 520)
		label.add_theme_font_size_override("font_size", 32)
		label.modulate = color
		label.z_index = 55
		parent.add_child(label)
		var tween = label.create_tween()
		tween.tween_property(label, "scale", Vector2(1.4, 1.4), 0.15)
		tween.tween_property(label, "scale", Vector2(1.0, 1.0), 0.1)
		tween.tween_interval(0.4)
		tween.tween_property(label, "modulate:a", 0.0, 0.3)
		tween.tween_callback(func(): if is_instance_valid(label): label.queue_free())

	# 角色跳起
	var tween2 = create_tween()
	if char_id == "usagi":
		tween2.tween_property(self, "position:y", position.y - 22, 0.10)
		tween2.parallel().tween_property(self, "rotation_degrees", 360.0, 0.25)
		tween2.tween_property(self, "position:y", position.y, 0.10)
		tween2.tween_property(self, "rotation_degrees", 0.0, 0.05)
	elif char_id == "hachiware":
		tween2.tween_property(self, "position:y", position.y - 20, 0.10)
		tween2.parallel().tween_property(self, "rotation_degrees", -30.0, 0.10)
		tween2.tween_interval(0.20)
		tween2.tween_property(self, "position:y", position.y, 0.12)
		tween2.parallel().tween_property(self, "rotation_degrees", 0.0, 0.12)
	else:
		tween2.tween_property(self, "position:y", position.y - 18, 0.10)
		tween2.tween_property(self, "position:y", position.y, 0.12)

func _on_player_updated(data: Dictionary) -> void:
	var char_name = GameManager.get_character_name()
	var color = CHAR_COLORS.get(data.get("character_id", "chiikawa"), Color.WHITE)
	if is_instance_valid(char_label):
		char_label.text = char_name
		char_label.modulate = color

func _show_hit_flash() -> void:
	if not is_instance_valid(cannon_sprite):
		return
	var tween = create_tween()
	tween.tween_property(cannon_sprite, "modulate", Color(2.5, 2.5, 2.5), 0.03)
	tween.tween_property(cannon_sprite, "modulate", Color.WHITE, 0.08)
