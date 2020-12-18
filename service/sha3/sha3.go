package sha3

import (
	"otter-cloud-ws/config"
	"otter-cloud-ws/pkg/sha3"
)

// Encrypt sha3 encrypt
func Encrypt(str string, lenght ...int) string {
	if len(lenght) > 0 {
		return sha3.Encrypt(str, lenght[0])
	}
	return sha3.Encrypt(str, config.Sha3Len)
}
