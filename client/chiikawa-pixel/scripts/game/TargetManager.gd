## TargetManager.gd — 目標物系統
## target-system-agent 負責維護
extends Node2D

# 目標物 Sprite 路徑
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
	# DAY-292 新增特殊目標
	"T106": "res://assets/sprites/targets/T106_chain_lightning.png",
	"T107": "res://assets/sprites/targets/T107_crab_torpedo.png",
	"T108": "res://assets/sprites/targets/T108_vortex_anemone.png",
	"T109": "res://assets/sprites/targets/T109_golden_dragon.png",
	"T110": "res://assets/sprites/targets/T110_thunder_lobster.png",
	# DAY-293 新增特殊目標
	"T111": "res://assets/sprites/targets/T111_awakened_phoenix.png",
	"T112": "res://assets/sprites/targets/T112_shockwave_bomb.png",
	# DAY-294 新增特殊目標
	"T113": "res://assets/sprites/targets/T113_drill_torpedo.png",
	"T114": "res://assets/sprites/targets/T114_time_freeze.png",
	"T115": "res://assets/sprites/targets/T115_chain_explosion.png",
	# DAY-295 新增特殊目標
	"T116": "res://assets/sprites/targets/T116_chain_long_king.png",
	"T117": "res://assets/sprites/targets/T117_dragon_shotgun.png",
	"T118": "res://assets/sprites/targets/T118_rocket_cannon.png",
	"T119": "res://assets/sprites/targets/T119_deep_whirlpool.png",
	"T120": "res://assets/sprites/targets/T120_vampire_mult.png",
	# DAY-296 新增特殊目標
	"T121": "res://assets/sprites/targets/T121_mirror_fish.png",
	"T122": "res://assets/sprites/targets/T122_golden_rain.png",
	"T123": "res://assets/sprites/targets/T123_freeze_bomb.png",
	"T124": "res://assets/sprites/targets/T124_thunder_storm.png",
	"T125": "res://assets/sprites/targets/T125_lucky_wheel.png",
	# DAY-301 新增特殊目標
	"T126": "res://assets/sprites/targets/T126_jackpot_fish.png",
	"T127": "res://assets/sprites/targets/T127_coop_fish.png",
	"T128": "res://assets/sprites/targets/T128_time_warp.png",
	# DAY-302 新增特殊目標
	"T129": "res://assets/sprites/targets/T129_chain_meteor.png",
	# DAY-303 新增特殊目標
	"T130": "res://assets/sprites/targets/T130_crash_fish.png",
	# DAY-304 新增特殊目標
	"T131": "res://assets/sprites/targets/T131_electric_eel.png",
	"T132": "res://assets/sprites/targets/T132_angler_fish.png",
	"T133": "res://assets/sprites/targets/T133_black_hole.png",
	"T134": "res://assets/sprites/targets/T134_bounty_hunter.png",
	"T135": "res://assets/sprites/targets/T135_tsunami.png",
	# DAY-305 新增特殊目標
	"T136": "res://assets/sprites/targets/T136_dragon_wrath_v2.png",
	"T137": "res://assets/sprites/targets/T137_humpback_whale.png",
	"T138": "res://assets/sprites/targets/T138_legend_dragon.png",
	"T139": "res://assets/sprites/targets/T139_guild_war.png",
	"T140": "res://assets/sprites/targets/T140_quality_fish.png",
	# DAY-306 新增特殊目標
	"T141": "res://assets/sprites/targets/T141_tornado.png",
	"T142": "res://assets/sprites/targets/T142_earthquake.png",
	"T143": "res://assets/sprites/targets/T143_volcano.png",
	"T144": "res://assets/sprites/targets/T144_cosmic_ray.png",
	"T145": "res://assets/sprites/targets/T145_divine_dragon.png",
	# DAY-307 新增特殊目標
	"T146": "res://assets/sprites/targets/T146_quantum.png",
	"T147": "res://assets/sprites/targets/T147_supernova.png",
	"T148": "res://assets/sprites/targets/T148_infinite.png",
	"T149": "res://assets/sprites/targets/T149_genesis.png",
	"T150": "res://assets/sprites/targets/T150_rebirth.png",
	# DAY-308 新增特殊目標
	"T151": "res://assets/sprites/targets/T151_awakened_croc.png",
	"T152": "res://assets/sprites/targets/T152_vampire_v2.png",
	"T153": "res://assets/sprites/targets/T153_super_awaken.png",
	"T154": "res://assets/sprites/targets/T154_giant_prize.png",
	"T155": "res://assets/sprites/targets/T155_immortal_boss.png",
	# DAY-309 新增
	"T156": "res://assets/sprites/targets/T156_ice_phoenix.png",
	"T157": "res://assets/sprites/targets/T157_dragon_fury.png",
	"T158": "res://assets/sprites/targets/T158_mult_cascade.png",
	"T159": "res://assets/sprites/targets/T159_awaken_boss_v2.png",
	"T160": "res://assets/sprites/targets/T160_ultimate_judgment.png",
	# DAY-310 新增
	"T161": "res://assets/sprites/targets/T161_combo_burst.png",
	"T162": "res://assets/sprites/targets/T162_time_bomb.png",
	"T163": "res://assets/sprites/targets/T163_elemental_fusion.png",
	"T164": "res://assets/sprites/targets/T164_treasure_hunter.png",
	"T165": "res://assets/sprites/targets/T165_myth_awaken.png",
}

