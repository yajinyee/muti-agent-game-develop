## BonusGame.gd
## 瘋狂拔草 Bonus Game 場景控制（規格書 9章）

extends Node2D

signal bonus_ended(score: int, multiplier: float, reward: int)

# 場景節點
@onready var timer_label: Label = $UI/TimerLabel
@onready var score_label: Label = $UI/ScoreLabel
@onready var weed_container: Node2D = $WeedContainer
@onready var bg_sprite: Sprite2D = $Background
@onready var result_panel: Control = $UI/ResultPanel
@onready var result_label: Label = $UI/ResultPanel/ResultLabel

# 狀態
var _is_active: bool = false
var _score: int = 0
var _time_left: float = 15.0
var _weed_nodes: Dictionary = {}  # instance_id -> Node2D
var _weed_hp: Dictionary = {}     # instance_id -> remaining_clicks（BG002 需要 2 次）

# 雜草顏色（依類型）
const WEED_COLORS = {
	"BG001": Color(0.3, 0.7, 0.2),   # 普通雜草
	"BG002": Color(0.2, 0.5, 0.1),   # 硬雜草（深綠）
	"BG003": Color(0.5, 1.0, 0.3),   # 發光雜草（亮綠）
	"BG004": Color(1.0, 0.85, 0.0),  # 金色雜草
	"BG005": Color(0.8, 0.2, 0.2),   # 搗亂怪草（紅色）
}

func _ready() -> void:
	visible = false
	result_panel.visible = false

	# 連接 GameManager 訊號
	GameManager.bonus_event.connect(_on_bonus_event)
	GameManager.target_spawned.connect(_on_target_spawned)
	GameManager.target_killed.connect(_on_target_killed)

func _process(delta: float) -> void:
	if not _is_active:
		return

	_time_left -= delta
	if _time_left < 0:
		_time_left = 0

	timer_label.text = "%.1f" % _time_left

	# 倒數最後 5 秒閃爍
	if _time_left <= 5.0:
		var flash = int(_time_left * 4) % 2 == 0
		timer_label.modulate = Color.RED if flash else Color.WHITE

func _on_bonus_event(event_data: Dictionary) -> void:
	var event = event_data.get("event", "")
	match event:
		"ready":
			_show_ready()
		"start":
			_start_bonus()
		"tick":
			_time_left = event_data.get("time_left", 0.0)
		"end":
			_end_bonus(event_data)

func _show_ready() -> void:
	visible = true
	# 顯示 Bonus Ready 提示
	var ready_label = Label.new()
	ready_label.text = "🌿 BONUS READY! 🌿"
	ready_label.position = Vector2(440, 300)
	ready_label.add_theme_font_size_override("font_size", 48)
	ready_label.modulate = Color(0.3, 1.0, 0.3)
	add_child(ready_label)

	var tween = create_tween()
	tween.tween_property(ready_label, "scale", Vector2(1.3, 1.3), 0.3)
	tween.tween_property(ready_label, "scale", Vector2(1.0, 1.0), 0.3)
	tween.tween_property(ready_label, "modulate:a", 0.0, 1.0)
	tween.tween_callback(ready_label.queue_free)

func _start_bonus() -> void:
	_is_active = true
	_score = 0
	_time_left = 15.0
	_weed_hp.clear()
	visible = true
	result_panel.visible = false
	score_label.text = "分數: 0"
	timer_label.modulate = Color.WHITE

	# 切換背景
	var bonus_bg_path = "res://assets/sprites/backgrounds/bonus_bg.png"
	if ResourceLoader.exists(bonus_bg_path):
		bg_sprite.texture = load(bonus_bg_path)
		bg_sprite.texture_filter = CanvasItem.TEXTURE_FILTER_NEAREST

