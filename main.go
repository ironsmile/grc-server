package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/ironsmile/grc-server/connpool"
)

type Message struct {
	Data string
	From *net.Conn
}

func main() {
	pool := connpool.New()

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Printf("Was not able to bind to port: %s", err)
		os.Exit(1)
	}

	output := make(chan *Message, 100)
	defer close(output)
	go writer(pool, output)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Was not able to accept connection: %s", err)
			continue
		}
		pool.AddConnection(&conn)
		go clientRoutine(&conn, output)
	}
}

func clientRoutine(conn *net.Conn, output chan<- *Message) {
	defer (*conn).Close()
	scaner := bufio.NewScanner(*conn)
	for scaner.Scan() {
		line := scaner.Text()
		log.Printf("Writing to the pool: %s", line)
		output <- &Message{Data: line, From: conn}
	}
}

func writer(pool *connpool.ConnPool, input <-chan *Message) {
	for msg := range input {
		line := msg.Data
		pool.MessageFrom(msg.From)
		_, err := fmt.Fprint(pool, line)
		if err != nil {
			log.Printf("Writing to the connection pool failed: %s", err)
		}
	}
}
