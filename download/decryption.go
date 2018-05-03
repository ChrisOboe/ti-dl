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
	"crypto/cipher"
	"encoding/base64"
	"github.com/ChrisOboe/rijndael256"
	"github.com/pkg/errors"
)

const masterKeyBase64 string = "UIlTTEMmmLfGowo/UC60x2H45W6MdGgTRfo/umg4754="

func fileKey(keyMaterialBase64 string) ([]byte, []byte, error) {
	masterKey, _ := base64.StdEncoding.DecodeString(masterKeyBase64)
	keyMaterial, _ := base64.StdEncoding.DecodeString(keyMaterialBase64)

	block, err := rijndael256.NewCipher(masterKey)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Couldn't create rijndael256 cipher for decrypting the fileKey with masterkey")
	}
	mode := cipher.NewCBCDecrypter(block, make([]byte, 16))
	mode.CryptBlocks(keyMaterial, keyMaterial)

	// [0:16] of keymaterial is junk
	key := keyMaterial[16:32]
	iv := append(keyMaterial[32:40], make([]byte, 8)...)

	return key, iv, nil
}
