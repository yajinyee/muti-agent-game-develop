## PerformanceMonitor.gd
## 監控遊戲效能，在低階設備自動降級以維持 30 FPS
## Autoload 單例

extends Node

# 效能等級
enum QualityLevel {
	HIGH,    # 全效果（預設）
	MEDIUM,  # 關閉部分特效
	LOW,     # 最低效果，確保 30 FPS
}

var current_quality: QualityLevel = QualityLevel.HIGH

# FPS 監控
var _fps_samples: Array[float] = []
const SAMPLE_COUNT = 60       # 取樣 60 幀（約 2 秒）
const CHECK_INTERVAL = 3.0    # 每 3 秒評估一次
var _check_timer: float = 0.0
var _degraded_count: int = 0  # 連續降級次數
var _upgraded_count: int = 0  # 連續升級次數

# 效能快照（供 HUD 讀取）
var snapshot_fps: float = 60.0
var snapshot_memory_mb: float = 0.0
var snapshot_draw_calls: int = 0
var snapshot_objects: int = 0
var snapshot_nodes: int = 0

# 效能門檻
const FPS_LOW_THRESHOLD = 28.0    # 低於此值降級
const FPS_HIGH_THRESHOLD = 50.0   # 高於此值升級
const DEGRADE_TRIGGER = 2         # 連續 N 次低 FPS 才降級（避免瞬間抖動）
const UPGRADE_TRIGGER = 5         # 連續 N 次高 FPS 才升級（保守升級）

# 效能設定
var _swim_animation_enabled: bool = true
var _screen_shake_enabled: bool = true
var _particle_count_scale: float = 1.0
var _outline_shader_enabled: bool = true

# Server 效能上報（DAY-045）
const PERF_REPORT_INTERVAL = 30.0  # 每 30 秒上報一次
var _perf_report_timer: float = 0.0

signal quality_changed(new_level: QualityLevel)

func _ready() -> void:
	# 初始化 FPS 樣本
	_fps_samples.resize(SAMPLE_COUNT)
	_fps_samples.fill(60.0)
	print("[PerfMon] Performance monitor started")

	# 加入 Godot Debugger 自訂監控器（DAY-045，只在 debug build 有效）
	# 讓開發時可以在 Debugger → Monitor 面板看到遊戲自訂指標
	if OS.is_debug_build():
		Performance.add_custom_monitor("game/fps_avg", get_current_fps)
		Performance.add_custom_monitor("game/memory_mb", func(): return snapshot_memory_mb)
		Performance.add_custom_monitor("game/draw_calls", func(): return float(snapshot_draw_calls))
		Performance.add_custom_monitor("game/node_count", func(): return float(snapshot_nodes))
		print("[PerfMon] Custom monitors registered in Debugger")

func _process(delta: float) -> void:
	# 收集 FPS 樣本
	var fps = Engine.get_frames_per_second()
	_fps_samples.push_back(fps)
	if _fps_samples.size() > SAMPLE_COUNT:
		_fps_samples.pop_front()

	# 更新效能快照（每幀）
	snapshot_fps = fps
	snapshot_memory_mb = Performance.get_monitor(Performance.MEMORY_STATIC) / 1048576.0
	snapshot_draw_calls = int(Performance.get_monitor(Performance.RENDER_TOTAL_DRAW_CALLS_IN_FRAME))
	snapshot_objects = int(Performance.get_monitor(Performance.OBJECT_COUNT))
	snapshot_nodes = int(Performance.get_monitor(Performance.OBJECT_NODE_COUNT))

	# 定期評估
	_check_timer += delta
	if _check_timer >= CHECK_INTERVAL:
		_check_timer = 0.0
		_evaluate_performance()

	# Server 效能上報（DAY-045）
	# 每 30 秒上報一次，讓 Server 能監控玩家端效能
	_perf_report_timer += delta
	if _perf_report_timer >= PERF_REPORT_INTERVAL:
		_perf_report_timer = 0.0
		_send_perf_report()

