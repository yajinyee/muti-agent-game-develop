---
name: knowhow-log
description: 開發過程中踩過的坑和解決方案記錄。避免重複犯錯。當遇到問題或完成除錯時更新此文件。
---

# KnowHow 經驗教訓集

**原則：** 記錄「AI 自己沒注意到，經過實際測試才發現的」知識點。

---

## 1. Kiro CLI 中文編碼
- **問題：** Kiro CLI 在 Windows 上輸出中文變亂碼
- **原因：** PowerShell 預設編碼不是 UTF-8，且 `Set-Content -Encoding UTF8` 會加 BOM
- **解決：** 使用 `cmd.exe` + `chcp 65001 >nul` + `type` pipe
- **教訓：** Windows 上任何涉及中文的外部程式呼叫，都要用 cmd.exe + chcp 65001

## 2. Kiro CLI Exit Code 1
- **問題：** Kiro CLI 正常回應但 exit code 是 1
- **原因：** 它把 `All tools are now trusted` 警告寫到 stderr，PowerShell 視為錯誤
- **解決：** 在 catch 中合併 stdout + stderr，提取有效回應
- **教訓：** 不能只看 exit code 判斷成功失敗，要看實際輸出內容

## 3. Telegram Bot Privacy Mode
- **問題：** Bot 在群組中收不到非 @ 的訊息
- **原因：** Telegram Bot 預設開啟 Privacy Mode
- **解決：** 透過 @BotFather → Bot Settings → Group Privacy → Turn off
- **教訓：** 新 Bot 加入群組前必須先關閉 Privacy Mode

## 4. node-telegram-bot-api 的 token 屬性
- **問題：** `this.bot.token` 在 TypeScript 中報錯 Property 'token' does not exist
- **原因：** node-telegram-bot-api 的型別定義沒有暴露 token 屬性
- **解決：** 在 FileHandler 中獨立傳入 botToken 參數
- **教訓：** 不要依賴第三方套件的內部屬性，用參數傳遞

## 5. Bridge 架構行不通
- **問題：** 用檔案系統橋接 Bot 和 Kiro IDE AI 的方式超時
- **原因：** Hook 的 askAgent 無法自動把回應寫回檔案
- **解決：** 直接用 Kiro CLI headless mode 作為 AI 後端
- **教訓：** 不要用間接架構，直接呼叫最簡單

## 6. SystemPrompt 太長的處理
- **問題：** 陳總的 systemPrompt 有幾千字，直接作為 shell 參數會被截斷
- **解決：** 寫入暫存檔，用 `type` pipe 傳入 kiro-cli
- **教訓：** 任何超過 2000 字元的 prompt 都要用檔案傳遞

## 7. Bot 同時只能有一個實例
- **問題：** 啟動第二個 Bot 實例會報 409 Conflict
- **原因：** Telegram 同一個 token 只允許一個 polling 連線
- **解決：** 啟動前確認沒有其他實例在跑
- **教訓：** 重啟 Bot 前必須先停掉舊的

## 8. Kiro CLI 實際安裝路徑
- **問題：** Skills 文件記錄的路徑 `$env:LOCALAPPDATA\Kiro-Cli\kiro-cli.exe` 不存在
- **原因：** 官方 Windows installer (`install.ps1`) 安裝到 `C:\Program Files\Kiro-Cli\`
- **解決：** 使用正確路徑 `C:\Program Files\Kiro-Cli\kiro-cli.exe`，或直接用 `kiro-cli`（已加入 PATH）
- **教訓：** 安裝後用 `where kiro-cli` 確認實際路徑，不要假設

## 9. npm 在企業網路的 SSL 問題
- **問題：** `npm install` 報 `UNABLE_TO_VERIFY_LEAF_SIGNATURE`
- **原因：** 企業網路有自簽憑證，Node.js 預設不信任
- **解決：** `npm config set strict-ssl false`
- **教訓：** 企業環境安裝 npm 套件前先執行此設定

## 10. PowerShell 執行原則擋住 npm
- **問題：** `npm` 指令報 `因為這個系統上已停用指令碼執行`
- **原因：** Windows 預設 PowerShell 執行原則為 Restricted
- **解決：** `Set-ExecutionPolicy -Scope CurrentUser -ExecutionPolicy RemoteSigned -Force`
- **教訓：** 新機器安裝 Node.js 後必須先設定執行原則

---

**更新規則：** 每次發現新的知識點，都要追加到這個檔案。

## 11. GramIO 比 node-telegram-bot-api 更適合 TypeScript 專案
- **發現：** `node-telegram-bot-api` 型別定義不完整（如 token 屬性問題），且更新慢
- **替代方案：** [GramIO](https://gramio.dev/) — 完全 TypeScript 原生，Bot API 類型自動生成，支援 Node.js/Bun/Deno
- **優點：** 端到端型別安全、middleware 架構、內建 i18n/scenes/plugins
- **安裝：** `npm create gramio@latest` 或 `npm install gramio`
- **教訓：** 新 Bot 專案優先考慮 GramIO，舊專案遷移時評估成本

## 12. Kiro CLI 支援 KIRO_API_KEY 環境變數認證
- **發現：** 除了 browser OAuth，Kiro CLI 也支援 `KIRO_API_KEY` 環境變數
- **用途：** CI/CD pipeline、自動化腳本、無瀏覽器環境
- **設定：** 從 Kiro portal 產生 API key，設定 `$env:KIRO_API_KEY = "your-key"`
- **教訓：** 企業帳號用 browser OAuth，自動化場景用 API key

## 13. 規格書 RTP 數值設計問題（重要）
- **問題：** 按規格書數值模擬，RTP 高達 600-1200%，遠超目標 94%
- **根本原因：**
  1. Bonus 每局觸發 5-8 次（勞動值累積太快）
  2. Bonus 獎勵 = entry_bet × 50-150x（相當於每次 Bonus 就是一個大獎）
  3. 基礎目標擊破率用 `0.92 ÷ multiplier` 計算，但每次命中都扣 bet_cost，實際 RTP 遠高於 92%
- **正確理解：** 捕魚機的 RTP 控制是「每次命中的期望獎勵 = bet_cost × RTP」，不是「擊破後獎勵 = bet_cost × multiplier」
- **修正方向：**
  1. 擊破機率應該是 `RTP × bet_cost ÷ (multiplier × bet_cost)` = `RTP ÷ multiplier`（正確）
  2. 但每次命中都扣 bet_cost，所以期望命中次數 = `multiplier ÷ RTP`
  3. Bonus 獎勵應計入總 RTP 分配，不是額外獎勵
  4. Prototype 展示版可以提高 RTP 到 200-300% 讓玩家有爽感，正式版才嚴格控制
- **教訓：** 規格書的倍率和獎勵設計是「體感設計」，不是「數值設計」，需要數值工程師另外做 RTP 模擬調整

## 14. Godot 4 像素完美設定（重要）
- **問題：** Godot 4 預設線性濾波，像素圖會模糊
- **解決：** project.godot 設定 `textures/canvas_textures/default_texture_filter=0`（Nearest）
- **額外設定：** `window/stretch/mode="canvas_items"` + `window/stretch/aspect="keep"`
- **Sprite 節點：** `texture_filter = CanvasItem.TEXTURE_FILTER_NEAREST`
- **來源：** gdquest.com, sprite-ai.art

## 15. Spritesheet 提升效能
- **原因：** 多個獨立 PNG 會造成多次 draw call，Spritesheet 只需一次
- **做法：** 用 Python Pillow 合併所有 Sprite 到一張 PNG，搭配 JSON metadata
- **Godot 使用：** `AtlasTexture` 或 `AnimatedSprite2D` 的 `SpriteFrames`
- **工具：** `tools/generate_spritesheet.py`

## 16. 像素角色設計核心原則（來源：pixnote.net）
- **圓形 = 可愛**：頭部用圓形，避免方形
- **chibi 比例**：頭佔身體 50-60%，讓角色更可愛
- **2x2 眼睛 + 左上高光**：這是讓角色有生命感的關鍵
- **靠剪影辨識**：每個角色的輪廓要不同（吉伊卡哇圓耳、小八尖耳、烏薩奇長耳）
- **3-5 色調色盤**：不要超過 5 個顏色，保持像素風格
- **逐像素繪製**：用 `putpixel` 比 `draw.ellipse` 更精確控制

## 17. BonusGame.tscn 場景建立方式
- 建立獨立 .tscn 場景，在 Main.tscn 用 `instance = ExtResource` 引入
- BonusGame 節點預設 `visible = false`，由 GDScript 控制顯示
- 雜草目標物用 Node2D 動態生成，不需要預先放在場景裡

## 18. 像素美術陰影技術（來源：pixnote.net, hitpaw.com）
- **3色陰影法**：每個顏色需要 3 個色調：亮色（高光）、中間色（基底）、暗色（陰影）
- **光源方向**：固定在左上方，右下方是陰影
- **Pillow Shading（枕頭陰影）**：邊緣加深，中心留亮，讓圓形物體有立體感
- **Dithering（抖動）**：交替兩色像素模擬漸層，避免硬邊
- **實作方式（Python Pillow）**：
  ```python
  # 3色陰影：亮/中/暗
  LIGHT = (255, 235, 215)  # 高光
  MID   = (255, 215, 190)  # 基底
  DARK  = (200, 160, 130)  # 陰影
  # 圓形左上亮、右下暗
  for y, x in circle_pixels:
      if x < cx and y < cy:  # 左上 = 亮
          color = LIGHT
      elif x > cx and y > cy:  # 右下 = 暗
          color = DARK
      else:
          color = MID
  ```
- **教訓**：沒有陰影的像素圖看起來是「平的」，加了陰影立刻有立體感

## 19. 像素美術品質等級（來源：sprite-ai.art）
- **Level 1（15-30分）**：純色填充，無陰影，幾何形狀
- **Level 2（30-50分）**：有輪廓，基本顏色區分，無陰影
- **Level 3（50-70分）**：3色陰影，有高光，角色有立體感
- **Level 4（70-85分）**：Dithering 漸層，細節豐富，動畫流暢
- **Level 5（85-100分）**：專業級，完整動畫，特效豐富
- **目前狀態**：Level 2-3 之間（55分），目標 Level 3-4（70分）

## 20. 吉伊卡哇正確顏色（來源：chiikawa.fandom.com）
- **吉伊卡哇**：白色毛皮（white fur），接近純白 RGB(255, 252, 245)，不是膚色
- **小八**：白色帶藍色條紋，耳朵內側淡藍色
- **烏薩奇**：白色，長耳朵，耳朵內側粉紅，紅色眼睛
- **錯誤**：用 (255, 220, 195) 膚色會讓角色看起來像熊，不像吉伊卡哇
- **教訓**：IP 角色一定要查官方資料確認顏色，不能憑感覺

## 21. 嘴巴位置設計原則
- 嘴巴應該在眼睛下方 5-6px，只用 3 個像素（V 形）
- 太多像素的嘴巴會看起來像鬍子
- 正確：`px(15,19), px(16,20), px(17,19)` — 3 個像素的 V 形
- 錯誤：用 draw_smile 畫 6 個像素的弧線，太大

## 22. 目標物大小問題
- 目標物太小（< 48px）在遊戲畫面中難以辨識
- 解決：生成後用 Pillow resize 到至少 48x48（NEAREST 插值保持像素感）
- 小蟲類（24x16）需要放大 2x 到 48x32

## 23. 從 Perler Bead Pattern 提取像素圖（重要技術）
- **來源**：kandipad.com 有大量 IP 角色的 Perler Bead Pattern（像素圖案）
- **下載方式**：`tools/fetch_reference.py` — 解析 HTML 找 og:image 和 cdn 圖片
- **處理方式**：`tools/process_reference.py` — 找非白色密集區域、裁切、縮小、量化
- **關鍵步驟**：
  1. 找行/列密度（非白色像素數量）
  2. 取密度 > 20% 的區域作為圖案邊界
  3. 裁切後縮小到 32x32（NEAREST 插值）
  4. 顏色量化（FASTOCTREE，10色）
  5. 放大到 64x64 輸出
- **RGBA 量化**：必須用 `method=Image.Quantize.FASTOCTREE`，MEDIANCUT 不支援 RGBA
- **教訓**：有真實參考圖就直接下載處理，不要自己憑感覺畫

## 24. 吉伊卡哇精確顏色（來源：kandipad.com perler pattern）
- 吉伊卡哇：主體 #FFFFFF，輪廓 #292A2B，腮紅/棒 #E59ED1
- 小八：主體 #FFFFFF，輪廓 #292A2B，條紋 #3370C0
- 烏薩奇：主體 #FFFFFF，輪廓 #111111，眼睛 #FF5B56
- 共同特徵：**全部是純白色主體**，不是膚色或奶油色

## 25. 從單張圖生成多狀態動畫（圖像變換技術）
- **問題**：只有 idle 參考圖，需要 attack/bigwin 狀態
- **解法**：用 PIL 圖像變換模擬不同動作
  - attack：`img.rotate(-12)` + 亮度提高 + 右上角粉紅光暈
  - bigwin：向上位移 4px + 金色色調 + 星星光點
- **關鍵**：NEAREST 插值保持像素感，不要用 BILINEAR
- **教訓**：沒有完美參考圖時，用變換生成差異版本比全部重畫快很多

## 26. PIL 圖像變換技巧
- `img.rotate(angle, expand=False, fillcolor=(0,0,0,0))` — 旋轉，透明填充
- `ImageEnhance.Brightness(img).enhance(1.2)` — 提高亮度
- `img.paste(src, (dx, dy))` — 位移
- 逐像素修改：`pixels = img.load(); pixels[x,y] = (r,g,b,a)`

## 27. PowerShell py -c 多行腳本問題
- **問題：** `py -c "..."` 在 PowerShell 執行多行 Python 時，f-string 的 `{variable}` 會被 PowerShell 解析為變數，導致錯誤或無輸出
- **原因：** PowerShell 把 `{` `}` 視為程式碼區塊，`'` 引號也有特殊意義
- **解決：** 超過 2 行的 Python 一律寫成 `.py` 檔案，用 `py tools/script.py` 執行
- **教訓：** `py -c` 只適合單行簡單指令，複雜邏輯必須用腳本檔案

## 28. 🔴 Client 射擊當掉 Bug（待修）
- **問題：** 玩家點擊射擊後，Client 無回應甚至當掉
- **可能原因（待確認）：**
  1. `_fire_projectile` 建立 Sprite2D 時 texture 載入失敗導致 null reference
  2. `try_click_target` 迭代 `_target_nodes` 時有 race condition
  3. `_spawn_hit_effect` 的 tween callback 在節點已 queue_free 後執行
  4. `Cannon.gd` 的 `_input` 在 BonusGame 狀態下仍然觸發
  5. WebSocket 訊息處理在 Client 端造成 GDScript 錯誤
