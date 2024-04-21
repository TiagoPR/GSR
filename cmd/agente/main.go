package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"gsr"
	"net"
)

type object struct {
	structure int
	object    int
	indice    int // devolve a instancia no respetivo indicie
	limite    int // se existir devolve de todos entre indice e limite
}

func sendResponse(conn *net.UDPConn, addr *net.UDPAddr) {
	_, err := conn.WriteToUDP([]byte("From Agent: Hello I got your message "), addr)
	if err != nil {
		fmt.Printf("Couldn't send response %v", err)
	}
}

func main() {
	p := make([]byte, 2048)
	addr := net.UDPAddr{
		Port: 1234,
		IP:   net.ParseIP("127.0.0.1"),
	}
	ser, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}

	for {

		n, remoteaddr, err := ser.ReadFromUDP(p)
		if err != nil {
			fmt.Printf("Some error %v", err)
			continue
		}
		fmt.Printf("Read a message from %v \n", remoteaddr)

		receivedPDU := gsr.PDU{}

		dec := gob.NewDecoder(bytes.NewReader(p[:n])) // Will read from network.
		err = dec.Decode(&receivedPDU)
		if err != nil {
			// Error decoding message: unexpected EOF [ERROR HERE]
			fmt.Printf("Error decoding message: %v\n", err)
			continue
		}

		// Print the received PDU.
		receivedPDU.Print()
		// go sendResponse(ser, remoteaddr)
	}
}
