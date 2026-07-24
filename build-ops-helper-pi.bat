@echo off
setlocal
echo Building PocoClinic Backup Helper for Raspberry Pi touchscreen...
cd /d "%~dp0ops-helper"
call npm install
if errorlevel 1 exit /b 1
call npm run build:pi
if errorlevel 1 exit /b 1

echo Copying Pi touch UI to ops-helper server static folder...
if exist "..\backend\cmd\ops-helper\static" rmdir /S /Q "..\backend\cmd\ops-helper\static"
mkdir "..\backend\cmd\ops-helper\static"
xcopy /E /I /Y dist\* ..\backend\cmd\ops-helper\static\
echo.
echo Done. On the Pi, run: run-ops-helper-pi.bat
echo Or use scripts/pi/start-backup-kiosk.sh for full-screen kiosk mode.
endlocal
