// Copyright 2018 Dmitry Kargashin <dkargashin3@gmail.com>
// Use of this source code is governed by GNU LGPLL
// license that can be found in the LICENSE file.

// utility.go contains some utility functions which are used by API handlers.

package api

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/nfnt/resize"

	"bmstu.codes/developers34/SBWeb/pkg/model"
)

// getIDfromCookie returns ID of user using cookie from request.
// This function must be used with checkSessionMiddleware because
// it doesn't handle any errors.
func getIDfromCookie(m *model.Model, r *http.Request) int64 {
	cookieSession, _ := r.Cookie("session_id")
	session, _ := m.CheckSession(&model.SessionID{
		ID: cookieSession.Value,
	})
	return session.ID
}

// loadImages process incoming request to upload images from it.
// ParseMultipartFrom must called before this function.
// It returns array of image's paths which were created.
func loadImages(r *http.Request) ([]string, error) {
	files := r.MultipartForm.File["images"]

	filenames := make([]string, 0)

	for _, item := range files {
		// open file header
		file, err := item.Open()
		defer file.Close()
		if err != nil {
			return nil, err
		}

		// make buf for first 512 bytes of file
		buf := make([]byte, 512)
		_, err = file.Read(buf)
		if err != nil {
			return nil, err
		}

		// detect type of file from request
		typeOfFile := http.DetectContentType(buf)
		if typeOfFile != "image/png" && typeOfFile != "image/jpeg" {
			return nil, errors.New("Trying to upload wrong file extension")
		}

		// decode image
		file.Seek(0, 0)
		var img image.Image

		if typeOfFile == "image/png" {
			img, err = png.Decode(file)
		} else {
			img, err = jpeg.Decode(file)
		}
		if err != nil {
			return nil, err
		}

		// resize image
		imgNew := resize.Thumbnail(1280, 720, img, resize.Lanczos3)

		// make md5 sum for file name
		hasher := md5.New()
		io.Copy(hasher, file)
		dataTime, _ := time.Now().MarshalBinary()
		randBuf := make([]byte, 32)
		rand.Read(randBuf)
		io.Copy(hasher, bytes.NewReader(dataTime))
		io.Copy(hasher, bytes.NewReader(randBuf))
		filename := hex.EncodeToString(hasher.Sum(nil))
		dst, err := os.Create("./images/" + filename + ".png")
		defer dst.Close()

		// save image
		err = png.Encode(dst, imgNew)
		if err != nil {
			return nil, err
		}

		filenames = append(filenames, "/images/"+filename+".png")

		// if we create or update user we need only one file
		if strings.Contains(r.URL.Path, "/users/") {
			break
		}
	}

	return filenames, nil
}

// deleteImages deletes files with filenames.
func deleteImages(filenames []string) error {
	for _, filename := range filenames {
		if filename == "" {
			break
		}
		err := os.Remove("." + filename)
		if err != nil {
			return err
		}
	}
	return nil
}
