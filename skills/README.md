# Skills 索引

> 維護者：Skill Librarian  
> 最後更新：2026-05-17  
> 總 Skill 數：6

---

## 概覽

Skills 目錄包含所有可重用的技術知識文件。每個 Skill 文件記錄了特定技術領域的最佳實踐、程式碼範例和已知問題解決方案。

---

## Skill 清單

### 美術與動畫

| Skill | 檔案 | 說明 | 適用 Agent |
|-------|------|------|-----------|
| Animation Consistency | `skill-animation-consistency.md` | 像素藝術動畫幀一致性技術（shared_scale、bottom_align、keep_largest_component）| Animation Agent, Sprite Generation Agent |
| ComfyUI Sprite Generation | `skill-comfyui-sprite-generation.md` | 使用 ComfyUI 生成像素藝術 Sprite 的完整流程 | Sprite Generation Agent |
| Godot Animation Import | `skill-godot-animation-import.md` | 將 spritesheet 匯入 Godot 4 AnimatedSprite2D 的方法 | Godot Client Agent, Animation Agent |
| Process Sprites | `skill-process-sprites.md` | Sprite 後處理技術（去背、縮放、合成）| Sprite Generation Agent, Animation Agent |

### 開發工具

| Skill | 檔案 | 說明 | 適用 Agent |
|-------|------|------|-----------|
| Git Windows Permissions | `skill-git-windows-permissions.md` | 解決 Windows 上 Git 權限問題（icacls、tmpdir）| 所有 Agent |
| RTP Simulation | `skill-rtp-simulation.md` | 捕魚機 RTP 模擬與校正方法 | Balance Agent |

---

## 使用指南

### 如何使用 Skill

1. 在開始任務前，先查看相關 Skill 文件
2. 遇到問題時，先搜尋 Skill 目錄是否有解決方案
3. 解決新問題後，更新或建立對應的 Skill 文件

### 如何建立新 Skill

1. 命名格式：`skill-<kebab-case-name>.md`
2. 必須包含：概覽、問題描述、解決方案、程式碼範例
3. 建立後更新本 README

### Skill 品質標準

- 每個 Skill 必須有可執行的程式碼範例
- 必須記錄已知問題和解決方案
- 定期更新（至少每月一次）

---

## 知識地圖

```
遊戲開發知識
├── 美術生成
│   ├── skill-comfyui-sprite-generation.md
│   ├── skill-process-sprites.md
│   └── skill-animation-consistency.md
├── 遊戲引擎
│   └── skill-godot-animation-import.md
├── 數值設計
│   └── skill-rtp-simulation.md
└── 開發工具
    └── skill-git-windows-permissions.md
```

---

## 待建立的 Skill

| 優先級 | Skill 名稱 | 說明 |
|-------|-----------|------|
| 🟠 P1 | skill-websocket-godot4.md | Godot 4 WebSocket 最佳實踐 |
| 🟠 P1 | skill-audio-sync-godot4.md | Godot 4 音效同步技術 |
| 🟡 P2 | skill-html5-optimization.md | HTML5 效能優化 |
| 🟡 P2 | skill-go-server-patterns.md | Go Server 設計模式 |

---

*本文件由 Skill Librarian 維護，每次新增 Skill 後更新*
