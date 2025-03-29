 @echo off
REM  drive letter 
set DRIVE_LETTER=%~d0

REM Construct the full paths for the script and app using the drive letter
set SCRIPT_PATH=%DRIVE_LETTER%\System\scripts\ps\link.ps1
set APP_PATH=%DRIVE_LETTER%\Apps\software\software.exe



REM if the script file exists
if not exist "%SCRIPT_PATH%" (
    echo PowerShell script not found at "%SCRIPT_PATH%"!
    pause
    exit /b 1
)

REM if the application exists
if not exist "%APP_PATH%" (
    echo Application not found at "%APP_PATH%"!
    pause
    exit /b 1
)

REM Run powerShell script 
powershell -ExecutionPolicy Bypass -File "%SCRIPT_PATH%" create Allusion

REM  powerShell script ran successfully non 0 exit code 
if %ERRORLEVEL% neq 0 (
    echo There was an error running the PowerShell script.
    pause
    exit /b %ERRORLEVEL%
)

REM run the software
start "" "%APP_PATH%"

REM everything works fine
exit