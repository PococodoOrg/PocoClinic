@echo off
echo Starting PocoClinic Application...

:: Start the backend in a new window
start cmd /k "run-backend.bat"

:: Wait a moment for the backend to start
timeout /t 2 /nobreak

:: Start the frontend in a new window
start cmd /k "run-frontend.bat" 