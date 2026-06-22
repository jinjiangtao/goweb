@echo off
chcp 65001 >nul
echo ========================================
echo   开发模式启动脚本
echo ========================================
echo.

cd /d "%~dp0backend"
if not exist "data" mkdir data
if not exist "uploads" mkdir uploads
set GOPROXY=https://goproxy.cn,direct
if not exist "go.sum" go mod tidy

cd /d "%~dp0frontend"
if not exist "node_modules" call npm install

echo.
echo 后端服务: http://localhost:8080
echo 前端开发: http://localhost:5173
echo.

start "Backend" cmd /k "cd /d %~dp0backend && set GOPROXY=https://goproxy.cn,direct && go run main.go"
start "Frontend" cmd /k "cd /d %~dp0frontend && npm run dev"

echo 服务已启动，请在浏览器中访问对应地址
pause
