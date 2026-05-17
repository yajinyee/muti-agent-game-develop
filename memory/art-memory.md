# 美術記憶 — 吉伊卡哇：像素大討伐

> 記錄所有美術相關的決策、參數、已知問題、生成歷史。由 Art Director 與 Sprite Generation Agent 共同維護。

**最後更新**：2025-01-01  
**更新者**：Art Director

---

## 視覺風格定義

### 核心風格
- **風格**：16-bit 像素藝術（Pixel Art）
- **參考**：吉伊卡哇官方像素周邊、復古 RPG 遊戲美術
- **整體調性**：可愛、治癒、略帶悲傷的日系風格

### 色彩規範
| 元素 | 最大色數 | 主色調 | 備註 |
|------|---------|-------|------|
| 角色 | 32 色 | 暖色系 | 依角色個性調整 |
| 目標物（普通）| 16 色 | 藍/綠色系 | 海洋生物感 |
| 目標物（特殊）| 24 色 | 金/紫色系 | 稀有感 |
| BOSS | 32 色 | 深色系 | 威脅感 |
| UI | 16 色 | 統一色板 | 不干擾遊戲畫面 |
| 背景 | 32 色 | 深藍色系 | 海底感 |

### 尺寸規範
| 元素 | 尺寸 | 格式 | 備註 |
|------|------|------|------|
| 吉伊卡哇 LV1-3 | 64x64 px | PNG（透明背景）| |
| 小八 LV4-7 | 96x96 px | PNG（透明背景）| |
| 烏薩奇 LV8-10 | 128x128 px | PNG（透明背景）| |
| 目標物（普通）| 64x64 px | PNG（透明背景）| |
| 目標物（特殊）| 64x64 px | PNG（透明背景）| |
| BOSS B001 | 256x256 px | PNG（透明背景）| |
| UI 按鈕 | 可變 | PNG（透明背景）| |

### 輪廓規範
- 所有角色與目標物必須有 **1px 黑色輪廓**
- 輪廓顏色：#1A1A1A（非純黑，避免過硬）
- 輪廓方式：外輪廓（不影響內部細節）

### 光源規範
- 光源方向：**左上方 45 度**
- 高光：白色或淺色（透明度 50-70%）
- 陰影：深色（透明度 30-50%）

---

## ComfyUI 生成參數記錄

### 角色生成（最佳參數）
```json
{
  "model": "chiikawa-pixel-v2.safetensors",
  "sampler": "DPM++ 2M Karras",
  "steps": 28,
  "cfg_scale": 7.5,
  "width": 512,
  "height": 512,
  "positive_prompt": "chiikawa character, pixel art, 16-bit style, cute, simple background, transparent background, 1px outline, [CHARACTER_SPECIFIC]",
  "negative_prompt": "realistic, 3d, blurry, low quality, watermark, text, extra limbs"
}
```

### 目標物生成（最佳參數）
```json
{
  "model": "pixel-fish-v1.safetensors",
  "sampler": "Euler a",
  "steps": 20,
  "cfg_scale": 6.0,
  "width": 256,
  "height": 256,
  "positive_prompt": "pixel art fish, cute, 16-bit, transparent background, simple design, [FISH_TYPE]",
  "negative_prompt": "realistic, complex, dark, scary"
}
```

### 後處理流程
1. 縮放到目標尺寸（Nearest Neighbor，保持像素感）
2. 調色板限制（最多 32 色）
3. 去背（透明背景）
4. 加輪廓（1px 黑色）
5. 輸出 PNG

---

## 角色設計記錄

### 吉伊卡哇（ちいかわ）
- **外觀**：小型白色動物，圓形頭部，大眼睛，小耳朵
- **等級差異**：LV1 最小最可愛，LV3 稍大稍強壯
- **主色**：白色（#F5F5F5）、粉色（#FFB3C1）
- **特色**：表情豐富，略帶憂鬱感
- **攻擊動作**：揮動小手/武器

### 小八（ハチワレ）
- **外觀**：貓咪，頭頂有八字形花紋
- **等級差異**：LV4 基礎，LV7 更有氣勢
- **主色**：白色（#FFFFFF）、灰色（#CCCCCC）、黑色（#333333）
- **特色**：八字形頭頂花紋是識別特徵
- **攻擊動作**：貓爪攻擊

### 烏薩奇（うさぎ）
- **外觀**：兔子，長耳朵，神秘感
- **等級差異**：LV8 基礎，LV10 最強形態
- **主色**：白色（#FAFAFA）、淡紫色（#E8D5F5）
- **特色**：長耳朵，眼神銳利
- **攻擊動作**：魔法攻擊

---

## 目標物設計記錄

### 設計原則
- 所有目標物都是海洋生物或海底相關物品
- 普通魚類：常見魚種，顏色鮮豔
- 中型魚：較大型魚類，有特色
- 大型魚：稀有魚種，視覺衝擊
- 特殊目標：非魚類（螃蟹、章魚、寶箱等）
- BOSS：巨大海洋生物，威脅感強

### 已完成目標物（部分記錄）
| ID | 名稱 | 類型 | 倍率 | 狀態 |
|----|------|------|------|------|
| T001 | 小金魚 | 普通魚 | 1x | ✅ |
| T002 | 小藍魚 | 普通魚 | 1x | ✅ |
| ... | ... | ... | ... | ... |
| B001 | 深海巨獸 | BOSS | 100-500x | ✅ |

---

## 已知美術問題

### 待修復
- 無已知嚴重問題

### 已修復
- 無記錄

---

## 美術資產目錄

```
client/chiikawa-pixel/assets/
├── characters/
│   ├── chiikawa_lv1.png ~ chiikawa_lv3.png
│   ├── hachiware_lv4.png ~ hachiware_lv7.png
│   └── usagi_lv8.png ~ usagi_lv10.png
├── targets/
│   ├── T001.png ~ T105.png
│   └── B001_boss.png
├── backgrounds/
│   ├── main_bg.png
│   └── bonus_bg.png
├── ui/
│   ├── btn_normal.png
│   ├── btn_active.png
│   └── btn_auto.png
└── audio/
    └── [音效檔案]
```

---

## 版本記錄

| 版本 | 日期 | 變更 |
|------|------|------|
| 1.0 | 2025-01-01 | 初始記錄，基於現有資產 |
