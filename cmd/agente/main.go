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
		1: []string{"00:1B:44:11:3A:B7"},         // id list
		2: []string{"Lights & A/C Conditioning"}, // type list
		3: []int{30},                             // beaconRate list
		4: []int{2},                              // nSensors list
		5: []int{2},                              // nActuators list
		6: []string{time.Now().Format("2006-01-02 15:04:05")},
		7: []string{"10:15:30"}, // upTime list
		8: []string{ // lastTimeUpdated list
			time.Now().Add(-5 * time.Minute).Format("2006-01-02 15:04:05"),
		},
		9:  []int{1}, // operationalStatus list
		10: []int{0}, // reset list
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
	responsePDU := messages.PDU{
		Tag:               receivedPDU.Tag,
		Tipo:              'R',
		Timestamp:         types.NewInfoTimestamp(),
		MessageIdentifier: "agente",
		Iid_list:          receivedPDU.Iid_list,
	}
	var valueElements []types.Tipo
	var errorElements []types.Tipo

	for _, iid := range receivedPDU.Iid_list.Elements {
		structure := iid.Value.Structure
		object := iid.Value.Objecto
		firstIndex := iid.Value.First_index
		secondIndex := iid.Value.Second_index

		if group, ok := mockLMIB[structure].(map[int]interface{}); ok {
			if list, ok := group[object]; ok {
				listValue := reflect.ValueOf(list)
				if listValue.Kind() != reflect.Slice {
					errorElements = append(errorElements, types.Tipo{
						Data_Type: 'S',
						Length:    1,
						Value:     "Object is not a list",
					})
					continue
				}

				// Handle different IID patterns
				if firstIndex == 0 {
					// Case x.y: Return all values
					for j := 0; j < listValue.Len(); j++ {
						value := listValue.Index(j).Interface()
						if tipo := processValue(value); tipo != nil {
							valueElements = append(valueElements, *tipo)
						}
					}
				} else {
					// Cases x.y.z or x.y.z.w
					idx := firstIndex - 1 // Convert to 0-based index
					if idx >= listValue.Len() {
						errMsg := fmt.Sprintf("Index %d out of range for object %d.%d", firstIndex, structure, object)
						errorElements = append(errorElements, types.Tipo{
							Data_Type: 'S',
							Length:    1,
							Value:     errMsg,
						})
						continue
					}

					if secondIndex == 0 {
						// Case x.y.z: Return single value
						value := listValue.Index(idx).Interface()
						if tipo := processValue(value); tipo != nil {
							valueElements = append(valueElements, *tipo)
						}
					} else {
						// Case x.y.z.w: Handle range request
						secondIdx := secondIndex - 1 // Convert to 0-based index

						// Always add the first value if it exists
						value := listValue.Index(idx).Interface()
						if tipo := processValue(value); tipo != nil {
							valueElements = append(valueElements, *tipo)
						}

						// Add error if second index is out of range
						if secondIdx >= listValue.Len() {
							errMsg := fmt.Sprintf("Index %d out of range for object %d.%d", secondIndex, structure, object)
							errorElements = append(errorElements, types.Tipo{
								Data_Type: 'S',
								Length:    1,
								Value:     errMsg,
							})
							continue
						}

						// Add remaining values if second index is valid
						for j := idx + 1; j <= secondIdx; j++ {
							value := listValue.Index(j).Interface()
							if tipo := processValue(value); tipo != nil {
								valueElements = append(valueElements, *tipo)
							}
						}
					}
				}
			} else {
				errMsg := fmt.Sprintf("Object %d not found in structure %d", object, structure)
				errorElements = append(errorElements, types.Tipo{
					Data_Type: 'S',
					Length:    1,
					Value:     errMsg,
				})
			}
		} else {
			errMsg := fmt.Sprintf("Structure %d not found", structure)
			errorElements = append(errorElements, types.Tipo{
				Data_Type: 'S',
				Length:    1,
				Value:     errMsg,
			})
		}
	}

	responsePDU.Value_list = types.Lists{
		N_Elements: len(valueElements),
		Elements:   valueElements,
	}
	responsePDU.Error_list = types.Lists{
		N_Elements: len(errorElements),
		Elements:   errorElements,
	}

	serializedPDU := responsePDU.SerializePDU()
	fmt.Printf("Serialized PDU: %s\n", serializedPDU)

	if _, err := conn.WriteToUDP([]byte(serializedPDU), addr); err != nil {
		fmt.Printf("Error sending response: %v\n", err)
	} else {
		fmt.Println("Response PDU sent successfully")
	}
}

// Helper function to process values and return appropriate Tipo
func processValue(value interface{}) *types.Tipo {
	var tipo types.Tipo
	switch v := value.(type) {
	case int:
		tipo = types.Tipo{
			Data_Type: 'I',
			Length:    1,
			Value:     strconv.Itoa(v),
		}
	case string:
		tipo = types.Tipo{
			Data_Type: 'S',
			Length:    1,
			Value:     v,
		}
	case time.Time:
		timeStr := v.Format(time.RFC3339)
		tipo = types.Tipo{
			Data_Type: 'T',
			Length:    len(timeStr),
			Value:     timeStr,
		}
	default:
		return nil
	}
	return &tipo
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
	pdu.Print()

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
