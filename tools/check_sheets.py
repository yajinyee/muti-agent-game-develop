import json

for char in ['chiikawa', 'hachiware', 'usagi']:
    with open(f'd:/Kiro/client/chiikawa-pixel/assets/sprites/sheets/{char}_animated.json') as f:
        data = json.load(f)
    print(f"{char}: cols={data['cols']}, rows={data['rows']}, frame_size={data['frame_size']}")
    for anim_name, anim in data['animations'].items():
        print(f"  {anim_name}: frames={anim['frames']}, fps={anim['fps']}")
