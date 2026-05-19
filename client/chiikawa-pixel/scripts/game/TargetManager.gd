## TargetManager.gd
## 管理畫面上的目標物節點（使用像素 Sprite）
## 掛載在 GameScene 上

extends Node2D

# 目標 Sprite 路徑對應表（備用，Spritesheet 載入失敗時使用）
const TARGET_SPRITES = {
	"T001": "res://assets/sprites/targets/T001_grass.png",
	"T002": "res://assets/sprites/targets/T002_bug_g.png",
	"T003": "res://assets/sprites/targets/T003_bug_r.png",
	"T004": "res://assets/sprites/targets/T004_bug_b.png",
	"T005": "res://assets/sprites/targets/T005_pudding.png",
	"T006": "res://assets/sprites/targets/T006_mushroom.png",
	"T101": "res://assets/sprites/targets/T101_mimic.png",
	"T102": "res://assets/sprites/targets/T102_chest.png",
	"T103": "res://assets/sprites/targets/T103_meteor.png",
	"T104": "res://assets/sprites/targets/T104_gold_grass.png",
	"T105": "res://assets/sprites/targets/T105_coin_fish.png",
	"B001": "res://assets/sprites/targets/B001_boss.png",
}

# 游泳動畫 Spritesheet（2幀橫排，128x64）
# 幀0：向上彎曲，幀1：向下彎曲，4fps 交替
const SWIM_SHEET_TARGETS = [
	"T001", "T002", "T003", "T004", "T005", "T006",
	"T101", "T102", "T103", "T104", "T105"
]
const SWIM_ANIM_FPS: float = 4.0  # 每 0.25 秒切換一幀

# B001 BOSS 動畫 Spritesheet（512x384，4幀×3狀態×128px）
const BOSS_SHEET_PATH = "res://assets/sprites/targets/B001_boss_sheet.png"
const BOSS_FRAME_SIZE = 128
const BOSS_COLS = 4
# Row 0: idle, Row 1: phase2, Row 2: death
const BOSS_ROW_IDLE   = 0
const BOSS_ROW_PHASE2 = 1
const BOSS_ROW_DEATH  = 2

# Spritesheet 中各目標的 UV 座標（來自 targets_sheet.json，cell_size=64）
const SHEET_REGIONS = {
	"T001": Rect2(0, 0, 64, 64),
	"T002": Rect2(64, 0, 64, 64),
	"T003": Rect2(128, 0, 64, 64),
	"T004": Rect2(192, 0, 64, 64),
	"T005": Rect2(0, 64, 64, 64),
	"T006": Rect2(64, 64, 64, 64),
	"T101": Rect2(128, 64, 64, 64),
	"T102": Rect2(192, 64, 64, 64),
	"T103": Rect2(0, 128, 64, 64),
	"T104": Rect2(64, 128, 64, 64),
	"T105": Rect2(128, 128, 64, 64),
}

# 命中特效 Sprite
const HIT_EFFECTS = {
	"chiikawa": "res://assets/sprites/effects/hit_chiikawa.png",
	"hachiware": "res://assets/sprites/effects/hit_hachiware.png",
	"usagi":    "res://assets/sprites/effects/hit_usagi.png",
}

# 投射物 Sprite
const PROJECTILE_SPRITES = {
	"chiikawa": "res://assets/sprites/effects/projectile_chiikawa.png",
	"hachiware": "res://assets/sprites/effects/projectile_hachiware.png",
	"usagi":    "res://assets/sprites/effects/projectile_usagi.png",
}

# 目標節點字典
var _target_nodes: Dictionary = {}  # instance_id -> Node2D

# ---- 資源快取（避免每次都 load，提升效能）----
var _cached_textures: Dictionary = {}   # path -> Texture2D
var _cached_outline_shader: Shader = null
var _cached_hit_flash_shader: Shader = null
var _cached_pixel_font: Font = null
var _targets_sheet: Texture2D = null    # 目標物 Spritesheet（減少 draw call）
var _atlas_textures: Dictionary = {}    # def_id -> AtlasTexture（快取裁切結果）
var _boss_sheet: Texture2D = null       # B001 BOSS 動畫 Spritesheet
var _boss_atlas_cache: Dictionary = {}  # "row_col" -> AtlasTexture

# BOSS 動畫狀態
var _boss_anim_timer: float = 0.0
var _boss_anim_frame: int = 0
var _boss_anim_row: int = BOSS_ROW_IDLE
const BOSS_ANIM_FPS: float = 4.0  # idle/phase2 4fps，death 8fps

# 游泳動畫狀態（全局計時器，所有目標物共用）
var _swim_anim_timer: float = 0.0
var _swim_anim_frame: int = 0  # 0 或 1
# 游泳動畫 AtlasTexture 快取：def_id -> [frame0_atlas, frame1_atlas]
var _swim_atlas_cache: Dictionary = {}

func _ready() -> void:
	GameManager.target_spawned.connect(_on_target_spawned)
	GameManager.target_updated.connect(_on_target_updated)
	GameManager.target_killed.connect(_on_target_killed)
	GameManager.boss_event.connect(_on_boss_event)
	# 初始化 TargetPool（預建立 24 個空殼節點，避免高頻 GC）
	TargetPool.init_pool(self)
	# 預載入常用資源
	_preload_resources()

