package httpServer

import (
	"net"
	"os"
	"io/ioutil"
	"strconv"
	"bytes"
  "strings"
  "log"
  "fmt"
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

  r.WriteString(CRLF)
  r.WriteString(body)

	res.NetConn.Write(r.Bytes())
	return
}

func (res *Response) SendFile(path string) {

	// check if we can access the file
	if fi, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			res.SetStatusCode(404)
		} else {
			res.SetStatusCode(500)
			// other error - maybe perms
		}
    res.Send()
    return
	} else {
		// it exists, may be a directory
		switch mode := fi.Mode(); {
      case mode.IsDir():
        res.SetStatusCode(405)
        res.Send()
        return
    }
	}

	res.sendFile(path)
}

func (res *Response) sendFile(path string) {
  dat, _ := ioutil.ReadFile(path)

  body := string(dat)

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

  r.WriteString(CRLF)
  r.WriteString(body)

  res.NetConn.Write(r.Bytes())
}

func (res *Response) SendStatic(matchedPath string, dir string, actualPath string) {

  path := strings.Replace(actualPath, matchedPath, dir, 1)

  // check if it's a file and we can access the file
  if fi, err := os.Stat(path); err != nil {
    if os.IsNotExist(err) {
      res.SetStatusCode(404)
    } else {
      res.SetStatusCode(500)
      // other error - maybe perms
    }
    res.Send()
    return
  } else {
    // it exists, may be a directory
    switch mode := fi.Mode(); {
    case mode.IsDir():
      fmt.Print(path)
      res.SendStaticListing(path)
      return
    }
  }

  res.sendFile(path)
}

func (res *Response) SendStaticListing(dir string) {
  files, err := ioutil.ReadDir(dir)
  if err != nil {
    log.Fatal(err)
  }
  html := "<!DOCTYPE html><head></head><body><p>"
  for _, f := range files {
    html += "<a href=\""+f.Name()+"\">"+f.Name()+"</a><br>"
  }
  html += "</p></body>"
  res.SendString(html)
}
