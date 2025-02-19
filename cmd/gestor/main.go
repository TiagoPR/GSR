package main

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gsr/messages"
	"gsr/messages/types"
	"gsr/security"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

var agents []string

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

	pdu := messages.DeserializePDU(serializedPdu)
	return pdu
}

func getRequest() messages.PDU {
	time := types.NewRequestTimestamp()
	b := make([]byte, 12)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Convert to base64
	randomString := base64.StdEncoding.EncodeToString(b)

	// Trim to exactly 16 characters
	randomString = strings.ReplaceAll(randomString, "/", "A")
	randomString = strings.ReplaceAll(randomString, "+", "B")
	randomString = randomString[:16]

	// This is the result - a random 16 character string
	messageIdentifier := randomString

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("How many IID do you want to send?")
	nIIDS, _ := reader.ReadString('\n')
	nIIDS = strings.TrimSpace(nIIDS)
	number, err := strconv.Atoi(nIIDS)

	if err != nil {
		fmt.Println("Wrong value, try again...")
		getRequest()
	}

	// Initialize the list as an empty slice
	iid_list := []types.IID_Tipo{}

	for i := 0; i < number; i++ {
		fmt.Println("Which IID?")
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		// Remove newline character
		text = strings.TrimSpace(text)

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

		// Create IID and IID_Tipo objects
		iid := types.NewIID(structure, object, firstIndex, secondIndex)
		iid_tipo := types.NewIID_Tipo(len(parts), iid)

		// Append the new IID_Tipo to the slice
		iid_list = append(iid_list, iid_tipo)
	}

	// If needed, create a new IID_List from the updated slice
	finalList := types.NewIID_List(len(iid_list), iid_list)

	pdu := messages.NewPDU('G', time, messageIdentifier, finalList, types.Lists{}, types.Lists{})
	fmt.Println("Get Request PDU:")
	pdu.Print()

	return pdu
}

func setRequest() messages.PDU {
	time := types.NewRequestTimestamp()
	b := make([]byte, 12)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Convert to base64
	randomString := base64.StdEncoding.EncodeToString(b)

	// Trim to exactly 16 characters
	randomString = strings.ReplaceAll(randomString, "/", "A")
	randomString = strings.ReplaceAll(randomString, "+", "B")
	randomString = randomString[:16]

	// This is the result - a random 16 character string
	messageIdentifier := randomString

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\nSET REQUEST")
	fmt.Println("Available read-write objects:")
	fmt.Println("1.3 - Beacon Rate")
	fmt.Println("1.6 - Date and Time")
	fmt.Println("1.10 - Reset")
	fmt.Println("3.3.X - Actuator Status")

	fmt.Println("\nHow many values do you want to set?")
	nIIDs, _ := reader.ReadString('\n')
	nIIDs = strings.TrimSpace(nIIDs)
	number, err := strconv.Atoi(nIIDs)
	if err != nil {
		fmt.Println("Wrong value, try again...")
		setRequest()
	}

	var iid_list []types.IID_Tipo
	var value_list []types.Tipo

	for i := 0; i < number; i++ {
		fmt.Printf("\nValue %d:\n", i+1)
		fmt.Println("Enter IID (format: x.y or x.y.z):")
		iidText, _ := reader.ReadString('\n')
		iidText = strings.TrimSpace(iidText)

		parts := strings.Split(iidText, ".")
		var structure, object, firstIndex, secondIndex int
		structure, _ = strconv.Atoi(parts[0])
		object, _ = strconv.Atoi(parts[1])
		if len(parts) > 2 {
			firstIndex, _ = strconv.Atoi(parts[2])
		}

		// Create and append IID
		iid := types.NewIID(structure, object, firstIndex, secondIndex)
		iid_tipo := types.NewIID_Tipo(len(parts), iid)
		iid_list = append(iid_list, iid_tipo)

		// Get value to set
		fmt.Println("Enter new value:")
		valueText, _ := reader.ReadString('\n')
		valueText = strings.TrimSpace(valueText)

		// Create value tipo based on the object type
		var tipo types.Tipo
		if structure == 1 && (object == 6) { // Date and time
			tipo = types.Tipo{
				Data_Type: 'S',
				Length:    1,
				Value:     valueText,
			}
		} else { // Numbers for other cases
			tipo = types.Tipo{
				Data_Type: 'I',
				Length:    1,
				Value:     valueText,
			}
		}
		value_list = append(value_list, tipo)
	}

	finalIIDList := types.NewIID_List(len(iid_list), iid_list)
	finalValueList := types.Lists{
		N_Elements: len(value_list),
		Elements:   value_list,
	}

	return messages.NewPDU('S', time, messageIdentifier, finalIIDList, finalValueList, types.Lists{})
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

	fmt.Println("Which message do you which to send?\n\t1 - GetRequest\n\t2 - SetRequest")

	reader = bufio.NewReader(os.Stdin)
	line, err = reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(line)

	line = strings.TrimSpace(line)

	var pdu messages.PDU

	line = strings.TrimSpace(line)

	if line == "1" {
		pdu = getRequest()
	} else if line == "2" {
		pdu = setRequest()
	} else {
		fmt.Println("Wrong value, try again...")
		send()
	}

	serializedPDU := pdu.SerializePDU()
	fmt.Println(serializedPDU)

	_, err = conn.Write([]byte(serializedPDU))
	if err != nil {
		fmt.Println("Error sending data:", err)
		return
	}

	fmt.Println("Sending PDU")
	receivedPDU := readPDU(conn)
	receivedPDU.Print()
}

