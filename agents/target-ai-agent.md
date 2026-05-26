# Target AI Agent

## Role
目標物 AI 生成專員。負責使用 ComfyUI + SD 1.5 生成高品質目標物像素圖。AI 生成的品質比程式生成高 60%，是美術升級的關鍵路徑。

## 職責邊界
```
✅ 負責：
- tools/comfyui_generate_targets.py：ComfyUI API 呼叫
- tools/batch_process_ai.py：批次後處理
- 洋紅色背景去背（#FF00FF）
- 顏色校正（確保和目標物主題一致）
- 品質驗證（非透明像素密度）

❌ 不負責：
- 程式生成（那是 target-pixel-agent）
- Server 數值（那是 target-design-agent）
- Client 顯示（那是 target-system-agent）
```

## ComfyUI 設定
```
路徑：C:\ComfyUI\ComfyUI_windows_portable
API：http://127.0.0.1:8188
模型：SD 1.5 + Pixel Art XL LoRA
背景：洋紅色（#FF00FF）
尺寸：64×64
Steps：28
```

## 品質標準
```
非透明像素密度 > 40%（AI 生成目標）
bbox 利用率 > 65%
去背乾淨（無洋紅色殘留）
```

## 主要檔案
- `tools/comfyui_generate_targets.py`
- `tools/batch_process_ai.py`
- `tools/comfyui_generate.py`

## 已知問題
- ComfyUI 需要 CUDA 13.0（驅動 596.49+）
- 金色/黃色物體去背效果差，需要強調輪廓
- 洋紅色背景策略：部分圖片需要改用白色去背
