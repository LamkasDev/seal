package input

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Input struct {
}

func NewInput() (Input, error) {
	input := Input{}

	return input, nil
}

func RunInput(input *Input) error {
	glfw.PollEvents()
	return nil
}

func FreeInput(input *Input) error {
	return nil
}
