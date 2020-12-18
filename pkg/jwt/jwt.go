package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"hash"
	"strings"
)

// AlgTyp jwt alg type
type AlgTyp string

const (
	// HS256 encrypt alg
	HS256 AlgTyp = "HS256"
	// HS384 encrypt alg
	HS384 AlgTyp = "HS384"
	// HS512 encrypt alg
	HS512 AlgTyp = "HS512"
)

// GenerateJWT generate jwt
func GenerateJWT(payload interface{}, alg, key string) (string, error) {
	// header
	jwtHeader := base64.StdEncoding.EncodeToString([]byte(`{"typ":"JWT","alg":"` + alg + `"}`))

	// payload
	bytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	jwtPayload := base64.StdEncoding.EncodeToString(bytes)

	// signature
	jwtSignature := encryptSignature(jwtHeader+"."+jwtPayload, alg, key)

	return jwtHeader + "." + jwtPayload + "." + jwtSignature, nil
}

// VerifyJWT verify jwt
func VerifyJWT(jwt, alg, key string) ([]byte, bool) {
	// jwts[0]=header, jwts[1]=payload, jwts[2]=signature
	jwts := strings.Split(jwt, ".")
	if len(jwts) != 3 {
		return nil, false
	}

	// signature
	jwtSignature := encryptSignature(jwts[0]+"."+jwts[1], alg, key)
	if jwts[2] != jwtSignature {
		return nil, false
	}

	bytes, err := base64.StdEncoding.DecodeString(jwts[1])
	if err != nil {
		return nil, false
	}
	return bytes, true
}

func encryptSignature(str, alg, key string) string {
	var h hash.Hash
	switch alg {
	case "HS256":
		h = hmac.New(sha256.New, []byte(key))
	case "HS384":
		h = hmac.New(sha512.New384, []byte(key))
	case "HS512":
		h = hmac.New(sha512.New, []byte(key))
	default:
		h = hmac.New(sha256.New, []byte(key))
	}
	h.Write([]byte(str))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
