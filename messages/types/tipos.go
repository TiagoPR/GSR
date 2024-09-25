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

func (t Tipo) Print() {
	fmt.Println("Data Type", t.Data_Type)
	fmt.Println("Length", t.Length)
	fmt.Println("Value", t.Value)
}

func (l IID_List) IIDListSerialize() string {
	line := fmt.Sprintf(`%d\0`, l.N_Elements)

	// Iterate over the Tipo slice and serialize each Tipo
	for _, t := range l.Elements {
		line += t.TipoIIDSerialize()
	}

	return line
}

func (t IID_Tipo) TipoIIDSerialize() string {
	return fmt.Sprintf(`%c\0%d\0%s\0`, t.Data_Type, t.Length, t.Value.IIDSerialize())
}

func (i IID) IIDSerialize() string {
	return fmt.Sprintf(`%d%d%d%d`, i.Structure, i.Objecto, i.First_index, i.Second_index)
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

	var tipos []Tipo
	for i := 1; i < len(elements); i += 3 {
		data_type, _ := strconv.Atoi(elements[i])
		byte := byte(data_type)
		length, _ := strconv.Atoi(elements[i+1])
		value := elements[i+2]
		tipos = append(tipos, Tipo{Data_Type: byte, Length: length, Value: value})
	}

	return Lists{N_Elements: nElements, Elements: tipos}
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
