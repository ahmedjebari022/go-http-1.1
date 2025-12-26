package header

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Headers map[string]string					

func NewHeaders() Headers{
	return Headers{}
}
const CRLF = "\r\n"
func (h Headers) Parse(data []byte) (n int, done bool, err error){
	i := bytes.Index(data, []byte(CRLF) )		
	if i == -1 {
		return 0, false, nil
	}else if i == 0 {
		return 2, true, nil
	}
	dataslice := strings.SplitN(string(data[:i]), ":", 2)
	if len(dataslice) < 2 {
		return 0, true, errors.New("Poorly formatted header")
	}
	fieldName := strings.TrimLeft(dataslice[0], " ")
	fieldValue := dataslice[1]
	if c := strings.ContainsFunc(fieldName, containsSeparator); c || len(fieldName) < 1{
		return 0, false, errors.New("Field name poorly formatted")
	}
	fieldName = strings.ToLower(strings.TrimRight(fieldName, " "))
	fieldValue = strings.Trim(fieldValue, " ")
	if _, ok := h[fieldName]; ok {
		fmt.Println("Found key map")
		h[fieldName] += fmt.Sprintf(", %s", fieldValue)
	}else{
		h[fieldName] = fieldValue
	}
	return i+2, false, nil
}


func containsSeparator(r rune ) bool {
	return r == ' ' || r == '/' || r == ',' || r == ';' || r == ':' || r == '\\'
}



func (h Headers) Get(key string) (int, error) {
	key = strings.ToLower(key)
	v, ok := h[key]
	if !ok {
		return -1, nil
	}
	bodyLength, err := strconv.Atoi(v)
	if err != nil {
		return 0, err
	}
	return bodyLength, nil

}