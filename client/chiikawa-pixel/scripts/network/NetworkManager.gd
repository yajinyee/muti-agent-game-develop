## NetworkManager.gd
## WebSocket 連線管理，負責與 Go Server 通訊
## Autoload 單例

extends Node
class_name NetworkManagerClass

# 訊號
signal connected()
signal disconnected()
signal message_received(type: String, payload: Dictionary)
signal connection_error(error: String)
signal rooms_fetched(rooms: Array)  # 房間列表取得（DAY-020）
signal spectator_snapshot_received(snapshot: Dictionary)  # 觀戰快照（DAY-024）

# 設定
# 本機開發（直連 Game Server，無 TLS）
const SERVER_URL_LOCAL = "ws://localhost:7777/ws"
const SPECTATE_URL_LOCAL = "ws://localhost:7777/spectate"
const HTTP_URL_LOCAL = "http://localhost:7777"

# 生產環境（透過 Nginx TLS 反向代理，wss://）
# 自動偵測：HTTPS 頁面用 wss://，HTTP 頁面用 ws://
# 自動偵測：hostname 動態取得，不硬編碼 IP

# 自動判斷：Web 平台讀 hostname + protocol，桌面版用 localhost
var SERVER_URL: String:
	get:
		if OS.has_feature("web"):
			var host = JavaScriptBridge.eval("window.location.hostname")
			var protocol = JavaScriptBridge.eval("window.location.protocol")
			if host != "localhost" and host != "127.0.0.1" and host != "":
				# HTTPS 頁面用 wss://（安全），HTTP 頁面用 ws://
				var ws_scheme = "wss" if protocol == "https:" else "ws"
				return ws_scheme + "://" + host + "/ws"
		return SERVER_URL_LOCAL

var SPECTATE_URL: String:
	get:
		if OS.has_feature("web"):
			var host = JavaScriptBridge.eval("window.location.hostname")
			var protocol = JavaScriptBridge.eval("window.location.protocol")
			if host != "localhost" and host != "127.0.0.1" and host != "":
				var ws_scheme = "wss" if protocol == "https:" else "ws"
				return ws_scheme + "://" + host + "/spectate"
		return SPECTATE_URL_LOCAL

var HTTP_BASE: String:
	get:
		if OS.has_feature("web"):
			var host = JavaScriptBridge.eval("window.location.hostname")
			var protocol = JavaScriptBridge.eval("window.location.protocol")
			if host != "localhost" and host != "127.0.0.1" and host != "":
				# HTTPS 頁面用 https://，HTTP 頁面用 http://
				return protocol + "//" + host
		return HTTP_URL_LOCAL

const RECONNECT_DELAY_MIN = 1.0   # 最短重連延遲（秒）
const RECONNECT_DELAY_MAX = 30.0  # 最長重連延遲（秒，exponential backoff 上限）
const RECONNECT_JITTER = 0.5      # 隨機抖動範圍（±0.5秒，防止 thundering herd）
const PING_INTERVAL = 30.0

var _socket: WebSocketPeer
var _player_id: String = ""
var _room_id: String = "room-001"  # 當前房間 ID（DAY-020）
var _connected: bool = false
var _reconnect_timer: float = 0.0
var _reconnect_attempt: int = 0    # 重連嘗試次數（用於 exponential backoff）
var _reconnect_delay: float = 1.0  # 當前重連延遲（動態計算）
var _ping_timer: float = 0.0
var _is_spectator: bool = false  # 是否為觀戰模式（DAY-024）

# Ping 延遲計算（DAY-036）
var _ping_sent_at: float = 0.0   # 發送 ping 的時間（Time.get_ticks_msec）
var _last_ping_ms: int = -1       # 最後一次 ping 延遲（ms），-1 表示尚未測量

# HTTP 請求（用於查詢房間列表）
var _http_request: HTTPRequest = null
# HTTP 請求（用於查詢觀戰快照）
var _snapshot_request: HTTPRequest = null

func _ready() -> void:
	_player_id = _generate_player_id()
	_socket = WebSocketPeer.new()
	# 建立 HTTPRequest 節點（用於查詢房間列表）
	_http_request = HTTPRequest.new()
	add_child(_http_request)
	_http_request.request_completed.connect(_on_rooms_response)
	# 建立 HTTPRequest 節點（用於查詢觀戰快照，DAY-024）
	_snapshot_request = HTTPRequest.new()
	add_child(_snapshot_request)
	_snapshot_request.request_completed.connect(_on_snapshot_response)
	# 預設連線到 room-001（向後相容）
	# 大廳 UI 可以呼叫 connect_to_room() 切換房間
	connect_to_server()

