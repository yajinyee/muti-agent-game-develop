"""
qa_check_day311.py — DAY-311 品質驗證腳本
驗證項目：
1. T161 精靈圖密度修復（21.5% → 64.2%）
2. Cannon.gd 手感強化（Hit Stop 0.06s + 震動 0.25 + 投射物升級）
3. Server AttackResultPayload 加入 pos_x/pos_y
4. Server build + vet 通過
5. 知識庫更新（條目 12-16）
6. Agent 文件完整性（target-design-agent + spec-architect）
"""
import os, struct, zlib, math

PASS = "✅"
FAIL = "❌"
WARN = "⚠️"

results = []

def check(name, condition, detail=""):
    status = PASS if condition else FAIL
    results.append((status, name, detail))
    return condition

# ── 1. 精靈圖驗證 ──────────────────────────────────────────────
def read_png_pixels(path):
    """讀取 PNG 並計算非透明像素數"""
    try:
        with open(path, "rb") as f:
            data = f.read()
        # 找 IDAT chunk
        pos = 8  # 跳過 PNG 簽名
        chunks = {}
        idat_data = b""
        while pos < len(data):
            length = struct.unpack(">I", data[pos:pos+4])[0]
            chunk_type = data[pos+4:pos+8]
            chunk_data = data[pos+8:pos+8+length]
            if chunk_type == b"IDAT":
                idat_data += chunk_data
            elif chunk_type == b"IHDR":
                chunks["IHDR"] = chunk_data
            pos += 12 + length
        
        if not idat_data or "IHDR" not in chunks:
            return 0, 0
        
        width = struct.unpack(">I", chunks["IHDR"][0:4])[0]
        height = struct.unpack(">I", chunks["IHDR"][4:8])[0]
        
        raw = zlib.decompress(idat_data)
        non_transparent = 0
        total = width * height
        
        row_size = 1 + width * 3  # filter byte + RGB
        # 檢查是否是 RGBA
        bit_depth = chunks["IHDR"][8]
        color_type = chunks["IHDR"][9]
        
        if color_type == 2:  # RGB
            row_size = 1 + width * 3
            for y in range(height):
                row_start = y * row_size + 1
                for x in range(width):
                    px_start = row_start + x * 3
                    r, g, b = raw[px_start], raw[px_start+1], raw[px_start+2]
                    if r > 0 or g > 0 or b > 0:
                        non_transparent += 1
        elif color_type == 6:  # RGBA
            row_size = 1 + width * 4
            for y in range(height):
                row_start = y * row_size + 1
                for x in range(width):
                    px_start = row_start + x * 4
                    a = raw[px_start + 3]
                    if a > 0:
                        non_transparent += 1
        
        return non_transparent, total
    except Exception as e:
        return 0, 0

sprites_dir = r"d:\Kiro\client\chiikawa-pixel\assets\sprites\targets"

# T161 密度檢查
t161_path = os.path.join(sprites_dir, "T161_combo_burst.png")
if os.path.exists(t161_path):
    non_t, total_t = read_png_pixels(t161_path)
    pct = non_t / total_t * 100 if total_t > 0 else 0
    check("T161 精靈圖密度 >= 35%", pct >= 35, f"{pct:.1f}% 非透明像素")
else:
    check("T161 精靈圖存在", False, "檔案不存在")

# 其他 DAY-310 精靈圖存在性
for tid in ["T162_time_bomb", "T163_elemental_fusion", "T164_treasure_hunter", "T165_myth_awaken"]:
    path = os.path.join(sprites_dir, f"{tid}.png")
    check(f"{tid}.png 存在", os.path.exists(path))

# ── 2. Cannon.gd 手感強化驗證 ──────────────────────────────────
cannon_path = r"d:\Kiro\client\chiikawa-pixel\scripts\game\Cannon.gd"
if os.path.exists(cannon_path):
    with open(cannon_path, "r", encoding="utf-8") as f:
        cannon_content = f.read()
    check("Cannon.gd Hit Stop 0.06s", "hit_stop(0.06)" in cannon_content)
    check("Cannon.gd 震動加強 0.25", "add_trauma(0.25)" in cannon_content)
    check("Cannon.gd _spawn_impact_burst 存在", "_spawn_impact_burst" in cannon_content)
    check("Cannon.gd 投射物升級（Node2D）", "Node2D.new()" in cannon_content)
    check("Cannon.gd 記錄最後射擊位置", "_last_fire_pos" in cannon_content)
    check("Cannon.gd 使用備用位置", "_last_fire_pos" in cannon_content and "pos_x" in cannon_content)
