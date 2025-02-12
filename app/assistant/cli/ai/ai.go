package ai

import (
	"github.com/west2-online/DomTok/app/assistant/model"
	"time"
)

// TODO: complete this file

func Example(input string, dialog model.IDialog) (err error) {
	defer dialog.Close()
	for i := range input {
		if string(input[i]) == "" {
			return nil
		}
		dialog.Send(string(input[i]) + "\n")

		time.Sleep(500 * time.Millisecond)
	}
	return nil
}
