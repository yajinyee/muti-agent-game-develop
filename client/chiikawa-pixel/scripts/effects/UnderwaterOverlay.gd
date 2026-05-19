## UnderwaterOverlay.gd
## 全螢幕水下視覺效果（色差 + 水波扭曲 + 藍色調 + 深度霧氣）
## 在 normal 背景（海底）狀態時啟用，提升沉浸感
## 掛載在 UnderwaterLayer（CanvasLayer layer=49）下的 ColorRect
## 注意：CanvasLayer 在 HUD（layer=1）之下，確保 UI 不受影響

extends ColorRect

const SHADER_PATH = "res://assets/shaders/underwater_overlay.gdshader"

var _mat: ShaderMaterial = null
var _active: bool = false
var _fade_tween: Tween = null

func _ready() -> void:
	# 全螢幕覆蓋（相對於 CanvasLayer）
	set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	# 初始顏色（shader 會覆蓋，這裡設白色讓 shader 正常採樣螢幕）
	color = Color(1, 1, 1, 1)
	# 不攔截滑鼠事件
	mouse_filter = Control.MOUSE_FILTER_IGNORE

	# 載入 shader
	if ResourceLoader.exists(SHADER_PATH):
		_mat = ShaderMaterial.new()
		_mat.shader = load(SHADER_PATH)
		# 初始 effect_alpha = 0（無效果，等待狀態切換）
		_mat.set_shader_parameter("effect_alpha", 0.0)
		material = _mat
	else:
		push_warning("[UnderwaterOverlay] Shader not found: " + SHADER_PATH)

	# 監聽遊戲狀態變化
	GameManager.game_state_changed.connect(_on_state_changed)

	# 初始狀態：海底（啟用效果）
	_set_active(true, false)

## 依遊戲狀態切換效果
func _on_state_changed(state: String) -> void:
	match state:
		"boss_warning", "boss_battle", "boss_result":
			# BOSS 狀態：關閉水下效果（BOSS 場景不在水下）
			_set_active(false, true)
		"bonus_game", "bonus_ready", "bonus_result":
			# Bonus 狀態：關閉水下效果（草地場景）
			_set_active(false, true)
		"normal_play", "special_target_event":
			# 正常遊戲：啟用水下效果
			_set_active(true, true)

## 設定效果啟用狀態（帶淡入淡出）
func _set_active(active: bool, animated: bool) -> void:
	if _active == active:
		return
	_active = active

	if _mat == null:
		return

	# 停止舊的 tween
	if is_instance_valid(_fade_tween):
		_fade_tween.kill()

	var current_alpha = _mat.get_shader_parameter("effect_alpha")
	if current_alpha == null:
		current_alpha = 0.0
	var target_alpha = 1.0 if active else 0.0

	if animated:
		_fade_tween = create_tween()
		_fade_tween.tween_method(
			func(v: float): _mat.set_shader_parameter("effect_alpha", v),
			current_alpha,
			target_alpha,
			0.5
		)
	else:
		_mat.set_shader_parameter("effect_alpha", target_alpha)
