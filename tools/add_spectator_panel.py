#!/usr/bin/env python3
"""在 Grafana dashboard 加入觀戰者計數面板（第 26 個）"""
import json
import re

path = r"d:\Kiro\monitoring\grafana\provisioning\dashboards\chiikawa-overview.json"

with open(path, "r", encoding="utf-8") as f:
    content = f.read()

data = json.loads(content)

# 更新 title
data["title"] = "吉伊卡哇：像素大討伐 — 監控總覽 (DAY-055)"

# 新增觀戰者計數面板
spectator_panel = {
    "id": 26,
    "title": "當前觀戰者數量",
    "type": "stat",
    "gridPos": {
        "x": 16,
        "y": 0,
        "w": 4,
        "h": 4
    },
    "targets": [
        {
            "expr": "chiikawa_connected_spectators",
            "legendFormat": "Spectators"
        }
    ],
    "options": {
        "colorMode": "background",
        "graphMode": "area",
        "textMode": "auto"
    },
    "fieldConfig": {
        "defaults": {
            "thresholds": {
                "steps": [
                    {"color": "blue", "value": None},
                    {"color": "green", "value": 1},
                    {"color": "yellow", "value": 5}
                ]
            }
        }
    }
}

# 加入到 panels 列表
data["panels"].append(spectator_panel)

with open(path, "w", encoding="utf-8") as f:
    json.dump(data, f, ensure_ascii=False, indent=2)

print(f"SUCCESS: Added spectator panel (total panels: {len(data['panels'])})")
