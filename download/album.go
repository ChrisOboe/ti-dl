// Copyright (c) 2018-2019 ChrisOboe <chris@oboe.email>
// SPDX-License-Identifier: GPL-3.0

package download

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func (d Download) Album(albumId int) error {
	album, err := d.api.Album(albumId)
	if err != nil {
		return errors.Wrap(err, "Couldn't get album informations.")
	}

	items, err := d.api.AlbumItems(albumId)
	if err != nil {
		return errors.Wrap(err, "Couldn't get album items.")
	}
	d.state <- fmt.Sprintf("Downloading album: %s (%d)", album.Title, albumId)

	_albumartist := stripBadChars(album.Artist.Name)
	_album := stripBadChars(album.Title)
	_releaseDate := stripBadChars(album.ReleaseDate)
	_releaseType := stripBadChars(album.Type)

	var downloaded []string
	for _, item := range items.Items {
		d.state <- fmt.Sprintf("Downloading track: %s (%d)", item.Item.Title, item.Item.ID)

		_tracknumber := strconv.Itoa(item.Item.TrackNumber)
		_title := stripBadChars(item.Item.Title)

		filePath := strings.Replace(d.settings.Destination, "${ALBUMARTIST}", _albumartist, -1)
		filePath = strings.Replace(filePath, "${RELEASETITLE}", _album, -1)
		filePath = strings.Replace(filePath, "${TRACKNUMBER}", _tracknumber, -1)
		filePath = strings.Replace(filePath, "${TRACKTITLE}", _title, -1)
		filePath = strings.Replace(filePath, "${RELEASEDATE}", _releaseDate, -1)
		filePath = strings.Replace(filePath, "${RELEASETYPE}", _releaseType, -1)

		err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
		if err != nil {
			return errors.Wrap(err, "error creating "+filepath.Dir(filePath))
		}
		file, err := d.trackDownload(item.Item.ID, filePath)
		if err != nil {
			fmt.Println(err)
		}
		downloaded = append(downloaded, file)
	}

	albumPath := strings.Replace(d.settings.Destination, "${ALBUMARTIST}", _albumartist, -1)
	albumPath = strings.Replace(albumPath, "${RELEASETITLE}", _album, -1)
	albumPath = strings.Replace(albumPath, "${TRACKNUMBER}", "00", -1)
	albumPath = strings.Replace(albumPath, "${TRACKTITLE}", "XYZ", -1)
	albumPath = strings.Replace(albumPath, "${RELEASEDATE}", _releaseDate, -1)
	albumPath = strings.Replace(albumPath, "${RELEASETYPE}", _releaseType, -1)

	if d.settings.Tag {
		args := []string{d.settings.Username, d.settings.Password, "album", strconv.Itoa(albumId), filepath.Dir(albumPath) + "/album.nfo", filepath.Dir(albumPath) + "/cover.jpg"}
		args = append(args, downloaded...)
		titag := exec.Command("ti-tag", args...)

		fmt.Println("Creating album metadata")
		stdout, _ := titag.StdoutPipe()
		err = titag.Start()

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		titag.Wait()
		if err != nil {
			return errors.Wrap(err, "Couldn't execute \""+titag.String()+"\"")
		}
	}

	return nil
}
