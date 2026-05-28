## LuckyChainEelPanel.gd — T209 幸運連鎖電鰻魚 UI
## lucky-panel-agent 負責維護
## DAY-323：連鎖電鰻系統 — 連鎖電擊 8 條魚，每條 ×40.0，全服 ×34.0 加成 68 秒
## 業界依據：Royal Fishing「Purple/Pink Lightning Eel chain reaction」升級版
extends BaseLuckyPanel

const PANEL_COLOR = Color(0.8, 0.0, 1.0)   # 深紫色（連鎖電鰻）
const PANEL_ICON = "⚡"
const PANEL_TITLE = "連鎖電鰻"

var _chain_label: Label = null
var _mult_label: Label = null
var _result_label: Label = null
var _global_label: Label = null

func _ready() -> void:
	super._ready()
	layer = 104
	_setup_chain_eel_ui()
	GameManager.lucky_chain_eel.connect(_on_lucky_chain_eel)

func _setup_chain_eel_ui() -> void:
	_chain_label = Label.new()
	_chain_label.text = "⚡ 連鎖電鰻啟動！"
	_chain_label.add_theme_color_override("font_color", Color(0.8, 0.0, 1.0))
	_chain_label.add_theme_font_size_override("font_size", 22)
	_chain_label.position = Vector2(20, 52)
	add_child(_chain_label)

	_mult_label = Label.new()
	_mult_label.text = "8 條連鎖電擊 × ×40.0"
	_mult_label.add_theme_color_override("font_color", Color(1.0, 0.5, 1.0))
	_mult_label.add_theme_font_size_override("font_size", 20)
	_mult_label.position = Vector2(20, 80)
	add_child(_mult_label)

	_result_label = Label.new()
	_result_label.text = "電擊中..."
	_result_label.add_theme_color_override("font_color", Color(1.0, 1.0, 1.0))
	_result_label.add_theme_font_size_override("font_size", 20)
	_result_label.position = Vector2(20, 108)
	add_child(_result_label)

	_global_label = Label.new()
	_global_label.text = "全服 ×34.0 加成 68 秒"
	_global_label.add_theme_color_override("font_color", Color(1.0, 1.0, 0.0))
	_global_label.add_theme_font_size_override("font_size", 16)
	_global_label.position = Vector2(20, 136)
	add_child(_global_label)

func _on_lucky_chain_eel(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"eel_start":
			var chain_count = data.get("chain_count", 8)
			var reward_mult = data.get("reward_mult", 40.0)
			show_banner(PANEL_ICON + " " + PANEL_TITLE + " 觸發！", PANEL_COLOR)
			show_panel()
			# 紫色閃屏
			flash_screen(Color(0.8, 0.0, 1.0))
			var tween = create_tween()
			tween.tween_interval(0.1)
			tween.tween_callback(func(): flash_screen(Color(1.0, 0.5, 1.0)))
			if is_instance_valid(_mult_label):
				_mult_label.text = "%d 條連鎖電擊 × ×%.1f" % [chain_count, reward_mult]
		"eel_complete":
			var hit_count = data.get("hit_count", 0)
			var reward_mult = data.get("reward_mult", 40.0)
			var is_perfect = data.get("is_perfect", false)
			var global_mult = data.get("global_mult", 34.0)
			var global_secs = data.get("global_secs", 68)
			var title = "⚡ 完美連鎖！" if is_perfect else "⚡ 連鎖電鰻完成！"
			show_settle(title,
				"電擊 %d 條（×%.1f）！全服 ×%.1f 加成 %d 秒！" % [hit_count, reward_mult, global_mult, global_secs],
				PANEL_COLOR)
			if is_perfect:
				flash_screen(Color(1.0, 1.0, 1.0))
			else:
				flash_screen(Color(0.8, 0.0, 1.0))
			hide_panel()
