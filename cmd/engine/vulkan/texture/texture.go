package texture

import (
	"image"

	"github.com/LamkasDev/seal/cmd/engine/vulkan/buffer"
	sealImage "github.com/LamkasDev/seal/cmd/engine/vulkan/image"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"

	_ "image/jpeg"
	_ "image/png"
)

type VulkanTexture struct {
	Image  sealImage.VulkanImage
	Buffer buffer.VulkanTextureBuffer
	Device *logical.VulkanLogicalDevice
}

func NewVulkanTexture(device *logical.VulkanLogicalDevice, source *image.RGBA) (VulkanTexture, error) {
	var err error
	texture := VulkanTexture{
		Device: device,
	}

	if texture.Image, err = sealImage.NewVulkanImage(device, uint32(source.Rect.Dx()), uint32(source.Rect.Dy())); err != nil {
		return texture, err
	}
	if texture.Buffer, err = buffer.NewVulkanTextureBuffer(device, source, &texture.Image); err != nil {
		return texture, err
	}

	return texture, nil
}

func FreeVulkanTexture(texture *VulkanTexture) error {
	if err := buffer.FreeVulkanTextureBuffer(&texture.Buffer); err != nil {
		return err
	}

	return sealImage.FreeVulkanImage(&texture.Image)
}
