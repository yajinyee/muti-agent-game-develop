## LuckyMutationPanel.gd — T181 幸運突變魚 UI
## lucky-panel-agent 負責維護
## DAY-315：突變系統 — 150種突變，最高 ×17.0，突變 ≥10x → 全服 ×16.0 加成 32 秒
extends BaseLuckyPanel

const PANEL_COLOR = Color(0.49, 0.11, 0.64)  # 深紫色（突變）
const PANEL_ICON = "🧬"
const PANEL_TITLE = "突變觸發"

var _mutation_label: Label = null
var _mult_label: Label = null

func _ready() -> void:
	super._ready()
	layer = 76
	_setup_mutation_ui()
	GameManager.lucky_mutation.connect(_on_lucky_mutation)

func _setup_mutation_ui() -> void:
	_mutation_label = Label.new()
	_mutation_label.text = "突變：等待中..."
	_mutation_label.add_theme_color_override("font_color", Color(0.8, 0.4, 1.0))
	_mutation_label.add_theme_font_size_override("font_size", 20)
	_mutation_label.position = Vector2(20, 80)
	add_child(_mutation_label)

	_mult_label = Label.new()
	_mult_label.text = "倍率：×?"
	_mult_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	_mult_label.add_theme_font_size_override("font_size", 18)
	_mult_label.position = Vector2(20, 110)
	add_child(_mult_label)

func _on_lucky_mutation(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"mutation_start":
			var mut_name = data.get("mutation_name", "未知突變")
			var mut_mult = data.get("mutation_mult", 1.0)
			if is_instance_valid(_mutation_label):
				_mutation_label.text = "突變：" + mut_name
			if is_instance_valid(_mult_label):
				_mult_label.text = "倍率：×" + str(mut_mult)
			show_banner(PANEL_ICON + " " + PANEL_TITLE + "！", PANEL_COLOR)
			show_panel()
			flash_screen(Color(0.5, 0.0, 0.8))
		"mutation_perfect":
			var mut_name = data.get("mutation_name", "傳說突變")
			var mut_mult = data.get("mutation_mult", 17.0)
			var boost_mult = data.get("boost_mult", 16.0)
			var boost_secs = data.get("boost_secs", 32)
			show_settle(PANEL_ICON + " 傳說突變！", "【" + mut_name + "】×" + str(mut_mult) + "！全服 ×" + str(boost_mult) + " 加成 " + str(boost_secs) + " 秒！", PANEL_COLOR)
			flash_screen(Color(0.6, 0.0, 1.0))
			hide_panel()
		"mutation_complete":
			var mut_name = data.get("mutation_name", "突變")
			var mut_mult = data.get("mutation_mult", 1.0)
			show_settle(PANEL_ICON + " 突變完成！", "【" + mut_name + "】×" + str(mut_mult) + "！", PANEL_COLOR)
			hide_panel()
