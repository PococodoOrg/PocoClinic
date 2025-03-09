@echo off
cd frontend
call npm install

echo Running tests...
call npm test run
if errorlevel 1 (
    echo Tests failed! Please fix the failing tests before running the application.
    pause
    exit /b 1
)

echo Tests passed! Starting the application...
call npm run dev 