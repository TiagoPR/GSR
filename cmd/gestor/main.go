package main

import (
	"bufio"
	"fmt"
	"gsr/messages"
	"gsr/messages/types"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
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

func readPDU(ser *net.UDPConn) messages.PDU {
	buf := make([]byte, 2048)
	n, remoteaddr, err := ser.ReadFromUDP(buf)
	if err != nil {
		fmt.Printf("Some error %v", err)
	}
	fmt.Printf("Read a message from %v \n", remoteaddr)

	// Print the received serialized PDU string
	fmt.Println("Received serialized PDU from agent:")
	serializedPdu := string(buf[:n])
	fmt.Println(serializedPdu)
	// Deserializing not working here!
	pdu := messages.DeserializePDU(serializedPdu)
	return pdu
}

func getRequest() messages.PDU {
	time := types.NewRequestTimestamp()
	fmt.Println(time)
	messageIdentifier := "gestor"
	fmt.Println(messageIdentifier)

	fmt.Println("Which IID?")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter text: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	// Remove newline character
	text = strings.Replace(text, "\n", "", -1)

	parts := strings.Split(text, ".")

	var structure, object, firstIndex, secondIndex int

	structure, _ = strconv.Atoi(parts[0])
	object, _ = strconv.Atoi(parts[1])

	if len(parts) > 2 {
		firstIndex, _ = strconv.Atoi(parts[2])
	}

	if len(parts) > 3 {
		secondIndex, _ = strconv.Atoi(parts[3])
	}

	fmt.Println("Structure: ", structure)
	fmt.Println("Object: ", object)
	fmt.Println("First Index: ", firstIndex)
	fmt.Println("Second Index: ", secondIndex)

	iid := types.NewIID(structure, object, firstIndex, secondIndex)
	iid_tipo := types.NewIID_Tipo(4, iid)
	iid_list := types.NewIID_List(1, []types.IID_Tipo{iid_tipo})

	pdu := messages.NewPDU('G', time, messageIdentifier, iid_list, types.Lists{}, types.Lists{})

	return pdu
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

	line = strings.TrimSpace(line) // Trim the input string

	raddr, err := net.ResolveUDPAddr("udp", line)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		// Couldn't create connection dial udp: lookup udp/1053: unknown port [ERROR HERE]
		fmt.Printf("Couldn't create connection %v", err)
		return
	}

	defer conn.Close()

	pdu := getRequest()
	serializedPDU := pdu.SerializePDU()

	_, err = conn.Write([]byte(serializedPDU))
	if err != nil {
		fmt.Println("Error sending data:", err)
		return
	}

	fmt.Println("Sending PDU")
	receivedPDU := readPDU(conn)
	receivedPDU.Print()
}

func main() {
	//p := make([]byte, 2048)
	if len(os.Args) < 2 {
		panic("Introduce the gestor's IP")
	}

	ip := os.Args[1]

	// listen to incoming udp packets
	pc, err := net.ListenPacket("udp", ip+":162")
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()

	go listenAgents(pc)

	time.Sleep(5)

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
