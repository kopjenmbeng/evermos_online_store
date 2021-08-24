package middleware

import (
	"crypto/sha256"
	b64 "encoding/base64"
	"encoding/binary"
	_ "fmt"
	"unicode/utf16"

	"golang.org/x/crypto/pbkdf2"
	// "golang.org/x/crypto/bcrypt"
)

func stringToUTF16Bytes(s string) []byte {
	runes := utf16.Encode([]rune(s))
	bytes := make([]byte, len(runes)*2)
	for i, r := range runes {
		binary.LittleEndian.PutUint16(bytes[i*2:], r)
	}
	return bytes
}

func HashPassword(PwdStr string, PasswordSecuritySalt string, PasswordSecurityIterations int, PasswordSecurityKeylen int) string {
	hashedPassword := pbkdf2.Key(stringToUTF16Bytes(PwdStr), stringToUTF16Bytes(PasswordSecuritySalt), PasswordSecurityIterations, PasswordSecurityKeylen, sha256.New)
	return b64.StdEncoding.EncodeToString(hashedPassword)
}
