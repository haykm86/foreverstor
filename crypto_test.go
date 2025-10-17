package main

import (
	"bytes"
	"testing"
)

func TestCopyEncryptDecript(t *testing.T) {
	payload:= "Foo not Bar"
	src := bytes.NewBuffer([]byte(payload))
	dst := new(bytes.Buffer)
	key := newEncryptionKey()
	_, err := copyEncrypt(key, src, dst)
	if err != nil {
		t.Error(err)
	}

	out := new(bytes.Buffer)
	if _, err := copyDecript(key, dst, out); err != nil {
		t.Error(err)
	}

	if out.String() != payload { 
		t.Errorf("decryption failed")
	}
}
