package main

import (
	"fmt"
	"gsr/messages"
	"gsr/messages/types"
	"net"
	"os"
	"reflect"
	"strconv"
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

func sendResponse(receivedPDU messages.PDU, conn *net.UDPConn, addr *net.UDPAddr) {
	responsePDU := receivedPDU
	var valueElements []types.Tipo

	// Process each OID in the IIDList
	for _, iid := range receivedPDU.Iid_list.Elements {
		structure := iid.Value.Structure
		object := iid.Value.Objecto
		firstIndex := iid.Value.First_index - 1 // Adjust for 0-based index
		secondIndex := -1                       // Default to -1 if not provided

		// Check if second index is provided
		if iid.Value.Second_index != 0 {
			secondIndex = iid.Value.Second_index - 1 // Adjust for 0-based index
		}

		fmt.Printf("\nProcessing OID: %d.%d.%d", structure, object, firstIndex+1)
		if secondIndex != -1 {
			fmt.Printf(".%d", secondIndex+1)
		}
		fmt.Println()

		if group, ok := mockLMIB[structure].(map[int]interface{}); ok {
			if list, ok := group[object]; ok {
				// Use reflection to handle different types of slices
				listValue := reflect.ValueOf(list)
				if listValue.Kind() != reflect.Slice {
					fmt.Printf("Object %d is not a list in structure %d\n", object, structure)
					continue
				}

				start := firstIndex
				end := listValue.Len()
				if secondIndex != -1 && secondIndex < end {
					end = secondIndex + 1
				}

				if start < 0 || start >= listValue.Len() {
					fmt.Printf("Invalid index range. List length: %d, Requested range: %d to %d\n", listValue.Len(), start, end-1)
					continue
				}

				fmt.Println("Retrieved values:")
				for i := start; i < end && i < listValue.Len(); i++ {
					value := listValue.Index(i).Interface()
					fmt.Printf("Index %d: %v (Type: %T)\n", i+1, value, value)
				}

				for j := start; j < end && j < listValue.Len(); j++ {
					value := listValue.Index(j).Interface()
					var tipo types.Tipo

					switch v := value.(type) {
					case int:
						tipo = types.Tipo{
							Data_Type: 'I',
							Length:    4, // Assuming 32-bit integer
							Value:     strconv.Itoa(v),
						}
					case string:
						tipo = types.Tipo{
							Data_Type: 'S',
							Length:    len(v),
							Value:     v,
						}
					case time.Time:
						tipo = types.Tipo{
							Data_Type: 'T',
							Length:    8, // Assuming 64-bit timestamp
							Value:     v.Format(time.RFC3339),
						}
					default:
						fmt.Printf("Unsupported type for value: %v\n", v)
						continue
					}

					valueElements = append(valueElements, tipo)
				}
			} else {
				fmt.Printf("Object %d doesn't exist in structure %d\n", object, structure)
			}
		} else {
			fmt.Printf("Structure %d doesn't exist in mockLMIB\n", structure)
		}
	}

	responsePDU.Timestamp = types.NewRequestTimestamp()
	responsePDU.Tipo = 'R'
	responsePDU.MessageIdentifier = "agente"
	responsePDU.Value_list = types.NewLists(len(valueElements), valueElements)

	serializedPDU := responsePDU.SerializePDU()
	_, err := conn.WriteToUDP([]byte(serializedPDU), addr)
	if err != nil {
		fmt.Printf("Error sending response: %v\n", err)
	} else {
		fmt.Println("Response PDU sent successfully")
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
	fmt.Println(pdu)

	if pdu.Tipo == 'G' {
		sendResponse(pdu, ser, remoteaddr)
	}

	readPDU(ser)
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
