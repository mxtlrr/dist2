/* TDC library. Small digit compression.
 * I've spent so much time thinking about how to do this
 * that it's driving me crazy. */
package tdc

import (
	"fmt"
	"strings"
)

// The tried and true method for digit compression is storing each digit
// in a nybble of a byte. This stores 2 digits/byte. Also, this will provide
// 50% compression for n bytes (C(n) = n/2)
func TDCEncodeString(data string) []byte {
	var res []byte
	// Ideally, data's length should be
	// divisible by 2
	fmt.Println(len(data))
	for i := 0; i < len(data); i += 2 {
		var (
			o  byte = data[i] - 0x30
			o2 byte = data[i+1] - 0x30
		)

		// Convert.
		newByte := (o << 4) | o2
		res = append(res, newByte)
	}
	return res
}

func TDCDecodeString(encode []byte) string {
	var ret strings.Builder
	for bt := range encode {
		var (
			top byte = encode[bt] >> 4  // High nibble
			low byte = encode[bt] & 0xF // Low nibble
		)

		ret.WriteString(fmt.Sprintf("%d%d", top, low))
	}
	return ret.String()
}
