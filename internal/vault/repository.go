package vault

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Repository interface {
	Save(key []byte, value io.Reader) error
	Retrieve(key []byte) (io.Reader, error)
}

func NewFSRepository(loc string) Repository {
	return &fsVaultRepository{
		VaultDir: loc,
	}
}

type fsVaultRepository struct {
	VaultDir string
}

// Retrieve implements Repository
func (v *fsVaultRepository) Retrieve(key []byte) (io.Reader, error) {
	keyStr := sha256EncodedKey(key)
	folderName, fileName := getFolderAndFileName(keyStr)

	folderPath := filepath.Join(v.VaultDir, folderName)
	err := os.MkdirAll(folderPath, 0766)
	if err != nil {
		err = fmt.Errorf("error when creating directory %s: %w", folderPath, err)
		return nil, err
	}
	filePath := filepath.Join(folderPath, fileName)
	f, err := os.Open(filePath)
	if err != nil {
		err = fmt.Errorf("error when opening file %s: %w", filePath, err)
		return nil, err
	}
	defer f.Close()
	zr, err := zlib.NewReader(f)
	if err != nil {
		err = fmt.Errorf("error reading the compressed value: %w", err)
		return nil, err
	}
	defer zr.Close()
	var buff bytes.Buffer
	io.Copy(&buff, zr)

	return &buff, nil
}

// Save implements Repository
func (v *fsVaultRepository) Save(key []byte, value io.Reader) error {
	keyStr := sha256EncodedKey(key)
	folderName, fileName := getFolderAndFileName(keyStr)

	folderPath := filepath.Join(v.VaultDir, folderName)
	err := os.MkdirAll(folderPath, 0766)
	if err != nil {
		err = fmt.Errorf("error when creating directory %s: %w", folderPath, err)
		return err
	}
	filePath := filepath.Join(folderPath, fileName)
	f, err := os.Create(filePath)
	if err != nil {
		err = fmt.Errorf("error when creating file %s: %w", filePath, err)
		return err
	}
	defer f.Close()
	zw := zlib.NewWriter(f)
	defer zw.Close()
	_, err = io.Copy(zw, value)
	if err != nil {
		err = fmt.Errorf("error compressing the value: %w", err)
		return err
	}
	// log.Printf("number of bytes written to the file %s is %d", filePath, written)

	return nil
}
