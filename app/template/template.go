// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package template

import (
	"fmt"
	"html/template"
	"strings"
	"sync"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/midoks/dztasks/internal/conf"
)

var (
	funcMap     []template.FuncMap
	funcMapOnce sync.Once
)

// FuncMap returns a list of user-defined template functions.
func FuncMap() []template.FuncMap {
	funcMapOnce.Do(func() {
		funcMap = []template.FuncMap{{
			// Application information functions
			"BuildCommit": func() string {
				return conf.App.Version
			},
			"Year": func() int {
				return time.Now().Year()
			},
			"AppSubURL": func() string {
				return conf.Web.Subpath
			},
			"AppName": func() string {
				return conf.App.Name
			},
			"AppVer": func() string {
				return conf.App.Version
			},

			// HTML processing functions
			"Safe":        Safe,
			"Str2HTML":    Str2HTML,
			"Sanitize":    bluemonday.UGCPolicy().Sanitize,
			"NewLine2br":  NewLine2br,
			"EscapePound": EscapePound,
			"Add": func(a, b int) int {
				return a + b
			},

			"SubStr": func(str string, start, length int) string {
				if len(str) == 0 {
					return ""
				}
				end := start + length
				if length == -1 {
					end = len(str)
				}
				if len(str) < end {
					return str
				}
				return str[start:end]
			},
			"LoadTimes": func(startTime time.Time) string {
				return fmt.Sprint(time.Since(startTime).Nanoseconds()/1e6) + "ms"
			},
			"Join": strings.Join,
			"DateFmtLong": func(t time.Time) string {
				return t.Format(time.RFC1123Z)
			},
			"DateFmtShort": func(t time.Time) string {
				// fmt.Println(t)
				return t.Format("Jan 02, 2006")
			},
		}}
	})
	return funcMap
}

func Safe(raw string) template.HTML {
	return template.HTML(raw)
}

func Str2HTML(raw string) template.HTML {
	return template.HTML(bluemonday.UGCPolicy().Sanitize(raw))
}

// NewLine2br simply replaces "\n" to "<br>".
func NewLine2br(raw string) string {
	return strings.Replace(raw, "\n", "<br>", -1)
}

// TODO: Use url.Escape.
func EscapePound(str string) string {
	return strings.NewReplacer("%", "%25", "#", "%23", " ", "%20", "?", "%3F").Replace(str)
}