# 目標物顏色（無 Sprite 時的備用顏色）
const TARGET_COLORS = {
	"T001": Color(0.2, 0.8, 0.2),   # 綠色雜草
	"T002": Color(0.3, 0.9, 0.3),   # 綠色小蟲
	"T003": Color(0.9, 0.3, 0.3),   # 紅色小蟲
	"T004": Color(0.3, 0.5, 0.9),   # 藍色小蟲
	"T005": Color(1.0, 0.9, 0.4),   # 黃色布丁
	"T006": Color(0.6, 0.4, 0.2),   # 棕色蘑菇
	"T101": Color(0.6, 0.6, 0.6),   # 灰色擬態
	"T102": Color(0.9, 0.7, 0.2),   # 金色寶箱
	"T103": Color(1.0, 1.0, 0.8),   # 白色流星
	"T104": Color(1.0, 0.85, 0.0),  # 金色雜草
	"T105": Color(1.0, 0.8, 0.2),   # 金色魚
	"B001": Color(0.8, 0.2, 0.2),   # 紅色 BOSS
	# DAY-292 新增特殊目標備用顏色
	"T106": Color(0.0, 0.9, 1.0),   # 青藍連鎖閃電
	"T107": Color(1.0, 0.4, 0.1),   # 橙紅螃蟹魚雷
	"T108": Color(0.5, 0.2, 0.8),   # 紫色渦旋海葵
	"T109": Color(1.0, 0.85, 0.0),  # 金色黃金龍魚
	"T110": Color(1.0, 0.3, 0.0),   # 火紅雷霆龍蝦
	# DAY-293 新增特殊目標備用顏色
	"T111": Color(1.0, 0.42, 0.21), # 火橙覺醒鳳凰
	"T112": Color(1.0, 0.27, 0.0),  # 深橙全場震盪
	# DAY-294 新增特殊目標備用顏色
	"T113": Color(1.0, 0.55, 0.15), # 橙色鑽頭魚雷
	"T114": Color(0.4, 0.85, 1.0),  # 冰藍時間凍結
	"T115": Color(0.9, 0.2, 0.15),  # 深紅連鎖爆炸
	# DAY-295 新增特殊目標備用顏色
	"T116": Color(1.0, 0.85, 0.0),  # 金色千龍王輪盤
	"T117": Color(0.8, 0.2, 0.9),   # 紫色龍力散彈
	"T118": Color(1.0, 0.3, 0.1),   # 火紅火箭砲
	"T119": Color(0.0, 0.6, 0.9),   # 深藍深海漩渦
	"T120": Color(0.5, 0.0, 0.5),   # 深紫吸血鬼
	# DAY-296 新增特殊目標備用顏色
	"T121": Color(0.88, 0.67, 1.0), # 淡紫鏡像魚
	"T122": Color(1.0, 0.85, 0.0),  # 金色黃金雨魚
	"T123": Color(0.0, 0.9, 1.0),   # 冰藍冰凍炸彈魚
	"T124": Color(1.0, 0.9, 0.2),   # 黃色雷暴魚
	"T125": Color(1.0, 0.42, 0.71), # 粉紅大轉盤魚
	# DAY-301 新增特殊目標備用顏色
	"T126": Color(1.0, 0.85, 0.0),  # 金色進階 Jackpot 魚
	"T127": Color(0.0, 0.9, 1.0),   # 青藍全服合作魚
	"T128": Color(0.55, 0.2, 0.86), # 紫色時間扭曲魚
	# DAY-302 新增特殊目標備用顏色
	"T129": Color(0.9, 0.4, 0.1),   # 火橙連鎖隕石魚
	# DAY-303 新增特殊目標備用顏色
	"T130": Color(0.8, 0.1, 0.1),   # 深紅崩潰魚
	# DAY-304 新增特殊目標備用顏色
	"T131": Color(1.0, 0.95, 0.0),  # 電黃電鰻魚
	"T132": Color(0.0, 0.9, 1.0),   # 青藍巨型安康魚
	"T133": Color(0.3, 0.0, 0.5),   # 深紫黑洞魚
	"T134": Color(1.0, 0.5, 0.0),   # 火橙賞金獵人魚
	"T135": Color(0.0, 0.5, 1.0),   # 深藍海嘯魚
	# DAY-305 新增特殊目標備用顏色
	"T136": Color(1.0, 0.27, 0.0),  # 火橙龍怒蓄積魚
	"T137": Color(0.0, 0.4, 0.8),   # 深藍座頭鯨魚
	"T138": Color(1.0, 0.5, 0.0),   # 金橙傳說龍魚
	"T139": Color(1.0, 0.85, 0.0),  # 金色公會戰魚
	"T140": Color(0.6, 0.0, 1.0),   # 紫色品質魚
	# DAY-306 新增特殊目標備用顏色
	"T141": Color(0.0, 0.9, 0.7),   # 青綠龍捲風魚
	"T142": Color(0.8, 0.4, 0.1),   # 棕橙地震魚
	"T143": Color(1.0, 0.27, 0.0),  # 火紅火山魚
	"T144": Color(0.6, 0.2, 1.0),   # 紫色星際魚
	"T145": Color(1.0, 0.85, 0.0),  # 金色神龍魚
	# DAY-307 新增特殊目標備用顏色
	"T146": Color(0.0, 0.9, 1.0),   # 青藍量子魚
	"T147": Color(1.0, 0.4, 0.1),   # 火橙超新星魚
	"T148": Color(0.6, 0.2, 1.0),   # 紫色無限魚
	"T149": Color(1.0, 0.85, 0.0),  # 金色創世魚
	"T150": Color(1.0, 0.27, 0.0),  # 火橙重生魚
	# DAY-308 新增特殊目標備用顏色
	"T151": Color(0.0, 0.8, 0.3),   # 深綠覺醒鱷魚
	"T152": Color(0.6, 0.0, 0.8),   # 深紫吸血鬼升級魚
	"T153": Color(1.0, 0.4, 0.0),   # 火橙超級覺醒魚
	"T154": Color(1.0, 0.85, 0.0),  # 金色巨型獎勵魚
	"T155": Color(0.9, 0.1, 0.1),   # 深紅不死 BOSS 魚
	# DAY-309 新增
	"T156": Color(0.0, 0.8, 1.0),   # 冰藍冰鳳凰魚
	"T157": Color(1.0, 0.3, 0.0),   # 火橙龍怒能量魚
	"T158": Color(0.1, 0.4, 1.0),   # 深藍倍率瀑布魚
	"T159": Color(1.0, 0.6, 0.0),   # 金橙覺醒 BOSS v2 魚
	"T160": Color(0.8, 0.0, 0.0),   # 深紅終極審判魚
	# DAY-310 新增
	"T161": Color(1.0, 0.4, 0.1),   # 火橙連擊爆發魚
	"T162": Color(0.9, 0.2, 0.1),   # 深紅時間炸彈魚
	"T163": Color(0.6, 0.2, 0.9),   # 深紫元素融合魚
	"T164": Color(1.0, 0.75, 0.1),  # 金色寶藏獵人魚
	"T165": Color(0.9, 0.8, 0.1),   # 神聖金色神話覺醒魚
}

