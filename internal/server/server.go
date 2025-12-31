package server

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"sync/atomic"

	"github.com/ahmedjebari022/go-http-1.1/internal/request"
	"github.com/ahmedjebari022/go-http-1.1/internal/response"
)
type State int
const (
	listeningState State = iota
	closedState 
)
type Server struct{
	state 		State
	listenr		net.Listener 	
	listening 	atomic.Bool
	handler 	Handler
}

type HandlerError struct {
	StatusCode 		response.StatusCode
	ErrorMessage 	string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError


func Serve(port int,handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d",port))															
	if err != nil {
		return nil, err	
	}
	server := &Server{
		listenr: listener,
		state: listeningState,
		handler: handler,
	}
	server.listening.Store(true)
	return server, nil
}

func (s *Server) Close(){
	s.listenr.Close()
	s.state = closedState
	s.listening.Store(false)
}



func (s *Server) Listen(){
	for {
		conn, err := s.listenr.Accept()
		if err != nil && s.listening.Load(){
			fmt.Printf("Error establishing connexion: %s\n", err.Error())
			break
		}
		fmt.Println("Connexion established successfully")
		go s.handle(conn)
	}
}


func (s *Server) handle(conn net.Conn){
	defer conn.Close()
	req, err := request.RequestFromReader(conn)	
	if err != nil {
		fmt.Printf("Error while parsin request: %s", err.Error())
		return 
	}
	var buff bytes.Buffer
	handlerError := s.handler(&buff, req)
	if handlerError != nil {
		handlerError.WriteError(conn)
		return
	}

	err = response.WriteStatusLine(conn, response.Success)	
	if err != nil {
		fmt.Printf("Error when Writing to response line: %s\n", err.Error())
		return
	}
	defaultHeader := response.GetDefaultHeaders(len(buff.Bytes()))
	err = response.WriteHeaders(conn, defaultHeader)
	if err != nil {
		fmt.Printf("Error when Writing to response headers: %s\n", err.Error())
		return
	}
	_, err = conn.Write(buff.Bytes())
	if err != nil {
		fmt.Printf("Error while writing body: %s\n",err.Error())
		return
	}
}

func (h *HandlerError) WriteError(w io.Writer) error {
	err := response.WriteStatusLine(w, h.StatusCode)
	if err != nil {
		return err
	}
	headers := response.GetDefaultHeaders(len(h.ErrorMessage))
	err = response.WriteHeaders(w, headers)
	if err != nil {
		return err
	} 
	_, err = w.Write([]byte(h.ErrorMessage))
	if err != nil {
		return err
	}
	return nil
}
