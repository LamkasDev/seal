package vulkan

import (
	"github.com/vkngwrapper/core/v2"
	"golang.org/x/exp/maps"
)

type VulkanInstanceCapabilities struct {
	AvailableLayers []string
}

func NewVulkanInstanceCapabilities(loader *core.VulkanLoader) (VulkanInstanceCapabilities, error) {
	capabilities := VulkanInstanceCapabilities{}
	layers, _, err := loader.AvailableLayers()
	if err != nil {
		return capabilities, err
	}
	capabilities.AvailableLayers = maps.Keys(layers)

	return capabilities, nil
}
