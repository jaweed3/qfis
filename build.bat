@echo off
echo === Building QFIS ===
echo.

echo [1/3] Building frontend...
cd frontend
call npm ci
call npm run build
cd ..
echo Done.

echo [2/3] Copying frontend dist to backend...
if exist backend\frontend\dist rmdir /s /q backend\frontend\dist
xcopy frontend\dist backend\frontend\dist\ /e /i /q >nul
echo Done.

echo [3/3] Cross-compiling for linux/amd64...
cd backend
set CGO_ENABLED=0
set GOOS=linux
set GOARCH=amd64
go build -o qfis-linux .
cd ..
echo Done.

echo.
echo === Build complete! ===
echo Binary: backend\qfis-linux
echo.
echo Deploy:
echo   scp backend\qfis-linux user@server:~/qfis/
