@echo off
setlocal
echo Building PocoClinic Linux release tarball...
node "%~dp0scripts\build-release.mjs" %*
endlocal