## 預載入常用資源（在 _ready 時一次性載入，避免遊戲中卡頓）
func _preload_resources() -> void:
	# 優先載入 Spritesheet（一張圖包含所有目標，減少 draw call）
	var sheet_path = "res://assets/sprites/sheets/targets_sheet.png"
	if ResourceLoader.exists(sheet_path):
		_targets_sheet = load(sheet_path)
		# 預建立所有 AtlasTexture（裁切快取）
		for def_id in SHEET_REGIONS:
			var atlas = AtlasTexture.new()
			atlas.atlas = _targets_sheet
			atlas.region = SHEET_REGIONS[def_id]
			_atlas_textures[def_id] = atlas

	# BOSS 動畫 Spritesheet 載入
	if ResourceLoader.exists(BOSS_SHEET_PATH):
		_boss_sheet = load(BOSS_SHEET_PATH)
		# 預建立所有 BOSS 幀的 AtlasTexture
		for row in range(3):
			for col in range(BOSS_COLS):
				var atlas = AtlasTexture.new()
				atlas.atlas = _boss_sheet
				atlas.region = Rect2(col * BOSS_FRAME_SIZE, row * BOSS_FRAME_SIZE,
									 BOSS_FRAME_SIZE, BOSS_FRAME_SIZE)
				_boss_atlas_cache["%d_%d" % [row, col]] = atlas
	else:
		# 備用：靜態 B001 PNG
		var boss_path = TARGET_SPRITES["B001"]
		if ResourceLoader.exists(boss_path):
			_cached_textures[boss_path] = load(boss_path)

	# 備用：預載入所有獨立目標 Sprite（Spritesheet 載入失敗時使用）
	if _targets_sheet == null:
		for def_id in TARGET_SPRITES:
			if def_id == "B001":
				continue  # BOSS 已單獨處理
			var path = TARGET_SPRITES[def_id]
			if ResourceLoader.exists(path):
				_cached_textures[path] = load(path)

	# 預載入 shader
	var outline_path = "res://assets/shaders/outline.gdshader"
	if ResourceLoader.exists(outline_path):
		_cached_outline_shader = load(outline_path)

	var hit_flash_path = "res://assets/shaders/hit_flash.gdshader"
	if ResourceLoader.exists(hit_flash_path):
		_cached_hit_flash_shader = load(hit_flash_path)

	# 預載入像素字體
	var font_path = "res://assets/fonts/pixel8.fnt"
	if ResourceLoader.exists(font_path):
		_cached_pixel_font = load(font_path)

	# 預載入游泳動畫 Spritesheet（2幀橫排，128x64）
	for def_id in SWIM_SHEET_TARGETS:
		var swim_path = "res://assets/sprites/targets/%s_swim.png" % def_id.to_lower()
		# 注意：路徑要和實際檔名一致
		var actual_path = _get_swim_sheet_path(def_id)
		if ResourceLoader.exists(actual_path):
			var swim_tex = load(actual_path)
			# 建立兩幀的 AtlasTexture
			var frame0 = AtlasTexture.new()
			frame0.atlas = swim_tex
			frame0.region = Rect2(0, 0, 64, 64)
			var frame1 = AtlasTexture.new()
			frame1.atlas = swim_tex
			frame1.region = Rect2(64, 0, 64, 64)
			_swim_atlas_cache[def_id] = [frame0, frame1]

## 取得游泳動畫 Spritesheet 路徑
func _get_swim_sheet_path(def_id: String) -> String:
	# 路徑格式：T001_grass_swim.png
	var name_map = {
		"T001": "T001_grass",
		"T002": "T002_bug_g",
		"T003": "T003_bug_r",
		"T004": "T004_bug_b",
		"T005": "T005_pudding",
		"T006": "T006_mushroom",
		"T101": "T101_mimic",
		"T102": "T102_chest",
		"T103": "T103_meteor",
		"T104": "T104_gold_grass",
		"T105": "T105_coin_fish",
	}
	var base_name = name_map.get(def_id, "")
	if base_name == "":
		return ""
	return "res://assets/sprites/targets/%s_swim.png" % base_name

## 取得目標 Texture（優先用 AtlasTexture，備用獨立 PNG）
func _get_target_texture(def_id: String) -> Texture2D:
	# BOSS 用動畫 Spritesheet 的第一幀（idle 幀0）
	if def_id == "B001":
		if _boss_atlas_cache.has("0_0"):
			return _boss_atlas_cache["0_0"]
		# 備用：靜態 PNG
		var path = TARGET_SPRITES.get("B001", "")
		return _cached_textures.get(path, null)
	# 其他目標優先用 Spritesheet AtlasTexture
	if _atlas_textures.has(def_id):
		return _atlas_textures[def_id]
	# 備用：獨立 PNG
	var path = TARGET_SPRITES.get(def_id, "")
	return _get_texture(path)

## 取得 BOSS 動畫幀 AtlasTexture
func _get_boss_frame(row: int, col: int) -> Texture2D:
	var key = "%d_%d" % [row, col]
	if _boss_atlas_cache.has(key):
		return _boss_atlas_cache[key]
	# 備用：靜態 PNG
	var path = TARGET_SPRITES.get("B001", "")
	return _cached_textures.get(path, null)

## 取得快取 Texture（避免重複 load）
func _get_texture(path: String) -> Texture2D:
	if _cached_textures.has(path):
		return _cached_textures[path]
	if ResourceLoader.exists(path):
		var tex = load(path)
		_cached_textures[path] = tex
		return tex
	return null

## 目標 HP 更新
func _on_target_updated(data: Dictionary) -> void:
	var instance_id = data.get("instance_id", "")
	update_target_hp(instance_id, data.get("hp", 0), data.get("max_hp", 1))

	# 受擊閃白效果（使用 shader）
	if _target_nodes.has(instance_id):
		var node = _target_nodes[instance_id]
		if is_instance_valid(node):
			_flash_hit(node)

	# T102 寶箱怪：受擊後加速逃跑（規格書 26.2）
	if data.get("is_fleeing", false):
		if _target_nodes.has(instance_id):
			var node = _target_nodes[instance_id]
			if is_instance_valid(node):
				var current_speed = node.get_meta("speed", 70.0)
				node.set_meta("flee_speed", current_speed * 2.5)  # 2.5x 加速
				node.set_meta("behavior", "flee")
				# 視覺反饋：閃爍紅色
				var tween = create_tween()
				tween.tween_property(node, "modulate", Color(2.0, 0.5, 0.5, 1.0), 0.06)
				tween.tween_property(node, "modulate", Color.WHITE, 0.06)
				tween.tween_property(node, "modulate", Color(2.0, 0.5, 0.5, 1.0), 0.06)
				tween.tween_property(node, "modulate", Color.WHITE, 0.06)

