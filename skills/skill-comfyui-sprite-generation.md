# Skill：ComfyUI 精靈圖生成

## 目的
使用 ComfyUI + Stable Diffusion 生成符合吉伊卡哇風格的像素藝術精靈圖，包含角色、目標物、UI 元素。

## 適用場景
- 生成新角色精靈圖
- 生成目標物圖像（T001-T105）
- 生成 BOSS 圖像（B001）
- 重新生成品質不足的圖像

## 前置條件
- ComfyUI 服務運行中（http://127.0.0.1:8188）
- 已安裝對應的 Checkpoint 模型
- Python 3.10+ 已安裝（用於後處理）
- PIL/Pillow 已安裝（`pip install Pillow`）

## 確認 ComfyUI 服務狀態
```bash
# 確認服務是否運行
curl http://127.0.0.1:8188/system_stats

# 若未運行，啟動 ComfyUI
cd C:\ComfyUI
python main.py --listen 127.0.0.1 --port 8188
```

## 使用方法

### 步驟 1：準備 Prompt

#### 角色 Prompt 模板
```
正向 Prompt：
chiikawa character, pixel art, 16-bit style, cute japanese mascot, 
simple design, transparent background, 1px black outline, 
white body, big round eyes, small ears, [CHARACTER_SPECIFIC],
game sprite, top-down view, centered composition

負向 Prompt：
realistic, 3d render, blurry, low quality, watermark, text, 
extra limbs, deformed, ugly, dark background, complex background,
photorealistic, anime style (too detailed)
```

#### 目標物 Prompt 模板
```
正向 Prompt：
pixel art [FISH_TYPE], cute, 16-bit game sprite, 
simple design, transparent background, 1px black outline,
colorful, side view, [COLOR_DESCRIPTION],
fish shooting game target, small size

負向 Prompt：
realistic, scary, dark, complex, text, watermark,
human, character, background
```

### 步驟 2：呼叫 ComfyUI API

```python
# tools/generate_sprite.py
import json
import urllib.request
import urllib.parse
import time
import os

COMFYUI_URL = "http://127.0.0.1:8188"

def queue_prompt(prompt_workflow):
    """提交工作流程到 ComfyUI"""
    data = json.dumps({"prompt": prompt_workflow}).encode('utf-8')
    req = urllib.request.Request(
        f"{COMFYUI_URL}/prompt",
        data=data,
        headers={'Content-Type': 'application/json'}
    )
    response = urllib.request.urlopen(req)
    return json.loads(response.read())

def get_history(prompt_id):
    """取得生成歷史"""
    with urllib.request.urlopen(f"{COMFYUI_URL}/history/{prompt_id}") as response:
        return json.loads(response.read())

def wait_for_completion(prompt_id, timeout=120):
    """等待生成完成"""
    start_time = time.time()
    while time.time() - start_time < timeout:
        history = get_history(prompt_id)
        if prompt_id in history:
            return history[prompt_id]
        time.sleep(2)
    raise TimeoutError(f"生成超時（{timeout}秒）")

def generate_sprite(positive_prompt, negative_prompt, output_path, 
                    width=512, height=512, steps=28, cfg=7.5, seed=-1):
    """生成精靈圖"""
    
    # 基礎工作流程（需根據實際 ComfyUI 設定調整）
    workflow = {
        "3": {
            "class_type": "KSampler",
            "inputs": {
                "seed": seed if seed != -1 else int(time.time()),
                "steps": steps,
                "cfg": cfg,
                "sampler_name": "dpm_2m_karras",
                "scheduler": "karras",
                "denoise": 1.0,
                "model": ["4", 0],
                "positive": ["6", 0],
                "negative": ["7", 0],
                "latent_image": ["5", 0]
            }
        },
        "4": {
            "class_type": "CheckpointLoaderSimple",
            "inputs": {
                "ckpt_name": "pixel-art-v1.safetensors"
            }
        },
        "5": {
            "class_type": "EmptyLatentImage",
            "inputs": {
                "width": width,
                "height": height,
                "batch_size": 1
            }
        },
        "6": {
            "class_type": "CLIPTextEncode",
            "inputs": {
                "text": positive_prompt,
                "clip": ["4", 1]
            }
        },
        "7": {
            "class_type": "CLIPTextEncode",
            "inputs": {
                "text": negative_prompt,
                "clip": ["4", 1]
            }
        },
        "8": {
            "class_type": "VAEDecode",
            "inputs": {
                "samples": ["3", 0],
                "vae": ["4", 2]
            }
        },
        "9": {
            "class_type": "SaveImage",
            "inputs": {
                "filename_prefix": "sprite",
                "images": ["8", 0]
            }
        }
    }
    
    result = queue_prompt(workflow)
    prompt_id = result['prompt_id']
    
    print(f"提交生成任務：{prompt_id}")
    history = wait_for_completion(prompt_id)
    
    # 取得輸出圖像路徑
    outputs = history['outputs']
    for node_id, node_output in outputs.items():
        if 'images' in node_output:
            for image in node_output['images']:
                image_path = os.path.join(
                    "C:\\ComfyUI\\output",
                    image['filename']
                )
                print(f"生成完成：{image_path}")
                return image_path
    
    return None

# 使用範例
if __name__ == "__main__":
    result = generate_sprite(
        positive_prompt="chiikawa character, pixel art, 16-bit style, cute, transparent background",
        negative_prompt="realistic, 3d, blurry",
        output_path="assets/pending/chiikawa_lv1.png",
        width=512,
        height=512
    )
    print(f"輸出：{result}")
```

