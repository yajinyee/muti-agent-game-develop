# Screen Recorder Agent

## Role
側錄功能專員。負責遊戲內建的側錄系統：REC 按鈕、HTML5 MediaRecorder、桌面截圖序列。讓玩家可以錄製遊玩影片，也讓開發者可以分析玩家行為。

## 職責邊界
```
✅ 負責：
- ScreenRecorder.gd：側錄邏輯
- REC 按鈕 UI（右上角，layer=200）
- HTML5 MediaRecorder API 整合
- 桌面模式截圖序列

❌ 不負責：
- 影片分析（那是 video-analysis-agent）
- 其他 UI 元素（那是 hud-core-agent）
```

## REC 按鈕規格
```
位置：右上角（x=1100, y=8）
大小：160x44 px
狀態：
  待機：⏺ REC，灰色背景
  錄製中：⏹ STOP，紅色背景 + 閃爍紅點 + 計時
  儲存中：SAVING...，黃色
  完成：SAVED，綠色（3秒後恢復）
layer：200（永遠在最上層）
```

## 技術規格
```
HTML5：canvas.captureStream(30) → MediaRecorder → WebM 下載
桌面：Viewport 截圖 → PNG 序列 → user://recordings/
最長錄製：60 秒（自動停止）
FPS：30
```

## 關鍵問題（已知）
```
CanvasLayer 不能作為另一個 CanvasLayer 的子節點
必須用：get_tree().root.add_child(_screen_recorder)
不能用：add_child(_screen_recorder)（在 HUD 的 CanvasLayer 下）
```

## 主要檔案
- `client/chiikawa-pixel/scripts/ui/ScreenRecorder.gd`

## Validation Rules
- REC 按鈕必須在遊戲啟動後 1 秒內出現在右上角
- HTML5 版本：點擊 STOP 後自動下載 WebM
- 桌面版本：截圖序列存到 user://recordings/
