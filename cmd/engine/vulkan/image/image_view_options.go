package image

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/vulkan-go/vulkan"
)

type VulkanImageViewOptions struct {
	CreateInfo vulkan.ImageViewCreateInfo
}

func NewVulkanImageViewOptions(device *logical.VulkanLogicalDevice, image *vulkan.Image, format vulkan.Format) VulkanImageViewOptions {
	options := VulkanImageViewOptions{
		CreateInfo: vulkan.ImageViewCreateInfo{
			SType:    vulkan.StructureTypeImageViewCreateInfo,
			Image:    *image,
			ViewType: vulkan.ImageViewType2d,
			Format:   format,
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
		},
	}

	return options
}
