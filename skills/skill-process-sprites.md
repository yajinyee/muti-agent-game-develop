# Skill：精靈圖後處理流程

## 目的
將 AI 生成的原始圖像（通常是 512x512 或更大）轉換為符合遊戲規格的像素藝術精靈圖，包含去背、縮放、調色板限制、加輪廓等步驟。

## 適用場景
- ComfyUI 生成圖像後的後處理
- 批次處理多張精靈圖
- 修復現有精靈圖的品質問題
- 統一調色板確保視覺一致性

## 前置條件
- Python 3.10+ 已安裝
- Pillow 已安裝：`pip install Pillow`
- NumPy 已安裝：`pip install numpy`
- 輸入圖像為 PNG 格式

## 完整後處理腳本

```python
# tools/process_sprites.py
"""
精靈圖後處理工具
用法：python tools/process_sprites.py --input [輸入路徑] --output [輸出路徑] --size [尺寸]
"""

import argparse
import os
from pathlib import Path
from PIL import Image, ImageFilter
import numpy as np

def remove_background(img: Image.Image, method: str = "white", threshold: int = 240) -> Image.Image:
    """
    去背景
    method: "white" = 移除白色背景, "alpha" = 保留現有透明度
    """
    img = img.convert("RGBA")
    data = np.array(img)
    
    if method == "white":
        r, g, b, a = data[:,:,0], data[:,:,1], data[:,:,2], data[:,:,3]
        # 移除接近白色的像素
        white_mask = (r > threshold) & (g > threshold) & (b > threshold)
        data[white_mask, 3] = 0
        # 移除接近白色的半透明像素
        near_white = (r > threshold - 20) & (g > threshold - 20) & (b > threshold - 20)
        data[near_white, 3] = np.minimum(data[near_white, 3], 
                                          255 - (r[near_white] - (threshold - 20)) * 5)
    
    return Image.fromarray(data)

def pixelate(img: Image.Image, target_size: tuple) -> Image.Image:
    """
    像素化縮放（保持像素藝術感）
    使用 NEAREST 插值確保不模糊
    """
    return img.resize(target_size, Image.NEAREST)

def limit_palette(img: Image.Image, max_colors: int = 32) -> Image.Image:
    """
    限制調色板顏色數
    保留透明度
    """
    # 分離 Alpha 通道
    if img.mode == "RGBA":
        r, g, b, a = img.split()
        rgb_img = Image.merge("RGB", (r, g, b))
    else:
        rgb_img = img.convert("RGB")
        a = None
    
    # 量化顏色
    quantized = rgb_img.quantize(colors=max_colors, method=Image.MEDIANCUT)
    result = quantized.convert("RGB")
    
    # 恢復透明度
    if a is not None:
        result = result.convert("RGBA")
        result.putalpha(a)
    
    return result

def add_outline(img: Image.Image, color: tuple = (26, 26, 26, 255), thickness: int = 1) -> Image.Image:
    """
    加外輪廓
    color: RGBA 顏色，預設深灰色 #1A1A1A
    thickness: 輪廓厚度（像素）
    """
    if img.mode != "RGBA":
        img = img.convert("RGBA")
    
    # 取得 Alpha 通道
    alpha = img.split()[3]
    
    # 膨脹 Alpha 通道（建立輪廓遮罩）
    dilated = alpha.filter(ImageFilter.MaxFilter(thickness * 2 + 1))
    
    # 建立輪廓圖層
    outline_layer = Image.new("RGBA", img.size, color)
    outline_layer.putalpha(dilated)
    
    # 合成：輪廓在底層，原圖在上層
    result = Image.alpha_composite(outline_layer, img)
    return result

def enhance_pixel_art(img: Image.Image) -> Image.Image:
    """
    增強像素藝術效果
    - 提高對比度
    - 確保邊緣清晰
    """
    from PIL import ImageEnhance
    
    # 提高對比度
    enhancer = ImageEnhance.Contrast(img.convert("RGB"))
    enhanced = enhancer.enhance(1.2)
    
    # 恢復透明度
    if img.mode == "RGBA":
        enhanced = enhanced.convert("RGBA")
        enhanced.putalpha(img.split()[3])
    
    return enhanced

def process_sprite(
    input_path: str,
    output_path: str,
    target_size: tuple = (64, 64),
    max_colors: int = 32,
    add_outline_flag: bool = True,
    outline_color: tuple = (26, 26, 26, 255),
    remove_bg: bool = True,
    bg_method: str = "white"
) -> str:
    """
    完整後處理流程
    """
    print(f"處理：{input_path} → {output_path}")
    
    # 載入圖像
    img = Image.open(input_path)
    print(f"  原始尺寸：{img.size}，模式：{img.mode}")
    
    # 1. 去背
    if remove_bg:
        img = remove_background(img, method=bg_method)
        print(f"  ✅ 去背完成")
    
    # 2. 縮放（像素化）
    img = pixelate(img, target_size)
    print(f"  ✅ 縮放到 {target_size}")
    
    # 3. 限制調色板
    img = limit_palette(img, max_colors)
    print(f"  ✅ 調色板限制到 {max_colors} 色")
    
    # 4. 加輪廓
    if add_outline_flag:
        img = add_outline(img, color=outline_color)
        print(f"  ✅ 加輪廓完成")
    
    # 確保輸出目錄存在
    os.makedirs(os.path.dirname(output_path), exist_ok=True)
    
    # 儲存
    img.save(output_path, "PNG", optimize=True)
    print(f"  ✅ 儲存完成：{output_path}")
    
    return output_path

def batch_process(input_dir: str, output_dir: str, **kwargs):
    """
    批次處理目錄中的所有 PNG 圖像
    """
    input_path = Path(input_dir)
    output_path = Path(output_dir)
    
    png_files = list(input_path.glob("*.png"))
    print(f"找到 {len(png_files)} 個 PNG 檔案")
    
    results = []
    for i, png_file in enumerate(png_files, 1):
        out_file = output_path / png_file.name
        print(f"\n[{i}/{len(png_files)}] 處理 {png_file.name}")
        try:
            result = process_sprite(str(png_file), str(out_file), **kwargs)
            results.append({"file": png_file.name, "status": "success", "output": result})
        except Exception as e:
            print(f"  ❌ 錯誤：{e}")
            results.append({"file": png_file.name, "status": "error", "error": str(e)})
    
    # 輸出摘要
    success = sum(1 for r in results if r["status"] == "success")
    print(f"\n批次處理完成：{success}/{len(png_files)} 成功")
    return results

# 命令列介面
if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="精靈圖後處理工具")
    parser.add_argument("--input", required=True, help="輸入路徑（檔案或目錄）")
    parser.add_argument("--output", required=True, help="輸出路徑（檔案或目錄）")
    parser.add_argument("--size", default="64x64", help="目標尺寸（格式：WxH）")
    parser.add_argument("--colors", type=int, default=32, help="最大顏色數")
    parser.add_argument("--no-outline", action="store_true", help="不加輪廓")
    parser.add_argument("--no-remove-bg", action="store_true", help="不去背")
    
    args = parser.parse_args()
    
    w, h = map(int, args.size.split("x"))
    target_size = (w, h)
    
    if os.path.isdir(args.input):
        batch_process(
            args.input, args.output,
            target_size=target_size,
            max_colors=args.colors,
            add_outline_flag=not args.no_outline,
            remove_bg=not args.no_remove_bg
        )
    else:
        process_sprite(
            args.input, args.output,
            target_size=target_size,
            max_colors=args.colors,
            add_outline_flag=not args.no_outline,
            remove_bg=not args.no_remove_bg
        )
```

