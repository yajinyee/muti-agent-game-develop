# Skill：像素藝術品質提升技術（2026 最新）

> 來源：[binary.ph - Pixel Art Mastery for Modern Games](https://binary.ph/2026/05/16/pixel-art-mastery-for-modern-games-color-clean-lines-and-lighting-that-scale/)  
> 研究日期：2026-05-21（DAY-005）  
> 記錄者：Research Agent / Skill Librarian

---

## 核心原則：三色陰影系統

每個材質必須有三個色調：

```python
# 正確的三色陰影設定
BASE    = (r, g, b)           # 主色
SHADOW  = (r-40, g-40, b-40)  # 陰影（定義體積）
HIGHLIGHT = (r+30, g+30, b+30) # 高光（定義邊緣）
```

**光源方向必須一致**：所有角色和目標物都用左上方光源。

---

## 輪廓清潔：避免「Doubles」問題

「Doubles」= 輪廓線上多餘的像素，讓形狀看起來厚重不清晰。

### 解決方法

1. **曲線控制序列**：`1-1-2-3-2-1-1`（像素數量）
   - 這個序列讓曲線看起來自然，不會有鋸齒感
   
2. **對稱性**：左右對稱的角色用鏡像工具

3. **最終縮放檢查**：在遊戲實際大小（96×96）下確認輪廓清晰

### 在 Python 中實作

```python
def draw_clean_outline(img, cx, cy, r, color):
    """使用 1-1-2-3-2-1-1 序列畫清潔輪廓"""
    # 避免在同一位置畫兩個像素（doubles）
    outline_pixels = set()
    for angle in range(0, 360, 1):
        rad = math.radians(angle)
        x = int(cx + r * math.cos(rad))
        y = int(cy + r * math.sin(rad))
        outline_pixels.add((x, y))
    for (x, y) in outline_pixels:
        img.putpixel((x, y), color)
```

---

## 調色板策略：系統化而非隨機

**原則**：調色板是「系統」，不是「一次性選擇」。

### 吉伊卡哇角色調色板（已驗證）

```python
# 吉伊卡哇
CHIIKAWA_PALETTE = {
    "body":      (255, 252, 245),  # 主體白色
    "shadow":    (220, 215, 205),  # 陰影
    "highlight": (255, 255, 255),  # 高光
    "outline":   (45, 25, 10),     # 輪廓深棕
    "blush":     (255, 155, 150),  # 腮紅
    "pink_rod":  (255, 130, 185),  # 討伐棒
}

# 小八
HACHIWARE_PALETTE = {
    "body":      (248, 248, 248),
    "shadow":    (210, 210, 210),
    "highlight": (255, 255, 255),
    "outline":   (25, 35, 65),
    "stripe":    (75, 115, 195),
    "blue_rod":  (95, 145, 235),
}

# 烏薩奇
USAGI_PALETTE = {
    "body":      (248, 248, 248),
    "shadow":    (210, 210, 210),
    "highlight": (255, 255, 255),
    "outline":   (55, 35, 55),
    "ear_pink":  (255, 170, 170),
    "yellow_rod":(255, 210, 25),
}
```

---

## 出貨前品質檢查清單

```
□ 所有材質使用一致的 base-shadow-highlight 三色方案？
□ 輪廓沒有 doubles 和多餘邊緣像素？
□ 光源方向在所有姿勢和幀中保持一致（左上方）？
□ 調色板在角色、目標物、UI 中統一應用？
□ Sprite 在遊戲實際大小（96×96）和移動時仍然清晰可讀？
```

---

## 應用到現有程式生成器

### generate_chars_v6.py 改善點

1. **加入 shadow 色**：目前只有 base 色，缺少陰影
2. **統一光源**：確保所有角色的高光在左上方
3. **輪廓清潔**：使用 `outline_circle` 而非 `fill_circle` 的邊緣

### generate_targets_v3.py 改善點

1. **目標物也需要三色陰影**：T001-T105 目前部分缺少陰影
2. **調色板一致性**：確保所有目標物的輪廓色統一

---

## 與現有 Skills 的關係

| Skill | 關聯 |
|-------|------|
| skill-process-sprites.md | 後處理時保留陰影資訊 |
| skill-comfyui-consistent-spritesheet-2025.md | AI 生成時指定調色板 |
| skill-animation-consistency.md | 動畫幀間陰影一致性 |

*Content was rephrased for compliance with licensing restrictions*  
*來源：binary.ph（2026-05-16）*
