package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	header "github.com/ahmedjebari022/go-http-1.1/internal/headers"
)

type state int 
const (
	Initialized state = iota
	Done 
	ParsingHeaders
	ParsingBody
)

type Request struct {
	RequestLine 	RequestLine
	State			state
	Header 			header.Headers
	Body			[]byte
}

type RequestLine struct {
	HttpVersion 	string
	RequestTarget 	string
	Method 			string
}

const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := Request{
		RequestLine: RequestLine{},
		Header : header.NewHeaders(),
		State: 0,
	}
	buf := make([]byte, bufferSize)
	readToIndex := 0
	for request.State != 1 {
		if len(buf) == readToIndex{
			cp := make([]byte, len(buf)*2)
			copy(cp,buf)
			buf = cp
		}
		nr, err := reader.Read(buf[readToIndex:])
		if err != nil {
			return &request, err
		}
		readToIndex += nr
		np, err := request.parse(buf[:readToIndex])
		if err != nil {
			return &request, err
		}
		newb := buf[np:]	
		copy(buf, newb)
		readToIndex -= np
	}
	return &request, nil
}



func parseRequestLine(data []byte) (RequestLine, error, int) {
	requestLineData := ""
	i := bytes.Index(data, []byte("\r\n")) 
	if i != -1 {
		requestLineData = string(data[:i])
	}else{
		return RequestLine{}, nil, 0
	}
	requestLineSlice := strings.Split(requestLineData, " ")
	if len(requestLineSlice) != 3 {
		return RequestLine{}, fmt.Errorf("Request line formulated badly") , 0
	}
	if vm := checkValidMethod(requestLineSlice[0]); !vm {
		return RequestLine{}, fmt.Errorf("Invalid http method"), 0
	}
	if vh := checkValidHttpVersion(strings.Split(requestLineSlice[2], "/")[1]); !vh{
		return RequestLine{}, fmt.Errorf("Invalid http Version"), 0
	}
	requestLine := RequestLine{
		Method: requestLineSlice[0],
		RequestTarget: requestLineSlice[1],
		HttpVersion: strings.Split(requestLineSlice[2], "/")[1],
	}
	return requestLine, nil, i+2
}

func checkValidHttpVersion(httpVersion string) bool {
	return httpVersion == "1.1"
}

func checkValidMethod(method string) bool {
	method = strings.ToUpper(method)
	return method == "POST" || method == "GET" || method == "PUT" || method == "DELETE" 
}


func (r *Request) parse(data []byte) (int, error){ 
	totalBytesParse := 0
	for r.State != 1 {
		n, err := r.parseSingle(data[totalBytesParse:])
		if err != nil {
			return 0, err
		}
		totalBytesParse += n
		if n == 0 {
			break
		}
	}
	return totalBytesParse, nil 	

}



func (r *Request) parseSingle(data []byte) (int, error){
	switch r.State {
	case 1 :
		return 0, errors.New("error: trying to read data in a done state")
	case 0 :
		requestLine, err, n := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}
		r.RequestLine = requestLine
		r.State = 2
		return n, nil
	case 2 :		
		n, done , err := r.Header.Parse(data)	
		if err != nil {
			return 0, err
		}
		if done {
			r.State = ParsingBody
			return n, nil
		}
		return n, nil
	case 3 :
		bodyLength, err := r.Header.Get("Content-length")
		if err != nil {
			return 0, err
		}
		if bodyLength == -1 {
			r.State = Done
			return 0, nil
		}
		r.Body = append(r.Body, data...)
		if len(r.Body) == bodyLength{
			r.State = Done										
		}else if len(r.Body) > bodyLength{
			return 0, fmt.Errorf("content length header and body length don't match")
		}
		fmt.Println("consumed the whole data")
		return len(data), nil
	default: 
		return 0, fmt.Errorf("unknown state")
	}
}