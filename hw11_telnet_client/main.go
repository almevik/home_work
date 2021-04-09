package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

const fTimeout = "timeout"

func main() {
	timeout := flag.Duration(fTimeout, 10*time.Second, "таймаут подключения к серверу")

	flag.Parse()
	if flag.NArg() < 2 {
		log.Fatal("некорректные аргументы: укажите адрес и порт")
	}

	host := flag.Arg(0)
	port := flag.Arg(1)
	addr := net.JoinHostPort(host, port)

	log.Printf("подключение к %s\n", addr)

	ctx, cancel := context.WithCancel(context.Background())

	client := NewTelnetClient(addr, *timeout, os.Stdin, os.Stdout, cancel)
	if err := client.Connect(); err != nil {
		log.Fatalf("не удалось подключиться к %v, %v", addr, err)
	}

	log.Printf("успешно подключен")
	defer client.Close()

	go listenStopSignal(ctx, cancel)
	go receive(client)
	go send(client)
}

// Слушатель сигнала остановки программы.
func listenStopSignal(ctx context.Context, cancel context.CancelFunc) {
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
		log.Println(fmt.Fprintln(os.Stderr, err))
		return
	}
}

func receive(client TelnetClient) {
	if err := client.Receive(); err != nil {
		log.Println(fmt.Fprintln(os.Stderr, err))
		return
	}
}
