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
	"hachiware": "尖尖哇嘎乃！",
	"usagi": "Yaha!"
}

# 像素字體
const PIXEL_FONT_PATH = "res://assets/fonts/pixel8.fnt"
var _pixel_font: Font = null

# 資源快取（避免每次射擊/大獎都重新 load）
var _cached_proj_textures: Dictionary = {}  # char_id -> Texture2D
var _cached_rainbow_shader: Shader = null

@onready var cannon_sprite: Sprite2D = $CannonSprite
@onready var attack_label: Label = $AttackLabel

func _ready() -> void:
	GameManager.attack_result.connect(_on_attack_result)
	GameManager.reward_received.connect(_on_reward_received)
	GameManager.player_updated.connect(_on_player_updated)
	# 預載入資源
	if ResourceLoader.exists(PIXEL_FONT_PATH):
		_pixel_font = load(PIXEL_FONT_PATH)
	# 預載入投射物 texture
	for char_id in PROJECTILE_SPRITES:
		var path = PROJECTILE_SPRITES[char_id]
		if ResourceLoader.exists(path):
			_cached_proj_textures[char_id] = load(path)
	# 預載入 rainbow shader
	var rainbow_path = "res://assets/shaders/rainbow_glow.gdshader"
	if ResourceLoader.exists(rainbow_path):
		_cached_rainbow_shader = load(rainbow_path)
	# 初始化 BulletPool（子彈節點加入遊戲場景）
	# 用 call_deferred 確保父節點已完全就緒
	call_deferred("_init_bullet_pool")

func _process(_delta: float) -> void:
	pass

func _init_bullet_pool() -> void:
	var parent = get_parent()
	if is_instance_valid(parent):
		BulletPool.init_pool(parent)

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
	var proj_speed = GameManager.player_data.get("projectile_speed", 700.0)
	if proj_speed <= 0:
		proj_speed = 700.0
	var dist = CANNON_POSITION.distance_to(target_pos)
	var flight_time = clamp(dist / proj_speed, 0.05, 0.25)

	# 從 BulletPool 取得子彈節點（Object Pooling）
	# 子彈永遠在遊戲場景中，Pool 只管理可用清單
	var texture = _cached_proj_textures.get(char_id, null)
	var proj: Node2D = BulletPool.acquire(char_id, texture)

	# 如果是降級節點（Pool 未初始化），加入場景
	if not proj.get_meta("pooled", false):
		parent.add_child(proj)

	proj.position = CANNON_POSITION

	# 如果沒有 texture，顯示備用 ColorRect
	var sprite = proj.get_node_or_null("Sprite")
	if is_instance_valid(sprite) and texture == null:
		sprite.visible = false
		# 確保有備用 ColorRect（只建立一次，之後重用）
		if not proj.has_node("FallbackRect"):
			var rect := ColorRect.new()
			rect.name = "FallbackRect"
			rect.size = Vector2(16, 10)
			rect.position = Vector2(-8, -5)
			proj.add_child(rect)
			var core := ColorRect.new()
			core.name = "FallbackCore"
			core.size = Vector2(8, 6)
			core.position = Vector2(-4, -3)
			core.color = Color.WHITE
			proj.add_child(core)
		var fb = proj.get_node_or_null("FallbackRect")
		if is_instance_valid(fb):
			fb.visible = true
			fb.color = CHAR_COLORS.get(char_id, Color.WHITE)
		var fc = proj.get_node_or_null("FallbackCore")
		if is_instance_valid(fc):
			fc.visible = true
	elif is_instance_valid(sprite):
		# 隱藏備用 ColorRect（如果存在）
		var fb = proj.get_node_or_null("FallbackRect")
		if is_instance_valid(fb):
			fb.visible = false
		var fc = proj.get_node_or_null("FallbackCore")
		if is_instance_valid(fc):
			fc.visible = false

	# 計算方向（防止零向量）
	var diff = target_pos - CANNON_POSITION
	if diff.length() > 1.0:
		proj.rotation = diff.angle()

	# 飛行動畫
	var tween = create_tween()
	BulletPool.register_tween(proj, tween)
	tween.tween_property(proj, "position", target_pos, flight_time)
	# 烏薩奇：旋轉殘影效果（規格書 2章）
	if char_id == "usagi":
		tween.parallel().tween_property(proj, "rotation_degrees", 720.0, flight_time)
	tween.tween_callback(func():
		if is_instance_valid(proj):
			HitEffect.spawn_hit(target_pos, char_id)
			# 歸還到 pool（不 queue_free）
			if proj.get_meta("pooled", false):
				BulletPool.release(proj)
			else:
				proj.queue_free()
	)

	# 拖尾效果
	_spawn_trail(parent, CANNON_POSITION, target_pos, flight_time, CHAR_COLORS.get(char_id, Color.WHITE))

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
	elif char_id == "hachiware":
		# 小八：高舉討伐棒（規格書 2章）— 向上旋轉 -30 度停頓，再回正
		tween2.tween_property(self, "position:y", position.y - 20, 0.10)
		tween2.parallel().tween_property(self, "rotation_degrees", -30.0, 0.10)
		tween2.tween_interval(0.20)  # 停頓（高舉姿勢）
		tween2.tween_property(self, "position:y", position.y, 0.12)
		tween2.parallel().tween_property(self, "rotation_degrees", 0.0, 0.12)
	else:
		# 吉伊卡哇：驚慌跳起（規格書 2章）— 快速跳起 + 輕微搖晃
		tween2.tween_property(self, "position:y", position.y - 18, 0.10)
		tween2.parallel().tween_property(self, "rotation_degrees", 8.0, 0.05)
		tween2.tween_property(self, "rotation_degrees", -8.0, 0.05)
		tween2.tween_property(self, "rotation_degrees", 0.0, 0.05)
		tween2.tween_property(self, "position:y", position.y, 0.12)

	# 大獎特效 + 強烈震動
	HitEffect.spawn_big_win(Vector2(640, 360), multiplier)
	ScreenShake.add_trauma(0.7)

	AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)

	# 大獎彩虹光暈（套用到砲台 Sprite，持續 1.5 秒後移除）
	if is_instance_valid(cannon_sprite) and _cached_rainbow_shader != null:
		var rainbow_mat = ShaderMaterial.new()
		rainbow_mat.shader = _cached_rainbow_shader
		rainbow_mat.set_shader_parameter("glow_intensity", 0.9)
		rainbow_mat.set_shader_parameter("glow_speed", 3.0)
		cannon_sprite.material = rainbow_mat
		# 1.5 秒後移除彩虹效果
		var timer = get_tree().create_timer(1.5)
		timer.timeout.connect(func():
			if is_instance_valid(cannon_sprite):
				cannon_sprite.material = null
		)

