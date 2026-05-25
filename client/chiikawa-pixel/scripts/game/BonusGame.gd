## BonusGame.gd — Bonus 遊戲（瘋狂拔草場景）
## bonus-game-agent 負責維護
## 當 bonus_event start 時顯示，結束後隱藏
extends CanvasLayer

# ── 常數 ─────────────────────────────────────────────────────
const WEED_COLORS = {
	"BG001": Color(0.2, 0.8, 0.2),    # 普通雜草 - 綠
	"BG002": Color(0.5, 0.35, 0.1),   # 硬雜草 - 棕
	"BG003": Color(0.3, 1.0, 0.5),    # 發光雜草 - 亮綠
	"BG004": Color(1.0, 0.85, 0.0),   # 金色雜草 - 金
	"BG005": Color(0.9, 0.2, 0.2),    # 搗亂怪草 - 紅
}
const WEED_SCORES = {
	"BG001": 1, "BG002": 3, "BG003": 8, "BG004": 20, "BG005": -5
}
const WEED_LABELS = {
	"BG001": "+1", "BG002": "+3", "BG003": "+8", "BG004": "+20", "BG005": "-5"
}
const WEED_EMOJIS = {
	"BG001": "🌿", "BG002": "🌾", "BG003": "✨", "BG004": "🌟", "BG005": "💀"
}

const WEED_SPAWN_INTERVAL = 0.6
const MAX_WEEDS = 18
const BONUS_DURATION = 15.0

# ── 節點引用 ──────────────────────────────────────────────────
var _overlay: Control = null
var _bg: ColorRect = null
var _timer_label: Label = null
var _score_label: Label = null
var _title_label: Label = null
var _result_panel: Control = null
var _weed_container: Control = null

# ── 狀態 ─────────────────────────────────────────────────────
var _active: bool = false
var _time_left: float = 0.0
var _score: int = 0
var _weed_nodes: Dictionary = {}  # weed_id -> Button
var _weed_hp: Dictionary = {}     # weed_id -> hp（BG002 需要 2 次）
var _spawn_timer: float = 0.0
var _weed_counter: int = 0

func _ready() -> void:
	_build_ui()
	GameManager.bonus_event.connect(_on_bonus_event)
	visible = false

func _process(delta: float) -> void:
	if not _active:
		return
	_time_left -= delta
	_update_timer()
	_spawn_timer += delta
	if _spawn_timer >= WEED_SPAWN_INTERVAL and _weed_nodes.size() < MAX_WEEDS:
		_spawn_timer = 0.0
		_spawn_weed()

# ── UI 建立 ───────────────────────────────────────────────────

func _build_ui() -> void:
	layer = 80

	# 全螢幕覆蓋
	_overlay = Control.new()
	_overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	add_child(_overlay)

	# 背景（草地綠）
	_bg = ColorRect.new()
	_bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_bg.color = Color(0.08, 0.22, 0.08, 0.97)
	_overlay.add_child(_bg)

	# 標題
	_title_label = Label.new()
	_title_label.text = "🌿 瘋狂拔草！BONUS GAME 🌿"
	_title_label.position = Vector2(0, 12)
	_title_label.size = Vector2(1280, 40)
	_title_label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	_title_label.add_theme_font_size_override("font_size", 26)
	_title_label.modulate = Color(0.3, 1.0, 0.4)
	_overlay.add_child(_title_label)

	# 計時器
	_timer_label = Label.new()
	_timer_label.text = "⏱ 15.0s"
	_timer_label.position = Vector2(900, 12)
	_timer_label.size = Vector2(200, 40)
	_timer_label.add_theme_font_size_override("font_size", 22)
	_timer_label.modulate = Color(1.0, 0.9, 0.2)
	_overlay.add_child(_timer_label)

	# 分數
	_score_label = Label.new()
	_score_label.text = "分數：0"
	_score_label.position = Vector2(80, 12)
	_score_label.size = Vector2(200, 40)
	_score_label.add_theme_font_size_override("font_size", 22)
	_score_label.modulate = Color(1.0, 1.0, 1.0)
	_overlay.add_child(_score_label)

	# 雜草容器
	_weed_container = Control.new()
	_weed_container.position = Vector2(0, 60)
	_weed_container.size = Vector2(1280, 620)
	_overlay.add_child(_weed_container)

	# 結算面板（初始隱藏）
	_result_panel = _build_result_panel()
	_overlay.add_child(_result_panel)

