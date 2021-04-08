package acrypto

import "crypto/aes"

func AesCBCEncrypt(plaintext, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return CBCEncrypt(block, plaintext, iv)
}

func AesCBCDecrypt(ciphertext, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return CBCDecrypt(block, ciphertext, iv)
}
