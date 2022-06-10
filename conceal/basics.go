package conceal

import (
	"encoding/hex"
	"math"
)

// StringRange used to define an acceptable concealed byte range. Values are inclusive
type StringRange struct {
	From byte
	To   byte
}

// Bytes conceals an array of bytes
func (c *Conceal) Bytes(source []byte) []byte {
	d := make([]byte, len(source))
	for i := range source {
		s := int(source[i])
		n := s + int(c.key[s])
		if n > 255 {
			n -= 256
		}
		d[i] = byte(n)
	}
	return d
}

// String conceals a printable string, with the resulting string also being printable
func (c *Conceal) String(source string) string {
	return c.StringRange(source, StringRange{From: 32, To: 126})
}

// StringRange conceals a printable string, with the resulting string also being printable and within a given ascii range
func (c *Conceal) StringRange(source string, sr StringRange) string {
	d := make([]byte, len(source))
	for i := range source {
		s := int(source[i])
		n := s + int(c.key[s])
		for n > int(sr.To) {
			n -= int(sr.To - sr.From + 1)
		}
		for n < int(sr.From) {
			n += int(sr.To - sr.From - 1)
		}
		d[i] = byte(n)
	}
	return string(d)
}

// StringHex conceals a string representation of a hexadecimal number (e.g. 0f93d8); odd numbers will be padded with a 0 in the front
func (c *Conceal) StringHex(source string) (string, error) {
	if len(source)%2 != 0 {
		source = "0" + source
	}
	sourceBytes, err := hex.DecodeString(source)
	if err != nil {
		return "", err
	}
	destBytes := c.Bytes(sourceBytes)
	return hex.EncodeToString(destBytes), nil
}

// StringRanges conceals a string, ensuring the concealed version falls within the restricted ranges
func (c *Conceal) StringRanges(source string, ranges []StringRange) string {
	d := make([]byte, len(source))
	for i := range source {
		s := int(source[i])
		n := s + int(c.key[s])
		inRange := false
		// find if n is in range
		for _, r := range ranges {
			if n >= int(r.From) && n <= int(r.To) {
				inRange = true
				break
			}
		}
		// n not in range
		if !inRange {
			rangeDist := math.MaxInt
			rangeNo := -1
			// find closest range to the n
			for ri, r := range ranges {
				if n > int(r.To) {
					dist := n - int(r.To)
					if dist < rangeDist {
						rangeNo = ri
					}
				} else if n < int(r.From) {
					dist := int(r.From) - n
					if dist < rangeDist {
						rangeNo = ri
					}
				}
			}
			// adjust n to be within the closest range
			sr := ranges[rangeNo]
			for n > int(sr.To) {
				n -= int(sr.To - sr.From + 1)
			}
			for n < int(sr.From) {
				n += int(sr.To - sr.From - 1)
			}
		}
		d[i] = byte(n)
	}
	return string(d)
}
