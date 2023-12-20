package buffer

import (
	"github.com/LamkasDev/seal/cmd/engine/vulkan/logical"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/physical"
	"github.com/vulkan-go/vulkan"
)

type VulkanBufferOptions struct {
	CreateInfo   vulkan.BufferCreateInfo
	AllocateInfo vulkan.MemoryAllocateInfo
}

func NewVulkanBufferOptions(size vulkan.DeviceSize, usage vulkan.BufferUsageFlags, sharingMode vulkan.SharingMode) VulkanBufferOptions {
	options := VulkanBufferOptions{
		CreateInfo: vulkan.BufferCreateInfo{
			SType:       vulkan.StructureTypeBufferCreateInfo,
			Size:        size,
			Usage:       usage,
			SharingMode: sharingMode,
		},
	}

	return options
}

func UpdateVulkanBufferOptions(options *VulkanBufferOptions, device *logical.VulkanLogicalDevice, requirements vulkan.MemoryRequirements, flags vulkan.MemoryPropertyFlags) error {
	memoryTypeIndex, err := physical.GetVulkanPhysicalDeviceMemoryTypeIndex(&device.Physical.Capabilities, requirements.MemoryTypeBits, flags)
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
