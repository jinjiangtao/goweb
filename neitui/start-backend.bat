@echo off
cd /d "%~dp0server"
echo Starting backend server...
go mod tidy
go run main.go
pause
