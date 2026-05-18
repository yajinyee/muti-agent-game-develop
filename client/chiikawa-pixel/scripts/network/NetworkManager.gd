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

# 設定
const SERVER_URL_LOCAL = "ws://localhost:7777/ws"
const SERVER_URL_REMOTE = "ws://220.137.205.22:7777/ws"
const HTTP_URL_LOCAL = "http://localhost:7777"
const HTTP_URL_REMOTE = "http://220.137.205.22:7777"

# 自動判斷：Web 平台讀 hostname，桌面版用 localhost
var SERVER_URL: String:
	get:
		if OS.has_feature("web"):
			var host = JavaScriptBridge.eval("window.location.hostname")
			if host != "localhost" and host != "127.0.0.1" and host != "":
				return SERVER_URL_REMOTE
		return SERVER_URL_LOCAL

var HTTP_BASE: String:
	get:
		if OS.has_feature("web"):
			var host = JavaScriptBridge.eval("window.location.hostname")
			if host != "localhost" and host != "127.0.0.1" and host != "":
				return HTTP_URL_REMOTE
		return HTTP_URL_LOCAL

const RECONNECT_DELAY = 3.0
const PING_INTERVAL = 30.0

var _socket: WebSocketPeer
var _player_id: String = ""
var _room_id: String = "room-001"  # 當前房間 ID（DAY-020）
var _connected: bool = false
var _reconnect_timer: float = 0.0
var _ping_timer: float = 0.0

# HTTP 請求（用於查詢房間列表）
var _http_request: HTTPRequest = null

func _ready() -> void:
	_player_id = _generate_player_id()
	_socket = WebSocketPeer.new()
	# 建立 HTTPRequest 節點（用於查詢房間列表）
	_http_request = HTTPRequest.new()
	add_child(_http_request)
	_http_request.request_completed.connect(_on_rooms_response)
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
				print("[Network] Connected to server")
				emit_signal("connected")

			# 處理收到的訊息
			while _socket.get_available_packet_count() > 0:
				var packet = _socket.get_packet()
				_handle_packet(packet)

			# Ping
			_ping_timer += delta
			if _ping_timer >= PING_INTERVAL:
				_ping_timer = 0.0
				send("ping", {})

		WebSocketPeer.STATE_CLOSED:
			if _connected:
				_connected = false
				print("[Network] Disconnected from server")
				emit_signal("disconnected")

			# 自動重連
			_reconnect_timer += delta
			if _reconnect_timer >= RECONNECT_DELAY:
				_reconnect_timer = 0.0
				connect_to_server()

		WebSocketPeer.STATE_CONNECTING:
			pass

		WebSocketPeer.STATE_CLOSING:
			pass

## 連線到 Server
func connect_to_server() -> void:
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
