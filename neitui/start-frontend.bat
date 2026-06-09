@echo off
cd /d "%~dp0admin"
echo Starting frontend server...
if not exist "node_modules" (
    echo Installing dependencies...
    call npm install
)
call npm run dev
pause
