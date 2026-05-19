# Dockerfile — 吉伊卡哇：像素大討伐 Game Server
# 多階段建置：builder 編譯，runner 執行（最小化映像大小）

# ── Stage 1: Builder ──────────────────────────────────────────
FROM golang:1.22-alpine AS builder

# 安裝必要工具
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /build

# 先複製 go.mod/go.sum，利用 Docker layer cache
COPY go.mod go.sum ./
RUN go mod download

# 複製原始碼
COPY . .

# 編譯（靜態連結，不依賴 libc）
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o /build/gameserver ./cmd/gameserver

# ── Stage 2: Runner ───────────────────────────────────────────
FROM alpine:3.19

# 安裝 wget（健康檢查用）和 CA 憑證
RUN apk add --no-cache wget ca-certificates tzdata

# 複製編譯好的 binary
COPY --from=builder /build/gameserver /gameserver

# 複製靜態檔案（HTML5 export）
COPY static/ /static/

# 設定工作目錄
WORKDIR /

# 環境變數預設值
ENV PORT=7777
ENV MAX_PLAYERS=10
ENV INITIAL_COINS=10000
ENV DEBUG=false
ENV REDIS_URL=""

# 開放 Port
EXPOSE 7777

# 健康檢查
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -qO- http://localhost:7777/health || exit 1

# 啟動
ENTRYPOINT ["/gameserver"]
