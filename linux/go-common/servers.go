package main

import (
	"fmt"
	"net/http"
)

func serviceStatusHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Service up ::HIT::<request><<%s>>::\n", req.URL)
}

func StartServer() {
	http.HandleFunc("/status", serviceStatusHandler)

	fmt.Println("Server starting on port 8082..")
	if err := http.ListenAndServe(":8082", nil); err != nil {
		fmt.Printf("HTTP server failed: %v\n", err)
	}
}
