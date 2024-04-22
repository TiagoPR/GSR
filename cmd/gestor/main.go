package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"gsr"
	"log"
	"net"
	"os"
)

var agents []string

func serve(pc net.PacketConn, addr net.Addr, buf []byte) {
	// 0 - 1: ID
	// 2: QR(1): Opcode(4)
	buf[2] |= 0x80 // Set QR bit

	pc.WriteTo(buf, addr)
}

func listenAgents(pc net.PacketConn) {
	buf := make([]byte, 1024)
	fmt.Println("Listening for agent's")
	_, addr, err := pc.ReadFrom(buf)
	agents = append(agents, addr.String())
	if err != nil {
		fmt.Println("Couldn't read for agent's")
	}
	println(agents[len(agents)-1])
	//go serve(pc, addr, buf[:n])

}

func send() {
	fmt.Println("Which agent you want to send a message")
	fmt.Printf("%v", agents)

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(line)

	conn, err := net.Dial("udp", line)
	if err != nil {
		// Couldn't create connection dial udp: lookup udp/1053: unknown port [ERROR HERE]
		fmt.Printf("Couldn't create connection %v", err)
		return
	}

	defer conn.Close()

	var network bytes.Buffer        // Stand-in for a network connection
	enc := gob.NewEncoder(&network) // Will write to network.
	// dec := gob.NewDecoder(&network) // Will read from network.

	pdu := gsr.NewPDU('G', "27/04/2000", "1", nil, nil, nil)
	err = enc.Encode(pdu)
	if err != nil {
		fmt.Printf("Couldn't encode data %v", err)
	}
	encodedData := network.Bytes() // Get the encoded data slice

	fmt.Println("Sending PDU")
	_, err = conn.Write(encodedData[:network.Len()]) // Send the encoded data
	if err != nil {
		fmt.Printf("Couldn't send data %v", err)
	}

}

func main() {
	//p := make([]byte, 2048)
	if len(os.Args) < 2 {
		panic("Introduce the gestor's IP")
	}

	ip := os.Args[1]

	// listen to incoming udp packets
	pc, err := net.ListenPacket("udp", ip+":1053")
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()

	go listenAgents(pc)

	for {
		send()
	}

	//_, err = bufio.NewReader(conn).Read(p)
	//if err == nil {
	//	fmt.Printf("%s\n", p)
	//} else {
	//	fmt.Printf("Some error %v\n", err)
	//}
}
