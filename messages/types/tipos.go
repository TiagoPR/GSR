package types

import "fmt"

type Lists struct {
	N_Elements int // ao passar na funcao que transforma em string por /0 no final
	Elements   []Tipo
}

type Tipo struct {
	Data_Type string // pode ser integer(I), timestamp(T), string(S)
	Length    int
	Value     string
}

type IID_List struct {
	N_Elements int
	Elements   []IID_Tipo
}

type IID_Tipo struct {
	Length int
	Value  IID
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

func (i IID) IIDSerialize() string {
	return fmt.Sprintf(`%d%d%d%d`, i.Structure, i.Objecto, i.First_index, i.Second_index)
}

func (t Tipo) TipoSerialize() string {
	return fmt.Sprintf(`%s\0%d\0%s\0`, t.Data_Type, t.Length, t.Value)
}

func (l Lists) ListsSerialize() string {
	line := fmt.Sprintf(`%d\0`, l.N_Elements)

	// Iterate over the Tipo slice and serialize each Tipo
	for _, t := range l.Elements {
		line += t.TipoSerialize()
	}

	return line
}

func newLists(n_elements int, elements []Tipo) Lists {
	return Lists{
		N_Elements: n_elements,
		Elements:   elements,
	}
}

func newTipo(data string, length int, value string) Tipo {
	return Tipo{
		Data_Type: data,
		Length:    length,
		Value:     value,
	}
}

func newIID_List(n_elements int, elements []IID) IID_List {
	return IID_List{
		N_Elements: n_elements,
		Elements:   elements,
	}
}

func newIID_Tipo(length int, value IID) IID_Tipo {
	return IID_Tipo{
		Length: length, // pode ser 2,3 ou 4
		Value:  value,
	}
}

func newIID(structure int, objeto int, first int, second int) IID {
	return IID{
		Structure:    structure,
		Objecto:      objeto,
		First_index:  first,
		Second_index: second,
	}
}

// Timestamps
func NewRequestTimestamp(value string) Tipo {
	return Tipo{
		Data_Type: "T",
		Length:    7,
		Value:     value,
	}
}

func NewInfoTimestamp(value string) Tipo {
	return Tipo{
		Data_Type: "T",
		Length:    5,
		Value:     value,
	}
}
