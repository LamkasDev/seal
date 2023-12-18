package window

import (
	"github.com/LamkasDev/seal/cmd/common/constants"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/vulkan-go/vulkan"
)

type Window struct {
	Handle  *glfw.Window
	Options WindowOptions
}

func NewWindow() (Window, error) {
	window := Window{
		Options: NewWindowOptions("Test", constants.DefaultResolution),
	}
	glfw.WindowHint(glfw.ClientAPI, glfw.NoAPI)
	windowRaw, err := glfw.CreateWindow(window.Options.Size.Width, window.Options.Size.Height, window.Options.Title, nil, nil)
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

func GetWindowImageExtent(window *Window) vulkan.Extent2D {
	w, h := window.Handle.GetFramebufferSize()
	return vulkan.Extent2D{
		Width:  uint32(w),
		Height: uint32(h),
	}
}
