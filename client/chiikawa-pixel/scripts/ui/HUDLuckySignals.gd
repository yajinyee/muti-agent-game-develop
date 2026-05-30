## HUDLuckySignals.gd — HUD Lucky 訊號連接模組
## hud-core-agent 負責維護
## DAY-335：從 HUD.gd 拆分出來，解決 2330+ 行技術債
## 職責：連接所有 148 個 Lucky 訊號，並委派給 LuckyPanelRegistry 或 LuckyEventSystem 處理
extends Node

## 依賴注入（由 HUD 在 _ready 中設定）
var hud: Node = null
var lucky_event_system: Node = null

## 連接所有 Lucky 訊號（由 HUD._ready 呼叫）
func connect_all_lucky_signals(hud_node: Node) -> void:
	hud = hud_node
	# DAY-292
	GameManager.lucky_chain_lightning.connect(_on_lucky_chain_lightning)
	GameManager.lucky_crab_torpedo.connect(_on_lucky_crab_torpedo)
	GameManager.lucky_vortex.connect(_on_lucky_vortex)
	GameManager.lucky_golden_dragon.connect(_on_lucky_golden_dragon)
	GameManager.lucky_thunder_lobster.connect(_on_lucky_thunder_lobster)
	GameManager.announce.connect(_on_announce)
	# DAY-293
	GameManager.lucky_awakened_phoenix.connect(_on_lucky_awakened_phoenix)
	GameManager.lucky_shockwave_bomb.connect(_on_lucky_shockwave_bomb)
	# DAY-294
	GameManager.lucky_drill_torpedo.connect(_on_lucky_drill_torpedo)
	GameManager.lucky_time_freeze.connect(_on_lucky_time_freeze)
	GameManager.lucky_chain_explosion.connect(_on_lucky_chain_explosion)
	# DAY-295
	GameManager.lucky_chain_long_king.connect(_on_lucky_chain_long_king)
	GameManager.lucky_dragon_shotgun.connect(_on_lucky_dragon_shotgun)
	GameManager.lucky_rocket_cannon.connect(_on_lucky_rocket_cannon)
	GameManager.lucky_deep_whirlpool.connect(_on_lucky_deep_whirlpool)
	GameManager.lucky_vampire_mult.connect(_on_lucky_vampire_mult)
	# DAY-296
	GameManager.lucky_mirror_fish.connect(_on_lucky_mirror_fish)
	GameManager.lucky_golden_rain.connect(_on_lucky_golden_rain)
	GameManager.lucky_freeze_bomb.connect(_on_lucky_freeze_bomb)
	GameManager.lucky_thunder_storm.connect(_on_lucky_thunder_storm)
	GameManager.lucky_lucky_wheel.connect(_on_lucky_lucky_wheel)
	# DAY-301
	GameManager.lucky_jackpot_fish.connect(_on_lucky_jackpot_fish)
	GameManager.lucky_coop_fish.connect(_on_lucky_coop_fish)
	GameManager.lucky_time_warp.connect(_on_lucky_time_warp)
	# DAY-302
	GameManager.lucky_chain_meteor.connect(_on_lucky_chain_meteor)
	# DAY-303
	GameManager.lucky_crash_fish.connect(_on_lucky_crash_fish)
	# DAY-304（Panel 自行連接，HUD 只做備用橫幅）
	GameManager.lucky_electric_eel.connect(_on_lucky_fallback.bind("electric_eel", "⚡ 電鰻觸發！"))
	GameManager.lucky_angler_fish.connect(_on_lucky_fallback.bind("angler_fish", "🐟 安康魚觸發！"))
	GameManager.lucky_black_hole.connect(_on_lucky_fallback.bind("black_hole", "🌑 黑洞觸發！"))
	GameManager.lucky_bounty_hunter.connect(_on_lucky_fallback.bind("bounty_hunter", "🎯 賞金獵人觸發！"))
	GameManager.lucky_tsunami.connect(_on_lucky_fallback.bind("tsunami", "🌊 海嘯觸發！"))
	print("[HUDLuckySignals] DAY-292~304 訊號連接完成")

## 委派給 LuckyEventSystem 的輔助函數
func _show_banner(text: String, color: Color, duration: float = 2.5) -> void:
	if is_instance_valid(lucky_event_system):
		lucky_event_system.show_banner(text, color, duration)
	elif is_instance_valid(hud) and hud.has_method("_show_fallback_banner"):
		hud._show_fallback_banner(text, color, duration)

