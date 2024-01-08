package font

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/logger"
)

const FONT_DEFAULT = "default"

var DefaultFonts = []string{
	FONT_DEFAULT,
}

type VulkanFontContainer struct {
	Device *logical.VulkanLogicalDevice
	Fonts  map[string]*VulkanFont
}

func NewVulkanFontContainer(device *logical.VulkanLogicalDevice) (VulkanFontContainer, error) {
	container := VulkanFontContainer{
		Device: device,
		Fonts:  map[string]*VulkanFont{},
	}
	for _, id := range DefaultFonts {
		if _, err := CreateVulkanFontWithContainer(&container, id); err != nil {
			return container, err
		}
	}
	logger.DefaultLogger.Debug("created new vulkan font container")

	return container, nil
}

func CreateVulkanFontWithContainer(container *VulkanFontContainer, id string) (VulkanFont, error) {
	Font, err := NewVulkanFont(container.Device, id)
	if err != nil {
		return Font, err
	}
	container.Fonts[id] = &Font

	return Font, nil
}

func FreeVulkanFontContainer(container *VulkanFontContainer) error {
	for _, Font := range container.Fonts {
		if err := FreeVulkanFont(Font); err != nil {
			return err
		}
	}

	return nil
}
