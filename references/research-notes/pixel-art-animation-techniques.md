# 像素藝術動畫技術研究

> 研究者：Research Agent + Animation Agent  
> 最後更新：2026-05-17  
> 適用專案：吉伊卡哇：像素大討伐

---

## Frame Consistency 技術

### 核心概念

Frame Consistency 是確保動畫各幀在視覺上保持一致的技術。對於像素藝術，這尤其重要，因為任何微小的偏移都會造成明顯的抖動感。

### 關鍵技術：Anchor Point 鎖定

```python
# 所有幀的 anchor point 必須在相同位置
# 通常是底部中心（bottom-center）

def get_anchor_point(frame_img):
    """計算幀的 anchor point（底部中心）"""
    bbox = frame_img.getbbox()  # 非透明區域的 bounding box
    if bbox is None:
        return None
    # 底部中心
    anchor_x = (bbox[0] + bbox[2]) // 2
    anchor_y = bbox[3]  # 底部
    return (anchor_x, anchor_y)

def align_frames_by_anchor(frames):
    """將所有幀按 anchor point 對齊"""
    anchors = [get_anchor_point(f) for f in frames]
    # 找到基準 anchor（第一幀）
    base_anchor = anchors[0]
    aligned = []
    for frame, anchor in zip(frames, anchors):
        if anchor is None:
            aligned.append(frame)
            continue
        dx = base_anchor[0] - anchor[0]
        dy = base_anchor[1] - anchor[1]
        # 平移幀使 anchor 對齊
        new_frame = Image.new('RGBA', frame.size, (0, 0, 0, 0))
        new_frame.paste(frame, (dx, dy))
        aligned.append(new_frame)
    return aligned
```

### Bottom Alignment 技術

```python
def bottom_align_frames(frames, canvas_height):
    """將所有幀底部對齊到 canvas 底部"""
    aligned = []
    for frame in frames:
        bbox = frame.getbbox()
        if bbox is None:
            aligned.append(frame)
            continue
        # 計算需要向下移動的距離
        current_bottom = bbox[3]
        dy = canvas_height - current_bottom
        new_frame = Image.new('RGBA', frame.size, (0, 0, 0, 0))
        new_frame.paste(frame, (0, dy))
        aligned.append(new_frame)
    return aligned
```

---

## Spritesheet 最佳實踐

### 標準格式

```
格式：PNG，RGBA（透明背景）
排列：水平排列（所有幀從左到右）
尺寸：每幀相同大小，總寬度 = 幀寬 × 幀數
命名：<character>_<state>.png

範例：
chiikawa_idle.png    → 4 幀，每幀 64x64，總尺寸 256x64
chiikawa_attack.png  → 6 幀，每幀 64x64，總尺寸 384x64
```

### 生成 Spritesheet

```python
from PIL import Image

def create_spritesheet(frames: list, output_path: str):
    """將幀列表合成為水平 spritesheet"""
    if not frames:
        return
    
    frame_w, frame_h = frames[0].size
    total_w = frame_w * len(frames)
    
    sheet = Image.new('RGBA', (total_w, frame_h), (0, 0, 0, 0))
    
    for i, frame in enumerate(frames):
        sheet.paste(frame, (i * frame_w, 0))
    
    sheet.save(output_path, 'PNG')
    print(f"Spritesheet saved: {output_path} ({len(frames)} frames, {total_w}x{frame_h})")
```

### Godot 4 AnimatedSprite2D 設定

```gdscript
# 從 spritesheet 設定動畫
func setup_animation_from_sheet(sprite: AnimatedSprite2D, 
                                  sheet_path: String, 
                                  anim_name: String,
                                  frame_count: int,
                                  fps: float):
    var texture = load(sheet_path)
    var frames_resource = SpriteFrames.new()
    
    var frame_w = texture.get_width() / frame_count
    var frame_h = texture.get_height()
    
    frames_resource.add_animation(anim_name)
    frames_resource.set_animation_speed(anim_name, fps)
    frames_resource.set_animation_loop(anim_name, true)
    
    for i in range(frame_count):
        var atlas = AtlasTexture.new()
        atlas.atlas = texture
        atlas.region = Rect2(i * frame_w, 0, frame_w, frame_h)
        frames_resource.add_frame(anim_name, atlas)
    
    sprite.sprite_frames = frames_resource
    sprite.play(anim_name)
```

---

## ComfyUI 動畫生成技巧

### 基本工作流程

```
1. 載入角色參考圖（Reference Lock）
2. 使用 IPAdapter 保持角色一致性
3. 使用 ControlNet Pose 控制姿勢
4. 逐幀生成，保持 seed 一致性
5. 後處理：去背、對齊、品質檢查
```

