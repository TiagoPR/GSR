package snmp

import (
	"bufio"
	"fmt"
	"net"
)

type pdu struct {
	tag               string
	tipo              byte
	timestamp         string
	messageIdentifier string
	iid_list          []string
	value_list        []string
	error_list        []string
}

func SetupGestor() {
	p := make([]byte, 2048)
	conn, err := net.Dial("udp", "127.0.0.1:1234")
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
	fmt.Fprintf(conn, "Olá agente, testing")
	_, err = bufio.NewReader(conn).Read(p)
	if err == nil {
		fmt.Printf("%s\n", p)
	} else {
		fmt.Printf("Some error %v\n", err)
	}
	conn.Close()
}
