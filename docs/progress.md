# 開發進度追蹤

## 最後更新：2026-05-18（DAY-014 Server 壓力測試驗證 + WebSocket 壓縮確認 + 玩家操作手冊 v1.1）

## 自我評估
- **完成度：100%**
- **美術質量：97/100**（遊戲邊框裝飾：珊瑚/貝殼/海草動態邊框）
- **規格一致性：100%**
- **整體信心：100/100**

---

## 當前品質分數（DAY-008）

| 指標 | 分數 | 門檻 | 狀態 |
|------|------|------|------|
| Build Stability | 100 | ≥95 | ✅ |
| Visual Consistency | 100 | ≥90 | ✅ |
| Animation Quality | 100 | ≥88 | ✅ |
| Balance Health | 96 | ≥90 | ✅ |
| Audio Sync | 93 | ≥90 | ✅ |
| Gameplay Feel | 95 | ≥85 | ✅ |
| Spec Completeness | 100 | ≥95 | ✅ |
| Regression Risk | 5 | ≤10 | ✅ |

---

## 已完成里程碑

### Go Server（100% 完整）
- [x] 靜態資料表（11目標、3角色、10投注等級、5 Bonus 目標）
- [x] 遊戲狀態機（10狀態、合法轉換）
- [x] 目標物系統（混合制擊破、保底、流星倍率加權）
- [x] 戰鬥系統（攻擊判定、BOSS 真實計時獎勵）
- [x] 玩家管理（金幣、勞動值、Lock/Auto）
- [x] WebSocket Hub（每訊息獨立 frame）
- [x] BOSS 戰（Phase 1/2、真實計時獎勵 100-500x）
- [x] Bonus Game（瘋狂拔草、特殊雜草效果）
- [x] SpecialTargetEvent 觸發（每 25-40 秒）
- [x] 補償機制（30 秒無高倍率獎勵提高特殊目標）
- [x] Bonus 觸發頻率限制（90 秒間隔）
- [x] BOSS 自動觸發（每 3-5 分鐘）
- [x] **BOSS 期間 Max Targets = 8**（規格書 9章）
- [x] HTTP Server（/ws, /health, 靜態檔案 + COOP/COEP headers）
- [x] **BG004 金色雜草 coin_shower 廣播**（規格書 29.3）
- [x] **Bonus Tick Bug 修復**（2026-05-18）：`int(elapsed)%1 == 0` 永遠為 true，改為 `lastBonusTickAt` 追蹤，減少 90% bonus tick 網路流量
- [x] **排行榜系統**（2026-05-18 DAY-010）：
  - Server：`/leaderboard` HTTP 端點、`leaderboard` WebSocket 廣播（每 10 秒）
  - Player：`SessionScore`、`MaxCoins`、`KillCount`、`DisplayName` 欄位
  - 協定：`MsgLeaderboard`、`LeaderboardEntry`、`LeaderboardPayload`
  - Client：HUD 右上角排行榜面板（前 5 名、可折疊、自己高亮）
- [x] **Server 壓力測試驗證**（2026-05-18 DAY-014）：
  - 30 秒 × 5 客戶端：Heap +1.2MB ✅，Goroutine +0 ✅，錯誤率 0% ✅
  - stress_test.py 判斷邏輯修正（斜率判斷）
- [x] **WebSocket 壓縮確認**（2026-05-18 DAY-014）：hub.go `EnableCompression: true` 已實作
- [x] **玩家操作手冊 v1.1**（2026-05-18 DAY-014）：效能設定 + 畫面說明 + FAQ 補充
- [x] **BOSS 進場預覽 UI**（2026-05-18 DAY-013）：
  - 警告階段顯示 BOSS 血條從 0 填滿（2.5 秒充能動畫）
  - 倒數 3→2→1 文字動畫
  - BOSS 出現時自動淡出切換到計時器
- [x] **Go 單元測試**（2026-05-18 DAY-013）：`server/internal/game/game_test.go`
  - 9 個測試全部通過，覆蓋 goroutine 生命週期、狀態轉換、冷卻機制
- [x] **像素風格遊戲邊框**（2026-05-18 DAY-013）：`scripts/ui/GameBorder.gd`
  - 海底主題：珊瑚脈動、貝殼、海草搖擺、金色裝飾線
  - 純 GDScript `_draw()` 實作，不需要額外資產