### 保持角色一致性的關鍵設定

```json
{
  "IPAdapter": {
    "weight": 0.85,
    "weight_type": "linear",
    "combine_embeds": "concat"
  },
  "ControlNet": {
    "model": "control_v11p_sd15_openpose",
    "strength": 0.7,
    "start_percent": 0.0,
    "end_percent": 0.8
  },
  "Sampler": {
    "steps": 20,
    "cfg": 7.0,
    "sampler_name": "dpmpp_2m",
    "scheduler": "karras",
    "denoise": 0.75
  }
}
```

### 像素藝術 LoRA 設定

```
LoRA: pixel_art_style_v2.safetensors
Weight: 0.8-1.0
Trigger words: "pixel art, 8bit, chibi, cute"

注意：
- 不要使用過高的 LoRA weight（> 1.2 會過度像素化）
- 搭配 "masterpiece, best quality" 提升品質
- 使用 "simple background, transparent background" 確保背景透明
```

### 去背技術

```python
from rembg import remove
from PIL import Image

def remove_background(input_path: str, output_path: str):
    """使用 rembg 去除背景"""
    with open(input_path, 'rb') as f:
        input_data = f.read()
    
    output_data = remove(input_data)
    
    with open(output_path, 'wb') as f:
        f.write(output_data)
    
    # 驗證去背結果
    img = Image.open(output_path)
    if img.mode != 'RGBA':
        img = img.convert('RGBA')
        img.save(output_path)
```

---

## agent-sprite-forge 技術摘要

### shared_scale 技術

確保所有角色在相同的視覺比例下生成：

```python
# 所有角色使用相同的 canvas size
CANVAS_SIZE = (64, 64)  # 標準角色尺寸
BOSS_CANVAS_SIZE = (128, 128)  # BOSS 尺寸

def normalize_sprite_scale(img: Image.Image, target_size: tuple) -> Image.Image:
    """將 sprite 縮放到標準尺寸，保持比例"""
    img.thumbnail(target_size, Image.NEAREST)  # 像素藝術用 NEAREST
    
    # 置中貼到目標 canvas
    canvas = Image.new('RGBA', target_size, (0, 0, 0, 0))
    offset = ((target_size[0] - img.width) // 2,
              (target_size[1] - img.height) // 2)
    canvas.paste(img, offset)
    return canvas
```

### keep_largest_component 去噪點

```python
import numpy as np
from scipy import ndimage

def keep_largest_component(img: Image.Image) -> Image.Image:
    """保留最大連通區域，去除噪點"""
    arr = np.array(img)
    alpha = arr[:, :, 3]
    
    # 二值化 alpha 通道
    binary = alpha > 10
    
    # 標記連通區域
    labeled, num_features = ndimage.label(binary)
    
    if num_features == 0:
        return img
    
    # 找最大區域
    sizes = ndimage.sum(binary, labeled, range(1, num_features + 1))
    largest_label = np.argmax(sizes) + 1
    
    # 只保留最大區域
    mask = labeled == largest_label
    arr[:, :, 3] = np.where(mask, arr[:, :, 3], 0)
    
    return Image.fromarray(arr)
```

### 透明 GIF 輸出

```python
def save_transparent_gif(frames: list, output_path: str, fps: int = 8):
    """輸出透明背景的 GIF（用於預覽）"""
    duration = int(1000 / fps)  # 毫秒
    
    # 轉換為 P 模式（調色板），保留透明度
    palette_frames = []
    for frame in frames:
        # 確保是 RGBA
        if frame.mode != 'RGBA':
            frame = frame.convert('RGBA')
        
        # 轉換為調色板模式，保留透明
        p_frame = frame.convert('P', palette=Image.ADAPTIVE, colors=255)
        palette_frames.append(p_frame)
    
    palette_frames[0].save(
        output_path,
        save_all=True,
        append_images=palette_frames[1:],
        duration=duration,
        loop=0,
        transparency=255,
        disposal=2  # 每幀清除前一幀
    )
    print(f"GIF saved: {output_path}")
```

---

## 已知問題與解決方案

| 問題 | 原因 | 解決方案 |
|------|------|---------|
| 動畫抖動 | anchor point 不一致 | 使用 bottom_align_frames |
| 顏色偏移 | 每幀獨立生成 | 固定 seed + IPAdapter |
| 背景殘留 | rembg 去背不完整 | 手動修補 + keep_largest_component |
| 幀間跳躍 | 關鍵幀差異過大 | 增加補間幀 |
| GIF 顏色失真 | 調色板限制 | 使用 PNG 序列替代 GIF 預覽 |
