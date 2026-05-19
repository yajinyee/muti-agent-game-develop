## BackgroundManager.gd
## 管理背景圖載入與切換，以及海底動態效果

extends Sprite2D

const BG_PATHS = {
	"normal": "res://assets/sprites/backgrounds/sea_bg.png",
	"boss":   "res://assets/sprites/backgrounds/boss_bg.png",
	"bonus":  "res://assets/sprites/backgrounds/bonus_bg.png",
}

const CAUSTICS_SHADER_PATH = "res://assets/shaders/underwater_caustics.gdshader"
const BubbleLayerScript = preload("res://scripts/game/BubbleLayer.gd")

var _current_bg: String = "normal"
var _bubble_layer: Node2D = null
var _caustics_mat: ShaderMaterial = null  # 焦散光效果（海底背景專用）
# _underwater_overlay 已移至 Main.tscn 的 UnderwaterLayer（DAY-031）

func _ready() -> void:
	position = Vector2(640, 360)
	texture_filter = CanvasItem.TEXTURE_FILTER_LINEAR
	_load_bg("normal")
	GameManager.game_state_changed.connect(_on_state_changed)

	# 建立氣泡/海草/光線層（在背景上方，目標物下方）
	_bubble_layer = BubbleLayerScript.new()
	_bubble_layer.name = "BubbleLayer"
	_bubble_layer.z_index = 1
	get_parent().call_deferred("add_child", _bubble_layer)

	# 預載入焦散 shader
	if ResourceLoader.exists(CAUSTICS_SHADER_PATH):
		_caustics_mat = ShaderMaterial.new()
		_caustics_mat.shader = load(CAUSTICS_SHADER_PATH)
		# 套用到背景（normal 狀態）
		material = _caustics_mat

	# 建立全螢幕水下色差效果（DAY-030）
	# 注意：DAY-031 已將 UnderwaterOverlay 加入 Main.tscn 的 UnderwaterLayer
	# BackgroundManager 不再需要建立第二個 overlay，避免重複效果
	# _underwater_overlay 保留為 null，由 Main.tscn 的節點負責

	# 啟動海底環境音（遊戲開始時）
	call_deferred("_start_initial_ambient")

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
			_switch_bg("boss", false)
		"bonus_game", "bonus_ready", "bonus_result":
			_switch_bg("bonus", false)
		_:
			_switch_bg("normal", true)
	# BGM 切換（依遊戲狀態）
	_switch_bgm(state)

## 帶像素化過場的背景切換
func _start_initial_ambient() -> void:
	if AudioManager != null:
		AudioManager.play_ambient("underwater")
		AudioManager.play_bgm(AudioManager.BGM.MAIN_GAME)

## BGM 切換（依遊戲狀態）
func _switch_bgm(state: String) -> void:
	if AudioManager == null:
		return
	match state:
		"boss_warning":
			# 警告時短暫停止 BGM，製造緊張感
			AudioManager.stop_bgm_briefly()
		"boss_battle":
			# BOSS Phase 1：戰鬥 BGM（循環）
			AudioManager.play_bgm(AudioManager.BGM.BOSS_BATTLE)
		"boss_result":
			# BOSS 結束：短暫靜音後回主 BGM
			AudioManager.stop_bgm_briefly()
		"bonus_ready":
			# Bonus 準備：短暫停止
			AudioManager.stop_bgm_briefly()
		"bonus_game":
			# Bonus 遊戲：歡樂 BGM
			AudioManager.play_bgm(AudioManager.BGM.BONUS_GAME)
		"bonus_result":
			AudioManager.stop_bgm_briefly()
		"normal_play", "special_target_event":
			# 正常遊戲：主 BGM
			AudioManager.play_bgm(AudioManager.BGM.MAIN_GAME)
func _switch_bg(key: String, bubbles_active: bool) -> void:
	if _current_bg == key:
		return
	_current_bg = key

	# 用像素化過場切換背景（更有像素遊戲感）
	HitEffect.pixelate_transition(0.15, 0.2, func():
		_load_bg(key)
		if is_instance_valid(_bubble_layer):
			_bubble_layer.set_active(bubbles_active)
		# 焦散效果只在海底背景時啟用
		if key == "normal" and _caustics_mat != null:
			material = _caustics_mat
		else:
			material = null
		# 環境音：海底狀態播放水聲，其他狀態停止
		if AudioManager != null:
			if key == "normal":
				AudioManager.play_ambient("underwater")
			else:
				AudioManager.stop_ambient()
	)
