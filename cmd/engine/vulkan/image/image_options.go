package image

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/physical"
	"github.com/vulkan-go/vulkan"
)

type VulkanImageOptions struct {
	CreateInfo   vulkan.ImageCreateInfo
	AllocateInfo vulkan.MemoryAllocateInfo
}

func NewVulkanImageOptions(format vulkan.Format, w uint32, h uint32, usage vulkan.ImageUsageFlags) VulkanImageOptions {
	options := VulkanImageOptions{
		CreateInfo: vulkan.ImageCreateInfo{
			SType:     vulkan.StructureTypeImageCreateInfo,
			ImageType: vulkan.ImageType2d,
			Extent: vulkan.Extent3D{
				Width:  w,
				Height: h,
				Depth:  1,
			},
			MipLevels:     1,
			ArrayLayers:   1,
			Format:        format,
			Tiling:        vulkan.ImageTilingOptimal,
			InitialLayout: vulkan.ImageLayoutUndefined,
			Usage:         usage,
			SharingMode:   vulkan.SharingModeExclusive,
			Samples:       vulkan.SampleCount1Bit,
		},
	}

	return options
}

func UpdateVulkanImageOptions(options *VulkanImageOptions, device *logical.VulkanLogicalDevice, requirements vulkan.MemoryRequirements) error {
	memoryTypeIndex, err := physical.GetVulkanPhysicalDeviceMemoryTypeIndex(&device.Physical.Capabilities, requirements.MemoryTypeBits, vulkan.MemoryPropertyFlags(vulkan.MemoryPropertyDeviceLocalBit))
	if err != nil {
		return err
	}
	options.AllocateInfo = vulkan.MemoryAllocateInfo{
		SType:           vulkan.StructureTypeMemoryAllocateInfo,
		AllocationSize:  requirements.Size,
		MemoryTypeIndex: memoryTypeIndex,
	}

	return nil
}
