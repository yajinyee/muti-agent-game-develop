## ScreenShake.gd — Trauma-based 螢幕震動
## screen-effect-agent 負責維護
extends Node

var _trauma: float = 0.0
var _camera: Camera2D = null

const DECAY = 2.5
const MAX_OFFSET = 12.0
const MAX_ROTATION = 0.04

func _ready() -> void:
	# 自動尋找場景中的 Camera2D
	call_deferred("_find_camera")

func _find_camera() -> void:
	var tree = get_tree()
	if tree == null:
		return
	var root = tree.get_root()
	if root == null:
		return
	_camera = _find_camera_in(root)

func _find_camera_in(node: Node) -> Camera2D:
	if node is Camera2D:
		return node
	for child in node.get_children():
		var result = _find_camera_in(child)
		if result != null:
			return result
	return null

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
