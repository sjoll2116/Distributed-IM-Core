# ============================================
# KamaChat 一键部署脚本 (Windows PowerShell)
# ============================================
# 用法: 右键以管理员身份运行，或在 PowerShell 中执行:
#   Set-ExecutionPolicy Bypass -Scope Process; .\deploy.ps1
# ============================================

param(
    [string]$MysqlUser     = "root",
    [string]$MysqlPassword = "123456",
    [string]$MysqlHost     = "127.0.0.1",
    [int]   $MysqlPort     = 3306,
    [string]$DbName        = "kama_chat",
    [string]$RedisHost     = "127.0.0.1",
    [int]   $RedisPort     = 6379,
    [int]   $BackendPort   = 8000,
    [int]   $FrontendPort  = 8080
)

$ErrorActionPreference = "Stop"
$ProjectRoot = Split-Path -Parent $MyInvocation.MyCommand.Path

function Write-Step($msg) {
    Write-Host ""
    Write-Host "=> $msg" -ForegroundColor Cyan
}

function Test-Command($cmd) {
    return [bool](Get-Command $cmd -ErrorAction SilentlyContinue)
}

# ============================================
# 1. 环境检查
# ============================================
Write-Host ""
Write-Host "========================================" -ForegroundColor Magenta
Write-Host "   KamaChat 一键部署" -ForegroundColor Magenta
Write-Host "========================================" -ForegroundColor Magenta

Write-Step "1/8 检查环境依赖..."

$missing = @()
if (-not (Test-Command "go"))    { $missing += "Go (https://go.dev/dl/)" }
if (-not (Test-Command "node"))  { $missing += "Node.js (https://nodejs.org/)" }
if (-not (Test-Command "npm"))   { $missing += "npm (随 Node.js 一起安装)" }
if (-not (Test-Command "mysql")) { $missing += "MySQL Client (https://dev.mysql.com/downloads/)" }

if ($missing.Count -gt 0) {
    Write-Host "  [!] 缺少以下工具，请先安装:" -ForegroundColor Red
    $missing | ForEach-Object { Write-Host "      - $_" -ForegroundColor Yellow }
    exit 1
}

$goVer  = (go version) -replace "go version ",""
$nodeVer = node --version
Write-Host "  [OK] Go:   $goVer" -ForegroundColor Green
Write-Host "  [OK] Node: $nodeVer" -ForegroundColor Green

# ============================================
# 2. MySQL 数据库创建
# ============================================
Write-Step "2/8 创建 MySQL 数据库 '$DbName'..."

try {
    $mysqlCmd = "CREATE DATABASE IF NOT EXISTS ``$DbName`` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
    $env:MYSQL_PWD = $MysqlPassword
    echo $mysqlCmd | mysql -h $MysqlHost -P $MysqlPort -u $MysqlUser 2>&1
    if ($LASTEXITCODE -ne 0) { throw "MySQL 命令执行失败" }
    Write-Host "  [OK] 数据库 '$DbName' 已就绪" -ForegroundColor Green
} catch {
    Write-Host "  [!] 无法连接 MySQL ($MysqlHost`:$MysqlPort)" -ForegroundColor Red
    Write-Host "      请确保 MySQL 服务已启动，且用户名密码正确。" -ForegroundColor Yellow
    Write-Host "      错误: $_" -ForegroundColor Yellow
    exit 1
} finally {
    Remove-Item Env:\MYSQL_PWD -ErrorAction SilentlyContinue
}

# ============================================
# 3. 检查 Redis
# ============================================
Write-Step "3/8 检查 Redis 连接..."

if (Test-Command "redis-cli") {
    try {
        $pong = redis-cli -h $RedisHost -p $RedisPort ping 2>&1
        if ($pong -match "PONG") {
            Write-Host "  [OK] Redis 已连接 ($RedisHost`:$RedisPort)" -ForegroundColor Green
        } else {
            Write-Host "  [WARN] Redis 响应异常: $pong" -ForegroundColor Yellow
        }
    } catch {
        Write-Host "  [WARN] 无法 ping Redis，分布式路由功能可能不可用" -ForegroundColor Yellow
    }
} else {
    Write-Host "  [WARN] redis-cli 未找到，跳过检查（确保 Redis 正在运行）" -ForegroundColor Yellow
}

# ============================================
# 4. 生成配置文件
# ============================================
Write-Step "4/8 生成 configs/config.toml..."

$configContent = @"
[mainConfig]
appName = "KamaChat"
host = "0.0.0.0"
port = $BackendPort

[mysqlConfig]
host = "$MysqlHost"
port = $MysqlPort
user = "$MysqlUser"
password = "$MysqlPassword"
databaseName = "$DbName"

[redisConfig]
host = "$RedisHost"
port = $RedisPort
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
"@

$configPath = Join-Path $ProjectRoot "configs\config.toml"
Set-Content -Path $configPath -Value $configContent -Encoding UTF8
Write-Host "  [OK] 配置已写入 $configPath" -ForegroundColor Green

# ============================================
# 5. 生成 TLS 证书
# ============================================
Write-Step "5/8 生成本地 TLS 证书..."