func _evaluate_performance() -> void:
	if _fps_samples.is_empty():
		return

	# 計算平均 FPS
	var avg_fps = 0.0
	for s in _fps_samples:
		avg_fps += s
	avg_fps /= _fps_samples.size()

	# 計算 P5 FPS（最差 5% 的幀）
	var sorted = _fps_samples.duplicate()
	sorted.sort()
	var p5_fps = sorted[int(sorted.size() * 0.05)]

	# 決定是否降級
	if p5_fps < FPS_LOW_THRESHOLD:
		_degraded_count += 1
		_upgraded_count = 0
		if _degraded_count >= DEGRADE_TRIGGER:
			_degrade_quality()
			_degraded_count = 0
	elif avg_fps > FPS_HIGH_THRESHOLD:
		_upgraded_count += 1
		_degraded_count = 0
		if _upgraded_count >= UPGRADE_TRIGGER:
			_upgrade_quality()
			_upgraded_count = 0

	# Debug 輸出（只在 DEBUG 模式）
	if OS.is_debug_build():
		print("[PerfMon] FPS avg=%.1f p5=%.1f quality=%s" % [avg_fps, p5_fps, _quality_name()])

func _degrade_quality() -> void:
	match current_quality:
		QualityLevel.HIGH:
			_set_quality(QualityLevel.MEDIUM)
		QualityLevel.MEDIUM:
			_set_quality(QualityLevel.LOW)
		QualityLevel.LOW:
			pass  # 已是最低，無法再降

func _upgrade_quality() -> void:
	match current_quality:
		QualityLevel.LOW:
			_set_quality(QualityLevel.MEDIUM)
		QualityLevel.MEDIUM:
			_set_quality(QualityLevel.HIGH)
		QualityLevel.HIGH:
			pass  # 已是最高

func _set_quality(level: QualityLevel) -> void:
	if level == current_quality:
		return

	var old = current_quality
	current_quality = level
	print("[PerfMon] Quality: %s → %s" % [_quality_name_for(old), _quality_name()])

	match level:
		QualityLevel.HIGH:
			_swim_animation_enabled = true
			_screen_shake_enabled = true
			_particle_count_scale = 1.0
			_outline_shader_enabled = true
			Engine.max_fps = 0  # 不限制 FPS

		QualityLevel.MEDIUM:
			_swim_animation_enabled = true
			_screen_shake_enabled = true
			_particle_count_scale = 0.5  # 粒子數量減半
			_outline_shader_enabled = true
			Engine.max_fps = 60

		QualityLevel.LOW:
			_swim_animation_enabled = false  # 關閉游泳動畫（最耗效能的 tween）
			_screen_shake_enabled = false    # 關閉畫面震動
			_particle_count_scale = 0.25     # 粒子數量 1/4
			_outline_shader_enabled = false  # 關閉 outline shader
			Engine.max_fps = 30              # 鎖定 30 FPS（穩定優先）

	emit_signal("quality_changed", level)

## 查詢介面（供其他系統使用）

func is_swim_animation_enabled() -> bool:
	return _swim_animation_enabled

func is_screen_shake_enabled() -> bool:
	return _screen_shake_enabled

func get_particle_count_scale() -> float:
	return _particle_count_scale

func is_outline_shader_enabled() -> bool:
	return _outline_shader_enabled

func get_current_fps() -> float:
	if _fps_samples.is_empty():
		return 60.0
	var sum = 0.0
	for s in _fps_samples:
		sum += s
	return sum / _fps_samples.size()

## get_bullet_pool_stats — 取得 BulletPool 統計（供 HUD 顯示）
func get_bullet_pool_stats() -> Dictionary:
	if Engine.has_singleton("BulletPool") or is_instance_valid(get_node_or_null("/root/BulletPool")):
		var pool = get_node_or_null("/root/BulletPool")
		if is_instance_valid(pool):
			return pool.get_stats()
	return {"active": 0, "pooled": 0, "total": 0}

func _quality_name() -> String:
	return _quality_name_for(current_quality)

func _quality_name_for(level: QualityLevel) -> String:
	match level:
		QualityLevel.HIGH: return "HIGH"
		QualityLevel.MEDIUM: return "MEDIUM"
		QualityLevel.LOW: return "LOW"
	return "UNKNOWN"

## _send_perf_report — 上報效能數據到 Server（DAY-045）
## 每 30 秒自動呼叫，讓 Grafana 能看到玩家端效能狀況
func _send_perf_report() -> void:
	# 確認 NetworkManager 已連線
	var nm = get_node_or_null("/root/NetworkManager")
	if not is_instance_valid(nm):
		return
	if not nm.is_connected_to_server():
		return

	var ping_ms = nm.get_ping_ms()
	if ping_ms < 0:
		ping_ms = 0  # 尚未測量時用 0

	nm.send_perf_report(
		snapshot_fps,
		snapshot_memory_mb,
		snapshot_draw_calls,
		snapshot_nodes,
		ping_ms,
		_quality_name()
	)
