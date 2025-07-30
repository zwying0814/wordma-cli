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
        
        REM 设置临时文件名和最终文件名
        set temp_name=temp-wordma-!GOOS!-!GOARCH!
        set final_name=wordma
        set archive_name=wordma-!GOOS!-!GOARCH!
        
        if "!GOOS!"=="windows" (
            set temp_name=!temp_name!.exe
            set final_name=!final_name!.exe
            set archive_name=!archive_name!.zip
        ) else (
            set archive_name=!archive_name!.tar.gz
        )
        
        echo Building for !GOOS!/!GOARCH!...
        
        set GOOS=!GOOS!
        set GOARCH=!GOARCH!
        go build -ldflags="!LDFLAGS!" -o "dist\!temp_name!" .
        
        if errorlevel 1 (
            echo Failed to build for !GOOS!/!GOARCH!
            exit /b 1
        )
        
        REM 创建压缩包
        pushd dist
        if "!GOOS!"=="windows" (
            REM 重命名并创建zip
            ren "!temp_name!" "!final_name!"
            powershell -command "Compress-Archive -Path '!final_name!' -DestinationPath '!archive_name!' -Force"
            del "!final_name!"
        ) else (
            REM 重命名并创建tar.gz (需要tar命令或使用PowerShell)
            ren "!temp_name!" "!final_name!"
            powershell -command "tar -czf '!archive_name!' '!final_name!'"
            del "!final_name!"
        )
        popd
        
        echo Created !archive_name!
    )
)

echo.
echo Build completed successfully!
echo Archives are available in the dist\ directory:
dir dist\*.zip dist\*.tar.gz 2>nul