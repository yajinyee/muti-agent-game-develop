## BackgroundManager.gd — 背景管理（海底/BOSS/Bonus 三種場景）
## environment-agent 負責維護
extends Node2D

# ── 背景狀態 ──────────────────────────────────────────────────
enum BgState { NORMAL, BOSS, BONUS }

var _current_state: BgState = BgState.NORMAL
var _bg_rect: ColorRect = null
var _bubble_layer: Node2D = null
var _overlay: ColorRect = null  # 狀態切換用的閃光覆蓋

# 背景顏色
const BG_COLORS = {
	BgState.NORMAL: Color(0.05, 0.10, 0.25),   # 深海藍
	BgState.BOSS:   Color(0.18, 0.03, 0.03),   # 深紅（BOSS 氛圍）
	BgState.BONUS:  Color(0.05, 0.18, 0.05),   # 深綠（草地）
}

# 氣泡顏色
const BUBBLE_COLORS = {
	BgState.NORMAL: Color(0.4, 0.7, 1.0, 0.25),
	BgState.BOSS:   Color(1.0, 0.3, 0.3, 0.2),
	BgState.BONUS:  Color(0.3, 1.0, 0.5, 0.2),
}

# ── 氣泡系統 ──────────────────────────────────────────────────
var _bubbles: Array = []
var _bubble_spawn_timer: float = 0.0
const BUBBLE_INTERVAL = 0.35
const MAX_BUBBLES = 20

func _ready() -> void:
	z_index = -10
	_build_background()
	_build_bubble_layer()
	_build_overlay()
	GameManager.game_state_changed.connect(_on_state_changed)

func _process(delta: float) -> void:
	_bubble_spawn_timer += delta
	if _bubble_spawn_timer >= BUBBLE_INTERVAL and _bubbles.size() < MAX_BUBBLES:
		_bubble_spawn_timer = 0.0
		_spawn_bubble()
	_update_bubbles(delta)

# ── 建立背景 ──────────────────────────────────────────────────

func _build_background() -> void:
	_bg_rect = ColorRect.new()
	_bg_rect.size = Vector2(1280, 720)
	_bg_rect.color = BG_COLORS[BgState.NORMAL]
	_bg_rect.z_index = -10
	add_child(_bg_rect)

	# 海底漸層（底部較深）
	var gradient = ColorRect.new()
	gradient.size = Vector2(1280, 200)
	gradient.position = Vector2(0, 520)
	gradient.color = Color(0.02, 0.05, 0.15, 0.6)
	gradient.z_index = -9
	add_child(gradient)

	# 海底沙地紋路（幾條橫線）
	for i in 4:
		var line = ColorRect.new()
		line.size = Vector2(1280, 2)
		line.position = Vector2(0, 560 + i * 30)
		line.color = Color(0.1, 0.15, 0.3, 0.4)
		line.z_index = -9
		add_child(line)

func _build_bubble_layer() -> void:
	_bubble_layer = Node2D.new()
	_bubble_layer.z_index = -5
	add_child(_bubble_layer)

func _build_overlay() -> void:
	_overlay = ColorRect.new()
	_overlay.size = Vector2(1280, 720)
	_overlay.color = Color(1.0, 1.0, 1.0, 0.0)
	_overlay.z_index = 5
	add_child(_overlay)

# ── 狀態切換 ──────────────────────────────────────────────────

func _on_state_changed(new_state: String) -> void:
	match new_state:
		"boss_warning", "boss_battle":
			_transition_to(BgState.BOSS)
		"bonus_game":
			_transition_to(BgState.BONUS)
		"normal_play", "boss_result", "bonus_result":
			_transition_to(BgState.NORMAL)

func _transition_to(new_state: BgState) -> void:
	if _current_state == new_state:
		return
	_current_state = new_state

	# 閃光過場
	var flash_color = Color(1.0, 1.0, 1.0, 0.5)
	match new_state:
		BgState.BOSS:   flash_color = Color(1.0, 0.2, 0.2, 0.6)
		BgState.BONUS:  flash_color = Color(0.2, 1.0, 0.3, 0.5)
		BgState.NORMAL: flash_color = Color(0.3, 0.6, 1.0, 0.4)

	_overlay.color = flash_color
	var tween = create_tween()
	tween.tween_property(_overlay, "color:a", 0.0, 0.5)

	# 背景顏色漸變
	var target_color = BG_COLORS[new_state]
	var bg_tween = _bg_rect.create_tween()
	bg_tween.tween_property(_bg_rect, "color", target_color, 0.8)

	# BOSS 狀態：螢幕震動
	if new_state == BgState.BOSS:
		ScreenShake.add_trauma(0.5)

# ── 氣泡系統 ──────────────────────────────────────────────────

func _spawn_bubble() -> void:
	var bubble = ColorRect.new()
	var size = randf_range(4, 14)
	bubble.size = Vector2(size, size)
	bubble.position = Vector2(randf_range(0, 1280), 720 + size)
	bubble.color = BUBBLE_COLORS[_current_state]
	bubble.z_index = -4
	_bubble_layer.add_child(bubble)

	var speed = randf_range(40, 120)
	var drift = randf_range(-20, 20)
	var lifetime = (720 + size) / speed

	bubble.set_meta("speed", speed)
	bubble.set_meta("drift", drift)
	bubble.set_meta("lifetime", lifetime)
	bubble.set_meta("age", 0.0)
	_bubbles.append(bubble)

	# 脈動動畫
	var pulse = bubble.create_tween().set_loops()
	pulse.tween_property(bubble, "modulate:a", 0.1, randf_range(0.5, 1.2))
	pulse.tween_property(bubble, "modulate:a", 1.0, randf_range(0.5, 1.2))

func _update_bubbles(delta: float) -> void:
	var to_remove: Array = []
	for bubble in _bubbles:
		if not is_instance_valid(bubble):
			to_remove.append(bubble)
			continue
		var age = bubble.get_meta("age", 0.0) + delta
		bubble.set_meta("age", age)
		var lifetime = bubble.get_meta("lifetime", 10.0)
		if age >= lifetime:
			to_remove.append(bubble)
			bubble.queue_free()
			continue
		var speed = bubble.get_meta("speed", 60.0)
		var drift = bubble.get_meta("drift", 0.0)
		bubble.position.y -= speed * delta
		bubble.position.x += drift * delta * 0.1
		# 顏色跟隨狀態
		bubble.color = BUBBLE_COLORS[_current_state]

	for b in to_remove:
		_bubbles.erase(b)
