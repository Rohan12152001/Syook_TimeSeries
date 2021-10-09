package tmp

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"github.com/itrepablik/itrlog"
)

// Initialize vector, which is the random bytes
var iv = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

// Keep this secret key with you.
const secretKey string = "abc&1*~#^2^#s0^=)^^7%b34"

func encodeBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func decodeBase64(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

// Encrypt method is to encrypt or hide any classified text
func Encrypt(text, secretKey string) (string, error) {
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", err
	}
	plainText := []byte(text)
	fmt.Println(">>", len([]byte(secretKey)), block.BlockSize(), len(iv))

	cfb := cipher.NewCTR(block, iv)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return encodeBase64(cipherText), nil


}

// Decrypt method is to extract back the encrypted text
func Decrypt(text, secretKey string) (string, error) {
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", err
	}
	cipherText := decodeBase64(text)
	cfb := cipher.NewCTR(block, iv)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}

func main() {
	fmt.Println("Hello Maharlikans!")

	phrase := "Hello World!"

	// To encrypt the phrase
	encText, err := Encrypt(phrase, secretKey)
	if err != nil {
		itrlog.Fatalw("error encrypting your classified text: ", err)
	}
	fmt.Println("encrypted text: ", encText)

	// To decrypt the original phrase
	decText, err := Decrypt(encText, secretKey)
	if err != nil {
		itrlog.Fatalw("error decrypting your encrypted text: ", err)
	}
	fmt.Println("decrypted text: ", decText)
}
