package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type PathTransFormFunc func(string) PathKey

var DefaultPathTransportFunc = func(key string) string {
	return key
}

func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashedStr := hex.EncodeToString(hash[:])
	blockSize := 10

	lenSlice := len(hashedStr) / blockSize
	paths := make([]string, lenSlice)
	for i := 0; i < lenSlice; i++ {
		from, to := i*blockSize, (i*blockSize)+blockSize
		paths[i] = hashedStr[from:to]
	}
	return PathKey{
		Pathname: strings.Join(paths, "/"),
		Filename: hashedStr,
	}

}

type PathKey struct {
	Pathname string
	Filename string
}

type StoreOpts struct {
	PathTransFormFunc PathTransFormFunc
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) Delete() {}

func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)
	return buf, err
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathkey := s.PathTransFormFunc(key)
	return os.Open(pathkey.Fullpath())

}

func (p PathKey) Fullpath() string {
	return fmt.Sprintf("%s/%s", p.Pathname, p.Filename)
}

func (s *Store) WriteStream(key string, r io.Reader) error {
	pathname := s.PathTransFormFunc(key)
	fmt.Printf("path is: %v", pathname)
	buf := new(bytes.Buffer)
	err := os.MkdirAll(pathname.Pathname, os.ModePerm)
	if err != nil {
		return err
	}
	io.Copy(buf, r)

	fullpath := pathname.Fullpath()
	//filename := "somefilename"

	f, err := os.Create(fullpath)
	if err != nil {
		return err
	}

	n, err := io.Copy(f, buf)
	if err != nil {
		return err
	}
	log.Printf("written %d bytes to the disk %s", n, fullpath)

	return nil
}
