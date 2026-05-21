# Skill: Gameplay Juice（遊戲手感強化）

**建立日期**：2026-05-22  
**建立者**：Research Agent + Godot Client Agent  
**版本**：1.0.0

---

## 概念

「Juice」是遊戲設計術語，指讓遊戲操作感覺更爽快的視覺/音效/時間反饋技術。
捕魚機遊戲的 Juice 重點：**命中感、擊殺爆炸感、大獎震撼感**。

---

## 核心技術

### 1. Screen Shake（畫面震動）

**Trauma 系統**（業界最佳實踐）：
- `trauma` 值 [0, 1]，加法累積，乘法衰減
- `shake = trauma²`（讓小 trauma 更柔和）
- 用 sin/cos 組合模擬平滑 noise（不需要 FastNoiseLite）

```gdscript
# 衰減
_trauma = max(0.0, _trauma - decay * delta)
# 震動強度
var shake = _trauma * _trauma
var ox = sin(_time * 1.7) * cos(_time * 2.3) * max_offset.x * shake
```

**trauma 建議值**：
| 事件 | trauma |
|------|--------|
| 普通命中 | 0.18 |
| 擊殺 | 0.35 |
| 大獎（20x+）| 0.7 |
| BOSS 登場 | 0.9 |
| BOSS Phase 2 | 0.6 |
| Bonus 觸發 | 0.4 |

**像素遊戲注意**：`pixel_perfect = true`，offset 取整數，避免模糊。

### 2. Hit Stop（時間凍結）

短暫凍結時間（0.03~0.08 秒）讓命中感更強烈：

```gdscript
Engine.time_scale = 0.0
await get_tree().create_timer(duration, true, false, true).timeout
Engine.time_scale = 1.0
```

**注意**：`create_timer` 第 4 個參數 `ignore_time_scale=true` 確保 timer 不受 time_scale 影響。

### 3. 命中特效層次

| 層次 | 效果 | 用途 |
|------|------|------|
| 閃光環 | 小圓形擴散 | 普通命中 |
| 衝擊波 | 向外擴散環 | 擊殺 |
| 粒子噴射 | 多方向粒子 | 擊殺/大獎 |
| 全畫面閃光 | CanvasLayer ColorRect | 大獎/BOSS |

### 4. 子彈拖尾

沿飛行路徑生成漸隱殘影，增加速度感：

```gdscript
# 每 30ms 一個殘影，透明度隨位置遞減
for i in steps:
    var t = float(i) / float(steps)
    var trail_pos = from.lerp(to, t)
    # 延遲生成 + 快速淡出
```

---

## Godot 4 實作架構

### Autoload 設計

```
project.godot [autoload]
  HitEffect = "*res://scripts/game/HitEffect.gd"   # 特效生成
  ScreenShake = "*res://scripts/game/ScreenShake.gd" # 震動控制
```

### ScreenShake 取得 Camera2D

```gdscript
# 延遲取得，避免場景未載入
func _find_camera() -> void:
    var scene = get_tree().current_scene
    if is_instance_valid(scene):
        _camera = scene.get_node_or_null("Camera2D")
```

### 全畫面閃光

```gdscript
# 用 CanvasLayer（layer=100）確保在最上層
var canvas = CanvasLayer.new()
canvas.layer = 100
var rect = ColorRect.new()
rect.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
```

---

## 效果評估

| 改善項目 | 改善前 | 改善後 |
|---------|--------|--------|
| 命中反饋 | 小閃光 0.1s | 閃光環 + 粒子 + Hit Stop |
| 擊殺反饋 | 縮放消失 | 爆炸 + 衝擊波 + 震動 |
| 大獎反饋 | 語音字卡 | 全畫面閃光 + 金色粒子雨 + 強震動 |
| BOSS 登場 | 無 | 全畫面紅色閃爍 + 最強震動 |
| 子彈 | 直線飛行 | 直線 + 拖尾殘影 |

預期 Gameplay Feel 從 88 提升到 92+。

---

## 踩坑記錄

1. **ScreenShake 不能繼承 Camera2D 作為 Autoload**：Autoload 是 Node，不能直接是 Camera2D。改用 Node 繼承，透過 `get_node_or_null` 找場景中的 Camera2D。
2. **Hit Stop 的 timer 要設 ignore_time_scale=true**：否則 time_scale=0 時 timer 永遠不會觸發。
3. **像素遊戲 offset 要取整數**：`round(ox)` 避免 sub-pixel 模糊。

---

*由 Research Agent 記錄，Godot Client Agent 實作驗證*