## 受擊閃白（shader 方式，只影響 Sprite2D）
func _flash_hit(node: Node2D) -> void:
	# 找到 Sprite2D 子節點
	for child in node.get_children():
		if child is Sprite2D:
			if _cached_hit_flash_shader != null:
				var mat = child.get_meta("hit_flash_mat", null)
				if mat == null:
					mat = ShaderMaterial.new()
					mat.shader = _cached_hit_flash_shader
					child.material = mat
					child.set_meta("hit_flash_mat", mat)
				# 閃白動畫
				var tween = create_tween()
				tween.tween_method(func(v): mat.set_shader_parameter("flash_amount", v), 1.0, 0.0, 0.12)
			else:
				# 備用：modulate 閃白
				var tween = create_tween()
				tween.tween_property(child, "modulate", Color(3.0, 3.0, 3.0, 1.0), 0.04)
				tween.tween_property(child, "modulate", Color.WHITE, 0.08)
			break

## BOSS 事件處理（Phase 2 視覺變化）
func _on_boss_event(event_data: Dictionary) -> void:
	var event = event_data.get("event", "")

	# BOSS 登場：全畫面特效 + 強烈震動（Server 廣播 "spawn"）
	if event == "spawn" or event == "boss_enter":
		# 取得 BOSS 節點的實際位置（如果已生成）
		var boss_pos = Vector2(1100, 360)  # 預設右側進場位置
		for id in _target_nodes:
			var n = _target_nodes[id]
			if is_instance_valid(n) and n.get_meta("target_type", "") == "boss":
				boss_pos = n.position
				break
		HitEffect.spawn_boss_enter(boss_pos)
		ScreenShake.add_trauma(0.9)
		return

	if event != "phase_change":
		return

	var instance_id = event_data.get("instance_id", "")
	if not _target_nodes.has(instance_id):
		return

	var node = _target_nodes[instance_id]
	if not is_instance_valid(node):
		return

	# Phase 2：BOSS 變紅 + 閃爍 + 放大（規格書 28.2）
	var tween = create_tween()
	# 閃爍 3 次
	for i in 3:
		tween.tween_property(node, "modulate", Color(3.0, 0.3, 0.3, 1.0), 0.08)
		tween.tween_property(node, "modulate", Color.WHITE, 0.08)
	# 最終變成紅色調（Phase 2 持續視覺）
	tween.tween_property(node, "modulate", Color(1.5, 0.5, 0.5, 1.0), 0.1)

	# 放大 10%（Phase 2 更有威脅感）
	var tween2 = create_tween()
	tween2.tween_property(node, "scale", Vector2(2.2, 2.2), 0.3)

	# Phase 2 震動
	ScreenShake.add_trauma(0.6)

	# 切換到 Phase 2 動畫行
	_boss_anim_row = BOSS_ROW_PHASE2
	_boss_anim_frame = 0
	_boss_anim_timer = 0.0

	# 顯示 Phase 2 警告文字
	var phase_label = Label.new()
	phase_label.text = "PHASE 2!"
	phase_label.position = node.position + Vector2(-30, -80)
	phase_label.add_theme_font_size_override("font_size", 22)
	phase_label.modulate = Color(1.0, 0.2, 0.2)
	add_child(phase_label)
	var tween3 = create_tween()
	tween3.tween_property(phase_label, "scale", Vector2(1.5, 1.5), 0.2)
	tween3.tween_property(phase_label, "scale", Vector2(1.0, 1.0), 0.1)
	tween3.tween_interval(0.5)
	tween3.tween_property(phase_label, "modulate:a", 0.0, 0.4)
	tween3.tween_callback(func():
		if is_instance_valid(phase_label):
			phase_label.queue_free()
	)

func _process(delta: float) -> void:
	_update_target_positions(delta)
	_update_escape_warnings()
	_update_boss_animation(delta)
	_update_swim_animation(delta)

## 游泳動畫幀更新（全局計時器，所有目標物共用同一幀）
func _update_swim_animation(delta: float) -> void:
	if _swim_atlas_cache.is_empty():
		return

	_swim_anim_timer += delta
	if _swim_anim_timer < 1.0 / SWIM_ANIM_FPS:
		return

	_swim_anim_timer = 0.0
	_swim_anim_frame = 1 - _swim_anim_frame  # 在 0 和 1 之間切換

	# 更新所有有游泳動畫的目標物
	for instance_id in _target_nodes:
		var node = _target_nodes[instance_id]
		if not is_instance_valid(node):
			continue

		var def_id = node.get_meta("def_id", "")
		if not _swim_atlas_cache.has(def_id):
			continue

		# 找 Sprite2D 子節點並更新 texture
		for child in node.get_children():
			if child is Sprite2D:
				var frames = _swim_atlas_cache[def_id]
				if frames.size() > _swim_anim_frame:
					child.texture = frames[_swim_anim_frame]
				break

## BOSS 動畫幀更新（每幀切換 AtlasTexture）
func _update_boss_animation(delta: float) -> void:
	if _boss_sheet == null:
		return

	var fps = BOSS_ANIM_FPS
	if _boss_anim_row == BOSS_ROW_DEATH:
		fps = 8.0  # death 動畫快一倍

	_boss_anim_timer += delta
	if _boss_anim_timer >= 1.0 / fps:
		_boss_anim_timer = 0.0
		_boss_anim_frame = (_boss_anim_frame + 1) % BOSS_COLS

		# 更新所有 BOSS 節點的 Sprite2D texture
		for instance_id in _target_nodes:
			var node = _target_nodes[instance_id]
			if not is_instance_valid(node):
				continue
			if node.get_meta("target_type", "") != "boss":
				continue
			var sprite = node.get_node_or_null("Sprite2D") if node.get_child_count() > 0 else null
			# 找 Sprite2D 子節點
			for child in node.get_children():
				if child is Sprite2D:
					var new_tex = _get_boss_frame(_boss_anim_row, _boss_anim_frame)
					if new_tex != null:
						child.texture = new_tex
					break

