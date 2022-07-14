package vault

import (
	"crypto/sha256"
	"encoding/hex"
)

func sha256EncodedKey(key []byte) string {
	hash := sha256.New()
	hash.Write(key)
	hashedKey := hash.Sum(nil)
	keyStr := hex.EncodeToString(hashedKey)
	return keyStr
}

func getFolderAndFileName(keyStr string) (folderName, fileName string) {
	folderName = keyStr[:2]
	fileName = keyStr[2:]
	return
}
