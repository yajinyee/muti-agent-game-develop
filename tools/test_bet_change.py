#!/usr/bin/env python3
"""測試 bet_change 是否正確觸發 player_update"""
import asyncio
import json
import websockets

async def test():
    url = "ws://localhost:7777/ws?player_id=test_bet"
    async with websockets.connect(url) as ws:
        print("Connected")
        # 收集前 5 個訊息
        msgs = []
        for i in range(5):
            try:
                raw = await asyncio.wait_for(ws.recv(), timeout=3.0)
                data = json.loads(raw)
                msgs.append(data)
                print(f"  [{i+1}] {data['type']}")
            except asyncio.TimeoutError:
                print(f"  [{i+1}] Timeout")
                break
        
        print("\nSending bet_change to level 5...")
        await ws.send(json.dumps({"type": "bet_change", "payload": {"bet_level": 5}}))
        
        # 等待 player_update
        for i in range(10):
            try:
                raw = await asyncio.wait_for(ws.recv(), timeout=2.0)
                data = json.loads(raw)
                print(f"  After bet_change [{i+1}]: {data['type']}", end="")
                if data['type'] == 'player_update':
                    print(f" -> bet_level={data['payload'].get('bet_level')}")
                    break
                else:
                    print()
            except asyncio.TimeoutError:
                print(f"  After bet_change [{i+1}]: Timeout")
                break

asyncio.run(test())