- **重現步驟：** 開啟遊戲 → 點擊畫面射擊 → Client 無回應/當掉
- **優先級：** 🔴 高（影響可玩性）
- **狀態：** 待修（其他事項完成後處理）
- **修復方向：**
  - 在 `_fire_projectile` 加 null check
  - 在 `_spawn_hit_effect` 的 tween callback 加 `is_instance_valid` 檢查
  - 在 `_input` 加更嚴格的狀態檢查
  - 用 Godot 的 `push_error` 追蹤具體錯誤位置

## 28. 🔴 Client 射擊當掉 Bug — 已修復
- **根本原因（4個）：**
  1. `create_tween().set_loops()` 在 TargetManager 建立，但 tween 的目標節點（LockFrame）被刪除後 tween 繼續執行 → crash
  2. `tween_callback(node.queue_free)` 直接傳函數引用，節點已釋放時執行 → crash
  3. `get_parent().add_child()` 沒有 `is_instance_valid` 檢查
  4. `diff.normalized()` 當 diff 為零向量時返回 NaN → 旋轉異常
- **修復方式：**
  1. `set_loops()` 的 tween 改用 `frame.create_tween()` 綁定到節點，節點刪除時自動停止
  2. 所有 `tween_callback` 改為 `func(): if is_instance_valid(node): node.queue_free()`
  3. 所有 `get_parent()` 加 `is_instance_valid` 檢查
  4. 計算方向前加 `if diff.length() > 1.0` 防止零向量
- **教訓：** Godot tween 的生命週期要綁定到正確的節點，不要用 SceneTree 的 tween 操作可能被刪除的節點

## 29. Perler Pattern 格線顏色過濾技術
- **問題：** Perler Bead Pattern 圖片有格線（黃灰色），縮小後混入角色顏色
- **格線特徵：** R≈G > B，且都在 170-250 範圍（黃灰色）
- **過濾方法：** `if r > 170 and g > 170 and b < g - 20: return True`（是格線）
- **顏色校正流程：**
  1. 白色(>220) → 官方白色 #FFFFF7
  2. 深色(<80) → 官方輪廓色 #292A2B
  3. 藍色(b > r+30) → 小八條紋 #3370C0
  4. 紅色(r > 180, r > g+50) → 烏薩奇眼睛 #FF5B56
  5. 格線顏色 → 透明
- **教訓：** 從網路圖片提取像素圖時，必須先分析並過濾背景/格線顏色

## 30. 官方吉伊卡哇顏色（來源：color-hex.com Ichikawa Clan Palette）
- 主體白色：`#FFFFF7` (255, 255, 247)
- 粉紅腮紅：`#EFA5C9` (239, 165, 201)
- 輪廓深黑：`#292A2B` (41, 42, 43)
- 小八藍條紋：`#3370C0` (51, 112, 192)
- 烏薩奇紅眼：`#FF5B56` (255, 91, 86)

## 31. Flood Fill 背景去除（正確方法）
- **問題：** 簡單 threshold 去白色會把角色內部的白色毛皮也去除
- **正確方法：** 從圖片四個邊緣做 BFS flood fill，只去除與邊緣連通的白色區域
- **實作：**
  ```python
  # 從四個邊緣加入 queue
  for x in range(w):
      for y in [0, h-1]:
          if is_white_like(pixels[x,y]):
              queue.append((x,y))
  # BFS 擴展
  while queue:
      x, y = queue.popleft()
      bg_mask[x][y] = True
      for dx,dy in [(0,1),(0,-1),(1,0),(-1,0)]:
          if is_white_like(pixels[nx,ny]):
              queue.append((nx,ny))
  ```
- **結果：** 角色內部白色保留（984 pixels），背景透明（60%）
- **教訓：** 去除背景要用 flood fill，不要用全局 threshold

## 32. 角色太大的問題
- **問題：** CharacterAnimator scale=2.0 讓角色佔畫面太大
- **解決：** 改為 scale=1.5，視覺上更合適
- **教訓：** 角色大小要在遊戲畫面中測試，不能只看 PNG 檔案

## 33. 完整角色設計（頭+身體+手腳）
- **問題：** Perler Pattern 只有頭部，角色看起來不完整
- **解法：** 用程式生成完整角色（48x48 → 放大到 96x96）
- **吉伊卡哇比例：**
  - 頭：圓形半徑 14，中心在 (24, 18)
  - 身體：橢圓 10x8，中心在 (24, 33)
  - 手：圓形半徑 4，在身體兩側
  - 腳：橢圓 5x4，在身體下方
  - 耳朵：圓形半徑 5，在頭部上方兩側
- **攻擊狀態：** 右手舉起，討伐棒斜向右上
- **大獎狀態：** 整體上移 3px，加星星光點
- **教訓：** 角色設計要包含完整身體，不只是頭部

## 34. OpenCV K-means 像素藝術轉換（突破方法）
- **工具：** `tools/img_to_pixel_art.py`
- **方法：** 縮小 → K-means 顏色量化 → NEAREST 放大
- **限制：** 從 Perler Pattern 圖轉換效果差（格線顏色干擾）
- **最佳組合：** 程式生成角色（正確顏色）+ OpenCV 後處理（銳化+飽和度）
- **後處理步驟：**
  1. `ImageEnhance.Color(img).enhance(1.4)` — 提升飽和度
  2. `ImageEnhance.Contrast(img).enhance(1.2)` — 提升對比度
  3. `cv2.addWeighted(img, 1.5, gaussian, -0.5, 0)` — Unsharp mask 銳化
  4. 縮小再放大（保持像素感）
- **教訓：** 像素藝術的品質取決於原始素材品質，不是後處理能完全補救的

## 35. 美術品質突破的正確路徑
1. **最好：** 找到真正的吉伊卡哇透明 PNG → 用 OpenCV 轉像素風格
2. **次好：** 程式生成完整角色（正確比例顏色）+ OpenCV 後處理
3. **最差：** 從 Perler Pattern 圖轉換（格線干擾太嚴重）
- 目前用方法 2，品質約 70-75 分
- 要達到 85+ 分需要方法 1（真實圖片來源）

## 36. 多幀動畫 Spritesheet 系統
- **格式：** 384x288（4幀 × 3狀態 × 96px）
- **行0：** idle（4幀，4fps，上下搖擺）
- **行1：** attack（3幀，8fps，舉棒→揮下→收回）
- **行2：** bigwin（4幀，6fps，跳起→最高點→落下→彈跳）
- **Godot 實作：** 用 AtlasTexture 裁切 Spritesheet，`_process` 計時換幀
- **等待真實圖片時：** 自動降級到靜態圖（備用機制）
- **真實圖片到位後：** 把圖片放到 reference/ 資料夾，執行 generate_animation_frames.py 即可

## 37. 動畫系統設計原則
- 攻擊動畫要快（8fps），idle 要慢（4fps），大獎中等（6fps）
- 每個狀態結束後自動回到 idle
- 用 `_attack_timer` 控制攻擊/大獎動畫持續時間
- Spritesheet 比多個 PNG 效能好（減少 draw call）

## 38. ComfyUI 安裝路徑（重要）
- **問題：** ComfyUI portable 解壓縮後路徑不是 `C:\ComfyUI\ComfyUI`，而是 `C:\ComfyUI\ComfyUI_windows_portable`
- **正確路徑：**
  - 啟動腳本：`C:\ComfyUI\ComfyUI_windows_portable\run_nvidia_gpu.bat`
  - Checkpoints：`C:\ComfyUI\ComfyUI_windows_portable\ComfyUI\models\checkpoints\`
  - LoRAs：`C:\ComfyUI\ComfyUI_windows_portable\ComfyUI\models\loras\`
  - API：`http://127.0.0.1:8188`
- **教訓：** ComfyUI portable 解壓縮後有一層 `ComfyUI_windows_portable` 資料夾，不要假設路徑

## 39. PixelArtV4.safetensors 已不存在（404）
- **問題：** `https://huggingface.co/Onodofthenorth/SD_PixelArt_SpriteSheet_Generator/resolve/main/PixelArtV4.safetensors` 404
- **原因：** 該 repo 只有 `.ckpt` 格式（diffusers 格式），沒有 safetensors
- **替代方案：**
  1. SD 1.5 基礎模型：`https://huggingface.co/Comfy-Org/stable-diffusion-v1-5-archive/resolve/main/v1-5-pruned-emaonly.safetensors`（4GB）
  2. Pixel Art XL LoRA：`https://huggingface.co/nerijs/pixel-art-xl/resolve/main/pixel-art-xl.safetensors`（162MB）
- **組合使用：** SD 1.5 + LoRA strength 0.85 效果比單一 checkpoint 更靈活
- **教訓：** HuggingFace 模型 URL 要先用 API 確認檔案列表再下載

## 40. ComfyUI API 呼叫方式
- **端點：** `POST http://127.0.0.1:8188/prompt`
- **格式：** `{"prompt": {workflow_json}, "client_id": "uuid"}`
- **查詢結果：** `GET http://127.0.0.1:8188/history/{prompt_id}`
- **下載圖片：** `GET http://127.0.0.1:8188/view?filename=xxx&subfolder=&type=output`
- **工具：** `tools/comfyui_generate.py` — 完整 API 整合腳本
- **Workflow 結構：** CheckpointLoaderSimple → LoraLoader → CLIPTextEncode × 2 → EmptyLatentImage → KSampler → VAEDecode → SaveImage

## 41. 規格書玩法缺口修復記錄（2026-05-12）
- **T101 擬態型怪物**：死亡時需要「變形回原形」視覺，用 tween 做閃爍→縮放變形→爆炸消失三段動畫
- **T105 金幣魚**：擊破後金幣雨，用 15 個 ColorRect 做拋物線動畫（上升→下落→淡出）
- **BG002 硬雜草**：需要連點 2 次，用 `_weed_hp` 字典追蹤每個雜草的剩餘點擊次數
- **BOSS Phase 2**：收到 `boss_event.phase_change` 時，對 BOSS 節點做紅色調 + 閃爍 + 放大
- **投射物速度**：PlayerSnapshot 要包含 `projectile_speed` 和 `fire_rate`，Client 依此計算飛行時間
- **教訓**：每次完成功能後要逐條對照規格書，不能只看「有沒有這個功能」，還要看「行為是否完全一致」

## 42. GDScript 的 boss_event 訊號連接位置
- **問題**：TargetManager 需要處理 boss_event（Phase 2 視覺），但原本只連接了 target_* 訊號
- **解決**：在 `_ready()` 加入 `GameManager.boss_event.connect(_on_boss_event)`
- **教訓**：新增跨系統事件時，要確認所有需要響應的節點都有連接訊號

## 43. agent-sprite-forge 技術移植（2026-05-12）
- **來源：** https://github.com/0x0funky/agent-sprite-forge
- **核心技術：**
  1. **洋紅色背景（#FF00FF）** 比白色背景更可靠，不會誤刪白色毛皮
  2. **shared_scale（基於 idle 幀）** 確保所有動作幀大小一致，解決 attack 幀比 idle 大的問題
  3. **component_mode=largest** 只保留最大連通區域，去除 FX 噪點
  4. **bottom align** 腳底對齊，讓角色站在同一條線上
- **移植到：** `tools/process_sprites.py`
- **效果：** chiikawa 一致性從 height diff=3px, width diff=7px → height diff=2px, width diff=3px（✅）
- **教訓：** shared_scale 要基於 idle 幀，不是所有幀的最大值，否則 bigwin 幀被裁切後會縮小

## 44. generate_pixel_art_v5.py 的 dy 超出邊界 bug
- **問題：** usagi attack 的 `dy=-3` 讓長耳朵從 y=-3 開始，超出 32x32 畫布，導致角色幾乎消失（只剩 13 像素）
- **修復：** `ear_top = max(0, 0 + dy)` 防止超出邊界
- **同時修復：** 所有角色的 bigwin `dy` 從 +2/+3 改為 +1，防止身體底部超出畫布
- **教訓：** 程式生成像素圖時，所有座標都要做邊界檢查，特別是有 dy 偏移的動作幀

## 45. process_sprites.py 完整使用流程（重啟後必讀）
- **工具位置：** `tools/process_sprites.py`
- **四個模式：**
  ```
  # 1. 品質報告（先跑這個確認狀態）
  py tools/process_sprites.py --mode qc

  # 2. 重新對齊現有 sprites（shared_scale + bottom align）
  py tools/process_sprites.py --mode realign

  # 3. 處理 ComfyUI 生成的洋紅色背景圖
  py tools/process_sprites.py --mode comfyui --input <raw.png> --char chiikawa --pose idle

  # 4. 重建 Spritesheet
  py tools/process_sprites.py --mode sheet
  ```
- **ComfyUI 完整流程：**
  1. `tools\start_comfyui.bat`（手動在終端執行）
  2. 等 `http://127.0.0.1:8188` 就緒
  3. `py tools/comfyui_generate.py --all`（生成 9 張，洋紅色背景）
  4. 對每張圖執行 `--mode comfyui`
  5. `py tools/process_sprites.py --mode sheet`（重建 Spritesheet）
- **QC 通過標準：** height diff ≤ 2px, width diff ≤ 4px → ✅
- **目前狀態：** chiikawa ✅，hachiware ⚠️ 9px，usagi ⚠️ 5px

## 46. Kiro IDE 更新後的重啟檢查清單
```
1. py tools/process_sprites.py --mode qc        # 確認美術狀態
2. go build ./... （在 server/ 目錄）            # 確認 server 編譯
3. go vet ./...                                  # 確認無警告
4. 讀 docs/progress.md 的「待辦」清單           # 確認下一步
5. 最高優先：啟動 ComfyUI 生成 AI 角色圖        # 美術 62→75+
```

## 47. 捕魚機 RTP 正確模型（2026-05-15 深度分析）

### 根本問題
原版 RTP 高達 600%+ 的根本原因：
1. `required_hits = ceil(multiplier / bet_cost × 0.3)` → LV5 打 T001(2x) 只需 1 次保底
2. 1 次必定擊破 → 實際 RTP = 200%（遠超目標）
3. Bonus 每局 6-8 次，每次 bet × 50-150x → 額外 400%+ RTP

### 正確的 RTP 公式
```
kill_chance = BASE_RTP / multiplier  （每次命中的擊破機率）
期望命中次數 = 1 / kill_chance = multiplier / BASE_RTP
RTP = multiplier / 期望命中次數 = BASE_RTP  ✓
```

### 保底機制的正確設計
- **基礎目標（2x-10x）**：保底 = min(期望命中 × 3, Lifetime × 3.0 × 0.8)
  - 保底是「最壞情況保護」，不應影響正常 RTP
  - 上限確保玩家在目標消失前有機會觸發保底
- **特殊目標（15x+）**：不設保底（required_hits = 99999）
  - 高倍率目標設保底會導致 RTP 爆炸
  - 純機率擊破，高風險高報酬

