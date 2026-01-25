@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

echo Checking current Go version...
go version

echo.
echo Updating all direct dependencies to stable versions...

REM Get list of direct dependencies
go list -m -f "{{if not .Indirect}}{{.Path}}{{end}}" all > temp_deps.txt

REM Update each dependency
for /f "usebackq tokens=*" %%m in ("temp_deps.txt") do (
    if not "%%m"=="" (
        if not "%%m"=="github.com/bifshteksex/hertzboard" (
            echo   Updating %%m@latest (stable)...
            go get "%%m@latest"
        )
    )
)

REM Clean up temp file
del temp_deps.txt

echo.
echo Updating indirect dependencies (patch level)...
go get -u=patch ./...

echo.
echo Cleaning up and optimizing go.mod...
go mod tidy

echo.
echo All dependencies updated to stable versions!
echo.
echo Checking updated dependencies:
go list -m -u all

echo.
echo Done! Don't forget to test the application after the update.
