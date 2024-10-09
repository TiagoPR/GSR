package main

import (
	"fmt"
	"gsr/messages"
	"net"
	"os"
	"sync"
	"time"
)

var mockLMIB = map[int]interface{}{
	1: map[int]interface{}{ // device Group
		1: []string{"00:1B:44:11:3A:B7", "00:1B:44:11:3A:C8"},       // id list
		2: []string{"Lights & A/C Conditioning", "Security System"}, // type list
		3: []int{30, 60},                                            // beaconRate list
		4: []int{2, 3},                                              // nSensors list
		5: []int{2, 2},                                              // nActuators list
		6: []string{ // dateAndTime list
			time.Now().Format("2006-01-02 15:04:05"),
			time.Now().Add(-1 * time.Hour).Format("2006-01-02 15:04:05"),
		},
		7: []string{"10:15:30", "05:30:15"}, // upTime list
		8: []string{ // lastTimeUpdated list
			time.Now().Add(-5 * time.Minute).Format("2006-01-02 15:04:05"),
			time.Now().Add(-10 * time.Minute).Format("2006-01-02 15:04:05"),
		},
		9:  []int{1, 1}, // operationalStatus list
		10: []int{0, 0}, // reset list
	},
	2: map[int]interface{}{ // sensors Table
		1: []string{"00:1B:44:11:3A:B8", "00:1B:44:11:3A:B9", "00:1B:44:11:3A:D1", "00:1B:44:11:3A:D2", "00:1B:44:11:3A:D3"}, // id list
		2: []string{"Light", "Temperature", "Motion", "Door", "Window"},                                                      // type list
		3: []int{75, 22, 0, 1, 0},                                                                                            // status list
		4: []int{0, -10, 0, 0, 0},                                                                                            // minValue list
		5: []int{100, 40, 1, 1, 1},                                                                                           // maxValue list
		6: []string{ // lastSamplingTime list
			time.Now().Add(-30 * time.Second).Format("2006-01-02 15:04:05"),
			time.Now().Add(-15 * time.Second).Format("2006-01-02 15:04:05"),
			time.Now().Add(-45 * time.Second).Format("2006-01-02 15:04:05"),
			time.Now().Add(-30 * time.Second).Format("2006-01-02 15:04:05"),
			time.Now().Add(-15 * time.Second).Format("2006-01-02 15:04:05"),
		},
	},
	3: map[int]interface{}{ // actuators Table
		1: []string{"00:1B:44:11:3A:C0", "00:1B:44:11:3A:C1", "00:1B:44:11:3A:E0", "00:1B:44:11:3A:E1"}, // id list
		2: []string{"Light Switch", "Temperature Control", "Alarm", "Door Lock"},                        // type list
		3: []int{1, 22, 0, 1},                                                                           // status list
		4: []int{0, 16, 0, 0},                                                                           // minValue list
		5: []int{1, 30, 1, 1},                                                                           // maxValue list
		6: []string{ // lastControlTime list
			time.Now().Add(-2 * time.Minute).Format("2006-01-02 15:04:05"),
			time.Now().Add(-2*time.Minute + 30*time.Second).Format("2006-01-02 15:04:05"),
			time.Now().Add(-3 * time.Minute).Format("2006-01-02 15:04:05"),
			time.Now().Add(-3*time.Minute + 30*time.Second).Format("2006-01-02 15:04:05"),
		},
	},
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
	serializedPdu := string(buf[:n])
	fmt.Println(serializedPdu)
	pdu := messages.DeserializePDU(serializedPdu)
	oid := []int{pdu.Iid_list.Elements[0].Value.Structure, pdu.Iid_list.Elements[0].Value.Objecto, pdu.Iid_list.Elements[0].Value.First_index - 1}
	value := mockLMIB[oid[0]].(map[int]interface{})[oid[1]].([]int)[oid[2]]
	print("Value: ", value)
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
