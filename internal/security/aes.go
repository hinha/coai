package security

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	b64 "encoding/base64"
	"errors"
	"github.com/hinha/coai/internal"
	"io"
)

var (
	// ErrEncryptShort indicates data is to short
	ErrEncryptShort = errors.New("encrypt too short")
	// KeySize for aes uses the generic key size
	KeySize = 32
)

type bearerCipher struct {
	key string
}

func NewCipher(key string) internal.Cipher {
	c := &bearerCipher{}
	c.key = string(c.HashTo32Bytes(key))
	return c
}

// HashTo32Bytes will compute a cryptographically useful hash of the input string.
func (b *bearerCipher) HashTo32Bytes(input string) []byte {
	data := sha256.Sum256([]byte(input))
	return data[0:]
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

// EncryptStandart to encrypt aes using cbc stantard
func (b *bearerCipher) EncryptStandart(plaintext string) (string, error) {
	if len(plaintext) == 0 {
		return "", ErrEncryptShort
	}

	block, err := aes.NewCipher([]byte(b.key))
	if err != nil {
		return "", err
	}

	ecb := cipher.NewCBCEncrypter(block, []byte(b.key)[16:])
	content := []byte(plaintext)
	content = pkcs5Padding(content, block.BlockSize())
	crypted := make([]byte, len(content))
	ecb.CryptBlocks(crypted, content)

	return b64.URLEncoding.EncodeToString(crypted), nil
}

func pkcs5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}

func (b *bearerCipher) DecryptStandart(src string) (string, error) {
	if len(src) == 0 {
		return "", ErrEncryptShort
	}

	enc, err := b64.URLEncoding.DecodeString(src)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(b.key))
	if err != nil {
		return "", err
	}

	ecb := cipher.NewCBCDecrypter(block, []byte(b.key)[16:])
	decrypted := make([]byte, len(enc))
	ecb.CryptBlocks(decrypted, enc)

	return string(pkcs5Trimming(decrypted)), nil
}

// EncryptText encrypt cbc plaintext using the given key with CTR encryption
func (b *bearerCipher) EncryptText(text string) (string, error) {
	if (len(b.key) < KeySize || len(b.key) > KeySize) || len(text) <= 0 {
		return "", ErrInvalidKeyLength
	}

	c, _ := aes.NewCipher([]byte(b.key))
	ct := make([]byte, aes.BlockSize+len(text))
	iv := ct[:aes.BlockSize]

	_, err := io.ReadFull(rand.Reader, iv)
	if err != nil {
		return "", err
	}

	cfb := cipher.NewCFBEncrypter(c, iv)
	cfb.XORKeyStream(ct[aes.BlockSize:], []byte(text))

	return b.EncodeString(ct), nil
}

// DecryptText dencrypt cbc plaintext using the given key with CTR encryption
func (b *bearerCipher) DecryptText(text string) (string, error) {
	ciphertext, err := b.DecodingString(text)
	if err != nil {
		return "", err
	}

	if (len(b.key) < KeySize || len(b.key) > KeySize) || len(string(ciphertext)) <= 0 {
		return "", ErrInvalidKeyLength
	}

	if len(string(ciphertext)) < aes.BlockSize {
		return "", ErrInvalidMessageShort
	}

	c, _ := aes.NewCipher([]byte(b.key))
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(c, iv)
	cfb.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

func (c *bearerCipher) EncodeString(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

func (c *bearerCipher) DecodingString(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}
