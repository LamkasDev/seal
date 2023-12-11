package common

type Size struct {
	Width  int
	Height int
}

func NewSize(width int, height int) Size {
	return Size{
		Width:  width,
		Height: height,
	}
}
