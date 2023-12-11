package vulkan

import (
	"github.com/samber/lo"
	"github.com/vulkan-go/vulkan"
)

type VulkanInstanceCapabilities struct {
	Layers     []vulkan.LayerProperties
	LayerNames []string
}

func NewVulkanInstanceCapabilities() (VulkanInstanceCapabilities, error) {
	capabilities := VulkanInstanceCapabilities{}

	var layerCount uint32
	vulkan.EnumerateInstanceLayerProperties(&layerCount, nil)
	capabilities.Layers = make([]vulkan.LayerProperties, layerCount)
	vulkan.EnumerateInstanceLayerProperties(&layerCount, capabilities.Layers)
	capabilities.LayerNames = lo.Map(capabilities.Layers, func(layer vulkan.LayerProperties, i int) string {
		return string(layer.LayerName[:])
	})

	return capabilities, nil
}
