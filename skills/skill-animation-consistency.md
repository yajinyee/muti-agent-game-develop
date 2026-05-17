# Skill: Animation Consistency

> 版本：1.0.0  
> 來源：agent-sprite-forge 實戰經驗  
> 最後更新：2026-05-17  
> 適用工具：animation_pipeline.py

---

## 概覽

確保像素藝術動畫各幀在視覺上保持一致的技術集合。這是捕魚機遊戲動畫品質的核心技能。

---

## 技術 1：shared_scale（共享比例）

### 問題
不同幀的角色大小不一致，導致動畫播放時角色忽大忽小。

### 解決方案

```python
from PIL import Image

# 所有角色使用統一的 canvas size
CANVAS_SIZES = {
    "chiikawa": (64, 64),
    "hachiware": (64, 64),
    "usagi": (64, 64),
    "boss": (128, 128),
    "fish_normal": (48, 48),
    "fish_medium": (64, 64),
    "fish_large": (80, 80),
}

def normalize_to_canvas(img: Image.Image, char_type: str) -> Image.Image:
    """將 sprite 縮放到標準 canvas，保持比例，底部對齊"""
    target_size = CANVAS_SIZES.get(char_type, (64, 64))
    
    # 取得非透明區域
    bbox = img.getbbox()
    if bbox is None:
        return Image.new('RGBA', target_size, (0, 0, 0, 0))
    
    # 裁切到非透明區域
    cropped = img.crop(bbox)
    
    # 計算縮放比例（保持比例，不超過 canvas 的 90%）
    max_w = int(target_size[0] * 0.9)
    max_h = int(target_size[1] * 0.9)
    
    scale = min(max_w / cropped.width, max_h / cropped.height)
    new_w = int(cropped.width * scale)
    new_h = int(cropped.height * scale)
    
    # 縮放（像素藝術用 NEAREST）
    scaled = cropped.resize((new_w, new_h), Image.NEAREST)
    
    # 建立 canvas，底部對齊
    canvas = Image.new('RGBA', target_size, (0, 0, 0, 0))
    x = (target_size[0] - new_w) // 2  # 水平置中
    y = target_size[1] - new_h          # 底部對齊
    canvas.paste(scaled, (x, y))
    
    return canvas
```

### 使用時機
- 從 ComfyUI 生成新幀後，立即套用 normalize_to_canvas
- 確保所有幀在合成 spritesheet 前都是相同 canvas size

---

## 技術 2：bottom_align（底部對齊）

### 問題
角色在不同幀的垂直位置不一致，導致動畫播放時角色上下跳動。

### 解決方案

```python
def bottom_align_frames(frames: list, canvas_height: int) -> list:
    """將所有幀底部對齊到 canvas 底部"""
    aligned = []
    
    for frame in frames:
        bbox = frame.getbbox()
        if bbox is None:
            aligned.append(frame)
            continue
        
        # 計算底部到 canvas 底部的距離
        current_bottom = bbox[3]
        dy = canvas_height - current_bottom
        
        if dy == 0:
            aligned.append(frame)
            continue
        
        # 建立新 canvas，將幀向下移動
        new_frame = Image.new('RGBA', frame.size, (0, 0, 0, 0))
        new_frame.paste(frame, (0, dy))
        aligned.append(new_frame)
    
    return aligned

# 使用範例
frames = [Image.open(f) for f in frame_files]
aligned_frames = bottom_align_frames(frames, canvas_height=64)
```

### 注意事項
- 底部對齊後，確認角色不會超出 canvas 頂部
- 如果角色太高，需要先縮小再對齊

---

## 技術 3：keep_largest_component（去噪點）

### 問題
AI 生成的圖像去背後，常有零散的噪點像素殘留，影響動畫品質。

### 解決方案

```python
import numpy as np
from scipy import ndimage

def keep_largest_component(img: Image.Image, min_alpha: int = 10) -> Image.Image:
    """保留最大連通區域，去除噪點"""
    arr = np.array(img.convert('RGBA'))
    alpha = arr[:, :, 3]
    
    # 二值化 alpha 通道
    binary = alpha > min_alpha
    
    if not binary.any():
        return img
    
    # 標記連通區域
    labeled, num_features = ndimage.label(binary)
    
    if num_features <= 1:
        return img  # 只有一個區域，不需要處理
    
    # 找最大區域
    sizes = ndimage.sum(binary, labeled, range(1, num_features + 1))
    largest_label = int(np.argmax(sizes)) + 1
    
    # 只保留最大區域
    mask = labeled == largest_label
    result_arr = arr.copy()
    result_arr[:, :, 3] = np.where(mask, arr[:, :, 3], 0)
    
    return Image.fromarray(result_arr)

# 使用範例
clean_frame = keep_largest_component(raw_frame)
```

