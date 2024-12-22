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
		1: []string{"00:1B:44:11:3A:B7"},                      // id read-only
		2: []string{"Lights & A/C Conditioning"},              // type read-only
		3: []int{30},                                          // beaconRate read-write
		4: []int{2},                                           // nSensors read-only
		5: []int{2},                                           // nActuators read-only
		6: []string{time.Now().Format("2006-01-02 15:04:05")}, // dateAndTime read-write
		7: []string{"10:15:30"},                               // upTime read-only
		8: []string{ // lastTimeUpdated list read-only
			time.Now().Add(-5 * time.Minute).Format("2006-01-02 15:04:05"),
		},
		9:  []int{1}, // operationalStatus read-only
		10: []int{0}, // reset read-write
	},
	2: map[int]interface{}{ // sensors Table
		1: []string{"00:1B:44:11:3A:B8", "00:1B:44:11:3A:B9", "00:1B:44:11:3A:D1", "00:1B:44:11:3A:D2", "00:1B:44:11:3A:D3"}, // id list read-only
		2: []string{"Light", "Temperature", "Motion", "Door", "Window"},                                                      // type list read-only
		3: []int{75, 22, 0, 1, 0},                                                                                            // status list read-only
		4: []int{0, -10, 0, 0, 0},                                                                                            // minValue list read-only
		5: []int{100, 40, 1, 1, 1},                                                                                           // maxValue list read-only
		6: []string{ // lastSamplingTime list read-only
			time.Now().Add(-30 * time.Second).Format("2006-01-02 15:04:05"),
			time.Now().Add(-15 * time.Second).Format("2006-01-02 15:04:05"),
			time.Now().Add(-45 * time.Second).Format("2006-01-02 15:04:05"),
			time.Now().Add(-30 * time.Second).Format("2006-01-02 15:04:05"),
			time.Now().Add(-15 * time.Second).Format("2006-01-02 15:04:05"),
		},
	},
	3: map[int]interface{}{ // actuators Table
		1: []string{"00:1B:44:11:3A:C0", "00:1B:44:11:3A:C1", "00:1B:44:11:3A:E0", "00:1B:44:11:3A:E1"}, // id list read-only
		2: []string{"Light Switch", "Temperature Control", "Alarm", "Door Lock"},                        // type list read-only
		3: []int{1, 22, 0, 1},                                                                           // status list read-write
		4: []int{0, 16, 0, 0},                                                                           // minValue list read-only
		5: []int{1, 30, 1, 1},                                                                           // maxValue list read-only
		6: []string{ // lastControlTime list read-only
			time.Now().Add(-2 * time.Minute).Format("2006-01-02 15:04:05"),
			time.Now().Add(-2*time.Minute + 30*time.Second).Format("2006-01-02 15:04:05"),
			time.Now().Add(-3 * time.Minute).Format("2006-01-02 15:04:05"),
			time.Now().Add(-3*time.Minute + 30*time.Second).Format("2006-01-02 15:04:05"),
		},
	},
}

func handleSetRequest(pdu messages.PDU, conn *net.UDPConn, addr *net.UDPAddr) {
	var responseValues []types.Tipo
	var responseErrors []types.Tipo

	// Map of read-write objects
	readWriteObjects := map[string]bool{
		"1.3":  true, // beaconRate
		"1.6":  true, // dateAndTime
		"1.10": true, // reset
		"3.3":  true, // actuator status
	}

	// Process each IID and its corresponding value
	for i, iid := range pdu.Iid_list.Elements {
		structure := iid.Value.Structure
		object := iid.Value.Objecto
		firstIndex := iid.Value.First_index

		// Check if it's a read-write object
		objectKey := fmt.Sprintf("%d.%d", structure, object)
		if !readWriteObjects[objectKey] {
			errorMsg := fmt.Sprintf("Object %s is read-only", objectKey)
			responseErrors = append(responseErrors, types.Tipo{
				Data_Type: 'S',
				Length:    1,
				Value:     errorMsg,
			})
			continue
		}

		newValue := pdu.Value_list.Elements[i]
		var updateError error

		// Handle each writable object type
		switch {
		case structure == 1 && object == 3: // beaconRate
			if val, err := strconv.Atoi(newValue.Value); err == nil {
				if val > 0 {
					if deviceGroup, ok := mockLMIB[1].(map[int]interface{}); ok {
						deviceGroup[3] = []int{val}
						updateError = nil
					} else {
						updateError = fmt.Errorf("internal error accessing device group")
					}
				} else {
					updateError = fmt.Errorf("beacon rate must be positive")
				}
			} else {
				updateError = fmt.Errorf("invalid beacon rate value")
			}

		case structure == 1 && object == 6: // dateAndTime
			if _, err := time.Parse("2006-01-02 15:04:05", newValue.Value); err == nil {
				if deviceGroup, ok := mockLMIB[1].(map[int]interface{}); ok {
					deviceGroup[6] = []string{newValue.Value}
					updateError = nil
				} else {
					updateError = fmt.Errorf("internal error accessing device group")
				}
			} else {
				updateError = fmt.Errorf("invalid date time format, use YYYY-MM-DD HH:MM:SS")
			}

		case structure == 1 && object == 10: // reset
			if val, err := strconv.Atoi(newValue.Value); err == nil {
				if val == 0 || val == 1 {
					if deviceGroup, ok := mockLMIB[1].(map[int]interface{}); ok {
						deviceGroup[10] = []int{val}
						updateError = nil
					} else {
						updateError = fmt.Errorf("internal error accessing device group")
					}
				} else {
					updateError = fmt.Errorf("reset value must be 0 or 1")
				}
			} else {
				updateError = fmt.Errorf("invalid reset value")
			}

		case structure == 3 && object == 3: // actuator status
			if val, err := strconv.Atoi(newValue.Value); err == nil {
				if actuatorGroup, ok := mockLMIB[3].(map[int]interface{}); ok {
					if firstIndex < 1 || firstIndex > 4 {
						updateError = fmt.Errorf("invalid actuator index (must be 1-4)")
						break
					}

					// Get min and max values for this actuator
					minValues := actuatorGroup[4].([]int)
					maxValues := actuatorGroup[5].([]int)
					actuatorIdx := firstIndex - 1

					if val >= minValues[actuatorIdx] && val <= maxValues[actuatorIdx] {
						statusList := actuatorGroup[3].([]int)
						statusList[actuatorIdx] = val
						actuatorGroup[3] = statusList

						// Update last control time
						timeList := actuatorGroup[6].([]string)
						timeList[actuatorIdx] = time.Now().Format("2006-01-02 15:04:05")
						actuatorGroup[6] = timeList

						updateError = nil
					} else {
						updateError = fmt.Errorf("value out of range (min: %d, max: %d)",
							minValues[actuatorIdx], maxValues[actuatorIdx])
					}
				} else {
					updateError = fmt.Errorf("internal error accessing actuator group")
				}
			} else {
				updateError = fmt.Errorf("invalid actuator status value")
			}
		}

		if updateError != nil {
			responseErrors = append(responseErrors, types.Tipo{
				Data_Type: 'S',
				Length:    1,
				Value:     updateError.Error(),
			})
		} else {
			responseValues = append(responseValues, types.Tipo{
				Data_Type: 'S',
				Length:    1,
				Value:     "Value set successfully",
			})
		}
	}

	// Create response PDU
	responsePDU := messages.NewPDU(
		'R',
		types.NewRequestTimestamp(),
		"agente",
		pdu.Iid_list,
		types.Lists{N_Elements: len(responseValues), Elements: responseValues},
		types.Lists{N_Elements: len(responseErrors), Elements: responseErrors},
	)

	serializedPDU := responsePDU.SerializePDU()
	fmt.Printf("Serialized PDU: %s\n", serializedPDU)

	if _, err := conn.WriteToUDP([]byte(serializedPDU), addr); err != nil {
		fmt.Printf("Error sending response: %v\n", err)
	} else {
		fmt.Println("Response PDU sent successfully")
	}
}

