// Package main — Go Server 入口
// server-infra-agent 負責維護
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"chiikawa-game/internal/game"
	"chiikawa-game/internal/ws"
)

const (
	Port    = "7777"
	Version = "3.0.0"
)

var startTime = time.Now()

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("🎮 吉伊卡哇：像素大討伐 Server v%s", Version)
	log.Printf("📡 Port: %s", Port)

	hub := ws.NewHub()
	g := game.NewGame(hub)

	hub.OnConnect = func(clientID string) {
		log.Printf("[WS] Connected: %s (total: %d)", clientID, hub.PlayerCount())
		g.AddPlayer(clientID)
	}
	hub.OnDisconnect = func(clientID string) {
		log.Printf("[WS] Disconnected: %s", clientID)
		g.RemovePlayer(clientID)
	}
	hub.OnMessage = func(clientID string, msgType string, payload json.RawMessage) {
		g.HandleMessage(clientID, msgType, payload)
	}

	g.Start()

	mux := http.NewServeMux()

	// WebSocket 端點
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		playerID := r.URL.Query().Get("player_id")
		if playerID == "" {
			playerID = fmt.Sprintf("player_%d", time.Now().UnixNano()%1000000)
		}
		hub.ServeWS(w, r, playerID)
	})

	// 靜態檔案（HTML5 build）
	mux.Handle("/", http.FileServer(http.Dir("./static")))

	// 健康檢查
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		uptime := int(time.Since(startTime).Seconds())
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":   "ok",
			"version":  Version,
			"uptime_s": uptime,
			"clients":  hub.PlayerCount(),
			"state":    g.GetState(),
		})
	})

	srv := &http.Server{
		Addr:    ":" + Port,
		Handler: mux,
	}

	// 優雅關閉
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		log.Println("Shutting down...")
		os.Exit(0)
	}()

	log.Printf("✅ Server ready at http://localhost:%s", Port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
