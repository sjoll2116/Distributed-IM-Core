#!/bin/bash
# ============================================
# KamaChat 一键部署脚本 (Linux/macOS)
# ============================================
# 用法: chmod +x deploy.sh && ./deploy.sh
# 可选参数:
#   ./deploy.sh --mysql-password=your_pwd --db-name=kama_chat
# ============================================

set -e

# 默认参数
MYSQL_USER="${MYSQL_USER:-root}"
MYSQL_PASSWORD="${MYSQL_PASSWORD:-123456}"
MYSQL_HOST="${MYSQL_HOST:-127.0.0.1}"
MYSQL_PORT="${MYSQL_PORT:-3306}"
DB_NAME="${DB_NAME:-kama_chat}"
REDIS_HOST="${REDIS_HOST:-127.0.0.1}"
REDIS_PORT="${REDIS_PORT:-6379}"
BACKEND_PORT="${BACKEND_PORT:-8000}"
FRONTEND_PORT="${FRONTEND_PORT:-8080}"

# 解析命令行参数
for arg in "$@"; do
  case $arg in
    --mysql-user=*)     MYSQL_USER="${arg#*=}" ;;
    --mysql-password=*) MYSQL_PASSWORD="${arg#*=}" ;;
    --mysql-host=*)     MYSQL_HOST="${arg#*=}" ;;
    --db-name=*)        DB_NAME="${arg#*=}" ;;
    --redis-host=*)     REDIS_HOST="${arg#*=}" ;;
    --backend-port=*)   BACKEND_PORT="${arg#*=}" ;;
    --frontend-port=*)  FRONTEND_PORT="${arg#*=}" ;;
  esac
done

PROJECT_ROOT="$(cd "$(dirname "$0")" && pwd)"
RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; CYAN='\033[0;36m'; NC='\033[0m'

step() { echo -e "\n${CYAN}=> $1${NC}"; }
ok()   { echo -e "  ${GREEN}[OK]${NC} $1"; }
warn() { echo -e "  ${YELLOW}[WARN]${NC} $1"; }
fail() { echo -e "  ${RED}[FAIL]${NC} $1"; exit 1; }

echo -e "\n${CYAN}========================================"
echo "   KamaChat 一键部署"
echo -e "========================================${NC}"

# ============================================
# 1. 环境检查
# ============================================
step "1/8 检查环境依赖..."

command -v go    >/dev/null 2>&1 || fail "Go 未安装 (https://go.dev/dl/)"
command -v node  >/dev/null 2>&1 || fail "Node.js 未安装 (https://nodejs.org/)"
command -v npm   >/dev/null 2>&1 || fail "npm 未安装"
command -v mysql >/dev/null 2>&1 || fail "MySQL Client 未安装"

ok "Go:   $(go version | sed 's/go version //')"
ok "Node: $(node --version)"

# ============================================
# 2. MySQL 数据库创建
# ============================================
step "2/8 创建 MySQL 数据库 '$DB_NAME'..."

MYSQL_PWD="$MYSQL_PASSWORD" mysql -h "$MYSQL_HOST" -P "$MYSQL_PORT" -u "$MYSQL_USER" \
  -e "CREATE DATABASE IF NOT EXISTS \`$DB_NAME\` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;" \
  2>&1 || fail "无法连接 MySQL ($MYSQL_HOST:$MYSQL_PORT)，请确保服务已启动且密码正确"

ok "数据库 '$DB_NAME' 已就绪"

# ============================================
# 3. 检查 Redis
# ============================================
step "3/8 检查 Redis 连接..."

if command -v redis-cli >/dev/null 2>&1; then
    pong=$(redis-cli -h "$REDIS_HOST" -p "$REDIS_PORT" ping 2>/dev/null || echo "FAIL")
    if [ "$pong" = "PONG" ]; then
        ok "Redis 已连接 ($REDIS_HOST:$REDIS_PORT)"
    else
        warn "Redis 无响应，分布式路由可能不可用"
    fi
else
    warn "redis-cli 未找到，跳过检查"
fi

# ============================================
# 4. 生成配置文件
# ============================================
step "4/8 生成 configs/config.toml..."

mkdir -p "$PROJECT_ROOT/configs"
cat > "$PROJECT_ROOT/configs/config.toml" << EOF
[mainConfig]
appName = "KamaChat"
host = "0.0.0.0"
port = $BACKEND_PORT

