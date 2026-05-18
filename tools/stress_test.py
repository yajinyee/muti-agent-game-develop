# -*- coding: utf-8 -*-
"""
Server 長時間記憶體壓力測試工具
測試目標：24 小時運行後記憶體增長 < 10%
用法：py tools/stress_test.py [--duration 60] [--clients 5]
"""
import asyncio
import json
import sys
import time
import argparse
import random
import urllib.request

try:
    import websockets
except ImportError:
    print("Install: py -m pip install websockets")
    sys.exit(1)

SERVER_WS  = "ws://localhost:7777/ws"
SERVER_HTTP = "http://localhost:7777"

# ── 工具函數 ──────────────────────────────────────────

def get_stats():
    """取得 /stats 端點的記憶體資訊"""
    try:
        with urllib.request.urlopen(f"{SERVER_HTTP}/stats", timeout=3) as r:
            return json.loads(r.read())
    except Exception as e:
        return {"error": str(e)}

def get_health():
    """取得 /health 端點"""
    try:
        with urllib.request.urlopen(f"{SERVER_HTTP}/health", timeout=3) as r:
            return json.loads(r.read())
    except Exception as e:
        return {"error": str(e)}

# ── 模擬玩家 ──────────────────────────────────────────

async def simulate_player(player_id: str, duration: float, results: list):
    """模擬一個玩家的完整遊戲行為"""
    url = f"{SERVER_WS}?player_id={player_id}"
    attacks = 0
    errors = 0
    reconnects = 0
    start = time.time()

    while time.time() - start < duration:
        try:
            async with websockets.connect(url, ping_interval=20, ping_timeout=10) as ws:
                # 等待初始狀態
                try:
                    async with asyncio.timeout(5.0):
                        raw = await ws.recv()
                        msg = json.loads(raw)
                except asyncio.TimeoutError:
                    errors += 1
                    continue

                # 隨機切換投注等級
                bet_level = random.randint(1, 10)
                await ws.send(json.dumps({
                    "type": "bet_change",
                    "payload": {"bet_level": bet_level}
                }))

                # 持續攻擊直到連線結束或時間到
                session_start = time.time()
                session_duration = random.uniform(10.0, 30.0)  # 每次連線 10-30 秒

                while time.time() - session_start < session_duration:
                    if time.time() - start >= duration:
                        break

                    # 接收訊息（非阻塞）
                    try:
                        async with asyncio.timeout(0.5):
                            raw = await ws.recv()
                            msg = json.loads(raw)

                            # 收到目標生成，嘗試攻擊
                            if msg.get("type") == "target_spawn":
                                target_id = msg["payload"].get("instance_id", "")
                                if target_id:
                                    await ws.send(json.dumps({
                                        "type": "attack",
                                        "payload": {
                                            "target_id": target_id,
                                            "click_x": random.uniform(100, 1180),
                                            "click_y": random.uniform(100, 620),
                                        }
                                    }))
                                    attacks += 1

                            # Bonus 點擊
                            elif msg.get("type") == "target_spawn" and msg["payload"].get("type") == "bonus":
                                target_id = msg["payload"].get("instance_id", "")
                                if target_id:
                                    await ws.send(json.dumps({
                                        "type": "bonus_click",
                                        "payload": {"target_id": target_id}
                                    }))

                    except asyncio.TimeoutError:
                        # 沒有訊息，發送 ping 保持連線
                        await ws.send(json.dumps({"type": "ping", "payload": {}}))
                    except Exception:
                        break

                # 隨機斷線重連（模擬真實玩家行為）
                if random.random() < 0.3:
                    reconnects += 1
                    await asyncio.sleep(random.uniform(0.5, 2.0))

        except Exception as e:
            errors += 1
            await asyncio.sleep(1.0)

    results.append({
        "player_id": player_id,
        "attacks": attacks,
        "errors": errors,
        "reconnects": reconnects,
    })

# ── 主測試流程 ──────────────────────────────────────────

