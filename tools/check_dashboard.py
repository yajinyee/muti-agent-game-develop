import json

with open('monitoring/grafana/provisioning/dashboards/chiikawa-overview.json', encoding='utf-8') as f:
    d = json.load(f)

print('panels:', len(d['panels']))
for p in d['panels']:
    print(f"  id={p['id']}: {p['title']} type={p['type']}")
