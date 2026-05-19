## UnderwaterOverlay.gd
## 全螢幕水下視覺效果（色差 + 水波扭曲 + 藍色調 + 深度霧氣）
## 在 normal 背景（海底）狀態時啟用，提升沉浸感
## 掛載在最高層 CanvasLayer 上的 ColorRect

extends ColorRect

const SHADER_PATH = "res://assets/shaders/underwater_overlay.gdshader"

var _mat: ShaderMaterial = null
var _active: bool = false
var _fade_tween: Tween = null

func _ready() -> void:
	# 全螢幕覆蓋
	set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	# 初始透明（不影響畫面）
	color = Color(0, 0, 0, 0)
	# 設定 z_index 讓它在所有遊戲元素上方，但在 HUD 下方
	z_index = 50
	mouse_filter = Control.MOUSE_FILTER_IGNORE  # 不攔截滑鼠事件

	# 載入 shader
	if ResourceLoader.exists(SHADER_PATH):
		_mat = ShaderMaterial.new()
		_mat.shader = load(SHADER_PATH)
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

	var target_alpha = 1.0 if active else 0.0

	if animated:
		_fade_tween = create_tween()
		_fade_tween.tween_method(
			func(v: float): _mat.set_shader_parameter("effect_alpha", v),
			_mat.get_shader_parameter("effect_alpha") if _mat.get_shader_parameter("effect_alpha") != null else 0.0,
			target_alpha,
			0.5
		)
	else:
		_mat.set_shader_parameter("effect_alpha", target_alpha)