func _process(delta: float) -> void:
	if _socket == null:
		return

	_socket.poll()
	var state = _socket.get_ready_state()

	match state:
		WebSocketPeer.STATE_OPEN:
			if not _connected:
				_connected = true
				_reconnect_attempt = 0   # 重置 backoff 計數器
				_reconnect_delay = RECONNECT_DELAY_MIN
				print("[Network] Connected to server")
				emit_signal("connected")

			# 處理收到的訊息
			while _socket.get_available_packet_count() > 0:
				var packet = _socket.get_packet()
				_handle_packet(packet)

			# Ping（帶時間戳，用於計算延遲）
			_ping_timer += delta
			if _ping_timer >= PING_INTERVAL:
				_ping_timer = 0.0
				_ping_sent_at = Time.get_ticks_msec()
				send("ping", {"t": int(_ping_sent_at)})

		WebSocketPeer.STATE_CLOSED:
			if _connected:
				_connected = false
				print("[Network] Disconnected from server")
				emit_signal("disconnected")

			# 自動重連（Exponential Backoff + Jitter，防止 thundering herd）
			_reconnect_timer += delta
			if _reconnect_timer >= _reconnect_delay:
				_reconnect_timer = 0.0
				_reconnect_attempt += 1
				# 計算下次延遲：min(base * 2^attempt, max) + jitter
				var base_delay = minf(RECONNECT_DELAY_MIN * pow(2.0, _reconnect_attempt - 1), RECONNECT_DELAY_MAX)
				_reconnect_delay = base_delay + randf_range(-RECONNECT_JITTER, RECONNECT_JITTER)
				_reconnect_delay = maxf(_reconnect_delay, RECONNECT_DELAY_MIN)
				print("[Network] Reconnecting (attempt %d, next delay %.1fs)..." % [_reconnect_attempt, _reconnect_delay])
				connect_to_server()

		WebSocketPeer.STATE_CONNECTING:
			pass

		WebSocketPeer.STATE_CLOSING:
			pass

## 連線到 Server
func connect_to_server() -> void:
	_is_spectator = false
	var url = SERVER_URL + "?player_id=" + _player_id + "&room_id=" + _room_id
	print("[Network] Connecting to: ", url)
	var err = _socket.connect_to_url(url)
	if err != OK:
		push_error("[Network] Connection failed: " + str(err))
		emit_signal("connection_error", "Connection failed: " + str(err))

## 設定房間並連線（由大廳呼叫，DAY-020）
func connect_to_room(room_id: String) -> void:
	_room_id = room_id
	connect_to_server()

## 以觀戰模式連線到指定房間（DAY-024）
## 觀戰者收到所有廣播，但無法發送遊戲指令
func spectate_room(room_id: String) -> void:
	_room_id = room_id
	_is_spectator = true
	var url = SPECTATE_URL + "?room_id=" + room_id
	print("[Network] Spectating room: ", room_id, " at ", url)
	var err = _socket.connect_to_url(url)
	if err != OK:
		push_error("[Network] Spectate connection failed: " + str(err))
		emit_signal("connection_error", "Spectate connection failed: " + str(err))

## 查詢觀戰快照（HTTP GET /spectate/snapshot，DAY-024）
func fetch_spectator_snapshot() -> void:
	if not is_instance_valid(_snapshot_request):
		return
	var url = HTTP_BASE + "/spectate/snapshot"
	print("[Network] Fetching spectator snapshot from: ", url)
	var err = _snapshot_request.request(url)
	if err != OK:
		push_error("[Network] Snapshot request failed: " + str(err))

## 處理觀戰快照回應（DAY-024）
func _on_snapshot_response(_result: int, response_code: int, _headers: PackedStringArray, body: PackedByteArray) -> void:
	if response_code != 200:
		push_error("[Network] Snapshot fetch failed: HTTP " + str(response_code))
		return
	var text = body.get_string_from_utf8()
	var json = JSON.new()
	if json.parse(text) != OK:
		push_error("[Network] Snapshot JSON parse error")
		return
	var data = json.get_data()
	if data is Dictionary:
		emit_signal("spectator_snapshot_received", data)

## 是否為觀戰模式（DAY-024）
func is_spectator() -> bool:
	return _is_spectator

## 查詢房間列表（HTTP GET /rooms，DAY-020）
func fetch_rooms() -> void:
	if not is_instance_valid(_http_request):
		return
	var url = HTTP_BASE + "/rooms"
	print("[Network] Fetching rooms from: ", url)
	var err = _http_request.request(url)
	if err != OK:
		push_error("[Network] HTTP request failed: " + str(err))

