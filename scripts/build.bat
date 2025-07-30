@echo off
REM Build script for Windows
REM Usage: scripts\build.bat [version]

setlocal enabledelayedexpansion

set VERSION=%1
if "%VERSION%"=="" set VERSION=dev

for /f "tokens=*" %%i in ('git rev-parse --short HEAD 2^>nul') do set GIT_COMMIT=%%i
if "%GIT_COMMIT%"=="" set GIT_COMMIT=unknown

for /f "tokens=*" %%i in ('powershell -command "Get-Date -Format 'yyyy-MM-ddTHH:mm:ssZ' -AsUTC"') do set BUILD_TIME=%%i

set LDFLAGS=-s -w -X main.Version=%VERSION% -X main.BuildTime=%BUILD_TIME% -X main.GitCommit=%GIT_COMMIT%

echo Building wordma CLI...
echo Version: %VERSION%
echo Build Time: %BUILD_TIME%
echo Git Commit: %GIT_COMMIT%
echo.

REM Create dist directory
if not exist dist mkdir dist

REM Build for multiple platforms
set platforms=linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

for %%p in (%platforms%) do (
    for /f "tokens=1,2 delims=/" %%a in ("%%p") do (
        set GOOS=%%a
        set GOARCH=%%b
        
        set output_name=wordma-!GOOS!-!GOARCH!
        if "!GOOS!"=="windows" set output_name=!output_name!.exe
        
        echo Building for !GOOS!/!GOARCH!...
        
        set GOOS=!GOOS!
        set GOARCH=!GOARCH!
        go build -ldflags="!LDFLAGS!" -o "dist\!output_name!" .
        
        if errorlevel 1 (
            echo Failed to build for !GOOS!/!GOARCH!
            exit /b 1
        )
    )
)

echo.
echo Build completed successfully!
echo Binaries are available in the dist\ directory:
dir dist\