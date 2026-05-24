## TargetManager.gd — 目標物系統
## target-system-agent 負責維護
extends Node2D

# 目標物 Sprite 路徑
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
	# DAY-292 新增特殊目標
	"T106": "res://assets/sprites/targets/T106_chain_lightning.png",
	"T107": "res://assets/sprites/targets/T107_crab_torpedo.png",
	"T108": "res://assets/sprites/targets/T108_vortex_anemone.png",
	"T109": "res://assets/sprites/targets/T109_golden_dragon.png",
	"T110": "res://assets/sprites/targets/T110_thunder_lobster.png",
}

# 目標物顏色（無 Sprite 時的備用顏色）
const TARGET_COLORS = {
	"T001": Color(0.2, 0.8, 0.2),   # 綠色雜草
	"T002": Color(0.3, 0.9, 0.3),   # 綠色小蟲
	"T003": Color(0.9, 0.3, 0.3),   # 紅色小蟲
	"T004": Color(0.3, 0.5, 0.9),   # 藍色小蟲
	"T005": Color(1.0, 0.9, 0.4),   # 黃色布丁
	"T006": Color(0.6, 0.4, 0.2),   # 棕色蘑菇
	"T101": Color(0.6, 0.6, 0.6),   # 灰色擬態
	"T102": Color(0.9, 0.7, 0.2),   # 金色寶箱
	"T103": Color(1.0, 1.0, 0.8),   # 白色流星
	"T104": Color(1.0, 0.85, 0.0),  # 金色雜草
	"T105": Color(1.0, 0.8, 0.2),   # 金色魚
	"B001": Color(0.8, 0.2, 0.2),   # 紅色 BOSS
	# DAY-292 新增特殊目標備用顏色
	"T106": Color(0.0, 0.9, 1.0),   # 青藍連鎖閃電
	"T107": Color(1.0, 0.4, 0.1),   # 橙紅螃蟹魚雷
	"T108": Color(0.5, 0.2, 0.8),   # 紫色渦旋海葵
	"T109": Color(1.0, 0.85, 0.0),  # 金色黃金龍魚
	"T110": Color(1.0, 0.3, 0.0),   # 火紅雷霆龍蝦
}

var _target_nodes: Dictionary = {}  # instance_id -> Node2D
var _cached_textures: Dictionary = {}

func _ready() -> void:
	GameManager.target_spawned.connect(_on_target_spawned)
	GameManager.target_updated.connect(_on_target_updated)
	GameManager.target_killed.connect(_on_target_killed)
	GameManager.boss_event.connect(_on_boss_event)

func _process(delta: float) -> void:
	_update_positions(delta)

func _update_positions(delta: float) -> void:
	for instance_id in _target_nodes:
		var node = _target_nodes[instance_id]
		if not is_instance_valid(node):
			continue
		var behavior = node.get_meta("behavior", "linear")
		var speed = node.get_meta("speed", 0.0)
		if node.get_meta("is_fleeing", false):
			speed = node.get_meta("flee_speed", speed * 2.5)

		match behavior:
			"linear", "flee", "fast":
				node.position.x -= speed * delta
			"sink":
				node.position.y += speed * 0.3 * delta
				node.position.x -= 10 * delta

		# 離開畫面則移除
		if node.position.x < -100:
			_target_nodes.erase(instance_id)
			node.queue_free()

# ── 目標物生成 ────────────────────────────────────────────────

func _on_target_spawned(data: Dictionary) -> void:
	var instance_id = data.get("instance_id", "")
	if _target_nodes.has(instance_id):
		return

	var node = _create_target_node(data)
	_target_nodes[instance_id] = node

	# 進場動畫
	node.scale = Vector2.ZERO
	var tween = node.create_tween()
	var target_type = data.get("type", "basic")
	if target_type == "boss":
		tween.tween_property(node, "scale", Vector2(1.0, 1.0), 0.4).set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_ELASTIC)
	else:
		tween.tween_property(node, "scale", Vector2(1.0, 1.0), 0.15).set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_BACK)

