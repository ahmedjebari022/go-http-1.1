package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
)

func main(){
	add, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatalf("Error while resolving the address %s\n",err.Error())
	}
	conn, err := net.DialUDP("udp", add, nil)
	if err != nil {
		log.Fatal("Error while establishing the connexion %s\n",err.Error())
	}
	bufio.NewReader(os.Stdin)
	
}