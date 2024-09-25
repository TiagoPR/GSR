package main

import (
	"fmt"
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

	// Print the received serialized PDU string
	fmt.Println("Received serialized PDU from manager:")
	fmt.Println(string(buf[:n]))
}

func main() {
	var wg sync.WaitGroup
	if len(os.Args) < 2 {
		panic("Introduce the agent's IP")
	}

	ip := os.Args[1]

	addr := net.UDPAddr{
		Port: 161,
		IP:   net.ParseIP(ip),
	}

	gestorAddr := net.UDPAddr{
		Port: 162,
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
