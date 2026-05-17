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