var _target_nodes: Dictionary = {}  # instance_id -> Node2D
var _cached_textures: Dictionary = {}
var _time_warp_speed_mult: float = 1.0  # 時間扭曲速度倍率（DAY-301）

func _ready() -> void:
	GameManager.target_spawned.connect(_on_target_spawned)
	GameManager.target_updated.connect(_on_target_updated)
	GameManager.target_killed.connect(_on_target_killed)
	GameManager.boss_event.connect(_on_boss_event)
	# DAY-301 時間扭曲訊號
	GameManager.lucky_time_warp.connect(_on_lucky_time_warp_for_speed)

func _process(delta: float) -> void:
	_update_positions(delta)

func _update_positions(delta: float) -> void:
	for instance_id in _target_nodes:
		var node = _target_nodes[instance_id]
		if not is_instance_valid(node):
			continue
		var behavior = node.get_meta("behavior", "linear")
		var speed = node.get_meta("speed", 0.0) * _time_warp_speed_mult
		if node.get_meta("is_fleeing", false):
			speed = node.get_meta("flee_speed", speed * 2.5)

		match behavior:
			"linear", "flee", "fast":
				node.position.x -= speed * delta
			"sink":
				node.position.y += speed * 0.3 * delta
				node.position.x -= 10 * delta

		# 離開畫面則移除
		if node.position.x < -100:
			_target_nodes.erase(instance_id)
			node.queue_free()

