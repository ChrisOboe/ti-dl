// Copyright (c) 2018-2019 ChrisOboe <chris@oboe.email>
// SPDX-License-Identifier: GPL-3.0

package download

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"github.com/ChrisOboe/ti-dl/settings"
	"github.com/ChrisOboe/tidapi"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Download struct {
	quality  tidapi.Audioquality
	api      tidapi.Tidal
	settings settings.Settings
	state    chan string
}

func New(state chan string, quality tidapi.Audioquality, api tidapi.Tidal, settings settings.Settings) Download {
	return Download{quality, api, settings, state}
}

// helpers
func stripBadChars(in string) string {
	in = strings.Replace(in, "/", "_", -1)
	in = strings.Replace(in, ":", "_", -1)
	return strings.Replace(in, "\\", "_", -1)
}

func download(url string, file string, key []byte, iv []byte) error {
	// check if file already exists
	_, err := os.Stat(file)
	if !os.IsNotExist(err) {
		return errors.New("File already exists")
	}

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

func (d Download) trackDownload(trackId int, file string) (string, error) {
	resp, err := d.api.TracksUrlpostpaywall(trackId, d.quality, tidapi.UrlUsageModeOffline)
	if err != nil {
		return "", errors.Wrap(err, "Couldn't download track "+strconv.Itoa(trackId))
	}

	switch resp.Codec {
	case "AAC":
		file += ".m4a"
	case "FLAC":
		file += ".flac"
	default:
		return "", errors.New(fmt.Sprintf("Codec %s doesn't have a known file extension.", resp.Codec))
	}

	var key []byte
	var iv []byte

	// some files aren't encrypted
	if resp.SecurityToken != "" {
		key, iv, err = fileKey(resp.SecurityToken)
		if err != nil {
			return "", errors.Wrap(err, "Getting decryption key for track "+string(trackId)+"failed.")
		}
	}

	err = download(resp.Urls[0], file, key, iv)
	if err != nil {
		return "", errors.Wrap(err, "Couldn't download track "+string(trackId))
	}

	/*
		if d.settings.Remux {
			ffmpeg := exec.Command("ffmpeg", "-y", "-v", "quiet", "-stats", "-i", file, "-codec", "copy", fileNoExt+".mka")

			fmt.Print("Remuxing to mka: ")
			stderr, _ := ffmpeg.StderrPipe()
			err = ffmpeg.Start()

			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				fmt.Println(scanner.Text())
			}
			ffmpeg.Wait()

			if err != nil {
				return errors.Wrap(err, "Couldn't execute \""+ffmpeg.String()+"\"")
			}
			err = os.Remove(file)
			if err != nil {
				return errors.Wrap(err, "Couldn't delete file.")
			}
		}
	*/

	if d.settings.Tag {
		titag := exec.Command("ti-tag", d.settings.Username, d.settings.Password, "track", strconv.Itoa(trackId), file)

		fmt.Println("Tagging file")
		stdout, _ := titag.StdoutPipe()
		err := titag.Start()

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		titag.Wait()

		if err != nil {
			return "", errors.Wrap(err, "Couldn't execute \""+titag.String()+"\"")
		}
	}

	return file, nil
}
