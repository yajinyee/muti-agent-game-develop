## BackgroundManager.gd
## 管理背景圖載入與切換

extends Sprite2D

const BG_PATHS = {
	"normal": "res://assets/sprites/backgrounds/sea_bg.png",
	"boss":   "res://assets/sprites/backgrounds/boss_bg.png",
	"bonus":  "res://assets/sprites/backgrounds/bonus_bg.png",
}

func _ready() -> void:
	position = Vector2(640, 360)
	texture_filter = CanvasItem.TEXTURE_FILTER_LINEAR  # 背景用線性濾波，更平滑
	_load_bg("normal")
	GameManager.game_state_changed.connect(_on_state_changed)

func _load_bg(key: String) -> void:
	var path = BG_PATHS.get(key, BG_PATHS["normal"])
	if ResourceLoader.exists(path):
		texture = load(path)
	else:
		# 備用：深藍色
		var img = Image.create(1280, 720, false, Image.FORMAT_RGB8)
		img.fill(Color(0.05, 0.1, 0.25))
		texture = ImageTexture.create_from_image(img)

func _on_state_changed(state: String) -> void:
	match state:
		"boss_warning", "boss_battle", "boss_result":
			_load_bg("boss")
		"bonus_game", "bonus_ready", "bonus_result":
			_load_bg("bonus")
		_:
			_load_bg("normal")
