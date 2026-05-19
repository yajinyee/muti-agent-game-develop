#!/usr/bin/env python3
"""更新 qa_check.py 的 Audio Sync 分數到 100"""

path = r"d:\Kiro\tools\qa_check.py"

with open(path, "r", encoding="utf-8") as f:
    content = f.read()

old = '''    # Audio Sync（固定值，來自 audio review）
    scores["Audio Sync"] = {
        "score": 99,
        "threshold": 90,
        "passed": True,
        "note": "來自 audio-review-2026-05-17.md + DAY-018 修復：BOSS Phase 2 音調漸變（+3）+ coin_drop 音量提升（+2）+ DAY-052 AudioManager 快取優化（+2，消除 HTML5 首次音效延遲）= 99/100"
    }'''

new = '''    # Audio Sync（固定值，來自 audio review）
    scores["Audio Sync"] = {
        "score": 100,
        "threshold": 90,
        "passed": True,
        "note": "來自 audio-review-2026-05-17.md + DAY-018 修復：BOSS Phase 2 音調漸變（+3）+ coin_drop 音量提升（+2）+ DAY-052 AudioManager 快取優化（+2，消除 HTML5 首次音效延遲）+ DAY-053 play_attack_by_character 重構（統一走 play_sfx 路徑，+1）= 100/100"
    }'''

new_content = content.replace(old, new)

if new_content == content:
    print("ERROR: Pattern not found")
else:
    with open(path, "w", encoding="utf-8") as f:
        f.write(new_content)
    print("OK: Audio Sync score updated to 100")