func _create_target_node(data: Dictionary) -> Node2D:
	var container = Node2D.new()
	container.position = Vector2(data.get("x", 0), data.get("y", 0))
	container.name = "Target_" + data.get("instance_id", "")
	add_child(container)

	var def_id = data.get("def_id", "T001")
	var target_type = data.get("type", "basic")
	var multiplier = data.get("multiplier", 2.0)

	# Sprite
	var sprite = Sprite2D.new()
	var tex = _get_texture(def_id)
	if tex != null:
		sprite.texture = tex
		sprite.texture_filter = CanvasItem.TEXTURE_FILTER_NEAREST
		# 目標物放大 2x，確保在畫面上清楚可見
		sprite.scale = Vector2(2.0, 2.0)
		if target_type == "boss":
			sprite.scale = Vector2(3.0, 3.0)
	else:
		# 備用：ColorRect
		var rect = ColorRect.new()
		var size = 64.0 if target_type != "boss" else 128.0
		rect.size = Vector2(size, size)
		rect.position = -rect.size / 2
		rect.color = TARGET_COLORS.get(def_id, Color(0.8, 0.2, 0.8))
		container.add_child(rect)
	container.add_child(sprite)

	# HP 條
	var hp_max = data.get("max_hp", 1)
	var hp_bg = ColorRect.new()
	hp_bg.size = Vector2(64, 6)
	hp_bg.position = Vector2(-32, -52)
	hp_bg.color = Color(0.2, 0.2, 0.2, 0.8)
	hp_bg.name = "HPBarBG"
	container.add_child(hp_bg)

	var hp_bar = ColorRect.new()
	hp_bar.size = Vector2(64, 6)
	hp_bar.position = Vector2(-32, -52)
	hp_bar.color = Color(0.2, 0.9, 0.2)
	hp_bar.name = "HPBar"
	container.add_child(hp_bar)

	# 倍率標籤
	if target_type != "boss":
		var mult_label = Label.new()
		mult_label.text = "x%.0f" % multiplier
		mult_label.position = Vector2(-20, 36)
		mult_label.add_theme_font_size_override("font_size", 14)
		mult_label.modulate = _mult_label_color(multiplier)
		mult_label.name = "MultLabel"
		container.add_child(mult_label)

	# 高倍率光暈
	if multiplier >= 30.0:
		_add_glow(container, multiplier)

	# 特殊搖晃（T103 流星、T104 金草）
	if def_id in ["T103", "T104"]:
		var wobble = container.create_tween().set_loops()
		var deg = 5.0 if def_id == "T103" else 3.0
		var dur = 0.15 if def_id == "T103" else 0.4
		wobble.tween_property(container, "rotation_degrees", deg, dur)
		wobble.tween_property(container, "rotation_degrees", -deg, dur)

	# 儲存 meta
	container.set_meta("instance_id", data.get("instance_id", ""))
	container.set_meta("def_id", def_id)
	container.set_meta("speed", data.get("speed", 0.0))
	container.set_meta("behavior", data.get("behavior", "linear"))
	container.set_meta("target_type", target_type)
	container.set_meta("multiplier", multiplier)
	container.set_meta("is_fleeing", false)

	return container

func _add_glow(node: Node2D, multiplier: float) -> void:
	var glow = ColorRect.new()
	glow.size = Vector2(80, 80)
	glow.position = Vector2(-40, -40)
	glow.z_index = -1
	if multiplier >= 50:
		glow.color = Color(1.0, 0.4, 0.1, 0.3)
	else:
		glow.color = Color(1.0, 0.85, 0.0, 0.25)
	node.add_child(glow)
	var tween = glow.create_tween().set_loops()
	tween.tween_property(glow, "modulate:a", 0.1, 0.4)
	tween.tween_property(glow, "modulate:a", 1.0, 0.4)

# ── 目標物更新 ────────────────────────────────────────────────

func _on_target_updated(data: Dictionary) -> void:
	var instance_id = data.get("instance_id", "")
	if not _target_nodes.has(instance_id):
		return
	var node = _target_nodes[instance_id]
	if not is_instance_valid(node):
		return

	# 更新 HP 條
	update_target_hp(instance_id, data.get("hp", 0), data.get("max_hp", 1))

	# 受擊閃白
	_flash_hit(node)

	# T102 逃跑
	if data.get("is_fleeing", false):
		node.set_meta("is_fleeing", true)
		node.set_meta("flee_speed", node.get_meta("speed", 70.0) * 2.5)
		var tween = node.create_tween()
		tween.tween_property(node, "modulate", Color(2.0, 0.5, 0.5), 0.06)
		tween.tween_property(node, "modulate", Color.WHITE, 0.06)
		tween.tween_property(node, "modulate", Color(2.0, 0.5, 0.5), 0.06)
		tween.tween_property(node, "modulate", Color.WHITE, 0.06)

func update_target_hp(instance_id: String, hp: int, max_hp: int) -> void:
	if not _target_nodes.has(instance_id):
		return
	var node = _target_nodes[instance_id]
	if not is_instance_valid(node):
		return
	var hp_bar = node.get_node_or_null("HPBar")
	var hp_bg = node.get_node_or_null("HPBarBG")
	if not is_instance_valid(hp_bar) or not is_instance_valid(hp_bg):
		return
	var pct = float(hp) / float(max_hp) if max_hp > 0 else 0.0
	hp_bar.size.x = hp_bg.size.x * pct
	if pct > 0.6:
		hp_bar.color = Color(0.2, 0.9, 0.2)
	elif pct > 0.3:
		hp_bar.color = Color(1.0, 0.8, 0.1)
	else:
		hp_bar.color = Color(1.0, 0.2, 0.2)

