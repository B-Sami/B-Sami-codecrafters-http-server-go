package main

import (
	"fmt"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	server := &Server{
		Network: "tcp",
		Addr:    ":4221",
		Handler: nil,
	}

	if err := server.ListenAndServe(server.Addr); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
