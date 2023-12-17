package engine

import (
	"time"

	"github.com/LamkasDev/seal/cmd/engine/input"
	"github.com/LamkasDev/seal/cmd/engine/progress"
	"github.com/LamkasDev/seal/cmd/engine/renderer"
	"github.com/LamkasDev/seal/cmd/engine/window"
	"github.com/go-gl/glfw/v3.3/glfw"
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
	ups, fps := 30, 60
	last := time.Now().UnixMilli()
	lastUpdate, lastFrame := last, last
	deltaUpdate, deltaFrame := float64(0), float64(0)
	targetUpdate, targetFrame := float64(1000/ups), float64(1000/fps)

	for !engine.Renderer.Window.Handle.ShouldClose() {
		glfw.PollEvents()

		now := time.Now().UnixMilli()
		diff := float64(now - last)
		deltaUpdate += diff / targetUpdate
		deltaFrame += diff / targetFrame
		if deltaUpdate > 1 {
			if err := input.RunInput(&engine.Input); err != nil {
				return err
			}
			_ = now - lastUpdate
			lastUpdate = now
			deltaUpdate--
		}
		if deltaFrame > 1 {
			if err := renderer.RunRenderer(&engine.Renderer); err != nil {
				return err
			}
			if err := window.RunWindow(&engine.Renderer.Window); err != nil {
				return err
			}
			_ = now - lastFrame
			lastFrame = now
			deltaFrame--
		}
		last = now
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
