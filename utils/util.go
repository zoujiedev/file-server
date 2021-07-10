package utils

import (
	"crypto"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
)

//sha1
func SHA1File(path string) (string, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return "", err
	}
	h := sha1.New()
	_, err = io.Copy(h, file)
	if err != nil {
		return "", err
	}
	//return fmt.Sprintf("%x",h.Sum(nil)), nil
	return hex.EncodeToString(h.Sum(nil)), nil
}

func Sha1Bytes(data []byte) string {
	_sha1 := sha1.New()
	_sha1.Write(data)
	return hex.EncodeToString(_sha1.Sum(nil))
}

func MD5(data []byte) string {
	hash := crypto.MD5.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}