func _build_result_panel() -> Control:
	var panel = Control.new()
	panel.position = Vector2(390, 220)
	panel.size = Vector2(500, 280)
	panel.visible = false
	panel.z_index = 10

	var bg = ColorRect.new()
	bg.size = Vector2(500, 280)
	bg.color = Color(0.05, 0.1, 0.05, 0.95)
	panel.add_child(bg)

	var border = ColorRect.new()
	border.size = Vector2(500, 4)
	border.color = Color(0.3, 1.0, 0.4)
	panel.add_child(border)

	var title = Label.new()
	title.name = "Title"
	title.text = "🌿 BONUS 結算！"
	title.position = Vector2(0, 20)
	title.size = Vector2(500, 40)
	title.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	title.add_theme_font_size_override("font_size", 28)
	title.modulate = Color(0.3, 1.0, 0.4)
	panel.add_child(title)

	var score_lbl = Label.new()
	score_lbl.name = "ScoreLabel"
	score_lbl.text = "拔草分數：0"
	score_lbl.position = Vector2(0, 80)
	score_lbl.size = Vector2(500, 36)
	score_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	score_lbl.add_theme_font_size_override("font_size", 22)
	score_lbl.modulate = Color.WHITE
	panel.add_child(score_lbl)

	var mult_lbl = Label.new()
	mult_lbl.name = "MultLabel"
	mult_lbl.text = "倍率：×20"
	mult_lbl.position = Vector2(0, 126)
	mult_lbl.size = Vector2(500, 36)
	mult_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	mult_lbl.add_theme_font_size_override("font_size", 22)
	mult_lbl.modulate = Color(1.0, 0.85, 0.0)
	panel.add_child(mult_lbl)

	var reward_lbl = Label.new()
	reward_lbl.name = "RewardLabel"
	reward_lbl.text = "獎勵：+0 💰"
	reward_lbl.position = Vector2(0, 172)
	reward_lbl.size = Vector2(500, 44)
	reward_lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	reward_lbl.add_theme_font_size_override("font_size", 30)
	reward_lbl.modulate = Color(1.0, 0.85, 0.0)
	panel.add_child(reward_lbl)

	var hint = Label.new()
	hint.name = "HintLabel"
	hint.text = "返回遊戲中..."
	hint.position = Vector2(0, 240)
	hint.size = Vector2(500, 28)
	hint.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	hint.add_theme_font_size_override("font_size", 14)
	hint.modulate = Color(0.6, 0.6, 0.6)
	panel.add_child(hint)

	return panel

# ── 事件處理 ──────────────────────────────────────────────────

func _on_bonus_event(event_data: Dictionary) -> void:
	match event_data.get("event", ""):
		"start":
			_start_bonus()
		"tick":
			_time_left = event_data.get("time_left", _time_left)
		"click":
			_score = event_data.get("score", _score)
			_update_score()
		"result":
			_show_result(event_data)
		"bonus_result", _:
			if event_data.get("event", "") == "bonus_result":
				_end_bonus()

func _start_bonus() -> void:
	_active = true
	_time_left = BONUS_DURATION
	_score = 0
	_spawn_timer = 0.0
	_weed_counter = 0
	_clear_weeds()
	_result_panel.visible = false
	visible = true
	_update_score()
	_update_timer()
	AudioManager.play_bgm(AudioManager.BGM.BONUS_GAME)
	AudioManager.play_sfx(AudioManager.SFX.BONUS_READY)
	# 進場動畫
	_overlay.modulate.a = 0.0
	var tween = create_tween()
	tween.tween_property(_overlay, "modulate:a", 1.0, 0.4)

func _end_bonus() -> void:
	_active = false
	var tween = create_tween()
	tween.tween_interval(2.5)
	tween.tween_property(_overlay, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func():
		visible = false
		_clear_weeds()
		AudioManager.play_bgm(AudioManager.BGM.MAIN_GAME)
	)

func _show_result(data: Dictionary) -> void:
	_active = false
	_clear_weeds()
	var score = data.get("score", _score)
	var mult = data.get("multiplier", 20.0)
	var reward = data.get("reward", 0)

	var sl = _result_panel.get_node_or_null("ScoreLabel")
	var ml = _result_panel.get_node_or_null("MultLabel")
	var rl = _result_panel.get_node_or_null("RewardLabel")
	if is_instance_valid(sl): sl.text = "拔草分數：%d" % score
	if is_instance_valid(ml): ml.text = "倍率：×%.1f" % mult
	if is_instance_valid(rl): rl.text = "獎勵：+%d 💰" % reward

	_result_panel.visible = true
	_result_panel.scale = Vector2.ZERO
	var tween = _result_panel.create_tween()
	tween.tween_property(_result_panel, "scale", Vector2(1.0, 1.0), 0.3).set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_BACK)
	ScreenShake.add_trauma(0.4)
	AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)

	# 2.5 秒後自動關閉
	tween.tween_interval(2.5)
	tween.tween_callback(_end_bonus)

# ── 雜草生成 ──────────────────────────────────────────────────

