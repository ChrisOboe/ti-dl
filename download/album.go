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
	"fmt"
	"github.com/pkg/errors"
	"os"
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
	for _, item := range items.Items {
		d.state <- fmt.Sprintf("Downloading track: %s (%d)", item.Item.Title, item.Item.ID)
		_albumartist := stripBadChars(album.Artist.Name)
		_album := stripBadChars(album.Title)
		_tracknumber := strconv.Itoa(item.Item.TrackNumber)
		_title := stripBadChars(item.Item.Title)
		_releaseDate := stripBadChars(album.ReleaseDate)
		_releaseType := stripBadChars(album.Type)

		filePath := strings.Replace(d.destination, "${ALBUMARTIST}", _albumartist, -1)
		filePath = strings.Replace(filePath, "${RELEASETITLE}", _album, -1)
		filePath = strings.Replace(filePath, "${TRACKNUMBER}", _tracknumber, -1)
		filePath = strings.Replace(filePath, "${TRACKTITLE}", _title, -1)
		filePath = strings.Replace(filePath, "${RELEASEDATE}", _releaseDate, -1)
		filePath = strings.Replace(filePath, "${RELEASETYPE}", _releaseType, -1)

		err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
		if err != nil {
			return errors.Wrap(err, "error creating "+filepath.Dir(filePath))
		}
		err = d.trackDownload(item.Item.ID, filePath)
		if err != nil {
			return errors.Wrap(err, "Problem with downloading + "+strconv.Itoa(albumId))
		}
	}

	return nil
}
