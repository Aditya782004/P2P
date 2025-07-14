package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestHash(t *testing.T) {
	key := "sampleKey"
	Pathname := CASPathTransformFunc(key)
	expectedOriginalKey := "a421c5f473312f80f7984da72bc9c05a7b48cf37"
	ExpectedPathname := "a421c5f473/312f80f798/4da72bc9c0/5a7b48cf37"
	if Pathname.Pathname != ExpectedPathname {
		fmt.Printf("have %s want %s ", Pathname.Pathname, ExpectedPathname)
	}
	if Pathname.Original != expectedOriginalKey {
		fmt.Printf("have %s want %s ", Pathname.Original, expectedOriginalKey)
	}

}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransFormFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)
	//key := "sampleKey"
	data := bytes.NewReader([]byte("some data"))
	err := s.WriteStream("specialPic", data)
	if err != nil {
		t.Errorf("some error occured")
	}
}
