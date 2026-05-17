# Skill：Godot 4 動畫匯入與設定

## 目的
在 Godot 4.6.2 中正確匯入精靈圖序列，設定 AnimationPlayer 與 SpriteFrames，確保動畫流暢播放且與音效同步。

## 適用場景
- 匯入新的角色動畫序列
- 匯入目標物游泳動畫
- 設定 BOSS 動畫狀態機
- 調整動畫 FPS 與循環設定

## 前置條件
- Godot 4.6.2 已安裝
- 精靈圖序列已準備好（PNG 格式，透明背景）
- 已通過 Art Director 審核

## 使用方法

### 方法 1：使用 SpriteFrames 資源（推薦）

#### 步驟 1：建立 SpriteFrames 資源
```gdscript
# 在 Godot 編輯器中：
# 1. 選擇 AnimatedSprite2D 節點
# 2. 在 Inspector 中點擊 Frames 屬性
# 3. 選擇 "New SpriteFrames"
# 4. 在 SpriteFrames 編輯器中新增動畫
```

#### 步驟 2：設定動畫參數
```gdscript
# 角色 idle 動畫設定
var frames = SpriteFrames.new()
frames.add_animation("idle")
frames.set_animation_loop("idle", true)  # 循環播放
frames.set_animation_speed("idle", 8.0)  # 8 FPS

# 新增幀
for i in range(4):  # 4 幀動畫
    var texture = load("res://assets/characters/chiikawa_idle_%02d.png" % i)
    frames.add_frame("idle", texture)

$AnimatedSprite2D.sprite_frames = frames
$AnimatedSprite2D.play("idle")
```

#### 步驟 3：攻擊動畫設定
```gdscript
frames.add_animation("attack")
frames.set_animation_loop("attack", false)  # 不循環
frames.set_animation_speed("attack", 16.0)  # 16 FPS

for i in range(6):  # 6 幀攻擊動畫
    var texture = load("res://assets/characters/chiikawa_attack_%02d.png" % i)
    frames.add_frame("attack", texture)

# 攻擊動畫完成後回到 idle
$AnimatedSprite2D.animation_finished.connect(func():
    if $AnimatedSprite2D.animation == "attack":
        $AnimatedSprite2D.play("idle")
)
```

### 方法 2：使用 AnimationPlayer（複雜動畫）

```gdscript
# 適用於需要同時控制多個屬性的動畫（位置、縮放、顏色等）

func setup_boss_animation():
    var anim_player = $AnimationPlayer
    
    # 建立進場動畫
    var animation = Animation.new()
    animation.length = 2.0  # 2 秒
    
    # 位置軌道
    var track_idx = animation.add_track(Animation.TYPE_VALUE)
    animation.track_set_path(track_idx, ".:position")
    animation.track_insert_key(track_idx, 0.0, Vector2(640, -200))  # 從螢幕外
    animation.track_insert_key(track_idx, 1.0, Vector2(640, 300))   # 進入畫面
    animation.track_insert_key(track_idx, 2.0, Vector2(640, 280))   # 輕微彈跳
    
    # 縮放軌道（進場時放大效果）
    var scale_track = animation.add_track(Animation.TYPE_VALUE)
    animation.track_set_path(scale_track, ".:scale")
    animation.track_insert_key(scale_track, 0.0, Vector2(0.5, 0.5))
    animation.track_insert_key(scale_track, 1.0, Vector2(1.2, 1.2))
    animation.track_insert_key(scale_track, 2.0, Vector2(1.0, 1.0))
    
    anim_player.add_animation("boss_enter", animation)
    anim_player.play("boss_enter")
```

### 方法 3：批次匯入腳本

```python
# tools/import_animations.py
# 批次將圖像序列轉換為 Godot 可用的格式

import os
import json
from pathlib import Path

def create_sprite_frames_config(animation_name, frames_dir, fps=8, loop=True):
    """
    建立 SpriteFrames 設定檔
    Godot 4 可以透過 .tres 資源檔載入
    """
    frames = []
    frame_files = sorted(Path(frames_dir).glob("*.png"))
    
    for frame_file in frame_files:
        frames.append(str(frame_file.relative_to(".")))
    
    config = {
        "animation": animation_name,
        "fps": fps,
        "loop": loop,
        "frames": frames
    }
    
    output_path = f"assets/animations/{animation_name}_config.json"
    with open(output_path, 'w') as f:
        json.dump(config, f, indent=2)
    
    print(f"建立動畫設定：{output_path}（{len(frames)} 幀，{fps} FPS）")
    return output_path

# 使用範例
create_sprite_frames_config(
    animation_name="chiikawa_idle",
    frames_dir="assets/characters/chiikawa/idle/",
    fps=8,
    loop=True
)
```

## 動畫規格

### 角色動畫規格
| 動畫 | 幀數 | FPS | 循環 | 說明 |
|------|------|-----|------|------|
| idle | 4-8 幀 | 8 | 是 | 待機動畫，必須無縫循環 |
| attack | 6-12 幀 | 16-24 | 否 | 攻擊動畫，播完回 idle |
| death | 4-8 幀 | 12 | 否 | 死亡/消失動畫 |
| special | 8-16 幀 | 16 | 否 | 特殊技能動畫 |

### 目標物動畫規格
| 動畫 | 幀數 | FPS | 循環 | 說明 |
|------|------|-----|------|------|
| swim | 4-8 幀 | 8-12 | 是 | 游泳動畫，必須無縫循環 |
| hit | 2-4 幀 | 16 | 否 | 被擊中閃白效果 |
| die | 4-8 幀 | 12 | 否 | 消滅動畫 |

### BOSS 動畫規格
| 動畫 | 幀數 | FPS | 循環 | 說明 |
|------|------|-----|------|------|
| enter | 8-16 幀 | 16 | 否 | 進場動畫 |
| idle | 4-8 幀 | 8 | 是 | 待機動畫 |
| attack | 8-16 幀 | 16 | 否 | 攻擊動畫 |
| hurt | 4-8 幀 | 16 | 否 | 受傷動畫 |
| die | 8-16 幀 | 12 | 否 | 死亡動畫 |

## 音效同步技巧

```gdscript
# 在特定幀觸發音效
func _on_AnimatedSprite2D_frame_changed():
    var sprite = $AnimatedSprite2D
    if sprite.animation == "attack":
        # 第 3 幀觸發攻擊音效
        if sprite.frame == 2:  # 0-indexed
            $AudioManager.play_attack_sound(current_character)
```

## 效能優化

```gdscript
# 使用 CanvasItem 的 visibility_changed 暫停不可見的動畫
func _on_visibility_changed():
    if visible:
        $AnimatedSprite2D.play()
    else:
        $AnimatedSprite2D.stop()

# 限制同時播放的動畫數量
const MAX_ACTIVE_ANIMATIONS = 20
var active_animations = 0
```

## 注意事項
- idle 動畫的第一幀和最後一幀必須相似（確保無縫循環）
- 攻擊音效必須在攻擊動畫的第 3 幀觸發（視覺上最有衝擊感的時機）
- HTML5 環境下，動畫數量過多會影響效能，超過 20 個同時播放需要優化
- 修改動畫後必須在 HTML5 環境測試（Godot 編輯器預覽可能有差異）

## 已知問題
- Godot 4 HTML5 匯出時，某些動畫可能有 1-2 幀的延遲，需要在觸發時機上補償
- SpriteFrames 資源在大量目標物時可能佔用較多記憶體，考慮使用 Atlas Texture

## 版本記錄
- 2025-01-01：初始版本，基於 Godot 4.6.2
