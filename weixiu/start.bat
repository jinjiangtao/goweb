@echo off
chcp 65001 >nul
echo ========================================
echo   设备运维工单系统 - 一键启动脚本 (Windows)
echo ========================================
echo.

echo [1/4] 检查 Go 环境...
where go >nul 2>&1
if %errorlevel% neq 0 (
    echo [错误] 未检测到 Go 环境，请先安装 Go 1.21+
    pause
    exit /b 1
)
for /f "tokens=3" %%i in ('go version') do set GO_VER=%%i
echo [OK] Go 版本: %GO_VER%

echo.
echo [2/4] 检查 Node.js 环境...
where node >nul 2>&1
if %errorlevel% neq 0 (
    echo [错误] 未检测到 Node.js 环境，请先安装 Node.js 18+
    pause
    exit /b 1
)
for /f "tokens=1" %%i in ('node -v') do set NODE_VER=%%i
echo [OK] Node.js 版本: %NODE_VER%

echo.
echo [3/4] 安装依赖...
cd /d "%~dp0backend"
set GOPROXY=https://goproxy.cn,direct
if not exist "go.sum" (
    echo 正在下载 Go 依赖包，请稍候...
    go mod tidy
    if %errorlevel% neq 0 (
        echo [错误] 后端依赖安装失败
        pause
        exit /b 1
    )
)
echo [OK] 后端依赖已就绪

cd /d "%~dp0frontend"
if not exist "node_modules" (
    echo 正在下载 npm 依赖包，请稍候（首次安装可能需要几分钟）...
    call npm install
    if %errorlevel% neq 0 (
        echo [错误] 前端依赖安装失败
        pause
        exit /b 1
    )
)
if not exist "dist" (
    echo 正在构建前端项目...
    call npm run build
    if %errorlevel% neq 0 (
        echo [错误] 前端构建失败
        pause
        exit /b 1
    )
)
echo [OK] 前端构建完成

echo.
echo ========================================
echo   启动服务 (端口: 8080)
echo ========================================
echo.
echo 默认账号：
echo   管理员: admin / 123456
echo   普通用户: user1 / 123456
echo   运维人员: tech1 / 123456
echo   运维人员: tech2 / 123456
echo.
echo 服务启动后，请访问: http://localhost:8080
echo.
echo 按 Ctrl+C 停止服务
echo.

cd /d "%~dp0backend"
if not exist "data" mkdir data
if not exist "uploads" mkdir uploads
set GOPROXY=https://goproxy.cn,direct
go run main.go

pause