## 使用範例

### 單張處理
```bash
# 處理單張圖像（64x64，32色，加輪廓）
python tools/process_sprites.py --input assets/pending/chiikawa_raw.png --output assets/characters/chiikawa_lv1.png --size 64x64 --colors 32

# 處理 BOSS 圖像（256x256，32色）
python tools/process_sprites.py --input assets/pending/boss_raw.png --output assets/targets/B001_boss.png --size 256x256 --colors 32
```

### 批次處理
```bash
# 批次處理所有待審核圖像
python tools/process_sprites.py --input assets/pending/ --output assets/processed/ --size 64x64 --colors 16
```

### Python 直接呼叫
```python
from tools.process_sprites import process_sprite, batch_process

# 單張
process_sprite(
    input_path="assets/pending/chiikawa_raw.png",
    output_path="assets/characters/chiikawa_lv1.png",
    target_size=(64, 64),
    max_colors=32
)

# 批次
batch_process(
    input_dir="assets/pending/",
    output_dir="assets/processed/",
    target_size=(64, 64),
    max_colors=16
)
```

## 各類型圖像的建議參數

| 類型 | 目標尺寸 | 最大色數 | 輪廓 | 去背方式 |
|------|---------|---------|------|---------|
| 角色 LV1-3 | 64x64 | 32 | 是 | white |
| 角色 LV4-7 | 96x96 | 32 | 是 | white |
| 角色 LV8-10 | 128x128 | 32 | 是 | white |
| 普通目標物 | 64x64 | 16 | 是 | white |
| 特殊目標物 | 64x64 | 24 | 是 | white |
| BOSS | 256x256 | 32 | 是 | white |
| UI 元素 | 可變 | 16 | 否 | white |

## 注意事項
- 去背後必須人工確認邊緣是否乾淨（自動去背可能有殘留）
- 調色板限制可能導致顏色失真，需要 Art Director 審核
- 輪廓加在縮放之後（確保輪廓是 1px）
- 輸出前確認透明背景正確（在棋盤格背景上預覽）

## 已知問題
- 某些 AI 生成圖像的背景不是純白色，需要調整 threshold 參數
- 調色板量化有時會讓相近顏色合併，導致細節丟失
- 解決方案：先手動調整 threshold，或使用 `--no-remove-bg` 後手動去背

## 版本記錄
- 2025-01-01：初始版本，完整後處理流程
