@echo off
setlocal EnableExtensions
echo Building PocoClinic Backup Helper for Raspberry Pi touchscreen...

where node >nul 2>&1
if errorlevel 1 (
  echo ERROR: Node.js is not installed or not on PATH.
  exit /b 1
)

cd /d "%~dp0ops-helper"
call npm install
if errorlevel 1 exit /b 1

call npm run build:pi
if errorlevel 1 exit /b 1

if not exist "dist\index.html" (
  echo ERROR: build completed but dist\index.html is missing.
  exit /b 1
)

echo Copying Pi touch UI to backend\cmd\ops-helper\static ...
if exist "..\backend\cmd\ops-helper\static" rmdir /S /Q "..\backend\cmd\ops-helper\static"
mkdir "..\backend\cmd\ops-helper\static"
xcopy /E /I /Y /Q dist\* ..\backend\cmd\ops-helper\static\
echo.
echo Done. On the Pi, run: run-ops-helper-pi.bat
echo Or use scripts/pi/start-backup-kiosk.sh for full-screen kiosk mode.
endlocal
