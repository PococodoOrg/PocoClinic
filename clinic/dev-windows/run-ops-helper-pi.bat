@echo off
setlocal
echo PocoClinic Backup Helper - Pi touchscreen mode
echo Opens at http://127.0.0.1:9090 (localhost only)
echo.

if "%DATABASE_URL%"=="" (
  echo Tip: set DATABASE_URL before starting. See .env.example
  echo.
)

if not exist "%~dp0..\..\backend\cmd\ops-helper\static\index.html" (
  echo Pi UI not built yet. Running build-ops-helper-pi.bat ...
  call "%~dp0..\..\build-ops-helper-pi.bat"
)

cd /d "%~dp0..\..\backend"
go run ./cmd/ops-helper
endlocal