### 步驟 3：後處理

```python
# tools/process_sprites.py
from PIL import Image
import numpy as np

def process_sprite(input_path, output_path, target_size=(64, 64), max_colors=32):
    """
    後處理精靈圖：
    1. 縮放到目標尺寸（Nearest Neighbor）
    2. 去背（移除白色/接近白色背景）
    3. 限制調色板
    4. 加輪廓
    """
    img = Image.open(input_path).convert("RGBA")
    
    # 1. 去背（如果背景是白色）
    img = remove_white_background(img)
    
    # 2. 縮放（保持像素感）
    img = img.resize(target_size, Image.NEAREST)
    
    # 3. 限制調色板
    img = limit_palette(img, max_colors)
    
    # 4. 加輪廓
    img = add_outline(img, color=(26, 26, 26, 255))  # #1A1A1A
    
    # 儲存
    img.save(output_path, "PNG")
    print(f"處理完成：{output_path}")
    return output_path

def remove_white_background(img, threshold=240):
    """移除白色背景，轉為透明"""
    data = np.array(img)
    r, g, b, a = data[:,:,0], data[:,:,1], data[:,:,2], data[:,:,3]
    mask = (r > threshold) & (g > threshold) & (b > threshold)
    data[mask, 3] = 0
    return Image.fromarray(data)

def limit_palette(img, max_colors):
    """限制調色板顏色數"""
    # 轉換為 P 模式（調色板模式）再轉回 RGBA
    rgb_img = img.convert("RGB")
    quantized = rgb_img.quantize(colors=max_colors, method=Image.MEDIANCUT)
    result = quantized.convert("RGBA")
    # 恢復透明度
    alpha = img.split()[3]
    result.putalpha(alpha)
    return result

def add_outline(img, color=(26, 26, 26, 255), thickness=1):
    """加外輪廓"""
    from PIL import ImageFilter
    # 簡單輪廓：膨脹後與原圖合成
    alpha = img.split()[3]
    dilated = alpha.filter(ImageFilter.MaxFilter(thickness * 2 + 1))
    outline = Image.new("RGBA", img.size, color)
    outline.putalpha(dilated)
    result = Image.alpha_composite(outline, img)
    return result
```

## 範例

### 生成吉伊卡哇 LV1
```python
from tools.generate_sprite import generate_sprite
from tools.process_sprites import process_sprite

# 生成
raw_path = generate_sprite(
    positive_prompt="chiikawa character, pixel art, 16-bit, cute white animal, big eyes, small ears, transparent background, game sprite",
    negative_prompt="realistic, 3d, blurry, dark",
    output_path="temp/chiikawa_raw.png",
    width=512, height=512, steps=28, cfg=7.5
)

# 後處理
final_path = process_sprite(
    input_path=raw_path,
    output_path="assets/pending/chiikawa_lv1.png",
    target_size=(64, 64),
    max_colors=32
)
```

## 注意事項
- 生成前必須確認 ComfyUI 服務運行中
- 所有生成的圖像先放到 `assets/pending/`，等待 Art Director 審核
- 記錄每次生成的 Prompt 和 Seed（確保可重現）
- 若 ComfyUI 不可用，記錄到 `failed-attempts/` 並通知 Research Agent

## 已知問題
- ComfyUI 有時會生成帶有白色背景的圖像，需要後處理去背
- 調色板限制可能導致顏色失真，需要 Art Director 審核
- 生成速度依 GPU 效能而異（RTX 3080 約 10-30 秒/張）

## 版本記錄
- 2025-01-01：初始版本，基於現有生成流程
