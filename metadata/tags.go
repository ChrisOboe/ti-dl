// Copyright (c) 2018-2019 ChrisOboe <chris@oboe.email>
// SPDX-License-Identifier: GPL-3.0

/*
package metadata

import (
	"github.com/pkg/errors"
	"github.com/wtolson/go-taglib"
	"strings"
)

type Tags struct {
	Title       string
	Album       string
	Artist      string
	Albumartist string
	Track       int
	Totaltracks int
	Disk        int
	Totaldisks  int
	Year        int
}

func WriteTags(filepath string, tags Tags) error {
	file, err := taglib.Read(filepath)
	if err != nil {
		return errors.Wrap(err, "Couldn't open file for tagging")
	}
	defer file.Close()

	file.SetTitle(tags.Title)
	file.SetAlbum(tags.Album)
	file.SetArtist(tags.Artist)
	
	file.SetYear(tags.Year)

	if strings.HasSuffix(filepath, "mp4") || strings.HasSuffix(filepath, "m4a") {

	} else {
    	file.SetTag(
file.SetTrack(tags.Track)
	}

	return nil
}

*/
