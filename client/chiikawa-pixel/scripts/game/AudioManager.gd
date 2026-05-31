## AudioManager.gd — 音效管理
## sfx-agent + bgm-agent 負責維護
extends Node

enum SFX {
	ATTACK_CHIIKAWA,
	ATTACK_HACHIWARE,
	ATTACK_USAGI,
	HIT,
	KILL,
	BIG_WIN,
	COIN_DROP,
	BOSS_WARNING,
	BOSS_ENTER,
	BONUS_READY,
	BONUS_GAME,
	WEED_PULL,
	# DAY-341 Combo 里程碑音效
	COMBO_5,
	COMBO_10,
	COMBO_20,
	COMBO_30,
}

enum BGM {
	MAIN_GAME,
	BOSS_ENTER,
	BOSS_RAGE,
	BONUS_GAME,
}

var _sfx_players: Array[AudioStreamPlayer] = []
var _bgm_player: AudioStreamPlayer = null
var _current_bgm: BGM = BGM.MAIN_GAME

const SFX_PATHS = {
	SFX.ATTACK_CHIIKAWA: "res://assets/audio/sfx/attack_fire.wav",
	SFX.ATTACK_HACHIWARE: "res://assets/audio/sfx/attack_fire_hachiware.wav",
	SFX.ATTACK_USAGI: "res://assets/audio/sfx/attack_fire_usagi.wav",
	SFX.HIT: "res://assets/audio/sfx/hit.wav",
	SFX.KILL: "res://assets/audio/sfx/kill.wav",
	SFX.BIG_WIN: "res://assets/audio/sfx/big_win.wav",
	SFX.COIN_DROP: "res://assets/audio/sfx/coin_drop.wav",
	SFX.BOSS_WARNING: "res://assets/audio/sfx/boss_warning.wav",
	SFX.BOSS_ENTER: "res://assets/audio/sfx/boss_enter.wav",
	SFX.BONUS_READY: "res://assets/audio/sfx/bonus_ready.wav",
	SFX.BONUS_GAME: "res://assets/audio/sfx/bonus_game.wav",
	SFX.WEED_PULL: "res://assets/audio/sfx/weed_pull.wav",
	# DAY-341 Combo 里程碑音效
	SFX.COMBO_5: "res://assets/audio/sfx/combo_5.wav",
	SFX.COMBO_10: "res://assets/audio/sfx/combo_10.wav",
	SFX.COMBO_20: "res://assets/audio/sfx/combo_20.wav",
	SFX.COMBO_30: "res://assets/audio/sfx/combo_30.wav",
}

const BGM_PATHS = {
	BGM.MAIN_GAME: "res://assets/audio/bgm/main_game.wav",
	BGM.BOSS_ENTER: "res://assets/audio/bgm/boss_enter.wav",
	BGM.BOSS_RAGE: "res://assets/audio/bgm/boss_rage.wav",
	BGM.BONUS_GAME: "res://assets/audio/bgm/bonus_game.wav",
}

func _ready() -> void:
	# 建立 BGM 播放器
	_bgm_player = AudioStreamPlayer.new()
	_bgm_player.bus = "Master"
	add_child(_bgm_player)
	# 建立 SFX 播放器池（8個）
	for i in 8:
		var p = AudioStreamPlayer.new()
		p.bus = "Master"
		add_child(p)
		_sfx_players.append(p)

func play_sfx(sfx: SFX) -> void:
	var path = SFX_PATHS.get(sfx, "")
	if path == "" or not ResourceLoader.exists(path):
		return
	# 找空閒的播放器
	for p in _sfx_players:
		if not p.playing:
			p.stream = load(path)
			p.play()
			return

func play_attack_by_character(char_id: String) -> void:
	match char_id:
		"hachiware": play_sfx(SFX.ATTACK_HACHIWARE)
		"usagi": play_sfx(SFX.ATTACK_USAGI)
		_: play_sfx(SFX.ATTACK_CHIIKAWA)

func play_bgm(bgm: BGM) -> void:
	if _current_bgm == bgm and _bgm_player.playing:
		return
	_current_bgm = bgm
	var path = BGM_PATHS.get(bgm, "")
	if path == "" or not ResourceLoader.exists(path):
		return
	var tween = create_tween()
	tween.tween_property(_bgm_player, "volume_db", -40.0, 0.3)
	tween.tween_callback(func():
		_bgm_player.stream = load(path)
		_bgm_player.play()
		var t2 = create_tween()
		t2.tween_property(_bgm_player, "volume_db", 0.0, 0.5)
	)

func stop_bgm() -> void:
	var tween = create_tween()
	tween.tween_property(_bgm_player, "volume_db", -40.0, 0.3)
	tween.tween_callback(func(): _bgm_player.stop())

## DAY-341 Combo 里程碑音效
func play_combo_milestone(combo_count: int) -> void:
	match combo_count:
		5:  play_sfx(SFX.COMBO_5)
		10: play_sfx(SFX.COMBO_10)
		20: play_sfx(SFX.COMBO_20)
		30: play_sfx(SFX.COMBO_30)
