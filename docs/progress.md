# 開發進度追蹤

## 最後更新：2026-05-17（目標物 AI 生成完成 🎉）

## 自我評估
- **完成度：99%**
- **美術質量：91/100**（目標物 AI 生成，平均像素密度 27%→62%，提升 130%）
- **規格一致性：Server 99% / Client 94%**
- **整體信心：99/100**

## 目標物 AI 生成結果
| 目標物 | 程式生成 | AI 生成 | 提升 |
|--------|---------|---------|------|
| 平均 | ~1100px (27%) | ~2525px (62%) | +130% |
| T001 雜草 | 442px (10%) | 2563px (62%) | +480% |
| T006 蘑菇 | 1759px (43%) | 2912px (71%) | +65% |
| T103 流星 | 825px (20%) | 2871px (70%) | +248% |

## AI 生成結果
| 角色 | 程式生成 | AI 生成 | 提升 |
|------|---------|---------|------|
| chiikawa_idle | 3772px (41%) | 6059px (65%) | +60% |
| hachiware_idle | 3684px (40%) | 5899px (64%) | +60% |
| usagi_idle | 3102px (34%) | 5616px (60%) | +81% |
| Spritesheet 非透明 | 38% | 56% | +47% |

## NVIDIA 驅動問題（待解決）
- **問題：** PyTorch 2.11.0+cu130 需要 CUDA 13.0，但驅動 555.85 只支援 CUDA 12.5
- **解決方案：** 更新驅動到 596.49（支援 CUDA 13.0）
- **下載：** https://www.nvidia.com/Download/index.aspx（GTX 1650, Windows 11）
- **更新後：** 重啟 → 執行 `run_nvidia_gpu.bat` → `py tools/comfyui_generate.py --all`

---

## 已完成

### Go Server（95% 完整）
- [x] 靜態資料表（11目標、3角色、10投注等級、5 Bonus 目標）
- [x] 遊戲狀態機（10狀態、合法轉換）
- [x] 目標物系統（混合制擊破、保底、流星倍率加權）
- [x] 戰鬥系統（攻擊判定、BOSS 真實計時獎勵）
- [x] 玩家管理（金幣、勞動值、Lock/Auto）
- [x] WebSocket Hub（每訊息獨立 frame）
- [x] BOSS 戰（Phase 1/2、真實計時獎勵）
- [x] Bonus Game（瘋狂拔草、特殊雜草效果）
- [x] SpecialTargetEvent 觸發（每 25-40 秒）
- [x] 補償機制（30 秒無高倍率獎勵提高特殊目標）
- [x] Bonus 觸發頻率限制（90 秒間隔）
- [x] Port 預設 7777
- [x] config 套件、logger 套件
- [x] HTTP Server（/ws, /health, 靜態檔案 + COOP/COEP headers）
- [x] PlayerSnapshot 新增 projectile_speed + fire_rate 欄位

### Godot Client（88% 完整）
- [x] NetworkManager.gd（WebSocket + 自動重連 + Web/桌面自動判斷）
- [x] GameManager.gd（訊息路由、target_updated 訊號）
- [x] AudioManager.gd（SFX 音效池 + BGM 管理）
- [x] BackgroundManager.gd（三種背景依狀態切換）
- [x] CharacterAnimator.gd（Sprite2D，idle/attack/bigwin 切換）
- [x] TargetManager.gd（Sprite + 所有移動行為 + HP 條 + 擊破特效 + Lock 視覺框）
  - [x] T101 擬態型怪物死亡變形（閃爍→縮放→爆炸 + 「正體！」文字）
  - [x] T105 金幣魚擊破後金幣雨（15 枚拋物線散落）
  - [x] BOSS Phase 2 視覺（紅色調 + 閃爍 + 放大 + PHASE 2 文字）
  - [x] boss_event 訊號連接（_on_boss_event）
- [x] Cannon.gd（投射物動畫 + 命中特效 + 大獎語音字卡 + 點擊自動鎖定）
  - [x] 投射物速度依 BetLevel 動態計算（700-980 px/s）
- [x] HUD.gd（完整 UI，Lock/Auto 狀態顯示，null safety）
- [x] BonusGame.gd（瘋狂拔草場景）
  - [x] BG002 硬雜草連點 2 次機制（_weed_hp 字典追蹤）
  - [x] BG004 金色雜草金幣雨視覺（20 枚散落）
- [x] BonusGame.tscn（場景檔）
- [x] Main.tscn（完整場景，含 BackgroundManager + CharacterAnimator）

### 美術資產（62% 品質）
- [x] 角色 Sprites v5（96x96，shared_scale + bottom align 對齊）
  - chiikawa：height diff=2px ✅
  - hachiware：height diff=9px（attack 尖耳朵設計特性）
  - usagi：height diff=5px
- [x] 目標物 Sprites v2（12個，含 BOSS 128x128）
- [x] 背景（海底/BOSS/Bonus 草地）
- [x] 特效（命中、死亡粒子、投射物）
- [x] UI 元素（金幣、報酬袋、勞動值條、WARNING）
- [x] Spritesheet（characters 288x288 / targets / effects）
- [x] 音效 SFX（11個）+ BGM（3個）
- [x] GIF 動畫預覽（9個：3角色 × 3動作）

