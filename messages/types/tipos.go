package types

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
	Elements   []IID
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
