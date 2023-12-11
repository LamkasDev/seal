package vulkan

import (
	"errors"
	"fmt"
	"slices"

	"github.com/LamkasDev/seal/cmd/common/arch"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/go-gl/glfw/v3.3/glfw"
	vulkan "github.com/vkngwrapper/core/v2/core1_0"
)

type VulkanInstanceOptions struct {
	CreateInfo vulkan.InstanceCreateInfo
}

func NewVulkanInstanceOptions(capabilities *VulkanInstanceCapabilities) (VulkanInstanceOptions, error) {
	options := VulkanInstanceOptions{}
	options.CreateInfo = vulkan.InstanceCreateInfo{
		ApplicationName:       "Seal Game",
		EngineName:            "Seal Engine",
		EnabledExtensionNames: glfw.GetCurrentContext().GetRequiredInstanceExtensions(),
		EnabledLayerNames:     []string{},
	}
	if arch.SealDebug {
		if err := EnableVulkanInstanceOptionsLayer(&options, capabilities, "VK_LAYER_KHRONOS_validation"); err != nil {
			logger.DefaultLogger.Warn(err.Error())
		}
		options.CreateInfo.EnabledExtensionNames = append(options.CreateInfo.EnabledExtensionNames, "VK_EXT_DEBUG_UTILS_EXTENSION_NAME")
	}

	return options, nil
}

func EnableVulkanInstanceOptionsLayer(options *VulkanInstanceOptions, capabilities *VulkanInstanceCapabilities, layer string) error {
	if !slices.Contains(capabilities.AvailableLayers, layer) {
		return errors.New(fmt.Sprintf("tried to enable layer '%s', but it is not available", layer))
	}
	options.CreateInfo.EnabledLayerNames = append(options.CreateInfo.EnabledLayerNames, layer)

	return nil
}
