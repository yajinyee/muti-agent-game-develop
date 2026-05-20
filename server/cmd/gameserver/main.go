// Package main Go Server 入口
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof" // pprof 監控端點（DEBUG 模式）
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/google/uuid"

	"digital-twin/server/internal/analytics"
	"digital-twin/server/internal/config"
	"digital-twin/server/internal/game"
	"digital-twin/server/internal/room"
	"digital-twin/server/internal/store"
	"digital-twin/server/internal/ws"
)

const (
	version = "0.1.0"
)

// serverStartTime 記錄 Server 啟動時間（用於 /health uptime 計算）
var serverStartTime = time.Now()

func main() {
	cfg := config.Load()

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("🎮 吉伊卡哇：像素大討伐 Server v%s", version)
	log.Printf("📡 Starting on port %s", cfg.Port)
	if cfg.DebugMode {
		log.Printf("🔧 Debug mode enabled")
	}

	// 建立 WebSocket Hub
	hub := ws.NewHub()

	// 初始化 Redis Pub/Sub 廣播代理（DAY-061，水平擴展）
	// 有 REDIS_URL 時啟用跨 Server 廣播，無 Redis 時自動降級為純本地廣播
	serverID := fmt.Sprintf("%s-%s", getHostname(), uuid.New().String()[:8])
	pubsubBroker := ws.NewPubSubBroker(cfg.RedisURL, "room-001", serverID, hub)
	if pubsubBroker != nil {
		pubsubBroker.Start()
		log.Printf("📡 Redis Pub/Sub enabled (serverID: %s)", serverID)
	} else {
		log.Printf("📡 Redis Pub/Sub disabled (single-instance mode, serverID: %s)", serverID)
	}

	// 初始化多房間管理器（DAY-019）
	roomMgr := room.NewManager()

	// 初始化玩家狀態 Store（DAY-026）
	// REDIS_URL 為空時自動降級到記憶體模式
	playerStore := store.New(cfg.RedisURL)
	if playerStore.IsRedis() {
		log.Printf("💾 Player state: Redis (%s)", cfg.RedisURL)
	} else {
		log.Printf("💾 Player state: Memory (set REDIS_URL for persistence)")
	}

	// 初始化數據埋點（日誌寫入 logs/ 目錄）
	tracker := analytics.Init("room-001", "./logs")
	defer tracker.Close()

	// 建立遊戲實例（帶 Store 和初始金幣設定，DAY-026）
	g := game.NewGameWithStore("room-001", hub, playerStore, cfg.InitialCoins)

	// 設定 WebSocket 事件處理
	hub.OnConnect = func(clientID string) {
		g.AddPlayer(clientID)
		// 埋點：玩家加入
		tracker.Track(analytics.EventPlayerJoin, clientID, map[string]interface{}{
			"room_id": "room-001",
		})
		hub.Send(clientID, &ws.Message{
			Type: ws.MsgGameState,
			Payload: ws.GameStatePayload{
				State:     g.GetState(),
				Timestamp: time.Now().UnixMilli(),
			},
		})
	}

	hub.OnDisconnect = func(clientID string) {
		g.RemovePlayer(clientID)
		// 埋點：玩家離開
		tracker.Track(analytics.EventPlayerLeave, clientID, map[string]interface{}{
			"room_id": "room-001",
		})
	}

	// 觀戰者斷線時廣播通知給玩家（DAY-055）
	hub.OnSpectatorDisconnect = func(spectatorID string) {
		remaining := hub.SpectatorCount()
		log.Printf("[Spectator] %s disconnected, remaining spectators: %d", spectatorID, remaining)
		hub.BroadcastToPlayers(&ws.Message{
			Type: ws.MsgSpectatorLeave,
			Payload: map[string]interface{}{
				"spectator_id":    spectatorID,
				"spectator_count": remaining,
			},
		})
	}

	hub.OnMessage = func(clientID string, msg *ws.Message) {
		g.HandleMessage(clientID, msg)
	}

	// 啟動遊戲循環
	g.Start()

	// HTTP 路由
	mux := http.NewServeMux()

	// WebSocket 端點
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		clientID := r.URL.Query().Get("player_id")
		if clientID == "" {
			clientID = uuid.New().String()
		}
		// 支援 room_id 參數（多房間架構，DAY-019）
		roomID := r.URL.Query().Get("room_id")
		if roomID == "" {
			roomID = "room-001"
		}
		// 記錄玩家加入的房間（向後相容：room-001 是預設房間）
		if _, ok := roomMgr.GetRoom(roomID); !ok {
			roomID = "room-001"
		}

		// 連線數硬限制（DAY-037）：超過 MaxPlayersPerRoom 時拒絕連線
		// 使用 hub.PlayerCount() 而非 roomMgr，因為目前單房間架構
		if hub.PlayerCount() >= cfg.MaxPlayersPerRoom {
			log.Printf("[WS] Room %s full (%d/%d), rejecting player %s",
				roomID, hub.PlayerCount(), cfg.MaxPlayersPerRoom, clientID)
			http.Error(w, `{"error":"room_full","message":"Room is full, please try again later"}`,
				http.StatusServiceUnavailable)
			return
		}

		log.Printf("[WS] Player %s connecting to room %s (%d/%d)",
			clientID, roomID, hub.PlayerCount()+1, cfg.MaxPlayersPerRoom)
		hub.ServeWS(w, r, clientID)
	})

	// 觀戰 WebSocket 端點（DAY-023）
	// 連線方式：ws://host:port/spectate?room_id=room-001
	// 觀戰者收到所有廣播，但無法發送遊戲指令（attack/bet_change 等）
	mux.HandleFunc("/spectate", func(w http.ResponseWriter, r *http.Request) {
		spectatorID := "spectator-" + uuid.New().String()[:8]
		roomID := r.URL.Query().Get("room_id")
		if roomID == "" {
			roomID = "room-001"
		}
		if _, ok := roomMgr.GetRoom(roomID); !ok {
			roomID = "room-001"
		}
		log.Printf("[WS] Spectator %s connecting to room %s", spectatorID, roomID)
		hub.ServeSpectatorWS(w, r, spectatorID)

		// 觀戰者連線後，非同步傳送遊戲快照（讓觀戰者立即看到當前狀態）
		go func() {
			time.Sleep(100 * time.Millisecond) // 等待連線建立
			snapshot := g.GetSpectatorSnapshot()
			hub.Send(spectatorID, &ws.Message{
				Type:    ws.MsgGameState,
				Payload: snapshot,
			})
			// 廣播所有現有目標給觀戰者
			for _, t := range snapshot.Targets {
				hub.Send(spectatorID, &ws.Message{
					Type:    ws.MsgTargetSpawn,
					Payload: t,
				})
			}
			// 傳送排行榜
			hub.Send(spectatorID, &ws.Message{
				Type:    ws.MsgLeaderboard,
				Payload: snapshot.Leaderboard,
			})
			log.Printf("[Spectator] %s initialized with %d targets", spectatorID, len(snapshot.Targets))

			// 廣播觀戰者加入通知給所有玩家（DAY-054d）
			// 讓玩家知道有人在觀戰，增加社交感
			hub.BroadcastToPlayers(&ws.Message{
				Type: ws.MsgSpectatorJoin,
				Payload: map[string]interface{}{
					"spectator_id":    spectatorID,
					"spectator_count": hub.SpectatorCount(),
				},
			})
		}()
	})

	// 觀戰快照 HTTP 端點（DAY-023）：供前端顯示房間預覽
	mux.HandleFunc("/spectate/snapshot", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		snapshot := g.GetSpectatorSnapshot()
		if err := json.NewEncoder(w).Encode(snapshot); err != nil {
			http.Error(w, "encode error", http.StatusInternalServerError)
		}
	})

	// 健康檢查
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		uptimeSec := int(time.Since(serverStartTime).Seconds())
		uptimeStr := fmt.Sprintf("%dh%dm%ds",
			uptimeSec/3600,
			(uptimeSec%3600)/60,
			uptimeSec%60,
		)
		// 任務重置時間（UTC+8）
		missionResetAt := g.GetMissionResetAt()
		missionResetStr := missionResetAt.Format("2006-01-02T15:04:05Z07:00")
		missionResetSec := int(time.Until(missionResetAt).Seconds())
		if missionResetSec < 0 {
			missionResetSec = 0
		}
		avgPing, _, _ := hub.GetPingStats()

		// Jackpot 池狀態（DAY-054）
		jackpotSnap := g.GetJackpotSnapshot()
		jackpotDaily := g.GetJackpotDailyStats()

		resp := map[string]interface{}{
			"status":               "ok",
			"version":              version,
			"uptime":               uptimeStr,
			"uptime_sec":           uptimeSec,
			"clients":              hub.PlayerCount(),
			"max_players":          cfg.MaxPlayersPerRoom,
			"spectators":           hub.SpectatorCount(),
			"game_state":           g.GetState(),
			"mission_reset_at":     missionResetStr,
			"mission_reset_in_sec": missionResetSec,
			"avg_ping_ms":          avgPing,
			"jackpot": map[string]interface{}{
				"mini":        jackpotSnap["mini"],
				"major":       jackpotSnap["major"],
				"grand":       jackpotSnap["grand"],
				"daily_wins":  jackpotDaily.TotalWins,
				"daily_payout": jackpotDaily.TotalPayout,
			},
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
	})

	// /livez — 存活探針（Kubernetes liveness probe）
	// 只要 Server 程序還活著就回傳 200，不檢查依賴服務
	// 用途：Kubernetes 判斷是否需要重啟 Pod
	mux.HandleFunc("/livez", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"alive","uptime_sec":%d}`,
			int(time.Since(serverStartTime).Seconds()))
	})

	// /readyz — 就緒探針（Kubernetes readiness probe）
	// 檢查 Server 是否準備好接受流量（遊戲循環已啟動 + Store 可用）
	// 用途：Kubernetes 判斷是否將流量路由到此 Pod
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// 就緒條件：Server 啟動超過 2 秒（遊戲循環已初始化）
		uptimeSec := int(time.Since(serverStartTime).Seconds())
		if uptimeSec < 2 {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(w, `{"status":"not_ready","reason":"initializing","uptime_sec":%d}`, uptimeSec)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ready","clients":%d,"game_state":%q,"uptime_sec":%d}`,
			hub.PlayerCount(), g.GetState(), uptimeSec)
	})

	// 統計端點（goroutine 數量、記憶體使用）
	mux.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		var ms runtime.MemStats
		runtime.ReadMemStats(&ms)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"goroutines":%d,"heap_alloc_mb":%.2f,"heap_sys_mb":%.2f,"gc_count":%d,"clients":%d,"spectators":%d}`,
			runtime.NumGoroutine(),
			float64(ms.HeapAlloc)/1024/1024,
			float64(ms.HeapSys)/1024/1024,
			ms.NumGC,
			hub.PlayerCount(),
			hub.SpectatorCount(),
		)
	})

	// 排行榜端點（HTTP GET，供外部查詢）
	mux.HandleFunc("/leaderboard", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		data := g.GetLeaderboardData()
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, "encode error", http.StatusInternalServerError)
		}
	})

	// 房間列表端點（HTTP GET，多房間架構 DAY-019）
	mux.HandleFunc("/rooms", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		rooms := roomMgr.ListRooms()
		if err := json.NewEncoder(w).Encode(rooms); err != nil {
			http.Error(w, "encode error", http.StatusInternalServerError)
		}
	})

	// 數據分析端點（HTTP GET，供運營查詢）
	mux.HandleFunc("/analytics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		stats := tracker.GetRoomStats()
		if err := json.NewEncoder(w).Encode(stats); err != nil {
			http.Error(w, "encode error", http.StatusInternalServerError)
		}
	})

	// Jackpot 狀態端點（DAY-048）
	// GET /jackpot — 回傳三個 Jackpot 池的當前金額和最近中獎記錄
	mux.HandleFunc("/jackpot", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		snap := g.GetJackpotSnapshot()
		history := g.GetJackpotHistory(5)
		daily := g.GetJackpotDailyStats()
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"mini":        snap["mini"],
			"major":       snap["major"],
			"grand":       snap["grand"],
			"history":     history,
			"daily_stats": daily,
			"timestamp":   time.Now().UnixMilli(),
		}); err != nil {
			http.Error(w, "encode error", http.StatusInternalServerError)
		}
	})

	// 週賽排名端點（DAY-066）
	mux.HandleFunc("/tournament", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		snap := g.GetTournamentSnapshot()
		if err := json.NewEncoder(w).Encode(snap); err != nil {
			http.Error(w, "encode error", http.StatusInternalServerError)
		}
	})

	// 公會系統端點（DAY-075）
	// GET /guilds — 取得所有公會列表
	// GET /guild?guild_id=xxx — 取得指定公會詳細資訊
	mux.HandleFunc("/guilds", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		allGuilds := g.Guild.GetAllGuilds()
		type GuildSummary struct {
			GuildID     string `json:"guild_id"`
			Name        string `json:"name"`
			Level       int    `json:"level"`
			MemberCount int    `json:"member_count"`
			OnlineCount int    `json:"online_count"`
			TotalKills  int    `json:"total_kills"`
			TotalCoins  int    `json:"total_coins"`
		}
		summaries := make([]GuildSummary, 0, len(allGuilds))
		for _, gd := range allGuilds {
			onlineCount := 0
			for _, m := range gd.Members {
				if m.IsOnline {
					onlineCount++
				}
			}
			summaries = append(summaries, GuildSummary{
				GuildID:     gd.ID,
				Name:        gd.Name,
				Level:       gd.Level,
				MemberCount: len(gd.Members),
				OnlineCount: onlineCount,
				TotalKills:  gd.TotalKills,
				TotalCoins:  gd.TotalCoins,
			})
		}
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"guilds":    summaries,
			"count":     len(summaries),
			"timestamp": time.Now().UnixMilli(),
		}); err != nil {
			http.Error(w, "encode error", http.StatusInternalServerError)
		}
	})

	// GET /guild-war — 取得公會戰當前排名（DAY-076）
	mux.HandleFunc("/guild-war", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		status, weekID, startAt, endAt := g.GuildWar.GetStatus()
		rankings := g.GuildWar.GetRankings()
		type RankEntry struct {
			Rank        int    `json:"rank"`
			GuildID     string `json:"guild_id"`
			GuildName   string `json:"guild_name"`
			GuildIcon   string `json:"guild_icon"`
			MemberCount int    `json:"member_count"`
			Score       int    `json:"score"`
			KillScore   int    `json:"kill_score"`
			BossScore   int    `json:"boss_score"`
			BonusScore  int    `json:"bonus_score"`
		}
		entries := make([]RankEntry, 0, len(rankings))
		for i, r := range rankings {
			entries = append(entries, RankEntry{
				Rank:        i + 1,
				GuildID:     r.GuildID,
				GuildName:   r.GuildName,
				GuildIcon:   r.GuildIcon,
				MemberCount: r.MemberCount,
				Score:       r.Score,
				KillScore:   r.KillScore,
				BossScore:   r.BossScore,
				BonusScore:  r.BonusScore,
			})
		}
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"week_id":       weekID,
			"status":        string(status),
			"start_at":      startAt.UnixMilli(),
			"end_at":        endAt.UnixMilli(),
			"rankings":      entries,
			"total_guilds":  g.GuildWar.GetParticipatingGuildCount(),
			"timestamp":     time.Now().UnixMilli(),
		}); err != nil {
			http.Error(w, "encode error", http.StatusInternalServerError)
		}
	})

	// GET /daily-boss — 取得每日 BOSS 挑戰狀態（DAY-077）
	mux.HandleFunc("/daily-boss", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		snap := g.DailyBoss.GetSnapshot()
		if snap == nil {
			http.Error(w, `{"error":"no daily boss"}`, http.StatusNotFound)
			return
		}
		topContribs := g.DailyBoss.GetTopContributors(10)
		type ContribEntry struct {
			Rank        int    `json:"rank"`
			PlayerID    string `json:"player_id"`
			DisplayName string `json:"display_name"`
			Damage      int    `json:"damage"`
			Reward      int    `json:"reward"`
		}
		contribs := make([]ContribEntry, 0, len(topContribs))
		for i, c := range topContribs {
			contribs = append(contribs, ContribEntry{
				Rank:        i + 1,
				PlayerID:    c.PlayerID,
				DisplayName: c.DisplayName,
				Damage:      c.Damage,
				Reward:      c.Reward,
			})
		}
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"date_id":        snap.DateID,
			"boss_id":        snap.BossType.ID,
			"boss_name":      snap.BossType.Name,
			"boss_icon":      snap.BossType.Icon,
			"max_hp":         snap.MaxHP,
			"current_hp":     snap.CurrentHP,
			"hp_percent":     g.DailyBoss.GetHPPercent(),
			"status":         string(snap.Status),
			"end_at":         snap.EndAt.UnixMilli(),
			"reward_pool":    snap.RewardPool,
			"total_damage":   snap.TotalDamage,
			"top_contribs":   contribs,
			"difficulty_mod": snap.DifficultyMod,
			"timestamp":      time.Now().UnixMilli(),
		}); err != nil {
			http.Error(w, "encode error", http.StatusInternalServerError)
		}
	})

	// GET /event — 取得當前限時活動狀態（DAY-079）
	mux.HandleFunc("/event", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		snap := g.Event.GetSnapshot()
		if err := json.NewEncoder(w).Encode(snap); err != nil {
			http.Error(w, "encode error", http.StatusInternalServerError)
		}
	})

	// GET /vip — 取得 VIP 等級資訊（DAY-078）
	mux.HandleFunc("/vip", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		playerID := r.URL.Query().Get("player_id")
		if playerID == "" {
			// 回傳 VIP 等級定義列表
			if err := json.NewEncoder(w).Encode(map[string]interface{}{
				"tiers":     g.VIP.GetTiers(),
				"timestamp": time.Now().UnixMilli(),
			}); err != nil {
				http.Error(w, "encode error", http.StatusInternalServerError)
			}
			return
		}
		snap := g.VIP.GetSnapshot(playerID)
		if err := json.NewEncoder(w).Encode(snap); err != nil {
			http.Error(w, "encode error", http.StatusInternalServerError)
		}
	})

	// 玩家個人資料端點（DAY-069）
	// GET /profile?player_id=xxx — 取得指定玩家的個人資料
	// GET /profiles — 取得所有在線玩家的個人資料摘要
	mux.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		playerID := r.URL.Query().Get("player_id")
		if playerID == "" {
			http.Error(w, `{"error":"player_id required"}`, http.StatusBadRequest)
			return
		}
		profile, ok := g.GetPlayerProfile(playerID)
		if !ok {
			http.Error(w, `{"error":"player not found"}`, http.StatusNotFound)
			return
		}
		if err := json.NewEncoder(w).Encode(profile); err != nil {
			http.Error(w, "encode error", http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/profiles", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		profiles := g.GetAllPlayerProfiles()
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"profiles":   profiles,
			"count":      len(profiles),
			"timestamp":  time.Now().UnixMilli(),
		}); err != nil {
			http.Error(w, "encode error", http.StatusInternalServerError)
		}
	})

	// Prometheus metrics 端點（DAY-040）
	// 格式：Prometheus text format（無外部依賴）
	// 用途：Grafana / Prometheus 監控，生產環境可視化
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		var ms runtime.MemStats
		runtime.ReadMemStats(&ms)

		uptimeSec := time.Since(serverStartTime).Seconds()
		analyticsStats := tracker.GetRoomStats()

		w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
		fmt.Fprintf(w, "# HELP chiikawa_uptime_seconds Server uptime in seconds\n")
		fmt.Fprintf(w, "# TYPE chiikawa_uptime_seconds gauge\n")
		fmt.Fprintf(w, "chiikawa_uptime_seconds %.2f\n\n", uptimeSec)

		fmt.Fprintf(w, "# HELP chiikawa_connected_players Current number of connected players\n")
		fmt.Fprintf(w, "# TYPE chiikawa_connected_players gauge\n")
		fmt.Fprintf(w, "chiikawa_connected_players %d\n\n", hub.PlayerCount())

		fmt.Fprintf(w, "# HELP chiikawa_connected_spectators Current number of connected spectators\n")
		fmt.Fprintf(w, "# TYPE chiikawa_connected_spectators gauge\n")
		fmt.Fprintf(w, "chiikawa_connected_spectators %d\n\n", hub.SpectatorCount())

		fmt.Fprintf(w, "# HELP chiikawa_max_players Maximum allowed players per room\n")
		fmt.Fprintf(w, "# TYPE chiikawa_max_players gauge\n")
		fmt.Fprintf(w, "chiikawa_max_players %d\n\n", cfg.MaxPlayersPerRoom)

		fmt.Fprintf(w, "# HELP chiikawa_goroutines Current number of goroutines\n")
		fmt.Fprintf(w, "# TYPE chiikawa_goroutines gauge\n")
		fmt.Fprintf(w, "chiikawa_goroutines %d\n\n", runtime.NumGoroutine())

		fmt.Fprintf(w, "# HELP chiikawa_heap_alloc_bytes Current heap allocation in bytes\n")
		fmt.Fprintf(w, "# TYPE chiikawa_heap_alloc_bytes gauge\n")
		fmt.Fprintf(w, "chiikawa_heap_alloc_bytes %d\n\n", ms.HeapAlloc)

		fmt.Fprintf(w, "# HELP chiikawa_heap_sys_bytes Total heap memory obtained from OS in bytes\n")
		fmt.Fprintf(w, "# TYPE chiikawa_heap_sys_bytes gauge\n")
		fmt.Fprintf(w, "chiikawa_heap_sys_bytes %d\n\n", ms.HeapSys)

		fmt.Fprintf(w, "# HELP chiikawa_gc_total Total number of GC cycles\n")
		fmt.Fprintf(w, "# TYPE chiikawa_gc_total counter\n")
		fmt.Fprintf(w, "chiikawa_gc_total %d\n\n", ms.NumGC)

		fmt.Fprintf(w, "# HELP chiikawa_total_players_joined Total players who joined since server start\n")
		fmt.Fprintf(w, "# TYPE chiikawa_total_players_joined counter\n")
		fmt.Fprintf(w, "chiikawa_total_players_joined %d\n\n", analyticsStats.TotalPlayers)

		fmt.Fprintf(w, "# HELP chiikawa_peak_concurrent_players Peak concurrent players since server start\n")
		fmt.Fprintf(w, "# TYPE chiikawa_peak_concurrent_players gauge\n")
		fmt.Fprintf(w, "chiikawa_peak_concurrent_players %d\n\n", analyticsStats.PeakPlayers)

		fmt.Fprintf(w, "# HELP chiikawa_total_attacks_fired Total attack events since server start\n")
		fmt.Fprintf(w, "# TYPE chiikawa_total_attacks_fired counter\n")
		fmt.Fprintf(w, "chiikawa_total_attacks_fired %d\n\n", analyticsStats.TotalAttacks)

		fmt.Fprintf(w, "# HELP chiikawa_total_kills Total kill events since server start\n")
		fmt.Fprintf(w, "# TYPE chiikawa_total_kills counter\n")
		fmt.Fprintf(w, "chiikawa_total_kills %d\n\n", analyticsStats.TotalKills)

		fmt.Fprintf(w, "# HELP chiikawa_total_coins_rewarded Total coins rewarded since server start\n")
		fmt.Fprintf(w, "# TYPE chiikawa_total_coins_rewarded counter\n")
		fmt.Fprintf(w, "chiikawa_total_coins_rewarded %d\n\n", analyticsStats.TotalReward)

		fmt.Fprintf(w, "# HELP chiikawa_overall_rtp Overall RTP (Return to Player) ratio\n")
		fmt.Fprintf(w, "# TYPE chiikawa_overall_rtp gauge\n")
		fmt.Fprintf(w, "chiikawa_overall_rtp %.4f\n\n", analyticsStats.OverallRTP)

		fmt.Fprintf(w, "# HELP chiikawa_boss_spawns_total Total boss spawn events since server start\n")
		fmt.Fprintf(w, "# TYPE chiikawa_boss_spawns_total counter\n")
		fmt.Fprintf(w, "chiikawa_boss_spawns_total %d\n\n", analyticsStats.BossSpawnCount)

		fmt.Fprintf(w, "# HELP chiikawa_bonus_games_total Total bonus game events since server start\n")
		fmt.Fprintf(w, "# TYPE chiikawa_bonus_games_total counter\n")
		fmt.Fprintf(w, "chiikawa_bonus_games_total %d\n\n", analyticsStats.BonusCount)

		fmt.Fprintf(w, "# HELP chiikawa_ws_messages_received_total Total WebSocket messages received from clients\n")
		fmt.Fprintf(w, "# TYPE chiikawa_ws_messages_received_total counter\n")
		fmt.Fprintf(w, "chiikawa_ws_messages_received_total %d\n\n", hub.MsgReceived.Load())

		fmt.Fprintf(w, "# HELP chiikawa_ws_messages_sent_total Total WebSocket messages sent to clients\n")
		fmt.Fprintf(w, "# TYPE chiikawa_ws_messages_sent_total counter\n")
		fmt.Fprintf(w, "chiikawa_ws_messages_sent_total %d\n\n", hub.MsgSent.Load())

		fmt.Fprintf(w, "# HELP chiikawa_ws_messages_dropped_total Total WebSocket messages dropped (buffer full or rate limited)\n")
		fmt.Fprintf(w, "# TYPE chiikawa_ws_messages_dropped_total counter\n")
		fmt.Fprintf(w, "chiikawa_ws_messages_dropped_total %d\n\n", hub.MsgDropped.Load())

		// WebSocket 壓縮統計（DAY-042）
		// gorilla/websocket 的 permessage-deflate 在 wire 層壓縮，無法直接取得壓縮後大小
		// 用原始位元組數 + 估算壓縮率（JSON 文字通常壓縮到 30-40%）來計算估算值
		bytesSentRaw := hub.BytesSentRaw.Load()
		msgSent := hub.MsgSent.Load()
		var avgMsgSizeBytes float64
		if msgSent > 0 {
			avgMsgSizeBytes = float64(bytesSentRaw) / float64(msgSent)
		}
		// permessage-deflate 對 JSON 的典型壓縮率：35%（壓縮後約為原始的 35%）
		// 估算節省的頻寬 = 原始大小 × (1 - 0.35) = 原始大小 × 0.65
		estimatedBytesSaved := int64(float64(bytesSentRaw) * 0.65)

		fmt.Fprintf(w, "# HELP chiikawa_ws_bytes_sent_raw_total Total raw (pre-compression) bytes sent via WebSocket\n")
		fmt.Fprintf(w, "# TYPE chiikawa_ws_bytes_sent_raw_total counter\n")
		fmt.Fprintf(w, "chiikawa_ws_bytes_sent_raw_total %d\n\n", bytesSentRaw)

		fmt.Fprintf(w, "# HELP chiikawa_ws_avg_message_size_bytes Average WebSocket message size in bytes (pre-compression)\n")
		fmt.Fprintf(w, "# TYPE chiikawa_ws_avg_message_size_bytes gauge\n")
		fmt.Fprintf(w, "chiikawa_ws_avg_message_size_bytes %.2f\n\n", avgMsgSizeBytes)

		fmt.Fprintf(w, "# HELP chiikawa_ws_estimated_bytes_saved_total Estimated bytes saved by permessage-deflate compression (65%% compression ratio assumed)\n")
		fmt.Fprintf(w, "# TYPE chiikawa_ws_estimated_bytes_saved_total counter\n")
		fmt.Fprintf(w, "chiikawa_ws_estimated_bytes_saved_total %d\n\n", estimatedBytesSaved)

		// WebSocket 訊息類型分布（DAY-043）
		// 讓 Grafana 能看到各類型訊息的頻率，識別高頻訊息類型
		msgTypeCounts := hub.GetMsgTypeCounts()
		if len(msgTypeCounts) > 0 {
			fmt.Fprintf(w, "# HELP chiikawa_ws_msg_type_total Total WebSocket messages sent by type\n")
			fmt.Fprintf(w, "# TYPE chiikawa_ws_msg_type_total counter\n")
			for msgType, count := range msgTypeCounts {
				fmt.Fprintf(w, "chiikawa_ws_msg_type_total{type=%q} %d\n", msgType, count)
			}
			fmt.Fprintf(w, "\n")
		}

		// DAY-041：active_targets 指標（讓 Grafana 能監控目標物數量）
		fmt.Fprintf(w, "# HELP chiikawa_active_targets Current number of active targets on screen\n")
		fmt.Fprintf(w, "# TYPE chiikawa_active_targets gauge\n")
		fmt.Fprintf(w, "chiikawa_active_targets %d\n\n", g.GetActiveTargetCount())

		// Ping latency 統計（DAY-044）
		// 追蹤 Server 發送 ping 到收到 pong 的往返延遲
		avgPingMs, maxPingMs, pingCount := hub.GetPingStats()
		fmt.Fprintf(w, "# HELP chiikawa_ws_ping_latency_avg_ms Average WebSocket ping/pong round-trip latency in milliseconds\n")
		fmt.Fprintf(w, "# TYPE chiikawa_ws_ping_latency_avg_ms gauge\n")
		fmt.Fprintf(w, "chiikawa_ws_ping_latency_avg_ms %.2f\n\n", avgPingMs)

		fmt.Fprintf(w, "# HELP chiikawa_ws_ping_latency_max_ms Maximum WebSocket ping/pong round-trip latency in milliseconds\n")
		fmt.Fprintf(w, "# TYPE chiikawa_ws_ping_latency_max_ms gauge\n")
		fmt.Fprintf(w, "chiikawa_ws_ping_latency_max_ms %d\n\n", maxPingMs)

		fmt.Fprintf(w, "# HELP chiikawa_ws_ping_samples_total Total number of ping/pong latency samples collected\n")
		fmt.Fprintf(w, "# TYPE chiikawa_ws_ping_samples_total counter\n")
		fmt.Fprintf(w, "chiikawa_ws_ping_samples_total %d\n\n", pingCount)

		// Per-client ping latency（DAY-044）
		clientLatencies := hub.GetClientPingLatencies()
		if len(clientLatencies) > 0 {
			fmt.Fprintf(w, "# HELP chiikawa_ws_client_ping_ms Latest ping/pong latency per client in milliseconds\n")
			fmt.Fprintf(w, "# TYPE chiikawa_ws_client_ping_ms gauge\n")
			for clientID, latMs := range clientLatencies {
				// 只顯示有效延遲（> 0 表示已收到至少一次 pong）
				if latMs > 0 {
					// 截短 clientID 避免 label 過長
					shortID := clientID
					if len(shortID) > 8 {
						shortID = shortID[:8]
					}
					fmt.Fprintf(w, "chiikawa_ws_client_ping_ms{client=%q} %d\n", shortID, latMs)
				}
			}
			fmt.Fprintf(w, "\n")
		}

		// Client 端效能數據（DAY-045）
		// 由 Client 每 30 秒上報，讓 Grafana 能看到玩家端的效能狀況
		perfSnapshots := hub.GetClientPerfSnapshots()
		if len(perfSnapshots) > 0 {
			fmt.Fprintf(w, "# HELP chiikawa_client_fps Client-side frames per second reported by player\n")
			fmt.Fprintf(w, "# TYPE chiikawa_client_fps gauge\n")
			for _, snap := range perfSnapshots {
				shortID := snap.ClientID
				if len(shortID) > 8 {
					shortID = shortID[:8]
				}
				fmt.Fprintf(w, "chiikawa_client_fps{client=%q,quality=%q} %.1f\n",
					shortID, snap.Quality, snap.FPS)
			}
			fmt.Fprintf(w, "\n")

			fmt.Fprintf(w, "# HELP chiikawa_client_memory_mb Client-side static memory usage in MB\n")
			fmt.Fprintf(w, "# TYPE chiikawa_client_memory_mb gauge\n")
			for _, snap := range perfSnapshots {
				shortID := snap.ClientID
				if len(shortID) > 8 {
					shortID = shortID[:8]
				}
				fmt.Fprintf(w, "chiikawa_client_memory_mb{client=%q} %.1f\n", shortID, snap.MemoryMB)
			}
			fmt.Fprintf(w, "\n")

			fmt.Fprintf(w, "# HELP chiikawa_client_draw_calls Client-side draw calls per frame\n")
			fmt.Fprintf(w, "# TYPE chiikawa_client_draw_calls gauge\n")
			for _, snap := range perfSnapshots {
				shortID := snap.ClientID
				if len(shortID) > 8 {
					shortID = shortID[:8]
				}
				fmt.Fprintf(w, "chiikawa_client_draw_calls{client=%q} %d\n", shortID, snap.DrawCalls)
			}
			fmt.Fprintf(w, "\n")

			// 計算所有 Client 的平均 FPS（供 Grafana 整體趨勢面板）
			var totalFPS float64
			for _, snap := range perfSnapshots {
				totalFPS += snap.FPS
			}
			avgClientFPS := totalFPS / float64(len(perfSnapshots))
			fmt.Fprintf(w, "# HELP chiikawa_client_avg_fps Average FPS across all connected clients\n")
			fmt.Fprintf(w, "# TYPE chiikawa_client_avg_fps gauge\n")
			fmt.Fprintf(w, "chiikawa_client_avg_fps %.1f\n\n", avgClientFPS)
		}

		// Client 端效能歷史統計（DAY-051）
		// 從 ring buffer 取最近 300 秒（5分鐘）的記錄，計算 p5/avg/p95 FPS
		perfHistory := hub.GetPerfHistory(300)
		if len(perfHistory) > 0 {
			var sumFPS float64
			minFPS := perfHistory[0].FPS
			maxFPS := perfHistory[0].FPS
			for _, e := range perfHistory {
				sumFPS += e.FPS
				if e.FPS < minFPS {
					minFPS = e.FPS
				}
				if e.FPS > maxFPS {
					maxFPS = e.FPS
				}
			}
			avgFPS := sumFPS / float64(len(perfHistory))
			fmt.Fprintf(w, "# HELP chiikawa_client_fps_history_avg Average FPS from last 5min history\n")
			fmt.Fprintf(w, "# TYPE chiikawa_client_fps_history_avg gauge\n")
			fmt.Fprintf(w, "chiikawa_client_fps_history_avg %.1f\n", avgFPS)
			fmt.Fprintf(w, "# HELP chiikawa_client_fps_history_min Min FPS from last 5min history\n")
			fmt.Fprintf(w, "# TYPE chiikawa_client_fps_history_min gauge\n")
			fmt.Fprintf(w, "chiikawa_client_fps_history_min %.1f\n", minFPS)
			fmt.Fprintf(w, "# HELP chiikawa_client_fps_history_max Max FPS from last 5min history\n")
			fmt.Fprintf(w, "# TYPE chiikawa_client_fps_history_max gauge\n")
			fmt.Fprintf(w, "chiikawa_client_fps_history_max %.1f\n", maxFPS)
			fmt.Fprintf(w, "# HELP chiikawa_client_fps_history_samples Total perf history samples in last 5min\n")
			fmt.Fprintf(w, "# TYPE chiikawa_client_fps_history_samples gauge\n")
			fmt.Fprintf(w, "chiikawa_client_fps_history_samples %d\n\n", len(perfHistory))
		}

		// Progressive Jackpot 指標（DAY-048）
		// 讓 Grafana 能監控三個 Jackpot 池的累積金額
		jackpotSnap := g.GetJackpotSnapshot()
		fmt.Fprintf(w, "# HELP chiikawa_jackpot_pool Current jackpot pool amount by level\n")
		fmt.Fprintf(w, "# TYPE chiikawa_jackpot_pool gauge\n")
		fmt.Fprintf(w, "chiikawa_jackpot_pool{level=\"mini\"} %d\n", jackpotSnap["mini"])
		fmt.Fprintf(w, "chiikawa_jackpot_pool{level=\"major\"} %d\n", jackpotSnap["major"])
		fmt.Fprintf(w, "chiikawa_jackpot_pool{level=\"grand\"} %d\n\n", jackpotSnap["grand"])

		// Jackpot 每日統計指標（DAY-049b）
		jackpotDaily := g.GetJackpotDailyStats()
		fmt.Fprintf(w, "# HELP chiikawa_jackpot_daily_wins Today's jackpot win count by level\n")
		fmt.Fprintf(w, "# TYPE chiikawa_jackpot_daily_wins counter\n")
		fmt.Fprintf(w, "chiikawa_jackpot_daily_wins{level=\"mini\"} %d\n", jackpotDaily.MiniCount)
		fmt.Fprintf(w, "chiikawa_jackpot_daily_wins{level=\"major\"} %d\n", jackpotDaily.MajorCount)
		fmt.Fprintf(w, "chiikawa_jackpot_daily_wins{level=\"grand\"} %d\n\n", jackpotDaily.GrandCount)
		fmt.Fprintf(w, "# HELP chiikawa_jackpot_daily_payout Today's total jackpot payout\n")
		fmt.Fprintf(w, "# TYPE chiikawa_jackpot_daily_payout counter\n")
		fmt.Fprintf(w, "chiikawa_jackpot_daily_payout %d\n\n", jackpotDaily.TotalPayout)
	})

	// pprof 監控端點（Debug 模式下啟用，用於記憶體/goroutine 分析）
	if cfg.DebugMode {
		mux.HandleFunc("/debug/pprof/", http.DefaultServeMux.ServeHTTP)
		log.Printf("🔍 pprof at http://localhost:%s/debug/pprof/", cfg.Port)
	}

	// 靜態檔案（Godot Web Export 用）
	// 注意：Godot HTML5 需要 SharedArrayBuffer，必須加 COOP/COEP headers
	// 支援 gzip 壓縮（wasm 從 36MB 壓縮到 9MB，減少 75% 下載量）
	staticHandler := http.FileServer(http.Dir("./static"))
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
		w.Header().Set("Cross-Origin-Resource-Policy", "cross-origin")

		// 對大型檔案啟用 gzip 壓縮（wasm、pck、js）
		path := r.URL.Path
		if isCompressiblePath(path) {
			// 只有在瀏覽器支援 gzip 時才提供壓縮版本
			acceptEncoding := r.Header.Get("Accept-Encoding")
			if strings.Contains(acceptEncoding, "gzip") {
				w.Header().Set("Content-Encoding", "gzip")
				// 設定正確的 Content-Type
				switch {
				case strings.HasSuffix(path, ".wasm"):
					w.Header().Set("Content-Type", "application/wasm")
				case strings.HasSuffix(path, ".pck"):
					w.Header().Set("Content-Type", "application/octet-stream")
				case strings.HasSuffix(path, ".js"):
					w.Header().Set("Content-Type", "application/javascript")
				}
				// 嘗試提供預壓縮的 .gz 版本
				gzPath := "./static" + path + ".gz"
				if _, err := os.Stat(gzPath); err == nil {
					http.ServeFile(w, r, gzPath)
					return
				}
			}
		}

		staticHandler.ServeHTTP(w, r)
	}))

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("✅ Server ready at http://localhost:%s", cfg.Port)
	log.Printf("🔌 WebSocket at ws://localhost:%s/ws", cfg.Port)
	log.Printf("👁️  Spectate at ws://localhost:%s/spectate", cfg.Port)
	log.Printf("📸 Snapshot at http://localhost:%s/spectate/snapshot", cfg.Port)
	log.Printf("❤️  Health at http://localhost:%s/health", cfg.Port)
	log.Printf("🟢 Liveness at http://localhost:%s/livez", cfg.Port)
	log.Printf("🔵 Readiness at http://localhost:%s/readyz", cfg.Port)
	log.Printf("📊 Stats at http://localhost:%s/stats", cfg.Port)
	log.Printf("🏆 Leaderboard at http://localhost:%s/leaderboard", cfg.Port)
	log.Printf("📈 Analytics at http://localhost:%s/analytics", cfg.Port)
	log.Printf("🏠 Rooms at http://localhost:%s/rooms", cfg.Port)
	log.Printf("📉 Metrics at http://localhost:%s/metrics", cfg.Port)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-quit
	log.Println("🛑 Shutting down server...")

	// 停止遊戲循環（清理 goroutine）
	g.Stop()

	// 停止 Redis Pub/Sub 代理（DAY-061）
	if pubsubBroker != nil {
		pubsubBroker.Stop()
	}

	// 關閉 Store（確保玩家狀態已儲存）
	if err := playerStore.Close(); err != nil {
		log.Printf("Store close error: %v", err)
	}

	// 給 HTTP 連線 5 秒完成
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}
	log.Println("✅ Server exited cleanly")
}

// isCompressiblePath 判斷路徑是否可以提供 gzip 壓縮版本
// Godot HTML5 的 wasm/pck/js 檔案壓縮效果顯著（wasm 75% 減少）
func isCompressiblePath(path string) bool {
	return strings.HasSuffix(path, ".wasm") ||
		strings.HasSuffix(path, ".pck") ||
		strings.HasSuffix(path, ".js")
}

// getHostname 取得主機名稱（用於 serverID 生成）
// 失敗時回傳 "unknown"
func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}
