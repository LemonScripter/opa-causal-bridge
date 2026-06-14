package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
)

/*
 * DCC Mock Service
 * Simulates the BioOS DCC kernel bridge for local testing.
 * Listens on a unix socket and returns 'Verified' (0x01) for valid IDs.
 */

const SocketPath = "dcc_test.sock"

func main() {
	if err := os.RemoveAll(SocketPath); err != nil {
		fmt.Printf("Error clearing socket: %v\n", err)
		return
	}

	l, err := net.Listen("unix", SocketPath)
	if err != nil {
		fmt.Printf("Listen error: %v\n", err)
		return
	}
	defer l.Close()

	fmt.Printf("DCC Mock Service listening on %s\n", SocketPath)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		os.Remove(SocketPath)
		os.Exit(0)
	}()

	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}

		go func(c net.Conn) {
			defer c.Close()
			buf := make([]byte, 1024)
			n, err := c.Read(buf)
			if err != nil {
				return
			}

			requestID := string(buf[:n])
			fmt.Printf("Verifying Request: %s\n", requestID)

			// Simple logic: IDs starting with 'VALID-' are verified
			if len(requestID) > 6 && requestID[:6] == "VALID-" {
				c.Write([]byte{0x01})
			} else {
				c.Write([]byte{0x00})
			}
		}(conn)
	}
}