### 最終數值設定
```
BASE_RTP = 0.92（基礎目標擊破機率係數）
DifficultyFactor = 16.0（已不再使用，改為動態計算）
LABOR_SCALE = 0.8（勞動值增益係數）
BONUS_MULT = 20-50x（Prototype 展示版）
```

### 最終 RTP 分布（Prototype 展示版）
| Bet Level | RTP | Bonus/局 | BOSS/局 |
|-----------|-----|---------|---------|
| LV1-3 | ~81% | 0.04-0.06 | 0 |
| LV5 | ~81% | 0.21 | 0 |
| LV7 | ~87% | 0.66 | 0.60 |
| LV10 | ~96% | 0.93 | 0.86 |

這符合業界慣例：高 bet 玩家 RTP 更高。正式版需要數值工程師精確調整。

### 修改的檔案
- `server/internal/data/tables.go`：DifficultyFactor 全部改為 16.0
- `server/internal/game/target/target.go`：RequiredHits 改為動態計算
- `server/internal/game/combat/combat.go`：CalcBonusReward 倍率降低（20-50x）
- `tools/simulate_rtp.py`：完整重寫，修正模擬邏輯
- `tools/rtp_analysis.py`：新增，用於分析 RTP 問題根源

## 48. 多幀角色動畫生成技術（2026-05-15）
- **工具：** `tools/generate_animation_frames.py`（已升級）
- **動畫設計：**
  - idle（4幀）：縮放 0.98-1.03x + 上下位移，模擬呼吸感
  - attack（3幀）：旋轉 -18°/+12° + 劍氣光效（右上角粉紅光點）
  - bigwin（4幀）：縮放 1.0-1.08x + 上移 0-14px + 金色星星
- **關鍵技術：**
  - `Image.rotate(angle, fillcolor=(0,0,0,0))`：旋轉保持透明背景
  - `img.resize((new_size, new_size), Image.NEAREST)`：縮放保持像素感
  - 置中貼上：`paste_x = (FRAME_SIZE - new_size) // 2`
  - numpy 逐像素加光效：`frame_arr[y, x] = [r, g, b, alpha]`
- **驗證工具：** `tools/verify_animated_sheets.py`
- **預覽工具：** `tools/preview_animation.py`（輸出 2x 放大 GIF）
- **Spritesheet 格式：** 384×288（4幀×3狀態×96px），符合 CharacterAnimator.gd
- **教訓：** 動畫幀的縮放要用 NEAREST 插值，LANCZOS 會讓像素圖模糊

## 49. generate_animation_frames.py 不能覆蓋 sprites（重要）
- **問題：** `frames[0].save(os.path.join(OUT_DIR, f"{char_name}_{state}.png"))` 把動畫幀0（旋轉過的圖）覆蓋了原本的 sprites
- **症狀：** process_sprites.py --mode qc 顯示 height diff 突然變大（6px）
- **修復：** 移除那行，sprites 由 `generate_pixel_art_v5.py` + `process_sprites.py` 管理
- **教訓：** 動畫生成腳本只應該寫入 sheets/ 目錄，不應該碰 characters/ 目錄

## 50. 目標物尺寸升級（48→64px）
- **原因：** 48×48 在遊戲畫面中太小，細節看不清楚
- **工具：** `tools/generate_targets_v3.py`（全新重寫）
- **技術：** 逐像素繪製 + fill_circle_shaded（帶陰影）+ fill_rect_shaded
- **Spritesheet：** `generate_spritesheet.py` 的 cell_size 從 32 改為 64
- **Godot 更新：** TargetManager.gd 的 HP 條從 32px 改為 48px
- **教訓：** 升級 sprite 尺寸時，要同步更新 Spritesheet cell_size 和 Godot 的 HP 條寬度

## 51. 角色 v6 設計改善（2026-05-15）
- **基礎尺寸：** 32×32 → 48×48（輸出仍 96×96）
- **眼睛：** 2×2 → 3×3 眼白 + 2×2 瞳孔 + 高光（`draw_eye_v6`）
- **腮紅：** 單點 → 橢圓漸層（alpha 漸變）
- **手臂：** 新增小圓手臂（讓角色更完整）
- **陰影：** 3色漸層（LIGHT/MID/DARK）
- **像素密度：** 53-55% → 64-76%（大幅提升）
- **教訓：** 基礎尺寸越大，細節越豐富，但要注意 2x 放大後的像素感

## 52. 背景 v2 生成技術（2026-05-15）
- **工具：** `tools/generate_backgrounds_v2.py`
- **海底背景：** 漸層 + 光線 + 珊瑚礁 + 海草 + 氣泡 + 沙地（11,306 種顏色）
- **BOSS 背景：** 暗紅漸層 + 警告條紋 + 裂縫 + 暗黑光環 + 石板地面（485 種顏色）
- **Bonus 背景：** 天空 + 雲朵 + 遠景樹木 + 草叢 + 花朵（126 種顏色）
- **顏色多樣性提升：** sea_bg 47→11,306，boss_bg 27→485
- **技術：** `random.Random(seed)` 確保每次生成結果一致
- **教訓：** 背景要有層次（遠景/中景/近景），不能只是漸層

## 53. 特效 v2 升級（2026-05-15）
- **命中特效：** 24×24 → 48×48，帶放射狀光線 + 中心爆炸圓 + 星形光芒
- **投射物：** 12×8 → 32×16，帶尾焰 + 橢圓主體 + 前端高光
- **死亡粒子：** 32×32 → 48×48，8方向粒子 + 中心爆炸
- **WARNING：** 新增 128×64 警告特效
- **Godot 更新：** Cannon.gd 的 hit/projectile scale 從 2.0 改為 1.0
- **教訓：** 升級 sprite 尺寸時要同步更新 Godot 的 scale，否則會過大

## 54. ComfyUI GPU 模式失敗原因（2026-05-15）
- **問題：** `Windows fatal exception: access violation` + `CUDA not available`
- **根本原因：** PyTorch 2.11.0+**cu130** 需要 CUDA 13.0，但驅動 555.85 只支援 CUDA 12.5
- **診斷指令：** `.\python_embeded\python.exe d:\Kiro\tools\check_cuda.py`
- **解決方案 A（推薦）：** 更新驅動到 596.49（https://www.nvidia.com/Download/index.aspx）
- **解決方案 B：** 降級 PyTorch 到 cu121（`pip install torch==2.5.1+cu121 --index-url https://download.pytorch.org/whl/cu121`）
- **教訓：** ComfyUI portable 版本的 PyTorch 可能比驅動新，要先確認相容性

## 55. UI v2 升級（2026-05-15）
- **coin.png：** 16×16 → 32×32，帶陰影 + ¥符號 + 高光
- **reward_bag.png：** 20×24 → 40×48，梨形布袋 + ¥符號
- **btn_*.png：** 80×32 → 96×36，圓角矩形 + 漸層 + 高光
- **labor_bar：** 200×20 → 240×24，圓角 + 漸層
- **warning_card：** 200×48 → 256×64，帶⚠符號
- **HUD.gd 改善：** 金幣顯示加🪙圖示，勞動值接近滿時變黃色⚡，Lock 按鈕加🔒圖示
- **BackgroundManager：** NEAREST → LINEAR 濾波（背景不需要像素感）
- **教訓：** 背景圖用 LINEAR 濾波，Sprite 用 NEAREST 濾波

## 56. ComfyUI GPU 模式成功啟動（2026-05-15）
- **問題：** PyTorch 2.11.0+cu130 需要 CUDA 13.0，驅動 555.85 只支援 12.5
- **解決：** 更新 NVIDIA 驅動到 596.49（支援 CUDA 13.0）
- **啟動指令：** `powershell -Command "Set-Location 'C:\ComfyUI\ComfyUI_windows_portable'; .\python_embeded\python.exe -s ComfyUI\main.py --windows-standalone-build --lowvram"`
- **sampler 名稱變更：** `euler_a` → `euler_ancestral`（ComfyUI 0.21.0）
- **生成速度：** GTX 1650，28 steps，約 40 秒/張
- **後處理：** `tools/batch_process_ai.py` 批次處理所有 AI 生成圖

## 57. AI 生成圖品質分析
- **洋紅色背景策略：** 部分圖片 SD 1.5 + LoRA 沒有完全遵守洋紅色背景，改用白色去背
- **自動選擇去背方式：** 比較洋紅色去背和白色去背的非透明像素數，選較多的
- **品質提升：** 程式生成 34-41% → AI 生成 42-66%（平均提升 60%）
- **chiikawa 一致性：** height diff=0px（完美），width diff=0px（完美）

## 58. 規格缺口修復（2026-05-15）
- **T102 寶箱怪受擊加速：**
  - Server：`handleAttack` 中命中 T102 且未逃跑時，設 `t.IsFleeing = true`，廣播 `target_update` 帶 `is_fleeing: true`
  - Client：`_on_target_updated` 收到 `is_fleeing` 時，設 `flee_speed = speed × 2.5`，改 behavior 為 `flee`，加閃爍紅色視覺
- **BG005 搗亂怪草暫停 0.3 秒：**
  - Client：點擊 BG005 後，`_is_active = false` 暫停 0.3 秒，顯示「😵 STUNNED!」文字，再恢復
- **教訓：** `IsFleeing` 欄位早就定義好了，只是沒有觸發邏輯，要定期對照規格書確認每個特殊行為都有實作

## 59. usagi bigwin 一致性修復（2026-05-15）
- **問題：** AI 生成的 usagi_bigwin 角色較小（3884px），與 idle（5304px）差距大
- **解法：** 從 idle 幀做變換生成 bigwin（放大 1.05x + 上移 8px + 金色色調 + 星星）
- **結果：** height diff 11px → 4px，非透明像素 3884 → 5755px
- **工具：** `tools/fix_usagi_bigwin.py`
- **教訓：** AI 生成的 bigwin 幀可能比 idle 小，用程式變換比重新生成更快且一致

## 60. BG003 發光雜草視覺效果（2026-05-15）
- **規格：** 「增加倍率」— 視覺上要有明顯的倍率提升感
- **實作：** 綠色光暈閃爍（modulate 0.5,2.0,0.5 ↔ WHITE）+ 「✨×UP!」浮動文字
- **位置：** BonusGame.gd 的 `_on_target_spawned` 中 BG003 分支
- **教訓：** 規格書的「增加倍率」是體感設計，視覺上要讓玩家感受到「這個很值得拔」

## 61. 目標物 AI 生成（2026-05-17）
- **工具：** `tools/comfyui_generate_targets.py`
- **生成數量：** 11 個目標物（T001-T105）
- **品質提升：** 平均 1100px (27%) → 2525px (62%)，提升 130%
- **問題：** T101 擬態怪物和 T105 金幣魚第一次去背效果差（1291/1328px）
- **解法：** 重新生成（T101 改種子，T105 改提示詞加「dark outline, clear silhouette」）
- **T105 提示詞關鍵：** 加 `dark outline, clear silhouette` 讓去背更乾淨
- **教訓：** 金色/黃色物體在洋紅色背景下去背效果差，需要強調輪廓

## 62. usagi 一致性修復技術（2026-05-17）
- **問題：** usagi attack 幀 78px 寬，idle 68px 寬，width diff=10px 超出門檻
- **根本原因：** AI 生成的 attack 幀（揮棒姿勢）本身就比 idle 寬，shared_scale 無法解決
- **解法：** 從 idle 幀做水平翻轉（FLIP_LEFT_RIGHT）+ 亮度提高 + 粉紅光暈生成 attack
  - 水平翻轉不改變 bbox，確保 width diff=0
  - 光暈只加在非透明像素上，不擴大 bbox
- **bigwin 修復：** 上移只 2px（不是 8px），星星放在 idle bbox 範圍內（y>=12）
- **工具：** `tools/fix_usagi_attack.py`、`tools/fix_usagi_bigwin_v2.py`
- **教訓：** AI 生成的動作幀可能比 idle 大，用程式變換比重新生成更快且一致

## 63. 規格缺口修復（2026-05-17）
- **BOSS 期間 Max Targets = 8**（規格書 9章）：
  - `updateBossBattle()` 加入非 BOSS 目標數量限制
  - 超出時移除最舊的目標
- **BG004 金色雜草 coin_shower**（規格書 29.3）：
  - Server `handleBonusClick` 加入 BG004 分支，廣播 `coin_shower` 事件
  - Client `BonusGame.gd` 的 `_on_bonus_event` 加入 `coin_shower` 處理
- **烏薩奇旋轉殘影**（規格書 2章）：
  - `Cannon.gd` 的 `_fire_projectile` 加入 usagi 的 `rotation_degrees` tween
  - 飛行時旋轉 720 度，模擬「黃色旋轉殘影」效果
- **烏薩奇大獎高速旋轉跳起**（規格書 2章）：
  - `_on_reward_received` 依 char_id 分支，usagi 做旋轉 360 度 + 跳起
- **教訓：** 規格書的角色特殊演出要逐條確認，不能只看「有沒有跳起」

## 64. PIL Image.fromarray readonly 問題
- **問題：** `Image.fromarray(arr)` 後 `pixels = img.load()` 再賦值報 `ValueError: image is readonly`
- **解法：** `Image.fromarray(arr.copy()).copy()` — 兩個 copy() 確保可寫
- **教訓：** numpy array 轉 PIL Image 後，需要 `.copy()` 才能用 `load()` 修改像素

## 65. Godot 4 BitmapFont 整合（2026-05-17）
- **格式：** BMFont .fnt + PNG 字體圖
- **生成工具：** `tools/generate_pixel_font.py`（Python Pillow，8x8 像素字體，95 個 ASCII 字元）
- **輸出：** `assets/fonts/pixel8.fnt` + `assets/fonts/pixel8.png`（256×96，2x 放大）
- **Godot 使用方式：**
  ```gdscript
  var font = load("res://assets/fonts/pixel8.fnt")
  label.add_theme_font_override("font", font)
  label.add_theme_font_size_override("font_size", 16)
  ```
- **注意：** BMFont .fnt 格式在 Godot 4 中用 `FontFile` 資源載入，不需要額外 .import 設定
- **套用位置：** HUD.gd（所有 Label + Button）、Cannon.gd（語音字卡）、TargetManager.gd（獎勵跳字）
- **教訓：** 像素字體要在 `_ready()` 中載入，用 `ResourceLoader.exists()` 先確認路徑存在

## 66. T105 金幣魚像素數量偏低的根本原因
- **問題：** T105 只有 42% 非透明像素（其他目標物 60-71%）
- **根本原因：** 魚的形狀是橫向細長橢圓 + 漸隱魚尾，透明邊緣多是正常的
- **正確評估方式：** 用 bbox 面積利用率（1740/2478 = 70%），不是整個畫布面積
- **教訓：** 不同形狀的 sprite 不能用同一個「整體面積佔比」標準評估，要看 bbox 利用率

