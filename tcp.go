package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"time"
)

type Server struct {
	d net.Dialer
	input chan []byte
	Port int
}

func WithPort(port int) func(s *Server) {
	return func(s *Server) {
		s.Port = port
	}
}

func NewServer(opts ...func(s *Server)) *Server {
	s := new(Server)

	// set default
	s.Port = 1234
	s.input = make(chan []byte, 10)
	
	for _, opt := range opts {
		opt(s)
	}
	
	return s
}

func (s *Server) Serve() {

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
	if err != nil {
		return
	}

	defer l.Close()
	
	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			return
		}

		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go func(c net.Conn) {
			defer c.Close()

			message, err := bufio.NewReader(c).ReadBytes('\n')
			if err != nil {
				return
			}

			// remove delimeter
			s.input <- message[:len(message)-1]
		}(conn)
	}
}

func (s *Server) Write(p []byte) (n int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	conn, err := s.d.DialContext(ctx, "tcp", fmt.Sprintf(":%d", s.Port))
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	
	// adding delimeter
	p = append(p, '\n')
	if _, err := conn.Write(p); err != nil {
		log.Fatal(err)
	}
	return
}

func (s *Server) Read(p []byte) (n int, err error) {
	msg := <- s.input
	input := bytes.NewBuffer(msg)
	n, err = input.Read(p)
	return
}