@echo off
setlocal
cd /d "%~dp0..\..\binder-printer"
if not exist node_modules (
  echo Installing binder-printer dependencies...
  call npm install
  if errorlevel 1 exit /b 1
)
echo Starting binder printer at http://localhost:5199
call npm run dev
