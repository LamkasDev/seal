package font

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"

	"os"

	"github.com/EngoEngine/glm"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/buffer"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/mesh"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/shader"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/texture"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/vertex"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

type VulkanFont struct {
	Options truetype.Options
	Device  *logical.VulkanLogicalDevice
	Font    *truetype.Font
	Texture texture.VulkanTexture
}

func NewVulkanFont(device *logical.VulkanLogicalDevice, id string) (VulkanFont, error) {
	var err error
	cfont := VulkanFont{
		Options: truetype.Options{
			Size: 128,
			DPI:  72,
		},
		Device: device,
	}

	fontBytes, err := os.ReadFile(fmt.Sprintf("../../resources/fonts/%s.ttf", id))
	if err != nil {
		return cfont, err
	}
	if cfont.Font, err = truetype.Parse(fontBytes); err != nil {
		return cfont, err
	}

	text := "Cube"
	face := truetype.NewFace(cfont.Font, &cfont.Options)
	textWidth := font.MeasureString(face, text).Ceil()
	textHeight := face.Metrics().Ascent.Ceil() + face.Metrics().Descent.Ceil()

	textColor, background := image.White, image.Transparent
	rgba := image.NewRGBA(image.Rect(0, 0, textWidth, textHeight))
	draw.Draw(rgba, rgba.Bounds(), background, image.Point{}, draw.Src)

	context := freetype.NewContext()
	context.SetFont(cfont.Font)
	context.SetClip(rgba.Bounds())
	context.SetDst(rgba)
	context.SetSrc(textColor)
	context.SetFontSize(cfont.Options.Size)
	context.SetDPI(cfont.Options.DPI)

	_, err = context.DrawString(text, freetype.Pt(0, int(cfont.Options.Size)))
	if err != nil {
		return cfont, err
	}
	if cfont.Texture, err = texture.NewVulkanTexture(device, rgba); err != nil {
		return cfont, err
	}
	f, _ := os.Create("../../resources/fonts/test.png")
	png.Encode(f, rgba)

	return cfont, nil
}

func CreateVulkanFontMesh(container *mesh.VulkanMeshContainer, font *VulkanFont) (mesh.VulkanMesh, error) {
	height := float32(font.Texture.Image.Options.CreateInfo.Extent.Width) / float32(font.Texture.Image.Options.CreateInfo.Extent.Height)
	return mesh.CreateVulkanMeshWithContainer(container, "font", shader.SHADER_BASIC, &font.Texture, buffer.VulkanMeshBufferOptions{
		Vertices: []vertex.VulkanVertex{
			{Position: glm.Vec3{0, 0, 0}, TexCoord: glm.Vec2{1, 1}},
			{Position: glm.Vec3{1, 0, 0}, TexCoord: glm.Vec2{1, 0}},
			{Position: glm.Vec3{1, height, 0}, TexCoord: glm.Vec2{0, 0}},
			{Position: glm.Vec3{0, height, 0}, TexCoord: glm.Vec2{0, 1}},
		},
		Indices: []uint16{
			0, 1, 2, 2, 3, 0,
		},
	})
}

func FreeVulkanFont(font *VulkanFont) error {
	return nil
}