func _end_bonus(data: Dictionary) -> void:
	_is_active = false
	var score = data.get("score", 0)
	var multiplier = data.get("multiplier", 50.0)
	var reward = data.get("reward", 0)

	# 顯示結算
	result_panel.visible = true
	result_label.text = "今日勞動報酬！\n分數: %d\n倍率: %.0fx\n獎勵: %d" % [score, multiplier, reward]

	# 慶祝動畫
	_play_celebration()

	# 3秒後隱藏
	await get_tree().create_timer(3.0).timeout
	visible = false
	result_panel.visible = false
	emit_signal("bonus_ended", score, multiplier, reward)

func _play_celebration() -> void:
	# 生成慶祝粒子
	for i in 20:
		var particle = ColorRect.new()
		particle.size = Vector2(8, 8)
		particle.color = [Color.YELLOW, Color.GREEN, Color.GOLD, Color.WHITE][i % 4]
		particle.position = Vector2(randf_range(100, 1180), randf_range(100, 620))
		add_child(particle)

		var tween = create_tween()
		tween.tween_property(particle, "position:y", particle.position.y - randf_range(50, 150), 1.0)
		tween.parallel().tween_property(particle, "modulate:a", 0.0, 1.0)
		tween.tween_callback(particle.queue_free)

## 目標生成（Bonus 雜草）
func _on_target_spawned(data: Dictionary) -> void:
	if data.get("type", "") != "bonus":
		return
	if not _is_active:
		return

	var instance_id = data.get("instance_id", "")
	var def_id = data.get("def_id", "BG001")

	var node = Node2D.new()
	node.position = Vector2(data.get("x", 0), data.get("y", 0))

	# 雜草 Sprite（用 ColorRect 代替）
	var rect = ColorRect.new()
	rect.size = Vector2(24, 32)
	rect.position = Vector2(-12, -32)
	rect.color = WEED_COLORS.get(def_id, Color.GREEN)
	node.add_child(rect)

	# 雜草形狀（三角形葉子）
	var leaf = ColorRect.new()
	leaf.size = Vector2(16, 20)
	leaf.position = Vector2(-8, -48)
	leaf.color = WEED_COLORS.get(def_id, Color.GREEN) * 1.2
	node.add_child(leaf)

	# 金色雜草特效
	if def_id == "BG004":
		rect.color = Color(1.0, 0.85, 0.0)
		leaf.color = Color(1.0, 0.95, 0.3)
		# 閃光效果
		var tween = create_tween().set_loops()
		tween.tween_property(node, "modulate", Color(1.5, 1.5, 0.5), 0.3)
		tween.tween_property(node, "modulate", Color.WHITE, 0.3)

	# 發光雜草特效（BG003）— 綠色光暈 + 倍率提升視覺
	if def_id == "BG003":
		rect.color = Color(0.4, 1.0, 0.2)
		leaf.color = Color(0.6, 1.0, 0.3)
		# 持續發光閃爍
		var tween_glow = create_tween().set_loops()
		tween_glow.tween_property(node, "modulate", Color(0.5, 2.0, 0.5, 1.0), 0.25)
		tween_glow.tween_property(node, "modulate", Color.WHITE, 0.25)
		# 顯示「×UP!」提示
		var up_label = Label.new()
		up_label.text = "✨×UP!"
		up_label.position = Vector2(-16, -70)
		up_label.add_theme_font_size_override("font_size", 14)
		up_label.modulate = Color(0.3, 1.0, 0.3)
		node.add_child(up_label)
		# 上下浮動動畫
		var tween_float = create_tween().set_loops()
		tween_float.tween_property(up_label, "position:y", -78.0, 0.5)
		tween_float.tween_property(up_label, "position:y", -70.0, 0.5)

	# 搗亂怪草（紅色警告）
	if def_id == "BG005":
		var warning = Label.new()
		warning.text = "⚠"
		warning.position = Vector2(-8, -60)
		node.add_child(warning)

	node.set_meta("instance_id", instance_id)
	node.set_meta("def_id", def_id)
	weed_container.add_child(node)
	_weed_nodes[instance_id] = node

	# BG002 硬雜草需要連點 2 次（規格書 29.3）
	if def_id == "BG002":
		_weed_hp[instance_id] = 2
	else:
		_weed_hp[instance_id] = 1

