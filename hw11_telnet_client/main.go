package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

var errLog = log.New(os.Stderr, "", 0)

func main() {
	timeoutFlag := flag.String("timeout", "10s", "таймаут подключения к серверу")
	timeout, err := time.ParseDuration(*timeoutFlag)

	flag.Parse()

	if err != nil {
		log.Fatalf("can't to parse flag: %v", err)
	}

	if flag.NArg() < 2 {
		log.Fatal("некорректные аргументы: укажите адрес и порт")
	}

	host := flag.Arg(0)
	port := flag.Arg(1)
	addr := net.JoinHostPort(host, port)

	log.Printf("подключение к %s\n", addr)

	ctx, cancel := context.WithCancel(context.Background())

	client := NewTelnetClient(addr, timeout, os.Stdin, os.Stdout, cancel)
	if err := client.Connect(); err != nil {
		log.Fatalf("не удалось подключиться к %v, %v", addr, err)
	}
	defer client.Close()

	go receive(client)
	go send(client)

	chSig := make(chan os.Signal, 1)
	signal.Notify(chSig, os.Interrupt)

	select {
	case <-chSig:
		cancel()
	case <-ctx.Done():
		close(chSig)
	}
}

func send(client TelnetClient) {
	if err := client.Send(); err != nil {
		errLog.Printf("не удалось отправить: %v", err)
		return
	}
}

func receive(client TelnetClient) {
	if err := client.Receive(); err != nil {
		errLog.Printf("не удалось прочитать: %v", err)
		return
	}
}
