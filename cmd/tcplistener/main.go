package main

import (
	"fmt"
	"ja_httpserver/internal/request"
	"log"
	"net"
	"os"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("error", "error", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("error", "error", err)
		}
		r, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal("error", "error", err)
			os.Exit(1)
		}
		fmt.Printf("Request :\n")
		fmt.Printf("- Method: %s\n", r.RequestLine.Method)
		fmt.Printf("- Target: %s\n", r.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", r.RequestLine.HttpVersion)

	}

	// f, err := os.Open("messages.txt")
	// if err != nil {
	// 	log.Fatal("error", "error", err)
	// }

	// lines := getLinesChannel(f)
	// for line := range lines {
	// 	fmt.Printf("read: %s \n", line)

	// }

	// fmt.Println("I hope I get the job!")
}
