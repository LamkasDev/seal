package swapchain

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/pipeline"
	"github.com/vulkan-go/vulkan"
)

type VulkanSwapchainOptions struct {
	CreateInfo vulkan.SwapchainCreateInfo
}

func NewVulkanSwapchainOptions(pipeline *pipeline.VulkanPipeline, surface *vulkan.Surface, old *VulkanSwapchain) VulkanSwapchainOptions {
	options := VulkanSwapchainOptions{}
	options.CreateInfo = vulkan.SwapchainCreateInfo{
		SType:            vulkan.StructureTypeSwapchainCreateInfo,
		Surface:          *surface,
		MinImageCount:    pipeline.Device.Physical.Capabilities.Surface.ImageCount,
		ImageFormat:      pipeline.Device.Physical.Capabilities.Surface.ImageFormats[pipeline.Device.Physical.Capabilities.Surface.ImageFormatIndex].Format,
		ImageColorSpace:  pipeline.Device.Physical.Capabilities.Surface.ImageFormats[pipeline.Device.Physical.Capabilities.Surface.ImageFormatIndex].ColorSpace,
		ImageExtent:      pipeline.Window.Data.Extent,
		ImageArrayLayers: 1,
		ImageUsage:       vulkan.ImageUsageFlags(vulkan.ImageUsageColorAttachmentBit),
		PreTransform:     pipeline.Device.Physical.Capabilities.Surface.Capabilities.CurrentTransform,
		CompositeAlpha:   vulkan.CompositeAlphaOpaqueBit,
		PresentMode:      pipeline.Device.Physical.Capabilities.Surface.PresentModes[pipeline.Device.Physical.Capabilities.Surface.PresentModeIndex],
		Clipped:          vulkan.True,
	}
	if old != nil {
		options.CreateInfo.OldSwapchain = old.Handle
	}
	if pipeline.Device.Physical.Capabilities.Queue.GraphicsIndex != pipeline.Device.Physical.Capabilities.Queue.PresentationIndex {
		options.CreateInfo.ImageSharingMode = vulkan.SharingModeConcurrent
		options.CreateInfo.QueueFamilyIndexCount = 2
		options.CreateInfo.PQueueFamilyIndices = []uint32{uint32(pipeline.Device.Physical.Capabilities.Queue.GraphicsIndex), uint32(pipeline.Device.Physical.Capabilities.Queue.PresentationIndex)}
	} else {
		options.CreateInfo.ImageSharingMode = vulkan.SharingModeExclusive
		options.CreateInfo.QueueFamilyIndexCount = 0
		options.CreateInfo.PQueueFamilyIndices = nil
	}

	return options
}