func _show_event(key: String, msg: String, duration: float = 2.5) -> void:
	if is_instance_valid(lucky_event_system):
		lucky_event_system.show_lucky_banner(key, msg, duration)
	else:
		_show_banner(msg, Color.WHITE, duration)

func _update_indicator(title: String, value: String, bar_pct: float = -1.0, color: Color = Color(1.0, 0.85, 0.0)) -> void:
	if is_instance_valid(lucky_event_system):
		lucky_event_system.update_indicator(title, value, bar_pct, color)

func _hide_indicator() -> void:
	if is_instance_valid(lucky_event_system):
		lucky_event_system.hide_indicator()

func _show_settle(lines: Array, duration: float = 3.5) -> void:
	if is_instance_valid(lucky_event_system):
		lucky_event_system.show_settle(lines, duration)

func _show_reward(amount: int, mult: float) -> void:
	if is_instance_valid(hud) and hud.has_method("_show_reward_popup"):
		hud._show_reward_popup(amount, mult)

## 通用備用橫幅（Panel 自行處理的 Lucky 魚，HUD 只顯示簡單提示）
func _on_lucky_fallback(data: Dictionary, _key: String, msg: String) -> void:
	var event = data.get("event", "")
	if event == "trigger":
		var name = data.get("trigger_name", "玩家")
		_show_banner("%s %s" % [msg, name], Color(1.0, 0.85, 0.0))
		AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)

# ── DAY-292 Lucky 事件處理 ────────────────────────────────────

func _on_lucky_chain_lightning(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			_show_event("chain_lightning", "⚡ %s 觸發連鎖閃電！" % name)
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.4)
		"chain_hit":
			var chain = data.get("chain_count", 0)
			var mult = data.get("multiplier", 1.0)
			_show_banner("⚡ 連鎖 %d！×%.1f" % [chain, mult], Color(0.0, 0.9, 1.0), 1.0)
		"settle":
			var reward = data.get("total_reward", 0)
			if reward > 0:
				_show_reward(reward, data.get("multiplier", 1.0))

func _on_lucky_crab_torpedo(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			_show_event("crab_torpedo", "🦀 %s 發射螃蟹魚雷！" % name)
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.35)
		"explosion":
			var no = data.get("explosion_no", 1)
			_show_banner("💥 魚雷爆炸 %d/3！" % no, Color(1.0, 0.6, 0.2), 0.8)
			ScreenShake.add_trauma(0.5)
		"settle":
			var reward = data.get("total_reward", 0)
			if reward > 0:
				_show_reward(reward, 3.0)

