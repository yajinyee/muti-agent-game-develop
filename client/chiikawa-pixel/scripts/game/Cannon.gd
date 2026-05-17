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

# 像素字體
const PIXEL_FONT_PATH = "res://assets/fonts/pixel8.fnt"
var _pixel_font: Font = null

@onready var cannon_sprite: Sprite2D = $CannonSprite
@onready var attack_label: Label = $AttackLabel

func _ready() -> void:
	GameManager.attack_result.connect(_on_attack_result)
	GameManager.reward_received.connect(_on_reward_received)
	GameManager.player_updated.connect(_on_player_updated)
	# 載入像素字體
	if ResourceLoader.exists(PIXEL_FONT_PATH):
		_pixel_font = load(PIXEL_FONT_PATH)

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
		s.scale = Vector2(1.0, 1.0)
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

	# 拖尾系統：每隔一段時間生成殘影
	var trail_color = CHAR_COLORS.get(char_id, Color.WHITE)
	var trail_timer = 0.0
	var trail_interval = 0.025  # 每 25ms 一個殘影

	# 飛行動畫（依實際速度計算時間）
	var tween = create_tween()
	tween.tween_property(proj, "position", target_pos, flight_time)
	# 烏薩奇：旋轉殘影效果（規格書 2章：黃色旋轉殘影）
	if char_id == "usagi":
		tween.parallel().tween_property(proj, "rotation_degrees", 720.0, flight_time)
	tween.tween_callback(func():
		if is_instance_valid(proj):
			# 命中特效（使用新的 HitEffect 系統）
			HitEffect.spawn_hit(target_pos, char_id)
			proj.queue_free()
	)

	# 拖尾協程（用 _spawn_trail_step 模擬）
	_spawn_trail(parent, CANNON_POSITION, target_pos, flight_time, trail_color)

func _on_attack_result(result: Dictionary) -> void:
	if result.get("is_hit", false):
		_show_hit_flash()
		AudioManager.play_sfx(AudioManager.SFX.HIT)
		# 命中震動（輕微）
		ScreenShake.add_trauma(0.18)
		# Hit Stop（增加打擊感）
		HitEffect.hit_stop(0.04)
	if result.get("is_kill", false):
		AudioManager.play_sfx(AudioManager.SFX.KILL)
		# 擊殺震動（中等）
		ScreenShake.add_trauma(0.35)

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
	if is_instance_valid(_pixel_font):
		label.add_theme_font_override("font", _pixel_font)
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

	# 角色跳起（規格書 2章：大獎演出）
	var tween2 = create_tween()
	if char_id == "usagi":
		# 烏薩奇：高速旋轉跳起（規格書 2章）
		tween2.tween_property(self, "position:y", position.y - 22, 0.10)
		tween2.parallel().tween_property(self, "rotation_degrees", 360.0, 0.25)
		tween2.tween_property(self, "position:y", position.y, 0.10)
		tween2.tween_property(self, "rotation_degrees", 0.0, 0.05)
	else:
		# 吉伊卡哇/小八：跳起
		tween2.tween_property(self, "position:y", position.y - 18, 0.12)
		tween2.tween_property(self, "position:y", position.y, 0.12)

	# 大獎特效 + 強烈震動
	HitEffect.spawn_big_win(Vector2(640, 360), multiplier)
	ScreenShake.add_trauma(0.7)

	AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)

func _show_hit_flash() -> void:
	if not is_instance_valid(cannon_sprite):
		return
	var tween = create_tween()
	tween.tween_property(cannon_sprite, "modulate", Color(2.0, 2.0, 2.0, 1.0), 0.04)
	tween.tween_property(cannon_sprite, "modulate", Color.WHITE, 0.06)

func _spawn_hit_effect(pos: Vector2, char_id: String) -> void:
	# 已由 HitEffect autoload 取代，保留空函式避免舊呼叫出錯
	HitEffect.spawn_hit(pos, char_id)

## 子彈拖尾：沿飛行路徑生成漸隱殘影
func _spawn_trail(parent: Node, from: Vector2, to: Vector2, duration: float, color: Color) -> void:
	if not is_instance_valid(parent):
		return

	var steps = int(duration / 0.03)  # 每 30ms 一個殘影
	steps = clamp(steps, 2, 8)

	for i in steps:
		var t = float(i) / float(steps)
		var trail_pos = from.lerp(to, t)
		var delay = t * duration * 0.7  # 殘影稍微落後

		# 用 SceneTreeTimer 延遲生成
		var timer = get_tree().create_timer(delay)
		timer.timeout.connect(func():
			if not is_instance_valid(parent):
				return
			var dot = ColorRect.new()
			dot.size = Vector2(5, 5)
			dot.position = trail_pos - Vector2(2.5, 2.5)
			dot.color = Color(color.r, color.g, color.b, 0.5 * (1.0 - t))
			dot.z_index = 5
			parent.add_child(dot)

			var tw = dot.create_tween()
			tw.tween_property(dot, "modulate:a", 0.0, 0.12)
			tw.tween_callback(func():
				if is_instance_valid(dot):
					dot.queue_free()
			)
		)

func _on_player_updated(data: Dictionary) -> void:
	var char_name = GameManager.get_character_name()
	var color = CHAR_COLORS.get(data.get("character_id", "chiikawa"), Color.WHITE)
	if is_instance_valid(attack_label):
		attack_label.text = char_name
		attack_label.modulate = color
