package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main(){
	add, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatalf("Error while resolving the address %s\n",err.Error())
	}
	conn, err := net.DialUDP("udp", nil, add)
	if err != nil {
		log.Fatalf("Error while establishing the connexion %s\n",err.Error())
	}
	defer conn.Close()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println(">")
		data, err := reader.ReadString(10)
		if err != nil {
			fmt.Printf("Error while reading data: %s \n",err.Error())	
		}
		_, err = conn.Write([]byte(data))
		if err != nil {
			fmt.Printf("Error while writing to connexion: %s\n",err.Error())
		}
	}
}