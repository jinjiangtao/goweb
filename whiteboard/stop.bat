@echo off
chcp 65001 >nul
echo 正在停止所有相关进程...
taskkill /F /IM whiteboard.exe 2>nul
taskkill /F /IM node.exe 2>nul
echo 所有进程已停止
pause
