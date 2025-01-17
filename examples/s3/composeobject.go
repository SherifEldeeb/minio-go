// +build ignore

/*
 * MinIO Go Library for Amazon S3 Compatible Cloud Storage
 * Copyright 2015-2017 MinIO, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"log"

	minio "github.com/SherifEldeeb/minio-go/v6"
	"github.com/SherifEldeeb/minio-go/v6/pkg/encrypt"
)

func main() {
	// Note: YOUR-ACCESSKEYID, YOUR-SECRETACCESSKEY, my-testfile, my-bucketname and
	// my-objectname are dummy values, please replace them with original values.

	// Requests are always secure (HTTPS) by default. Set secure=false to enable insecure (HTTP) access.
	// This boolean value is the last argument for New().

	// New returns an Amazon S3 compatible client object. API compatibility (v2 or v4) is automatically
	// determined based on the Endpoint value.
	s3Client, err := minio.New("s3.amazonaws.com", "YOUR-ACCESSKEYID", "YOUR-SECRETACCESSKEY", true)
	if err != nil {
		log.Fatalln(err)
	}

	// Enable trace.
	// s3Client.TraceOn(os.Stderr)

	// Prepare source decryption key (here we assume same key to
	// decrypt all source objects.)
	decKey, _ := encrypt.NewSSEC([]byte{1, 2, 3})

	// Source objects to concatenate. We also specify decryption
	// key for each
	src1 := minio.NewSourceInfo("bucket1", "object1", decKey)
	src1.SetMatchETagCond("31624deb84149d2f8ef9c385918b653a")

	src2 := minio.NewSourceInfo("bucket2", "object2", decKey)
	src2.SetMatchETagCond("f8ef9c385918b653a31624deb84149d2")

	src3 := minio.NewSourceInfo("bucket3", "object3", decKey)
	src3.SetMatchETagCond("5918b653a31624deb84149d2f8ef9c38")

	// Create slice of sources.
	srcs := []minio.SourceInfo{src1, src2, src3}

	// Prepare destination encryption key
	encKey, _ := encrypt.NewSSEC([]byte{8, 9, 0})

	// Create destination info
	dst, err := minio.NewDestinationInfo("bucket", "object", encKey, nil)
	if err != nil {
		log.Fatalln(err)
	}

	err = s3Client.ComposeObject(dst, srcs)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Composed object successfully.")
}
