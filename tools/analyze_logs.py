#!/usr/bin/env python3
"""
analyze_logs.py — 解析 analytics JSONL 日誌，輸出統計報告
用法：py tools/analyze_logs.py [--log logs/events-YYYY-MM-DD.jsonl] [--all]
"""

import json
import os
import sys
import argparse
from datetime import datetime
from collections import defaultdict

def load_events(log_path: str) -> list:
    """載入 JSONL 日誌檔案"""
    events = []
    if not os.path.exists(log_path):
        print(f"[ERROR] Log file not found: {log_path}")
        return events
    with open(log_path, 'r', encoding='utf-8') as f:
        for line in f:
            line = line.strip()
            if not line:
                continue
            try:
                events.append(json.loads(line))
            except json.JSONDecodeError:
                pass
    return events

def analyze(events: list) -> dict:
    """分析事件列表，回傳統計報告"""
    stats = {
        "total_events": len(events),
        "time_range": {"start": None, "end": None},
        "players": {
            "total": 0,
            "unique_ids": set(),
            "avg_session_sec": 0,
            "sessions": [],
        },
        "attacks": {
            "total": 0,
            "by_bet_level": defaultdict(int),
            "auto_ratio": 0,
        },
        "kills": {
            "total": 0,
            "by_target": defaultdict(int),
            "by_type": defaultdict(int),
        },
        "rewards": {
            "total_amount": 0,
            "max_single": 0,
            "by_multiplier_bucket": defaultdict(int),
        },
        "boss": {
            "spawn_count": 0,
            "kill_count": 0,
            "kill_rate": 0,
        },
        "bonus": {
            "start_count": 0,
        },
        "rtp": {
            "total_bet": 0,
            "total_reward": 0,
            "overall": 0,
        },
    }

    auto_attacks = 0
    total_attacks = 0
    session_durations = []

    for event in events:
        ts = event.get("ts", 0)
        event_type = event.get("event", "")
        data = event.get("data", {})

        # 時間範圍
        if stats["time_range"]["start"] is None or ts < stats["time_range"]["start"]:
            stats["time_range"]["start"] = ts
        if stats["time_range"]["end"] is None or ts > stats["time_range"]["end"]:
            stats["time_range"]["end"] = ts

        if event_type == "player_join":
            stats["players"]["unique_ids"].add(event.get("player_id", ""))

        elif event_type == "session_summary":
            sess = data.get("stats", {})
            if sess:
                stats["players"]["sessions"].append(sess)
                dur = sess.get("duration_sec", 0)
                if dur > 0:
                    session_durations.append(dur)
                # 累積 RTP 數據
                stats["rtp"]["total_bet"] += sess.get("total_bet", 0)
                stats["rtp"]["total_reward"] += sess.get("total_reward", 0)

        elif event_type == "attack":
            total_attacks += 1
            stats["attacks"]["total"] += 1
            bet_level = data.get("bet_level", 0)
            stats["attacks"]["by_bet_level"][bet_level] += 1
            if data.get("is_auto", False):
                auto_attacks += 1
            # 累積投注
            stats["rtp"]["total_bet"] += data.get("bet_cost", 0)

        elif event_type == "kill":
            stats["kills"]["total"] += 1
            def_id = data.get("def_id", "unknown")
            target_type = data.get("target_type", "normal")
            stats["kills"]["by_target"][def_id] += 1
            stats["kills"]["by_type"][target_type] += 1

        elif event_type == "reward":
            amount = data.get("amount", 0)
            stats["rewards"]["total_amount"] += amount
            stats["rtp"]["total_reward"] += amount
            if amount > stats["rewards"]["max_single"]:
                stats["rewards"]["max_single"] = amount
            # 倍率分桶
            mult = data.get("multiplier", 1.0)
            if mult < 5:
                stats["rewards"]["by_multiplier_bucket"]["1x-5x"] += 1
            elif mult < 20:
                stats["rewards"]["by_multiplier_bucket"]["5x-20x"] += 1
            elif mult < 100:
                stats["rewards"]["by_multiplier_bucket"]["20x-100x"] += 1
            else:
                stats["rewards"]["by_multiplier_bucket"]["100x+"] += 1

        elif event_type == "boss_spawn":
            stats["boss"]["spawn_count"] += 1

        elif event_type == "boss_kill":
            stats["boss"]["kill_count"] += 1

        elif event_type == "bonus_start":
            stats["bonus"]["start_count"] += 1

        elif event_type == "room_summary":
            room = data.get("stats", {})
            if room:
                stats["rtp"]["total_bet"] = max(stats["rtp"]["total_bet"], room.get("total_bet", 0))
                stats["rtp"]["total_reward"] = max(stats["rtp"]["total_reward"], room.get("total_reward", 0))

    # 計算衍生指標
    stats["players"]["total"] = len(stats["players"]["unique_ids"])
    stats["players"]["unique_ids"] = list(stats["players"]["unique_ids"])
    if session_durations:
        stats["players"]["avg_session_sec"] = sum(session_durations) / len(session_durations)

    if total_attacks > 0:
        stats["attacks"]["auto_ratio"] = auto_attacks / total_attacks

    if stats["boss"]["spawn_count"] > 0:
        stats["boss"]["kill_rate"] = stats["boss"]["kill_count"] / stats["boss"]["spawn_count"]

    if stats["rtp"]["total_bet"] > 0:
        stats["rtp"]["overall"] = stats["rtp"]["total_reward"] / stats["rtp"]["total_bet"]

    # 時間範圍格式化
    if stats["time_range"]["start"]:
        stats["time_range"]["start"] = datetime.fromtimestamp(
            stats["time_range"]["start"] / 1000).strftime("%Y-%m-%d %H:%M:%S")
    if stats["time_range"]["end"]:
        stats["time_range"]["end"] = datetime.fromtimestamp(
            stats["time_range"]["end"] / 1000).strftime("%Y-%m-%d %H:%M:%S")

    return stats

