package connection

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"strconv"
	"strings"
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
	scanner := bufio.NewScanner(conn)
	var pos int
	var index int

	headers := make(map[string]string)
	scanner.Scan()
	line := scanner.Text()

	// find method
	index = strings.Index(line, " ")
	headers["Method"] = line[0:index]

	// find path
	pos = index + 1
	index = strings.Index(line[pos:], " ") + pos
	headers["Path"] = line[pos:index]

	// find version
	pos = index + 1
	headers["Version"] = line[pos:]

	// breaks on the first empty line
	// technically the message could have a body below!
	for scanner.Scan() {
		line = scanner.Text()
		if line != "" {
			keyIndex := strings.Index(line, ":")
			headers[line[0:keyIndex]] = line[keyIndex+2:]
		} else {
			break
		}
	}

	return headers
}
