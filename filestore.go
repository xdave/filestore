// Copyright 2015 Dave Gradwell
// Under BSD-style license (see LICENSE file)

package filestore

import (
	"io"
	"net/url"
)

type FileStore interface {
	Bucket(name string) Bucket
}

type Bucket interface {
	Key(name string) Key
}

type Key interface {
	GetURI() *url.URL
	Upload(io.Reader) (string, error)
}
