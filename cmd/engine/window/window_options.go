package window

import (
	"github.com/LamkasDev/seal/cmd/common"
)

type WindowOptions struct {
	Title string
	Size  common.Size
}

func NewWindowOptions(title string, size common.Size) WindowOptions {
	return WindowOptions{
		Title: title,
		Size:  size,
	}
}
