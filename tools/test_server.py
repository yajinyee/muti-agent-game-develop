# -*- coding: utf-8 -*-
"""
Go Server WebSocket integration test
"""
import asyncio
import json
import sys

try:
    import websockets
except ImportError:
    print("Install: py -m pip install websockets")
    sys.exit(1)

SERVER_URL = "ws://localhost:7777/ws?player_id=test_player_001"
RESULTS = []

async def send(ws, msg_type, payload=None):
    msg = {"type": msg_type, "payload": payload or {}}
    await ws.send(json.dumps(msg))

async def recv_until(ws, expected_type, timeout=5.0):
    try:
        async with asyncio.timeout(timeout):
            while True:
                raw = await ws.recv()
                msg = json.loads(raw)
                if msg.get("type") == expected_type:
                    return msg
    except asyncio.TimeoutError:
        return None

async def run_tests():
    print("=" * 50)
    print("Chiikawa Game Server - Integration Test")
    print("=" * 50)

    try:
        async with websockets.connect(SERVER_URL) as ws:
            print("OK WebSocket connected")
            RESULTS.append(("WebSocket connect", True))

            # 1. initial game_state
            msg = await recv_until(ws, "game_state", timeout=3.0)
            if msg:
                state = msg["payload"].get("state", "")
                ok = state == "normal_play"
                print(f"{'OK' if ok else 'NG'} Initial state: {state}")
                RESULTS.append(("Initial state", ok))
                # 如果不是 normal_play，等待轉換
                if not ok:
                    for _ in range(10):
                        msg2 = await recv_until(ws, "game_state", timeout=5.0)
                        if msg2 and msg2["payload"].get("state") == "normal_play":
                            print("OK Transitioned to normal_play")
                            break
            else:
                print("NG No initial state")
                RESULTS.append(("Initial state", False))

            # 2. target spawn
            await asyncio.sleep(1.5)
            msg = await recv_until(ws, "target_spawn", timeout=5.0)
            target_id = ""
            if msg:
                target_id = msg["payload"].get("instance_id", "")
                target_name = msg["payload"].get("name", "?")
                print(f"OK Target spawned: {target_name} ({target_id[:8]}...)")
                RESULTS.append(("Target spawn", True))

                # 3. attack
                await send(ws, "attack", {
                    "target_id": target_id,
                    "click_x": 640.0,
                    "click_y": 360.0
                })
                result = await recv_until(ws, "attack_result", timeout=3.0)
                if result:
                    is_hit = result["payload"].get("is_hit", False)
                    print(f"{'OK' if is_hit else 'NG'} Attack hit={is_hit}")
                    RESULTS.append(("Attack hit", is_hit))
                else:
                    print("NG No attack result")
                    RESULTS.append(("Attack hit", False))
            else:
                print("NG No target spawned")
                RESULTS.append(("Target spawn", False))

            # 4. bet change
            await send(ws, "bet_change", {"bet_level": 5})
            bet_ok = False
            for _ in range(5):
                msg = await recv_until(ws, "player_update", timeout=3.0)
                if msg and msg["payload"].get("bet_level", 0) == 5:
                    char_id = msg["payload"].get("character_id", "")
                    print(f"OK Bet LV5 char={char_id}")
                    RESULTS.append(("Bet change", char_id == "hachiware"))
                    bet_ok = True
                    break
            if not bet_ok:
                print("NG Bet change failed")
                RESULTS.append(("Bet change", False))

            # 5. trigger boss (需要在 normal_play 狀態)
            await asyncio.sleep(0.5)
            await send(ws, "trigger_boss", {})
            msg = await recv_until(ws, "boss_event", timeout=8.0)
            if msg:
                event = msg["payload"].get("event", "")
                print(f"{'OK' if event == 'warning' else 'NG'} Boss event: {event}")
                RESULTS.append(("Boss trigger", event == "warning"))
            else:
                print("NG No boss event")
                RESULTS.append(("Boss trigger", False))

            # 6. ping/pong
            await send(ws, "ping", {})
            msg = await recv_until(ws, "pong", timeout=3.0)
            print(f"{'OK' if msg else 'NG'} Ping/Pong")
            RESULTS.append(("Ping/Pong", msg is not None))

    except ConnectionRefusedError:
        print("NG Cannot connect to server")
        RESULTS.append(("WebSocket connect", False))
    except Exception as e:
        print(f"NG Error: {e}")
        RESULTS.append(("Test run", False))

    print("\n" + "=" * 50)
    passed = sum(1 for _, ok in RESULTS if ok)
    total = len(RESULTS)
    for name, ok in RESULTS:
        print(f"  {'OK' if ok else 'NG'} {name}")
    print(f"\nPassed: {passed}/{total}")
    return passed == total

if __name__ == "__main__":
    success = asyncio.run(run_tests())
    sys.exit(0 if success else 1)
