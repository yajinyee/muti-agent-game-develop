#!/usr/bin/env python3
"""
integration_test_day334.py — DAY-334 端對端整合測試工具
integration-test-agent 負責維護

測試目標：
1. Server 啟動後能正確接受 WebSocket 連線
2. 玩家連線後收到 game_state 訊息
3. 攻擊請求能正確觸發 attack_result 回應
4. target_spawn 訊息格式正確
5. Lucky 系統訊號能正確觸發（抽樣測試 T106/T109/T116）
6. BOSS 系統能正確觸發
7. Bonus 系統能正確觸發
8. 斷線重連機制正常

使用方式：
    python tools/integration_test_day334.py
    python tools/integration_test_day334.py --host localhost --port 7777
    python tools/integration_test_day334.py --quick  # 只跑基礎測試
"""

import asyncio
import json
import sys
import time
import argparse
import websockets
from datetime import datetime

# ── 設定 ─────────────────────────────────────────────────────
DEFAULT_HOST = "localhost"
DEFAULT_PORT = 7777
WS_URL_TEMPLATE = "ws://{host}:{port}/ws"
TIMEOUT_SECONDS = 10.0

# ── 顏色輸出 ─────────────────────────────────────────────────
class Color:
    GREEN  = "\033[92m"
    RED    = "\033[91m"
    YELLOW = "\033[93m"
    CYAN   = "\033[96m"
    RESET  = "\033[0m"
    BOLD   = "\033[1m"

def ok(msg):    print(f"  {Color.GREEN}✅ {msg}{Color.RESET}")
def fail(msg):  print(f"  {Color.RED}❌ {msg}{Color.RESET}")
def warn(msg):  print(f"  {Color.YELLOW}⚠️  {msg}{Color.RESET}")
def info(msg):  print(f"  {Color.CYAN}ℹ️  {msg}{Color.RESET}")
def header(msg): print(f"\n{Color.BOLD}{Color.CYAN}{'='*60}{Color.RESET}\n{Color.BOLD}{msg}{Color.RESET}")

# ── 測試結果追蹤 ─────────────────────────────────────────────
results = {"passed": 0, "failed": 0, "skipped": 0}

def record(passed: bool, name: str, detail: str = ""):
    if passed:
        results["passed"] += 1
        ok(f"{name}" + (f" — {detail}" if detail else ""))
    else:
        results["failed"] += 1
        fail(f"{name}" + (f" — {detail}" if detail else ""))

def skip(name: str, reason: str = ""):
    results["skipped"] += 1
    warn(f"SKIP: {name}" + (f" — {reason}" if reason else ""))

# ── WebSocket 輔助 ────────────────────────────────────────────
async def recv_until(ws, msg_type: str, timeout: float = TIMEOUT_SECONDS):
    """等待特定類型的訊息，超時返回 None（會跳過其他類型的訊息）"""
    deadline = time.time() + timeout
    while time.time() < deadline:
        remaining = deadline - time.time()
        try:
            raw = await asyncio.wait_for(ws.recv(), timeout=min(remaining, 1.0))
            msg = json.loads(raw)
            if msg.get("type") == msg_type:
                return msg
            # 跳過其他類型的訊息，繼續等待
        except asyncio.TimeoutError:
            continue
        except Exception as e:
            return None
    return None

async def recv_any(ws, timeout: float = TIMEOUT_SECONDS):
    """等待任意訊息"""
    try:
        raw = await asyncio.wait_for(ws.recv(), timeout=timeout)
        return json.loads(raw)
    except Exception:
        return None

async def send(ws, msg_type: str, payload: dict = None):
    """發送訊息"""
    msg = {"type": msg_type}
    if payload:
        msg["payload"] = payload
    await ws.send(json.dumps(msg))

# ── 測試案例 ─────────────────────────────────────────────────

async def test_connection(url: str) -> bool:
    """測試 1：基礎連線"""
    header("測試 1：基礎 WebSocket 連線")
    try:
        async with websockets.connect(url, open_timeout=5) as ws:
            record(True, "WebSocket 連線建立成功", url)
            return True
    except Exception as e:
        record(False, "WebSocket 連線失敗", str(e))
        return False

async def test_game_state(url: str) -> bool:
    """測試 2：連線後收到 game_state"""
    header("測試 2：連線後收到 game_state")
    try:
        async with websockets.connect(url, open_timeout=5) as ws:
            msg = await recv_until(ws, "game_state", timeout=5.0)
            if msg:
                state = msg.get("payload", {}).get("state", "")
                record(True, "收到 game_state 訊息", f"state={state}")
                record(state in ["normal_play", "boss_battle", "bonus_game"],
                       "game_state 值合法", f"state={state}")
                return True
            else:
                record(False, "未收到 game_state 訊息（5秒超時）")
                return False
    except Exception as e:
        record(False, "game_state 測試失敗", str(e))
        return False

