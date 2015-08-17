// Copyright 2015 Dave Gradwell
// Under BSD-style license (see LICENSE file)

package s3store

import (
	"io"
	"net/url"
	"time"
)

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

import (
	"github.com/xdave/filestore"
)

type FileStore struct {
	s3 s3iface.S3API
}

func NewFileStore(s3interface s3iface.S3API) *FileStore {
	return &FileStore{
		s3: s3interface,
	}
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

func (k *Key) GetURI() *url.URL {
	req, _ := k.bucket.store.s3.GetObjectRequest(&s3.GetObjectInput{
		Bucket: &k.bucket.name,
		Key:    &k.name,
	})
	str, _ := req.Presign(time.Hour)
	uri, _ := url.Parse(str)
	return uri
}

func (k *Key) Upload(body io.Reader) (string, error) {
	opts := s3manager.UploadOptions{S3: k.bucket.store.s3.(*s3.S3)}
	uploader := s3manager.NewUploader(&opts)
	input := s3manager.UploadInput{
		Bucket: &k.bucket.name,
		Key:    &k.name,
		Body:   body,
	}
	out, err := uploader.Upload(&input)
	if err != nil {
		return "", err
	}
	return out.Location, nil
}
