package vault

import (
	"fmt"
	"io"
)

type Encrypter interface {
	Encrypt(plaintext []byte) ([]byte, error)
}

type Decrypter interface {
	Decrypt(ciphertext []byte) ([]byte, error)
}

type EncryptDecrypter interface {
	Encrypter
	Decrypter
}

type Service interface {
	Save(key []byte, value io.Reader) error
	Retrieve(key []byte) (io.Reader, error)
}

func NewService(repo Repository) Service {
	return &fsVaultService{
		Repository: repo,
	}
}

type fsVaultService struct {
	Repository Repository
}

// Retrieve implements Service
func (v *fsVaultService) Retrieve(key []byte) (io.Reader, error) {
	resp, err := v.Repository.Retrieve(key)
	if err != nil {
		err = fmt.Errorf("unable to retrieve the data: %w", err)
		return nil, err
	}

	// TODO: Before returning resp the value should be decrypted
	return resp, nil
}

// Save implements Service
func (v *fsVaultService) Save(key []byte, value io.Reader) error {
	// TODO: Before calling Repository.Save the value should be encrypted
	err := v.Repository.Save(key, value)
	if err != nil {
		err = fmt.Errorf("unable to save the data: %w", err)
		return err
	}

	return nil
}
