package messages

import (
	"fmt"
	"gsr/messages/types"
)

type PDU struct {
	Tag               string
	Tipo              string
	Timestamp         types.Tipo // gerada no envio da mensagem
	MessageIdentifier string
	Iid_list          types.IID_List
	Value_list        types.Lists
	Error_list        types.Lists
}

// constructor of pdu with tag default value
func NewPDU(tipo string, timestamp types.Tipo, messageIdentifier string, iid_list types.IID_List, value_list types.Lists, error_list types.Lists) PDU {
	pdu := PDU{}
	pdu.Tag = "kdk847ufh84jg87g\\0"
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
