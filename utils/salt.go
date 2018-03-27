package utils

import (
	"crypto/md5"
	"fmt"
)

const Salt = "lzx"

func SaltPassword(pass_md5 string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(Salt + pass_md5)))
}

func IsLowerMD5(str string) bool {
	if len(str) != 32 {
		return false
	}

	for _, val := range []byte(str) {
		if val >= '0' && val <= '9'  || val >= 'a' && val <= 'f'{
			continue
		}
		return false
	}

	return true
}
