package glfw

import "github.com/go-gl/glfw/v3.3/glfw"

func StartGLFW() error {
	if err := glfw.Init(); err != nil {
		return err
	}

	return nil
}

func EndGLFW() error {
	glfw.Terminate()

	return nil
}
