@echo off
echo ========================================
echo  在线拖拽式流程审批系统 - 一键启动
echo ========================================
echo.

echo 正在启动后端服务...
start "审批系统-后端" cmd /k "%~dp0start-server.bat"

echo 等待后端服务启动...
timeout /t 5 /nobreak >nul

echo.
echo 正在启动前端服务...
start "审批系统-前端" cmd /k "%~dp0start-web.bat"

echo.
echo ========================================
echo  启动完成！
echo ========================================
echo.
echo 后端服务地址: http://localhost:8080
echo 前端服务地址: http://localhost:5173
echo.
echo 默认测试账号:
echo   admin / 123456 (管理员)
echo   zhangsan / 123456
echo   lisi / 123456
echo   wangwu / 123456
echo   zhaoliu / 123456
echo.
echo 请在浏览器中打开 http://localhost:5173
echo.
pause
