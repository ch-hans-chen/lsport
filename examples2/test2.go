
package main

import (
	"flag"
	"github.com/ch-hans-chen/lsport"
)

var (
	address  string
	baudrate int
	databits int
	stopbits int
	parity   string
	message string
)

func main() {
	flag.StringVar(&address, "a", "/dev/ttyS1", "address")
	flag.IntVar(&baudrate, "b", 9600, "baud rate")
	flag.IntVar(&databits, "d", 8, "data bits")
	flag.IntVar(&stopbits, "s", 1, "stop bits")
	flag.StringVar(&parity, "p", "N", "parity (N/E/O)")
	flag.StringVar(&message, "m", "lsport", "message")
	flag.Parse()

	s := lsport.Conf{}
	lsport.Init(&s, address)
	lsport.SetParams(&s, baudrate, databits, stopbits, parity)
	lsport.Write(s.Port, []byte(message))
	lsport.Close(&s)
	return
}
