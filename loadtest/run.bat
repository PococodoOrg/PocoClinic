@echo off
REM Thin wrapper — all logic lives in run.py
cd /d "%~dp0"
if not exist ".venv\Scripts\python.exe" (
  python -m venv .venv 2>nul
)
.venv\Scripts\python.exe run.py %*
exit /b %ERRORLEVEL%
