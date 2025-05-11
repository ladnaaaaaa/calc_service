@echo off
set "PROJECT_ROOT=%~dp0.."
set "PROTOC=%PROJECT_ROOT%\tools\protoc\bin\protoc.exe"
set "PROTO_FILE=%PROJECT_ROOT%\api\calculator.proto"

echo Installing protoc-gen-go...
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

echo Generating gRPC code...
"%PROTOC%" --proto_path=. --go_out=. --go_opt=paths=source_relative ^
    --go-grpc_out=. --go-grpc_opt=paths=source_relative ^
    "%PROTO_FILE%"

echo Done! 