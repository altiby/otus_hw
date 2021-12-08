package main

import (
	"context"
	"fmt"
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
	ctx        context.Context
	cancelFunc context.CancelFunc
	address    string
	timeout    time.Duration
	in         io.ReadCloser
	out        io.Writer
	conn       net.Conn
}

func (t *telnetClient) Connect() error {
	timoutContext, cancelFunc := context.WithTimeout(t.ctx, t.timeout)
	dialer := net.Dialer{}
	conn, err := dialer.DialContext(timoutContext, "tcp", t.address)
	if err != nil {
		defer cancelFunc()
		return fmt.Errorf("coonect %w", err)
	}
	t.conn = conn
	t.cancelFunc = cancelFunc
	return nil
}

func (t *telnetClient) Close() error {
	if t.conn == nil {
		return nil
	}
	t.cancelFunc()
	return t.conn.Close()
}

func (t *telnetClient) Send() error {
	if _, err := io.Copy(t.conn, t.in); err != nil {
		return err
	}
	return nil
}

func (t *telnetClient) Receive() error {
	if _, err := io.Copy(t.out, t.conn); err != nil {
		return err
	}
	return nil
}

func NewTelnetClient(ctx context.Context, address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{ctx: ctx, address: address, timeout: timeout, in: in, out: out}
}

// Place your code here.
// P.S. Author's solution takes no more than 50 lines.
