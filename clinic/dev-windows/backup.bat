@echo off
setlocal
cd /d "%~dp0..\..\backend"

if "%DATABASE_URL%"=="" (
  echo DATABASE_URL is not set. Example:
  echo   set DATABASE_URL=./data/pococlinic.db
  exit /b 1
)

go run ./cmd/backup
endlocal
