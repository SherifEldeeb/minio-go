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

package minio

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/SherifEldeeb/minio-go/v6/pkg/s3utils"
)

// GetBucketLifecycle - get bucket lifecycle.
func (c Client) GetBucketLifecycle(bucketName string) (string, error) {
	// Input validation.
	if err := s3utils.CheckValidBucketName(bucketName); err != nil {
		return "", err
	}
	bucketLifecycle, err := c.getBucketLifecycle(bucketName)
	if err != nil {
		errResponse := ToErrorResponse(err)
		if errResponse.Code == "NoSuchLifecycleConfiguration" {
			return "", nil
		}
		return "", err
	}
	return bucketLifecycle, nil
}

// Request server for current bucket lifecycle.
func (c Client) getBucketLifecycle(bucketName string) (string, error) {
	// Get resources properly escaped and lined up before
	// using them in http request.
	urlValues := make(url.Values)
	urlValues.Set("lifecycle", "")

	// Execute GET on bucket to get lifecycle.
	resp, err := c.executeMethod(context.Background(), "GET", requestMetadata{
		bucketName:  bucketName,
		queryValues: urlValues,
	})

	defer closeResponse(resp)
	if err != nil {
		return "", err
	}

	if resp != nil {
		if resp.StatusCode != http.StatusOK {
			return "", httpRespToErrorResponse(resp, bucketName, "")
		}
	}

	bucketLifecycleBuf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	lifecycle := string(bucketLifecycleBuf)
	return lifecycle, err
}
