package swapchain

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	sealWindow "github.com/LamkasDev/seal/cmd/engine/window"
	"github.com/vulkan-go/vulkan"
)

type VulkanSwapchainOptions struct {
	CreateInfo vulkan.SwapchainCreateInfo
}

func NewVulkanSwapchainOptions(device *logical.VulkanLogicalDevice, window *sealWindow.Window, surface *vulkan.Surface, old *VulkanSwapchain) VulkanSwapchainOptions {
	options := VulkanSwapchainOptions{}
	options.CreateInfo = vulkan.SwapchainCreateInfo{
		SType:            vulkan.StructureTypeSwapchainCreateInfo,
		Surface:          *surface,
		MinImageCount:    device.Physical.Capabilities.Surface.ImageCount,
		ImageFormat:      device.Physical.Capabilities.Surface.ImageFormats[device.Physical.Capabilities.Surface.ImageFormatIndex].Format,
		ImageColorSpace:  device.Physical.Capabilities.Surface.ImageFormats[device.Physical.Capabilities.Surface.ImageFormatIndex].ColorSpace,
		ImageExtent:      window.Data.Extent,
		ImageArrayLayers: 1,
		ImageUsage:       vulkan.ImageUsageFlags(vulkan.ImageUsageColorAttachmentBit),
		PreTransform:     device.Physical.Capabilities.Surface.Capabilities.CurrentTransform,
		CompositeAlpha:   vulkan.CompositeAlphaOpaqueBit,
		PresentMode:      device.Physical.Capabilities.Surface.PresentModes[device.Physical.Capabilities.Surface.PresentModeIndex],
		Clipped:          vulkan.True,
	}
	if old != nil {
		options.CreateInfo.OldSwapchain = old.Handle
	}
	if device.Physical.Capabilities.Queue.GraphicsIndex != device.Physical.Capabilities.Queue.PresentationIndex {
		options.CreateInfo.ImageSharingMode = vulkan.SharingModeConcurrent
		options.CreateInfo.QueueFamilyIndexCount = 2
		options.CreateInfo.PQueueFamilyIndices = []uint32{uint32(device.Physical.Capabilities.Queue.GraphicsIndex), uint32(device.Physical.Capabilities.Queue.PresentationIndex)}
	} else {
		options.CreateInfo.ImageSharingMode = vulkan.SharingModeExclusive
		options.CreateInfo.QueueFamilyIndexCount = 0
		options.CreateInfo.PQueueFamilyIndices = nil
	}

	return options
}
