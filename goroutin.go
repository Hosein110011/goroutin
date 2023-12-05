package main

import (
	"fmt"
	"time"
)


type Message struct {
	FROM string 
	Payload string
}

type Server struct {
	msgch chan Message
	quitch chan struct{}
}

func (s *Server) StartAndListen() {
free:
	for {
		select {
		case msg := <- s.msgch:
			fmt.Printf("received message from : %s payload %s\n", msg.FROM, msg.Payload)
		case <- s.quitch:
			fmt.Println("the server is doing a gracefull shutdown")
			break free
		default:
		}
	}
}

func sendMessageToServer(msgch chan Message, payload string) {
	msg := Message{
		FROM: "you",
		Payload: payload,
	}

	msgch <- msg
}

func graceFullQuitServer(quitch chan struct{}) {
	close(quitch)
}

func main() {
	s := &Server{
		msgch: make(chan Message),
		quitch: make(chan struct{}),
	}

	go s.StartAndListen()
	
	// for i := 0; i < 1000; i++ {
	go func() {
		time.Sleep(2 * time.Second)
		sendMessageToServer(s.msgch, "Hello")
	}()
	
	go func() {
		time.Sleep(4 * time.Second)
		graceFullQuitServer(s.quitch)
	} ()

	select {}
}

