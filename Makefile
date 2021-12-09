proto:
	protoc greet/greetpb/greet.proto --go_out=plugins=grpc:.

run:
	go run greet/greet_server/server.go