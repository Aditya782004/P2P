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

var DefaultPathTransportFunc = func(key string) PathKey {
	return PathKey{
		Pathname: key,
		Filename: key,
	}
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
	Root              string
	PathTransFormFunc PathTransFormFunc
}

type Store struct {
	StoreOpts
}

const DefaultRootDir = "root"

func NewStore(opts StoreOpts) *Store {
	if opts.PathTransFormFunc == nil {
		opts.PathTransFormFunc = DefaultPathTransportFunc
	}
	if len(opts.Root) == 0 {
		opts.Root = DefaultRootDir
	}
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) Delete(key string) error {
	fmt.Println("key is: ", key)
	pathkey := s.PathTransFormFunc(key)
	// if len(pathkey.Fullpath()) == 0 {
	// 	return nil
	// }
	fmt.Println("Full path:", pathkey.Fullpath())
	//fullStr := fmt.Sprint("%s/%s", s.Root, pathkey.Fullpath())
	SplitStr := strings.Split(pathkey.Fullpath(), "/")
	err := os.RemoveAll(SplitStr[0])
	if err != nil {
		return err
	}
	defer func() {
		fmt.Printf("deleted [%s] from the disk", pathkey.Fullpath())
	}()
	return nil
}

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
	//fullPathWithStr := fmt.Sprintf("%s/%s", s.Root, pathkey.Fullpath())
	return os.Open(pathkey.Fullpath())

}

func (p PathKey) Fullpath() string {
	return fmt.Sprintf("%s/%s/%s", DefaultRootDir, p.Pathname, p.Filename)
}

func (s *Store) WriteStream(key string, r io.Reader) error {
	pathname := s.PathTransFormFunc(key)
	fmt.Printf("path is: %v", pathname)
	buf := new(bytes.Buffer)
	err := os.MkdirAll(s.Root+"/"+pathname.Pathname, os.ModePerm)
	if err != nil {
		return err
	}
	io.Copy(buf, r)

	//fullpath := pathname.Fullpath()
	//filename := "somefilename"

	f, err := os.Create(pathname.Fullpath())
	if err != nil {
		return err
	}
	defer f.Close()
	n, err := io.Copy(f, buf)
	if err != nil {
		return err
	}
	log.Printf("written %d bytes to the disk %s", n, pathname.Fullpath())

	return nil
}