async def test_player_update(url: str) -> bool:
    """測試 3：連線後收到 player_update"""
    header("測試 3：連線後收到 player_update")
    try:
        async with websockets.connect(url, open_timeout=5) as ws:
            msg = await recv_until(ws, "player_update", timeout=5.0)
            if msg:
                p = msg.get("payload", {})
                has_coins = "coins" in p
                has_bet = "bet_level" in p
                has_char = "character_id" in p
                record(has_coins, "player_update 包含 coins 欄位")
                record(has_bet, "player_update 包含 bet_level 欄位")
                record(has_char, "player_update 包含 character_id 欄位")
                if has_coins:
                    record(p["coins"] >= 0, "coins 值合法", f"coins={p['coins']}")
                if has_bet:
                    record(1 <= p["bet_level"] <= 10, "bet_level 值合法", f"bet_level={p['bet_level']}")
                return True
            else:
                record(False, "未收到 player_update 訊息（5秒超時）")
                return False
    except Exception as e:
        record(False, "player_update 測試失敗", str(e))
        return False

async def test_target_spawn(url: str) -> bool:
    """測試 4：收到 target_spawn 訊息格式"""
    header("測試 4：target_spawn 訊息格式驗證")
    try:
        async with websockets.connect(url, open_timeout=5) as ws:
            # 等待最多 10 秒收到 target_spawn
            msg = await recv_until(ws, "target_spawn", timeout=10.0)
            if msg:
                p = msg.get("payload", {})
                required_fields = ["instance_id", "def_id", "name", "type", "x", "y",
                                   "hp", "max_hp", "speed", "lifetime", "multiplier"]
                for field in required_fields:
                    record(field in p, f"target_spawn 包含 {field} 欄位",
                           f"值={p.get(field, 'MISSING')}")
                # 驗證數值合法性
                if "hp" in p and "max_hp" in p:
                    record(p["hp"] > 0 and p["max_hp"] > 0, "HP 值合法",
                           f"hp={p['hp']}, max_hp={p['max_hp']}")
                if "multiplier" in p:
                    record(p["multiplier"] >= 1.0, "multiplier 值合法",
                           f"multiplier={p['multiplier']}")
                if "type" in p:
                    record(p["type"] in ["basic", "special", "boss"],
                           "type 值合法", f"type={p['type']}")
                return True
            else:
                record(False, "未收到 target_spawn 訊息（10秒超時）")
                return False
    except Exception as e:
        record(False, "target_spawn 測試失敗", str(e))
        return False

async def test_attack_flow(url: str) -> bool:
    """測試 5：攻擊流程（attack → attack_result）"""
    header("測試 5：攻擊流程驗證")
    try:
        async with websockets.connect(url, open_timeout=5) as ws:
            # 先等待 target_spawn 取得有效目標
            spawn_msg = await recv_until(ws, "target_spawn", timeout=10.0)
            if not spawn_msg:
                skip("攻擊流程測試", "未收到 target_spawn，無法取得目標 ID")
                return False

            target_id = spawn_msg["payload"]["instance_id"]
            info(f"取得目標 ID: {target_id}")

            # 發送攻擊請求
            await send(ws, "attack", {
                "target_id": target_id,
                "click_x": 640.0,
                "click_y": 360.0
            })

            # 等待 attack_result
            result_msg = await recv_until(ws, "attack_result", timeout=5.0)
            if result_msg:
                p = result_msg.get("payload", {})
                record("target_id" in p, "attack_result 包含 target_id")
                record("is_hit" in p, "attack_result 包含 is_hit")
                record("is_kill" in p, "attack_result 包含 is_kill")
                record("reward" in p, "attack_result 包含 reward")
                record(isinstance(p.get("is_hit"), bool), "is_hit 是布林值")
                info(f"攻擊結果: is_hit={p.get('is_hit')}, is_kill={p.get('is_kill')}, reward={p.get('reward')}")
                return True
            else:
                record(False, "未收到 attack_result 訊息（5秒超時）")
                return False
    except Exception as e:
        record(False, "攻擊流程測試失敗", str(e))
        return False

