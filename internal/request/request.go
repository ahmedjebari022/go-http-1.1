package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

type state int 
const (
	initialized = iota
	done 
)

type Request struct {
	RequestLine 	RequestLine
	State			state
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
		State: 0,
	}
	buf := make([]byte, bufferSize, bufferSize)
	readToIndex := 0
	for request.State == 0 {
		if len(buf) == cap(buf){
			cp := make([]byte, len(buf)*2, cap(buf)*2)
			copy(cp,buf)
			buf = cp
		}
		nr, err := reader.Read(buf[readToIndex:])
		if err != nil {
			return &request, err
		}
		readToIndex += nr
		np, err := request.parse(buf)
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
	return requestLine, nil, i+1
}

func checkValidHttpVersion(httpVersion string) bool {
	return httpVersion == "1.1"
}

func checkValidMethod(method string) bool {
	method = strings.ToUpper(method)
	return method == "POST" || method == "GET" || method == "PUT" || method == "DELETE" 
}


func (r *Request) parse(data []byte) (int, error){
	if r.State == 1{
		return 0, errors.New("error: trying to read data in a done state")
	}
	requestLine, err, n := parseRequestLine(data)
	if err != nil {
		return 0, err
	}
	if n == 0 {
		return 0, nil
	}
	r.RequestLine = requestLine
	r.State = 1
	return n, nil


}