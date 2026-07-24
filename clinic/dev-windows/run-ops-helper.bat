@echo off
setlocal EnableExtensions
set "ROOT=%~dp0..\.."
set "BACKEND=%ROOT%\backend"
set "OPS_STATIC=%BACKEND%\cmd\ops-helper\static\index.html"
set "OPS_DIST=%ROOT%\ops-helper\dist\index.html"
set "HELPER_URL=http://127.0.0.1:9090"

echo PocoClinic Backup Helper (localhost only)
echo.

where go >nul 2>&1
if errorlevel 1 (
  echo ERROR: Go is not installed or not on PATH.
  echo Install Go from https://go.dev/dl/ and reopen your terminal.
  exit /b 1
)

cd /d "%BACKEND%"

if not exist ".env" (
  if "%DATABASE_URL%"=="" (
    echo ERROR: DATABASE_URL is not set and backend\.env was not found.
    echo.
    echo   copy backend\.env.example backend\.env
    echo   edit DATABASE_URL=./data/pococlinic.db
    echo   run migrate.bat
    echo.
    exit /b 1
  )
) else if "%DATABASE_URL%"=="" (
  echo Config: backend\.env will supply DATABASE_URL and backup paths.
)

if not exist "%OPS_STATIC%" (
  if exist "%OPS_DIST%" (
    echo Syncing ops-helper UI from dist...
    if not exist "%BACKEND%\cmd\ops-helper\static" mkdir "%BACKEND%\cmd\ops-helper\static"
    xcopy /E /I /Y /Q "%ROOT%\ops-helper\dist\*" "%BACKEND%\cmd\ops-helper\static\"
  ) else (
    echo UI not built yet. Running build-ops-helper.bat ...
    call "%ROOT%\build-ops-helper.bat"
    if errorlevel 1 (
      echo ERROR: Backup Helper UI build failed.
      exit /b 1
    )
  )
)

if not exist "%OPS_STATIC%" (
  echo ERROR: Backup Helper UI is still missing after build.
  exit /b 1
)

if not exist "data" mkdir data
if not exist "backups" mkdir backups

echo.
echo Starting helper at %HELPER_URL%
echo Backup folder: .\backups
echo Press Ctrl+C to stop.
echo.

start /B powershell -NoProfile -WindowStyle Hidden -Command ^
  "$url='%HELPER_URL%/api/status'; $deadline=(Get-Date).AddSeconds(45); " ^
  "while((Get-Date) -lt $deadline){try{$r=Invoke-WebRequest -Uri $url -UseBasicParsing -TimeoutSec 2; if($r.StatusCode -eq 200){Start-Process '%HELPER_URL%'; exit 0}}catch{}; Start-Sleep -Milliseconds 400}; exit 1"

go run ./cmd/ops-helper
set "EXIT_CODE=%ERRORLEVEL%"
endlocal & exit /b %EXIT_CODE%
