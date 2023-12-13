package window

import (
	"github.com/LamkasDev/seal/cmd/common"
	"github.com/go-gl/glfw/v3.3/glfw"
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

type Window struct {
	Handle *glfw.Window
}

func NewWindow(options WindowOptions) (Window, error) {
	window := Window{}
	glfw.WindowHint(glfw.ClientAPI, glfw.NoAPI)
	windowRaw, err := glfw.CreateWindow(options.Size.Width, options.Size.Height, options.Title, nil, nil)
	if err != nil {
		return window, err
	}
	window.Handle = windowRaw

	return window, nil
}

func RunWindow(window *Window) error {
	return nil
}

func FreeWindow(window *Window) error {
	window.Handle.Destroy()

	return nil
}
