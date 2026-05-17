## TargetManager.gd
## 管理畫面上的目標物節點（使用像素 Sprite）
## 掛載在 GameScene 上

extends Node2D

# 目標 Sprite 路徑對應表
const TARGET_SPRITES = {
	"T001": "res://assets/sprites/targets/T001_grass.png",
	"T002": "res://assets/sprites/targets/T002_bug_g.png",
	"T003": "res://assets/sprites/targets/T003_bug_r.png",
	"T004": "res://assets/sprites/targets/T004_bug_b.png",
	"T005": "res://assets/sprites/targets/T005_pudding.png",
	"T006": "res://assets/sprites/targets/T006_mushroom.png",
	"T101": "res://assets/sprites/targets/T101_mimic.png",
	"T102": "res://assets/sprites/targets/T102_chest.png",
	"T103": "res://assets/sprites/targets/T103_meteor.png",
	"T104": "res://assets/sprites/targets/T104_gold_grass.png",
	"T105": "res://assets/sprites/targets/T105_coin_fish.png",
	"B001": "res://assets/sprites/targets/B001_boss.png",
}

# 命中特效 Sprite
const HIT_EFFECTS = {
	"chiikawa": "res://assets/sprites/effects/hit_chiikawa.png",
	"hachiware": "res://assets/sprites/effects/hit_hachiware.png",
	"usagi":    "res://assets/sprites/effects/hit_usagi.png",
}

# 投射物 Sprite
const PROJECTILE_SPRITES = {
	"chiikawa": "res://assets/sprites/effects/projectile_chiikawa.png",
	"hachiware": "res://assets/sprites/effects/projectile_hachiware.png",
	"usagi":    "res://assets/sprites/effects/projectile_usagi.png",
}

# 目標節點字典
var _target_nodes: Dictionary = {}  # instance_id -> Node2D

func _ready() -> void:
	GameManager.target_spawned.connect(_on_target_spawned)
	GameManager.target_updated.connect(_on_target_updated)
	GameManager.target_killed.connect(_on_target_killed)
	GameManager.boss_event.connect(_on_boss_event)

## 目標 HP 更新
func _on_target_updated(data: Dictionary) -> void:
	var instance_id = data.get("instance_id", "")
	update_target_hp(instance_id, data.get("hp", 0), data.get("max_hp", 1))

	# T102 寶箱怪：受擊後加速逃跑（規格書 26.2）
	if data.get("is_fleeing", false):
		if _target_nodes.has(instance_id):
			var node = _target_nodes[instance_id]
			if is_instance_valid(node):
				var current_speed = node.get_meta("speed", 70.0)
				node.set_meta("flee_speed", current_speed * 2.5)  # 2.5x 加速
				node.set_meta("behavior", "flee")
				# 視覺反饋：閃爍紅色
				var tween = create_tween()
				tween.tween_property(node, "modulate", Color(2.0, 0.5, 0.5, 1.0), 0.06)
				tween.tween_property(node, "modulate", Color.WHITE, 0.06)
				tween.tween_property(node, "modulate", Color(2.0, 0.5, 0.5, 1.0), 0.06)
				tween.tween_property(node, "modulate", Color.WHITE, 0.06)

