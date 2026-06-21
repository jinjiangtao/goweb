@echo off
chcp 65001 >nul
echo ========================================
echo     在线协同白板系统 - 启动脚本
echo ========================================
echo.

echo [1/2] 启动后端服务...
cd /d "%~dp0backend"
if not exist "whiteboard.exe" (
    echo 正在编译后端...
    go build -o whiteboard.exe ./cmd/main.go
)
start "Whiteboard Backend" cmd /k whiteboard.exe

echo.
echo [2/2] 启动前端服务...
cd /d "%~dp0frontend"
start "Whiteboard Frontend" cmd /k npm run dev

echo.
echo ========================================
echo  服务启动完成！
echo  前端地址: http://localhost:3000
echo  后端地址: http://localhost:8080
echo  WebSocket: ws://localhost:8080/ws/:id
echo ========================================
echo.
echo 按任意键退出...
pause >nul
