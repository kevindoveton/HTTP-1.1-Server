package connection

import (
	"bufio"
	"bytes"
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

	// get a map of headers
	headers := parseRequest(conn)

	// sample response
	response := response(conn, headers)

	// send new string back to client
	conn.Write([]byte(response))
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

func parseRequest(conn net.Conn) map[string]string {
	// will listen for message to process ending in newline (\n)
	// req, _ := bufio.NewReader(conn).ReadString('\n')
	// scanner := bufio.NewScanner(conn)
	// var pos int
	// var index int

	// scanner.Scan()
	// line := scanner.Text()

	// // find method
	// index = strings.Index(line, " ")
	// headers["Method"] = line[0:index]

	// // find path
	// pos = index + 1
	// index = strings.Index(line[pos:], " ") + pos
	// headers["Path"] = line[pos:index]

	// // find version
	// pos = index + 1
	// headers["Version"] = line[pos:]

	headers := make(map[string]string)
	req := bufio.NewReader(conn)

	char, _ := req.ReadBytes(' ')
	headers["Method"] = string(char)[:len(char)-1]

	char, _ = req.ReadBytes(' ')
	headers["Path"] = string(char)[:len(char)-1]

	char, _ = req.ReadBytes('\n')
	headers["Version"] = string(char)[:len(char)-1]

	foundClrf := false
	for !foundClrf {
		key := make([]byte, 0)
		value := make([]byte, 0)
		foundKey := false
		foundValue := false

		// find the key
		for !foundKey {
			char, _ := req.ReadByte()

			if char == ':' {
				foundKey = true
			} else {
				key = AppendByte(key, char)
			}
		}

		// find the value
		// first check if it is a space
		char, _ := req.ReadByte()
		if char == '\n' {
			foundValue = true
		} else if char == ' ' {
		}
		for !foundValue {
			char, _ := req.ReadByte()
			if char == '\n' {
				foundValue = true
				value = AppendByte(value, char)
			} else {
				value = AppendByte(value, char)
			}
		}

		foundClrf = true
		fmt.Println(string(key))
		fmt.Println(string(value))
	}

	return headers
}

// stripLeadingWhitespace - remove all leading whitespace
func stripLeadingWhitespace(reader bufio.Reader) byte {
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
