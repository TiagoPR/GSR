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
	parts := strings.Split(serialized, `\0`)

	nElements, _ := strconv.Atoi(parts[0])

	var iidTypes []IID_Tipo

	if nElements == 0 {
		return IID_List{N_Elements: 0, Elements: iidTypes}
	}

	// Skip the count (parts[0])
	remaining := strings.Join(parts[1:], `\0`)

	// Process each IID
	for i := 0; i < nElements; i++ {

		// Split for current IID
		iidParts := strings.SplitN(remaining, `\0`, 4)
		if len(iidParts) < 3 {
			fmt.Println("Not enough parts for IID")
			break
		}

		// Create IID_Tipo from the first three parts
		tipoStr := iidParts[0] + `\0` + iidParts[1] + `\0` + iidParts[2]
		iidType := DeserializeIID_Tipo(tipoStr)
		iidTypes = append(iidTypes, iidType)

		// If there's more to process, update remaining
		if len(iidParts) > 3 {
			remaining = iidParts[3]
		} else {
			break
		}

	}

	result := IID_List{
		N_Elements: nElements,
		Elements:   iidTypes,
	}
	return result
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
		Value:     DeserializeIID(elements[2], length),
	}
}

func (i IID) IIDSerialize() string {
	return fmt.Sprintf(`%d%d%d%d`, i.Structure, i.Objecto, i.First_index, i.Second_index)
}

func DeserializeIID(serialized string, length int) IID {
	fmt.Println("Deserialize IID: ", serialized)
	var structure, object, firstIndex, secondIndex int

	// Always take first digit as structure
	if len(serialized) >= 1 {
		structure, _ = strconv.Atoi(serialized[:1])
	}

	// If length is 2, format is x.yy.0 (like 1.10.0)
	if length == 2 {
		object, _ = strconv.Atoi(serialized[1:3])
		firstIndex, _ = strconv.Atoi(serialized[3:])
	} else { // if length is 3, format is x.y.z (like 1.1.0)
		object, _ = strconv.Atoi(serialized[1:2])
		firstIndex, _ = strconv.Atoi(serialized[2:3])
	}

	return IID{
		Structure:    structure,
		Objecto:      object,
		First_index:  firstIndex,
		Second_index: secondIndex,
	}
}

func (t Tipo) TipoSerialize() string {
	return fmt.Sprintf(`%c\0%d\0%s\0`, t.Data_Type, t.Length, t.Value)
}

func DeserializeTipo(serialized string) Tipo {
	elements := strings.Split(serialized, `\0`)

	data := []byte(elements[0])[0] // Get the actual character byte
	length, _ := strconv.Atoi(elements[1])

	return Tipo{
		Data_Type: data,
		Length:    length,
		Value:     elements[2],
	}
}

func (l Lists) ListsSerialize() string {
	if l.N_Elements == 0 {
		return `\0`
	}

	line := fmt.Sprintf(`%d\0`, l.N_Elements)
	for _, t := range l.Elements {
		line += t.TipoSerialize()
	}
	return line
}

func DeserializeLists(serialized string) Lists {
	if serialized == "" {
		return Lists{N_Elements: 0, Elements: []Tipo{}}
	}

	elements := strings.Split(serialized, `\0`)
	nElements, _ := strconv.Atoi(elements[0])

	var tipos []Tipo
	currentIndex := 1

	// Process exactly nElements
	for i := 0; i < nElements && currentIndex+2 < len(elements); i++ {
		// Make sure we get the Data_Type character
		dataType := []byte(elements[currentIndex])[0]
		length, _ := strconv.Atoi(elements[currentIndex+1])
		value := elements[currentIndex+2]

		tipos = append(tipos, Tipo{
			Data_Type: dataType,
			Length:    length,
			Value:     value,
		})
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
func NewInfoTimestamp() Tipo {
	t := time.Now()
	// Format: days hours mins secs ms
	timestamp := fmt.Sprintf("%02d%02d%02d%02d%03d",
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
		t.Nanosecond()/1000000) // Convert nanoseconds to milliseconds
	return Tipo{
		Data_Type: 'T',
		Length:    5,
		Value:     timestamp,
	}
}

func NewRequestTimestamp() Tipo {
	t := time.Now()
	// Format: day month year hours mins secs ms
	timestamp := fmt.Sprintf("%02d%02d%04d%02d%02d%02d%03d",
		t.Day(),
		t.Month(),
		t.Year(),
		t.Hour(),
		t.Minute(),
		t.Second(),
		t.Nanosecond()/1000000) // Convert nanoseconds to milliseconds
	return Tipo{
		Data_Type: 'T',
		Length:    7,
		Value:     timestamp,
	}
}
