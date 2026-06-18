@echo off
echo ========================================
echo  启动审批系统后端服务
echo ========================================
echo.

cd /d "%~dp0server"

echo [1/3] 检查 Go 环境...
where go >nul 2>&1
if %errorlevel% neq 0 (
    echo 错误: 未找到 Go 环境，请先安装 Go 1.20+
    pause
    exit /b 1
)
echo Go 环境检查通过
echo.

echo [2/3] 下载依赖...
go mod download
if %errorlevel% neq 0 (
    echo 错误: 依赖下载失败
    pause
    exit /b 1
)
echo 依赖下载完成
echo.

echo [3/3] 启动后端服务 (端口: 8080)...
echo.
echo 服务启动后，请不要关闭此窗口
echo 按 Ctrl+C 停止服务
echo.
go run main.go

pause
