/*
Copyright 2024 The west2-online Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import (
	"net/url"
	"strings"
)

func UriEncode(uri string) string {
	uris := strings.Split(uri, "/")
	for i := 0; i < len(uris); i++ {
		uris[i] = url.PathEscape(uris[i])
	}
	return strings.Join(uris, "/")
}

// HrefEncode 解析rawQuery以获得未解码的url(hertz会自动解析未解码的url)
func HrefEncode(rawQuery string) string {
	params := strings.Split(rawQuery, "&")
	for _, param := range params {
		if strings.HasPrefix(param, "href") {
			return strings.Replace(param, "href=", "", 1)
		}
	}
	return ""
}
