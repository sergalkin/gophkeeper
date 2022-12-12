//go:generate mockgen -source=./crypt.go -destination=./mock/crypt.go -package=cryptmock
package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
)

type Crypter interface {
	Encode(payload string) string
	Decode(sha string) (string, error)
}

var (
	key = []byte{4, 51, 71, 14, 63, 8, 95, 100, 44, 4, 19, 85, 57, 54, 23, 54, 26, 59, 24, 44, 47, 52, 63, 1, 84, 24,
		23, 51, 3, 88, 72, 73}
	nonce = []byte{4, 51, 71, 14, 63, 8, 95, 100, 44, 4, 19, 85}
)

type crypt struct {
	aesGCM   cipher.AEAD
	aesBlock cipher.Block
	nonce    []byte
}

// NewCrypt - creates new Crypter instance.
func NewCrypt() (*crypt, error) {
	aesBlock, errBlock := aes.NewCipher(key)
	if errBlock != nil {
		return nil, fmt.Errorf("error in creating new cipher: %w", errBlock)
	}

	aesGCM, errGCM := cipher.NewGCM(aesBlock)
	if errGCM != nil {
		return nil, fmt.Errorf("error in creating GCM: %w", errGCM)
	}

	return &crypt{
		aesGCM:   aesGCM,
		aesBlock: aesBlock,
		nonce:    nonce,
	}, nil

}

// Encode - returns sha of sealed payload by aesGCM.
func (c *crypt) Encode(payload string) string {
	src := []byte(payload)

	dst := c.aesGCM.Seal(nil, c.nonce, src, nil)

	sha := hex.EncodeToString(dst)

	return sha
}

// Decode - returns decoded string by aesGCM from sha.
func (c *crypt) Decode(sha string) (string, error) {
	dst, errDecode := hex.DecodeString(sha)
	if errDecode != nil {
		return "", fmt.Errorf("hex decode error: %w", errDecode)
	}

	src, errGCM := c.aesGCM.Open(nil, c.nonce, dst, nil)
	if errGCM != nil {
		return "", errDecode
	}

	return string(src), nil
}
