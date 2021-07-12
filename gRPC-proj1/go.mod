//module name of mod file must be the project folder name
//go mod project_folder_name
//go mod tidy, to automatically import all dependecies

module gRPC-proj1

go 1.16

require (
	github.com/golang/protobuf v1.5.2
	google.golang.org/grpc v1.39.0
	google.golang.org/protobuf v1.27.1 // indirect

)
