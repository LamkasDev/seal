package physical

import (
	"github.com/LamkasDev/seal/cmd/common/ctool"
	"github.com/LamkasDev/seal/cmd/logger"
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
	ImageFormats     []vulkan.SurfaceFormat
	ImageFormatIndex int
	PresentModes     []vulkan.PresentMode
	PresentModeIndex int
	ImageExtent      vulkan.Extent2D
	ImageCount       uint32
}

func NewVulkanPhysicalDeviceCapabilities(handle vulkan.PhysicalDevice, window *glfw.Window, surface *vulkan.Surface) (VulkanPhysicalDeviceCapabilities, error) {
	capabilities := VulkanPhysicalDeviceCapabilities{}

	var extensionsCount uint32
	if res := vulkan.EnumerateDeviceExtensionProperties(handle, "", &extensionsCount, nil); res != vulkan.Success {
		logger.DefaultLogger.Errorf("vulkan error: %d", int32(res))
	}
	capabilities.Extensions = make([]vulkan.ExtensionProperties, extensionsCount)
	capabilities.ExtensionNames = make([]string, extensionsCount)
	if res := vulkan.EnumerateDeviceExtensionProperties(handle, "", &extensionsCount, capabilities.Extensions); res != vulkan.Success {
		logger.DefaultLogger.Errorf("vulkan error: %d", int32(res))
	}
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
			if res := vulkan.GetPhysicalDeviceSurfaceSupport(handle, uint32(i), *surface, &support); res != vulkan.Success {
				logger.DefaultLogger.Errorf("vulkan error: %d", int32(res))
			}
			if support == 1 {
				capabilities.Queue.PresentationIndex = i
			}
		}
	}

	if res := vulkan.GetPhysicalDeviceSurfaceCapabilities(handle, *surface, &capabilities.Surface.Capabilities); res != vulkan.Success {
		logger.DefaultLogger.Errorf("vulkan error: %d", int32(res))
	}
	capabilities.Surface.Capabilities.Deref()

	var surfaceFormatsCount uint32
	if res := vulkan.GetPhysicalDeviceSurfaceFormats(handle, *surface, &surfaceFormatsCount, nil); res != vulkan.Success {
		logger.DefaultLogger.Errorf("vulkan error: %d", int32(res))
	}
	capabilities.Surface.ImageFormats = make([]vulkan.SurfaceFormat, surfaceFormatsCount)
	if res := vulkan.GetPhysicalDeviceSurfaceFormats(handle, *surface, &surfaceFormatsCount, capabilities.Surface.ImageFormats); res != vulkan.Success {
		logger.DefaultLogger.Errorf("vulkan error: %d", int32(res))
	}
	for i := 0; i < len(capabilities.Surface.ImageFormats); i++ {
		capabilities.Surface.ImageFormats[i].Deref()
		if capabilities.Surface.ImageFormatIndex == -1 {
			if capabilities.Surface.ImageFormats[i].Format == vulkan.FormatB8g8r8a8Srgb && capabilities.Surface.ImageFormats[i].ColorSpace == vulkan.ColorSpaceSrgbNonlinear {
				capabilities.Surface.ImageFormatIndex = i
			}
		}
	}

	var surfacePresentModesCount uint32
	if res := vulkan.GetPhysicalDeviceSurfacePresentModes(handle, *surface, &surfacePresentModesCount, nil); res != vulkan.Success {
		logger.DefaultLogger.Errorf("vulkan error: %d", int32(res))
	}
	capabilities.Surface.PresentModes = make([]vulkan.PresentMode, surfacePresentModesCount)
	if res := vulkan.GetPhysicalDeviceSurfacePresentModes(handle, *surface, &surfacePresentModesCount, capabilities.Surface.PresentModes); res != vulkan.Success {
		logger.DefaultLogger.Errorf("vulkan error: %d", int32(res))
	}
	for i := 0; i < len(capabilities.Surface.PresentModes); i++ {
		if capabilities.Surface.PresentModes[i] == vulkan.PresentModeMailbox {
			capabilities.Surface.PresentModeIndex = i
			break
		}
	}

	capabilities.Surface.ImageExtent = capabilities.Surface.Capabilities.CurrentExtent
	capabilities.Surface.ImageExtent.Deref()
	if capabilities.Surface.ImageExtent.Width == vulkan.MaxUint32 {
		var w, h int
		w, h = window.GetFramebufferSize()
		capabilities.Surface.ImageExtent.Width = lo.Clamp(uint32(w), capabilities.Surface.Capabilities.MinImageExtent.Width, capabilities.Surface.Capabilities.MaxImageExtent.Width)
		capabilities.Surface.ImageExtent.Height = lo.Clamp(uint32(h), capabilities.Surface.Capabilities.MinImageExtent.Height, capabilities.Surface.Capabilities.MaxImageExtent.Height)
	}

	capabilities.Surface.ImageCount = capabilities.Surface.Capabilities.MinImageCount + 1
	if capabilities.Surface.Capabilities.MaxImageCount > 0 && capabilities.Surface.ImageCount > capabilities.Surface.Capabilities.MaxImageCount {
		capabilities.Surface.ImageCount = capabilities.Surface.Capabilities.MaxImageCount
	}

	return capabilities, nil
}