# ── 目標物生成 ────────────────────────────────────────────────

func _on_target_spawned(data: Dictionary) -> void:
	var instance_id = data.get("instance_id", "")
	if _target_nodes.has(instance_id):
		return

	var node = _create_target_node(data)
	_target_nodes[instance_id] = node

	# 進場動畫
	node.scale = Vector2.ZERO
	var tween = node.create_tween()
	var target_type = data.get("type", "basic")
	if target_type == "boss":
		tween.tween_property(node, "scale", Vector2(1.0, 1.0), 0.4).set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_ELASTIC)
	else:
		tween.tween_property(node, "scale", Vector2(1.0, 1.0), 0.15).set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_BACK)

func _create_target_node(data: Dictionary) -> Node2D:
	var container = Node2D.new()
	container.position = Vector2(data.get("x", 0), data.get("y", 0))
	container.name = "Target_" + data.get("instance_id", "")
	add_child(container)

	var def_id = data.get("def_id", "T001")
	var target_type = data.get("type", "basic")
	var multiplier = data.get("multiplier", 2.0)

	# Sprite
	var sprite = Sprite2D.new()
	var tex = _get_texture(def_id)
	if tex != null:
		sprite.texture = tex
		sprite.texture_filter = CanvasItem.TEXTURE_FILTER_NEAREST
		# 目標物放大：基礎目標 2.5x，特殊目標 2.8x，確保在畫面上清楚可見
		if target_type == "boss":
			sprite.scale = Vector2(3.0, 3.0)
		elif target_type == "special":
			sprite.scale = Vector2(2.8, 2.8)
		else:
			sprite.scale = Vector2(2.5, 2.5)
	else:
		# 備用：ColorRect
		var rect = ColorRect.new()
		var size = 64.0 if target_type != "boss" else 128.0
		rect.size = Vector2(size, size)
		rect.position = -rect.size / 2
		rect.color = TARGET_COLORS.get(def_id, Color(0.8, 0.2, 0.8))
		container.add_child(rect)
	container.add_child(sprite)

	# HP 條
	var hp_max = data.get("max_hp", 1)
	var hp_bg = ColorRect.new()
	hp_bg.size = Vector2(64, 6)
	hp_bg.position = Vector2(-32, -52)
	hp_bg.color = Color(0.2, 0.2, 0.2, 0.8)
	hp_bg.name = "HPBarBG"
	container.add_child(hp_bg)

	var hp_bar = ColorRect.new()
	hp_bar.size = Vector2(64, 6)
	hp_bar.position = Vector2(-32, -52)
	hp_bar.color = Color(0.2, 0.9, 0.2)
	hp_bar.name = "HPBar"
	container.add_child(hp_bar)

	# 倍率標籤
	if target_type != "boss":
		var mult_label = Label.new()
		mult_label.text = "x%.0f" % multiplier
		mult_label.position = Vector2(-20, 36)
		mult_label.add_theme_font_size_override("font_size", 14)
		mult_label.modulate = _mult_label_color(multiplier)
		mult_label.name = "MultLabel"
		container.add_child(mult_label)

	# 高倍率光暈
	if multiplier >= 30.0:
		_add_glow(container, multiplier)

	# Lucky 特殊魚標記（T106-T160）
	if def_id.begins_with("T1") and def_id.length() == 4:
		var tid_num = int(def_id.substr(1))
		if tid_num >= 106 and tid_num <= 155:
			_add_lucky_badge(container, def_id)

	# 特殊搖晃（T103 流星、T104 金草）
	if def_id in ["T103", "T104"]:
		var wobble = container.create_tween().set_loops()
		var deg = 5.0 if def_id == "T103" else 3.0
		var dur = 0.15 if def_id == "T103" else 0.4
		wobble.tween_property(container, "rotation_degrees", deg, dur)
		wobble.tween_property(container, "rotation_degrees", -deg, dur)

	# 儲存 meta
	container.set_meta("instance_id", data.get("instance_id", ""))
	container.set_meta("def_id", def_id)
	container.set_meta("speed", data.get("speed", 0.0))
	container.set_meta("behavior", data.get("behavior", "linear"))
	container.set_meta("target_type", target_type)
	container.set_meta("multiplier", multiplier)
	container.set_meta("is_fleeing", false)

	return container

