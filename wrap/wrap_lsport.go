/*
Package serial provides a wrapper for lsport to modbus.
*/
package serial

import (
	"errors"
	"fmt"
	"github.com/ch-hans-chen/lsport"
	"io"
	"time"
)

var (
	// ErrTimeout is occurred when timing out.
	ErrTimeout = errors.New("serial:timeout")
)

// Config is common configuration for serial port.
type Config struct {
	// Device path (/dev/ttyS0)
	Address string
	// Baud rate (default 19200)
	BaudRate int
	// Data bits: 5, 6, 7 or 8 (default 8)
	DataBits int
	// Stop bits: 1 or 2 (default 1)
	StopBits int
	// Parity: N - None, E - Even, O - Odd (default E)
	// (The use of no parity requires 2 stop bits.)
	Parity string
	// Read (Write) timeout.
	Timeout time.Duration
}

// Port is the interface for controlling serial port.
type Port interface {
	io.ReadWriteCloser
	// Connect connects to the serial port.
	Open(*Config) error
}

// port implements Port interface.
type port struct {
	fh *lsport.Conf
	timeout time.Duration
}

// New allocates and returns a new serial port controller.
func New() Port {
	return &port{}
}

// Open opens a serial port.
func Open(c *Config) (p Port, err error) {
	p = New()
	err = p.Open(c)
	return
}

func (p *port) Open(c *Config) (err error) {
	s := lsport.Conf{}
	_, err = lsport.Init(&s, c.Address)
	if err != nil {
		fmt.Printf("Open %s fail, Reason: %s\n", c.Address, err)
		return
	}
	_, err = lsport.SetParams(&s, c.BaudRate, c.DataBits, c.StopBits, c.Parity)
	if err != nil {
		fmt.Printf("Open %v fail, Reason: %s\n", c, err)
		return
	}

	p.fh = &s
	p.timeout = c.Timeout
	return
}

func (p *port) Close() (err error) {
	lsport.Close(p.fh)
	return nil
}

// blocking read
func (p *port) Read(b []byte) (n int, err error) {
	var nn int32
	var timeout_ms uint

	timeout_ms = uint(p.timeout) / uint(1*time.Millisecond)
	nn, err = lsport.BlockingRead(p.fh.Port, b, timeout_ms)
	if err != nil {
		//fmt.Printf("BlockingRead fail, Reason: %s\n", err)
		err = ErrTimeout
		return
	}
	n = int(nn)
	return
}

// blocking write
func (p *port) Write(b []byte) (n int, err error) {
	_, err = lsport.BlockingWrite(p.fh.Port, b, uint16(p.timeout))
	if err != nil {
		fmt.Printf("BlockingWrite fail, Reason: %s\n", err)
		err = ErrTimeout
		return
	}
	return
}
