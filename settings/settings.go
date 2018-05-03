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

package settings

import (
	"bytes"
	"github.com/BurntSushi/toml"
	"github.com/ChrisOboe/dirs"
	"github.com/alexflint/go-arg"
	"github.com/pkg/errors"
	"os"
)

const appname = "ti-dl"
const defaultDest = "./${ALBUMARTIST}/${RELEASETYPE}/${RELEASEDATE} - ${RELEASETITLE}/${TRACKNUMBER} - ${TRACKTITLE}"

type Settings struct {
	Configfile  string
	Url         string
	Username    string
	Password    string
	Destination string
}

type args struct {
	Configfile    string `arg:"-c" help:"The path of the configfile to use."`
	Url           string `arg:"positional" help:"A url to an album, an artist, a playlist or a user."`
	Username      string `arg:"-u" help:"The username of your tidal account."`
	Password      string `arg:"-p" help:"The password of your tidal account."`
	Destination   string `arg:"-d" help:"The path were music gets downloaded."`
	Defaultconfig bool   `help:"Prints an example config file."`
}

type configfile struct {
	User struct {
		Username string
		Password string
	}
	Paths struct {
		Destination string
	}
}

func Usage() string {
	var a args
	p := arg.MustParse(&a)
	buf := new(bytes.Buffer)
	p.WriteHelp(buf)
	return buf.String()
}

func Generate() (string, error) {
	c := configfile{}
	c.User.Username = "dummyuser"
	c.User.Password = "dummypassword"
	c.Paths.Destination = defaultDest

	buf := new(bytes.Buffer)
	err := toml.NewEncoder(buf).Encode(c)
	if err != nil {
		return "", errors.Wrap(err, "couldn't encode config to toml")
	}
	return buf.String(), nil
}

func Get() (Settings, error) {
	a := GetArgs()
	var c configfile

	_, err := os.Stat(a.Configfile)
	if err == nil {
		c, err = getConf(a.Configfile)
		if err != nil {
			return Settings{}, errors.Wrap(err, "problem with configfile")
		}
	}

	var s Settings
	s.Configfile = a.Configfile
	s.Url = a.Url

	if a.Username != "" {
		s.Username = a.Username
	} else if c.User.Username != "" {
		s.Username = c.User.Username
	} else {
		return Settings{}, errors.New("no username given")
	}

	if a.Password != "" {
		s.Password = a.Password
	} else if c.User.Password != "" {
		s.Password = c.User.Password
	} else {
		return Settings{}, errors.New("no password given")
	}

	if a.Destination != "" {
		s.Destination = a.Destination
	} else if c.Paths.Destination != "" {
		s.Destination = c.Paths.Destination
	} else {
		s.Destination = defaultDest
	}

	return s, nil
}

func GetArgs() args {
	var a args
	a.Configfile = dirs.Get(appname).Config + "config.toml"
	arg.MustParse(&a)
	return a
}

func getConf(file string) (configfile, error) {
	var c configfile
	_, err := toml.DecodeFile(file, &c)
	if err != nil {
		return configfile{}, errors.Wrap(err, "can't read configfile")
	}
	return c, nil
}