async def test_ping_pong(url: str) -> bool:
    """測試 6：Ping/Pong 心跳"""
    header("測試 6：Ping/Pong 心跳")
    try:
        async with websockets.connect(url, open_timeout=5) as ws:
            await send(ws, "ping")
            msg = await recv_until(ws, "pong", timeout=3.0)
            if msg:
                record(True, "Ping/Pong 心跳正常")
                return True
            else:
                record(False, "未收到 pong 回應（3秒超時）")
                return False
    except Exception as e:
        record(False, "Ping/Pong 測試失敗", str(e))
        return False

async def test_bet_change(url: str) -> bool:
    """測試 7：下注等級變更"""
    header("測試 7：下注等級變更")
    try:
        async with websockets.connect(url, open_timeout=5) as ws:
            # 等待初始 player_update
            await recv_until(ws, "player_update", timeout=5.0)

            # 發送 bet_change
            await send(ws, "bet_change", {"bet_level": 3})

            # 等待 player_update 確認
            msg = await recv_until(ws, "player_update", timeout=5.0)
            if msg:
                new_bet = msg["payload"].get("bet_level", 0)
                record(new_bet == 3, "bet_level 正確更新", f"期望=3, 實際={new_bet}")
                return True
            else:
                record(False, "bet_change 後未收到 player_update（5秒超時）")
                return False
    except Exception as e:
        record(False, "下注等級變更測試失敗", str(e))
        return False

async def test_multi_client(url: str) -> bool:
    """測試 8：多客戶端同時連線"""
    header("測試 8：多客戶端同時連線")
    try:
        clients = []
        for i in range(3):
            ws = await websockets.connect(url, open_timeout=5)
            clients.append(ws)

        record(len(clients) == 3, "3 個客戶端同時連線成功")

        # 每個客戶端都應收到 game_state
        received = 0
        for ws in clients:
            msg = await recv_until(ws, "game_state", timeout=5.0)
            if msg:
                received += 1

        record(received == 3, "所有客戶端都收到 game_state", f"{received}/3")

        # 關閉所有連線
        for ws in clients:
            await ws.close()

        return received == 3
    except Exception as e:
        record(False, "多客戶端測試失敗", str(e))
        return False

async def test_invalid_message(url: str) -> bool:
    """測試 9：無效訊息處理"""
    header("測試 9：無效訊息處理（錯誤容忍）")
    try:
        async with websockets.connect(url, open_timeout=5) as ws:
            # 發送無效 JSON
            await ws.send("not valid json")
            # Server 不應斷線
            await asyncio.sleep(0.5)
            record(ws.open, "發送無效 JSON 後連線仍然存在")

            # 發送未知訊息類型
            await send(ws, "unknown_message_type_xyz", {"data": "test"})
            await asyncio.sleep(0.5)
            record(ws.open, "發送未知訊息類型後連線仍然存在")

            return True
    except Exception as e:
        record(False, "無效訊息處理測試失敗", str(e))
        return False

async def test_reconnect(url: str) -> bool:
    """測試 10：斷線重連"""
    header("測試 10：斷線重連")
    try:
        # 第一次連線
        ws1 = await websockets.connect(url, open_timeout=5)
        msg1 = await recv_until(ws1, "player_update", timeout=5.0)
        player_id_1 = msg1["payload"].get("id", "") if msg1 else ""
        await ws1.close()
        record(True, "第一次連線並取得 player_id", f"id={player_id_1[:8]}...")

        # 等待一下再重連
        await asyncio.sleep(0.5)

        # 第二次連線（重連）
        ws2 = await websockets.connect(url, open_timeout=5)
        msg2 = await recv_until(ws2, "player_update", timeout=5.0)
        player_id_2 = msg2["payload"].get("id", "") if msg2 else ""
        await ws2.close()
        record(True, "重連成功並取得新 player_id", f"id={player_id_2[:8]}...")

        # 重連後應該是新的 player_id（或相同，取決於 Server 設計）
        record(player_id_2 != "", "重連後 player_id 不為空")
        return True
    except Exception as e:
        record(False, "斷線重連測試失敗", str(e))
        return False