func _add_glow(node: Node2D, multiplier: float) -> void:
	var glow = ColorRect.new()
	glow.size = Vector2(80, 80)
	glow.position = Vector2(-40, -40)
	glow.z_index = -1
	if multiplier >= 50:
		glow.color = Color(1.0, 0.4, 0.1, 0.3)
	else:
		glow.color = Color(1.0, 0.85, 0.0, 0.25)
	node.add_child(glow)
	var tween = glow.create_tween().set_loops()
	tween.tween_property(glow, "modulate:a", 0.1, 0.4)
	tween.tween_property(glow, "modulate:a", 1.0, 0.4)

# Lucky 特殊魚標記（T106-T125）— 左上角 LUCKY 徽章 + 脈動光環
func _add_lucky_badge(node: Node2D, def_id: String) -> void:
	# 脈動光環（比普通光暈更大更亮）
	var ring = ColorRect.new()
	ring.size = Vector2(96, 96)
	ring.position = Vector2(-48, -48)
	ring.z_index = -1
	# 依倍率範圍選顏色
	var tid_num = int(def_id.substr(1))
	var ring_color: Color
	if tid_num >= 141:
		ring_color = Color(1.0, 1.0, 0.5, 0.60)    # 超亮金（T141+，最高階）
	elif tid_num >= 131:
		ring_color = Color(1.0, 0.95, 0.0, 0.50)   # 亮金（T131-T140）
	elif tid_num >= 126:
		ring_color = Color(1.0, 0.85, 0.0, 0.40)   # 金色（T126-T130）
	elif tid_num >= 121:
		ring_color = Color(0.88, 0.67, 1.0, 0.35)  # 淡紫（T121-T125）
	elif tid_num >= 116:
		ring_color = Color(1.0, 0.85, 0.0, 0.35)   # 金色（T116-T120）
	elif tid_num >= 111:
		ring_color = Color(1.0, 0.42, 0.21, 0.35)  # 火橙（T111-T115）
	else:
		ring_color = Color(0.0, 0.9, 1.0, 0.35)    # 青藍（T106-T110）
	ring.color = ring_color
	node.add_child(ring)

	# 脈動動畫（比普通光暈快）
	var tween = ring.create_tween().set_loops()
	tween.tween_property(ring, "modulate:a", 0.05, 0.25)
	tween.tween_property(ring, "modulate:a", 1.0, 0.25)

	# LUCKY 徽章（左上角小標籤）
	var badge = Label.new()
	badge.text = "✨"
	badge.position = Vector2(-48, -68)
	badge.add_theme_font_size_override("font_size", 14)
	badge.z_index = 10
	node.add_child(badge)

	# 徽章浮動動畫
	var badge_tween = badge.create_tween().set_loops()
	badge_tween.tween_property(badge, "position:y", -72.0, 0.5)
	badge_tween.tween_property(badge, "position:y", -68.0, 0.5)

