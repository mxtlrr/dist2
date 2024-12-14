package main

import (
	"fmt"
	"bufio"
	"os"
	"strconv"

	"github.com/mxtlrr/dist2/src/server/tdc"
)
func main() {
	var (
		digPath  string
		digCount string
	)
	fmt.Print("Enter full path of digit file: ")
	fmt.Scanln(&digPath)

	fmt.Print("Enter digit count: ")
	fmt.Scanln(&digCount)

	realDigCount, _ := strconv.Atoi(digCount)

	fileP, err := os.Open(digPath)
	if err != nil {
		panic(err)
	}
  r4 := bufio.NewReader(fileP)
  b4, err := r4.Peek(realDigCount/2)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Decompressed data:\n\n%s\n", tdc.TDCDecodeString(b4))
}
