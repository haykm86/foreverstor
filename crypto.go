package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
)

func generateID() string {
	buf := make([]byte, 32)
	io.ReadFull(rand.Reader, buf)
	return hex.EncodeToString(buf)
}

func hashKey(key string) string {
	hash := md5.Sum([]byte(key))
	return hex.EncodeToString(hash[:])
}

func newEncryptionKey() []byte {
	keyBuf := make([]byte, 32)
	io.ReadFull(rand.Reader, keyBuf)
	return keyBuf
}

func copyDecript(key []byte, src io.Reader, dst io.Writer) (int, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return 0, err
	}

	// Read IV from the begining of our src
	iv := make([]byte, block.BlockSize()) // 16 bytes
	if _, err := io.ReadFull(src, iv); err != nil {
		return 0, err
	}

	var (
		buf    = make([]byte, 32*1024)
		stream = cipher.NewCTR(block, iv)
		total  = 0
	)

	for {
		n, err := src.Read(buf)
		if n > 0 {
			stream.XORKeyStream(buf[:n], buf[:n])
			written, writeErr := dst.Write(buf[:n])
			if writeErr != nil {
				return total, writeErr
			}
			if written != n {
				return total + written, io.ErrShortWrite
			}
			total += written
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return total, err
		}
	}

	return total, nil
}

func copyEncrypt(key []byte, src io.Reader, dst io.Writer) (int, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return 0, err
	}

	iv := make([]byte, block.BlockSize()) // 16 bytes
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return 0, err
	}

	// prepend the IV to the file.
	total, err := dst.Write(iv)
	if err != nil {
		return 0, err
	}

	var (
		buf    = make([]byte, 32*1024)
		stream = cipher.NewCTR(block, iv)
	)

	for {
		n, err := src.Read(buf)
		if n > 0 {
			stream.XORKeyStream(buf[:n], buf[:n])
			written, writeErr := dst.Write(buf[:n])
			if writeErr != nil {
				return total, writeErr
			}
			if written != n {
				return total + written, io.ErrShortWrite
			}
			total += written
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return total, err
		}
	}

	return total, nil
}