## 目標擊破（拔草）
func _on_target_killed(data: Dictionary) -> void:
	var instance_id = data.get("instance_id", "")
	if not _weed_nodes.has(instance_id):
		return

	var node = _weed_nodes[instance_id]
	_weed_nodes.erase(instance_id)

	# 拔草動畫
	var tween = create_tween()
	tween.tween_property(node, "position:y", node.position.y - 60, 0.3)
	tween.parallel().tween_property(node, "scale", Vector2(0, 0), 0.3)
	tween.tween_callback(node.queue_free)

	# 更新分數
	var labor_gain = data.get("labor_gain", 1)
	_score += labor_gain
	score_label.text = "分數: %d" % _score

## 處理點擊（拔草）
func _input(event: InputEvent) -> void:
	if not _is_active:
		return
	if not (event is InputEventMouseButton and event.pressed):
		return

	var click_pos = event.position

	# 找最近的雜草
	var closest_id = ""
	var closest_dist = 50.0

	for instance_id in _weed_nodes:
		var node = _weed_nodes[instance_id]
		var dist = node.position.distance_to(click_pos)
		if dist < closest_dist:
			closest_dist = dist
			closest_id = instance_id

	if closest_id != "":
		var def_id = ""
		if _weed_nodes.has(closest_id):
			def_id = _weed_nodes[closest_id].get_meta("def_id", "")

		# BG002 硬雜草：需要連點 2 次（規格書 29.3）
		if def_id == "BG002" and _weed_hp.has(closest_id):
			_weed_hp[closest_id] -= 1
			# 第一次點擊：搖晃動畫，不送 server
			if _weed_hp[closest_id] > 0:
				var node = _weed_nodes[closest_id]
				if is_instance_valid(node):
					var tween = create_tween()
					tween.tween_property(node, "rotation_degrees", 15, 0.05)
					tween.tween_property(node, "rotation_degrees", -15, 0.05)
					tween.tween_property(node, "rotation_degrees", 0, 0.05)
				AudioManager.play_sfx(AudioManager.SFX.HIT)
				return  # 不送 server，等第二次點擊

		NetworkManager.send_bonus_click(closest_id, click_pos.x, click_pos.y)
		AudioManager.play_sfx(AudioManager.SFX.WEED_PULL)

		# BG004 金色雜草：觸發金幣雨視覺（規格書 29.3）
		if def_id == "BG004":
			_spawn_coin_shower(click_pos)

		# BG005 搗亂怪草：暫停操作 0.3 秒（規格書 29.3）
		if def_id == "BG005":
			_is_active = false
			var stun_label = Label.new()
			stun_label.text = "😵 STUNNED!"
			stun_label.position = Vector2(500, 300)
			stun_label.add_theme_font_size_override("font_size", 36)
			stun_label.modulate = Color(1.0, 0.3, 0.3)
			add_child(stun_label)
			await get_tree().create_timer(0.3).timeout
			if is_instance_valid(stun_label):
				stun_label.queue_free()
			_is_active = true

## BG004 金色雜草金幣雨（規格書 29.3）
func _spawn_coin_shower(origin: Vector2) -> void:
	for i in 20:
		var coin = ColorRect.new()
		coin.size = Vector2(8, 8)
		coin.color = Color(1.0, 0.85, 0.0)
		coin.position = origin + Vector2(randf_range(-15, 15), 0)
		add_child(coin)

		var tx = origin.x + randf_range(-100, 100)
		var ty = origin.y + randf_range(60, 180)
		var tween = create_tween()
		tween.tween_property(coin, "position", Vector2(tx, ty), 0.5)
		tween.parallel().tween_property(coin, "modulate:a", 0.0, 0.5)
		tween.tween_callback(func():
			if is_instance_valid(coin):
				coin.queue_free()
		)

	AudioManager.play_sfx(AudioManager.SFX.COIN_DROP)