### 使用時機
- 在 rembg 去背後立即執行
- 在合成 spritesheet 前執行

---

## 技術 4：透明 GIF 輸出

### 問題
標準 GIF 不支援真正的透明度（只有 1-bit 透明），直接輸出會有白色背景。

### 解決方案

```python
def save_preview_gif(frames: list, output_path: str, fps: int = 8,
                     bg_color: tuple = (180, 180, 180)):
    """輸出帶有固定背景色的預覽 GIF（用於審核）"""
    duration = int(1000 / fps)
    
    palette_frames = []
    for frame in frames:
        if frame.mode != 'RGBA':
            frame = frame.convert('RGBA')
        
        # 合成到背景色
        bg = Image.new('RGBA', frame.size, bg_color + (255,))
        bg.paste(frame, mask=frame.split()[3])
        
        # 轉換為調色板模式
        p_frame = bg.convert('P', palette=Image.ADAPTIVE, colors=255)
        palette_frames.append(p_frame)
    
    if not palette_frames:
        return
    
    palette_frames[0].save(
        output_path,
        save_all=True,
        append_images=palette_frames[1:],
        duration=duration,
        loop=0,
        optimize=False
    )

# 如果需要真正透明的 GIF（用於網頁展示）
def save_transparent_gif_apng(frames: list, output_path: str, fps: int = 8):
    """輸出 APNG（支援真正透明度，比 GIF 更好）"""
    # 需要 apng 套件：pip install apng
    try:
        from apng import APNG, PNG
        duration = int(1000 / fps)
        
        png_frames = []
        for i, frame in enumerate(frames):
            tmp_path = f"/tmp/frame_{i}.png"
            frame.save(tmp_path, 'PNG')
            png_frames.append(PNG.from_file(tmp_path))
        
        anim = APNG()
        for png_frame in png_frames:
            anim.append(png_frame, delay=duration)
        anim.save(output_path)
    except ImportError:
        print("[WARN] apng 套件未安裝，改用標準 GIF")
        save_preview_gif(frames, output_path, fps)
```

---

## 已知問題與解決方案

| 問題 | 原因 | 解決方案 |
|------|------|---------|
| 動畫抖動 | anchor point 不一致 | 使用 bottom_align_frames |
| 角色大小不一 | 生成時比例不固定 | 使用 normalize_to_canvas |
| 噪點殘留 | rembg 去背不完整 | 使用 keep_largest_component |
| GIF 顏色失真 | 調色板限制 256 色 | 使用 APNG 或 PNG 序列 |
| 幀間跳躍 | 關鍵幀差異過大 | 增加補間幀，降低 denoise |

---

## 完整工作流程

```python
# 標準動畫幀處理流程
def process_animation_frame(raw_path: str, char_type: str) -> Image.Image:
    """處理單個動畫幀的完整流程"""
    from rembg import remove
    
    # Step 1: 去背
    with open(raw_path, 'rb') as f:
        raw_data = f.read()
    clean_data = remove(raw_data)
    img = Image.open(io.BytesIO(clean_data)).convert('RGBA')
    
    # Step 2: 去噪點
    img = keep_largest_component(img)
    
    # Step 3: 標準化比例
    img = normalize_to_canvas(img, char_type)
    
    return img

# 處理整個動畫
def process_animation(frame_paths: list, char_type: str, 
                       output_path: str, fps: int = 8):
    frames = [process_animation_frame(p, char_type) for p in frame_paths]
    frames = bottom_align_frames(frames, CANVAS_SIZES[char_type][1])
    
    # 合成 spritesheet
    create_spritesheet(frames, output_path)
    
    # 生成預覽 GIF
    gif_path = output_path.replace('.png', '_preview.gif')
    save_preview_gif(frames, gif_path, fps)
    
    return frames
```

---

## 相關工具

- `tools/animation_pipeline.py` — 主要動畫工具
- `tools/qa_check.py` — QA 檢查（包含 Sprite QC）
- `docs/feature-specs/animation-pipeline-spec.md` — 完整規格
