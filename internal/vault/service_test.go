package vault_test

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/alwindoss/ark/internal/vault"
)

func TestVaultService(t *testing.T) {
	loc := t.TempDir()
	repo := vault.NewFSRepository(loc)
	svc := vault.NewService(repo)
	key := []byte("key1")
	valueStr := "This is the value1 for key1"
	value := strings.NewReader(valueStr)
	err := svc.Save(key, value)
	if err != nil {
		t.Logf("expected error to be nil but was %v", err)
		t.FailNow()
	}
	respReader, err := svc.Retrieve(key)
	if err != nil {
		t.Logf("expected error to be nil but was %v", err)
		t.FailNow()
	}
	respData, err := ioutil.ReadAll(respReader)
	if err != nil {
		t.Logf("expected error to be nil but was %v", err)
		t.FailNow()
	}
	if string(respData) != valueStr {
		t.Logf("expected the retrieved value to be '%s', but was '%s'", valueStr, string(respData))
		t.FailNow()
	}

}
