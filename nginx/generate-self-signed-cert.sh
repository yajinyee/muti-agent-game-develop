#!/bin/bash
# generate-self-signed-cert.sh
# 生成自簽 TLS 憑證（開發/測試用）
# 生產環境請使用 Let's Encrypt（見 certbot-setup.sh）

set -e

CERT_DIR="$(dirname "$0")/ssl"
mkdir -p "$CERT_DIR"

echo "生成自簽 TLS 憑證（開發用）..."

openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
    -keyout "$CERT_DIR/key.pem" \
    -out "$CERT_DIR/cert.pem" \
    -subj "/C=TW/ST=Taiwan/L=Taipei/O=Chiikawa Game/CN=localhost" \
    -addext "subjectAltName=DNS:localhost,IP:127.0.0.1"

echo "✅ 憑證已生成："
echo "   cert.pem: $CERT_DIR/cert.pem"
echo "   key.pem:  $CERT_DIR/key.pem"
echo ""
echo "⚠️  這是自簽憑證，瀏覽器會顯示安全警告。"
echo "   生產環境請使用 Let's Encrypt。"
