package vulkan

import (
	"errors"
	"fmt"
	"slices"

	"github.com/LamkasDev/seal/cmd/common/arch"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/vulkan-go/vulkan"
)

type VulkanInstanceOptions struct {
	ApplicationInfo vulkan.ApplicationInfo
	CreateInfo      vulkan.InstanceCreateInfo
}

func NewVulkanInstanceOptions(instance *VulkanInstance) (VulkanInstanceOptions, error) {
	options := VulkanInstanceOptions{}
	options.ApplicationInfo = vulkan.ApplicationInfo{
		SType:            vulkan.StructureTypeApplicationInfo,
		PApplicationName: "Seal Game",
		PEngineName:      "Seal Engine",
		ApiVersion:       vulkan.ApiVersion11,
	}
	options.CreateInfo = vulkan.InstanceCreateInfo{
		SType:                   vulkan.StructureTypeInstanceCreateInfo,
		PApplicationInfo:        &options.ApplicationInfo,
		PpEnabledExtensionNames: glfw.GetCurrentContext().GetRequiredInstanceExtensions(),
		PpEnabledLayerNames:     []string{},
	}
	if arch.SealDebug {
		if err := EnableVulkanInstanceOptionsLayer(&options, &instance.Capabilities, "VK_LAYER_KHRONOS_validation"); err != nil {
			logger.DefaultLogger.Warn(err.Error())
		}
	}
	options.CreateInfo.EnabledExtensionCount = uint32(len(options.CreateInfo.PpEnabledExtensionNames))
	options.CreateInfo.EnabledLayerCount = uint32(len(options.CreateInfo.PpEnabledLayerNames))

	return options, nil
}

func EnableVulkanInstanceOptionsLayer(options *VulkanInstanceOptions, capabilities *VulkanInstanceCapabilities, layer string) error {
	if !slices.Contains(capabilities.LayerNames, layer) {
		return errors.New(fmt.Sprintf("tried to enable layer '%s', but it is not available", layer))
	}
	options.CreateInfo.PpEnabledLayerNames = append(options.CreateInfo.PpEnabledLayerNames, layer)

	return nil
}
