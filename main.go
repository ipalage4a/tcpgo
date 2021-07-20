package main

import (
	"bytes"
	"log"
)


func main() {
	server := NewServer(WithPort(8888))
	go server.Serve()
	log.Printf("Server started at ::%d\n", server.Port)
	for {
		var b  = make([]byte, 100, 1500)
		if _, err := server.Read(b); err != nil {
			panic(err)
		}

		str := string(bytes.Trim(b, "\x00"))
		log.Println("msg received:", str)
	}
}