package window

import (
	"unsafe"

	"github.com/LamkasDev/seal/cmd/common/constants"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/vulkan-go/vulkan"
)

type Window struct {
	Handle  *glfw.Window
	Data    *WindowData
	Options WindowOptions
}

type WindowData struct {
	Resized bool
	Extent  vulkan.Extent2D
}

func NewWindow() (Window, error) {
	data := WindowData{
		Resized: false,
		Extent:  constants.DefaultExtent,
	}
	window := Window{
		Data:    &data,
		Options: NewWindowOptions("Test", constants.DefaultExtent),
	}
	glfw.WindowHint(glfw.ClientAPI, glfw.NoAPI)

	windowRaw, err := glfw.CreateWindow(int(window.Options.Extent.Width), int(window.Options.Extent.Height), window.Options.Title, nil, nil)
	if err != nil {
		return window, err
	}
	window.Handle = windowRaw
	window.Handle.SetUserPointer(unsafe.Pointer(&data))
	window.Handle.SetFramebufferSizeCallback(ResizeWindowCallback)

	return window, nil
}

func ResizeWindowCallback(w *glfw.Window, width, height int) {
	data := (*WindowData)(w.GetUserPointer())
	data.Resized = true
	data.Extent = vulkan.Extent2D{
		Width:  uint32(width),
		Height: uint32(height),
	}
}

func FreeWindow(window *Window) error {
	window.Handle.Destroy()
	return nil
}
