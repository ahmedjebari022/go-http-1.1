package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ahmedjebari022/go-http-1.1/internal/request"
	"github.com/ahmedjebari022/go-http-1.1/internal/response"
	"github.com/ahmedjebari022/go-http-1.1/internal/server"
)



const port = 42069
func main(){
	server, err := server.Serve(port,testHandler)	
	if err != nil {
		log.Fatal("Error starting server : %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)
	server.Listen()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stoped")
}


func testHandler(w io.Writer, req *request.Request) *server.HandlerError{
	requestTarget := req.RequestLine.RequestTarget
	switch requestTarget{
	case "/yourproblem": 
		return &server.HandlerError{
            StatusCode:   response.ClientError,
            ErrorMessage: "Your problem is not my problem\n",
        }
	case "/myproblem":
		return &server.HandlerError{
            StatusCode:   response.ServerError,
            ErrorMessage:  "Woopsie, my bad\n",
        }
	default :
		w.Write([]byte("All good, frfr\n"))
		return nil
	}
	
}