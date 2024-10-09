package messages

import (
	"fmt"
	"gsr/messages/types"
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
	elements := strings.SplitN(serialized, `\0`, 7)

	pdu := PDU{}
	pdu.Tag = elements[0]
	pdu.Tipo = elements[1][0]
	pdu.Timestamp = types.DeserializeTipo(elements[2] + `\0` + elements[3] + `\0` + elements[4])
	pdu.MessageIdentifier = elements[5]

	remaining := elements[6]

	// Deserialize IID_List
	iidListParts := strings.SplitN(remaining, `\0`, 2)
	pdu.Iid_list = types.DeserializeIID_List(iidListParts[0] + `\0` + iidListParts[1])

	// Move to Value_list
	valueListStart := strings.Index(iidListParts[1], `\00\0`) + 3
	remaining = iidListParts[1][valueListStart:]

	// Deserialize Value_list
	valueListParts := strings.SplitN(remaining, `\00\0`, 2)
	pdu.Value_list = types.DeserializeLists(valueListParts[0])

	// Deserialize Error_list
	if len(valueListParts) > 1 {
		pdu.Error_list = types.DeserializeLists(valueListParts[1])
	} else {
		pdu.Error_list = types.Lists{N_Elements: 0, Elements: []types.Tipo{}}
	}

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

func (p PDU) Print() {
	fmt.Println("Tag:", p.Tag)
	fmt.Println("Tipo:", p.Tipo)
	fmt.Println("Timestamp:", p.Timestamp)
	fmt.Println("Message Identifier:", p.MessageIdentifier)
	fmt.Println("Iid List:", p.Iid_list)
	fmt.Println("Value List:", p.Value_list)
	fmt.Println("Error List:", p.Error_list)
}
