set GOOS=windows
set GOARCH=amd64
go build -o build/capuchin_windows_tray_amd64.exe -trimpath -ldflags "-s -w -H=windowsgui" -tags systray
go build -o build/capuchin_windows_amd64.exe -trimpath -ldflags "-s -w"

set GOOS=linux
set GOARCH=amd64
go build -o build/capuchin_linux_amd64 -trimpath -ldflags "-s -w"
