package main

import (
	//"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"gsr"
	"net"
)

func main() {
	//p := make([]byte, 2048)
	conn, err := net.Dial("udp", "127.0.0.1:1234")
	if err != nil {
		fmt.Printf("Couldn't create connection %v", err)
		return
	}

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
	//_, err = bufio.NewReader(conn).Read(p)
	//if err == nil {
	//	fmt.Printf("%s\n", p)
	//} else {
	//	fmt.Printf("Some error %v\n", err)
	//}
	conn.Close()
}
