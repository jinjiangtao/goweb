@echo off
chcp 65001 >nul
echo ========================================
echo  个人笔记系统 - 启动脚本
echo ========================================
echo.

cd /d "%~dp0"

echo [1/2] 启动后端服务...
cd server
if not exist "go.sum" (
    echo 正在下载后端依赖...
    go mod tidy
)
echo 启动后端服务 (首次编译可能需要1-2分钟)...
start "后端服务 - 端口8080" cmd /k "go run main.go"
cd ..

timeout /t 5 /nobreak >nul

echo [2/2] 启动前端服务...
cd web
if not exist "node_modules" (
    echo 正在下载前端依赖...
    npm install
)
start "前端服务 - 端口5173" cmd /k "npm run dev"
cd ..

echo.
echo ========================================
echo  服务启动中...
echo  前端地址: http://localhost:5173
echo  后端地址: http://localhost:8080
echo ========================================
echo.
echo 注意事项:
echo   - 首次启动后端编译SQLite驱动可能较慢
echo   - 请等待后端显示 "Listening and serving HTTP on :8080"
echo   - 然后在浏览器打开 http://localhost:5173
echo ========================================
pause
