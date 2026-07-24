@echo off
setlocal EnableExtensions
cd /d "%~dp0"

echo === PocoClinic: run all tests ===
echo.

powershell -NoProfile -ExecutionPolicy Bypass -File "%~dp0scripts\check-no-secrets.ps1"
if errorlevel 1 (
  pause
  exit /b 1
)
echo.

where go >nul 2>&1
if errorlevel 1 (
  echo ERROR: Go is not on PATH. Install Go and reopen the terminal, then retry.
  pause
  exit /b 1
)

where node >nul 2>&1
if errorlevel 1 (
  echo ERROR: Node.js is not on PATH. Install Node.js and reopen the terminal, then retry.
  pause
  exit /b 1
)

echo [1/3] Backend (go test ./...)
pushd backend
go test ./... -count=1
set "BACKEND_EXIT=%ERRORLEVEL%"
popd
if not "%BACKEND_EXIT%"=="0" (
  echo.
  echo Backend tests FAILED.
  pause
  exit /b 1
)
echo Backend tests passed.
echo.

echo [2/3] Frontend (typecheck + lint + test + build)
pushd frontend
call npx tsc --noEmit
if errorlevel 1 goto :frontend_fail
call npm run lint
if errorlevel 1 goto :frontend_fail
call npm test
if errorlevel 1 goto :frontend_fail
call npm run build
if errorlevel 1 goto :frontend_fail
popd
echo Frontend tests passed.
echo.
goto :loadtest

:frontend_fail
popd
echo.
echo Frontend checks FAILED.
pause
exit /b 1

:loadtest
echo [3/3] Loadtest helpers (pytest)
if exist "loadtest\.venv\Scripts\python.exe" (
  loadtest\.venv\Scripts\python.exe -m pytest -q loadtest
) else (
  python -m pytest -q loadtest
)
if errorlevel 1 (
  echo.
  echo Loadtest tests FAILED.
  pause
  exit /b 1
)
echo Loadtest tests passed.
echo.

echo === All tests passed ===
exit /b 0
