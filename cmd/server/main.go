package main

import (
	"fmt"
	"os"
	"os/signal"
	"rcsp/internal/server"
	"syscall"
)

func HandleInterrupt(s *server.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT)
	go func() {
		<-c
		fmt.Printf("\r")
		s.Down()
		os.Exit(0)
	}()
}

func main() {
	server, err := server.NewServer("config.json")
	if err != nil {
		panic(err)
	}
	HandleInterrupt(server)
	if err := server.Up(); err != nil {
		panic(err)
	}
}
