@echo off
set "CURR_DRIVE=%~d0"


echo %CURR_DRIVE%
echo %DRIVE_LETTER%


REM  current drive with environment variable DRIVE_LETTER ??
if /I not "%DRIVE_LETTER%"=="%CURR_DRIVE%" (
    echo Running initialization script...
    call init.bat setx
    exit /b 0
)


REM Construct the full paths for the script and application using the drive letter
set "SCRIPT_PATH=%CURR_DRIVE%\System\scripts\ps\link.ps1"
set "APP_PATH=%CURR_DRIVE%\Apps\Allusion\Allusion.exe"

REM Check if the PowerShell script exists
if not exist "%SCRIPT_PATH%" (
    echo PowerShell script not found at "%SCRIPT_PATH%"!
    pause
    exit /b 1
)

REM Check if the application exists
if not exist "%APP_PATH%" (
    echo Application not found at "%APP_PATH%"!
    pause
    exit /b 1
)

REM Run PowerShell script with dynamic drive letter
powershell -ExecutionPolicy Bypass -File "%SCRIPT_PATH%" create Allusion

REM Check if the PowerShell script ran successfully
if %ERRORLEVEL% neq 0 (
    echo There was an error running the PowerShell script.
    pause
    exit /b %ERRORLEVEL%
)

REM Run the application with dynamic drive letter
start "" "%APP_PATH%"

REM Everything works fine; exit the script
exit
