package utils

type Cipher interface {
	EncryptStandart(plaintext string) (string, error)
	DecryptStandart(src string) (string, error)
	EncryptText(text string) (string, error)
	DecryptText(text string) (string, error)
}
