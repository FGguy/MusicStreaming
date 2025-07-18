package auth

import (
	"crypto/md5"
	"encoding/hex"
)

func ValidatePassword(hashedPass string, salt string, pass string) bool {
	h := md5.Sum([]byte(pass + salt))
	return hashedPass == hex.EncodeToString(h[:])
}