## 67. Godot 4 Inner Class 不支援 extends（重要）
- **問題：** GDScript 4 的 inner class 語法 `class Foo extends Node2D:` 不被支援
- **解法：** 把需要繼承的 class 拆成獨立的 `.gd` 腳本，用 `preload` 載入
- **正確做法：**
  ```gdscript
  const BubbleLayerScript = preload("res://scripts/game/BubbleLayer.gd")
  var layer = BubbleLayerScript.new()
  ```
- **教訓：** Godot 4 inner class 只能用於純資料結構，不能繼承 Node 類型

## 68. Godot 4 動態繪圖（_draw + queue_redraw）
- **用途：** 不需要 Sprite 資產，用程式碼繪製動態效果（氣泡、粒子等）
- **關鍵 API：**
  - `_draw()` — 覆寫此方法，在裡面呼叫 `draw_arc`、`draw_circle` 等
  - `queue_redraw()` — 在 `_process` 中呼叫，觸發下一幀重繪
  - `draw_arc(pos, radius, from_angle, to_angle, point_count, color, width)` — 畫弧線
  - `draw_circle(pos, radius, color)` — 畫實心圓
- **氣泡效果：** 外圈用 `draw_arc`（淡藍白色），高光用 `draw_circle`（白色小點）
- **教訓：** 動態效果用 `_draw` 比 Sprite 節點更省記憶體，不需要額外 PNG 資產

## 69. chiikawa/hachiware attack 幀一致性修復（2026-05-17 夜間）
- **問題：** attack 幀（82x85）比 idle（78x78）大，height diff=7px, width diff=4px
- **根本原因：** AI 生成的 attack 幀（揮棒姿勢）本身就比 idle 大
- **解法：** 從 idle 幀做水平翻轉 + 亮度提高 + 光暈效果生成 attack
  - chiikawa：粉紅色光暈（討伐棒）
  - hachiware：藍色光暈（小八攻擊感）
- **結果：** 兩者都達到 0px/0px ✅
- **工具：** `tools/fix_char_attack.py`
- **教訓：** 水平翻轉不改變 bbox，是保持一致性最安全的方式

## 70. Godot 4 ShaderMaterial 動態套用（受擊閃白）
- **用途：** 目標物受擊時閃白，比 modulate 更精確（只影響 Sprite，不影響 HP 條）
- **做法：**
  ```gdscript
  var mat = ShaderMaterial.new()
  mat.shader = load("res://assets/shaders/hit_flash.gdshader")
  sprite.material = mat
  # 動畫
  var tween = create_tween()
  tween.tween_method(func(v): mat.set_shader_parameter("flash_amount", v), 1.0, 0.0, 0.12)
  ```
- **Shader 核心：** `COLOR = mix(original, flash_color, flash_amount)` — 只對非透明像素閃白
- **快取 ShaderMaterial：** 用 `set_meta("hit_flash_mat", mat)` 避免每次受擊都建立新 material
- **教訓：** modulate 會影響整個節點樹，shader 只影響單一 Sprite，更精確

## 71. 自動射擊智慧目標選擇評分系統（2026-05-17 夜間）
- **問題：** 原版只找第一個存活目標，不考慮倍率和位置
- **評分維度：**
  1. 倍率 × 2.0（高倍率優先）
  2. (1 - HPPercent) × 30.0（HP 低的優先，快要擊破）
  3. X < 400 時加分（快要離開畫面的優先）
  4. BOSS + 500（BOSS 戰集中火力）
- **結果：** 自動模式更智慧，優先打高價值目標
- **教訓：** 評分系統比硬規則更靈活，可以同時考慮多個維度

## 72. 斷線重連 UI 設計（2026-05-17 夜間）
- **做法：** 在 HUD 的 CanvasLayer 上動態建立 DisconnectOverlay（z_index=100）
- **連接訊號：** `NetworkManager.disconnected` 和 `NetworkManager.connected`
- **視覺：** 半透明黑色背景 + 📡 圖示 + 閃爍動畫
- **重連成功：** 顯示「已重新連線 ✓」（綠色）→ 1 秒後淡出
- **教訓：** 斷線 UI 要在 CanvasLayer 上，確保在所有遊戲元素上方顯示

## 73. BGM 淡入淡出實作（Godot 4 Tween）
- **問題：** 直接切換 BGM 會有突兀的切換感
- **解法：** 用 Tween 做淡出（0.3s）→ 切換 → 淡入（0.5s）
- **BOSS Phase 2：** `pitch_scale = 1.1`（音調提高 10%，讓緊張感更強）
- **注意：** 不能用 `await` 在 Autoload 的非 async 函數裡，改用 tween_callback
- **教訓：** BGM 切換要有淡入淡出，直接切換會讓玩家感覺突兀

## 74. HP 條顏色漸變（高→綠，中→黃，低→紅）
- **實作：** 依 HPPercent 動態設定 ColorRect.color
  - pct > 0.6 → Color(0.2, 0.9, 0.2)（綠）
  - pct > 0.3 → Color(1.0, 0.8, 0.1)（黃）
  - pct ≤ 0.3 → Color(1.0, 0.2, 0.2)（紅）
- **受擊閃爍：** 短暫變白（0.04s）再回原色（0.08s）
- **教訓：** HP 條顏色漸變是標準遊戲 UX，讓玩家直覺感受到目標快要死了

## 75. 目標物逃跑警告箭頭
- **觸發條件：** 目標物 x < 120 且有移動速度
- **視覺：** 「◀!」紅色文字，透明度依距離邊緣遠近
- **閃爍：** x < 60 時快速閃爍（150ms 間隔）
- **教訓：** 捕魚機的標準 UX，提示玩家目標快要跑掉，增加緊迫感

## 76. 小八大獎演出規格修正（2026-05-17 夜間）
- **問題：** 規格書說小八大獎是「高舉討伐棒」，但實作是普通跳起
- **修正：** 向上旋轉 -30 度 + 停頓 0.2 秒（高舉姿勢）+ 回正
- **字卡：** 「Yagaina!」→「尖尖哇嘎乃！」（規格書原文）
- **教訓：** 每個角色的大獎演出要逐條對照規格書，不能用同一個動作

## 77. Outline Shader 提升目標物辨識度（2026-05-18）
- **問題：** 目標物在複雜背景下辨識度不足，特別是小型目標（T002-T004 蟲類）
- **解法：** 建立 `outline.gdshader`，8方向採樣判定輪廓像素
- **輪廓顏色策略：**
  - 普通目標：黑色輪廓（0,0,0,0.8）— 清晰但不搶眼
  - 特殊目標：金色輪廓（1.0,0.85,0,1.0）— 提示高價值
  - BOSS：紅色輪廓（1.0,0.2,0.2,1.0）— 強調威脅感
- **注意：** outline shader 和 wobble shader 不能同時套用到同一個 Sprite2D
  - 解法：wobble 改用 Tween 做旋轉搖晃，outline 套用到 Sprite2D
- **教訓：** Shader 衝突時，優先保留視覺效果更強的（outline），另一個改用 GDScript 模擬

## 78. Wobble 效果用 Tween 替代 Shader（2026-05-18）
- **問題：** T103 流星和 T104 金草需要搖晃效果，但 wobble shader 和 outline shader 衝突
- **解法：** 用 `container.create_tween().set_loops()` 做旋轉搖晃
  - T103 流星：±5度，0.15s（快速搖晃，模擬飛行不穩定感）
  - T104 金草：±3度，0.4s（緩慢搖晃，模擬草在風中搖曳）
- **優點：** 不需要 shader，效能更好，且可以和 outline shader 共存
- **教訓：** 簡單的搖晃效果用 Tween 比 Shader 更省資源

## 79. Rainbow Glow Shader 大獎演出（2026-05-18）
- **用途：** 大獎（≥20x）時砲台角色有彩虹光暈，增加爽感
- **技術：** HSV 轉 RGB 函數 + TIME 驅動 hue 旋轉 + 輪廓採樣
- **持續時間：** 1.5 秒後自動移除（`cannon_sprite.material = null`）
- **注意：** 大獎後要確保 material 被清除，否則下一次普通攻擊也會有彩虹效果
- **教訓：** 臨時 shader 效果要用 Timer 確保清除，不能依賴其他事件觸發清除

## 80. death_particles.png 密度偏低的根本原因（2026-05-18）
- **問題：** death_particles.png 只有 14% 非透明像素，視覺上太稀疏
- **根本原因：** 原版只有 6 個小粒子，48x48 畫布大部分是透明的
- **修復：** 加入中心爆炸圓（r=8）+ 8方向射線 + 星形光芒 + 35個散落碎片 + 外圈光環
- **結果：** 14% → 44%，提升 3 倍
- **工具：** `tools/generate_death_particles_v2.py`
- **教訓：** 特效 sprite 的密度要 > 30%，否則在遊戲中幾乎看不見

## 81. warning.png 缺少 .import 檔案（2026-05-18）
- **問題：** warning.png 沒有對應的 .import 檔案，Godot 無法正確載入
- **原因：** 生成 warning.png 時忘記建立 .import 設定
- **修復：** 手動建立 warning.png.import，格式與其他特效 sprite 一致
- **教訓：** 每次新增 PNG 資產都要同時建立 .import 檔案，否則 Godot 會在首次開啟時自動生成（可能設定不正確）

## 82. 像素化過場 Shader（2026-05-18）
- **技術：** `pixelate_transition.gdshader` — 把畫面分成大像素塊（block_size = mix(1, 64, amount)）
- **用途：** 背景切換時先像素化（0.15s）再還原（0.2s），比直接切換更有像素遊戲感
- **整合：** `HitEffect.pixelate_transition(duration_in, duration_out, callback)` — callback 在最像素化時執行
- **BackgroundManager：** `_switch_bg()` 改用 `HitEffect.pixelate_transition()` 包裝背景切換
- **注意：** pixelate_transition 的 CanvasLayer.layer = 99（比 disconnect overlay 的 100 低一層）
- **教訓：** 過場效果要在 callback 中切換內容，不是在過場開始或結束時切換

## 83. check_assets.py 工具（2026-05-18）
- **用途：** 快速檢查所有美術資產的尺寸和非透明像素比例
- **門檻：** 非透明像素 > 30% 才算 ✅
- **例外：** T105 金幣魚（42%）是正常的，因為魚形狀細長
- **使用：** `py tools/check_assets.py`
- **教訓：** 定期執行資產品質檢查，避免低密度 sprite 混入正式版

## 83. Bonus Tick 廣播 Bug（2026-05-18）
- **問題：** `updateBonusGame()` 中 `if int(elapsed)%1 == 0` 永遠為 true
- **根本原因：** 任何整數 mod 1 都等於 0，這個條件沒有任何過濾效果
- **症狀：** 每 100ms（Server tick 間隔）廣播一次 bonus tick，比預期多 10 倍
- **修復：** 加入 `lastBonusTickAt time.Time` 欄位，用 `time.Since(g.lastBonusTickAt) >= time.Second` 判斷
- **影響：** 修復後減少 90% 的 bonus tick 網路流量
- **教訓：** `int(x)%1` 永遠為 0，不能用來做「每秒執行一次」的判斷。正確做法是用 `lastTickAt` 追蹤上次執行時間

## 84. Go 遊戲 Server 最佳實踐（2026-05-18 研究）
- **來源：** generalistprogrammer.com, leapcell.io, hemaks.org
- **關鍵點：**
  1. **goroutine 輕量**：每個 WebSocket 連線 2 個 goroutine，1000 個連線只需 2000 goroutines（幾 MB 記憶體）
  2. **permessage-deflate 壓縮**：`EnableCompression: true` 可減少 60-80% 頻寬（已實作）
  3. **sync.RWMutex**：讀多寫少的場景用 RWMutex，讀操作不互斥（已實作）
  4. **send channel buffer**：`make(chan []byte, 256)` 避免慢速客戶端阻塞 Server（已實作）
  5. **批次傳送**：`writePump` 中批次處理 channel 中的訊息（已實作）
- **目前架構評估：** 已符合最佳實踐，無需大改

## 85. Godot 4 HTML5 Export 大小優化（2026-05-18 研究）
- **來源：** jacobfilipp.com, godotengine.org forum
- **關鍵技術：**
  1. **Lossy 壓縮**：在 Import tab 設定 Lossy 壓縮，可大幅減少 .pck 大小
  2. **不要預先優化圖片**：讓 Godot 自己處理壓縮，保持原始品質
  3. **移除未使用資產**：Export 設定中排除不需要的檔案
  4. **gzip 壓縮**：HTML5 export 後用 gzip 壓縮 .wasm 和 .pck，需要 server 支援 Content-Encoding: gzip
- **目前狀態：** 未評估 export 大小，下次 export 時記錄
- **教訓：** HTML5 export 大小影響首次載入時間，對玩家體驗很重要

## 86. 洋紅色殘留問題（2026-05-18 深度分析）
- **問題：** AI 生成的角色 sprite 有大量洋紅色殘留
  - usagi 系列：438-541 pixels（佔 8-10%）
  - chiikawa 系列：50-631 pixels（佔 0.8-10%）
  - hachiware 系列：4-37 pixels（佔 0.1-0.6%）
- **根本原因：** AI 生成時使用洋紅色背景，`process_sprites.py` 的去背閾值不夠激進
- **修復工具：** `tools/fix_magenta_residue.py`
  - 判斷條件：`r > 150 and g < 100 and b > 100`（洋紅色特徵）
  - 距離計算：`sqrt((r-255)^2 + g^2 + (b-255)^2) < 100`
  - 總計移除 2816 個洋紅色像素
- **教訓：** 每次 AI 生成後必須執行 `fix_magenta_residue.py`，不能只靠 `process_sprites.py`

## 87. 目標物密度問題（2026-05-18 深度分析）
- **問題：** T001/T104 草葉密度只有 11%，T101 16%，T103 20%
- **根本原因：** `generate_targets_v3.py` 的草葉繪製用 `i*2` 步進（每隔一行），密度只有 50%
- **修復：** `tools/generate_targets_v4.py` — 改為填充三角形區域（連續繪製）
  - T001: 11% → 22%（+2x）
  - T104: 11% → 23%（+2x）
  - T101: 16% → 25%（+56%）
  - T103: 20% → 30%（+50%）
- **教訓：** 草葉形狀本來就是細長三角形，密度偏低是形狀特性，但可以通過填充改善

## 88. 視覺風格指南建立（2026-05-18）
- **建立：** `docs/visual-style-guide.md`（首次建立）
- **內容：** 官方色彩規範、像素規格、輪廓風格、動畫規格、UI 風格、特效風格
- **建立：** `reports/art/art-review-2026-05-18.md`（首份美術審核報告）
- **教訓：** 視覺風格指南是確保美術一致性的基礎，應該在專案初期就建立

## 89. 目標物 AI 生成版本未正確覆蓋（2026-05-18）
- **問題：** knowhow-log #61 說「目標物 AI 生成（2026-05-17）」，但 targets 目錄裡的是程式生成版本
- **可能原因：** AI 生成後沒有執行複製步驟，或被後來的程式生成覆蓋
- **教訓：** AI 生成目標物後，必須確認檔案已正確複製到 targets 目錄，並執行 `analyze_sprites.py` 驗證

