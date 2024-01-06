package input

import (
	"github.com/LamkasDev/seal/cmd/engine/renderer"
	"github.com/LamkasDev/seal/cmd/engine/scene"
	"github.com/LamkasDev/seal/cmd/engine/time"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Input struct {
	Scene *scene.Scene
}

func NewInput(scene *scene.Scene) (Input, error) {
	input := Input{
		Scene: scene,
	}

	return input, nil
}

func RunInput(input *Input) error {
	mov := time.DeltaTime * 10
	if renderer.RendererInstance.Window.Handle.GetKey(glfw.KeyD) == glfw.Press {
		renderer.RendererInstance.Pipeline.Camera.Position[0] += mov
	}
	if renderer.RendererInstance.Window.Handle.GetKey(glfw.KeyA) == glfw.Press {
		renderer.RendererInstance.Pipeline.Camera.Position[0] -= mov
	}
	if renderer.RendererInstance.Window.Handle.GetKey(glfw.KeySpace) == glfw.Press {
		renderer.RendererInstance.Pipeline.Camera.Position[1] += mov
	}
	if renderer.RendererInstance.Window.Handle.GetKey(glfw.KeyLeftControl) == glfw.Press {
		renderer.RendererInstance.Pipeline.Camera.Position[1] -= mov
	}
	if renderer.RendererInstance.Window.Handle.GetKey(glfw.KeyW) == glfw.Press {
		renderer.RendererInstance.Pipeline.Camera.Position[2] -= mov
	}
	if renderer.RendererInstance.Window.Handle.GetKey(glfw.KeyS) == glfw.Press {
		renderer.RendererInstance.Pipeline.Camera.Position[2] += mov
	}

	return nil
}

func FreeInput(input *Input) error {
	return nil
}
