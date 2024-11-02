package types

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Lists struct {
	N_Elements int // ao passar na funcao que transforma em string por /0 no final
	Elements   []Tipo
}

type Tipo struct {
	Data_Type byte // pode ser integer(I), timestamp(T), string(S)
	Length    int
	Value     string
}

type IID_List struct {
	N_Elements int
	Elements   []IID_Tipo
}

type IID_Tipo struct {
	Data_Type byte // vai ser sempre IID(D)
	Length    int
	Value     IID
}

type IID struct {
	Structure    int
	Objecto      int
	First_index  int // corresponde ao primeiro elemento que queremos
	Second_index int // corresponde até onde queremos as informacoes dos elementos
}

// Tipo printer
func (t Tipo) Print() {
	fmt.Printf("{Type: %c, Length: %d, Value: %s}", t.Data_Type, t.Length, t.Value)
}

// IID_List printer
func (il IID_List) Print() {
	if il.N_Elements == 0 {
		fmt.Printf("[]")
		return
	}
	fmt.Printf("[\n")
	for i, elem := range il.Elements {
		fmt.Printf("    ") // Indent
		elem.Print()
		if i < len(il.Elements)-1 {
			fmt.Printf(",\n")
		}
	}
	fmt.Printf("\n  ]")
}

// IID_Tipo printer
func (it IID_Tipo) Print() {
	fmt.Printf("{Type: %c, Length: %d, Value: ", it.Data_Type, it.Length)
	it.Value.Print()
	fmt.Printf("}")
}

// IID printer
func (iid IID) Print() {
	if iid.Second_index != 0 {
		fmt.Printf("{%d.%d.%d.%d}", iid.Structure, iid.Objecto, iid.First_index, iid.Second_index)
	} else {
		fmt.Printf("{%d.%d.%d}", iid.Structure, iid.Objecto, iid.First_index)
	}
}

// Lists printer
func (l Lists) Print() {
	if l.N_Elements == 0 {
		fmt.Printf("[]")
		return
	}
	fmt.Printf("[\n")
	for i, elem := range l.Elements {
		fmt.Printf("    ") // Indent
		elem.Print()
		if i < len(l.Elements)-1 {
			fmt.Printf(",\n")
		}
	}
	fmt.Printf("\n  ]")
}

func (l IID_List) IIDListSerialize() string {
	line := fmt.Sprintf(`%d\0`, l.N_Elements)

	// Iterate over the Tipo slice and serialize each Tipo
	for _, t := range l.Elements {
		line += t.TipoIIDSerialize()
	}

	return line
}

func DeserializeIID_List(serialized string) IID_List {
	elements := strings.Split(serialized, `\0`)
	nElements, _ := strconv.Atoi(elements[0])

	if nElements == 0 {
		return IID_List{N_Elements: 0, Elements: []IID_Tipo{}}
	}

	var iidTypes []IID_Tipo

	// For each IID in the list
	currentIndex := 1
	for i := 0; i < nElements && currentIndex < len(elements)-2; i++ {
		// Get Data_Type, Length, and Value parts
		dataType := elements[currentIndex]
		length := elements[currentIndex+1]
		value := elements[currentIndex+2]

		// Create the serialized string for one IID_Tipo
		tipoStr := dataType + `\0` + length + `\0` + value

		iidType := DeserializeIID_Tipo(tipoStr)
		iidTypes = append(iidTypes, iidType)

		currentIndex += 3
	}

	return IID_List{
		N_Elements: nElements,
		Elements:   iidTypes,
	}
}

func (t IID_Tipo) TipoIIDSerialize() string {
	return fmt.Sprintf(`%c\0%d\0%s\0`, t.Data_Type, t.Length, t.Value.IIDSerialize())
}

func DeserializeIID_Tipo(serialized string) IID_Tipo {
	elements := strings.Split(serialized, `\0`)
	length, _ := strconv.Atoi(elements[1])
	return IID_Tipo{
		Data_Type: elements[0][0],
		Length:    length,
		Value:     DeserializeIID(elements[2]),
	}
}

func (i IID) IIDSerialize() string {
	return fmt.Sprintf(`%d%d%d%d`, i.Structure, i.Objecto, i.First_index, i.Second_index)
}

func DeserializeIID(serialized string) IID {
	// For an input like "1112"
	var iid IID
	if len(serialized) >= 1 {
		iid.Structure, _ = strconv.Atoi(serialized[0:1])
	}
	if len(serialized) >= 2 {
		iid.Objecto, _ = strconv.Atoi(serialized[1:2])
	}
	if len(serialized) >= 3 {
		iid.First_index, _ = strconv.Atoi(serialized[2:3])
	}
	if len(serialized) >= 4 {
		iid.Second_index, _ = strconv.Atoi(serialized[3:4])
	}
	return iid
}

func (t Tipo) TipoSerialize() string {
	return fmt.Sprintf(`%c\0%d\0%s\0`, t.Data_Type, t.Length, t.Value)
}

func DeserializeTipo(serialized string) Tipo {
	elements := strings.Split(serialized, `\0`)
	data, _ := strconv.Atoi(elements[0])
	byte := byte(data)
	length, _ := strconv.Atoi(elements[1])
	return Tipo{
		Data_Type: byte,
		Length:    length,
		Value:     elements[2],
	}
}

func (l Lists) ListsSerialize() string {
	line := fmt.Sprintf(`%d\0`, l.N_Elements)

	// Iterate over the Tipo slice and serialize each Tipo
	for _, t := range l.Elements {
		line += t.TipoSerialize()
	}

	return line
}

func DeserializeLists(serialized string) Lists {
	elements := strings.Split(serialized, `\0`)
	nElements, _ := strconv.Atoi(elements[0])

	if nElements == 0 {
		return Lists{N_Elements: 0, Elements: []Tipo{}}
	}

	var tipos []Tipo
	currentIndex := 1

	for i := 0; i < nElements; i++ {
		if currentIndex+2 >= len(elements) {
			break
		}

		// Reconstruct the tipo string with \0 separators
		tipoStr := elements[currentIndex] + `\0` + elements[currentIndex+1] + `\0` + elements[currentIndex+2]
		tipo := DeserializeTipo(tipoStr)
		tipos = append(tipos, tipo)
		currentIndex += 3
	}

	return Lists{
		N_Elements: nElements,
		Elements:   tipos,
	}
}

func NewLists(n_elements int, elements []Tipo) Lists {
	return Lists{
		N_Elements: n_elements,
		Elements:   elements,
	}
}

func NewTipo(data byte, length int, value string) Tipo {
	return Tipo{
		Data_Type: data,
		Length:    length,
		Value:     value,
	}
}

func NewIID_List(n_elements int, elements []IID_Tipo) IID_List {
	return IID_List{
		N_Elements: n_elements,
		Elements:   elements,
	}
}

func NewIID_Tipo(length int, value IID) IID_Tipo {
	return IID_Tipo{
		Data_Type: 'D',
		Length:    length, // pode ser 2,3 ou 4
		Value:     value,
	}
}

func NewIID(structure int, objeto int, first int, second int) IID {
	return IID{
		Structure:    structure,
		Objecto:      objeto,
		First_index:  first,
		Second_index: second,
	}
}

// Timestamps
func NewRequestTimestamp() Tipo {
	t := time.Now()
	timestamp := t.Format("02:01:2006:15:04:05.000")
	return Tipo{
		Data_Type: 'T',
		Length:    7,
		Value:     timestamp,
	}
}

func NewInfoTimestamp(value string) Tipo {
	return Tipo{
		Data_Type: 'T',
		Length:    5,
		Value:     value,
	}
}