func _show_hit_flash() -> void:
	if not is_instance_valid(cannon_sprite):
		return
	# 砲台閃白
	var tween = create_tween()
	tween.tween_property(cannon_sprite, "modulate", Color(2.5, 2.5, 2.5, 1.0), 0.03)
	tween.tween_property(cannon_sprite, "modulate", Color.WHITE, 0.08)
	# 砲台輕微縮放（打擊感）
	var tween2 = create_tween()
	tween2.tween_property(cannon_sprite, "scale", Vector2(1.15, 1.15), 0.03)
	tween2.tween_property(cannon_sprite, "scale", Vector2(1.0, 1.0), 0.08)

func _spawn_hit_effect(pos: Vector2, char_id: String) -> void:
	# 已由 HitEffect autoload 取代，保留空函式避免舊呼叫出錯
	HitEffect.spawn_hit(pos, char_id)

## 子彈拖尾：沿飛行路徑生成漸隱殘影（用 tween interval，比多個 timer 效能更好）
func _spawn_trail(parent: Node, from: Vector2, to: Vector2, duration: float, color: Color) -> void:
	if not is_instance_valid(parent):
		return

	var steps = clamp(int(duration / 0.025), 3, 10)
	var interval = duration * 0.7 / float(steps)

	# 用單一 tween 序列依序生成殘影（避免多個 timer 的 overhead）
	var seq_tween = create_tween()
	for i in steps:
		var t = float(i) / float(steps)
		var trail_pos = from.lerp(to, t)
		var dot_size = lerp(8.0, 3.0, t)
		var dot_alpha = 0.6 * (1.0 - t)

		seq_tween.tween_interval(interval)
		seq_tween.tween_callback(func():
			if not is_instance_valid(parent):
				return
			var dot = ColorRect.new()
			dot.size = Vector2(dot_size, dot_size)
			dot.position = trail_pos - Vector2(dot_size / 2, dot_size / 2)
			dot.color = Color(color.r, color.g, color.b, dot_alpha)
			dot.z_index = 5
			parent.add_child(dot)

			var tw = dot.create_tween()
			tw.tween_property(dot, "modulate:a", 0.0, 0.10)
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
