build:
	go build -i -v -ldflags "-s -w" -o go-program.exe .

lint:
	go fmt .
	go vet .

proto:
    cd go-rpc/proto/ && protoc --go_out=plugins=grpc:. handler.proto
    cd go-rpc/proto-fast/ && protoc --gofast_out=. handler.proto

.PHONY: proto lint build