$sslDir = Join-Path $ProjectRoot "pkg\ssl"
$certFile = Join-Path $sslDir "127.0.0.1+2.pem"
$keyFile  = Join-Path $sslDir "127.0.0.1+2-key.pem"

if (-not (Test-Path $sslDir)) { New-Item -ItemType Directory -Path $sslDir -Force | Out-Null }

if (Test-Path $certFile) {
    Write-Host "  [OK] 证书已存在，跳过生成" -ForegroundColor Green
} elseif (Test-Command "mkcert") {
    Push-Location $sslDir
    mkcert -install 2>&1 | Out-Null
    mkcert -cert-file "127.0.0.1+2.pem" -key-file "127.0.0.1+2-key.pem" 127.0.0.1 localhost "::1" 2>&1
    Pop-Location
    Write-Host "  [OK] 证书已通过 mkcert 生成" -ForegroundColor Green
} else {
    Write-Host "  [!] mkcert 未安装，正在生成自签名证书（浏览器会提示不安全）..." -ForegroundColor Yellow
    # 使用 Go 内置工具生成自签名证书
    $certGenScript = @"
package main

import (
    "crypto/ecdsa"
    "crypto/elliptic"
    "crypto/rand"
    "crypto/x509"
    "crypto/x509/pkix"
    "encoding/pem"
    "math/big"
    "net"
    "os"
    "time"
)

func main() {
    priv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
    tmpl := &x509.Certificate{
        SerialNumber: big.NewInt(1),
        Subject:      pkix.Name{Organization: []string{"KamaChat Dev"}},
        NotBefore:    time.Now(),
        NotAfter:     time.Now().Add(365 * 24 * time.Hour),
        KeyUsage:     x509.KeyUsageDigitalSignature,
        ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
        IPAddresses:  []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("::1")},
        DNSNames:     []string{"localhost"},
    }
    certDER, _ := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &priv.PublicKey, priv)
    certOut, _ := os.Create("$($certFile -replace '\\','/')")
    pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certDER})
    certOut.Close()
    keyBytes, _ := x509.MarshalECPrivateKey(priv)
    keyOut, _ := os.Create("$($keyFile -replace '\\','/')")
    pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: keyBytes})
    keyOut.Close()
}
"@
    $tmpCertGen = Join-Path $env:TEMP "kama_certgen.go"
    Set-Content -Path $tmpCertGen -Value $certGenScript -Encoding UTF8
    go run $tmpCertGen
    Remove-Item $tmpCertGen -ErrorAction SilentlyContinue
    Write-Host "  [OK] 自签名证书已生成" -ForegroundColor Green
}

# ============================================
# 6. 创建必要目录
# ============================================
Write-Step "6/8 创建静态资源目录..."

@("static\avatars", "static\files", "logs") | ForEach-Object {
    $dir = Join-Path $ProjectRoot $_
    if (-not (Test-Path $dir)) { New-Item -ItemType Directory -Path $dir -Force | Out-Null }
}
Write-Host "  [OK] 目录已就绪" -ForegroundColor Green

# ============================================
# 7. 安装依赖 & 编译后端
# ============================================
Write-Step "7/8 安装 Go 依赖并编译后端..."

Push-Location $ProjectRoot
go mod tidy 2>&1 | Out-Null
go build -o kama_chat_server.exe ./cmd/kama_chat_server
if ($LASTEXITCODE -ne 0) {
    Write-Host "  [!] 后端编译失败" -ForegroundColor Red
    Pop-Location
    exit 1
}
Write-Host "  [OK] 后端编译完成: kama_chat_server.exe" -ForegroundColor Green
Pop-Location

# ============================================
# 8. 安装前端依赖
# ============================================
Write-Step "8/8 安装前端依赖..."

$frontendDir = Join-Path $ProjectRoot "web\chat-server"
Push-Location $frontendDir
if (-not (Test-Path "node_modules")) {
    npm install 2>&1 | Out-Null
    Write-Host "  [OK] 前端依赖安装完成" -ForegroundColor Green
} else {
    Write-Host "  [OK] node_modules 已存在，跳过安装" -ForegroundColor Green
}
Pop-Location

# ============================================
# 启动服务
# ============================================
Write-Host ""
Write-Host "========================================" -ForegroundColor Green
Write-Host "   部署完成! 启动服务..." -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green
Write-Host ""
Write-Host "  后端: https://localhost:$BackendPort" -ForegroundColor White
Write-Host "  前端: http://localhost:$FrontendPort" -ForegroundColor White
Write-Host ""
Write-Host "  按 Ctrl+C 停止所有服务" -ForegroundColor DarkGray
Write-Host ""

# 启动后端（后台进程）
$backendJob = Start-Job -ScriptBlock {
    param($root)
    Set-Location $root
    & "$root\kama_chat_server.exe"
} -ArgumentList $ProjectRoot

# 等待后端启动
Start-Sleep -Seconds 2

# 启动前端（前台，方便看日志）
Push-Location $frontendDir
try {
    npm run serve
} finally {
    # 清理后端进程
    Write-Host "`n正在停止后端服务..." -ForegroundColor Yellow
    Stop-Job $backendJob -ErrorAction SilentlyContinue
    Remove-Job $backendJob -ErrorAction SilentlyContinue
    Write-Host "所有服务已停止" -ForegroundColor Green
    Pop-Location
}
