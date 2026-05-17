# Animation Audit Report

**日期**：2026-05-17  
**執行者**：Animation Agent  
**整體分數**：87/100  
**門檻**：>= 88（⚠️ 接近門檻，需持續改善）

---

## 摘要

| 項目 | 數量 |
|------|------|
| 總動畫數 | 24（3 角色 × 8 狀態）|
| 通過（>= 85）| 18 ✅ |
| 未通過（< 85）| 3 ❌ |
| 缺失 | 3 ⚠️ |
| 平均分數 | 87/100 |

---

## 各角色詳細結果

### chiikawa（吉伊卡哇）

| 動畫狀態 | 狀態 | 分數 | 幀數 | 問題 |
|---------|------|------|------|------|
| idle | ✅ passed | 91 | 4 | 無 |
| attack | ✅ passed | 89 | 6 | 輕微 jitter（±2px）|
| hit | ✅ passed | 92 | 4 | 無 |
| hurt | ✅ passed | 88 | 3 | 無 |
| bigwin | ✅ passed | 90 | 8 | 無 |
| skill | ✅ passed | 87 | 8 | 輕微 color drift（ΔE=4.2）|
| bonus | ✅ passed | 91 | 6 | 無 |
| fail | ✅ passed | 88 | 4 | 無 |

**角色平均分數**：89.5/100 ✅

### hachiware（小八）

| 動畫狀態 | 狀態 | 分數 | 幀數 | 問題 |
|---------|------|------|------|------|
| idle | ✅ passed | 90 | 4 | 無 |
| attack | ✅ passed | 88 | 6 | 無 |
| hit | ✅ passed | 91 | 4 | 無 |
| hurt | ❌ failed | 78 | 3 | bottom_alignment 偏差 5px；jitter ±4px |
| bigwin | ✅ passed | 86 | 8 | 輕微 silhouette 差異（82%）|
| skill | ⚠️ missing | 0 | - | spritesheet 不存在 |
| bonus | ✅ passed | 89 | 6 | 無 |
| fail | ⚠️ missing | 0 | - | spritesheet 不存在 |

**角色平均分數**：74.0/100 ⚠️（含缺失項目）/ 87.0（僅計算存在項目）

### usagi（烏薩奇）

| 動畫狀態 | 狀態 | 分數 | 幀數 | 問題 |
|---------|------|------|------|------|
| idle | ✅ passed | 93 | 4 | 無 |
| attack | ✅ passed | 91 | 6 | 無 |
| hit | ✅ passed | 90 | 4 | 無 |
| hurt | ❌ failed | 82 | 3 | anchor_point 偏移 3px |
| bigwin | ❌ failed | 80 | 8 | deformation 12%；color drift ΔE=6.1 |
| skill | ✅ passed | 88 | 8 | 無 |
| bonus | ✅ passed | 90 | 6 | 無 |
| fail | ⚠️ missing | 0 | - | spritesheet 不存在 |

**角色平均分數**：76.8/100 ⚠️（含缺失）/ 87.6（僅計算存在項目）

---

## 已知問題

### 🔴 高優先級（需立即修復）

1. **hachiware hurt 動畫**（分數 78）
   - bottom_alignment 偏差 5px（超過容差 3px）
   - jitter ±4px（超過容差 2px）
   - 修復方案：重新對齊 anchor point，使用 bottom_align_frames 工具

2. **usagi bigwin 動畫**（分數 80）
   - deformation 12%（超過容差 10%）
   - color drift ΔE=6.1（超過容差 5.0）
   - 修復方案：重新生成，固定 seed 和 IPAdapter 設定

3. **usagi hurt 動畫**（分數 82）
   - anchor_point 偏移 3px（超過容差 2px）
   - 修復方案：手動調整幀位置

### 🟡 中優先級（本週修復）

4. **hachiware skill 動畫**（缺失）
   - 需要生成 8 幀 skill 動畫
   - 預計工時：2 小時

5. **hachiware fail 動畫**（缺失）
   - 需要生成 4 幀 fail 動畫
   - 預計工時：1 小時

6. **usagi fail 動畫**（缺失）
   - 需要生成 4 幀 fail 動畫
   - 預計工時：1 小時

### 🟢 低優先級（有時間再做）

7. **chiikawa idle 幀數提升**
   - 目前 4 幀，目標 8 幀
   - 更流暢的待機動畫
   - 預計工時：3 小時

---

## 改善建議

1. **優先修復 3 個未通過動畫**（hachiware hurt、usagi bigwin、usagi hurt）
2. **補齊 3 個缺失動畫**（hachiware skill/fail、usagi fail）
3. **修復後重新執行 audit**，確認分數達到 >= 88
4. **長期目標**：所有動畫分數 >= 90，整體 Animation Quality >= 90

---

## Animation Quality 分數

| 計算方式 | 分數 |
|---------|------|
| 所有動畫平均（含缺失=0）| 82.5/100 |
| 僅計算存在的動畫 | 87.0/100 |
| **官方 Animation Quality** | **87/100** |
| 門檻 | >= 88 |
| 狀態 | ⚠️ 接近門檻 |

> 注意：Animation Quality 計算方式為「僅計算存在的動畫的平均分數」。缺失動畫不計入分數，但會在報告中標記。

---

*報告生成時間：2026-05-17 09:00:00*  
*下次審查建議：修復高優先級問題後立即重新審查*
