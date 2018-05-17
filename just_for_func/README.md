Commands required to run this project:

// after making changes to the todo.proto, 
- generate the new todo.pb.go by cd into todo and run the following command
	protoc -I . todo.proto --go_out=plugins=grpc:.

// to run the server
- build the server using the command 
	"go build"
- move the generated "server" executable binary at the mydb.pb level
- execute the "server" binary using to start the server at port :8888
	./server

// to run the client
- navigate to the cmd/todo folder
- run the main.go file present in the cmd/todo folder using the command 
	go run main.go list
- to add to the list run the command 
	go run main.go add hello world
