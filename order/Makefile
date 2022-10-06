.PHONY: test
test:
	go test -v ./...

.PHONY: run
run:
	go run cmd/termit.go

.PHONY: gen-proto
gen-proto:
	mkdir -p pb && protoc --proto_path=proto proto/*.proto --go_out=:pb --go-grpc_out=:pb
