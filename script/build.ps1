go build -ldflags="-s -w -H windowsgui" -o ..\bin\main.exe github.com/itning/DouBanReptile
cd ../bin
upx main.exe
echo SUCCESS
pause