## 目標生成
func _on_target_spawned(data: Dictionary) -> void:
	var instance_id = data.get("instance_id", "")
	if _target_nodes.has(instance_id):
		return

	# 建立目標節點（使用 TargetPool 重用節點，避免 GC 壓力）
	var node = _create_target_node(data)
	# 注意：TargetPool.acquire() 已在 _create_target_node 內呼叫，
	# 節點已加入場景（由 TargetPool.init_pool 時加入），不需要再 add_child
	_target_nodes[instance_id] = node

	# 進場動畫（scale 0 → 1，彈入效果）
	# 讓目標物出現更有存在感，而不是突然冒出來
	var target_type = data.get("type", "basic")
	var multiplier = data.get("multiplier", 0.0)
	node.scale = Vector2.ZERO
	var spawn_tween = node.create_tween()

	if target_type == "boss":
		# BOSS：慢速放大（0.4s），更有威壓感
		spawn_tween.tween_property(node, "scale", Vector2(1.0, 1.0), 0.4).set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_ELASTIC)
	elif multiplier >= 30.0:
		# 高倍率特殊目標（30x+）：彈入 + 輕微過衝（1.15x → 1.0x）
		spawn_tween.tween_property(node, "scale", Vector2(1.15, 1.15), 0.18).set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_BACK)
		spawn_tween.tween_property(node, "scale", Vector2(1.0, 1.0), 0.08)
		# 高倍率目標進場音效（coin_drop 短促提示，讓玩家注意到高價值目標出現）
		AudioManager.play_sfx(AudioManager.SFX.COIN_DROP)
	else:
		# 普通目標：快速彈入（0.12s）
		spawn_tween.tween_property(node, "scale", Vector2(1.0, 1.0), 0.12).set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_BACK)

## 目標擊破
func _on_target_killed(data: Dictionary) -> void:
	var instance_id = data.get("instance_id", "")
	if not _target_nodes.has(instance_id):
		return

	var node = _target_nodes[instance_id]
	# 先從字典移除，防止 _update_target_positions 繼續操作它
	_target_nodes.erase(instance_id)

	# 再播放特效（節點已從字典移除，不會被 update 干擾）
	if is_instance_valid(node):
		_play_kill_effect(node, data)
	# 注意：_play_kill_effect 內部會在動畫結束後呼叫 TargetPool.release(node)
	# 對於沒有動畫的情況，直接 release

## 建立目標節點（使用 TargetPool 重用節點，避免 GC 壓力）
func _create_target_node(data: Dictionary) -> Node2D:
	# 從 TargetPool 取出空殼節點（已在場景中，不需要 add_child）
	var container = TargetPool.acquire()
	container.position = Vector2(data.get("x", 0), data.get("y", 0))
	container.name = "Target_" + data.get("instance_id", "")

	var def_id = data.get("def_id", "T001")
	var target_type = data.get("type", "basic")

	# 使用像素 Sprite（優先用游泳動畫幀0，其次 Spritesheet AtlasTexture，減少 draw call）
	var sprite = Sprite2D.new()
	var tex: Texture2D = null

	# 優先使用游泳動畫 spritesheet 的幀0（有動畫效果）
	if _swim_atlas_cache.has(def_id):
		var frames = _swim_atlas_cache[def_id]
		if frames.size() > 0:
			tex = frames[0]

	# 備用：靜態 Spritesheet 或獨立 PNG
	if tex == null:
		tex = _get_target_texture(def_id)

	if tex != null:
		sprite.texture = tex
		# 像素完美縮放（關閉濾波）
		sprite.texture_filter = CanvasItem.TEXTURE_FILTER_NEAREST

		# 套用 outline shader（使用快取 shader，避免重複 load）
		if _cached_outline_shader != null and PerformanceMonitor.is_outline_shader_enabled():
			var mat = ShaderMaterial.new()
			mat.shader = _cached_outline_shader
			# 依目標類型設定輪廓顏色
			match target_type:
				"boss":
					mat.set_shader_parameter("outline_color", Color(1.0, 0.2, 0.2, 1.0))  # 紅色輪廓
					mat.set_shader_parameter("outline_width", 2.0)
				"special":
					mat.set_shader_parameter("outline_color", Color(1.0, 0.85, 0.0, 1.0))  # 金色輪廓
					mat.set_shader_parameter("outline_width", 1.5)
				_:
					mat.set_shader_parameter("outline_color", Color(0.0, 0.0, 0.0, 0.8))  # 黑色輪廓
					mat.set_shader_parameter("outline_width", 1.0)
			sprite.material = mat
			sprite.set_meta("outline_mat", mat)

		# 特殊目標物加 wobble tween（T103 流星、T104 金草）
		if def_id in ["T103", "T104"]:
			var tween = container.create_tween().set_loops()
			if def_id == "T103":
				# 流星：快速搖晃
				tween.tween_property(container, "rotation_degrees", 5.0, 0.15)
				tween.tween_property(container, "rotation_degrees", -5.0, 0.15)
			else:
				# 金草：緩慢搖晃
				tween.tween_property(container, "rotation_degrees", 3.0, 0.4)
				tween.tween_property(container, "rotation_degrees", -3.0, 0.4)
			TargetPool.register_tween(container, tween)  # 追蹤 tween，release 時自動 kill
	else:
		# 備用：ColorRect（texture 載入失敗時）
		var rect = ColorRect.new()
		rect.size = Vector2(32, 32) if target_type != "boss" else Vector2(64, 64)
		rect.position = -rect.size / 2
		rect.color = Color(0.8, 0.2, 0.8)
		container.add_child(rect)

	# BOSS 放大（B001 是 96x96，不需要額外放大）
	if target_type == "boss":
		sprite.scale = Vector2(1.5, 1.5)  # 輕微放大讓 BOSS 更有存在感

	container.add_child(sprite)

	# HP 條（規格書 8.1：受擊反饋）
	var hp_bar_bg = ColorRect.new()
	hp_bar_bg.size = Vector2(48, 5)
	hp_bar_bg.position = Vector2(-24, -38)
	hp_bar_bg.color = Color(0.2, 0.2, 0.2, 0.8)
	hp_bar_bg.name = "HPBarBG"
	container.add_child(hp_bar_bg)

	var hp_bar = ColorRect.new()
	hp_bar.size = Vector2(48, 5)
	hp_bar.position = Vector2(-24, -38)
	hp_bar.color = Color(0.2, 0.9, 0.2)
	hp_bar.name = "HPBar"
	container.add_child(hp_bar)

	# 儲存資料
	container.set_meta("instance_id", data.get("instance_id", ""))
	container.set_meta("def_id", def_id)
	container.set_meta("speed", data.get("speed", 0.0))
	container.set_meta("behavior", data.get("behavior", "linear"))
	container.set_meta("spawn_time", Time.get_ticks_msec())
	container.set_meta("target_type", target_type)

	# 倍率標籤（捕魚機標準 UX：讓玩家一眼看出目標價值）
	var multiplier = data.get("multiplier", 0.0)
	if multiplier > 0.0 and target_type != "boss":
		_add_multiplier_label(container, multiplier, target_type)

	# 高倍率目標光暈效果（30x+ 有金色光暈閃爍，50x 有更強烈的橙紅光暈）
	# 讓玩家一眼識別高價值目標，增加遊戲爽感
	if multiplier >= 30.0 and target_type != "boss":
		_add_high_value_glow(container, multiplier)

	# 游泳動畫：輕微上下搖擺 + 旋轉傾斜（讓目標物有生命感）
	# T103/T104 已有旋轉搖晃，不再加上下搖擺
	# BOSS 不加（有自己的移動邏輯）
	if target_type not in ["boss"] and def_id not in ["T103", "T104"] and PerformanceMonitor.is_swim_animation_enabled():
		var swim_amp = randf_range(3.0, 7.0)   # 搖擺幅度（像素）
		var swim_dur = randf_range(0.6, 1.2)   # 搖擺週期（秒）
		var swim_rot = randf_range(2.0, 5.0)   # 旋轉幅度（度）
		var swim_phase = randf_range(0.0, 1.0) # 隨機相位（避免所有魚同步）

		# Y 軸搖擺（主要游泳動作）
		var swim_tween = container.create_tween().set_loops()
		swim_tween.tween_property(sprite, "position:y", swim_amp, swim_dur * (0.5 + swim_phase * 0.1))
		swim_tween.tween_property(sprite, "position:y", -swim_amp, swim_dur * (0.5 + swim_phase * 0.1))
		TargetPool.register_tween(container, swim_tween)  # 追蹤 tween，release 時自動 kill

		# 旋轉傾斜（模擬魚游泳時身體傾斜）
		var rot_tween = container.create_tween().set_loops()
		rot_tween.tween_property(container, "rotation_degrees", swim_rot, swim_dur * 0.55)
		rot_tween.tween_property(container, "rotation_degrees", -swim_rot, swim_dur * 0.55)
		TargetPool.register_tween(container, rot_tween)  # 追蹤 tween

		# 縮放呼吸感（特殊目標更明顯）
		if target_type == "special":
			var scale_tween = container.create_tween().set_loops()
			scale_tween.tween_property(sprite, "scale", Vector2(1.05, 0.95), swim_dur * 0.5)
			scale_tween.tween_property(sprite, "scale", Vector2(0.95, 1.05), swim_dur * 0.5)
			TargetPool.register_tween(container, scale_tween)  # 追蹤 tween

	return container

