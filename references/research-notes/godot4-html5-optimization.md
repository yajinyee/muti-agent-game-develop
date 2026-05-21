# Godot 4 HTML5 優化研究筆記

> 研究者：Research Agent  
> 最後更新：2026-05-17  
> 適用版本：Godot 4.6.2

---

## 效能優化技巧

### Texture Compression（紋理壓縮）

#### 推薦設定
```
Project Settings > Rendering > Textures > Default Texture Filter: Nearest
（像素藝術必須使用 Nearest，避免模糊）

匯入設定（.import 檔案）：
compress/mode = 0  # Lossless（PNG 像素藝術用這個）
compress/high_quality = false
mipmaps/generate = false  # 像素藝術不需要 mipmap
```

#### Spritesheet 優化
- 使用 2 的冪次方尺寸（256x256, 512x512, 1024x1024）
- 合併小圖到 Atlas Texture，減少 draw call
- 建議工具：`tools/animation_pipeline.py` 的 pack 功能

#### HTML5 特別注意
- HTML5 不支援 S3TC/BPTC 壓縮，使用 ETC2 或 ASTC
- 在 Export 設定中選擇 `VRAM Compressed (Lossy)` 或 `Lossless`
- 像素藝術建議 Lossless，避免壓縮失真

---

### Audio Streaming（音效串流）

#### 設定方式
```gdscript
# 大型 BGM 檔案（> 1MB）使用串流
# 在 .import 設定中：
# AudioStreamWAV: loop/mode = 1（循環）
# 或使用 AudioStreamOggVorbis 節省空間

# 小型 SFX 不需要串流，直接載入記憶體
```

#### HTML5 音效限制
- **重要**：HTML5 需要使用者互動後才能播放音效（瀏覽器政策）
- 解決方案：在第一次點擊時解鎖 AudioServer
```gdscript
func _ready():
    # 等待使用者互動
    get_viewport().gui_focus_changed.connect(_on_first_interaction)

func _on_first_interaction(_node):
    AudioServer.set_bus_mute(AudioServer.get_bus_index("Master"), false)
```

#### 音效格式建議
- WAV：短音效（< 1 秒），無壓縮，最低延遲
- OGG：長 BGM，壓縮比好，HTML5 支援良好
- MP3：避免使用（Godot 4 HTML5 有相容性問題）

---

### Draw Call Reduction（減少繪製呼叫）

#### 使用 CanvasItem 批次渲染
```gdscript
# 同一個 CanvasLayer 內的相同材質會自動批次
# 確保所有魚類使用同一個 Atlas Texture

# 錯誤做法：每條魚用不同 Texture
var fish = Sprite2D.new()
fish.texture = load("res://assets/fish_001.png")  # 每個都不同 = 多個 draw call

# 正確做法：使用 Atlas
var fish = Sprite2D.new()
fish.texture = atlas_texture  # 共用 Atlas = 批次渲染
fish.region_enabled = true
fish.region_rect = Rect2(0, 0, 32, 32)  # 指定 Atlas 中的位置
```

#### 粒子效果優化
```gdscript
# 使用 GPUParticles2D 而非 CPUParticles2D
# HTML5 上 GPU 粒子效能更好

# 限制同時存在的粒子數量
$GPUParticles2D.amount = 20  # 不要超過 50
```

#### Z-Index 管理
```
Layer 0: 背景（海底）
Layer 1: 目標物（魚類）
Layer 2: 子彈
Layer 3: 特效（命中/爆炸）
Layer 4: 角色
Layer 5: UI
Layer 6: 彈出視窗
```

---

## WebSocket 最佳實踐

### 連線管理
```gdscript
# 使用 WebSocketPeer（Godot 4 推薦）
var ws = WebSocketPeer.new()

func _ready():
    ws.connect_to_url("ws://localhost:7777/ws")

func _process(delta):
    ws.poll()
    var state = ws.get_ready_state()
    if state == WebSocketPeer.STATE_OPEN:
        while ws.get_available_packet_count() > 0:
            var packet = ws.get_packet()
            _handle_message(packet.get_string_from_utf8())
    elif state == WebSocketPeer.STATE_CLOSED:
        _handle_disconnect()
```

