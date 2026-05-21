## ScreenShake.gd
## 畫面震動系統（Trauma-based，平滑震動）
## Autoload 單例（Node）
## 透過操作場景中的 Camera2D 節點實現震動
##
## 使用方式：
##   ScreenShake.add_trauma(0.18)  # 輕微（命中）
##   ScreenShake.add_trauma(0.35)  # 中等（擊殺）
##   ScreenShake.add_trauma(0.7)   # 強烈（大獎）
##   ScreenShake.add_trauma(0.9)   # 最強（BOSS 登場）

extends Node

# ── 參數 ──────────────────────────────────────────────
## trauma 衰減速度（每秒）
var decay: float = 2.5
## 最大位移（像素）
var max_offset: Vector2 = Vector2(7.0, 5.0)
## 最大旋轉（弧度）
var max_rotation: float = 0.025
## 像素遊戲：限制到整數像素
var pixel_perfect: bool = true

# ── 內部狀態 ──────────────────────────────────────────
var _trauma: float = 0.0
var _time: float = 0.0
var _camera: Camera2D = null

# ── 公開 API ──────────────────────────────────────────

func add_trauma(amount: float) -> void:
	# 低效能模式下關閉震動
	if not PerformanceMonitor.is_screen_shake_enabled():
		return
	_trauma = clamp(_trauma + amount, 0.0, 1.0)

func set_trauma(amount: float) -> void:
	_trauma = clamp(amount, 0.0, 1.0)

func stop() -> void:
	_trauma = 0.0
	_apply_offset(Vector2.ZERO, 0.0)

# ── 內部邏輯 ──────────────────────────────────────────

func _process(delta: float) -> void:
	# 延遲取得 Camera2D（場景載入後才存在）
	if not is_instance_valid(_camera):
		_find_camera()

	if _trauma <= 0.0:
		if is_instance_valid(_camera) and (_camera.offset != Vector2.ZERO or _camera.rotation != 0.0):
			_apply_offset(Vector2.ZERO, 0.0)
		return

	_time += delta * 9.0

	var shake = _trauma * _trauma  # 平方讓小 trauma 更柔和

	var ox = sin(_time * 1.7) * cos(_time * 2.3) * max_offset.x * shake
	var oy = cos(_time * 1.3) * sin(_time * 2.7) * max_offset.y * shake
	var rot = sin(_time * 2.1) * max_rotation * shake

	if pixel_perfect:
		_apply_offset(Vector2(round(ox), round(oy)), rot)
	else:
		_apply_offset(Vector2(ox, oy), rot)

	_trauma = max(0.0, _trauma - decay * delta)

func _apply_offset(off: Vector2, rot: float) -> void:
	if is_instance_valid(_camera):
		_camera.offset = off
		_camera.rotation = rot

func _find_camera() -> void:
	var scene = get_tree().current_scene
	if is_instance_valid(scene):
		_camera = scene.get_node_or_null("Camera2D")
