## CharacterAnimator.gd
## 角色動畫控制 - 支援多幀動畫 Spritesheet
## 掛載在 Cannon 的 CannonSprite Sprite2D 上

extends Sprite2D

enum AnimState { IDLE, ATTACK, BIGWIN }

# 動畫設定（對應 upgrade_idle_8frames.py 的 metadata）
# idle 升級：4幀 → 8幀，fps 4 → 8（更流暢的呼吸感）
const ANIM_CONFIG = {
	"idle":   {"row": 0, "frames": 8, "fps": 8.0},
	"attack": {"row": 1, "frames": 3, "fps": 8.0},
	"bigwin": {"row": 2, "frames": 4, "fps": 6.0},
}

const FRAME_SIZE = 96
const COLS = 8  # 升級：4 → 8 cols

var _current_char: String = "chiikawa"
var _current_state: AnimState = AnimState.IDLE
var _current_frame: int = 0
var _frame_timer: float = 0.0
var _attack_timer: float = 0.0
const ATTACK_DURATION = 0.4

# 各角色的 Spritesheet 路徑
const CHAR_SHEETS = {
	"chiikawa":  "res://assets/sprites/sheets/chiikawa_animated.png",
	"hachiware": "res://assets/sprites/sheets/hachiware_animated.png",
	"usagi":     "res://assets/sprites/sheets/usagi_animated.png",
}

# 備用：單幀靜態圖
const CHAR_SPRITES = {
	"chiikawa": {
		"idle":   "res://assets/sprites/characters/chiikawa_idle.png",
		"attack": "res://assets/sprites/characters/chiikawa_attack.png",
		"bigwin": "res://assets/sprites/characters/chiikawa_bigwin.png",
	},
	"hachiware": {
		"idle":   "res://assets/sprites/characters/hachiware_idle.png",
		"attack": "res://assets/sprites/characters/hachiware_attack.png",
		"bigwin": "res://assets/sprites/characters/hachiware_bigwin.png",
	},
	"usagi": {
		"idle":   "res://assets/sprites/characters/usagi_idle.png",
		"attack": "res://assets/sprites/characters/usagi_attack.png",
		"bigwin": "res://assets/sprites/characters/usagi_bigwin.png",
	},
}

var _use_spritesheet: bool = false
var _sheet_texture: Texture2D = null

func _ready() -> void:
	texture_filter = CanvasItem.TEXTURE_FILTER_NEAREST
	scale = Vector2(1.0, 1.0)
	
	GameManager.player_updated.connect(_on_player_updated)
	GameManager.attack_result.connect(_on_attack_result)
	GameManager.reward_received.connect(_on_reward_received)
	
	_load_character("chiikawa")

func _process(delta: float) -> void:
	# 攻擊計時
	if _attack_timer > 0:
		_attack_timer -= delta
		if _attack_timer <= 0:
			_set_state(AnimState.IDLE)
	
	# 動畫幀更新
	if _use_spritesheet:
		_update_animation(delta)

func _update_animation(delta: float) -> void:
	var state_name = _get_state_name()
	var config = ANIM_CONFIG.get(state_name, ANIM_CONFIG["idle"])
	var fps = config["fps"]
	var frame_count = config["frames"]
	var row = config["row"]
	
	_frame_timer += delta
	if _frame_timer >= 1.0 / fps:
		_frame_timer = 0.0
		_current_frame = (_current_frame + 1) % frame_count
		_update_sprite_region(row, _current_frame)

func _update_sprite_region(row: int, col: int) -> void:
	if not _use_spritesheet or _sheet_texture == null:
		return
	
	# 設定 AtlasTexture 裁切區域
	var atlas = AtlasTexture.new()
	atlas.atlas = _sheet_texture
	atlas.region = Rect2(col * FRAME_SIZE, row * FRAME_SIZE, FRAME_SIZE, FRAME_SIZE)
	atlas.filter_clip = true
	texture = atlas

func _load_character(char_id: String) -> void:
	_current_char = char_id
	_current_frame = 0
	_frame_timer = 0.0
	
	# 嘗試載入 Spritesheet
	var sheet_path = CHAR_SHEETS.get(char_id, "")
	if sheet_path != "" and ResourceLoader.exists(sheet_path):
		_sheet_texture = load(sheet_path)
		_use_spritesheet = true
		_update_sprite_region(0, 0)
		print("[CharAnim] Using spritesheet for: ", char_id)
	else:
		# 備用：靜態圖
		_use_spritesheet = false
		_load_static_sprite()
		print("[CharAnim] Using static sprite for: ", char_id)

func _load_static_sprite() -> void:
	var state_name = _get_state_name()
	var char_sprites = CHAR_SPRITES.get(_current_char, CHAR_SPRITES["chiikawa"])
	var path = char_sprites.get(state_name, "")
	if path != "" and ResourceLoader.exists(path):
		texture = load(path)

func _get_state_name() -> String:
	match _current_state:
		AnimState.ATTACK: return "attack"
		AnimState.BIGWIN: return "bigwin"
		_: return "idle"

func _set_state(new_state: AnimState) -> void:
	if _current_state == new_state:
		return
	_current_state = new_state
	_current_frame = 0
	_frame_timer = 0.0
	
	if _use_spritesheet:
		var config = ANIM_CONFIG.get(_get_state_name(), ANIM_CONFIG["idle"])
		_update_sprite_region(config["row"], 0)
	else:
		_load_static_sprite()

func _on_player_updated(data: Dictionary) -> void:
	var char_id = data.get("character_id", "chiikawa")
	if char_id != _current_char:
		_load_character(char_id)

func _on_attack_result(result: Dictionary) -> void:
	if result.get("is_hit", false):
		_set_state(AnimState.ATTACK)
		_attack_timer = ATTACK_DURATION

func _on_reward_received(reward: Dictionary) -> void:
	if reward.get("multiplier", 1.0) >= 20:
		_set_state(AnimState.BIGWIN)
		_attack_timer = 1.5
