package physical

import (
	"slices"

	"github.com/LamkasDev/seal/cmd/engine/window"
	"github.com/LamkasDev/seal/cmd/logger"
	"github.com/vulkan-go/vulkan"
)

type VulkanPhysicalDevice struct {
	Handle       vulkan.PhysicalDevice
	Window       *window.Window
	Surface      *vulkan.Surface
	Capabilities VulkanPhysicalDeviceCapabilities
	Properties   vulkan.PhysicalDeviceProperties
	Features     vulkan.PhysicalDeviceFeatures
}

func NewVulkanPhysicalDevice(handle vulkan.PhysicalDevice, cwindow *window.Window, surface *vulkan.Surface) (VulkanPhysicalDevice, error) {
	var err error
	device := VulkanPhysicalDevice{
		Handle:  handle,
		Window:  cwindow,
		Surface: surface,
	}

	if device.Capabilities, err = NewVulkanPhysicalDeviceCapabilities(handle, cwindow, surface); err != nil {
		return device, err
	}
	logger.DefaultLogger.Debug("created new vulkan physical device capabilities")

	vulkan.GetPhysicalDeviceProperties(device.Handle, &device.Properties)
	device.Properties.Deref()
	vulkan.GetPhysicalDeviceFeatures(device.Handle, &device.Features)
	device.Features.Deref()

	return device, nil
}

func ResizeVulkanPhysicalDevice(device *VulkanPhysicalDevice) error {
	return ResizeVulkanPhysicalDeviceCapabilities(&device.Capabilities, device.Handle, device.Window, device.Surface)
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
