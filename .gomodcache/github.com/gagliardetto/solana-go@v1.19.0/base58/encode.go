package base58

import (
	"encoding/binary"
	"unsafe"
)

// Encode32 encodes a 32-byte array to a base58 string.
//
// Allocates exactly one []byte of the encoded length. For zero-allocation
// hot paths, prefer AppendEncode32 which writes into a caller-owned buffer.
func Encode32(src *[32]byte) string {
	var raw [raw58Sz32]byte
	skip := encodeRaw32(src, &raw)
	outLen := raw58Sz32 - skip
	out := make([]byte, outLen)
	for i := range outLen {
		out[i] = base58Chars[raw[skip+i]]
	}
	return unsafe.String(unsafe.SliceData(out), len(out))
}

// Encode64 encodes a 64-byte array to a base58 string.
//
// Allocates exactly one []byte of the encoded length. For zero-allocation
// hot paths, prefer AppendEncode64.
func Encode64(src *[64]byte) string {
	var raw [raw58Sz64]byte
	skip := encodeRaw64(src, &raw)
	outLen := raw58Sz64 - skip
	out := make([]byte, outLen)
	for i := range outLen {
		out[i] = base58Chars[raw[skip+i]]
	}
	return unsafe.String(unsafe.SliceData(out), len(out))
}

// AppendEncode32 appends the base58 encoding of src to dst and returns the
// extended buffer. It allocates only if dst has insufficient capacity.
func AppendEncode32(dst []byte, src *[32]byte) []byte {
	var raw [raw58Sz32]byte
	skip := encodeRaw32(src, &raw)
	outLen := raw58Sz32 - skip
	// Grow dst in place if possible; otherwise allocate.
	total := len(dst) + outLen
	if cap(dst) < total {
		grown := make([]byte, total)
		copy(grown, dst)
		dst = grown
	} else {
		dst = dst[:total]
	}
	out := dst[total-outLen:]
	for i := range outLen {
		out[i] = base58Chars[raw[skip+i]]
	}
	return dst
}

// AppendEncode64 appends the base58 encoding of src to dst and returns the
// extended buffer. It allocates only if dst has insufficient capacity.
func AppendEncode64(dst []byte, src *[64]byte) []byte {
	var raw [raw58Sz64]byte
	skip := encodeRaw64(src, &raw)
	outLen := raw58Sz64 - skip
	total := len(dst) + outLen
	if cap(dst) < total {
		grown := make([]byte, total)
		copy(grown, dst)
		dst = grown
	} else {
		dst = dst[:total]
	}
	out := dst[total-outLen:]
	for i := range outLen {
		out[i] = base58Chars[raw[skip+i]]
	}
	return dst
}

// encodeRaw32 fills raw with the raw base-58 digits for a 32-byte input and
// returns the number of leading digits to skip when producing the final output.
func encodeRaw32(src *[32]byte, raw *[raw58Sz32]byte) int {
	var intermediate [intermediateSz32]uint64
	encodeMatMul32(src, &intermediate)

	for i := intermediateSz32 - 1; i >= 1; i-- {
		intermediate[i-1] += intermediate[i] / r1div
		intermediate[i] %= r1div
	}

	for i := range intermediateSz32 {
		v := uint32(intermediate[i])
		raw[5*i+4] = byte(v % 58)
		v /= 58
		raw[5*i+3] = byte(v % 58)
		v /= 58
		raw[5*i+2] = byte(v % 58)
		v /= 58
		raw[5*i+1] = byte(v % 58)
		v /= 58
		raw[5*i+0] = byte(v)
	}

	inLeading0s := 0
	for _, b := range src {
		if b != 0 {
			break
		}
		inLeading0s++
	}

	rawLeading0s := 0
	for _, b := range raw {
		if b != 0 {
			break
		}
		rawLeading0s++
	}

	return rawLeading0s - inLeading0s
}

// encodeRaw64 fills raw with the raw base-58 digits for a 64-byte input and
// returns the number of leading digits to skip.
//
// The accumulation uses plain uint64 arithmetic. Each product is u32×u32 so
// it fits in u64. After the first 8 input limbs a mini-reduction prevents
// overflow before adding the remaining 8 limbs (matches Firedancer).
func encodeRaw64(src *[64]byte, raw *[raw58Sz64]byte) int {
	var bin [binarySz64]uint32
	for i := range binarySz64 {
		bin[i] = binary.BigEndian.Uint32(src[i*4 : i*4+4])
	}

	var intermediate [intermediateSz64]uint64

	// First 8 limbs.
	for i := 0; i < 8; i++ {
		for k := range intermediateSz64 - 1 {
			intermediate[k+1] += uint64(bin[i]) * uint64(encTable64[i][k])
		}
	}

	// Mini-reduction to prevent overflow before the second half.
	intermediate[intermediateSz64-3] += intermediate[intermediateSz64-2] / r1div
	intermediate[intermediateSz64-2] %= r1div

	// Last 8 limbs.
	for i := 8; i < binarySz64; i++ {
		for k := range intermediateSz64 - 1 {
			intermediate[k+1] += uint64(bin[i]) * uint64(encTable64[i][k])
		}
	}

	// Full carry propagation.
	for i := intermediateSz64 - 1; i >= 1; i-- {
		intermediate[i-1] += intermediate[i] / r1div
		intermediate[i] %= r1div
	}

	for i := range intermediateSz64 {
		v := uint32(intermediate[i])
		raw[5*i+4] = byte(v % 58)
		v /= 58
		raw[5*i+3] = byte(v % 58)
		v /= 58
		raw[5*i+2] = byte(v % 58)
		v /= 58
		raw[5*i+1] = byte(v % 58)
		v /= 58
		raw[5*i+0] = byte(v)
	}

	inLeading0s := 0
	for _, b := range src {
		if b != 0 {
			break
		}
		inLeading0s++
	}

	rawLeading0s := 0
	for _, b := range raw {
		if b != 0 {
			break
		}
		rawLeading0s++
	}

	return rawLeading0s - inLeading0s
}
