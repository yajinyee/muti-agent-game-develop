## ScreenShake.gd — Trauma-based 螢幕震動
## screen-effect-agent 負責維護
extends Node

var _trauma: float = 0.0
var _camera: Camera2D = null

const DECAY = 2.5
const MAX_OFFSET = 12.0
const MAX_ROTATION = 0.04

func _process(delta: float) -> void:
	if _trauma <= 0:
		return
	_trauma = maxf(0.0, _trauma - DECAY * delta)
	var shake = _trauma * _trauma  # trauma² 讓小震動更柔和
	if is_instance_valid(_camera):
		_camera.offset = Vector2(
			randf_range(-MAX_OFFSET, MAX_OFFSET) * shake,
			randf_range(-MAX_OFFSET, MAX_OFFSET) * shake
		)
		_camera.rotation = randf_range(-MAX_ROTATION, MAX_ROTATION) * shake
	if _trauma <= 0 and is_instance_valid(_camera):
		_camera.offset = Vector2.ZERO
		_camera.rotation = 0.0

func add_trauma(amount: float) -> void:
	_trauma = minf(1.0, _trauma + amount)

func set_camera(cam: Camera2D) -> void:
	_camera = cam
