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
		fmt.Fprintf(w, `{"status":"ok","version":"%s","uptime":"%s","uptime_sec":%d,"clients":%d,"max_players":%d,"spectators":%d,"game_state":"%s","mission_reset_at":"%s","mission_reset_in_sec":%d}`,
			version,
			uptimeStr,
			uptimeSec,
			hub.PlayerCount(),
			cfg.MaxPlayersPerRoom,
			hub.SpectatorCount(),
			g.GetState(),
			missionResetStr,
			missionResetSec,
		)
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
