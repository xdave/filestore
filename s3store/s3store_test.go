// Copyright 2015 Dave Gradwell
// Under BSD-style license (see LICENSE file)

package s3store_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

import (
	. "github.com/smartystreets/goconvey/convey"
)

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/s3"
)

import (
	"github.com/xdave/filestore/s3store"
)

var testbucket = "testbucket"
var testkey = "testkey"

func s3TestHandler(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			t.Log("GetObject", r.URL.Path)
		case "PUT":
			if strings.Contains(r.URL.Path, testkey) {
				t.Log("PutObject", r.URL.Path)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		default:
			t.Log("Unhandled action:", r.Method, r.URL.Path)
		}
		fmt.Fprintf(w, r.URL.Path)
	}
}

func TestS3(t *testing.T) {
	s3server := httptest.NewServer(http.HandlerFunc(s3TestHandler(t)))
	defer s3server.Close()

	config := aws.Config{
		Credentials:      credentials.NewStaticCredentials("akia", "zing", ""),
		Region:           aws.String("us-east-1"),  // This seems to be required
		Endpoint:         aws.String(s3server.URL), // Here's the magic
		S3ForcePathStyle: aws.Bool(true),           // And here
	}
	service := s3.New(&config)

	Convey("NewFileStore() should return a non-nil FileStore", t, func() {
		fs := s3store.NewFileStore(service)
		So(fs, ShouldNotBeNil)
		Convey("Bucket() should return a non-nil Bucket", func() {
			bucket := fs.Bucket(testbucket)
			So(bucket, ShouldNotBeNil)
			Convey("Bucket.Key() should return a non-nil Key", func() {
				key := bucket.Key(testkey)
				So(key, ShouldNotBeNil)
				Convey("Key.GetURI() should return a proper URI", func() {
					uri := key.GetURI()
					expected := filepath.Join(testbucket, testkey)
					So(uri, ShouldNotBeNil)
					So(uri.Path, ShouldContainSubstring, expected)
				})
				Convey("Key.Upload() should work", func() {
					testdata := []byte("This is a test")
					reader := bytes.NewReader(testdata)
					uri, err := key.Upload(reader)
					So(err, ShouldBeNil)
					So(uri, ShouldNotBeEmpty)
				})
			})
			Convey("Key.Upload() upload failure", func() {
				key := bucket.Key("")
				testdata := []byte("asdf")
				reader := bytes.NewReader(testdata)
				uri, err := key.Upload(reader)
				So(err, ShouldNotBeNil)
				So(uri, ShouldBeEmpty)
			})
		})
	})
}
