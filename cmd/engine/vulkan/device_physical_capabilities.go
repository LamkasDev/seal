package vulkan

import (
	"github.com/LamkasDev/seal/cmd/common/ctool"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/samber/lo"
	"github.com/vulkan-go/vulkan"
)

type VulkanPhysicalDeviceCapabilities struct {
	Extensions     []vulkan.ExtensionProperties
	ExtensionNames []string
	Queue          VulkanPhysicalDeviceQueueCapabilities
	Surface        VulkanPhysicalDeviceSurfaceCapabilities
}

type VulkanPhysicalDeviceQueueCapabilities struct {
	Families          []vulkan.QueueFamilyProperties
	GraphicsIndex     int
	PresentationIndex int
}

type VulkanPhysicalDeviceSurfaceCapabilities struct {
	Capabilities     vulkan.SurfaceCapabilities
	Formats          []vulkan.SurfaceFormat
	FormatIndex      int
	PresentModes     []vulkan.PresentMode
	PresentModeIndex int
	Extent           vulkan.Extent2D
}

func NewVulkanPhysicalDeviceCapabilities(handle vulkan.PhysicalDevice, window *glfw.Window, surface *vulkan.Surface) (VulkanPhysicalDeviceCapabilities, error) {
	capabilities := VulkanPhysicalDeviceCapabilities{}

	var extensionsCount uint32
	vulkan.EnumerateDeviceExtensionProperties(handle, "", &extensionsCount, nil)
	capabilities.Extensions = make([]vulkan.ExtensionProperties, extensionsCount)
	capabilities.ExtensionNames = make([]string, extensionsCount)
	vulkan.EnumerateDeviceExtensionProperties(handle, "", &extensionsCount, capabilities.Extensions)
	for i := 0; i < len(capabilities.Extensions); i++ {
		capabilities.Extensions[i].Deref()
		capabilities.ExtensionNames[i] = ctool.ByteArray256ToString(capabilities.Extensions[i].ExtensionName)
	}

	var queueFamiliesCount uint32
	vulkan.GetPhysicalDeviceQueueFamilyProperties(handle, &queueFamiliesCount, nil)
	capabilities.Queue.Families = make([]vulkan.QueueFamilyProperties, queueFamiliesCount)
	vulkan.GetPhysicalDeviceQueueFamilyProperties(handle, &queueFamiliesCount, capabilities.Queue.Families)
	for i := 0; i < len(capabilities.Queue.Families); i++ {
		capabilities.Queue.Families[i].Deref()
		if capabilities.Queue.GraphicsIndex == -1 {
			if vulkan.QueueFlagBits(capabilities.Queue.Families[i].QueueFlags)&vulkan.QueueGraphicsBit == vulkan.QueueGraphicsBit {
				capabilities.Queue.GraphicsIndex = i
			}
		}
		if capabilities.Queue.PresentationIndex == -1 {
			var support vulkan.Bool32
			vulkan.GetPhysicalDeviceSurfaceSupport(handle, uint32(i), *surface, &support)
			if support == 1 {
				capabilities.Queue.PresentationIndex = i
			}
		}
	}

	vulkan.GetPhysicalDeviceSurfaceCapabilities(handle, *surface, &capabilities.Surface.Capabilities)
	capabilities.Surface.Capabilities.Deref()

	var surfaceFormatsCount uint32
	vulkan.GetPhysicalDeviceSurfaceFormats(handle, *surface, &surfaceFormatsCount, nil)
	capabilities.Surface.Formats = make([]vulkan.SurfaceFormat, surfaceFormatsCount)
	vulkan.GetPhysicalDeviceSurfaceFormats(handle, *surface, &surfaceFormatsCount, capabilities.Surface.Formats)
	for i := 0; i < len(capabilities.Surface.Formats); i++ {
		capabilities.Surface.Formats[i].Deref()
		if capabilities.Surface.FormatIndex == -1 {
			if capabilities.Surface.Formats[i].Format == vulkan.FormatB8g8r8a8Srgb && capabilities.Surface.Formats[i].ColorSpace == vulkan.ColorSpaceSrgbNonlinear {
				capabilities.Surface.FormatIndex = i
			}
		}
	}

	var surfacePresentModesCount uint32
	vulkan.GetPhysicalDeviceSurfacePresentModes(handle, *surface, &surfacePresentModesCount, nil)
	capabilities.Surface.PresentModes = make([]vulkan.PresentMode, surfacePresentModesCount)
	vulkan.GetPhysicalDeviceSurfacePresentModes(handle, *surface, &surfacePresentModesCount, capabilities.Surface.PresentModes)
	for i := 0; i < len(capabilities.Surface.PresentModes); i++ {
		if capabilities.Surface.PresentModes[i] == vulkan.PresentModeMailbox {
			capabilities.Surface.PresentModeIndex = i
			break
		}
	}

	capabilities.Surface.Extent = capabilities.Surface.Capabilities.CurrentExtent
	if capabilities.Surface.Extent.Width == vulkan.MaxUint32 {
		var w, h int
		w, h = window.GetFramebufferSize()
		capabilities.Surface.Extent.Width = lo.Clamp(uint32(w), capabilities.Surface.Capabilities.MinImageExtent.Width, capabilities.Surface.Capabilities.MaxImageExtent.Width)
		capabilities.Surface.Extent.Height = lo.Clamp(uint32(h), capabilities.Surface.Capabilities.MinImageExtent.Height, capabilities.Surface.Capabilities.MaxImageExtent.Height)
	}

	return capabilities, nil
}
