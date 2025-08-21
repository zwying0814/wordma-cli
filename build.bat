@echo off
REM Simple build script for local development
REM Automatically detects version from git tags

setlocal enabledelayedexpansion

REM Get version from git tag
for /f "tokens=*" %%i in ('git describe --tags --exact-match HEAD 2^>nul') do set VERSION=%%i
if "!VERSION!"=="" (
    for /f "tokens=*" %%i in ('git describe --tags --abbrev=0 2^>nul') do set GIT_TAG=%%i
    if "!GIT_TAG!"=="" (
        set VERSION=dev
    ) else (
        set VERSION=!GIT_TAG!-dev
    )
)

REM Get build info
for /f "tokens=*" %%i in ('git rev-parse --short HEAD 2^>nul') do set GIT_COMMIT=%%i
if "%GIT_COMMIT%"=="" set GIT_COMMIT=unknown

for /f "tokens=*" %%i in ('powershell -command "(Get-Date).ToUniversalTime().ToString('yyyy-MM-ddTHH:mm:ssZ')"') do set BUILD_TIME=%%i

set LDFLAGS=-s -w -X main.Version=%VERSION% -X main.BuildTime=%BUILD_TIME% -X main.GitCommit=%GIT_COMMIT%

echo Building wordma CLI...
echo Version: %VERSION%
echo Build Time: %BUILD_TIME%
echo Git Commit: %GIT_COMMIT%
echo.

go build -ldflags="%LDFLAGS%" -o wordma.exe .

if errorlevel 1 (
    echo Build failed!
    exit /b 1
)

echo Build completed successfully!
echo Executable: wordma.exe