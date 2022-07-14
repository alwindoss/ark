package vault

import (
	"bytes"
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

func NewService(repo Repository, masterPwd string) Service {
	return &fsVaultService{
		Repository:     repo,
		MasterPassword: masterPwd,
	}
}

type fsVaultService struct {
	Repository     Repository
	MasterPassword string
}

// Retrieve implements Service
func (v *fsVaultService) Retrieve(key []byte) (io.Reader, error) {
	resp, err := v.Repository.Retrieve(key)
	if err != nil {
		err = fmt.Errorf("unable to retrieve the data: %w", err)
		return nil, err
	}

	var buff bytes.Buffer
	io.Copy(&buff, resp)
	plaintext, err := Decrypt([]byte(v.MasterPassword), buff.Bytes())
	if err != nil {
		err = fmt.Errorf("unable to save the data: %w", err)
		return nil, err
	}
	respReader := bytes.NewBuffer(plaintext)

	// TODO: Before returning resp the value should be decrypted
	return respReader, nil
}

// Save implements Service
func (v *fsVaultService) Save(key []byte, value io.Reader) error {
	// TODO: Before calling Repository.Save the value should be encrypted
	var buff bytes.Buffer
	io.Copy(&buff, value)
	cipher, err := Encrypt([]byte(v.MasterPassword), buff.Bytes())
	if err != nil {
		err = fmt.Errorf("unable to save the data: %w", err)
		return err
	}
	cipherReader := bytes.NewBuffer(cipher)
	err = v.Repository.Save(key, cipherReader)
	if err != nil {
		err = fmt.Errorf("unable to save the data: %w", err)
		return err
	}

	return nil
}