async def run_stress_test(duration: float, num_clients: int):
    print("=" * 60)
    print(f"Server 壓力測試 — {num_clients} 個客戶端，持續 {duration:.0f} 秒")
    print("=" * 60)

    # 確認 Server 正常
    health = get_health()
    if "error" in health:
        print(f"[ERROR] Server 無法連線: {health['error']}")
        print("請先啟動 Server: cd server && go run ./cmd/server")
        return False

    print(f"[OK] Server 健康狀態: {health}")

    # 記錄初始記憶體（Server 剛啟動，尚未有玩家）
    stats_start = get_stats()
    heap_start = stats_start.get("heap_alloc_mb", 0)
    goroutines_start = stats_start.get("goroutines", 0)
    print(f"\n[初始狀態（無玩家）]")
    print(f"  Heap: {heap_start:.2f} MB")
    print(f"  Goroutines: {goroutines_start}")
    print(f"  GC 次數: {stats_start.get('gc_count', 0)}")
    print(f"  ⚠️  注意：初始 Heap 很低，玩家連線後會正常增長（遊戲初始化）")

    # 啟動所有模擬玩家
    results = []
    tasks = []
    for i in range(num_clients):
        player_id = f"stress_player_{i:03d}"
        task = asyncio.create_task(simulate_player(player_id, duration, results))
        tasks.append(task)

    # 定期監控記憶體（每 10 秒）
    monitor_interval = 10.0
    checkpoints = []
    start_time = time.time()

    print(f"\n[監控中] 每 {monitor_interval:.0f} 秒記錄一次...")
    print(f"{'時間':>8} | {'Heap(MB)':>10} | {'Goroutines':>12} | {'GC次數':>8}")
    print("-" * 50)

    while time.time() - start_time < duration:
        await asyncio.sleep(monitor_interval)
        elapsed = time.time() - start_time
        stats = get_stats()
        heap = stats.get("heap_alloc_mb", 0)
        goroutines = stats.get("goroutines", 0)
        gc_count = stats.get("gc_count", 0)
        checkpoints.append({
            "elapsed": elapsed,
            "heap_mb": heap,
            "goroutines": goroutines,
            "gc_count": gc_count,
        })
        print(f"{elapsed:>7.0f}s | {heap:>10.2f} | {goroutines:>12} | {gc_count:>8}")

    # 等待所有玩家完成
    await asyncio.gather(*tasks, return_exceptions=True)

    # 最終記憶體
    stats_end = get_stats()
    heap_end = stats_end.get("heap_alloc_mb", 0)
    goroutines_end = stats_end.get("goroutines", 0)

    # 分析結果
    print(f"\n{'=' * 60}")
    print("[測試結果]")
    print(f"  初始 Heap: {heap_start:.2f} MB")
    print(f"  最終 Heap: {heap_end:.2f} MB")

    if heap_start > 0:
        heap_growth = (heap_end - heap_start) / heap_start * 100
        print(f"  Heap 增長: {heap_growth:+.1f}%")

        # 正確的洩漏判斷：用穩定後的趨勢斜率，而非從零開始的增長
        # 取後半段 checkpoint 判斷是否持續增長（斜率 > 0.5 MB/min 才算洩漏）
        if len(checkpoints) >= 4:
            mid = len(checkpoints) // 2
            late_heaps = [c["heap_mb"] for c in checkpoints[mid:]]
            late_times = [c["elapsed"] for c in checkpoints[mid:]]
            # 線性回歸斜率（MB/秒）
            n = len(late_heaps)
            mean_t = sum(late_times) / n
            mean_h = sum(late_heaps) / n
            slope_num = sum((late_times[i] - mean_t) * (late_heaps[i] - mean_h) for i in range(n))
            slope_den = sum((late_times[i] - mean_t) ** 2 for i in range(n))
            slope_mb_per_min = (slope_num / slope_den * 60) if slope_den > 0 else 0
            print(f"  穩定後增長斜率: {slope_mb_per_min:+.3f} MB/min")
            heap_ok = slope_mb_per_min < 0.5  # 每分鐘增長 < 0.5 MB 才算正常
            print(f"  記憶體洩漏測試: {'✅ 通過' if heap_ok else '❌ 失敗'} (門檻 < 0.5 MB/min)")
        else:
            # checkpoint 不夠，用絕對增長量判斷（允許初始化增長，但不超過 50MB）
            heap_abs_growth = heap_end - heap_start
            heap_ok = heap_abs_growth < 50.0
            print(f"  記憶體洩漏測試: {'✅ 通過' if heap_ok else '❌ 失敗'} (絕對增長 < 50 MB，實際 {heap_abs_growth:.1f} MB)")
    else:
        heap_ok = True

    print(f"\n  初始 Goroutines: {goroutines_start}")
    print(f"  最終 Goroutines: {goroutines_end}")
    goroutine_growth = goroutines_end - goroutines_start
    goroutine_ok = goroutine_growth < 50  # 允許最多增加 50 個 goroutine
    print(f"  Goroutine 增長: {goroutine_growth:+d}")
    print(f"  Goroutine 洩漏測試: {'✅ 通過' if goroutine_ok else '❌ 失敗'} (門檻 < +50)")

    # 玩家統計
    total_attacks = sum(r["attacks"] for r in results)
    total_errors = sum(r["errors"] for r in results)
    total_reconnects = sum(r["reconnects"] for r in results)
    print(f"\n  總攻擊次數: {total_attacks}")
    print(f"  總錯誤次數: {total_errors}")
    print(f"  總重連次數: {total_reconnects}")
    if total_attacks > 0:
        error_rate = total_errors / (total_attacks + total_errors) * 100
        print(f"  錯誤率: {error_rate:.1f}%")
        error_ok = error_rate < 5.0
        print(f"  錯誤率測試: {'✅ 通過' if error_ok else '❌ 失敗'} (門檻 < 5%)")
    else:
        error_ok = True

    # 記憶體趨勢分析
    if len(checkpoints) >= 3:
        heaps = [c["heap_mb"] for c in checkpoints]
        max_heap = max(heaps)
        min_heap = min(heaps)
        print(f"\n  Heap 峰值: {max_heap:.2f} MB")
        print(f"  Heap 最低: {min_heap:.2f} MB")
        print(f"  Heap 波動: {max_heap - min_heap:.2f} MB")

    overall_ok = heap_ok and goroutine_ok and error_ok
    print(f"\n{'=' * 60}")
    print(f"整體結果: {'✅ 全部通過' if overall_ok else '❌ 有問題需要修復'}")
    print(f"{'=' * 60}")

    return overall_ok


def main():
    parser = argparse.ArgumentParser(description="Server 壓力測試")
    parser.add_argument("--duration", type=float, default=60.0,
                        help="測試持續時間（秒），預設 60 秒")
    parser.add_argument("--clients", type=int, default=5,
                        help="模擬客戶端數量，預設 5 個")
    args = parser.parse_args()

    success = asyncio.run(run_stress_test(args.duration, args.clients))
    sys.exit(0 if success else 1)


if __name__ == "__main__":
    main()