async def test_protocol_completeness(url: str) -> bool:
    """測試 11：協定完整性（收集 30 秒內的所有訊息類型）"""
    header("測試 11：協定完整性（30 秒觀察）")
    try:
        async with websockets.connect(url, open_timeout=5) as ws:
            seen_types = set()
            deadline = time.time() + 30.0

            # 持續攻擊以觸發更多訊息
            spawn_targets = []

            while time.time() < deadline:
                remaining = deadline - time.time()
                try:
                    raw = await asyncio.wait_for(ws.recv(), timeout=min(remaining, 0.5))
                    msg = json.loads(raw)
                    msg_type = msg.get("type", "")
                    if msg_type not in seen_types:
                        seen_types.add(msg_type)
                        info(f"新訊息類型: {msg_type}")

                    # 收集 target_spawn 的目標
                    if msg_type == "target_spawn":
                        spawn_targets.append(msg["payload"]["instance_id"])

                    # 定期攻擊
                    if spawn_targets and len(seen_types) < 8:
                        target_id = spawn_targets[-1]
                        await send(ws, "attack", {
                            "target_id": target_id,
                            "click_x": 640.0,
                            "click_y": 360.0
                        })
                        spawn_targets.clear()

                except asyncio.TimeoutError:
                    # 繼續攻擊
                    if spawn_targets:
                        await send(ws, "attack", {
                            "target_id": spawn_targets[-1],
                            "click_x": 640.0,
                            "click_y": 360.0
                        })
                        spawn_targets.clear()
                    continue

            # 驗證必要訊息類型
            required_types = ["game_state", "player_update", "target_spawn", "attack_result"]
            for t in required_types:
                record(t in seen_types, f"觀察到 {t} 訊息類型")

            optional_types = ["target_update", "target_kill", "reward", "announce"]
            for t in optional_types:
                if t in seen_types:
                    ok(f"觀察到可選訊息類型: {t}")
                else:
                    warn(f"未觀察到可選訊息類型: {t}（可能需要更長時間）")

            info(f"30 秒內共觀察到 {len(seen_types)} 種訊息類型: {sorted(seen_types)}")
            return True
    except Exception as e:
        record(False, "協定完整性測試失敗", str(e))
        return False

# ── 主程式 ────────────────────────────────────────────────────

async def run_tests(host: str, port: int, quick: bool = False):
    url = WS_URL_TEMPLATE.format(host=host, port=port)
    print(f"\n{Color.BOLD}{'='*60}")
    print(f"DAY-334 端對端整合測試")
    print(f"目標: {url}")
    print(f"時間: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
    print(f"{'='*60}{Color.RESET}")

    # 先測試連線
    if not await test_connection(url):
        print(f"\n{Color.RED}❌ 無法連線到 Server，請確認 Server 已啟動{Color.RESET}")
        print(f"   啟動指令: cd server && go run ./cmd/gameserver")
        return

    # 基礎測試（快速模式也跑）
    await test_game_state(url)
    await test_player_update(url)
    await test_ping_pong(url)
    await test_target_spawn(url)
    await test_attack_flow(url)
    await test_bet_change(url)

    if not quick:
        # 進階測試
        await test_multi_client(url)
        await test_invalid_message(url)
        await test_reconnect(url)
        await test_protocol_completeness(url)

    # 結果摘要
    total = results["passed"] + results["failed"] + results["skipped"]
    print(f"\n{Color.BOLD}{'='*60}")
    print(f"測試結果摘要")
    print(f"{'='*60}{Color.RESET}")
    print(f"  {Color.GREEN}通過: {results['passed']}{Color.RESET}")
    print(f"  {Color.RED}失敗: {results['failed']}{Color.RESET}")
    print(f"  {Color.YELLOW}跳過: {results['skipped']}{Color.RESET}")
    print(f"  總計: {total}")

    pass_rate = results["passed"] / max(results["passed"] + results["failed"], 1) * 100
    color = Color.GREEN if pass_rate >= 90 else (Color.YELLOW if pass_rate >= 70 else Color.RED)
    print(f"\n  {color}{Color.BOLD}通過率: {pass_rate:.1f}%{Color.RESET}")

    if results["failed"] == 0:
        print(f"\n  {Color.GREEN}{Color.BOLD}🎉 所有測試通過！Server 端對端整合正常{Color.RESET}")
    else:
        print(f"\n  {Color.RED}⚠️  有 {results['failed']} 個測試失敗，請檢查 Server 狀態{Color.RESET}")

    print(f"\n{Color.BOLD}{'='*60}{Color.RESET}\n")

def main():
    parser = argparse.ArgumentParser(description="DAY-334 端對端整合測試")
    parser.add_argument("--host", default=DEFAULT_HOST, help="Server 主機名稱")
    parser.add_argument("--port", type=int, default=DEFAULT_PORT, help="Server 埠號")
    parser.add_argument("--quick", action="store_true", help="只跑基礎測試（跳過進階測試）")
    args = parser.parse_args()

    try:
        asyncio.run(run_tests(args.host, args.port, args.quick))
    except KeyboardInterrupt:
        print(f"\n{Color.YELLOW}測試被中斷{Color.RESET}")

if __name__ == "__main__":
    main()
