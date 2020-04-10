@echo off &TITLE Generation Protobuf Code For Go
echo.
mode con cols=100 lines=30
color 0D
@rem cls

setlocal
echo Copyright 2020 micro-community authors.
echo Generate the Go code for .proto files
echo @rem choco install -y protoc
echo @rem  go get -v github.com/micro/protoc-gen-micro
echo @rem  go get -v github.com/golang/protobuf/proto
echo @rem  go get -v github.com/golang/protobuf/protoc-gen-go
echo @rem  visit https://github.com/micro/protoc-gen-micro for this tools
echo.
@rem enter this directory of bat
cd /d %~dp0

echo.
echo ## Current Dir: %cd%
echo.

@rem compile *.proto to go code

protoc -I. --micro_out=. --go_out=. protocol_contract.proto


echo.
echo..........work had been done.
echo.
echo..........code had been generated to :  %cd%
echo.
echo..........press any key to exit
pause >nul