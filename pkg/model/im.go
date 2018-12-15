package model

import (
	"io"
)

// Copyright 2018 Dmitry Kargashin <dkargashin3@gmail.com>
// Use of this source code is governed by GNU LGPL
// license that can be found in the LICENSE file.

// IM interface describes interface for interaction with amazon s3
type IM interface {
	IsExist(key string) bool
	UploadImage(key string, body io.Reader) (string, error)
	DeleteImage(key string) error
	DownloadImage(key string) error
}
