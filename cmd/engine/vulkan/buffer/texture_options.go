package buffer

import (
	"github.com/vulkan-go/vulkan"
)

type VulkanTextureBufferOptions struct {
	Pixels []uint8
}

func NewVulkanTextureBufferOptions(pixels []uint8) VulkanTextureBufferOptions {
	data := VulkanTextureBufferOptions{
		Pixels: pixels,
	}

	return data
}

func GetVulkanTextureBufferOptionsPixelsOffset(options *VulkanTextureBufferOptions) vulkan.DeviceSize {
	return 0
}

func GetVulkanTextureBufferOptionsPixelsSize(options *VulkanTextureBufferOptions) vulkan.DeviceSize {
	return vulkan.DeviceSize(len(options.Pixels))
}

func GetVulkanTextureBufferOptionsSize(options *VulkanTextureBufferOptions) vulkan.DeviceSize {
	return GetVulkanTextureBufferOptionsPixelsSize(options)
}