def print_report(stats: dict, log_path: str):
    """輸出格式化報告"""
    print("=" * 60)
    print(f"📊 遊戲數據分析報告")
    print(f"   來源：{log_path}")
    print(f"   時間：{stats['time_range']['start']} ~ {stats['time_range']['end']}")
    print(f"   總事件數：{stats['total_events']}")
    print("=" * 60)

    print("\n👥 玩家統計")
    print(f"   唯一玩家數：{stats['players']['total']}")
    avg_sec = stats['players']['avg_session_sec']
    print(f"   平均遊戲時長：{avg_sec:.0f} 秒（{avg_sec/60:.1f} 分鐘）")

    print("\n⚔️  攻擊統計")
    print(f"   總攻擊次數：{stats['attacks']['total']}")
    print(f"   自動模式比例：{stats['attacks']['auto_ratio']*100:.1f}%")
    if stats['attacks']['by_bet_level']:
        print("   投注等級分布：")
        for lv in sorted(stats['attacks']['by_bet_level'].keys()):
            count = stats['attacks']['by_bet_level'][lv]
            pct = count / stats['attacks']['total'] * 100 if stats['attacks']['total'] > 0 else 0
            print(f"     LV{lv}: {count} 次 ({pct:.1f}%)")

    print("\n💀 擊破統計")
    print(f"   總擊破數：{stats['kills']['total']}")
    if stats['kills']['by_type']:
        print("   目標類型分布：")
        for t, count in sorted(stats['kills']['by_type'].items(), key=lambda x: -x[1]):
            print(f"     {t}: {count} 次")
    if stats['kills']['by_target']:
        print("   熱門目標 Top 5：")
        top5 = sorted(stats['kills']['by_target'].items(), key=lambda x: -x[1])[:5]
        for def_id, count in top5:
            print(f"     {def_id}: {count} 次")

    print("\n💰 獎勵統計")
    print(f"   總獎勵金幣：{stats['rewards']['total_amount']:,}")
    print(f"   單次最高獎勵：{stats['rewards']['max_single']:,}")
    if stats['rewards']['by_multiplier_bucket']:
        print("   倍率分布：")
        for bucket, count in stats['rewards']['by_multiplier_bucket'].items():
            print(f"     {bucket}: {count} 次")

    print("\n👹 BOSS 統計")
    print(f"   BOSS 出現次數：{stats['boss']['spawn_count']}")
    print(f"   BOSS 擊敗次數：{stats['boss']['kill_count']}")
    print(f"   BOSS 擊敗率：{stats['boss']['kill_rate']*100:.1f}%")

    print("\n🌿 Bonus 統計")
    print(f"   Bonus 觸發次數：{stats['bonus']['start_count']}")

    print("\n📈 RTP 分析")
    print(f"   總投入金幣：{stats['rtp']['total_bet']:,}")
    print(f"   總獲得金幣：{stats['rtp']['total_reward']:,}")
    rtp_pct = stats['rtp']['overall'] * 100
    rtp_status = "✅" if 85 <= rtp_pct <= 110 else ("⚠️" if rtp_pct < 85 else "🔴")
    print(f"   整體 RTP：{rtp_pct:.2f}% {rtp_status}")
    if rtp_pct < 85:
        print("   ⚠️  RTP 偏低，玩家可能流失")
    elif rtp_pct > 110:
        print("   🔴 RTP 過高，需要調整數值")

    print("\n" + "=" * 60)

def main():
    parser = argparse.ArgumentParser(description="分析遊戲 analytics 日誌")
    parser.add_argument("--log", default=None, help="指定日誌檔案路徑")
    parser.add_argument("--all", action="store_true", help="分析 logs/ 目錄下所有日誌")
    parser.add_argument("--json", action="store_true", help="輸出 JSON 格式（供程式處理）")
    args = parser.parse_args()

    log_dir = "logs"

    if args.all:
        # 分析所有日誌
        if not os.path.exists(log_dir):
            print(f"[ERROR] logs/ directory not found")
            sys.exit(1)
        log_files = [f for f in os.listdir(log_dir) if f.endswith(".jsonl")]
        if not log_files:
            print("[INFO] No log files found in logs/")
            sys.exit(0)
        all_events = []
        for lf in sorted(log_files):
            path = os.path.join(log_dir, lf)
            all_events.extend(load_events(path))
        stats = analyze(all_events)
        if args.json:
            print(json.dumps(stats, indent=2, ensure_ascii=False, default=str))
        else:
            print_report(stats, f"logs/*.jsonl ({len(log_files)} files)")
    else:
        # 分析單一日誌
        if args.log:
            log_path = args.log
        else:
            # 預設用今天的日誌
            today = datetime.now().strftime("%Y-%m-%d")
            log_path = os.path.join(log_dir, f"events-{today}.jsonl")

        events = load_events(log_path)
        if not events:
            print(f"[INFO] No events found in {log_path}")
            print("[INFO] Start the server and play to generate analytics data")
            sys.exit(0)

        stats = analyze(events)
        if args.json:
            print(json.dumps(stats, indent=2, ensure_ascii=False, default=str))
        else:
            print_report(stats, log_path)

if __name__ == "__main__":
    main()
