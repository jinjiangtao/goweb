#!/bin/bash
clear
echo "========================================"
echo "  设备运维工单系统 - 一键启动脚本 (Linux/Mac)"
echo "========================================"
echo ""

echo "[1/4] 检查 Go 环境..."
if ! command -v go &> /dev/null; then
    echo "[错误] 未检测到 Go 环境，请先安装 Go 1.21+"
    exit 1
fi
echo "[OK] Go 版本: $(go version | awk '{print $3}')"

echo ""
echo "[2/4] 检查 Node.js 环境..."
if ! command -v node &> /dev/null; then
    echo "[错误] 未检测到 Node.js 环境，请先安装 Node.js 18+"
    exit 1
fi
echo "[OK] Node.js 版本: $(node -v)"

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
export GOPROXY=https://goproxy.cn,direct

echo ""
echo "[3/4] 安装依赖..."
cd "$SCRIPT_DIR/backend"
if [ ! -f "go.sum" ]; then
    echo "正在下载 Go 依赖包..."
    go mod tidy
    if [ $? -ne 0 ]; then
        echo "[错误] 后端依赖安装失败"
        exit 1
    fi
fi
echo "[OK] 后端依赖已就绪"

cd "$SCRIPT_DIR/frontend"
if [ ! -d "node_modules" ]; then
    echo "正在下载 npm 依赖包（首次安装可能需要几分钟）..."
    npm install
    if [ $? -ne 0 ]; then
        echo "[错误] 前端依赖安装失败"
        exit 1
    fi
fi
if [ ! -d "dist" ]; then
    echo "正在构建前端项目..."
    npm run build
    if [ $? -ne 0 ]; then
        echo "[错误] 前端构建失败"
        exit 1
    fi
fi
echo "[OK] 前端构建完成"

echo ""
echo "========================================"
echo "  启动服务 (端口: 8080)"
echo "========================================"
echo ""
echo "默认账号："
echo "  管理员: admin / 123456"
echo "  普通用户: user1 / 123456"
echo "  运维人员: tech1 / 123456"
echo "  运维人员: tech2 / 123456"
echo ""
echo "服务启动后，请访问: http://localhost:8080"
echo ""
echo "按 Ctrl+C 停止服务"
echo ""

cd "$SCRIPT_DIR/backend"
mkdir -p data uploads
go run main.go
