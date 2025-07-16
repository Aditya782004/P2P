package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

// func TestHash(t *testing.T) {
// 	key := "sampleKey"
// 	Pathname := CASPathTransformFunc(key)
// 	expectedOriginalKey := "a421c5f473312f80f7984da72bc9c05a7b48cf37"
// 	ExpectedPathname := "a421c5f473/312f80f798/4da72bc9c0/5a7b48cf37"
// 	if Pathname.Pathname != ExpectedPathname {
// 		fmt.Printf("have %s want %s ", Pathname.Pathname, ExpectedPathname)
// 	}
// 	if Pathname.Filename != expectedOriginalKey {
// 		fmt.Printf("have %s want %s ", Pathname.Filename, expectedOriginalKey)
// 	}

// }

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransFormFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)
	key := "sampleKey"
	data := []byte("some data")
	err := s.WriteStream(key, bytes.NewReader(data))
	if err != nil {
		t.Errorf("some error occured")
	}
	r, err := s.Read(key)
	if err != nil {
		t.Error(err)
	}

	b, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("io.ReadAll failed: %v", err)
	}
	fmt.Println(string(b))
	if string(b) != string(data) {
		t.Errorf("have %s wnat %s", b, data)
	}
	errr := s.Delete(key)
	if errr != nil {
		t.Error(errr)
	}
}
