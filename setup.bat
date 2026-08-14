@echo off
setlocal
cd /d "%~dp0"
echo.
echo Starting PocoClinic setup wizard...
echo.
node "%~dp0scripts\setup.mjs" %*
set EXITCODE=%ERRORLEVEL%
if %EXITCODE% neq 0 (
  echo.
  echo Setup exited with an error. Is Node.js installed? https://nodejs.org/
  echo.
  exit /b %EXITCODE%
)
endlocal