## BOSS 事件處理（Phase 2 視覺變化）
func _on_boss_event(event_data: Dictionary) -> void:
	var event = event_data.get("event", "")

	# BOSS 登場：全畫面特效 + 強烈震動
	if event == "boss_enter":
		HitEffect.spawn_boss_enter()
		ScreenShake.add_trauma(0.9)
		return

	if event != "phase_change":
		return

	var instance_id = event_data.get("instance_id", "")
	if not _target_nodes.has(instance_id):
		return

	var node = _target_nodes[instance_id]
	if not is_instance_valid(node):
		return

	# Phase 2：BOSS 變紅 + 閃爍 + 放大（規格書 28.2）
	var tween = create_tween()
	# 閃爍 3 次
	for i in 3:
		tween.tween_property(node, "modulate", Color(3.0, 0.3, 0.3, 1.0), 0.08)
		tween.tween_property(node, "modulate", Color.WHITE, 0.08)
	# 最終變成紅色調（Phase 2 持續視覺）
	tween.tween_property(node, "modulate", Color(1.5, 0.5, 0.5, 1.0), 0.1)

	# 放大 10%（Phase 2 更有威脅感）
	var tween2 = create_tween()
	tween2.tween_property(node, "scale", Vector2(2.2, 2.2), 0.3)

	# Phase 2 震動
	ScreenShake.add_trauma(0.6)

	# 顯示 Phase 2 警告文字
	var phase_label = Label.new()
	phase_label.text = "PHASE 2!"
	phase_label.position = node.position + Vector2(-30, -80)
	phase_label.add_theme_font_size_override("font_size", 22)
	phase_label.modulate = Color(1.0, 0.2, 0.2)
	add_child(phase_label)
	var tween3 = create_tween()
	tween3.tween_property(phase_label, "scale", Vector2(1.5, 1.5), 0.2)
	tween3.tween_property(phase_label, "scale", Vector2(1.0, 1.0), 0.1)
	tween3.tween_interval(0.5)
	tween3.tween_property(phase_label, "modulate:a", 0.0, 0.4)
	tween3.tween_callback(func():
		if is_instance_valid(phase_label):
			phase_label.queue_free()
	)

func _process(delta: float) -> void:
	_update_target_positions(delta)

## 目標生成
func _on_target_spawned(data: Dictionary) -> void:
	var instance_id = data.get("instance_id", "")
	if _target_nodes.has(instance_id):
		return

	# 建立目標節點（暫用 ColorRect 代替 Sprite）
	var node = _create_target_node(data)
	add_child(node)
	_target_nodes[instance_id] = node

## 目標擊破
func _on_target_killed(data: Dictionary) -> void:
	var instance_id = data.get("instance_id", "")
	if not _target_nodes.has(instance_id):
		return

	var node = _target_nodes[instance_id]
	# 先從字典移除，防止 _update_target_positions 繼續操作它
	_target_nodes.erase(instance_id)

	# 再播放特效（節點已從字典移除，不會被 update 干擾）
	if is_instance_valid(node):
		_play_kill_effect(node, data)

## 建立目標節點（使用像素 Sprite）
func _create_target_node(data: Dictionary) -> Node2D:
	var container = Node2D.new()
	container.position = Vector2(data.get("x", 0), data.get("y", 0))
	container.name = "Target_" + data.get("instance_id", "")

	var def_id = data.get("def_id", "T001")
	var target_type = data.get("type", "basic")

	# 使用像素 Sprite
	var sprite = Sprite2D.new()
	var sprite_path = TARGET_SPRITES.get(def_id, "")
	if sprite_path != "" and ResourceLoader.exists(sprite_path):
		sprite.texture = load(sprite_path)
		# 像素完美縮放（關閉濾波）
		sprite.texture_filter = CanvasItem.TEXTURE_FILTER_NEAREST
	else:
		# 備用：ColorRect
		var rect = ColorRect.new()
		rect.size = Vector2(32, 32) if target_type != "boss" else Vector2(64, 64)
		rect.position = -rect.size / 2
		rect.color = Color(0.8, 0.2, 0.8)
		container.add_child(rect)

	# BOSS 放大（B001 是 96x96，不需要額外放大）
	if target_type == "boss":
		sprite.scale = Vector2(1.5, 1.5)  # 輕微放大讓 BOSS 更有存在感

	container.add_child(sprite)

	# HP 條（規格書 8.1：受擊反饋）
	var hp_bar_bg = ColorRect.new()
	hp_bar_bg.size = Vector2(48, 5)
	hp_bar_bg.position = Vector2(-24, -38)
	hp_bar_bg.color = Color(0.2, 0.2, 0.2, 0.8)
	hp_bar_bg.name = "HPBarBG"
	container.add_child(hp_bar_bg)

	var hp_bar = ColorRect.new()
	hp_bar.size = Vector2(48, 5)
	hp_bar.position = Vector2(-24, -38)
	hp_bar.color = Color(0.2, 0.9, 0.2)
	hp_bar.name = "HPBar"
	container.add_child(hp_bar)

	# 儲存資料
	container.set_meta("instance_id", data.get("instance_id", ""))
	container.set_meta("def_id", def_id)
	container.set_meta("speed", data.get("speed", 0.0))
	container.set_meta("behavior", data.get("behavior", "linear"))
	container.set_meta("spawn_time", Time.get_ticks_msec())
	container.set_meta("target_type", target_type)

	return container

