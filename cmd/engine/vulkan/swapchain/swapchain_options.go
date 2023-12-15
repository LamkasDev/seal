package swapchain

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/vulkan-go/vulkan"
)

type VulkanSwapchainOptions struct {
	CreateInfo vulkan.SwapchainCreateInfo
}

func NewVulkanSwapchainOptions(device *logical.VulkanLogicalDevice, surface *vulkan.Surface) (VulkanSwapchainOptions, error) {
	options := VulkanSwapchainOptions{}
	options.CreateInfo = vulkan.SwapchainCreateInfo{
		SType:            vulkan.StructureTypeSwapchainCreateInfo,
		Surface:          *surface,
		MinImageCount:    device.Physical.Capabilities.Surface.ImageCount,
		ImageFormat:      device.Physical.Capabilities.Surface.ImageFormats[device.Physical.Capabilities.Surface.ImageFormatIndex].Format,
		ImageColorSpace:  device.Physical.Capabilities.Surface.ImageFormats[device.Physical.Capabilities.Surface.ImageFormatIndex].ColorSpace,
		ImageExtent:      device.Physical.Capabilities.Surface.ImageExtent,
		ImageArrayLayers: 1,
		ImageUsage:       vulkan.ImageUsageFlags(vulkan.ImageUsageColorAttachmentBit),
		PreTransform:     device.Physical.Capabilities.Surface.Capabilities.CurrentTransform,
		CompositeAlpha:   vulkan.CompositeAlphaOpaqueBit,
		PresentMode:      device.Physical.Capabilities.Surface.PresentModes[device.Physical.Capabilities.Surface.PresentModeIndex],
		Clipped:          vulkan.True,
		OldSwapchain:     nil,
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

	return options, nil
}
