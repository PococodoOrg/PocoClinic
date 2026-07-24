@echo off
setlocal EnableExtensions
set "ROOT=%~dp0..\.."
set "BACKEND=%ROOT%\backend"
set "OPS_STATIC=%BACKEND%\cmd\ops-helper\static\index.html"
set "OPS_DIST=%ROOT%\ops-helper\dist\index.html"
set "HELPER_URL=http://127.0.0.1:9090/?pi=1"

echo PocoClinic Backup Helper - Pi touchscreen mode
echo Opens at %HELPER_URL% (localhost only)
echo.

where go >nul 2>&1
if errorlevel 1 (
  echo ERROR: Go is not installed or not on PATH.
  exit /b 1
)

cd /d "%BACKEND%"

if not exist ".env" (
  if "%DATABASE_URL%"=="" (
    echo ERROR: DATABASE_URL is not set and backend\.env was not found.
    echo Copy backend\.env.example to backend\.env first.
    exit /b 1
  )
) else if "%DATABASE_URL%"=="" (
  echo Config: backend\.env will supply DATABASE_URL and backup paths.
)

if not exist "%OPS_STATIC%" (
  echo Pi UI not built yet. Running build-ops-helper-pi.bat ...
  call "%ROOT%\build-ops-helper-pi.bat"
  if errorlevel 1 exit /b 1
)

if not exist "data" mkdir data
if not exist "backups" mkdir backups

echo.
echo Starting helper at http://127.0.0.1:9090
echo Press Ctrl+C to stop.
echo.

start /B powershell -NoProfile -WindowStyle Hidden -Command ^
  "$url='http://127.0.0.1:9090/api/status'; $deadline=(Get-Date).AddSeconds(45); " ^
  "while((Get-Date) -lt $deadline){try{$r=Invoke-WebRequest -Uri $url -UseBasicParsing -TimeoutSec 2; if($r.StatusCode -eq 200){Start-Process '%HELPER_URL%'; exit 0}}catch{}; Start-Sleep -Milliseconds 400}; exit 1"

go run ./cmd/ops-helper
set "EXIT_CODE=%ERRORLEVEL%"
endlocal & exit /b %EXIT_CODE%
