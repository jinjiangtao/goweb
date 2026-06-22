@echo off
chcp 65001 > nul
echo ========================================
echo    云端相册系统 - 一键启动脚本
echo ========================================
echo.

echo [1/4] 检查后端依赖...
cd backend
if not exist go.sum (
    echo 正在下载后端依赖...
    go mod tidy
)
echo 后端依赖检查完成
cd ..

echo.
echo [2/4] 检查前端依赖...
cd frontend
if not exist node_modules (
    echo 正在安装前端依赖...
    call npm install
)
echo 前端依赖检查完成
cd ..

echo.
echo [3/4] 启动后端服务...
cd backend
start "Xiangce Backend" cmd /k "title 后端服务 - 端口 8080 && go run main.go"
cd ..

echo.
echo [4/4] 启动前端服务...
cd frontend
start "Xiangce Frontend" cmd /k "title 前端服务 - 端口 3000 && npm run dev"
cd ..

echo.
echo ========================================
echo    系统启动中，请稍候...
echo    前端地址: http://localhost:3000
echo    后端地址: http://localhost:8080
echo.
echo    首次启动可能需要几分钟时间
echo ========================================
echo.
pause
