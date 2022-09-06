package security

import "errors"

var (
	// ErrFailedTokenCreation indicates JWT Token failed to create, reason unknown
	ErrFailedTokenCreation = errors.New("failed to create JWT Token")
	// ErrInvalidSigningAlgorithm indicates signing algorithm is invalid, needs to be HS256, HS384, HS512, RS256, RS384 or RS512
	ErrInvalidSigningAlgorithm = errors.New("invalid signing algorithm")
	// ErrNoPrivKeyFile indicates that the given private key is unreadable
	ErrNoPrivKeyFile = errors.New("private key file unreadable")
	// ErrNoPubKeyFile indicates that the given public key is unreadable
	ErrNoPubKeyFile = errors.New("public key file unreadable")
	// ErrInvalidKeyLength occurs when a key has been used with an invalid length
	ErrInvalidKeyLength = errors.New("cipher: invalid key length")
	// ErrInvalidMessageShort occurs when a text less than BlockSize
	ErrInvalidMessageShort = errors.New("cipher: ciphertext too short")
)
