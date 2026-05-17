## BackgroundManager.gd
## 管理背景圖載入與切換，以及海底氣泡動畫

extends Sprite2D

const BG_PATHS = {
	"normal": "res://assets/sprites/backgrounds/sea_bg.png",
	"boss":   "res://assets/sprites/backgrounds/boss_bg.png",
	"bonus":  "res://assets/sprites/backgrounds/bonus_bg.png",
}

const BubbleLayerScript = preload("res://scripts/game/BubbleLayer.gd")

var _current_bg: String = "normal"
var _bubble_layer: Node2D = null

func _ready() -> void:
	position = Vector2(640, 360)
	texture_filter = CanvasItem.TEXTURE_FILTER_LINEAR
	_load_bg("normal")
	GameManager.game_state_changed.connect(_on_state_changed)

	# 建立氣泡層（在背景上方，目標物下方）
	_bubble_layer = BubbleLayerScript.new()
	_bubble_layer.name = "BubbleLayer"
	_bubble_layer.z_index = 1
	get_parent().call_deferred("add_child", _bubble_layer)

func _load_bg(key: String) -> void:
	var path = BG_PATHS.get(key, BG_PATHS["normal"])
	if ResourceLoader.exists(path):
		texture = load(path)
	else:
		var img = Image.create(1280, 720, false, Image.FORMAT_RGB8)
		img.fill(Color(0.05, 0.1, 0.25))
		texture = ImageTexture.create_from_image(img)

func _on_state_changed(state: String) -> void:
	match state:
		"boss_warning", "boss_battle", "boss_result":
			_load_bg("boss")
			_current_bg = "boss"
			if is_instance_valid(_bubble_layer):
				_bubble_layer.set_active(false)
		"bonus_game", "bonus_ready", "bonus_result":
			_load_bg("bonus")
			_current_bg = "bonus"
			if is_instance_valid(_bubble_layer):
				_bubble_layer.set_active(false)
		_:
			_load_bg("normal")
			_current_bg = "normal"
			if is_instance_valid(_bubble_layer):
				_bubble_layer.set_active(true)
