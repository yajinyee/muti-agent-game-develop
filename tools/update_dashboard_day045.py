"""
DAY-045: 更新 Grafana Dashboard，加入 Client 端效能面板（19-21）
- Panel 19: Client 平均 FPS stat
- Panel 20: Client 記憶體使用 stat
- Panel 21: Client 端效能趨勢 timeseries（FPS + 記憶體）
"""
import json

DASHBOARD_PATH = 'monitoring/grafana/provisioning/dashboards/chiikawa-overview.json'

with open(DASHBOARD_PATH, encoding='utf-8') as f:
    dashboard = json.load(f)

# 確認目前面板數
print(f"目前面板數: {len(dashboard['panels'])}")

# 計算新面板的 y 位置（每行 8 高，每個 stat 面板 4 寬）
# 現有 18 個面板，最後一行是 16-18（ping 相關）
# 新面板放在下一行

new_panels = [
    {
        "id": 19,
        "title": "Client 平均 FPS",
        "type": "stat",
        "gridPos": {"h": 4, "w": 6, "x": 0, "y": 40},
        "datasource": {"type": "prometheus", "uid": "prometheus"},
        "fieldConfig": {
            "defaults": {
                "color": {"mode": "thresholds"},
                "thresholds": {
                    "mode": "absolute",
                    "steps": [
                        {"color": "red", "value": None},
                        {"color": "yellow", "value": 30},
                        {"color": "green", "value": 50}
                    ]
                },
                "unit": "short",
                "displayName": "Avg FPS"
            }
        },
        "options": {
            "reduceOptions": {"calcs": ["lastNotNull"]},
            "orientation": "auto",
            "textMode": "auto",
            "colorMode": "background",
            "graphMode": "area"
        },
        "targets": [
            {
                "datasource": {"type": "prometheus", "uid": "prometheus"},
                "expr": "chiikawa_client_avg_fps",
                "legendFormat": "平均 FPS",
                "refId": "A"
            }
        ]
    },
    {
        "id": 20,
        "title": "Client 記憶體使用（MB）",
        "type": "stat",
        "gridPos": {"h": 4, "w": 6, "x": 6, "y": 40},
        "datasource": {"type": "prometheus", "uid": "prometheus"},
        "fieldConfig": {
            "defaults": {
                "color": {"mode": "thresholds"},
                "thresholds": {
                    "mode": "absolute",
                    "steps": [
                        {"color": "green", "value": None},
                        {"color": "yellow", "value": 150},
                        {"color": "red", "value": 250}
                    ]
                },
                "unit": "decmbytes",
                "displayName": "Memory"
            }
        },
        "options": {
            "reduceOptions": {"calcs": ["lastNotNull"]},
            "orientation": "auto",
            "textMode": "auto",
            "colorMode": "background",
            "graphMode": "area"
        },
        "targets": [
            {
                "datasource": {"type": "prometheus", "uid": "prometheus"},
                "expr": "avg(chiikawa_client_memory_mb)",
                "legendFormat": "平均記憶體 MB",
                "refId": "A"
            }
        ]
    },
    {
        "id": 21,
        "title": "Client 端效能趨勢（FPS + 記憶體）",
        "type": "timeseries",
        "gridPos": {"h": 8, "w": 12, "x": 12, "y": 40},
        "datasource": {"type": "prometheus", "uid": "prometheus"},
        "fieldConfig": {
            "defaults": {
                "color": {"mode": "palette-classic"},
                "custom": {
                    "lineWidth": 2,
                    "fillOpacity": 10,
                    "showPoints": "never"
                }
            },
            "overrides": [
                {
                    "matcher": {"id": "byName", "options": "記憶體 MB"},
                    "properties": [
                        {"id": "custom.axisPlacement", "value": "right"},
                        {"id": "unit", "value": "decmbytes"}
                    ]
                }
            ]
        },
        "options": {
            "tooltip": {"mode": "multi"},
            "legend": {"displayMode": "list", "placement": "bottom"}
        },
        "targets": [
            {
                "datasource": {"type": "prometheus", "uid": "prometheus"},
                "expr": "chiikawa_client_avg_fps",
                "legendFormat": "平均 FPS",
                "refId": "A"
            },
            {
                "datasource": {"type": "prometheus", "uid": "prometheus"},
                "expr": "avg(chiikawa_client_memory_mb)",
                "legendFormat": "記憶體 MB",
                "refId": "B"
            }
        ]
    }
]

# 加入新面板
dashboard['panels'].extend(new_panels)

# 更新 title 加入版本標記
dashboard['title'] = '吉伊卡哇：像素大討伐 — 監控總覽 (DAY-045)'

print(f"更新後面板數: {len(dashboard['panels'])}")

with open(DASHBOARD_PATH, 'w', encoding='utf-8') as f:
    json.dump(dashboard, f, ensure_ascii=False, indent=2)

print("✅ Grafana dashboard 已更新（21 個面板）")
