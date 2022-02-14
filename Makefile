.PHONY:test
test:
	go test ./... -v

cover:
	go test ./... -coverprofile cover.out && go tool cover -html=cover.out

devcert:
	go run $(GOROOT)/src/crypto/tls/generate_cert.go --host localhost
	@echo "\n----\nRemember to allow the use of the unsigned certificate (from the organization Acme Co) in the browser."
	@echo "\nYou can manually visit localhost:50051 and change the trust settings\n---\n"

protos:
	#Backend
	protoc -I proto --proto_path=./proto  --go_opt=paths=source_relative \
	--go_out=plugins=grpc,paths=source_relative:./torqrpc proto/torq.proto
	# Frontend
	protoc -I proto proto/torq.proto --plugin=./frontend/node_modules/.bin/protoc-gen-ts_proto \
	--ts_proto_opt=esModuleInterop=true,env=browser,forceLong=long,outputClientImpl=grpc-web \
	--ts_proto_out=./frontend/src/torqrpc

# TODO: when front end is transfered


