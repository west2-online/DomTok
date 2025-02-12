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

package ai

import (
	"time"

	"github.com/west2-online/DomTok/app/assistant/model"
)

// TODO: complete this file

func Example(input string, dialog model.IDialog) (err error) {
	defer dialog.Close()
	for i := range input {
		if string(input[i]) == "" {
			return nil
		}
		dialog.Send(string(input[i]) + "\n")

		time.Sleep(time.Second)
	}
	return nil
}
