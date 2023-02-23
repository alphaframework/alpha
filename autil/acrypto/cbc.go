package acrypto

import "crypto/cipher"

func CBCEncrypt(block cipher.Block, plaintext, iv []byte) ([]byte, error) {
	blockSize := block.BlockSize()
	plaintext = PKCS7Padding(plaintext, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	ciphertext := make([]byte, len(plaintext))
	blockMode.CryptBlocks(ciphertext, plaintext)
	return ciphertext, nil
}

func CBCDecrypt(block cipher.Block, ciphertext, iv []byte) ([]byte, error) {
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, iv[:blockSize])
	plaintext := make([]byte, len(ciphertext))
	blockMode.CryptBlocks(plaintext, ciphertext)
	plaintext = PKCS7UnPadding(plaintext)
	return plaintext, nil
}
