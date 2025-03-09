@echo off
echo Running Backend Tests...
cd backend
go mod download
go test ./... -v
if %ERRORLEVEL% neq 0 (
    echo Tests failed! Please fix the failing tests before running the application.
    pause
    exit /b 1
)

echo All tests passed! Starting Backend Application...
go run cmd/main.go 