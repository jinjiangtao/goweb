@echo off
chcp 65001 > nul
echo ========================================
echo    云端相册系统 - 生产构建脚本
echo ========================================
echo.

echo [1/3] 构建后端...
cd backend
if not exist go.sum (
    go mod tidy
)
go build -o xiangce.exe .
echo 后端构建完成
cd ..

echo.
echo [2/3] 构建前端...
cd frontend
if not exist node_modules (
    call npm install
)
call npm run build
echo 前端构建完成
cd ..

echo.
echo [3/3] 整理部署文件...
if not exist dist (
    mkdir dist
)
xcopy /E /I /Y frontend\dist dist\web
copy backend\xiangce.exe dist\
echo 部署文件已生成到 dist 目录

echo.
echo ========================================
echo    构建完成！
echo    部署目录: dist\
echo    启动方式: xiangce.exe
echo ========================================
echo.
pause
