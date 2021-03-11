// Copyright 2021 PurpleSec Team
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package parseurl is a simple wrapper that fixes some of the weird issues that the
// standard Golang 'url.Parse' function does.
//
// Fixes things such as
// - "localhost:8080"
//   url.Parse: Host == "", Scheme == "localhost:8080"
// - "10.10.10.10/url/"
//   url.Parse: Returns an error <why?>
// - "localhost:"
//   url.Parse: Host == "localhost:" <doesn't strip the ':'>
//
// This package adds in checks for invalid values returned by 'url.Parse' such as the Host field being empty.
// All non-standard errors wrap the error "ErrInvalidURL" to assist in indication of the error.
//
// This library is a drop-in replacement for the "url.Parse" function.
// Just import "github.com/PurpleSec/parseurl" and go!
//
package parseurl

import (
	"errors"
	"net/url"
	"strings"
)

// ErrInvalidURL is a error that is used to indicate a non-standard error caught by the 'Parse' function. The
// returned errors will wrap this error.
var ErrInvalidURL = errors.New("URL is invalid")

type errStr string

func (errStr) Unwrap() error {
	return ErrInvalidURL
}
func (e errStr) Error() string {
	return string(e)
}

// Parse parses rawurl into a URL structure.
//
// The rawurl may be relative (a path, without a host) or absolute
// (starting with a scheme). Trying to parse a hostname and path
// without a scheme is invalid but may not necessarily return an
// error, due to parsing ambiguities.
//
// This function is a modified version of the standard 'url.Parse' function that will handle and fix any errors
// that occur during Parsing. This function also includes additional error checks that will prevent some common
// formatting issues from occuring without an error.
func Parse(rawurl string) (*url.URL, error) {
	var (
		i   = strings.IndexRune(rawurl, '/')
		u   *url.URL
		err error
	)
	if i == 0 && len(rawurl) > 2 && rawurl[1] != '/' {
		u, err = url.Parse("/" + rawurl)
	} else if i == -1 || i+1 >= len(rawurl) || rawurl[i+1] != '/' {
		u, err = url.Parse("//" + rawurl)
	} else {
		u, err = url.Parse(rawurl)
	}
	if err != nil {
		return nil, err
	}
	if len(u.Host) == 0 {
		return nil, errStr(`parse "` + rawurl + `": empty host field`)
	}
	if u.Host[len(u.Host)-1] == ':' {
		return nil, errStr(`parse "` + rawurl + `": invalid port specified`)
	}
	return u, nil
}
