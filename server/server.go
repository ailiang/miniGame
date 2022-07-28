package main

import (
	"bufio"
	"io"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:9999")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go Process(conn)
	}
}

func Process(conn net.Conn) {
	defer conn.Close()
	for {
		var buf [128]byte
		r := bufio.NewReader(conn)
		n, err := r.Read(buf[:])
		if err != nil {
			if err == io.EOF {
				log.Print("normal end")
			} else {
				log.Print("excpet end")
			}
			break
		}
		log.Print("rev msg", string(buf[:n]))
	}
}
