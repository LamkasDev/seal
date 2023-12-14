package vulkan

import (
	"slices"

	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/vulkan-go/vulkan"
)

type VulkanPhysicalDevice struct {
	Handle       vulkan.PhysicalDevice
	Capabilities VulkanPhysicalDeviceCapabilities
	Properties   vulkan.PhysicalDeviceProperties
	Features     vulkan.PhysicalDeviceFeatures
}

func NewVulkanPhysicalDevice(handle vulkan.PhysicalDevice, window *glfw.Window, surface *vulkan.Surface) (VulkanPhysicalDevice, error) {
	var err error
	device := VulkanPhysicalDevice{
		Handle: handle,
	}

	if device.Capabilities, err = NewVulkanPhysicalDeviceCapabilities(handle, window, surface); err != nil {
		logger.DefaultLogger.Panic(err.Error())
	}
	logger.DefaultLogger.Debug("created new vulkan instance capabilities")

	vulkan.GetPhysicalDeviceProperties(device.Handle, &device.Properties)
	device.Properties.Deref()
	vulkan.GetPhysicalDeviceFeatures(device.Handle, &device.Features)
	device.Features.Deref()

	return device, nil
}

func CompareVulkanPhysicalDevice(device *VulkanPhysicalDevice) int {
	return int(device.Properties.Limits.MaxImageDimension2D)
}

func IsVulkanPhysicalDeviceSupported(device *VulkanPhysicalDevice) bool {
	return (device.Features.GeometryShader == 1) &&
		device.Capabilities.Queue.GraphicsIndex != -1 &&
		device.Capabilities.Queue.PresentationIndex != -1 &&
		device.Capabilities.Surface.ImageFormatIndex != -1 &&
		len(device.Capabilities.Surface.PresentModes) > 0 &&
		slices.Contains(device.Capabilities.ExtensionNames, vulkan.KhrSwapchainExtensionName)
}
