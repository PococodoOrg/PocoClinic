@echo off
setlocal
echo Building PocoClinic Backup Helper UI...
cd /d "%~dp0ops-helper"
call npm install
if errorlevel 1 exit /b 1
call npm run build
if errorlevel 1 exit /b 1

echo Copying build to ops-helper server static folder...
if exist "..\backend\cmd\ops-helper\static" rmdir /S /Q "..\backend\cmd\ops-helper\static"
mkdir "..\backend\cmd\ops-helper\static"
xcopy /E /I /Y dist\* ..\backend\cmd\ops-helper\static\
echo Done. Start with run-ops-helper.bat
endlocal
