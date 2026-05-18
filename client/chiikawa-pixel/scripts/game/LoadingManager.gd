## LoadingManager.gd
## 資產預載入管理（Autoload 單例）
## 在遊戲啟動時背景預載入所有重要資產，避免首次使用時的卡頓
##
## 使用方式：
##   LoadingManager.preload_all()          # 啟動時呼叫
##   LoadingManager.get_progress()         # 取得進度 0.0~1.0
##   LoadingManager.is_ready()             # 是否全部載入完成
##   LoadingManager.get_texture(path)      # 取得已快取的 Texture2D
##   LoadingManager.get_audio(path)        # 取得已快取的 AudioStream

extends Node

signal loading_progress(progress: float)
signal loading_complete()

# 快取字典
var _texture_cache: Dictionary = {}   # path -> Texture2D
var _audio_cache: Dictionary = {}     # path -> AudioStream
var _shader_cache: Dictionary = {}    # path -> Shader

var _total_assets: int = 0
var _loaded_assets: int = 0
var _is_ready: bool = false
var _loading_started: bool = false

# ── 需要預載入的資產清單 ──────────────────────────────────────

const TEXTURES_TO_PRELOAD = [
	# 角色 Spritesheet
	"res://assets/sprites/characters/chiikawa_sheet.png",
	"res://assets/sprites/characters/hachiware_sheet.png",
	"res://assets/sprites/characters/usagi_sheet.png",
	# 目標物 Spritesheet
	"res://assets/sprites/targets/targets_sheet.png",
	# 背景
	"res://assets/sprites/backgrounds/sea_bg.png",
	"res://assets/sprites/backgrounds/boss_bg.png",
	"res://assets/sprites/backgrounds/bonus_bg.png",
	# 特效
	"res://assets/sprites/effects/hit_effect.png",
	"res://assets/sprites/effects/death_particles.png",
	"res://assets/sprites/effects/warning.png",
	# 投射物
	"res://assets/sprites/effects/projectile_chiikawa.png",
	"res://assets/sprites/effects/projectile_hachiware.png",
	"res://assets/sprites/effects/projectile_usagi.png",
	# UI
	"res://assets/sprites/ui/coin.png",
	"res://assets/sprites/ui/reward_bag.png",
]

const AUDIO_TO_PRELOAD = [
	# SFX
	"res://assets/audio/sfx/attack_fire.wav",
	"res://assets/audio/sfx/attack_fire_hachiware.wav",
	"res://assets/audio/sfx/attack_fire_usagi.wav",
	"res://assets/audio/sfx/hit.wav",
	"res://assets/audio/sfx/kill.wav",
	"res://assets/audio/sfx/coin_drop.wav",
	"res://assets/audio/sfx/reward_bag.wav",
	"res://assets/audio/sfx/boss_warning.wav",
	"res://assets/audio/sfx/bonus_ready.wav",
	"res://assets/audio/sfx/weed_pull.wav",
	"res://assets/audio/sfx/big_win.wav",
	"res://assets/audio/sfx/bubble_pop.wav",
	"res://assets/audio/sfx/bonus_trigger.wav",
	"res://assets/audio/sfx/bonus_end.wav",
	# BGM（只預載入 main_game，其他按需載入）
	"res://assets/audio/bgm/main_game.wav",
]

const SHADERS_TO_PRELOAD = [
	"res://assets/shaders/hit_flash.gdshader",
	"res://assets/shaders/outline.gdshader",
	"res://assets/shaders/rainbow_glow.gdshader",
	"res://assets/shaders/shockwave_distortion.gdshader",
	"res://assets/shaders/pixelate_transition.gdshader",
	"res://assets/shaders/underwater_caustics.gdshader",
	"res://assets/shaders/water_surface.gdshader",
]

# ── 公開 API ──────────────────────────────────────────────────

## 啟動預載入（在 Main.gd 的 _ready() 中呼叫）
func preload_all() -> void:
	if _loading_started:
		return
	_loading_started = true

	# 計算總資產數
	_total_assets = TEXTURES_TO_PRELOAD.size() + AUDIO_TO_PRELOAD.size() + SHADERS_TO_PRELOAD.size()
	_loaded_assets = 0

	print("[LoadingManager] 開始預載入 %d 個資產..." % _total_assets)

	# 使用 ResourceLoader.load_threaded_request 背景載入
	# 先載入 Textures（最重要）
	for path in TEXTURES_TO_PRELOAD:
		if ResourceLoader.exists(path):
			ResourceLoader.load_threaded_request(path)
		else:
			_loaded_assets += 1  # 不存在的資產直接跳過

	for path in AUDIO_TO_PRELOAD:
		if ResourceLoader.exists(path):
			ResourceLoader.load_threaded_request(path)
		else:
			_loaded_assets += 1

	for path in SHADERS_TO_PRELOAD:
		if ResourceLoader.exists(path):
			ResourceLoader.load_threaded_request(path)
		else:
			_loaded_assets += 1

