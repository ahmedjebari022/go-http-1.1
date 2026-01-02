package server

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync/atomic"

	header "github.com/ahmedjebari022/go-http-1.1/internal/headers"
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

type Handler func(w *response.Writer, req *request.Request)

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
		go s.proxyHandle(conn)
	}
}

func (s *Server) proxyHandle(conn net.Conn){
	defer conn.Close()
	req, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Printf("Error while parsin request: %s", err.Error())
		return
	}
	requestLine := req.RequestLine.RequestTarget
	if strings.HasPrefix(requestLine, "/httpbin/"){
		n := strings.TrimPrefix(requestLine, "/httpbin/")
		res, err := http.Get("https://httpbin.org/" + n )
		if err != nil {
			fmt.Printf("Error: %s\n",err.Error())
			return
		}
		defer res.Body.Close()
		writer := response.NewWriter(conn)
		writer.Header.Set("Transfer-Encoding", "chunked")
		writer.Header.Set("Trailer", "X-Content-Sha256, X-Content-Length")
		delete(writer.Header, "Content-Length")
		err = writer.WriteStatusLine(response.Success)
		if err != nil {
			return
		}
		err = writer.WriteHeaders(writer.Header)
		p := make([]byte, 1024)
		var body []byte
		len := 0
		for {
			nr, err := res.Body.Read(p)
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}
			fmt.Printf("n: %d\n", nr)
			_, err = writer.WriteChunkedBody(p[:nr])
			if err != nil {
				fmt.Printf("Error: %s", err.Error())
				return
			}
			len += nr
			body = append(body, p[:nr]...)
			p = make([]byte, 1024)
		}
		_, err = writer.WriteChunkedBodyDone()
		if err != nil {
			fmt.Printf("Error: %s", err.Error())
			return
		}

		trailer := header.NewHeaders()
		sum := sha256.Sum256(body)
		trailer["X-Content-Sha256"] = hex.EncodeToString(sum[:])
		trailer["X-Content-Length"] = strconv.Itoa(len)
		err = writer.WriteTrailers(trailer)
		if err != nil {
			fmt.Printf("Error: %s\n",err.Error())
			return
		}
	}else if requestLine == "/video"{
		writer := response.NewWriter(conn)
		delete(writer.Header, "Content-Length")	
		writer.Header.Set("Content-Type","video/mp4")	
		writer.WriteStatusLine(response.Success)
		writer.WriteHeaders(writer.Header)
		video, err := os.ReadFile("./assets/vim.mp4")
		if err != nil {
			return 
		}
		writer.WriteBody(video)
	}
}


func (s *Server) handle(conn net.Conn){
	defer conn.Close()
	req, err := request.RequestFromReader(conn)	
	if err != nil {
		fmt.Printf("Error while parsin request: %s", err.Error())
		return 
	}
	writer := response.NewWriter(conn)	
	s.handler(&writer, req)
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
