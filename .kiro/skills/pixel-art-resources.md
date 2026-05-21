---
name: pixel-art-resources
description: 免費像素美術資源清單與下載方式。當需要提升遊戲美術質量時使用。包含角色、怪物、背景、UI 的免費素材來源。
---

# 免費像素美術資源 Skill

## 核心原則
網路上有大量免費高品質像素素材，不需要全部自己畫。
優先找 CC0 或 CC-BY 授權的素材，可商用且不需標註。

## 主要素材來源

### OpenGameArt.org（最推薦，大量 CC0）
- 網址：https://opengameart.org
- 搜尋關鍵字：`pixel art cute character 16x16`、`fish monster sprite`
- 授權：多數 CC0 或 CC-BY，可商用
- 下載格式：PNG spritesheet

### itch.io 免費素材
- 網址：https://itch.io/game-assets/free/tag-16-bit/tag-sprites
- 魚類素材：https://itch.io/game-assets/newest/tag-fish/tag-sprites
- 怪物素材：https://blacis.itch.io/pixel-monsters-mega-pack
- 注意：確認每個素材的授權條款

### 推薦具體素材
| 素材 | 網址 | 授權 | 用途 |
|------|------|------|------|
| 16-bit Cute Character | https://opengameart.org/content/16-bit-cute-character | CC-BY 3.0 | 角色基底 |
| Cute Sprites Pack 1 | https://opengameart.org/content/cute-sprites-pack-1 | CC0 | 可愛怪物 |
| Pixel Monsters Megapack | https://blacis.itch.io/pixel-monsters-mega-pack | 免費商用 | 怪物目標 |
| Fish Sprite Bundle | https://llenpix.itch.io/fish-sprite | 免費 | 魚類目標 |
| LPC Character Generator | https://pflat.itch.io/lpc-character-generator | CC-BY | 角色生成 |

## 下載工具
使用 `tools/download_assets.py` 自動下載並整理素材。

## 使用流程
1. 從上方來源下載 PNG spritesheet
2. 放入 `client/chiikawa-pixel/assets/sprites/downloads/`
3. 執行 `tools/process_assets.py` 裁切並轉換格式
4. 替換 `assets/sprites/characters/` 或 `assets/sprites/targets/` 中的佔位圖
5. 重新執行 `tools/generate_spritesheet.py` 更新 Spritesheet
6. 重新匯出 HTML5

## 美術優化優先順序
1. **角色（最重要）**：吉伊卡哇/小八/烏薩奇 — 找可愛圓頭角色素材替換
2. **BOSS**：找大型怪物素材
3. **目標物**：找各種小怪物/生物素材
4. **背景**：找海底場景 tileset
5. **UI**：找像素風 UI 框架

## KnowHow
- OpenGameArt 的 LPC（Liberated Pixel Cup）系列品質最高，且全部 CC-BY
- itch.io 免費素材需要逐一確認授權，不是全部都能商用
- 32x32 比 16x16 更清晰，Godot 縮放後效果更好
- 下載 spritesheet 後用 Pillow 裁切比手動裁切快很多
