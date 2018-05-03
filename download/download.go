// Copyright (c) 2018 ChrisOboe
//
// This file is part of ti-dl
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package download

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"github.com/ChrisOboe/tidapi"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"os"
	"strings"
)

type Download struct {
	quality     tidapi.Audioquality
	api         tidapi.Tidal
	destination string
	state       chan string
}

func New(state chan string, quality tidapi.Audioquality, api tidapi.Tidal, destination string) Download {
	return Download{quality, api, destination, state}
}

// helpers
func stripBadChars(in string) string {
	in = strings.Replace(in, "/", "_", -1)
	in = strings.Replace(in, ":", "_", -1)
	return strings.Replace(in, "\\", "_", -1)
}

func download(url string, file string, key []byte, iv []byte) error {
	// file stuff
	handler, err := os.Create(file)
	if err != nil {
		return errors.Wrap(err, "Couldn't create "+file)
	}
	defer handler.Close()

	// web stuff
	response, err := http.Get(url)
	if err != nil {
		return errors.Wrap(err, "Error on getting "+url)
	}
	defer response.Body.Close()

	var source io.Reader

	// crypto stuff
	if len(key) != 0 {
		block, err := aes.NewCipher(key)
		if err != nil {
			return errors.Wrap(err, "Couldn't init cipher.")
		}
		stream := cipher.NewCTR(block, iv)
		source = &cipher.StreamReader{S: stream, R: response.Body}
	} else {
		source = response.Body
	}

	_, err = io.Copy(handler, source)
	if err != nil {
		return errors.Wrap(err, "Error while downloading "+url)
	}

	return nil
}

func (d Download) trackDownload(trackId int, file string) error {
	resp, err := d.api.TracksUrlpostpaywall(trackId, d.quality, tidapi.UrlUsageModeOffline)
	if err != nil {
		return errors.Wrap(err, "Couldn't download track "+string(trackId))
	}

	switch resp.Codec {
	case "AAC":
		file += ".mp4"
	case "FLAC":
		file += ".flac"
	default:
		return errors.New(fmt.Sprintf("Codec %s doesn't have a known file extension.", resp.Codec))
		file += ".???"
	}

	var key []byte
	var iv []byte

	// some files aren't encrypted
	if resp.SecurityToken != "" {
		key, iv, err = fileKey(resp.SecurityToken)
		if err != nil {
			return errors.Wrap(err, "Getting decryption key for track "+string(trackId)+"failed.")
		}
	}

	err = download(resp.Urls[0], file, key, iv)
	if err != nil {
		return errors.Wrap(err, "Couldn't download track "+string(trackId))
	}

	return nil
}
