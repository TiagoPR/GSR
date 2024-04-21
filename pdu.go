package gsr

import (
	"fmt"
)

type PDU struct {
	Tag               string
	Tipo              byte
	Timestamp         string
	MessageIdentifier string
	Iid_list          []string
	Value_list        []string
	Error_list        []string
}

// constructor of pdu with tag default value
func NewPDU(tipo byte, timestamp string, messageIdentifier string, iid_list []string, value_list []string, error_list []string) PDU {
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
