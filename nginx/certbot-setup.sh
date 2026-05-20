#!/bin/bash
# certbot-setup.sh
# 使用 Let's Encrypt 取得免費 TLS 憑證（生產環境）
#
# 前置需求：
#   - 已有網域名稱（如 game.example.com）
#   - 網域 DNS 已指向此 Server IP
#   - Port 80 已開放（ACME challenge 需要）
#
# 使用方式：
#   chmod +x certbot-setup.sh
#   sudo ./certbot-setup.sh game.example.com

set -e

DOMAIN="${1:-}"
if [ -z "$DOMAIN" ]; then
    echo "用法: $0 <domain>"
    echo "範例: $0 game.example.com"
    exit 1
fi

CERT_DIR="$(dirname "$0")/ssl"
mkdir -p "$CERT_DIR"

echo "為 $DOMAIN 申請 Let's Encrypt 憑證..."

# 使用 certbot standalone 模式（需要暫停 Nginx）
# 或使用 webroot 模式（Nginx 繼續運行）

# 方法 1：Docker certbot（推薦，不需要安裝 certbot）
docker run --rm \
    -v "$(pwd)/nginx/ssl:/etc/letsencrypt/live/$DOMAIN" \
    -v "$(pwd)/nginx/certbot-webroot:/var/www/certbot" \
    certbot/certbot certonly \
    --webroot \
    --webroot-path=/var/www/certbot \
    --email admin@$DOMAIN \
    --agree-tos \
    --no-eff-email \
    -d $DOMAIN

echo "✅ 憑證已取得："
echo "   cert.pem: $CERT_DIR/cert.pem"
echo "   key.pem:  $CERT_DIR/key.pem"
echo ""
echo "憑證有效期 90 天，設定自動更新："
echo "   0 0 1 * * docker run --rm certbot/certbot renew"

# 複製憑證到 nginx/ssl/
cp "/etc/letsencrypt/live/$DOMAIN/fullchain.pem" "$CERT_DIR/cert.pem"
cp "/etc/letsencrypt/live/$DOMAIN/privkey.pem" "$CERT_DIR/key.pem"

echo "✅ 憑證已複製到 nginx/ssl/"
