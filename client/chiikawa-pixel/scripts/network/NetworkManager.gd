## NetworkManager.gd — WebSocket 連線管理
## network-agent 負責維護
extends Node

signal connected()
signal disconnected()
signal message_received(type: String, payload: Dictionary)

const SERVER_URL_LOCAL = "ws://localhost:7777/ws"
const RECONNECT_DELAY_MIN = 1.0
const RECONNECT_DELAY_MAX = 30.0
const PING_INTERVAL = 30.0

var _socket: WebSocketPeer
var _player_id: String = ""
var _connected: bool = false
var _reconnect_timer: float = 0.0
var _reconnect_attempt: int = 0
var _reconnect_delay: float = 1.0
var _ping_timer: float = 0.0

func _ready() -> void:
	_player_id = "player_%06d" % (randi() % 999999)
	_socket = WebSocketPeer.new()
	_connect()

func _connect() -> void:
	var url = _get_server_url() + "?player_id=" + _player_id
	print("[Network] Connecting to: ", url)
	_socket.connect_to_url(url)

func _get_server_url() -> String:
	if OS.has_feature("web"):
		var host = JavaScriptBridge.eval("window.location.hostname")
		var proto = JavaScriptBridge.eval("window.location.protocol")
		if host != "localhost" and host != "127.0.0.1" and host != "":
			var scheme = "wss" if proto == "https:" else "ws"
			return scheme + "://" + host + "/ws"
	return SERVER_URL_LOCAL

func _process(delta: float) -> void:
	if _socket == null:
		return
	_socket.poll()
	var state = _socket.get_ready_state()

	match state:
		WebSocketPeer.STATE_OPEN:
			if not _connected:
				_connected = true
				_reconnect_attempt = 0
				_reconnect_delay = RECONNECT_DELAY_MIN
				print("[Network] Connected!")
				emit_signal("connected")
			while _socket.get_available_packet_count() > 0:
				_handle_packet(_socket.get_packet())
			_ping_timer += delta
			if _ping_timer >= PING_INTERVAL:
				_ping_timer = 0.0
				send("ping", {})

		WebSocketPeer.STATE_CLOSED:
			if _connected:
				_connected = false
				print("[Network] Disconnected")
				emit_signal("disconnected")
			_reconnect_timer += delta
			if _reconnect_timer >= _reconnect_delay:
				_reconnect_timer = 0.0
				_reconnect_attempt += 1
				var base = minf(RECONNECT_DELAY_MIN * pow(2.0, _reconnect_attempt - 1), RECONNECT_DELAY_MAX)
				_reconnect_delay = base + randf_range(-0.5, 0.5)
				_reconnect_delay = maxf(_reconnect_delay, RECONNECT_DELAY_MIN)
				print("[Network] Reconnecting (attempt %d)..." % _reconnect_attempt)
				_connect()

func _handle_packet(packet: PackedByteArray) -> void:
	var text = packet.get_string_from_utf8()
	var json = JSON.new()
	if json.parse(text) != OK:
		return
	var data = json.get_data()
	if not data is Dictionary:
		return
	var type = data.get("type", "")
	var payload = data.get("payload", {})
	if not payload is Dictionary:
		payload = {}
	emit_signal("message_received", type, payload)

func send(type: String, payload: Dictionary) -> void:
	if not _connected:
		return
	var msg = JSON.stringify({"type": type, "payload": payload})
	_socket.send_text(msg)

func send_attack(target_id: String, x: float, y: float) -> void:
	send("attack", {"target_id": target_id, "click_x": x, "click_y": y})

func send_auto_toggle() -> void:
	send("auto_toggle", {})

func send_bet_change(level: int) -> void:
	send("bet_change", {"bet_level": level})

func send_lock(target_id: String) -> void:
	send("lock", {"target_id": target_id})

func send_bonus_click(target_id: String, x: float, y: float) -> void:
	send("bonus_click", {"target_id": target_id, "click_x": x, "click_y": y})

func send_trigger_boss() -> void:
	send("trigger_boss", {})

func send_trigger_bonus() -> void:
	send("trigger_bonus", {})

func send_collect_golden_coin(coin_id: int) -> void:
	send("collect_golden_coin", {"coin_id": coin_id})

func is_connected_to_server() -> bool:
	return _connected

func get_player_id() -> String:
	return _player_id
