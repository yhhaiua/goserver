@echo off
del /q /s "..\..\protocol\*.go"
packetget.exe
go fmt ..\..\protocol\
pause