## 更新目標位置（依行為模式移動）
func _update_target_positions(delta: float) -> void:
	# 先收集要移除的 ID，避免迭代中修改 Dictionary
	var to_remove: Array[String] = []

	for instance_id in _target_nodes:
		var node = _target_nodes[instance_id]
		if not is_instance_valid(node):
			to_remove.append(instance_id)
			continue

		var speed = node.get_meta("speed", 0.0)
		var behavior = node.get_meta("behavior", "linear")

		if speed > 0:
			var t = Time.get_ticks_msec() * 0.001
			match behavior:
				"linear":
					node.position.x -= speed * delta
				"curve":
					node.position.x -= speed * delta
					node.position.y += sin(t * 2.0) * 30 * delta
				"jump":
					node.position.x -= speed * delta * 0.5
					node.position.y += cos(t * 4.0) * 50 * delta
				"meteor":
					node.position.x -= speed * delta
					node.position.y += speed * 0.3 * delta
				"sway":
					node.position.x -= speed * 0.3 * delta
					node.position.x += sin(t) * 20 * delta
				"static_sway":
					node.position.x += sin(t * 2.0) * 5 * delta
				"sink":
					node.position.y += speed * 0.4 * delta
				"flee":
					var flee_spd = node.get_meta("flee_speed", speed)
					node.position.x -= flee_spd * delta
				"coin_rain":
					node.position.x -= speed * delta
				"mimic":
					node.position.x -= speed * delta
					node.position.y += sin(t * 1.5) * 15 * delta
				"boss_phases":
					# BOSS 左右移動
					node.position.x += sin(t * 0.5) * speed * delta

		# 離開畫面標記移除
		if node.position.x < -150 or node.position.x > 1450:
			to_remove.append(instance_id)
		else:
			# 可見性剔除：畫面外的目標物設 visible=false，減少 draw call
			# 畫面範圍：x 0~1280，y 0~720（加 64px 緩衝避免閃爍）
			var in_screen = (node.position.x > -64 and node.position.x < 1344 and
							 node.position.y > -64 and node.position.y < 784)
			if node.visible != in_screen:
				node.visible = in_screen
			to_remove.append(instance_id)

	# 統一移除
	for id in to_remove:
		if _target_nodes.has(id):
			var node = _target_nodes[id]
			if is_instance_valid(node):
				TargetPool.release(node)  # 歸還到 pool，不 queue_free
			_target_nodes.erase(id)