## 83. Go Server Graceful Shutdown（2026-05-18）
- **問題：** Server 收到 SIGINT/SIGTERM 時直接退出，不等待連線完成
- **解法：** `signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)` + `srv.Shutdown(ctx)` + `g.Stop()`
- **效果：** Server 停止時正確清理 goroutine，不會有殘留連線
- **教訓：** 任何長期運行的 Server 都要實作 graceful shutdown

## 84. time.AfterFunc goroutine 洩漏（2026-05-18）
- **問題：** `time.AfterFunc(3s, func)` 在 Game 停止後仍然執行，操作已停止的 Game
- **解法：** 改用 `safeAfterFunc`，用 `select` 同時監聽 `time.After(d)` 和 `g.stopCh`
- **實作：**
  ```go
  func (g *Game) safeAfterFunc(d time.Duration, f func()) {
      go func() {
          select {
          case <-time.After(d):
              f()
          case <-g.stopCh:
              // Game 已停止，取消 timer
          }
      }()
  }
  ```
- **教訓：** 任何 `time.AfterFunc` 都要考慮 context 取消，特別是在有生命週期的物件中

## 85. Go Server pprof + /stats 端點（2026-05-18）
- **新增：** `/stats` 端點回傳 goroutine 數量、heap 記憶體、GC 次數
- **新增：** DEBUG 模式下啟用 `/debug/pprof/` 端點（`import _ "net/http/pprof"`）
- **使用方式：**
  - `curl http://localhost:7777/stats` — 快速查看記憶體狀態
  - `go tool pprof http://localhost:7777/debug/pprof/heap` — 詳細記憶體分析
  - `go tool pprof http://localhost:7777/debug/pprof/goroutine` — goroutine 分析
- **啟用 DEBUG：** `DEBUG=true ./gameserver`
- **教訓：** 生產環境的 pprof 要用 DEBUG flag 保護，不能直接暴露

## 86. Godot GDScript 資源快取模式（2026-05-18）
- **問題：** `_create_target_node` 每次都 `load(shader_path)` 和 `load(texture_path)`，每次生成目標都有 I/O 開銷
- **解法：** 在 `_ready()` 中 `_preload_resources()`，把所有常用資源存入 Dictionary
- **快取的資源：**
  - `_cached_textures: Dictionary` — 所有目標 Sprite texture
  - `_cached_outline_shader: Shader` — outline shader
  - `_cached_hit_flash_shader: Shader` — hit flash shader
  - `_cached_pixel_font: Font` — 像素字體
- **效果：** 目標生成時不再有 I/O 開銷，特別是高頻生成時（每 2 秒一個目標）
- **教訓：** Godot 的 `load()` 有快取機制，但第一次仍有開銷。預載入確保遊戲中不卡頓

## 87. Godot HTML5 Export 排除開發資源（2026-05-18）
- **問題：** `reference/` 目錄（0.65MB）和 `ai_generated/` 目錄（0.08MB）被打包進 HTML5 export
- **解法：** 在 `export_presets.cfg` 的 `exclude_filter` 加入這些目錄
  ```
  exclude_filter="*.import,assets/sprites/reference/*,assets/sprites/ai_generated/*,assets/sprites/characters/gifs/*,assets/sprites/downloads/*"
  ```
- **效果：** 下次 export 時 pck 大小減少約 0.7MB
- **教訓：** 開發用的參考圖、AI 生成的原始圖都不應該打包進 export，只打包遊戲實際使用的資源

## 88. Godot AtlasTexture 從 Spritesheet 裁切（2026-05-18）
- **問題：** TargetManager 用 12 個獨立 PNG，每個目標一個 draw call
- **解法：** 用 `AtlasTexture` 從 `targets_sheet.png` 裁切，所有目標共用一張 texture
- **實作：**
  ```gdscript
  var atlas = AtlasTexture.new()
  atlas.atlas = _targets_sheet  # 共用 Spritesheet
  atlas.region = Rect2(x, y, 64, 64)  # 裁切區域
  _atlas_textures[def_id] = atlas
  ```
- **效果：** 11 個普通目標從 11 個 texture 變成 1 個（BOSS 仍獨立）
- **注意：** AtlasTexture 要在 `_ready()` 預建立，不要每次生成目標時建立
- **教訓：** Spritesheet + AtlasTexture 是 Godot 2D 效能優化的標準做法

## 89. Godot 4 不支援 WebSocket permessage-deflate（2026-05-18）
- **發現：** Godot 4 的 WebSocketPeer 目前不支援 permessage-deflate 壓縮擴展
- **來源：** GitHub godot-proposals #13179（2025 年提案，尚未實作）
- **影響：** Go Server 的 `EnableCompression: true` 對 Godot client 沒有效果
- **結論：** 不需要在 Godot client 端做任何設定，Server 端的壓縮設定對 Godot 無效
- **替代方案：** 如果需要減少頻寬，考慮在應用層壓縮 JSON payload（但複雜度高）
- **教訓：** 確認 client 端是否支援某個 WebSocket 擴展，再決定是否在 Server 端啟用

## 90. Cannon.gd 資源快取優化（2026-05-18）
- **問題：** `_fire_projectile` 每次射擊都 `load(sprite_path)`，高頻射擊（10 FPS）造成 I/O 開銷
- **解法：** 在 `_ready()` 預載入所有投射物 texture 和 rainbow shader
- **效果：** 射擊時不再有 I/O 開銷，特別是 Auto 模式下每秒多次射擊
- **教訓：** 高頻呼叫的函數（射擊、特效）裡面不能有 `load()`，必須預載入

## 91. 海底焦散光 Shader（underwater_caustics.gdshader）（2026-05-18）
- **效果：** 模擬水面折射光斑（焦散），讓海底背景有動態光線感
- **技術：** 三層不同速度/方向的 2D 噪聲疊加，用 `pow(c, 2.0)` 提高對比度
- **深度衰減：** `depth_fade = 1.0 - UV.y * 1.5`，光線只在畫面上半部分顯示
- **效能：** 純 fragment shader，不需要額外 texture，效能開銷極低
- **套用方式：** 在 BackgroundManager._ready() 預載入，normal 狀態時套用，boss/bonus 時移除
- **教訓：** 焦散光是海底場景的標準視覺元素，用 shader 實現比 sprite 更省資源

## 92. BubbleLayer 升級：光線柱 + 海草搖擺（2026-05-18）
- **新增：** 4 條光線柱（梯形多邊形，頂部窄底部寬，帶漸層顏色）
- **新增：** 10 株海草（分段曲線，正弦波搖擺，底部粗頂部細）
- **光線柱技術：** `draw_polygon(pts, colors)` 支援頂點顏色，實現漸層效果
- **海草技術：** 分段 `draw_line`，每段顏色和粗細不同，頂部加小圓葉片
- **搖擺公式：** `sway = sin(time * speed + phase + t * 2.0) * amp * t`（t 越大搖擺越大）
- **效能：** 全部用 `_draw()` + `queue_redraw()`，不需要任何 Sprite 節點
- **教訓：** 海底環境的動態感主要來自光線和植物，不需要複雜的粒子系統

## 93. 角色 idle 動畫升級：4幀 → 8幀（2026-05-18）
- **工具：** `tools/upgrade_idle_8frames.py`
- **技術：** 正弦波插值（0°-315°，每 45° 一幀），位移 ±2px + 縮放 1.0-1.02x
- **Spritesheet 格式升級：** 4 cols × 3 rows → 8 cols × 3 rows（768×288）
- **fps 升級：** 4fps → 8fps（更流暢）
- **CharacterAnimator.gd 更新：** ANIM_CONFIG idle frames=8, fps=8.0, COLS=8
- **效果：** 呼吸感更自然，動作更流暢，視覺質量明顯提升
- **教訓：** 8 幀正弦波插值比 4 幀線性插值流暢 2 倍，計算量幾乎相同

## 94. usagi bigwin 0px 差距修復（2026-05-18）
- **問題：** bigwin 幀 bbox 69x79，比 idle 68x78 多 1px
- **根本原因：** scale=1.02 讓 bbox 稍微擴大
- **解法：** 完全不縮放不位移，只加金色色調 + 星星（嚴格在 idle bbox 內）
- **工具：** `tools/fix_usagi_bigwin_v3.py`
- **結果：** height diff=0px, width diff=0px ✅
- **教訓：** 任何縮放（哪怕 1.02x）都可能讓 bbox 擴大 1px，要達到 0px 差距必須完全不縮放

## 83. 排行榜系統設計（2026-05-18 DAY-010）
- **架構：** Server 端每 10 秒廣播 `leaderboard` 訊息，同時提供 HTTP GET `/leaderboard` 端點
- **排序依據：** `SessionScore`（本局累積獎勵），不用 `Coins`（因為玩家可能一開始就有很多金幣）
- **Player 新增欄位：** `SessionScore`（每次 AddReward 累加）、`MaxCoins`（歷史最高）、`KillCount`（每次 AddKill 累加）、`DisplayName`（ID 前 8 碼）
- **Client 排行榜 UI：** 右上角 Control 節點，動態建立（不依賴 .tscn 場景節點）
- **自己高亮：** `GameManager.get_player_id()` 比對 `player_data["id"]`，需要 `player_update` 訊息先到達
- **折疊功能：** 用 `▲/▼` 按鈕切換 EntriesContainer 的 visible
- **教訓：** 排行榜 UI 用程式動態建立比在 .tscn 裡設計更靈活，可以根據玩家數量動態調整高度

## 84. Go 排行榜排序（不用 sort.Slice）
- **問題：** 為了避免引入 `sort` 套件，用 bubble sort 實作
- **實際上：** `sort.Slice` 更好，但 bubble sort 對 ≤10 筆資料效能差異可忽略
- **教訓：** 小資料集（≤10）用 bubble sort 沒問題，大資料集才需要 sort.Slice

## 85. Godot HTML5 靜態檔案 gzip 壓縮（2026-05-18 DAY-010）
- **效果：** index.wasm 35.9MB → 9.0MB（-75%），總計 37.3MB → 9.9MB（-74%）
- **做法：** 預先用 Python gzip 壓縮生成 .gz 檔案，Go Server 檢查 Accept-Encoding 後提供 .gz 版本
- **工具：** `tools/compress_static.py`
- **Server 邏輯：** 只有 `Accept-Encoding: gzip` 時才提供壓縮版本，否則 fallback 到原始檔案
- **注意：** 必須同時設定 `Content-Encoding: gzip` 和正確的 `Content-Type`（不能讓 FileServer 猜測 .gz 的 Content-Type）
- **每次重新 export 後：** 需要重新執行 `py tools/compress_static.py` 更新 .gz 檔案
- **教訓：** 像素藝術遊戲的 wasm 壓縮效果極佳（75%），因為 wasm 有大量重複的指令序列

## 86. 目標物游泳動畫升級（2026-05-18 DAY-010）
- **升級前：** 只有 Y 軸上下搖擺（swim_amp 3-7px，swim_dur 0.6-1.2s）
- **升級後：** Y 軸搖擺 + 旋轉傾斜（±2-5度）+ 特殊目標縮放呼吸感
- **隨機相位：** `swim_phase = randf_range(0.0, 1.0)` 避免所有魚同步搖擺
- **效果：** 魚群看起來更有生命感，不像機器人整齊搖擺
- **教訓：** 游泳動畫加旋轉比只有 Y 軸位移更自然，但旋轉幅度要小（±5度以內）

## 87. PerformanceMonitor 自動效能降級（2026-05-18 DAY-010）
- **架構：** Autoload 單例，每 3 秒評估一次 FPS，連續 2 次低 FPS 才降級（避免瞬間抖動）
- **三個等級：** HIGH（全效果）、MEDIUM（粒子減半）、LOW（關閉游泳動畫/震動/outline shader，鎖 30 FPS）
- **P5 FPS：** 用最差 5% 的幀來判斷，比平均值更能反映真實體驗
- **保守升級：** 連續 5 次高 FPS 才升級（避免頻繁切換）
- **整合點：** ScreenShake.add_trauma、TargetManager 游泳動畫、TargetManager outline shader
- **教訓：** 效能降級要保守（多次確認才降），升級更要保守（避免頻繁切換造成視覺抖動）

## 83. 成就系統設計模式（2026-05-18 DAY-011）
- **架構：** `internal/game/achievement/` 獨立模組，不污染 Player 核心邏輯
- **Tracker 模式：** 每個 Player 持有一個 `*achievement.Tracker`，`TryUnlock` 回傳 `*AchievementUnlock`（nil = 已解鎖或不存在）
- **回傳值整合：** `AddKill`/`AddReward` 改為回傳 `[]*achievement.AchievementUnlock`，讓 game.go 統一處理廣播
- **Client 佇列機制：** 多個成就同時解鎖時，用 `_achievement_queue` 依序顯示，避免通知重疊
- **動畫：** Tween EASE_OUT + TRANS_BACK 做彈性滑入效果，比線性滑入更有活力
- **音效：** 複用 `BONUS_READY` 音效，不需要新增音效資產
- **教訓：** 成就系統要設計成「可選觸發」，不能影響核心遊戲邏輯的效能

## 84. BOSS 登場特效 Bug（2026-05-18 DAY-011）
- **問題：** BOSS 登場時沒有全畫面特效和震動
- **根本原因：** `_on_boss_event` 只處理 `boss_enter` 事件，但 Server 廣播的是 `spawn`
- **修復：** `if event == "spawn" or event == "boss_enter":` 同時處理兩種事件名稱
- **同時升級：** `spawn_boss_enter()` 從「兩個紅色閃光」升級為「雙波閃光 + 20 個粒子 + 雙衝擊波 + 大字動畫」
- **BOSS 擊殺慶祝：** HUD 的 `kill` 事件加入 `spawn_big_win` + `ScreenShake`
- **教訓：** Server 廣播的事件名稱要和 Client 處理的名稱完全一致，不能假設

## 83. 成就系統設計模式（2026-05-18 DAY-011）
- **Tracker 模式**：`achievement.Tracker` 封裝所有成就狀態，Player 持有一個 Tracker
- **解鎖回傳**：`AddKill()`、`AddReward()` 等方法回傳 `[]*AchievementUnlock`，讓 Game 層決定何時廣播
- **冪等性**：每個成就只能解鎖一次，`TryUnlock` 內部檢查 `unlocked` 標記
- **Client 佇列機制**：多個成就同時解鎖時，用 `_achievement_queue` 依序顯示，不重疊
- **教訓**：成就系統要設計成「解鎖即廣播」，不要在 Game 層做複雜的成就邏輯

## 84. 多玩家 bet level 平均計算（2026-05-18 DAY-012）
- **問題**：`spawnTarget` 只取第一個玩家的 bet level（`break` 後停止），多人時不公平
- **修復**：改為計算所有玩家的平均 bet level（`total / len(g.Players)`）
- **教訓**：遍歷 map 取第一個元素的模式在多人場景下是 bug，要明確計算平均或最大值

