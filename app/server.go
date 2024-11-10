package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
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

func NewServerResponse(statusCode int, statusPhrase string, headers map[string]string, body string) *HTTPResponse {
	return &HTTPResponse{
		StatusLine: StatusLine{
			HTTPVersion:  "HTTP/1.1",
			StatusCode:   statusCode,
			StatusPhrase: statusPhrase,
		},
		Headers:      headers,
		ResponseBody: body,
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

func createHeaders(body string) map[string]string {
	headers := map[string]string{
		"Content-Type":   "text/plain",
		"Content-Length": fmt.Sprintf("%d", len(body)),
	}
	return headers
}
func handleRequest(conn net.Conn) {
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading request:", err)
		conn.Close()
		return
	}

	// Parse the incoming request
	request := string(buf)
	lines := strings.Split(request, "\r\n")
	methodAndPath := strings.Split(lines[0], " ")
	if len(methodAndPath) < 3 {
		fmt.Println("Invalid request format")
		conn.Close()
		return
	}

	method := methodAndPath[0]
	path := methodAndPath[1]

	// Handle the GET method
	if method == "GET" {
		var responseBody string
		var statusCode int
		var statusPhrase string

		// Handle root path "/"
		if path == "/" {
			responseBody = ""
			statusCode = 200
			statusPhrase = "OK"
		} else if strings.HasPrefix(path, "/echo/") {
			echoStr := strings.TrimPrefix(path, "/echo/")
			if echoStr != "" {
				responseBody = echoStr
				statusCode = 200
				statusPhrase = "OK"
			} else {
				// If the echo string is empty, return 404
				responseBody = "Echo path is empty"
				statusCode = 404
				statusPhrase = "Not Found"
			}
		} else {
			// For any other undefined path, return 404
			responseBody = "Page not found"
			statusCode = 404
			statusPhrase = "Not Found"
		}

		// Prepare headers
		headers := map[string]string{
			"Content-Type":   "text/plain",
			"Content-Length": fmt.Sprintf("%d", len(responseBody)),
		}

		// Create response and send it
		response := NewServerResponse(statusCode, statusPhrase, headers, responseBody)
		_, err := conn.Write([]byte(response.FullResponse()))
		if err != nil {
			fmt.Println("Error sending response:", err)
		}
	} else {
		headers := map[string]string{
			"Content-Type":   "text/plain",
			"Content-Length": "0",
		}
		response := NewServerResponse(405, "Method Not Allowed", headers, "")
		_, err := conn.Write([]byte(response.FullResponse()))
		if err != nil {
			fmt.Println("Error sending response:", err)
		}
	}

	conn.Close()
}