## 擊破特效（規格書 8.2）
func _play_kill_effect(node: Node2D, data: Dictionary) -> void:
	if not is_instance_valid(node):
		return

	# 先記錄位置，再做動畫
	var kill_pos = node.position
	var reward = data.get("reward", 0)
	var multiplier = data.get("multiplier", 1.0)
	var def_id = data.get("def_id", "")

	# T101 擬態型怪物：死亡時先變形回原形（規格書 26.2）
	if def_id == "T101":
		_play_mimic_death(node, kill_pos, reward, multiplier)
		return

	# T105 巨大金幣魚：擊破後金幣雨（規格書 26.2）
	if def_id == "T105":
		_spawn_coin_rain(kill_pos)

	# B001 BOSS：播放 death 動畫後消失
	if def_id == "B001":
		_play_boss_death(node, kill_pos, reward, multiplier)
		return

	# 閃白所有子節點
	for child in node.get_children():
		if child is Sprite2D or child is ColorRect:
			child.modulate = Color.WHITE * 3.0

	# 縮放爆炸動畫
	var tween = create_tween()
	tween.tween_property(node, "scale", Vector2(1.5, 1.5), 0.08)
	tween.tween_property(node, "scale", Vector2(0.0, 0.0), 0.15)
	tween.tween_callback(func():
		TargetPool.release(node)  # 歸還到 pool，不 queue_free
	)

	# 獎勵跳字（用記錄的位置，不依賴 node）
	if reward > 0:
		_spawn_reward_text(kill_pos, reward, multiplier)

	# 使用 HitEffect 系統（取代舊的 _spawn_death_particles）
	HitEffect.spawn_kill(kill_pos, multiplier)

	# 震動（依倍率）
	var trauma = clamp(0.2 + multiplier * 0.005, 0.2, 0.6)
	ScreenShake.add_trauma(trauma)

## B001 BOSS 死亡動畫（播放 death 幀後消失）
func _play_boss_death(node: Node2D, kill_pos: Vector2, reward: int, multiplier: float) -> void:
	if not is_instance_valid(node):
		return

	# 切換到 death 動畫行
	_boss_anim_row = BOSS_ROW_DEATH
	_boss_anim_frame = 0
	_boss_anim_timer = 0.0

	# 播放 death 動畫（4幀 × 8fps = 0.5秒）
	var death_duration = float(BOSS_COLS) / 8.0  # 0.5 秒

	# 縮放爆炸動畫（配合 death 幀）
	var tween = create_tween()
	tween.tween_property(node, "scale", Vector2(2.5, 2.5), death_duration * 0.3)
	tween.tween_property(node, "scale", Vector2(0.0, 0.0), death_duration * 0.7)
	tween.tween_callback(func():
		TargetPool.release(node)  # 歸還到 pool
		# 重置 BOSS 動畫狀態
		_boss_anim_row = BOSS_ROW_IDLE
		_boss_anim_frame = 0
	)

	# 獎勵跳字
	if reward > 0:
		_spawn_reward_text(kill_pos, reward, multiplier)

	# BOSS 死亡特效（大爆炸）
	HitEffect.spawn_kill(kill_pos, multiplier)
	ScreenShake.add_trauma(1.0)

## T101 擬態型怪物死亡變形（規格書 26.2）
func _play_mimic_death(node: Node2D, kill_pos: Vector2, reward: int, multiplier: float) -> void:
	if not is_instance_valid(node):
		return

	# 第一階段：閃爍（偽裝破碎）
	var tween1 = create_tween()
	tween1.tween_property(node, "modulate", Color(2.0, 0.5, 0.5, 1.0), 0.08)
	tween1.tween_property(node, "modulate", Color.WHITE, 0.08)
	tween1.tween_property(node, "modulate", Color(2.0, 0.5, 0.5, 1.0), 0.08)
	tween1.tween_property(node, "modulate", Color.WHITE, 0.08)

	# 第二階段：縮小再放大（變形）
	var tween2 = create_tween()
	tween2.tween_property(node, "scale", Vector2(0.3, 1.5), 0.15)
	tween2.tween_property(node, "scale", Vector2(1.5, 0.3), 0.15)

	# 第三階段：爆炸消失
	var tween3 = create_tween()
	tween3.tween_interval(0.35)
	tween3.tween_property(node, "scale", Vector2(2.0, 2.0), 0.1)
	tween3.parallel().tween_property(node, "modulate:a", 0.0, 0.1)
	tween3.tween_callback(func():
		TargetPool.release(node)  # 歸還到 pool
	)

	# 生成「真面目」文字
	var reveal_label = Label.new()
	reveal_label.text = "正體！"
	reveal_label.position = kill_pos + Vector2(-20, -40)
	reveal_label.add_theme_font_size_override("font_size", 18)
	reveal_label.modulate = Color(1.0, 0.3, 0.3)
	add_child(reveal_label)
	var tween4 = create_tween()
	tween4.tween_property(reveal_label, "position:y", kill_pos.y - 90, 0.8)
	tween4.parallel().tween_property(reveal_label, "modulate:a", 0.0, 0.8)
	tween4.tween_callback(func():
		if is_instance_valid(reveal_label):
			reveal_label.queue_free()
	)

	if reward > 0:
		_spawn_reward_text(kill_pos, reward, multiplier)
	_spawn_death_particles(kill_pos)

## T105 金幣魚擊破後金幣雨（規格書 26.2）
func _spawn_coin_rain(origin: Vector2) -> void:
	# 生成 18 枚金幣從擊破位置散落（升級版：用 Node2D + _draw 繪製真實金幣）
	for i in 18:
		var coin = _create_pixel_coin()
		coin.position = origin + Vector2(randf_range(-20, 20), randf_range(-10, 10))
		add_child(coin)

		# 拋物線軌跡（更自然的弧線）
		var target_x = origin.x + randf_range(-150, 150)
		var target_y = origin.y + randf_range(100, 220)
		var peak_y = origin.y - randf_range(70, 140)
		var mid_x = origin.x + (target_x - origin.x) * 0.5

		var tween = coin.create_tween()
		# 上升（帶旋轉）
		tween.tween_property(coin, "position", Vector2(mid_x, peak_y), 0.22).set_ease(Tween.EASE_OUT)
		tween.parallel().tween_property(coin, "rotation_degrees", randf_range(90, 270), 0.22)
		# 下落（加速）
		tween.tween_property(coin, "position", Vector2(target_x, target_y), 0.32).set_ease(Tween.EASE_IN)
		tween.parallel().tween_property(coin, "rotation_degrees", randf_range(360, 540), 0.32)
		tween.parallel().tween_property(coin, "modulate:a", 0.0, 0.32)
		tween.tween_callback(func():
			if is_instance_valid(coin):
				coin.queue_free()
		)

	# 播放金幣音效
	AudioManager.play_sfx(AudioManager.SFX.COIN_DROP)

## 建立像素金幣節點（靜態腳本，效能比動態 GDScript 好）
const PixelCoinScript = preload("res://scripts/effects/PixelCoin.gd")