func _spawn_weed() -> void:
	var weed_id = _pick_weed_type()
	var uid = "weed_%d" % _weed_counter
	_weed_counter += 1

	var btn = Button.new()
	btn.text = WEED_EMOJIS.get(weed_id, "🌿")
	btn.size = Vector2(72, 72)
	btn.position = Vector2(
		randf_range(40, 1200),
		randf_range(20, 520)
	)
	btn.add_theme_font_size_override("font_size", 28)
	btn.modulate = WEED_COLORS.get(weed_id, Color.WHITE)
	btn.z_index = 5
	btn.set_meta("weed_id", weed_id)
	btn.set_meta("uid", uid)

	# BG002 硬雜草需要 2 次點擊
	if weed_id == "BG002":
		_weed_hp[uid] = 2
	else:
		_weed_hp[uid] = 1

	btn.pressed.connect(func(): _on_weed_clicked(btn))
	_weed_container.add_child(btn)
	_weed_nodes[uid] = btn

	# 進場動畫
	btn.scale = Vector2.ZERO
	var tween = btn.create_tween()
	tween.tween_property(btn, "scale", Vector2(1.0, 1.0), 0.15).set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_BACK)

	# 自動消失（8 秒後）
	var auto_tween = btn.create_tween()
	auto_tween.tween_interval(8.0)
	auto_tween.tween_property(btn, "modulate:a", 0.0, 0.3)
	auto_tween.tween_callback(func():
		if is_instance_valid(btn):
			_weed_nodes.erase(uid)
			_weed_hp.erase(uid)
			btn.queue_free()
	)

func _on_weed_clicked(btn: Button) -> void:
	if not is_instance_valid(btn):
		return
	var uid = btn.get_meta("uid", "")
	var weed_id = btn.get_meta("weed_id", "BG001")

	# HP 扣減
	var hp = _weed_hp.get(uid, 1)
	hp -= 1
	_weed_hp[uid] = hp

	# 點擊動畫
	var tween = btn.create_tween()
	tween.tween_property(btn, "scale", Vector2(1.3, 1.3), 0.06)
	tween.tween_property(btn, "scale", Vector2(1.0, 1.0), 0.06)

	if hp > 0:
		# BG002 第一次點擊：變色提示
		btn.modulate = Color(1.0, 0.6, 0.2)
		AudioManager.play_sfx(AudioManager.SFX.HIT)
		return

	# 擊破
	_weed_nodes.erase(uid)
	_weed_hp.erase(uid)

	# 分數跳字
	var score_val = WEED_SCORES.get(weed_id, 1)
	_spawn_score_text(btn.position + btn.size / 2, score_val)

	# 消失動畫
	var kill_tween = btn.create_tween()
	kill_tween.tween_property(btn, "scale", Vector2(1.5, 1.5), 0.08)
	kill_tween.parallel().tween_property(btn, "modulate:a", 0.0, 0.15)
	kill_tween.tween_callback(func(): if is_instance_valid(btn): btn.queue_free())

	# 通知 Server
	NetworkManager.send_bonus_click(weed_id, btn.position.x, btn.position.y)
	AudioManager.play_sfx(AudioManager.SFX.WEED_PULL)

func _spawn_score_text(pos: Vector2, score: int) -> void:
	var lbl = Label.new()
	lbl.text = ("+%d" % score) if score > 0 else ("%d" % score)
	lbl.position = pos + Vector2(-20, -20)
	lbl.add_theme_font_size_override("font_size", 22)
	lbl.modulate = Color(0.3, 1.0, 0.4) if score > 0 else Color(1.0, 0.3, 0.3)
	lbl.z_index = 20
	_weed_container.add_child(lbl)
	var tween = lbl.create_tween()
	tween.tween_property(lbl, "position:y", pos.y - 60, 0.6)
	tween.parallel().tween_property(lbl, "modulate:a", 0.0, 0.6)
	tween.tween_callback(func(): if is_instance_valid(lbl): lbl.queue_free())

func _pick_weed_type() -> String:
	var weights = {"BG001": 180, "BG002": 80, "BG003": 35, "BG004": 10, "BG005": 20}
	var total = 0
	for w in weights.values():
		total += w
	var r = randi() % total
	var acc = 0
	for id in weights:
		acc += weights[id]
		if r < acc:
			return id
	return "BG001"

# ── 輔助 ─────────────────────────────────────────────────────

func _update_timer() -> void:
	if not is_instance_valid(_timer_label):
		return
	var t = max(0.0, _time_left)
	_timer_label.text = "⏱ %.1fs" % t
	if t <= 5.0:
		_timer_label.modulate = Color(1.0, 0.3, 0.3)
	elif t <= 10.0:
		_timer_label.modulate = Color(1.0, 0.7, 0.2)
	else:
		_timer_label.modulate = Color(1.0, 0.9, 0.2)

func _update_score() -> void:
	if not is_instance_valid(_score_label):
		return
	_score_label.text = "分數：%d" % _score

func _clear_weeds() -> void:
	for uid in _weed_nodes:
		var node = _weed_nodes[uid]
		if is_instance_valid(node):
			node.queue_free()
	_weed_nodes.clear()
	_weed_hp.clear()
