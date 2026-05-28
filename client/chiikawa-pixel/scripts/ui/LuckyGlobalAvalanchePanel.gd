## LuckyGlobalAvalanchePanel.gd — T215 幸運全服雪崩魚 UI
## lucky-panel-agent 負責維護
## DAY-324：全服雪崩系統 — 5 波全服連鎖消除，每波 ×8.0，全服 ×38.0 加成 76 秒（新史上最高）
## 業界依據：Avalanche Reels + Global Multiplier 組合（2026 最新趨勢）
extends BaseLuckyPanel

const PANEL_COLOR = Color(0.53, 0.81, 0.98)   # 天藍色（全服雪崩）
const PANEL_ICON = "❄️"
const PANEL_TITLE = "全服雪崩"

var _wave_label: Label = null
var _killed_label: Label = null
var _total_label: Label = null
var _boost_label: Label = null
var _record_label: Label = null

func _ready() -> void:
	super._ready()
	layer = 110
	_setup_global_avalanche_ui()
	GameManager.lucky_global_avalanche.connect(_on_lucky_global_avalanche)

func _setup_global_avalanche_ui() -> void:
	_record_label = Label.new()
	_record_label.text = "🏆 新史上最高全服倍率 ×38.0"
	_record_label.add_theme_color_override("font_color", Color(1.0, 1.0, 0.0))
	_record_label.add_theme_font_size_override("font_size", 15)
	_record_label.position = Vector2(20, 52)
	add_child(_record_label)

	_wave_label = Label.new()
	_wave_label.text = "波次: 0 / 5"
	_wave_label.add_theme_color_override("font_color", Color(0.53, 0.81, 0.98))
	_wave_label.add_theme_font_size_override("font_size", 22)
	_wave_label.position = Vector2(20, 72)
	add_child(_wave_label)

	_killed_label = Label.new()
	_killed_label.text = "本波消滅: 0 個"
	_killed_label.add_theme_color_override("font_color", Color(1.0, 1.0, 1.0))
	_killed_label.add_theme_font_size_override("font_size", 18)
	_killed_label.position = Vector2(20, 98)
	add_child(_killed_label)

	_total_label = Label.new()
	_total_label.text = "累計消滅: 0 個"
	_total_label.add_theme_color_override("font_color", Color(0.9, 0.9, 1.0))
	_total_label.add_theme_font_size_override("font_size", 16)
	_total_label.position = Vector2(20, 122)
	add_child(_total_label)

	_boost_label = Label.new()
	_boost_label.text = "全服 ×38.0 加成 76 秒（新史上最高）"
	_boost_label.add_theme_color_override("font_color", Color(1.0, 1.0, 0.0))
	_boost_label.add_theme_font_size_override("font_size", 13)
	_boost_label.position = Vector2(20, 145)
	add_child(_boost_label)

func _on_lucky_global_avalanche(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"global_avalanche_start":
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 開始！", PANEL_COLOR)
			show_panel()
			# 三重閃屏：冰藍 → 白 → 冰藍
			flash_screen(PANEL_COLOR)
			var tween = create_tween()
			tween.tween_interval(0.1)
			tween.tween_callback(func(): flash_screen(Color(1.0, 1.0, 1.0)))
			tween.tween_interval(0.1)
			tween.tween_callback(func(): flash_screen(PANEL_COLOR))
		"global_wave_hit":
			var wave = data.get("wave", 0)
			var killed = data.get("killed", 0)
			var total_killed = data.get("total_killed", 0)
			if is_instance_valid(_wave_label):
				_wave_label.text = "波次: " + str(wave) + " / 5"
			if is_instance_valid(_killed_label):
				_killed_label.text = "本波消滅: " + str(killed) + " 個"
			if is_instance_valid(_total_label):
				_total_label.text = "累計消滅: " + str(total_killed) + " 個"
		"global_avalanche_complete":
			var total_killed = data.get("total_killed", 0)
			var global_mult = data.get("global_mult", 38.0)
			var global_secs = data.get("global_secs", 76)
			show_settle(PANEL_ICON + " 全服雪崩完成！",
				"消滅 " + str(total_killed) + " 個！全服 ×" + str(global_mult) + " 加成 " + str(global_secs) + " 秒！（新史上最高）",
				PANEL_COLOR)
			# 最強閃屏
			flash_screen(Color(1.0, 1.0, 1.0))
			var tween2 = create_tween()
			tween2.tween_interval(0.1)
			tween2.tween_callback(func(): flash_screen(PANEL_COLOR))
			tween2.tween_interval(0.1)
			tween2.tween_callback(func(): flash_screen(Color(1.0, 1.0, 1.0)))
			hide_panel()