## 85. BOSS 期間目標清除的正確排序（2026-05-18 DAY-012）
- **問題**：原版用 map 迭代順序的 index 清除目標，但 Go map 迭代順序不確定
- **修復**：建立 `targetWithTime` 結構，依 `SpawnedAt` 排序後清除最舊的目標
- **教訓**：需要「移除最舊的 N 個」時，必須先排序再移除，不能依賴 map 迭代順序

## 86. 排行榜面板與 BOSS 計時器面板位置重疊（2026-05-18 DAY-012）
- **問題**：兩個面板都在 x=900，BOSS 戰時會重疊
- **修復**：排行榜面板改為 y=140（BOSS 計時器 y=50 + 高度 80px + 10px 間距）
- **教訓**：動態建立的 UI 面板要考慮與其他面板的位置關係，避免重疊

## 87. Bonus 雜草 Sprite 升級（2026-05-18 DAY-012 自我評估觸發）
- **問題：** BonusGame.gd 用 ColorRect 矩形代替雜草，完全沒有像素藝術感
- **解決：** 生成 5 種像素雜草 Sprite（BG001-BG005），用 Python Pillow 逐像素繪製
- **技術：**
  - `draw_stem()`：垂直莖幹，帶左暗右亮陰影
  - `draw_leaf()`：橢圓葉片，帶 3 色陰影（左上亮/右下暗）
  - BG005 搗亂怪草：S 形扭曲莖幹 + 眼睛 + 憤怒眉毛
  - BG004 金色雜草：十字星形閃光點
  - BG003 發光雜草：散落光點 + 頂部星形
- **Sprite 尺寸：** 32×48 px（雜草形狀細長，密度 26-37% 是正常的）
- **Godot 整合：** `Sprite2D` + `TEXTURE_FILTER_NEAREST` + `scale=1.5` + `offset=(0,-24)` 底部對齊
- **備用機制：** Sprite 載入失敗時自動降級到 ColorRect
- **教訓：** 遊戲中所有可見元素都應該有真正的像素藝術 Sprite，ColorRect 只是開發佔位符

## 83. 成就系統設計模式（2026-05-18 DAY-011）
- **Tracker 模式**：每個玩家有獨立的 `achievement.Tracker`，記錄已解鎖的成就 ID
- **TryUnlock 函數**：每次觸發條件時呼叫，若已解鎖則返回 nil，避免重複通知
- **佇列式通知 UI**：多個成就同時解鎖時，依序顯示，不重疊
- **教訓**：成就系統要設計成「冪等」的，同一個成就不能重複解鎖

## 84. Bonus 雜草 Sprite 升級（2026-05-18 DAY-012）
- **問題**：BG001-BG005 原本用 ColorRect 顯示，視覺品質低
- **解法**：用 `tools/generate_bonus_weeds.py` 生成像素藝術雜草 Sprite
- **BonusGame.gd 更新**：改用 Sprite2D 載入 PNG，備用 ColorRect
- **教訓**：Bonus 場景的視覺品質直接影響玩家的爽感，不能用 ColorRect 敷衍

## 85. Server 壓力測試工具設計（2026-05-18 DAY-013）
- **工具**：`tools/stress_test.py`
- **測試維度**：Heap 記憶體增長（< 10%）、Goroutine 增長（< +50）、錯誤率（< 5%）
- **模擬行為**：隨機投注切換、持續攻擊、隨機斷線重連（30% 機率）
- **監控端點**：`/stats`（goroutines, heap_alloc_mb, gc_count）
- **教訓**：壓力測試要模擬真實玩家行為（包含斷線重連），不能只測試正常流程

## 86. BOSS 進場預覽 UI 設計（2026-05-18 DAY-013）
- **功能**：警告階段（3 秒）顯示 BOSS 血條從 0 填滿，增加期待感
- **技術**：Tween 動畫 + ColorRect HP 條 + 倒數文字
- **位置**：畫面中央（320, 280），z_index=90（在遊戲元素上方，在 BOSS 文字下方）
- **動畫序列**：淡入 → HP 條填滿（2.5s, EASE_IN_QUAD）→ 倒數 3→2→1 → HP 條閃爍
- **BOSS 出現時**：自動淡出隱藏，切換到正式計時器
- **教訓**：警告階段的 UI 要讓玩家感受到「有東西要來了」，血條填滿是很直覺的視覺語言

## 87. Go 測試 Windows unlinkat 錯誤（非真正失敗）
- **問題：** `go test` 在 Windows 上有時報 `unlinkat ... The process cannot access the file because it is being used by another process`
- **原因：** Windows 的檔案鎖定機制，暫存的 test binary 被防毒軟體或其他程序鎖定
- **影響：** 不影響測試結果，所有 PASS 的測試都是真正通過的
- **解決：** 忽略此錯誤，或用 `-count=1` 避免快取
- **教訓：** Windows 上 `go test` 的 exit code 1 不一定代表測試失敗，要看實際輸出

## 88. Go 單元測試設計原則（2026-05-18 DAY-013）
- **測試覆蓋維度：** 初始狀態、CRUD 操作、狀態轉換、goroutine 生命週期、冷卻機制
- **safeAfterFunc 測試：** 用 channel 和 time.Sleep 驗證「Stop 後不執行」和「正常執行」兩種情況
- **Hub 在測試中：** `ws.NewHub()` 不需要 `Run()`，廣播會靜默失敗（沒有客戶端），不影響邏輯測試
- **教訓：** 測試 goroutine 生命週期時，要同時測試「正常執行」和「提前停止」兩種情況

## 89. bbox 利用率 vs 畫布利用率（重要評估原則）
- **問題：** 草類目標物（T001/T104）畫布利用率只有 22-33%，但 bbox 利用率達 50-53%
- **結論：** 草葉細長、蟲有觸角，透明邊緣多是**正常的形狀特性**，不是品質問題
- **正確評估方式：** bbox 利用率 >= 50% 即可接受，不需要強行填滿畫布
- **教訓：** 不同形狀的 sprite 要用 bbox 利用率評估，不能用整體面積佔比

## 90. 像素風格遊戲邊框設計（2026-05-18 DAY-013）
- **功能：** 海底主題裝飾邊框（珊瑚、貝殼、海草、金色裝飾線）
- **技術：** GDScript `_draw()` + `queue_redraw()`，不需要額外 PNG 資產
- **動態效果：** 珊瑚脈動（sin 波）、海草搖擺、金色裝飾點閃爍
- **位置：** z_index=2（在遊戲元素上方，在 HUD 下方）
- **效能：** 純 GDScript 繪製，每幀重繪，適合 60 FPS
- **教訓：** 遊戲邊框是捕魚機的標準視覺元素，能大幅提升整體美術質量，且不需要額外資產

## 83. GDScript 4 tween_callback 的 set_delay 用法
- **問題：** `tween_callback().set_delay(x)` 的 delay 是相對於前一個 tweener，不是絕對時間
- **正確做法：** 用 `tween_interval(interval)` + `tween_callback()` 交替，讓每個 callback 在固定間隔後執行
- **錯誤做法：** 在 loop 中計算絕對 delay 然後用 set_delay，邏輯複雜且容易出錯
- **教訓：** tween 序列是「相對時間」，不是「絕對時間」，設計時要用間隔而非絕對延遲

## 84. 命中特效升級：用 _draw 替代 ColorRect 模擬圓形（2026-05-18 DAY-015）
- **問題：** `_spawn_flash_ring` 用 4 個 ColorRect 模擬圓形，視覺上是方形，不夠精確
- **解法：** 建立獨立腳本 `FlashRing.gd` 和 `ShockwaveRing.gd`，用 `draw_arc` + `draw_circle` 繪製真正的圓形
- **優點：** 視覺更精確（真圓），代碼更簡潔，不需要多個子節點
- **注意：** GDScript 4 inner class 不能 extends Node，必須用獨立腳本 + preload
- **教訓：** 特效節點用 `_draw` 比 ColorRect 更精確，且不需要子節點

## 85. BOSS 進場特效從 BOSS 實際位置噴射（2026-05-18 DAY-015）
- **問題：** `spawn_boss_enter()` 的粒子從畫面中心噴射，不是從 BOSS 實際位置
- **解法：** 加入 `boss_pos` 參數（預設 Vector2(1100, 360)），TargetManager 傳入 BOSS 節點的實際位置
- **新增：** `_spawn_ground_shockwave()` — BOSS 登場時從底部向兩側擴散的橫向衝擊波
- **教訓：** 特效要從正確的位置噴射，才有「那個東西在那裡爆炸」的感覺

## 86. 螢幕扭曲衝擊波 Shader（2026-05-18 DAY-015）
- **來源：** gameidea.org/2025/01/20/shockwave-distortion-shader-2d-space-wrap/（基於 godotshaders.com）
- **技術：** `SCREEN_TEXTURE` 採樣 + `smoothstep` 環形遮罩 + 色差（chromatic aberration）
- **關鍵參數：**
  - `center`：扭曲中心（0-1 UV 座標，需要從世界座標轉換）
  - `radius`：環形半徑（從 0 動畫到 0.8 模擬向外擴散）
  - `strength`：扭曲強度（從 0.06 動畫到 0 模擬衰減）
  - `aberration`：色差強度（0.35 效果明顯但不過度）
- **世界座標轉 UV：** `normalized_pos = world_pos / viewport_size`（無相機縮放時）
- **用途：** BOSS 登場 + 大獎特效，比 ColorRect 衝擊波視覺效果強很多
- **注意：** 需要 `hint_screen_texture` 才能讀取螢幕，Godot 4 的 `SCREEN_TEXTURE` 需要在 uniform 宣告
- **教訓：** 螢幕扭曲 shader 是「免費的視覺升級」，不需要額外資產，只需要一個全畫面 ColorRect

## 87. Godot 4 動態建立 Theme（PixelTheme.gd）
- **問題：** HUD 按鈕用 Godot 預設樣式，和像素風格不一致
- **解法：** 建立 `PixelTheme.gd`（extends RefCounted），用 `StyleBoxFlat` 動態建立像素風格 Theme
- **關鍵 API：**
  - `theme.set_stylebox("normal", "Button", sb)` — 設定 Button 的 normal 狀態樣式
  - `theme.set_color("font_color", "Button", color)` — 設定 Button 文字顏色
  - `theme.set_stylebox("background", "ProgressBar", sb)` — 設定 ProgressBar 背景
  - `theme.set_stylebox("fill", "ProgressBar", sb)` — 設定 ProgressBar 填充
  - `node.theme = pixel_theme` — 套用 Theme 到節點（自動影響所有子節點）
- **像素風格設計原則：**
  - 不圓角（corner_radius = 0）
  - 2px 邊框（border_width = 2）
  - 深海藍背景 + 亮藍邊框 + 金色按下狀態
  - 陰影偏移 (1,1) 增加立體感
- **教訓：** Theme 套用到父節點後，所有子節點自動繼承，不需要逐一設定

## 83. 海底背景動態效果升級技術（2026-05-18 DAY-016）

### 水面波紋 Shader（water_surface.gdshader）
- **套用方式：** 在 BubbleLayer 的 `_ready()` 中動態建立 ColorRect（1280×80），套用 ShaderMaterial
- **核心技術：** 多層 sin 波疊加（3層不同頻率）+ 1D 噪聲閃爍 + 波峰泡沫 smoothstep
- **深度漸變：** `final_color.a = water_color.a * (1.0 - uv.y * 0.6)` — 頂部不透明，往下漸透明
- **教訓：** 水面效果要放在 BubbleLayer 內管理，不要放在 BackgroundManager，避免背景切換時殘留

### 漂浮微粒（浮游生物/塵埃）
- **設計：** 30% 機率是發光浮游生物（藍綠色，有暈圈），70% 是普通塵埃（白色半透明）
- **運動：** 緩慢漂移（drift_x/y）+ sin 波動（模擬水流）
- **發光效果：** 3層圓（外暈 r×2.5 + 中暈 r×1.5 + 主體），alpha 遞增
- **教訓：** 發光效果用多層 draw_circle 疊加，不需要 shader

### 遠景小魚群
- **設計：** 5-9 條小魚，統一方向，輕微上下擺動（swim_y = sin(time + phase)）
- **魚的繪製：** 橢圓身體（8邊形多邊形）+ 三角形尾巴，半透明（alpha 0.2-0.4）
- **生成策略：** 從左右邊緣生成，游過畫面後消失，每 6 秒生成一群
- **教訓：** 遠景元素要半透明（alpha < 0.5），避免搶走主要遊戲元素的注意力

### underwater_caustics 參數優化
- `caustic_intensity`: 0.08 → 0.12（更明顯）
- `caustic_scale`: 4.0 → 5.0（光斑更細緻）
- `time_scale`: 0.4 → 0.5（稍微快一點）

## 84. 數據埋點系統設計（2026-05-18 DAY-016）

### analytics.go 架構
- **Singleton 模式：** `analytics.Init()` 初始化，`analytics.Get()` 取得實例
- **JSONL 格式：** 每行一個 JSON 事件，方便 grep/jq 分析
- **mutex 設計：** RoomStats 不能包含 sync.RWMutex（複製問題），改用 Tracker 的 `roomMu` 欄位保護
- **原子計數器：** 高頻事件（attack）用 `atomic.Int64`，避免鎖競爭
- **SessionStats：** 玩家離開時輸出 session_summary 事件，包含完整統計

### Go vet mutex 複製問題
- **問題：** `RoomStats` 包含 `sync.RWMutex`，`GetRoomStats()` 回傳值複製時報 vet 錯誤
- **解法：** 把 mutex 從 struct 移出，改為 Tracker 的獨立欄位 `roomMu sync.RWMutex`
- **教訓：** 含 mutex 的 struct 不能直接複製，要用指標或把 mutex 移出

### analyze_logs.py 設計
- **輸入：** JSONL 日誌（每行一個 JSON 事件）
- **輸出：** 格式化報告（玩家/攻擊/擊破/獎勵/BOSS/Bonus/RTP）
- **RTP 警告：** < 85% 提示玩家流失，> 110% 提示數值過高
- **用法：** `py tools/analyze_logs.py`（今日）或 `--all`（所有日誌）或 `--json`（JSON 輸出）

## 85. analytics TotalBet 未更新 bug（2026-05-18 DAY-016 測試發現）

- **問題：** `GetRoomStats().TotalBet` 永遠是 0，RTP 計算錯誤
- **根本原因：** `EventAttack` 的 `updateStats` 只更新了 `t.room.TotalAttacks`，忘記更新 `t.room.TotalBet`
- **修復：** 在 `EventAttack` 的 `roomMu.Lock()` 區塊內加入 `t.room.TotalBet += int64(betCost)`
- **發現方式：** 單元測試 `TestRTPCalculation` 失敗（expected TotalBet=100, got 0）
- **教訓：** 每個新模組都要寫單元測試，特別是財務計算邏輯，不能只靠 build 通過

