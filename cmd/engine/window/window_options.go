package window

import (
	"github.com/vulkan-go/vulkan"
)

type WindowOptions struct {
	Title  string
	Extent vulkan.Extent2D
}

func NewWindowOptions(title string, extent vulkan.Extent2D) WindowOptions {
	return WindowOptions{
		Title:  title,
		Extent: extent,
	}
}
