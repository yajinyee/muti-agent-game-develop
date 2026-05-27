## LuckyJackpotGrandPanel.gd — T174 幸運 Grand Jackpot 魚 UI
## progressive-jackpot-agent 負責維護
## DAY-313：Grand Jackpot 系統 — 擊破後直接觸發 Grand Jackpot（5000x 起跳累積獎池）
## 最高層級 Jackpot，觸發時全螢幕慶典演出
extends BaseLuckyPanel

const PANEL_COLOR = Color(1.0, 0.85, 0.0)  # 金色（Grand）
const PANEL_ICON = "🎰👑"
const PANEL_TITLE = "GRAND JACKPOT"

var _pool_label: Label = null
var _celebration_timer: float = 0.0
var _is_celebrating: bool = false

func _ready() -> void:
	super._ready()
	layer = 69
	_setup_grand_jackpot_ui()
	GameManager.lucky_jackpot_pool.connect(_on_lucky_jackpot_pool)

func _setup_grand_jackpot_ui() -> void:
	_pool_label = Label.new()
	_pool_label.text = "Grand Pool: 5000x"
	_pool_label.add_theme_color_override("font_color", Color(1.0, 0.85, 0.0))
	_pool_label.add_theme_font_size_override("font_size", 22)
	_pool_label.position = Vector2(20, 80)
	add_child(_pool_label)

func _process(delta: float) -> void:
	if _is_celebrating:
		_celebration_timer -= delta
		if _celebration_timer <= 0:
			_is_celebrating = false

func _on_lucky_jackpot_pool(data: Dictionary) -> void:
	var event = data.get("event", "")
	match event:
		"pool_update":
			var grand = data.get("grand", 5000.0)
			if is_instance_valid(_pool_label):
				_pool_label.text = "Grand Pool: %.0fx" % grand
		"jackpot_win":
			var tier = data.get("tier", "")
			if tier != "grand":
				return
			var reward = data.get("reward", 0)
			var player_name = data.get("player_name", "玩家")
			_trigger_grand_celebration(player_name, reward)

func _trigger_grand_celebration(player_name: String, reward: int) -> void:
	_is_celebrating = true
	_celebration_timer = 5.0

	# 全螢幕金色閃光（3次）
	flash_screen(Color(1.0, 0.85, 0.0))
	flash_screen(Color(1.0, 0.85, 0.0))
	flash_screen(Color(1.0, 0.85, 0.0))

	# 強烈震動
	ScreenShake.add_trauma(1.0)

	# 橫幅
	show_banner(PANEL_ICON + " ★ GRAND JACKPOT ★ " + player_name + " 獲得 " + str(reward) + " 金幣！", PANEL_COLOR, 5.0)
	show_panel()

	# 結算彈窗
	show_settle(
		"👑 GRAND JACKPOT！",
		player_name + "\n獲得 " + str(reward) + " 金幣！\n全場恭喜！",
		PANEL_COLOR
	)

	# 延遲隱藏
	var tween = create_tween()
	tween.tween_interval(4.0)
	tween.tween_callback(func(): hide_panel())