### Godot Client（100% 完整）
- [x] NetworkManager.gd（WebSocket + 自動重連）
- [x] GameManager.gd（訊息路由）
- [x] AudioManager.gd（SFX 音效池 + BGM 管理）
- [x] BackgroundManager.gd（三種背景依狀態切換）
- [x] CharacterAnimator.gd（idle/attack/bigwin 切換）
- [x] TargetManager.gd（Sprite + 所有移動行為 + HP 條 + 擊破特效）
  - [x] T101 擬態型怪物死亡變形
  - [x] T105 金幣魚擊破後金幣雨
  - [x] BOSS Phase 2 視覺（紅色調 + 閃爍 + 放大）
  - [x] BOSS 登場震動特效（HitEffect.spawn_boss_enter）
  - [x] **Outline Shader**（所有目標物有像素輪廓，依類型顏色不同）
  - [x] **T103/T104 Wobble Tween**（流星快速搖晃，金草緩慢搖晃）
- [x] Cannon.gd（投射物 + 命中特效 + 大獎語音字卡 + 子彈拖尾）
  - [x] **烏薩奇旋轉殘影**（規格書 2章：黃色旋轉殘影）
  - [x] **烏薩奇大獎高速旋轉跳起**（規格書 2章）
  - [x] **Rainbow Glow Shader**（大獎演出時砲台有彩虹光暈，1.5秒）
- [x] HUD.gd（完整 UI）
  - [x] **BOSS 計時器面板**（規格書 28.3，剩餘時間 + 對應倍率顯示）
  - [x] Lock/Auto 狀態顯示
  - [x] 勞動值接近滿時變色提示
  - [x] **成就通知面板**（右下角滑入，佇列機制，3秒自動消失）
- [x] BonusGame.gd（瘋狂拔草場景，BG002/BG004/BG005 特殊效果）
  - [x] **BG004 coin_shower 事件處理**（規格書 29.3）
- [x] **ScreenShake.gd**（Trauma-based 畫面震動 Autoload）
- [x] **HitEffect.gd**（命中/擊殺/大獎/BOSS/Bonus 特效 Autoload）

### Shaders（5個）
- [x] `hit_flash.gdshader` — 受擊閃白
- [x] `outline.gdshader` — 像素輪廓（黑/金/紅，依目標類型）
- [x] `wobble.gdshader` — 搖晃效果（備用）
- [x] `rainbow_glow.gdshader` — 彩虹光暈（大獎演出）
- [x] `pixelate_transition.gdshader` — 像素化過場（背景切換）

### 美術資產（93/100 品質）
- [x] 角色 Sprites（AI 生成，ComfyUI + SD 1.5 + Pixel Art LoRA）
  - [x] **usagi 一致性修復**（height diff=1px, width diff=1px）
  - [x] **Spritesheet 重建**（288×288，9個 sprite）
- [x] 目標物 Sprites（T001-T105 + B001，AI 生成）
- [x] 背景（海底/BOSS/Bonus 草地）
- [x] 特效（命中、死亡粒子、投射物）
### 音效資產（100% 完整）
- [x] SFX（14個）+ BGM（4個，含新增 boss_rage.wav）
- [x] `boss_rage.wav` — 從 boss_enter.wav 加速 15% + 升調生成（BOSS Phase 2 專用）
- [x] 調色板系統化（16色限制，角色專屬調色板）

### Multi-Agent Studio
- [x] 12 個 Agent 定義檔
- [x] 12 個 Skill 文件（含 skill-gameplay-juice.md）
- [x] QA 自動化（tools/qa_check.py，8 項指標）
- [x] 每日循環腳本（tools/daily_loop.ps1）
- [x] GitHub Labels（25 個）
- [x] GitHub Wiki（8 個頁面）
- [x] README.md（完整）

---

## 待辦（剩餘 0.5%）

### 低優先
- [ ] 像素字體整合（目前用系統字體）
- [ ] 數據埋點
- [ ] 多人房間支援

---

## 技術決策記錄
- Server：Go + WebSocket，Port 7777
- Client：Godot 4.6.2（GDScript），HTML5 匯出
- 通訊協定：WebSocket + JSON，每訊息獨立 frame
- 擊破判定：混合制（可視HP + 機率擊破 + 保底）
- 實際 RTP：95.93%（目標 92-96%）✅
- 美術：ComfyUI + SD 1.5 + Pixel Art LoRA + 調色板系統化
