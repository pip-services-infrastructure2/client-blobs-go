.PHONY: all build clean install uninstall fmt simplify check run test

install:
	@go install main.go

run: install
	@go run main.go

test:
	@go clean -testcache && go test -v ./test/...

protogen:
	protoc --go_out=plugins=grpc:. protos/blobs_v1.proto