func sendResponse(receivedPDU messages.PDU, conn *net.UDPConn, addr *net.UDPAddr) {
	responsePDU := messages.PDU{
		Tag:               receivedPDU.Tag,
		Tipo:              'R',
		Timestamp:         types.NewRequestTimestamp(),
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

func StartNotificationSender(conn *net.UDPConn, addr *net.UDPAddr) {
	// Get the device group for beacon rate
	deviceGroup, ok := mockLMIB[1].(map[int]interface{})
	if !ok {
		fmt.Println("Error: Could not access device group")
		return
	}

	// Get the beacon rate
	beaconRates, ok := deviceGroup[3].([]int)
	if !ok || len(beaconRates) == 0 {
		fmt.Println("Error: Could not get beacon rate")
		return
	}
	rate := beaconRates[0]

	// Create ticker for the beacon rate
	ticker := time.NewTicker(time.Duration(rate) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Get temperature sensor value
		sensorsTable, ok := mockLMIB[2].(map[int]interface{})
		if !ok {
			fmt.Println("Error: Could not access sensors table")
			continue
		}

		statusList, ok := sensorsTable[3].([]int)
		if !ok || len(statusList) < 2 {
			fmt.Println("Error: Could not access temperature sensor status")
			continue
		}

		temperatureValue := statusList[1] // Index 1 is the temperature sensor

		// Create IID list for temperature sensor status
		iidList := types.IID_List{
			N_Elements: 1,
			Elements: []types.IID_Tipo{
				{
					Data_Type: 'D',
					Length:    4,
					Value: types.IID{
						Structure:    2, // Sensors table
						Objecto:      3, // Status list
						First_index:  2, // Temperature sensor (index 2)
						Second_index: 0,
					},
				},
			},
		}

		// Create value list with temperature
		valueList := types.Lists{
			N_Elements: 1,
			Elements: []types.Tipo{
				{
					Data_Type: 'I',
					Length:    1,
					Value:     strconv.Itoa(temperatureValue),
				},
			},
		}

		// Create notification PDU
		pdu := messages.NewPDU(
			'N', // Notification type
			types.NewInfoTimestamp(),
			"agente",
			iidList,
			valueList,
			types.Lists{}, // Empty error list
		)

		// Serialize and send PDU
		serializedPDU := pdu.SerializePDU()
		if _, err := conn.WriteToUDP([]byte(serializedPDU), addr); err != nil {
			fmt.Printf("Error sending notification: %v\n", err)
		} else {
			fmt.Printf("Notification sent - Temperature: %d°C\n", temperatureValue)
		}
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
	pdu.Print()

	if pdu.Tipo == 'G' {
		sendResponse(pdu, ser, remoteaddr)
	}
	if pdu.Tipo == 'S' {
		handleSetRequest(pdu, ser, remoteaddr)
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

	go func() {
		// Start the notification sender in a goroutine
		go StartNotificationSender(ser, &gestorAddr)
	}()

	wg.Add(1)
	go readPDU(ser)
	wg.Wait()
}
