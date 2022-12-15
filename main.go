package main

import (
	"fmt"

	"github.com/fluent/go-wasm-sockets/wasi/socket"
)

func main() {
	sockfd, err := socket.Open(socket.AF_INET, socket.SOCK_STREAM)
	if err != nil {
		fmt.Printf("OPEN: %s\n", err.Error())
		return
	}
	err = socket.Connect(sockfd, socket.SocketAddressInet{Address: "127.0.0.1", Port: 9999})
	if err != nil {
		fmt.Printf("CONNECT: %s\n", err.Error())
		return
	}
	buf := make([]byte, 4096)
	recv, err := socket.Recv(sockfd, buf)
	if err != nil {
		fmt.Printf("RECV: %s\n", err.Error())
		return
	}
	fmt.Printf("RECV[%d] = %s\n", recv, buf[0:recv])
}
