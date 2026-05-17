## AudioManager.gd
## 音效管理（規格書 10章）
## Autoload 單例

extends Node
class_name AudioManagerClass

# 音效類型（規格書 10.2）
enum SFX {
	ATTACK_FIRE,
	HIT,
	KILL,
	COIN_DROP,
	REWARD_BAG,
	BOSS_WARNING,
	BONUS_READY,
	WEED_PULL,
	BIG_WIN,
}

# BGM 類型（規格書 10.1）
enum BGM {
	MAIN_GAME,
	BOSS_ENTER,
	BOSS_RAGE,
	BONUS_GAME,
	BIG_WIN_FANFARE,
}

var _sfx_players: Array[AudioStreamPlayer] = []
var _bgm_player: AudioStreamPlayer
var _current_bgm: BGM = BGM.MAIN_GAME

const SFX_POOL_SIZE = 8

func _ready() -> void:
	# 建立 SFX 音效池
	for i in SFX_POOL_SIZE:
		var player = AudioStreamPlayer.new()
		player.bus = "SFX"
		add_child(player)
		_sfx_players.append(player)

	# BGM 播放器
	_bgm_player = AudioStreamPlayer.new()
	_bgm_player.bus = "BGM"
	add_child(_bgm_player)

## 播放音效
func play_sfx(sfx: SFX) -> void:
	var path = _get_sfx_path(sfx)
	if path == "" or not ResourceLoader.exists(path):
		return
	# 找空閒的播放器
	for player in _sfx_players:
		if not player.playing:
			player.stream = load(path)
			player.play()
			return
	# 所有播放器都忙，用第一個（覆蓋）
	_sfx_players[0].stream = load(path)
	_sfx_players[0].play()

## 播放 BGM
func play_bgm(bgm: BGM) -> void:
	if _current_bgm == bgm and _bgm_player.playing:
		return
	_current_bgm = bgm
	var path = _get_bgm_path(bgm)
	if path == "" or not ResourceLoader.exists(path):
		return
	_bgm_player.stream = load(path)
	_bgm_player.play()

## 停止 BGM（BOSS 出場前 0.5 秒靜音）
func stop_bgm_briefly() -> void:
	_bgm_player.stop()
	await get_tree().create_timer(0.5).timeout
	play_bgm(_current_bgm)

## 依角色播放攻擊音效
func play_attack_by_character(character_id: String) -> void:
	match character_id:
		"chiikawa":
			play_sfx(SFX.ATTACK_FIRE)
		"hachiware":
			# 小八用獨立音效
			var path = "res://assets/audio/sfx/attack_fire_hachiware.wav"
			if ResourceLoader.exists(path):
				for player in _sfx_players:
					if not player.playing:
						player.stream = load(path)
						player.play()
						return
		"usagi":
			# 烏薩奇用獨立音效
			var path = "res://assets/audio/sfx/attack_fire_usagi.wav"
			if ResourceLoader.exists(path):
				for player in _sfx_players:
					if not player.playing:
						player.stream = load(path)
						player.play()
						return

func _get_sfx_path(sfx: SFX) -> String:
	match sfx:
		SFX.ATTACK_FIRE: return "res://assets/audio/sfx/attack_fire.wav"
		SFX.HIT: return "res://assets/audio/sfx/hit.wav"
		SFX.KILL: return "res://assets/audio/sfx/kill.wav"
		SFX.COIN_DROP: return "res://assets/audio/sfx/coin_drop.wav"
		SFX.REWARD_BAG: return "res://assets/audio/sfx/reward_bag.wav"
		SFX.BOSS_WARNING: return "res://assets/audio/sfx/boss_warning.wav"
		SFX.BONUS_READY: return "res://assets/audio/sfx/bonus_ready.wav"
		SFX.WEED_PULL: return "res://assets/audio/sfx/weed_pull.wav"
		SFX.BIG_WIN: return "res://assets/audio/sfx/big_win.wav"
	return ""

func _get_bgm_path(bgm: BGM) -> String:
	match bgm:
		BGM.MAIN_GAME: return "res://assets/audio/bgm/main_game.wav"
		BGM.BOSS_ENTER: return "res://assets/audio/bgm/boss_enter.wav"
		BGM.BOSS_RAGE: return "res://assets/audio/bgm/boss_enter.wav"  # 暫用同一首
		BGM.BONUS_GAME: return "res://assets/audio/bgm/bonus_game.wav"
		BGM.BIG_WIN_FANFARE: return "res://assets/audio/sfx/big_win.wav"
	return ""
