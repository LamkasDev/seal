package vulkan

import (
	"github.com/LamkasDev/seal/cmd/common/ctool"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type VulkanInstanceCapabilities struct {
	Layers     []vulkan.LayerProperties
	LayerNames []string
}

func NewVulkanInstanceCapabilities() (VulkanInstanceCapabilities, error) {
	capabilities := VulkanInstanceCapabilities{}

	var layerCount uint32
	if res := vulkan.EnumerateInstanceLayerProperties(&layerCount, nil); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
	}
	capabilities.Layers = make([]vulkan.LayerProperties, layerCount)
	capabilities.LayerNames = make([]string, layerCount)
	if res := vulkan.EnumerateInstanceLayerProperties(&layerCount, capabilities.Layers); res != vulkan.Success {
		logger.DefaultLogger.Error(vulkan.Error(res))
	}
	for i := 0; i < len(capabilities.Layers); i++ {
		capabilities.Layers[i].Deref()
		capabilities.LayerNames[i] = ctool.ByteArray256ToString(capabilities.Layers[i].LayerName)
	}

	return capabilities, nil
}
