@echo off
chcp 65001 >nul
echo ========================================
echo     在线协同白板系统 - 编译脚本
echo ========================================
echo.

echo [1/2] 编译后端...
cd /d "%~dp0backend"
go mod tidy

where gcc >nul 2>nul
if %errorlevel% equ 0 (
    echo 检测到 GCC，使用 CGO 编译...
    set CGO_ENABLED=1
    go build -o whiteboard.exe ./cmd/main.go
) else (
    echo 警告: 未检测到 GCC，SQLite 需要 CGO 支持
    echo 请安装 TDM-GCC 或 MinGW-w64 后重新编译
    echo 下载地址: https://jmeubank.github.io/tdm-gcc/
    echo.
    echo 尝试使用纯 Go SQLite 驱动编译...
    echo 正在下载纯 Go SQLite 驱动...
    go get modernc.org/sqlite@v1.27.0
    
    echo 正在替换数据库驱动...
    powershell -Command "(Get-Content internal/models/database.go) -replace '\"gorm.io/driver/sqlite\"', \"gorm.io/driver/sqlite\"`n`t_ \"modernc.org/sqlite\"\" | Set-Content internal/models/database.go"
    powershell -Command "(Get-Content internal/models/database.go) -replace 'gorm.Open\(sqlite.Open\(dbPath\)', 'gorm.Open(sqlite.Dialector{DriverName: \"sqlite\", DSN: dbPath}' | Set-Content internal/models/database.go"
    
    go build -o whiteboard.exe ./cmd/main.go
    
    echo 正在恢复原代码...
    git checkout -- internal/models/database.go 2>nul
)

if %errorlevel% equ 0 (
    echo 后端编译成功！
) else (
    echo 后端编译失败！
    echo 请确保已安装 GCC 并添加到 PATH 环境变量
    pause
    exit /b 1
)

echo.
echo [2/2] 编译前端...
cd /d "%~dp0frontend"
if not exist "node_modules" (
    echo 正在安装依赖...
    npm install
)
npm run build
if %errorlevel% equ 0 (
    echo 前端编译成功！
) else (
    echo 前端编译失败！
    pause
    exit /b 1
)

echo.
echo ========================================
echo  编译完成！
echo  后端可执行文件: backend/whiteboard.exe
echo  前端静态文件: frontend/dist/
echo ========================================
echo.
pause