### 心跳機制
```gdscript
var heartbeat_timer: float = 0.0
const HEARTBEAT_INTERVAL = 30.0

func _process(delta):
    heartbeat_timer += delta
    if heartbeat_timer >= HEARTBEAT_INTERVAL:
        heartbeat_timer = 0.0
        _send_heartbeat()

func _send_heartbeat():
    var msg = {"type": "ping", "timestamp": Time.get_unix_time_from_system()}
    ws.send_text(JSON.stringify(msg))
```

### 訊息佇列（避免 HTML5 阻塞）
```gdscript
var message_queue: Array = []

func _send_message(msg: Dictionary):
    message_queue.append(JSON.stringify(msg))

func _process(delta):
    # 每幀最多發送 5 條訊息，避免阻塞
    var sent = 0
    while message_queue.size() > 0 and sent < 5:
        ws.send_text(message_queue.pop_front())
        sent += 1
```

---

## 已知問題與解決方案

### 問題 1：HTML5 音效無法自動播放
- **症狀**：遊戲開始時沒有音效
- **原因**：瀏覽器安全政策，需要使用者互動
- **解決**：在 UI 點擊事件中解鎖 AudioServer（見上方程式碼）

### 問題 2：HTML5 WebSocket 連線到 localhost 失敗
- **症狀**：在 HTTPS 頁面無法連線 ws://localhost
- **原因**：Mixed Content 安全限制
- **解決**：開發時使用 HTTP 頁面，或使用 wss:// + HTTPS

### 問題 3：像素藝術在高 DPI 螢幕模糊
- **症狀**：Retina/4K 螢幕上像素藝術模糊
- **原因**：瀏覽器縮放使用雙線性插值
- **解決**：
```gdscript
# Project Settings > Display > Window > Stretch > Mode = canvas_items
# Project Settings > Display > Window > Stretch > Aspect = keep
# 在 CSS 中加入：
# canvas { image-rendering: pixelated; }
```

### 問題 4：HTML5 Build 載入時間過長
- **症狀**：初始載入超過 10 秒
- **原因**：WASM 檔案過大
- **解決**：
  - 啟用 GZip 壓縮（需要 Web Server 支援）
  - 減少不必要的 GDScript 模組
  - 使用 Export 的 `Exclude Files` 排除不需要的資源

### 問題 5：Godot 4.6.2 HTML5 WebSocket 斷線重連
- **症狀**：網路不穩時斷線後無法重連
- **解決**：
```gdscript
func _handle_disconnect():
    await get_tree().create_timer(3.0).timeout
    ws = WebSocketPeer.new()
    ws.connect_to_url(SERVER_URL)
    reconnect_count += 1
    if reconnect_count > 5:
        show_error_dialog("連線失敗，請重新整理頁面")
```

---

## 效能基準

| 指標 | 目標 | 測試設備 |
|------|------|---------|
| 初始載入時間 | < 5 秒 | 中階電腦 + 100Mbps |
| 穩定 FPS | >= 60 | 中階電腦 |
| 最低 FPS | >= 30 | 低階電腦 |
| 記憶體使用 | < 200MB | 任何設備 |
| WebSocket 延遲 | < 50ms | 本地網路 |

---

## 參考資源

- [Godot 4 HTML5 Export 官方文件](https://docs.godotengine.org/en/stable/tutorials/export/exporting_for_web.html)
- [Godot 4 WebSocketPeer 文件](https://docs.godotengine.org/en/stable/classes/class_websocketpeer.html)
- [Pixel Art Rendering in Godot](https://docs.godotengine.org/en/stable/tutorials/2d/2d_sprite_animation.html)
