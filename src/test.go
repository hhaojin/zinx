package main

import (
	"fmt"
	"net"
)

func main() {

	conn, err := net.Dial("tcp", ":8999")
	if err != nil {
		panic(err)
	}
	var n int
	for {
		conn.Write([]byte("hello world"))
		buf := make([]byte, 1024)

		conn.Read(buf)
		fmt.Println(string(buf))
		n++
		if n == 3 {
			break
		}
	}

	conn.Close()

}
