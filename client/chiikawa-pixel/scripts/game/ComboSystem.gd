## ComboSystem.gd — 連擊計數器與爽感系統
## combo-system-agent 負責維護
## DAY-322：新建，提升「每 30 秒有幾次讓玩家哇的時刻」指標
##
## 設計原則：
## - 連擊計數：每次命中 +1，超過 2 秒未命中重置
## - 連擊里程碑：5/10/20/50 連擊有特殊演出
## - 連擊倍率加成：連擊越高，獎勵跳字越大（視覺效果）
## - 不影響 Server 端 RTP，純視覺/音效效果
extends Node

signal combo_updated(count: int)
signal combo_milestone(count: int)
signal combo_reset(final_count: int)

# 連擊設定
const COMBO_TIMEOUT = 2.5          # 超過此秒數未命中則重置
const MILESTONE_COUNTS = [5, 10, 20, 50, 100]  # 里程碑連擊數

# 里程碑顏色
const MILESTONE_COLORS = {
	5:   Color(1.0, 1.0, 0.4),    # 黃色
	10:  Color(1.0, 0.7, 0.0),    # 橙色
	20:  Color(1.0, 0.3, 0.1),    # 火紅
	50:  Color(0.8, 0.0, 1.0),    # 深紫
	100: Color(1.0, 0.0, 0.5),    # 宇宙粉紅
}

# 里程碑文字
const MILESTONE_TEXTS = {
	5:   "COMBO x5!",
	10:  "🔥 COMBO x10!",
	20:  "💥 COMBO x20!!",
	50:  "⚡ COMBO x50!!!",
	100: "🌟 COMBO x100!!!!",
}

var _combo_count: int = 0
var _combo_timer: float = 0.0
var _is_active: bool = false
var _scene_root: Node = null

func _ready() -> void:
	GameManager.attack_result.connect(_on_attack_result)
	call_deferred("_find_scene_root")

func _find_scene_root() -> void:
	var tree = get_tree()
	if tree == null:
		return
	var root = tree.get_root()
	for child in root.get_children():
		if child is Node2D:
			_scene_root = child
			return
	_scene_root = root

func _process(delta: float) -> void:
	if not _is_active:
		return
	_combo_timer += delta
	if _combo_timer >= COMBO_TIMEOUT:
		_reset_combo()

func _on_attack_result(result: Dictionary) -> void:
	if not result.get("is_hit", false):
		return
	_register_hit()

func _register_hit() -> void:
	_combo_count += 1
	_combo_timer = 0.0
	_is_active = true
	combo_updated.emit(_combo_count)

	# 檢查里程碑
	if _combo_count in MILESTONE_COUNTS:
		_trigger_milestone(_combo_count)

func _reset_combo() -> void:
	if _combo_count > 0:
		combo_reset.emit(_combo_count)
	_combo_count = 0
	_combo_timer = 0.0
	_is_active = false
	combo_updated.emit(0)

func _trigger_milestone(count: int) -> void:
	combo_milestone.emit(count)
	_spawn_milestone_effect(count)

	# 螢幕震動（里程碑越高震動越強）
	var trauma = 0.3 + float(count) / 200.0
	ScreenShake.add_trauma(clamp(trauma, 0.3, 0.9))

	# 音效（用 BIG_WIN 音效）
	if count >= 20:
		AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)

func _spawn_milestone_effect(count: int) -> void:
	if not is_instance_valid(_scene_root):
		return

	var color = MILESTONE_COLORS.get(count, Color(1.0, 0.85, 0.0))
	var text = MILESTONE_TEXTS.get(count, "COMBO x%d!" % count)

	# 中央大字
	var label = Label.new()
	label.text = text
	var font_size = 28 + count / 5
	font_size = clamp(font_size, 28, 52)
	label.add_theme_font_size_override("font_size", font_size)
	label.modulate = color
	label.position = Vector2(640 - 100, 260)
	label.z_index = 65
	_scene_root.add_child(label)

	# 彈跳進場 + 停留 + 淡出
	var tween = label.create_tween()
	tween.tween_property(label, "scale", Vector2(1.5, 1.5), 0.08)
	tween.tween_property(label, "scale", Vector2(1.0, 1.0), 0.06)
	tween.tween_interval(0.6)
	tween.tween_property(label, "modulate:a", 0.0, 0.4)
	tween.tween_callback(func(): if is_instance_valid(label): label.queue_free())

	# 放射狀粒子爆炸
	var particle_count = 6 + count / 5
	particle_count = clamp(particle_count, 6, 20)
	for i in particle_count:
		var dot = ColorRect.new()
		var size = randf_range(6, 12)
		dot.size = Vector2(size, size)
		dot.color = color
		dot.position = Vector2(640, 300)
		dot.z_index = 60
		_scene_root.add_child(dot)
		var angle = (float(i) / float(particle_count)) * TAU
		var dist = randf_range(60, 120 + count * 2)
		var target = Vector2(640, 300) + Vector2(cos(angle), sin(angle)) * dist
		var ptween = dot.create_tween()
		ptween.tween_property(dot, "position", target, 0.4)
		ptween.parallel().tween_property(dot, "modulate:a", 0.0, 0.4)
		ptween.tween_callback(func(): if is_instance_valid(dot): dot.queue_free())

	# 50+ 連擊：全螢幕閃光
	if count >= 50:
		_spawn_screen_flash(color)

func _spawn_screen_flash(color: Color) -> void:
	if not is_instance_valid(_scene_root):
		return
	var flash = ColorRect.new()
	flash.size = Vector2(1280, 720)
	flash.position = Vector2.ZERO
	flash.color = Color(color.r, color.g, color.b, 0.3)
	flash.z_index = 70
	_scene_root.add_child(flash)
	var tween = flash.create_tween()
	tween.tween_property(flash, "modulate:a", 0.0, 0.2)
	tween.tween_callback(func(): if is_instance_valid(flash): flash.queue_free())

## 取得當前連擊數（供 HUD 顯示）
func get_combo_count() -> int:
	return _combo_count

## 取得連擊視覺強度（0.0-1.0，供 HUD 顯示動畫）
func get_combo_intensity() -> float:
	if _combo_count == 0:
		return 0.0
	return clamp(float(_combo_count) / 50.0, 0.0, 1.0)
