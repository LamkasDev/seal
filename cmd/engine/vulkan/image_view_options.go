package vulkan

import (
	"github.com/vulkan-go/vulkan"
)

type VulkanImageViewOptions struct {
	CreateInfo vulkan.ImageViewCreateInfo
}

func NewVulkanImageViewOptions(device *VulkanLogicalDevice, image *vulkan.Image) (VulkanImageViewOptions, error) {
	options := VulkanImageViewOptions{}
	options.CreateInfo = vulkan.ImageViewCreateInfo{
		SType:    vulkan.StructureTypeImageViewCreateInfo,
		Image:    *image,
		ViewType: vulkan.ImageViewType2d,
		Format:   device.Physical.Capabilities.Surface.ImageFormats[device.Physical.Capabilities.Surface.ImageFormatIndex].Format,
		Components: vulkan.ComponentMapping{
			R: vulkan.ComponentSwizzleIdentity,
			G: vulkan.ComponentSwizzleIdentity,
			B: vulkan.ComponentSwizzleIdentity,
			A: vulkan.ComponentSwizzleIdentity,
		},
		SubresourceRange: vulkan.ImageSubresourceRange{
			AspectMask:     vulkan.ImageAspectFlags(vulkan.ImageAspectColorBit),
			BaseMipLevel:   0,
			LevelCount:     1,
			BaseArrayLayer: 0,
			LayerCount:     1,
		},
	}

	return options, nil
}
