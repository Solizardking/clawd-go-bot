package base58

import (
	"encoding/binary"
	"errors"
)

var (
	ErrInvalidChar   = errors.New("base58: invalid base58 character")
	ErrInvalidLength = errors.New("base58: invalid encoded length")
	ErrValueTooLarge = errors.New("base58: decoded value too large for output size")
	ErrLeadingZeros  = errors.New("base58: leading '1' count does not match leading zero bytes")
)

// Decode32 decodes a base58 string into a 32-byte array.
func Decode32(encoded string, dst *[32]byte) error {
	encLen := len(encoded)
	if encLen == 0 || encLen > raw58Sz32 {
		return ErrInvalidLength
	}

	var raw [raw58Sz32]byte
	offset := raw58Sz32 - encLen
	for i := range encLen {
		c := encoded[i]
		if c < '1' || c > 'z' {
			return ErrInvalidChar
		}
		digit := base58Inverse[c-'1']
		if digit == base58InvalidDigit {
			return ErrInvalidChar
		}
		raw[offset+i] = digit
	}

	var intermediate [intermediateSz32]uint64
	for i := range intermediateSz32 {
		intermediate[i] = uint64(raw[5*i+0])*11316496 +
			uint64(raw[5*i+1])*195112 +
			uint64(raw[5*i+2])*3364 +
			uint64(raw[5*i+3])*58 +
			uint64(raw[5*i+4])
	}

	// Matrix-vector multiply (assembly on arm64, Go on other archs).
	var bin [binarySz32]uint64
	decodeMatMul32(&intermediate, &bin)

	for i := binarySz32 - 1; i >= 1; i-- {
		bin[i-1] += bin[i] >> 32
		bin[i] &= 0xFFFFFFFF
	}

	if bin[0] > 0xFFFFFFFF {
		return ErrValueTooLarge
	}

	for i := range binarySz32 {
		binary.BigEndian.PutUint32(dst[i*4:i*4+4], uint32(bin[i]))
	}

	return validateLeadingZeros(encoded, dst[:])
}

// Decode64 decodes a base58 string into a 64-byte array.
func Decode64(encoded string, dst *[64]byte) error {
	encLen := len(encoded)
	if encLen == 0 || encLen > raw58Sz64 {
		return ErrInvalidLength
	}

	var raw [raw58Sz64]byte
	offset := raw58Sz64 - encLen
	for i := range encLen {
		c := encoded[i]
		if c < '1' || c > 'z' {
			return ErrInvalidChar
		}
		digit := base58Inverse[c-'1']
		if digit == base58InvalidDigit {
			return ErrInvalidChar
		}
		raw[offset+i] = digit
	}

	var intermediate [intermediateSz64]uint64
	for i := range intermediateSz64 {
		intermediate[i] = uint64(raw[5*i+0])*11316496 +
			uint64(raw[5*i+1])*195112 +
			uint64(raw[5*i+2])*3364 +
			uint64(raw[5*i+3])*58 +
			uint64(raw[5*i+4])
	}

	// Plain uint64 accumulation — each product is ≤ 2^62 and the sum
	// of 18 terms stays under 2^64 (verified by Firedancer analysis).
	var bin [binarySz64]uint64
	for k := range binarySz64 {
		var acc uint64
		for i := range intermediateSz64 {
			acc += intermediate[i] * uint64(decTable64[i][k])
		}
		bin[k] = acc
	}

	for i := binarySz64 - 1; i >= 1; i-- {
		bin[i-1] += bin[i] >> 32
		bin[i] &= 0xFFFFFFFF
	}

	if bin[0] > 0xFFFFFFFF {
		return ErrValueTooLarge
	}

	for i := range binarySz64 {
		binary.BigEndian.PutUint32(dst[i*4:i*4+4], uint32(bin[i]))
	}

	return validateLeadingZeros(encoded, dst[:])
}

// validateLeadingZeros verifies that the number of leading '1' characters in
// the encoded input equals the number of leading zero bytes in the decoded
// output. This is a required invariant of base58: each leading zero byte in
// the raw value is represented by exactly one '1' in the encoding.
func validateLeadingZeros(encoded string, dst []byte) error {
	inLeading1s := 0
	for i := 0; i < len(encoded) && encoded[i] == '1'; i++ {
		inLeading1s++
	}

	outLeading0s := 0
	for _, b := range dst {
		if b != 0 {
			break
		}
		outLeading0s++
	}

	if inLeading1s != outLeading0s {
		return ErrLeadingZeros
	}
	return nil
}
