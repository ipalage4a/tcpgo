package main

import (
	"bytes"
	"testing"
)

func Test_Tcp(t *testing.T) {
	server := NewServer()
	go server.Serve()
	
	for _, in := range []string{
		"1111", "2222",
	} {
		if _, err := server.Write([]byte(in)); err != nil {
			t.Fatal(err)
		}

		var b  = make([]byte, 100, 1500)
		if _, err := server.Read(b); err != nil {
			t.Fatal(err)
		}

		str := string(bytes.Trim(b, "\x00"))
		if str != in {
			t.Errorf("sent data not equal response: %s != %s", string(b), in )
		}
	}
}