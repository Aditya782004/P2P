package main

import (
	"bytes"
	"crypto/md5"
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
		Original: hashedStr,
	}

}

type PathKey struct {
	Pathname string
	Original string
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

func (s *StoreOpts) WriteStream(key string, r io.Reader) error {
	pathname := s.PathTransFormFunc(key)
	fmt.Printf("path is: %v", pathname)
	buf := new(bytes.Buffer)
	err := os.MkdirAll(pathname.Pathname, os.ModePerm)
	if err != nil {
		return err
	}
	io.Copy(buf, r)
	filenameBytes := md5.Sum(buf.Bytes())
	filename := hex.EncodeToString(filenameBytes[:])
	//filename := "somefilename"

	f, err := os.Create(pathname.Pathname + "/" + filename)
	if err != nil {
		return err
	}

	n, err := io.Copy(f, buf)
	if err != nil {
		return err
	}
	log.Printf("written %d bytes to the disk", n)

	return nil
}
