set GOOS=windows
set GOARCH=amd64
go build -o build/capuchin_windows_amd64.exe -trimpath -ldflags "-s -w"
