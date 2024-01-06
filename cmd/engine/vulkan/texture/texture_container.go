package texture

import (
	"fmt"

	"github.com/LamkasDev/seal/cmd/engine/vulkan/buffer"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/logger"
	"neilpa.me/go-stbi"
)

const TEXTURE_BASIC = "basic"

var DefaultTextures = []string{
	TEXTURE_BASIC,
}

type VulkanTextureContainer struct {
	Device   *logical.VulkanLogicalDevice
	Textures map[string]VulkanTexture
}

func NewVulkanTextureContainer(device *logical.VulkanLogicalDevice) (VulkanTextureContainer, error) {
	container := VulkanTextureContainer{
		Device:   device,
		Textures: map[string]VulkanTexture{},
	}
	for _, texture := range DefaultTextures {
		if _, err := CreateVulkanTextureWithContainer(&container, texture); err != nil {
			return container, err
		}
	}
	logger.DefaultLogger.Debug("created new vulkan texture container")

	return container, nil
}

func CreateVulkanTextureWithContainer(container *VulkanTextureContainer, id string) (VulkanTexture, error) {
	textureImage, err := stbi.Load(fmt.Sprintf("../../resources/textures/%s.jpg", id))
	if err != nil {
		return VulkanTexture{}, err
	}

	texture, err := NewVulkanTexture(container.Device, buffer.NewVulkanTextureBufferOptions(textureImage.Pix))
	if err != nil {
		return texture, err
	}
	container.Textures[id] = texture

	return texture, nil
}

func FreeVulkanTextureContainer(container *VulkanTextureContainer) error {
	for _, texture := range container.Textures {
		if err := FreeVulkanTexture(&texture); err != nil {
			return err
		}
	}

	return nil
}
