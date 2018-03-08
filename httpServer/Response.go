package httpServer

import (
	"net"
	"os"
	"io/ioutil"
	"strconv"
	"bytes"
  "strings"
)

type Response struct {
	NetConn net.Conn
	Status string
}

func (res *Response) Send() {
  var r bytes.Buffer
  r.WriteString("HTTP/1.1 ")
  r.WriteString(res.Status)
  r.WriteString(CRLF)
  res.NetConn.Write(r.Bytes())
}

func (res *Response) Send404() {
  res.SetStatusCode(404)
  res.Send()
}

func (res *Response) SetStatusCode(statusCode int) {
	res.Status = strconv.Itoa(statusCode)
	return
}

func (res *Response) SendString(body string)  {
  body = strings.Join([]string{body, "\n"}, "")
	var r bytes.Buffer
	headers := make(map[string]string)

	headers["Content-Length"] = strconv.Itoa(len(body))
	headers["Content-Type"] = "text/html; charset=utf-8"

	r.WriteString("HTTP/1.1 ")
	r.WriteString(res.Status)
	r.WriteString(CRLF)

	for k := range headers {
		r.WriteString(k)
		r.WriteString(": ")
		r.WriteString(headers[k])
		r.WriteString("\n")
	}

	//method, methodOk := reqHeaders["Method"]
	//if methodOk && method != "HEAD" {
		r.WriteString(CRLF)
		r.WriteString(body)
  //}
	res.NetConn.Write(r.Bytes())
	return
}

func sendFile(conn net.Conn, reqHeaders map[string]string) (string, int) {
	var res bytes.Buffer

	statusCode := "200 OK"
	//path := strings.Join([]string{server.WebRoot, reqHeaders["Path"]}, "")
	path := ""

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