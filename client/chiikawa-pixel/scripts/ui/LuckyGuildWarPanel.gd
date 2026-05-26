## LuckyGuildWarPanel.gd — T139 幸運公會戰魚 UI
## lucky-panel-agent 負責維護
## 業界依據：Fishing Frenzy Chapter 3「Guild Wars — territory control and leaderboard rankings」
## 視覺主題：金色 + 積分進度條 + 公會勝利演出
extends CanvasLayer

const LAYER_Z = 39

const COLOR_GUILD    = Color(1.0, 0.85, 0.0)   # 金色（公會主色）
const COLOR_PROGRESS = Color(0.0, 0.9, 0.3)    # 翠綠（進度）
const COLOR_VICTORY  = Color(1.0, 0.5, 0.0)    # 火橙（勝利）
const COLOR_BG       = Color(0.05, 0.04, 0.0, 0.90)

var _banner: Control = null
var _indicator: Control = null
var _progress_bar: ColorRect = null
var _score_label: Label = null
var _time_label: Label = null
var _flash_overlay: ColorRect = null
var _war_timer: float = 0.0
var _target_points: int = 15
var _current_points: int = 0

func _ready() -> void:
	layer = LAYER_Z
	_create_flash_overlay()
	_create_indicator()
	GameManager.lucky_guild_war.connect(_on_lucky_guild_war)

func _process(delta: float) -> void:
	if _war_timer > 0.0:
		_war_timer -= delta
		if _war_timer < 0.0:
			_war_timer = 0.0
		if is_instance_valid(_time_label):
			_time_label.text = "剩餘：%.0fs" % _war_timer

func _create_flash_overlay() -> void:
	_flash_overlay = ColorRect.new()
	_flash_overlay.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	_flash_overlay.color = Color(1.0, 0.85, 0.0, 0.0)
	_flash_overlay.mouse_filter = Control.MOUSE_FILTER_IGNORE
	add_child(_flash_overlay)

func _create_indicator() -> void:
	_indicator = Control.new()
	_indicator.position = Vector2(960, 10)
	_indicator.size = Vector2(300, 110)
	_indicator.visible = false
	add_child(_indicator)

	var bg = ColorRect.new()
	bg.size = _indicator.size
	bg.color = COLOR_BG
	_indicator.add_child(bg)

	var title = Label.new()
	title.text = "⚔️ 公會戰進行中"
	title.position = Vector2(8, 4)
	title.add_theme_font_size_override("font_size", 13)
	title.modulate = COLOR_GUILD
	_indicator.add_child(title)

	_score_label = Label.new()
	_score_label.text = "積分：0/15"
	_score_label.position = Vector2(8, 28)
	_score_label.add_theme_font_size_override("font_size", 18)
	_score_label.modulate = COLOR_PROGRESS
	_indicator.add_child(_score_label)

	_time_label = Label.new()
	_time_label.text = "剩餘：30s"
	_time_label.position = Vector2(8, 56)
	_time_label.add_theme_font_size_override("font_size", 14)
	_time_label.modulate = Color(0.8, 0.8, 0.8)
	_indicator.add_child(_time_label)

	var bar_bg = ColorRect.new()
	bar_bg.position = Vector2(8, 82)
	bar_bg.size = Vector2(280, 14)
	bar_bg.color = Color(0.2, 0.2, 0.2)
	_indicator.add_child(bar_bg)

	_progress_bar = ColorRect.new()
	_progress_bar.position = Vector2(8, 82)
	_progress_bar.size = Vector2(0, 14)
	_progress_bar.color = COLOR_PROGRESS
	_indicator.add_child(_progress_bar)

func _flash(color: Color, alpha: float = 0.4, duration: float = 0.2) -> void:
	if not is_instance_valid(_flash_overlay):
		return
	_flash_overlay.color = Color(color.r, color.g, color.b, alpha)
	var tween = create_tween()
	tween.tween_property(_flash_overlay, "color:a", 0.0, duration)

func _show_banner(text: String, color: Color) -> void:
	if is_instance_valid(_banner):
		_banner.queue_free()
	_banner = Control.new()
	_banner.position = Vector2(0, 100)
	_banner.size = Vector2(1280, 70)
	add_child(_banner)
	var bg = ColorRect.new()
	bg.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	bg.color = Color(0, 0, 0, 0.85)
	_banner.add_child(bg)
	var lbl = Label.new()
	lbl.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	lbl.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	lbl.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	lbl.text = text
	lbl.add_theme_font_size_override("font_size", 26)
	lbl.modulate = color
	_banner.add_child(lbl)
	var tween = _banner.create_tween()
	tween.tween_interval(2.5)
	tween.tween_property(_banner, "modulate:a", 0.0, 0.5)
	tween.tween_callback(func():
		if is_instance_valid(_banner):
			_banner.queue_free()
			_banner = null
	)

func _on_lucky_guild_war(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"war_start":
			_war_timer = data.get("duration", 30.0)
			_target_points = data.get("target_points", 15)
			_current_points = 0
			_indicator.visible = true
			if is_instance_valid(_score_label):
				_score_label.text = "積分：0/%d" % _target_points
			if is_instance_valid(_progress_bar):
				_progress_bar.size.x = 0
			_flash(COLOR_GUILD, 0.4, 0.2)
			_show_banner("⚔️ 公會戰開始！30 秒內達成 %d 積分！" % _target_points, COLOR_GUILD)

		"war_progress":
			_current_points = data.get("current_points", 0)
			if is_instance_valid(_score_label):
				_score_label.text = "積分：%d/%d" % [_current_points, _target_points]
			if is_instance_valid(_progress_bar) and _target_points > 0:
				_progress_bar.size.x = 280.0 * float(_current_points) / float(_target_points)
			_flash(COLOR_PROGRESS, 0.2, 0.1)

		"war_victory":
			_indicator.visible = false
			_flash(COLOR_VICTORY, 0.7, 0.3)
			_flash(COLOR_VICTORY, 0.7, 0.3)
			_flash(COLOR_VICTORY, 0.7, 0.3)
			var points = data.get("current_points", 0)
			_show_banner("⚔️ 公會勝利！達成 %d 積分！全服 ×4.5 加成 10 秒！" % points, COLOR_VICTORY)

		"war_victory_end":
			pass

		"war_timeout":
			_indicator.visible = false