# ── 目標物更新 ────────────────────────────────────────────────

func _on_target_updated(data: Dictionary) -> void:
	var instance_id = data.get("instance_id", "")
	if not _target_nodes.has(instance_id):
		return
	var node = _target_nodes[instance_id]
	if not is_instance_valid(node):
		return

	# 更新 HP 條
	update_target_hp(instance_id, data.get("hp", 0), data.get("max_hp", 1))

	# 受擊閃白
	_flash_hit(node)

	# T102 逃跑
	if data.get("is_fleeing", false):
		node.set_meta("is_fleeing", true)
		node.set_meta("flee_speed", node.get_meta("speed", 70.0) * 2.5)
		var tween = node.create_tween()
		tween.tween_property(node, "modulate", Color(2.0, 0.5, 0.5), 0.06)
		tween.tween_property(node, "modulate", Color.WHITE, 0.06)
		tween.tween_property(node, "modulate", Color(2.0, 0.5, 0.5), 0.06)
		tween.tween_property(node, "modulate", Color.WHITE, 0.06)

func update_target_hp(instance_id: String, hp: int, max_hp: int) -> void:
	if not _target_nodes.has(instance_id):
		return
	var node = _target_nodes[instance_id]
	if not is_instance_valid(node):
		return
	var hp_bar = node.get_node_or_null("HPBar")
	var hp_bg = node.get_node_or_null("HPBarBG")
	if not is_instance_valid(hp_bar) or not is_instance_valid(hp_bg):
		return
	var pct = float(hp) / float(max_hp) if max_hp > 0 else 0.0
	hp_bar.size.x = hp_bg.size.x * pct
	if pct > 0.6:
		hp_bar.color = Color(0.2, 0.9, 0.2)
	elif pct > 0.3:
		hp_bar.color = Color(1.0, 0.8, 0.1)
	else:
		hp_bar.color = Color(1.0, 0.2, 0.2)

func _flash_hit(node: Node2D) -> void:
	for child in node.get_children():
		if child is Sprite2D:
			var tween = node.create_tween()
			tween.tween_property(child, "modulate", Color(3.0, 3.0, 3.0), 0.04)
			tween.tween_property(child, "modulate", Color.WHITE, 0.08)
			break

# ── 目標物擊破 ────────────────────────────────────────────────

