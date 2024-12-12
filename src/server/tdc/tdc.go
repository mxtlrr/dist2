/* TDC library. Small digit compression.
 * I've spent so much time thinking about how to do this
 * that it's driving me crazy. */
package tdc

// The tried and true method for digit compression is storing each digit
// in a nybble of a byte. This stores 2 digits/byte. Also, this will provide
// 50% compression for n bytes (C(n) = n/2)
func TDCEncodeString(data string) []byte {
	var res []byte
	// Ideally, data's length should be
	// divisible by 2
	for i := 0; i < len(data); i += 2 {
		var (
			o  byte = data[i]
			o2 byte = data[i+1]
		)

		// Convert.
		newByte := (o << 4) | o2
		res = append(res, newByte)
	}
	return res
}

func TDCDecodeString(encode []byte) string {
	// TODO...
}
