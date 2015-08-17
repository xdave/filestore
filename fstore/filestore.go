// Copyright 2015 Dave Gradwell
// Under BSD-style license (see LICENSE file)

package fstore

import (
	"io"
	"net/url"
	"os"
	"path/filepath"
)
import (
	"github.com/xdave/filestore"
)

type FileStore struct {
}

func NewFileStore() filestore.FileStore {
	return &FileStore{}
}

func (fs *FileStore) Bucket(name string) filestore.Bucket {
	return &Bucket{
		store: fs,
		name:  name,
	}
}

type Bucket struct {
	store *FileStore
	name  string
}

func (b *Bucket) Key(name string) filestore.Key {
	return &Key{
		bucket: b,
		name:   name,
	}
}

type Key struct {
	bucket *Bucket
	name   string
}

// Returns the full path of the file
func (k *Key) GetURI() *url.URL {
	path := filepath.Join(k.bucket.name, k.name)
	uri, _ := url.Parse(path)
	return uri
}

// Functions effectively as a file copy
func (k *Key) Upload(body io.Reader) (path string, err error) {
	path = k.GetURI().String()
	dir := filepath.Dir(path)
	_ = os.MkdirAll(dir, 0777)
	dest, err := os.Create(path)
	if err != nil {
		return
	}
	defer dest.Close()
	_, err = io.Copy(dest, body)
	if err != nil {
		return
	}
	return
}
