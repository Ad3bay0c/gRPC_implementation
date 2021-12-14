proto:
	protoc greet/greetpb/greet.proto --go_out=plugins=grpc:.

proto-calculator:
	protoc calculator/calculatorpb/calculator.proto --go_out=plugins=grpc:.

proto-blog:
	protoc blog/blogpb/blog.proto --go_out=plugins=grpc:.

run:
	go run greet/greet_server/server.go

run-client:
	go run greet/greet_client/client.go