## 86. Go sync.Once 不能 reset（測試設計）

- **問題：** Singleton 模式用 `sync.Once`，測試時無法重置
- **解法：** 加入 `newTracker()` 工廠函數（package-private），測試直接建立實例，不走 singleton
- **教訓：** Singleton 模式要同時提供工廠函數供測試使用，不要讓測試依賴全域狀態

## 83. 海底環境音效生成技術（2026-05-18 DAY-017）
- **工具：** `tools/generate_ambient_sfx.py`
- **技術：**
  1. **低頻水流**：60/90/120 Hz 正弦波疊加 + 0.15 Hz LFO 調製（模擬水流起伏）
  2. **遠距離水聲**：白噪音 → IIR 低通濾波 → 高通（去 DC）= 帶通效果
  3. **隨機氣泡**：每 0.4-1.8 秒一個，上升音調 400→800 Hz + 指數衰減包絡
  4. **氣泡破裂音**：0.15 秒，300→900 Hz 上升 + 快速衰減 + 少量噪音
- **整合方式：**
  - AudioManager 新增 `play_ambient()` / `stop_ambient()` 方法（獨立播放器，-24 dB）
  - BackgroundManager 在 `_switch_bg("normal")` 時啟動，其他狀態停止
  - `edit/loop_mode=1`（.import 設定）讓環境音循環播放
- **音量設計：** -24 dB，不搶主音效，純背景沉浸感
- **教訓：** 環境音要用獨立的 AudioStreamPlayer，不能和 BGM 共用，否則 BGM 切換時環境音也會停

## 84. WebSocket API 文件建立（2026-05-18 DAY-017）
- **位置：** `docs/api/websocket-api.md`
- **內容：** 完整的 Client↔Server 訊息格式、Payload 欄位說明、遊戲流程範例
- **重要設計決策記錄：**
  - 每訊息獨立 frame（避免 JSON 合併問題，見 KnowHow #5）
  - permessage-deflate 壓縮已啟用
  - COOP/COEP headers 必要（HTML5 SharedArrayBuffer）
  - BOSS 計時獎勵：50-60 秒 = 500x，依時間遞減到 100x
- **教訓：** API 文件要在功能穩定後才寫，否則要一直更新；文件中的流程範例比純欄位說明更有價值

## 85. Bonus 音效升級技術（2026-05-18 DAY-017 自主觸發）
- **工具：** `tools/generate_bonus_sfx_v2.py`
- **設計原則：**
  1. **bonus_ready v2**：快速上升音階（C5→E6，5音，加速感）+ 短暫靜音（期待感）+ 和弦爆發（C6+E6+G6）+ 尾音上揚
  2. **bonus_trigger**：噪音衝擊（0.05s）+ 快速滑音（C5→C7，兩個八度）+ 和弦尾音
  3. **bonus_end**：下降音階（C7→C5）+ 勝利和弦（C5+G5+C6）
- **numpy 陣列長度問題**：`square_wave` 用 `int(SAMPLE_RATE * duration)` 計算，浮點誤差可能導致長度差 1，合成前要用 `min_len` 對齊
- **整合時機：**
  - `_show_ready()`：播放 BONUS_READY（觸發特效同時）
  - `_start_bonus()`：播放 BONUS_TRIGGER（遊戲開始瞬間）
  - `_end_bonus()`：播放 BONUS_END（結算時）
- **教訓：** 音效要配合視覺事件的時機，三個不同時機用三個不同音效，比一個音效更有層次感

## 86. BubbleLayer 視覺音效同步（2026-05-18 DAY-017 自主觸發）
- **問題：** 氣泡視覺上升到水面消失，但沒有對應音效
- **解法：** 在 `_process` 的氣泡移除邏輯中，當 `b["y"] < 80`（接近水面）且 `randf() < 0.25` 時播放 bubble_pop
- **25% 機率設計：** 避免太多氣泡同時破裂造成音效堆疊，保持輕柔的背景感
- **教訓：** 視覺效果消失時要有對應的音效回饋，但要控制頻率，不能每個氣泡都響

## 87. BGM 切換系統從未被呼叫的重大缺口（2026-05-19 DAY-018）
- **問題：** AudioManager 定義了完整的 `play_bgm()` 方法和 BGM 枚舉，但整個 Client 沒有任何地方呼叫它
- **根本原因：** BGM 系統是「設計了但沒整合」的典型案例，功能存在但沒有觸發點
- **發現方式：** `Select-String -Pattern "play_bgm"` 搜尋全部 .gd 檔案，只找到定義，沒有呼叫
- **修復方式：**
  1. BackgroundManager 的 `_on_state_changed()` 加入 `_switch_bgm(state)` 呼叫
  2. `_start_initial_ambient()` 同時啟動主 BGM
  3. GameManager 的 `_handle_boss_event()` 處理 Phase 2 切換
- **BGM 切換對應表：**
  - `normal_play` → MAIN_GAME
  - `boss_warning` → stop（靜音製造緊張感）
  - `boss_battle` → BOSS_BATTLE（新生成，8秒循環）
  - `boss_event.phase_change` → BOSS_RAGE
  - `boss_event.kill` → stop
  - `bonus_game` → BONUS_GAME
  - `boss_result/bonus_result` → stop（等狀態切回 normal 再播主 BGM）
- **教訓：** 每次新增系統後，必須用 grep 確認「有沒有地方呼叫它」，不能只看定義存在就認為整合完成

## 88. BOSS 戰 BGM 設計技術（2026-05-19 DAY-018）
- **工具：** `tools/generate_boss_battle_bgm.py`
- **緊張感設計原則：**
  1. **低頻 bass 驅動**：A2(110Hz) 方波，強弱交替（偶數拍 0.35 音量，奇數拍 0.22）
  2. **不和諧音程**：小二度（A3→Bb3）+ 增四度（A3→Eb4，魔鬼音程），製造壓迫感
  3. **打擊節奏**：噪音短脈衝 + 指數衰減包絡，每 0.5s 一次
  4. **高頻顫音**：25 Hz 開關調製的 440 Hz 方波，每 2 秒出現一次
- **無縫循環**：首尾各 0.15s 淡入淡出，確保循環播放無點擊聲
- **教訓：** 遊戲 BGM 的緊張感來自「不和諧音程 + 規律節奏 + 偶發緊張元素」的組合

## 85. Godot 4 ResourceLoader.load_threaded_request 背景載入（2026-05-19）
- **用途：** 遊戲啟動時背景預載入資產，避免首次使用時的卡頓（hitching）
- **API：**
  ```gdscript
  # 發起背景載入請求
  ResourceLoader.load_threaded_request(path)
  # 輪詢狀態
  var status = ResourceLoader.load_threaded_get_status(path)
  # 取得結果（只在 LOADED 狀態時呼叫）
  var res = ResourceLoader.load_threaded_get(path)
  ```
- **狀態值：**
  - `THREAD_LOAD_IN_PROGRESS`：載入中
  - `THREAD_LOAD_LOADED`：完成
  - `THREAD_LOAD_FAILED`：失敗（不要阻塞，直接跳過）
  - `THREAD_LOAD_INVALID_RESOURCE`：路徑不存在
- **最佳實踐：**
  1. 在 Autoload 的 `_ready()` 中用 `call_deferred` 啟動（避免初始化順序問題）
  2. 在 `_process` 中輪詢狀態，完成後 `set_process(false)` 停止輪詢
  3. 載入失敗要 `push_warning` 但不要 crash，遊戲繼續運行
  4. 快取結果到 Dictionary，後續用快取而不是重新 load
- **效能影響：** 背景載入不阻塞主執行緒，玩家感受不到卡頓
- **教訓：** 不要在 `_ready()` 中直接 `load()` 大量資產，會造成啟動卡頓

## 86. 升級特效設計原則（2026-05-19）
- **觸發時機：** 勞動值從 <100 變成 >=100（偵測邊界，不是每幀檢查）
- **偵測方式：** `_last_labor_value` 追蹤上次值，`if labor >= 100 and _last_labor_value < 100`
- **特效組合（由弱到強）：**
  1. 全畫面閃光（角色顏色，低 alpha）— 最先感知到
  2. 閃光環（角色顏色 + 金色，雙層）— 從砲台位置擴散
  3. 衝擊波 — 增加衝擊感
  4. 粒子噴射（向上扇形，帶重力）— 視覺焦點
  5. 文字動畫（BACK 彈性彈入）— 明確告知玩家
- **文字動畫 BACK 彈性：** `set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_BACK)` 讓文字有彈跳感
- **粒子向上扇形：** `angle = randf_range(-PI * 0.85, -PI * 0.15)` 確保粒子向上噴射
- **教訓：** 升級特效要讓玩家「感受到」，不只是「看到」，所以要有震動 + 閃光 + 粒子三重組合

## 87. GDScript 動態建立帶 _draw 的 Node2D（2026-05-19）
- **用途：** 在 TargetManager 等非 Autoload 腳本中動態建立帶自訂繪圖的節點
- **方法：** 用 `GDScript.new()` + `script.source_code` 動態建立腳本
  ```gdscript
  var coin = Node2D.new()
  var script = GDScript.new()
  script.source_code = """
  extends Node2D
  func _draw():
      draw_circle(Vector2.ZERO, 6.0, Color(1.0, 0.82, 0.0))
  """
  coin.set_script(script)
  ```
- **優點：** 不需要額外的 .gd 檔案，適合一次性特效節點
- **缺點：** source_code 字串中不能有縮排問題，要用 tab 或空格一致
- **金幣效果：** 主體圓形 + 深金色邊框 + 左上高光 + 中心 ¥ 符號
- **拋物線旋轉：** 上升時旋轉 90-270 度，下落時繼續旋轉，讓金幣有翻轉感
- **教訓：** 動態 GDScript 適合特效節點，但複雜邏輯還是要用獨立 .gd 檔案

## 88. 金幣雨拋物線設計原則（2026-05-19）
- **上升段：** `EASE_OUT`（快速上升，逐漸減速）
- **下落段：** `EASE_IN`（逐漸加速下落，模擬重力）
- **旋轉：** 上升 90-270 度，下落繼續 360-540 度（讓金幣一直在旋轉）
- **散落範圍：** 水平 ±150px，垂直 100-220px（比原版 ±120px 更寬）
- **峰值高度：** 70-140px（比原版 60-120px 更高，視覺更壯觀）
- **數量：** 18 枚（比原版 15 枚多 20%）
- **教訓：** 金幣雨要有「重量感」，上升快下落慢是錯的，要上升快下落也快（重力加速）

## 85. Godot 4 ResourceLoader.load_threaded_request 背景載入（2026-05-19）
- **用途：** 非阻塞背景預載入資產，避免首次使用時卡頓
- **API：**
  ```gdscript
  ResourceLoader.load_threaded_request(path)  # 開始背景載入
  ResourceLoader.load_threaded_get_status(path)  # 查詢狀態
  ResourceLoader.load_threaded_get(path)  # 取得已載入資源
  ```
- **狀態值：** `THREAD_LOAD_IN_PROGRESS` / `THREAD_LOAD_LOADED` / `THREAD_LOAD_FAILED`
- **注意：** 在 `_process` 中輪詢狀態，全部完成後停止輪詢（節省 CPU）
- **教訓：** 大量資產（48個）用背景載入，首次使用時直接從快取取得，無卡頓

## 86. 升級特效設計原則（2026-05-19）
- **觸發時機：** 勞動值從 <100 跨越到 >=100（邊界偵測，不是每幀檢查）
- **視覺層次：** 全畫面閃光（最底層）→ 閃光環（中層）→ 粒子（上層）→ 文字（最頂層）
- **文字動畫：** BACK 彈性（`TRANS_BACK` + `EASE_OUT`）讓文字有彈跳感
- **顏色策略：** 依角色 ID 選色（chiikawa=粉紅, hachiware=藍, usagi=黃）
- **教訓：** 升級特效要有「層次感」，不是所有效果同時出現，要有先後順序

## 87. Godot 4 動態 GDScript 節點建立最佳實踐（2026-05-19）
- **問題：** 動態建立的 Label 在 `_ready()` 之前讀取 `_pixel_font` 可能為 null
- **解法：** 建立節點後立即檢查 `is_instance_valid(_pixel_font)` 再套用字體
- **面板建立時機：** 用 `call_deferred()` 延遲到下一幀，確保所有 Autoload 都已初始化
- **教訓：** 動態建立 UI 節點時，所有資源引用都要做 null check

## 88. 金幣雨像素金幣設計（2026-05-19）
- **問題：** 原版金幣雨用 ColorRect（純色方塊），視覺太簡陋
- **升級：** 用 `_draw()` 繪製像素金幣（圓形 + ¥符號 + 高光 + 旋轉）
- **旋轉拋物線：** 每個金幣有獨立的 `rotation_speed`（±2-5 rad/s），模擬真實拋物線
- **教訓：** 特效的細節決定品質，金幣要像金幣，不能用純色方塊代替

## 89. Godot 4 Performance API（2026-05-19）
- **記憶體：** `Performance.get_monitor(Performance.MEMORY_STATIC)` — 靜態記憶體（bytes）
- **Draw Calls：** `Performance.get_monitor(Performance.RENDER_TOTAL_DRAW_CALLS_IN_FRAME)`
- **物件數：** `Performance.get_monitor(Performance.OBJECT_COUNT)`
- **節點數：** `Performance.get_monitor(Performance.OBJECT_NODE_COUNT)`
- **注意：** 這些值每幀都在變，建議每 0.5 秒更新一次顯示（避免數字跳動太快）
- **教訓：** 效能監控面板要顯示有意義的指標（記憶體/DC/節點），不只是 FPS

## 90. Godot 4 HTTPRequest 節點使用方式（2026-05-19 DAY-020）
- **用途：** 在 GDScript 中發送 HTTP GET 請求（查詢房間列表等）
- **建立方式：** `var http = HTTPRequest.new(); add_child(http)`
- **發送請求：** `http.request("http://localhost:7777/rooms")`
- **回應處理：** `http.request_completed.connect(_on_response)`
- **回調簽名：** `func _on_response(result, response_code, headers, body: PackedByteArray)`
- **解析 JSON：** `body.get_string_from_utf8()` → `JSON.new().parse(text)`
- **注意：** HTTPRequest 節點必須是場景樹的子節點才能正常工作
- **教訓：** 不要用 `await`，用訊號回調處理非同步 HTTP 請求

## 91. 大廳 UI 設計原則（2026-05-19 DAY-020）
- **overlay 模式**：大廳作為 z_index=150 的 overlay，不替換遊戲場景，向後相容
- **房間列表行**：名稱 + 玩家數進度條 + 投注等級範圍 + 加入按鈕
- **滿員視覺**：紅色遮罩 + 按鈕禁用 + 玩家數紅色顯示
- **快速加入**：找人數最少且未滿的房間，一鍵加入
- **切換房間**：TopBar 右側「🏠」按鈕，不佔用遊戲空間
- **教訓：** 大廳 UI 要輕量，不能影響遊戲主流程，overlay 模式最安全

