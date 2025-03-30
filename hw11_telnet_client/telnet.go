package main

import (
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	address string
	timeout time.Duration
	conn    net.Conn
	in      io.ReadCloser
	out     io.Writer
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (c *telnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return err
	}
	c.conn = conn

	return nil
}

func (c *telnetClient) Send() error {
	_, err := io.Copy(c.conn, c.in)
	if err != nil {
		return err
	}

	return nil
}

func (c *telnetClient) Receive() error {
	_, err := io.Copy(c.out, c.conn)
	if err != nil {
		return err
	}

	return nil
}

func (c *telnetClient) Close() error {
	err := c.conn.Close()
	if err != nil {
		return err
	}

	return nil
}
