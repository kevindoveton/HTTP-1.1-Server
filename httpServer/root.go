package httpServer

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
)

// CRLF - Carriage Return, Line Feed
const CRLF = "\r\n"

var webRoot = ""

// SetWebRoot - set the web root dir
func SetWebRoot(path string) {
	webRoot = path
}

// GetWebRoot - return set web root
func GetWebRoot() string {
	return webRoot
}

// Run - Run the server
func Run(port int) {
	fmt.Println("Server opened on port " + strconv.Itoa(port) + ".")

	// listen on all interfaces
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))

	if (err != nil) {
		fmt.Println(err)
		return
	}

	// run loop forever (or more likely until ctrl-c)
	for {
		// wait for a connection
		conn, err := ln.Accept()

		if err != nil {
			fmt.Println(err)
			break
		}
		go HandleConnection(conn)
	}
}


// HandleConnection - Handle a connection
func HandleConnection(conn net.Conn) {

	// Close connection when this function ends
	defer func() {
		//fmt.Println("Closing connection...")
		conn.Close()
	}()

	// output message received
	//fmt.Println("Message Received")

	req := bufio.NewReader(conn)

	// make sure its not an empty request
	// postman was doing some funny things
	r, _ := req.Peek(1)
	if string(r) == "" {
		response := sendError(500)
		conn.Write([]byte(response))
		return
	}

	// get a map of headers
	headers := parseRequest(req)

	// get the body of request
	_ = parseBody(req, headers)

	// sample response
	response, err := response(conn, headers)
	if err != 0 {
		response := sendError(err)
		conn.Write([]byte(response))
		return
	}

	// send new string back to client
	conn.Write([]byte(response))

	//fmt.Println()
	//fmt.Println("Headers:")
	//PrettyPrint(headers)
	//fmt.Println()
	//fmt.Println("Body:")
	//fmt.Println(body)
	//fmt.Println()
}

func response(conn net.Conn, reqHeaders map[string]string) (string, int) {
	var res bytes.Buffer

	statusCode := "200 OK"
	path := strings.Join([]string{webRoot, reqHeaders["Path"]}, "")

	// check if we can access the file
	if fi, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return "", 404
		} else {
			return "", 500
			// other error - maybe perms
		}
	} else {
		// it exists, may be a directory
		switch mode := fi.Mode(); {
		case mode.IsDir():
			return "", 405
		}
	}

	dat, _ := ioutil.ReadFile(path)

	bodyContent := string(dat)

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

	return res.String(), 0
}

func parseRequest(req *bufio.Reader) map[string]string {
	headers := make(map[string]string)

	char, _ := req.ReadBytes(' ')
	//fmt.Println(string(char))
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
func parseBody(req *bufio.Reader, headers map[string]string) string {
	contentLength, _ := strconv.Atoi(headers["Content-Length"])
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

func PrettyPrint(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	println(string(b))
}