func _on_lucky_vortex(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			_show_event("vortex", "🌀 %s 召喚渦旋海葵！" % name)
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
		"pull":
			var tl = data.get("time_left", 0.0)
			_update_indicator("🌀 渦旋海葵", "%.0fs" % tl, tl / 8.0, Color(0.7, 0.4, 1.0))
		"end":
			_hide_indicator()
			_show_banner("🌀 渦旋爆炸！全場 HP -20%！", Color(0.8, 0.5, 1.0))
			ScreenShake.add_trauma(0.6)

func _on_lucky_golden_dragon(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			_show_event("golden_dragon", "🐉 %s 觸發黃金龍魚輪盤！" % name)
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
		"spin":
			var inner = data.get("inner_mult", 1.0)
			var outer = data.get("outer_mult", 1.0)
			var final_m = data.get("final_mult", 1.0)
			_show_banner("🐉 內環 ×%.0f × 外環 ×%.0f = ×%.0f！" % [inner, outer, final_m], Color(1.0, 0.85, 0.0), 3.0)
		"result":
			var reward = data.get("reward", 0)
			var final_m = data.get("final_mult", 1.0)
			if reward > 0:
				_show_reward(reward, final_m)
			if final_m >= 100:
				ScreenShake.add_trauma(0.8)

func _on_lucky_thunder_lobster(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			_show_event("thunder_lobster", "🦞⚡ %s 觸發雷霆龍蝦！15 秒免費射擊！" % name)
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
		"auto_fire":
			var tl = data.get("time_left", 0.0)
			var kills = data.get("kill_count", 0)
			_update_indicator("🦞⚡ 雷霆模式", "%.0fs | %d 條" % [tl, kills], tl / 15.0, Color(1.0, 0.5, 0.2))
		"end":
			_hide_indicator()
			var reward = data.get("total_reward", 0)
			var kills = data.get("kill_count", 0)
			_show_banner("🦞 雷霆結束！擊破 %d 條，獎勵 %d！" % [kills, reward], Color(1.0, 0.7, 0.3))
			if reward > 0:
				_show_reward(reward, float(kills))

func _on_announce(data: Dictionary) -> void:
	var msg = data.get("message", "")
	var color_str = data.get("color", "#FFFFFF")
	var color = Color.WHITE
	if color_str.begins_with("#") and color_str.length() == 7:
		var r = color_str.substr(1, 2).hex_to_int() / 255.0
		var g = color_str.substr(3, 2).hex_to_int() / 255.0
		var b = color_str.substr(5, 2).hex_to_int() / 255.0
		color = Color(r, g, b)
	var priority = data.get("priority", "normal")
	var duration = 2.0
	match priority:
		"high": duration = 3.0
		"critical": duration = 4.0
	_show_banner(msg, color, duration)

# ── DAY-293 ──────────────────────────────────────────────────

func _on_lucky_awakened_phoenix(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"awaken_start":
			_show_event("awakened_phoenix", "🔥 %s 觸發覺醒鳳凰！下 5 次攻擊 Power Up！" % name)
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.35)
		"power_up":
			var mult = data.get("power_up_mult", 6.0)
			var shots = data.get("shots_left", 0)
			_update_indicator("🔥 覺醒鳳凰", "×%.0f | 剩 %d 次" % [mult, shots], float(shots) / 5.0, Color(1.0, 0.6, 0.2))
		"perfect_awaken":
			_hide_indicator()
			_show_banner("🔥✨ 完美覺醒！%s 全服 ×2.0 加成 8 秒！" % name, Color(1.0, 0.85, 0.0), 3.5)
			ScreenShake.add_trauma(0.6)
		"awaken_end":
			_hide_indicator()
			var reward = data.get("total_reward", 0)
			var hits = data.get("hit_count", 0)
			_show_banner("🔥 覺醒結束！命中 %d 次，獎勵 %d！" % [hits, reward], Color(1.0, 0.7, 0.3))
			if reward > 0:
				_show_reward(reward, float(hits))

func _on_lucky_shockwave_bomb(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"shockwave_start":
			_show_event("shockwave_bomb", "💥 %s 觸發全場震盪！全場 HP -35%！" % name)
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.7)
		"shockwave_hit":
			var hits = data.get("hit_count", 0)
			var reward = data.get("total_reward", 0)
			_show_banner("💥 震盪命中 %d 個目標！獎勵 %d！" % [hits, reward], Color(1.0, 0.5, 0.2))
			if reward > 0:
				_show_reward(reward, float(hits) * 0.5)
		"super_shockwave":
			_show_banner("💥🌊 超級震盪！%s 全服 ×1.8 加成 6 秒！" % name, Color(1.0, 0.42, 0.21), 3.5)
			ScreenShake.add_trauma(0.8)

# ── DAY-294 ──────────────────────────────────────────────────

func _on_lucky_drill_torpedo(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			_show_event("drill_torpedo", "🚀 %s 發射鑽頭魚雷！穿透最多 5 個目標！" % name)
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.4)
		"penetrate":
			var cnt = data.get("penetrate_cnt", 0)
			var mult = data.get("accum_mult", 1.0)
			_update_indicator("🚀 鑽頭魚雷", "穿透 %d | ×%.1f" % [cnt, mult], float(cnt) / 5.0, Color(1.0, 0.55, 0.15))
		"explode":
			_hide_indicator()
			_show_banner("💥 魚雷終點爆炸！AOE 傷害！", Color(1.0, 0.4, 0.1))
			ScreenShake.add_trauma(0.55)
		"perfect":
			_show_banner("🚀💥 完美穿透！%s 全服 ×2.2 加成 6 秒！" % name, Color(1.0, 0.85, 0.0), 3.5)
			ScreenShake.add_trauma(0.65)

func _on_lucky_time_freeze(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"freeze_start":
			_show_event("time_freeze", "❄️ %s 觸發時間凍結！全場凍結 8 秒！傷害 ×1.8！" % name)
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.3)
		"freeze_end":
			_hide_indicator()
			_show_banner("❄️💥 冰裂爆炸！全場 HP -25%！", Color(0.6, 0.9, 1.0))
			ScreenShake.add_trauma(0.5)
		"perfect_freeze":
			var kills = data.get("kill_count", 0)
			_show_banner("❄️✨ 完美凍結！%s 擊破 %d 條！全服 ×2.0 加成 5 秒！" % [name, kills], Color(0.0, 0.9, 1.0), 3.5)
			ScreenShake.add_trauma(0.6)

func _on_lucky_chain_explosion(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"chain_start":
			_show_event("chain_explosion", "💥 %s 觸發連鎖爆炸！12 秒連鎖模式！" % name)
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.4)
		"chain_explode":
			var cnt = data.get("chain_count", 0)
			var mult = data.get("accum_mult", 1.0)
			_update_indicator("💥 連鎖爆炸", "×%.1f | %d 次" % [mult, cnt], -1.0, Color(0.9, 0.2, 0.15))
			ScreenShake.add_trauma(0.25)
		"chain_burst":
			_hide_indicator()
			_show_banner("💥🔥 連鎖爆發！%s 全服 ×2.5 加成 6 秒！" % name, Color(1.0, 0.85, 0.0), 3.5)
			ScreenShake.add_trauma(0.7)
		"chain_end":
			_hide_indicator()
			var cnt = data.get("chain_count", 0)
			var reward = data.get("total_reward", 0)
			_show_banner("💥 連鎖結束！%d 次連鎖，獎勵 %d！" % [cnt, reward], Color(1.0, 0.6, 0.3))
			if reward > 0:
				_show_reward(reward, float(cnt) * 0.5)

# ── DAY-295 ──────────────────────────────────────────────────

func _on_lucky_chain_long_king(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			_show_event("chain_long_king", "🐉👑 %s 觸發千龍王輪盤！最高 1000x！" % name)
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.5)
		"spin":
			var inner = data.get("inner_mult", 1.0)
			var outer = data.get("outer_mult", 1.0)
			var final_m = data.get("final_mult", 1.0)
			_show_banner("🐉 內環 ×%.0f × 外環 ×%.0f = ×%.0f！" % [inner, outer, final_m], Color(1.0, 0.85, 0.0), 3.0)
		"result":
			var reward = data.get("reward", 0)
			var final_m = data.get("final_mult", 1.0)
			if reward > 0:
				_show_reward(reward, final_m)
			if final_m >= 200:
				ScreenShake.add_trauma(0.7)
		"mega_win":
			var final_m = data.get("final_mult", 1.0)
			_show_banner("🐉👑✨ MEGA WIN！×%.0f！千龍王降臨！" % final_m, Color(1.0, 0.85, 0.0), 4.0)
			ScreenShake.add_trauma(1.0)

func _on_lucky_dragon_shotgun(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			_show_event("dragon_shotgun", "🐲💥 %s 觸發龍力散彈！8 方向攻擊！" % name)
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.45)
		"shotgun_fire":
			var dir = data.get("direction", 0)
			var hits = data.get("total_hits", 0)
			_show_banner("🐲 方向 %d 命中！總計 %d 個！" % [dir + 1, hits], Color(0.9, 0.4, 1.0), 0.6)
			ScreenShake.add_trauma(0.2)
		"settle":
			var reward = data.get("total_reward", 0)
			var hits = data.get("total_hits", 0)
			_show_banner("🐲 散彈結束！命中 %d 個，獎勵 %d！" % [hits, reward], Color(0.8, 0.5, 1.0))
			if reward > 0:
				_show_reward(reward, float(hits) * 0.4)

func _on_lucky_rocket_cannon(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			_show_event("rocket_cannon", "🚀💥 %s 召喚火箭砲！3 枚火箭！" % name)
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.4)
		"rocket_launch":
			var no = data.get("rocket_no", 1)
			_show_banner("🚀 第 %d 枚火箭發射！" % no, Color(1.0, 0.5, 0.2), 0.6)
		"rocket_explode":
			var no = data.get("rocket_no", 1)
			var hits = data.get("hit_targets", [])
			_show_banner("💥 火箭 %d 爆炸！命中 %d 個！" % [no, hits.size()], Color(1.0, 0.4, 0.1), 0.8)
			ScreenShake.add_trauma(0.5)
		"settle":
			var reward = data.get("total_reward", 0)
			_show_banner("🚀 火箭砲結束！獎勵 %d！" % reward, Color(1.0, 0.6, 0.3))
			if reward > 0:
				_show_reward(reward, 3.0)

func _on_lucky_deep_whirlpool(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			_show_event("deep_whirlpool", "🌊🌀 %s 觸發深海漩渦！全場 HP -50%！6 秒！" % name)
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.5)
		"whirlpool_damage":
			var hits = data.get("hit_count", 0)
			_show_banner("🌀 漩渦傷害！命中 %d 個！" % hits, Color(0.2, 0.7, 1.0), 0.7)
		"settle":
			var reward = data.get("total_reward", 0)
			_show_banner("🌊 深海漩渦結束！獎勵 %d！" % reward, Color(0.0, 0.8, 1.0))
			if reward > 0:
				_show_reward(reward, 5.0)

func _on_lucky_vampire_mult(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			_show_event("vampire_mult", "🧛 %s 觸發吸血鬼！每次擊破吸收倍率！" % name)
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.3)
		"absorb":
			var cnt = data.get("absorb_count", 0)
			var mult = data.get("current_mult", 1.0)
			_show_banner("🧛 吸收 %d 次！當前 ×%.1f" % [cnt, mult], Color(0.7, 0.2, 0.8), 0.7)
		"mult_mode":
			var mult = data.get("current_mult", 5.0)
			_show_banner("🧛✨ %s 進入倍率模式！×%.1f！10 秒！" % [name, mult], Color(0.8, 0.0, 0.8), 3.5)
			ScreenShake.add_trauma(0.6)
		"settle":
			var cnt = data.get("absorb_count", 0)
			var mult = data.get("current_mult", 1.0)
			_show_banner("🧛 吸血結束！吸收 %d 次，最終 ×%.1f！" % [cnt, mult], Color(0.6, 0.1, 0.7))

# ── DAY-296 ──────────────────────────────────────────────────

func _on_lucky_mirror_fish(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			var shots = data.get("shots_left", 3)
			_show_event("mirror_fish", "🪞 %s 觸發鏡像魚！下 %d 次攻擊自動複製！" % [name, shots])
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.3)
		"mirror_hit":
			var shots = data.get("shots_left", 0)
			var hits = data.get("hit_count", 0)
			_show_banner("🪞 鏡像命中！已命中 %d 次，剩餘 %d 次" % [hits, shots], Color(0.88, 0.67, 1.0), 0.8)
		"perfect_mirror":
			_show_banner("🪞✨ 完美鏡像！%s 全服 ×1.8 加成 5 秒！" % name, Color(0.88, 0.67, 1.0), 3.5)
			ScreenShake.add_trauma(0.5)
		"settle":
			var hits = data.get("hit_count", 0)
			var reward = data.get("total_reward", 0)
			_show_banner("🪞 鏡像結束！命中 %d 次，獎勵 %d！" % [hits, reward], Color(0.8, 0.6, 1.0))
			if reward > 0:
				_show_reward(reward, float(hits))

func _on_lucky_golden_rain(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			var total = data.get("total_coins", 10)
			_show_event("golden_rain", "🌧️💰 %s 觸發黃金雨！%d 個黃金幣！" % [name, total])
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.35)
		"coin_collect":
			var collected = data.get("collected_coins", 0)
			_show_banner("💰 收集 %d 個黃金幣！" % collected, Color(1.0, 0.85, 0.0), 0.6)
		"golden_harvest":
			var collected = data.get("collected_coins", 0)
			var reward = data.get("total_reward", 0)
			_show_banner("💰✨ 黃金豐收！%s 收集 %d 個！全服 ×2.0 加成 6 秒！" % [name, collected], Color(1.0, 0.85, 0.0), 3.5)
			ScreenShake.add_trauma(0.6)
			if reward > 0:
				_show_reward(reward, 2.0)
		"settle":
			var collected = data.get("collected_coins", 0)
			var reward = data.get("total_reward", 0)
			_show_banner("🌧️ 黃金雨結束！收集 %d 個，獎勵 %d！" % [collected, reward], Color(1.0, 0.85, 0.0))
			if reward > 0:
				_show_reward(reward, float(collected) * 0.3)

func _on_lucky_freeze_bomb(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"freeze_start":
			var frozen = data.get("frozen_targets", [])
			_show_event("freeze_bomb", "❄️💣 %s 投擲冰凍炸彈！%d 個目標凍結！" % [name, frozen.size()])
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.3)
		"bomb_explode":
			var hits = data.get("hit_count", 0)
			_show_banner("❄️💥 冰凍炸彈爆炸！命中 %d 個！HP -60%！" % hits, Color(0.4, 0.9, 1.0))
			ScreenShake.add_trauma(0.65)
		"perfect_freeze":
			var hits = data.get("hit_count", 0)
			_show_banner("❄️💥✨ 冰爆完美！%s 命中 %d 個！全服 ×2.2 加成 5 秒！" % [name, hits], Color(0.0, 0.9, 1.0), 3.5)
			ScreenShake.add_trauma(0.7)

func _on_lucky_thunder_storm(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"storm_start":
			var count = data.get("lightning_count", 6)
			_show_event("thunder_storm", "⛈️ %s 召喚雷暴！%d 道閃電！10 秒！" % [name, count])
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.4)
		"lightning_strike":
			var no = data.get("strike_no", 1)
			var hits = data.get("hit_targets", [])
			var mult = data.get("accum_mult", 1.0)
			_show_banner("⚡ 第 %d 道閃電！命中 %d 個！累積 ×%.1f" % [no, hits.size(), mult], Color(1.0, 0.9, 0.2), 0.7)
			if hits.size() > 0:
				ScreenShake.add_trauma(0.2)
		"perfect_storm":
			var strikes = data.get("hit_strikes", 6)
			_show_banner("⛈️✨ 雷暴完美！%s %d 道全命中！全服 ×2.3 加成 6 秒！" % [name, strikes], Color(1.0, 0.85, 0.0), 3.5)
			ScreenShake.add_trauma(0.7)
		"storm_end":
			var strikes = data.get("hit_strikes", 0)
			var mult = data.get("accum_mult", 1.0)
			var reward = data.get("total_reward", 0)
			_show_banner("⛈️ 雷暴結束！%d 道命中，累積 ×%.1f，獎勵 %d！" % [strikes, mult, reward], Color(1.0, 0.85, 0.0))
			if reward > 0:
				_show_reward(reward, mult)

func _on_lucky_lucky_wheel(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			var pool = data.get("pool_size", 20000)
			_show_event("lucky_wheel", "🎡 %s 觸發幸運大轉盤！大獎池 %d！" % [name, pool])
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.35)
		"spin_result":
			var slot_name = data.get("slot_name", "×2")
			var slot_type = data.get("slot_type", "mult")
			var reward = data.get("reward", 0)
			match slot_type:
				"jackpot":
					_show_banner("🎡🏆 %s 中大獎！%s！獎勵 %d！" % [name, slot_name, reward], Color(1.0, 0.85, 0.0), 4.0)
					ScreenShake.add_trauma(0.8)
					if reward > 0:
						_show_reward(reward, 100.0)
				"aoe":
					_show_banner("🎡💥 %s 轉到 %s！全場 HP -50%！" % [name, slot_name], Color(1.0, 0.42, 0.71))
					ScreenShake.add_trauma(0.5)
				"mult":
					var mult = data.get("slot_mult", 2.0)
					_show_banner("🎡 %s 轉到 %s！獎勵 %d！" % [name, slot_name, reward], Color(1.0, 0.42, 0.71))
					if reward > 0:
						_show_reward(reward, mult)

# ── DAY-301 ──────────────────────────────────────────────────

func _on_lucky_jackpot_fish(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"trigger":
			_show_event("jackpot_fish", "🏆 %s 觸發進階 Jackpot！" % name)
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.4)
		"jackpot_result":
			var tier = data.get("tier_name", "Mini")
			var reward = data.get("reward", 0)
			var tier_idx = data.get("tier_idx", 0)
			var colors = [Color(0.7, 0.4, 0.2), Color(0.8, 0.8, 0.9), Color(1.0, 0.55, 0.0), Color(1.0, 0.85, 0.0)]
			_show_banner("🏆 %s 中 %s Jackpot！獲得 %d！" % [name, tier, reward], colors[clamp(tier_idx, 0, 3)], 3.0)
			if reward > 0:
				_show_reward(reward, float(tier_idx + 1) * 10.0)
			if tier_idx == 3:
				ScreenShake.add_trauma(0.8)
		"grand_boost":
			var mult = data.get("boost_mult", 3.0)
			var secs = data.get("boost_secs", 10)
			_show_banner("🏆✨ GRAND JACKPOT！%s 全服 ×%.0f 加成 %d 秒！" % [name, mult, secs], Color(1.0, 0.85, 0.0), 4.0)
			ScreenShake.add_trauma(0.9)

func _on_lucky_coop_fish(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"coop_start":
			var target = data.get("target_points", 8)
			_show_event("coop_fish", "🤝 %s 發起全服合作！目標 %d 點！20 秒！" % [name, target])
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.3)
		"coop_progress":
			var current = data.get("current_points", 0)
			var target = data.get("target_points", 8)
			var tl = data.get("time_left", 0.0)
			_update_indicator("🤝 全服合作", "%d/%d 點" % [current, target], float(current) / float(max(target, 1)), Color(0.0, 0.9, 1.0))
		"coop_success":
			_hide_indicator()
			var boost = data.get("boost_mult", 4.0)
			var secs = data.get("boost_secs", 8)
			_show_banner("🤝✨ 全服合作成功！全服 ×%.0f 加成 %d 秒！" % [boost, secs], Color(0.0, 1.0, 0.5), 3.5)
			ScreenShake.add_trauma(0.7)
		"coop_timeout":
			_hide_indicator()
			var current = data.get("current_points", 0)
			var target = data.get("target_points", 8)
			_show_banner("🤝 合作挑戰時間到！達成 %d/%d 點" % [current, target], Color(0.6, 0.6, 0.6), 2.0)

func _on_lucky_time_warp(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"warp_start":
			var duration = data.get("duration", 10.0)
			var dmg = data.get("damage_mult", 2.0)
			_show_event("time_warp", "⏰ %s 觸發時間扭曲！全場慢速 %.0f 秒！傷害 ×%.0f！" % [name, duration, dmg])
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.35)
		"warp_end":
			_hide_indicator()
			var kills = data.get("kill_count", 0)
			_show_banner("⏰💥 時間扭曲結束！全場 HP -20%！擊破 %d 條！" % kills, Color(0.55, 0.2, 0.86))
			ScreenShake.add_trauma(0.5)
		"time_collapse":
			var kills = data.get("kill_count", 0)
			var boost = data.get("boost_mult", 2.5)
			var secs = data.get("boost_secs", 6)
			_show_banner("⏰💥 時間崩潰！%s 擊破 %d 條！全服 ×%.0f 加成 %d 秒！" % [name, kills, boost, secs], Color(1.0, 0.85, 0.0), 3.5)
			ScreenShake.add_trauma(0.7)

# ── DAY-302 ──────────────────────────────────────────────────

func _on_lucky_chain_meteor(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("player_name", "玩家")
	match event:
		"meteor_start":
			_show_event("chain_meteor", "☄️ %s 觸發連鎖隕石！" % name)
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.5)
		"meteor_hit":
			var no = data.get("meteor_no", 1)
			var hits = data.get("hit_count", 0)
			_show_banner("☄️ 隕石 %d 命中 %d 個！" % [no, hits], Color(1.0, 0.6, 0.2), 0.7)
			ScreenShake.add_trauma(0.3)
		"settle":
			var reward = data.get("total_reward", 0)
			_show_banner("☄️ 連鎖隕石結束！獎勵 %d！" % reward, Color(1.0, 0.7, 0.3))
			if reward > 0:
				_show_reward(reward, 5.0)

# ── DAY-303 ──────────────────────────────────────────────────

func _on_lucky_crash_fish(data: Dictionary) -> void:
	var event = data.get("event", "")
	var name = data.get("trigger_name", "玩家")
	match event:
		"crash_start":
			_show_event("crash_fish", "💥 %s 觸發崩潰魚！倍率持續上升！" % name)
			AudioManager.play_sfx(AudioManager.SFX.BIG_WIN)
			ScreenShake.add_trauma(0.3)
		"crash_update":
			var mult = data.get("current_mult", 1.0)
			_update_indicator("💥 崩潰倍率", "×%.1f" % mult, -1.0, Color(1.0, 0.3, 0.1))
		"crash_end":
			_hide_indicator()
			var mult = data.get("final_mult", 1.0)
			var reward = data.get("total_reward", 0)
			_show_banner("💥 崩潰！最終 ×%.1f，獎勵 %d！" % [mult, reward], Color(1.0, 0.4, 0.1))
			if reward > 0:
				_show_reward(reward, mult)
