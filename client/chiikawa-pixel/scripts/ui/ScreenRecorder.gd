## ScreenRecorder.gd — 遊戲內側錄系統
## 功能：玩家可在遊戲中點選開始/停止側錄，錄製畫面存為 WebM（HTML5）或 PNG 序列（桌面）
## 技術：HTML5 用 MediaRecorder API（透過 JavaScriptBridge），桌面用 Viewport 截圖序列
## 位置：CanvasLayer layer=200（最上層，不被其他 UI 遮擋）
extends CanvasLayer

# ---- 狀態 ----
var _is_recording: bool = false
var _record_button: Button = null
var _status_label: Label = null
var _elapsed_label: Label = null
var _elapsed: float = 0.0
var _blink_timer: float = 0.0

# ---- 桌面模式（非 HTML5）----
var _frame_buffer: Array = []          # 存 PNG bytes
var _frame_timer: float = 0.0
const FRAME_INTERVAL: float = 1.0 / 30.0  # 30fps
const MAX_FRAMES: int = 1800               # 最多 60 秒 × 30fps

# ---- HTML5 模式 ----
var _js_recorder = null   # JavaScript MediaRecorder 物件
var _is_html5: bool = false

func _ready() -> void:
	layer = 200
	_is_html5 = OS.get_name() == "Web"
	_build_ui()
	if _is_html5:
		_init_html5_recorder()

# ---- UI 建立 ----
func _build_ui() -> void:
	# 右上角浮動面板
	var panel = PanelContainer.new()
	panel.name = "RecorderPanel"
	panel.position = Vector2(1100, 8)
	panel.size = Vector2(160, 44)

	var style = StyleBoxFlat.new()
	style.bg_color = Color(0.05, 0.05, 0.05, 0.82)
	style.border_color = Color(0.4, 0.4, 0.4, 1.0)
	style.set_border_width_all(1)
	style.set_corner_radius_all(4)
	panel.add_theme_stylebox_override("panel", style)
	add_child(panel)

	var hbox = HBoxContainer.new()
	hbox.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	hbox.add_theme_constant_override("separation", 4)
	panel.add_child(hbox)

	# 錄製按鈕
	_record_button = Button.new()
	_record_button.text = "⏺ REC"
	_record_button.custom_minimum_size = Vector2(72, 36)
	_record_button.add_theme_font_size_override("font_size", 12)
	_record_button.pressed.connect(_on_record_pressed)
	_set_button_idle_style()
	hbox.add_child(_record_button)

	# 右側資訊欄
	var vbox = VBoxContainer.new()
	vbox.add_theme_constant_override("separation", 0)
	hbox.add_child(vbox)

	_status_label = Label.new()
	_status_label.text = "READY"
	_status_label.add_theme_font_size_override("font_size", 10)
	_status_label.modulate = Color(0.6, 0.6, 0.6)
	vbox.add_child(_status_label)

	_elapsed_label = Label.new()
	_elapsed_label.text = "00:00"
	_elapsed_label.add_theme_font_size_override("font_size", 11)
	_elapsed_label.modulate = Color(0.8, 0.8, 0.8)
	vbox.add_child(_elapsed_label)

func _set_button_idle_style() -> void:
	var s = StyleBoxFlat.new()
	s.bg_color = Color(0.18, 0.18, 0.18)
	s.set_corner_radius_all(3)
	_record_button.add_theme_stylebox_override("normal", s)
	_record_button.modulate = Color(0.85, 0.85, 0.85)

func _set_button_recording_style() -> void:
	var s = StyleBoxFlat.new()
	s.bg_color = Color(0.7, 0.1, 0.1)
	s.set_corner_radius_all(3)
	_record_button.add_theme_stylebox_override("normal", s)
	_record_button.modulate = Color(1.0, 1.0, 1.0)

# ---- 按鈕事件 ----
func _on_record_pressed() -> void:
	if _is_recording:
		_stop_recording()
	else:
		_start_recording()

# ---- 開始錄製 ----
func _start_recording() -> void:
	_is_recording = true
	_elapsed = 0.0
	_blink_timer = 0.0
	_frame_buffer.clear()
	_frame_timer = 0.0

	_record_button.text = "⏹ STOP"
	_set_button_recording_style()
	_status_label.text = "REC"
	_status_label.modulate = Color(1.0, 0.3, 0.3)

	if _is_html5:
		_start_html5_recording()
	else:
		_status_label.text = "REC (PNG)"

	print("[ScreenRecorder] 開始錄製")

# ---- 停止錄製 ----
func _stop_recording() -> void:
	_is_recording = false
	_record_button.text = "⏺ REC"
	_set_button_idle_style()
	_status_label.text = "SAVING..."
	_status_label.modulate = Color(1.0, 0.85, 0.2)

	if _is_html5:
		_stop_html5_recording()
	else:
		_save_png_sequence()

	print("[ScreenRecorder] 停止錄製，共 %d 幀" % _frame_buffer.size())