func _create_pixel_coin() -> Node2D:
	var coin = PixelCoinScript.new()
	coin.z_index = 12
	return coin

## 生成死亡粒子
func _spawn_death_particles(pos: Vector2) -> void:
	for i in 6:
		var particle = ColorRect.new()
		particle.size = Vector2(6, 6)
		particle.color = [Color.GOLD, Color.YELLOW, Color.WHITE, Color(1,0.5,0)][i % 4]
		particle.position = pos + Vector2(randf_range(-10, 10), randf_range(-10, 10))
		add_child(particle)
		var tween = create_tween()
		var target_pos = pos + Vector2(randf_range(-40, 40), randf_range(-60, -10))
		tween.tween_property(particle, "position", target_pos, 0.4)
		tween.parallel().tween_property(particle, "modulate:a", 0.0, 0.4)
		tween.tween_callback(func():
			if is_instance_valid(particle):
				particle.queue_free()
		)

## 生成獎勵跳字
func _spawn_reward_text(pos: Vector2, amount: int, multiplier: float) -> void:
	var label = Label.new()
	label.text = "+%d" % amount
	label.position = pos
	label.add_theme_font_size_override("font_size", 16)

	# 套用像素字體（使用快取，避免重複 load）
	if _cached_pixel_font != null:
		label.add_theme_font_override("font", _cached_pixel_font)

	# 依倍率設定顏色（規格書 8.3）
	if multiplier >= 100:
		label.modulate = Color(1.0, 0.2, 0.2)
	elif multiplier >= 20:
		label.modulate = Color(1.0, 0.8, 0.0)
	else:
		label.modulate = Color.WHITE

	add_child(label)

	var tween = create_tween()
	tween.tween_property(label, "position:y", pos.y - 80, 1.0)
	tween.parallel().tween_property(label, "modulate:a", 0.0, 1.0)
	tween.tween_callback(func():
		if is_instance_valid(label):
			label.queue_free()
	)

## 點擊目標（玩家點擊畫面時呼叫）
func try_click_target(click_pos: Vector2) -> String:
	var closest_id = ""
	var closest_dist = 70.0  # 點擊判定範圍（像素）

	for instance_id in _target_nodes:
		var node = _target_nodes[instance_id]
		if not is_instance_valid(node):
			continue
		var dist = node.position.distance_to(click_pos)
		if dist < closest_dist:
			closest_dist = dist
			closest_id = instance_id

	return closest_id

## 顯示 Lock 視覺框
func show_lock_indicator(instance_id: String) -> void:
	# 先清除舊的 lock 框
	for id in _target_nodes:
		var node = _target_nodes[id]
		if is_instance_valid(node):
			var old_lock = node.get_node_or_null("LockFrame")
			if old_lock:
				old_lock.queue_free()

	if instance_id == "":
		return

	if not _target_nodes.has(instance_id):
		return

	var node = _target_nodes[instance_id]
	if not is_instance_valid(node):
		return

	# 建立像素準星框
	var frame = Node2D.new()
	frame.name = "LockFrame"
	node.add_child(frame)

	var draw = func():
		pass  # 用 ColorRect 模擬準星

	# 四個角落的 L 形準星
	var size = 20
	var corners = [
		Vector2(-size, -size), Vector2(size, -size),
		Vector2(-size, size),  Vector2(size, size)
	]
	for corner in corners:
		var h = ColorRect.new()
		h.size = Vector2(8, 2)
		h.color = Color(1.0, 0.9, 0.0, 0.9)
		h.position = corner + Vector2(-4 if corner.x < 0 else -4, -1)
		frame.add_child(h)

		var v = ColorRect.new()
		v.size = Vector2(2, 8)
		v.color = Color(1.0, 0.9, 0.0, 0.9)
		v.position = corner + Vector2(-1, -4 if corner.y < 0 else -4)
		frame.add_child(v)

	# 閃爍動畫（綁定到 frame 節點，節點刪除時自動停止）
	var tween = frame.create_tween().set_loops()
	tween.tween_property(frame, "modulate:a", 0.4, 0.4)
	tween.tween_property(frame, "modulate:a", 1.0, 0.4)

## 更新目標 HP 條
func update_target_hp(instance_id: String, hp: int, max_hp: int) -> void:
	if not _target_nodes.has(instance_id):
		return
	var node = _target_nodes[instance_id]
	var hp_bar = node.get_node_or_null("HPBar")
	if hp_bar and max_hp > 0:
		var new_width = 48.0 * (float(hp) / float(max_hp))
		hp_bar.size.x = new_width
		# HP 條顏色依血量變化（高 → 綠，中 → 黃，低 → 紅）
		var pct = float(hp) / float(max_hp)
		if pct > 0.6:
			hp_bar.color = Color(0.2, 0.9, 0.2)
		elif pct > 0.3:
			hp_bar.color = Color(1.0, 0.8, 0.1)
		else:
			hp_bar.color = Color(1.0, 0.2, 0.2)
		# 受擊閃爍（HP 條短暫變白再回原色）
		var orig_color = hp_bar.color
		var tween = create_tween()
		tween.tween_property(hp_bar, "color", Color.WHITE, 0.04)
		tween.tween_property(hp_bar, "color", orig_color, 0.08)
		# 低血量脈動效果（HP < 30%）：啟動或停止脈動 tween
		var pulse_tween_key = "hp_pulse_tween"
		if pct <= 0.3:
			# 若尚未有脈動 tween，建立一個
			if not node.has_meta(pulse_tween_key):
				var pulse = create_tween().set_loops()
				pulse.tween_property(hp_bar, "modulate:a", 0.4, 0.25).set_ease(Tween.EASE_IN_OUT)
				pulse.tween_property(hp_bar, "modulate:a", 1.0, 0.25).set_ease(Tween.EASE_IN_OUT)
				node.set_meta(pulse_tween_key, pulse)
		else:
			# HP 回到 30% 以上，停止脈動
			if node.has_meta(pulse_tween_key):
				var pulse = node.get_meta(pulse_tween_key)
				if pulse is Tween:
					pulse.kill()
				node.remove_meta(pulse_tween_key)
				hp_bar.modulate.a = 1.0

