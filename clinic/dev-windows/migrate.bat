@echo off
echo Running database migrations...
cd /d "%~dp0..\..\backend"
if "%DATABASE_URL%"=="" (
    if exist ".env" (
        echo DATABASE_URL not set in the shell - migrate will read backend\.env if present.
    ) else (
        echo ERROR: DATABASE_URL is not set and backend\.env was not found.
        echo Example: set DATABASE_URL=./data/pococlinic.db
        echo Or copy backend\.env.example to backend\.env
        exit /b 1
    )
)
go run ./cmd/migrate
exit /b %ERRORLEVEL%