## 更新目標位置（依行為模式移動）
func _update_target_positions(delta: float) -> void:
	# 先收集要移除的 ID，避免迭代中修改 Dictionary
	var to_remove: Array[String] = []

	for instance_id in _target_nodes:
		var node = _target_nodes[instance_id]
		if not is_instance_valid(node):
			to_remove.append(instance_id)
			continue

		var speed = node.get_meta("speed", 0.0)
		var behavior = node.get_meta("behavior", "linear")

		if speed > 0:
			var t = Time.get_ticks_msec() * 0.001
			match behavior:
				"linear":
					node.position.x -= speed * delta
				"curve":
					node.position.x -= speed * delta
					node.position.y += sin(t * 2.0) * 30 * delta
				"jump":
					node.position.x -= speed * delta * 0.5
					node.position.y += cos(t * 4.0) * 50 * delta
				"meteor":
					node.position.x -= speed * delta
					node.position.y += speed * 0.3 * delta
				"sway":
					node.position.x -= speed * 0.3 * delta
					node.position.x += sin(t) * 20 * delta
				"static_sway":
					node.position.x += sin(t * 2.0) * 5 * delta
				"sink":
					node.position.y += speed * 0.4 * delta
				"flee":
					var flee_spd = node.get_meta("flee_speed", speed)
					node.position.x -= flee_spd * delta
				"coin_rain":
					node.position.x -= speed * delta
				"mimic":
					node.position.x -= speed * delta
					node.position.y += sin(t * 1.5) * 15 * delta
				"boss_phases":
					# BOSS 左右移動
					node.position.x += sin(t * 0.5) * speed * delta

		# 離開畫面標記移除
		if node.position.x < -150 or node.position.x > 1450:
			to_remove.append(instance_id)

	# 統一移除
	for id in to_remove:
		if _target_nodes.has(id):
			var node = _target_nodes[id]
			if is_instance_valid(node):
				node.queue_free()
			_target_nodes.erase(id)

## 擊破特效（規格書 8.2）
func _play_kill_effect(node: Node2D, data: Dictionary) -> void:
	if not is_instance_valid(node):
		return

	# 先記錄位置，再做動畫
	var kill_pos = node.position
	var reward = data.get("reward", 0)
	var multiplier = data.get("multiplier", 1.0)
	var def_id = data.get("def_id", "")

	# T101 擬態型怪物：死亡時先變形回原形（規格書 26.2）
	if def_id == "T101":
		_play_mimic_death(node, kill_pos, reward, multiplier)
		return

	# T105 巨大金幣魚：擊破後金幣雨（規格書 26.2）
	if def_id == "T105":
		_spawn_coin_rain(kill_pos)

	# 閃白所有子節點
	for child in node.get_children():
		if child is Sprite2D or child is ColorRect:
			child.modulate = Color.WHITE * 3.0

	# 縮放爆炸動畫
	var tween = create_tween()
	tween.tween_property(node, "scale", Vector2(1.5, 1.5), 0.08)
	tween.tween_property(node, "scale", Vector2(0.0, 0.0), 0.15)
	tween.tween_callback(func():
		if is_instance_valid(node):
			node.queue_free()
	)

	# 獎勵跳字（用記錄的位置，不依賴 node）
	if reward > 0:
		_spawn_reward_text(kill_pos, reward, multiplier)

	# 使用 HitEffect 系統（取代舊的 _spawn_death_particles）
	HitEffect.spawn_kill(kill_pos, multiplier)

	# 震動（依倍率）
	var trauma = clamp(0.2 + multiplier * 0.005, 0.2, 0.6)
	ScreenShake.add_trauma(trauma)

