package font

import (
	"image"
	"image/draw"

	"os"

	"github.com/LamkasDev/seal/cmd/engine/vulkan/texture"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

type VulkanFont struct {
	Font    truetype.Font
	Texture texture.VulkanTexture
}

func NewVulkanFont(file string) (VulkanFont, error) {
	font := VulkanFont{}

	// load font file and typeface
	fontBytes, err := os.ReadFile("../../resources/fonts/default.ttf")
	if err != nil {
		return font, err
	}
	f, err := truetype.Parse(fontBytes)
	opts := truetype.Options{
		Size: 12,
	}
	face := truetype.NewFace(f, &opts) // sampling of some of the options that are set
	fg, bg := image.White, image.Black
	rgba := image.NewRGBA(image.Rect(0, 0, IMAGE_WIDTH, IMAGE_HEIGHT))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetFont(f)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)

	// Single-line case// calculate full length of string
	textWidth := font.MeasureString(face, text).Ceil()
	textHeight := face.Metrics().Ascent.Ceil() + face.Metrics().Descent.Ceil() // to center text in your image
	x := (IMAGE_WIDTH - textWidth) / 2
	y := (IMAGE_HEIGHT - textHeight) / 2
	pt := freetype.Pt(x, y) // draw the string
	_, err = c.DrawString(text, pt)
}
