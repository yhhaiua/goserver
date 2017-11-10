@echo off
del /q /s "..\..\protocol\*.go"
packetget.exe
go fmt
pause