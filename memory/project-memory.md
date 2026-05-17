# 專案記憶 — 吉伊卡哇：像素大討伐

> 本文件是整個 Multi-Agent Studio 的共享記憶。記錄專案的當前狀態、重要決策、已知問題、技術細節。所有 Agent 都應在開始工作前讀取本文件。

**最後更新**：2025-01-01  
**更新者**：Game Director

---

## 專案基本資訊

| 項目 | 內容 |
|------|------|
| 遊戲名稱 | 吉伊卡哇：像素大討伐 |
| 遊戲類型 | 捕魚機（Fish Shooting Game）|
| 開發狀態 | 99% 完成 |
| 美術質量 | 91/100 |
| 規格一致性 | 95% |
| Server 技術 | Go + WebSocket，Port 7777 |
| Client 技術 | Godot 4.6.2，HTML5 匯出 |
| 目標平台 | Web（HTML5）|

---

## 技術架構

### Server（Go）
- **語言**：Go 1.21+
- **通訊**：WebSocket（RFC 6455）
- **Port**：7777
- **路徑**：`server/`
- **入口**：`server/main.go`
- **狀態**：運行中（需確認最新編譯狀態）

### Client（Godot）
- **引擎**：Godot 4.6.2
- **語言**：GDScript
- **匯出**：HTML5
- **路徑**：`client/chiikawa-pixel/`
- **主場景**：`Main.tscn`
- **Bonus 場景**：`BonusGame.tscn`

### 通訊協定
- **版本**：1.0.0
- **格式**：UTF-8 JSON
- **心跳**：30 秒
- **詳細規格**：`docs/protocol-change-policy.md`

---

## 遊戲內容

### 角色系統
| 角色 | 等級 | 特色 | 攻擊音效 |
|------|------|------|---------|
| 吉伊卡哇 | LV1-3 | 基礎角色，可愛小動物 | attack_fire.wav |
| 小八（ハチワレ）| LV4-7 | 中階角色，貓咪 | attack_fire_hachiware.wav |
| 烏薩奇（うさぎ）| LV8-10 | 高階角色，兔子 | attack_fire_usagi.wav |

### 目標物系統
| 類型 | ID 範圍 | 數量 | 倍率範圍 | 出現頻率 |
|------|---------|------|---------|---------|
| 普通魚 | T001-T030 | 30 種 | 1-3x | 高 |
| 中型魚 | T031-T060 | 30 種 | 3-8x | 中 |
| 大型魚 | T061-T090 | 30 種 | 8-20x | 低 |
| 特殊目標 | T091-T105 | 15 種 | 20-50x | 極低 |
| BOSS | B001 | 1 種 | 100-500x | 稀有 |

**總計**：105 種普通目標物 + 1 種 BOSS = 106 種

### 音效資產
| 檔案 | 用途 | 狀態 |
|------|------|------|
| attack_fire.wav | 吉伊卡哇攻擊 | ✅ 完成 |
| attack_fire_hachiware.wav | 小八攻擊 | ✅ 完成 |
| attack_fire_usagi.wav | 烏薩奇攻擊 | ✅ 完成 |
| big_win.wav | 大獎 | ✅ 完成 |
| bonus_game.wav | Bonus BGM | ✅ 完成 |
| bonus_ready.wav | Bonus 準備 | ✅ 完成 |
| boss_enter.wav | BOSS 登場 | ✅ 完成 |
| boss_warning.wav | BOSS 警告 | ✅ 完成 |
| coin_drop.wav | 硬幣掉落 | ✅ 完成 |
| hit.wav | 命中 | ✅ 完成 |
| kill.wav | 擊殺 | ✅ 完成 |
| main_game.wav | 主遊戲 BGM | ✅ 完成 |
| reward_bag.wav | 獎勵袋 | ✅ 完成 |
| weed_pull.wav | 拔草 | ✅ 完成 |

---

## 已完成的重要里程碑

| 里程碑 | 完成日期 | 說明 |
|-------|---------|------|
| AI 角色圖生成 | - | 吉伊卡哇、小八、烏薩奇各等級圖像 |
| 目標物 AI 生成 | - | T001-T105 全部 11 種類型 |
| RTP 校正 | - | 使用蒙地卡羅模擬校正至 92-96% |
| 多幀動畫 | - | 角色與目標物動畫實作完成 |
| WebSocket 通訊 | - | Server-Client 雙向通訊完成 |
| Bonus 遊戲 | - | 基礎流程完成 |
| BOSS 戰 | - | B001 完整流程完成 |
| HTML5 匯出 | - | Godot 4 HTML5 匯出設定完成 |

---

## 當前品質分數

| 指標 | 分數 | 門檻 | 狀態 |
|------|------|------|------|
| Spec Completeness | 95 | >=95 | ✅ |
| Build Stability | 待測 | >=95 | ⏳ |
| Visual Consistency | 91 | >=90 | ✅ |
| Animation Quality | 待測 | >=88 | ⏳ |
| Audio Sync | 待測 | >=90 | ⏳ |
| Gameplay Feel | 待測 | >=85 | ⏳ |
| Balance Health | 待測 | >=90 | ⏳ |
| Regression Risk | 待測 | <=10 | ⏳ |

---

## 已知問題與注意事項

### 技術注意事項
1. **Kiro CLI 路徑**：`C:\Program Files\Kiro-Cli\kiro-cli.exe`（已驗證 2026-05-07）
2. **中文編碼**：使用 cmd.exe + chcp 65001 確保中文正確
3. **Go 編譯**：每次修改後必須執行 `go build ./...` + `go vet ./...`
4. **HTML5 測試**：必須在 Chrome/Firefox 最新版測試

### 待解決問題
- 目前無已知嚴重問題

---

## 重要決策記錄

| 日期 | 決策 | 理由 | 決策者 |
|------|------|------|-------|
| 2025-01-01 | 建立 Multi-Agent Studio 架構 | 提升開發效率與品質管理 | Game Director |

---

## 下一步重點

1. 執行完整品質評估，建立所有指標的基準分數
2. 根據評估結果，優先修復最低分項目
3. 目標：美術質量從 91 提升到 95+

---

*本文件由所有 Agent 共同維護，每次重要變更後更新*
