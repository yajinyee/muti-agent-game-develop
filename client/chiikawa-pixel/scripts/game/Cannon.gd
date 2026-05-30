## Cannon.gd — 射擊系統
## cannon-agent 負責維護
extends Node2D

const CANNON_POS = Vector2(640, 630)

const CHAR_COLORS = {
	"chiikawa": Color(1.0, 0.6, 0.8),
	"hachiware": Color(0.4, 0.6, 1.0),
	"usagi": Color(1.0, 0.9, 0.2),
}

# DAY-339 其他玩家的投射物顏色（稍微暗一點，區分自己的）
const OTHER_PLAYER_COLORS = {
	"chiikawa": Color(0.8, 0.4, 0.6, 0.7),
	"hachiware": Color(0.2, 0.4, 0.8, 0.7),
	"usagi": Color(0.8, 0.7, 0.1, 0.7),
}

const VOICE_TEXTS = {
	"chiikawa": "YaDa!",
	"hachiware": "尖尖哇嘎乃！",
	"usagi": "Yaha!",
}

@onready var cannon_sprite: Sprite2D = $CannonSprite
@onready var char_label: Label = $CharLabel

var _auto_timer: float = 0.0
# DAY-311 記錄最後射擊位置（用於命中特效）
var _last_fire_pos: Vector2 = Vector2(640, 300)

func _ready() -> void:
	GameManager.attack_result.connect(_on_attack_result)
	GameManager.reward_received.connect(_on_reward_received)
	GameManager.player_updated.connect(_on_player_updated)
	# DAY-339 多人投射物顯示
	GameManager.other_player_attack.connect(_on_other_player_attack)

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
	_last_fire_pos = target_pos  # DAY-311 記錄射擊位置
	NetworkManager.send_attack(target_id, target_pos.x, target_pos.y)
	var char_id = GameManager.get_character_id()
	# DAY-338 打擊感優化：攻擊音效在射擊時播放（即時反饋）
	AudioManager.play_attack_by_character(char_id)
	# 觸發角色攻擊動畫
	var animator = get_node_or_null("CharacterAnimator")
	if is_instance_valid(animator):
		animator.play_attack()
	# DAY-338：投射物到達時觸發命中特效（本地預測，不等 Server 回應）
	_spawn_projectile_with_impact(target_pos, target_id, char_id)

func _spawn_projectile(target_pos: Vector2) -> void:
	_spawn_projectile_with_impact(target_pos, "", GameManager.get_character_id())