func _on_target_killed(data: Dictionary) -> void:
	var instance_id = data.get("instance_id", "")
	if not _target_nodes.has(instance_id):
		return
	var node = _target_nodes[instance_id]
	_target_nodes.erase(instance_id)

	if not is_instance_valid(node):
		return

	var reward = data.get("reward", 0)
	var multiplier = data.get("multiplier", 1.0)

	# 擊破特效
	HitEffect.spawn_kill(node.position, multiplier)
	if reward > 0:
		HitEffect.spawn_reward_text(node.position, reward, multiplier)

	# 大獎演出
	if multiplier >= 20:
		HitEffect.spawn_big_win(node.position, multiplier)
		AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
	else:
		AudioManager.play_sfx(AudioManager.SFX.KILL)

	# T105 金幣魚：金幣雨
	if data.get("def_id", "") == "T105":
		_spawn_coin_rain(node.position)

	# 消失動畫
	var tween = node.create_tween()
	tween.tween_property(node, "scale", Vector2(1.5, 1.5), 0.1)
	tween.parallel().tween_property(node, "modulate:a", 0.0, 0.15)
	tween.tween_callback(func(): if is_instance_valid(node): node.queue_free())

func _spawn_coin_rain(pos: Vector2) -> void:
	for i in 10:
		var coin = ColorRect.new()
		coin.size = Vector2(12, 12)
		coin.color = Color(1.0, 0.85, 0.0)
		coin.position = pos
		coin.z_index = 45
		get_parent().add_child(coin)
		var angle = randf_range(-PI/2 - 0.5, -PI/2 + 0.5)
		var dist = randf_range(40, 100)
		var target = pos + Vector2(cos(angle), sin(angle)) * dist
		var tween = coin.create_tween()
		tween.tween_property(coin, "position", target, 0.4)
		tween.tween_property(coin, "position:y", target.y + 60, 0.3)
		tween.parallel().tween_property(coin, "modulate:a", 0.0, 0.3)
		tween.tween_callback(func(): if is_instance_valid(coin): coin.queue_free())

# ── BOSS 事件 ─────────────────────────────────────────────────

func _on_boss_event(event_data: Dictionary) -> void:
	var event = event_data.get("event", "")
	if event == "phase_change":
		var phase = event_data.get("phase", 2)
		var instance_id = event_data.get("instance_id", "")
		if _target_nodes.has(instance_id):
			var node = _target_nodes[instance_id]
			if is_instance_valid(node):
				if phase == 3:
					# Phase 3 絕望模式：更快閃爍 + 更大縮放 + 深紅色調
					var tween = node.create_tween()
					# 5次快速閃爍（比 Phase 2 更快，0.04s vs 0.06s）
					for i in 5:
						tween.tween_property(node, "modulate", Color(5.0, 0.1, 0.1), 0.04)
						tween.tween_property(node, "modulate", Color(0.3, 0.05, 0.05), 0.04)
					# 放大到 1.3x（比 Phase 2 的 1.2x 更大）
					tween.tween_property(node, "scale", Vector2(1.3, 1.3), 0.4).set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_ELASTIC)
					# 持續深紅色調（比 Phase 2 更深）
					tween.tween_property(node, "modulate", Color(2.5, 0.2, 0.2), 0.2)
					# 加入 Phase 3 標記
					_add_phase3_indicator(node)
					ScreenShake.add_trauma(1.0)
					AudioManager.play_bgm(AudioManager.BGM.BOSS_RAGE)
				else:
					# Phase 2 進入：強烈閃爍 + 放大 + 持續紅色調
					var tween = node.create_tween()
					# 5次強烈閃爍
					for i in 5:
						tween.tween_property(node, "modulate", Color(4.0, 0.2, 0.2), 0.06)
						tween.tween_property(node, "modulate", Color(0.5, 0.1, 0.1), 0.06)
					# 放大到 1.2x（Phase 2 BOSS 更大更威脅）
					tween.tween_property(node, "scale", Vector2(1.2, 1.2), 0.4).set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_ELASTIC)
					# 持續紅色調（Phase 2 標誌）
					tween.tween_property(node, "modulate", Color(1.8, 0.4, 0.4), 0.2)
					# 加入 Phase 2 標記
					_add_phase2_indicator(node)
					ScreenShake.add_trauma(0.8)
					AudioManager.play_bgm(AudioManager.BGM.BOSS_RAGE)

