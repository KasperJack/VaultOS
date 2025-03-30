@echo off
:: Get the current script's drive letter
set DRIVE=%~d0

:: Define environment variables
set PACKAGE_DIR=%DRIVE%\system\package
set SOFTWARE_YAML=%DRIVE%\system\config\software.yaml
set JUNCTIONS_JSON=%DRIVE%\system\config\junctions.json
set APPS_DIR=%DRIVE%\system\software\apps
set GAMES_DIR=%DRIVE%\system\software\games
set APPS_SHORTCUTS_DIR=%DRIVE%\apps
set GAMES_SHORTCUTS_DIR=%DRIVE%\games
set PS_SCRIPTS_DIR=%DRIVE%\system\scripts
set LINK_PS=%PS_SCRIPTS_DIR%\link.ps1

:: Set the environment variables permanently for future processes
setx PACKAGE_DIR "%PACKAGE_DIR%"
setx SOFTWARE_YAML "%SOFTWARE_YAML%"
setx JUNCTIONS_JSON "%JUNCTIONS_JSON%"
setx APPS_DIR "%APPS_DIR%"
setx GAMES_DIR "%GAMES_DIR%"
setx APPS_SHORTCUTS_DIR "%APPS_SHORTCUTS_DIR%"
setx GAMES_SHORTCUTS_DIR "%GAMES_SHORTCUTS_DIR%"
setx PS_SCRIPTS_DIR "%PS_SCRIPTS_DIR%"
setx LINK_PS "%LINK_PS%"

echo Environment variables initialized,
