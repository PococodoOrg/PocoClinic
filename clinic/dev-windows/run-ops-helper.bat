@echo off
setlocal
echo PocoClinic Backup Helper (localhost only)
echo.
if "%DATABASE_URL%"=="" (
  echo Tip: set DATABASE_URL for backup/restore. See .env.example
  echo   set DATABASE_URL=./data/pococlinic.db
  echo.
)

if not exist "%~dp0..\..\backend\cmd\ops-helper\static\index.html" (
  echo UI not built yet. Running build-ops-helper.bat first...
  call "%~dp0..\..\build-ops-helper.bat"
)

cd /d "%~dp0..\..\backend"
echo Starting helper at http://127.0.0.1:9090
start http://127.0.0.1:9090
go run ./cmd/ops-helper
endlocal
