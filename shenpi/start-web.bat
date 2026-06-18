@echo off
echo ========================================
echo  启动审批系统前端服务
echo ========================================
echo.

cd /d "%~dp0web"

echo [1/3] 检查 Node.js 环境...
where node >nul 2>&1
if %errorlevel% neq 0 (
    echo 错误: 未找到 Node.js 环境，请先安装 Node.js 16+
    pause
    exit /b 1
)
echo Node.js 环境检查通过
echo.

echo [2/3] 安装依赖...
if not exist "node_modules" (
    call npm install
    if %errorlevel% neq 0 (
        echo 错误: 依赖安装失败
        pause
        exit /b 1
    )
    echo 依赖安装完成
) else (
    echo 依赖已存在，跳过安装
)
echo.

echo [3/3] 启动前端服务 (端口: 5173)...
echo.
echo 服务启动后，请不要关闭此窗口
echo 按 Ctrl+C 停止服务
echo.
call npm run dev

pause
