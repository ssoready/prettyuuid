package prettyuuid

import (
	"fmt"
	"math"
	"math/big"
	"strings"
)

type Format struct {
	Prefix   string
	Alphabet string
}

func (f *Format) Format(uuid [16]byte) string {
	var n big.Int
	n.SetBytes(uuid[:])

	b := big.NewInt(int64(len(f.Alphabet)))

	s := make([]byte, f.len())
	copy(s, f.Prefix)
	for i := len(f.Prefix); i < f.len(); i++ {
		s[i] = f.Alphabet[0]
	}

	i := len(s) - 1
	for n.BitLen() > 0 {
		var r big.Int
		n.QuoRem(&n, b, &r)

		s[i] = f.Alphabet[r.Int64()]
		i--
	}

	return string(s)
}

func (f *Format) Parse(s string) ([16]byte, error) {
	if len(s) != f.len() {
		return [16]byte{}, fmt.Errorf("%q does not have expected length %v", s, f.len())
	}

	var n big.Int
	for i := len(f.Prefix); i < len(s); i++ {
		d := strings.IndexByte(f.Alphabet, s[i])
		if d == -1 {
			return [16]byte{}, fmt.Errorf("%q contains illegal char at position %v", s, i)
		}

		n.Mul(&n, big.NewInt(int64(len(f.Alphabet))))
		n.Add(&n, big.NewInt(int64(d)))
	}

	b := n.Bytes()
	if len(b) > 16 {
		panic(fmt.Errorf("prettyuuid invariant failure"))
	}

	var out [16]byte
	copy(out[16-len(b):], b)
	return out, nil
}

func (f *Format) len() int {
	// digits required is ceil(log_{radix}(max_uuid))
	// log_{radix}(max_uuid) = log2(max_uuid) / log2(radix) = 128 / log2(radix)
	return len(f.Prefix) + int(math.Ceil(128.0/math.Log2(float64(len(f.Alphabet)))))
}
