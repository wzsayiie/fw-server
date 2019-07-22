@echo off

cd /d %~dp0

::build every target

set target_os=windows
set file_name=server.exe
call :build

set target_os=darwin
set file_name=server
call :build

goto :end

::functions

:build
cd cmdserver
set GOOS=%target_os%
set GOARCH=amd64
go build -o ..\BUILD\%GOOS%_%GOARCH%\%file_name%
set GOOS=
set GOARCH=
cd ..

:end