## T101 擬態型怪物死亡變形（規格書 26.2）
func _play_mimic_death(node: Node2D, kill_pos: Vector2, reward: int, multiplier: float) -> void:
	if not is_instance_valid(node):
		return

	# 第一階段：閃爍（偽裝破碎）
	var tween1 = create_tween()
	tween1.tween_property(node, "modulate", Color(2.0, 0.5, 0.5, 1.0), 0.08)
	tween1.tween_property(node, "modulate", Color.WHITE, 0.08)
	tween1.tween_property(node, "modulate", Color(2.0, 0.5, 0.5, 1.0), 0.08)
	tween1.tween_property(node, "modulate", Color.WHITE, 0.08)

	# 第二階段：縮小再放大（變形）
	var tween2 = create_tween()
	tween2.tween_property(node, "scale", Vector2(0.3, 1.5), 0.15)
	tween2.tween_property(node, "scale", Vector2(1.5, 0.3), 0.15)

	# 第三階段：爆炸消失
	var tween3 = create_tween()
	tween3.tween_interval(0.35)
	tween3.tween_property(node, "scale", Vector2(2.0, 2.0), 0.1)
	tween3.parallel().tween_property(node, "modulate:a", 0.0, 0.1)
	tween3.tween_callback(func():
		if is_instance_valid(node):
			node.queue_free()
	)

	# 生成「真面目」文字
	var reveal_label = Label.new()
	reveal_label.text = "正體！"
	reveal_label.position = kill_pos + Vector2(-20, -40)
	reveal_label.add_theme_font_size_override("font_size", 18)
	reveal_label.modulate = Color(1.0, 0.3, 0.3)
	add_child(reveal_label)
	var tween4 = create_tween()
	tween4.tween_property(reveal_label, "position:y", kill_pos.y - 90, 0.8)
	tween4.parallel().tween_property(reveal_label, "modulate:a", 0.0, 0.8)
	tween4.tween_callback(func():
		if is_instance_valid(reveal_label):
			reveal_label.queue_free()
	)

	if reward > 0:
		_spawn_reward_text(kill_pos, reward, multiplier)
	_spawn_death_particles(kill_pos)

## T105 金幣魚擊破後金幣雨（規格書 26.2）
func _spawn_coin_rain(origin: Vector2) -> void:
	# 生成 15 枚金幣從擊破位置散落
	for i in 15:
		var coin = ColorRect.new()
		coin.size = Vector2(10, 10)
		coin.color = Color(1.0, 0.85, 0.0)  # 金色
		coin.position = origin + Vector2(randf_range(-20, 20), randf_range(-10, 10))
		add_child(coin)

		# 拋物線軌跡
		var target_x = origin.x + randf_range(-120, 120)
		var target_y = origin.y + randf_range(80, 200)
		var peak_y = origin.y - randf_range(60, 120)

		var tween = create_tween()
		# 上升
		tween.tween_property(coin, "position", Vector2(
			origin.x + (target_x - origin.x) * 0.5,
			peak_y
		), 0.25)
		# 下落
		tween.tween_property(coin, "position", Vector2(target_x, target_y), 0.35)
		tween.parallel().tween_property(coin, "modulate:a", 0.0, 0.35)
		tween.tween_callback(func():
			if is_instance_valid(coin):
				coin.queue_free()
		)

	# 播放金幣音效
	AudioManager.play_sfx(AudioManager.SFX.COIN_DROP)

## 生成死亡粒子
func _spawn_death_particles(pos: Vector2) -> void:
	for i in 6:
		var particle = ColorRect.new()
		particle.size = Vector2(6, 6)
		particle.color = [Color.GOLD, Color.YELLOW, Color.WHITE, Color(1,0.5,0)][i % 4]
		particle.position = pos + Vector2(randf_range(-10, 10), randf_range(-10, 10))
		add_child(particle)
		var tween = create_tween()
		var target_pos = pos + Vector2(randf_range(-40, 40), randf_range(-60, -10))
		tween.tween_property(particle, "position", target_pos, 0.4)
		tween.parallel().tween_property(particle, "modulate:a", 0.0, 0.4)
		tween.tween_callback(func():
			if is_instance_valid(particle):
				particle.queue_free()
		)

