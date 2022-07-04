clean:
	#rm pto/v1/*
	#rm pto/*
	#rm swagger/*

gv1:
	protoc --proto_path=proto proto/*.proto --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative --go_out=:pto/v1 --go-grpc_out=:pto/v1
	#protoc --proto_path=proto proto/v1/*.proto --go_opt=paths=source_relative --go_out=:pto --go-grpc_out=:pto --grpc-gateway_out=:pto --openapiv2_out=:swagger --validate_out="lang=go:pto"
	#protoc --proto_path=./proto --proto_path=third_party --go_out=internal/pb --go-grpc_opt=paths=source_relative --go-grpc_out=internal/pb --go-grpc_opt=paths=source_relative ./proto/collector.proto

generate:
	#protoc --proto_path=proto proto/*.proto  --go_out=:pto --go-grpc_out=:pto --grpc-gateway_out=:pto --openapiv2_out=:swagger --validate_out="lang=go:pto"

mock:
	mkdir -p ./pto/mocks
	mockgen -package mocks -source=./pto/v1/tgo_grpc.pb.go TgoClient > pto/mocks/tgo_mock.go
	mkdir -p ./store/mocks
	mockgen -package mocks -source=./store/lc.go LC > store/mocks/lc.go
