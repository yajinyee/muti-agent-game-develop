// Package admin 管理後台 Dashboard（DAY-104）
// 提供即時遊戲監控 Web UI，由 Go Server 直接服務
// 訪問：http://localhost:7777/admin
package admin

import (
	_ "embed"
	"net/http"
)

//go:embed dashboard.html
var dashboardHTML []byte

// Handler 回傳 Admin Dashboard HTTP handler
func Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.WriteHeader(http.StatusOK)
		w.Write(dashboardHTML)
	}
}
