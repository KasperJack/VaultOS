@echo off
REM Get current drive letter
set CURRENT_DRIVE=%~d0

REM Set environment variables for your portable system
setx PORT_DRIVE "%CURRENT_DRIVE%"
setx PORT_ROOT "%CURRENT_DRIVE%\YourSystemRoot"
setx PORT_BIN "%CURRENT_DRIVE%\YourSystemRoot\bin"
setx PORT_LIB "%CURRENT_DRIVE%\YourSystemRoot\lib"
setx PORT_CONFIG "%CURRENT_DRIVE%\YourSystemRoot\config"
setx PORT_DATA "%CURRENT_DRIVE%\YourSystemRoot\data"

REM Notify user
echo Environment variables set for portable system.
echo Please open a new command prompt window to use your portable system.