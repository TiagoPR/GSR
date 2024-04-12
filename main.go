package main

import (
	"fmt"
	"gsr/snmp"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go snmp.SetupAgente()
	go snmp.SetupGestor()

	wg.Wait()
	fmt.Println("Hello, World")
}
