package engine

import (
	"github.com/LamkasDev/seal/cmd/engine/input"
	"github.com/LamkasDev/seal/cmd/engine/renderer"
	"github.com/LamkasDev/seal/cmd/engine/window"
)

type Engine struct {
	Renderer renderer.Renderer
	Input    input.Input
}

func NewEngine(rendererOptions renderer.RendererOptions) (Engine, error) {
	var err error
	engine := Engine{}

	engine.Renderer, err = renderer.NewRenderer(rendererOptions)
	if err != nil {
		return engine, err
	}

	engine.Input, err = input.NewInput()
	if err != nil {
		return engine, err
	}

	return engine, err
}

func RunEngine(engine *Engine) error {
	for !engine.Renderer.Options.Window.Handle.ShouldClose() {
		if err := renderer.RunRenderer(&engine.Renderer); err != nil {
			return err
		}
		if err := window.RunWindow(engine.Renderer.Options.Window); err != nil {
			return err
		}
		if err := input.RunInput(&engine.Input); err != nil {
			return err
		}
	}

	return nil
}

func FreeEngine(engine *Engine) error {
	if err := renderer.FreeRenderer(&engine.Renderer); err != nil {
		return err
	}
	if err := input.FreeInput(&engine.Input); err != nil {
		return err
	}

	return nil
}
