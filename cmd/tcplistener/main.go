package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"slices"

	"github.com/ahmedjebari022/go-http-1.1/internal/request"
)



func main(){
	fmt.Println("I hope I get the job!")
	listener, err := net.Listen("tcp",":42069")
	if err != nil {
		log.Fatalf("Error while Opening tcp connection %s\n",err.Error())
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting the connection")
			break
		}
		fmt.Println("Connection has been accepted")
		r, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Printf("Error while parsin the request: %s\n", err.Error())
			break
		}
		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n",r.RequestLine.Method,r.RequestLine.RequestTarget,r.RequestLine.HttpVersion)
		fmt.Println("Connection has been closed")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	b := make([]byte, 8)
	p := make([]byte, 8)
	c := make(chan string)
	line := ""
	go func(){
		defer close(c)
		for{
			n, err := f.Read(b)
			p = b[:n]
			if err == io.EOF {
				c <- line
				return
			}	
			if i := slices.Index(p, 10); i != -1{
				bnl := p[:i]
				line += string(bnl)
				c <- line
				anl := p[i+1:]
				line = string(anl)
			}else{
				line += string(p)
			}
		}
	}()
	return c
}
