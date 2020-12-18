package sha3

import (
	"encoding/hex"

	"golang.org/x/crypto/sha3"
)

// Encrypt SHA-3 crypto
func Encrypt(str string, long int) string {
	switch long {
	case 224:
		byte28 := sha3.Sum224([]byte(str))
		return hex.EncodeToString(byte28[:])

	case 256:
		byte32 := sha3.Sum256([]byte(str))
		return hex.EncodeToString(byte32[:])

	case 384:
		byte48 := sha3.Sum384([]byte(str))
		return hex.EncodeToString(byte48[:])

	case 512:
		byte64 := sha3.Sum512([]byte(str))
		return hex.EncodeToString(byte64[:])

	default:
		byte48 := sha3.Sum256([]byte(str))
		return hex.EncodeToString(byte48[:])
	}
}
