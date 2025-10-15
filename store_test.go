package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "mombestpicture"
	pathKey := CASPathTransformFunc(key)
	expectedPathName := "cf5d4/b01c4/d9438/c22c5/6c832/f83bd/3e8c6/304f9"
	expectedFileName := "cf5d4b01c4d9438c22c56c832f83bd3e8c6304f9"
	if pathKey.PathName != expectedPathName {
		t.Errorf("have %s want %s", pathKey.PathName, expectedPathName)
	}

	if pathKey.FileName != expectedFileName {
		t.Errorf("have %s want %s", pathKey.FileName, expectedFileName)
	}
}

// func TestStoreDeleteKey(t *testing.T) {
// 	opts := StoreOpts{
// 		PathTransformFunc: CASPathTransformFunc,
// 	}
// 	s := NewStore(opts)
// 	key := "momsspecials"
// 	data := []byte("some jpg bytes")

// 	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
// 		t.Error(err)
// 	}

// 	if err := s.Delete(key); err != nil {
// 		t.Error(err)
// 	}
// }

func TestStore(t *testing.T) {
	s := newStore()
	defer tearDown(t, s)
	for i := 0; i < 50; i++ {

		key := fmt.Sprintf("foo_%d", i)
		data := []byte("some jpg bytes")

		if _, err := s.writeStream(key, bytes.NewReader(data)); err != nil {
			t.Error(err)
		}

		if ok := s.Has(key); !ok {
			t.Errorf("expexted to have key %s", key)
		}

		r, err := s.Read(key)
		if err != nil {
			t.Error(err)
		}

		b, _ := io.ReadAll(r)

		fmt.Println(string(b))

		if string(b) != string(data) {
			t.Errorf("want %s have %s", data, b)
		}
		if err := s.Delete(key); err != nil {
			t.Error(err)
		}

		if ok := s.Has(key); ok {
			t.Errorf("expexted to Not have key %s", key)
		}
	}
}

func newStore() *Store {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	return NewStore(opts)
}

func tearDown(t *testing.T, s *Store) {
	if err := s.Clear(); err != nil {
		t.Error(err)
	}
}
