package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"gsr"
	"net"
	"os"
	"sync"
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

func sendPing(conn *net.UDPConn, addr *net.UDPAddr, ip string) {
	_, err := conn.WriteToUDP([]byte(ip), addr)
	if err != nil {
		fmt.Printf("Couldn't send ping %v", err)
	}
}

func readPDU(ser *net.UDPConn) {
	buf := make([]byte, 2048)
	n, remoteaddr, err := ser.ReadFromUDP(buf)
	if err != nil {
		fmt.Printf("Some error %v", err)
	}
	fmt.Printf("Read a message from %v \n", remoteaddr)

	receivedPDU := gsr.PDU{}

	dec := gob.NewDecoder(bytes.NewReader(buf[:n])) // Will read from network.
	err = dec.Decode(&receivedPDU)
	if err != nil {
		// Error decoding message: unexpected EOF [ERROR HERE]
		fmt.Printf("Error decoding message: %v\n", err)
	}

	// Print the received PDU.
	receivedPDU.Print()
}

func main() {
	var wg sync.WaitGroup
	if len(os.Args) < 2 {
		panic("Introduce the agent's IP")
	}

	ip := os.Args[1]

	addr := net.UDPAddr{
		Port: 1053,
		IP:   net.ParseIP(ip),
	}

	gestorAddr := net.UDPAddr{
		Port: 1053,
		IP:   net.ParseIP("127.0.0.1"),
	}

	ser, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}
	fmt.Println("Sending ping to gestor")
	// gestor needs to identify the agent
	sendPing(ser, &gestorAddr, ip)

	wg.Add(1)
	go readPDU(ser)
	wg.Wait()
}
