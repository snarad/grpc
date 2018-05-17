package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/snarad/grpc/just_for_func/todo"
	"golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

func main() {
	srv := grpc.NewServer()
	var tasks taskServer
	todo.RegisterTasksServer(srv, tasks)
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("could not listen to :8888: %v", err)
	}
	log.Fatal(srv.Serve(l))
}

const dbPath = "mydb.pb"

func (s taskServer) Add(ctx context.Context, text *todo.Text) (*todo.Task, error) {
	task := &todo.Task{
		Text: text.Text,
		Done: false,
	}
	b, err := proto.Marshal(task)
	if err != nil {
		return nil, fmt.Errorf("could not encode task: %v ", err)
	}

	f, err := os.OpenFile(dbPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("could not open %s : %v", dbPath, err)
	}

	if err := gob.NewEncoder(f).Encode(int64(len(b))); err != nil {
		return nil, fmt.Errorf("could not endode length of message: %v", err)
	}

	_, err = f.Write(b)
	if err != nil {
		return nil, fmt.Errorf("could not write task to file: %v", err)
	}

	if err := f.Close(); err != nil {
		return nil, fmt.Errorf("could not close file %s : %v", dbPath, err)
	}

	fmt.Println(proto.MarshalTextString(task))
	return task, nil
}

type taskServer struct{}

func (s taskServer) List(ctx context.Context, void *todo.Void) (*todo.TaskList, error) {
	// return nil, fmt.Errorf("not implemented")

	b, err := ioutil.ReadFile(dbPath)
	if err != nil {
		return nil, fmt.Errorf("could not read %s: %v", dbPath, err)
	}

	var tasks todo.TaskList

	for {
		if len(b) == 0 {
			return &tasks, nil
		} else if len(b) < 4 {
			return nil, fmt.Errorf("remaining odd %d bytes, what to do? ", len(b))
		}
		var length int64
		if err := gob.NewDecoder(bytes.NewReader(b[:4])).Decode(&length); err != nil {
			return nil, fmt.Errorf("could not decode message length: %v", err)
		}

		b = b[4:]

		var task todo.Task
		if err := proto.Unmarshal(b[:length], &task); err != nil {
			return nil, fmt.Errorf("could not read task: %v", err)
		}

		b = b[length:]
		tasks.Tasks = append(tasks.Tasks, &task)

		if task.Done {
			fmt.Printf("done emoji")
		} else {
			fmt.Printf("not done emoji")
		}
		fmt.Printf(" %s\n", task.Text)
	}

}
