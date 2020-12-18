package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io/ioutil"
)

// LoadPublicKey load public key
func LoadPublicKey(publicKeyFilePath string) (*rsa.PublicKey, error) {
	byteBuffer, err := ioutil.ReadFile(publicKeyFilePath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(byteBuffer)
	if block == nil {
		return nil, errors.New("load public key error")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	publickey := pub.(*rsa.PublicKey)

	return publickey, nil
}

// LoadPrivateKey load private key
func LoadPrivateKey(privateKeyFilePath string) (*rsa.PrivateKey, error) {
	byteBuffer, err := ioutil.ReadFile(privateKeyFilePath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode([]byte(byteBuffer))
	if block == nil {
		return nil, errors.New("load private key error")
	}

	privatekey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privatekey, nil
}

// Encrypt encrypt
func Encrypt(text string, publickey *rsa.PublicKey) (string, error) {
	encryptedTextByteBuffer, err := rsa.EncryptOAEP(sha512.New(), rand.Reader, publickey, []byte(text), nil)
	if err != nil {
		return "", err
	}

	encryptedText := base64.StdEncoding.EncodeToString(encryptedTextByteBuffer)
	return encryptedText, nil
}

// Decrypt decrypt
func Decrypt(encryptedText string, privatekey *rsa.PrivateKey) (string, error) {
	encryptedTextByteBuffer, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	decryptedText, err := rsa.DecryptOAEP(sha512.New(), rand.Reader, privatekey, encryptedTextByteBuffer, nil)
	if err != nil {
		return "", err
	}

	return string(decryptedText), nil
}
