package messages

import (
	"fmt"
	"gsr/messages/types"
	"strconv"
	"strings"
)

type PDU struct {
	Tag               string
	Tipo              byte
	Timestamp         types.Tipo // gerada no envio da mensagem
	MessageIdentifier string
	Iid_list          types.IID_List
	Value_list        types.Lists
	Error_list        types.Lists
}

func (p PDU) SerializePDU() string {
	return fmt.Sprintf(`%s\0%c\0%s%s\0%s%s%s`, p.Tag, p.Tipo, p.Timestamp.TipoSerialize(), p.MessageIdentifier, p.Iid_list.IIDListSerialize(), p.Value_list.ListsSerialize(), p.Error_list.ListsSerialize())
}

func DeserializePDU(serialized string) PDU {
	fmt.Println("\nDeserializing PDU...")
	elements := strings.SplitN(serialized, `\0`, 7)
	pdu := PDU{}
	pdu.Tag = elements[0]
	pdu.Tipo = elements[1][0]
	pdu.Timestamp = types.DeserializeTipo(elements[2] + `\0` + elements[3] + `\0` + elements[4])
	pdu.MessageIdentifier = elements[5]
	remaining := elements[6]

	// First get number of IID elements
	parts := strings.SplitN(remaining, `\0`, 2)
	nIIDElements, _ := strconv.Atoi(parts[0])
	if len(parts) < 2 {
		return pdu
	}
	remaining = parts[1]

	// Calculate where IID_List ends based on number of elements
	iidParts := make([]string, 0)
	iidParts = append(iidParts, strconv.Itoa(nIIDElements))

	// Process each IID
	currentRemaining := remaining
	for i := 0; i < nIIDElements; i++ {
		currentParts := strings.SplitN(currentRemaining, `\0`, 4)
		if len(currentParts) < 3 {
			break
		}
		iidParts = append(iidParts, currentParts[0], currentParts[1], currentParts[2])
		if len(currentParts) > 3 {
			currentRemaining = currentParts[3]
		} else {
			currentRemaining = ""
			break
		}
	}

	// Deserialize IID_List
	pdu.Iid_list = types.DeserializeIID_List(strings.Join(iidParts, `\0`))

	fmt.Println("Processing Value/Error lists...")
	fmt.Printf("Current remaining: %s\n", currentRemaining)

	// Initialize empty lists
	pdu.Value_list = types.Lists{N_Elements: 0, Elements: []types.Tipo{}}
	pdu.Error_list = types.Lists{N_Elements: 0, Elements: []types.Tipo{}}

	// First check if we have a value list by looking at the first character
	firstChar := currentRemaining[0]
	if firstChar == '0' {
		// No value list, everything after first \0 is error list
		parts := strings.SplitN(currentRemaining, `\0`, 2)
		if len(parts) > 1 {
			fmt.Println("Found error list only")
			fmt.Printf("Error list string: %s\n", parts[1])
			pdu.Error_list = types.DeserializeLists(parts[1])
		}
	} else {
		// We have a value list, split by \00\0
		listParts := strings.Split(currentRemaining, `\00\0`)
		fmt.Printf("Parts after \\00\\0 split: %v\n", listParts)

		if len(listParts) > 0 {
			fmt.Printf("Value list string: %s\n", listParts[0])
			pdu.Value_list = types.DeserializeLists(listParts[0])

			if len(listParts) > 1 {
				fmt.Printf("Error list string: %s\n", listParts[1])
				pdu.Error_list = types.DeserializeLists(listParts[1])
			}
		}
	}

	fmt.Printf("Final Value list: %+v\n", pdu.Value_list)
	fmt.Printf("Final Error list: %+v\n", pdu.Error_list)
	return pdu
}

// constructor of pdu with tag default value
func NewPDU(tipo byte, timestamp types.Tipo, messageIdentifier string, iid_list types.IID_List, value_list types.Lists, error_list types.Lists) PDU {
	pdu := PDU{}
	pdu.Tag = "kdk847ufh84jg87g"
	pdu.Tipo = tipo
	pdu.Timestamp = timestamp
	pdu.MessageIdentifier = messageIdentifier
	pdu.Iid_list = iid_list
	pdu.Value_list = value_list
	pdu.Error_list = error_list

	return pdu
}

// PDU printer
func (p PDU) Print() {
	fmt.Println("PDU {")
	fmt.Printf("  Tag: %s\n", p.Tag)
	fmt.Printf("  Type: %c\n", p.Tipo)
	fmt.Printf("  Timestamp: ")
	p.Timestamp.Print()
	fmt.Printf("\n")
	fmt.Printf("  Message Identifier: %s\n", p.MessageIdentifier)
	fmt.Printf("  IID List: ")
	p.Iid_list.Print()
	fmt.Printf("\n")
	fmt.Printf("  Value List: ")
	p.Value_list.Print()
	fmt.Printf("\n")
	fmt.Printf("  Error List: ")
	p.Error_list.Print()
	fmt.Printf("\n")
	fmt.Println("}")
}
