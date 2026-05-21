# 夜間報告模板

> 複製此模板，命名為 `nightly-report-YYYY-MM-DD.md`，由 QA Agent 在每日工作結束後填寫。

---

# 夜間報告 — [DATE]

**報告時間**：[TIME]  
**報告人**：QA Playtest Agent  
**審閱人**：Game Director

---

## 今日整體狀態

| 指標 | 昨日 | 今日 | 變化 | 狀態 |
|------|------|------|------|------|
| 完成度 | XX% | XX% | +X% | ✅/⚠️/❌ |
| 美術質量 | XX | XX | +X | ✅/⚠️/❌ |
| 規格一致性 | XX% | XX% | +X% | ✅/⚠️/❌ |

---

## 品質分數儀表板

| 指標 | 分數 | 門檻 | 狀態 | 趨勢 |
|------|------|------|------|------|
| Spec Completeness | XX | >=95 | ✅/❌ | ↑/↓/→ |
| Build Stability | XX | >=95 | ✅/❌ | ↑/↓/→ |
| Visual Consistency | XX | >=90 | ✅/❌ | ↑/↓/→ |
| Animation Quality | XX | >=88 | ✅/❌ | ↑/↓/→ |
| Audio Sync | XX | >=90 | ✅/❌ | ↑/↓/→ |
| Gameplay Feel | XX | >=85 | ✅/❌ | ↑/↓/→ |
| Balance Health | XX | >=90 | ✅/❌ | ↑/↓/→ |
| Regression Risk | XX | <=10 | ✅/❌ | ↑/↓/→ |

**整體評級**：🟢 良好 / 🟡 注意 / 🔴 警告

---

## 今日完成工作

### Game Director
- [完成項目]

### Spec Architect
- [完成項目]

### Art Director
- [完成項目]

### Sprite Generation Agent
- [完成項目]

### Animation Agent
- [完成項目]

### Audio Director
- [完成項目]

### Godot Client Agent
- [完成項目]

### Go Server Agent
- [完成項目]

### Balance Agent
- [完成項目]

### QA Playtest Agent
- [完成項目]

### Research Agent
- [完成項目]

### Skill Librarian
- [完成項目]

---

## 今日發現的問題

### 🔴 嚴重問題（需立即處理）
1. [問題描述]
   - 影響：[影響範圍]
   - 指派：[Agent]
   - 預計修復：[時間]

### 🟠 重要問題（明日優先）
1. [問題描述]
   - 影響：[影響範圍]
   - 指派：[Agent]

### 🟡 一般問題（本週內）
1. [問題描述]

---

## Build 狀態

### Go Server
```
go build ./... : ✅ 成功 / ❌ 失敗
go vet ./...   : ✅ 成功 / ❌ 失敗
go test ./...  : XX/XX 通過
```

### Godot Client
```
HTML5 匯出    : ✅ 成功 / ❌ 失敗
場景載入      : ✅ 成功 / ❌ 失敗
WebSocket 連線: ✅ 成功 / ❌ 失敗
```

---

## RTP 今日模擬結果

- 模擬局數：XX 萬局
- 實際 RTP：XX%（目標 92-96%）
- 偏差：±XX%
- 狀態：✅ 正常 / ⚠️ 偏高 / ⚠️ 偏低 / ❌ 異常

---

## 明日計畫建議

### 最高優先（P0）
1. [建議任務] — 理由：[原因]

### 重要任務（P1）
1. [建議任務] — 指派：[Agent]

### 一般任務（P2）
1. [建議任務]

---

## 今日學習與知識更新

### 新增 Skills
- [Skill 名稱]：[簡述]

### 更新 Skills
- [Skill 名稱]：[更新內容]

### 研究發現
- [發現]：[來源]

---

## Game Director 批注

> *由 Game Director 在審閱後填寫*

[批注內容]

**明日重點**：[重點說明]

---

*報告結束 — [DATE] [TIME]*
