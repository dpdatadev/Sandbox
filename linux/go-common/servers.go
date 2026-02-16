package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type outputResponse struct {
	title  string
	output string
}

func (or *outputResponse) Stringify() string {
	return "SERVICE :" + or.title + " :: STDOUT :: " + or.output
}

// todo figure out the best way to incorporate HTTP for the command runners
func serviceEchoHandler(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	store := NewInMemoryStore()
	exec := NewLocalExecutor()

	svc := NewCommandService(store, exec)

	cmd := NewCommand(
		"ip",
		[]string{"neigbor"},
		"test command",
	)

	if err := svc.Run(ctx, cmd); err != nil {
		panic(err) //todo remove all panics from framework code
	}
	outPutString := &outputResponse{"ECHO SERVICE TEST: " + req.Method + " :: ", cmd.Stdout}
	fmt.Fprintf(w, "%s", outPutString.Stringify())
}

func serviceStatusHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Service up ::HIT::<request><<%s>>::\n", req.URL)
}

func StartServer() {
	http.HandleFunc("/status", serviceStatusHandler)
	http.HandleFunc("/service/echo", serviceEchoHandler)

	fmt.Println("Server starting on port 8082..")
	if err := http.ListenAndServe(":8082", nil); err != nil {
		fmt.Printf("HTTP server failed: %v\n", err)
	}
}
