#!/bin/sh

mkdir build
mkdir release

echo "正在构建 Windows x32"
export GOOS=windows
export GOARCH=386
go build -o build/FiveM_Host_Repair_windows_i686.exe main.go

echo "正在构建 Windows x64"
export GOOS=windows
export GOARCH=amd64
go build -o build/FiveM_Host_Repair_windows_amd64.exe main.go

echo "正在构建 Linux x32"
export GOOS=linux
export GOARCH=386
go build -o build/FiveM_Host_Repair_linux_i686 main.go

echo "正在构建 Linux x64"
export GOOS=linux
export GOARCH=amd64
go build -o build/FiveM_Host_Repair_linux_amd64 main.go

echo "构建完成，正在打包 Windows 版本..."
zip release/FiveM_Host_Repair_windows.zip build/FiveM_Host_Repair_windows_amd64.exe build/FiveM_Host_Repair_windows_i686.exe

echo "构建完成，正在打包 Linux 版本..."
zip release/FiveM_Host_Repair_linux.zip build/FiveM_Host_Repair_linux_amd64 build/FiveM_Host_Repair_linux_i686

echo "打包完成，正在清理目录..."

rm -rf build/

echo "全部构建已完成，构建的文件已保存至 release/ 目录下"
