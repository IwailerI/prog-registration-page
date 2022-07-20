@echo off
go build -o server.exe main.go 
server run -d
pause