func listenForMessages(pc net.PacketConn, config security.SecurityConfig) {
	for {
		buf := make([]byte, 2048)
		n, addr, err := pc.ReadFrom(buf)
		if err != nil {
			fmt.Printf("Error reading from connection: %v\n", err)
			continue
		}

		// Try to unmarshal as SecurePing first
		var securePing security.SecurePing
		if err := json.Unmarshal(buf[:n], &securePing); err == nil {
			// This is a ping message
			message := securePing.IP + securePing.Time
			if !security.VerifyHMAC(message, securePing.HMAC, config.SecretKey) {
				fmt.Printf("Invalid HMAC from %v\n", addr)
				continue
			}
			// Add to agents list only if HMAC is valid
			if !contains(agents, addr.String()) {
				agents = append(agents, addr.String())
				fmt.Printf("New authenticated agent connected from: %s\n", addr.String())
				fmt.Println("Which agent you want to send a message")
				fmt.Printf("%v", agents)
			}
			continue
		}

		// If not a ping, try to handle as PDU
		serializedPdu := string(buf[:n])
		pdu := messages.DeserializePDU(serializedPdu)
		if pdu.Tipo == 'N' { // Notification
			fmt.Printf("\nReceived Notification from %v\n", addr)
			pdu.Print()
			continue
		}

		// If we get here, the message format was not recognized
		fmt.Printf("Received unknown message format from %v\n", addr)
	}
}

func main() {
	if len(os.Args) < 3 {
		panic("Usage: program <gestor_IP> <secret_key>")
	}

	ip := os.Args[1]
	hexKey := os.Args[2]

	// Decode the hexadecimal key into raw bytes
	secretKey, err := hex.DecodeString(hexKey)
	if err != nil {
		log.Fatalf("Error decoding secret key: %v", err)
	}

	if len(secretKey) != security.HMAC_KEY_SIZE {
		log.Fatalf("Secret key must be 32 bytes long, but got %d bytes", len(secretKey))
	}

	config := security.SecurityConfig{
		SecretKey: secretKey,
	}
	// Listen for UDP packets
	pc, err := net.ListenPacket("udp", ip+":162")
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()

	// Start listener in a goroutine
	go listenForMessages(pc, config)

	// Handle user input in the main goroutine
	for {
		send()
	}
}

// Helper function to check if a string is in a slice
func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}
