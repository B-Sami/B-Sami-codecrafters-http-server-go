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
    Handler http.Handler
}

type StatusLine struct {
	HTTPVersion  string
	StatusCode   int   
	StatusPhrase string
}
type HTTPResponse struct {
	StatusLine StatusLine
	Headers map[string]string
	ResponseBody string
}

func NewServerResponse() *HTTPResponse {
	return &HTTPResponse{
		StatusLine: StatusLine{
			HTTPVersion:  "HTTP/1.1",
			StatusCode:   200,
			StatusPhrase: "OK",
		},
		Headers:      make(map[string]string),
		ResponseBody: "",
	}
}

func (r *HTTPResponse) FormatStatusLine() string {
	return fmt.Sprintf("%s %d %s\r\n", r.StatusLine.HTTPVersion, r.StatusLine.StatusCode, r.StatusLine.StatusPhrase)
}

func (r *HTTPResponse) FormatHeaders() string {
	var headers string
	for key, value := range r.Headers {
		headers += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	if headers == "" {
		return "\r\n"
	}
	return headers + "\r\n"
}

func (r *HTTPResponse) FormatBody() string {
	return r.ResponseBody
}

func (r *HTTPResponse) FullResponse() string {
	return r.FormatStatusLine() + r.FormatHeaders() + r.FormatBody()
}

func main() {
	fmt.Println("Logs from your program will appear here!")

    server := &Server{
		Network: "tcp",
		Addr:    ":4221",
		Handler: http.DefaultServeMux,
	}

	if err := server.ListenAndServe(server.Addr); err != nil {
		fmt.Println("Error starting server:", err)
	}
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

func handleRequest(conn net.Conn) {
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading request:", err)
		conn.Close()
		return
	}

	response := NewServerResponse()

	_, err = conn.Write([]byte(response.FullResponse()))
	if err != nil {
		fmt.Println("Error sending response:", err)
		conn.Close()
		return
	}

	conn.Close()
}