func _add_phase2_indicator(boss_node: Node2D) -> void:
	# 移除舊的 Phase 2 標記（如果有）
	var old = boss_node.get_node_or_null("Phase2Label")
	if is_instance_valid(old):
		old.queue_free()
	# 新增 PHASE 2 標籤
	var lbl = Label.new()
	lbl.name = "Phase2Label"
	lbl.text = "⚠ PHASE 2"
	lbl.position = Vector2(-40, -80)
	lbl.add_theme_font_size_override("font_size", 14)
	lbl.modulate = Color(1.0, 0.2, 0.2)
	lbl.z_index = 15
	boss_node.add_child(lbl)
	# 脈動動畫
	var tween = lbl.create_tween().set_loops()
	tween.tween_property(lbl, "modulate:a", 0.3, 0.3)
	tween.tween_property(lbl, "modulate:a", 1.0, 0.3)

func _add_phase3_indicator(boss_node: Node2D) -> void:
	# 移除舊的 Phase 2 / Phase 3 標記（如果有）
	for label_name in ["Phase2Label", "Phase3Label"]:
		var old = boss_node.get_node_or_null(label_name)
		if is_instance_valid(old):
			old.queue_free()
	# 新增 PHASE 3 標籤（紅色，更大字體）
	var lbl = Label.new()
	lbl.name = "Phase3Label"
	lbl.text = "💀 PHASE 3 絕望模式！"
	lbl.position = Vector2(-72, -96)
	lbl.add_theme_font_size_override("font_size", 18)
	lbl.modulate = Color(1.0, 0.0, 0.0)
	lbl.z_index = 20
	boss_node.add_child(lbl)
	# 快速脈動動畫（比 Phase 2 更快）
	var tween = lbl.create_tween().set_loops()
	tween.tween_property(lbl, "modulate:a", 0.1, 0.15)
	tween.tween_property(lbl, "modulate:a", 1.0, 0.15)

# ── 點擊目標 ──────────────────────────────────────────────────

func try_click_target(click_pos: Vector2) -> String:
	var best_id = ""
	var best_dist = 80.0  # 點擊半徑
	for instance_id in _target_nodes:
		var node = _target_nodes[instance_id]
		if not is_instance_valid(node):
			continue
		var dist = node.position.distance_to(click_pos)
		if dist < best_dist:
			best_dist = dist
			best_id = instance_id
	return best_id

# ── 輔助 ──────────────────────────────────────────────────────

func _get_texture(def_id: String) -> Texture2D:
	var path = TARGET_SPRITES.get(def_id, "")
	if path == "":
		return null
	if _cached_textures.has(path):
		return _cached_textures[path]
	if ResourceLoader.exists(path):
		var tex = load(path)
		_cached_textures[path] = tex
		return tex
	return null

func _mult_label_color(mult: float) -> Color:
	if mult >= 100: return Color(1.0, 0.3, 0.1)   # 火紅（超高倍率）
	if mult >= 50:  return Color(1.0, 0.4, 0.1)   # 橙紅
	if mult >= 30:  return Color(1.0, 0.85, 0.0)  # 金色
	if mult >= 15:  return Color(0.8, 0.9, 1.0)   # 淡藍
	return Color(1.0, 1.0, 1.0)                   # 白色

# ── DAY-301 時間扭曲速度效果 ──────────────────────────────────

func _on_lucky_time_warp_for_speed(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"warp_start":
			_time_warp_speed_mult = data.get("speed_mult", 0.3)
			# 視覺提示：所有目標物變藍色調
			for instance_id in _target_nodes:
				var node = _target_nodes[instance_id]
				if is_instance_valid(node):
					var tween = node.create_tween()
					tween.tween_property(node, "modulate", Color(0.7, 0.8, 1.2), 0.3)
		"warp_end", "time_collapse", "collapse_end":
			_time_warp_speed_mult = 1.0
			# 恢復正常顏色
			for instance_id in _target_nodes:
				var node = _target_nodes[instance_id]
				if is_instance_valid(node):
					var tween = node.create_tween()
					tween.tween_property(node, "modulate", Color.WHITE, 0.3)
