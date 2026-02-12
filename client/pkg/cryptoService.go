package pkg

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"hash"
)

var (
	ErrCheckCryptoSignature = errors.New("content has not passed the digital signature verification")
	ErrSizeSignedArray      = errors.New("error in the size of the array with data and digital signature")
	ErrCryptoSignature      = errors.New("error signing data with a cryptographic key")
)

type Crypter interface {
	Start(key []byte)
	Write(src []byte)
	Sum(array []byte) []byte
	Reset()

	AddCrypter
	CheckCrypter
}

type CheckCrypter interface {
	CheckSignature(expectedSignature []byte) bool
}

type AddCrypter interface {
	AddSumToHeadOfArray(array []byte) ([]byte, error)
	AddSumToTailOfArray(array []byte) ([]byte, error)
}

type CryptoService struct {
	HMAC hash.Hash
}

func NewCryptoService() *CryptoService {
	return &CryptoService{
		HMAC: nil,
	}
}

func (c *CryptoService) Start(key []byte) {
	c.HMAC = hmac.New(sha256.New, key)
}

func (c *CryptoService) Write(src []byte) {
	c.HMAC.Write(src)
}

func (c *CryptoService) Reset() {
	c.HMAC.Reset()
}

func (c *CryptoService) Sum(array []byte) []byte {
	return c.HMAC.Sum(nil)
}

func (c *CryptoService) AddSumToHeadOfArray(array []byte) ([]byte, error) {
	if len(array) < sha256.Size {
		return nil, ErrSizeSignedArray
	}
	copy(array, c.HMAC.Sum(nil))
	return array, nil
}

func (c *CryptoService) AddSumToTailOfArray(array []byte) ([]byte, error) {
	if len(array) < sha256.Size {
		return nil, ErrSizeSignedArray
	}
	// Начальная позиция массива, куда мы хотим вставить данные.
	insertPosition := len(array) - sha256.Size
	// Вставка данных.
	copy(array[insertPosition:], c.HMAC.Sum(nil))
	return array, nil
}

func (c *CryptoService) CheckSignature(expectedSignature []byte) bool {
	return hmac.Equal(c.HMAC.Sum(nil), expectedSignature)
}

// GenKey
func GenKey(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// GenEncodeKey
func GenEncodeKey(size int) (string, error) {
	b, err := GenKey(size)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
