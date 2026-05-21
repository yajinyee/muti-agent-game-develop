# Animation Pipeline 完整規格

> 版本：1.0.0  
> 維護者：Animation Agent + Art Director  
> 最後更新：2026-05-17

---

## 概覽

本規格定義吉伊卡哇：像素大討伐的動畫生產流程，確保所有角色動畫達到 Animation Quality >= 88 的品質門檻。

---

## 8 步驟動畫 Pipeline

### Step 1：Reference Lock（參考鎖定）
- 確認角色的 canonical sprite（正式參考圖）
- 鎖定 canvas size（不得在後續步驟更改）
- 確認 anchor point（通常為底部中心）
- 輸出：`references/char-<name>-reference.png`

### Step 2：Pose Plan（姿勢規劃）
- 列出該動畫所需的所有關鍵姿勢
- 標記每個姿勢的幀號
- 確認動畫循環點（loop point）
- 輸出：`references/char-<name>-<state>-pose-plan.md`

### Step 3：Keyframe（關鍵幀製作）
- 依照 Pose Plan 製作每個關鍵幀
- 每幀必須通過 Frame Consistency 檢查
- 使用 ComfyUI 生成或手工繪製
- 輸出：`assets/sprites/chars/<name>/<state>/keyframe_*.png`

### Step 4：In-between（補間幀製作）
- 在關鍵幀之間插入補間幀
- 確保動作流暢，無跳幀感
- 補間幀同樣需通過 Frame Consistency 檢查
- 輸出：`assets/sprites/chars/<name>/<state>/frame_*.png`

### Step 5：Consistency Check（一致性檢查）
- 執行 `py tools/animation_pipeline.py --check <sheet_path>`
- 一致性分數必須 >= 85 才能進入下一步
- 不通過則回到 Step 3 修正
- 輸出：`reports/animation/consistency-<name>-<state>.json`

### Step 6：Spritesheet（精靈圖合成）
- 將所有幀合成為單一 spritesheet
- 格式：PNG，透明背景，水平排列
- 命名：`<name>_<state>.png`
- 輸出：`client/chiikawa-pixel/assets/sprites/chars/<name>/<name>_<state>.png`

### Step 7：Preview GIF（預覽 GIF 生成）
- 執行 `py tools/animation_pipeline.py --gif <name>`
- 生成預覽 GIF 供 Art Director 審核
- 輸出：`reports/animation/preview/<name>_<state>.gif`

### Step 8：Godot Import（Godot 匯入）
- 在 Godot 中設定 AnimatedSprite2D
- 設定正確的 FPS（idle: 8fps, attack: 12fps, hit: 10fps）
- 確認 loop 設定正確
- 執行遊戲內預覽確認效果
- 輸出：更新 `.tscn` 場景檔案

---

## 每個角色必備動畫狀態

| 狀態 | 幀數 | FPS | 循環 | 說明 |
|------|------|-----|------|------|
| `idle` | 4-8 幀 | 8 | ✅ | 待機動畫 |
| `attack` | 6-8 幀 | 12 | ❌ | 攻擊動畫 |
| `hit` | 4 幀 | 10 | ❌ | 命中目標 |
| `hurt` | 3 幀 | 10 | ❌ | 受傷（BOSS 反擊）|
| `bigwin` | 8-12 幀 | 10 | ✅ | 大獎慶祝 |
| `skill` | 8 幀 | 12 | ❌ | 技能釋放 |
| `bonus` | 6 幀 | 10 | ✅ | Bonus 遊戲特殊動畫 |
| `fail` | 4 幀 | 8 | ❌ | 失敗/空彈 |

### 角色清單
- `chiikawa`（吉伊卡哇）：LV1-3
- `hachiware`（小八）：LV4-7
- `usagi`（烏薩奇）：LV8-10

---

## Frame Consistency 檢查清單

每幀必須通過以下所有檢查：

| 項目 | 說明 | 容差 |
|------|------|------|
| canvas_size | 所有幀的畫布尺寸必須相同 | 0px |
| transparent_bg | 背景必須完全透明（alpha=0）| 0 |
| anchor_point | 底部中心 anchor 位置一致 | ±2px |
| bottom_alignment | 角色底部對齊基準線 | ±3px |
| silhouette | 輪廓形狀與參考圖一致 | 相似度 >= 85% |
| head_ratio | 頭部佔全身比例一致 | ±5% |
| weapon_position | 武器/道具位置相對角色一致 | ±4px |
| color_drift | 主要顏色不得偏移 | ΔE < 5 |
| deformation | 無異常形變（拉伸/壓縮）| 面積差 < 10% |
| jitter | 相鄰幀位置抖動 | ±2px |

---

## Animation Work Report 格式

每次完成動畫工作後，輸出以下格式的報告：

```markdown
# Animation Work Report

**角色**：<name>
**動畫狀態**：<state>
**完成日期**：YYYY-MM-DD
**執行者**：Animation Agent

## 幀資訊
- 總幀數：N
- Canvas Size：WxH px
- FPS：N
- 循環：是/否

## Consistency 分數
- 整體分數：XX/100
- canvas_size：通過/失敗
- transparent_bg：通過/失敗
- anchor_point：通過/失敗（偏差：±Xpx）
- bottom_alignment：通過/失敗（偏差：±Xpx）
- silhouette：通過/失敗（相似度：XX%）
- head_ratio：通過/失敗（偏差：±X%）
- weapon_position：通過/失敗（偏差：±Xpx）
- color_drift：通過/失敗（ΔE：X.X）
- deformation：通過/失敗（面積差：X%）
- jitter：通過/失敗（最大抖動：±Xpx）

## 已知問題
- [問題描述]

## 改善建議
- [建議內容]

## 審核結果
- Art Director 審核：通過/待修正
- 備註：[備註]
```

---

## 工具指令

```bash
# 審查所有動畫
py tools/animation_pipeline.py --audit

# 生成特定角色的預覽 GIF
py tools/animation_pipeline.py --gif chiikawa
py tools/animation_pipeline.py --gif hachiware
py tools/animation_pipeline.py --gif usagi

# 檢查單個 spritesheet
py tools/animation_pipeline.py --check client/chiikawa-pixel/assets/sprites/chars/chiikawa/chiikawa_idle.png

# 生成完整報告
py tools/animation_pipeline.py --audit --report
```

---

## 品質門檻

| 指標 | 門檻 | 說明 |
|------|------|------|
| Consistency Score | >= 85 | 單個動畫的幀一致性 |
| Animation Quality | >= 88 | 整體動畫品質（所有角色平均）|
| 禁止 merge 條件 | < 88 | 低於門檻禁止合併到主分支 |
