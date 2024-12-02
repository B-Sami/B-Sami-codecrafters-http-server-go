package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net"
	"os"
	"path"
	"strings"
)

func createHeaders(body string, contentType string) map[string]string {
	headers := map[string]string{
		"Content-Type":   contentType,
		"Content-Length": fmt.Sprintf("%d", len(body)),
	}
	return headers
}

func handleRequest(conn net.Conn) {
	buf := make([]byte, 1024)
	if _, err := conn.Read(buf); err != nil {
		fmt.Println("Error reading request:", err)
		conn.Close()
		return
	}

	request := string(buf)
	lines := strings.Split(request, "\r\n")
	methodAndPath := strings.Split(lines[0], " ")
	if len(methodAndPath) < 3 {
		fmt.Println("Invalid request format")
		conn.Close()
		return
	}

	method := methodAndPath[0]
	urlPath := methodAndPath[1]

	headers := parseHeaders(lines)

	supportsGzip := strings.Contains(headers["Accept-Encoding"], "gzip")
	userAgent := getUserAgent(lines)

	var responseBody string
	var statusCode int
	var statusPhrase string
	var contentType string = "text/plain"

	switch method {
	case "GET":
		responseBody, statusCode, statusPhrase, contentType = handleGetRequest(urlPath, userAgent, supportsGzip)
	case "POST":
		responseBody, statusCode, statusPhrase = handlePostRequest(urlPath, buf)
	default:
		responseBody = ""
		statusCode = 405
		statusPhrase = "Method Not Allowed"
	}

	responseHeaders := createHeaders(responseBody, contentType)
	if supportsGzip && method == "GET" && strings.HasPrefix(urlPath, "/echo/") && responseBody != "" {
		responseHeaders["Content-Encoding"] = "gzip"
	}
	response := NewServerResponse(statusCode, statusPhrase, responseHeaders, responseBody)
	if _, err := conn.Write([]byte(response.FullResponse())); err != nil {
		fmt.Println("Error sending response:", err)
	}

	conn.Close()
}

func handleGetRequest(urlPath string, userAgent string, supportsGzip bool) (string, int, string, string) {
	var responseBody string
	var statusCode int = 404 // Default to Not Found
	var statusPhrase = "Not Found"
	var contentType string = "text/plain"

	switch urlPath {
	case "/":
		responseBody = ""
		statusCode = 200
		statusPhrase = "OK"
	case "/user-agent":
		responseBody = userAgent
		statusCode = 200
		statusPhrase = "OK"
	default:
		if strings.HasPrefix(urlPath, "/echo/") {
			echoStr := strings.TrimPrefix(urlPath, "/echo/")
			if echoStr != "" && supportsGzip {
				var buffer bytes.Buffer
				w := gzip.NewWriter(&buffer)
				w.Write([]byte(echoStr))
				w.Close()
				echoStr = buffer.String()
			}
			responseBody = echoStr
			if echoStr != "" {
				statusCode = 200
				statusPhrase = "OK"
			} else {
				responseBody = "Echo path is empty"
			}
		} else if strings.HasPrefix(urlPath, "/files/") {
			dir := os.Args[2]
			fileName := strings.TrimPrefix(urlPath, "/files/")
			contentFile, err := os.ReadFile(path.Join(dir, fileName))
			if err == nil { // File found and read successfully.
				responseBody = string(contentFile[:])
				statusCode = 200
				statusPhrase = "OK"
				contentType = "application/octet-stream" // Correctly set content type for files
			} else { // File not found or error reading.
				responseBody = "File not found or error reading."
			}
		} else { // Other paths not found.
			responseBody = "Page not found."
		}
	}

	return responseBody, statusCode, statusPhrase, contentType
}

func parseHeaders(lines []string) map[string]string {
	headers := make(map[string]string)
	for _, line := range lines[1:] {
        if line == "" {
            break
        }
        parts := strings.SplitN(line, ": ", 2)
        if len(parts) == 2 {
            headers[parts[0]] = parts[1]
        }
    }
	return headers
}

func getUserAgent(lines []string) string {
	for _, line := range lines[1:] {
        if strings.HasPrefix(line, "User-Agent:") {
            return strings.TrimPrefix(line, "User-Agent: ")
        }
    }
	return ""
}

func handlePostRequest(urlPath string, buf []byte) (string,int,string){
	var responseBody string 
	var statusCode int 
	var statusPhrase string 
	switch{
	case strings.HasPrefix(urlPath,"/files/"):
	content:=strings.Trim(strings.Split(string(buf),"\r\n")[len(strings.Split(string(buf),"\r\n"))-1],"\x00")
	dir:=os.Args[2]
	err:=os.WriteFile(path.Join(dir,urlPath[7:]),[]byte(content),0644)
	if err!=nil{
	responseBody="Error writing file."
	statusCode=404 
	statusPhrase="Not Found"}else{
	responseBody="HTTP/1.1 201 Created\r\n\r\n"
	statusCode=201 
	statusPhrase="Created"}
	default:
	responseBody="Page not found."
	statusCode=404 
	statusPhrase="Not Found"}
	return responseBody,statusCode,statusPhrase 
}
