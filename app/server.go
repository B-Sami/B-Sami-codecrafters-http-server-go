package main

import (
	"fmt"
	"net"
	"os"
)

// Ensures gofmt doesn't remove the "net" and "os" imports above (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	// Print connection details
	fmt.Println("Accepted new connection:")
	fmt.Println("Remote Address:", conn.RemoteAddr()) // Remote (client) address
	fmt.Println("Local Address:", conn.LocalAddr())   // Local (server) address

	// Optional: Print the connection directly (may provide limited details)
	fmt.Printf("Conn object: %v\n", conn)
}
