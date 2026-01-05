package main

import (
	"flag"
	"fmt"
	"net"
	"time"
)

var port = flag.Int("port", 8080, "port number")

func main() {

	flag.Parse()

	if *port <= 0 {
		fmt.Println("invalid port number")
		return
	}
	addr := fmt.Sprintf("%s:%d", "localhost", *port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("error listening to port %d: %v\n", *port, err)
		return
	}
	fmt.Printf("server started listening at :%d\n", *port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("connection error: ", err)
			return
		}
		go handleConnection(conn)
	}

}
func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("client connected from: %v\n", conn.LocalAddr())
	for {
		var data string 
		if conn.LocalAddr().String() == "127.0.0.1:9010" {
			data = time.Now().Add(time.Minute).Format(time.RFC1123)			
		} else {
			data = time.Now().Format(time.RFC1123)
		}
		fmt.Fprintf(conn, "%s\n", data)
		time.Sleep(time.Second)

	}
}
