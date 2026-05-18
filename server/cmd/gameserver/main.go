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
	"digital-twin/server/internal/ws"
)

const (
	defaultPort = "8080"
	version     = "0.1.0"
)

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

	// 初始化數據埋點（日誌寫入 logs/ 目錄）
	tracker := analytics.Init("room-001", "./logs")
	defer tracker.Close()

	// 建立遊戲實例（單一房間 Prototype）
	g := game.NewGame("room-001", hub)

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
		hub.ServeWS(w, r, clientID)
	})

	// 健康檢查
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok","version":"%s","clients":%d}`,
			version, hub.ClientCount())
	})

	// 統計端點（goroutine 數量、記憶體使用）
	mux.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		var ms runtime.MemStats
		runtime.ReadMemStats(&ms)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"goroutines":%d,"heap_alloc_mb":%.2f,"heap_sys_mb":%.2f,"gc_count":%d,"clients":%d}`,
			runtime.NumGoroutine(),
			float64(ms.HeapAlloc)/1024/1024,
			float64(ms.HeapSys)/1024/1024,
			ms.NumGC,
			hub.ClientCount(),
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

	// 數據分析端點（HTTP GET，供運營查詢）
	mux.HandleFunc("/analytics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		stats := tracker.GetRoomStats()
		if err := json.NewEncoder(w).Encode(stats); err != nil {
			http.Error(w, "encode error", http.StatusInternalServerError)
		}
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
	log.Printf("❤️  Health at http://localhost:%s/health", cfg.Port)
	log.Printf("📊 Stats at http://localhost:%s/stats", cfg.Port)
	log.Printf("🏆 Leaderboard at http://localhost:%s/leaderboard", cfg.Port)
	log.Printf("📈 Analytics at http://localhost:%s/analytics", cfg.Port)

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
