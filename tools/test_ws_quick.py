#!/usr/bin/env python3
"""快速 WebSocket 連線測試"""
import asyncio
import json
import websockets

async def test():
    url = "ws://localhost:7777/ws?player_id=test_quick"
    print(f"Connecting to {url}")
    async with websockets.connect(url) as ws:
        print("Connected!")
        for i in range(15):
            try:
                msg = await asyncio.wait_for(ws.recv(), timeout=2.0)
                data = json.loads(msg)
                msg_type = data.get("type", "unknown")
                print(f"  [{i+1}] Received: {msg_type}")
                if msg_type in ["game_state", "player_update", "target_spawn"]:
                    print(f"       payload keys: {list(data.get('payload', {}).keys())}")
            except asyncio.TimeoutError:
                print(f"  [{i+1}] Timeout (no more messages)")
                break

asyncio.run(test())
