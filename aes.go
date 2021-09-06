package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// encryptFile encrypts the file specified by filename with the given key,
// placing the result in outFilename (or filename + ".enc" if outFilename is
// empty). The key has to be 16, 24 or 32 bytes long to select between AES-128,
// AES-192 or AES-256. Returns the name of the output file if successful.
func encryptFile(key []byte, input io.Reader, outFilename string) ([]byte, error) {
	plaintext, err := ioutil.ReadAll(input)
	if err != nil {
		return nil, err
	}

	of, err := os.Create(outFilename)
	if err != nil {
		return nil, err
	}
	defer of.Close()

	// Write the original plaintext size into the output file first, encoded in
	// a 8-byte integer.
	origSize := uint64(len(plaintext))
	if err = binary.Write(of, binary.LittleEndian, origSize); err != nil {
		return nil, err
	}

	// Pad plaintext to a multiple of BlockSize with random padding.
	if len(plaintext)%aes.BlockSize != 0 {
		bytesToPad := aes.BlockSize - (len(plaintext) % aes.BlockSize)
		padding := make([]byte, bytesToPad)
		if _, err := rand.Read(padding); err != nil {
			return nil, err
		}
		plaintext = append(plaintext, padding...)
	}

	// Generate random IV and write it to the output file.
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}
	if _, err = of.Write(iv); err != nil {
		return nil, err
	}

	// Ciphertext has the same size as the padded plaintext.
	ciphertext := make([]byte, len(plaintext))

	// Use AES implementation of the cipher.Block interface to encrypt the whole
	// file in CBC mode.
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintext)

	return ciphertext,nil
}

// decryptFile decrypts the file specified by filename with the given key. See
// doc for encryptFile for more details.
func decryptFile(key []byte, input io.Reader) ([]byte, error) {
	ciphertext, err := ioutil.ReadAll(input)
	if err != nil {
		return nil, err
	}

	// cipertext has the original plaintext size in the first 8 bytes, then IV
	// in the next 16 bytes, then the actual ciphertext in the rest of the buffer.
	// Read the original plaintext size, and the IV.
	var origSize uint64
	buf := bytes.NewReader(ciphertext)
	if err = binary.Read(buf, binary.LittleEndian, &origSize); err != nil {
		return nil, err
	}
	iv := make([]byte, aes.BlockSize)
	if _, err = buf.Read(iv); err != nil {
		return nil, err
	}

	// The remaining ciphertext has size=paddedSize.
	paddedSize := len(ciphertext) - 8 - aes.BlockSize
	if paddedSize%aes.BlockSize != 0 {
		return nil, fmt.Errorf("want padded plaintext size to be aligned to block size")
	}
	plaintext := make([]byte, paddedSize)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext[8+aes.BlockSize:])

	return plaintext[:origSize], nil
}

