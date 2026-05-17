---
name: pixel-art-drawing
description: 像素美術繪製技巧與常見錯誤。當需要用 Python Pillow 生成像素圖，或修正美術問題時使用。
---

# 像素美術繪製 Skill

## 核心設計原則（來源：pixnote.net）

### 角色比例
- **chibi 比例**：頭佔身體 50-60%，讓角色更可愛
- **32x32 標準**：頭 16px，身體 8px，腳 8px
- **64x64 詳細版**：頭 32px，身體 16px，腳 16px

### 眼睛（最重要）
```python
def draw_eye_2x2(img, x, y, pupil_color):
    """2x2 眼睛 + 左上高光"""
    px(img, x,   y,   pupil_color)
    px(img, x+1, y,   pupil_color)
    px(img, x,   y+1, pupil_color)
    px(img, x+1, y+1, pupil_color)
    px(img, x,   y,   (255,255,255,255))  # 左上高光
```
- 8x8 角色：1x1 眼睛
- 16x16 角色：1x2 眼睛
- 32x32 角色：2x2 眼睛 + 高光
- 高光永遠在左上角

### 嘴巴（常見錯誤）
```python
# 正確：微笑 = 兩端高，中間低
def draw_smile(img, x, y, w, color):
    px(img, x,     y,   color)  # 左端（高）
    px(img, x+w-1, y,   color)  # 右端（高）
    for i in range(1, w-1):
        px(img, x+i, y+1, color)  # 中間（低）

# 錯誤：兩端低中間高 = 看起來像鬍子或皺眉
```

### 輪廓（剪影辨識）
- 每個角色靠剪影就能辨識
- 吉伊卡哇：圓耳 + 圓頭
- 小八：尖耳 + 條紋
- 烏薩奇：長耳 + 長身

### 顏色
- 每個角色 3-5 個顏色
- 必須有：主色、暗色（輪廓）、亮色（高光）
- 腮紅用粉紅色小圓點（2px 半徑）

## 常見錯誤與修正

| 問題 | 原因 | 修正 |
|------|------|------|
| 嘴巴像鬍子 | draw_smile 方向錯誤 | 兩端高中間低 |
| 角色太大擋畫面 | scale 設太高 | 32x32 用 2x，不要 3x |
| 背景沒顯示 | tscn 沒設 texture | 在 _ready() 動態載入 |
| 眼睛沒生命感 | 缺少高光 | 左上角加 1px 白色高光 |
| 像素模糊 | 線性濾波 | TEXTURE_FILTER_NEAREST |

## Python Pillow 像素繪製模板

```python
from PIL import Image

def new_img(w, h):
    return Image.new("RGBA", (w, h), (0, 0, 0, 0))

def px(img, x, y, color):
    """安全設定像素"""
    if 0 <= x < img.width and 0 <= y < img.height:
        img.putpixel((x, y), color)

def fill_circle(img, cx, cy, r, color):
    """填充圓形"""
    for y in range(cy-r, cy+r+1):
        for x in range(cx-r, cx+r+1):
            if (x-cx)**2 + (y-cy)**2 <= r**2:
                px(img, x, y, color)
```

## 驗證流程
1. 生成 PNG 後，在 Godot 實際跑起來截圖
2. 截圖確認：角色辨識度、嘴巴方向、眼睛高光、大小比例
3. 發現問題立刻修正，不要等到最後
