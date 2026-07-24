@echo off
echo PocoClinic backend
echo.
echo Persistent data: set DATABASE_URL or put it in backend\.env
echo   Example: DATABASE_URL=./data/pococlinic.db
echo   Then: migrate.bat   (or migrations run automatically in development)
echo.
echo Running Backend Tests...
cd /d "%~dp0backend"
go mod download
go test ./... -v
if %ERRORLEVEL% neq 0 (
    echo Tests failed! Please fix the failing tests before running the application.
    pause
    exit /b 1
)

echo All tests passed! Starting Backend Application...
go run ./cmd/main.go
