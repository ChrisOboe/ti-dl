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

package main

import (
	"fmt"
	"github.com/ChrisOboe/ti-dl/download"
	"github.com/ChrisOboe/ti-dl/settings"
	"github.com/ChrisOboe/tidapi"
	"github.com/pkg/errors"
	"os"
	"regexp"
	"strconv"
)

var albumRegex *regexp.Regexp
var artistRegex *regexp.Regexp
var playlistRegex *regexp.Regexp
var userRegex *regexp.Regexp

func main() {
	albumRegex = regexp.MustCompile(`(https?:\\)?(play\.wimpmusic\.com|listen\.tidal\.com)/album/(?P<albumId>\d+)`)
	artistRegex = regexp.MustCompile(`(https?:\\)?(play\.wimpmusic\.com|listen\.tidal\.com)/artist/(?P<artistId>\d+)`)
	playlistRegex = regexp.MustCompile(`(https?:\\)?(play\.wimpmusic\.com|listen\.tidal\.com)/playlist/(?P<playlistId>[0-9a-f]+-[0-9a-f]+-[0-9a-f]+-[0-9a-f]+-[0-9a-f]+)`)
	userRegex = regexp.MustCompile(`(https?:\\)?(play\.wimpmusic\.com|listen\.tidal\.com)/user/(?P<userId>\d+)`)

	a := settings.GetArgs()
	if a.Defaultconfig {
		d, err := settings.Generate()
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		fmt.Println(d)
		os.Exit(0)
	}

	s, err := settings.Get()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if s.Url == "" {
		fmt.Println(settings.Usage())
		os.Exit(0)
	}

	if !(albumRegex.MatchString(s.Url) || artistRegex.MatchString(s.Url) || playlistRegex.MatchString(s.Url) || userRegex.MatchString(s.Url)) {
		fmt.Println("The given url isn't supported.")
		os.Exit(0)
	}

	api := tidapi.New()
	loginData, err := api.LoginUsername(s.Username, s.Password)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	api.SetSessionId(loginData.SessionId)
	api.SetCountryCode(loginData.CountryCode)

	ids, err := albumIds(api, s.Url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	state := make(chan string)
	go statePrinter(state)

	d := download.New(state, tidapi.AudioqualityHiRes, api, s.Destination)
	for _, id := range ids {
		iid, _ := strconv.Atoi(id)
		err := d.Album(iid)
		if err != nil {
			fmt.Println(err)
		}
	}

	close(state)
	os.Exit(0)
}

func albumIds(api tidapi.Tidal, url string) ([]string, error) {
	albumIds := []string{}

	if albumRegex.MatchString(url) {
		albumId := albumRegex.FindStringSubmatch(url)[3]
		albumIds = []string{albumId}
	} else if artistRegex.MatchString(url) {
		artistId := artistRegex.FindStringSubmatch(url)[3]
		iArtistId, _ := strconv.Atoi(artistId)
		// get album ids
		albums, err := api.ArtistAlbums(iArtistId)
		if err != nil {
			return []string{}, errors.Wrap(err, "Couldn't get album releases")
		}
		for _, release := range albums.Items {
			albumIds = append(albumIds, strconv.Itoa(release.ID))
		}
		// get ep ids
		eps, err := api.ArtistEpsAndSingles(iArtistId)
		if err != nil {
			return []string{}, errors.Wrap(err, "Couldn't get album releases")
		}
		for _, release := range eps.Items {
			albumIds = append(albumIds, strconv.Itoa(release.ID))
		}
	} else if playlistRegex.MatchString(url) {
		//playlistId := playlistRegex.FindStringSubmatch(url)[3]
		fmt.Println("Currently unsupported. :(")
		os.Exit(0)
	} else if userRegex.MatchString(url) {
		userId := userRegex.FindStringSubmatch(url)[3]
		iUserId, _ := strconv.Atoi(userId)
		favs, err := api.UserFavorites(iUserId)
		if err != nil {
			return []string{}, errors.Wrap(err, "Couldn't get user favourites")
		}
		albumIds = favs.Album
	}
	return albumIds, nil
}

func statePrinter(state chan string) {
	for {
		fmt.Println(<-state)
	}
}
