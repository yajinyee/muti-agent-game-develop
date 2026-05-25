## CharacterAnimator.gd — 角色動畫管理
## character-animation-agent 負責維護
## 管理三個角色的 idle/attack/bigwin 動畫狀態
extends Node2D

# ── 角色精靈路徑 ──────────────────────────────────────────────
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

const CHAR_COLORS = {
	"chiikawa": Color(1.0, 0.6, 0.8),
	"hachiware": Color(0.4, 0.6, 1.0),
	"usagi": Color(1.0, 0.9, 0.2),
}

# ── 動畫狀態 ──────────────────────────────────────────────────
enum AnimState { IDLE, ATTACK, BIGWIN }

var _sprite: Sprite2D = null
var _current_char: String = "chiikawa"
var _anim_state: AnimState = AnimState.IDLE
var _anim_timer: float = 0.0
var _idle_bob_timer: float = 0.0

# 待機呼吸動畫
const IDLE_BOB_SPEED = 1.8
const IDLE_BOB_AMOUNT = 3.0

func _ready() -> void:
	_sprite = Sprite2D.new()
	_sprite.texture_filter = CanvasItem.TEXTURE_FILTER_NEAREST
	_sprite.scale = Vector2(3.0, 3.0)
	add_child(_sprite)
	GameManager.player_updated.connect(_on_player_updated)
	_load_character("chiikawa")

func _process(delta: float) -> void:
	_idle_bob_timer += delta * IDLE_BOB_SPEED
	match _anim_state:
		AnimState.IDLE:
			_update_idle(delta)
		AnimState.ATTACK:
			_anim_timer -= delta
			if _anim_timer <= 0:
				_set_state(AnimState.IDLE)
		AnimState.BIGWIN:
			_anim_timer -= delta
			if _anim_timer <= 0:
				_set_state(AnimState.IDLE)

func _update_idle(_delta: float) -> void:
	# 上下呼吸動畫
	position.y = sin(_idle_bob_timer) * IDLE_BOB_AMOUNT

# ── 公開 API ──────────────────────────────────────────────────

func play_attack() -> void:
	_set_state(AnimState.ATTACK)
	_anim_timer = 0.15
	# 攻擊時向前傾
	var tween = create_tween()
	tween.tween_property(self, "rotation_degrees", -15.0, 0.06)
	tween.tween_property(self, "rotation_degrees", 0.0, 0.09)

func play_bigwin() -> void:
	_set_state(AnimState.BIGWIN)
	_anim_timer = 1.2
	# 大獎跳躍
	var tween = create_tween()
	tween.tween_property(self, "position:y", position.y - 25, 0.15)
	tween.tween_property(self, "position:y", position.y, 0.15)
	tween.tween_property(self, "position:y", position.y - 15, 0.10)
	tween.tween_property(self, "position:y", position.y, 0.10)
	# 閃光
	if is_instance_valid(_sprite):
		var flash = _sprite.create_tween()
		flash.tween_property(_sprite, "modulate", Color(3.0, 3.0, 1.0), 0.08)
		flash.tween_property(_sprite, "modulate", Color.WHITE, 0.15)

# ── 內部 ──────────────────────────────────────────────────────

func _on_player_updated(data: Dictionary) -> void:
	var char_id = data.get("character_id", "chiikawa")
	if char_id != _current_char:
		_load_character(char_id)

func _load_character(char_id: String) -> void:
	_current_char = char_id
	_set_state(AnimState.IDLE)

func _load_sprite_for_state(state_name: String) -> void:
	if not is_instance_valid(_sprite):
		return
	var char_sprites = CHAR_SPRITES.get(_current_char, CHAR_SPRITES["chiikawa"])
	var path = char_sprites.get(state_name, char_sprites["idle"])
	if ResourceLoader.exists(path):
		_sprite.texture = load(path)
	else:
		# 備用：嘗試 idle
		var idle_path = char_sprites.get("idle", "")
		if ResourceLoader.exists(idle_path):
			_sprite.texture = load(idle_path)
		else:
			_sprite.texture = null

func _set_state(s: AnimState) -> void:
	_anim_state = s
	match s:
		AnimState.IDLE:
			_load_sprite_for_state("idle")
			if is_instance_valid(_sprite):
				_sprite.modulate = CHAR_COLORS.get(_current_char, Color.WHITE)
		AnimState.ATTACK:
			_load_sprite_for_state("attack")
			if is_instance_valid(_sprite):
				_sprite.modulate = Color.WHITE
		AnimState.BIGWIN:
			_load_sprite_for_state("bigwin")
			if is_instance_valid(_sprite):
				_sprite.modulate = Color(1.5, 1.5, 0.5)
