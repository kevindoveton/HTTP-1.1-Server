package connection

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
)

// CRLF - Carriage Return, Line Feed
const CRLF = "\r\n"

// HandleConnection - Handle a connection
func HandleConnection(conn net.Conn) {

	// Close connection when this function ends
	defer func() {
		fmt.Println("Closing connection...")
		conn.Close()
	}()

	// output message received
	fmt.Println("Message Received")

	req := bufio.NewReader(conn)

	// get a map of headers
	headers := parseRequest(req)

	// get the body of request
	body := parseBody(req, headers)

	// sample response
	response := response(conn, headers)

	// send new string back to client
	conn.Write([]byte(response))

	fmt.Println(body)
}

func response(conn net.Conn, reqHeaders map[string]string) string {
	var res bytes.Buffer

	statusCode := "200 OK"
	bodyContent := "i ♥ u"

	headers := make(map[string]string)

	headers["Content-Length"] = strconv.Itoa(len(bodyContent))
	headers["Content-Type"] = "text/html; charset=utf-8"

	res.WriteString("HTTP/1.1 ")
	res.WriteString(statusCode)
	res.WriteString(CRLF)

	for k := range headers {
		res.WriteString(k)
		res.WriteString(": ")
		res.WriteString(headers[k])
		res.WriteString("\n")
	}

	method, methodOk := reqHeaders["Method"]
	if methodOk && method != "HEAD" {
		res.WriteString(CRLF)
		res.WriteString(bodyContent)
	}

	return res.String()
}

func parseRequest(req *bufio.Reader) map[string]string {
	headers := make(map[string]string)

	char, _ := req.ReadBytes(' ')
	headers["Method"] = string(char)[:len(char)-1]

	char, _ = req.ReadBytes(' ')
	headers["Path"] = string(char)[:len(char)-1]

	char, _ = req.ReadBytes('\n')
	headers["Version"] = string(char)[:len(char)-2]

	foundClrf := false
	for !foundClrf {
		key := make([]byte, 0)
		value := make([]byte, 0)
		foundKey := false
		foundValue := false

		charPeek, _ := req.Peek(2)
		if string(charPeek) == "\r\n" {
			// found clrf
			req.Discard(2)
			foundClrf = true
			break
		}

		// find the key
		for !foundKey {
			char, _ := req.ReadByte()
			if char == ':' {
				foundKey = true
				break
			} else {
				key = AppendByte(key, char)
			}
		}

		// find the value
		// first check if it is a space
		char := stripLeadingWhitespace(req)
		value = AppendByte(value, char)

		for !foundValue {
			char, _ := req.ReadByte()
			if char == '\r' {
				charPeek, _ := req.Peek(1)
				if string(charPeek) == "\n" {
					req.Discard(1) // discard \n
					break
				}
			}
			if char == '\n' {
				foundValue = true
				break
			} else {
				value = AppendByte(value, char)
			}
		}
		headers[string(key)] = string(value)
	}

	PrettyPrint(headers)

	return headers
}

func parseBody(req *bufio.Reader, headers map[string]string) string {
	contentLength, _ := strconv.Atoi(headers["Content-Length"])
	body, _ := req.Peek(contentLength)
	return string(body)
}

// stripLeadingWhitespace - remove all leading whitespace
// returns the first non space char
func stripLeadingWhitespace(reader *bufio.Reader) byte {
	for true {
		char, _ := reader.ReadByte()
		if char != ' ' {
			return char
		}
	}
	return 0
}

// AppendByte - grow array
func AppendByte(slice []byte, data ...byte) []byte {
	m := len(slice)
	n := m + len(data)
	if n > cap(slice) { // if necessary, reallocate
		// allocate double what's needed, for future growth.
		newSlice := make([]byte, (n+1)*2)
		copy(newSlice, slice)
		slice = newSlice
	}
	slice = slice[0:n]
	copy(slice[m:n], data)
	return slice
}

func PrettyPrint(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	println(string(b))
}
