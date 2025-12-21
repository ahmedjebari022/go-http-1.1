package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"slices"
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
		c := getLinesChannel(conn)
		for v := range c {
			fmt.Printf("%s\n",v)
		}
		fmt.Println("Connection has been closed")

	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	p := make([]byte, 8)
	c := make(chan string)
	line := ""
	go func(){
		defer close(c)
		for{
			_, err := f.Read(p)
			if err == io.EOF {
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
