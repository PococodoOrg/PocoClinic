@echo off
setlocal
cd /d "%~dp0..\..\backend"

if "%DATABASE_URL%"=="" (
  echo DATABASE_URL is not set. Example:
  echo   set DATABASE_URL=./data/pococlinic.db
  exit /b 1
)

echo WARNING: This replaces all database contents with the selected backup.
echo Pass a filename as the first argument, or leave blank for the latest backup.
echo.

set FILE=%~1
if "%FILE%"=="" (
  go run ./cmd/restore --confirm
) else (
  go run ./cmd/restore --confirm --file "%FILE%"
)
endlocal
