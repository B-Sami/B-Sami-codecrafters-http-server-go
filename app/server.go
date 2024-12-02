package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
)

type Server struct {
	Network string
	Addr    string
	Handler http.Handler // Although not used, kept for potential future use
}

func (s *Server) ListenAndServe(address string) error {
	fmt.Println("Starting server on", address)
	listener, err := net.Listen(s.Network, address)
	if err != nil {
		fmt.Println("Failed to bind to", address, "Error:", err)
		os.Exit(1)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		fmt.Println("Connection accepted:", conn.RemoteAddr())
		go handleRequest(conn)
	}
}
