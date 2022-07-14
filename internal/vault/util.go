package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/crypto/scrypt"
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

func Encrypt(pwdKey, data []byte) ([]byte, error) {
	pwdKey, salt, err := DeriveKey(pwdKey, nil)
	if err != nil {
		return nil, err
	}

	blockCipher, err := aes.NewCipher(pwdKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	ciphertext = append(ciphertext, salt...)

	return ciphertext, nil
}

func Decrypt(pwdKey, data []byte) ([]byte, error) {
	salt, data := data[len(data)-32:], data[:len(data)-32]

	pwdKey, _, err := DeriveKey(pwdKey, salt)
	if err != nil {
		return nil, err
	}

	blockCipher, err := aes.NewCipher(pwdKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func DeriveKey(password, salt []byte) ([]byte, []byte, error) {
	if salt == nil {
		salt = make([]byte, 32)
		if _, err := rand.Read(salt); err != nil {
			return nil, nil, err
		}
	}

	key, err := scrypt.Key(password, salt, 1048576, 8, 1, 32)
	if err != nil {
		return nil, nil, err
	}

	return key, salt, nil
}