## 處理房間列表回應
func _on_rooms_response(_result: int, response_code: int, _headers: PackedStringArray, body: PackedByteArray) -> void:
	if response_code != 200:
		push_error("[Network] Rooms fetch failed: HTTP " + str(response_code))
		emit_signal("rooms_fetched", [])
		return
	var text = body.get_string_from_utf8()
	var json = JSON.new()
	if json.parse(text) != OK:
		push_error("[Network] Rooms JSON parse error")
		emit_signal("rooms_fetched", [])
		return
	var data = json.get_data()
	if data is Array:
		emit_signal("rooms_fetched", data)
	else:
		emit_signal("rooms_fetched", [])

## 取得當前房間 ID
func get_room_id() -> String:
	return _room_id

## 傳送訊息
func send(type: String, payload: Dictionary) -> void:
	if not _connected:
		return
	var msg = {"type": type, "payload": payload}
	var json_str = JSON.stringify(msg)
	_socket.send_text(json_str)

## 攻擊
func send_attack(target_id: String, click_x: float, click_y: float) -> void:
	send("attack", {
		"target_id": target_id,
		"click_x": click_x,
		"click_y": click_y
	})

## 鎖定目標
func send_lock(target_id: String) -> void:
	send("lock", {"target_id": target_id})

## 切換自動攻擊
func send_auto_toggle() -> void:
	send("auto_toggle", {})

## 切換投注
func send_bet_change(bet_level: int) -> void:
	send("bet_change", {"bet_level": bet_level})

## Bonus 點擊
func send_bonus_click(target_id: String, click_x: float, click_y: float) -> void:
	send("bonus_click", {
		"target_id": target_id,
		"click_x": click_x,
		"click_y": click_y
	})

## 觸發 BOSS（Prototype 展示用）
func send_trigger_boss() -> void:
	send("trigger_boss", {})

## 觸發 Bonus（Prototype 展示用）
func send_trigger_bonus() -> void:
	send("trigger_bonus", {})

## 設定顯示名稱（DAY-021）
func send_set_display_name(display_name: String) -> void:
	send("set_display_name", {"display_name": display_name})

## 上報 Client 端效能數據（DAY-045）
## 每 30 秒由 PerformanceMonitor 呼叫，讓 Server 能監控玩家端效能
func send_perf_report(fps: float, memory_mb: float, draw_calls: int, node_count: int, ping_ms: int, quality: String) -> void:
	send("client_perf", {
		"fps": fps,
		"memory_mb": memory_mb,
		"draw_calls": draw_calls,
		"node_count": node_count,
		"ping_ms": ping_ms,
		"quality": quality,
		"timestamp": int(Time.get_unix_time_from_system() * 1000)
	})

## 購買特殊武器（DAY-089）
func send_buy_special_weapon(weapon_type: String) -> void:
	send("buy_special_weapon", {"weapon_type": weapon_type})

## 使用特殊武器（DAY-089）
func send_use_special_weapon(weapon_type: String, click_x: float, click_y: float) -> void:
	send("use_special_weapon", {
		"weapon_type": weapon_type,
		"click_x": click_x,
		"click_y": click_y
	})

## 查詢特殊武器狀態（DAY-089）
func send_get_special_weapons() -> void:
	send("get_special_weapons", {})

## 處理收到的封包
func _handle_packet(packet: PackedByteArray) -> void:
	var text = packet.get_string_from_utf8()
	var json = JSON.new()
	var err = json.parse(text)
	if err != OK:
		push_error("[Network] JSON parse error: " + text)
		return

	var data = json.get_data()
	if not data is Dictionary:
		return

	var type = data.get("type", "")
	var payload = data.get("payload", {})

	# 計算 ping 延遲（DAY-036）
	if type == "pong" and _ping_sent_at > 0:
		_last_ping_ms = int(Time.get_ticks_msec() - _ping_sent_at)
		_ping_sent_at = 0.0

	emit_signal("message_received", type, payload)

## 產生玩家 ID
func _generate_player_id() -> String:
	return "player_" + str(randi() % 999999).pad_zeros(6)

## 取得連線狀態
func is_connected_to_server() -> bool:
	return _connected

## 取得玩家 ID
func get_player_id() -> String:
	return _player_id

## 取得最後一次 ping 延遲（ms），-1 表示尚未測量（DAY-036）
func get_ping_ms() -> int:
	return _last_ping_ms
