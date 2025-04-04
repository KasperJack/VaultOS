@echo off

if "%~1"=="" (
  echo Usage: %0 [setx^|unsetx]
  goto :eof
)

set "DRIVE=%~d0"

REM  variables 
set "PACKAGE_DIR=%DRIVE%\system\package"
set "SOFTWARE_YAML=%DRIVE%\system\config\software.yaml"
set "JUNCTIONS_JSON=%DRIVE%\system\config\junctions.json"
set "APPS_DIR=%DRIVE%\system\software\apps"
set "GAMES_DIR=%DRIVE%\system\software\games"
set "APPS_SHORTCUTS_DIR=%DRIVE%\apps"
set "GAMES_SHORTCUTS_DIR=%DRIVE%\games"
set "PS_SCRIPTS_DIR=%DRIVE%\system\scripts\ps"
set "LINK_PS=%PS_SCRIPTS_DIR%\link.ps1"

if /I "%~1"=="setx" (
  echo Setting environment variables...
  setx DRIVE_LETTER "%DRIVE%"
  setx PACKAGE_DIR "%PACKAGE_DIR%"
  setx SOFTWARE_YAML "%SOFTWARE_YAML%"
  setx JUNCTIONS_JSON "%JUNCTIONS_JSON%"
  setx APPS_DIR "%APPS_DIR%"
  setx GAMES_DIR "%GAMES_DIR%"
  setx APPS_SHORTCUTS_DIR "%APPS_SHORTCUTS_DIR%"
  setx GAMES_SHORTCUTS_DIR "%GAMES_SHORTCUTS_DIR%"
  setx PS_SCRIPTS_DIR "%PS_SCRIPTS_DIR%"
  setx LINK_PS "%LINK_PS%"
  echo Environment variables have been initialized.
) else if /I "%~1"=="unsetx" (
  echo Removing environment variables...
  reg delete "HKCU\Environment" /f /v DRIVE_LETTER
  reg delete "HKCU\Environment" /f /v PACKAGE_DIR
  reg delete "HKCU\Environment" /f /v SOFTWARE_YAML
  reg delete "HKCU\Environment" /f /v JUNCTIONS_JSON
  reg delete "HKCU\Environment" /f /v APPS_DIR
  reg delete "HKCU\Environment" /f /v GAMES_DIR
  reg delete "HKCU\Environment" /f /v APPS_SHORTCUTS_DIR
  reg delete "HKCU\Environment" /f /v GAMES_SHORTCUTS_DIR
  reg delete "HKCU\Environment" /f /v PS_SCRIPTS_DIR
  reg delete "HKCU\Environment" /f /v LINK_PS
  echo Environment variables have been removed.
) else (
  echo Invalid argument: %1.  use "setx" to set or "unsetx" to remove the variables.
)
