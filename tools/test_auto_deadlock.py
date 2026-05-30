#!/usr/bin/env python3
"""
test_auto_deadlock.py — 測試 AUTO 射擊死鎖修復
DAY-338 修復驗證：確認開啟 AUTO 後，新連線玩家仍能正常收到訊息
"""
import asyncio
import json
import websockets
import time

async def test():
    url = "ws://localhost:7777/ws"
    
    print("=== AUTO 射擊死鎖修復驗證 ===\n")
    
    # 步驟 1：連線並開啟 AUTO
    print("[1] 連線玩家 A 並開啟 AUTO 射擊...")
    ws_a = await websockets.connect(url + "?player_id=player_auto_test_A")
    
    # 等待初始訊息
    for i in range(5):
        try:
            raw = await asyncio.wait_for(ws_a.recv(), timeout=2.0)
            data = json.loads(raw)
            print(f"  玩家 A 收到: {data['type']}")
        except asyncio.TimeoutError:
            break
    
    # 開啟 AUTO
    await ws_a.send(json.dumps({"type": "auto_toggle", "payload": {}}))
    print("  玩家 A AUTO 已開啟")
    
    # 等待 AUTO 射擊觸發（3 秒）
    print("  等待 3 秒讓 AUTO 射擊觸發...")
    await asyncio.sleep(3)
    
    # 步驟 2：新玩家連線，確認能收到訊息
    print("\n[2] 新玩家 B 連線（AUTO 射擊進行中）...")
    ws_b = await websockets.connect(url + "?player_id=player_auto_test_B")
    
    received_types = []
    for i in range(10):
        try:
            raw = await asyncio.wait_for(ws_b.recv(), timeout=2.0)
            data = json.loads(raw)
            received_types.append(data['type'])
            print(f"  玩家 B 收到: {data['type']}")
        except asyncio.TimeoutError:
            print(f"  玩家 B 超時（已收到 {len(received_types)} 個訊息）")
            break
    
    # 驗證
    print("\n=== 驗證結果 ===")
    has_game_state = "game_state" in received_types or "player_update" in received_types
    if has_game_state:
        print("✅ 死鎖修復成功！AUTO 射擊進行中，新玩家仍能正常收到訊息")
    else:
        print("❌ 死鎖仍然存在！新玩家無法收到訊息")
    
    print(f"   收到的訊息類型: {received_types}")
    
    await ws_a.close()
    await ws_b.close()

asyncio.run(test())