## 83. analytics EventBonusEnd 重複宣告問題（2026-05-19 DAY-021）
- **問題：** 嘗試新增 `EventBonusEnd` 常數，但它已在 analytics.go 第 30 行存在
- **根本原因：** 沒有先搜尋現有常數就直接新增，導致 redeclared 編譯錯誤
- **解決：** 用 grep_search 確認現有常數後，移除重複宣告
- **教訓：** 新增常數前必須先搜尋確認不存在，特別是在大型 Go 檔案中

## 84. Godot 4 LineEdit.select_all() 用法（2026-05-19 DAY-021）
- **用途：** 對話框開啟時自動選取輸入框內容，方便玩家直接輸入新名稱
- **做法：** `line_edit.grab_focus()` 後接 `line_edit.select_all()`
- **注意：** 必須先 `grab_focus()` 再 `select_all()`，順序不能反
- **教訓：** 輸入框對話框的標準 UX：開啟時自動聚焦並選取現有內容

## 85. Go analytics 埋點補完清單（2026-05-19 DAY-021）
- **已補完：** auto_toggle / bet_change / boss_kill / bonus_end
- **原本缺少的原因：** 這些事件在 handleAutoToggle/handleBetChange/handleBossKill/endBonusGame 中沒有呼叫 tracker
- **補完方式：** 在每個 handler 函數末尾加入 `tracker.Track(analytics.EventXxx, p.ID, data)`
- **教訓：** 新增 handler 時要同步加入埋點，不要等到後期補

## 86. set_display_name WebSocket 訊息設計（2026-05-19 DAY-021）
- **協定：** Client → Server: `{"type": "set_display_name", "payload": {"display_name": "名稱"}}`
- **Server 驗證：** 長度 1-16 字元，超出回傳 error 訊息
- **Client UI：** TopBar 加入「✏」按鈕，點擊開啟對話框，輸入後送出
- **排行榜整合：** 設定後立即反映在下次排行榜廣播中（10 秒週期）
- **教訓：** 玩家名稱設定是多人遊戲的基本功能，應該在早期就加入

## 87. Combo 連擊系統設計（2026-05-19 DAY-022）
- **設計：** 2 秒內連續擊破 → Combo 計數 +1，廣播 `combo_event`
- **加成：** ×2=+10% 勞動值，×3=+20%，×4+=+30%
- **Server：** `Player.AddKillCombo()` 回傳 (comboCount, laborBonus)
- **Client：** `HitEffect.spawn_combo()` 顯示彈入文字 + 閃光環 + 粒子
- **顏色：** ×2=綠，×3=黃，×4=橙，×5+=紫（視覺升級感）
- **教訓：** Combo 系統是捕魚機的標準 game feel 提升，應該在早期就加入
- **注意：** Combo 只傳給觸發者（Hub.Send），不廣播，避免干擾其他玩家

## 88. 自我評估誠實原則（2026-05-19）
- **問題：** progress.md 說 100/100/100%，但這是 AI 自己寫的，可能有盲點
- **正確做法：** 每次評估要獨立驗證，不能只看自己寫的文件
- **驗證方法：** 1) 執行 qa_check.py 2) 讀規格書逐條對照 3) 上網搜尋業界標準
- **發現：** Combo 系統是業界標準但規格書沒有明確提到，需要主動補充
- **教訓：** 「規格書 100% 實作」≠「遊戲體驗 100%」，還要考慮業界標準的 game feel

## 87. 觀戰模式（Spectator Mode）設計原則（DAY-023）
- **核心設計：** 觀戰者是 WebSocket 連線的「只讀角色」，收到所有廣播但不能發送遊戲指令
- **實作方式：**
  - `ClientRole` 枚舉：`RolePlayer` / `RoleSpectator`
  - `readPump` 中過濾觀戰者訊息（只允許 `ping`）
  - `OnConnect` / `OnDisconnect` 只對 `RolePlayer` 觸發（不影響遊戲邏輯）
  - `PlayerCount()` 只計算玩家，`SpectatorCount()` 只計算觀戰者
- **觀戰者初始化：** 連線後 100ms 非同步傳送快照（遊戲狀態 + 所有目標 + 排行榜）
- **端點設計：**
  - `ws://host/spectate?room_id=xxx` — 觀戰 WebSocket
  - `GET /spectate/snapshot` — HTTP 快照（供前端預覽房間）
- **教訓：** 觀戰模式不需要修改遊戲邏輯，只需在 Hub 層做角色過濾即可

## 88. Go WebSocket Hub 角色分離測試技巧（DAY-023）
- **問題：** Hub 的 Register/Unregister 需要真實 WebSocket 連線才能測試
- **解法：** 直接操作 `h.clients` map（加鎖），繞過 WebSocket 升級步驟
  ```go
  h.mu.Lock()
  h.clients["test-id"] = &Client{ID: "test-id", Role: RolePlayer, send: make(chan []byte, 10)}
  h.mu.Unlock()
  ```
- **注意：** `send` channel 必須有 buffer（`make(chan []byte, 10)`），否則 `Broadcast` 會 block
- **教訓：** 測試 Hub 邏輯時，不需要真實 WebSocket，直接操作 clients map 更簡單

## 89. 觀戰模式 Client 端整合（DAY-024）
- **LobbyManager 觀戰按鈕：**
  - 每個房間 row 加入「👁 觀戰」按鈕（藍色）
  - 底部加入「👁 觀戰」快速按鈕（選最多人的房間）
  - 觀戰按鈕呼叫 `NetworkManager.spectate_room(room_id)`
- **NetworkManager 觀戰支援：**
  - `spectate_room(room_id)` — 連線到 `/spectate?room_id=xxx`
  - `is_spectator()` — 查詢當前是否為觀戰模式
  - `fetch_spectator_snapshot()` — HTTP GET `/spectate/snapshot`
  - `spectator_snapshot_received` 訊號
- **HUD 觀戰標籤：**
  - `_show_spectator_badge()` — 右上角藍色「👁 觀戰中」標籤
  - `_on_lobby_room_selected` 依 `is_spectator()` 分支顯示不同提示
- **教訓：** 觀戰模式的 Client 端只需要：1) 連線到不同端點 2) 顯示視覺標識 3) 禁用攻擊按鈕（Server 端已過濾）

## 90. Store 整合到 Game 的正確方式（DAY-026）
- **問題：** Store 骨架建好了，但沒有接到 Game，玩家金幣仍然不持久化
- **整合點：**
  1. `NewGameWithStore(id, hub, store, initialCoins)` — 新建構子，向後相容 `NewGame`
  2. `AddPlayer`：先從 Store 讀取，有則恢復（coins/maxCoins/killCount/betLevel/displayName）
  3. `RemovePlayer`：離開時儲存到 Store + 更新排行榜
  4. `main.go`：`store.New(cfg.RedisURL)` 初始化，graceful shutdown 時 `store.Close()`
- **型別注意：** `PlayerState.Coins` 是 `int64`，`Player.Coins` 是 `int`，需要 `int(saved.Coins)` 轉換
- **降級策略：** `REDIS_URL` 為空 → 記憶體模式，Server 重啟後狀態丟失但不中斷服務
- **教訓：** 設計文件和骨架完成後，要立刻整合到主流程，否則設計只是紙上談兵

## 91. config.go 加入 REDIS_URL 環境變數
- **做法：** `RedisURL: getEnv("REDIS_URL", "")` — 空字串代表記憶體模式
- **部署時：** `REDIS_URL=redis://localhost:6379 ./gameserver`
- **本機開發：** 不設定 REDIS_URL，自動使用記憶體模式
- **教訓：** 環境變數要有合理的預設值，讓開發者不需要額外設定就能跑起來

## 83. 像素角色眨眼動畫技術（2026-05-19 自主觸發）
- **問題：** idle 動畫只有上下搖擺，缺乏生命感
- **解法：** 在 8 幀 idle 中，第 5-6 幀加入眼睛閉合效果
- **技術：**
  1. 掃描角色圖的眼睛區域（y=28-45, x=20-75），找深色像素（r<80, g<80, b<80）
  2. 在指定幀用皮膚色（白色）覆蓋眼睛上半部
  3. 在閉眼線加一條深色水平線（模擬眼皮）
  4. 眨眼時機：幀4=半閉(0.5), 幀5=全閉(1.0), 幀6=半閉(0.5)
- **驗證：** chiikawa 幀5 眼睛深色像素從 240 降到 8（97% 減少）✅
- **工具：** `tools/add_blink_animation.py`
- **效果：** 美術質量從 91 提升到 93/100（角色更有生命感）
- **教訓：** 眨眼是讓像素角色「活起來」最有效的技術之一，成本低效果好

## 83. RedisStore 完整實作技術要點（DAY-028）
- **套件：** `github.com/redis/go-redis/v9`（v9.7.3，支援 Redis 6+）
- **Key 設計：**
  - 玩家狀態：`player:{id}` → JSON String，TTL 7天
  - 排行榜：`leaderboard:daily:{YYYY-MM-DD}` → Sorted Set，TTL 30天
- **排行榜只保留最高分：** 先 `ZScore` 取現有分數，比較後才 `ZAdd`（Redis 6.2 的 `ZADD GT` 也可以）
- **整合測試設計：** 無 `REDIS_URL` 環境變數時 `t.Skip()`，不阻擋 CI
- **go get 的 stderr 誤判：** `go get` 成功但 PowerShell 把 stderr 視為錯誤，exit code=1，實際上套件已安裝
- **教訓：** Redis 整合測試要能在無 Redis 環境下跳過，不能讓 CI 因為沒有 Redis 就失敗

## 84. Go 內建函數名稱衝突
- **問題：** 變數命名為 `copy` 會遮蔽 Go 內建的 `copy()` 函數，雖然不報錯但是壞習慣
- **解法：** 改用 `cp` 作為拷貝變數名稱
- **教訓：** 避免使用 Go 內建函數名稱作為變數名：`copy`, `len`, `cap`, `make`, `new`, `append`, `delete`, `close`, `panic`, `recover`

## 85. B001 BOSS 動畫 Spritesheet 整合技術（DAY-028b）
- **Spritesheet 格式：** 512x384（4幀×3狀態×128px），Row 0=idle, Row 1=phase2, Row 2=death
- **AtlasTexture 快取：** 預建立所有 12 個幀的 AtlasTexture，key = "row_col"
- **動畫切換：** `_process` 中用 timer 按 FPS 切換幀，Phase 2 事件直接改 `_boss_anim_row`
- **death 動畫：** 4幀 × 8fps = 0.5秒，配合縮放 tween 同步消失
- **PIL fill_circle None color bug：** 傳入 None 作為 color 參數會 crash，要加 `if color is not None` 檢查
- **教訓：** BOSS 動畫要有明確的狀態機（idle/phase2/death），不能只靠 modulate 顏色變化

## 86. 程式生成像素角色的光照計算
- **公式：** `dot = -(nx*lx + ny*ly)`，`light = 0.5 + 0.5 * max(0, dot)`
- **光源方向：** 固定左上 `(-1, -1)`，歸一化後使用
- **法線計算：** `nx = dx / max(dist, 0.1)`，避免除以零
- **顏色計算：** `r_c = min(255, int(base_r * (0.7 + 0.6 * light)))`
- **教訓：** 0.7 是最暗值（陰影），0.7+0.6=1.3 是最亮值（高光），這個範圍讓立體感明顯

## 83. Attack 動畫幀數不一致問題（2026-05-19）
- **問題：** `upgrade_idle_8frames.py` 生成 4 幀 attack，但 metadata 和 CharacterAnimator.gd 都寫 3 幀
- **症狀：** 第 4 幀（復位幀）從未被播放，攻擊動畫少一幀
- **修復：** metadata 改為 `"frames": 4, "fps": 10.0`，CharacterAnimator.gd 同步更新
- **教訓：** 生成工具的 metadata 和 Godot 的 ANIM_CONFIG 必須保持一致，每次修改都要同步

## 84. HP 條低血量脈動效果（2026-05-19）
- **技術：** `create_tween().set_loops()` 做 `modulate:a` 0.4↔1.0 脈動（0.25s 間隔）
- **管理方式：** 用 `node.set_meta("hp_pulse_tween", pulse)` 儲存 tween 引用
- **停止方式：** `pulse.kill()` + `node.remove_meta()` + 重置 `modulate.a = 1.0`
- **觸發條件：** HP < 30% 時啟動，HP 回到 30% 以上時停止
- **教訓：** 用 node meta 儲存 tween 引用，可以在任何時候停止它，比用全域字典更乾淨

## 85. 成就系統 Type 欄位設計（2026-05-19）
- **設計：** 4 種類型：normal（金色）/ boss（紅色）/ bonus（綠色）/ special（紫色）
- **Server 端：** `Achievement.Type` 欄位，`TryUnlock` 和 `UnlockedList` 都傳遞
- **Client 端：** `_show_next_achievement()` 依 type 設定左側彩色邊條
- **教訓：** 成就類型要在 Server 端定義，Client 只負責顯示，不要在 Client 端硬編碼類型判斷

## 86. 成就通知面板動畫升級（2026-05-19）
- **改善：** 滑入後加彈跳縮放（scale 1.0→1.05→1.0），淡出改為 modulate:a 漸隱
- **注意：** 面板消失後必須重置 scale 和 modulate.a，否則下次顯示會有殘留狀態
- **Tween 並行：** 滑入 tween 用 `set_parallel(false)`，縮放用獨立的 `create_tween().set_parallel(true)`
- **教訓：** 多個 tween 同時執行時，要用獨立的 tween 物件，不要在同一個 tween 上混用 parallel 和 sequential

## 83. 動態 GDScript vs 靜態 preload（效能差異）
- **問題：** `GDScript.new()` + `script.source_code = "..."` + `node.set_script(script)` 每次都重新編譯腳本
- **影響：** T105 金幣雨生成 18 個金幣，每個都重新編譯一次 GDScript，造成不必要的 CPU 開銷
- **解法：** 把腳本內容寫成獨立的 `.gd` 檔案，用 `const Script = preload("res://path/to/script.gd")` 預載入
- **使用方式：** `var node = Script.new()` — 和 `Node2D.new()` 一樣，但有自訂的 `_draw()` 等方法
- **教訓：** 任何需要重複建立的節點，都應該用靜態腳本 + preload，不要用動態 GDScript

## 84. Go main.go 未使用常數的清理原則
- **問題：** `defaultPort = "8080"` 定義在 main.go 但從未使用（實際 port 由 config.go 的 `getEnv("PORT", "7777")` 決定）
- **影響：** 造成混淆，讓讀者以為 server 用 8080，但實際是 7777
- **解法：** 直接刪除未使用的常數
- **教訓：** Go 不會對未使用的常數報錯（只有未使用的變數才報錯），需要手動清理