## 目標物接近左邊緣時顯示逃跑警告箭頭
func _update_escape_warnings() -> void:
	for instance_id in _target_nodes:
		var node = _target_nodes[instance_id]
		if not is_instance_valid(node):
			continue

		var speed = node.get_meta("speed", 0.0)
		if speed <= 0:
			continue  # 靜止目標不需要警告

		var x = node.position.x
		var warning = node.get_node_or_null("EscapeWarning")

		# 目標物 x < 120 且有速度時顯示警告
		if x < 120 and x > 0:
			if warning == null:
				# 建立警告箭頭
				warning = Label.new()
				warning.name = "EscapeWarning"
				warning.text = "◀!"
				warning.position = Vector2(-30, -20)
				warning.add_theme_font_size_override("font_size", 14)
				warning.modulate = Color(1.0, 0.3, 0.3, 0.0)
				node.add_child(warning)

			# 透明度依距離邊緣的遠近（越近越明顯）
			var alpha = clamp((120.0 - x) / 100.0, 0.0, 1.0)
			warning.modulate.a = alpha

			# 閃爍（最後 60px 快速閃爍）
			if x < 60:
				var flash = int(Time.get_ticks_msec() / 150) % 2 == 0
				warning.modulate.a = alpha if flash else alpha * 0.3
		else:
			# 離開警告區域，移除警告
			if warning != null:
				warning.queue_free()

## 加入倍率標籤（捕魚機標準 UX）
## 在目標物上方顯示倍率，讓玩家一眼看出目標價值
func _add_multiplier_label(container: Node2D, multiplier: float, target_type: String) -> void:
	var label = Label.new()
	label.name = "MultiplierLabel"

	# 格式化倍率文字
	if multiplier == int(multiplier):
		label.text = "%dx" % int(multiplier)
	else:
		label.text = "%.0fx" % multiplier

	# 位置：目標物上方
	label.position = Vector2(-20, -52)
	label.size = Vector2(40, 16)
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER

	# 字體大小（依倍率大小調整）
	var font_size: int
	if multiplier >= 50:
		font_size = 14
	elif multiplier >= 20:
		font_size = 13
	else:
		font_size = 11
	label.add_theme_font_size_override("font_size", font_size)

	# 套用像素字體
	if _cached_pixel_font != null:
		label.add_theme_font_override("font", _cached_pixel_font)

	# 顏色（依倍率和類型）
	var color: Color
	if target_type == "special":
		if multiplier >= 50:
			color = Color(1.0, 0.3, 0.1)   # 橙紅（最高倍率）
		elif multiplier >= 30:
			color = Color(1.0, 0.85, 0.0)  # 金色
		else:
			color = Color(0.4, 1.0, 0.8)   # 青色
	else:
		if multiplier >= 10:
			color = Color(1.0, 0.9, 0.3)   # 黃色
		elif multiplier >= 5:
			color = Color(0.8, 1.0, 0.5)   # 淡綠
		else:
			color = Color(0.9, 0.9, 0.9)   # 白灰

	label.modulate = color

	# 文字陰影（增加可讀性）
	label.add_theme_color_override("font_shadow_color", Color(0.0, 0.0, 0.0, 0.8))
	label.add_theme_constant_override("shadow_offset_x", 1)
	label.add_theme_constant_override("shadow_offset_y", 1)

	container.add_child(label)

	# 輕微上下浮動動畫（讓標籤更有生命感）
	var float_tween = label.create_tween().set_loops()
	float_tween.tween_property(label, "position:y", -56.0, 0.8).set_ease(Tween.EASE_IN_OUT).set_trans(Tween.TRANS_SINE)
	float_tween.tween_property(label, "position:y", -52.0, 0.8).set_ease(Tween.EASE_IN_OUT).set_trans(Tween.TRANS_SINE)

## 高倍率目標光暈效果（30x+ 金色光暈，50x 橙紅光暈）
## 用 ColorRect 模擬光暈，不需要額外 shader
func _add_high_value_glow(container: Node2D, multiplier: float) -> void:
	# 光暈顏色和強度依倍率決定
	var glow_color: Color
	var glow_size: float
	var pulse_speed: float

	if multiplier >= 50.0:
		# 50x：橙紅光暈，強烈閃爍（T105 金幣魚、T103 流星最高倍率）
		glow_color = Color(1.0, 0.5, 0.1, 0.35)
		glow_size = 52.0
		pulse_speed = 0.4
	elif multiplier >= 30.0:
		# 30x：金色光暈，中等閃爍（T104 金草、T102 寶箱怪）
		glow_color = Color(1.0, 0.85, 0.0, 0.25)
		glow_size = 44.0
		pulse_speed = 0.6

	# 建立光暈 ColorRect（圓形用大正方形 + 透明度模擬）
	var glow = ColorRect.new()
	glow.name = "HighValueGlow"
	glow.size = Vector2(glow_size, glow_size)
	glow.position = Vector2(-glow_size / 2.0, -glow_size / 2.0)
	glow.color = glow_color
	glow.z_index = -1  # 在 Sprite 後面
	container.add_child(glow)
	container.move_child(glow, 0)  # 移到最底層

	# 脈動閃爍動畫（讓光暈有呼吸感）
	var glow_tween = glow.create_tween().set_loops()
	glow_tween.tween_property(glow, "modulate:a", 0.3, pulse_speed).set_ease(Tween.EASE_IN_OUT).set_trans(Tween.TRANS_SINE)
	glow_tween.tween_property(glow, "modulate:a", 1.0, pulse_speed).set_ease(Tween.EASE_IN_OUT).set_trans(Tween.TRANS_SINE)

	# 50x 額外：縮放脈動（更強烈的視覺衝擊）
	if multiplier >= 50.0:
		var scale_tween = glow.create_tween().set_loops()
		scale_tween.tween_property(glow, "scale", Vector2(1.15, 1.15), pulse_speed * 0.8).set_ease(Tween.EASE_IN_OUT)
		scale_tween.tween_property(glow, "scale", Vector2(0.9, 0.9), pulse_speed * 0.8).set_ease(Tween.EASE_IN_OUT)
