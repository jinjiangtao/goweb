@echo off
echo ========================================
echo  在线简历编辑系统 - 启动脚本
echo ========================================
echo.

echo [1/3] 启动后端服务 (Go Gin)...
cd /d "%~dp0server"
start "Backend Server" cmd /k "jianli-server.exe"

timeout /t 3 /nobreak >nul

echo [2/3] 启动前端服务 (Vue3)...
cd /d "%~dp0web"
start "Frontend Server" cmd /k "npm run dev"

timeout /t 5 /nobreak >nul

echo [3/3] 服务启动完成！
echo.
echo 后端服务: http://localhost:8080
echo 前端服务: http://localhost:5173
echo.
echo 请在浏览器中访问 http://localhost:5173
echo.
pause
