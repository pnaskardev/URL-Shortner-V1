package utils

import "math/big"

var base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func EncodeToBase62(data []byte) string {
	if len(data) == 0 {
		return ""
	}

	// Treat data as a big-endian number
	num := new(big.Int).SetBytes(data)
	var result []byte
	zero := big.NewInt(0)
	base := big.NewInt(62)

	for num.Cmp(zero) > 0 {
		var remainder big.Int
		num.DivMod(num, base, &remainder)
		result = append(result, base62Chars[remainder.Int64()])
	}

	// Reverse the result as we built it backwards
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result)
}