func _flash_hit(node: Node2D) -> void:
	for child in node.get_children():
		if child is Sprite2D:
			var tween = node.create_tween()
			tween.tween_property(child, "modulate", Color(3.0, 3.0, 3.0), 0.04)
			tween.tween_property(child, "modulate", Color.WHITE, 0.08)
			break

# ── 目標物擊破 ────────────────────────────────────────────────

func _on_target_killed(data: Dictionary) -> void:
	var instance_id = data.get("instance_id", "")
	if not _target_nodes.has(instance_id):
		return
	var node = _target_nodes[instance_id]
	_target_nodes.erase(instance_id)

	if not is_instance_valid(node):
		return

	var reward = data.get("reward", 0)
	var multiplier = data.get("multiplier", 1.0)

	# 擊破特效
	HitEffect.spawn_kill(node.position, multiplier)
	if reward > 0:
		HitEffect.spawn_reward_text(node.position, reward, multiplier)

	# 大獎演出
	if multiplier >= 20:
		HitEffect.spawn_big_win(node.position, multiplier)
		AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
	else:
		AudioManager.play_sfx(AudioManager.SFX.KILL)

	# T105 金幣魚：金幣雨
	if data.get("def_id", "") == "T105":
		_spawn_coin_rain(node.position)

	# 消失動畫
	var tween = node.create_tween()
	tween.tween_property(node, "scale", Vector2(1.5, 1.5), 0.1)
	tween.parallel().tween_property(node, "modulate:a", 0.0, 0.15)
	tween.tween_callback(func(): if is_instance_valid(node): node.queue_free())

func _spawn_coin_rain(pos: Vector2) -> void:
	for i in 10:
		var coin = ColorRect.new()
		coin.size = Vector2(12, 12)
		coin.color = Color(1.0, 0.85, 0.0)
		coin.position = pos
		coin.z_index = 45
		get_parent().add_child(coin)
		var angle = randf_range(-PI/2 - 0.5, -PI/2 + 0.5)
		var dist = randf_range(40, 100)
		var target = pos + Vector2(cos(angle), sin(angle)) * dist
		var tween = coin.create_tween()
		tween.tween_property(coin, "position", target, 0.4)
		tween.tween_property(coin, "position:y", target.y + 60, 0.3)
		tween.parallel().tween_property(coin, "modulate:a", 0.0, 0.3)
		tween.tween_callback(func(): if is_instance_valid(coin): coin.queue_free())

# ── BOSS 事件 ─────────────────────────────────────────────────

func _on_boss_event(event_data: Dictionary) -> void:
	var event = event_data.get("event", "")
	if event == "phase_change":
		var instance_id = event_data.get("instance_id", "")
		if _target_nodes.has(instance_id):
			var node = _target_nodes[instance_id]
			if is_instance_valid(node):
				var tween = node.create_tween()
				for i in 3:
					tween.tween_property(node, "modulate", Color(3.0, 0.3, 0.3), 0.08)
					tween.tween_property(node, "modulate", Color.WHITE, 0.08)
				tween.tween_property(node, "modulate", Color(1.5, 0.5, 0.5), 0.1)
				tween.parallel().tween_property(node, "scale", Vector2(1.1, 1.1), 0.3)
				ScreenShake.add_trauma(0.6)

# ── 點擊目標 ──────────────────────────────────────────────────

func try_click_target(click_pos: Vector2) -> String:
	var best_id = ""
	var best_dist = 80.0  # 點擊半徑
	for instance_id in _target_nodes:
		var node = _target_nodes[instance_id]
		if not is_instance_valid(node):
			continue
		var dist = node.position.distance_to(click_pos)
		if dist < best_dist:
			best_dist = dist
			best_id = instance_id
	return best_id

# ── 輔助 ──────────────────────────────────────────────────────

func _get_texture(def_id: String) -> Texture2D:
	var path = TARGET_SPRITES.get(def_id, "")
	if path == "":
		return null
	if _cached_textures.has(path):
		return _cached_textures[path]
	if ResourceLoader.exists(path):
		var tex = load(path)
		_cached_textures[path] = tex
		return tex
	return null

func _mult_label_color(mult: float) -> Color:
	if mult >= 50: return Color(1.0, 0.4, 0.1)
	if mult >= 30: return Color(1.0, 0.85, 0.0)
	if mult >= 15: return Color(0.8, 0.9, 1.0)
	return Color(1.0, 1.0, 1.0)