## 生成獎勵跳字
func _spawn_reward_text(pos: Vector2, amount: int, multiplier: float) -> void:
	var label = Label.new()
	label.text = "+%d" % amount
	label.position = pos
	label.add_theme_font_size_override("font_size", 16)

	# 套用像素字體
	var font_path = "res://assets/fonts/pixel8.fnt"
	if ResourceLoader.exists(font_path):
		label.add_theme_font_override("font", load(font_path))

	# 依倍率設定顏色（規格書 8.3）
	if multiplier >= 100:
		label.modulate = Color(1.0, 0.2, 0.2)
	elif multiplier >= 20:
		label.modulate = Color(1.0, 0.8, 0.0)
	else:
		label.modulate = Color.WHITE

	add_child(label)

	var tween = create_tween()
	tween.tween_property(label, "position:y", pos.y - 80, 1.0)
	tween.parallel().tween_property(label, "modulate:a", 0.0, 1.0)
	tween.tween_callback(func():
		if is_instance_valid(label):
			label.queue_free()
	)

## 點擊目標（玩家點擊畫面時呼叫）
func try_click_target(click_pos: Vector2) -> String:
	var closest_id = ""
	var closest_dist = 70.0  # 點擊判定範圍（像素）

	for instance_id in _target_nodes:
		var node = _target_nodes[instance_id]
		if not is_instance_valid(node):
			continue
		var dist = node.position.distance_to(click_pos)
		if dist < closest_dist:
			closest_dist = dist
			closest_id = instance_id

	return closest_id

## 顯示 Lock 視覺框
func show_lock_indicator(instance_id: String) -> void:
	# 先清除舊的 lock 框
	for id in _target_nodes:
		var node = _target_nodes[id]
		if is_instance_valid(node):
			var old_lock = node.get_node_or_null("LockFrame")
			if old_lock:
				old_lock.queue_free()

	if instance_id == "":
		return

	if not _target_nodes.has(instance_id):
		return

	var node = _target_nodes[instance_id]
	if not is_instance_valid(node):
		return

	# 建立像素準星框
	var frame = Node2D.new()
	frame.name = "LockFrame"
	node.add_child(frame)

	var draw = func():
		pass  # 用 ColorRect 模擬準星

	# 四個角落的 L 形準星
	var size = 20
	var corners = [
		Vector2(-size, -size), Vector2(size, -size),
		Vector2(-size, size),  Vector2(size, size)
	]
	for corner in corners:
		var h = ColorRect.new()
		h.size = Vector2(8, 2)
		h.color = Color(1.0, 0.9, 0.0, 0.9)
		h.position = corner + Vector2(-4 if corner.x < 0 else -4, -1)
		frame.add_child(h)

		var v = ColorRect.new()
		v.size = Vector2(2, 8)
		v.color = Color(1.0, 0.9, 0.0, 0.9)
		v.position = corner + Vector2(-1, -4 if corner.y < 0 else -4)
		frame.add_child(v)

	# 閃爍動畫（綁定到 frame 節點，節點刪除時自動停止）
	var tween = frame.create_tween().set_loops()
	tween.tween_property(frame, "modulate:a", 0.4, 0.4)
	tween.tween_property(frame, "modulate:a", 1.0, 0.4)

## 更新目標 HP 條
func update_target_hp(instance_id: String, hp: int, max_hp: int) -> void:
	if not _target_nodes.has(instance_id):
		return
	var node = _target_nodes[instance_id]
	var hp_bar = node.get_node_or_null("HPBar")
	if hp_bar and max_hp > 0:
		hp_bar.size.x = 48.0 * (float(hp) / float(max_hp))
