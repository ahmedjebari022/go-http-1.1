package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/ahmedjebari022/go-http-1.1/internal/request"
	"github.com/ahmedjebari022/go-http-1.1/internal/response"
	"github.com/ahmedjebari022/go-http-1.1/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, refactoredHandler)
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

func testHandler(w io.Writer, req *request.Request) *server.HandlerError {
	requestTarget := req.RequestLine.RequestTarget
	switch requestTarget {
	case "/yourproblem":
		return &server.HandlerError{
			StatusCode:   response.ClientError,
			ErrorMessage: "Your problem is not my problem\n",
		}
	case "/myproblem":
		return &server.HandlerError{
			StatusCode:   response.ServerError,
			ErrorMessage: "Woopsie, my bad\n",
		}
	default:
		w.Write([]byte("All good, frfr\n"))
		return nil
	}

}

func refactoredHandler(w *response.Writer, req *request.Request) {
	requestTarget := req.RequestLine.RequestTarget
	switch requestTarget {
	case "/yourproblem":
		err := w.WriteStatusLine(response.ClientError)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			return
		}
		w.Header.Set("Content-Type", "text/html")
		resBody := ` 
                <html>
                <head>
                    <title>400 Bad Request</title>
                </head>
                <body>
                    <h1>Bad Request</h1>
                    <p>Your request honestly kinda sucked.</p>
                </body>
                </html>
            `
		w.Header.Set("Content-Length", strconv.Itoa(len(resBody)))
		w.WriteHeaders(w.Header)
		_, err = w.WriteBody([]byte(resBody))
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			return
		}
	case "/myproblem":
		err := w.WriteStatusLine(response.ServerError)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			return
		}
		w.Header.Set("Content-Type", "text/html")
		resBody := ` 
                <html>
                <head>
                    <title>500 Internal Server Error</title>
                </head>
                <body>
                    <h1>Internal Server Error</h1>
                    <p>Okay, you know what? This one is on me.</p>
                </body>
                </html>
            `
		w.Header.Set("Content-Length", strconv.Itoa(len(resBody)))
		w.WriteHeaders(w.Header)
		_, err = w.WriteBody([]byte(resBody))
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			return
		}
	default:
		err := w.WriteStatusLine(response.Success)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			return
		}
		resBody := `
				<html>
  				<head>
    				<title>200 OK</title>
  				</head>
  				<body>
    				<h1>Success!</h1>
    				<p>Your request was an absolute banger.</p>
  				</body>
				</html>
			`
		w.Header.Set("Content-Length", strconv.Itoa(len(resBody)))
		w.Header.Set("Content-Type", "text/html")
		w.WriteHeaders(w.Header)
		_, err = w.WriteBody([]byte(resBody))
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			return
		}
	}
}
