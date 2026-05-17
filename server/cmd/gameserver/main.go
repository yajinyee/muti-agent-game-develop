// Package main Go Server 入口
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

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

	// 建立遊戲實例（單一房間 Prototype）
	g := game.NewGame("room-001", hub)

	// 設定 WebSocket 事件處理
	hub.OnConnect = func(clientID string) {
		g.AddPlayer(clientID)
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

	// 靜態檔案（Godot Web Export 用）
	// 注意：Godot HTML5 需要 SharedArrayBuffer，必須加 COOP/COEP headers
	staticHandler := http.FileServer(http.Dir("./static"))
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
		w.Header().Set("Cross-Origin-Resource-Policy", "cross-origin")
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

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
