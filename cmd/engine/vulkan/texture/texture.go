package texture

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/buffer"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"

	_ "image/jpeg"
	_ "image/png"
)

type VulkanTexture struct {
	Buffer buffer.VulkanTextureBuffer
	Device *logical.VulkanLogicalDevice
}

func NewVulkanTexture(device *logical.VulkanLogicalDevice, options buffer.VulkanTextureBufferOptions) (VulkanTexture, error) {
	var err error
	texture := VulkanTexture{
		Device: device,
	}

	if texture.Buffer, err = buffer.NewVulkanTextureBuffer(device, options); err != nil {
		return texture, err
	}

	return texture, nil
}

func FreeVulkanTexture(texture *VulkanTexture) error {
	return nil
}