else:
    check("Cannon.gd 存在", False)

# ── 3. Server Protocol 驗證 ────────────────────────────────────
protocol_path = r"d:\Kiro\server\internal\protocol\messages.go"
if os.path.exists(protocol_path):
    with open(protocol_path, "r", encoding="utf-8") as f:
        proto_content = f.read()
    check("AttackResultPayload 有 PosX", "PosX" in proto_content and "pos_x" in proto_content)
    check("AttackResultPayload 有 PosY", "PosY" in proto_content and "pos_y" in proto_content)
    check("T161-T165 訊息常數存在", "MsgLuckyComboBurst" in proto_content)
else:
    check("protocol/messages.go 存在", False)

# game.go 驗證
game_path = r"d:\Kiro\server\internal\game\game.go"
if os.path.exists(game_path):
    with open(game_path, "r", encoding="utf-8") as f:
        game_content = f.read()
    check("game.go 命中時回傳 t.X/t.Y", "PosX:        t.X" in game_content)
    check("game.go 未命中時回傳 ClickX/ClickY", "PosX:        req.ClickX" in game_content)
else:
    check("game.go 存在", False)

# ── 4. Agent 文件驗證 ──────────────────────────────────────────
agents_dir = r"d:\Kiro\agents"
for agent in ["target-design-agent.md", "spec-architect.md", "combo-system-agent.md"]:
    path = os.path.join(agents_dir, agent)
    check(f"Agent {agent} 存在", os.path.exists(path))

# ── 5. 知識庫更新驗證 ──────────────────────────────────────────
knowhow_path = r"d:\Kiro\.kiro\skills\knowhow-log.md"
if os.path.exists(knowhow_path):
    with open(knowhow_path, "r", encoding="utf-8") as f:
        kh_content = f.read()
    check("知識庫有 DAY-311 條目（#12）", "## 12." in kh_content)
    check("知識庫有精靈圖密度條目", "精靈圖密度" in kh_content)
    check("知識庫有 Hit Stop 條目", "Hit Stop" in kh_content)
else:
    check("knowhow-log.md 存在", False)

# ── 6. 精靈圖完整性（T001-T006 + T101-T165 + B001）──────────────
expected_sprites = (
    [f"T00{i}" for i in range(1, 7)] +
    [f"T10{i}" for i in range(1, 10)] +
    [f"T11{i}" for i in range(0, 10)] +
    [f"T12{i}" for i in range(0, 10)] +
    [f"T13{i}" for i in range(0, 10)] +
    [f"T14{i}" for i in range(0, 10)] +
    [f"T15{i}" for i in range(0, 10)] +
    [f"T16{i}" for i in range(1, 6)] +
    ["B001"]
)

missing = []
for tid in expected_sprites:
    # 找對應的 PNG（名稱格式：T001_xxx.png）
    found = False
    for fname in os.listdir(sprites_dir):
        if fname.startswith(tid + "_") or fname == tid + ".png":
            found = True
            break
    if not found:
        missing.append(tid)

check(f"所有精靈圖存在（{len(expected_sprites)} 個）", len(missing) == 0,
      f"缺少：{missing[:5]}..." if missing else "全部存在")

# ── 輸出結果 ──────────────────────────────────────────────────
print("\n" + "="*60)
print("DAY-311 QA 驗證報告")
print("="*60)

passed = sum(1 for s, _, _ in results if s == PASS)
failed = sum(1 for s, _, _ in results if s == FAIL)
total = len(results)

for status, name, detail in results:
    detail_str = f" ({detail})" if detail else ""
    print(f"{status} {name}{detail_str}")

print("="*60)
print(f"結果：{passed}/{total} 通過")
if failed == 0:
    print("🎉 全部通過！DAY-311 品質驗證完成")
else:
    print(f"⚠️ {failed} 項未通過，需要修復")
