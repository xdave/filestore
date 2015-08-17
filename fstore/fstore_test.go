// Copyright 2015 Dave Gradwell
// Under BSD-style license (see LICENSE file)

package fstore_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

import (
	. "github.com/smartystreets/goconvey/convey"
)

import (
	"github.com/xdave/filestore/fstore"
)

type FakeReadErrorCreator struct{}

func (r FakeReadErrorCreator) Read(p []byte) (int, error) {
	return 0, fmt.Errorf("Fake read error!")
}

func TestFS(t *testing.T) {
	testbucket := "testbucket"
	testkey := "testkey"
	Convey("NewFileStore() should return a non-nil FileStore", t, func() {
		store := fstore.NewFileStore()
		So(store, ShouldNotBeNil)
		tmpdir, err := ioutil.TempDir("", testbucket)
		if err != nil {
			t.Error("Wat?", err)
		}
		defer func() {
			_ = os.RemoveAll(tmpdir)
		}()
		Convey("Bucket() should return a non-nil Bucket", func() {
			bucket := store.Bucket(tmpdir)
			So(bucket, ShouldNotBeNil)
			Convey("Bucket.Key() should return a non-nil Key", func() {
				key := bucket.Key(testkey)
				So(key, ShouldNotBeNil)
				Convey("Key.GetURI() should return a proper URI", func() {
					uri := key.GetURI()
					So(uri, ShouldNotBeNil)
					So(uri.String(),
						ShouldEqual,
						filepath.Join(tmpdir, testkey))
				})
				Convey("Key.Upload() should work", func() {
					data := []byte("test data")
					buf := bytes.NewReader(data)
					location, err := key.Upload(buf)
					So(err, ShouldBeNil)
					So(location, ShouldNotBeEmpty)
				})
			})
		})
		Convey("Key.Upload() err when cannot create file", func() {
			tmpdir, err := ioutil.TempDir("", testbucket)
			if err != nil {
				t.Error("Wat?", err)
			}
			defer func() {
				_ = os.RemoveAll(tmpdir)
			}()
			bucket := store.Bucket(tmpdir)
			key := bucket.Key(testkey)
			path := key.GetURI().String()
			f, err := os.Create(path) // Pre-create destination
			if err != nil {
				t.Error("Wat?", err)
			}
			defer f.Close()
			os.Chmod(path, 0444) // Make readonly
			_, err = key.Upload(f)
			So(err, ShouldNotBeNil)
		})
		Convey("Key.Upload() err when cannot read source", func() {
			tmpdir, err := ioutil.TempDir("", testbucket)
			if err != nil {
				t.Error("Wat?", err)
			}
			defer func() {
				_ = os.RemoveAll(tmpdir)
			}()
			bucket := store.Bucket(tmpdir)
			key := bucket.Key(testkey)
			_, err = key.Upload(&FakeReadErrorCreator{}) // Some read error
			So(err, ShouldNotBeNil)
		})
	})
}
