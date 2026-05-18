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
	BUBBLE_POP,  # 氣泡破裂音（對應 BubbleLayer 視覺）
	BONUS_TRIGGER,  # Bonus 遊戲開始爆發音（v2）
	BONUS_END,      # Bonus 結算音效（v2）
}

# BGM 類型（規格書 10.1）
enum BGM {
	MAIN_GAME,
	BOSS_ENTER,
	BOSS_BATTLE,   # BOSS Phase 1 戰鬥 BGM（循環，緊張低頻）
	BOSS_RAGE,     # BOSS Phase 2 BGM（加速 15% + 升調）
	BONUS_GAME,
	BIG_WIN_FANFARE,
	UNDERWATER_AMBIENT,  # 海底環境音（循環，低音量背景）
}

var _sfx_players: Array[AudioStreamPlayer] = []
var _bgm_player: AudioStreamPlayer
var _ambient_player: AudioStreamPlayer  # 環境音專用播放器（獨立於 BGM）
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

	# 環境音播放器（獨立，不受 BGM 切換影響）
	_ambient_player = AudioStreamPlayer.new()
	_ambient_player.bus = "SFX"
	_ambient_player.volume_db = -24.0  # 很低音量，純背景感
	add_child(_ambient_player)

## 播放音效
func play_sfx(sfx: SFX) -> void:
	var path = _get_sfx_path(sfx)
	if path == "" or not ResourceLoader.exists(path):
		return
	# 找空閒的播放器
	for player in _sfx_players:
		if not player.playing:
			player.stream = load(path)
			# coin_drop 音量提升 2 dB（audio-review 建議：-7dB → -5dB）
			if sfx == SFX.COIN_DROP:
				player.volume_db = 2.0  # 相對提升 2 dB
			else:
				player.volume_db = 0.0
			player.play()
			return
	# 所有播放器都忙，用第一個（覆蓋）
	_sfx_players[0].stream = load(path)
	_sfx_players[0].volume_db = 2.0 if sfx == SFX.COIN_DROP else 0.0
	_sfx_players[0].play()

## 播放環境音（循環，低音量）
func play_ambient(ambient_type: String = "underwater") -> void:
	var path = ""
	match ambient_type:
		"underwater":
			path = "res://assets/audio/sfx/underwater_ambient.wav"
	if path == "" or not ResourceLoader.exists(path):
		return
	if _ambient_player.playing:
		return  # 已在播放，不重複
	_ambient_player.stream = load(path)
	_ambient_player.play()

## 停止環境音
func stop_ambient() -> void:
	if _ambient_player.playing:
		var tween = create_tween()
		tween.tween_property(_ambient_player, "volume_db", -60.0, 0.5)
		tween.tween_callback(func(): _ambient_player.stop(); _ambient_player.volume_db = -24.0)

## 播放 BGM（帶淡入淡出）
func play_bgm(bgm: BGM, fade_in: float = 0.5) -> void:
	if _current_bgm == bgm and _bgm_player.playing:
		return
	_current_bgm = bgm
	var path = _get_bgm_path(bgm)
	if path == "" or not ResourceLoader.exists(path):
		return

	# 淡出舊 BGM
	if _bgm_player.playing:
		var tween = create_tween()
		tween.tween_property(_bgm_player, "volume_db", -60.0, 0.3)
		tween.tween_callback(func():
			_bgm_player.stop()
			_bgm_player.volume_db = _get_bgm_volume(bgm)
			_bgm_player.stream = load(path)
			_bgm_player.play()
			# 淡入新 BGM
			_bgm_player.volume_db = -60.0
			var tween2 = create_tween()
			tween2.tween_property(_bgm_player, "volume_db", _get_bgm_volume(bgm), fade_in)
			# BOSS Phase 2：音調漸變（避免突兀切換）
			# 從當前音調漸變到目標音調（0.5 秒），比直接跳變更自然
			var target_pitch = _get_bgm_pitch(bgm)
			if bgm == BGM.BOSS_RAGE:
				_bgm_player.pitch_scale = 1.0  # 從正常音調開始
				tween2.parallel().tween_property(_bgm_player, "pitch_scale", target_pitch, 0.5)
			else:
				_bgm_player.pitch_scale = target_pitch
		)
	else:
		_bgm_player.volume_db = -60.0
		_bgm_player.pitch_scale = _get_bgm_pitch(bgm)
		_bgm_player.stream = load(path)
		_bgm_player.play()
		var tween = create_tween()
		tween.tween_property(_bgm_player, "volume_db", _get_bgm_volume(bgm), fade_in)

## 停止 BGM（BOSS 出場前 0.5 秒靜音）
func stop_bgm_briefly() -> void:
	var tween = create_tween()
	tween.tween_property(_bgm_player, "volume_db", -60.0, 0.3)
	tween.tween_callback(func():
		_bgm_player.stop()
	)

## 取得 BGM 目標音量
func _get_bgm_volume(bgm: BGM) -> float:
	match bgm:
		BGM.MAIN_GAME: return -12.0
		BGM.BOSS_ENTER: return -8.0
		BGM.BOSS_BATTLE: return -8.0  # Phase 1 戰鬥，和進場同音量
		BGM.BOSS_RAGE: return -6.0   # Phase 2 更大聲
		BGM.BONUS_GAME: return -10.0
		BGM.BIG_WIN_FANFARE: return -4.0
		BGM.UNDERWATER_AMBIENT: return -24.0  # 環境音很低
	return -10.0

## 取得 BGM 音調（BOSS Phase 2 提高音調）
func _get_bgm_pitch(bgm: BGM) -> float:
	match bgm:
		BGM.BOSS_RAGE: return 1.1   # Phase 2 音調 +10%（規格書 audio-map.json）
		_: return 1.0

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
		SFX.BUBBLE_POP: return "res://assets/audio/sfx/bubble_pop.wav"
		SFX.BONUS_TRIGGER: return "res://assets/audio/sfx/bonus_trigger.wav"
		SFX.BONUS_END: return "res://assets/audio/sfx/bonus_end.wav"
	return ""

func _get_bgm_path(bgm: BGM) -> String:
	match bgm:
		BGM.MAIN_GAME: return "res://assets/audio/bgm/main_game.wav"
		BGM.BOSS_ENTER: return "res://assets/audio/bgm/boss_enter.wav"
		BGM.BOSS_BATTLE: return "res://assets/audio/bgm/boss_battle.wav"  # Phase 1 戰鬥循環
		BGM.BOSS_RAGE: return "res://assets/audio/bgm/boss_rage.wav"   # Phase 2：加速 15% + 升調
		BGM.BONUS_GAME: return "res://assets/audio/bgm/bonus_game.wav"
		BGM.BIG_WIN_FANFARE: return "res://assets/audio/sfx/big_win.wav"
	return ""