## 取得載入進度（0.0 ~ 1.0）
func get_progress() -> float:
	if _total_assets == 0:
		return 1.0
	return float(_loaded_assets) / float(_total_assets)

## 是否全部載入完成
func is_ready() -> bool:
	return _is_ready

## 取得已快取的 Texture2D（如果未快取則同步載入）
func get_texture(path: String) -> Texture2D:
	if _texture_cache.has(path):
		return _texture_cache[path]
	if ResourceLoader.exists(path):
		var tex = load(path) as Texture2D
		if tex:
			_texture_cache[path] = tex
		return tex
	return null

## 取得已快取的 AudioStream（如果未快取則同步載入）
func get_audio(path: String) -> AudioStream:
	if _audio_cache.has(path):
		return _audio_cache[path]
	if ResourceLoader.exists(path):
		var audio = load(path) as AudioStream
		if audio:
			_audio_cache[path] = audio
		return audio
	return null

## 取得已快取的 Shader
func get_shader(path: String) -> Shader:
	if _shader_cache.has(path):
		return _shader_cache[path]
	if ResourceLoader.exists(path):
		var shader = load(path) as Shader
		if shader:
			_shader_cache[path] = shader
		return shader
	return null

# ── 內部處理 ──────────────────────────────────────────────────

func _process(_delta: float) -> void:
	if _is_ready or not _loading_started:
		return

	# 輪詢背景載入狀態
	var all_done = true
	var newly_loaded = 0

	for path in TEXTURES_TO_PRELOAD:
		if not ResourceLoader.exists(path):
			continue
		if _texture_cache.has(path):
			continue
		var status = ResourceLoader.load_threaded_get_status(path)
		match status:
			ResourceLoader.THREAD_LOAD_LOADED:
				var res = ResourceLoader.load_threaded_get(path)
				if res is Texture2D:
					_texture_cache[path] = res
				newly_loaded += 1
			ResourceLoader.THREAD_LOAD_IN_PROGRESS:
				all_done = false
			ResourceLoader.THREAD_LOAD_FAILED:
				# 載入失敗，跳過（不阻塞遊戲）
				newly_loaded += 1
				push_warning("[LoadingManager] 載入失敗: " + path)

	for path in AUDIO_TO_PRELOAD:
		if not ResourceLoader.exists(path):
			continue
		if _audio_cache.has(path):
			continue
		var status = ResourceLoader.load_threaded_get_status(path)
		match status:
			ResourceLoader.THREAD_LOAD_LOADED:
				var res = ResourceLoader.load_threaded_get(path)
				if res is AudioStream:
					_audio_cache[path] = res
				newly_loaded += 1
			ResourceLoader.THREAD_LOAD_IN_PROGRESS:
				all_done = false
			ResourceLoader.THREAD_LOAD_FAILED:
				newly_loaded += 1
				push_warning("[LoadingManager] 載入失敗: " + path)

	for path in SHADERS_TO_PRELOAD:
		if not ResourceLoader.exists(path):
			continue
		if _shader_cache.has(path):
			continue
		var status = ResourceLoader.load_threaded_get_status(path)
		match status:
			ResourceLoader.THREAD_LOAD_LOADED:
				var res = ResourceLoader.load_threaded_get(path)
				if res is Shader:
					_shader_cache[path] = res
				newly_loaded += 1
			ResourceLoader.THREAD_LOAD_IN_PROGRESS:
				all_done = false
			ResourceLoader.THREAD_LOAD_FAILED:
				newly_loaded += 1
				push_warning("[LoadingManager] 載入失敗: " + path)

	if newly_loaded > 0:
		_loaded_assets = min(_loaded_assets + newly_loaded, _total_assets)
		emit_signal("loading_progress", get_progress())

	if all_done:
		_is_ready = true
		_loaded_assets = _total_assets
		print("[LoadingManager] 預載入完成！快取: %d textures, %d audio, %d shaders" % [
			_texture_cache.size(), _audio_cache.size(), _shader_cache.size()
		])
		emit_signal("loading_complete")
		set_process(false)  # 完成後停止輪詢
