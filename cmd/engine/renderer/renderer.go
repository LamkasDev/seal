package renderer

type Renderer struct {
}

func NewRenderer() (Renderer, error) {
	renderer := Renderer{}

	return renderer, nil
}

func RunRenderer(renderer *Renderer) error {
	return nil
}

func FreeRenderer(renderer *Renderer) error {
	return nil
}