[mysqlConfig]
host = "$MYSQL_HOST"
port = $MYSQL_PORT
user = "$MYSQL_USER"
password = "$MYSQL_PASSWORD"
databaseName = "$DB_NAME"

[redisConfig]
host = "$REDIS_HOST"
port = $REDIS_PORT
password = ""
db = 0

[authCodeConfig]
accessKeyID = ""
accessKeySecret = ""
signName = ""
templateCode = ""

[logConfig]
logPath = "./logs"

[kafkaConfig]
messageMode = "channel"
hostPort = "127.0.0.1:9092"
loginTopic = "login"
chatTopic = "chat_message"
logoutTopic = "logout"
partition = 0
timeout = 1

[staticSrcConfig]
staticAvatarPath = "./static/avatars"
staticFilePath = "./static/files"
EOF

ok "配置已写入 configs/config.toml"

# ============================================
# 5. 生成 TLS 证书
# ============================================
step "5/8 生成本地 TLS 证书..."

SSL_DIR="$PROJECT_ROOT/pkg/ssl"
CERT_FILE="$SSL_DIR/127.0.0.1+2.pem"
KEY_FILE="$SSL_DIR/127.0.0.1+2-key.pem"

mkdir -p "$SSL_DIR"

if [ -f "$CERT_FILE" ]; then
    ok "证书已存在，跳过"
elif command -v mkcert >/dev/null 2>&1; then
    cd "$SSL_DIR"
    mkcert -install 2>/dev/null
    mkcert -cert-file "127.0.0.1+2.pem" -key-file "127.0.0.1+2-key.pem" 127.0.0.1 localhost ::1
    cd "$PROJECT_ROOT"
    ok "证书已通过 mkcert 生成"
else
    warn "mkcert 未安装，使用 openssl 生成自签名证书..."
    openssl req -x509 -newkey ec -pkeyopt ec_paramgen_curve:P-256 \
        -keyout "$KEY_FILE" -out "$CERT_FILE" -days 365 -nodes \
        -subj "/O=KamaChat Dev" \
        -addext "subjectAltName=IP:127.0.0.1,IP:::1,DNS:localhost" 2>/dev/null
    ok "自签名证书已生成"
fi

# ============================================
# 6. 创建必要目录
# ============================================
step "6/8 创建静态资源目录..."
mkdir -p "$PROJECT_ROOT/static/avatars" "$PROJECT_ROOT/static/files" "$PROJECT_ROOT/logs"
ok "目录已就绪"

# ============================================
# 7. 安装依赖 & 编译后端
# ============================================
step "7/8 安装 Go 依赖并编译后端..."

cd "$PROJECT_ROOT"
go mod tidy >/dev/null 2>&1
go build -o kama_chat_server ./cmd/kama_chat_server || fail "后端编译失败"
ok "后端编译完成: kama_chat_server"

# ============================================
# 8. 安装前端依赖
# ============================================
step "8/8 安装前端依赖..."

FRONTEND_DIR="$PROJECT_ROOT/web/chat-server"
cd "$FRONTEND_DIR"
if [ ! -d "node_modules" ]; then
    npm install >/dev/null 2>&1
    ok "前端依赖安装完成"
else
    ok "node_modules 已存在，跳过"
fi

# ============================================
# 启动服务
# ============================================
echo -e "\n${GREEN}========================================"
echo "   部署完成! 启动服务..."
echo -e "========================================${NC}"
echo ""
echo -e "  后端: https://localhost:$BACKEND_PORT"
echo -e "  前端: http://localhost:$FRONTEND_PORT"
echo ""
echo -e "  ${YELLOW}按 Ctrl+C 停止所有服务${NC}"
echo ""

# 启动后端
cd "$PROJECT_ROOT"
./kama_chat_server &
BACKEND_PID=$!

# 退出时清理
cleanup() {
    echo -e "\n${YELLOW}正在停止服务...${NC}"
    kill $BACKEND_PID 2>/dev/null
    wait $BACKEND_PID 2>/dev/null
    echo -e "${GREEN}所有服务已停止${NC}"
    exit 0
}
trap cleanup SIGINT SIGTERM

# 等待后端启动
sleep 2

# 启动前端（前台）
cd "$FRONTEND_DIR"
npx vue-cli-service serve --port $FRONTEND_PORT
