---
name: comfyui-pixel-art
description: 使用 ComfyUI 生成吉伊卡哇像素藝術 Sprite 的完整流程
---

# ComfyUI 像素藝術生成 Skill

## 安裝狀態（2026-05-12）

| 項目 | 狀態 | 路徑 |
|------|------|------|
| ComfyUI Portable | ✅ 已安裝 | `C:\ComfyUI\ComfyUI_windows_portable\` |
| SD 1.5 基礎模型 | ✅ 已下載 | `...\models\checkpoints\v1-5-pruned-emaonly.safetensors` |
| Pixel Art LoRA | ✅ 已下載 | `...\models\loras\pixel_art_lora.safetensors` |
| 啟動腳本 | ✅ 已建立 | `tools/start_comfyui.bat` |
| API 整合腳本 | ✅ 已建立 | `tools/comfyui_generate.py` |

## 使用流程

### 1. 啟動 ComfyUI
```
tools\start_comfyui.bat
```
或直接：
```
C:\ComfyUI\ComfyUI_windows_portable\run_nvidia_gpu.bat
```
等待出現 `To see the GUI go to: http://127.0.0.1:8188`

### 2. 生成角色 Sprite

**生成單個角色：**
```
py tools/comfyui_generate.py --character chiikawa --pose idle
py tools/comfyui_generate.py --character hachiware --pose attack
py tools/comfyui_generate.py --character usagi --pose bigwin
```

**生成全部（3角色 × 3動作 = 9張）：**
```
py tools/comfyui_generate.py --all
```

輸出位置：`client/chiikawa-pixel/assets/sprites/ai_generated/`

### 3. 後處理（自動執行）
- 縮小到 96×96（NEAREST 插值）
- Flood fill 去除白色背景
- 儲存為透明 PNG

### 4. 整合到遊戲
生成的圖片放到 `assets/sprites/ai_generated/` 後，
執行 `py tools/generate_animation_frames.py` 合成 Spritesheet。

## Workflow 設定

| 參數 | 值 |
|------|-----|
| 基礎模型 | v1-5-pruned-emaonly.safetensors |
| LoRA | pixel_art_lora.safetensors（強度 0.85）|
| 解析度 | 512×512 |
| Steps | 28 |
| CFG | 7.5 |
| Sampler | euler_a |

## 提示詞模板

### 吉伊卡哇
```
chiikawa character, white fluffy round creature, tiny cute chibi, big black eyes with white highlight, 
small pink blush marks, pixel art sprite, 16-bit retro game style, transparent background, 
white fur, simple clean design, game sprite sheet
```

### 小八
```
hachiware character, white cat with blue stripes on head, pointed ears, cute chibi, big black eyes, 
pixel art sprite, 16-bit retro game style, transparent background, white fur with blue markings
```

### 烏薩奇
```
usagi rabbit character, white rabbit with long ears, red eyes, cute chibi, pixel art sprite, 
16-bit retro game style, transparent background, white fur, pink inner ears
```

## 常見問題

**Q: ComfyUI 啟動後沒有 GPU 加速？**
A: 確認使用 `run_nvidia_gpu.bat`，不是 `run_cpu.bat`

**Q: 生成的圖片背景不透明？**
A: `comfyui_generate.py` 的後處理會自動 flood fill 去背，確認 PIL 已安裝

**Q: 模型載入失敗？**
A: 確認 `v1-5-pruned-emaonly.safetensors` 在 checkpoints 資料夾，`pixel_art_lora.safetensors` 在 loras 資料夾

## 模型來源
- SD 1.5：[Comfy-Org/stable-diffusion-v1-5-archive](https://huggingface.co/Comfy-Org/stable-diffusion-v1-5-archive)
- Pixel Art LoRA：[nerijs/pixel-art-xl](https://huggingface.co/nerijs/pixel-art-xl)
