.PHONY: test
test:
	@go test ./... -coverprofile .cover.txt
	@go tool cover -func .cover.txt
	@rm .cover.txt

.PHONY: build
build:
	rm -rf bin
	mkdir -p bin
	CGO_ENABLED=0 go build -o bin/lark -ldflags "-X lark/pod.appRelease=${release}" main.go

.PHONY: ci-build
ci-build:
	rm -rf bin
	mkdir -p bin
	CGO_ENABLED=0 go build -o bin/lark -ldflags "-X lark/pod.appRelease=${release}" main.go
	cp /usr/local/go/lib/time/zoneinfo.zip bin/zoneinfo.zip

.PHONY: protoc
protoc:
	protoc --go_out=. --go-grpc_out=. proto/*.proto