## DAY-338 打擊感優化：投射物到達時觸發命中特效（本地預測）
func _spawn_projectile_with_impact(target_pos: Vector2, target_id: String, char_id: String) -> void:
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

	# DAY-311 投射物升級：圓形核心 + 光暈
	var proj_root = Node2D.new()
	proj_root.position = CANNON_POS
	proj_root.z_index = 20
	parent.add_child(proj_root)

	# 外層光暈
	var glow = ColorRect.new()
	glow.size = Vector2(20, 20)
	glow.position = -Vector2(10, 10)
	glow.color = Color(color.r, color.g, color.b, 0.4)
	proj_root.add_child(glow)

	# 核心
	var core = ColorRect.new()
	core.size = Vector2(12, 12)
	core.position = -Vector2(6, 6)
	core.color = color
	proj_root.add_child(core)

	# 高光點
	var highlight = ColorRect.new()
	highlight.size = Vector2(4, 4)
	highlight.position = -Vector2(8, 8)
	highlight.color = Color(1.0, 1.0, 1.0, 0.8)
	proj_root.add_child(highlight)

	# 方向旋轉
	var diff = target_pos - CANNON_POS
	if diff.length() > 1.0:
		proj_root.rotation = diff.angle()

	# 飛行動畫
	var tween = proj_root.create_tween()
	tween.tween_property(proj_root, "position", target_pos, flight_time).set_ease(Tween.EASE_IN).set_trans(Tween.TRANS_QUAD)
	# 烏薩奇旋轉
	if char_id == "usagi":
		tween.parallel().tween_property(proj_root, "rotation_degrees", proj_root.rotation_degrees + 720.0, flight_time)
	# 飛行中縮放脈動
	tween.parallel().tween_property(glow, "modulate:a", 0.8, flight_time * 0.5)
	tween.tween_callback(func():
		if is_instance_valid(proj_root):
			HitEffect.spawn_hit(target_pos, char_id)
			# DAY-338 打擊感優化：投射物到達時播放命中音效 + Hit Stop（本地預測）
			if target_id != "":
				AudioManager.play_sfx(AudioManager.SFX.HIT)
				ScreenShake.add_trauma(0.2)
				HitEffect.hit_stop(0.05)
			proj_root.queue_free()
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
		# DAY-338 打擊感優化：Server 確認命中時，加強震動（本地預測已播放基礎音效）
		# 只在 is_kill 時加強，避免重複音效
		if result.get("is_kill", false):
			ScreenShake.add_trauma(0.35)
			HitEffect.hit_stop(0.07)
		else:
			# 命中但未擊破：輕微額外震動（疊加在本地預測的震動上）
			ScreenShake.add_trauma(0.1)
		# 命中時投射物縮放彈跳（視覺衝擊）
		_spawn_impact_burst(result.get("pos_x", _last_fire_pos.x), result.get("pos_y", _last_fire_pos.y))

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
	# DAY-311 手感強化：更強的白閃 + 彈跳縮放
	tween.tween_property(cannon_sprite, "modulate", Color(3.0, 3.0, 3.0), 0.02)
	tween.parallel().tween_property(cannon_sprite, "scale", Vector2(1.15, 1.15), 0.02)
	tween.tween_property(cannon_sprite, "modulate", Color.WHITE, 0.10)
	tween.parallel().tween_property(cannon_sprite, "scale", Vector2(1.0, 1.0), 0.10)

## DAY-311 命中衝擊爆發（在命中點生成小爆炸）
func _spawn_impact_burst(hit_x: float, hit_y: float) -> void:
	var parent = get_parent()
	if not is_instance_valid(parent):
		return
	var pos = Vector2(hit_x, hit_y)
	var char_id = GameManager.get_character_id()
	var color = CHAR_COLORS.get(char_id, Color.WHITE)
	# 4個方向小粒子
	for i in 4:
		var dot = ColorRect.new()
		dot.size = Vector2(6, 6)
		dot.color = color
		dot.position = pos - dot.size / 2
		dot.z_index = 45
		parent.add_child(dot)
		var angle = (float(i) / 4.0) * TAU + PI / 4.0
		var target = pos + Vector2(cos(angle), sin(angle)) * randf_range(12, 20)
		var tween = dot.create_tween()
		tween.tween_property(dot, "position", target - dot.size / 2, 0.15)
		tween.parallel().tween_property(dot, "scale", Vector2(0.1, 0.1), 0.15)
		tween.parallel().tween_property(dot, "modulate:a", 0.0, 0.15)
		tween.tween_callback(func(): if is_instance_valid(dot): dot.queue_free())

## DAY-339 多人投射物顯示：顯示其他玩家的投射物
func _on_other_player_attack(data: Dictionary) -> void:
	var char_id = data.get("character_id", "chiikawa")
	var target_x = data.get("target_x", 640.0)
	var target_y = data.get("target_y", 300.0)
	var target_pos = Vector2(target_x, target_y)
	var color = OTHER_PLAYER_COLORS.get(char_id, Color(0.7, 0.7, 0.7, 0.6))

	var parent = get_parent()
	if not is_instance_valid(parent):
		return

	var speed = 700.0
	var dist = CANNON_POS.distance_to(target_pos)
	var flight_time = clamp(dist / speed, 0.05, 0.3)

	# 其他玩家的投射物（比自己的小一點，透明度低一點）
	var proj_root = Node2D.new()
	proj_root.position = CANNON_POS
	proj_root.z_index = 15  # 比自己的投射物低一層
	parent.add_child(proj_root)

	# 外層光暈（較小）
	var glow = ColorRect.new()
	glow.size = Vector2(14, 14)
	glow.position = -Vector2(7, 7)
	glow.color = Color(color.r, color.g, color.b, 0.3)
	proj_root.add_child(glow)

	# 核心（較小）
	var core = ColorRect.new()
	core.size = Vector2(8, 8)
	core.position = -Vector2(4, 4)
	core.color = color
	proj_root.add_child(core)

	# 方向旋轉
	var diff = target_pos - CANNON_POS
	if diff.length() > 1.0:
		proj_root.rotation = diff.angle()

	# 飛行動畫
	var tween = proj_root.create_tween()
	tween.tween_property(proj_root, "position", target_pos, flight_time).set_ease(Tween.EASE_IN).set_trans(Tween.TRANS_QUAD)
	# 烏薩奇旋轉
	if char_id == "usagi":
		tween.parallel().tween_property(proj_root, "rotation_degrees", proj_root.rotation_degrees + 720.0, flight_time)
	tween.tween_callback(func():
		if is_instance_valid(proj_root):
			proj_root.queue_free()
	)
