@echo off
setlocal EnableExtensions
echo Building PocoClinic Backup Helper UI...

where node >nul 2>&1
if errorlevel 1 (
  echo ERROR: Node.js is not installed or not on PATH.
  echo Install Node.js LTS from https://nodejs.org/ and reopen your terminal.
  exit /b 1
)

cd /d "%~dp0ops-helper"
if not exist "package.json" (
  echo ERROR: ops-helper folder not found.
  exit /b 1
)

call npm install
if errorlevel 1 exit /b 1

call npm run build
if errorlevel 1 exit /b 1

if not exist "dist\index.html" (
  echo ERROR: build completed but dist\index.html is missing.
  exit /b 1
)

echo Copying build to backend\cmd\ops-helper\static ...
if exist "..\backend\cmd\ops-helper\static" rmdir /S /Q "..\backend\cmd\ops-helper\static"
mkdir "..\backend\cmd\ops-helper\static"
xcopy /E /I /Y /Q dist\* ..\backend\cmd\ops-helper\static\
echo Done. Start with run-ops-helper.bat
endlocal
