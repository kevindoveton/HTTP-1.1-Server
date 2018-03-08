package httpServer

import (
	"net"
	"bufio"
	"strconv"
	"strings"
)

type Connection struct {
	req Request
	res Response
	netConn net.Conn
	router *Router
}

// Handle - handle the connection
func (conn *Connection) Handle() {

	req := bufio.NewReader(conn.netConn)

	// make sure its not an empty request
	// postman was doing some funny things
	r, _ := req.Peek(1)
	if string(r) == "" {
		response := sendError(500)
		conn.netConn.Write([]byte(response))
		return
	}

	// create the request
	headers := parseRequest(req)
	contentLength, _ :=strconv.Atoi(headers["Content-Length"])
	body := parseBody(req, contentLength)

	conn.req = Request{
		headers: headers,
		body: body,
	}

	conn.res = Response{
	  conn.netConn,
	  "200 OK",
  }

	conn.router.GetResponse(&conn.req, &conn.res)

	// create the response

	// sample response
	//response, err := response(conn.netConn, headers)
	//if err != 0 {
	//	response := sendError(err)
	//	conn.netConn.Write([]byte(response))
	//	return
	//}
	//
	//// send new string back to client
	//conn.netConn.Write([]byte(response))
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
		if string(charPeek) == CRLF {
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

	return headers
}

// parseBody - Parse the body of the message
func parseBody(req *bufio.Reader, contentLength int) string {
	body, _ := req.Peek(contentLength)
	return string(body)
}

func sendError(status int) string {
	return strings.Join([]string{"HTTP/1.1", strconv.Itoa(status), CRLF, strconv.Itoa(status)}, " ")
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

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}