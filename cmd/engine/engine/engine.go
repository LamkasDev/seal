package engine

import (
	"github.com/LamkasDev/seal/cmd/engine/input"
	"github.com/LamkasDev/seal/cmd/engine/progress"
	"github.com/LamkasDev/seal/cmd/engine/renderer"
	"github.com/LamkasDev/seal/cmd/engine/window"
)

type Engine struct {
	Renderer renderer.Renderer
	Input    input.Input
}

func NewEngine() (Engine, error) {
	var err error
	engine := Engine{}

	progress.AdvanceLoading()
	if engine.Renderer, err = renderer.NewRenderer(); err != nil {
		return engine, err
	}

	progress.AdvanceLoading()
	if engine.Input, err = input.NewInput(); err != nil {
		return engine, err
	}

	return engine, err
}

func RunEngine(engine *Engine) error {
	for !engine.Renderer.Window.Handle.ShouldClose() {
		if err := renderer.RunRenderer(&engine.Renderer); err != nil {
			return err
		}
		if err := window.RunWindow(&engine.Renderer.Window); err != nil {
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
