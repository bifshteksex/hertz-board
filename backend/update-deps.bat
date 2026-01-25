@echo off
chcp 65001 >nul

echo Checking current Go version...
go version

echo.
echo Updating all direct dependencies to stable versions...

REM Use go get to update all dependencies
go get -u ./...

echo.
echo Cleaning up and optimizing go.mod...
go mod tidy

echo.
echo All dependencies updated!
echo.
echo Checking updated dependencies:
go list -m -u all

echo.
echo Done! Don't forget to test the application after the update.
