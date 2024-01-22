package main

import (
	"fmt"

	"golang.org/x/text/encoding/japanese"
)

func main() {
	buf := []byte{0xe4, 0xaa}
	dec := japanese.ShiftJIS.NewDecoder()
	decBuf, err := dec.Bytes(buf)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(decBuf))
}
