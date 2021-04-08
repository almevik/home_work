package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type client struct {
	addr    string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &client{
		addr:    address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (c *client) Connect() error {
	var err error

	c.conn, err = net.DialTimeout("tcp", c.addr, c.timeout)
	if err != nil {
		return fmt.Errorf("ошибка подключения: %w", err)
	}

	return nil
}

func (c *client) Close() error {
	return c.conn.Close()
}

func (c *client) Send() error {
	_, err := io.Copy(c.conn, c.in)
	if err != nil {
		return fmt.Errorf("ошибка отправки: %w", err)
	}

	log.Println("EOF")

	return nil
}

func (c *client) Receive() error {
	_, err := io.Copy(c.out, c.conn)
	if err != nil {
		return fmt.Errorf("ошибка при получении: %w", err)
	}
	fmt.Fprintln(os.Stderr, "соединение закрыто")

	return nil
}
