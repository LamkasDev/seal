package image

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/vulkan-go/vulkan"
)

type VulkanImageOptions struct {
	CreateInfo vulkan.ImageCreateInfo
}

func NewVulkanImageOptions(device *logical.VulkanLogicalDevice, image *vulkan.Image) VulkanImageOptions {
	options := VulkanImageOptions{
		CreateInfo: vulkan.ImageCreateInfo{
			SType:     vulkan.StructureTypeImageCreateInfo,
			ImageType: vulkan.ImageType2d,
			Extent: ,
		},
	}

	return options
}