# ---- _process：計時 + 截圖（桌面模式）----
func _process(delta: float) -> void:
	if not _is_recording:
		return

	_elapsed += delta

	# 更新計時顯示
	var mins = int(_elapsed) / 60
	var secs = int(_elapsed) % 60
	_elapsed_label.text = "%02d:%02d" % [mins, secs]

	# 閃爍紅點
	_blink_timer += delta
	if _blink_timer >= 0.5:
		_blink_timer = 0.0
		if _status_label.modulate.a > 0.5:
			_status_label.modulate.a = 0.2
		else:
			_status_label.modulate.a = 1.0

	# 桌面模式：每幀截圖
	if not _is_html5:
		_frame_timer += delta
		if _frame_timer >= FRAME_INTERVAL:
			_frame_timer -= FRAME_INTERVAL
			_capture_frame()

	# 超過最大長度自動停止
	if _elapsed >= 60.0:
		_stop_recording()

# ---- 桌面模式：截圖 ----
func _capture_frame() -> void:
	if _frame_buffer.size() >= MAX_FRAMES:
		return
	var img: Image = get_viewport().get_texture().get_image()
	_frame_buffer.append(img.save_png_to_buffer())

# ---- 桌面模式：儲存 PNG 序列 ----
func _save_png_sequence() -> void:
	if _frame_buffer.is_empty():
		_status_label.text = "EMPTY"
		_status_label.modulate = Color(0.6, 0.6, 0.6)
		return

	# 儲存到 user:// 目錄
	var dir_path = "user://recordings/"
	DirAccess.make_dir_recursive_absolute(dir_path)
	var timestamp = Time.get_datetime_string_from_system().replace(":", "-").replace(" ", "_")
	var folder = dir_path + "rec_" + timestamp + "/"
	DirAccess.make_dir_recursive_absolute(folder)

	for i in range(_frame_buffer.size()):
		var path = folder + "frame_%04d.png" % i
		var fa = FileAccess.open(path, FileAccess.WRITE)
		if fa:
			fa.store_buffer(_frame_buffer[i])
			fa.close()

	var real_path = ProjectSettings.globalize_path(folder)
	_status_label.text = "SAVED"
	_status_label.modulate = Color(0.3, 1.0, 0.3)
	_elapsed_label.text = "%d frm" % _frame_buffer.size()
	print("[ScreenRecorder] 已儲存 %d 幀到 %s" % [_frame_buffer.size(), real_path])

	# 3 秒後恢復 READY 狀態
	await get_tree().create_timer(3.0).timeout
	_status_label.text = "READY"
	_status_label.modulate = Color(0.6, 0.6, 0.6)
	_elapsed_label.text = "00:00"
	_frame_buffer.clear()

# ---- HTML5 模式：初始化 MediaRecorder ----
func _init_html5_recorder() -> void:
	# 注入 JavaScript 輔助函數
	# MediaRecorder 錄製 canvas 串流，存為 WebM
	var js_code = """
window._kiroRecorder = null;
window._kiroChunks = [];

window.kiroStartRecording = function() {
	try {
		var canvas = document.querySelector('canvas');
		if (!canvas) { console.error('[KiroRec] 找不到 canvas'); return false; }
		var stream = canvas.captureStream(30);
		window._kiroChunks = [];
		var options = { mimeType: 'video/webm;codecs=vp8' };
		try {
			window._kiroRecorder = new MediaRecorder(stream, options);
		} catch(e) {
			window._kiroRecorder = new MediaRecorder(stream);
		}
		window._kiroRecorder.ondataavailable = function(e) {
			if (e.data && e.data.size > 0) window._kiroChunks.push(e.data);
		};
		window._kiroRecorder.start(1000);
		console.log('[KiroRec] 開始錄製');
		return true;
	} catch(e) {
		console.error('[KiroRec] 錯誤:', e);
		return false;
	}
};

window.kiroStopRecording = function() {
	if (!window._kiroRecorder) return;
	window._kiroRecorder.onstop = function() {
		var blob = new Blob(window._kiroChunks, { type: 'video/webm' });
		var url = URL.createObjectURL(blob);
		var a = document.createElement('a');
		a.href = url;
		var ts = new Date().toISOString().replace(/[:.]/g, '-').slice(0,19);
		a.download = 'chiikawa_rec_' + ts + '.webm';
		document.body.appendChild(a);
		a.click();
		document.body.removeChild(a);
		URL.revokeObjectURL(url);
		console.log('[KiroRec] 錄製完成，已下載');
	};
	window._kiroRecorder.stop();
};
"""
	JavaScriptBridge.eval(js_code)
	print("[ScreenRecorder] HTML5 MediaRecorder 初始化完成")

# ---- HTML5 模式：開始錄製 ----
func _start_html5_recording() -> void:
	var ok = JavaScriptBridge.eval("window.kiroStartRecording()")
	if ok:
		_status_label.text = "REC"
	else:
		_status_label.text = "ERR"
		_status_label.modulate = Color(1.0, 0.5, 0.0)
		_is_recording = false
		_record_button.text = "⏺ REC"
		_set_button_idle_style()

# ---- HTML5 模式：停止錄製（自動下載 WebM）----
func _stop_html5_recording() -> void:
	JavaScriptBridge.eval("window.kiroStopRecording()")
	_status_label.text = "↓ DL"
	_status_label.modulate = Color(0.3, 1.0, 0.3)
	await get_tree().create_timer(2.0).timeout
	_status_label.text = "READY"
	_status_label.modulate = Color(0.6, 0.6, 0.6)
	_elapsed_label.text = "00:00"
