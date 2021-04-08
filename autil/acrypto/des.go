package acrypto

import "crypto/des"

func DesCBCEncrypt(plaintext, key, iv []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return CBCEncrypt(block, plaintext, iv)
}

func DesCBCDecrypt(ciphertext, key, iv []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return CBCDecrypt(block, ciphertext, iv)
}