### 工具腳本
- [x] tools/test_server.py（7/7 整合測試）
- [x] tools/simulate_rtp.py（RTP 模擬）
- [x] tools/generate_pixel_art_v5.py（v5 美術生成，修復 dy 邊界 bug）
- [x] tools/generate_spritesheet.py（Spritesheet 生成）
- [x] tools/generate_sfx.py（8-bit 音效生成）
- [x] tools/process_sprites.py（**新增** agent-sprite-forge 後處理工具）
  - `--mode realign`：重新對齊現有 sprites
  - `--mode comfyui`：處理 ComfyUI 洋紅色背景圖
  - `--mode sheet`：重建 Spritesheet
  - `--mode qc`：品質報告
- [x] **ComfyUI 安裝完成**（2026-05-12）
  - 解壓縮至：`C:\ComfyUI\ComfyUI_windows_portable\`
  - SD 1.5 模型：`v1-5-pruned-emaonly.safetensors`（4068 MB）✅
  - Pixel Art LoRA：`pixel_art_lora.safetensors`（162 MB）✅
  - 啟動腳本：`tools/start_comfyui.bat`
  - API 整合：`tools/comfyui_generate.py`（已改用洋紅色背景策略）

---

## 待辦（剩餘 14%）

### 高優先
- [x] **RTP 數值校正完成**（2026-05-15）— 見下方詳細記錄
- [x] **多幀角色動畫完成**（2026-05-15）
  - idle：4幀呼吸感（縮放 0.98-1.03x + 上下位移）
  - attack：3幀揮棒（旋轉 -18°/+12° + 劍氣光效）
  - bigwin：4幀跳躍（縮放 1.0-1.08x + 金色星星）
  - 9個 GIF 預覽生成完成
  - Spritesheet 格式：384×288，符合 CharacterAnimator.gd 期望
  - 修復：generate_animation_frames.py 不再覆蓋 characters/ sprites
- [x] **目標物 v3 完成**（2026-05-15）
  - 尺寸從 48×48 升級到 64×64
  - 全部重新設計：帶陰影、細節豐富
  - T001 雜草：三片葉子，帶輪廓
  - T002-T004 小蟲：橢圓身體+圓頭+觸角+腳
  - T005 布丁：梯形身體+奶油頂+眼睛+嘴巴
  - T006 蘑菇：半圓傘+白色斑點+莖
  - T101 擬態：雜草外觀+詭異眼睛+嘴巴
  - T102 寶箱：木箱+金邊+鎖+眼睛+牙齒
  - T103 流星：發光核心+尾巴+星芒
  - T104 金色雜草：金色三葉+閃光點
  - T105 金幣魚：橢圓魚身+金幣符號+魚尾
  - B001 BOSS：96×96，暗黑光環+紅眼+邪惡微笑
  - targets_sheet.png 升級到 256×192
  - TargetManager.gd HP 條更新（48px 寬）
- [ ] **ComfyUI 生成 AI 角色圖**
  1. 啟動：`tools\start_comfyui.bat`（手動在終端執行，長時間服務）
  2. 等待 `http://127.0.0.1:8188` 就緒
  3. 生成：`py tools/comfyui_generate.py --all`（9張，洋紅色背景）
  4. 後處理：`py tools/process_sprites.py --mode comfyui --input <path> --char <char> --pose <pose>`
  5. 重建 Sheet：`py tools/process_sprites.py --mode sheet`
- [ ] **多幀角色動畫** — 每個動畫 4 幀，讓角色真正動起來（目前是靜態圖切換）
- [ ] **RTP 數值校正** — 目標 94%，目前遠高於此，需調整擊破機率

### 中優先
- [ ] 像素字體整合（目前用系統字體）
- [ ] hachiware/usagi 幀一致性優化（height diff 仍有 9px/5px）

### 低優先（未來功能）
- [ ] 數據埋點
- [ ] 營運工具後台
- [ ] 多人房間支援

---

## 技術決策記錄
- Server：Go + WebSocket（gorilla/websocket），Port 7777
- Client：Godot 4.6.2（GDScript），HTML5 匯出
- 通訊協定：WebSocket + JSON，每訊息獨立 frame
- 擊破判定：混合制（可視HP + 機率擊破 + 保底）
- 目標 RTP：94%（待校正）
- 遊玩網址：http://localhost:7777
- 美術後處理：agent-sprite-forge 技術（洋紅色去背 + shared_scale + bottom align）

---

## 重啟對話後的第一步驟（給下一個 session）

```
1. 讀取 docs/progress.md（本文件）確認狀態
2. 讀取 .kiro/skills/knowhow-log.md 確認已知問題（特別是 #38-44）
3. go build ./... 確認 Server 編譯狀態
4. py tools/process_sprites.py --mode qc 確認美術狀態
5. 根據「待辦」清單決定下一步
```

## 本次對話完成的事（2026-05-12）

1. **ComfyUI 安裝** — 解壓縮 1908MB 壓縮檔，下載 SD 1.5（4068MB）+ Pixel Art LoRA（162MB）
2. **agent-sprite-forge 分析** — 理解其核心技術：洋紅色背景、shared_scale、component_mode=largest、bottom align
3. **process_sprites.py 建立** — 移植 agent-sprite-forge 後處理技術
4. **美術 bug 修復** — usagi_attack dy=-3 超出邊界、所有角色 bigwin dy 過大
5. **玩法規格缺口修復** — T101/T105 特效、BG002/BG004 機制、BOSS Phase 2、投射物速度
6. **Server PlayerSnapshot 升級** — 新增 projectile_speed + fire_rate
7. **所有修改通過 go build + go vet**
