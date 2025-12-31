package response

import (
	"fmt"
	"io"
	"strconv"

	header "github.com/ahmedjebari022/go-http-1.1/internal/headers"
)


type Response struct{

}

type StatusCode int
const (
	Success StatusCode = iota
	ClientError
	ServerError
)


func (c StatusCode) String() string {
	cm := map[StatusCode]string {
		Success: "200",
		ClientError: "400",
		ServerError: "500",
	}
	if v, ok := cm[c]; ok {
		return v
	}
	return ""
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error{
	statusLine := "HTTP/1.1 "
	switch statusCode{
	case Success:
		statusLine += statusCode.String() + " OK"
	case ClientError:
		statusLine +=  statusCode.String() + " Bad Request"	
	case ServerError:
		statusLine += statusCode.String() + " Internal Server Error"	
	default: 
		statusLine += statusCode.String()
	}
	fmt.Println(statusLine)
	_, err := w.Write([]byte(statusLine + "\r\n"))
	if err != nil {
		fmt.Println("error when writing status line to conn")
		return err
	}
	return nil
}

func GetDefaultHeaders(contentLen int) header.Headers{
	headers := header.NewHeaders()
	contentLenString := strconv.Itoa(contentLen)
	headers["Content-Length"] = contentLenString
	headers["Connection"] = "close"
	headers["Content-Type"] = "text/plain"
	return headers
}


func WriteHeaders(w io.Writer, headers header.Headers) error {
	header := ""
	for key, value := range headers{
		header += fmt.Sprintf("%s: %s\r\n",key,value)
	}
	header += "\r\n"
	fmt.Printf("headers: %s", header)
	_, err := w.Write([]byte(header))
	if err != nil {
		fmt.Println("error when writing header to conn")
		return err
	}
	return nil
}